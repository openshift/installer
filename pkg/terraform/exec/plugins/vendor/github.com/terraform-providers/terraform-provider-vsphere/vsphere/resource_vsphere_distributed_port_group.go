package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/dvportgroup"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVSphereDistributedPortGroup() *schema.Resource {
	s := map[string]*schema.Schema{
		"distributed_virtual_switch_uuid": {
			Type:        schema.TypeString,
			Description: "The UUID of the DVS to attach this port group to.",
			Required:    true,
			ForceNew:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The generated UUID of the portgroup.",
			Computed:    true,
		},
		// Tagging
		vSphereTagAttributeKey: tagsSchema(),
		// Custom Attributes
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}

	structure.MergeSchema(s, schemaDVPortgroupConfigSpec())

	return &schema.Resource{
		Create: resourceVSphereDistributedPortGroupCreate,
		Read:   resourceVSphereDistributedPortGroupRead,
		Update: resourceVSphereDistributedPortGroupUpdate,
		Delete: resourceVSphereDistributedPortGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereDistributedPortGroupImport,
		},
		Schema: s,
	}
}

func resourceVSphereDistributedPortGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	tagsClient, err := tagsClientIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	dvsID := d.Get("distributed_virtual_switch_uuid").(string)
	dvs, err := dvsFromUUID(client, dvsID)
	if err != nil {
		return fmt.Errorf("could not find DVS %q: %s", dvsID, err)
	}

	spec := expandDVPortgroupConfigSpec(d)
	task, err := dvportgroup.Create(client, dvs, spec)
	if err != nil {
		return fmt.Errorf("error creating portgroup: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	info, err := task.WaitForResult(tctx, nil)
	if err != nil {
		return fmt.Errorf("error waiting for portgroup creation to complete: %s", err)
	}
	pg, err := dvportgroup.FromMOID(client, info.Result.(types.ManagedObjectReference).Value)
	if err != nil {
		return fmt.Errorf("error fetching portgroup after creation: %s", err)
	}
	props, err := dvportgroup.Properties(pg)
	if err != nil {
		return fmt.Errorf("error fetching portgroup properties after creation: %s", err)
	}

	d.SetId(pg.Reference().Value)
	d.Set("key", props.Key)

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, object.NewReference(client.Client, pg.Reference())); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	// Set custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(object.NewReference(client.Client, pg.Reference())); err != nil {
			return err
		}
	}

	return resourceVSphereDistributedPortGroupRead(d, meta)
}

func resourceVSphereDistributedPortGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	pgID := d.Id()
	pg, err := dvportgroup.FromMOID(client, pgID)
	if err != nil {
		return fmt.Errorf("could not find portgroup %q: %s", pgID, err)
	}
	props, err := dvportgroup.Properties(pg)
	if err != nil {
		return fmt.Errorf("error fetching portgroup properties: %s", err)
	}

	d.Set("key", props.Key)

	if err := flattenDVPortgroupConfigInfo(d, props.Config); err != nil {
		return err
	}

	if tagsClient, _ := meta.(*VSphereClient).TagsClient(); tagsClient != nil {
		if err := readTagsForResource(tagsClient, pg, d); err != nil {
			return fmt.Errorf("error reading tags: %s", err)
		}
	}

	if customattribute.IsSupported(client) {
		customattribute.ReadFromResource(client, props.Entity(), d)
	}

	return nil
}

func resourceVSphereDistributedPortGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	tagsClient, err := tagsClientIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	pgID := d.Id()
	pg, err := dvportgroup.FromMOID(client, pgID)
	if err != nil {
		return fmt.Errorf("could not find portgroup %q: %s", pgID, err)
	}
	spec := expandDVPortgroupConfigSpec(d)
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	task, err := pg.Reconfigure(ctx, spec)
	if err != nil {
		return fmt.Errorf("error reconfiguring portgroup: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return fmt.Errorf("error waiting for portgroup update to complete: %s", err)
	}

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, object.NewReference(client.Client, pg.Reference())); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	// Update custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(object.NewReference(client.Client, pg.Reference())); err != nil {
			return fmt.Errorf("error updating custom attributes: %s", err)
		}
	}

	return resourceVSphereDistributedPortGroupRead(d, meta)
}

func resourceVSphereDistributedPortGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	pgID := d.Id()
	pg, err := dvportgroup.FromMOID(client, pgID)
	if err != nil {
		return fmt.Errorf("could not find portgroup %q: %s", pgID, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	task, err := pg.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("error deleting portgroup: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return fmt.Errorf("error waiting for portgroup deletion to complete: %s", err)
	}
	return nil
}

func resourceVSphereDistributedPortGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// We use the inventory path to the portgroup to import. There is not
	// checking to make sure that it belongs to the configured DVS, but on
	// subsequent plans, if it is not, the resource will be in an unusable state
	// as all query calls for DVS CRUD calls require the correct DVS UUID in
	// addition to the portgroup UUID.
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	p := d.Id()
	pg, err := dvportgroup.FromPath(client, p, nil)
	if err != nil {
		return nil, fmt.Errorf("error locating portgroup: %s", err)
	}
	props, err := dvportgroup.Properties(pg)
	if err != nil {
		return nil, fmt.Errorf("error fetching portgroup properties: %s", err)
	}
	d.SetId(props.Key)

	// We need to populate the DVS UUID here as well or else our read calls will
	// fail.
	dvsID := props.Config.DistributedVirtualSwitch.Value
	dvs, err := dvsFromMOID(client, dvsID)
	if err != nil {
		return nil, fmt.Errorf("error getting DVS with ID %q: %s", dvsID, err)
	}
	dvProps, err := dvsProperties(dvs)
	if err != nil {
		return nil, fmt.Errorf("error fetching DVS properties: %s", err)
	}

	d.Set("distributed_virtual_switch_uuid", dvProps.Uuid)

	return []*schema.ResourceData{d}, nil
}
