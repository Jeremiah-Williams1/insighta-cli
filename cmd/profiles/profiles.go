package profiles

import (
	"insighta-cli/cmd"

	"github.com/spf13/cobra"
)

var ProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
}

func init() {
	cmd.AddCommand(ProfilesCmd)
}
