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