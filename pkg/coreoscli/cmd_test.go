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
		name                   string
		streamFlag             string
		releaseVersionInjected bool
		expected               types.OSImageStream
		expectedError          string
	}{
		{
			name:                   "uninjected release version without stream",
			releaseVersionInjected: false,
			expectedError:          fmt.Sprintf("release version metadata was not injected into the installer; specify --stream with one of %v", validStreams),
		},
		{
			name:                   "uninjected release version with explicit stream",
			streamFlag:             string(validStreams[0]),
			releaseVersionInjected: false,
			expected:               validStreams[0],
		},
		{
			name:                   "invalid explicit stream",
			streamFlag:             "invalid",
			releaseVersionInjected: false,
			expectedError:          fmt.Sprintf("invalid value %q for --stream; must be one of %v", "invalid", validStreams),
		},
		{
			name:                   "injected release version without stream",
			releaseVersionInjected: true,
			expected:               rhcos.BuildDefaultOSImageStream(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := selectOSImageStream(tt.streamFlag, tt.releaseVersionInjected)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
