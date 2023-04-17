package process_settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessSettingsSingleton(t *testing.T) {
	t.Run("Get returns an error with the singleton instance has not been set", func(t *testing.T) {
		SetGlobalProcessSettings(nil)
		_, err := Get("settings", "log_stream", "sip")

		assert.Equal(t, "The global process settings have not been set", err.Error())
	})

	t.Run("SafeGet returns nil with the singleton instance has not been set", func(t *testing.T) {
		SetGlobalProcessSettings(nil)
		value, _ := SafeGet("settings", "log_stream", "sip")

		assert.Nil(t, value)
	})
}
