package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	expectedConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = ""
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-8pxks
cluster-default-provider = g2
region = us-east
g2Credentials = /etc/vpc/ibmcloud_api_key
g2VpcName = ocp4-8pxks-vpc
g2workerServiceAccountID = 1e1f75646aef447814a6d907cc83fb3c
g2VpcSubnetNames = ocp4-8pxks-subnet-compute-us-east-1,ocp4-8pxks-subnet-compute-us-east-2,ocp4-8pxks-subnet-compute-us-east-3,ocp4-8pxks-subnet-control-plane-us-east-1,ocp4-8pxks-subnet-control-plane-us-east-2,ocp4-8pxks-subnet-control-plane-us-east-3

`
	zones := []string{"us-east-1", "us-east-2", "us-east-3"}
	actualConfig, err := CloudProviderConfig("ocp4-8pxks", "1e1f75646aef447814a6d907cc83fb3c", "us-east", zones, zones)
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
