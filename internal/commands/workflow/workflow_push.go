package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var force bool

var workflowPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push workflows from local to suprsend",
	Long:  `push workflows from local to suprsend dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")
		workspace, _ := cmd.Flags().GetString("workspace")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			log.WithError(err).Errorf("Directory '%s' does not exist. Exiting.\n", dirPath)
			return
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local workflows directory")
			return
		}

		fmt.Printf("Pushing workflows to %s\n", workspace)
		mgmntClient := utils.GetSuprSendMgmntClient()

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			slug := strings.TrimSuffix(file.Name(), ".json")
			path := filepath.Join(dirPath, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				log.WithError(err).Errorf("Failed to read file %s", file.Name())
				continue
			}

			var workflow map[string]any
			if err := json.Unmarshal(data, &workflow); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
				continue
			}

			err = mgmntClient.PushWorkflow(workspace, slug, workflow)
			if err != nil {
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				continue
			}
		}
	},
}

func init() {
	workflowPushCmd.Flags().BoolVarP(&force, "force-dir", "d", false, "Create workflow directory without the permission")
	WorkflowCmd.AddCommand(workflowPushCmd)
}
