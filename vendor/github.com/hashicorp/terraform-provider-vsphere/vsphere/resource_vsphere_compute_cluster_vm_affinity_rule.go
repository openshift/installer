package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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

const resourceVSphereComputeClusterVMAffinityRuleName = "vsphere_compute_cluster_vm_affinity_rule"

func resourceVSphereComputeClusterVMAffinityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterVMAffinityRuleCreate,
		Read:   resourceVSphereComputeClusterVMAffinityRuleRead,
		Update: resourceVSphereComputeClusterVMAffinityRuleUpdate,
		Delete: resourceVSphereComputeClusterVMAffinityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterVMAffinityRuleImport,
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
			"virtual_machine_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The UUIDs of the virtual machines to run on the same host together.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceVSphereComputeClusterVMAffinityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterVMAffinityRuleIDString(d))

	cluster, _, err := resourceVSphereComputeClusterVMAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAffinityRuleSpec(d, meta)
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

	info, err = resourceVSphereComputeClusterVMAffinityRuleFindEntryByName(cluster, info.Name)
	if err != nil {
		return err
	}

	id, err := resourceVSphereComputeClusterVMAffinityRuleFlattenID(cluster, info.Key)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterVMAffinityRuleIDString(d))
	return resourceVSphereComputeClusterVMAffinityRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMAffinityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterVMAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterVMAffinityRuleFindEntry(cluster, key)
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

	if err = flattenClusterAffinityRuleSpec(d, meta, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterVMAffinityRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMAffinityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterVMAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAffinityRuleSpec(d, meta)
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

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterVMAffinityRuleIDString(d))
	return resourceVSphereComputeClusterVMAffinityRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMAffinityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterVMAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAffinityRuleObjects(d, meta)
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

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterVMAffinityRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMAffinityRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereComputeClusterVMAffinityRuleClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterVMAffinityRuleFindEntryByName(cluster, name)
	if err != nil {
		return nil, err
	}

	id, err := resourceVSphereComputeClusterVMAffinityRuleFlattenID(cluster, info.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterAffinityRuleSpec reads certain ResourceData keys and returns a
// ClusterAffinityRuleSpec.
func expandClusterAffinityRuleSpec(d *schema.ResourceData, meta interface{}) (*types.ClusterAffinityRuleSpec, error) {
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

	obj := &types.ClusterAffinityRuleSpec{
		ClusterRuleInfo: types.ClusterRuleInfo{
			Enabled:     structure.GetBool(d, "enabled"),
			Mandatory:   structure.GetBool(d, "mandatory"),
			Name:        d.Get("name").(string),
			UserCreated: structure.BoolPtr(true),
		},
		Vm: results.ManagedObjectReferences(),
	}
	return obj, nil
}

// flattenClusterAffinityRuleSpec saves a ClusterAffinityRuleSpec into the supplied ResourceData.
func flattenClusterAffinityRuleSpec(d *schema.ResourceData, meta interface{}, obj *types.ClusterAffinityRuleSpec) error {
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
		"enabled":             obj.Enabled,
		"mandatory":           obj.Mandatory,
		"name":                obj.Name,
		"virtual_machine_ids": results.UUIDs(),
	})
}

// resourceVSphereComputeClusterVMAffinityRuleIDString prints a friendly string for the
// vsphere_compute_cluster_vm_affinity_rule resource.
func resourceVSphereComputeClusterVMAffinityRuleIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterVMAffinityRuleName)
}

// resourceVSphereComputeClusterVMAffinityRuleFlattenID makes an ID for the
// vsphere_compute_cluster_vm_affinity_rule resource.
func resourceVSphereComputeClusterVMAffinityRuleFlattenID(cluster *object.ClusterComputeResource, key int32) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, strconv.Itoa(int(key))}, ":"), nil
}

// resourceVSphereComputeClusterVMAffinityRuleParseID parses an ID for the
// vsphere_compute_cluster_vm_affinity_rule and outputs its parts.
func resourceVSphereComputeClusterVMAffinityRuleParseID(id string) (string, int32, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", 0, fmt.Errorf("bad ID %q", id)
	}

	key, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("bad key in ID %q: %s", parts[1], err)
	}

	return parts[0], int32(key), nil
}

// resourceVSphereComputeClusterVMAffinityRuleFindEntry attempts to locate an
// existing VM affinity rule in a cluster's configuration by key. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereComputeClusterVMAffinityRuleFindEntry(
	cluster *object.ClusterComputeResource,
	key int32,
) (*types.ClusterAffinityRuleSpec, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Key == key {
			if vmAffinityRuleInfo, ok := info.(*types.ClusterAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM affinity rule key %d in cluster %q", key, cluster.Name())
				return vmAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule key %d in cluster %q is not a VM affinity rule", key, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No VM affinity rule key %d found in cluster %q", key, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterVMAffinityRuleFindEntryByName attempts to locate an
// existing VM affinity rule in a cluster's configuration by name. It differs from
// the standard resourceVSphereComputeClusterVMAffinityRuleFindEntry in that we
// don't allow missing entries, as it's designed to be used in places where we
// don't want to allow for missing entries, such as during creation and import.
func resourceVSphereComputeClusterVMAffinityRuleFindEntryByName(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterAffinityRuleSpec, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Name == name {
			if vmAffinityRuleInfo, ok := info.(*types.ClusterAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM affinity rule %q in cluster %q", name, cluster.Name())
				return vmAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule %q in cluster %q is not a VM affinity rule", name, cluster.Name())
		}
	}

	return nil, fmt.Errorf("no VM affinity rule %q found in cluster %q", name, cluster.Name())
}

// resourceVSphereComputeClusterVMAffinityRuleObjects handles the fetching of the
// cluster and rule key depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, only the cluster is retrieved from compute_cluster_id. -1 is
// returned for the key.
func resourceVSphereComputeClusterVMAffinityRuleObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterVMAffinityRuleObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterVMAffinityRuleObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterVMAffinityRuleObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	return resourceVSphereComputeClusterVMAffinityRuleFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		-1,
	)
}

func resourceVSphereComputeClusterVMAffinityRuleObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, key, err := resourceVSphereComputeClusterVMAffinityRuleParseID(d.Id())
	if err != nil {
		return nil, 0, err
	}

	return resourceVSphereComputeClusterVMAffinityRuleFetchObjects(meta, clusterID, key)
}

// resourceVSphereComputeClusterVMAffinityRuleFetchObjects fetches the "objects"
// for a cluster rule. This is currently just the cluster object as the rule
// key a static value and a pass-through - this is to keep its workflow
// consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterVMAffinityRuleFetchObjects(
	meta interface{},
	clusterID string,
	key int32,
) (*object.ClusterComputeResource, int32, error) {
	client, err := resourceVSphereComputeClusterVMAffinityRuleClient(meta)
	if err != nil {
		return nil, 0, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, key, nil
}

func resourceVSphereComputeClusterVMAffinityRuleClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
