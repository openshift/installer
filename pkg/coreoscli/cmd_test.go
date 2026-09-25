package coreoscli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

func TestSelectOSImageStream(t *testing.T) {
	validStreams := types.OSImageStreamValues()
	tests := []struct {
		name          string
		streamFlag    string
		expected      types.OSImageStream
		expectedError string
	}{
		{
			name:          "invalid explicit stream",
			streamFlag:    "invalid",
			expectedError: fmt.Sprintf("invalid value %q for --stream; must be one of %v", "invalid", validStreams),
		},
		{
			name:     "default stream",
			expected: rhcos.BuildDefaultOSImageStream(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := selectOSImageStream(tt.streamFlag)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}

	for _, stream := range validStreams {
		t.Run(fmt.Sprintf("explicit stream %q", stream), func(t *testing.T) {
			actual, err := selectOSImageStream(string(stream))
			assert.NoError(t, err)
			assert.Equal(t, stream, actual)
		})
	}
}
