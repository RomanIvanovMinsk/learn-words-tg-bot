package main

import (
	"fmt"
	"net/http"
	"strings"

	helper "./helpers"
	wr "./models"
	clienthelper "./clienthelper"
	
)

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(strings.ToLower(req.URL.Path), "uploadlist") {
		body := &wr.WordsList{}
		helper.DecodeListForSave(req, body);
		fmt.Println(body);
		clienthelper.UploadList(body, res);
	}

	fmt.Println("Start handler")
	body := &wr.WebhookReqBody{}
	helper.DecodeRequestBody(req, body)
	helper.SelectAction(body)

}

// Finally, the main function starts our server on port 3000
func main() {
	http.ListenAndServe(":88", http.HandlerFunc(Handler))
}
