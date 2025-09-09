package category

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

type CategoryDisplayItem struct {
	SerialNo     int    `json:"serial_no"`
	RootCategory string `json:"root_category"`
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List categories",
	Long:  "List preferences categories in a workspace",
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
		if p != nil {
			p.Stop(fmt.Sprintf("Listed %d categories from %s", len(categories.Categories), workspace))
		}

		outputType, _ := cmd.Flags().GetString("output")
		displayItems := make([]CategoryDisplayItem, len(categories.Categories))
		for i, category := range categories.Categories {
			displayItems[i] = CategoryDisplayItem{
				SerialNo:     i + 1,
				RootCategory: category.RootCategory,
			}
		}

		if len(displayItems) == 0 && utils.IsOutputPiped() {
			utils.OutputData([]interface{}{}, outputType)
			return
		}

		utils.OutputData(displayItems, outputType)
	},
}

func init() {
	categoryListCmd.PersistentFlags().StringP("mode", "m", "live", "Mode of preferences to list")
	CategoryCmd.AddCommand(categoryListCmd)
}
