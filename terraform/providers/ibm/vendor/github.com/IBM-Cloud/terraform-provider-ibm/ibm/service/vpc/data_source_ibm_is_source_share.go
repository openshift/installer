// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsSourceShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsSourceShareRead,

		Schema: map[string]*schema.Schema{
			"share_replica": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The replica file share identifier.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this share.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the share.",
			},
		},
	}
}

func dataSourceIbmIsSourceShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	replicaShareId := d.Get("share_replica").(string)

	getShareSourceOptions := &vpcv1.GetShareSourceOptions{
		ShareID: &replicaShareId,
	}

	share, response, err := vpcClient.GetShareSourceWithContext(context, getShareSourceOptions)
	if err != nil {
		log.Printf("[DEBUG] GetShareSourceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetShareSourceWithContext failed %s\n%s", err, response))
	}

	d.SetId(*share.ID)

	if err = d.Set("crn", share.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("href", share.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("name", share.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("resource_type", share.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}
