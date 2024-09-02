package parser

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Nevoral/SpecGenzo/internal/model"
)

// CreateTomlFile creates a TOML file from the provided config struct
func CreateTomlSpec(filename string, w *model.WebSpecification) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(w); err != nil {
		return fmt.Errorf("could not encode config to TOML: %v", err)
	}

	return nil
}

// LoadFromTomlFile loads the TOML data from the specified file into the WebConfig struct.
func LoadFromTomlFile(filename string) (*model.WebSpecification, error) {
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
