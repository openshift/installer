package virtualmachine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/vcenter"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/vappcontainer"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const powerOnWaitMilli = 500

var errGuestShutdownTimeout = errors.New("the VM did not power off within the specified amount of time")

// vmUUIDSearchIndexVersion denotes the minimum version we use the SearchIndex
// VM UUID search for. All versions lower than this use ContainerView to find
// the VM.
var vmUUIDSearchIndexVersion = viapi.VSphereVersion{
	Major: 6,
	Minor: 5,
}

// UUIDNotFoundError is an error type that is returned when a
// virtual machine could not be found by UUID.
type UUIDNotFoundError struct {
	s string
}

// VCenterDeploy containss everything required to create a VM from a content
// library item, and reduces the number of arguments passed between functions.
type VCenterDeploy struct {
	VCenterManager *vcenter.Manager

	VMName          string
	Annotation      string
	FolderID        string
	DatastoreID     string
	ResourcePoolID  string
	HostSystemID    string
	StoragePolicyID string
	DiskType        string
	NetworkMap      []vcenter.NetworkMapping
	VAppProperties  []vcenter.Property
	LibraryItem     *library.Item
}

// Error implements error for UUIDNotFoundError.
func (e *UUIDNotFoundError) Error() string {
	return e.s
}

// newUUIDNotFoundError returns a new UUIDNotFoundError with the
// text populated.
func newUUIDNotFoundError(s string) *UUIDNotFoundError {
	return &UUIDNotFoundError{
		s: s,
	}
}

// IsUUIDNotFoundError returns true if the error is a UUIDNotFoundError.
func IsUUIDNotFoundError(err error) bool {
	_, ok := err.(*UUIDNotFoundError)
	return ok
}

func List(client *govmomi.Client) ([]*object.VirtualMachine, error) {
	return vmsByPath(client, "/*")
}

func vmsByPath(client *govmomi.Client, path string) ([]*object.VirtualMachine, error) {
	ctx := context.TODO()
	var vms []*object.VirtualMachine
	finder := find.NewFinder(client.Client, false)
	es, err := finder.ManagedObjectListChildren(ctx, path+"/*", "vm", "folder", "pool")
	if err != nil {
		return nil, err
	}
	for _, id := range es {
		switch {
		case id.Object.Reference().Type == "VirtualMachine":
			vm, err := FromMOID(client, id.Object.Reference().Value)
			if err != nil {
				return nil, err
			}
			vms = append(vms, vm)
		case id.Object.Reference().Type == "Folder" || id.Object.Reference().Type == "ResourcePool":
			newRPs, err := vmsByPath(client, id.Path)
			if err != nil {
				return nil, err
			}
			vms = append(vms, newRPs...)
		default:
			continue
		}
	}
	return vms, nil
}

// FromUUID locates a virtualMachine by its UUID.
func FromUUID(client *govmomi.Client, uuid string) (*object.VirtualMachine, error) {
	log.Printf("[DEBUG] Locating virtual machine with UUID %q", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	var result object.Reference
	var err error
	version := viapi.ParseVersionFromClient(client)
	expected := vmUUIDSearchIndexVersion
	expected.Product = version.Product
	if version.Older(expected) {
		result, err = virtualMachineFromContainerView(ctx, client, uuid)
	} else {
		result, err = virtualMachineFromSearchIndex(ctx, client, uuid)
	}

	if err != nil {
		return nil, err
	}

	// We need to filter our object through finder to ensure that the
	// InventoryPath field is populated, or else functions that depend on this
	// being present will fail.
	finder := find.NewFinder(client.Client, false)

	vm, err := finder.ObjectReference(ctx, result.Reference())
	if err != nil {
		return nil, err
	}

	// Should be safe to return here. If our reference returned here and is not a
	// VM, then we have bigger problems and to be honest we should be panicking
	// anyway.
	log.Printf("[DEBUG] VM %q found for UUID %q", vm.(*object.VirtualMachine).InventoryPath, uuid)
	return vm.(*object.VirtualMachine), nil
}

// virtualMachineFromSearchIndex gets the virtual machine reference via the
// SearchIndex MO and is the method used to fetch UUIDs on newer versions of
// vSphere.
func virtualMachineFromSearchIndex(ctx context.Context, client *govmomi.Client, uuid string) (object.Reference, error) {
	log.Printf("[DEBUG] Using SearchIndex to look up UUID %q", uuid)
	search := object.NewSearchIndex(client.Client)
	result, err := search.FindByUuid(ctx, nil, uuid, true, structure.BoolPtr(false))
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, newUUIDNotFoundError(fmt.Sprintf("virtual machine with UUID %q not found", uuid))
	}

	return result, nil
}

