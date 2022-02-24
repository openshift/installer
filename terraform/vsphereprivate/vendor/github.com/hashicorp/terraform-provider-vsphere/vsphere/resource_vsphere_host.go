package vsphere

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	gtask "github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVsphereHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceVsphereHostCreate,
		Read:   resourceVsphereHostRead,
		Update: resourceVsphereHostUpdate,
		Delete: resourceVsphereHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vSphere datacenter the host will belong to.",
			},
			"cluster": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of the vSphere cluster the host will belong to.",
				ConflictsWith: []string{"cluster_managed"},
			},
			"cluster_managed": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "Must be set if host is a member of a managed compute_cluster resource.",
				ConflictsWith: []string{"cluster"},
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FQDN or IP address of the host.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the administration account of the host.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of the administration account of the host.",
				Sensitive:   true,
			},
			"thumbprint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host's certificate SHA-1 thumbprint. If not set then the CA that signed the host's certificate must be trusted.",
			},
			"license": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "License key that will be applied to this host.",
			},
			"force": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Force add the host to vsphere, even if it's already managed by a different vSphere instance.",
				Default:     false,
			},
			"connected": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the state of the host. If set to false then the host will be asked to disconnect.",
				Default:     true,
			},
			"maintenance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the host's maintenance mode. Default is false",
				Default:     false,
			},
			"lockdown": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Set the host's lockdown status. Default is disabled. Valid options are 'disabled', 'normal', 'strict'",
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "normal", "strict"}, true),
			},
		},
	}
}

func resourceVsphereHostRead(d *schema.ResourceData, meta interface{}) error {

	// NOTE: Destroying the host without telling vsphere about it will result in us not
	// knowing that the host does not exist any more.

	// Look for host
	client := meta.(*VSphereClient).vimClient
	hostID := d.Id()

	// Find host and get reference to it.
	hs, err := hostsystem.FromID(client, hostID)
	if err != nil {
		if viapi.IsManagedObjectNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while searching host %s. Error: %s ", hostID, err)
	}

	maintenanceState, err := hostsystem.HostInMaintenance(hs)
	if err != nil {
		return fmt.Errorf("error while checking maintenance status for host %s. Error: %s", hostID, err)
	}
	d.Set("maintenance", maintenanceState)

	// Retrieve host's properties.
	log.Printf("[DEBUG] Got host %s", hs.String())
	host, err := hostsystem.Properties(hs)
	if err != nil {
		return fmt.Errorf("error while retrieving properties for host %s. Error: %s", hostID, err)
	}

	if host.Parent != nil && host.Parent.Type == "ClusterComputeResource" && !d.Get("cluster_managed").(bool) {
		d.Set("cluster", host.Parent.Value)
	} else {
		d.Set("cluster", "")
	}

	connectionState, err := hostsystem.GetConnectionState(hs)
	if err != nil {
		return fmt.Errorf("error while getting connection state for host %s. Error: %s", hostID, err)
	}

	if connectionState == types.HostSystemConnectionStateDisconnected {
		// Config and LicenseManager cannot be used while the host is
		// disconnected.
		d.Set("connected", false)
		return nil
	}
	d.Set("connected", true)

	lockdownMode, err := hostLockdownString(host.Config.LockdownMode)
	if err != nil {
		return err
	}

	log.Printf("Setting lockdown to %s", lockdownMode)
	d.Set("lockdown", lockdownMode)

	licenseKey := d.Get("license").(string)
	if licenseKey != "" {
		licFound, err := isLicenseAssigned(client.Client, hostID, licenseKey)
		if err != nil {
			return fmt.Errorf("error while checking license assignment for host %s. Error: %s", hostID, err)
		}

		if !licFound {
			d.Set("license", "")
		}
	}

	return nil
}

