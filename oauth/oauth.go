package oauth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"school-list-project/database"
	"strconv"
	"time"
)

var (
	GoogleOauthConfig *oauth2.Config
)

func ConfigOauth() {
	if os.Getenv("ENV") == "DEV" {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8080/oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	} else {
		GoogleOauthConfig = &oauth2.Config{
			RedirectURL:  "https://" + os.Getenv("HOST") + "/oauth/v1/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	}
}

type Error struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/login", HandleGoogleLogin)
	router.Get("/callback", HandleGoogleCallback)
	router.Get("/logout", HandleGoogleLogout)
	router.Get("/profile", GetProfile)
	router.Get("/refresh", RefreshSession)
	return router
}

func GetSeed() int64 {
	seed := time.Now().UnixNano() // A new random seed (independent from state)
	rand.Seed(seed)
	return seed
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := database.SessionStore.Get(r, "state")
	PanicOnErr(err)

	stateString := strconv.FormatInt(GetSeed(),10)
	state.Values["state"] = stateString

	err = state.Save(r, w)
	if err != nil {
		fmt.Println("Unable to store state data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := GoogleOauthConfig.AuthCodeURL(stateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling google callback")

	stateSession, err := database.SessionStore.Get(r, "state")
	if err != nil {
		RespondWithError(w, r, 500, "The server was unable to retrieve session state")
		return
	}

	userData, err := GetUserInfo(r.FormValue("state"), r.FormValue("code"), r)
	fmt.Println(userData)
	if err != nil {
		fmt.Println("Error getting content: " + err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	stateSession.Options.MaxAge = -1
	_ = stateSession.Save(r, w)
	// Add a user to user database if they don't exist
	// otherwise replace the previous access token field
	// with the new one

	if !UserExists(userData.Email) {
		CreateUser(userData)
	} else {
		ReplaceAccessToken(userData)
	}

	// set the user information
	session, err := database.SessionStore.Get(r, "session")
	PanicOnErr(err)

	session.Values["GoogleId"] = userData.GoogleId
	session.Values["Email"] = userData.Email
	session.Values["Name"] = userData.Name

	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Unable to store session data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Redirect user back to app
	fmt.Println(session, userData)
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
func CreateUser(userData User) {
	// Prepare the sql query for later
	insert, err := database.DBCon.Prepare(`INSERT INTO account (email, access_token, refresh_token, google_id, expires_in) VALUES ($1, $2, $3, $4, $5)`)
	// if there is an error inserting, handle it
	PanicOnErr(err)

	//Execute the previous sql query using data from the
	// userData struct being passed into the function
	_, err = insert.Exec(userData.Email, userData.AccessToken, userData.RefreshToken, userData.GoogleId, userData.ExpiresIn)
	PanicOnErr(err)
}

func UserExists(email string) bool {
	fmt.Println("Checking if user exist: ", email)

	// Prepare the sql query for later
	rows, err := database.DBCon.Query("SELECT COUNT(*) as count FROM account WHERE email = $1", email)
	PanicOnErr(err)

	return CheckCount(rows) > 0
}

func CheckCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		PanicOnErr(err)
	}
	return count
}

func ReplaceAccessToken(userData User) {
	fmt.Println("Replacing access token for ", userData.Email)

	_, err := database.DBCon.Query("UPDATE account SET access_token=$1, refresh_token=$2, expires_in=$3 WHERE email = $4", userData.AccessToken, userData.RefreshToken, userData.ExpiresIn, userData.Email)
	if err != nil {
		fmt.Println("Unable to update access token")
	}
}

type User struct {
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	GoogleId     string    `json:"id"`
	ExpiresIn    time.Time `json:"expires_in"`
	AccessToken  string
	RefreshToken string
}

func GetUserInfo(state string, code string, r *http.Request) (User, error) {
	var userData User
	stateSession, err := database.SessionStore.Get(r, "state")
	if err != nil {
		return userData, err
	}

	//Check if the oauth state google returned matches the one saved
	if state != stateSession.Values["state"] {
		return userData, fmt.Errorf("invalid oauth state")
	}

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return userData, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	//Send access token to google's user api in return for a users data!
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return userData, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return userData, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	err = json.Unmarshal(contents, &userData)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(userData.Email)
	userData.ExpiresIn = token.Expiry
	userData.AccessToken = token.AccessToken
	userData.RefreshToken = token.RefreshToken

	return userData, nil
}

func HandleGoogleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Attempting to expire session")

	session, err := database.SessionStore.Get(r, "session")
	if err != nil {
		RespondWithError(w, r, 500, "The server was unable to retrieve this session")
		return
	}

	fmt.Println("current session: ", session)
	fmt.Println("Is session new? ", session.IsNew)

	if session.ID != "" {
		fmt.Println("session id: ", session.ID)
		session.Options.MaxAge = -1

		err = session.Save(r, w)

		if err != nil {
			RespondWithError(w, r, 500, "The server was unable to expire this session")
		} else {
			render.JSON(w, r, `{"successful logout"}`)
		}

	} else {
		http.Redirect(w, r, "./", http.StatusTemporaryRedirect)
	}
}

type Profile struct {
	Email string
	Name  string
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	session, err := database.SessionStore.Get(r, "session")
	if err != nil {
		RespondWithError(w, r, 500, "The server was unable to retrieve this session")
		return
	}

	if session.ID != "" {
		fmt.Println("Getting cookies for profile")
		// get some session values
		Email := session.Values["Email"]
		EmailStr := fmt.Sprintf("%v", Email)
		Name := session.Values["Name"]
		NameStr := fmt.Sprintf("%v", Name)

		userData := Profile{EmailStr, NameStr}

		render.JSON(w, r, userData)
	} else {
		_, _ = w.Write([]byte("no session found"))
	}
}

func RefreshSession(w http.ResponseWriter, r *http.Request) {
	session, err := database.SessionStore.Get(r, "session")
	if err != nil {
		RespondWithError(w, r, 500, "The server was unable to retrieve this session")
		return
	}

	fmt.Println("current session: ", session)
	fmt.Println("Is session new? ", session.IsNew)

	if session.ID != "" {
		fmt.Println("session id: ", session.Values["GoogleId"])
		session.Options.MaxAge = 3600

		err = session.Save(r, w)
		if err != nil {
			RespondWithError(w, r, 500, "The server was unable to refresh this session")
		} else {
			render.JSON(w, r, `{"successful refresh"}`)
		}
	} else {
		http.Redirect(w, r, "./login", http.StatusTemporaryRedirect)
	}
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RespondWithError(w http.ResponseWriter, r *http.Request, status_code int, error_msg string) {
	render.Status(r, status_code)
	render.JSON(w, r, Error{
		StatusCode:   status_code,
		ErrorMessage: error_msg,
	})
}
