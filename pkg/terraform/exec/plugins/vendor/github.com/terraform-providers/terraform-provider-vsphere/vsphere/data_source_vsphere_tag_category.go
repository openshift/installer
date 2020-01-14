package vsphere

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceVSphereTagCategory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereTagCategoryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the category.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the category.",
			},
			"cardinality": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The associated cardinality of the category. Can be one of SINGLE (object can only be assigned one tag in this category) or MULTIPLE (object can be assigned multiple tags in this category).",
			},
			"associable_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Object types to which this category's tags can be attached.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVSphereTagCategoryRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*VSphereClient).TagsClient()
	if err != nil {
		return err
	}

	id, err := tagCategoryByName(client, d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceVSphereTagCategoryRead(d, meta)
}
