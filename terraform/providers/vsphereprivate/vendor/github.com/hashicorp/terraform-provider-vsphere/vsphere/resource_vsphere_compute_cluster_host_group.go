package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterHostGroupName = "vsphere_compute_cluster_host_group"

func resourceVSphereComputeClusterHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterHostGroupCreate,
		Read:   resourceVSphereComputeClusterHostGroupRead,
		Update: resourceVSphereComputeClusterHostGroupUpdate,
		Delete: resourceVSphereComputeClusterHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterHostGroupImport,
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
			"host_system_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The managed object IDs of the hosts.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVSphereComputeClusterHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterHostGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterHostGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterHostGroup(d, name)
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

	id, err := resourceVSphereComputeClusterHostGroupFlattenID(cluster, name)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterHostGroupIDString(d))
	return resourceVSphereComputeClusterHostGroupRead(d, meta)
}

func resourceVSphereComputeClusterHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterHostGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterHostGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterHostGroupFindEntry(cluster, name)
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

	if err = flattenClusterHostGroup(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterHostGroupIDString(d))
	return nil
}

func resourceVSphereComputeClusterHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterHostGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterHostGroupObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterHostGroup(d, name)
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

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterHostGroupIDString(d))
	return resourceVSphereComputeClusterHostGroupRead(d, meta)
}

func resourceVSphereComputeClusterHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterHostGroupIDString(d))

	cluster, name, err := resourceVSphereComputeClusterHostGroupObjects(d, meta)
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

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterHostGroupIDString(d))
	return nil
}

func resourceVSphereComputeClusterHostGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereComputeClusterHostGroupClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterHostGroupFindEntry(cluster, name)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("cluster group entry %q does not exist in cluster %q", name, cluster.Name())
	}

	id, err := resourceVSphereComputeClusterHostGroupFlattenID(cluster, name)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterHostGroup reads certain ResourceData keys and returns a
// ClusterHostGroup.
func expandClusterHostGroup(d *schema.ResourceData, name string) (*types.ClusterHostGroup, error) {
	obj := &types.ClusterHostGroup{
		ClusterGroupInfo: types.ClusterGroupInfo{
			Name:        name,
			UserCreated: structure.BoolPtr(true),
		},
		Host: structure.SliceInterfacesToManagedObjectReferences(d.Get("host_system_ids").(*schema.Set).List(), "HostSystem"),
	}
	return obj, nil
}

// flattenClusterHostGroup saves a ClusterHostGroup into the supplied ResourceData.
func flattenClusterHostGroup(d *schema.ResourceData, obj *types.ClusterHostGroup) error {
	var hostIDs []string
	for _, v := range obj.Host {
		hostIDs = append(hostIDs, v.Value)
	}

	return structure.SetBatch(d, map[string]interface{}{
		"host_system_ids": hostIDs,
	})
}

// resourceVSphereComputeClusterHostGroupIDString prints a friendly string for the
// vsphere_cluster_host_group resource.
func resourceVSphereComputeClusterHostGroupIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterHostGroupName)
}

// resourceVSphereComputeClusterHostGroupFlattenID makes an ID for the
// vsphere_cluster_host_group resource.
func resourceVSphereComputeClusterHostGroupFlattenID(cluster *object.ClusterComputeResource, name string) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, name}, ":"), nil
}

// resourceVSphereComputeClusterHostGroupParseID parses an ID for the
// vsphere_cluster_host_group and outputs its parts.
func resourceVSphereComputeClusterHostGroupParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereComputeClusterHostGroupFindEntry attempts to locate an
// existing host group in a cluster's configuration. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereComputeClusterHostGroupFindEntry(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterHostGroup, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Group {
		if info.GetClusterGroupInfo().Name == name {
			if hostInfo, ok := info.(*types.ClusterHostGroup); ok {
				log.Printf("[DEBUG] Found host group %q in cluster %q", name, cluster.Name())
				return hostInfo, nil
			}
			return nil, fmt.Errorf("unique group name %q in cluster %q is not a host group", name, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No host group name %q found in cluster %q", name, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterHostGroupObjects handles the fetching of the
// cluster and group name depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the compute_cluster_id and name attributes.
func resourceVSphereComputeClusterHostGroupObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterHostGroupObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterHostGroupObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterHostGroupObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	return resourceVSphereComputeClusterHostGroupFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		d.Get("name").(string),
	)
}

func resourceVSphereComputeClusterHostGroupObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, string, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, name, err := resourceVSphereComputeClusterHostGroupParseID(d.Id())
	if err != nil {
		return nil, "", err
	}

	return resourceVSphereComputeClusterHostGroupFetchObjects(meta, clusterID, name)
}

// resourceVSphereComputeClusterHostGroupFetchObjects fetches the "objects" for
// a cluster host group. This is currently just the cluster object as the name
// of the group is a static value and a pass-through - this is to keep its
// workflow consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterHostGroupFetchObjects(
	meta interface{},
	clusterID string,
	name string,
) (*object.ClusterComputeResource, string, error) {
	client, err := resourceVSphereComputeClusterHostGroupClient(meta)
	if err != nil {
		return nil, "", err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, "", fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, name, nil
}

func resourceVSphereComputeClusterHostGroupClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
