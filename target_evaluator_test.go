package process_settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetEvaluation(t *testing.T) {
	tests := []struct {
		name             string
		settingsFile     SettingsFile
		targetingContext map[string]interface{}
		expectedResult   bool
	}{
		{
			name: "The target does not match when the settings file is metadata",
			settingsFile: SettingsFile{
				Metadata: SettingsMetadata{
					Version: 17,
					End:     true,
				},
			},
			targetingContext: map[string]interface{}{},
			expectedResult:   false,
		},
		{
			name: "The target matches when the settings file does not have a target",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
			},
			targetingContext: map[string]interface{}{},
			expectedResult:   true,
		},
		{
			name: "The target matches when the settings file has a target that matches the targeting context",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": "test",
				},
			},
			targetingContext: map[string]interface{}{
				"test": "test",
			},
			expectedResult: true,
		},
		{
			name: "The target does not match when the settings file has a target that does not match the targeting context",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": "test",
				},
			},
			targetingContext: map[string]interface{}{
				"test": "not test",
			},
			expectedResult: false,
		},
		{
			name: "The target matches when the settings file target uses an array and the targeting context value is in the array",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": []interface{}{"test", "not test"},
				},
			},
			targetingContext: map[string]interface{}{
				"test": "test",
			},
			expectedResult: true,
		},
		{
			name: "The target does not match when the settings file target uses an array and the targeting context value is not in the array",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": []interface{}{"not test"},
				},
			},
			targetingContext: map[string]interface{}{
				"test": "test",
			},
			expectedResult: false,
		},
		{
			name: "The target matches when the settings file target and static context are matching nested maps",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": map[string]interface{}{
						"test": "test",
					},
				},
			},
			targetingContext: map[string]interface{}{
				"test": map[string]interface{}{
					"test": "test",
				},
			},
			expectedResult: true,
		},
		{
			name: "The target does not match when the settings target expects a map and the targeting context value is not a map",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": map[string]interface{}{
						"test": "test",
					},
				},
			},
			targetingContext: map[string]interface{}{
				"test": "test",
			},
			expectedResult: false,
		},
		{
			name: "The target does not match when the targeting context is missing a key that the settings target expects",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": map[string]interface{}{
						"test": "test",
					},
				},
			},
			targetingContext: map[string]interface{}{
				"test": map[string]interface{}{},
			},
			expectedResult: false,
		},
		{
			name: "The target matches for very complex cases",
			settingsFile: SettingsFile{
				FileName: "test",
				Settings: map[string]interface{}{
					"test": "test",
				},
				Target: map[string]interface{}{
					"test": map[string]interface{}{
						"test": []interface{}{
							"test",
						},
					},
					"test2": map[string]interface{}{
						"test": 1,
					},
				},
			},
			targetingContext: map[string]interface{}{
				"test": map[string]interface{}{
					"test": "test",
				},
				"test2": map[string]interface{}{
					"test": 1,
				},
			},
			expectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := TargetEvaluator{test.targetingContext}
			assert.Equal(t, test.expectedResult, evaluator.IsTargetMatch(test.settingsFile))
		})
	}
}
