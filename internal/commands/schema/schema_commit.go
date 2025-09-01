package schema

import (
	"fmt"
	"os"

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
			log.Error("Schema slug argument is required. Example: suprsend schema commit <slug>")
			return
		}
		slug := args[0]

		workspace, _ := cmd.Flags().GetString("workspace")
		mgmnt_client := utils.GetSuprSendMgmntClient()
		var p *pin.Pin

		if !utils.IsOutputPiped() {
			p = pin.New("Committing schema...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
		}

		err := mgmnt_client.FinalizeSchema(workspace, slug, true)
		if err != nil {
			log.Error(err.Error())
			return
		}

		if p != nil {
			p.Stop(fmt.Sprintf("Successfully committed schema '%s' to live mode", slug))
		} else {
			fmt.Fprintf(os.Stdout, "Successfully committed schema '%s' to live mode\n", slug)
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaCommitCmd)
}
