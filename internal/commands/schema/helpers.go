package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/mgmnt"
)

func writeSchemasToFiles(schemasResp *mgmnt.SchemasResponse, dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", dirPath, err)
			}
		} else {
			return fmt.Errorf("error accessing '%s': %w", dirPath, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("'%s' exists but is not a directory", dirPath)
	}

	for _, schema := range schemasResp.Results {
		obj, ok := schema.(map[string]any)
		if !ok {
			continue
		}

		slug, _ := obj["slug"].(string)

		filename := filepath.Join(dirPath, fmt.Sprintf("%s.json", slug))

		if _, err := os.Stat(filename); err == nil {
			log.Infof("Skipped (already exists): %s\n", filename)
			continue
		} else if err != nil && !os.IsNotExist(err) {
			log.Errorf("Error checking file '%s': %v\n", filename, err)
			continue
		}

		fileData, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			log.Errorf("Error marshalling schema '%s': %v\n", slug, err)
			continue
		}

		if err := os.WriteFile(filename, fileData, 0644); err != nil {
			log.Errorf("Error writing file '%s': %v\n", filename, err)
			continue
		}

		log.Printf("Wrote: %s\n", filename)
	}

	return nil
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
