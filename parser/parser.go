package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Nevoral/SpecGenzo/model"
)

// CreateJsonFile creates a JSON file from the provided config struct
func CreateJsonSpec(filename string, w *model.WebSpecification) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	spec, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling user to JSON:", err)
	}

	_, err = file.Write(spec)
	if err != nil {
		return fmt.Errorf("Error writing JSON data to file:", err)
	}

	return nil
}

// LoadFromJsonFile loads the JSON data from the specified file into the WebConfig struct.
func LoadFromJsonFile(filename string) (*model.WebSpecification, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("could not close file:", err)
		}
	}(file)
	var w model.WebSpecification
	if _, err = toml.NewDecoder(file).Decode(w); err != nil {
		return nil, fmt.Errorf("could not decode TOML: %v", err)
	}

	return &w, nil
}
