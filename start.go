package main

import (
	"fmt"
	"net/http"
	"strings"

	helper "./helpers"
	wr "./models"
)

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(strings.ToLower(req.Url), "uploadlist") {
		//actions.UploadList(req.Body)
	}

	fmt.Println("Start handler")
	body := &wr.WebhookReqBody{}
	helper.DecodeRequestBody(req, body)
	helper.SelectAction(body)

}

// FInally, the main funtion starts our server on port 3000
func main() {
	http.ListenAndServe(":88", http.HandlerFunc(Handler))
}
