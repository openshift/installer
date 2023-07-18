package aws

const (
	// AvailabilityZoneType is the type of regular zone placed on the region.
	AvailabilityZoneType = "availability-zone"
	// LocalZoneType is the type of Local zone placed on the metropolitan areas.
	LocalZoneType = "local-zone"
)

// Zones stores the map of Zone attributes indexed by Zone Name.
type Zones map[string]*Zone

// Zone stores the Availability or Local Zone attributes used to set machine attributes, and to
// feed VPC resources as a source for for terraform variables.
type Zone struct {

	// Name is the availability, local or wavelength zone name.
	Name string `json:"name"`

	// ZoneType is the type of subnet's availability zone.
	// The valid values are availability-zone and local-zone.
	Type string `json:"type"`

	// ZoneGroupName is the AWS zone group name.
	// For Availability Zones, this parameter has the same value as the Region name.
	//
	// For Local Zones, the name of the associated group, for example us-west-2-lax-1.
	GroupName string `json:"groupName"`

	// ParentZoneName is the name of the zone that handles some of the Local Zone
	// control plane operations, such as API calls.
	ParentZoneName string `json:"parentZoneName"`

	// PreferredInstanceType is the offered instance type on the subnet's zone.
	// It's used for the edge pools which does not offer the same type across different zone groups.
	PreferredInstanceType string `json:"preferredInstanceType"`
}
