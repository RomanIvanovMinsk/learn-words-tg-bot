package models

type WebhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Usage
type Usage struct {
	Usage string
}

// Word
type Word struct {
	Word   string
	Stem   string
	Lang   string
	Usages []Usage
}


type WordsList struct {
	ChatID int64
	Words []Word
}
