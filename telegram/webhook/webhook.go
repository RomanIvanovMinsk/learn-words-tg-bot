package webhook

import (
	"WordsBot/models/telegram"
	"WordsBot/telegram/commands/callback"
	"WordsBot/telegram/commands/givemetheword"
	"WordsBot/telegram/commands/import_utility"
	"WordsBot/telegram/commands/myid"
	"WordsBot/telegram/commands/start"
	"strings"
)

func HandleWebhook(body *telegram.WebhookReqBody) error {
	command := strings.ToLower(body.Message.Text)

	if command == "/myid" {
		err := myid.Handle(body)
		return err
	}

	if command == start.Name() {
		err := start.Handle(body)
		return err
	}
	if command == givemetheword.Name() {
		err := givemetheword.Handle(body)
		return err
	}

	if command == import_utility.Name() {
		err := import_utility.Handle(body)
		return err
	}

	if body.Callback.Data != "" {
		err := callback.Handle(body)
		return err
	}
	return nil
}
