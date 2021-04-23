package telegram

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "I'm telegram route")
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", GetRoot)

	return r
}
