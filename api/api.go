package api

import (
	"WordsBot/api/seed"
	"WordsBot/api/telegram"
	"WordsBot/api/wordsImport"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetApiRouter() http.Handler {
	r := chi.NewRouter()

	r.Mount("/telegram", telegram.NewRouter())
	r.Mount("/wordsImport", wordsImport.NewRouter())
	r.Mount("/seed", seed.NewRouter())

	return r
}
