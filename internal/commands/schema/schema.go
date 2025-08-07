package schema

import "github.com/spf13/cobra"

var SchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage schema",
	Long:  `Manage schema`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
