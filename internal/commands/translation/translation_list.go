package translation

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var translationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Translations",
	Long:  "List Translations",
	Run: func(cmd *cobra.Command, args []string) {
		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}
		mode, _ := cmd.Flags().GetString("mode")
		workspace, _ := cmd.Flags().GetString("workspace")
		mgmntClient := utils.GetSuprSendMgmntClient()
		translations, err := mgmntClient.ListTranslations(workspace, mode)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch translations")
			return
		}
		msg := fmt.Sprintf("Listed %d translation files from %s in %s mode", len(translations.Results), workspace, mode)
		if p != nil {
			p.Stop(msg)
		}
		outputType, _ := cmd.Flags().GetString("output")
		if len(translations.Results) == 0 && utils.IsOutputPiped() {
			utils.OutputData([]interface{}{}, outputType)
			return
		}
		utils.OutputData(translations.Results, outputType)
	},
}

func init() {
	translationListCmd.Flags().StringP("mode", "m", "live", "Mode to list translations for")
	translationListCmd.Flags().StringP("output", "o", "pretty", "Output type (pretty, yaml, json)")
	TranslationCmd.PersistentFlags().StringP("workspace", "w", "staging", "Workspace to list translations for")
	TranslationCmd.PersistentFlags().StringP("service-token", "s", "", "Service token (default: $SUPRSEND_SERVICE_TOKEN)")
	TranslationCmd.AddCommand(translationListCmd)
}
