package workflow

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var workflowPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull workflows from workspace to local directory",
	Long:  `pull workflows from workspace to local directory of each workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		mode, _ := cmd.Flags().GetString("mode")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir == "" {
			outputDir = promptForOutputDirectory()
			if outputDir == "" {
				fmt.Fprintf(os.Stdout, "No output directory specified. Exiting.\n")
				return
			}
		}
		if err := ensureOutputDirectory(outputDir); err != nil {
			fmt.Fprintf(os.Stdout, "Error with output directory: %v\n", err)
			return
		}

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
		workflows_resp, err := mgmnt_client.GetWorkflows(workspace, mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get workflows: %v\n", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled %d workflows from %s", len(workflows_resp.Results), workspace))
		}

		stats, err := WriteWorkflowsToFiles(*workflows_resp, outputDir)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save workflows: %v\n", err)
			return
		}

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
	},
}

func init() {
	workflowPullCmd.PersistentFlags().StringP("mode", "m", "live", "Mode of workflows to pull from (draft, live), default: live")
	workflowPullCmd.PersistentFlags().StringP("output-dir", "d", "", "Output directory for workflows")
	WorkflowCmd.AddCommand(workflowPullCmd)
}
