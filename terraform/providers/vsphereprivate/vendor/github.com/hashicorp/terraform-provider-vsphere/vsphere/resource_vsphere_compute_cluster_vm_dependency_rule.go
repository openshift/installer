package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterVMDependencyRuleName = "vsphere_compute_cluster_vm_dependency_rule"

func resourceVSphereComputeClusterVMDependencyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterVMDependencyRuleCreate,
		Read:   resourceVSphereComputeClusterVMDependencyRuleRead,
		Update: resourceVSphereComputeClusterVMDependencyRuleUpdate,
		Delete: resourceVSphereComputeClusterVMDependencyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterVMDependencyRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"compute_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the virtual machine group in the cluster.",
			},
			"dependency_vm_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the VM group that this rule depends on. The VMs defined in the group specified by vm_group_name will not be started until the VMs in this group are started.",
			},
			"vm_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the VM group that is the subject of this rule. The VMs defined in this group will not be started until the VMs in the group specified by dependency_vm_group_name are started.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable this rule in the cluster.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When true, prevents any virtual machine operations that may violate this rule.",
			},
		},
	}
}

func resourceVSphereComputeClusterVMDependencyRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterVMDependencyRuleIDString(d))

	cluster, _, err := resourceVSphereComputeClusterVMDependencyRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDependencyRuleInfo(d)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		RulesSpec: []types.ClusterRuleSpec{
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

	info, err = resourceVSphereComputeClusterVMDependencyRuleFindEntryByName(cluster, info.Name)
	if err != nil {
		return err
	}

	id, err := resourceVSphereComputeClusterVMDependencyRuleFlattenID(cluster, info.Key)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterVMDependencyRuleIDString(d))
	return resourceVSphereComputeClusterVMDependencyRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMDependencyRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterVMDependencyRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMDependencyRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterVMDependencyRuleFindEntry(cluster, key)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the compute_cluster_id. This is ForceNew, but we set these for
	// completeness on import so that if the wrong cluster/VM combo was used, it
	// will be noted.
	if err = d.Set("compute_cluster_id", cluster.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"compute_cluster_id\": %s", err)
	}

	if err = flattenClusterDependencyRuleInfo(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterVMDependencyRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMDependencyRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterVMDependencyRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMDependencyRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDependencyRuleInfo(d)
	if err != nil {
		return err
	}
	info.Key = key

	spec := &types.ClusterConfigSpecEx{
		RulesSpec: []types.ClusterRuleSpec{
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

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterVMDependencyRuleIDString(d))
	return resourceVSphereComputeClusterVMDependencyRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMDependencyRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterVMDependencyRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMDependencyRuleObjects(d, meta)
	if err != nil {
		return err
	}

	spec := &types.ClusterConfigSpecEx{
		RulesSpec: []types.ClusterRuleSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: key,
				},
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterVMDependencyRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMDependencyRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereComputeClusterVMDependencyRuleClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterVMDependencyRuleFindEntryByName(cluster, name)
	if err != nil {
		return nil, err
	}

	id, err := resourceVSphereComputeClusterVMDependencyRuleFlattenID(cluster, info.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterDependencyRuleInfo reads certain ResourceData keys and returns a
// ClusterDependencyRuleInfo.
func expandClusterDependencyRuleInfo(d *schema.ResourceData) (*types.ClusterDependencyRuleInfo, error) {
	obj := &types.ClusterDependencyRuleInfo{
		ClusterRuleInfo: types.ClusterRuleInfo{
			Enabled:     structure.GetBool(d, "enabled"),
			Mandatory:   structure.GetBool(d, "mandatory"),
			Name:        d.Get("name").(string),
			UserCreated: structure.BoolPtr(true),
		},
		DependsOnVmGroup: d.Get("dependency_vm_group_name").(string),
		VmGroup:          d.Get("vm_group_name").(string),
	}
	return obj, nil
}

// flattenClusterDependencyRuleInfo saves a ClusterDependencyRuleInfo into the supplied ResourceData.
func flattenClusterDependencyRuleInfo(d *schema.ResourceData, obj *types.ClusterDependencyRuleInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"enabled":                  obj.Enabled,
		"mandatory":                obj.Mandatory,
		"name":                     obj.Name,
		"dependency_vm_group_name": obj.DependsOnVmGroup,
		"vm_group_name":            obj.VmGroup,
	})
}

// resourceVSphereComputeClusterVMDependencyRuleIDString prints a friendly string for the
// vsphere_compute_cluster_vm_dependency_rule resource.
func resourceVSphereComputeClusterVMDependencyRuleIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterVMDependencyRuleName)
}

