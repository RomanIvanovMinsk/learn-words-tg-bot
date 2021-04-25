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
    UserTelegramId nvarchar(30)      not null unique,
)


create table Words
(
    Id     uniqueidentifier not null primary key default (NEWID()),
    UserId uniqueidentifier not null,
    Word   nvarchar(50)      not null,
    Stem   nvarchar(50)      not null,
    Lang   nvarchar(6)       not null,

    constraint FK_UserId foreign key (UserId) references Profile (UserId),
    constraint unique_stem_lang unique (Stem, Lang)
)

create table Usage
(
    Id     uniqueidentifier primary key default (NEWID()),
    WordId uniqueidentifier not null,
    Text   nvarchar(500)     not null,

    constraint FK_WordId foreign key (WordId) references Words (Id)
)