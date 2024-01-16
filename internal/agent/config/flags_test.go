package config

import (
	"os"
	"testing"
)

func Test_parseFlags(t *testing.T) {
	t.Run("agent parseFlags test", func(t *testing.T) {
		os.Args = append(os.Args, "-a", "localhost:3000")
		if got := New(); got.Endpoint != "localhost:3000" {
			t.Errorf("parseFlags works incorrect, expected %s, got %s", "localhost:3000", got.Endpoint)
		}
	})
}
