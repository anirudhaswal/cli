package schema

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit schema from draft to live",
	Long:  `Commit schema from draft to live in a workspace`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("schema slug argument is required for schemas.")
		}
		slug := args[0]

		workspace, _ := cmd.Flags().GetString("workspace")
		mgmnt_client := utils.GetSuprSendMgmntClient()
		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		err := mgmnt_client.FinalizeSchema(workspace, slug, true)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Committed schema %s", slug))
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaCommitCmd)
}
