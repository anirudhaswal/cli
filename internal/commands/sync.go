package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/commands/event"
	"github.com/suprsend/cli/internal/commands/schema"
	"github.com/suprsend/cli/internal/commands/workflow"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/cli/mgmnt"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all your SuprSend assests locally",
	Long:  `Sync all your SuprSend assets locally with the server`,
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		fromWorkspace, _ := cmd.Flags().GetString("from")
		toWorkspace, _ := cmd.Flags().GetString("to")
		assets, _ := cmd.Flags().GetString("assets")

		var assetsToSync []string
		switch assets {
		case "all":
			assetsToSync = []string{"workflows", "schemas"}
		case "workflows":
			assetsToSync = []string{"workflows"}
		case "schemas":
			assetsToSync = []string{"schemas"}
		case "assets":
			assetsToSync = []string{"assets"}
		}

		fmt.Printf("Syncing assets from %s to %s ... \n", fromWorkspace, toWorkspace)
		fmt.Printf("Assets to sync: %v\n", assetsToSync)

		mgmnt_client := utils.GetSuprSendMgmntClient()

		for _, assetType := range assetsToSync {
			switch assetType {
			case "workflows":
				syncWorkflows(fromWorkspace, toWorkspace, mode)
			case "schemas":
				syncSchemas(mgmnt_client, fromWorkspace, toWorkspace)
			case "assets":
				syncEvents(mgmnt_client, fromWorkspace, toWorkspace)
			default:
				log.Errorf("Invalid asset type: %s", assetType)
			}
		}

		fmt.Println("Sync complete")
	},
}

func syncWorkflows(fromWorkspace, toWorkspace, mode string) {
	dirPath := filepath.Join(".", "suprsend", "workflow")

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
}

func syncSchemas(mgmnt_client *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace string) {
	dirPath := filepath.Join(".", "suprsend", "schema")

	fmt.Printf("Pulling schemas from %s ... \n", fromWorkspace)
	schemas_resp, err := mgmnt_client.GetSchemas(fromWorkspace)
	if err != nil {
		log.WithError(err).Error("Error getting schemas")
		return
	}

	_, err = schema.WriteSchemasToFiles(schemas_resp, dirPath)
	if err != nil {
		log.WithError(err).Error("Error saving schemas")
		return
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.WithError(err).Errorf("Failed to read local schemas directory")
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

		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
			return
		}

		err = mgmnt_client.PushSchema(toWorkspace, slug, schema)
		if err != nil {
			log.WithError(err).Errorf("Failed to push schema %s", slug)
			continue
		}
	}
}

func syncEvents(mgmnt_client *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace string) {
	dirPath := filepath.Join(".", "suprsend", "event")

	fmt.Printf("Pulling events from %s ... \n", fromWorkspace)
	events_resp, err := mgmnt_client.GetEvents(fromWorkspace)
	if err != nil {
		log.WithError(err).Error("Error getting events")
		return
	}

	_, err = event.WriteEventsToFiles(events_resp, dirPath)
	if err != nil {
		log.WithError(err).Error("Error saving events")
		return
	}

	fmt.Printf("Pushing events to %s ... \n", toWorkspace)
	err = mgmnt_client.PushEvents(toWorkspace)
	if err != nil {
		log.WithError(err).Error("Failed to push events")
		return
	}
}
