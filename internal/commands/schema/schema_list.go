package schema

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var schemaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schemas",
	Long:  `List schemas in a workspace`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("workspace argument is required for schemas.")
		}

		workspace := args[0]
		mgmnt_client := utils.GetSuprSendMgmntClient()

		schemas, err := mgmnt_client.ListSchema(workspace)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(schemas.Results, outputType)
	},
}

func init() {
	SchemaCmd.Flags().StringP("workspace", "w", "staging", "Workspace to pull the schemas from")
	SchemaCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Parent().HelpFunc()(cmd, args)
	})
	SchemaCmd.AddCommand(schemaListCmd)
}
