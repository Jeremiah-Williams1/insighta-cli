package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "insighta",
	Short: "Insighta Labs+ CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
