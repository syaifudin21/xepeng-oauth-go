package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Config Config
}

func NewClient(config Config) *Client {
	// Merge with defaults
	if config.BaseURL == "" {
		config.BaseURL = DefaultConfig.BaseURL
	}
	if config.APIBaseURL == "" {
		config.APIBaseURL = DefaultConfig.APIBaseURL
	}
	if len(config.Scopes) == 0 {
		config.Scopes = DefaultConfig.Scopes
	}
	if config.RefreshBuffer == 0 {
		config.RefreshBuffer = DefaultConfig.RefreshBuffer
	}

	return &Client{Config: config}
}

// GetAuthorizationURL returns the URL to redirect the user to, along with the state and code verifier
func (c *Client) GetAuthorizationURL() (string, string, string, error) {
	state := GenerateState()
	verifier, err := GenerateCodeVerifier(64)
	if err != nil {
		return "", "", "", err
	}
	challenge := GenerateCodeChallenge(verifier)

	params := url.Values{}
	params.Add("client_id", c.Config.ClientID)
	params.Add("redirect_uri", c.Config.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", strings.Join(c.Config.Scopes, " "))
	params.Add("state", state)
	params.Add("code_challenge", challenge)
	params.Add("code_challenge_method", "S256")

	authURL := fmt.Sprintf("%s/oauth/authorize?%s", c.Config.BaseURL, params.Encode())
	return authURL, state, verifier, nil
}

// ExchangeCode exchanges an authorization code for tokens
func (c *Client) ExchangeCode(code, codeVerifier string) (*TokenResponse, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", c.Config.RedirectURI)
	params.Add("client_id", c.Config.ClientID)
	params.Add("client_secret", c.Config.ClientSecret)
	params.Add("code_verifier", codeVerifier)

	return c.fetchToken(params)
}

// RefreshToken refreshes tokens using a refresh token
func (c *Client) RefreshToken(refreshToken string) (*TokenResponse, error) {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", refreshToken)
	params.Add("client_id", c.Config.ClientID)
	params.Add("client_secret", c.Config.ClientSecret)

	return c.fetchToken(params)
}

// GetUserInfo fetches user info using the access token
func (c *Client) GetUserInfo(accessToken string) (map[string]interface{}, error) {
	apiURL := c.Config.APIBaseURL
	if apiURL == "" {
		apiURL = c.Config.BaseURL
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/userinfo", apiURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &OAuthError{
			Message:    "Failed to fetch user info",
			Code:       "userinfo_failed",
			StatusCode: resp.StatusCode,
		}
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// RevokeToken revokes the access token
func (c *Client) RevokeToken(accessToken string) error {
	apiURL := c.Config.APIBaseURL
	if apiURL == "" {
		apiURL = c.Config.BaseURL
	}

	params := url.Values{}
	params.Add("client_id", c.Config.ClientID)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/revoke", apiURL), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &OAuthError{
			Message:    "Failed to revoke token",
			Code:       "revoke_failed",
			StatusCode: resp.StatusCode,
		}
	}
	return nil
}

func (c *Client) fetchToken(params url.Values) (*TokenResponse, error) {
	apiURL := c.Config.APIBaseURL
	if apiURL == "" {
		apiURL = c.Config.BaseURL
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/oauth/token", apiURL), params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errData)
		
		msg := "Token request failed"
		if m, ok := errData["message"].(string); ok {
			msg = m
		}
		code := "token_error"
		if c, ok := errData["error"].(string); ok {
			code = c
		}

		return nil, &OAuthError{
			Message:    msg,
			Code:       code,
			StatusCode: resp.StatusCode,
		}
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	tokenResp.ExpiresAt = time.Now().Unix() + int64(tokenResp.ExpiresIn)
	return &tokenResp, nil
}

func (c *Client) ShouldRefreshToken(tokens *TokenResponse) bool {
	if tokens == nil {
		return true
	}
	buffer := int64(c.Config.RefreshBuffer)
	return time.Now().Unix() >= (tokens.ExpiresAt - buffer)
}
