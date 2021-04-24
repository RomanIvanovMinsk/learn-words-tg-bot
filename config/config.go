package config

type AppConfig struct {
	Bot  Bot    `json:"Bot" cfg:"Bot"`
	Host string `json:"Host" cfg:"Host"`
}

type Bot struct {
	Token string `json:"Token" cfg:"Token"`
	Host  string `json:"Host" cfg:"Host"`
}

func GetBotHost(config *AppConfig) string {
	return config.Bot.Host + "/bot" + config.Bot.Token
}
