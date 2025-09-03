package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push schemas",
	Long:  "Push schemas in a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")

		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir == "" {
			outputDir = promptForOutputDirectory()
		}

		if err := ensureOutputDirectory(outputDir); err != nil {
			fmt.Fprintf(os.Stdout, "Error with output directory: %v\n", err)
			return
		}

		files, err := os.ReadDir(outputDir)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local schema directory")
			return
		}

		fmt.Printf("Pushing schemas to %s\n", workspace)
		mgmntClient := utils.GetSuprSendMgmntClient()
		hasError := false
		var p *pin.Pin
		var cancel context.CancelFunc

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			slug := strings.TrimSuffix(file.Name(), ".json")
			if !hasError && !utils.IsOutputPiped() {
				p = pin.New(fmt.Sprintf("Pushing %s...", slug),
					pin.WithSpinnerColor(pin.ColorCyan),
					pin.WithTextColor(pin.ColorYellow),
				)
				cancel = p.Start(context.Background())
			}
			path := filepath.Join(outputDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				if p != nil && cancel != nil {
					p.Stop("")
					cancel()
					p = nil
					cancel = nil
				}
				hasError = true
				log.WithError(err).Errorf("Failed to read file %s", file.Name())
				continue
			}

			var schema map[string]any
			if err := json.Unmarshal(data, &schema); err != nil {
				if p != nil && cancel != nil {
					p.Stop("")
					cancel()
					p = nil
					cancel = nil
				}
				hasError = true
				log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
				continue
			}

			err = mgmntClient.PushSchema(workspace, slug, schema)
			if err != nil {
				if p != nil && cancel != nil {
					p.Stop("")
					cancel()
					p = nil
					cancel = nil
				}
				hasError = true
				log.WithError(err).Errorf("Failed to push schema %s", slug)
				continue
			}

			if !hasError && p != nil && cancel != nil {
				p.Stop(fmt.Sprintf("Pushed workflow: %s", slug))
				cancel()
				p = nil
				cancel = nil
			} else {
				fmt.Fprintf(os.Stdout, "Pushed workflow: %s\n", slug)
			}
			hasError = false
		}
	},
}

func init() {
	schemaPushCmd.PersistentFlags().StringP("output-dir", "d", "", "Output directory for schemas")
	SchemaCmd.AddCommand(schemaPushCmd)
}
