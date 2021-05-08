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

type SettingsModel struct {
	Settings Settings `json:"settings"`
}

type Settings struct {
	NumberOfShards int `json:"number_of_shards"`
	Index          struct {
		analysis Analysis
	} `json:"index"`
}

type Analysis struct {
}

type IndexedModel struct {
	Index IndexInner `json:"index"`
}

type IndexInner struct {
	Index string `json:"_index" `
	Type  string `json:"_type"`
	Id    int    `json:"_id"`
}

type DictionaryWord struct {
	Key   string `json:"key"`
	Value string `json:value`
}