// virtualMachineFromContainerView is a compatability method that is
// used when the version of vSphere is too old to support using SearchIndex's
// FindByUuid method correctly. This is mainly to facilitate the ability to use
// FromUUID to find both templates in addition to virtual machines, which
// historically was not supported by FindByUuid.
func virtualMachineFromContainerView(ctx context.Context, client *govmomi.Client, uuid string) (object.Reference, error) {
	log.Printf("[DEBUG] Using ContainerView to look up UUID %q", uuid)
	m := view.NewManager(client.Client)

	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = v.Destroy(ctx); err != nil {
			log.Printf("[DEBUG] virtualMachineFromContainerView: Unexpected error destroying container view: %s", err)
		}
	}()

	var vms, results []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"config.uuid"}, &results)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.Config == nil {
			continue
		}
		if result.Config.Uuid == uuid {
			vms = append(vms, result)
		}
	}

	switch {
	case len(vms) < 1:
		return nil, newUUIDNotFoundError(fmt.Sprintf("virtual machine with UUID %q not found", uuid))
	case len(vms) > 1:
		return nil, fmt.Errorf("multiple virtual machines with UUID %q found", uuid)
	}

	return object.NewReference(client.Client, vms[0].Self), nil
}

// FromMOID locates a virtualMachine by its managed
// object reference ID.
func FromMOID(client *govmomi.Client, id string) (*object.VirtualMachine, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "VirtualMachine",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	vm, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	// Should be safe to return here. If our reference returned here and is not a
	// VM, then we have bigger problems and to be honest we should be panicking
	// anyway.
	return vm.(*object.VirtualMachine), nil
}

