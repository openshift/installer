// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileLinkRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_link",
					"profile_id"),
			},
			"link_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the link.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the link.",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_link", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getLinkOptions := &iamidentityv1.GetLinkOptions{}

	getLinkOptions.SetProfileID(d.Get("profile_id").(string))
	getLinkOptions.SetLinkID(d.Get("link_id").(string))

	profileLink, _, err := iamIdentityClient.GetLinkWithContext(context, getLinkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetLinkWithContext failed: %s", err.Error()), "(Data) ibm_iam_trusted_profile_link", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*profileLink.ID)

	if err = d.Set("entity_tag", profileLink.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-entity_tag").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(profileLink.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("modified_at", flex.DateTimeToString(profileLink.ModifiedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting modified_at: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-modified_at").GetDiag()
	}

	if !core.IsNil(profileLink.Name) {
		if err = d.Set("name", profileLink.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-name").GetDiag()
		}
	}

	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cr_type: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-cr_type").GetDiag()
	}

	link := []map[string]interface{}{}
	linkMap, err := DataSourceIBMIamTrustedProfileLinkProfileLinkLinkToMap(profileLink.Link)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_link", "read", "link-to-map").GetDiag()
	}
	link = append(link, linkMap)
	if err = d.Set("link", link); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting link: %s", err), "(Data) ibm_iam_trusted_profile_link", "read", "set-link").GetDiag()
	}

	return nil
}

func DataSourceIBMIamTrustedProfileLinkProfileLinkLinkToMap(model *iamidentityv1.ProfileLinkLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Namespace != nil {
		modelMap["namespace"] = *model.Namespace
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}
