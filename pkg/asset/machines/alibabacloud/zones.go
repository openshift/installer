package alibabacloud

import (
	// "context"

	"github.com/openshift/installer/pkg/asset/installconfig/alibabacloud"
)

// GetAvailabilityZones returns a list of supported zones for the specified region.
func GetAvailabilityZones(region string) ([]string, error) {
	// ctx := context.TODO()
	client, err := alibabacloud.NewClient(region)
	if err != nil {
		return nil, err
	}

	response, err := client.DescribeAvailableResource("Zone")
	if err != nil {
		return nil, err
	}

	var zones []string
	for _, zone := range response.AvailableZones.AvailableZone {
		if zone.Status == "Available" {
			zones = append(zones, zone.ZoneId)
		}
	}
	return zones, nil
}
