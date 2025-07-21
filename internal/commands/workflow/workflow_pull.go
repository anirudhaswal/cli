package workflow

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var workflowPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull workflows from workspace to local directory",
	Long:  `pull workflows from workspace to local directory of each workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")

		workspace, _ := cmd.Flags().GetString("workspace")
		mode, _ := cmd.Flags().GetString("mode")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if force {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					log.WithError(err).Error("Failed to create directory")
					return
				}
				log.Printf("Directory created at: %s\n", dirPath)
			} else {
				log.Printf("Directory '%s' does not exist. Do you want to create it? (y/n): ", dirPath)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(input))

				if input == "y" || input == "yes" {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						log.WithError(err).Error("Failed to create directory")
						return
					}
					log.Infof("Directory created at: %s", dirPath)
				} else {
					log.Error("Directory not created. Exiting.")
					return
				}
			}
		} else {
			log.Infof("Directory already exists: %s\n", dirPath)
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		workflows_resp, err := mgmnt_client.GetWorkflows(workspace, mode)

		if err != nil {
			log.WithError(err).Error("Error getting workflows")
			return
		}

		log.Infoln("Pulling workflows...")
		if err := WriteWorkflowsToFiles(*workflows_resp, dirPath); err != nil {
			log.WithError(err).Error("Error saving workflows")
			return
		}
	},
}

func init() {
	WorkflowCmd.AddCommand(workflowPullCmd)
}
