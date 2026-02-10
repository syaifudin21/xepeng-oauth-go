package oauth

import "fmt"

type Config struct {
	BaseURL      string
	APIBaseURL   string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	AutoRefresh  bool
	RefreshBuffer int // in seconds
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ExpiresAt    int64  `json:"-"`
}

type OAuthError struct {
	Message    string
	Code       string
	StatusCode int
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("OAuthError: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
}

var DefaultConfig = Config{
	BaseURL:       "https://staging-app.xepeng.com",
	APIBaseURL:    "https://staging-api.xepeng.com",
	Scopes:        []string{"profile", "email"},
	AutoRefresh:   true,
	RefreshBuffer: 300,
}
