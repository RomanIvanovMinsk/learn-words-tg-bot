package sqlService

import (
	"WordsBot/models"
	"context"
	"errors"
	"fmt"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"log"
	"net/url"
)

var db *sqlx.DB

type DbWord struct {
	Id    mssql.UniqueIdentifier `db:"Id"`
	Word  string                 `db:"Word"`
	Stem  string                 `db:"Stem"`
	Lang  string                 `db:"Lang"`
	Usage string                 `db:"Usage"`
}

func OpenConnection(url *url.URL) error {
	var err error

	db, err = sqlx.Open("sqlserver", url.String())
	if err != nil {
		log.Println("Error creating connection pool: ", err.Error())
		return err
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Db Connected!\n")

	return nil
}

func CreateProfile(telegramId string) (string, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	m := map[string]interface{}{"TelegramId": telegramId}
	rows, err := db.NamedQuery(`exec dbo.CreateProfile :TelegramId`, m)

	userId, err := readUserId(err, rows)
	rows.Close()

	if err != nil {
		return "", err
	}

	return userId.String(), nil
}

func readUserId(err error, rows *sqlx.Rows) (mssql.UniqueIdentifier, error) {
	if err != nil {
		log.Println(err.Error())
		return mssql.UniqueIdentifier{}, err
	}

	hasRows := rows.Next()
	if hasRows == false {
		log.Println("Error executing query")
		return mssql.UniqueIdentifier{}, errors.New("no results")
	}
	var userId = mssql.UniqueIdentifier{}
	rows.Scan(&userId)
	return userId, nil
}

func GetUserIdByTelegramId(telegramId string) (string, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	m := map[string]interface{}{"TelegramId": telegramId}
	rows, err := db.NamedQuery(`exec dbo.GetUserIdByTelegramId :TelegramId`, m)

	userId, err := readUserId(err, rows)
	rows.Close()

	if err != nil {
		return "", err
	}

	return userId.String(), nil
}

func AddWord(userId string, word models.Word) error {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	preparedUserId := getUserId(userId)
	tx := db.MustBegin()
	stmt, err := tx.Preparex(`exec dbo.AddWord @p1, @p2, @p3, @p4, @p5`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	for _, usage := range word.Usages {
		stmt.MustExec(preparedUserId, word.Stem, word.Word, word.Lang, usage.Usage)
	}
	if word.Usages == nil || len(word.Usages) == 0 {
		stmt.MustExec(preparedUserId, word.Stem, word.Word, word.Lang, nil)
	}
	tx.Commit()

	return nil
}

func getUserId(userId string) mssql.UniqueIdentifier {
	preparedUserId := mssql.UniqueIdentifier{}
	preparedUserId.Scan(userId)
	return preparedUserId
}

func Queue() error {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = db.Exec(`exec dbo.QueueWords`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetIntervalWords(userId string) (*models.Word, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	preparedUserId := getUserId(userId)

	words := make([]DbWord, 0)
	tsql := fmt.Sprintf("exec dbo.GetIntervalWords '%s'", preparedUserId)
	err = db.Select(&words, tsql, preparedUserId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	word := models.Word{
		Id:     words[0].Id.String(),
		Word:   words[0].Word,
		Stem:   words[0].Stem,
		Lang:   words[0].Lang,
		Usages: make([]models.Usage, 0, len(words)),
	}

	for _, v := range words {
		word.Usages = append(word.Usages, models.Usage{Usage: v.Usage})
	}

	return &word, nil
}

func Answer(userId string, remember bool) error {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	preparedUserId := getUserId(userId)

	_, err = db.Exec("exec dbo.Answer @p1, @p2", preparedUserId, remember)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetUsages(userId string, word string, offset int) ([]string, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//preparedUserId := getUserId(userId)

	usages := make([]string, 0)

	wordId := mssql.UniqueIdentifier{}
	wordId.Scan(word)

	err = db.Select(&usages, "exec dbo.GetWordUsages @p1, @p2", wordId, offset)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return usages, nil
}
