package api

import (
	"WordsBot/api/telegram"
	"WordsBot/api/wordsImport"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetApiRouter() http.Handler {
	r := chi.NewRouter()

	r.Mount("/telegram", telegram.NewRouter())
	r.Mount("/wordsImport", wordsImport.NewRouter())

	return r
}
