package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVSphereDistributedVirtualSwitch() *schema.Resource {
	s := map[string]*schema.Schema{
		"datacenter_id": {
			Type:        schema.TypeString,
			Description: "The ID of the datacenter to create this virtual switch in.",
			Required:    true,
			ForceNew:    true,
		},
		"folder": {
			Type:        schema.TypeString,
			Description: "The folder to create this virtual switch in, relative to the datacenter.",
			Optional:    true,
			ForceNew:    true,
		},
		"network_resource_control_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether or not to enable network resource control, enabling advanced traffic shaping and resource control features.",
			Optional:    true,
		},
		// Tagging
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	structure.MergeSchema(s, schemaDVSCreateSpec())

	return &schema.Resource{
		Create: resourceVSphereDistributedVirtualSwitchCreate,
		Read:   resourceVSphereDistributedVirtualSwitchRead,
		Update: resourceVSphereDistributedVirtualSwitchUpdate,
		Delete: resourceVSphereDistributedVirtualSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereDistributedVirtualSwitchImport,
		},
		Schema: s,
	}
}

func resourceVSphereDistributedVirtualSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	dc, err := datacenterFromID(client, d.Get("datacenter_id").(string))
	if err != nil {
		return fmt.Errorf("cannot locate datacenter: %s", err)
	}
	fo, err := folder.FromPath(client, d.Get("folder").(string), folder.VSphereFolderTypeNetwork, dc)
	if err != nil {
		return fmt.Errorf("cannot locate folder: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	spec := expandDVSCreateSpec(d)
	task, err := fo.CreateDVS(ctx, spec)
	if err != nil {
		return fmt.Errorf("error creating DVS: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	info, err := task.WaitForResult(tctx, nil)
	if err != nil {
		return fmt.Errorf("error waiting for DVS creation to complete: %s", err)
	}

	dvs, err := dvsFromMOID(client, info.Result.(types.ManagedObjectReference).Value)
	if err != nil {
		return fmt.Errorf("error fetching DVS after creation: %s", err)
	}
	props, err := dvsProperties(dvs)
	if err != nil {
		return fmt.Errorf("error fetching DVS properties after creation: %s", err)
	}

	d.SetId(props.Uuid)

	// Enable network resource I/O control if it needs to be enabled
	if d.Get("network_resource_control_enabled").(bool) {
		enableDVSNetworkResourceManagement(client, dvs, true)
	}

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, object.NewReference(client.Client, dvs.Reference())); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	// Set custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(dvs); err != nil {
			return err
		}
	}

	return resourceVSphereDistributedVirtualSwitchRead(d, meta)
}

func resourceVSphereDistributedVirtualSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	id := d.Id()
	dvs, err := dvsFromUUID(client, id)
	if err != nil {
		return fmt.Errorf("could not find DVS %q: %s", id, err)
	}
	props, err := dvsProperties(dvs)
	if err != nil {
		return fmt.Errorf("error fetching DVS properties: %s", err)
	}

	// Set the datacenter ID, for completion's sake when importing
	dcp, err := folder.RootPathParticleNetwork.SplitDatacenter(dvs.InventoryPath)
	if err != nil {
		return fmt.Errorf("error parsing datacenter from inventory path: %s", err)
	}
	dc, err := getDatacenter(client, dcp)
	if err != nil {
		return fmt.Errorf("error locating datacenter: %s", err)
	}
	d.Set("datacenter_id", dc.Reference().Value)

	// Set the folder
	f, err := folder.RootPathParticleNetwork.SplitRelativeFolder(dvs.InventoryPath)
	if err != nil {
		return fmt.Errorf("error parsing DVS path %q: %s", dvs.InventoryPath, err)
	}
	d.Set("folder", folder.NormalizePath(f))

	// Read in config info
	if err := flattenVMwareDVSConfigInfo(d, props.Config.(*types.VMwareDVSConfigInfo)); err != nil {
		return err
	}

	// Read tags if we have the ability to do so
	if tagsClient, _ := meta.(*VSphereClient).TagsManager(); tagsClient != nil {
		if err := readTagsForResource(tagsClient, dvs, d); err != nil {
			return fmt.Errorf("error reading tags: %s", err)
		}
	}

	// Read set custom attributes
	if customattribute.IsSupported(client) {
		customattribute.ReadFromResource(client, props.Entity(), d)
	}

	return nil
}

func resourceVSphereDistributedVirtualSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	id := d.Id()
	dvs, err := dvsFromUUID(client, id)
	if err != nil {
		return fmt.Errorf("could not find DVS %q: %s", id, err)
	}

	// If we have a pending version upgrade, do that first.
	if d.HasChange("version") {
		old, new := d.GetChange("version")
		var ovi, nvi int
		for n, v := range dvsVersions {
			if old.(string) == v {
				ovi = n
			}
			if new.(string) == v {
				nvi = n
			}
		}
		if nvi < ovi {
			return fmt.Errorf("downgrading dvSwitches are not allowed (old: %s new: %s)", old, new)
		}
		if err := upgradeDVS(client, dvs, new.(string)); err != nil {
			return fmt.Errorf("could not upgrade DVS: %s", err)
		}
		props, err := dvsProperties(dvs)
		if err != nil {
			return fmt.Errorf("could not get DVS properties after upgrade: %s", err)
		}
		// ConfigVersion increments after a DVS upgrade, which means this needs to
		// be updated before the post-update read to ensure that we don't run into
		// ConcurrentAccess errors on the update operation below.
		d.Set("config_version", props.Config.(*types.VMwareDVSConfigInfo).ConfigVersion)
	}

	spec := expandVMwareDVSConfigSpec(d)
	if err := updateDVSConfiguration(client, dvs, spec); err != nil {
		return fmt.Errorf("could not update DVS: %s", err)
	}

	// Modify network I/O control if necessary
	if d.HasChange("network_resource_control_enabled") {
		enableDVSNetworkResourceManagement(client, dvs, d.Get("network_resource_control_enabled").(bool))
	}

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, object.NewReference(client.Client, dvs.Reference())); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	// Apply custom attribute updates
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(dvs); err != nil {
			return err
		}
	}

	return resourceVSphereDistributedVirtualSwitchRead(d, meta)
}

func resourceVSphereDistributedVirtualSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return err
	}
	id := d.Id()
	dvs, err := dvsFromUUID(client, id)
	if err != nil {
		return fmt.Errorf("could not find DVS %q: %s", id, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	task, err := dvs.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("error deleting DVS: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return fmt.Errorf("error waiting for DVS deletion to complete: %s", err)
	}

	return nil
}

func resourceVSphereDistributedVirtualSwitchImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Due to the relative difficulty in trying to fetch a DVS's UUID, we use the
	// inventory path to the DVS instead, and just run it through finder. A full
	// path is required unless the default datacenter can be utilized.
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	p := d.Id()
	dvs, err := dvsFromPath(client, p, nil)
	if err != nil {
		return nil, fmt.Errorf("error locating DVS: %s", err)
	}
	props, err := dvsProperties(dvs)
	if err != nil {
		return nil, fmt.Errorf("error fetching DVS properties: %s", err)
	}
	d.SetId(props.Uuid)
	return []*schema.ResourceData{d}, nil
}
