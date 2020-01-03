package database

import (
	"database/sql"
	"github.com/antonlindstrom/pgstore"
	"golang.org/x/oauth2"
)

var (
	DBCon             *sql.DB
	SessionStore      *pgstore.PGStore
	GoogleOauthConfig *oauth2.Config
)
