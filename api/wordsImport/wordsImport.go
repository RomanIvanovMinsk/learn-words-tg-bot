package wordsImport

import (
	"WordsBot/services/wordsManager"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	clienthelper "WordsBot/clienthelper"
	wr "WordsBot/models"
)

func RenderYearChart(w http.ResponseWriter, r *http.Request) {
	words := make([]wr.Word, 0, 1)
	words = append(words, wr.Word{
		Word:   "Good",
		Stem:   "good",
		Lang:   "en",
		Usages: []wr.Usage{{Usage: "that is Good 1"}},
	})
	wordsManager.AddWordsList("729006239", words)
	render.PlainText(w, r, "I'm wordImport route")
}

func importWordsList(w http.ResponseWriter, req *http.Request) {
	body := &wr.WordsList{}
	err := render.DecodeJSON(req.Body, body)
	if err != nil {
		fmt.Println("could not decode request body", err)
	}
	fmt.Println(body)
	clienthelper.UploadList(body, w)
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", RenderYearChart)
	r.Post("/importList", importWordsList)

	return r
}
