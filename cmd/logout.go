package cmd

import (
	"fmt"
	"insighta-cli/auth"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of insighta-cli",
	Run:   runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) {
	err := auth.ClearTokens()
	if err != nil {
		fmt.Println("Error Logging out, try again")
		return
	}
	fmt.Println("Logged out successfull")

}
