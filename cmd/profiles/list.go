package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"insighta-cli/auth"

	"github.com/spf13/cobra"
)

var (
	flagGender   string
	flagCountry  string
	flagAgeGroup string
	flagMinAge   string
	flagMaxAge   string
	flagSortBy   string
	flagOrder    string
	flagPage     string
	flagLimit    string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles with optional filters",
	Run:   runList,
}

func init() {
	ProfilesCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&flagGender, "gender", "", "Filter by gender")
	listCmd.Flags().StringVar(&flagCountry, "country", "", "Filter by country code")
	listCmd.Flags().StringVar(&flagAgeGroup, "age-group", "", "Filter by age group")
	listCmd.Flags().StringVar(&flagMinAge, "min-age", "", "Minimum age")
	listCmd.Flags().StringVar(&flagMaxAge, "max-age", "", "Maximum age")
	listCmd.Flags().StringVar(&flagSortBy, "sort-by", "", "Sort by field")
	listCmd.Flags().StringVar(&flagOrder, "order", "", "Sort order (asc/desc)")
	listCmd.Flags().StringVar(&flagPage, "page", "1", "Page number")
	listCmd.Flags().StringVar(&flagLimit, "limit", "10", "Results per page")
}

func runList(cmd *cobra.Command, args []string) {
	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	// Build query string from flags
	params := fmt.Sprintf("?page=%s&limit=%s", flagPage, flagLimit)
	if flagGender != "" {
		params += "&gender=" + flagGender
	}
	if flagCountry != "" {
		params += "&country_id=" + flagCountry
	}
	if flagAgeGroup != "" {
		params += "&age_group=" + flagAgeGroup
	}
	if flagMinAge != "" {
		params += "&min_age=" + flagMinAge
	}
	if flagMaxAge != "" {
		params += "&max_age=" + flagMaxAge
	}
	if flagSortBy != "" {
		params += "&sort_by=" + flagSortBy
	}
	if flagOrder != "" {
		params += "&order=" + flagOrder
	}

	resp, err := auth.MakeRequest("GET", apiBase+"/api/profiles"+params, nil)
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

	// Print as table
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
