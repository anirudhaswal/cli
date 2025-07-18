package schema

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var schemaCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit schema from draft to live",
	Long:  `Commit schema from draft to live in a workspace`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("workspace & schema slug argument is required for schemas.")
		}
		workspace := args[0]
		slug := args[1]

		mgmnt_client := utils.GetSuprSendMgmntClient()

		err := mgmnt_client.FinalizeSchema(workspace, slug, true)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		log.Printf("Committed schema: %s\n", slug)
	},
}

func init() {
	SchemaCmd.AddCommand(schemaCommitCmd)
}
