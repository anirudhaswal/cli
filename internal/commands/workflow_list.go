/*
Copyright © 2025 SuprSend
*/
package commands

import (
	"bufio"
	"context"
	"encoding/json"
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

		workspace, _ := cmd.Flags().GetString("workspace")

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
		workflows_resp, err := mgmnt_client.GetWorkflows(workspace, 20, 0, "live")
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

var workflowPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push workflows from local to suprsend",
	Long:  `push workflows from local to suprsend dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")
		workspace, _ := cmd.Flags().GetString("workspace")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("Directory '%s' does not exist. Exiting.\n", dirPath)
			return
		}

		mgmntClient := utils.GetSuprSendMgmntClient()
		resp, err := mgmntClient.GetWorkflows("production", 20, 0, "live")
		if err != nil {
			log.Errorf("Failed to get workflows: %v", err)
			return
		}

		existingSlugs := make(map[string]bool)
		for _, wf := range resp.Results {
			existingSlugs[wf.Slug] = true
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.Errorf("Failed to read local workflows directory: %v", err)
			return
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			slug := strings.TrimSuffix(file.Name(), ".json")
			if _, exists := existingSlugs[slug]; exists {
				fmt.Printf("Skipping '%s.json' (already exists on server)\n", slug)
				continue
			}

			path := filepath.Join(dirPath, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				log.Errorf("Failed to read file %s: %v", file.Name(), err)
				continue
			}

			var workflow map[string]any
			if err := json.Unmarshal(data, &workflow); err != nil {
				log.Errorf("Failed to parse JSON for %s: %v", file.Name(), err)
				continue
			}

			err = mgmntClient.PushWorkflow(workspace, slug, workflow)
			if err != nil {
				log.Errorf("Failed to push workflow %s: %v", slug, err)
				continue
			}

			fmt.Printf("Pushed workflow: %s\n", slug)
		}
	},
}

var workflowCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit a draft workflow to live.",
	Long:  "Commits a draft workflow to live",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Category Slug is required")
		}
		workspace, _ := cmd.Flags().GetString("workspace")

		slug := args[0]

		mgmntClient := utils.GetSuprSendMgmntClient()

		commmitMessage, _ := cmd.Flags().GetString("commit-message")

		err := mgmntClient.CommitWorkflow(workspace, slug, commmitMessage)
		if err != nil {
			log.Errorf("Failed to commit workflow %s: %v", slug, err)
		}

		fmt.Printf("Committed workflow: %s\n", slug)
	},
}

func init() {
	workflowListCmd.Flags().IntP("limit", "l", 20, "Limit the number of workflows to list")
	workflowListCmd.Flags().IntP("offset", "f", 0, "Offset the number of workflows to list (default: 0)")
	// add flag to set mode which can be one of draft, live with validation of the flag
	workflowListCmd.Flags().StringP("mode", "m", "live", "Mode to list workflows (draft, live)")
	workflowListCmd.Flags().StringP("commit-message", "c", "", "Commit Message for making a workflow live")
	workflowListCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Parent().HelpFunc()(cmd, args)
	})
	workflowPullCmd.Flags().BoolVarP(&force, "force-dir", "f", false, "Create workflow directory without the permission")
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowPullCmd)
	workflowCmd.AddCommand(workflowPushCmd)
	workflowCmd.AddCommand(workflowCommitCmd)
}
