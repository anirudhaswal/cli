package schema

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/cli/mgmnt"
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
			p.Stop(fmt.Sprintf("Listed %d schemas from %s", len(schemas.Results), workspace))
		}

		outputType, _ := cmd.Flags().GetString("output")
		if len(schemas.Results) == 0 && utils.IsOutputPiped() {
			emptyResponse := mgmnt.ListSchemaResponse{
				Results: []mgmnt.SchemaResponse{},
				Meta: struct {
					Count  int `json:"count"`
					Limit  int `json:"limit"`
					Offset int `json:"offset"`
				}{
					Count:  0,
					Limit:  20,
					Offset: 0,
				},
			}
			utils.OutputData(emptyResponse, outputType)
			return
		}
		utils.OutputData(schemas.Results, outputType)
	},
}

func init() {
	SchemaCmd.Flags().StringP("workspace", "w", "staging", "Workspace to pull the schemas from")
	SchemaCmd.AddCommand(schemaListCmd)
}
