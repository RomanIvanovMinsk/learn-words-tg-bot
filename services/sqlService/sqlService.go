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
		log.Println(err.Error())
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
	preparedUserId := mssql.UniqueIdentifier{}
	preparedUserId.Scan(userId)
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
