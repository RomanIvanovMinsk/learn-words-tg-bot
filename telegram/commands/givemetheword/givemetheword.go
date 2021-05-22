package givemetheword

import (
	wr "WordsBot/models"
	"WordsBot/models/telegram"
	"WordsBot/services/wordsManager"
	"WordsBot/telegram/api"
	"encoding/json"
	"fmt"
	"strconv"
)

func Name() string {
	return "/givemetheword"
}

func Description() string {
	return "get next word"
}

func Handle(body *telegram.WebhookReqBody) error {
	chatId := body.Message.Chat.ID
	word, err := getTheWord(chatId)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return err
	}

	if err := SendQuestion(chatId, word); err != nil {
		fmt.Println("error in sending reply:", err)
		return err
	}

	return err
}

func SendQuestion(chatID int64, word *wr.Word) error {
	buttons := buttonsForKeyboard(word)

	keyboard := &telegram.InlineKeyboardMarkup{}
	keyboard.Keyboard = buttons

	reqBody := &telegram.SendMessageReqBodyReply{
		ChatID:    chatID,
		Text:      fmt.Sprintf("Please select one answer: *%s*", word.Stem),
		Reply:     *keyboard,
		ParseMode: "Markdown",
	}

	err := api.SendMessage(reqBody)

	if err != nil {
		return err
	}

	return nil
}

func buttonsForKeyboard(word *wr.Word) [][]telegram.InlineKeyboardButton {
	buttons := [][]telegram.InlineKeyboardButton{{}}
	buttons[0] = append(buttons[0], *setUpButton(wr.Answer{Command: "Answer", Word: word.Stem, Remember: true}, "I remember the word"))
	buttons[0] = append(buttons[0], *setUpButton(wr.Answer{Command: "Answer", Word: word.Stem, Remember: false}, "I do not remember word"))
	buttons[0] = append(buttons[0], *setUpButton(wr.GetUsages{
		Id:     word.Id,
		Offset: 0}, "Show usages"))
	return buttons
}

func setUpButton(word interface{}, comment string) *telegram.InlineKeyboardButton {
	button := &telegram.InlineKeyboardButton{}
	button.Text = comment
	b, _ := json.Marshal(word)
	button.CallbackData = string(b)

	return button
}

func getTheWord(chatID int64) (*wr.Word, error) {
	word, err := wordsManager.GetIntervalWords(strconv.FormatInt(chatID, 10))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return word, nil
}
