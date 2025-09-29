package workflow

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var worklowEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables a workflow.",
	Long:  "Enables a workflow to activate. Example: suprsend workflow enable <slug>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("Category slug is required.")
			return
		}
		workspace, _ := cmd.Flags().GetString("workspace")
		slug := args[0]

		mgmntClient := utils.GetSuprSendMgmntClient()
		err := mgmntClient.ChangeStatusWorkflow(workspace, slug, true)
		if err != nil {
			log.Error(err.Error())
			return
		}

		fmt.Printf("Enabled workflow: %s\n", slug)
	},
}

func init() {
	WorkflowCmd.AddCommand(worklowEnableCmd)
}
