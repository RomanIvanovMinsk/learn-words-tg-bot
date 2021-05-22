package start

import (
	"WordsBot/models/telegram"
	"WordsBot/services/sqlService"
	"fmt"
	"strconv"
)

func Name() string {
	return "/start"
}

func Handle(body *telegram.WebhookReqBody) error {
	fmt.Println("Start creating profile")
	userId, err := sqlService.CreateProfile(strconv.FormatInt(body.Message.Chat.ID, 10))
	if err != nil {
		return err
	} else {
		fmt.Printf("UserId %s\n", userId)
	}
	return nil
}
