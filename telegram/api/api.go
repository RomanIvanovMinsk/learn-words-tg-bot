package api

import (
	"WordsBot/models/telegram"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var botHost string

func Configure(host string) {
	botHost = host
}

func SendMessage(reqBody *telegram.SendMessageReqBodyReply) error {
	return post("/sendMessage", reqBody)
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
