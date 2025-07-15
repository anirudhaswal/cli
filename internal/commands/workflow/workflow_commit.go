package workflow

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var workflowCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit a draft workflow to live.",
	Long:  "Commits a draft workflow to live",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("Category slug is required.")
		}
		workspace, _ := cmd.Flags().GetString("workspace")

		slug := args[0]

		mgmntClient := utils.GetSuprSendMgmntClient()

		commmitMessage, _ := cmd.Flags().GetString("commit-message")

		err := mgmntClient.CommitWorkflow(workspace, slug, commmitMessage)
		if err != nil {
			log.WithError(err).Errorf("Failed to commit workflow %s", slug)
		}

		log.Printf("Committed workflow: %s\n", slug)
	},
}

func init() {
	workflowListCmd.Flags().StringP("commit-message", "c", "", "Commit Message for making a workflow live")
	WorkflowCmd.AddCommand(workflowCommitCmd)
}
