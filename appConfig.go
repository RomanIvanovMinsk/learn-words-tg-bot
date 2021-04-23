package main

type appConfig struct {
	Bot bot `json:"Bot" cfg:"Bot"`
}

type bot struct {
	token string `json:"token" cfg:"token"`
}
