package models

type Usage struct {
	Usage string
}

type Word struct {
	Id     string
	Word   string
	Stem   string
	Lang   string
	Usages []Usage
}

type WordsList struct {
	TelegramUserId string
	Words          []Word
	ChatID         int64
}

type Answer struct {
	Command  string
	Word     string
	Remember bool
}

type GetUsages struct {
	Id     string
	Offset int
}
