package vsphere

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/vappcontainer"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereVAppEntityName = "vsphere_vapp_entity"

func resourceVSphereVAppEntity() *schema.Resource {
	s := map[string]*schema.Schema{
		"target_id": {
			Type:        schema.TypeString,
			Description: "Managed object ID of the entity to power on or power off. This can be a virtual machine or a vApp.",
			Required:    true,
			ForceNew:    true,
		},
		"container_id": {
			Type:        schema.TypeString,
			Description: "Managed object ID of the vApp container the entity is a member of.",
			Required:    true,
			ForceNew:    true,
		},
		"start_action": {
			Type:        schema.TypeString,
			Description: "How to start the entity. Valid settings are none or powerOn. If set to none, then the entity does not participate in auto-start.",
			Default:     "powerOn",
			Optional:    true,
		},
		"start_delay": {
			Type:        schema.TypeInt,
			Description: "Delay in seconds before continuing with the next entity in the order of entities to be started.",
			Default:     120,
			Optional:    true,
		},
		"stop_action": {
			Type:        schema.TypeString,
			Description: "Defines the stop action for the entity. Can be set to none, powerOff, guestShutdown, or suspend. If set to none, then the entity does not participate in auto-stop.",
			Default:     "powerOff",
			Optional:    true,
		},
		"stop_delay": {
			Type:        schema.TypeInt,
			Description: "Delay in seconds before continuing with the next entity in the order of entities to be stopped.",
			Default:     120,
			Optional:    true,
		},
		"start_order": {
			Type:        schema.TypeInt,
			Description: "Order to start and stop target in vApp.",
			Optional:    true,
			Default:     1,
		},
		"wait_for_guest": {
			Type:        schema.TypeBool,
			Description: "Determines if the VM should be marked as being started when VMware Tools are ready instead of waiting for start_delay. This property has no effect for vApps.",
			Optional:    true,
			Default:     false,
		},
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	return &schema.Resource{
		Create: resourceVSphereVAppEntityCreate,
		Read:   resourceVSphereVAppEntityRead,
		Update: resourceVSphereVAppEntityUpdate,
		Delete: resourceVSphereVAppEntityDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereVAppEntityImport,
		},
		Schema: s,
	}
}

func resourceVSphereVAppEntityImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := resourceVSphereVAppEntityClient(meta)
	if err != nil {
		return nil, err
	}
	vcID, vmID, err := resourceVSphereVAppEntitySplitID(d.Id())
	if err != nil {
		return nil, err
	}
	vc, err := vappcontainer.FromID(client, vcID)
	if err != nil {
		return nil, err
	}
	vcProps, err := vappcontainer.Properties(vc)
	if err != nil {
		return nil, err
	}
	ve := resourceVSphereVAppEntityFromKey(vmID, vcProps)
	if ve == nil {
		return nil, fmt.Errorf("unable to locate vapp_entity %s", d.Id())
	}
	return []*schema.ResourceData{d}, nil
}

func resourceVSphereVAppEntityCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereVAppEntityIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	container, err := vappcontainer.FromID(client, d.Get("container_id").(string))
	if err != nil {
		return err
	}
	cProps, err := vappcontainer.Properties(container)
	if err != nil {
		return err
	}
	d.SetId(resourceVSphereVAppEntityIDFromKeys(d.Get("container_id").(string), d.Get("target_id").(string)))
	entityConfig, err := expandVAppEntityConfigSpec(client, d)
	if err != nil {
		return err
	}
	cProps.VAppConfig.EntityConfig = append(cProps.VAppConfig.EntityConfig, *entityConfig)
	updateSpec := types.VAppConfigSpec{
		EntityConfig: cProps.VAppConfig.EntityConfig,
	}
	if err = vappcontainer.Update(container, updateSpec); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereVAppEntityIDString(d))
	return resourceVSphereVAppEntityRead(d, meta)
}

func resourceVSphereVAppEntityRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereVAppEntityIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	vcID, vmID, err := resourceVSphereVAppEntitySplitID(d.Id())
	if err != nil {
		return err
	}
	if err = d.Set("container_id", vcID); err != nil {
		return err
	}
	if err = d.Set("target_id", vmID); err != nil {
		return err
	}
	entity, err := resourceVSphereVAppEntityFind(client, d.Id())
	if err != nil {
		return err
	}
	if entity == nil {
		log.Printf("[DEBUG] %s: Resource has been deleted", resourceVSphereVAppEntityIDString(d))
		d.SetId("")
		return nil
	}
	err = flattenVAppEntityConfigSpec(client, d, entity)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Read finished successfully", resourceVSphereVAppEntityIDString(d))
	return nil
}

func resourceVSphereVAppEntityUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereVAppEntityIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	container, err := vappcontainer.FromID(client, d.Get("container_id").(string))
	if err != nil {
		return err
	}
	cProps, err := vappcontainer.Properties(container)
	if err != nil {
		return err
	}
	entityConfig, err := expandVAppEntityConfigSpec(client, d)
	if err != nil {
		return err
	}
	cProps.VAppConfig.EntityConfig = []types.VAppEntityConfigInfo{*entityConfig}
	updateSpec := types.VAppConfigSpec{
		EntityConfig: cProps.VAppConfig.EntityConfig,
	}
	if err = vappcontainer.Update(container, updateSpec); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereVAppEntityIDString(d))
	return nil
}

func resourceVSphereVAppEntityDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereVAppEntityIDString(d))
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	container, err := vappcontainer.FromID(client, d.Get("container_id").(string))
	if err != nil {
		return err
	}
	cProps, err := vappcontainer.Properties(container)
	if err != nil {
		return err
	}
	vm, err := virtualmachine.FromMOID(client, d.Get("target_id").(string))
	if err != nil {
		return err
	}
	var el []types.VAppEntityConfigInfo
	for _, e := range cProps.VAppConfig.EntityConfig {
		if *e.Key != vm.Reference() {
			el = append(el, e)
		}
	}
	updateSpec := types.VAppConfigSpec{
		EntityConfig: el,
	}
	if err = vappcontainer.Update(container, updateSpec); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Delete finished successfully", resourceVSphereVAppEntityIDString(d))
	return nil
}

func resourceVSphereVAppEntityFind(client *govmomi.Client, id string) (*types.VAppEntityConfigInfo, error) {
	parts := strings.SplitN(id, ":", 2)
	cid := parts[0]
	eid := parts[1]
	container, err := vappcontainer.FromID(client, cid)
	if err != nil {
		return nil, err
	}
	props, err := vappcontainer.Properties(container)
	if err != nil {
		return nil, err
	}
	entity := resourceVSphereVAppEntityFromKey(eid, props)
	return entity, nil
}

// resourceVSphereVAppEntityIDString prints a friendly string for the
// vapp_entity resource.
func resourceVSphereVAppEntityIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereVAppEntityName)
}

func flattenVAppEntityConfigSpec(_ *govmomi.Client, d *schema.ResourceData, obj *types.VAppEntityConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"start_action":   obj.StartAction,
		"start_delay":    obj.StartDelay,
		"start_order":    obj.StartOrder,
		"stop_action":    obj.StopAction,
		"stop_delay":     obj.StopDelay,
		"wait_for_guest": obj.WaitingForGuest,
	})
}

func expandVAppEntityConfigSpec(client *govmomi.Client, d *schema.ResourceData) (*types.VAppEntityConfigInfo, error) {
	_, vmID, err := resourceVSphereVAppEntitySplitID(d.Id())
	if err != nil {
		return nil, err
	}
	vm, err := virtualmachine.FromMOID(client, vmID)
	if err != nil {
		return nil, err
	}
	mor := vm.Reference()
	return &types.VAppEntityConfigInfo{
		Key:             &mor,
		StartAction:     d.Get("start_action").(string),
		StartDelay:      int32(d.Get("start_delay").(int)),
		StartOrder:      int32(d.Get("start_order").(int)),
		StopAction:      d.Get("stop_action").(string),
		StopDelay:       int32(d.Get("stop_delay").(int)),
		WaitingForGuest: structure.GetBoolPtr(d, "wait_for_guest"),
	}, nil
}

func resourceVSphereVAppEntityClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}

func resourceVSphereVAppEntitySplitID(id string) (string, string, error) {
	key := strings.Split(id, ":")
	if len(key) != 2 {
		return "", "", fmt.Errorf("error parsing vApp Entity ID: %s", id)
	}
	return key[0], key[1], nil
}

func resourceVSphereVAppEntityIDFromKeys(vcID, vmID string) string {
	return fmt.Sprintf("%s:%s", vcID, vmID)
}

func resourceVSphereVAppEntityFromKey(key string, c *mo.VirtualApp) *types.VAppEntityConfigInfo {
	log.Printf("[DEBUG] Locating VApp entity with key %s", key)
	for _, e := range c.VAppConfig.EntityConfig {
		if e.Key.Value == key {
			log.Printf("[DEBUG] vApp entity found: %s", key)
			return &e
		}
	}
	return nil
}
