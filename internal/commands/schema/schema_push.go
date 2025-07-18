package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var schemaPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push schemas",
	Long:  "Push schemas in a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("workspace argument is required for schemas.")
		}
		workspace := args[0]

		dirPath := filepath.Join(".", "suprsend", "schema")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			log.WithError(err).Errorf("Directory '%s' does not exist. Exiting.\n", dirPath)
			return
		}

		mgmntClient := utils.GetSuprSendMgmntClient()
		resp, err := mgmntClient.GetSchemas(workspace)
		if err != nil {
			log.WithError(err).Error("Failed to get schemas")
			return
		}

		existingSlugs := make(map[string]bool)
		for _, schema := range resp.Results {
			obj, ok := schema.(map[string]any)
			if !ok {
				continue
			}

			slug, _ := obj["slug"].(string)
			existingSlugs[slug] = true
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local schema directory")
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

			var schema map[string]any
			if err := json.Unmarshal(data, &schema); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
				return
			}

			err = mgmntClient.PushSchema(workspace, slug, schema)
			if err != nil {
				log.WithError(err).Errorf("Failed to push schema %s", slug)
				return
			}

			log.Printf("Pushed schema: %s\n", slug)
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaPushCmd)
}
