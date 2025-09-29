package translation

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

var translationPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push workflows from local to suprsend",
	Long:  "push workflows from local to suprsend",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		outputDir, _ := cmd.Flags().GetString("dir")
		commit, _ := cmd.Flags().GetString("commit")
		commitMessage, _ := cmd.Flags().GetString("commit-message")

		if outputDir == "" {
			outputDir = filepath.Join(".", "suprsend", "translation")
		}

		files, err := os.ReadDir(outputDir)
		if err != nil {
			log.WithError(err).Errorf("Failed to read local translation directory")
			return
		}

		fmt.Printf("Pushing translations to %s\n", workspace)
		mgmntClient := utils.GetSuprSendMgmntClient()

		var translations []map[string]any

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			path := filepath.Join(outputDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				log.WithError(err).Errorf("Failed to read file %s", file.Name())
				return
			}

			var translation map[string]any
			if err := json.Unmarshal(data, &translation); err != nil {
				log.WithError(err).Errorf("Failed to parse JSON for %s", file.Name())
				return
			}

			translations = append(translations, translation)
		}

		if len(translations) > 0 {
			var p *pin.Pin
			var cancel context.CancelFunc

			if !utils.IsOutputPiped() {
				p = pin.New("Pushing translations...",
					pin.WithSpinnerColor(pin.ColorCyan),
					pin.WithTextColor(pin.ColorYellow),
				)
				cancel = p.Start(context.Background())
			}

			err = mgmntClient.PushTranslation(workspace, translations, commit, commitMessage)
			if p != nil && cancel != nil {
				if err != nil {
					p.Stop("")
					cancel()
					log.WithError(err).Errorf("Failed to push translations")
				} else {
					p.Stop(fmt.Sprintf("Pushed %d translations", len(translations)))
					cancel()
				}
			} else {
				if err != nil {
					log.WithError(err).Errorf("Failed to push translations")
				} else {
					fmt.Fprintf(os.Stdout, "Pushed %d translations\n", len(translations))
				}
			}
		} else {
			fmt.Printf("No translation files found in %s\n", outputDir)
		}
	},
}

func init() {
	translationPushCmd.Flags().StringP("commit", "c", "true", "Commit the translations (--commit=true)")
	translationPushCmd.Flags().StringP("commit-message", "m", "", "Commit message")
	translationPushCmd.Flags().StringP("dir", "d", "", "Directory for translations pull to (default: ./suprsend/translation)")
	TranslationCmd.AddCommand(translationPushCmd)
}
