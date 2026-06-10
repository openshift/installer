package types

// MachineInfo holds the quota-relevant fields for a machine,
// abstracting over both MAPI and CAPI representations.
type MachineInfo struct {
	InstanceType     string
	AvailabilityZone string
	Replicas         int64
}
