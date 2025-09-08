package category

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteToFile(data interface{}, filename string) error {
	os.MkdirAll("./suprsend/category", 0755)
	filename = "./suprsend/category/" + filename
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return os.WriteFile(filename, jsonData, 0644)
}

func ReadFromFile(filepath string) (interface{}, error) {
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}
