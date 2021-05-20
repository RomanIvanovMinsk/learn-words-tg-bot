package actions

import (
	wr "WordsBot/models"
	"WordsBot/models/telegram"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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
	return post("/sendMessage", reqBody)
}

func SendResponseWithReply(reqBody *telegram.SendMessageReqBodyReply) error {
	return post("/sendMessage", reqBody)
}

func SendQuestion(chatID int64, word *wr.Word) error {
	buttons := ButtonsForKeyboard(word)

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

func ButtonsForKeyboard(word *wr.Word) [][]telegram.InlineKeyboardButton {
	buttons := [][]telegram.InlineKeyboardButton{{}}
	buttons[0] = append(buttons[0], *SetUpButton(wr.Answer{Command: "Answer", Word: word.Stem, Remember: true}, "I remember the word"))
	buttons[0] = append(buttons[0], *SetUpButton(wr.Answer{Command: "Answer", Word: word.Stem, Remember: false}, "I do not remember word"))
	buttons[0] = append(buttons[0], *SetUpButton(wr.GetUsages{
		Id:     word.Id,
		Offset: 0}, "Show usages"))
	return buttons
}

func ProcessAnswer(chatID int64, callback *telegram.WebhookReqBody) (wr.Answer, error) {
	answer := wr.Answer{}
	json.Unmarshal([]byte(callback.Callback.Data), &answer)
	return answer, nil
}

func SetUpButton(word interface{}, comment string) *telegram.InlineKeyboardButton {
	button := &telegram.InlineKeyboardButton{}
	button.Text = comment
	b, _ := json.Marshal(word)
	button.CallbackData = string(b)

	return button
}

func EditMessageReplyMarkup(reqBody *telegram.EditMessageReplyMarkupRequest) error {
	return post("/editMessageReplyMarkup", reqBody)
}

func AnswerCallbackQuery(reqBody *telegram.AnswerCallbackQuery) error {
	return post("/answerCallbackQuery", reqBody)
}

func SetMyCommands(reqBody *telegram.SetMyCommandsRequest) error {
	return post("/setMyCommands", reqBody)
}

func post(url string, payload interface{}) error {
	reqBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: LoggingRoundTripper{http.DefaultTransport},
		Timeout:   30 * time.Second,
	}
	res, err := client.Post(botHost+url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	fmt.Printf("Sending request to %v\n", req.URL)

	res, e = lrt.Proxied.RoundTrip(req)

	// Handle the result.
	if e != nil {
		fmt.Printf("Error: %v", e)
	} else if res.StatusCode >= 400 && res.StatusCode <= 500 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Received %s error\n", string(body))
	}

	return
}
