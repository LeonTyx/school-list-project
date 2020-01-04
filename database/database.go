package database

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

var (
	DBCon             *sql.DB
	SessionStore      *pgstore.PGStore
	GoogleOauthConfig *oauth2.Config
)

func InitOauthStore() {
	var err error

	SessionStore, err = pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("DATABASE_SECRET")))
	if err != nil {
		panic(err)
	}

	SessionStore.MaxAge(1800)
	if os.Getenv("ENV") == "DEV" {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8080/oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}

		SessionStore.Options.Secure = false
	} else {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "https://"+os.Getenv("HOST")+"oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}

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
