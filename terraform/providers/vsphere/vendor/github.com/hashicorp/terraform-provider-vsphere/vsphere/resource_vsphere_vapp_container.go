package vsphere

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/resourcepool"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/vappcontainer"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereVAppContainerName = "vsphere_vapp_container"

var vAppContainerCPUSharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

var vAppContainerMemorySharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

func resourceVSphereVAppContainer() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the vApp container.",
		},
		"parent_resource_pool_id": {
			Type:        schema.TypeString,
			Description: "The managed object ID of the parent resource pool or the compute resource the vApp container is in.",
			Required:    true,
		},
		"parent_folder_id": {
			Type:        schema.TypeString,
			Description: "The ID of the parent VM folder.",
			Optional:    true,
		},
		"cpu_share_level": {
			Type:         schema.TypeString,
			Description:  "The allocation level. The level is a simplified view of shares. Levels map to a pre-determined set of numeric values for shares. Can be one of low, normal, high, or custom.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(vAppContainerCPUSharesLevelAllowedValues, false),
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
			Description: "Amount of CPU (MHz) that is guaranteed available to the vApp container.",
			Optional:    true,
			Default:     0,
		},
		"cpu_expandable": {
			Type:        schema.TypeBool,
			Description: "Determines if the reservation on a vApp container can grow beyond the specified value, if the parent resource pool has unreserved resources.",
			Optional:    true,
			Default:     true,
		},
		"cpu_limit": {
			Type:        schema.TypeInt,
			Description: "The utilization of a vApp container will not exceed this limit, even if there are available resources. Set to -1 for unlimited.",
			Optional:    true,
			Default:     -1,
		},
		"memory_share_level": {
			Type:         schema.TypeString,
			Description:  "The allocation level. The level is a simplified view of shares. Levels map to a pre-determined set of numeric values for shares. Can be one of low, normal, high, or custom.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(vAppContainerMemorySharesLevelAllowedValues, false),
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
			Description: "Amount of memory (MB) that is guaranteed available to the vApp container.",
			Optional:    true,
			Default:     0,
		},
		"memory_expandable": {
			Type:        schema.TypeBool,
			Description: "Determines if the reservation on a vApp container can grow beyond the specified value, if the parent resource pool has unreserved resources.",
			Optional:    true,
			Default:     true,
		},
		"memory_limit": {
			Type:        schema.TypeInt,
			Description: "The utilization of a vApp container will not exceed this limit, even if there are available resources. Set to -1 for unlimited.",
			Optional:    true,
			Default:     -1,
		},
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	return &schema.Resource{
		Create: resourceVSphereVAppContainerCreate,
		Read:   resourceVSphereVAppContainerRead,
		Update: resourceVSphereVAppContainerUpdate,
		Delete: resourceVSphereVAppContainerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereVAppContainerImport,
		},
		Schema: s,
	}
}

func resourceVSphereVAppContainerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return nil, err
	}
	vc, err := vappcontainer.FromPath(client, d.Id(), nil)
	if err != nil {
		return nil, err
	}
	d.SetId(vc.Reference().Value)
	return []*schema.ResourceData{d}, nil
}

func resourceVSphereVAppContainerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereVAppContainerIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	prp, err := resourcepool.FromID(client, d.Get("parent_resource_pool_id").(string))
	if err != nil {
		return err
	}
	rpSpec := expandVAppContainerConfigSpec(d)
	vcSpec := &types.VAppConfigSpec{}
	var f *object.Folder

	if vappcontainer.IsVApp(client, d.Get("parent_resource_pool_id").(string)) {
		f = nil
	} else {
		if pf, ok := d.GetOk("parent_folder_id"); ok {
			f, err = folder.FromID(client, pf.(string))
			if err != nil {
				return err
			}
		} else {
			p := strings.Split(prp.InventoryPath, "/")
			if len(p) < 2 {
				return fmt.Errorf("unable to locate datacenter name from parent resource pool")
			}
			var dc *object.Datacenter
			dc, err = getDatacenter(client, p[1])
			if err != nil {
				return err
			}
			f, err = folder.FromPath(client, "", folder.VSphereFolderTypeVM, dc)
			if err != nil {
				return err
			}
		}
	}

	vc, err := vappcontainer.Create(prp, d.Get("name").(string), rpSpec, vcSpec, f)
	if err != nil {
		return err
	}
	if err = resourceVSphereVAppContainerApplyTags(d, meta, vc); err != nil {
		return err
	}
	d.SetId(vc.Reference().Value)
	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereVAppContainerIDString(d))
	return resourceVSphereVAppContainerRead(d, meta)
}

func resourceVSphereVAppContainerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereVAppContainerIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	vc, err := vappcontainer.FromID(client, d.Id())
	if err != nil {
		if viapi.IsManagedObjectNotFoundError(err) {
			log.Printf("[DEBUG] %s: Resource has been deleted", resourceVSphereVAppContainerIDString(d))
			d.SetId("")
			return nil
		}
		return err
	}
	if err = resourceVSphereVAppContainerReadTags(d, meta, vc); err != nil {
		return err
	}
	if err = d.Set("name", vc.Name()); err != nil {
		return err
	}
	vcProps, err := vappcontainer.Properties(vc)
	if err != nil {
		return err
	}
	if vcProps.Parent != nil {
		if err = d.Set("parent_resource_pool_id", vcProps.Parent.Value); err != nil {
			return err
		}
	}

	if vcProps.ParentFolder != nil {
		if err = d.Set("parent_folder_id", vcProps.ParentFolder.Value); err != nil {
			return err
		}
	}
	if err = flattenVAppContainerConfigSpec(d, vcProps.Config); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Read finished successfully", resourceVSphereVAppContainerIDString(d))
	return nil
}

func resourceVSphereVAppContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereVAppContainerIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	vc, err := vappcontainer.FromID(client, d.Id())
	if err != nil {
		return err
	}
	if err = resourceVSphereVAppContainerApplyTags(d, meta, vc); err != nil {
		return err
	}
	op, np := d.GetChange("parent_resource_pool_id")
	if op != np {
		log.Printf("[DEBUG] %s: Parent resource pool has changed. Moving from %s, to %s", resourceVSphereVAppContainerIDString(d), op, np)
		var p *object.ResourcePool
		p, err = resourcepool.FromID(client, np.(string))
		if err != nil {
			return err
		}
		if err = resourcepool.MoveIntoResourcePool(p, vc.Reference()); err != nil {
			return err
		}
		log.Printf("[DEBUG] %s: Move finished successfully", resourceVSphereVAppContainerIDString(d))
	}

	vcSpec := types.VAppConfigSpec{}
	if err = vappcontainer.Update(vc, vcSpec); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereVAppContainerIDString(d))
	return nil
}

func resourceVSphereVAppContainerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereVAppContainerIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	vc, err := vappcontainer.FromID(client, d.Id())
	if err != nil {
		return err
	}
	if err = resourceVSphereVAppContainerValidateEmpty(vc); err != nil {
		return err
	}
	if err = vappcontainer.Delete(vc); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereVAppContainerIDString(d))
	return nil
}

// resourceVSphereVAppContainerIDString prints a friendly string for the
// vsphere_vapp_container resource.
func resourceVSphereVAppContainerIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereVAppContainerName)
}

func flattenVAppContainerConfigSpec(d *schema.ResourceData, obj types.ResourceConfigSpec) error {
	if err := flattenVAppContainerMemoryAllocation(d, obj.MemoryAllocation); err != nil {
		return err
	}
	return flattenVAppContainerCPUAllocation(d, obj.CpuAllocation)
}

func flattenVAppContainerCPUAllocation(d *schema.ResourceData, obj types.ResourceAllocationInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"cpu_reservation": obj.Reservation,
		"cpu_expandable":  obj.ExpandableReservation,
		"cpu_limit":       obj.Limit,
		"cpu_shares":      obj.Shares.Shares,
		"cpu_share_level": obj.Shares.Level,
	})
}

func flattenVAppContainerMemoryAllocation(d *schema.ResourceData, obj types.ResourceAllocationInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"memory_reservation": obj.Reservation,
		"memory_expandable":  obj.ExpandableReservation,
		"memory_limit":       obj.Limit,
		"memory_shares":      obj.Shares.Shares,
		"memory_share_level": obj.Shares.Level,
	})
}

func expandVAppContainerConfigSpec(d *schema.ResourceData) *types.ResourceConfigSpec {
	return &types.ResourceConfigSpec{
		CpuAllocation:    expandResourcePoolCPUAllocation(d),
		MemoryAllocation: expandResourcePoolMemoryAllocation(d),
	}
}

func resourceVSphereVAppContainerClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}

func resourceVSphereVAppContainerValidateEmpty(va *object.VirtualApp) error {
	ne, err := vappcontainer.HasChildren(va)
	if err != nil {
		return fmt.Errorf("error checking contents of vApp container: %s", err)
	}
	if ne {
		return fmt.Errorf("vApp container %q still has children resources. Please move or remove all items before deleting", va.InventoryPath)
	}
	return nil
}

// resourceVSphereVAppContainerApplyTags processes the tags step for both create
// and update for vsphere_vapp_container.
func resourceVSphereVAppContainerApplyTags(d *schema.ResourceData, meta interface{}, va *object.VirtualApp) error {
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}

	// Apply any pending tags now.
	if tagsClient == nil {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping", resourceVSphereComputeClusterIDString(d))
		return nil
	}

	log.Printf("[DEBUG] %s: Applying any pending tags", resourceVSphereVAppContainerIDString(d))
	return processTagDiff(tagsClient, d, va)
}

// resourceVSphereVAppContainerReadTags reads the tags for
// vsphere_vapp_container.
func resourceVSphereVAppContainerReadTags(d *schema.ResourceData, meta interface{}, va *object.VirtualApp) error {
	if tagsClient, _ := meta.(*Client).TagsManager(); tagsClient != nil {
		log.Printf("[DEBUG] %s: Reading tags", resourceVSphereVAppContainerIDString(d))
		if err := readTagsForResource(tagsClient, va, d); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] %s: Tags unsupported on this connection, skipping tag read", resourceVSphereVAppContainerIDString(d))
	}
	return nil
}
