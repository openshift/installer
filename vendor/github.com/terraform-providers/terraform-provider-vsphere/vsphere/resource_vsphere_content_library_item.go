package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"log"
	"strings"
)

func resourceVSphereContentLibraryItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereContentLibraryItemCreate,
		Delete: resourceVSphereContentLibraryItemDelete,
		Read:   resourceVSphereContentLibraryItemRead,
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
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID of source VM of content library item.",
			},
			"type": {
				Type:        schema.TypeString,
				Default:     "ovf",
				Optional:    true,
				ForceNew:    true,
				Description: "Type of content library item.",
			},
		},
	}
}

func resourceVSphereContentLibraryItemRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemRead : Reading Content Library item (%s)", d.Id())
	rc := meta.(*VSphereClient).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("type", item.Type)
	d.Set("library_id", item.LibraryID)
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemRead : Content Library item (%s) read is complete", d.Id())
	return nil
}

func resourceVSphereContentLibraryItemCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemCreate : Beginning Content Library item (%s) creation", d.Get("name").(string))
	rc := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromID(rc, d.Get("library_id").(string))
	if err != nil {
		return err
	}
	files := d.Get("file_url").(*schema.Set)
	id, err := contentlibrary.CreateLibraryItem(rc, lib, d.Get("name").(string), d.Get("description").(string), d.Get("type").(string), files.List())
	if err != nil {
		return err
	}
	d.SetId(id)
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemCreate : Content Library item (%s) creation complete", d.Get("name").(string))
	return resourceVSphereContentLibraryItemRead(d, meta)
}

func resourceVSphereContentLibraryItemDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryItemDelete : Deleting Content Library item (%s)", d.Id())
	rc := meta.(*VSphereClient).restClient
	item, err := contentlibrary.ItemFromID(rc, d.Id())
	if err != nil {
		return err
	}
	return contentlibrary.DeleteLibraryItem(rc, item)
}

func resourceVSphereContentLibraryItemImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*VSphereClient).restClient
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
