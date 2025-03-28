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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileLinkListRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_links",
					"profile_id"),
			},
			"links": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of links to a trusted profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unique identifier of the link.",
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
				},
			},
		},
	}
}

func DataSourceIBMIamTrustedProfileLinksValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIamTrustedProfileLinksValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_links", Schema: validateSchema}
	return &iBMIamTrustedProfileLinksValidator
}

func dataSourceIBMIamTrustedProfileLinkListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_links", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listLinksOptions := &iamidentityv1.ListLinksOptions{}

	listLinksOptions.SetProfileID(d.Get("profile_id").(string))

	profileLinkList, _, err := iamIdentityClient.ListLinksWithContext(context, listLinksOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListLinksWithContext failed: %s", err.Error()), "(Data) ibm_iam_trusted_profile_links", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIamTrustedProfileLinksID(d))

	links := []map[string]interface{}{}
	for _, linksItem := range profileLinkList.Links {
		linksItemMap, err := DataSourceIBMIamTrustedProfileLinksProfileLinkToMap(&linksItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_links", "read", "links-to-map").GetDiag()
		}
		links = append(links, linksItemMap)
	}
	if err = d.Set("links", links); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting links: %s", err), "(Data) ibm_iam_trusted_profile_links", "read", "set-links").GetDiag()
	}

	return nil
}

// dataSourceIBMIamTrustedProfileLinksID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfileLinksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIamTrustedProfileLinksProfileLinkToMap(model *iamidentityv1.ProfileLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["entity_tag"] = *model.EntityTag
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["modified_at"] = model.ModifiedAt.String()
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	modelMap["cr_type"] = *model.CrType
	linkMap, err := DataSourceIBMIamTrustedProfileLinksProfileLinkLinkToMap(model.Link)
	if err != nil {
		return modelMap, err
	}
	modelMap["link"] = []map[string]interface{}{linkMap}
	return modelMap, nil
}

func DataSourceIBMIamTrustedProfileLinksProfileLinkLinkToMap(model *iamidentityv1.ProfileLinkLink) (map[string]interface{}, error) {
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