// FromPath returns a VirtualMachine via its supplied path.
func FromPath(client *govmomi.Client, path string, dc *object.Datacenter) (*object.VirtualMachine, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		finder.SetDatacenter(dc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return finder.VirtualMachine(ctx, path)
}

// Properties is a convenience method that wraps fetching the
// VirtualMachine MO from its higher-level object.
func Properties(vm *object.VirtualMachine) (*mo.VirtualMachine, error) {
	log.Printf("[DEBUG] Fetching properties for VM %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.VirtualMachine
	if err := vm.Properties(ctx, vm.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// ConfigOptions is a convenience method that wraps fetching the VirtualMachine ConfigOptions
// as returned by QueryConfigOption.
func ConfigOptions(vm *object.VirtualMachine) (*types.VirtualMachineConfigOption, error) {

	// First grab the properties so that we can sneak the EnvironmentBrowser out of it
	props, err := Properties(vm)
	if err != nil {
		return nil, err
	}

	// Make a context so we can timeout according to the provider configuration
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	// Build a request for the config option, and then query for configuration options
	log.Printf("[DEBUG] Fetching configuration options for VM %q", vm.InventoryPath)
	request := types.QueryConfigOption{This: props.EnvironmentBrowser}

	response, err := methods.QueryConfigOption(ctx, vm.Client(), &request)
	if err != nil {
		return nil, err
	}

	return response.Returnval, nil
}

// WaitForGuestIP waits for a virtual machine to have an IP address.
//
// The timeout is specified in minutes. If zero or a negative value is passed,
// the waiter returns without error immediately.
func WaitForGuestIP(client *govmomi.Client, vm *object.VirtualMachine, timeout int, ignoredGuestIPs []interface{}) error {
	if timeout < 1 {
		log.Printf("[DEBUG] Skipping IP waiter for VM %q", vm.InventoryPath)
		return nil
	}
	log.Printf(
		"[DEBUG] Waiting for an available IP address on VM %q (timeout = %dm)",
		vm.InventoryPath,
		timeout,
	)

	p := client.PropertyCollector()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()

	err := property.Wait(ctx, p, vm.Reference(), []string{"guest.ipAddress"}, func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			if c.Val == nil {
				continue
			}

			ip := net.ParseIP(c.Val.(string))
			if skipIPAddrForWaiter(ip, ignoredGuestIPs) {
				continue
			}

			return true
		}

		return false
	})

	if err != nil {
		// Provide a friendly error message if we timed out waiting for a routable IP.
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("timeout waiting for an available IP address")
		}
		return err
	}

	log.Printf("[DEBUG] IP address is now available for VM %q", vm.InventoryPath)
	return nil
}

// WaitForGuestNet waits for a virtual machine to have routable network
// access. This is denoted as a gateway, and at least one IP address that can
// reach that gateway. This function supports both IPv4 and IPv6, and returns
// the moment either stack is routable - it doesn't wait for both.
//
// The timeout is specified in minutes. If zero or a negative value is passed,
// the waiter returns without error immediately.
func WaitForGuestNet(client *govmomi.Client, vm *object.VirtualMachine, routable bool, timeout int, ignoredGuestIPs []interface{}) error {
	if timeout < 1 {
		log.Printf("[DEBUG] Skipping network waiter for VM %q", vm.InventoryPath)
		return nil
	}
	log.Printf(
		"[DEBUG] Waiting for an available IP address on VM %q (routable= %t, timeout = %dm)",
		vm.InventoryPath,
		routable,
		timeout,
	)
	var v4gw, v6gw net.IP

	p := client.PropertyCollector()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()

	err := property.Wait(ctx, p, vm.Reference(), []string{"guest.net", "guest.ipStack"}, func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			switch v := c.Val.(type) {
			case types.ArrayOfGuestStackInfo:
				for _, s := range v.GuestStackInfo {
					if s.IpRouteConfig != nil {
						for _, r := range s.IpRouteConfig.IpRoute {
							switch r.Network {
							case "0.0.0.0":
								v4gw = net.ParseIP(r.Gateway.IpAddress)
							case "::":
								v6gw = net.ParseIP(r.Gateway.IpAddress)
							}
						}
					}
				}
			case types.ArrayOfGuestNicInfo:
				for _, n := range v.GuestNicInfo {
					if n.IpConfig != nil {
						for _, addr := range n.IpConfig.IpAddress {
							ip := net.ParseIP(addr.IpAddress)
							if skipIPAddrForWaiter(ip, ignoredGuestIPs) {
								continue
							}
							if !routable {
								// We are done. The rest of this block concerns itself with
								// checking for a routable address, but the waiter has been
								// flagged to not wait for one.
								return true
							}
							var mask net.IPMask
							if ip.To4() != nil {
								mask = net.CIDRMask(int(addr.PrefixLength), 32)
							} else {
								mask = net.CIDRMask(int(addr.PrefixLength), 128)
							}
							if ip.Mask(mask).Equal(v4gw.Mask(mask)) || ip.Mask(mask).Equal(v6gw.Mask(mask)) {
								return true
							}
						}
					}
				}
			}
		}

		return false
	})

	if err != nil {
		// Provide a friendly error message if we timed out waiting for a routable IP.
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("timeout waiting for an available IP address")
		}
		return err
	}

	log.Printf("[DEBUG] IP address(es) is/are now available for VM %q", vm.InventoryPath)
	return nil
}

func skipIPAddrForWaiter(ip net.IP, ignoredGuestIPs []interface{}) bool {
	switch {
	case ip.IsLinkLocalMulticast():
		fallthrough
	case ip.IsLinkLocalUnicast():
		fallthrough
	case ip.IsLoopback():
		fallthrough
	case ip.IsMulticast():
		return true
	default:
		// ignoredGuestIPs pre-validated by Schema!
		for _, ignoredGuestIP := range ignoredGuestIPs {
			if strings.Contains(ignoredGuestIP.(string), "/") {
				_, ignoredIPNet, _ := net.ParseCIDR(ignoredGuestIP.(string))
				if ignoredIPNet.Contains(ip) {
					return true
				}
			} else if net.ParseIP(ignoredGuestIP.(string)).Equal(ip) {
				return true
			}
		}
	}

	return false
}

