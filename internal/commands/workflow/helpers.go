package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/mgmnt"
)

func WriteWorkflowsToFiles(resp mgmnt.WorkflowsResponse, outputDir string) error {
	info, err := os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", outputDir, err)
			}
		} else {
			return fmt.Errorf("error accessing '%s': %w", outputDir, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("'%s' exists but is not a directory", outputDir)
	}

	for _, wf := range resp.Results {
		obj, ok := wf.(map[string]any)
		if !ok {
			continue
		}

		slug, _ := obj["slug"].(string)
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", slug))

		if _, err := os.Stat(filename); err == nil {
			log.Printf("Skipped (already exists): %s\n", filename)
			continue
		} else if err != nil && !os.IsNotExist(err) {
			log.Errorf("Error checking file '%s': %v\n", filename, err)
			return err
		}

		fileData, err := json.MarshalIndent(wf, "", "  ")
		if err != nil {
			log.Errorf("Error marshalling workflow '%s': %v\n", slug, err)
			continue
		}

		if err := os.WriteFile(filename, fileData, 0644); err != nil {
			log.Errorf("Error writing file '%s': %v\n", filename, err)
			continue
		}

		log.Infof("Wrote: %s\n", filename)
	}

	return nil
}
