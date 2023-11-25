package ibmcloud

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
)

// AvailabilityZones returns a list of supported zones for the specified region.
func AvailabilityZones(region string, serviceEndpoints []configv1.IBMCloudServiceEndpoint) ([]string, error) {
	ctx := context.TODO()

	client, err := ibmcloud.NewClient(serviceEndpoints)
	if err != nil {
		return nil, err
	}

	return client.GetVPCZonesForRegion(ctx, region)
}
