package seed

import (
	"WordsBot/services/sqlService"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/queue", queue)

	return r
}

func queue(writer http.ResponseWriter, request *http.Request) {
	sqlService.Queue()
}
