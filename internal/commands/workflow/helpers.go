package workflow

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

type WorkflowWriteStats struct {
	Total   int
	Success int
	Failed  int
	Errors  []string
}

func isDebugMode() bool {
	return viper.GetBool("debug")
}

func promptForOutputDirectory() string {
	reader := bufio.NewReader(os.Stdin)
	defaultDir := filepath.Join(".", "suprsend", "workflow")
	fmt.Fprintf(os.Stdout, "Where would you like to save the workflows?\n")
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

func WriteWorkflowsToFiles(resp mgmnt.WorkflowsResponse, outputDir string) (*WorkflowWriteStats, error) {
	stats := &WorkflowWriteStats{
		Total:  len(resp.Results),
		Errors: []string{},
	}

	info, err := os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return stats, err
			}
		} else {
			errMsg := fmt.Sprintf("error accessing '%s': %v", outputDir, err)
			return stats, fmt.Errorf(errMsg)
		}
	} else if !info.IsDir() {
		return stats, err
	}

	for _, wf := range resp.Results {
		obj, ok := wf.(map[string]any)
		if !ok {
			stats.Failed++
			stats.Errors = append(stats.Errors, "Invalid workflow format")
			continue
		}

		slug, _ := obj["slug"].(string)
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", slug))

		fileData, err := json.MarshalIndent(wf, "", "  ")
		if err != nil {
			debugErrorLog("Error: %s", err)
			fmt.Fprintf(os.Stdout, "Error: Failed to marshal workflow '%s': %v\n", slug, err)
			stats.Failed++
			stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to marshal workflow '%s': %v", slug, err))
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
		fmt.Fprintf(os.Stdout, "Wrote workflow to %s\n", filename)
		stats.Success++
	}

	return stats, nil
}
