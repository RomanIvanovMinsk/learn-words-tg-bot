package main

import (
	"WordsBot/actions"
	"WordsBot/api"
	"WordsBot/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
)

var Config *config.AppConfig

// Finally, the main function starts our server on port 3000
func main() {

	Config = &config.AppConfig{}
	goconfig.File = "config.json"
	err := goconfig.Parse(Config)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	configureWebhooks(Config)
	actions.Configure(config.GetBotHost(Config))

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

func configureWebhooks(config *config.AppConfig) {
	log.Println("Configuring webhook")
	_, err := http.Get(config.Bot.Host + "/bot" + config.Bot.Token + "/setWebhook?url=" + config.Host + "/telegram")
	if err != nil {
		log.Fatal(err)
	}
}
