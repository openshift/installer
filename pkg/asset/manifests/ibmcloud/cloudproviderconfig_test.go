package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
)

func TestCloudProviderConfig(t *testing.T) {
	accountID := "1e1f75646aef447814a6d907cc83fb3c"
	existingSubnets := []string{
		"existing-subnet-control-plane-eu-gb-1",
		"existing-subnet-control-plane-eu-gb-2",
		"existing-subnet-control-plane-eu-gb-3",
		"existing-subnet-compute-eu-gb-1",
		"existing-subnet-compute-eu-gb-2",
		"existing-subnet-compute-eu-gb-3",
	}
	defaultConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = ""
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-8pxks
cluster-default-provider = g2
region = us-east
g2Credentials = /etc/vpc/ibmcloud_api_key
g2ResourceGroupName = ocp4-8pxks-rg
g2VpcName = ocp4-8pxks-vpc
g2workerServiceAccountID = 1e1f75646aef447814a6d907cc83fb3c
g2VpcSubnetNames = ocp4-8pxks-subnet-compute-us-east-1,ocp4-8pxks-subnet-compute-us-east-2,ocp4-8pxks-subnet-compute-us-east-3,ocp4-8pxks-subnet-control-plane-us-east-1,ocp4-8pxks-subnet-control-plane-us-east-2,ocp4-8pxks-subnet-control-plane-us-east-3


`

	existingSubnetConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = ""
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-hf4vtt
cluster-default-provider = g2
region = eu-gb
g2Credentials = /etc/vpc/ibmcloud_api_key
g2ResourceGroupName = ocp4-hf4vtt-rg
g2VpcName = ocp4-hf4vtt-vpc
g2workerServiceAccountID = 1e1f75646aef447814a6d907cc83fb3c
g2VpcSubnetNames = existing-subnet-control-plane-eu-gb-1,existing-subnet-control-plane-eu-gb-2,existing-subnet-control-plane-eu-gb-3,existing-subnet-compute-eu-gb-1,existing-subnet-compute-eu-gb-2,existing-subnet-compute-eu-gb-3


`

	singleEndpointOverrideConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = ""
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-ghs4s3
cluster-default-provider = g2
region = eu-gb
g2Credentials = /etc/vpc/ibmcloud_api_key
g2ResourceGroupName = ocp4-ghs4s3-rg
g2VpcName = ocp4-ghs4s3-vpc
g2workerServiceAccountID = 1e1f75646aef447814a6d907cc83fb3c
g2VpcSubnetNames = existing-subnet-control-plane-eu-gb-1,existing-subnet-control-plane-eu-gb-2,existing-subnet-control-plane-eu-gb-3,existing-subnet-compute-eu-gb-1,existing-subnet-compute-eu-gb-2,existing-subnet-compute-eu-gb-3
g2EndpointOverride = https://ibmcloud.vpc.override.endpoint.test


`

	multiEndpointOverrideConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = ""
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-ppcj33
cluster-default-provider = g2
region = eu-gb
g2Credentials = /etc/vpc/ibmcloud_api_key
g2ResourceGroupName = ocp4-ppcj33-rg
g2VpcName = ocp4-ppcj33-vpc
g2workerServiceAccountID = 1e1f75646aef447814a6d907cc83fb3c
g2VpcSubnetNames = existing-subnet-control-plane-eu-gb-1,existing-subnet-control-plane-eu-gb-2,existing-subnet-control-plane-eu-gb-3,existing-subnet-compute-eu-gb-1,existing-subnet-compute-eu-gb-2,existing-subnet-compute-eu-gb-3
iamEndpointOverride = https://ibmcloud.iam.override.endpoint.test
g2EndpointOverride = https://ibmcloud.vpc.override.endpoint.test
rmEndpointOverride = https://ibmcloud.resource-manager.override.endpoint.test


`

	eugbZones := []string{"eu-gb-1", "eu-gb-2", "eu-gb-3"}
	useastZones := []string{"us-east-1", "us-east-2", "us-east-3"}

	cases := []struct {
		name              string
		infraID           string
		accountID         string
		region            string
		resourceGroupName string
		vpcName           string
		subnets           []string
		cpZones           []string
		computeZones      []string
		serviceEndpoints  []configv1.IBMCloudServiceEndpoint
		expectedConfig    string
	}{
		{
			name:              "default subnet config",
			infraID:           "ocp4-8pxks",
			accountID:         accountID,
			region:            "us-east",
			resourceGroupName: "ocp4-8pxks-rg",
			vpcName:           "ocp4-8pxks-vpc",
			subnets:           []string{},
			cpZones:           useastZones,
			computeZones:      useastZones,
			expectedConfig:    defaultConfig,
		},
		{
			name:              "existing subnet config",
			infraID:           "ocp4-hf4vtt",
			accountID:         accountID,
			region:            "eu-gb",
			resourceGroupName: "ocp4-hf4vtt-rg",
			vpcName:           "ocp4-hf4vtt-vpc",
			subnets:           existingSubnets,
			cpZones:           eugbZones,
			computeZones:      eugbZones,
			expectedConfig:    existingSubnetConfig,
		},
		{
			name:              "single endpoint override config",
			infraID:           "ocp4-ghs4s3",
			accountID:         accountID,
			region:            "eu-gb",
			resourceGroupName: "ocp4-ghs4s3-rg",
			vpcName:           "ocp4-ghs4s3-vpc",
			subnets:           existingSubnets,
			cpZones:           eugbZones,
			computeZones:      eugbZones,
			serviceEndpoints: []configv1.IBMCloudServiceEndpoint{
				{
					Name: configv1.IBMCloudServiceVPC,
					URL:  "https://ibmcloud.vpc.override.endpoint.test",
				},
			},
			expectedConfig: singleEndpointOverrideConfig,
		},
		{
			name:              "multiple endpoint override config",
			infraID:           "ocp4-ppcj33",
			accountID:         accountID,
			region:            "eu-gb",
			resourceGroupName: "ocp4-ppcj33-rg",
			vpcName:           "ocp4-ppcj33-vpc",
			subnets:           existingSubnets,
			cpZones:           eugbZones,
			computeZones:      eugbZones,
			serviceEndpoints: []configv1.IBMCloudServiceEndpoint{
				{
					Name: configv1.IBMCloudServiceCOS,
					URL:  "https://ibmcloud.cos.override.endpoint.test",
				},
				{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://ibmcloud.iam.override.endpoint.test",
				},
				{
					Name: configv1.IBMCloudServiceVPC,
					URL:  "https://ibmcloud.vpc.override.endpoint.test",
				},
				{
					Name: configv1.IBMCloudServiceResourceManager,
					URL:  "https://ibmcloud.resource-manager.override.endpoint.test",
				},
			},
			expectedConfig: multiEndpointOverrideConfig,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualConfig, err := CloudProviderConfig(tc.infraID, tc.accountID, tc.region, tc.resourceGroupName, tc.vpcName, tc.subnets, tc.cpZones, tc.computeZones, tc.serviceEndpoints)
			assert.NoError(t, err, "failed to create cloud provider config")
			assert.Equal(t, tc.expectedConfig, actualConfig, "unexpected cloud provider config")
		})
	}
}
