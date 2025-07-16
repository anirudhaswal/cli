package workflow

import (
	"encoding/json"
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

		mgmntClient := utils.GetSuprSendMgmntClient()
		resp, err := mgmntClient.GetWorkflows("production", 20, 0, "live")
		if err != nil {
			log.WithError(err).Error("Failed to get workflows")
			return
		}

		existingSlugs := make(map[string]bool)
		for _, wf := range resp.Results {
			existingSlugs[wf.Slug] = true
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local workflows directory")
			return
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			slug := strings.TrimSuffix(file.Name(), ".json")
			if _, exists := existingSlugs[slug]; exists {
				log.Printf("Skipping '%s.json' (already exists on server)\n", slug)
				continue
			}

			path := filepath.Join(dirPath, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				log.WithError(err).Errorf("Failed to read file %s", file.Name())
				return
			}

			var workflow map[string]any
			if err := json.Unmarshal(data, &workflow); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
				return
			}

			err = mgmntClient.PushWorkflow(workspace, slug, workflow)
			if err != nil {
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				return
			}

			log.Printf("Pushed workflow: %s\n", slug)
		}
	},
}

func init() {
	workflowPullCmd.Flags().BoolVarP(&force, "force-dir", "f", false, "Create workflow directory without the permission")
	WorkflowCmd.AddCommand(workflowPushCmd)
}
