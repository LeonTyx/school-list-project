package oauth

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type Profile struct {
	Email string
	Name  string
}

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/profile", GetProfile)
	return router
}

func GetProfile(w http.ResponseWriter, r *http.Request){
	profile := Profile{
		Email:  "userID",
		Name: "Senor Namey",
	}
	render.JSON(w, r, profile) // A chi router helper for serializing and returning json
}