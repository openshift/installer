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

func ResourceIBMIAMTrustedProfileLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamTrustedProfileLinkCreate,
		ReadContext:   resourceIBMIamTrustedProfileLinkRead,
		DeleteContext: resourceIBMIamTrustedProfileLinkDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_link",
					"profile_id"),
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional name of the Link.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"link": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Link details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CRN of the compute resource.",
						},
						"namespace": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.",
						},
					},
				},
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
			"link_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the unique identifier of the link.",
			},
		},
	}
}

func ResourceIBMIAMTrustedProfileLinkValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIAMTrustedProfileLinkValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_link", Schema: validateSchema}
	return &iBMIAMTrustedProfileLinkValidator
}

func resourceIBMIamTrustedProfileLinkCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createLinkOptions := &iamidentityv1.CreateLinkOptions{}

	createLinkOptions.SetProfileID(d.Get("profile_id").(string))
	createLinkOptions.SetCrType(d.Get("cr_type").(string))
	linkModel, err := ResourceIBMIamTrustedProfileLinkMapToCreateProfileLinkRequestLink(d.Get("link.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "create", "parse-link").GetDiag()
	}
	createLinkOptions.SetLink(linkModel)
	if _, ok := d.GetOk("name"); ok {
		createLinkOptions.SetName(d.Get("name").(string))
	}

	profileLink, _, err := iamIdentityClient.CreateLinkWithContext(context, createLinkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateLinkWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_link", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createLinkOptions.ProfileID, *profileLink.ID))

	return resourceIBMIamTrustedProfileLinkRead(context, d, meta)
}

func resourceIBMIamTrustedProfileLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getLinkOptions := &iamidentityv1.GetLinkOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "sep-id-parts").GetDiag()
	}

	getLinkOptions.SetProfileID(parts[0])
	getLinkOptions.SetLinkID(parts[1])

	profileLink, response, err := iamIdentityClient.GetLinkWithContext(context, getLinkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetLinkWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_link", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(profileLink.Name) {
		if err = d.Set("name", profileLink.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		err = fmt.Errorf("Error setting cr_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-cr_type").GetDiag()
	}
	if err = d.Set("profile_id", getLinkOptions.ProfileID); err != nil {
		err = fmt.Errorf("Error setting profile_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-profile_id").GetDiag()
	}
	linkMap, err := ResourceIBMIamTrustedProfileLinkProfileLinkLinkToMap(profileLink.Link)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "link-to-map").GetDiag()
	}
	if err = d.Set("link", []map[string]interface{}{linkMap}); err != nil {
		err = fmt.Errorf("Error setting link: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-link").GetDiag()
	}
	if err = d.Set("entity_tag", profileLink.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-entity_tag").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(profileLink.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("modified_at", flex.DateTimeToString(profileLink.ModifiedAt)); err != nil {
		err = fmt.Errorf("Error setting modified_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-modified_at").GetDiag()
	}
	if err = d.Set("link_id", profileLink.ID); err != nil {
		err = fmt.Errorf("Error setting link_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "read", "set-link_id").GetDiag()
	}

	return nil
}

func resourceIBMIamTrustedProfileLinkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteLinkOptions := &iamidentityv1.DeleteLinkOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_link", "delete", "sep-id-parts").GetDiag()
	}

	deleteLinkOptions.SetProfileID(parts[0])
	deleteLinkOptions.SetLinkID(parts[1])

	_, err = iamIdentityClient.DeleteLinkWithContext(context, deleteLinkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteLinkWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_link", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIamTrustedProfileLinkMapToCreateProfileLinkRequestLink(modelMap map[string]interface{}) (*iamidentityv1.CreateProfileLinkRequestLink, error) {
	model := &iamidentityv1.CreateProfileLinkRequestLink{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	model.Namespace = core.StringPtr(modelMap["namespace"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIamTrustedProfileLinkProfileLinkLinkToMap(model *iamidentityv1.ProfileLinkLink) (map[string]interface{}, error) {
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
