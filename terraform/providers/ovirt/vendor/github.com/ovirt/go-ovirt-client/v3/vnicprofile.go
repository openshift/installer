package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "VnicProfile" -n "VNIC profile" -o "VNICProfile" -s "Profile" -T VNICProfileID

// VNICProfileID is the ID of the VNIC profile.
type VNICProfileID string

// VNICProfileClient defines the methods related to dealing with virtual NIC profiles.
type VNICProfileClient interface {
	// CreateVNICProfile creates a new VNIC profile with the specified name and network ID.
	CreateVNICProfile(name string, networkID NetworkID, params OptionalVNICProfileParameters, retries ...RetryStrategy) (VNICProfile, error)
	// GetVNICProfile returns a single VNIC Profile based on the ID
	GetVNICProfile(id VNICProfileID, retries ...RetryStrategy) (VNICProfile, error)
	// ListVNICProfiles lists all VNIC Profiles.
	ListVNICProfiles(retries ...RetryStrategy) ([]VNICProfile, error)
	// RemoveVNICProfile removes a VNIC profile
	RemoveVNICProfile(id VNICProfileID, retries ...RetryStrategy) error
}

// OptionalVNICProfileParameters is a set of parameters for creating VNICProfiles that are optional.
type OptionalVNICProfileParameters interface{}

// BuildableVNICProfileParameters is a buildable version of OptionalVNICProfileParameters.
type BuildableVNICProfileParameters interface {
	OptionalVNICProfileParameters
}

// CreateVNICProfileParams creats a buildable set of optional parameters for VNICProfile creation.
func CreateVNICProfileParams() BuildableVNICProfileParameters {
	return &vnicProfileParams{}
}

type vnicProfileParams struct{}

// VNICProfileData is the core of VNICProfile, providing only data access functions.
type VNICProfileData interface {
	// ID returns the identifier of the VNICProfile.
	ID() VNICProfileID
	// Name returns the human-readable name of the VNIC profile.
	Name() string
	// NetworkID returns the network ID the VNICProfile is attached to.
	NetworkID() NetworkID
}

// VNICProfile is a collection of settings that can be applied to individual virtual network interface cards in the
// Engine.
type VNICProfile interface {
	VNICProfileData

	// Network fetches the network object from the oVirt engine. This is an API call and may be slow.
	Network(retries ...RetryStrategy) (Network, error)
	// Remove removes the current VNIC profile.
	Remove(retries ...RetryStrategy) error
}

func convertSDKVNICProfile(sdkObject *ovirtsdk.VnicProfile, client Client) (VNICProfile, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("VNICProfile", "ID")
	}
	name, ok := sdkObject.Name()
	if !ok {
		return nil, newFieldNotFound("VNICProfile", "name")
	}
	network, ok := sdkObject.Network()
	if !ok {
		return nil, newFieldNotFound("VNICProfile", "network")
	}
	networkID, ok := network.Id()
	if !ok {
		return nil, newFieldNotFound("Network on VNICProfile", "ID")
	}

	return &vnicProfile{
		client: client,

		id:        VNICProfileID(id),
		name:      name,
		networkID: NetworkID(networkID),
	}, nil
}

type vnicProfile struct {
	client Client

	id        VNICProfileID
	networkID NetworkID
	name      string
}

func (v vnicProfile) Remove(retries ...RetryStrategy) error {
	return v.client.RemoveVNICProfile(v.id, retries...)
}

func (v vnicProfile) Network(retries ...RetryStrategy) (Network, error) {
	return v.client.GetNetwork(v.networkID, retries...)
}

func (v vnicProfile) Name() string {
	return v.name
}

func (v vnicProfile) NetworkID() NetworkID {
	return v.networkID
}

func (v vnicProfile) ID() VNICProfileID {
	return v.id
}
