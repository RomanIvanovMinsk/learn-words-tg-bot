package config

type AppConfig struct {
	Bot  Bot       `json:"Bot" cfg:"Bot"`
	Host string    `json:"Host" cfg:"Host"`
	Sql  SqlConfig `json:"Sql" cfg:"Sql"`
	Port string    `json:"Port" cfg:"Port"`
}

type Bot struct {
	Token string `json:"Token" cfg:"Token"`
	Host  string `json:"Host" cfg:"Host"`
}

type SqlConfig struct {
	Host     string `json:"host" cfg:"host"`
	Database string `json:"database" cfg:"database"`
	User     string `json:"user" cfg:"user"`
	Password string `json:"password" cfg:"password"`
}

func GetBotHost(config *AppConfig) string {
	return config.Bot.Host + "/bot" + config.Bot.Token
}
