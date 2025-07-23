package profiles

import (
	"github.com/spf13/cobra"
)

var ProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage Profiles",
	Long:  "Manage Profiles and store credentials",
}
