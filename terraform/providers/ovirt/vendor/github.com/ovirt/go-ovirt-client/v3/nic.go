package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// NICID is the ID for a network interface.
type NICID string

// NICClient defines the methods related to dealing with network interfaces.
type NICClient interface {
	// CreateNIC adds a new NIC to a VM specified in vmid.
	CreateNIC(
		vmid VMID,
		vnicProfileID VNICProfileID,
		name string,
		optional OptionalNICParameters,
		retries ...RetryStrategy,
	) (NIC, error)
	// UpdateNIC allows updating the NIC.
	UpdateNIC(
		vmid VMID,
		nicID NICID,
		params UpdateNICParameters,
		retries ...RetryStrategy,
	) (NIC, error)
	// GetNIC returns one specific NIC with the ID specified in id, attached to a VM with the ID specified in vmid.
	GetNIC(vmid VMID, id NICID, retries ...RetryStrategy) (NIC, error)
	// ListNICs lists all NICs attached to the VM specified in vmid.
	ListNICs(vmid VMID, retries ...RetryStrategy) ([]NIC, error)
	// RemoveNIC removes the network interface specified.
	RemoveNIC(vmid VMID, id NICID, retries ...RetryStrategy) error
}

// OptionalNICParameters is an interface that declares the source of optional parameters for NIC creation.
type OptionalNICParameters interface {
	// represent mac_address for NIC
	Mac() string
}

// BuildableNICParameters is a modifiable version of OptionalNICParameters. You can use CreateNICParams() to create a
// new copy, or implement your own.
type BuildableNICParameters interface {
	OptionalNICParameters

	// WithMac sets macAddress for the NIC.
	WithMac(mac string) (BuildableNICParameters, error)

	// MustWithMac is the same as WithMac, but panics instead of returning an error.
	MustWithMac(mac string) BuildableNICParameters
}

// CreateNICParams returns a buildable structure of OptionalNICParameters.
func CreateNICParams() BuildableNICParameters {
	return &nicParams{}
}

type nicParams struct {
	mac string
}

func (c *nicParams) Mac() string {
	return c.mac
}

func (c *nicParams) WithMac(mac string) (BuildableNICParameters, error) {
	c.mac = mac
	return c, nil
}

func (c *nicParams) MustWithMac(mac string) BuildableNICParameters {
	builder, err := c.WithMac(mac)
	if err != nil {
		panic(err)
	}
	return builder
}

// UpdateNICParameters is an interface that declares methods of changeable parameters for NIC's. Each
// method can return nil to leave an attribute unchanged, or a new value for the attribute.
type UpdateNICParameters interface {
	// Name potentially returns a changed name for a NIC.
	Name() *string

	// VNICProfileID potentially returns a change VNIC profile for a NIC.
	VNICProfileID() *VNICProfileID

	// Mac potentially returns a change MacAddress for a nic
	Mac() *string
}

// BuildableUpdateNICParameters is a buildable version of UpdateNICParameters.
type BuildableUpdateNICParameters interface {
	UpdateNICParameters

	// WithName sets the name of a NIC for the UpdateNIC method.
	WithName(name string) (BuildableUpdateNICParameters, error)
	// MustWithName is identical to WithName, but panics instead of returning an error.
	MustWithName(name string) BuildableUpdateNICParameters

	// WithVNICProfileID sets the VNIC profile ID of a NIC for the UpdateNIC method.
	WithVNICProfileID(id VNICProfileID) (BuildableUpdateNICParameters, error)
	// MustWithVNICProfileID is identical to WithVNICProfileID, but panics instead of returning an error.
	MustWithVNICProfileID(id VNICProfileID) BuildableUpdateNICParameters

	// WithMac sets MaAddress of a NIC for the UpdateNIC method.
	WithMac(mac string) (BuildableUpdateNICParameters, error)
	// MustWithMac is identical to WithMac, but panics instead of returning an error.
	MustWithMac(mac string) BuildableUpdateNICParameters
}

// UpdateNICParams creates a buildable UpdateNICParameters.
func UpdateNICParams() BuildableUpdateNICParameters {
	return &updateNICParams{}
}

type updateNICParams struct {
	name          *string
	vnicProfileID *VNICProfileID
	mac           *string
}

func (u *updateNICParams) Name() *string {
	return u.name
}

func (u *updateNICParams) VNICProfileID() *VNICProfileID {
	return u.vnicProfileID
}

func (u *updateNICParams) Mac() *string {
	return u.mac
}

func (u *updateNICParams) WithName(name string) (BuildableUpdateNICParameters, error) {
	u.name = &name
	return u, nil
}

func (u *updateNICParams) MustWithName(name string) BuildableUpdateNICParameters {
	b, err := u.WithName(name)
	if err != nil {
		panic(err)
	}
	return b
}

func (u *updateNICParams) WithVNICProfileID(id VNICProfileID) (BuildableUpdateNICParameters, error) {
	u.vnicProfileID = &id
	return u, nil
}

