package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/commands/category"
	"github.com/suprsend/cli/internal/commands/event"
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
			assetsToSync = []string{"workflow", "schema", "category", "event"}
		case "workflow":
			assetsToSync = []string{"workflow"}
		case "schema":
			assetsToSync = []string{"schema"}
		case "event":
			assetsToSync = []string{"event"}
		case "category":
			assetsToSync = []string{"category"}
		default:
			log.Errorf("Invalid asset type: '%s'. Valid options are: all, workflow, schema, event, category", assets)
			return
		}

		log.Infof("Syncing assets from %s to %s ...", fromWorkspace, toWorkspace)
		log.Infof("Assets to sync: %v", assetsToSync)

		mgmntClient := utils.GetSuprSendMgmntClient()
		hasErrors := false

		for _, assetType := range assetsToSync {
			switch assetType {
			case "workflow":
				err := syncWorkflows(mgmntClient, fromWorkspace, toWorkspace, mode)
				if err != nil {
					hasErrors = true
				}
			case "schema":
				err := syncSchemas(mgmntClient, fromWorkspace, toWorkspace, mode)
				if err != nil {
					hasErrors = true
				}
			case "event":
				err := syncEvents(mgmntClient, fromWorkspace, toWorkspace)
				if err != nil {
					hasErrors = true
				}
			case "category":
				err := syncCategories(mgmntClient, fromWorkspace, toWorkspace, mode)
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

func init() {
	rootCmd.AddCommand(syncCmd)

	// Flags consumed in Run
	syncCmd.Flags().StringP("from", "f", "staging", "Source workspace (required)")
	syncCmd.Flags().StringP("to", "t", "production", "Destination workspace (required)")
	syncCmd.Flags().StringP("mode", "m", "live", "Mode to sync assets (draft, live), default: live")
	syncCmd.Flags().StringP("assets", "a", "all", "Assets to sync (all, workflow, schema, event, category)")
}

func syncWorkflows(mgmntClient *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace, mode string) error {
	dirPath := filepath.Join(".", "suprsend", "workflow")

	workflows_resp, err := mgmntClient.GetWorkflows(fromWorkspace, mode)
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

	var errors []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".json")
		path := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("error reading file %s: %v", file.Name(), err))
			continue
		}

		var wf map[string]any
		if err := json.Unmarshal(data, &wf); err != nil {
			errors = append(errors, fmt.Sprintf("error unmarshalling JSON for %s: %v", file.Name(), err))
			continue
		}

		err = mgmntClient.PushWorkflow(toWorkspace, slug, wf, false, "")
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to push workflow %s: %v", slug, err))
			log.WithError(err).Errorf("Failed to push workflow %s", slug)
			continue
		}

		log.Infof("Pushed workflow: %s\n", slug)
	}
	if len(errors) > 0 {
		return fmt.Errorf("one or more workflows failed to sync:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

func syncSchemas(mgmntClient *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace, mode string) error {
	dirPath := filepath.Join(".", "suprsend", "schema")

	log.Infof("Pulling schemas from %s ...", fromWorkspace)
	schemas_resp, err := mgmntClient.GetSchemas(fromWorkspace, mode)
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

	var errors []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".json")
		path := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("error reading file %s: %v", file.Name(), err))
			continue
		}

		var sch map[string]any
		if err := json.Unmarshal(data, &sch); err != nil {
			errors = append(errors, fmt.Sprintf("error unmarshalling JSON for %s: %v", file.Name(), err))
			continue
		}

		err = mgmntClient.PushSchema(toWorkspace, slug, sch)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to push schema %s: %v", slug, err))
			log.WithError(err).Errorf("Failed to push schema %s", slug)
			continue
		}

		log.Infof("Pushed schema: %s\n", slug)
	}
	if len(errors) > 0 {
		return fmt.Errorf("one or more schemas failed to sync:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

func syncEvents(mgmntClient *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace string) error {
	dirPath := filepath.Join(".", "suprsend", "event")

	log.Infof("Pulling events from %s ...", fromWorkspace)
	events_resp, err := mgmntClient.GetEvents(fromWorkspace)
	if err != nil {
		return fmt.Errorf("error getting events: %w", err)
	}

	_, err = event.WriteEventsToFiles(events_resp, dirPath)
	if err != nil {
		return fmt.Errorf("error writing events to files: %w", err)
	}

	fmt.Printf("Pushing events to %s ...", toWorkspace)
	err = mgmntClient.PushEvents(toWorkspace)
	if err != nil {
		return fmt.Errorf("error pushing events: %w", err)
	}
	return nil
}

func syncCategories(mgmntClient *mgmnt.SS_MgmntClient, fromWorkspace, toWorkspace, mode string) error {
	dirPath := filepath.Join(".", "suprsend", "category")
	categoriesResp, err := mgmntClient.ListCategories(fromWorkspace, mode)
	if err != nil {
		return fmt.Errorf("error listing categories: %w", err)
	}
	log.Infof("Pulling categories from %s ...", fromWorkspace)
	err = category.WriteToFile(categoriesResp, "categories_preferences.json")
	if err != nil {
		return fmt.Errorf("error writing categories to files: %w", err)
	}
	filePath := filepath.Join(dirPath, "categories_preferences.json")
	categories, err := category.ReadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading categories from file: %w", err)
	}
	err = mgmntClient.PushCategories(toWorkspace, categories)
	if err != nil {
		return fmt.Errorf("error pushing categories: %w", err)
	}
	log.Printf("Pushed categories to %s", toWorkspace)
	return nil
}
