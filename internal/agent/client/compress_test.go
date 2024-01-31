package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompress(t *testing.T) {
	t.Run("compress test", func(t *testing.T) {
		_, err := Compress([]byte("test string"))
		require.NoError(t, err, "data compress error")
	})
}
