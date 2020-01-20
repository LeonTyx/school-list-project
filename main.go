package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"school-list-project/api/school-list"
	"school-list-project/api/supplies"
	supply_list "school-list-project/api/supply-list"
	"school-list-project/database"
	"school-list-project/oauth"
	_ "school-list-project/oauth"
	"strings"
)

func ForceSsl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "production" {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}
		}

		next.ServeHTTP(w, r)
	})

}
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
		ForceSsl,
	)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/schools", school_list.Routes())
		r.Mount("/supply_lists", supply_list.Routes())
		r.Mount("/supplies", supplies.Routes())
	})

	router.Route("/oauth/v1", func(r chi.Router) {
		r.Mount("/", oauth.Routes())
	})

	workDir, _ := os.Getwd()
	frontendDir := filepath.Join(workDir, "frontend/build")

	FileServer(router, "/", http.Dir(frontendDir))

	return router
}

func InitEnv() {
	if _, err := os.Stat("environment.env"); err == nil {
		err := godotenv.Load("environment.env")
		if err != nil {
			fmt.Println("Error loading environment.env")
		}
		fmt.Println("Current environment:", os.Getenv("ENV"))
	}
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	fmt.Println("Running on port:", port)
	return ":" + port
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func main() {
	InitEnv()
	oauth.ConfigOauth()

	database.InitOauthStore()
	database.InitDB()

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

	log.Fatal(http.ListenAndServe(GetPort(), router))
}