func resourceVsphereHostCreate(d *schema.ResourceData, meta interface{}) error {
	err := validateFields(d)
	if err != nil {
		return err
	}

	client := meta.(*VSphereClient).vimClient

	hcs := buildHostConnectSpec(d)

	licenseKey := d.Get("license").(string)

	if licenseKey != "" {
		licFound, err := licenseExists(client.Client, licenseKey)
		if err != nil {
			return fmt.Errorf("error while looking for license key. Error: %s", err)
		}

		if !licFound {
			return fmt.Errorf("license key supplied (%s) did not match against known license keys", licenseKey)
		}
	}

	var connectedState bool
	val := d.Get("connected")
	if val == nil {
		connectedState = true
	} else {
		connectedState = val.(bool)
	}

	var task *object.Task
	clusterID := d.Get("cluster").(string)
	if clusterID != "" {
		ccr, err := clustercomputeresource.FromID(client, clusterID)
		if err != nil {
			return fmt.Errorf("error while searching cluster %s. Error: %s", clusterID, err)
		}

		task, err = ccr.AddHost(context.TODO(), hcs, connectedState, &licenseKey, nil)
		if err != nil {
			return fmt.Errorf("error while adding host with hostname %s to cluster %s.  Error: %s", d.Get("hostname").(string), clusterID, err)
		}
	} else {
		dcId := d.Get("datacenter").(string)
		dc, err := datacenterFromID(client, dcId)
		if err != nil {
			return fmt.Errorf("error while retrieving datacenter object for datacenter: %s. Error: %s", dcId, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
		defer cancel()
		var dcProps mo.Datacenter
		if err := dc.Properties(ctx, dc.Reference(), nil, &dcProps); err != nil {
			return fmt.Errorf("error while retrieving properties for datacenter %s. Error: %s", dcId, err)
		}

		hostFolder := object.NewFolder(client.Client, dcProps.HostFolder)
		task, err = hostFolder.AddStandaloneHost(context.TODO(), hcs, connectedState, &licenseKey, nil)
		if err != nil {
			return fmt.Errorf("error while adding standalone host %s. Error: %s", hcs.HostName, err)
		}
	}

	p := property.DefaultCollector(client.Client)
	res, err := gtask.Wait(context.TODO(), task.Reference(), p, nil)
	if err != nil {
		return fmt.Errorf("host addition failed. %s", err)
	}
	taskResult := res.Result

	var hostID string
	taskResultType := taskResult.(types.ManagedObjectReference).Type
	switch taskResultType {
	case "ComputeResource":
		computeResource := object.NewComputeResource(client.Client, taskResult.(types.ManagedObjectReference))
		crHosts, err := computeResource.Hosts(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to retrieve created computeResource Hosts. Error: %s", err)
		}
		hostID = crHosts[0].Reference().Value
		log.Printf("[DEBUG] standalone hostID: %s", hostID)
	case "HostSystem":
		hostID = taskResult.(types.ManagedObjectReference).Value
	default:
		return fmt.Errorf("unexpected task result type encountered. Got %s while waiting ComputeResourceType or Hostsystem", taskResultType)
	}
	log.Printf("[DEBUG] Host added with ID %s", hostID)
	d.SetId(hostID)

	host, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("failed while retrieving host object for host %s. Error: %s", hostID, err)
	}

	lockdownModeString := d.Get("lockdown").(string)
	lockdownMode, err := hostLockdownType(lockdownModeString)
	if err != nil {
		return err
	}

	if connectedState {
		hostProps, err := hostsystem.Properties(host)
		if err != nil {
			return fmt.Errorf("error while retrieving properties for host %s. Error: %s", hostID, err)
		}

		hamRef := hostProps.ConfigManager.HostAccessManager.Reference()
		ham := NewHostAccessManager(client.Client, hamRef)
		err = ham.ChangeLockdownMode(context.TODO(), lockdownMode)
		if err != nil {
			return fmt.Errorf("error while changing lockdown mode for host %s. Error: %s", hostID, err)
		}
	}

	maintenanceMode := d.Get("maintenance").(bool)
	if maintenanceMode {
		err = hostsystem.EnterMaintenanceMode(host, provider.DefaultAPITimeout, true)
	} else {
		err = hostsystem.ExitMaintenanceMode(host, provider.DefaultAPITimeout)
	}
	if err != nil {
		return fmt.Errorf("error while toggling maintenance mode for host %s. Error: %s", hostID, err)
	}

	return resourceVsphereHostRead(d, meta)
}

func resourceVsphereHostUpdate(d *schema.ResourceData, meta interface{}) error {
	err := validateFields(d)
	if err != nil {
		return err
	}

	client := meta.(*VSphereClient).vimClient

	// First let's establish where we are and where we want to go
	var desiredConnectionState bool
	if d.HasChange("connected") {
		_, newVal := d.GetChange("connected")
		desiredConnectionState = newVal.(bool)
	} else {
		desiredConnectionState = d.Get("connected").(bool)
	}

	hostID := d.Id()
	hostObject, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem object for host ID %s. Error: %s", hostID, err)
	}

	actualConnectionState, err := hostsystem.GetConnectionState(hostObject)
	if err != nil {
		return fmt.Errorf("error while retrieving connection state for host %s. Error: %s", hostID, err)
	}

	// Have there been any changes that warrant a reconnect?
	reconnect := false
	connectionKeys := []string{"hostname", "username", "password", "thumbprint"}
	for _, k := range connectionKeys {
		if d.HasChange(k) {
			reconnect = true
			break
		}
	}

	// Decide if we're going to reconnect or not
	reconnectNeeded, err := shouldReconnect(d, meta, actualConnectionState, desiredConnectionState, reconnect)
	if err != nil {
		return err
	}

	switch reconnectNeeded {
	case 1:
		err := resourceVSphereHostReconnect(d, meta)
		if err != nil {
			return fmt.Errorf("error while reconnecting host %s. Error: %s", hostID, err)
		}
	case -1:
		err := resourceVSphereHostDisconnect(d, meta)
		if err != nil {
			return fmt.Errorf("error while disconnecting host %s. Error: %s", hostID, err)
		}
	case 0:
		break
	}

	mutableKeys := map[string]func(*schema.ResourceData, interface{}, interface{}, interface{}) error{
		"license":     resourceVSphereHostUpdateLicense,
		"cluster":     resourceVSphereHostUpdateCluster,
		"maintenance": resourceVSphereHostUpdateMaintenanceMode,
		"lockdown":    resourceVSphereHostUpdateLockdownMode,
		"thumbprint":  resourceVSphereHostUpdateThumbprint,
	}
	for k, v := range mutableKeys {
		log.Printf("[DEBUG] Checking if key %s changed", k)
		if !d.HasChange(k) {
			continue
		}
		log.Printf("[DEBUG] Key %s has change, processing", k)
		old, newVal := d.GetChange(k)
		err := v(d, meta, old, newVal)
		if err != nil {
			return fmt.Errorf("error while updating %s: %s", k, err)
		}
	}
	return resourceVsphereHostRead(d, meta)
}

func resourceVsphereHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	hostID := d.Id()

	hs, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem object for host ID %s. Error: %s", hostID, err)
	}

	connectionState, err := hostsystem.GetConnectionState(hs)
	if err != nil {
		return fmt.Errorf("error while retrieving connection state for host %s. Error: %s", hostID, err)
	}

	if connectionState != types.HostSystemConnectionStateDisconnected {
		// We cannot put a disconnected server in maintenance mode.
		err = resourceVSphereHostDisconnect(d, meta)
		if err != nil {
			return fmt.Errorf("error while disconnecting host: %s", err.Error())
		}
	}

	hostProps, err := hostsystem.Properties(hs)
	if err != nil {
		return fmt.Errorf("error while retrieving properties fort host %s. Error: %s", hostID, err)
	}

	// If this is a standalone host we need to destroy the ComputeResource object
	// and not the Hostsystem itself.
	var task *object.Task
	if hostProps.Parent.Type == "ComputeResource" {
		cr := object.NewComputeResource(client.Client, *hostProps.Parent)
		task, err = cr.Destroy(context.TODO())
		if err != nil {
			return fmt.Errorf("error while submitting destroy task for compute resource %s. Error: %s", hostProps.Parent.Value, err)
		}
	} else {
		task, err = hs.Destroy(context.TODO())
		if err != nil {
			return fmt.Errorf("error while submitting destroy task for host system %s. Error: %s", hostProps.Parent.Value, err)
		}
	}
	p := property.DefaultCollector(client.Client)
	_, err = gtask.Wait(context.TODO(), task.Reference(), p, nil)
	if err != nil {
		return fmt.Errorf("error while waiting for host (%s) to be removed: %s", hostID, err)
	}
	return nil
}

