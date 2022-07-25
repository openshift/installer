package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/storagepod"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereStorageDrsVMOverrideName = "vsphere_storage_drs_vm_override"

func resourceVSphereStorageDrsVMOverride() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereStorageDrsVMOverrideCreate,
		Read:   resourceVSphereStorageDrsVMOverrideRead,
		Update: resourceVSphereStorageDrsVMOverrideUpdate,
		Delete: resourceVSphereStorageDrsVMOverrideDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereStorageDrsVMOverrideImport,
		},

		Schema: map[string]*schema.Schema{
			"datastore_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the datastore cluster.",
			},
			"virtual_machine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the virtual machine.",
			},
			"sdrs_enabled": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Overrides the default Storage DRS setting for this virtual machine.",
				ValidateFunc: structure.ValidateBoolStringPtr(),
				StateFunc:    structure.BoolStringPtrState,
			},
			"sdrs_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Overrides any Storage DRS automation levels for this virtual machine.",
				ValidateFunc: validation.StringInSlice(storageDrsPodConfigInfoBehaviorAllowedValues, false),
			},
			"sdrs_intra_vm_affinity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Overrides the intra-VM affinity setting for this virtual machine.",
				ValidateFunc: structure.ValidateBoolStringPtr(),
				StateFunc:    structure.BoolStringPtrState,
			},
		},
	}
}

func resourceVSphereStorageDrsVMOverrideCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereStorageDrsVMOverrideIDString(d))

	client, err := resourceVSphereStorageDrsVMOverrideClient(meta)
	if err != nil {
		return err
	}

	pod, vm, err := resourceVSphereStorageDrsVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandStorageDrsVMConfigInfo(d, vm)
	if err != nil {
		return err
	}
	spec := types.StorageDrsConfigSpec{
		VmConfigSpec: []types.StorageDrsVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: info,
			},
		},
	}

	if err = storagepod.ApplyDRSConfiguration(client, pod, spec); err != nil {
		return err
	}

	id, err := resourceVSphereStorageDrsVMOverrideFlattenID(pod, vm)
	if err != nil {
		return fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereStorageDrsVMOverrideIDString(d))
	return resourceVSphereStorageDrsVMOverrideRead(d, meta)
}

func resourceVSphereStorageDrsVMOverrideRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereStorageDrsVMOverrideIDString(d))

	pod, vm, err := resourceVSphereStorageDrsVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereStorageDrsVMOverrideFindEntry(pod, vm)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the datastore_cluster_id and virtual_machine_id here. These are
	// ForceNew, but we set these for completeness on import so that if the wrong
	// datastore cluster/VM combo was used, it will be noted.
	if err = d.Set("datastore_cluster_id", pod.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"datastore_cluster_id\": %s", err)
	}

	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error getting properties of virtual machine: %s", err)
	}
	if err = d.Set("virtual_machine_id", props.Config.Uuid); err != nil {
		return fmt.Errorf("error setting attribute \"virtual_machine_id\": %s", err)
	}

	if err = flattenStorageDrsVMConfigInfo(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereStorageDrsVMOverrideIDString(d))
	return nil
}

func resourceVSphereStorageDrsVMOverrideUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereStorageDrsVMOverrideIDString(d))

	client, err := resourceVSphereStorageDrsVMOverrideClient(meta)
	if err != nil {
		return err
	}

	pod, vm, err := resourceVSphereStorageDrsVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandStorageDrsVMConfigInfo(d, vm)
	if err != nil {
		return err
	}
	spec := types.StorageDrsConfigSpec{
		VmConfigSpec: []types.StorageDrsVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					// NOTE: ArrayUpdateOperationAdd here replaces existing entries,
					// versus adding duplicates. This is not documented in
					// StorageDrsVmConfigSpec but it is in the counterpart
					// ClusterDrsVmConfigSpec (used for compute DRS).
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: info,
			},
		},
	}

	if err := storagepod.ApplyDRSConfiguration(client, pod, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereStorageDrsVMOverrideIDString(d))
	return resourceVSphereStorageDrsVMOverrideRead(d, meta)
}

func resourceVSphereStorageDrsVMOverrideDelete(d *schema.ResourceData, meta interface{}) error {
	resourceIDString := resourceVSphereStorageDrsVMOverrideIDString(d)
	log.Printf("[DEBUG] %s: Beginning delete", resourceIDString)

	client, err := resourceVSphereStorageDrsVMOverrideClient(meta)
	if err != nil {
		return err
	}

	pod, vm, err := resourceVSphereStorageDrsVMOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	spec := types.StorageDrsConfigSpec{
		VmConfigSpec: []types.StorageDrsVmConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: vm.Reference(),
				},
			},
		},
	}

	if err := storagepod.ApplyDRSConfiguration(client, pod, spec); err != nil {
		return err
	}

	d.SetId("")

	log.Printf("[DEBUG] %s: Deleted successfully", resourceIDString)
	return nil
}

func resourceVSphereStorageDrsVMOverrideImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	podPath, ok := data["datastore_cluster_path"]
	if !ok {
		return nil, errors.New("missing datastore_cluster_path in input data")
	}
	vmPath, ok := data["virtual_machine_path"]
	if !ok {
		return nil, errors.New("missing virtual_machine_path in input data")
	}

	client, err := resourceVSphereStorageDrsVMOverrideClient(meta)
	if err != nil {
		return nil, err
	}

	pod, err := storagepod.FromPath(client, podPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate datastore cluster %q: %s", podPath, err)
	}

	vm, err := virtualmachine.FromPath(client, vmPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate virtual machine %q: %s", vmPath, err)
	}

	id, err := resourceVSphereStorageDrsVMOverrideFlattenID(pod, vm)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandStorageDrsVMConfigInfo reads certain ResourceData keys and returns a
