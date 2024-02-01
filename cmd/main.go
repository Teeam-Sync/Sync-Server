package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// "os"

	grpc_handler "github.com/Teeam-Sync/Sync-Server/api/handler"
	database "github.com/Teeam-Sync/Sync-Server/internal/database/mongodb"
	"github.com/Teeam-Sync/Sync-Server/internal/logger"
	"github.com/gorilla/pat"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	// "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

const (
	ENV_FILE = ".env"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name: "oauthstate", Value: state, Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthstate.Value {
		logger.Debug("invalid google oauth state cookie:", oauthstate.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		logger.Debug(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data))
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange,", err.Error())
	}

	res, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo,", err.Error())
	}

	return io.ReadAll(res.Body)
}

func main() {
	initialize()
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)
	// app := fiber.New()
	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
	// port := os.Getenv("PORT")
	// app.Post("/google_login", controllers.GoogleLogin)
	// app.Post("/google_callback", controllers.GoogleCallback)
	// app.Listen(fmt.Sprintf(":%s", port))
}

func init() {
	defer func() { // panic 발생 시 recover
		if r := recover(); r != nil {
			log.Println("Recovered from panic during initialization:", r)
		}
	}()

	err := godotenv.Load(ENV_FILE)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initialize() {
	database.Initialize()
	grpc_handler.Initialize()
}
