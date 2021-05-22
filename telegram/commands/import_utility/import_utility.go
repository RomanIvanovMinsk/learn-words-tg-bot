package import_utility

import (
	"WordsBot/models/telegram"
	"WordsBot/telegram/api"
	"fmt"
)

func Name() string {
	return "/import_utility"
}

func Description() string {
	return "return user id"
}

func Handle(body *telegram.WebhookReqBody) error {
	chatId := body.Message.Chat.ID
	reqBody := &telegram.SendMessageReqBodyReply{
		ChatID: chatId,
		Text:   fmt.Sprintf("Here the link\n%s", "https://github.com/xtergs/WordsImport/releases"),
	}
	api.SendMessage(reqBody)

	return nil
}
