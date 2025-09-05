package commands

import (
	"context"
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
	"github.com/yarlson/pin"
)

var generateTypesCmd = &cobra.Command{
	Use:   "generate-types",
	Short: "Generate Types from JSON Schema",
	Long:  "Generate types from JSON schema for various programming languages",
}

var generateTypesPythonCmd = &cobra.Command{
	Use:   "python [flags] <output-file>",
	Short: "Generate Python types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run:   generateTypesForLanguage("python"),
}

var generateTypesJavaCmd = &cobra.Command{
	Use:   "java [flags] <package-name>",
	Short: "Generate Java types from JSON Schema",
	Long:  "Generate Java types from JSON Schema with specified package name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		buildFlags, _ := cmd.Flags().GetString("build-flags")
		mode, _ := cmd.Flags().GetString("mode")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		packageName := args[0]
		if packageName == "" {
			log.Error("Package name argument is required")
			return
		}
		if outputDir == "" {
			log.Error("Output directory is required for Java generation")
			return
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.WithError(err).Errorf("Failed to create output directory: %s", outputDir)
			return
		}
		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Generating Java types...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}
		mgmntClient := utils.GetSuprSendMgmntClient()
		schemasResp, err := mgmntClient.GetSchemas(workspace, mode)
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
			if utils.IsEmptySchema(schemaResp.JSONSchema.Properties) {
				continue
			}
			validSchemas = append(validSchemas, &schemaResp)
		}
		if len(validSchemas) == 0 {
			if p != nil {
				p.Stop("No valid schemas found")
			}
			fmt.Println("No valid schemas found with meaningful JSON schema content")
			return
		}

		generatedCount := 0
		for _, targetSchema := range validSchemas {
			schemaName := targetSchema.Name + "Data"
			fileName := filepath.Join(outputDir, schemaName+".java")
			if _, err := os.Stat(fileName); err == nil {
				if err := os.WriteFile(fileName, []byte(""), 0644); err != nil {
					log.WithError(err).Errorf("Failed to clear existing file: %s", fileName)
					continue
				}
			}
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
			javaFlags := buildFlags
			if javaFlags != "" {
				javaFlags += ",package=" + packageName
			} else {
				javaFlags = "package=" + packageName
			}
			err = runTypeMorph("java", string(schemaBytes), schemaName, fileName, javaFlags)
			if err != nil {
				log.WithError(err).Errorln("Could not generate types for schema: " + targetSchema.Name)
			} else {
				generatedCount++
			}
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Generated %d Java type files in %s", generatedCount, outputDir))
		}
	},
}

var generateTypesTypeScriptCmd = &cobra.Command{
	Use:   "typescript [flags] <output-file>",
	Short: "Generate TypeScript types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zod, _ := cmd.Flags().GetString("zod")
		if zod == "true" {
			generateTypesForLanguage("typescript-zod")(cmd, args)
		} else {
			generateTypesForLanguage("typescript")(cmd, args)
		}
	},
}

var generateTypesGoCmd = &cobra.Command{
	Use:   "go [flags] <output-file>",
	Short: "Generate Go types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run:   generateTypesForLanguage("go"),
}

var generateTypesKotlinCmd = &cobra.Command{
	Use:   "kotlin [flags] <output-file>",
	Short: "Generate Kotlin types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run:   generateTypesForLanguage("kotlin"),
}

var generateTypesSwiftCmd = &cobra.Command{
	Use:   "swift [flags] <output-file>",
	Short: "Generate Swift types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run:   generateTypesForLanguage("swift"),
}

var generateTypesDartCmd = &cobra.Command{
	Use:   "dart [flags] <output-file>",
	Short: "Generate Dart types from JSON Schema",
	Args:  cobra.ExactArgs(1),
	Run:   generateTypesForLanguage("dart"),
}

func generateTypesForLanguage(targetLang string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		buildFlags, _ := cmd.Flags().GetString("build-flags")
		mode, _ := cmd.Flags().GetString("mode")
		fileName := args[0]
		if fileName == "" {
			log.Error("File name argument is required")
			return
		}

		expectedExt := utils.GetExtensionForLanguage(targetLang)
		fileExtension := strings.ToLower(filepath.Ext(fileName))
		if fileExtension != expectedExt {
			log.Errorf("File extension %s doesn't match expected extension %s for %s", fileExtension, expectedExt, targetLang)
			return
		}

		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New(fmt.Sprintf("Generating %s types...", strings.Title(targetLang)),
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		mgmntClient := utils.GetSuprSendMgmntClient()
		schemasResp, err := mgmntClient.GetSchemas(workspace, mode)
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
			if utils.IsEmptySchema(schemaResp.JSONSchema.Properties) {
				continue
			}

			validSchemas = append(validSchemas, &schemaResp)
		}

		if len(validSchemas) == 0 {
			if p != nil {
				p.Stop("No valid schemas found")
			}
			fmt.Println("No valid schemas found with meaningful JSON schema content")
			return
		}

		if _, err := os.Stat(fileName); err == nil {
			if err := os.WriteFile(fileName, []byte(""), 0644); err != nil {
				log.WithError(err).Errorf("Failed to clear existing file: %s", fileName)
				return
			}
		}

		generatedCount := 0
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
				generatedCount++
			}
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Generated %d %s types in %s", generatedCount, targetLang, fileName))
		}
	}
}

func init() {
	commonFlags := []*cobra.Command{
		generateTypesPythonCmd,
		generateTypesTypeScriptCmd,
		generateTypesGoCmd,
		generateTypesKotlinCmd,
		generateTypesSwiftCmd,
		generateTypesDartCmd,
		generateTypesJavaCmd,
	}

	for _, cmd := range commonFlags {
		cmd.Flags().String("workspace", "staging", "Workspace to get schemas from.")
		cmd.Flags().String("mode", "live", "Mode of schema to fetch.")
		cmd.Flags().String("build-flags", "", "Flags to generate types in a certain way.")
	}

	generateTypesJavaCmd.Flags().String("output-dir", "", "Output directory for generated Java files (required)")
	generateTypesJavaCmd.MarkFlagRequired("output-dir")
	generateTypesCmd.AddCommand(generateTypesJavaCmd)
	generateTypesCmd.AddCommand(generateTypesPythonCmd)
	generateTypesTypeScriptCmd.Flags().String("zod", "false", "Generate Zod types for TypeScript")
	generateTypesCmd.AddCommand(generateTypesTypeScriptCmd)
	generateTypesCmd.AddCommand(generateTypesGoCmd)
	generateTypesCmd.AddCommand(generateTypesKotlinCmd)
	generateTypesCmd.AddCommand(generateTypesSwiftCmd)
	generateTypesCmd.AddCommand(generateTypesDartCmd)
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
