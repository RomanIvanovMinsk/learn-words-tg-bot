package clienthelper

import (
	wr "WordsBot/models"
	"WordsBot/services/wordsManager"
	"github.com/go-chi/render"
	"net/http"
)

func UploadList(list *wr.WordsList, w http.ResponseWriter, req *http.Request) {
	err := wordsManager.AddWordsList(list.TelegramUserId, list.Words)
	if err != nil {
		render.Status(req, 500)
		return
	}
	render.PlainText(w, req, "Ok")
}
