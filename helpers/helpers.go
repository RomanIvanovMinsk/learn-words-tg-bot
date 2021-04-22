package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	action "../actions"
	wr "../models"
)

// Decode web request body
func DecodeRequestBody(req *http.Request, body *wr.WebhookReqBody) {
	// First, decode the JSON response body

	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
	}

	fmt.Println(body.Message.Text)
}

func SelectAction(body *wr.WebhookReqBody) {
	// Check if the message contains the word "marco"
	// if not, return without doing anything
	if strings.Contains(strings.ToLower(body.Message.Text), "marco") {
		// If the text contains marco, call the `sayPolo` function, which
		// is defined below
		if err := action.SayPolo(body.Message.Chat.ID); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}

		// log a confirmation message if the message is sent successfully
		fmt.Println("reply sent")
	}

	if strings.Contains(strings.ToLower(body.Message.Text), "myid") {
		fmt.Println("Start my id")
		if err := action.GetMyId(body.Message.Chat.ID); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	}
}
