package process_settings

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ProcessSettings struct {
	FilePath      string
	StaticContext map[string]interface{}
	Settings      *[]SettingsFile
}

func NewProcessSettingsFromFile(filePath string, staticContext map[string]interface{}) (*ProcessSettings, error) {
	settings, err := loadSettingsFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return &ProcessSettings{
		FilePath:      filePath,
		StaticContext: staticContext,
		Settings:      settings,
	}, nil
}

func loadSettingsFromFile(filePath string) (*[]SettingsFile, error) {
	var settings []SettingsFile
	err := loadYamlFile(filePath, &settings)
	if err != nil {
		return nil, err
	}

	for i, setting := range settings {
		valid, err := setting.IsValid()
		if !valid {
			return nil, errors.New(fmt.Sprintf("Invalid settings file at index %d: %s => %v", i, err.Error(), setting))
		}
	}

	if settings[len(settings)-1].Metadata.End != true {
		return nil, errors.New("The settings file does not have the END metadata")
	}
	return &settings, nil
}

func loadYamlFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}
	return nil
}
