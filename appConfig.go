package main

type appConfig struct {
	Bot bot `json:"Bot" cfg:"Bot"`
}

type bot struct {
	id string `json:"id" cfg:"id"`
}
