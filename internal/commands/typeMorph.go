package commands

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

// GET v1/:ws/schema -- bulk list
// GET v1/:ws/schema/:slug_schema -- get a single schema

var typeMorphCmd = &cobra.Command{
	Use:   "type-morph",
	Short: "Generate Types from JSON Schema from a single slug",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("Please pass the required arguments [schema_name] [target_lang]")
		}
		schemaSlug := args[0]
		targetLang := args[1]
		workspace, _ := cmd.Flags().GetString("workspace")
		fileName, err := cmd.Flags().GetString("file")
		if err != nil {
			log.WithError(err).Error("Could not read file_name flag")
			return
		}

		mgmntClient := utils.GetSuprSendMgmntClient()

		resp, err := mgmntClient.GetSchema(workspace, schemaSlug)
		if err != nil {
			log.WithError(err).Error("Couldn't fetch schema")
			return
		}

		schemaName := resp.Name

		schemaJSON := map[string]interface{}{
			"type":       resp.JSONSchema.Type,
			"properties": resp.JSONSchema.Properties,
			"required":   resp.JSONSchema.Required,
		}

		schemaBytes, err := json.MarshalIndent(schemaJSON, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal schema: %v", err)
		}

		err = runTypeMorph(targetLang, string(schemaBytes), schemaName, fileName)
		if err != nil {
			log.WithError(err).Error("Could not generate types")
		}
	},
}

func init() {
	flags := typeMorphCmd.Flags()
	flags.String("workspace", "staging", "Workspace to get schemas from.")
	flags.String("file", "", "Target file to generate types into. (required)")

	rootCmd.AddCommand(typeMorphCmd)
}

func runTypeMorph(language, schema, schemaName, fileName string) error {
	binaryPath, err := writeTempExecutable(utils.TypeMorphBin)
	if err != nil {
		panic(fmt.Errorf("failed to write temp binary: %w", err))
	}
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, language, schema, schemaName, fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run typemorph: %w", err)
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
