package schema

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var schemaResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset schema from live to draft",
	Long:  `Reset schema from live to draft in a workspace`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("schema slug argument is required for schemas.")
		}
		slug := args[0]

		workspace, _ := cmd.Flags().GetString("workspace")

		mgmnt_client := utils.GetSuprSendMgmntClient()

		err := mgmnt_client.FinalizeSchema(workspace, slug, false)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaResetCmd)
}
