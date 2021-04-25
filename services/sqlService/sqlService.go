package sqlService

import (
	"context"
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
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	defer rows.Close()
	hasRows := rows.Next()
	if hasRows == false {
		log.Println("Error executing query")
		return "", err
	}
	var userId = mssql.UniqueIdentifier{}
	rows.Scan(&userId)

	return userId.String(), nil
}
