package schema

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull schemas",
	Long:  `Pull schemas in a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir, _ := cmd.Flags().GetString("output-dir")
		mode, _ := cmd.Flags().GetString("mode")
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

		workspace, _ := cmd.Flags().GetString("workspace")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to create directory: %v\n", err)
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
		schemas, err := mgmnt_client.GetSchemas(workspace, mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled %d schemas", len(schemas.Results)))
		}
		stats, err := WriteSchemasToFiles(schemas, outputDir)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save schemas: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "%d schemas saved\n", stats.Success)
		fmt.Fprintf(os.Stdout, "%d schemas failed\n", stats.Failed)
		fmt.Fprintf(os.Stdout, "%d schemas skipped\n", stats.Total-stats.Success-stats.Failed)
		fmt.Fprintf(os.Stdout, "%d schemas already exist\n", stats.Total-stats.Success-stats.Failed)
	},
}

func init() {
	schemaPullCmd.PersistentFlags().StringP("output-dir", "d", "", "Output directory for schemas")
	schemaPullCmd.PersistentFlags().StringP("mode", "m", "live", "Mode of schemas to pull (draft, live), default: live")
	SchemaCmd.AddCommand(schemaPullCmd)
}
