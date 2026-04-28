package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"insighta-cli/auth"
	"strings"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Shows the UserId and Role",
	Run:   runwhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func runwhoami(cmd *cobra.Command, args []string) {
	cred, err := auth.LoadTokens()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	claims, err := decodeJWT(cred.AccessToken)
	fmt.Printf("ID: %s\nRole: %s\n", claims["id"], claims["role"])
}

func decodeJWT(tokenString string) (map[string]any, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	// Decode the payload (middle part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims map[string]any
	err = json.Unmarshal(payload, &claims)
	return claims, err
}
