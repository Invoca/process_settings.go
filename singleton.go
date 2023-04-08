package process_settings

import "errors"

var instance *ProcessSettings

func SetGlobalProcessSettings(settings *ProcessSettings) {
	instance = settings
}

func Get(settingPath ...interface{}) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("The global process settings have not been set")
	}
	return instance.Get(settingPath...)
}

func SafeGet(settingPath ...interface{}) interface{} {
	if instance == nil {
		return nil
	}
	return instance.SafeGet(settingPath...)
}
