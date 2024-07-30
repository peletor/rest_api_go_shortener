package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 10",
			size: 10,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.size, len(NewRandomString(tt.size)))
	}
}
