package main

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"school-list-project/database"
	"school-list-project/oauth"
	supply_list "school-list-project/supply-list"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "school-list-project/oauth"
	"school-list-project/school-list"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/schools", school_list.Routes())
		r.Mount("/api/supply_lists", supply_list.Routes())
		r.Mount("/api/oauth", oauth.Routes())
	})

	return router
}

func initEnv() {
	if _, err := os.Stat("environment.env"); err == nil {
		err := godotenv.Load("environment.env")
		if err != nil {
			fmt.Println("Error loading .env file")
		}
		fmt.Println("ENV:", os.Getenv("ENV"))
	}
	fmt.Println("ENV:", os.Getenv("ENV"))
}

func initOauthStore() {
	var err error

	database.SessionStore, err = pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("DATABASE_SECRET")))
	if err != nil {
		panic(err)
	}

	database.SessionStore.MaxAge(1800)
	if os.Getenv("ENV") == "DEV" {
		database.GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8080/v1/api/oauth/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}

		database.SessionStore.Options.Secure = false
	} else {
		database.GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "https://safe-brook-30495.herokuapp.com/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}

		database.SessionStore.Options.Secure = true
		database.SessionStore.Options.Domain = "safe-brook-30495.herokuapp.com"
	}
	fmt.Println("Successful oauth store connection!", database.SessionStore)
}

func initDB() {
	var err error

	// Open up our database connection.
	database.DBCon, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println("Cannot open SQL connection")
		panic(err.Error())
	}
	fmt.Println("Successful database connection ", database.DBCon)
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	initEnv()
	initOauthStore()
	initDB()

	defer database.DBCon.Close()
	defer database.SessionStore.Close()

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(GetPort(), router)) // Note, the port is usually gotten from the environment.

}
