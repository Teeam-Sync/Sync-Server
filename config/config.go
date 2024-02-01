package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {
	err := os.Getenv(".env")
	if err != "" {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://localhost:8080/google_callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
