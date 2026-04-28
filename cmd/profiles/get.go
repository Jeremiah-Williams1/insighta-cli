// cmd/profiles/get.go
package profiles

import (
	"encoding/json"
	"fmt"
	"insighta-cli/auth"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a profile by ID",
	Run:   runGet,
}

func init() {
	ProfilesCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: profile ID required")
		os.Exit(1)
	}

	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	resp, err := auth.MakeRequest("GET", apiBase+"/api/profiles/"+args[0], nil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Data   struct {
			ID                 string    `json:"id"`
			Name               string    `json:"name"`
			Gender             string    `json:"gender"`
			GenderProbability  float64   `json:"gender_probability"`
			Age                int       `json:"age"`
			AgeGroup           string    `json:"age_group"`
			CountryID          string    `json:"country_id"`
			CountryName        string    `json:"country_name"`
			CountryProbability float64   `json:"country_probability"`
			CreatedAt          time.Time `json:"created_at"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)
	p := result.Data

	fmt.Printf("ID:                  %s\n", p.ID)
	fmt.Printf("Name:                %s\n", p.Name)
	fmt.Printf("Gender:              %s (%.0f%%)\n", p.Gender, p.GenderProbability*100)
	fmt.Printf("Age:                 %d (%s)\n", p.Age, p.AgeGroup)
	fmt.Printf("Country:             %s (%s)\n", p.CountryName, p.CountryID)
	fmt.Printf("Created:             %s\n", p.CreatedAt.Format("2006-01-02"))
}
