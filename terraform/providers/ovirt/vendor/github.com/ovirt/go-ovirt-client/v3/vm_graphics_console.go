package ovirtclient

import ovirtsdk "github.com/ovirt/go-ovirt"

// VMGraphicsConsoleID is the identifier for graphics consoles on a VM.
type VMGraphicsConsoleID string

// GraphicsConsoleClient lists the methods to access and manipulate graphics consoles on VMs.
type GraphicsConsoleClient interface {
	ListVMGraphicsConsoles(vmID VMID, retries ...RetryStrategy) ([]VMGraphicsConsole, error)
	RemoveVMGraphicsConsole(vmID VMID, graphicsConsoleID VMGraphicsConsoleID, retries ...RetryStrategy) error
}

// VMGraphicsConsoleData contains the data for VMGraphicsConsole objects.
type VMGraphicsConsoleData interface {
	ID() VMGraphicsConsoleID
	VMID() VMID
}

// VMGraphicsConsole is an object representing a graphics console on a virtual machine.
type VMGraphicsConsole interface {
	VMGraphicsConsoleData

	// Remove removes the graphics console.
	Remove(retries ...RetryStrategy) error
}

type vmGraphicsConsole struct {
	client Client

	id   VMGraphicsConsoleID
	vmID VMID
}

func (v *vmGraphicsConsole) Remove(retries ...RetryStrategy) error {
	return v.client.RemoveVMGraphicsConsole(v.vmID, v.id, retries...)
}

func (v *vmGraphicsConsole) ID() VMGraphicsConsoleID {
	return v.id
}

func (v *vmGraphicsConsole) VMID() VMID {
	return v.vmID
}

func convertSDKGraphicsConsole(sdkObject *ovirtsdk.GraphicsConsole, client Client) (VMGraphicsConsole, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("graphics console", "id")
	}
	vm, ok := sdkObject.Vm()
	if !ok {
		return nil, newFieldNotFound("graphics console", "vm")
	}
	vmID, ok := vm.Id()
	if !ok {
		return nil, newFieldNotFound("vm on graphics console", "id")
	}

	return &vmGraphicsConsole{
		client,
		VMGraphicsConsoleID(id),
		VMID(vmID),
	}, nil
}
