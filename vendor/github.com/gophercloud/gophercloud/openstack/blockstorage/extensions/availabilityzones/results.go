package availabilityzones

import (
	"github.com/gophercloud/gophercloud/pagination"
)

// ZoneState represents the current state of the availability zone.
type ZoneState struct {
	// Returns true if the availability zone is available
	Available bool `json:"available"`
}

// AvailabilityZone contains all the information associated with an OpenStack
// AvailabilityZone.
type AvailabilityZone struct {
	// The availability zone name
	ZoneName  string    `json:"zoneName"`
	ZoneState ZoneState `json:"zoneState"`
}

type AvailabilityZonePage struct {
	pagination.SinglePageBase
}

// ExtractAvailabilityZones returns a slice of AvailabilityZones contained in a
// single page of results.
func ExtractAvailabilityZones(r pagination.Page) ([]AvailabilityZone, error) {
	var s struct {
		AvailabilityZoneInfo []AvailabilityZone `json:"availabilityZoneInfo"`
	}
	err := (r.(AvailabilityZonePage)).ExtractInto(&s)
	return s.AvailabilityZoneInfo, err
}
