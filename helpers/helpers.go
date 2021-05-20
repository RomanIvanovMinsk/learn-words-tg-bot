package helpers

import (
	"WordsBot/models/telegram"
	"WordsBot/services/sqlService"
	"WordsBot/services/wordsManager"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	wr "WordsBot/models"

	action "WordsBot/actions"
)

func DecodeRequestBody(req *http.Request, body *telegram.WebhookReqBody) {
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
	}
}

func SelectAction(body *telegram.WebhookReqBody) {
	// Check if the message contains the word "marco"
	// if not, return without doing anything
	command := strings.ToLower(body.Message.Text)
	if strings.Contains(command, "marco") {
		// If the text contains marco, call the `sayPolo` function, which
		// is defined below
		if err := action.SayPolo(body.Message.Chat.ID); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}

		// log a confirmation message if the message is sent successfully
		fmt.Println("reply sent")
	}

	if command == "/myid" {
		fmt.Println("Start my id")
		if err := action.GetMyId(body.Message.Chat.ID); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	}

	if command == "/start" {
		fmt.Println("Start creating profile")
		userId, err := sqlService.CreateProfile(strconv.FormatInt(body.Message.Chat.ID, 10))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("UserId %s\n", userId)
		}
		return
	}
	if command == "/givemetheword" {
		chatId := body.Message.Chat.ID
		word, err := getTheWord(chatId)
		if err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}

		if err := action.SendQuestion(chatId, word); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	}

	if command == "/import_utility" {
		chatId := body.Message.Chat.ID
		reqBody := &telegram.SendMessageReqBody{
			ChatID: chatId,
			Text:   fmt.Sprintf("Here the link\n%s", "https://github.com/xtergs/WordsImport/releases"),
		}
		action.SendResponse(reqBody)
	}

	if body.Callback.Data != "" {
		fmt.Println("Start process user answer")
		chatId := body.Callback.Message.Chat.ID
		answer, err := action.ProcessAnswer(chatId, body)
		if err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
		if answer.Command != "Answer" {
			answer := &wr.GetUsages{}
			json.Unmarshal([]byte(body.Callback.Data), &answer)
			usages, _ := wordsManager.GetUsages(strconv.FormatInt(chatId, 10), answer.Id, answer.Offset)
			fmt.Printf("%s", usages)
			err = action.SendResponse(&telegram.SendMessageReqBody{
				ChatID: chatId,
				Text:   strings.Join(usages, "\n\n"),
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("We get the answer and %d answer is %t", chatId, answer.Remember)
			wordsManager.Answer(strconv.FormatInt(chatId, 10), answer.Remember)

			err = action.AnswerCallbackQuery(&telegram.AnswerCallbackQuery{CallbackQueryId: body.Callback.ID})
			if err != nil {
				fmt.Println("error in sending reply:", err)
				return
			}

			err = action.EditMessageReplyMarkup(&telegram.EditMessageReplyMarkupRequest{
				ChatId:      body.Callback.Message.Chat.ID,
				MessageId:   body.Callback.Message.MessageId,
				ReplyMarkup: telegram.InlineKeyboardMarkup{Keyboard: [][]telegram.InlineKeyboardButton{}},
			})

			if err != nil {
				fmt.Println("error in sending reply:", err)
				return
			}
		}

	}

}

func getTheWord(chatID int64) (*wr.Word, error) {
	word, err := wordsManager.GetIntervalWords(strconv.FormatInt(chatID, 10))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return word, nil
}
