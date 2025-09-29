package event

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

type EventWriteStats struct {
	Total   int
	Success int
	Failed  int
	Errors  []string
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
	defaultDir := filepath.Join(".", "suprsend", "event")
	fmt.Fprintf(os.Stdout, "Where would you like to save the events?\n")
	fmt.Fprintf(os.Stdout, "Default: %s\n", defaultDir)
	fmt.Fprintf(os.Stdout, "Enter directory path (or press Enter for default): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultDir
	}
	return input
}

func WriteEventsToFiles(events_resp *mgmnt.EventsResponse, dirPath string) (*EventWriteStats, error) {
	stats := &EventWriteStats{
		Total:  len(events_resp.Results),
		Errors: []string{},
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return stats, err
			}
		} else {
			return stats, err
		}
	} else if !info.IsDir() {
		return stats, err
	}

	eventSchemaMapping := map[string]interface{}{
		"events": events_resp.Results,
	}

	filename := filepath.Join(dirPath, "event_schema_mapping.json")
	fileData, err := json.MarshalIndent(eventSchemaMapping, "", "  ")
	if err != nil {
		debugErrorLog("Failed to marshal event schema mapping: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to marshal event schema mapping: %v\n", err)
		stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to marshal events: %v", err))
		return stats, err
	}
	if err := os.WriteFile(filename, fileData, 0644); err != nil {
		debugErrorLog("Failed to write event schema mapping to file: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to write event schema mapping to file: %v\n", err)
		stats.Failed++
		stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to write events: %v", err))
		return stats, err
	}

	debugLog("Successfully wrote %d events to %s", len(events_resp.Results), filename)
	fmt.Fprintf(os.Stdout, "Successfully wrote events to %s\n", filename)
	stats.Success = len(events_resp.Results)
	return stats, nil
}