func blockUntilReadyForMethod(method string, vm *object.VirtualMachine, ctx context.Context) error {
	log.Printf("[DEBUG] blockUntilReadyForMethod: Going to block until %q is no longer in the Disabled Methods list for vm %s", method, vm.Reference().Value)

	for {
		vprops, err := Properties(vm)
		if err != nil {
			return fmt.Errorf("cannot fetch properties of created virtual machine: %s", err)
		}
		stillPending := false
		for _, methodName := range vprops.DisabledMethod {
			if methodName == method {
				stillPending = true
				break
			}
		}

		if !stillPending {
			log.Printf("[DEBUG] blockUntilReadyForMethod: %q no longer disabled for vm %s", method, vm.Reference().Value)
			break
		}

		select {
		case <-time.After(5 * time.Second):
			log.Printf("[DEBUG] blockUntilReadyForMethod: %q still disabled for vm %s, about to check again", method, vm.Reference().Value)
		case <-ctx.Done():
			return fmt.Errorf("blockUntilReadyForMethod: timed out while waiting for %q to become available for vm %s", method, vm.Reference().Value)
		}
	}

	return nil
}

// Create wraps the creation of a virtual machine and the subsequent waiting of
// the task. A higher-level virtual machine object is returned.
func Create(c *govmomi.Client, f *object.Folder, s types.VirtualMachineConfigSpec, p *object.ResourcePool,
	h *object.HostSystem, timeout time.Duration) (*object.VirtualMachine, error) {
	log.Printf("[DEBUG] Creating virtual machine %q", fmt.Sprintf("%s/%s", f.InventoryPath, s.Name))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var task *object.Task
	// Check to see if the resource pool is a vApp
	vc, err := vappcontainer.FromID(c, p.Reference().Value)
	if err != nil {
		if !viapi.IsManagedObjectNotFoundError(err) {
			return nil, err
		}
		task, err = f.CreateVM(ctx, s, p, h)
	} else {
		task, err = vc.CreateChildVM(ctx, s, h)
	}
	if err != nil {
		return nil, err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), timeout)
	defer tcancel()
	result, err := task.WaitForResult(tctx, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Virtual machine %q: creation complete (MOID: %q)", fmt.Sprintf("%s/%s", f.InventoryPath, s.Name), result.Result.(types.ManagedObjectReference).Value)
	return FromMOID(c, result.Result.(types.ManagedObjectReference).Value)
}

// Clone wraps the creation of a virtual machine and the subsequent waiting of
// the task. A higher-level virtual machine object is returned.
func Clone(c *govmomi.Client, src *object.VirtualMachine, f *object.Folder, name string, spec types.VirtualMachineCloneSpec, timeout int) (*object.VirtualMachine, error) {
	log.Printf("[DEBUG] Cloning virtual machine %q", fmt.Sprintf("%s/%s", f.InventoryPath, name))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	task, err := src.Clone(ctx, f, name, spec)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			err = errors.New("timeout waiting for clone to complete")
		}
		return nil, err
	}
	result, err := task.WaitForResult(ctx, nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			err = errors.New("timeout waiting for clone to complete")
		}
		return nil, err
	}
	log.Printf("[DEBUG] Virtual machine %q: clone complete (MOID: %q)", fmt.Sprintf("%s/%s", f.InventoryPath, name), result.Result.(types.ManagedObjectReference).Value)
	return FromMOID(c, result.Result.(types.ManagedObjectReference).Value)
}

// Deploy clones a virtual machine from a content library item.
func Deploy(deployData *VCenterDeploy) (*types.ManagedObjectReference, error) {
	log.Printf("[DEBUG] virtualmachine.Deploy: Deploying VM from Content Library item.")
	// Get OVF mappings for NICs

	switch deployData.LibraryItem.Type {
	case library.ItemTypeOVF:
		return deployData.deployOvf()
	case library.ItemTypeVMTX:
		return deployData.deployVmtx()
	default:
		return nil, fmt.Errorf("unsupported library item type: %s", deployData.LibraryItem.Type)
	}
}