func resourceVSphereHostUpdateLockdownMode(d *schema.ResourceData, meta, old, newVal interface{}) error {
	client := meta.(*VSphereClient).vimClient
	hostID := d.Id()
	host, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem object for host ID %s. Error: %s", hostID, err)
	}
	lockdownModeString := newVal.(string)
	lockdownMode, err := hostLockdownType(lockdownModeString)
	if err != nil {
		return err
	}

	var hostProps mo.HostSystem
	err = host.Properties(context.TODO(), host.ConfigManager().Reference(), []string{"configManager.hostAccessManager"}, &hostProps)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem properties for host ID %s. Error: %s", hostID, err)

	}

	hamRef := hostProps.ConfigManager.HostAccessManager.Reference()
	ham := NewHostAccessManager(client.Client, hamRef)
	err = ham.ChangeLockdownMode(context.TODO(), lockdownMode)
	if err != nil {
		return fmt.Errorf("error while changing lonckdown mode for host ID %s to %s. Error: %s", hostID, lockdownMode, err)

	}

	return nil
}

func resourceVSphereHostUpdateMaintenanceMode(d *schema.ResourceData, meta, old, newVal interface{}) error {
	client := meta.(*VSphereClient).vimClient
	hostID := d.Id()

	host, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem object for host ID %s. Error: %s", hostID, err)
	}

	maintenanceMode := newVal.(bool)
	if maintenanceMode {
		err = hostsystem.EnterMaintenanceMode(host, provider.DefaultAPITimeout, true)
	} else {
		err = hostsystem.ExitMaintenanceMode(host, provider.DefaultAPITimeout)
	}
	if err != nil {
		return fmt.Errorf("error while toggling maintenance mode for host %s. Error: %s", host.Name(), err)
	}
	return nil
}

func resourceVSphereHostUpdateLicense(d *schema.ResourceData, meta, old, newVal interface{}) error {
	client := meta.(*VSphereClient).vimClient
	lm := license.NewManager(client.Client)
	lam, err := lm.AssignmentManager(context.TODO())
	if err != nil {
		return fmt.Errorf("error while accessing License Assignment Manager endpoint. Error: %s", err)
	}
	_, err = lam.Update(context.TODO(), d.Id(), newVal.(string), "")
	if err != nil {
		return fmt.Errorf("error while updating license. error: %s", err)
	}
	return nil
}

func resourceVSphereHostUpdateCluster(d *schema.ResourceData, meta, old, newVal interface{}) error {
	client := meta.(*VSphereClient).vimClient
	hostID := d.Id()
	newClusterID := newVal.(string)

	newCluster, err := clustercomputeresource.FromID(client, newClusterID)
	if err != nil {
		return fmt.Errorf("error while searching newVal cluster %s. Error: %s", newClusterID, err)
	}

	hs, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return fmt.Errorf("error while retrieving HostSystem object for host ID %s. Error: %s", hostID, err)
	}

	err = hostsystem.EnterMaintenanceMode(hs, provider.DefaultAPITimeout, true)
	if err != nil {
		return fmt.Errorf("error while putting host to maintenance mode: %s", err.Error())
	}

	task, err := newCluster.MoveInto(context.TODO(), hs)
	if err != nil {
		return fmt.Errorf("error while moving HostSystem with ID %s to new cluster. Error: %s", hostID, err)
	}
	p := property.DefaultCollector(client.Client)
	_, err = gtask.Wait(context.TODO(), task.Reference(), p, nil)
	if err != nil {
		return fmt.Errorf("error while moving host to new cluster (%s): %s", newClusterID, err)
	}

	err = hostsystem.ExitMaintenanceMode(hs, provider.DefaultAPITimeout)
	if err != nil {
		return fmt.Errorf("error while taking host out of maintenance mode: %s", err.Error())
	}

	return nil
}

