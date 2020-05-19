package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"log"
	"strings"
)

func resourceVSphereContentLibrary() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereContentLibraryCreate,
		Delete: resourceVSphereContentLibraryDelete,
		Read:   resourceVSphereContentLibraryRead,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereContentLibraryImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional description of the content library.",
			},
			"storage_backing": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the content library.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVSphereContentLibraryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryRead : Beginning Content Library (%s) read", d.Id())
	c := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromID(c, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	d.SetId(lib.ID)
	sb := contentlibrary.FlattenStorageBackings(lib.Storage)
	d.Set("name", lib.Name)
	d.Set("description", lib.Description)
	err = d.Set("storage_backing", sb)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] resourceVSphereContentLibraryRead : Content Library (%s) read is complete", d.Id())
	return nil
}

func resourceVSphereContentLibraryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryCreate : Beginning Content Library (%s) creation", d.Get("name").(string))
	vc := meta.(*VSphereClient).vimClient
	rc := meta.(*VSphereClient).restClient
	backings, err := contentlibrary.ExpandStorageBackings(vc, d)
	if err != nil {
		return err
	}
	id, err := contentlibrary.CreateLibrary(rc, d.Get("name").(string), d.Get("description").(string), backings)
	if err != nil {
		return err
	}
	d.SetId(id)
	log.Printf("[DEBUG] resourceVSphereContentLibraryCreate : Content Library (%s) creation is complete", d.Get("name").(string))
	return resourceVSphereContentLibraryRead(d, meta)
}

func resourceVSphereContentLibraryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryDelete : Deleting Content Library (%s)", d.Id())
	c := meta.(*VSphereClient).restClient
	lib, err := contentlibrary.FromID(c, d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] resourceVSphereContentLibraryDelete : Content Library (%s) deleted", d.Id())
	return contentlibrary.DeleteLibrary(c, lib)
}

func resourceVSphereContentLibraryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*VSphereClient).restClient
	_, err := contentlibrary.FromID(client, d.Id())
	if err != nil {
		return nil, err
	}
	err = resourceVSphereContentLibraryRead(d, meta)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
