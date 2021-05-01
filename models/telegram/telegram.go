package telegram

type WebhookReqBody struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`

	Callback CallbackQuery `json:"callback_query"`
}

type User struct {
	Id    int64 `json:"id"`
	IsBot bool  `json:"is_bot"`
}

type Message struct {
	MessageId int64           `json:"message_id"`
	Text      string          `json:"text"`
	From      User            `json:"from"`
	Chat      Chat            `json:"chat"`
	Entities  []MessageEntity `json:"entities"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type Callback struct {
	ID      string  `json:"id"`
	Info    string  `json:"data"`
	Message Message `json:"message"`
}

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type SendMessageReqBodyReply struct {
	ChatID    int64                `json:"chat_id"`
	Text      string               `json:"text"`
	Reply     InlineKeyboardMarkup `json:"reply_markup"`
	ParseMode string               `json:"parse_mode"`
}

// Keyboard button
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	Data    string  `json:"data"`
	Message Message `json:"message"`
	From    User    `json:"from"`
}

type EditMessageReplyMarkupRequest struct {
	ChatId      int64                `json:"chat_id"`
	MessageId   int64                `json:"message_id"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type AnswerCallbackQuery struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
}

type SetMyCommandsRequest struct {
	Commands []BotCommand `json:"commands"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
