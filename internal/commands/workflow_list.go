/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

// listCmd represents the list command
var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflows for a workspace",
	Long:  `List workflows for a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		mgmnt_client := utils.GetSuprSendMgmntClient()

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		mode, _ := cmd.Flags().GetString("mode")
		workflows, err := mgmnt_client.GetWorkflows(workspace, limit, offset, mode)
		if err != nil {
			log.Errorf("Error getting workflows: %s", err)
			return
		}
		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(workflows, outputType)
	},
}

func init() {
	workflowListCmd.Flags().IntP("limit", "l", 20, "Limit the number of workflows to list")
	workflowListCmd.Flags().IntP("offset", "f", 0, "Offset the number of workflows to list (default: 0)")
	// add flag to set mode which can be one of draft, live with validation of the flag
	workflowListCmd.Flags().StringP("mode", "m", "live", "Mode to list workflows (draft, live)")
	workflowListCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Parent().HelpFunc()(cmd, args)
	})
	workflowCmd.AddCommand(workflowListCmd)
}
