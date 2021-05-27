package ibmcloud

import (
	"context"

	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
)

// AvailabilityZones returns a list of supported zones for the specified region.
func AvailabilityZones(region string) ([]string, error) {
	ctx := context.TODO()

	client, err := ibmcloud.NewClient()
	if err != nil {
		return nil, err
	}

	return client.GetVPCZonesForRegion(ctx, region)
}