// StorageDrsVmConfigInfo.
func expandStorageDrsVMConfigInfo(d *schema.ResourceData, vm *object.VirtualMachine) (*types.StorageDrsVmConfigInfo, error) {
	enabled, err := structure.GetBoolStringPtr(d, "sdrs_enabled")
	if err != nil {
		return nil, fmt.Errorf("error parsing field \"sdrs_enabled\": %s", err)
	}
	intraVMAffinity, err := structure.GetBoolStringPtr(d, "sdrs_intra_vm_affinity")
	if err != nil {
		return nil, fmt.Errorf("error parsing field \"sdrs_intra_vm_affinity\": %s", err)
	}

	obj := &types.StorageDrsVmConfigInfo{
		Behavior:        d.Get("sdrs_automation_level").(string),
		Enabled:         enabled,
		IntraVmAffinity: intraVMAffinity,
		Vm:              types.NewReference(vm.Reference()),
	}

	return obj, nil
}

// flattenStorageDrsVmConfigInfo saves a StorageDrsVmConfigInfo into the
// supplied ResourceData.
func flattenStorageDrsVMConfigInfo(d *schema.ResourceData, obj *types.StorageDrsVmConfigInfo) error {
	if err := d.Set("sdrs_automation_level", obj.Behavior); err != nil {
		return fmt.Errorf("error setting attribute \"sdrs_automation_level\": %s", err)
	}
	if err := structure.SetBoolStringPtr(d, "sdrs_enabled", obj.Enabled); err != nil {
		return fmt.Errorf("error setting attribute \"sdrs_enabled\": %s", err)
	}
	if err := structure.SetBoolStringPtr(d, "sdrs_intra_vm_affinity", obj.IntraVmAffinity); err != nil {
		return fmt.Errorf("error setting attribute \"sdrs_automation_level\": %s", err)
	}

	return nil
}

// resourceVSphereStorageDrsVMOverrideIDString prints a friendly string for the
// vsphere_storage_drs_vm_override resource.
func resourceVSphereStorageDrsVMOverrideIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereStorageDrsVMOverrideName)
}

// resourceVSphereStorageDrsVMOverrideFlattenID makes an ID for the
// vsphere_storage_drs_vm_override resource.
func resourceVSphereStorageDrsVMOverrideFlattenID(pod *object.StoragePod, vm *object.VirtualMachine) (string, error) {
	podID := pod.Reference().Value
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return "", fmt.Errorf("cannot compute ID off of properties of virtual machine: %s", err)
	}
	vmID := props.Config.Uuid
	return strings.Join([]string{podID, vmID}, ":"), nil
}

// resourceVSphereStorageDrsVMOverrideParseID parses an ID for the
// vsphere_storage_drs_vm_override and outputs its parts.
func resourceVSphereStorageDrsVMOverrideParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereStorageDrsVMOverrideFindEntry attempts to locate an existing VM
// config in a Storage Pod's DRS configuration. It's used by the resource's
// read functionality and tests. nil is returned if the entry cannot be found.
func resourceVSphereStorageDrsVMOverrideFindEntry(
	pod *object.StoragePod,
	vm *object.VirtualMachine,
) (*types.StorageDrsVmConfigInfo, error) {
	props, err := storagepod.Properties(pod)
	if err != nil {
		return nil, fmt.Errorf("error fetching datastore cluster properties: %s", err)
	}

	for _, info := range props.PodStorageDrsEntry.StorageDrsConfig.VmConfig {
		if *info.Vm == vm.Reference() {
			log.Printf("[DEBUG] Found storage DRS config info for VM %q in datastore cluster %q", vm.Name(), pod.Name())
			return &info, nil
		}
	}

	log.Printf("[DEBUG] No storage DRS config info found for VM %q in datastore cluster %q", vm.Name(), pod.Name())
	return nil, nil
}

// resourceVSphereStorageDrsVMOverrideObjects handles the fetching of the
// datastore cluster and virtual machine depending on what attributes are
// available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the datastore_cluster_id and virtual_machine_id
// attributes.
func resourceVSphereStorageDrsVMOverrideObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.StoragePod, *object.VirtualMachine, error) {
	if d.Id() != "" {
		return resourceVSphereStorageDrsVMOverrideObjectsFromID(d, meta)
	}
	return resourceVSphereStorageDrsVMOverrideObjectsFromAttributes(d, meta)
}

func resourceVSphereStorageDrsVMOverrideObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.StoragePod, *object.VirtualMachine, error) {
	return resourceVSphereStorageDrsVMOverrideFetchObjects(
		meta,
		d.Get("datastore_cluster_id").(string),
		d.Get("virtual_machine_id").(string),
	)
}

func resourceVSphereStorageDrsVMOverrideObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.StoragePod, *object.VirtualMachine, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	podID, vmID, err := resourceVSphereStorageDrsVMOverrideParseID(d.Id())
	if err != nil {
		return nil, nil, err
	}

	return resourceVSphereStorageDrsVMOverrideFetchObjects(meta, podID, vmID)
}

func resourceVSphereStorageDrsVMOverrideFetchObjects(
	meta interface{},
	podID string,
	vmID string,
) (*object.StoragePod, *object.VirtualMachine, error) {
	client, err := resourceVSphereStorageDrsVMOverrideClient(meta)
	if err != nil {
		return nil, nil, err
	}

	pod, err := storagepod.FromID(client, podID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate datastore cluster: %s", err)
	}

	vm, err := virtualmachine.FromUUID(client, vmID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate virtual machine: %s", err)
	}

	return pod, vm, nil
}

func resourceVSphereStorageDrsVMOverrideClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
