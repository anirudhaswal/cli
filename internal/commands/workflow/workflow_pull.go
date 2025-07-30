package workflow

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var workflowPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull workflows from workspace to local directory",
	Long:  `pull workflows from workspace to local directory of each workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "workflow")

		workspace, _ := cmd.Flags().GetString("workspace")
		mode, _ := cmd.Flags().GetString("mode")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if force {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					fmt.Fprintf(os.Stdout, "Error: Failed to create directory: %v\n", err)
					return
				}
				fmt.Fprintf(os.Stdout, "Success: Directory created at %s\n", dirPath)
			} else {
				fmt.Fprintf(os.Stdout, "Directory '%s' does not exist. Do you want to create it? (y/n): ", dirPath)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(input))

				if input == "y" || input == "yes" {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						fmt.Fprintf(os.Stdout, "Error: Failed to create directory: %v\n", err)
						return
					}
					fmt.Fprintf(os.Stdout, "Success: Directory created at %s\n", dirPath)
				} else {
					fmt.Fprintf(os.Stdout, "Error: Directory not created. Exiting.\n")
					return
				}
			}
		} else {
			fmt.Fprintf(os.Stdout, "Directory already exists: %s\n", dirPath)
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		workflows_resp, err := mgmnt_client.GetWorkflows(workspace, mode)

		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get workflows: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Pulling workflows...\n")

		stats, err := WriteWorkflowsToFiles(*workflows_resp, dirPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save workflows: %v\n", err)
			return
		}

		// Final summary
		fmt.Fprintf(os.Stdout, "\n=== Workflow Pull Summary ===\n")
		fmt.Fprintf(os.Stdout, "Total workflows processed: %d\n", stats.Total)
		fmt.Fprintf(os.Stdout, "Successfully written: %d\n", stats.Success)
		fmt.Fprintf(os.Stdout, "Failed to write: %d\n", stats.Failed)

		if stats.Failed > 0 {
			fmt.Fprintf(os.Stdout, "\nFailed workflows:\n")
			for _, errorMsg := range stats.Errors {
				fmt.Fprintf(os.Stdout, "  - %s\n", errorMsg)
			}
		}

		if stats.Success > 0 {
			fmt.Fprintf(os.Stdout, "\nSuccess: Workflows pulled successfully to %s\n", dirPath)
		}
	},
}

func init() {
	WorkflowCmd.AddCommand(workflowPullCmd)
}
