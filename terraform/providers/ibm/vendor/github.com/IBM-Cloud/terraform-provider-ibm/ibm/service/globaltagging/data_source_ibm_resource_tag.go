// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMResourceTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceTagRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_tag", "resource_id"),
				Description:  "CRN of the resource on which the tags should be attached",
			},
			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_resource_tag", tags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags associated with resource instance",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource type on which the tags should be fetched",
			},
			"tag_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_tag", "tag_type"),
				Description:  "Tag type on which the tags should be fetched",
				Default:      "user",
			},
		},
	}
}

func dataSourceIBMResourceTagRead(d *schema.ResourceData, meta interface{}) error {
	var rID, rType string

	if r, ok := d.GetOk("resource_id"); ok && r != nil {
		rID = r.(string)
	}
	if v, ok := d.GetOk(resourceType); ok && v != nil {
		rType = v.(string)
	}
	tType := ""
	if t, ok := d.GetOk("tag_type"); ok && t != nil {
		tType = t.(string)
	}

	tags, err := flex.GetGlobalTagsUsingCRN(meta, rID, rType, tType)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error on get of resource tags (%s) tags: %s", d.Id(), err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("resource_id", rID)
	d.Set("resource_type", rType)
	d.Set("tags", tags)
	d.Set("tag_type", tType)
	return nil
}
