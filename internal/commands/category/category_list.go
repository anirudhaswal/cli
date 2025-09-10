package category

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

type CategoryTableRow struct {
	RootCategory             string `json:"root_category"`
	Section                  string `json:"section"`
	CategoryName             string `json:"category_name"`
	DefaultPreference        string `json:"default_preference"`
	DefaultMandatoryChannels string `json:"default_mandatory_channels"`
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
		outputType, _ := cmd.Flags().GetString("output")

		// Create flattened table rows
		var tableRows []CategoryTableRow
		for _, category := range categories.Categories {
			for _, section := range category.Sections {
				for _, subcategory := range section.Subcategories {
					tableRows = append(tableRows, CategoryTableRow{
						RootCategory:             category.RootCategory,
						Section:                  section.Name,
						CategoryName:             subcategory.Name,
						DefaultPreference:        subcategory.DefaultPreference,
						DefaultMandatoryChannels: strings.Join(subcategory.DefaultMandatoryChannels, ", "),
					})
				}
			}
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Listed %d categories from %s", len(tableRows), workspace))
		}

		if len(tableRows) == 0 && utils.IsOutputPiped() {
			utils.OutputData([]interface{}{}, outputType)
			return
		}

		utils.OutputData(tableRows, outputType)
	},
}

func init() {
	categoryListCmd.PersistentFlags().StringP("mode", "m", "live", "Mode of preferences to list")
	CategoryCmd.AddCommand(categoryListCmd)
}
