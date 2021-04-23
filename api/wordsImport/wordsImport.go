package wordsImport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func RenderYearChart(w http.ResponseWriter, r *http.Request) {

	render.PlainText(w, r, "I'm wordImport route")
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", RenderYearChart)

	return r
}
