package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
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
			"publication": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   "Publication configuration for content library.",
				ConflictsWith: []string{"subscription"},
				MaxItems:      1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"authentication_method": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "NONE",
					},
					"username": {
						Type:     schema.TypeString,
						ForceNew: true,
						Computed: true,
						Optional: true,
					},
					"password": {
						Type:     schema.TypeString,
						ForceNew: true,
						Computed: true,
						Optional: true,
					},
					"published": {
						Type:     schema.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  false,
					},
					"publish_url": {
						Type:     schema.TypeString,
						ForceNew: true,
						Computed: true,
					},
				},
				},
			},
			"subscription": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Description:   "Publication configuration for content library.",
				ConflictsWith: []string{"publication"},
				MaxItems:      1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"authentication_method": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "NONE",
					},
					"subscription_url": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "NONE",
					},
					"username": {
						Type:     schema.TypeString,
						ForceNew: true,
						Computed: true,
						Optional: true,
					},
					"password": {
						Type:     schema.TypeString,
						ForceNew: true,
						Computed: true,
						Optional: true,
					},
					"on_demand": {
						Type:     schema.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  true,
					},
					"automatic_sync": {
						Type:     schema.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  false,
					},
				},
				},
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
		}
		return err
	}
	d.SetId(lib.ID)
	if err = contentlibrary.FlattenPublication(d, lib.Publication); err != nil {
		return err
	}
	if err = contentlibrary.FlattenSubscription(d, lib.Subscription); err != nil {
		return err
	}
	if err = contentlibrary.FlattenStorageBackings(d, lib.Storage); err != nil {
		return err
	}
	d.Set("name", lib.Name)
	d.Set("description", lib.Description)
	log.Printf("[DEBUG] resourceVSphereContentLibraryRead : Content Library (%s) read is complete", d.Id())
	return nil
}

func resourceVSphereContentLibraryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceVSphereContentLibraryCreate : Beginning Content Library (%s) creation", d.Get("name").(string))
	vimClient := meta.(*VSphereClient).vimClient
	restClient := meta.(*VSphereClient).restClient
	backings, err := contentlibrary.ExpandStorageBackings(vimClient, d)
	if err != nil {
		return err
	}
	id, err := contentlibrary.CreateLibrary(d, restClient, backings)
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
