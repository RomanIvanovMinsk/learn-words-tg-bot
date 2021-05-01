package models

type Usage struct {
	Usage string
}

type Word struct {
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
	Word     string
	Remember bool
}
