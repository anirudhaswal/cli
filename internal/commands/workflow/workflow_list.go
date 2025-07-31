/*
Copyright © 2025 SuprSend
*/
package workflow

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

// listCmd represents the list command
var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflows for a workspace",
	Long:  `List workflows for a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}
		workspace, _ := cmd.Flags().GetString("workspace")
		mgmnt_client := utils.GetSuprSendMgmntClient()

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		mode, _ := cmd.Flags().GetString("mode")
		workflows, err := mgmnt_client.ListWorkflows(workspace, limit, offset, mode)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch workflows")
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Showing %d workflows out of %d from workspace %s \n", len(workflows.Results), workflows.Meta.Count, workspace))
		}
		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(workflows.Results, outputType)
	},
}

func init() {
	workflowListCmd.PersistentFlags().IntP("limit", "l", 20, "Limit the number of workflows to list")
	workflowListCmd.PersistentFlags().IntP("offset", "f", 0, "Offset the number of workflows to list (default: 0)")

	workflowListCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Parent().HelpFunc()(cmd, args)
	})
	WorkflowCmd.AddCommand(workflowListCmd)
}