func resourceVSphereHostUpdateThumbprint(d *schema.ResourceData, meta, old, newVal interface{}) error {
	return resourceVSphereHostReconnect(d, meta)
}

func resourceVSphereHostReconnect(d *schema.ResourceData, meta interface{}) error {
	hostID := d.Id()
	client := meta.(*VSphereClient).vimClient
	host := object.NewHostSystem(client.Client, types.ManagedObjectReference{Type: "HostSystem", Value: d.Id()})
	hcs := buildHostConnectSpec(d)

	task, err := host.Reconnect(context.TODO(), &hcs, nil)
	if err != nil {
		return fmt.Errorf("error while reconnecting host with ID %s. Error: %s", hostID, err)
	}

	p := property.DefaultCollector(client.Client)
	_, err = gtask.Wait(context.TODO(), task.Reference(), p, nil)
	if err != nil {
		return fmt.Errorf("error while reconnecting host(%s): %s", hostID, err)
	}

	maintenanceState, err := hostsystem.HostInMaintenance(host)
	if err != nil {
		return fmt.Errorf("error while retrieving host maintenance status for host %s. Error: %s", host.Name(), err)
	}

	maintenanceConfig := d.Get("maintenance").(bool)
	if maintenanceState && !maintenanceConfig {
		err := hostsystem.ExitMaintenanceMode(host, provider.DefaultAPITimeout)
		if err != nil {
			return fmt.Errorf("error while taking host %s out of maintenance mode. Error: %s", host.Name(), err)
		}
	}
	return nil
}

func resourceVSphereHostDisconnect(d *schema.ResourceData, meta interface{}) error {
	hostID := d.Id()
	client := meta.(*VSphereClient).vimClient
	host := object.NewHostSystem(client.Client, types.ManagedObjectReference{Type: "HostSystem", Value: d.Id()})
	task, err := host.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error while disconnecting host %s. Error: %s", host.Name(), err)
	}

	p := property.DefaultCollector(client.Client)
	_, err = gtask.Wait(context.TODO(), task.Reference(), p, nil)
	if err != nil {
		return fmt.Errorf("error while disconnecting host(%s): %s", hostID, err)
	}
	return nil
}

func shouldReconnect(d *schema.ResourceData, meta interface{}, actual types.HostSystemConnectionState, desired, shouldReconnect bool) (int, error) {
	log.Printf("[DEBUG] Figuring out if we need to do something about the host's connection")

	// desired state is connected and one of the connectionKeys has changed
	if shouldReconnect && desired {
		log.Printf("[DEBUG] Desired state is connected and one of the settings relevant to the connection changed. Reconnecting")
		return 1, nil
	}

	// desired state is connected and actual state is disconnected
	if desired && (actual != types.HostSystemConnectionStateConnected) {
		log.Printf("[DEBUG] Desired state is connected but host is not connected. Reconnecting")
		return 1, nil
	}

	// desired state is connected and actual state is connected (or host is missing heartbeats) and
	// none of the connectionKeys have changed.
	if desired && (actual != types.HostSystemConnectionStateDisconnected) && !shouldReconnect {
		log.Printf("[DEBUG] Desired state is connected and host is connected and no changes in config. Noop")
		return 0, nil
	}

	// desired state is disconnected and host is disconnected
	if !desired && (actual == types.HostSystemConnectionStateDisconnected) {
		log.Printf("[DEBUG] Desired state is disconnected and host is disconnected")
		return 0, nil
	}

	if !desired && (actual != types.HostSystemConnectionStateDisconnected) {
		log.Printf("[DEBUG] Desired state is disconnected but host is not disconnected. Disconnecting")
		return -1, nil
	}

	log.Printf("[DEBUG] Unexpected combination of desired and actual states, not sure how to handle. Please submit a bug report.")
	return 255, fmt.Errorf("unexpected combination of connection states")
}

