package schema

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

var schemaPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push schemas",
	Long:  "Push schemas in a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")

		dirPath := filepath.Join(".", "suprsend", "schema")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			log.WithError(err).Errorf("Directory '%s' does not exist. Exiting.\n", dirPath)
			return
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local schema directory")
			return
		}

		fmt.Printf("Pushing schemas to %s\n", workspace)
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

			fmt.Fprintf(os.Stdout, "Pushed schema: %s\n", slug)
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaPushCmd)
}
