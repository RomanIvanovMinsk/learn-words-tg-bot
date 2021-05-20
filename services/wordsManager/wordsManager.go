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

func GetIntervalWords(telegramUserId string) (*models.Word, error) {
	userId, err := sqlService.GetUserIdByTelegramId(telegramUserId)
	if err != nil {
		return nil, err
	}

	return sqlService.GetIntervalWords(userId)
}

func Answer(telegramUserId string, remember bool) error {
	userId, err := sqlService.GetUserIdByTelegramId(telegramUserId)
	if err != nil {
		return err
	}

	return sqlService.Answer(userId, remember)
}

func GetUsages(telegramUserId string, word string, offset int) ([]string, error) {
	userId, err := sqlService.GetUserIdByTelegramId(telegramUserId)
	if err != nil {
		return nil, err
	}

	return sqlService.GetUsages(userId, word, offset)
}
