package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereHAVMOverrideName = "vsphere_ha_vm_override"

var vmOverrideClusterDasVMSettingsIsolationResponseAllowedValues = []string{
	string(types.ClusterDasVmSettingsIsolationResponseClusterIsolationResponse),
	string(types.ClusterDasVmSettingsIsolationResponseNone),
	string(types.ClusterDasVmSettingsIsolationResponsePowerOff),
	string(types.ClusterDasVmSettingsIsolationResponseShutdown),
}

var vmOverrideClusterDasConfigInfoServiceStateAllowedValues = []string{
	string(types.ClusterDasVmSettingsRestartPriorityClusterRestartPriority),
	string(types.ClusterDasVmSettingsRestartPriorityLowest),
	string(types.ClusterDasVmSettingsRestartPriorityLow),
	string(types.ClusterDasVmSettingsRestartPriorityMedium),
	string(types.ClusterDasVmSettingsRestartPriorityHigh),
	string(types.ClusterDasVmSettingsRestartPriorityHighest),
}

var vmOverrideClusterVMStorageProtectionForPDLAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionWarning),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive),
}

var vmOverrideClusterVMStorageProtectionForAPDAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionDisabled),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionWarning),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartConservative),
	string(types.ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive),
}

var vmOverrideClusterVMReactionOnAPDClearedAllowedValues = []string{
	string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedUseClusterDefault),
	string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedNone),
	string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedReset),
}

func resourceVSphereHAVMOverride() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereHAVMOverrideCreate,
		Read:   resourceVSphereHAVMOverrideRead,
		Update: resourceVSphereHAVMOverrideUpdate,
		Delete: resourceVSphereHAVMOverrideDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereHAVMOverrideImport,
		},

		Schema: map[string]*schema.Schema{
			"compute_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the cluster.",
			},
			"virtual_machine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the virtual machine.",
			},
			// Host monitoring - VM restarts
			"ha_vm_restart_priority": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasVmSettingsRestartPriorityClusterRestartPriority),
				Description:  "The restart priority for this virtual machine when vSphere detects a host failure. Can be one of clusterRestartPriority, lowest, low, medium, high, or highest.",
				ValidateFunc: validation.StringInSlice(vmOverrideClusterDasConfigInfoServiceStateAllowedValues, false),
			},
			"ha_vm_restart_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The maximum time, in seconds, that vSphere HA will wait for the virtual machine to be ready. Use -1 to use the cluster default.",
			},
			// Host monitoring - host isolation
			"ha_host_isolation_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasVmSettingsIsolationResponseClusterIsolationResponse),
				Description:  "The action to take on this virtual machine when a host is isolated from the rest of the cluster. Can be one of clusterIsolationResponse, none, powerOff, or shutdown.",
				ValidateFunc: validation.StringInSlice(vmOverrideClusterDasVMSettingsIsolationResponseAllowedValues, false),
			},
			// VM component protection - datastore monitoring - Permanent Device Loss
			"ha_datastore_pdl_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault),
				Description:  "Controls the action to take on this virtual machine when the cluster has detected a permanent device loss to a relevant datastore. Can be one of clusterDefault, disabled, warning, or restartAggressive.",
				ValidateFunc: validation.StringInSlice(vmOverrideClusterVMStorageProtectionForPDLAllowedValues, false),
			},
			// VM component protection - datastore monitoring - All Paths Down
			"ha_datastore_apd_response": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault),
				Description:  "Controls the action to take on this virtual machine when the cluster has detected loss to all paths to a relevant datastore. Can be one of clusterDefault, disabled, warning, restartConservative, or restartAggressive.",
				ValidateFunc: validation.StringInSlice(vmOverrideClusterVMStorageProtectionForAPDAllowedValues, false),
			},
			"ha_datastore_apd_recovery_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedUseClusterDefault),
				Description:  "Controls the action to take on this virtual machine if an APD status on an affected datastore clears in the middle of an APD event. Can be one of useClusterDefault, none or reset.",
				ValidateFunc: validation.StringInSlice(vmOverrideClusterVMReactionOnAPDClearedAllowedValues, false),
			},
			"ha_datastore_apd_response_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Controls the delay in minutes to wait after an APD timeout event to execute the response action defined in ha_datastore_apd_response. Specify -1 to use the cluster setting.",
			},
			// VM monitoring
			"ha_vm_monitoring_use_cluster_defaults": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Determines whether or not the cluster's default settings or the VM override settings specified in this resource are used for virtual machine monitoring. The default is true (use cluster defaults) - set to false to have overrides take effect.",
			},
			"ha_vm_monitoring": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.ClusterDasConfigInfoVmMonitoringStateVmMonitoringDisabled),
				Description:  "The type of virtual machine monitoring to use for this virtual machine. Can be one of vmMonitoringDisabled, vmMonitoringOnly, or vmAndAppMonitoring.",
				ValidateFunc: validation.StringInSlice(clusterDasConfigInfoVMMonitoringStateAllowedValues, false),
			},
			"ha_vm_failure_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "If a heartbeat from this virtual machine is not received within this configured interval, the virtual machine is marked as failed. The value is in seconds.",
			},
			"ha_vm_minimum_uptime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "The time, in seconds, that HA waits after powering on this virtual machine before monitoring for heartbeats.",
			},
			"ha_vm_maximum_resets": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The maximum number of resets that HA will perform to this virtual machine when responding to a failure event.",
			},
			"ha_vm_maximum_failure_window": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The length of the reset window in which ha_vm_maximum_resets can operate. When this window expires, no more resets are attempted regardless of the setting configured in ha_vm_maximum_resets. -1 means no window, meaning an unlimited reset time is allotted.",
			},
		},
	}
}

func resourceVSphereHAVMOverrideCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereHAVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereHAVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDasVMConfigInfo(d, meta, vm)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DasVmConfigSpec: []types.ClusterDasVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: info,
			},
		},
	}

	if err = clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	id, err := resourceVSphereHAVMOverrideFlattenID(cluster, vm)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereHAVMOverrideIDString(d))
	return resourceVSphereHAVMOverrideRead(d, meta)
}

func resourceVSphereHAVMOverrideRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereHAVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereHAVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereHAVMOverrideFindEntry(cluster, vm)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the compute_cluster_id and virtual_machine_id here. These are
	// ForceNew, but we set these for completeness on import so that if the wrong
	// cluster/VM combo was used, it will be noted.
	if err = d.Set("compute_cluster_id", cluster.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"compute_cluster_id\": %s", err)
	}

	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error getting properties of virtual machine: %s", err)
	}
	if err = d.Set("virtual_machine_id", props.Config.Uuid); err != nil {
		return fmt.Errorf("error setting attribute \"virtual_machine_id\": %s", err)
	}

	if err = flattenClusterDasVMConfigInfo(d, meta, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereHAVMOverrideIDString(d))
	return nil
}

func resourceVSphereHAVMOverrideUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereHAVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereHAVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDasVMConfigInfo(d, meta, vm)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DasVmConfigSpec: []types.ClusterDasVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					// NOTE: Unlike other overrides, this needs to be an
					// ArrayUpdateOperationEdit, or else an "parameter incorrect" error
					// is given.
					Operation: types.ArrayUpdateOperationEdit,
				},
				Info: info,
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereHAVMOverrideIDString(d))
	return resourceVSphereHAVMOverrideRead(d, meta)
}

func resourceVSphereHAVMOverrideDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereHAVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereHAVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	spec := &types.ClusterConfigSpecEx{
		DasVmConfigSpec: []types.ClusterDasVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: vm.Reference(),
				},
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereHAVMOverrideIDString(d))
	return nil
}

func resourceVSphereHAVMOverrideImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	clusterPath, ok := data["compute_cluster_path"]
	if !ok {
		return nil, errors.New("missing compute_cluster_path in input data")
	}
	vmPath, ok := data["virtual_machine_path"]
	if !ok {
		return nil, errors.New("missing virtual_machine_path in input data")
	}

	client, err := resourceVSphereHAVMOverrideClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	vm, err := virtualmachine.FromPath(client, vmPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate virtual machine %q: %s", vmPath, err)
	}

	id, err := resourceVSphereHAVMOverrideFlattenID(cluster, vm)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterDasVMConfigInfo reads certain ResourceData keys and returns a
