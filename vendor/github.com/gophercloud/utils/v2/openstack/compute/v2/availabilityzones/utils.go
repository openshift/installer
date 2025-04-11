package availabilityzones

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
)

// ListAvailableAvailabilityZones is a convenience function that return a slice of available Availability Zones.
func ListAvailableAvailabilityZones(ctx context.Context, client *gophercloud.ServiceClient) ([]string, error) {
	var ret []string
	allPages, err := availabilityzones.List(client).AllPages(ctx)
	if err != nil {
		return ret, err
	}

	availabilityZoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return ret, err
	}

	for _, zoneInfo := range availabilityZoneInfo {
		if zoneInfo.ZoneState.Available {
			ret = append(ret, zoneInfo.ZoneName)
		}
	}
	return ret, nil
}
