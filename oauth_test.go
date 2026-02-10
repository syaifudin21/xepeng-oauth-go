package oauth

import (
	"net/url"
	"testing"
)

func TestPKCE(t *testing.T) {
	verifier, err := GenerateCodeVerifier(64)
	if err != nil {
		t.Fatalf("Failed to generate verifier: %v", err)
	}
	if len(verifier) != 64 {
		t.Errorf("Expected length 64, got %d", len(verifier))
	}

	challenge := GenerateCodeChallenge(verifier)
	if challenge == "" {
		t.Error("Challenge should not be empty")
	}

	state := GenerateState()
	if len(state) != 43 {
		t.Errorf("Expected state length 43, got %d", len(state))
	}
}

func TestGetAuthorizationURL(t *testing.T) {
	config := Config{
		BaseURL:     "https://test.com",
		ClientID:    "client123",
		RedirectURI: "https://app.com/callback",
		Scopes:      []string{"openid", "profile"},
	}
	client := NewClient(config)

	authURL, state, verifier, err := client.GetAuthorizationURL()
	if err != nil {
		t.Fatalf("Failed to get auth URL: %v", err)
	}

	if state == "" || verifier == "" {
		t.Error("State and verifier should not be empty")
	}

	parsedURL, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("Failed to parse auth URL: %v", err)
	}

	if parsedURL.Host != "test.com" {
		t.Errorf("Expected host test.com, got %s", parsedURL.Host)
	}

	q := parsedURL.Query()
	if q.Get("client_id") != "client123" {
		t.Errorf("Expected client123, got %s", q.Get("client_id"))
	}
	if q.Get("redirect_uri") != "https://app.com/callback" {
		t.Errorf("Expected redirect_uri, got %s", q.Get("redirect_uri"))
	}
	if q.Get("scope") != "openid profile" {
		t.Errorf("Expected scope 'openid profile', got %s", q.Get("scope"))
	}
	if q.Get("state") != state {
		t.Errorf("Expected state %s, got %s", state, q.Get("state"))
	}
	
	expectedChallenge := GenerateCodeChallenge(verifier)
	if q.Get("code_challenge") != expectedChallenge {
		t.Errorf("Expected challenge %s, got %s", expectedChallenge, q.Get("code_challenge"))
	}
}

func TestShouldRefreshToken(t *testing.T) {
	client := NewClient(Config{RefreshBuffer: 300})
	
	tokens := &TokenResponse{
		ExpiresAt: 0, // Expired long ago
	}
	if !client.ShouldRefreshToken(tokens) {
		t.Error("Should need refresh for expired tokens")
	}

	// Not expired
	tokens.ExpiresAt = 1e15 // Far in the future
	if client.ShouldRefreshToken(tokens) {
		t.Error("Should not need refresh for future tokens")
	}
}