func (deployData *VCenterDeploy) deployVmtx() (*types.ManagedObjectReference, error) {
	storage := &vcenter.DiskStorage{
		Datastore: deployData.DatastoreID,
		StoragePolicy: &vcenter.StoragePolicy{
			Policy: deployData.StoragePolicyID,
			Type:   "USE_SOURCE_POLICY",
		},
	}
	if deployData.StoragePolicyID != "" {
		storage.StoragePolicy.Type = "USE_SPECIFIED_POLICY"
	}

	deploy := vcenter.DeployTemplate{
		Name:          deployData.VMName,
		Description:   deployData.Annotation,
		DiskStorage:   storage,
		VMHomeStorage: storage,
		Placement: &vcenter.Placement{
			ResourcePool: deployData.ResourcePoolID,
			Host:         deployData.HostSystemID,
			Folder:       deployData.FolderID,
		},
	}
	ctx := context.TODO()
	return deployData.VCenterManager.DeployTemplateLibraryItem(ctx, deployData.LibraryItem.ID, deploy)
}

func (deployData *VCenterDeploy) deployOvf() (*types.ManagedObjectReference, error) {
	deploy := vcenter.Deploy{
		DeploymentSpec: vcenter.DeploymentSpec{
			Name:               deployData.VMName,
			DefaultDatastoreID: deployData.DatastoreID,
			AcceptAllEULA:      true,
			Annotation:         deployData.Annotation,
			AdditionalParams: []vcenter.AdditionalParams{
				{
					Class: vcenter.ClassDeploymentOptionParams,
					Type:  vcenter.TypeDeploymentOptionParams,
				},
				{
					Class:      vcenter.ClassPropertyParams,
					Type:       vcenter.TypePropertyParams,
					Properties: deployData.VAppProperties,
				},
			},
			NetworkMappings:     deployData.NetworkMap,
			StorageProvisioning: deployData.DiskType,
			StorageProfileID:    deployData.StoragePolicyID,
		},
		Target: vcenter.Target{
			ResourcePoolID: deployData.ResourcePoolID,
			HostID:         deployData.HostSystemID,
			FolderID:       deployData.FolderID,
		},
	}
	ctx := context.TODO()
	return deployData.VCenterManager.DeployLibraryItem(ctx, deployData.LibraryItem.ID, deploy)
}

// VAppProperties converts the vApp properties from the configuration and
// converts them into a slice of vcenter.Properties to be used while deploying
// content library items as VMs.
func VAppProperties(propertyMap map[string]interface{}) []vcenter.Property {
	properties := []vcenter.Property{}
	for key, value := range propertyMap {
		property := vcenter.Property{
			ID:    key,
			Label: "",
			Value: value.(string),
		}
		properties = append(properties, property)
	}
	return properties
}

// DiskType converts standard disk type labels into labels used in content
// library deployment specs.
func DiskType(d *schema.ResourceData) string {
	thin := d.Get("disk.0.thin_provisioned").(bool)
	eagerlyScrub := d.Get("disk.0.eagerly_scrub").(bool)
	switch {
	case thin:
		return "thin"
	case eagerlyScrub:
		return "eagerZeroedThick"
	default:
		return "thick"
	}
}

