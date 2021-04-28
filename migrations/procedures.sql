Use WordsBot;
go

drop function if exists dbo.GetWordId
go
create or
alter function dbo.GetWordId(
    @UserId uniqueidentifier,
    @Stem nvarchar(50),
    @Lang nvarchar(6)
)
    returns uniqueidentifier
as
begin
    DECLARE @id uniqueidentifier


    select top (1) @id = Id
    from Words
    where Stem = @Stem
      and Lang = @Lang
      and UserId = @UserId

    return @id;
end
go


create or
alter function dbo.GetWordUsageId(
    @WordId uniqueidentifier,
    @Usage nvarchar(500)
)
    returns uniqueidentifier
as
begin
    declare @id uniqueidentifier

    select top (1) @id = Id
    from dbo.Usage
    where WordId = @wordId
      and Text = @Usage

    return @id;
end
go



create or
alter procedure CreateProfile @TelegramId nvarchar(30)
as
begin
    insert into Profile(UserTelegramId)
    OUTPUT Inserted.UserId
    values (@TelegramId)
end
go



create or
alter procedure GetUserIdByTelegramId @TelegramId nvarchar(30)
as
begin
    select UserId
    from Profile
    where UserTelegramId = @TelegramId
end
go



create or
alter procedure AddWord @UserId uniqueidentifier,
                        @Stem nvarchar(50),
                        @Word nvarchar(50),
                        @Lang nvarchar(6),
                        @Usage nvarchar(500)
as
begin
    declare @wordId as uniqueidentifier;
    exec @wordId = dbo.GetWordId @UserId, @Stem, @Lang;
    if @wordId is null
        begin
            insert into Words(UserId, Word, Stem, Lang)
            values (@UserId, @Word, @Stem, @Lang);
            exec @wordId = dbo.GetWordId @UserId, @Stem, @Lang;
        end

    if @Usage is not null
        begin
            if dbo.GetWordUsageId(@wordId, @Usage) is null
                insert into Usage(WordId, Text)
                values (@wordId, @Usage)
        end
end
go



create or
alter procedure QueueWords
as
begin
--     declare @WordsInQueue as int;
--     select @WordsInQueue = WordsInQueue from Settings;

    With C0 as (
        select UserId, (max(WordsInQueue) - Count(Q.WordId)) as toAdd
        from Queue as Q
                 right join dbo.Words on Q.WordId = Words.Id
        cross join Settings
        group by UserId
    ),
         C1 as (
             select distinct WordId,
                    toAdd,
                    count(U.WordId) over (partition by WordId) as cnt,
                    W.UserId
             from Words as W
                      inner join C0 on W.UserId = C0.UserId
                      left join Usage U on W.Id = U.WordId

             where WordId not in (select WordId from Queue) and
                   W.Learned is null
         ),
         C2 as (
             select *,
                    row_number() over (partition by UserId order by cnt desc) as rnb
             from C1
         )
    insert
    into Queue(WordId, Due, Stage, IsWaiting)

    select WordId, SYSUTCDATETIME() as Due, 0 as Stage, 0 as IsWaiting
    from C2
    where rnb <= toAdd;
end
go



create or
alter procedure GetIntervalWords @UserId uniqueidentifier
as
begin

    begin transaction
        declare @ToReturn table
                          (
                              Id        uniqueidentifier,
                              UserId    uniqueidentifier,
                              Word      nvarchar(50),
                              Stem      nvarchar(50),
                              Lang      nvarchar(6),
                              IsWaiting bit
                          )
        insert into @ToReturn(Id, UserId, Word, Stem, Lang, IsWaiting)
        select top (1) W2.Id, W2.Userid, W2.Word, W2.Stem, W2.Lang, IsWaiting
        from Queue
                 join Words W2 on W2.Id = Queue.WordId
        where Due <= SYSUTCDATETIME()
          and UserId = @UserId
        order by IsWaiting desc, Stage asc, Due asc

        declare @isWaiting as bit;
        select @isWaiting = IsWaiting from @ToReturn;

        if @isWaiting = 0
            begin
                update Queue
                Set IsWaiting = 1
                where wordId in (select WordId from @ToReturn)
            end

        select * from @ToReturn
    commit transaction
end
go


create or
alter function dbo.GetNextDueDate(@Stage int)
    returns datetime2
as
begin
    declare @CurrentTime as datetime2 = SYSUTCDATETIME();
    return case @Stage
               when 1 then DATEADD(day, 1, @CurrentTime)
               when 2 then DATEADD(day, 2, @CurrentTime)
               when 3 then DATEADD(day, 4, @CurrentTime)
               when 4 then DATEADD(day, 7, @CurrentTime)
               else DATEADD(day, 14, @CurrentTime)
        end

end
go



create or
alter procedure Answer @UserId uniqueidentifier, @IsRemember bit
as
begin
    begin transaction
        declare @CurrentWord as uniqueidentifier;
        declare @CurrentStage as int;
        select @CurrentWord = WordId, @CurrentStage = Stage
        from Queue
                 inner join Words W on W.Id = Queue.WordId
        where W.UserId = @UserId
          and IsWaiting = 1;

        insert into Answers(WordId, Stage, IsRemember)
        select @CurrentWord, Stage, @IsRemember
        from Queue
        where WordId = @CurrentWord;

        declare @UpdatedWord table
                             (
                                 Stage int
                             )

        update Queue
        set Stage = IIF(@IsRemember = 1, Stage + 1, Stage),
            Due   = dbo.GetNextDueDate(IIF(@IsRemember = 1, Stage + 1, Stage))
        output Inserted.Stage into @UpdatedWord
        where WordId = @CurrentWord;

        select @CurrentStage = max(Stage) from @UpdatedWord;
        if @CurrentStage = 6
            begin
                delete
                from Queue
                where WordId = @CurrentWord;
                update Words
                set Learned = SYSUTCDATETIME();
            end
    commit transaction
end