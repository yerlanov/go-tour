package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	Google *oauth2.Config
}

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func LoadAuthConfig(config Config) *Auth {
	return &Auth{
		Google: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}