// Customize wraps the customization of a virtual machine and the subsequent
// waiting of the task.
func Customize(vm *object.VirtualMachine, spec types.CustomizationSpec) error {
	log.Printf("[DEBUG] Sending customization spec to virtual machine %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vm.Customize(ctx, spec)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return task.Wait(tctx)
}

// PowerOn wraps powering on a VM and the waiting for the subsequent task.
func PowerOn(vm *object.VirtualMachine, pTimeout time.Duration) error {
	vmPath := vm.InventoryPath
	log.Printf("[DEBUG] Powering on virtual machine %q", vmPath)
	var ctxTimeout time.Duration
	if pTimeout > provider.DefaultAPITimeout {
		ctxTimeout = pTimeout
	} else {
		ctxTimeout = provider.DefaultAPITimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := blockUntilReadyForMethod("PowerOnVM_Task", vm, ctx)
	if err != nil {
		return err
	}

	// This is the controversial part. Although we take every precaution to make sure the VM
	// is in a state that can be started we have noticed that vsphere will randomly fail to
	// power on the vm with "InvalidState" errors.
	//
	// We're adding a small loop that will try to power on the VM until we hit a timeout
	// or manage to call PowerOnVM_Task successfully.

powerLoop:
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			vprops, err := Properties(vm)
			if err != nil {
				return fmt.Errorf("cannot fetch properties of created virtual machine: %s", err)
			}
			if vprops.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
				log.Printf("[DEBUG] VM %q is powered off, attempting to power on.", vmPath)
				task, err := vm.PowerOn(ctx)
				if err != nil {
					log.Printf("[DEBUG] Failed to submit PowerOn task for vm %q. Error: %s", vmPath, err)
					return fmt.Errorf("failed to submit poweron task for vm %q: %s", vmPath, err)
				}
				err = task.Wait(ctx)
				if err != nil {
					if err.Error() == "The operation is not allowed in the current state." {
						log.Printf("[DEBUG] vm %q cannot be powered on in the current state", vmPath)
						continue powerLoop
					} else {
						log.Printf("[DEBUG] PowerOn task for vm %q failed. Error: %s", vmPath, err)
						return fmt.Errorf("powerOn task for vm %q failed: %s", vmPath, err)
					}
				}
				log.Printf("[DEBUG] PowerOn task for VM %q was successful.", vmPath)
				break powerLoop
			}
		case <-ctx.Done():
			return fmt.Errorf("timed out while trying to power on vm %q", vmPath)
		}
	}
	return nil
}

// PowerOff wraps powering off a VM and the waiting for the subsequent task.
func PowerOff(vm *object.VirtualMachine) error {
	log.Printf("[DEBUG] Forcing power off of virtual machine of %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vm.PowerOff(ctx)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return task.Wait(tctx)
}

// ShutdownGuest wraps the graceful shutdown of a guest VM, and then waiting an
// appropriate amount of time for the guest power state to go to powered off.
// If the VM does not power off in the shutdown period specified by timeout (in
// minutes), an error is returned.
//
// The minimum value for timeout is 1 minute - setting to a 0 or negative value
// is not allowed and will just reset the timeout to the minimum.
func ShutdownGuest(client *govmomi.Client, vm *object.VirtualMachine, timeout int) error {
	log.Printf("[DEBUG] Attempting guest shutdown of virtual machine %q", vm.InventoryPath)
	sctx, scancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer scancel()
	if err := vm.ShutdownGuest(sctx); err != nil {
		return err
	}

	// We now wait on VM power state to be powerOff, via a property collector that waits on power state.
	p := client.PropertyCollector()
	if timeout < 1 {
		timeout = 1
	}
	pctx, pcancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer pcancel()

	err := property.Wait(pctx, p, vm.Reference(), []string{"runtime.powerState"}, func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			switch v := c.Val.(type) {
			case types.VirtualMachinePowerState:
				if v == types.VirtualMachinePowerStatePoweredOff {
					return true
				}
			}
		}

		return false
	})

	if err != nil {
		// Provide a friendly error message if we timed out waiting for a shutdown.
		if pctx.Err() == context.DeadlineExceeded {
			return errGuestShutdownTimeout
		}
		return err
	}
	return nil
}

// GracefulPowerOff is a meta-operation that handles powering down of virtual
// machines. A graceful shutdown is attempted first if possible (VMware tools
// is installed, and the guest state is not suspended), and then, if allowed, a
// power-off is forced if that fails.
func GracefulPowerOff(client *govmomi.Client, vm *object.VirtualMachine, timeout int, force bool) error {
	vprops, err := Properties(vm)
	if err != nil {
		return err
	}
	// First we attempt a guest shutdown if we have VMware tools and if the VM is
	// actually powered on (we don't expect that a graceful shutdown would
	// complete on a suspended VM, so there's really no point in trying).
	if vprops.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn && vprops.Guest != nil && vprops.Guest.ToolsRunningStatus == string(types.VirtualMachineToolsRunningStatusGuestToolsRunning) {
		if err := ShutdownGuest(client, vm, timeout); err != nil {
			if err == errGuestShutdownTimeout && !force {
				return err
			}
		} else {
			return nil
		}
	}
	// If the guest shutdown failed (and we were allowed to proceed), or
	// conditions did not satisfy the criteria for a graceful shutdown, do a full
	// power-off of the VM.
	return PowerOff(vm)
}

