package ovirt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	expectedConfig := `{
	"storageDomainId": "dd3ec3e5-e38b-4f02-9947-c669368cde56",
	"clusterId": "dd3ec3e5-e38b-4f02-9947-c669368cde57",
	"networkName": "production"
}
`
	actualConfig, err := CloudProviderConfig("dd3ec3e5-e38b-4f02-9947-c669368cde56", "dd3ec3e5-e38b-4f02-9947-c669368cde57", "production")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
