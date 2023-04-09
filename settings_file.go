package process_settings

import "errors"

type SettingsMetadata struct {
	Version int  `yaml:"version"`
	End     bool `yaml:"END"`
}

type SettingsFile struct {
	FileName string                 `yaml:"filename"`
	Target   map[string]interface{} `yaml:"target"`
	Settings map[string]interface{} `yaml:"settings"`
	Metadata SettingsMetadata       `yaml:"meta"`
}

func (s *SettingsFile) IsValid() (bool, error) {
	if s.Metadata != (SettingsMetadata{}) {
		if s.FileName != "" || s.Target != nil || s.Settings != nil {
			return false, errors.New("The settings file must only have settings or metadata, not both")
		}
		return true, nil
	}

	if s.FileName == "" || s.Settings == nil {
		return false, errors.New("The settings file must have a filename and settings")
	}

	return true, nil
}
