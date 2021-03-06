package main

import (
	"WordsBot/api"
	"WordsBot/config"
	"WordsBot/models/telegram"
	"WordsBot/services/sqlService"
	telegramApi "WordsBot/telegram/api"
	"WordsBot/telegram/commands/givemetheword"
	"WordsBot/telegram/commands/import_utility"
	"WordsBot/telegram/commands/myid"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/httplog"
	"github.com/gosidekick/goconfig"
	_ "github.com/gosidekick/goconfig/json"
)

var Config *config.AppConfig

func main() {

	Config = &config.AppConfig{}
	goconfig.File = "config.json"
	err := goconfig.Parse(Config)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	initDb(Config.Sql)

	configureWebhooks(Config)
	telegramApi.Configure(config.GetBotHost(Config))
	err = telegramApi.SetMyCommands(&telegram.SetMyCommandsRequest{
		Commands: []telegram.BotCommand{
			//{Command: "/start", Description: "create user profile"},
			{Command: myid.Name(), Description: myid.Description()},
			{Command: givemetheword.Name(), Description: givemetheword.Description()},
			{Command: import_utility.Name(), Description: import_utility.Description()},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON),
		middleware.DefaultLogger,
		middleware.RequestID,
		middleware.CleanPath)
	router.Use(middleware.Recoverer)

	router.Mount("/api/", api.GetApiRouter())

	server := http.Server{
		Addr:    ":" + Config.Port,
		Handler: router,
	}

	go func() {
		fmt.Println("Server is listening...")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	gracefulShutDown(server)
}

func gracefulShutDown(server http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func configureWebhooks(config *config.AppConfig) {
	log.Println("Configuring webhook")
	_, err := http.Get(config.Bot.Host + "/bot" + config.Bot.Token + "/setWebhook?url=" + config.Host + "/api/telegram")
	if err != nil {
		log.Fatal(err)
	}
}

func initDb(config config.SqlConfig) {
	query := url.Values{}
	query.Add("database", config.Database)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.User, config.Password),
		Host:     config.Host,
		RawQuery: query.Encode(),
	}
	err := sqlService.OpenConnection(u)

	if err != nil {
		fmt.Println(err.Error())
	}
}
