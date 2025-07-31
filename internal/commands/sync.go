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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flag := cmd.InheritedFlags().Lookup("workspace"); flag != nil {
			flag.Hidden = true
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")

		mode, _ := cmd.Flags().GetString("mode")
		fromWorkspace, _ := cmd.Flags().GetString("from")
		toWorkspace, _ := cmd.Flags().GetString("to")

		log.Infof("Syncing workflows from %s to %s", fromWorkspace, toWorkspace)

		mgmnt_client := utils.GetSuprSendMgmntClient()
		workflows_resp, err := mgmnt_client.GetWorkflows(fromWorkspace, mode)

		if err != nil {
			log.WithError(err).Error("Error getting workflows")
			return
		}

		log.Infof("Pulling workflows from %s ... \n", fromWorkspace)
		_, err = workflow.WriteWorkflowsToFiles(*workflows_resp, dirPath)
		if err != nil {
			log.WithError(err).Error("Error saving workflows")
			return
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

			err = mgmnt_client.PushWorkflow(toWorkspace, slug, workflow)
			if err != nil {
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				continue
			}

			log.Printf("Pushed workflow: %s\n", slug)
		}

		log.Println("Sync Completed")
	},
}
