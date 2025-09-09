package category

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var categoryPullCmd = &cobra.Command{
	Use:   "pull",
	Long:  "Pull categories from a workspace",
	Short: "Pull categories from a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		mode, _ := cmd.Flags().GetString("mode")

		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		categories, err := mgmnt_client.ListCategories(workspace, mode)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch categories")
			return
		}
		err = WriteToFile(categories, "categories_preferences.json")
		if err != nil {
			log.WithError(err).Error("Couldn't write categories to file")
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled categories from %s", workspace))
		}
	},
}

func init() {
	categoryPullCmd.PersistentFlags().String("mode", "live", "Mode to pull categories from")
	CategoryCmd.AddCommand(categoryPullCmd)
}