// ClusterDasVmConfigInfo.
func expandClusterDasVMConfigInfo(
	d *schema.ResourceData,
	meta interface{},
	vm *object.VirtualMachine,
) (*types.ClusterDasVmConfigInfo, error) {
	client, err := resourceVSphereHAVMOverrideClient(meta)
	if err != nil {
		return nil, err
	}
	version := viapi.ParseVersionFromClient(client)

	obj := &types.ClusterDasVmConfigInfo{
		DasSettings: expandClusterDasVMSettings(d, version),
		Key:         vm.Reference(),
	}

	// Expand ha_vm_monitoring_use_cluster_defaults here as it's not included in
	// the base vsphere_compute_cluster resource.
	obj.DasSettings.VmToolsMonitoringSettings.ClusterSettings = structure.GetBool(d, "ha_vm_monitoring_use_cluster_defaults")

	return obj, nil
}

// flattenClusterDasVmConfigInfo saves a ClusterDasVmConfigInfo into the
// supplied ResourceData.
func flattenClusterDasVMConfigInfo(d *schema.ResourceData, meta interface{}, obj *types.ClusterDasVmConfigInfo) error {
	client, err := resourceVSphereHAVMOverrideClient(meta)
	if err != nil {
		return err
	}
	version := viapi.ParseVersionFromClient(client)

	// Set ha_vm_monitoring_use_cluster_defaults here as it's not included in the
	// base vsphere_compute_cluster resource.
	if err := d.Set("ha_vm_monitoring_use_cluster_defaults", obj.DasSettings.VmToolsMonitoringSettings.ClusterSettings); err != nil {
		return err
	}
	return flattenClusterDasVMSettings(d, obj.DasSettings, version)
}

// resourceVSphereHAVMOverrideIDString prints a friendly string for the
// vsphere_ha_vm_override resource.
func resourceVSphereHAVMOverrideIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereHAVMOverrideName)
}

// resourceVSphereHAVMOverrideFlattenID makes an ID for the
// vsphere_ha_vm_override resource.
func resourceVSphereHAVMOverrideFlattenID(cluster *object.ClusterComputeResource, vm *object.VirtualMachine) (string, error) {
	clusterID := cluster.Reference().Value
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return "", fmt.Errorf("cannot compute ID off of properties of virtual machine: %s", err)
	}
	vmID := props.Config.Uuid
	return strings.Join([]string{clusterID, vmID}, ":"), nil
}

// resourceVSphereHAVMOverrideParseID parses an ID for the
// vsphere_ha_vm_override and outputs its parts.
func resourceVSphereHAVMOverrideParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereHAVMOverrideFindEntry attempts to locate an existing
// VM-specific HA config in a cluster's configuration. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereHAVMOverrideFindEntry(
	cluster *object.ClusterComputeResource,
	vm *object.VirtualMachine,
) (*types.ClusterDasVmConfigInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).DasVmConfig {
		if info.Key == vm.Reference() {
			log.Printf("[DEBUG] Found HA config info for VM %q in cluster %q", vm.Name(), cluster.Name())
			return &info, nil
		}
	}

	log.Printf("[DEBUG] No HA config info found for VM %q in cluster %q", vm.Name(), cluster.Name())
	return nil, nil
}

// resourceVSphereHAVMOverrideObjects handles the fetching of the cluster and
// virtual machine depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the compute_cluster_id and virtual_machine_id
// attributes.
func resourceVSphereHAVMOverrideObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	if d.Id() != "" {
		return resourceVSphereHAVMOverrideObjectsFromID(d, meta)
	}
	return resourceVSphereHAVMOverrideObjectsFromAttributes(d, meta)
}

func resourceVSphereHAVMOverrideObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	return resourceVSphereHAVMOverrideFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		d.Get("virtual_machine_id").(string),
	)
}

func resourceVSphereHAVMOverrideObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, vmID, err := resourceVSphereHAVMOverrideParseID(d.Id())
	if err != nil {
		return nil, nil, err
	}

	return resourceVSphereHAVMOverrideFetchObjects(meta, clusterID, vmID)
}

func resourceVSphereHAVMOverrideFetchObjects(
	meta interface{},
	clusterID string,
	vmID string,
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	client, err := resourceVSphereHAVMOverrideClient(meta)
	if err != nil {
		return nil, nil, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate cluster: %s", err)
	}

	vm, err := virtualmachine.FromUUID(client, vmID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate virtual machine: %s", err)
	}

	return cluster, vm, nil
}

func resourceVSphereHAVMOverrideClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
