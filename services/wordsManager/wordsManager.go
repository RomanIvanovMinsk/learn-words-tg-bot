package wordsManager

import (
	"WordsBot/models"
	"WordsBot/services/sqlService"
)

func AddWordsList(telegramUserId string, words []models.Word) error {

	userId, err := sqlService.GetUserIdByTelegramId(telegramUserId)
	if err != nil {
		return err
	}

	for _, wordInfo := range words {
		err := sqlService.AddWord(userId, wordInfo)
		if err != nil {
			return err
		}
	}

	return nil
}
