package vsphere

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceVSphereContentLibraryItemResourceV0() *schema.Resource {
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
