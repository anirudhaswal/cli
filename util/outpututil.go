package util

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"

	// log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// OutputData chooses the output format based on the flag
func OutputData(data interface{}, format string) {
	switch format {
	case "json":
		outputJSON(data)
	case "yaml":
		outputYAML(data)
	default:
		outputTable(data)
	}
}

// outputJSON prints data in JSON format
func outputJSON(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error creating JSON output:", err)
		return
	}
	fmt.Println(string(jsonData))
}

// outputYAML prints data in YAML format
func outputYAML(data interface{}) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Error creating YAML output:", err)
		return
	}
	fmt.Println(string(yamlData))
}

// outputText prints data in plain text format
// func outputText(data interface{}) {
// 	fmt.Printf("Data: %v\n", data)
// }

// FetchData simulates data retrieval
func FetchData() interface{} {
	return map[string]string{"key": "value"}
}

func outputTable(data interface{}) {
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
			fmt.Println("No data to display")
			return
		}
		elemType := val.Index(0).Type()
		if elemType.Kind() != reflect.Struct {
			fmt.Println("Slice elements must be structs")
			return
		}
		values := make([]reflect.Value, val.Len())
		for i := 0; i < val.Len(); i++ {
			values[i] = val.Index(i)
		}
		printStructAsTable(values)
		return
	}

	fmt.Println("Input must be a struct or a slice of structs")
}

func printStructAsTable(values []reflect.Value) {
	if len(values) == 0 {
		fmt.Println("No data to display")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	// Set headers based on struct field names
	elemType := values[0].Type()
	var headers []string
	for i := 0; i < elemType.NumField(); i++ {
		headers = append(headers, elemType.Field(i).Name)
	}
	table.SetHeader(headers)

	// Add rows
	for _, val := range values {
		var row []string
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			row = append(row, formatValue(field))
		}
		table.Append(row)
	}

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
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
