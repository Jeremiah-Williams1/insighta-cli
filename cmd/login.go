package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"insighta-cli/auth"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via GitHub OAuth",
	Run:   runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	// 1. Generate state and PKCE verifier
	state := generateRandomString()
	codeVerifier := generateRandomString()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// 2. Start a local server to catch the callback
	codeChan := make(chan string, 1)
	stateChan := make(chan string, 1)

	server := &http.Server{Addr: ":9999"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		stateChan <- r.URL.Query().Get("state")
		codeChan <- r.URL.Query().Get("code")
		fmt.Fprintln(w, "Login successful! You can close this tab.")
		go server.Close()
	})

	go server.ListenAndServe()

	// 3. Build GitHub OAuth URL and open in browser
	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	params := url.Values{}
	params.Set("client_id", os.Getenv("CLIENT_ID"))
	params.Set("redirect_uri", "http://localhost:9999/callback")
	params.Set("scope", "user:email")
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	authURL := "https://github.com/login/oauth/authorize?" + params.Encode()
	fmt.Println("Opening browser for GitHub login...")
	openBrowser(authURL)

	// 4. Wait for callback
	returnedState := <-stateChan
	code := <-codeChan

	// 5. Validate state
	if returnedState != state {
		fmt.Println("Error: state mismatch, possible CSRF attack")
		os.Exit(1)
	}

	// 6. Send code + verifier to your backend
	resp, err := http.PostForm(apiBase+"/auth/github/callback", url.Values{
		"code":          {code},
		"code_verifier": {codeVerifier},
	})
	if err != nil {
		fmt.Println("Error contacting backend:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	// 7. Save tokens
	err = auth.SaveTokens(result.AccessToken, result.RefreshToken)
	if err != nil {
		fmt.Println("Error saving tokens:", err)
		os.Exit(1)
	}

	fmt.Println("Logged in successfully!")
}

func generateRandomString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func openBrowser(url string) {
	exec.Command("xdg-open", url).Start()
	fmt.Printf("If browser didn't open, visit:\n%s\n", url)
}
