// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsShareProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareProfileRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share profile name.",
			},
			"family": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this share profile belongs to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share profile.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIbmIsShareProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareProfileOptions := &vpcbetav1.GetShareProfileOptions{}

	getShareProfileOptions.SetName(d.Get("name").(string))

	shareProfile, response, err := vpcClient.GetShareProfileWithContext(context, getShareProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] GetShareProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*shareProfile.Name)
	if err = d.Set("family", shareProfile.Family); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting family: %s", err))
	}
	if err = d.Set("href", shareProfile.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("resource_type", shareProfile.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

// dataSourceIbmIsShareProfileID returns a reasonable ID for the list.
func dataSourceIbmIsShareProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
