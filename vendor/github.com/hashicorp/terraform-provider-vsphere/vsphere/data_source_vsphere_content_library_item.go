package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
)

func dataSourceVSphereContentLibraryItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereContentLibraryItemRead,
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
				Description: "ID of the content library to contain item.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of content library item.",
			},
		},
	}
}

func dataSourceVSphereContentLibraryItemRead(d *schema.ResourceData, meta interface{}) error {
	rc := meta.(*VSphereClient).restClient
	lib, _ := contentlibrary.FromID(rc, d.Get("library_id").(string))
	item, err := contentlibrary.ItemFromName(rc, lib, d.Get("name").(string))
	if err != nil {
		return provider.ProviderError(d.Get("name").(string), "dataSourceVSphereContentLibraryItemRead", err)
	}
	d.Set("type", item.Type)
	d.SetId(item.ID)
	return nil
}
