package aws

const (
	// AvailabilityZoneType is the type of regular zone placed on the region.
	AvailabilityZoneType = "availability-zone"
	// LocalZoneType is the type of AWS Local Zones placed on the metropolitan area.
	LocalZoneType = "local-zone"
	// WavelengthZoneType is the type of AWS Wavelength Zones placed on the telecommunications
	// providers’ data centers at the edge of the 5G network.
	WavelengthZoneType = "wavelength-zone"
	// ZoneOptInStatusOptedIn is the opt-in status of the zone.
	// For Availability Zones, this parameter always has the value of opt-in-not-required.
	// For Local Zones and Wavelength Zones, this parameter is the opt-in status.
	ZoneOptInStatusOptedIn = "opted-in"
)
