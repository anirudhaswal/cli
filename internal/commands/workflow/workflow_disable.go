package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var workflowDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a workflow",
	Long:  "Disable a workflow to deactivate.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("Category slug is required.")
			return
		}
		workspace, _ := cmd.Flags().GetString("workspace")

		slug := args[0]

		mgmntClient := utils.GetSuprSendMgmntClient()

		err := mgmntClient.ChangeStatusWorkflow(workspace, slug, false)
		if err != nil {
			log.WithError(err).Errorf("Failed to disable workflow %s", slug)
			return
		}

		log.Printf("Disabled workflow: %s\n", slug)
	},
}

func init() {
	WorkflowCmd.AddCommand(workflowDisableCmd)
}
