package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//CloudProviderConfig is the azure cloud provider config
type CloudProviderConfig struct {
	TenantID       string
	SubscriptionID string
	GroupLocation  string
	ResourcePrefix string
}

// JSON generates the cloud provider json config for the azure platform.
func (params CloudProviderConfig) JSON() (string, error) {
	resourceGroupName := params.ResourcePrefix + "-rg"
	config := config{
		authConfig: authConfig{
			Cloud:                       "AzurePublicCloud",
			TenantID:                    params.TenantID,
			SubscriptionID:              params.SubscriptionID,
			UseManagedIdentityExtension: true,
			UserAssignedIdentityID: fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s",
				params.SubscriptionID,
				resourceGroupName,
				params.ResourcePrefix+"-identity"),
		},
		ResourceGroup:          resourceGroupName,
		Location:               params.GroupLocation,
		SubnetName:             params.ResourcePrefix + "-node-subnet",
		SecurityGroupName:      params.ResourcePrefix + "-nsg",
		VnetName:               params.ResourcePrefix + "-vnet",
		VnetResourceGroup:      resourceGroupName,
		RouteTableName:         params.ResourcePrefix + "-routetable",
		CloudProviderBackoff:   true,
		CloudProviderRateLimit: true,

		UseInstanceMetadata: true,
		//default to standard load balancer, supports tcp resets on idle
		//https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-tcp-reset
		LoadBalancerSku: "standard",
	}
	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(config); err != nil {
		return "", err
	}
	return buff.String(), nil
}
