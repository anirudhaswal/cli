package category

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir == "" {
			outputDir = filepath.Join(".", "suprsend", "category")
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				outputDir = promptForOutputDirectory()
			}
			if outputDir == "" {
				fmt.Fprintf(os.Stdout, "No output directory specified. Exiting.\n")
				return
			}
		}
		if err := ensureOutputDirectory(outputDir); err != nil {
			fmt.Fprintf(os.Stdout, "Error with output directory: %v\n", err)
			return
		}
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
		err = WriteToFileWithPath(categories, outputDir, "categories_preferences.json")
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
	categoryPullCmd.PersistentFlags().StringP("mode", "m", "live", "Mode to pull categories from")
	categoryPullCmd.PersistentFlags().StringP("output-dir", "d", "", "Output directory for categories")
	CategoryCmd.AddCommand(categoryPullCmd)
}