func hostLockdownType(lockdownMode string) (types.HostLockdownMode, error) {
	lockdownModes := map[string]types.HostLockdownMode{
		"disabled": types.HostLockdownModeLockdownDisabled,
		"normal":   types.HostLockdownModeLockdownNormal,
		"strict":   types.HostLockdownModeLockdownStrict,
	}

	log.Printf("Looking for mode %s in lockdown modes %#v", lockdownMode, lockdownModes)
	if modeString, ok := lockdownModes[lockdownMode]; ok {
		log.Printf("Found match for %s. Returning %s.", lockdownMode, modeString)
		return modeString, nil
	}
	return "", fmt.Errorf("unknown Lockdown mode encountered")
}

func hostLockdownString(lockdownMode types.HostLockdownMode) (string, error) {
	lockdownModes := map[types.HostLockdownMode]string{
		types.HostLockdownModeLockdownDisabled: "disabled",
		types.HostLockdownModeLockdownNormal:   "normal",
		types.HostLockdownModeLockdownStrict:   "strict",
	}

	log.Printf("Looking for mode %s in lockdown modes %#v", lockdownMode, lockdownModes)
	if modeString, ok := lockdownModes[lockdownMode]; ok {
		log.Printf("Found match for %s. Returning %s.", lockdownMode, modeString)
		return modeString, nil
	}
	return "", fmt.Errorf("unknown Lockdown mode encountered")
}

func buildHostConnectSpec(d *schema.ResourceData) types.HostConnectSpec {
	hcs := types.HostConnectSpec{
		HostName:      d.Get("hostname").(string),
		UserName:      d.Get("username").(string),
		Password:      d.Get("password").(string),
		SslThumbprint: d.Get("thumbprint").(string),
		Force:         d.Get("force").(bool),
	}
	return hcs
}

func isLicenseAssigned(client *vim25.Client, hostID, licenseKey string) (bool, error) {
	ctx := context.TODO()
	lm := license.NewManager(client)
	am, err := lm.AssignmentManager(ctx)
	if err != nil {
		return false, err
	}

	licenses, err := am.QueryAssigned(ctx, hostID)
	if err != nil {
		return false, err
	}

	licFound := false
	for _, lic := range licenses {
		if licenseKey == lic.AssignedLicense.LicenseKey {
			licFound = true
			break
		}
	}
	return licFound, nil
}

func licenseExists(client *vim25.Client, licenseKey string) (bool, error) {
	ctx := context.TODO()
	lm := license.NewManager(client)
	ll, err := lm.List(ctx)
	if err != nil {
		return false, err
	}

	licFound := false
	for _, l := range ll {
		if l.LicenseKey == licenseKey {
			licFound = true
			break
		}
	}
	return licFound, nil
}

// Make sure input makes sense
func validateFields(d *schema.ResourceData) error {
	_, dcSet := d.GetOk("datacenter")
	_, clusterSet := d.GetOk("cluster")
	if dcSet && clusterSet {
		return fmt.Errorf("datacenter and cluster arguments are mutually exclusive")
	}
	return nil
}

// --------------
// Implementing stuff govmomi should provide for us
//
//

type HostAccessManager struct {
	object.Common
}

func NewHostAccessManager(c *vim25.Client, ref types.ManagedObjectReference) *HostAccessManager {
	return &HostAccessManager{
		Common: object.NewCommon(c, ref),
	}
}

func (h HostAccessManager) ChangeLockdownMode(ctx context.Context, mode types.HostLockdownMode) error {
	req := types.ChangeLockdownMode{
		This: h.Reference(),
		Mode: mode,
	}
	_, err := methods.ChangeLockdownMode(ctx, h.Client(), &req)
	return err
}
