package commands

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var typeMorphCmd = &cobra.Command{
	Use:   "type-morph",
	Short: "Generate Types from JSON Schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		targetLang, err := cmd.Flags().GetString("language")
		if err != nil {
			return fmt.Errorf("could not read 'language' flag: %w", err)
		}

		schema, err := cmd.Flags().GetString("schema")
		if err != nil {
			return fmt.Errorf("could not read 'schema' flag: %w", err)
		}

		schemaName, err := cmd.Flags().GetString("schemaName")
		if err != nil {
			return fmt.Errorf("could not read 'schemaName' flag: %w", err)
		}

		fileName, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("could not read 'file' flag: %w", err)
		}

		return runTypeMorph(targetLang, schema, schemaName, fileName)
	},
}

func init() {
	flags := typeMorphCmd.Flags()
	flags.String("language", "", "Target language to generate types for. (required)")
	flags.String("schema", "", "JSON Schema file to generate types from.(required)")
	flags.String("schemaName", "SchemaType", "SchemaName for the types generated.")
	flags.String("file", "", "Target file to generate types into. (required)")

	typeMorphCmd.MarkFlagRequired("language")
	typeMorphCmd.MarkFlagRequired("schema")
	typeMorphCmd.MarkFlagRequired("file")

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
