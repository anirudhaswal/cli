package schema

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/suprsend/cli/mgmnt"
)

type SchemaWriteStats struct {
	Total   int
	Success int
	Failed  int
	Errors  []string
}

type FilteredSchema struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func isDebugMode() bool {
	return viper.GetBool("debug")
}

func debugLog(format string, args ...interface{}) {
	if isDebugMode() {
		log.Infof(format, args...)
	}
}

func debugErrorLog(format string, args ...interface{}) {
	if isDebugMode() {
		log.Errorf(format, args...)
	}
}

func promptForOutputDirectory() string {
	reader := bufio.NewReader(os.Stdin)
	defaultDir := filepath.Join(".", "suprsend", "schema")
	fmt.Fprintf(os.Stdout, "Where would you like to save the schema?\n")
	fmt.Fprintf(os.Stdout, "Default: %s\n", defaultDir)
	fmt.Fprintf(os.Stdout, "Enter directory path (or press Enter for default): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultDir
	}
	return input
}

func ensureOutputDirectory(dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stdout, "Creating directory: %s\n", dirPath)
			return os.MkdirAll(dirPath, 0755)
		}
		return fmt.Errorf("error checking directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", dirPath)
	}
	if info.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("directory '%s' is not writable", dirPath)
	}
	return nil
}

func validateInputDirectory(dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("error checking directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", dirPath)
	}
	if info.Mode().Perm()&0400 == 0 {
		return fmt.Errorf("directory '%s' is not readable", dirPath)
	}
	return nil
}

func filterSchemaData(schemas []mgmnt.SchemaResponse) []FilteredSchema {
	filtered := make([]FilteredSchema, len(schemas))
	for i, schema := range schemas {
		filtered[i] = FilteredSchema{
			Slug:        schema.Slug,
			Name:        schema.Name,
			Description: schema.Description,
		}
	}
	return filtered
}

func WriteSchemasToFiles(schemasResp *mgmnt.SchemasResponse, dirPath string) (*SchemaWriteStats, error) {
	stats := &SchemaWriteStats{
		Total:  len(schemasResp.Results),
		Errors: []string{},
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return stats, err
			}
		} else {
			errMsg := fmt.Sprintf("error accessing '%s': %v", dirPath, err)
			return stats, fmt.Errorf(errMsg)
		}
	} else if !info.IsDir() {
		return stats, err
	}

	for _, schema := range schemasResp.Results {
		obj, ok := schema.(map[string]any)
		if !ok {
			stats.Failed++
			stats.Errors = append(stats.Errors, "Invalid schema format")
			continue
		}

		slug, _ := obj["slug"].(string)
		filename := filepath.Join(dirPath, fmt.Sprintf("%s.json", slug))

		fileData, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			debugErrorLog("Error: %s", err)
			fmt.Fprintf(os.Stdout, "Error: Failed to marshal schema '%s': %v\n", slug, err)
			stats.Failed++
			stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to marshal schema '%s': %v", slug, err))
			continue
		}

		if err := os.WriteFile(filename, fileData, 0644); err != nil {
			debugErrorLog("Error: %s", err)
			fmt.Fprintf(os.Stdout, "Error: Failed to write file '%s': %v\n", filename, err)
			stats.Failed++
			stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to write file '%s': %v", filename, err))
			continue
		}

		debugLog("Wrote: %s", filename)
		fmt.Fprintf(os.Stdout, "Wrote schema to %s\n", filename)
		stats.Success++
	}

	return stats, nil
}

type SchemasResponse struct {
	Schemas []SchemaResponse `json:"schemas"`
}

type SchemaResponse struct {
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsEnabled   bool       `json:"is_enabled"`
	JSONSchema  JSONSchema `json:"json_schema"`
}

type JSONSchema struct {
	Type       string              `json:"type"`
	Title      string              `json:"title"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type string `json:"type"`
}
