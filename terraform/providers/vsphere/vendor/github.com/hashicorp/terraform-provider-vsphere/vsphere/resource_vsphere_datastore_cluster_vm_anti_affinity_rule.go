package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/storagepod"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereDatastoreClusterVMAntiAffinityRuleName = "vsphere_datastore_cluster_vm_anti_affinity_rule"

func resourceVSphereDatastoreClusterVMAntiAffinityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDatastoreClusterVMAntiAffinityRuleCreate,
		Read:   resourceVSphereDatastoreClusterVMAntiAffinityRuleRead,
		Update: resourceVSphereDatastoreClusterVMAntiAffinityRuleUpdate,
		Delete: resourceVSphereDatastoreClusterVMAntiAffinityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereDatastoreClusterVMAntiAffinityRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"datastore_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the datastore cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the virtual machine group in the cluster.",
			},
			"virtual_machine_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The UUIDs of the virtual machines to run on different datastores from each other.",
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

func resourceVSphereDatastoreClusterVMAntiAffinityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceVSphereDatastoreClusterVMAntiAffinityRuleValidateRuleVMCount(d); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))

	pod, _, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAntiAffinityRuleSpec(d, meta)
	if err != nil {
		return err
	}
	spec := types.StorageDrsConfigSpec{
		PodConfigSpec: &types.StorageDrsPodConfigSpec{
			Rule: []types.ClusterRuleSpec{
				{
					ArrayUpdateSpec: types.ArrayUpdateSpec{
						Operation: types.ArrayUpdateOperationAdd,
					},
					Info: info,
				},
			},
		},
	}

	if err = resourceVSphereDatastoreClusterVMAntiAffinityRuleApplySDRSConfigSpec(pod, spec); err != nil {
		return err
	}

	info, err = resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntryByName(pod, info.Name)
	if err != nil {
		return err
	}

	id, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleFlattenID(pod, info.Key)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))
	return resourceVSphereDatastoreClusterVMAntiAffinityRuleRead(d, meta)
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))

	pod, key, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntry(pod, key)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the datastore_cluster_id. This is ForceNew, but we set these for
	// completeness on import so that if the wrong pod/VM combo was used, it
	// will be noted.
	if err = d.Set("datastore_cluster_id", pod.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"datastore_cluster_id\": %s", err)
	}

	if err = flattenClusterAntiAffinityRuleSpec(d, meta, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))
	return nil
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceVSphereDatastoreClusterVMAntiAffinityRuleValidateRuleVMCount(d); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))

	pod, key, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterAntiAffinityRuleSpec(d, meta)
	if err != nil {
		return err
	}
	info.Key = key

	spec := types.StorageDrsConfigSpec{
		PodConfigSpec: &types.StorageDrsPodConfigSpec{
			Rule: []types.ClusterRuleSpec{
				{
					ArrayUpdateSpec: types.ArrayUpdateSpec{
						Operation: types.ArrayUpdateOperationEdit,
					},
					Info: info,
				},
			},
		},
	}

	if err := resourceVSphereDatastoreClusterVMAntiAffinityRuleApplySDRSConfigSpec(pod, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))
	return resourceVSphereDatastoreClusterVMAntiAffinityRuleRead(d, meta)
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))

	pod, key, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects(d, meta)
	if err != nil {
		return err
	}

	spec := types.StorageDrsConfigSpec{
		PodConfigSpec: &types.StorageDrsPodConfigSpec{
			Rule: []types.ClusterRuleSpec{
				{
					ArrayUpdateSpec: types.ArrayUpdateSpec{
						Operation: types.ArrayUpdateOperationRemove,
						RemoveKey: key,
					},
				},
			},
		},
	}

	if err := resourceVSphereDatastoreClusterVMAntiAffinityRuleApplySDRSConfigSpec(pod, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d))
	return nil
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	podPath, ok := data["datastore_cluster_path"]
	if !ok {
		return nil, errors.New("missing datastore_cluster_path in input data")
	}
	name, ok := data["name"]
	if !ok {
		return nil, errors.New("missing name in input data")
	}

	client, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleClient(meta)
	if err != nil {
		return nil, err
	}

	pod, err := storagepod.FromPath(client, podPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate datastore cluster %q: %s", podPath, err)
	}

	info, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntryByName(pod, name)
	if err != nil {
		return nil, err
	}

	id, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleFlattenID(pod, info.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString prints a friendly string for the
// vsphere_datastore_cluster_vm_anti_affinity_rule resource.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereDatastoreClusterVMAntiAffinityRuleName)
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleFlattenID makes an ID for the
// vsphere_datastore_cluster_vm_anti_affinity_rule resource.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleFlattenID(pod *object.StoragePod, key int32) (string, error) {
	podID := pod.Reference().Value
	return strings.Join([]string{podID, strconv.Itoa(int(key))}, ":"), nil
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleParseID parses an ID for the
// vsphere_datastore_cluster_vm_anti_affinity_rule and outputs its parts.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleParseID(id string) (string, int32, error) {
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

// resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntry attempts to
// locate an existing VM anti-affinity rule in a datastore cluster's
// configuration by key. It's used by the resource's read functionality and
// tests. nil is returned if the entry cannot be found.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntry(
	pod *object.StoragePod,
	key int32,
) (*types.ClusterAntiAffinityRuleSpec, error) {
	props, err := storagepod.Properties(pod)
	if err != nil {
		return nil, fmt.Errorf("error fetching datastore cluster properties: %s", err)
	}

	for _, info := range props.PodStorageDrsEntry.StorageDrsConfig.PodConfig.Rule {
		if info.GetClusterRuleInfo().Key == key {
			if vmAntiAffinityRuleInfo, ok := info.(*types.ClusterAntiAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM anti-affinity rule key %d in datastore cluster %q", key, pod.Name())
				return vmAntiAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule key %d in datastore cluster %q is not a VM anti-affinity rule", key, pod.Name())
		}
	}

	log.Printf("[DEBUG] No VM anti-affinity rule key %d found in datastore cluster %q", key, pod.Name())
	return nil, nil
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntryByName attempts to
// locate an existing VM anti-affinity rule in a datastore cluster's
// configuration by name. It differs from the standard
// resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntry in that we don't
// allow missing entries, as it's designed to be used in places where we don't
// want to allow for missing entries, such as during creation and import.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleFindEntryByName(
	pod *object.StoragePod,
	name string,
) (*types.ClusterAntiAffinityRuleSpec, error) {
	props, err := storagepod.Properties(pod)
	if err != nil {
		return nil, fmt.Errorf("error fetching datastore cluster properties: %s", err)
	}

	for _, info := range props.PodStorageDrsEntry.StorageDrsConfig.PodConfig.Rule {
		if info.GetClusterRuleInfo().Name == name {
			if vmAntiAffinityRuleInfo, ok := info.(*types.ClusterAntiAffinityRuleSpec); ok {
				log.Printf("[DEBUG] Found VM anti-affinity rule %q in datastore cluster %q", name, pod.Name())
				return vmAntiAffinityRuleInfo, nil
			}
			return nil, fmt.Errorf("rule %q in datastore cluster %q is not a VM anti-affinity rule", name, pod.Name())
		}
	}

	return nil, fmt.Errorf("no VM anti-affinity rule %q found in datastore cluster %q", name, pod.Name())
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects handles the
// fetching of the cluster and rule key depending on what attributes are
// available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, only the cluster is retrieved from datastore_cluster_id. -1 is
// returned for the key.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.StoragePod, int32, error) {
	if d.Id() != "" {
		return resourceVSphereDatastoreClusterVMAntiAffinityRuleObjectsFromID(d, meta)
	}
	return resourceVSphereDatastoreClusterVMAntiAffinityRuleObjectsFromAttributes(d, meta)
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.StoragePod, int32, error) {
	return resourceVSphereDatastoreClusterVMAntiAffinityRuleFetchObjects(
		meta,
		d.Get("datastore_cluster_id").(string),
		-1,
	)
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.StoragePod, int32, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	podID, key, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleParseID(d.Id())
	if err != nil {
		return nil, 0, err
	}

	return resourceVSphereDatastoreClusterVMAntiAffinityRuleFetchObjects(meta, podID, key)
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleFetchObjects fetches the
// "objects" for a cluster rule. This is currently just the cluster object as
// the rule key a static value and a pass-through - this is to keep its
// workflow consistent with other cluster-dependent resources that derive from
// ArrayUpdateSpec that have managed object as keys, such as VM and host
// overrides.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleFetchObjects(
	meta interface{},
	podID string,
	key int32,
) (*object.StoragePod, int32, error) {
	client, err := resourceVSphereDatastoreClusterVMAntiAffinityRuleClient(meta)
	if err != nil {
		return nil, 0, err
	}

	pod, err := storagepod.FromID(client, podID)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot locate datastore cluster: %s", err)
	}

	return pod, key, nil
}

