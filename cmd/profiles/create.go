// cmd/profiles/create.go
package profiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"insighta-cli/auth"
	"os"

	"github.com/spf13/cobra"
)

var flagName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new profile (admin only)",
	Run:   runCreate,
}

func init() {
	ProfilesCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&flagName, "name", "", "Name for the profile (required)")
	createCmd.MarkFlagRequired("name")
}

func runCreate(cmd *cobra.Command, args []string) {
	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	body, _ := json.Marshal(map[string]string{"name": flagName})

	resp, err := auth.MakeRequest("POST", apiBase+"/api/profiles", bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		Status string         `json:"status"`
		Data   map[string]any `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != "success" {
		fmt.Println("Error creating profile")
		os.Exit(1)
	}

	fmt.Printf("Profile created:\n")
	fmt.Printf("ID:      %s\n", result.Data["id"])
	fmt.Printf("Name:    %s\n", result.Data["name"])
	fmt.Printf("Gender:  %s\n", result.Data["gender"])
	fmt.Printf("Age:     %v\n", result.Data["age"])
	fmt.Printf("Country: %s\n", result.Data["country_name"])
}
