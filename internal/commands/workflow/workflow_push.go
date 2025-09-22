package workflow

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

var workflowPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push workflows from local to suprsend",
	Long:  `push workflows from local to suprsend dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		path, _ := cmd.Flags().GetString("path")
		commit, _ := cmd.Flags().GetBool("commit")
		commitMessage, _ := cmd.Flags().GetString("commit-message")
		slug, _ := cmd.Flags().GetString("slug")

		if path == "" {
			path = filepath.Join(".", "suprsend", "workflow")
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
			log.WithError(err).Errorf("Failed to read local workflows directory")
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
				log.WithError(err).Errorf("Failed to find workflow file %s", filePath)
				return
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.WithError(err).Errorf("Failed to read workflow file %s", filePath)
			}

			if !hasError && !utils.IsOutputPiped() {
				p = pin.New(fmt.Sprintf("Pushing %s...", slug),
					pin.WithSpinnerColor(pin.ColorCyan),
					pin.WithTextColor(pin.ColorYellow),
				)
				cancel = p.Start(context.Background())
			}
			var workflow map[string]any
			if err := json.Unmarshal(data, &workflow); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", filePath)
				return
			}

			err = mgmntClient.PushWorkflow(workspace, slug, workflow, commit, commitMessage)
			if err != nil {
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				return
			}

			if p != nil && cancel != nil {
				p.Stop(fmt.Sprintf("Pushed workflow: %s", slug))
				cancel()
				p = nil
				cancel = nil
			} else {
				fmt.Fprintf(os.Stdout, "Pushed workflow: %s\n", slug)
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
			filePath := filepath.Join(path, file.Name())
			data, err := os.ReadFile(filePath)
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

			var workflow map[string]any
			if err := json.Unmarshal(data, &workflow); err != nil {
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

			err = mgmntClient.PushWorkflow(workspace, slug, workflow, commit, commitMessage)
			if err != nil {
				if p != nil && cancel != nil {
					p.Stop("")
					cancel()
					p = nil
					cancel = nil
				}
				hasError = true
				log.WithError(err).Errorf("Failed to push workflow %s", slug)
				continue
			}

			if p != nil && cancel != nil {
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
	workflowPushCmd.PersistentFlags().StringP("path", "p", "", "Output directory for workflows")
	workflowPushCmd.PersistentFlags().BoolP("commit", "c", true, "Commit the workflows")
	workflowPushCmd.PersistentFlags().StringP("commit-message", "m", "", "Commit message for the workflows")
	workflowPushCmd.PersistentFlags().StringP("slug", "g", "", "Slug of the workflow to push")
	WorkflowCmd.AddCommand(workflowPushCmd)
}
