package schema

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset schema from live to draft",
	Long:  `Reset schema from live to draft in a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("schema slug argument is required for schemas. Example: suprsend schema reset <slug>")
			return
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

		err := mgmnt_client.FinalizeSchema(workspace, slug, false)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Reset schema %s", slug))
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaResetCmd)
}
