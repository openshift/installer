package envbrowse

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// EnvironmentBrowser is a higher-level interface to a specific object's
// environment browser.
//
// This essentially fills the role of such functionality lacking in govmomi at
// this point in time and may serve as the basis for a respective PR at a later
// point in time.
type EnvironmentBrowser struct {
	object.Common
}

// NewEnvironmentBrowser initializes a new EnvironmentBrowser based off the
// supplied managed object reference.
func NewEnvironmentBrowser(c *vim25.Client, ref types.ManagedObjectReference) *EnvironmentBrowser {
	return &EnvironmentBrowser{
		Common: object.NewCommon(c, ref),
	}
}

// DefaultDevices loads a satisfactory default device list for the optionally
// supplied host and descriptor key. The result is returned as a higher-level
// VirtualDeviceList object. This can be used as an initial VirtualDeviceList
// when building a device list and VirtualDeviceConfigSpec list for new virtual
// machines.
//
// Appropriate options for key can be loaded by running
// QueryConfigOptionDescriptor, which will return a list of
// VirtualMachineConfigOptionDescriptor which will contain the appropriate key
// for the virtual machine version needed. If no key is supplied, the results
// generally reflect the most recent VM hardware version.
func (b *EnvironmentBrowser) DefaultDevices(ctx context.Context, key string, host *object.HostSystem) (object.VirtualDeviceList, error) {
	var eb mo.EnvironmentBrowser

	err := b.Properties(ctx, b.Reference(), nil, &eb)
	if err != nil {
		return nil, err
	}

	req := types.QueryConfigOption{
		This: b.Reference(),
		Key:  key,
	}
	if host != nil {
		ref := host.Reference()
		req.Host = &ref
	}
	res, err := methods.QueryConfigOption(ctx, b.Client(), &req)
	if err != nil {
		return nil, err
	}
	if res.Returnval == nil {
		return nil, errors.New("no config options were found for the supplied criteria")
	}
	return object.VirtualDeviceList(res.Returnval.DefaultDevice), nil
}

// OSFamily fetches the operating system family for the supplied guest ID.
func (b *EnvironmentBrowser) OSFamily(ctx context.Context, guest string) (string, error) {
	var eb mo.EnvironmentBrowser

	err := b.Properties(ctx, b.Reference(), nil, &eb)
	if err != nil {
		return "", err
	}

	req := types.QueryConfigOption{
		This: b.Reference(),
	}
	res, err := methods.QueryConfigOption(ctx, b.Client(), &req)
	if err != nil {
		return "", err
	}
	if res.Returnval == nil {
		return "", errors.New("no config options were found for the supplied criteria")
	}
	for _, osd := range res.Returnval.GuestOSDescriptor {
		if osd.Id == guest {
			family := osd.Family
			log.Printf("[DEBUG] OSFamily: family for %q is %q", guest, family)
			return family, nil
		}
	}
	return "", fmt.Errorf("could not find guest ID %q", guest)
}

// QueryConfigOptionDescriptor returns a list the list of ConfigOption keys
// available on the environment that this browser targets. The keys can be used
// as query options for DefaultDevices and other functions, facilitating the
// specification of results specific to a certain VM version.
func (b *EnvironmentBrowser) QueryConfigOptionDescriptor(ctx context.Context) ([]types.VirtualMachineConfigOptionDescriptor, error) {
	req := types.QueryConfigOptionDescriptor{
		This: b.Reference(),
	}
	res, err := methods.QueryConfigOptionDescriptor(ctx, b.Client(), &req)
	if err != nil {
		return nil, err
	}
	return res.Returnval, nil
}
