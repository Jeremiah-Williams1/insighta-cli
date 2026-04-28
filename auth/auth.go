package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// This struct is what gets saved to disk as JSON
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// credentialsPath returns ~/.insighta/credentials.json
func credentialsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".insighta", "credentials.json"), nil
}

// SaveTokens writes tokens to disk
func SaveTokens(accessToken, refreshToken string) error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	os.MkdirAll(filepath.Dir(path), 0700)

	creds := Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// LoadTokens reads tokens from disk
func LoadTokens() (Credentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return Credentials{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}, err
	}

	var creds Credentials
	err = json.Unmarshal(data, &creds)
	return creds, err
}

// ClearTokens deletes the credentials file (logout)
func ClearTokens() error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func MakeRequest(method, url string, body io.Reader) (*http.Response, error) {
	creds, err := LoadTokens()
	if err != nil {
		return nil, fmt.Errorf("not logged in, run: insighta login")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	req.Header.Set("X-API-Version", "1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
