// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIBMIamTrustedProfileLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileLinkRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
			},
			"link_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the link.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the claim rule.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional name of the Link.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"link": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the compute resource.",
						},
						"namespace": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfileLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getLinkOptions := &iamidentityv1.GetLinkOptions{}

	getLinkOptions.SetProfileID(d.Get("profile_id").(string))
	getLinkOptions.SetLinkID(d.Get("link_id").(string))

	profileLink, response, err := iamIdentityClient.GetLink(getLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLink failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getLinkOptions.ProfileID, *getLinkOptions.LinkID))
	if err = d.Set("entity_tag", profileLink.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(profileLink.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", dateTimeToString(profileLink.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}
	if err = d.Set("name", profileLink.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cr_type: %s", err))
	}

	if profileLink.Link != nil {
		err = d.Set("link", dataSourceProfileLinkFlattenLink(*profileLink.Link))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting link %s", err))
		}
	}

	return nil
}

func dataSourceProfileLinkFlattenLink(result iamidentityv1.ProfileLinkLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceProfileLinkLinkToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceProfileLinkLinkToMap(linkItem iamidentityv1.ProfileLinkLink) (linkMap map[string]interface{}) {
	linkMap = map[string]interface{}{}

	if linkItem.CRN != nil {
		linkMap["crn"] = linkItem.CRN
	}
	if linkItem.Namespace != nil {
		linkMap["namespace"] = linkItem.Namespace
	}
	if linkItem.Name != nil {
		linkMap["name"] = linkItem.Name
	}

	return linkMap
}
