package process_settings

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProcessSettings struct {
	FilePath        string
	Settings        *[]SettingsFile
	TargetEvaluator TargetEvaluator
}

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

func (ps *ProcessSettings) SafeGet(settingPath ...interface{}) interface{} {
	var value interface{}

	for _, settingsFile := range *ps.Settings {
		if ps.TargetEvaluator.IsTargetMatch(settingsFile) {
			if fileValue := dig(settingsFile.Settings, settingPath...); fileValue != nil {
				value = fileValue
			}
		}
	}

	return value
}

func dig(settings interface{}, settingPath ...interface{}) interface{} {
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
