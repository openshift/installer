package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	expectedConfig := `[global]
project-id = test-project-id
regional   = true
multizone  = true
node-tags  = uid-master,uid-worker

`
	actualConfig, err := CloudProviderConfig("uid", "test-project-id")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
