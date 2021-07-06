// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMResourceTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceTagRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_resource_tag", resourceID),
				Description:  "CRN of the resource on which the tags should be attached",
			},
			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_resource_tag", tags)},
				Set:         resourceIBMVPCHash,
				Description: "List of tags associated with resource instance",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource type on which the tags should be fetched",
			},
		},
	}
}

func dataSourceIBMResourceTagRead(d *schema.ResourceData, meta interface{}) error {
	var rID, rType string
	rID = d.Get("resource_id").(string)
	if v, ok := d.GetOk(resourceType); ok && v != nil {
		rType = v.(string)
	}

	tags, err := GetGlobalTagsUsingCRN(meta, rID, rType, "")
	if err != nil {
		return fmt.Errorf(
			"Error on get of resource tags (%s) tags: %s", d.Id(), err)
	}

	d.SetId(rID)
	d.Set("resource_id", rID)
	d.Set("resource_type", rType)
	d.Set("tags", tags)

	return nil
}
