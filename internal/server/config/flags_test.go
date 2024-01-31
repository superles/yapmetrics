package config

import (
	"flag"
	"os"
	"testing"
)

func Test_parseFlags(t *testing.T) {
	t.Run("server parseFlags test", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = os.Args[:1]
		os.Args = append(os.Args, "-a", "localhost:3000")
		if got := parseFlags(); got.Endpoint != "localhost:3000" {
			t.Errorf("parseFlags works incorrect, expected %s, got %s", "localhost:3000", got.Endpoint)
		}
	})
}
