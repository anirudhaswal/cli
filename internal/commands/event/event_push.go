package event

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var eventPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push linked events",
	Long:  "Push linked events in schemas",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")

		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Pushing events...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		err := mgmnt_client.PushEvents(workspace)
		if err != nil {
			if p != nil {
				p.Stop("")
			}
			log.WithError(err).Error("Failed to push events")
			return
		}
		if p != nil {
			p.Stop("Successfully pushed events")
		} else {
			fmt.Fprintf(os.Stdout, "Successfully pushed events\n")
		}
	},
}

func init() {
	eventPushCmd.Flags().StringP("workspace", "w", "staging", "Workspace to push events to")
	EventCmd.AddCommand(eventPushCmd)
}
