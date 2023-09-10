package agent

import (
	"testing"
)

func Test_sendAll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive test #1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//err := sendAll()
			//assert.Error(t, err)
		})
	}
}
