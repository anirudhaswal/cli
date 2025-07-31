package schema

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
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
		} else {
			log.Infof("Directory already exists: %s\n", dirPath)
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		schemas, err := mgmnt_client.GetSchemas(workspace)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Pulling schemas...\n")
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
