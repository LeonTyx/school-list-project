package authorization

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"school-list-project/database"
)

func ValidSession(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := database.SessionStore.Get(r, "session")
		if err != nil {
			RespondWithError(w, r, 500, "The server was unable to retrieve this session")
			return
		}

		if session.ID != "" {
			fmt.Println("Getting cookies for profile")
			googleID := session.Values["GoogleId"]

			fmt.Println(googleID)
		} else {
			RespondWithError(w, r, 401, "This user has no current session. Use of this endpoint is thus unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}


func RespondWithError(w http.ResponseWriter, r *http.Request, status_code int, error_msg string) {
	render.Status(r, status_code)
	render.JSON(w, r, Error{
		StatusCode:   status_code,
		ErrorMessage: error_msg,
	})
}

type Error struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
}
