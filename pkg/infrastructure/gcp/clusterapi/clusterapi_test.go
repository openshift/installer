package clusterapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProviderTimeouts(t *testing.T) {
	provider := Provider{}

	assert.Equal(t, 30*time.Minute, provider.NetworkTimeout())
	assert.Equal(t, 15*time.Minute, provider.ProvisionTimeout())
}
