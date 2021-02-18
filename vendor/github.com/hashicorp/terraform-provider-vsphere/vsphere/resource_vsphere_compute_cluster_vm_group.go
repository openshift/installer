package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterVMGroupName = "vsphere_compute_cluster_vm_group"

func resourceVSphereComputeClusterVMGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterVMGroupCreate,
		Read:   resourceVSphereComputeClusterVMGroupRead,
		Update: resourceVSphereComputeClusterVMGroupUpdate,
		Delete: resourceVSphereComputeClusterVMGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterVMGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique name of the virtual machine group in the cluster.",
			},
			"compute_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the cluster.",
			},
			"virtual_machine_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The UUIDs of the virtual machines in this group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVSphereComputeClusterVMGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterVMGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterVMGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterVMGroup(d, meta, name)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
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

	id, err := resourceVSphereComputeClusterVMGroupFlattenID(cluster, name)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterVMGroupIDString(d))
	return resourceVSphereComputeClusterVMGroupRead(d, meta)
}

func resourceVSphereComputeClusterVMGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterVMGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterVMGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterVMGroupFindEntry(cluster, name)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the compute_cluster_id and name here. These are
	// ForceNew, but we set these for completeness on import so that if the wrong
	// cluster/VM combo was used, it will be noted.
	if err = d.Set("compute_cluster_id", cluster.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"compute_cluster_id\": %s", err)
	}

	// This is the "correct" way to set name here, even if it's a bit
	// superfluous.
	if err = d.Set("name", info.Name); err != nil {
		return fmt.Errorf("error setting attribute \"name\": %s", err)
	}

	if err = flattenClusterVMGroup(d, meta, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterVMGroupIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterVMGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterVMGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterVMGroup(d, meta, name)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationEdit,
				},
				Info: info,
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterVMGroupIDString(d))
	return resourceVSphereComputeClusterVMGroupRead(d, meta)
}

func resourceVSphereComputeClusterVMGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterVMGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterVMGroupObjects(d, meta)
	if err != nil {
		return err
	}

	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: name,
				},
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterVMGroupIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	clusterPath, ok := data["compute_cluster_path"]
	if !ok {
		return nil, errors.New("missing compute_cluster_path in input data")
	}
	name, ok := data["name"]
	if !ok {
		return nil, errors.New("missing name in input data")
	}

	client, err := resourceVSphereComputeClusterVMGroupClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterVMGroupFindEntry(cluster, name)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("cluster group entry %q does not exist in cluster %q", name, cluster.Name())
	}

	id, err := resourceVSphereComputeClusterVMGroupFlattenID(cluster, name)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterVMGroup reads certain ResourceData keys and returns a
// ClusterVmGroup.
func expandClusterVMGroup(d *schema.ResourceData, meta interface{}, name string) (*types.ClusterVmGroup, error) {
	client, err := resourceVSphereComputeClusterVMGroupClient(meta)
	if err != nil {
		return nil, err
	}

	results, err := virtualmachine.MOIDsForUUIDs(
		client,
		structure.SliceInterfacesToStrings(d.Get("virtual_machine_ids").(*schema.Set).List()),
	)
	if err != nil {
		return nil, err
	}

	obj := &types.ClusterVmGroup{
		ClusterGroupInfo: types.ClusterGroupInfo{
			Name:        name,
			UserCreated: structure.BoolPtr(true),
		},
		Vm: results.ManagedObjectReferences(),
	}
	return obj, nil
}

// flattenClusterVmGroup saves a ClusterVmGroup into the supplied ResourceData.
func flattenClusterVMGroup(d *schema.ResourceData, meta interface{}, obj *types.ClusterVmGroup) error {
	client, err := resourceVSphereComputeClusterVMGroupClient(meta)
	if err != nil {
		return err
	}

	results, err := virtualmachine.UUIDsForManagedObjectReferences(
		client,
		obj.Vm,
	)
	if err != nil {
		return err
	}

	return structure.SetBatch(d, map[string]interface{}{
		"virtual_machine_ids": results.UUIDs(),
	})
}

// resourceVSphereComputeClusterVMGroupIDString prints a friendly string for the
// vsphere_cluster_vm_group resource.
func resourceVSphereComputeClusterVMGroupIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterVMGroupName)
}

// resourceVSphereComputeClusterVMGroupFlattenID makes an ID for the
// vsphere_cluster_vm_group resource.
func resourceVSphereComputeClusterVMGroupFlattenID(cluster *object.ClusterComputeResource, name string) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, name}, ":"), nil
}

// resourceVSphereComputeClusterVMGroupParseID parses an ID for the
// vsphere_cluster_vm_group and outputs its parts.
func resourceVSphereComputeClusterVMGroupParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereComputeClusterVMGroupFindEntry attempts to locate an existing
// VM group config in a cluster's configuration. It's used by the resource's
// read functionality and tests. nil is returned if the entry cannot be found.
func resourceVSphereComputeClusterVMGroupFindEntry(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterVmGroup, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Group {
		if info.GetClusterGroupInfo().Name == name {
			if vmInfo, ok := info.(*types.ClusterVmGroup); ok {
				log.Printf("[DEBUG] Found VM group %q in cluster %q", name, cluster.Name())
				return vmInfo, nil
			}
			return nil, fmt.Errorf("unique group name %q in cluster %q is not a VM group", name, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No VM group name %q found in cluster %q", name, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterVMGroupObjects handles the fetching of the cluster and
// group name depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the compute_cluster_id and name attributes.
func resourceVSphereComputeClusterVMGroupObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterVMGroupObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterVMGroupObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterVMGroupObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	return resourceVSphereComputeClusterVMGroupFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		d.Get("name").(string),
	)
}

func resourceVSphereComputeClusterVMGroupObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, name, err := resourceVSphereComputeClusterVMGroupParseID(d.Id())
	if err != nil {
		return nil, "", err
	}

	return resourceVSphereComputeClusterVMGroupFetchObjects(meta, clusterID, name)
}

// resourceVSphereComputeClusterVMGroupFetchObjects fetches the "objects" for a
// cluster VM group. This is currently just the cluster object as the name of
// the group is a static value and a pass-through - this is to keep its
// workflow consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterVMGroupFetchObjects(
	meta interface{},
	clusterID string,
	name string,
) (*object.ClusterComputeResource, string, error) {
	client, err := resourceVSphereComputeClusterVMGroupClient(meta)
	if err != nil {
		return nil, "", err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, "", fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, name, nil
}

func resourceVSphereComputeClusterVMGroupClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
