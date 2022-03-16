package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
)

func dataSourceVSphereFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereFolderRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Description: "The absolute path of the folder.",
				Required:    true,
			},
		},
	}
}

func dataSourceVSphereFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	fo, err := folder.FromAbsolutePath(client, d.Get("path").(string))
	if err != nil {
		return fmt.Errorf("cannot locate folder: %s", err)
	}

	d.SetId(fo.Reference().Value)

	return nil
}
