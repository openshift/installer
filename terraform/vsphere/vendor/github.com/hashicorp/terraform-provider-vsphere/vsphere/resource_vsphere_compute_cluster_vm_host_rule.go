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
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereComputeClusterVMHostRuleName = "vsphere_compute_cluster_vm_host_rule"

func resourceVSphereComputeClusterVMHostRule() *schema.Resource {
	return &schema.Resource{
		Create:        resourceVSphereComputeClusterVMHostRuleCreate,
		Read:          resourceVSphereComputeClusterVMHostRuleRead,
		Update:        resourceVSphereComputeClusterVMHostRuleUpdate,
		Delete:        resourceVSphereComputeClusterVMHostRuleDelete,
		CustomizeDiff: resourceVSphereComputeClusterVMHostRuleCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereComputeClusterVMHostRuleImport,
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
			"vm_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the virtual machine group to use with this rule.",
			},
			"affinity_host_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"anti_affinity_host_group_name"},
				Description:   "When this field is used, virtual machines defined in vm_group_name will be run on the hosts defined in this host group.",
			},
			"anti_affinity_host_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"affinity_host_group_name"},
				Description:   "When this field is used, virtual machines defined in vm_group_name will not be run on the hosts defined in this host group.",
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

func resourceVSphereComputeClusterVMHostRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereComputeClusterVMHostRuleIDString(d))

	cluster, _, err := resourceVSphereComputeClusterVMHostRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterVMHostRuleInfo(d)
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

	info, err = resourceVSphereComputeClusterVMHostRuleFindEntryByName(cluster, info.Name)
	if err != nil {
		return err
	}

	id, err := resourceVSphereComputeClusterVMHostRuleFlattenID(cluster, info.Key)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereComputeClusterVMHostRuleIDString(d))
	return resourceVSphereComputeClusterVMHostRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMHostRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereComputeClusterVMHostRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMHostRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereComputeClusterVMHostRuleFindEntry(cluster, key)
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

	if err = flattenClusterVMHostRuleInfo(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereComputeClusterVMHostRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMHostRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereComputeClusterVMHostRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMHostRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterVMHostRuleInfo(d)
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

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereComputeClusterVMHostRuleIDString(d))
	return resourceVSphereComputeClusterVMHostRuleRead(d, meta)
}

func resourceVSphereComputeClusterVMHostRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereComputeClusterVMHostRuleIDString(d))

	cluster, key, err := resourceVSphereComputeClusterVMHostRuleObjects(d, meta)
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

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereComputeClusterVMHostRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMHostRuleCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning diff customization and validation", resourceVSphereComputeClusterVMHostRuleIDString(d))

	if err := resourceVSphereComputeClusterVMHostRuleValidateHostRulesSpecified(d); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Diff customization and validation complete", resourceVSphereComputeClusterVMHostRuleIDString(d))
	return nil
}

func resourceVSphereComputeClusterVMHostRuleValidateHostRulesSpecified(d *schema.ResourceDiff) error {
	log.Printf(
		"[DEBUG] %s: Validating presence of one of affinity_host_group_name or anti_affinity_host_group_name",
		resourceVSphereComputeClusterVMHostRuleIDString(d),
	)
	_, affinityOk := d.GetOk("affinity_host_group_name")
	affinityKnown := d.NewValueKnown("affinity_host_group_name")
	_, antiOk := d.GetOk("anti_affinity_host_group_name")
	antiKnown := d.NewValueKnown("anti_affinity_host_group_name")

	if !affinityOk && affinityKnown && !antiOk && antiKnown {
		return errors.New("one of affinity_host_group_name or anti_affinity_host_group_name must be specified")
	}

	return nil
}

func resourceVSphereComputeClusterVMHostRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

	client, err := resourceVSphereComputeClusterVMHostRuleClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	info, err := resourceVSphereComputeClusterVMHostRuleFindEntryByName(cluster, name)
	if err != nil {
		return nil, err
	}

	id, err := resourceVSphereComputeClusterVMHostRuleFlattenID(cluster, info.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterVMHostRuleInfo reads certain ResourceData keys and returns a
// ClusterVmHostRuleInfo.
func expandClusterVMHostRuleInfo(d *schema.ResourceData) (*types.ClusterVmHostRuleInfo, error) {
	obj := &types.ClusterVmHostRuleInfo{
		ClusterRuleInfo: types.ClusterRuleInfo{
			Enabled:     structure.GetBool(d, "enabled"),
			Mandatory:   structure.GetBool(d, "mandatory"),
			Name:        d.Get("name").(string),
			UserCreated: structure.BoolPtr(true),
		},
		AffineHostGroupName:     d.Get("affinity_host_group_name").(string),
		AntiAffineHostGroupName: d.Get("anti_affinity_host_group_name").(string),
		VmGroupName:             d.Get("vm_group_name").(string),
	}
	return obj, nil
}

// flattenClusterVMHostRuleInfo saves a ClusterVmHostRuleInfo into the supplied ResourceData.
func flattenClusterVMHostRuleInfo(d *schema.ResourceData, obj *types.ClusterVmHostRuleInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"enabled":                       obj.Enabled,
		"mandatory":                     obj.Mandatory,
		"name":                          obj.Name,
		"affinity_host_group_name":      obj.AffineHostGroupName,
		"anti_affinity_host_group_name": obj.AntiAffineHostGroupName,
		"vm_group_name":                 obj.VmGroupName,
	})
}

