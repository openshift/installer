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
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterVMAntiAffinityRuleName = "vsphere_compute_cluster_vm_anti_affinity_rule"

func resourceVSphereComputeClusterVMAntiAffinityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereComputeClusterVMAntiAffinityRuleCreate,
		Read:   resourceVSphereComputeClusterVMAntiAffinityRuleRead,
		Update: resourceVSphereComputeClusterVMAntiAffinityRuleUpdate,
		Delete: resourceVSphereComputeClusterVMAntiAffinityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterVMAntiAffinityRuleImport,
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
				Description: "The UUIDs of the virtual machines to run on hosts different from each other.",
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

func resourceVSphereComputeClusterVMAntiAffinityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))

	cluster, _, err := resourceVSphereComputeClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAntiAffinityRuleSpec(d, meta)
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

	info, err = resourceVSphereComputeClusterVMAntiAffinityRuleFindEntryByName(cluster, info.Name)
	if err != nil {
		return err
	}

	id, err := resourceVSphereComputeClusterVMAntiAffinityRuleFlattenID(cluster, info.Key)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))
	return resourceVSphereComputeClusterVMAntiAffinityRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMAntiAffinityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterVMAntiAffinityRuleFindEntry(cluster, key)
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

	if err = flattenClusterAntiAffinityRuleSpec(d, meta, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMAntiAffinityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAntiAffinityRuleSpec(d, meta)
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

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))
	return resourceVSphereComputeClusterVMAntiAffinityRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMAntiAffinityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMAntiAffinityRuleObjects(d, meta)
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

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMAntiAffinityRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereComputeClusterVMAntiAffinityRuleClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterVMAntiAffinityRuleFindEntryByName(cluster, name)
	if err != nil {
		return nil, err
	}

	id, err := resourceVSphereComputeClusterVMAntiAffinityRuleFlattenID(cluster, info.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterAntiAffinityRuleSpec reads certain ResourceData keys and returns a
// ClusterAntiAffinityRuleSpec.
func expandClusterAntiAffinityRuleSpec(d *schema.ResourceData, meta interface{}) (*types.ClusterAntiAffinityRuleSpec, error) {
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

	obj := &types.ClusterAntiAffinityRuleSpec{
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

// flattenClusterAntiAffinityRuleSpec saves a ClusterAntiAffinityRuleSpec into the supplied ResourceData.
func flattenClusterAntiAffinityRuleSpec(d *schema.ResourceData, meta interface{}, obj *types.ClusterAntiAffinityRuleSpec) error {
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

// resourceVSphereComputeClusterVMAntiAffinityRuleIDString prints a friendly string for the
// vsphere_compute_cluster_vm_anti_affinity_rule resource.
func resourceVSphereComputeClusterVMAntiAffinityRuleIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterVMAntiAffinityRuleName)
}

// resourceVSphereComputeClusterVMAntiAffinityRuleFlattenID makes an ID for the
// vsphere_compute_cluster_vm_anti_affinity_rule resource.
func resourceVSphereComputeClusterVMAntiAffinityRuleFlattenID(cluster *object.ClusterComputeResource, key int32) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, strconv.Itoa(int(key))}, ":"), nil
}

// resourceVSphereComputeClusterVMAntiAffinityRuleParseID parses an ID for the
// vsphere_compute_cluster_vm_anti_affinity_rule and outputs its parts.
func resourceVSphereComputeClusterVMAntiAffinityRuleParseID(id string) (string, int32, error) {
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

// resourceVSphereComputeClusterVMAntiAffinityRuleFindEntry attempts to locate an
// existing VM anti-affinity rule in a cluster's configuration by key. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereComputeClusterVMAntiAffinityRuleFindEntry(
	cluster *object.ClusterComputeResource,
	key int32,
) (*types.ClusterAntiAffinityRuleSpec, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Key == key {
			if vmAffinityRuleInfo, ok := info.(*types.ClusterAntiAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM anti-affinity rule key %d in cluster %q", key, cluster.Name())
				return vmAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule key %d in cluster %q is not a VM anti-affinity rule", key, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No VM anti-affinity rule key %d found in cluster %q", key, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterVMAntiAffinityRuleFindEntryByName attempts to locate an
// existing VM anti-affinity rule in a cluster's configuration by name. It differs from
// the standard resourceVSphereComputeClusterVMAntiAffinityRuleFindEntry in that we
// don't allow missing entries, as it's designed to be used in places where we
// don't want to allow for missing entries, such as during creation and import.
func resourceVSphereComputeClusterVMAntiAffinityRuleFindEntryByName(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterAntiAffinityRuleSpec, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Name == name {
			if vmAffinityRuleInfo, ok := info.(*types.ClusterAntiAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM anti-affinity rule %q in cluster %q", name, cluster.Name())
				return vmAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule %q in cluster %q is not a VM anti-affinity rule", name, cluster.Name())
		}
	}

	return nil, fmt.Errorf("no VM anti-affinity rule %q found in cluster %q", name, cluster.Name())
}

// resourceVSphereComputeClusterVMAntiAffinityRuleObjects handles the fetching of the
// cluster and rule key depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, only the cluster is retrieved from compute_cluster_id. -1 is
// returned for the key.
func resourceVSphereComputeClusterVMAntiAffinityRuleObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterVMAntiAffinityRuleObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterVMAntiAffinityRuleObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterVMAntiAffinityRuleObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	return resourceVSphereComputeClusterVMAntiAffinityRuleFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		-1,
	)
}

func resourceVSphereComputeClusterVMAntiAffinityRuleObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, key, err := resourceVSphereComputeClusterVMAntiAffinityRuleParseID(d.Id())
	if err != nil {
		return nil, 0, err
	}

	return resourceVSphereComputeClusterVMAntiAffinityRuleFetchObjects(meta, clusterID, key)
}

// resourceVSphereComputeClusterVMAntiAffinityRuleFetchObjects fetches the "objects"
// for a cluster rule. This is currently just the cluster object as the rule
// key a static value and a pass-through - this is to keep its workflow
// consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterVMAntiAffinityRuleFetchObjects(
	meta interface{},
	clusterID string,
	key int32,
) (*object.ClusterComputeResource, int32, error) {
	client, err := resourceVSphereComputeClusterVMAntiAffinityRuleClient(meta)
	if err != nil {
		return nil, 0, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, key, nil
}

func resourceVSphereComputeClusterVMAntiAffinityRuleClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