// resourceVSphereComputeClusterVMDependencyRuleFlattenID makes an ID for the
// vsphere_compute_cluster_vm_dependency_rule resource.
func resourceVSphereComputeClusterVMDependencyRuleFlattenID(cluster *object.ClusterComputeResource, key int32) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, strconv.Itoa(int(key))}, ":"), nil
}

// resourceVSphereComputeClusterVMDependencyRuleParseID parses an ID for the
// vsphere_compute_cluster_vm_dependency_rule and outputs its parts.
func resourceVSphereComputeClusterVMDependencyRuleParseID(id string) (string, int32, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", 0, fmt.Errorf("bad ID %q", id)
	}

	key, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("while converting key in ID %q to int32: %s", parts[1], err)
	}

	return parts[0], int32(key), nil
}

// resourceVSphereComputeClusterVMDependencyRuleFindEntry attempts to locate an
// existing VM dependency rule in a cluster's configuration by key. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereComputeClusterVMDependencyRuleFindEntry(
	cluster *object.ClusterComputeResource,
	key int32,
) (*types.ClusterDependencyRuleInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Key == key {
			if vmDependencyRuleInfo, ok := info.(*types.ClusterDependencyRuleInfo); ok {
				log.Printf("[DEBUG] Found VM dependency rule key %d in cluster %q", key, cluster.Name())
				return vmDependencyRuleInfo, nil
			}
			return nil, fmt.Errorf("rule key %d in cluster %q is not a VM dependency rule", key, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No VM dependency rule key %d found in cluster %q", key, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterVMDependencyRuleFindEntryByName attempts to locate an
// existing VM dependency rule in a cluster's configuration by name. It differs from
// the standard resourceVSphereComputeClusterVMDependencyRuleFindEntry in that we
// don't allow missing entries, as it's designed to be used in places where we
// don't want to allow for missing entries, such as during creation and import.
func resourceVSphereComputeClusterVMDependencyRuleFindEntryByName(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterDependencyRuleInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Name == name {
			if vmDependencyRuleInfo, ok := info.(*types.ClusterDependencyRuleInfo); ok {
				log.Printf("[DEBUG] Found VM dependency rule %q in cluster %q", name, cluster.Name())
				return vmDependencyRuleInfo, nil
			}
			return nil, fmt.Errorf("rule %q in cluster %q is not a VM dependency rule", name, cluster.Name())
		}
	}

	return nil, fmt.Errorf("no VM dependency rule %q found in cluster %q", name, cluster.Name())
}

// resourceVSphereComputeClusterVMDependencyRuleObjects handles the fetching of the
// cluster and rule key depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, only the cluster is retrieved from compute_cluster_id. -1 is
// returned for the key.
func resourceVSphereComputeClusterVMDependencyRuleObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterVMDependencyRuleObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterVMDependencyRuleObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterVMDependencyRuleObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	return resourceVSphereComputeClusterVMDependencyRuleFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		-1,
	)
}

func resourceVSphereComputeClusterVMDependencyRuleObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, key, err := resourceVSphereComputeClusterVMDependencyRuleParseID(d.Id())
	if err != nil {
		return nil, 0, err
	}

	return resourceVSphereComputeClusterVMDependencyRuleFetchObjects(meta, clusterID, key)
}

// resourceVSphereComputeClusterVMDependencyRuleFetchObjects fetches the "objects"
// for a cluster rule. This is currently just the cluster object as the rule
// key a static value and a pass-through - this is to keep its workflow
// consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterVMDependencyRuleFetchObjects(
	meta interface{},
	clusterID string,
	key int32,
) (*object.ClusterComputeResource, int32, error) {
	client, err := resourceVSphereComputeClusterVMDependencyRuleClient(meta)
	if err != nil {
		return nil, 0, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, key, nil
}

func resourceVSphereComputeClusterVMDependencyRuleClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
