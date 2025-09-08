package category

import (
	"context"
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var categoryPushCmd = &cobra.Command{
	Use:   "push",
	Long:  "Push categories to a workspace",
	Short: "Push categories to a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		dirPath := path.Join(".", "suprsend", "category")
		categories, err := ReadFromFile(path.Join(dirPath, "categories_preferences.json"))
		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		if err != nil {
			log.WithError(err).Error("Couldn't read categories from file")
			return
		}
		mgmnt_client := utils.GetSuprSendMgmntClient()
		err = mgmnt_client.PushCategories(workspace, categories)
		if err != nil {
			log.WithError(err).Error("Couldn't push categories")
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pushed categories to %s", workspace))
		}
	},
}

func init() {
	categoryPushCmd.PersistentFlags().String("workspace", "staging", "Workspace to push categories to")
	CategoryCmd.AddCommand(categoryPushCmd)
}
