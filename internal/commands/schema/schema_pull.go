package schema

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var force bool

var schemaPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull schemas",
	Long:  `Pull schemas in a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "schema")

		workspace, _ := cmd.Flags().GetString("workspace")

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
		schemas, err := mgmnt_client.GetSchemas(workspace)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled %d schemas", len(schemas.Results)))
		}
		stats, err := WriteSchemasToFiles(schemas, dirPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save schemas: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Pulled %d schemas\n", stats.Total)
		fmt.Fprintf(os.Stdout, "%d schemas saved\n", stats.Success)
		fmt.Fprintf(os.Stdout, "%d schemas failed\n", stats.Failed)
		fmt.Fprintf(os.Stdout, "%d schemas skipped\n", stats.Total-stats.Success-stats.Failed)
		fmt.Fprintf(os.Stdout, "%d schemas already exist\n", stats.Total-stats.Success-stats.Failed)
	},
}

func init() {
	SchemaCmd.AddCommand(schemaPullCmd)
}
