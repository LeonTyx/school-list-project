package database

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"os"
)

var (
	DBCon        *sql.DB
	SessionStore *pgstore.PGStore
)

func InitOauthStore() {
	var err error

	SessionStore, err = pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("DATABASE_SECRET")))
	if err != nil {
		panic(err)
	}

	SessionStore.MaxAge(1800)
	SessionStore.Options.HttpOnly = true
	if os.Getenv("ENV") == "DEV" {
		SessionStore.Options.Secure = false
	} else {
		SessionStore.Options.Secure = true
		SessionStore.Options.Domain = os.Getenv("DOMAIN")
	}
	fmt.Println("Successful oauth store connection!")
}

func InitDB() {
	var err error

	// Open up our database connection.
	DBCon, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println("Cannot open SQL connection")
		panic(err.Error())
	}
	fmt.Println("Successful database connection!")
}
