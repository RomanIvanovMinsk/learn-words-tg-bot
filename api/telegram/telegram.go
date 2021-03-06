package telegram

import (
	"WordsBot/models/telegram"
	"WordsBot/telegram/webhook"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "I'm telegram route")
}

func PostWebhook(w http.ResponseWriter, r *http.Request) {

	body := &telegram.WebhookReqBody{}
	err := render.DecodeJSON(r.Body, body)
	if err != nil {
		fmt.Println("could not decode request body", err)
	}

	fmt.Println(body.Message.Text)

	webhook.HandleWebhook(body)
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", GetRoot)
	r.Post("/", PostWebhook)

	return r
}
