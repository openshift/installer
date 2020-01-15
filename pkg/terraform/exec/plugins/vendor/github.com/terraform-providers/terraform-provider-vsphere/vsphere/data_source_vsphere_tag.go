package vsphere

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceVSphereTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereTagRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the tag.",
				Required:    true,
			},
			"category_id": {
				Type:        schema.TypeString,
				Description: "The unique identifier of the parent category for this tag.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tag.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVSphereTagRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*VSphereClient).TagsClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	categoryID := d.Get("category_id").(string)

	tagID, err := tagByName(client, name, categoryID)
	if err != nil {
		return err
	}

	d.SetId(tagID)
	return resourceVSphereTagRead(d, meta)
}
