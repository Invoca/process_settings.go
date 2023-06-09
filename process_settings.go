// Package process_settings implements dynamic settings for the process.
package process_settings

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// A ProcessSettings is a collection of settings files and a target evaluator
// that can be used to get the value of a settings based on the current targeting.
type ProcessSettings struct {
	FilePath            string            // The path to the settings file that was used to create the ProcessSettings
	Settings            *[]SettingsFile   // The settings files that make up the ProcessSettings
	TargetEvaluator     TargetEvaluator   // The target evaluator that is used to determine which settings files are applicable
	Monitor             *fsnotify.Watcher // The file monitor that is used to detect changes to the settings file
	WhenUpdatedRegistry []func()          // A list of functions to call when the settings are updated
}

type SettingNotFound struct {
	SettingPath []string
}

func (e *SettingNotFound) Error() string {
	return fmt.Sprintf("The setting '%s' was not found", dotDelimitedSettingsPath(e.SettingPath))
}

// NewProcessSettingsFromFile creates a new instance of ProcessSettings by
// loading the settings from a specified file path and using the specified
// static context to evaluate the targeting.
func NewProcessSettingsFromFile(filePath string, staticContext map[string]interface{}) (*ProcessSettings, error) {
	settings, err := loadSettingsFromFile(filePath)
	if err != nil {
		return nil, err
	}

	monitor, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = monitor.Add(filePath)
	if err != nil {
		return nil, err
	}

	return &ProcessSettings{
		FilePath:        filePath,
		Settings:        settings,
		TargetEvaluator: TargetEvaluator{staticContext},
		Monitor:         monitor,
	}, nil
}

// Get returns the value of a setting based on the current targeting.
// If the requested setting is not found, an error is returned.
func (ps *ProcessSettings) Get(settingPath ...string) (interface{}, error) {
	var value interface{}

	valueFound := false
	for _, settingsFile := range *ps.Settings {
		if ps.TargetEvaluator.isTargetMatch(settingsFile) {
			if fileValue, keyExists := dig(settingsFile.Settings, settingPath...); keyExists {
				value = fileValue
				valueFound = true
			}
		}
	}

	if !valueFound {
		return nil, &SettingNotFound{settingPath}
	}

	return value, nil
}

// SafeGet returns the value of a setting based on the current targeting.
// If the requested setting is not found, nil is returned.
func (ps *ProcessSettings) SafeGet(settingPath ...string) (interface{}, error) {
	value, err := ps.Get(settingPath...)
	if err != nil && errors.Is(err, &SettingNotFound{}) {
		return nil, err
	}
	return value, nil
}

// StartMonitor starts a goroutine that monitors the settings file for changes
func (ps *ProcessSettings) StartMonitor() {
	go func() {
		defer ps.Monitor.Close()

		for {
			select {
			case event, ok := <-ps.Monitor.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					settings, err := loadSettingsFromFile(ps.FilePath)
					if err != nil {
						log.Println("Error processing new version of the process settings file:", err)
					}
					ps.Settings = settings
					for _, fn := range ps.WhenUpdatedRegistry {
						fn()
					}
				}
			case err, ok := <-ps.Monitor.Errors:
				if !ok {
					return
				}
				log.Println("Error reported from fsnotify:", err)
			}
		}
	}()
}

// WhenUpdated registers a function to be called when the settings are updated and by default calls the function immediately.
// Optionally false can be passed as the second argument to not call the function immediately.
// The function returns an index that can be used to cancel the function using CancelWhenUpdated.
func (ps *ProcessSettings) WhenUpdated(fn func(), initial_update ...bool) int {
	ps.WhenUpdatedRegistry = append(ps.WhenUpdatedRegistry, fn)
	if len(initial_update) == 0 || initial_update[0] == true {
		fn()
	}
	return len(ps.WhenUpdatedRegistry) - 1
}

// CancelWhenUpdated cancels a function that was registered using WhenUpdated.
func (ps *ProcessSettings) CancelWhenUpdated(index int) {
	ps.WhenUpdatedRegistry[index] = func() {}
}

func dig(settings interface{}, settingPath ...string) (interface{}, bool) {
	if settings == nil {
		return nil, false
	}

	settingsType := reflect.TypeOf(settings).Kind()

	if len(settingPath) == 1 {
		switch settingsType {
		case reflect.Map:
			value, keyExists := settings.(map[string]interface{})[settingPath[0]]
			return value, keyExists
		default:
			return nil, false
		}
	}

	switch reflect.TypeOf(settings).Kind() {
	case reflect.Map:
		return dig(settings.(map[string]interface{})[settingPath[0]], settingPath[1:]...)
	default:
		return nil, false
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

func dotDelimitedSettingsPath(settingPath []string) string {
	stringifiedSettingPath := make([]string, len(settingPath))
	for i, path := range settingPath {
		stringifiedSettingPath[i] = fmt.Sprintf("%v", path)
	}
	return strings.Join(stringifiedSettingPath, ".")
}
