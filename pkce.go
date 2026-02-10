package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
)

const (
	charset                   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	codeVerifierMinLength     = 43
	codeVerifierMaxLength     = 128
	defaultCodeVerifierLength = 64
)

func GenerateCodeVerifier(length int) (string, error) {
	if length == 0 {
		length = defaultCodeVerifierLength
	}
	if length < codeVerifierMinLength || length > codeVerifierMaxLength {
		return "", fmt.Errorf("code verifier length must be between %d and %d", codeVerifierMinLength, codeVerifierMaxLength)
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func GenerateCodeChallenge(verifier string) string {
	sha := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sha[:])
}

func GenerateState() string {
	s, _ := GenerateCodeVerifier(43)
	return s
}

// Base64UrlEncode is a helper if needed elsewhere, 
// though base64.RawURLEncoding usually suffices for PKCE.
func base64UrlEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
