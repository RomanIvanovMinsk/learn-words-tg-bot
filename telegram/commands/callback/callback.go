package callback

import (
	wr "WordsBot/models"
	"WordsBot/models/telegram"
	"WordsBot/services/wordsManager"
	"WordsBot/telegram/api"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func Handle(body *telegram.WebhookReqBody) error {
	fmt.Println("Start process user answer")
	chatId := body.Callback.Message.Chat.ID
	answer, err := ProcessAnswer(chatId, body)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return err
	}
	if answer.Command != "Answer" {
		answer := &wr.GetUsages{}
		json.Unmarshal([]byte(body.Callback.Data), &answer)
		usages, _ := wordsManager.GetUsages(strconv.FormatInt(chatId, 10), answer.Id, answer.Offset)
		fmt.Printf("%s", usages)
		text := strings.Join(usages, "\n\n")
		if len(text) <= 0 {
			text = "No usages found for the word"
		}
		err = api.SendMessage(&telegram.SendMessageReqBodyReply{
			ChatID: chatId,
			Text:   text,
		})
		if err != nil {
			fmt.Println("error in sending reply:", err)
		}
	} else {
		fmt.Printf("We get the answer and %d answer is %t", chatId, answer.Remember)
		wordsManager.Answer(strconv.FormatInt(chatId, 10), answer.Remember)

		err = api.AnswerCallbackQuery(&telegram.AnswerCallbackQuery{CallbackQueryId: body.Callback.ID})
		if err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}

		err = api.EditMessageReplyMarkup(&telegram.EditMessageReplyMarkupRequest{
			ChatId:      body.Callback.Message.Chat.ID,
			MessageId:   body.Callback.Message.MessageId,
			ReplyMarkup: telegram.InlineKeyboardMarkup{Keyboard: [][]telegram.InlineKeyboardButton{}},
		})

		if err != nil {
			fmt.Println("error in sending reply:", err)
			return err
		}
	}
	return nil
}

func ProcessAnswer(chatID int64, callback *telegram.WebhookReqBody) (wr.Answer, error) {
	answer := wr.Answer{}
	json.Unmarshal([]byte(callback.Callback.Data), &answer)
	return answer, nil
}
