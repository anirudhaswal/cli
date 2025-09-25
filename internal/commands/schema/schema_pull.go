package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var schemaPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull schemas",
	Long:  `Pull schemas in a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir, _ := cmd.Flags().GetString("dir")
		mode, _ := cmd.Flags().GetString("mode")
		slug, _ := cmd.Flags().GetString("slug")
		if outputDir == "" {
			outputDir = filepath.Join(".", "suprsend", "schema")
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				outputDir = promptForOutputDirectory()
			}
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
		if slug != "" {
			schema, err := mgmnt_client.GetSchemaBySlug(workspace, slug, mode)
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error: Failed to get schema: %v\n", err)
				return
			}
			schemaData, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error: Failed to marshal schema: %v\n", err)
				return
			}
			if p != nil {
				p.Stop(fmt.Sprintf("Pulled %s from %s", slug, workspace))
			}
			os.WriteFile(filepath.Join(outputDir, slug+".json"), schemaData, 0644)
			return
		}
		schemas, err := mgmnt_client.GetSchemas(workspace, mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled %d schemas", len(schemas.Results)))
		}
		_, err = WriteSchemasToFiles(schemas, outputDir)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save schemas: %v\n", err)
			return
		}
	},
}

func init() {
	schemaPullCmd.Flags().StringP("dir", "d", "", "Directory to pull schemas (default: ./suprsend/schema)")
	schemaPullCmd.PersistentFlags().StringP("mode", "m", "live", "Mode of schemas to pull (draft, live), default: live")
	schemaPullCmd.PersistentFlags().StringP("slug", "g", "", "Slug of schema to pull")
	SchemaCmd.AddCommand(schemaPullCmd)
}
