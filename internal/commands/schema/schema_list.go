package schema

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schemas",
	Long:  `List schemas in a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")

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

		schemas, err := mgmnt_client.ListSchema(workspace)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		if p != nil {
			p.Stop(fmt.Sprintf("Showing %d schemas out of %d from workspace %s \n", len(schemas.Results), schemas.Meta.Count, workspace))
		}

		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(schemas.Results, outputType)
	},
}

func init() {
	SchemaCmd.Flags().StringP("workspace", "w", "staging", "Workspace to pull the schemas from")
	SchemaCmd.AddCommand(schemaListCmd)
}
