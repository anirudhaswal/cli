package profiles

import (
	"github.com/spf13/cobra"
)

var ProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Profile",
	Long:  "Manage Profile and store credentials",
}
