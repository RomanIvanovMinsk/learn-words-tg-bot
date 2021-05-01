package actions

import (
	"WordsBot/models/telegram"
	"WordsBot/services/wordsManager"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	wr "WordsBot/models"
)

var botHost string

func Configure(host string) {
	botHost = host
}

//The below code deals with the process of sending a response message
// to the user
func SayPolo(chatID int64) error {
	// Create the request body struct
	reqBody := &telegram.SendMessageReqBody{
		ChatID: chatID,
		Text:   "Polo!!",
	}

	err := SendResponse(reqBody)
	if err != nil {
		return err
	}

	return nil
}

// get userId and return it
func GetMyId(chatID int64) error {
	fmt.Println("chat id %d", strconv.FormatInt(chatID, 10))
	reqBody := &telegram.SendMessageReqBody{
		ChatID: chatID,
		Text:   strconv.FormatInt(chatID, 10),
	}

	err := SendResponse(reqBody)

	if err != nil {
		return err
	}

	return nil
}

func SendResponse(reqBody *telegram.SendMessageReqBody) error {
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post(botHost+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func SendResponseWithReply(reqBody *telegram.SendMessageReqBodyReply) error {
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post(botHost+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func SendQuestion(chatID int64, word *wr.Word) error {
	buttons := [][]telegram.InlineKeyboardButton{{}}
	buttons[0] = append(buttons[0], *SetUpButton(wr.Answer{Word: word.Stem, Remember: true}, "I remember the word"))
	buttons[0] = append(buttons[0], *SetUpButton(wr.Answer{Word: word.Stem, Remember: false}, "I do not remember word"))

	keyboard := &telegram.InlineKeyboardMarkup{}
	keyboard.Keyboard = buttons

	reqBody := &telegram.SendMessageReqBodyReply{
		ChatID:    chatID,
		Text:      fmt.Sprintf("Please select one answer: *%s*", word.Stem),
		Reply:     *keyboard,
		ParseMode: "Markdown",
	}

	err := SendResponseWithReply(reqBody)

	if err != nil {
		return err
	}

	return nil
}

func ProcessAnswer(chatID int64, callback *telegram.WebhookReqBody) error {
	answer := wr.Answer{}
	json.Unmarshal([]byte(callback.Callback.Data), &answer)
	fmt.Printf("We get the answer and %d answer is %t", chatID, answer.Remember)
	wordsManager.Answer(strconv.FormatInt(chatID, 10), answer.Remember)
	return nil
}

func SetUpButton(word wr.Answer, comment string) *telegram.InlineKeyboardButton {
	button := &telegram.InlineKeyboardButton{}
	button.Text = comment
	b, _ := json.Marshal(word)
	button.CallbackData = string(b)

	return button
}

func EditMessageReplyMarkup(reqBody *telegram.EditMessageReplyMarkupRequest) error {

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(botHost+"/editMessageReplyMarkup", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		var b []byte
		res.Body.Read(b)
		fmt.Println(string(b))
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func AnswerCallbackQuery(reqBody *telegram.AnswerCallbackQuery) error {

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(botHost+"/answerCallbackQuery", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		var b []byte
		res.Body.Read(b)
		fmt.Println(string(b))
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func SetMyCommands(reqBody *telegram.SetMyCommandsRequest) error {

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(botHost+"/setMyCommands", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		var b []byte
		res.Body.Read(b)
		fmt.Println(string(b))
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
