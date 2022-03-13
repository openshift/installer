package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/resourcepool"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereResourcePoolName = "vsphere_resource_pool"

var resourcePoolCPUSharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

var resourcePoolMemorySharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

func resourceVSphereResourcePool() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of resource pool.",
		},
		"parent_resource_pool_id": {
			Type:        schema.TypeString,
			Description: "The ID of the root resource pool of the compute resource the resource pool is in.",
			Required:    true,
		},
		"cpu_share_level": {
			Type:         schema.TypeString,
			Description:  "The allocation level. The level is a simplified view of shares. Levels map to a pre-determined set of numeric values for shares. Can be one of low, normal, high, or custom.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(resourcePoolCPUSharesLevelAllowedValues, false),
			Default:      "normal",
		},
		"cpu_shares": {
			Type:        schema.TypeInt,
			Description: "The number of shares allocated. Used to determine resource allocation in case of resource contention. If this is set, cpu_share_level must be custom.",
			Computed:    true,
			Optional:    true,
		},
		"cpu_reservation": {
			Type:        schema.TypeInt,
			Description: "Amount of CPU (MHz) that is guaranteed available to the resource pool.",
			Optional:    true,
			Default:     0,
		},
		"cpu_expandable": {
			Type:        schema.TypeBool,
			Description: "Determines if the reservation on a resource pool can grow beyond the specified value, if the parent resource pool has unreserved resources.",
			Optional:    true,
			Default:     true,
		},
		"cpu_limit": {
			Type:        schema.TypeInt,
			Description: "The utilization of a resource pool will not exceed this limit, even if there are available resources. Set to -1 for unlimited.",
			Optional:    true,
			Default:     -1,
		},
		"memory_share_level": {
			Type:         schema.TypeString,
			Description:  "The allocation level. The level is a simplified view of shares. Levels map to a pre-determined set of numeric values for shares. Can be one of low, normal, high, or custom.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(resourcePoolMemorySharesLevelAllowedValues, false),
			Default:      "normal",
		},
		"memory_shares": {
			Type:        schema.TypeInt,
			Description: "The number of shares allocated. Used to determine resource allocation in case of resource contention. If this is set, memory_share_level must be custom.",
			Computed:    true,
			Optional:    true,
		},
		"memory_reservation": {
			Type:        schema.TypeInt,
			Description: "Amount of memory (MB) that is guaranteed available to the resource pool.",
			Optional:    true,
			Default:     0,
		},
		"memory_expandable": {
			Type:        schema.TypeBool,
			Description: "Determines if the reservation on a resource pool can grow beyond the specified value, if the parent resource pool has unreserved resources.",
			Optional:    true,
			Default:     true,
		},
		"memory_limit": {
			Type:        schema.TypeInt,
			Description: "The utilization of a resource pool will not exceed this limit, even if there are available resources. Set to -1 for unlimited.",
			Optional:    true,
			Default:     -1,
		},
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	return &schema.Resource{
		Create: resourceVSphereResourcePoolCreate,
		Read:   resourceVSphereResourcePoolRead,
		Update: resourceVSphereResourcePoolUpdate,
		Delete: resourceVSphereResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereResourcePoolImport,
		},
		Schema: s,
	}
}

func resourceVSphereResourcePoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := resourceVSphereResourcePoolClient(meta)
	if err != nil {
		return nil, err
	}
	rp, err := resourcepool.FromPathOrDefault(client, d.Id(), nil)
	if err != nil {
		return nil, err
	}
	d.SetId(rp.Reference().Value)
	return []*schema.ResourceData{d}, nil
}

func resourceVSphereResourcePoolCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereResourcePoolIDString(d))
	client, err := resourceVSphereResourcePoolClient(meta)
	if err != nil {
		return err
	}
	prp, err := resourcepool.FromID(client, d.Get("parent_resource_pool_id").(string))
	if err != nil {
		return err
	}
	rpSpec := expandResourcePoolConfigSpec(d)
	rp, err := resourcepool.Create(prp, d.Get("name").(string), rpSpec)
	if err != nil {
		return err
	}
	if err = resourceVSphereResourcePoolApplyTags(d, meta, rp); err != nil {
		return err
	}
	d.SetId(rp.Reference().Value)
	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereResourcePoolIDString(d))
	return resourceVSphereResourcePoolRead(d, meta)
}

func resourceVSphereResourcePoolRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereResourcePoolIDString(d))
	client, err := resourceVSphereResourcePoolClient(meta)
	if err != nil {
		return err
	}
	rp, err := resourcepool.FromID(client, d.Id())
	if err != nil {
		if viapi.IsManagedObjectNotFoundError(err) {
			log.Printf("[DEBUG] %s: Resource has been deleted", resourceVSphereResourcePoolIDString(d))
			d.SetId("")
			return nil
		}
		return err
	}
	if err = resourceVSphereResourcePoolReadTags(d, meta, rp); err != nil {
		return err
	}
	err = d.Set("name", rp.Name())
	if err != nil {
		return err
	}
	rpProps, err := resourcepool.Properties(rp)
	if err != nil {
		return err
	}
	if err = d.Set("parent_resource_pool_id", rpProps.Parent.Value); err != nil {
		return err
	}
	err = flattenResourcePoolConfigSpec(d, rpProps.Config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Read finished successfully", resourceVSphereResourcePoolIDString(d))
	return nil
}

func resourceVSphereResourcePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereResourcePoolIDString(d))
	client, err := resourceVSphereResourcePoolClient(meta)
	if err != nil {
		return err
	}
	rp, err := resourcepool.FromID(client, d.Id())
	if err != nil {
		return err
	}
	if err = resourceVSphereResourcePoolApplyTags(d, meta, rp); err != nil {
		return err
	}
	op, np := d.GetChange("parent_resource_pool_id")
	if op != np {
		log.Printf("[DEBUG] %s: Parent resource pool has changed. Moving from %s, to %s", resourceVSphereResourcePoolIDString(d), op, np)
		p, err := resourcepool.FromID(client, np.(string))
		if err != nil {
			return err
		}
		err = resourcepool.MoveIntoResourcePool(p, rp.Reference())
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] %s: Move finished successfully", resourceVSphereResourcePoolIDString(d))
	}

	rpSpec := expandResourcePoolConfigSpec(d)
	err = resourcepool.Update(rp, d.Get("name").(string), rpSpec)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereResourcePoolIDString(d))
	return nil
}

func resourceVSphereResourcePoolDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereResourcePoolIDString(d))
	client, err := resourceVSphereResourcePoolClient(meta)
	if err != nil {
		return err
	}
	rp, err := resourcepool.FromID(client, d.Id())
	if err != nil {
		return err
	}
	err = resourceVSphereResourcePoolValidateEmpty(rp)
	if err != nil {
		return err
	}
	err = resourcepool.Delete(rp)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereResourcePoolIDString(d))
	return nil
}

// resourceVSphereResourcePoolIDString prints a friendly string for the
// vsphere_virtual_machine resource.
func resourceVSphereResourcePoolIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereResourcePoolName)
}

func flattenResourcePoolConfigSpec(d *schema.ResourceData, obj types.ResourceConfigSpec) error {
	err := flattenResourcePoolMemoryAllocation(d, obj.MemoryAllocation)
	if err != nil {
		return err
	}
	return flattenResourcePoolCPUAllocation(d, obj.CpuAllocation)
}

func flattenResourcePoolCPUAllocation(d *schema.ResourceData, obj types.ResourceAllocationInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"cpu_reservation": obj.Reservation,
		"cpu_expandable":  obj.ExpandableReservation,
		"cpu_limit":       obj.Limit,
		"cpu_shares":      obj.Shares.Shares,
		"cpu_share_level": obj.Shares.Level,
	})
}

func flattenResourcePoolMemoryAllocation(d *schema.ResourceData, obj types.ResourceAllocationInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"memory_reservation": obj.Reservation,
		"memory_expandable":  obj.ExpandableReservation,
		"memory_limit":       obj.Limit,
		"memory_shares":      obj.Shares.Shares,
		"memory_share_level": obj.Shares.Level,
	})
}

func expandResourcePoolConfigSpec(d *schema.ResourceData) *types.ResourceConfigSpec {
	return &types.ResourceConfigSpec{
		CpuAllocation:    expandResourcePoolCPUAllocation(d),
		MemoryAllocation: expandResourcePoolMemoryAllocation(d),
	}
}

func expandResourcePoolCPUAllocation(d *schema.ResourceData) types.ResourceAllocationInfo {
	return types.ResourceAllocationInfo{
		Reservation:           structure.GetInt64Ptr(d, "cpu_reservation"),
		ExpandableReservation: structure.GetBoolPtr(d, "cpu_expandable"),
		Limit:                 structure.GetInt64Ptr(d, "cpu_limit"),
		Shares: &types.SharesInfo{
			Level:  types.SharesLevel(d.Get("cpu_share_level").(string)),
			Shares: int32(d.Get("cpu_shares").(int)),
		},
	}
}

func expandResourcePoolMemoryAllocation(d *schema.ResourceData) types.ResourceAllocationInfo {
	return types.ResourceAllocationInfo{
		Reservation:           structure.GetInt64Ptr(d, "memory_reservation"),
		ExpandableReservation: structure.GetBoolPtr(d, "memory_expandable"),
		Limit:                 structure.GetInt64Ptr(d, "memory_limit"),
		Shares: &types.SharesInfo{
			Shares: int32(d.Get("memory_shares").(int)),
			Level:  types.SharesLevel(d.Get("memory_share_level").(string)),
		},
	}
}

func resourceVSphereResourcePoolClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}

func resourceVSphereResourcePoolValidateEmpty(rp *object.ResourcePool) error {
	ne, err := resourcepool.HasChildren(rp)
	if err != nil {
		return fmt.Errorf("error checking contents of resource pool: %s", err)
	}
	if ne {
		return fmt.Errorf("resource pool %q still has children resources. Please move or remove all items before deleting", rp.InventoryPath)
	}
	return nil
}

// resourceVSphereResourcePoolApplyTags processes the tags step for both create
// and update for vsphere_resource_pool.
func resourceVSphereResourcePoolApplyTags(d *schema.ResourceData, meta interface{}, rp *object.ResourcePool) error {
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}

	// Apply any pending tags now.
	if tagsClient == nil {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping", resourceVSphereComputeClusterIDString(d))
		return nil
	}

	log.Printf("[DEBUG] %s: Applying any pending tags", resourceVSphereResourcePoolIDString(d))
	return processTagDiff(tagsClient, d, rp)
}

// resourceVSphereResourcePoolReadTags reads the tags for
// vsphere_resource_pool.
func resourceVSphereResourcePoolReadTags(d *schema.ResourceData, meta interface{}, rp *object.ResourcePool) error {
	if tagsClient, _ := meta.(*VSphereClient).TagsManager(); tagsClient != nil {
		log.Printf("[DEBUG] %s: Reading tags", resourceVSphereResourcePoolIDString(d))
		if err := readTagsForResource(tagsClient, rp, d); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping tag read", resourceVSphereResourcePoolIDString(d))
	}
	return nil
}
