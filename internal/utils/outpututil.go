/*
Copyright © 2025 SuprSend
*/
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
)

func supportsColor() bool {
	// check if output is redirected to a file
	fileInfo, _ := os.Stdout.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	// check if stdout is a tty
	if viper.GetBool("NO_COLOR") {
		return false
	}

	return true
}

// OutputData chooses the output format based on the flag
func OutputData(data any, format string) {
	switch format {
	case "json":
		outputJSON(data)
	case "yaml":
		outputYAML(data)
	default:
		// outputJSON(data)
		outputTable(data)
	}
}

// outputJSON prints data in JSON format
func outputJSON(data any) {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		log.Fatal("Error creating JSON output:", err)
		return
	}

	if supportsColor() {
		fmt.Println(string(pretty.Color(jsonData, nil)))
	} else {
		fmt.Println(string(jsonData))
	}
}

func colorizeYAML(yamlString string) string {
	lines := []string{}
	for _, line := range strings.Split(yamlString, "\n") {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}
			lines = append(lines, color.HiBlueString(key)+":"+color.GreenString(value))
		} else {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

// outputYAML prints data in YAML format
func outputYAML(data any) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal("Error creating YAML output:", err)
		return
	}
	// trim the yaml data
	yamlString := strings.TrimSpace(string(yamlData))
	if supportsColor() {
		fmt.Println(colorizeYAML(yamlString))
	} else {
		fmt.Println(yamlString)
	}
}

func outputTable(data any) {
	val := reflect.ValueOf(data)

	// If the input is a pointer, get the underlying element
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle single struct
	if val.Kind() == reflect.Struct {
		printStructAsTable([]reflect.Value{val})
		return
	}

	// Handle slice of structs
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			log.Fatal("No data to display")
			return
		}
		elemType := val.Index(0).Type()
		if elemType.Kind() != reflect.Struct {
			log.Fatal("Slice elements must be structs")
			return
		}
		values := make([]reflect.Value, val.Len())
		for i := 0; i < val.Len(); i++ {
			values[i] = val.Index(i)
		}
		printStructAsTable(values)
		return
	}

	log.Fatal("Input must be a struct or a slice of structs")
}

func printStructAsTable(values []reflect.Value) {
	if len(values) == 0 {
		log.Fatal("No data to display")
		return
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{
				Separators: tw.Separators{BetweenRows: tw.On},
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{Alignment: tw.AlignCenter},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					MergeMode: tw.MergeHierarchical,
					Alignment: tw.AlignLeft,
				},
			},
		}),
	)

	// Set headers based on struct field names
	elemType := values[0].Type()
	var headers []string
	for i := 0; i < elemType.NumField(); i++ {
		headers = append(headers, elemType.Field(i).Name)
	}
	table.Header(headers)

	// Add rows
	var rows [][]any
	for _, val := range values {
		var row []any
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			row = append(row, formatValue(field))
		}
		rows = append(rows, row)
	}
	table.Bulk(rows)

	table.Render()
}

func formatValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 2, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Map:
		// Check if it's map[string]any
		if v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.Interface {
			m := v.Interface()
			if b, err := json.Marshal(m); err == nil {
				return string(b)
			}
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
