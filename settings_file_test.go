package process_settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingsFileValidation(t *testing.T) {
	tests := []struct {
		name          string
		settingsFile  SettingsFile
		expectedValid bool
		expectedError string
	}{
		{
			name: "The settings file is valid when just metadata is present",
			settingsFile: SettingsFile{
				Metadata: SettingsMetadata{
					Version: 17,
					End:     true,
				},
			},
			expectedValid: true,
		},
		{
			name: "The settings files is valid when just the filename and settings are present",
			settingsFile: SettingsFile{
				FileName: "honeypot.yml",
				Settings: map[string]interface{}{
					"honeypot": map[string]interface{}{
						"answer_odds": 100,
					},
				},
			},
			expectedValid: true,
		},
		{
			name: "The settings file is invalid when the filename is empty",
			settingsFile: SettingsFile{
				FileName: "",
				Settings: map[string]interface{}{
					"honeypot": map[string]interface{}{
						"answer_odds": 100,
					},
				},
			},
			expectedValid: false,
			expectedError: "The settings file must have a filename and settings",
		},
		{
			name: "The settings file is invalid when the filename is missing",
			settingsFile: SettingsFile{
				Settings: map[string]interface{}{
					"honeypot": map[string]interface{}{
						"answer_odds": 100,
					},
				},
			},
			expectedValid: false,
			expectedError: "The settings file must have a filename and settings",
		},
		{
			name: "The settings file is invalid when the settings are missing",
			settingsFile: SettingsFile{
				FileName: "honeypot.yml",
			},
			expectedValid: false,
			expectedError: "The settings file must have a filename and settings",
		},
		{
			name: "The settings file is invalid when it has settings and metadata",
			settingsFile: SettingsFile{
				FileName: "honeypot.yml",
				Settings: map[string]interface{}{
					"honeypot": map[string]interface{}{
						"answer_odds": 100,
					},
				},
				Metadata: SettingsMetadata{
					Version: 17,
					End:     true,
				},
			},
			expectedValid: false,
			expectedError: "The settings file must only have settings or metadata, not both",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid, err := test.settingsFile.isValid()
			assert.Equal(t, test.expectedValid, valid)
			if !valid {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), test.expectedError)
			}
		})
	}
}
