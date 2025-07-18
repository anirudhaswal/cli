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
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("workspace & schema slug argument is required for schemas.")
		}
		workspace := args[0]
		slug := args[1]

		mgmnt_client := utils.GetSuprSendMgmntClient()

		err := mgmnt_client.FinalizeSchema(workspace, slug, false)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		log.Printf("Reset schema: %s\n", slug)
	},
}

func init() {
	SchemaCmd.AddCommand(schemaResetCmd)
}