func (u *updateNICParams) MustWithVNICProfileID(id VNICProfileID) BuildableUpdateNICParameters {
	b, err := u.WithVNICProfileID(id)
	if err != nil {
		panic(err)
	}
	return b
}

func (u *updateNICParams) WithMac(mac string) (BuildableUpdateNICParameters, error) {
	u.mac = &mac
	return u, nil
}

func (u *updateNICParams) MustWithMac(mac string) BuildableUpdateNICParameters {
	b, err := u.WithMac(mac)
	if err != nil {
		panic(err)
	}
	return b
}

// NICData is the core of NIC which only provides data-access functions.
type NICData interface {
	// ID is the identifier for this network interface.
	ID() NICID
	// Name is the user-given name of the network interface.
	Name() string
	// VMID is the identified of the VM this NIC is attached to. May be nil if the NIC is not attached.
	VMID() VMID
	// VNICProfileID returns the ID of the VNIC profile in use by the NIC.
	VNICProfileID() VNICProfileID
	// Mac returns a MacAddress for a nic
	Mac() string
}

// NIC represents a network interface.
type NIC interface {
	NICData

	// GetVM fetches an up to date copy of the virtual machine this NIC is attached to. This involves an API call and
	// may be slow.
	GetVM(retries ...RetryStrategy) (VM, error)
	// GetVNICProfile retrieves the VNIC profile associated with this NIC. This involves an API call and may be slow.
	GetVNICProfile(retries ...RetryStrategy) (VNICProfile, error)
	// Update updates the NIC with the specified parameters. It returns the updated NIC as a response. You can use
	// UpdateNICParams() to obtain a buildable parameter structure.
	Update(params UpdateNICParameters, retries ...RetryStrategy) (NIC, error)
	// Remove removes the current network interface. This involves an API call and may be slow.
	Remove(retries ...RetryStrategy) error
}

func convertSDKNIC(sdkObject *ovirtsdk.Nic, cli Client) (NIC, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("id", "NIC")
	}
	name, ok := sdkObject.Name()
	if !ok {
		return nil, newFieldNotFound("name", "NIC")
	}
	vm, ok := sdkObject.Vm()
	if !ok {
		return nil, newFieldNotFound("vm", "NIC")
	}
	vmid, ok := vm.Id()
	if !ok {
		return nil, newFieldNotFound("VM in NIC", "ID")
	}
	vnicProfile, ok := sdkObject.VnicProfile()
	if !ok {
		return nil, newFieldNotFound("VM", "vNIC Profile")
	}
	vnicProfileID, ok := vnicProfile.Id()
	if !ok {
		return nil, newFieldNotFound("vNIC Profile on VM", "ID")
	}
	mac, ok := sdkObject.Mac()
	if !ok {
		return nil, newFieldNotFound("mac", "NIC")
	}
	macAddr, ok := mac.Address()
	if !ok {
		return nil, newFieldNotFound("address", "mac")
	}
	return &nic{
		cli,
		NICID(id),
		name,
		VMID(vmid),
		VNICProfileID(vnicProfileID),
		macAddr,
	}, nil
}

type nic struct {
	client Client

	id            NICID
	name          string
	vmid          VMID
	vnicProfileID VNICProfileID
	mac           string
}

func (n nic) Update(params UpdateNICParameters, retries ...RetryStrategy) (NIC, error) {
	return n.client.UpdateNIC(n.vmid, n.id, params, retries...)
}

func (n nic) GetVM(retries ...RetryStrategy) (VM, error) {
	return n.client.GetVM(n.vmid, retries...)
}

func (n nic) GetVNICProfile(retries ...RetryStrategy) (VNICProfile, error) {
	return n.client.GetVNICProfile(n.vnicProfileID, retries...)
}

func (n nic) VNICProfileID() VNICProfileID {
	return n.vnicProfileID
}

func (n nic) ID() NICID {
	return n.id
}

func (n nic) Name() string {
	return n.name
}

func (n nic) VMID() VMID {
	return n.vmid
}

func (n nic) Mac() string {
	return n.mac
}

func (n nic) Remove(retries ...RetryStrategy) error {
	return n.client.RemoveNIC(n.vmid, n.id, retries...)
}

func (n nic) withName(name string) *nic {
	return &nic{
		client:        n.client,
		id:            n.id,
		name:          name,
		vmid:          n.vmid,
		vnicProfileID: n.vnicProfileID,
		mac:           n.mac,
	}
}

func (n nic) withVNICProfileID(vnicProfileID VNICProfileID) *nic {
	return &nic{
		client:        n.client,
		id:            n.id,
		name:          n.name,
		vmid:          n.vmid,
		vnicProfileID: vnicProfileID,
		mac:           n.mac,
	}
}

func (n nic) withMac(mac string) *nic {
	return &nic{
		client:        n.client,
		id:            n.id,
		name:          n.name,
		vmid:          n.vmid,
		vnicProfileID: n.vnicProfileID,
		mac:           mac,
	}
}
