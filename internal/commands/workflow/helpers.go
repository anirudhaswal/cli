package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suprsend/cli/mgmnt"
)

func writeWorkflowsToFiles(resp mgmnt.WorkflowAPIResponse, outputDir string) error {
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
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", wf.Slug))

		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("Skipped (already exists): %s\n", filename)
			continue
		} else if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error checking file '%s': %v\n", filename, err)
			continue
		}

		fileData, err := json.MarshalIndent(wf, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling workflow '%s': %v\n", wf.Slug, err)
			continue
		}

		if err := os.WriteFile(filename, fileData, 0644); err != nil {
			fmt.Printf("Error writing file '%s': %v\n", filename, err)
			continue
		}

		fmt.Printf("Wrote: %s\n", filename)
	}

	return nil
}
