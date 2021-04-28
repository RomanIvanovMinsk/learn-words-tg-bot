-- use master
-- drop database WordsBot;
-- go

IF NOT EXISTS(SELECT *
              FROM sys.databases
              WHERE name = 'WordsBot')
    BEGIN
        CREATE DATABASE [WordsBot]
    END
go
use WordsBot
go
create table Profile
(
    UserId         uniqueidentifier not null primary key default (NEWID()),
    UserTelegramId nvarchar(30)     not null unique,
)


create table Words
(
    Id      uniqueidentifier not null primary key default (NEWID()),
    UserId  uniqueidentifier not null,
    Word    nvarchar(50)     not null,
    Stem    nvarchar(50)     not null,
    Lang    nvarchar(6)      not null,
    Learned datetime2        null,
    Created datetime2         not null             default (SYSUTCDATETIME()),

    constraint FK_Words_UserId foreign key (UserId) references Profile (UserId),
    constraint unique_stem_lang unique (Stem, Lang)
)

create table Usage
(
    Id     uniqueidentifier primary key default (NEWID()),
    WordId uniqueidentifier not null,
    Text   nvarchar(500)    not null,

    constraint FK_Usage_WordId foreign key (WordId) references Words (Id)
)

create table Queue
(
    WordId    uniqueidentifier primary key,
    Due       DateTime2 not null,
    Stage     int       not null default (0),
    IsWaiting bit       not null default (0),

    constraint FK_Queue_WordId foreign key (WordId) references Words (Id)
)

create table Answers
(
    WordId     uniqueidentifier not null,
    Stage      int              not null,
    IsRemember bit              not null,
    Date       datetime2        not null default (SYSUTCDATETIME())
)

create table Settings
(
    WordsInQueue int not null default(15)
)

insert into Settings(WordsInQueue)
values(default);