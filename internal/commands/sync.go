package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/commands/workflow"
	"github.com/suprsend/cli/internal/utils"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all your SuprSend assests locally",
	Long:  `Sync all your SuprSend assets locally with the server`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")

		workspace, _ := cmd.Flags().GetString("workspace")
		mode, _ := cmd.Flags().GetString("mode")

		mgmnt_client := utils.GetSuprSendMgmntClient()
		workflows_resp, err := mgmnt_client.GetWorkflows(workspace, mode)

		if err != nil {
			log.WithError(err).Error("Error getting workflows")
			return
		}

		log.Infoln("Pulling workflows...")
		if err := workflow.WriteWorkflowsToFiles(*workflows_resp, dirPath); err != nil {
			log.WithError(err).Error("Error saving workflows")
			return
		}

		existingSlugs := make(map[string]bool)
		for _, wf := range workflows_resp.Results {
			obj, ok := wf.(map[string]any)
			if !ok {
				continue
			}

			slug, _ := obj["slug"].(string)
			existingSlugs[slug] = true
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

			err = mgmnt_client.PushWorkflow(workspace, slug, workflow)
			if err != nil {
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				return
			}

			log.Printf("Pushed workflow: %s\n", slug)
		}

		log.Println("Sync Completed")
	},
}
