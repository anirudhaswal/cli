package category

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var categoryCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit categories",
	Long:  "Commit categories to a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		commitMsg, _ := cmd.Flags().GetString("commit-message")
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
		err := mgmnt_client.FinalizeCategories(workspace, commitMsg)
		if err != nil {
			log.WithError(err).Error("Couldn't commit categories")
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Committed categories to %s", workspace))
		}
	},
}

func init() {
	categoryCommitCmd.PersistentFlags().String("workspace", "staging", "Workspace to commit categories to")
	categoryCommitCmd.PersistentFlags().String("commit-message", "", "Commit message")
	CategoryCmd.AddCommand(categoryCommitCmd)
}
