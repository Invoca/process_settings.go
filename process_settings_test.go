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
			fileName:      "testdata/does_not_exist.yml",
			expectedError: []string{"open testdata/does_not_exist.yml: no such file or directory"},
		},
		{
			name:          "The file is not a valid yaml",
			fileName:      "testdata/invalid_yaml.yml",
			expectedError: []string{"yaml: unmarshal errors:\n  line 5: cannot unmarshal !!map into []process_settings.SettingsFile"},
		},
		{
			name:          "The file doesn't have metadata as the last element of the array",
			fileName:      "testdata/invalid_metadata.yml",
			expectedError: []string{"The settings file does not have the END metadata"},
		},
		{
			name:          "The file has an invalid settings file",
			fileName:      "testdata/invalid_settings.yml",
			expectedError: []string{"Invalid settings file at index 0: The settings file must only have settings or metadata, not both => {honeypot.yml map[] map[honeypot:map[answer_odds:100 max_recording_seconds:600 status_change_min_days:7]] {17 true}}"},
		},
		{
			name:            "The file is valid",
			fileName:        "testdata/combined_process_settings.yml",
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
		settings, _ := NewProcessSettingsFromFile("testdata/combined_process_settings.yml", nil)
		assert.IsType(t, &[]SettingsFile{}, settings.Settings)
		assert.Equal(t, "telecom", (*settings.Settings)[1].Target["app"])
		assert.Equal(t, "caller_id_privacy", (*settings.Settings)[3].Settings["log_stream"].(map[string]interface{})["sip"])
		assert.Equal(t, "+12755554321", (*settings.Settings)[3].Target["caller_id"].([]interface{})[1])
	})
}

var (
	honeypotWithoutLogStream = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"answer_odds": 100,
				},
			},
		},
	}
	honeypotWithLogStream = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"answer_odds": 100,
					"log_stream":  "sip",
				},
			},
		},
	}
	honeypotWithTarget = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Target: map[string]interface{}{
				"app": "telecom",
			},
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"answer_odds": 100,
					"log_stream":  "sip",
				},
			},
		},
	}
	honeypotWithTargetedOverride = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"log_stream": "original",
				},
			},
		},
		{
			FileName: "honeypot_override.yml",
			Target: map[string]interface{}{
				"app": "telecom",
			},
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"answer_odds": 100,
					"log_stream":  "override",
				},
			},
		},
	}
	honeypotWithSettingsArray = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"certs": []interface{}{
						map[string]interface{}{
							"path": "original_0",
						},
						map[string]interface{}{
							"path": "original_1",
						},
					},
				},
			},
		},
	}
	complexHoneypotWithSettingsOnlyInTarget = &[]SettingsFile{
		{
			FileName: "honeypot.yml",
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"certs": []interface{}{
						map[string]interface{}{
							"path": "original_0",
						},
						map[string]interface{}{
							"path": "original_1",
						},
					},
				},
			},
		},
		{
			FileName: "honeypot_override.yml",
			Target: map[string]interface{}{
				"app": "telecom",
			},
			Settings: map[string]interface{}{
				"honeypot": map[string]interface{}{
					"log_stream": map[string]interface{}{
						"telecom": "something",
					},
				},
			},
		},
	}
)

var getAndSafeGetTests = []struct {
	name            string
	processSettings ProcessSettings
	settingPath     []string
	expectedError   string
	expectedValue   interface{}
}{
	{
		name: "Returns an error when the setting is not found",
		processSettings: ProcessSettings{
			Settings: honeypotWithoutLogStream,
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedError: "The setting 'honeypot.log_stream' was not found",
	},
	{
		name: "Returns the value when the setting is found",
		processSettings: ProcessSettings{
			Settings: honeypotWithLogStream,
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedValue: "sip",
	},
	{
		name: "Does not find the setting when the targeting does not match",
		processSettings: ProcessSettings{
			Settings: honeypotWithTarget,
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedError: "The setting 'honeypot.log_stream' was not found",
	},
	{
		name: "Finds the setting when the targeting does not match",
		processSettings: ProcessSettings{
			Settings: honeypotWithTarget,
			TargetEvaluator: TargetEvaluator{
				targetingContext: map[string]interface{}{
					"app": "telecom",
				},
			},
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedValue: "sip",
	},
	{
		name: "Ignores overridden settings when the targeting does not match",
		processSettings: ProcessSettings{
			Settings: honeypotWithTargetedOverride,
			TargetEvaluator: TargetEvaluator{
				targetingContext: map[string]interface{}{
					"app": "not telecom",
				},
			},
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedValue: "original",
	},
	{
		name: "Returns the overridden settings when the targeting matches",
		processSettings: ProcessSettings{
			Settings: honeypotWithTargetedOverride,
			TargetEvaluator: TargetEvaluator{
				targetingContext: map[string]interface{}{
					"app": "telecom",
				},
			},
		},
		settingPath:   []string{"honeypot", "log_stream"},
		expectedValue: "override",
	},
	{
		name: "Returns nil when the nested setting doesn't exist due to targeting",
		processSettings: ProcessSettings{
			Settings: complexHoneypotWithSettingsOnlyInTarget,
		},
		settingPath:   []string{"honeypot", "log_stream", "telecom"},
		expectedError: "The setting 'honeypot.log_stream.telecom' was not found",
	},
	{
		name: "Returns the setting value when the nested setting exists due to targeting",
		processSettings: ProcessSettings{
			Settings: complexHoneypotWithSettingsOnlyInTarget,
			TargetEvaluator: TargetEvaluator{
				targetingContext: map[string]interface{}{
					"app": "telecom",
				},
			},
		},
		settingPath:   []string{"honeypot", "log_stream", "telecom"},
		expectedValue: "something",
	},
}

func TestProcessSettings_SafeGet(t *testing.T) {
	for _, test := range getAndSafeGetTests {
		t.Run(test.name, func(t *testing.T) {
			value := test.processSettings.SafeGet(test.settingPath...)
			if test.expectedError == "" {
				assert.Equal(t, test.expectedValue, value)
			} else {
				assert.Nil(t, value)
			}
		})
	}
}

func TestProcessSettings_Get(t *testing.T) {
	for _, test := range getAndSafeGetTests {
		t.Run(test.name, func(t *testing.T) {
			value, err := test.processSettings.Get(test.settingPath...)
			if test.expectedError == "" {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedValue, value)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			}
		})
	}
}
