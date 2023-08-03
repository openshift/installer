// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileLinkRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_link",
					"profile_id"),
			},
			"link_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the link.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the claim rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional name of the Link.",
			},
			"cr_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"link": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the compute resource.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
						},
						"name": {
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

func DataSourceIBMIamTrustedProfileLinkValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIamTrustedProfileLinkValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_link", Schema: validateSchema}
	return &iBMIamTrustedProfileLinkValidator
}

func dataSourceIBMIamTrustedProfileLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(profileLink.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", flex.DateTimeToString(profileLink.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting modified_at: %s", err))
	}
	if err = d.Set("name", profileLink.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting cr_type: %s", err))
	}

	if profileLink.Link != nil {
		err = d.Set("link", dataSourceProfileLinkFlattenLink(*profileLink.Link))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting link %s", err))
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
