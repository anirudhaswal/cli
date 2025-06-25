/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		workflows, err := mgmnt_client.GetWorkflows(workspace, limit, offset, mode)
		if err != nil {
			log.Errorf("Error getting workflows: %s", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Showing %d workflows out of %d from workspace %s \n", len(workflows.Results), workflows.Meta.Count, workspace))
		}
		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(workflows.Results, outputType)
	},
}

var force bool

var workflowPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull workflows from workspace to local directory",
	Long:  `pull workflows from workspace to local directory of each workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if force {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					fmt.Println("Failed to create directory:", err)
					os.Exit(1)
				}
				fmt.Println("Directory created at:", dirPath)
			} else {
				fmt.Printf("Directory '%s' does not exist. Do you want to create it? (y/n): ", dirPath)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(input))

				if input == "y" || input == "yes" {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						fmt.Println("Failed to create directory:", err)
						os.Exit(1)
					}
					fmt.Println("Directory created at:", dirPath)
				} else {
					fmt.Println("Directory not created. Exiting.")
					return
				}
			}
		} else {
			fmt.Println("Directory already exists:", dirPath)
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		workflows_resp, err := mgmnt_client.GetWorkflows("production", 20, 0, "live")
		if err != nil {
			log.Errorf("Error getting workflows: %s", err)
			return
		}

		fmt.Println("Pulling workflows...")
		if err := utils.WriteWorkflowsToFiles(*workflows_resp, "./suprsend/workflow"); err != nil {
			fmt.Println("Error saving workflows:", err)
		}
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
	workflowPullCmd.Flags().BoolVarP(&force, "force-dir", "f", false, "Create workflow directory without the permission")
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowPullCmd)
}
