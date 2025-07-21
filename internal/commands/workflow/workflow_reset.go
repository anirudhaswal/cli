package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var workflowResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset a live workflow to draft",
	Long:  "Commits a live workflow to draft",
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
			log.WithError(err).Errorf("Failed to commit workflow %s", slug)
			return
		}

		log.Printf("Reset workflow: %s\n", slug)
	},
}

func init() {
	WorkflowCmd.AddCommand(workflowResetCmd)
}
