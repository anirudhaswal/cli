package commands

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/cli/mgmnt"
)

var generateTypesCmd = &cobra.Command{
	Use:   "generate-types",
	Short: "Generate Types from JSON Schema from a single slug",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		buildFlags, _ := cmd.Flags().GetString("build-flags")
		fileName := args[0]
		if fileName == "" {
			log.Error("File argument is required")
			return
		}

		targetLang := detectLanguageFromFile(fileName)
		if targetLang == "" {
			fileExtension := strings.ToLower(filepath.Ext(fileName))
			log.Errorf("Unsupported file extension: %s. We currently support .ts, .go, .py, .kt, .swift files only.", fileExtension)
			return
		}

		mgmntClient := utils.GetSuprSendMgmntClient()

		schemasResp, err := mgmntClient.GetSchemas(workspace)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schemas")
			return
		}

		var validSchemas []*mgmnt.SchemaResponse
		for _, schema := range schemasResp.Results {
			schemaBytes, err := json.Marshal(schema)
			if err != nil {
				continue
			}
			var schemaResp mgmnt.SchemaResponse
			if err := json.Unmarshal(schemaBytes, &schemaResp); err != nil {
				continue
			}

			if !schemaResp.IsEnabled {
				continue
			}
			if utils.IsEmptySchema(schemaResp.JSONSchema.Properties) {
				continue
			}

			validSchemas = append(validSchemas, &schemaResp)
		}

		if len(validSchemas) == 0 {
			fmt.Println("No valid schemas found with meaningful JSON schema content")
			return
		}

		for _, targetSchema := range validSchemas {
			schemaName := targetSchema.Name + "Data"

			schemaJSON := map[string]interface{}{
				"type":       targetSchema.JSONSchema.Type,
				"properties": targetSchema.JSONSchema.Properties,
			}
			if targetSchema.JSONSchema.Required != nil {
				schemaJSON["required"] = targetSchema.JSONSchema.Required
			}
			if targetSchema.JSONSchema.Defs != nil {
				schemaJSON["$defs"] = targetSchema.JSONSchema.Defs
			}
			schemaBytes, err := json.MarshalIndent(schemaJSON, "", "  ")
			if err != nil {
				log.Fatalf("Failed to marshal schema: %v", err)
			}

			err = runTypeMorph(targetLang, string(schemaBytes), schemaName, fileName, buildFlags)
			if err != nil {
				log.WithError(err).Errorln("Could not generate types for schema: " + targetSchema.Name)
			} else {
				fmt.Printf("Generated types for %s at %s\n", schemaName, fileName)
			}
		}
	},
}

func detectLanguageFromFile(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	return utils.LanguageMap[ext]
}

func init() {
	flags := generateTypesCmd.Flags()
	flags.String("workspace", "staging", "Workspace to get schemas from.")
	flags.String("mode", "live", "Mode of schema to fetch.")
	flags.String("build-flags", "", "Flags to generate types in a certain way.")

	rootCmd.AddCommand(generateTypesCmd)
}

func runTypeMorph(language, schema, schemaName, fileName, buildFlags string) error {
	binaryPath, err := writeTempExecutable(utils.TypeMorphBin)
	if err != nil {
		return fmt.Errorf("failed to initialize type generator: %w", err)
	}
	defer os.Remove(binaryPath)

	args := []string{language, schema, schemaName, fileName}
	if buildFlags != "" {
		args = append(args, "--build-flags="+buildFlags)
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("type generation failed for %s: %w", fileName, err)
	}
	return nil
}

func writeTempExecutable(data []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "typemorph-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		return "", err
	}

	// Make it executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
