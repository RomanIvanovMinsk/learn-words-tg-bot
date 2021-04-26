package helpers

import (
	"WordsBot/services/sqlService"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	wr "WordsBot/models"

	action "WordsBot/actions"
)

// Decode web request body
func DecodeRequestBody(req *http.Request, body *wr.WebhookReqBody) {
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
	}
}

func DecodeListForSave(req *http.Request, body *wr.WordsList) {
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
	}
}

func SelectAction(body *wr.WebhookReqBody) {
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

	if strings.Contains(command, "myid") {
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
}