func resourceVSphereDatastoreClusterVMAntiAffinityRuleClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleApplySDRSConfigSpec
// applying a SDRS config spec for the
// vsphere_datastore_cluster_vm_anti_affinity_rule resource.
//
// This is wrapped to abstract the fact that we are deriving the client from
// the StoragePod. This is because helper workflows that have been created more
// recently (ie: cluster helpers) do this, and more than likely the storagepod
// helper will do it eventually as well. If there is ever an issue with this,
// it can be changed here.  There should be no issue though as govmomi.Client
// is mainly just vim25.Client with some additional session helper bits that is
// not normally needed during normal operation.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleApplySDRSConfigSpec(
	pod *object.StoragePod,
	spec types.StorageDrsConfigSpec,
) error {
	return storagepod.ApplyDRSConfiguration(
		&govmomi.Client{
			Client: pod.Client(),
		},
		pod,
		spec,
	)
}

// resourceVSphereDatastoreClusterVMAntiAffinityRuleValidateRuleVMCount ensures
// that the VM count in the anti-affinity rule at any point in time before it's
// created or updated is a length of at least 2.
//
// This validation is necessary as a rule of only 1 VM here is a no-op and
// ultimately will result in a broken resource (the rule will not exist after
// creation, or example). Unfortunately, this needs to happen at apply time
// right now due to issues with TF core and how it processes lists when values
// are computed. Once these issues are fixed with TF core, the validation here
// should be removed and moved to schema.
func resourceVSphereDatastoreClusterVMAntiAffinityRuleValidateRuleVMCount(d *schema.ResourceData) error {
	if d.Get("virtual_machine_ids").(*schema.Set).Len() < 2 {
		return errors.New("length of virtual_machine_ids must be 2 or more")
	}
	return nil
}
