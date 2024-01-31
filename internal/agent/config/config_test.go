package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("agent new config test", func(t *testing.T) {
		if got := New(); got == nil {
			t.Error("config not init")
		}
	})
}
