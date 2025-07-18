package schema

import (
	"bufio"
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := filepath.Join(".", "suprsend", "schema")

		if len(args) < 1 {
			log.Error("workspace argument is required for schemas.")
		}
		workspace := args[0]

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if force {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					log.WithError(err).Error("Failed to create directory")
					return
				}
				log.Printf("Directory created at: %s\n", dirPath)
			} else {
				log.Printf("Directory '%s' does not exist. Do you want to create it? (y/n): ", dirPath)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(input))

				if input == "y" || input == "yes" {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						log.WithError(err).Error("Failed to create directory")
						return
					}
					log.Infof("Directory created at: %s", dirPath)
				} else {
					log.Error("Directory not created. Exiting.")
					return
				}
			}
		} else {
			log.Infof("Directory already exists: %s\n", dirPath)
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()

		schemas, err := mgmnt_client.GetSchemas(workspace)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		log.Infoln("Pulling schemas...")
		if err := writeSchemasToFiles(schemas, dirPath); err != nil {
			log.WithError(err).Error("Error saving schemas")
			return
		}
	},
}

func init() {
	SchemaCmd.AddCommand(schemaPullCmd)
}
