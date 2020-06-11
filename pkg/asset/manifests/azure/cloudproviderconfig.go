package azure

import (
	"bytes"
	"encoding/json"
)

//CloudProviderConfig is the azure cloud provider config
type CloudProviderConfig struct {
	TenantID                 string
	SubscriptionID           string
	GroupLocation            string
	ResourcePrefix           string
	NetworkResourceGroupName string
	NetworkSecurityGroupName string
	VirtualNetworkName       string
	SubnetName               string
}

// JSON generates the cloud provider json config for the azure platform.
// managed resource names are matching the convention defined by capz
func (params CloudProviderConfig) JSON() (string, error) {
	resourceGroupName := params.ResourcePrefix + "-rg"
	config := config{
		authConfig: authConfig{
			Cloud:                       "AzurePublicCloud",
			TenantID:                    params.TenantID,
			SubscriptionID:              params.SubscriptionID,
			UseManagedIdentityExtension: true,
			// The cloud provider needs the clientID which is only known after terraform has run.
			// When left empty, the existing managed identity on the VM will be used.
			// By leaving it empty, we don't have to create the identity before running the installer.
			// We only need to know that there will be one assigned to the VM, and we control this.
			// ref: https://github.com/kubernetes/kubernetes/blob/4b7c607ba47928a7be77fadef1550d6498397a4c/staging/src/k8s.io/legacy-cloud-providers/azure/auth/azure_auth.go#L69
			UserAssignedIdentityID: "",
		},
		ResourceGroup:     resourceGroupName,
		Location:          params.GroupLocation,
		SubnetName:        params.SubnetName,
		SecurityGroupName: params.NetworkSecurityGroupName,
		VnetName:          params.VirtualNetworkName,
		VnetResourceGroup: params.NetworkResourceGroupName,
		RouteTableName:    params.ResourcePrefix + "-node-routetable",
		// client side rate limiting is problematic for scaling operations. We disable it by default.
		// https://github.com/kubernetes-sigs/cloud-provider-azure/issues/247
		// https://bugzilla.redhat.com/show_bug.cgi?id=1782516#c7
		CloudProviderRateLimit:       false,
		CloudProviderBackoff:         true,
		CloudProviderBackoffDuration: 6,

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
