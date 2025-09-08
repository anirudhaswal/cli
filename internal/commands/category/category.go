package category

import "github.com/spf13/cobra"

var CategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage preference categories",
	Long:  "Manage preference categories",
}
