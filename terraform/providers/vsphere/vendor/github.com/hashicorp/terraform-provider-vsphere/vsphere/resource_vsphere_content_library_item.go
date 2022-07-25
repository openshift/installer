package vsphere

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
)

func resourceVSphereContentLibraryItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereContentLibraryItemCreate,
		Delete: resourceVSphereContentLibraryItemDelete,
		Read:   resourceVSphereContentLibraryItemRead,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourceVSphereContentLibraryItemResourceV0().CoreConfigSchema().ImpliedType(),
			Upgrade: resourceVSphereContentLibraryItemUpgradeV0,
		}},
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereContentLibraryItemImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library item.",
			},
			"library_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the content library to contain item",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional description of the content library item.",
			},
			"file_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "ID of source VM of content library item.",
				ConflictsWith: []string{"source_uuid"},
			},
			"type": {
				Type:        schema.TypeString,
				Default:     "ovf",
				Optional:    true,
				ForceNew:    true,
				Description: "Type of content library item.",
			},
			"source_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The managed object ID of an existing VM to be cloned to the content library.",
				ConflictsWith: []string{"file_url"},
			},
		},
	}
}

func resourceVSphereContentLibraryItemUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if len(rawState["file_url"].([]interface{})) < 1 {
		rawState["file_url"] = interface{}("")
		return rawState, nil
	}

	for _, file := range rawState["file_url"].([]interface{}) {
		if strings.HasSuffix(file.(string), "ova") || strings.HasSuffix(file.(string), "ovf") || strings.HasSuffix(file.(string), "iso") {
			rawState["file_url"] = interface{}(file.(string))
			return rawState, nil
		}
	}

	rawState["file_url"] = rawState["file_url"].([]interface{})[0]
	return rawState, nil
}

func resourceVSphereContentLibraryItemRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemRead : Reading Content Library item (%s)", d.Id())
	rc := meta.(*Client).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("name", item.Name)
	_ = d.Set("description", item.Description)
	_ = d.Set("type", item.Type)
	_ = d.Set("library_id", item.LibraryID)
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemRead : Content Library item (%s) read is complete", d.Id())
	return nil
}

func resourceVSphereContentLibraryItemCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemCreate : Beginning Content Library item (%s) creation", d.Get("name").(string))
	rc := meta.(*Client).restClient
	lib, err := contentlibrary.FromID(rc, d.Get("library_id").(string))
	if err != nil {
		return err
	}
	var moid virtualmachine.MOIDForUUIDResult
	if uuid, ok := d.GetOk("source_uuid"); ok {
		moid, err = virtualmachine.MOIDForUUID(meta.(*Client).vimClient, uuid.(string))
		if err != nil {
			return err
		}
	}
	id, err := contentlibrary.CreateLibraryItem(rc, lib, d.Get("name").(string), d.Get("description").(string), d.Get("type").(string), d.Get("file_url").(string), moid.MOID)
	if err != nil {
		return err
	}
	d.SetId(*id)
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemCreate : Content Library item (%s) creation complete", d.Get("name").(string))
	return resourceVSphereContentLibraryItemRead(d, meta)
}

func resourceVSphereContentLibraryItemDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemDelete : Deleting Content Library item (%s)", d.Id())
	rc := meta.(*Client).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		return err
	}
	return contentlibrary.DeleteLibraryItem(rc, item)
}

func resourceVSphereContentLibraryItemImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client).restClient
	_, err := contentlibrary.ItemFromID(client, d.Id())
	if err != nil {
		return nil, err
	}
	err = resourceVSphereContentLibraryItemRead(d, meta)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
