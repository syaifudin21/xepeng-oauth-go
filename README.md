# Xepeng OAuth Go SDK

SDK resmi Go untuk integrasi dengan Xepeng OAuth Service. Paket ini mendukung alur OAuth 2.0 dengan PKCE (Proof Key for Code Exchange) untuk keamanan maksimal.

## Fitur

- ✅ OAuth 2.0 Authorization Code Flow dengan PKCE (S256).
- ✅ Penukaran Authorization Code menjadi Access & Refresh Token.
- ✅ Refresh Access Token.

## Instalasi

Gunakan `go get` untuk menambahkan paket ini ke proyek Anda:

```bash
go get github.com/syaifudin21/xepeng-oauth-go
```

## Cara Penggunaan

### 1. Inisialisasi Klien

```go
import "github.com/syaifudin21/xepeng-oauth-go"

config := oauth.Config{
    ClientID:     "CLIENT_ID_ANDA",
    ClientSecret: "CLIENT_SECRET_ANDA",
    RedirectURI:  "https://aplikasi-anda.com/callback",
    Scopes:       []string{"profile", "email"},
}

client := oauth.NewClient(config)
```

### 2. Mendapatkan URL Otorisasi

Gunakan ini untuk mengarahkan pengguna ke halaman login Xepeng.

```go
authURL, state, codeVerifier, err := client.GetAuthorizationURL()
if err != nil {
    // Tangani error
}

// SIMPAN 'state' dan 'codeVerifier' di session atau storage aman (cookie/redis)
// Kemudian redirect user ke 'authURL'
```

### 3. Menukar Code dengan Token

Setelah pengguna diarahkan kembali ke `RedirectURI` Anda, ambil parameter `code` dari URL.

```go
// Ambil 'code' dari query parameter
// Ambil 'codeVerifier' yang sebelumnya disimpan

tokens, err := client.ExchangeCode(code, codeVerifier)
if err != nil {
    // Tangani error
}

fmt.Printf("Access Token: %s\n", tokens.AccessToken)
fmt.Printf("Refresh Token: %s\n", tokens.RefreshToken)
fmt.Printf("Scope: %s\n", tokens.Scope)
fmt.Printf("Client ID: %s\n", tokens.ClientID)
fmt.Printf("Client Secret: %s\n", tokens.ClientSecret) // Berisi hashed secret atau secret aplikasi
```

### 4. Refresh Token

```go
newTokens, err := client.RefreshToken(tokens.RefreshToken)
if err != nil {
    // Tangani error
}
```

## Lisensi

MIT
