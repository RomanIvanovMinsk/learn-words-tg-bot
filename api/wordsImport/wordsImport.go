package wordsImport

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	clienthelper "WordsBot/clienthelper"
	helper "WordsBot/helpers"
	wr "WordsBot/models"
)

func RenderYearChart(w http.ResponseWriter, r *http.Request) {

	render.PlainText(w, r, "I'm wordImport route")
}

func importWordsList(w http.ResponseWriter, req *http.Request) {
	body := &wr.WordsList{}
	helper.DecodeListForSave(req, body)
	fmt.Println(body)
	clienthelper.UploadList(body, w)
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", RenderYearChart)
	r.Post("/importList", importWordsList)

	return r
}
