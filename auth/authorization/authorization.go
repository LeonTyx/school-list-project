package authorization

import (
	"context"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"school-list-project/database"
)

type Resource struct {
	Resource string
	Policy   Policy
}

type Policy struct {
	CanAdd    bool
	CanView   bool
	CanEdit   bool
	CanDelete bool
}

func ValidSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := database.SessionStore.Get(r, "session")
		if err != nil {
			RespondWithError(w, r, 500, "The server was unable to retrieve this session")
			return
		}

		if session.ID == "" {
			RespondWithError(w, r, 401, "This user has no current session. Use of this endpoint is thus unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := r.Context().Value("resource")
		session, err := database.SessionStore.Get(r, "session")
		if err != nil {
			RespondWithError(w, r, 500, "The server was unable to retrieve this session")
			return
		}
		googleID := session.Values["GoogleId"]

		rows, err := database.DBCon.Query("SELECT rrb.can_add, rrb.can_view, rrb.can_edit, rrb.can_delete  FROM account a INNER JOIN role ro ON a.role_id = ro.role_id LEFT JOIN resource_role_bridge rrb ON rrb.role_id = a.role_id INNER JOIN resource rsc ON rrb.resource_id = rsc.resource_id WHERE google_id = $1 AND rsc.resource_name=$2", googleID, resource)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var resourcePolicy Resource
		for rows.Next() {
			err := rows.Scan(&resourcePolicy.Policy.CanAdd, &resourcePolicy.Policy.CanView,
				&resourcePolicy.Policy.CanEdit, &resourcePolicy.Policy.CanDelete)
			if err != nil {
				RespondWithError(w, r, 500, "The server was unable to retrieve permission")
				return
			}
		}

		ctx := context.WithValue(r.Context(), "resourcePolicy", resourcePolicy)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CanView(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourcePolicy := r.Context().Value("resourcePolicy").(Resource)

		if !resourcePolicy.Policy.CanView {
			RespondWithError(w, r, 403, "You do not have access to this endpoint.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CanAdd(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourcePolicy := r.Context().Value("resourcePolicy").(Resource)

		if !resourcePolicy.Policy.CanAdd {
			RespondWithError(w, r, 403, "You do not have access to this endpoint.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CanEdit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourcePolicy := r.Context().Value("resourcePolicy").(Resource)

		if !resourcePolicy.Policy.CanEdit {
			RespondWithError(w, r, 403, "You do not have access to this endpoint.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CanDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourcePolicy := r.Context().Value("resourcePolicy").(Resource)

		if !resourcePolicy.Policy.CanDelete {
			RespondWithError(w, r, 403, "You do not have access to this endpoint.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RespondWithError(w http.ResponseWriter, r *http.Request, statusCode int, errorMsg string) {
	render.Status(r, statusCode)
	render.JSON(w, r, Error{
		StatusCode:   statusCode,
		ErrorMessage: errorMsg,
	})
}

type Error struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
}