// MoveToFolder moves a virtual machine to the specified folder.
func MoveToFolder(client *govmomi.Client, vm *object.VirtualMachine, relative string) error {
	log.Printf("[DEBUG] Moving virtual %q to VM path %q", vm.InventoryPath, relative)
	f, err := folder.VirtualMachineFolderFromObject(client, vm, relative)
	if err != nil {
		return err
	}
	return folder.MoveObjectTo(vm.Reference(), f)
}

// Reconfigure wraps the Reconfigure task and the subsequent waiting for
// the task to complete.
func Reconfigure(vm *object.VirtualMachine, spec types.VirtualMachineConfigSpec) error {
	log.Printf("[DEBUG] Reconfiguring virtual machine %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vm.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return task.Wait(tctx)
}

// Relocate wraps the Relocate task and the subsequent waiting for the task to
// complete.
func Relocate(vm *object.VirtualMachine, spec types.VirtualMachineRelocateSpec, timeout int) error {
	log.Printf("[DEBUG] Beginning migration of virtual machine %q (timeout %d)", vm.InventoryPath, timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	task, err := vm.Relocate(ctx, spec, "")
	if err != nil {
		return err
	}
	if err := task.Wait(ctx); err != nil {
		// Provide a friendly error message if we timed out waiting for the migration.
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("timeout waiting for migration to complete")
		}
	}
	return nil
}

// Destroy wraps the Destroy task and the subsequent waiting for the task to
// complete.
func Destroy(vm *object.VirtualMachine) error {
	log.Printf("[DEBUG] Deleting virtual machine %q", vm.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vm.Destroy(ctx)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return task.Wait(tctx)
}

// MOIDForUUIDResult is a struct that holds a virtual machine UUID -> MOID
// association, designed to be used as a helper for mass returning the results
// of translating multiple UUIDs to managed object IDs for various virtual
// machine operations.
type MOIDForUUIDResult struct {
	// The UUID of a virtual machine.
	UUID string

	// The matching managed object reference ID for the virtual machine at the ID
	// referenced by UUID.
	MOID string
}

// MOIDForUUIDResults is a slice that holds multiple MOIDForUUIDResult structs.
type MOIDForUUIDResults []MOIDForUUIDResult

// MOIDForUUID returns the managed object reference ID for a specific virtual
// machine UUID and returns a MOIDForUUIDResult with the appropriate
// association.
func MOIDForUUID(client *govmomi.Client, uuid string) (MOIDForUUIDResult, error) {
	vm, err := FromUUID(client, uuid)
	if err != nil {
		return MOIDForUUIDResult{}, err
	}
	return MOIDForUUIDResult{
		UUID: uuid,
		MOID: vm.Reference().Value,
	}, nil
}

// UUIDForMOID returns the managed object reference ID for a specific virtual
// machine MOID and returns a MOIDForUUIDResult with the appropriate
// association.
func UUIDForMOID(client *govmomi.Client, moid string) (MOIDForUUIDResult, error) {
	vm, err := FromMOID(client, moid)
	if err != nil {
		return MOIDForUUIDResult{}, err
	}
	props, err := Properties(vm)
	if err != nil {
		return MOIDForUUIDResult{}, err
	}
	return MOIDForUUIDResult{
		UUID: props.Config.Uuid,
		MOID: vm.Reference().Value,
	}, nil
}

// MOIDsForUUIDs returns a MOIDForUUIDResults for a list of UUIDs. If one UUID
// cannot be found, an error is returned. There are no partial results
// returned.
func MOIDsForUUIDs(client *govmomi.Client, uuids []string) (MOIDForUUIDResults, error) {
	var results MOIDForUUIDResults
	for _, uuid := range uuids {
		result, err := MOIDForUUID(client, uuid)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// UUIDsForMOIDs returns a MOIDForUUIDResults for a list of MOIDs. If one MOID
// cannot be found, an error is returned. There are no partial results
// returned.
func UUIDsForMOIDs(client *govmomi.Client, moids []string) (MOIDForUUIDResults, error) {
	var results MOIDForUUIDResults
	for _, uuid := range moids {
		result, err := UUIDForMOID(client, uuid)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// UUIDsForManagedObjectReferences returns a MOIDForUUIDResults for a list of
// ManagedObjectReferences. If one cannot be found, an error is returned. There
// are no partial results returned.
func UUIDsForManagedObjectReferences(client *govmomi.Client, refs []types.ManagedObjectReference) (MOIDForUUIDResults, error) {
	var moids []string
	for _, ref := range refs {
		moids = append(moids, ref.Value)
	}
	return UUIDsForMOIDs(client, moids)
}

// MOIDs returns all MOIDs in a MOIDForUUIDResults.
func (r MOIDForUUIDResults) MOIDs() []string {
	var moids []string
	for _, result := range r {
		moids = append(moids, result.MOID)
	}
	return moids
}

// ManagedObjectReferences returns all MOIDs in a MOIDForUUIDResults, as
// ManagedObjectReferences as type VirtualMachine.
func (r MOIDForUUIDResults) ManagedObjectReferences() []types.ManagedObjectReference {
	var refs []types.ManagedObjectReference
	for _, result := range r {
		refs = append(refs, types.ManagedObjectReference{
			Type:  "VirtualMachine",
			Value: result.MOID,
		})
	}
	return refs
}

// UUIDs returns all UUIDs in a MOIDForUUIDResults.
func (r MOIDForUUIDResults) UUIDs() []string {
	var uuids []string
	for _, result := range r {
		uuids = append(uuids, result.UUID)
	}
	return uuids
}

// GetHardwareVersionID gets the hardware version string from integer
func GetHardwareVersionID(vint int) string {
	// hardware_version isn't set, so return an empty string.
	if vint == 0 {
		return ""
	}
	return fmt.Sprintf("vmx-%d", vint)
}

// GetHardwareVersionNumber gets the hardware version number from string.
func GetHardwareVersionNumber(vstring string) int {
	vstring = strings.TrimPrefix(vstring, "vmx-")
	v, err := strconv.Atoi(vstring)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse hardware version: %s", vstring)
	}
	return v
}

// SetHardwareVersion sets the virtual machine's hardware version. The virtual
// machine must be powered off, and the version can only be increased.
func SetHardwareVersion(vm *object.VirtualMachine, target int) error {

	// First query for the configuration options of the vm
	copts, err := ConfigOptions(vm)
	if err != nil {
		return err
	}

	// Now we can grab its version to compare against the target
	current := int(copts.HardwareOptions.HwVersion)
	log.Printf("[DEBUG] Found current hardware version: %d", current)

	// If the hardware version matches, then we're done here and can leave.
	if current == target || target == 0 {
		return nil
	}

	// Otherwise we need to validate it to ensure we're not downgrading
	// the hardware version.
	log.Printf("[DEBUG] Validating the target hardware version: %d", target)
	if err := ValidateHardwareVersion(current, target); err != nil {
		return err
	}

	// We can now proceed to upgrade the hardware version on the vm
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()

	log.Printf("[DEBUG] Upgrading VM from hw version %d to hw version %d", current, target)
	task, err := vm.UpgradeVM(ctx, GetHardwareVersionID(target))
	_, err = task.WaitForResult(ctx, nil)
	return err
}

// ValidateHardwareVersion checks that the target hardware version is equal to
// or greater than the current hardware version.
func ValidateHardwareVersion(current, target int) error {
	switch {
	case target == 0:
		return nil
	case target < current:
		return fmt.Errorf("Cannot downgrade virtual machine hardware version. current: %d, target: %d", current, target)
	}
	return nil
}
