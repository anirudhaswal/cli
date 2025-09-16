package event

import "github.com/spf13/cobra"

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Manage events",
	Long:  "Manage events",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
