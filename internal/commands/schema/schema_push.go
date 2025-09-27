package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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
		slug, _ := cmd.Flags().GetString("slug")
		commit, _ := cmd.Flags().GetString("commit")
		commitMessage, _ := cmd.Flags().GetString("commit-message")
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			path = filepath.Join(".", "suprsend", "schema")
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Errorf("Directory %s does not exist", path)
			return
		}

		if err := validateInputDirectory(path); err != nil {
			log.Errorf("Error with input directory: %v\n", err)
			return
		}

		files, err := os.ReadDir(path)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local schema directory")
			return
		}

		mgmntClient := utils.GetSuprSendMgmntClient()
		hasError := false
		var p *pin.Pin
		var cancel context.CancelFunc

		if slug != "" {
			fileName := fmt.Sprintf("%s.json", slug)
			filePath := filepath.Join(path, fileName)

			if _, err := os.Stat(filePath); err != nil {
				log.WithError(err).Errorf("Failed to find schema file %s", filePath)
				return
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.WithError(err).Errorf("Failed to read schema file %s", filePath)
				return
			}

			if !hasError && !utils.IsOutputPiped() {
				p = pin.New(fmt.Sprintf("Pushing %s...", slug),
					pin.WithSpinnerColor(pin.ColorCyan),
					pin.WithTextColor(pin.ColorYellow),
				)
				cancel = p.Start(context.Background())
			}

			var schema map[string]any
			if err := json.Unmarshal(data, &schema); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", filePath)
				return
			}

			urlEncodedCommitMessage := url.QueryEscape(commitMessage)
			err = mgmntClient.PushSchema(workspace, slug, schema, commit, urlEncodedCommitMessage)
			if err != nil {
				log.WithError(err).Errorf("Failed to push schema %s", slug)
				return
			}

			if p != nil && cancel != nil {
				p.Stop(fmt.Sprintf("Pushed schema: %s", slug))
				cancel()
				p = nil
				cancel = nil
			} else {
				fmt.Fprintf(os.Stdout, "Pushed schema: %s\n", slug)
			}
			return
		}

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
			path := filepath.Join(path, file.Name())
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

			err = mgmntClient.PushSchema(workspace, slug, schema, commit, commitMessage)
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

			if p != nil && cancel != nil {
				p.Stop(fmt.Sprintf("Pushed schema: %s", slug))
				cancel()
				p = nil
				cancel = nil
			} else {
				fmt.Fprintf(os.Stdout, "Pushed schema: %s\n", slug)
			}
			hasError = false
		}
	},
}

func init() {
	schemaPushCmd.Flags().StringP("dir", "p", "", "Directory for schemas pull to (default: ./suprsend/schema)")
	schemaPushCmd.Flags().StringP("commit", "c", "true", "Commit the schemas (--commit=true)")
	schemaPushCmd.Flags().StringP("commit-message", "m", "", "Commit message describing the changes for --commit=true")
	schemaPushCmd.PersistentFlags().StringP("slug", "g", "", "Slug of schema to push")
	SchemaCmd.AddCommand(schemaPushCmd)
}
