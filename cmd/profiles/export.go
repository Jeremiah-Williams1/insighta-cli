// cmd/profiles/export.go
package profiles

import (
	"fmt"
	"insighta-cli/auth"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	flagExportFormat  string
	flagExportGender  string
	flagExportCountry string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export profiles as CSV",
	Run:   runExport,
}

func init() {
	ProfilesCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&flagExportFormat, "format", "", "Export format (csv)")
	exportCmd.Flags().StringVar(&flagExportGender, "gender", "", "Filter by gender")
	exportCmd.Flags().StringVar(&flagExportCountry, "country", "", "Filter by country code")
	exportCmd.MarkFlagRequired("format")
}

func runExport(cmd *cobra.Command, args []string) {
	if flagExportFormat != "csv" {
		fmt.Println("Error: only --format csv is supported")
		os.Exit(1)
	}

	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	params := "?format=csv"
	if flagExportGender != "" {
		params += "&gender=" + flagExportGender
	}
	if flagExportCountry != "" {
		params += "&country_id=" + flagExportCountry
	}

	resp, err := auth.MakeRequest("GET", apiBase+"/api/profiles/export"+params, nil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("profiles_%s.csv", time.Now().Format("20060102_150405"))
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		os.Exit(1)
	}

	fmt.Printf("Exported to %s\n", filename)
}
