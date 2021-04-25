package actions

import (
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
	reqBody := &wr.SendMessageReqBody{
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
	reqBody := &wr.SendMessageReqBody{
		ChatID: chatID,
		Text:   strconv.FormatInt(chatID, 10),
	}

	err := SendResponse(reqBody)

	if err != nil {
		return err
	}

	return nil
}

func SendResponse(reqBody *wr.SendMessageReqBody) error {
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

func ProcessAnswer(chatID int64, callback *wr.Callback) {

}
