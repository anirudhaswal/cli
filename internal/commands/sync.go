package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/commands/schema"
	"github.com/suprsend/cli/internal/commands/workflow"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/cli/mgmnt"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all your SuprSend assets locally",
	Long:  `Sync all your SuprSend assets locally with the server`,
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		fromWorkspace, _ := cmd.Flags().GetString("from")
		toWorkspace, _ := cmd.Flags().GetString("to")
		assets, _ := cmd.Flags().GetString("assets")
		if fromWorkspace == toWorkspace {
			log.Error("Cannot sync within the same workspace. Source and destination workspaces must be different.")
			return
		}

		var assetsToSync []string
		switch assets {
		case "all":
			assetsToSync = []string{"workflow", "schema"}
		case "workflow":
			assetsToSync = []string{"workflow"}
		case "schema":
			assetsToSync = []string{"schema"}
		default:
			log.Errorf("Invalid asset type: '%s'. Valid options are: all, workflow, schema", assets)
			return
		}

		log.Infof("Syncing assets from %s to %s ... \n", fromWorkspace, toWorkspace)
		log.Infof("Assets to sync: %v\n", assetsToSync)

		mgmnt_client := utils.GetSuprSendMgmntClient()
		hasErrors := false

		for _, assetType := range assetsToSync {
			switch assetType {
			case "workflow":
				err := syncWorkflows(fromWorkspace, toWorkspace, mode)
				if err != nil {
					hasErrors = true
				}
			case "schema":
				err := syncSchemas(mgmnt_client, fromWorkspace, toWorkspace, mode)
				if err != nil {
					hasErrors = true
				}
			default:
				log.Errorf("Invalid asset type: %s", assetType)
			}
		}
		if hasErrors {
			log.Error("Sync complete with errors")
		} else {
			log.Info("Sync complete")
		}
	},
}

func syncWorkflows(fromWorkspace, toWorkspace, mode string) error {
	dirPath := filepath.Join(".", "suprsend", "workflow")

	mgmnt_client := utils.GetSuprSendMgmntClient()
	workflows_resp, err := mgmnt_client.GetWorkflows(fromWorkspace, mode)
	if err != nil {
		return fmt.Errorf("error getting workflows: %w", err)
	}

	log.Infof("Pulling workflows from %s ... \n", fromWorkspace)
	_, err = workflow.WriteWorkflowsToFiles(*workflows_resp, dirPath)
	if err != nil {
		return fmt.Errorf("error writing workflows to files: %w", err)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading local workflows directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".json")
		path := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", file.Name(), err)
		}

		var workflow map[string]any
		if err := json.Unmarshal(data, &workflow); err != nil {
			return fmt.Errorf("error unmarshalling JSON for %s: %w", file.Name(), err)
		}

		err = mgmnt_client.PushWorkflow(toWorkspace, slug, workflow, false, "")
		if err != nil {
			log.WithError(err).Errorf("Failed to push workflow %s", slug)
			continue
		}

		log.Infof("Pushed workflow: %s\n", slug)
	}
	return nil
}

func syncSchemas(mgmnt_client *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace, mode string) error {
	dirPath := filepath.Join(".", "suprsend", "schema")

	fmt.Printf("Pulling schemas from %s ... \n", fromWorkspace)
	schemas_resp, err := mgmnt_client.GetSchemas(fromWorkspace, mode)
	if err != nil {
		return fmt.Errorf("error getting schemas: %w", err)
	}

	_, err = schema.WriteSchemasToFiles(schemas_resp, dirPath)
	if err != nil {
		return fmt.Errorf("error writing schemas to files: %w", err)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading local schemas directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".json")
		path := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", file.Name(), err)
		}

		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			return fmt.Errorf("error unmarshalling JSON for %s: %w", file.Name(), err)
		}

		err = mgmnt_client.PushSchema(toWorkspace, slug, schema)
		if err != nil {
			log.WithError(err).Errorf("Failed to push schema %s", slug)
			continue
		}

		log.Infof("Pushed schema: %s\n", slug)
	}
	return nil
}
