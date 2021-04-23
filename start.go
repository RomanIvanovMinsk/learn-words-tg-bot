package main

import (
	"WordsBot/api"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gosidekick/goconfig"
)

// Handler This handler is called everytime telegram sends us a webhook event
// func Handler(res http.ResponseWriter, req *http.Request) {
// 	if strings.Contains(strings.ToLower(req.URL.String()), "uploadlist") {
// 		//actions.UploadList(req.Body)
// 	}

// 	fmt.Println("Start handler")
// 	body := &wr.WebhookReqBody{}
// 	helper.DecodeRequestBody(req, body)
// 	helper.SelectAction(body)

// }

// FInally, the main funtion starts our server on port 3000
func main() {

	config := appConfig{}
	goconfig.File = "config.json"
	err := goconfig.Parse(&config)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := chi.NewRouter()
	router.Mount("/", api.GetApiRouter())

	server := http.Server{
		Addr:    ":8821",
		Handler: router,
	}

	go func() {
		fmt.Println("Server is listening...")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
