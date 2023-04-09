// Package process_settings implements dynamic settings for the process.
package process_settings

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// A ProcessSettings is a collection of settings files and a target evaluator
// that can be used to get the value of a settings based on the current targeting.
type ProcessSettings struct {
	FilePath        string          // The path to the settings file that was used to create the ProcessSettings
	Settings        *[]SettingsFile // The settings files that make up the ProcessSettings
	TargetEvaluator TargetEvaluator // The target evaluator that is used to determine which settings files are applicable
}

// NewProcessSettingsFromFile creates a new instance of ProcessSettings by
// loading the settings from a specified file path and using the specified
// static context to evaluate the targeting.
func NewProcessSettingsFromFile(filePath string, staticContext map[string]interface{}) (*ProcessSettings, error) {
	settings, err := loadSettingsFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return &ProcessSettings{
		FilePath:        filePath,
		Settings:        settings,
		TargetEvaluator: TargetEvaluator{staticContext},
	}, nil
}

// Get returns the value of a setting based on the current targeting.
// If the requested setting is not found, an error is returned.
func (ps *ProcessSettings) Get(settingPath ...interface{}) (interface{}, error) {
	value := ps.SafeGet(settingPath...)

	if value == nil {
		stringifiedSettingPath := make([]string, len(settingPath))
		for i := range settingPath {
			stringifiedSettingPath[i] = fmt.Sprintf("%v", settingPath[i])
		}
		return nil, errors.New(fmt.Sprintf("The setting '%s' was not found", strings.Join(stringifiedSettingPath, ".")))
	}

	return value, nil
}

// SafeGet returns the value of a setting based on the current targeting.
// If the requested setting is not found, nil is returned.
func (ps *ProcessSettings) SafeGet(settingPath ...interface{}) interface{} {
	var value interface{}

	for _, settingsFile := range *ps.Settings {
		if ps.TargetEvaluator.isTargetMatch(settingsFile) {
			if fileValue := dig(settingsFile.Settings, settingPath...); fileValue != nil {
				value = fileValue
			}
		}
	}

	return value
}

func dig(settings interface{}, settingPath ...interface{}) interface{} {
	if settings == nil {
		return nil
	}

	settingsType := reflect.TypeOf(settings).Kind()

	if len(settingPath) == 1 {
		switch settingsType {
		case reflect.Map:
			nextKey := settingPath[0].(string)
			return settings.(map[string]interface{})[nextKey]
		case reflect.Slice:
			if reflect.TypeOf(settingPath[0]).Kind() == reflect.Int {
				return settings.([]interface{})[settingPath[0].(int)]
			}
		default:
			return nil
		}
	}

	switch reflect.TypeOf(settings).Kind() {
	case reflect.Map:
		nextKey := settingPath[0].(string)
		return dig(settings.(map[string]interface{})[nextKey], settingPath[1:]...)
	case reflect.Slice:
		if reflect.TypeOf(settingPath[0]).Kind() == reflect.Int {
			return dig(settings.([]interface{})[settingPath[0].(int)], settingPath[1:]...)
		}
		return nil
	default:
		return nil
	}
}

func loadSettingsFromFile(filePath string) (*[]SettingsFile, error) {
	var settings []SettingsFile
	err := loadYamlFile(filePath, &settings)
	if err != nil {
		return nil, err
	}

	for i, setting := range settings {
		valid, err := setting.isValid()
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
