package profiles

import (
	"encoding/json"
	"fmt"
	"insighta-cli/auth"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search profiles using natural language",
	Run:   runSearch,
}

func init() {
	ProfilesCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: search query required")
		os.Exit(1)
	}

	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	query := url.QueryEscape(args[0])
	resp, err := auth.MakeRequest("GET", apiBase+"/api/profiles/search?q="+query, nil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		Status     string           `json:"status"`
		Page       int              `json:"page"`
		Limit      int              `json:"limit"`
		Total      int              `json:"total"`
		TotalPages int              `json:"total_pages"`
		Data       []map[string]any `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tGENDER\tAGE\tCOUNTRY\n")
	fmt.Fprintf(w, "--\t----\t------\t---\t-------\n")
	for _, p := range result.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
			p["id"], p["name"], p["gender"], p["age"], p["country_name"])
	}
	w.Flush()

	fmt.Printf("\nPage %d of %d | Total: %d\n", result.Page, result.TotalPages, result.Total)
}
