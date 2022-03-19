package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereDRSVMOverrideName = "vsphere_drs_vm_override"

func resourceVSphereDRSVMOverride() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDRSVMOverrideCreate,
		Read:   resourceVSphereDRSVMOverrideRead,
		Update: resourceVSphereDRSVMOverrideUpdate,
		Delete: resourceVSphereDRSVMOverrideDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereDRSVMOverrideImport,
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
			"drs_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DRS for this virtual machine.",
			},
			"drs_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.DrsBehaviorManual),
				Description:  "The automation level for this virtual machine in the cluster. Can be one of manual, partiallyAutomated, or fullyAutomated.",
				ValidateFunc: validation.StringInSlice(drsBehaviorAllowedValues, false),
			},
		},
	}
}

func resourceVSphereDRSVMOverrideCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereDRSVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereDRSVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDrsVMConfigInfo(d, vm)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DrsVmConfigSpec: []types.ClusterDrsVmConfigSpec{
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

	id, err := resourceVSphereDRSVMOverrideFlattenID(cluster, vm)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereDRSVMOverrideIDString(d))
	return resourceVSphereDRSVMOverrideRead(d, meta)
}

func resourceVSphereDRSVMOverrideRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereDRSVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereDRSVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereDRSVMOverrideFindEntry(cluster, vm)
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

	if err = flattenClusterDrsVMConfigInfo(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereDRSVMOverrideIDString(d))
	return nil
}

func resourceVSphereDRSVMOverrideUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereDRSVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereDRSVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDrsVMConfigInfo(d, vm)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DrsVmConfigSpec: []types.ClusterDrsVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					// NOTE: ArrayUpdateOperationAdd here replaces existing entries,
					// versus adding duplicates or "merging" old settings with new ones
					// that have missing fields.
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: info,
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereDRSVMOverrideIDString(d))
	return resourceVSphereDRSVMOverrideRead(d, meta)
}

func resourceVSphereDRSVMOverrideDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereDRSVMOverrideIDString(d))

	cluster, vm, err := resourceVSphereDRSVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	spec := &types.ClusterConfigSpecEx{
		DrsVmConfigSpec: []types.ClusterDrsVmConfigSpec{
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

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereDRSVMOverrideIDString(d))
	return nil
}

func resourceVSphereDRSVMOverrideImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereDRSVMOverrideClient(meta)
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

	id, err := resourceVSphereDRSVMOverrideFlattenID(cluster, vm)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterDrsVMConfigInfo reads certain ResourceData keys and returns a
// ClusterDrsVmConfigInfo.
func expandClusterDrsVMConfigInfo(d *schema.ResourceData, vm *object.VirtualMachine) (*types.ClusterDrsVmConfigInfo, error) {
	obj := &types.ClusterDrsVmConfigInfo{
		Behavior: types.DrsBehavior(d.Get("drs_automation_level").(string)),
		Enabled:  structure.GetBool(d, "drs_enabled"),
		Key:      vm.Reference(),
	}

	return obj, nil
}

// flattenClusterDrsVmConfigInfo saves a ClusterDrsVmConfigInfo into the
// supplied ResourceData.
func flattenClusterDrsVMConfigInfo(d *schema.ResourceData, obj *types.ClusterDrsVmConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"drs_automation_level": obj.Behavior,
		"drs_enabled":          obj.Enabled,
	})
}

// resourceVSphereDRSVMOverrideIDString prints a friendly string for the
// vsphere_drs_vm_override resource.
func resourceVSphereDRSVMOverrideIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereDRSVMOverrideName)
}

// resourceVSphereDRSVMOverrideFlattenID makes an ID for the
// vsphere_drs_vm_override resource.
func resourceVSphereDRSVMOverrideFlattenID(cluster *object.ClusterComputeResource, vm *object.VirtualMachine) (string, error) {
	clusterID := cluster.Reference().Value
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return "", fmt.Errorf("cannot compute ID off of properties of virtual machine: %s", err)
	}
	vmID := props.Config.Uuid
	return strings.Join([]string{clusterID, vmID}, ":"), nil
}

// resourceVSphereDRSVMOverrideParseID parses an ID for the
// vsphere_drs_vm_override and outputs its parts.
func resourceVSphereDRSVMOverrideParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereDRSVMOverrideFindEntry attempts to locate an existing DRS VM
// config in a cluster's configuration. It's used by the resource's read
// functionality and tests. nil is returned if the entry cannot be found.
func resourceVSphereDRSVMOverrideFindEntry(
	cluster *object.ClusterComputeResource,
	vm *object.VirtualMachine,
) (*types.ClusterDrsVmConfigInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).DrsVmConfig {
		if info.Key == vm.Reference() {
			log.Printf("[DEBUG] Found DRS config info for VM %q in cluster %q", vm.Name(), cluster.Name())
			return &info, nil
		}
	}

	log.Printf("[DEBUG] No DRS config info found for VM %q in cluster %q", vm.Name(), cluster.Name())
	return nil, nil
}

// resourceVSphereDRSVMOverrideObjects handles the fetching of the cluster and
// virtual machine depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the compute_cluster_id and virtual_machine_id
// attributes.
func resourceVSphereDRSVMOverrideObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	if d.Id() != "" {
		return resourceVSphereDRSVMOverrideObjectsFromID(d, meta)
	}
	return resourceVSphereDRSVMOverrideObjectsFromAttributes(d, meta)
}

func resourceVSphereDRSVMOverrideObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	return resourceVSphereDRSVMOverrideFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		d.Get("virtual_machine_id").(string),
	)
}

func resourceVSphereDRSVMOverrideObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, vmID, err := resourceVSphereDRSVMOverrideParseID(d.Id())
	if err != nil {
		return nil, nil, err
	}

	return resourceVSphereDRSVMOverrideFetchObjects(meta, clusterID, vmID)
}

func resourceVSphereDRSVMOverrideFetchObjects(
	meta interface{},
	clusterID string,
	vmID string,
) (*object.ClusterComputeResource, *object.VirtualMachine, error) {
	client, err := resourceVSphereDRSVMOverrideClient(meta)
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

func resourceVSphereDRSVMOverrideClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
