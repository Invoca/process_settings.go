package process_settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProcessSettingsFromFile(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		expectedError   []string
		expectedSize    int
		expectedVersion int
	}{
		{
			name:          "The file does not exist",
			fileName:      "fixtures/does_not_exist.yml",
			expectedError: []string{"open fixtures/does_not_exist.yml: no such file or directory"},
		},
		{
			name:          "The file is not a valid yaml",
			fileName:      "fixtures/invalid_yaml.yml",
			expectedError: []string{"yaml: unmarshal errors:\n  line 5: cannot unmarshal !!map into []process_settings.SettingsFile"},
		},
		{
			name:          "The file doesn't have metadata as the last element of the array",
			fileName:      "fixtures/invalid_metadata.yml",
			expectedError: []string{"The settings file does not have the END metadata"},
		},
		{
			name:          "The file has an invalid settings file",
			fileName:      "fixtures/invalid_settings.yml",
			expectedError: []string{"Invalid settings file at index 0: The settings file must only have settings or metadata, not both => {honeypot.yml map[] map[honeypot:map[answer_odds:100 max_recording_seconds:600 status_change_min_days:7]] {17 true}}"},
		},
		{
			name:            "The file is valid",
			fileName:        "fixtures/combined_process_settings.yml",
			expectedSize:    6,
			expectedVersion: 17,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			settings, err := NewProcessSettingsFromFile(test.fileName, nil)
			if test.expectedError == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedSize, len(*settings.Settings))
				assert.Equal(t, test.expectedVersion, (*settings.Settings)[test.expectedSize-1].Metadata.Version)
			} else {
				assert.Error(t, err)
				assert.Contains(t, test.expectedError, err.Error())
			}
		})
	}

	t.Run("The settings are accessible when loaded from the file", func(t *testing.T) {
		settings, _ := NewProcessSettingsFromFile("fixtures/combined_process_settings.yml", nil)
		assert.IsType(t, &[]SettingsFile{}, settings.Settings)
		assert.Equal(t, "telecom", (*settings.Settings)[1].Target["app"])
		assert.Equal(t, "caller_id_privacy", (*settings.Settings)[3].Settings["log_stream"].(map[string]interface{})["sip"])
		assert.Equal(t, "+12755554321", (*settings.Settings)[3].Target["caller_id"].([]interface{})[1])
	})
}
