// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIamAccessTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIamAccessTag,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_tag", "name"),
				Set:          flex.ResourceIBMVPCHash,
				Description:  "Name of the access tag to be fetched",
			},
			tagType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the tag(access)",
			},
		},
	}
}

func dataSourceIBMIamAccessTag(d *schema.ResourceData, meta interface{}) error {
	tagName := ""
	if t, ok := d.GetOk("name"); ok && t != nil {
		tagName = t.(string)
	}
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return err
	}
	accessTagType := "access"
	listTagsOptions := &globaltaggingv1.ListTagsOptions{
		TagType: &accessTagType,
	}
	taggingResult, _, err := gtClient.ListTags(listTagsOptions)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error retrieving access tag (%s): %s", tagName, err)
	}

	var taglist []string
	for _, item := range taggingResult.Items {
		taglist = append(taglist, *item.Name)
	}
	existingAccessTags := flex.NewStringSet(flex.ResourceIBMVPCHash, taglist)
	if !existingAccessTags.Contains(tagName) {
		return flex.FmtErrorf("[ERROR] Access tag %s not found", tagName)
	}
	d.SetId(tagName)
	d.Set("name", tagName)
	d.Set(tagType, accessTagType)
	return nil
}
