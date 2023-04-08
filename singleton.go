package process_settings

import "errors"

var instance *ProcessSettings

// SetGlobalProcessSettings sets the global process settings instance
// to be used by the rest of the process.
func SetGlobalProcessSettings(settings *ProcessSettings) {
	instance = settings
}

// Get returns the value of a setting based on the current targeting.
// If the global instance has not been set, or the requested setting is not found,
// an error is returned.
func Get(settingPath ...interface{}) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("The global process settings have not been set")
	}
	return instance.Get(settingPath...)
}

// SafeGet returns the value of a setting based on the current targeting.
// If the global instance has not been set, or the requested setting is not found,
// nil is returned.
func SafeGet(settingPath ...interface{}) interface{} {
	if instance == nil {
		return nil
	}
	return instance.SafeGet(settingPath...)
}
