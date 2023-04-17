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
func Get(settingPath ...string) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("The global process settings have not been set")
	}
	return instance.Get(settingPath...)
}

// SafeGet returns the value of a setting based on the current targeting.
// If the global instance has not been set, or the requested setting is not found,
// nil is returned.
func SafeGet(settingPath ...string) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("The global process settings have not been set")
	}
	return instance.SafeGet(settingPath...)
}

// WhenUpdated registers a function to be called when the settings are updated on the global ProcessSettings instance.
// If the global instance has not been set, an error is returned.
func WhenUpdated(fn func(), initial_update ...bool) (int, error) {
	if instance == nil {
		return 0, errors.New("The global process settings have not been set")
	}
	return instance.WhenUpdated(fn, initial_update...), nil
}

// CancelWhenUpdated cancels a function that was registered on the global ProcessSettings instance using WhenUpdated.
// If the global instance has not been set, an error is returned.
func CancelWhenUpdated(index int) error {
	if instance == nil {
		return errors.New("The global process settings have not been set")
	}
	instance.CancelWhenUpdated(index)
	return nil
}
