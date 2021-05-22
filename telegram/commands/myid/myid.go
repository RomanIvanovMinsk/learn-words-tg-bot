package myid

import (
	"WordsBot/models/telegram"
	"WordsBot/telegram/api"
	"fmt"
	"strconv"
)

func Name() string {
	return "/myid"
}

func Description() string {
	return "return user id"
}

func Handle(body *telegram.WebhookReqBody) error {
	fmt.Println("Start my id")
	fmt.Printf("chat id %s\n", strconv.FormatInt(body.Message.Chat.ID, 10))
	reqBody := &telegram.SendMessageReqBodyReply{
		ChatID: body.Message.Chat.ID,
		Text:   strconv.FormatInt(body.Message.Chat.ID, 10),
	}

	err := api.SendMessage(reqBody)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return err
	}

	return nil
}
