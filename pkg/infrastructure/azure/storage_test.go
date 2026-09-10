package azure

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/stretchr/testify/assert"
)

func TestBlockBlobClientOptions(t *testing.T) {
	tests := []struct {
		name                 string
		allowSharedKeyAccess bool
		expectedStatusCodes  []int
	}{
		{
			name: "token credential",
			expectedStatusCodes: []int{
				http.StatusRequestTimeout,
				http.StatusTooManyRequests,
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
				http.StatusForbidden,
			},
		},
		{
			name:                 "shared key",
			allowSharedKeyAccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := blockBlobClientOptions(&arm.ClientOptions{}, tt.allowSharedKeyAccess)

			if tt.expectedStatusCodes == nil {
				assert.Nil(t, options.Retry.StatusCodes)
			} else {
				assert.Equal(t, tt.expectedStatusCodes, options.Retry.StatusCodes)
			}
			assert.Zero(t, options.Retry.MaxRetries)
			assert.Zero(t, options.Retry.TryTimeout)
			assert.Zero(t, options.Retry.RetryDelay)
			assert.Zero(t, options.Retry.MaxRetryDelay)
			assert.Nil(t, options.Retry.ShouldRetry)
		})
	}
}