// resourceVSphereComputeClusterVMHostRuleIDString prints a friendly string for the
// vsphere_compute_cluster_vm_host_rule resource.
func resourceVSphereComputeClusterVMHostRuleIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereComputeClusterVMHostRuleName)
}

// resourceVSphereComputeClusterVMHostRuleFlattenID makes an ID for the
// vsphere_compute_cluster_vm_host_rule resource.
func resourceVSphereComputeClusterVMHostRuleFlattenID(cluster *object.ClusterComputeResource, key int32) (string, error) {
	clusterID := cluster.Reference().Value
	return strings.Join([]string{clusterID, strconv.Itoa(int(key))}, ":"), nil
}

// resourceVSphereComputeClusterVMHostRuleParseID parses an ID for the
// vsphere_compute_cluster_vm_host_rule and outputs its parts.
func resourceVSphereComputeClusterVMHostRuleParseID(id string) (string, int32, error) {
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

// resourceVSphereComputeClusterVMHostRuleFindEntry attempts to locate an
// existing VM/host rule in a cluster's configuration by key. It's used by the
// resource's read functionality and tests. nil is returned if the entry cannot
// be found.
func resourceVSphereComputeClusterVMHostRuleFindEntry(
	cluster *object.ClusterComputeResource,
	key int32,
) (*types.ClusterVmHostRuleInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Key == key {
			if vmHostRuleInfo, ok := info.(*types.ClusterVmHostRuleInfo); ok {
				log.Printf("[DEBUG] Found VM/host rule key %d in cluster %q", key, cluster.Name())
				return vmHostRuleInfo, nil
			}
			return nil, fmt.Errorf("rule key %d in cluster %q is not a VM/host rule", key, cluster.Name())
		}
	}

	log.Printf("[DEBUG] No VM/host rule key %d found in cluster %q", key, cluster.Name())
	return nil, nil
}

// resourceVSphereComputeClusterVMHostRuleFindEntryByName attempts to locate an
// existing VM/host rule in a cluster's configuration by name. It differs from
// the standard resourceVSphereComputeClusterVMHostRuleFindEntry in that we
// don't allow missing entries, as it's designed to be used in places where we
// don't want to allow for missing entries, such as during creation and import.
func resourceVSphereComputeClusterVMHostRuleFindEntryByName(
	cluster *object.ClusterComputeResource,
	name string,
) (*types.ClusterVmHostRuleInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).Rule {
		if info.GetClusterRuleInfo().Name == name {
			if vmHostRuleInfo, ok := info.(*types.ClusterVmHostRuleInfo); ok {
				log.Printf("[DEBUG] Found VM/host rule %q in cluster %q", name, cluster.Name())
				return vmHostRuleInfo, nil
			}
			return nil, fmt.Errorf("rule %q in cluster %q is not a VM/host rule", name, cluster.Name())
		}
	}

	return nil, fmt.Errorf("no VM/host rule %q found in cluster %q", name, cluster.Name())
}

// resourceVSphereComputeClusterVMHostRuleObjects handles the fetching of the
// cluster and rule key depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, only the cluster is retrieved from compute_cluster_id. -1 is
// returned for the key.
func resourceVSphereComputeClusterVMHostRuleObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	if d.Id() != "" {
		return resourceVSphereComputeClusterVMHostRuleObjectsFromID(d, meta)
	}
	return resourceVSphereComputeClusterVMHostRuleObjectsFromAttributes(d, meta)
}

func resourceVSphereComputeClusterVMHostRuleObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	return resourceVSphereComputeClusterVMHostRuleFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		-1,
	)
}

func resourceVSphereComputeClusterVMHostRuleObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, int32, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, key, err := resourceVSphereComputeClusterVMHostRuleParseID(d.Id())
	if err != nil {
		return nil, 0, err
	}

	return resourceVSphereComputeClusterVMHostRuleFetchObjects(meta, clusterID, key)
}

// resourceVSphereComputeClusterVMHostRuleFetchObjects fetches the "objects"
// for a cluster rule. This is currently just the cluster object as the rule
// key a static value and a pass-through - this is to keep its workflow
// consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereComputeClusterVMHostRuleFetchObjects(
	meta interface{},
	clusterID string,
	key int32,
) (*object.ClusterComputeResource, int32, error) {
	client, err := resourceVSphereComputeClusterVMHostRuleClient(meta)
	if err != nil {
		return nil, 0, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot locate cluster: %s", err)
	}

	return cluster, key, nil
}

func resourceVSphereComputeClusterVMHostRuleClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
