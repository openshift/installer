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
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_link",
					"profile_id"),
			},
			"cr_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"link": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Link details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the compute resource.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional name of the Link.",
			},
			"link_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of this link.",
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
		return diag.FromErr(err)
	}

	createLinkOptions := &iamidentityv1.CreateLinkOptions{}
	profile := d.Get("profile_id").(string)
	createLinkOptions.SetProfileID(profile)
	createLinkOptions.SetCrType(d.Get("cr_type").(string))
	link := resourceIBMIamTrustedProfileLinkMapToCreateProfileLinkRequestLink(d.Get("link.0").(map[string]interface{}))
	createLinkOptions.SetLink(&link)
	if _, ok := d.GetOk("name"); ok {
		createLinkOptions.SetName(d.Get("name").(string))
	}

	profileLink, response, err := iamIdentityClient.CreateLink(createLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateLink failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", profile, *profileLink.ID))

	return resourceIBMIamTrustedProfileLinkRead(context, d, meta)
}

func resourceIBMIamTrustedProfileLinkMapToCreateProfileLinkRequestLink(createProfileLinkRequestLinkMap map[string]interface{}) iamidentityv1.CreateProfileLinkRequestLink {
	createProfileLinkRequestLink := iamidentityv1.CreateProfileLinkRequestLink{}

	createProfileLinkRequestLink.CRN = core.StringPtr(createProfileLinkRequestLinkMap["crn"].(string))
	createProfileLinkRequestLink.Namespace = core.StringPtr(createProfileLinkRequestLinkMap["namespace"].(string))
	if createProfileLinkRequestLinkMap["name"] != nil {
		createProfileLinkRequestLink.Name = core.StringPtr(createProfileLinkRequestLinkMap["name"].(string))
	}

	return createProfileLinkRequestLink
}

func resourceIBMIamTrustedProfileLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Invalid ID %s", err))
	}
	getLinkOptions := &iamidentityv1.GetLinkOptions{}

	getLinkOptions.SetProfileID(parts[0])
	getLinkOptions.SetLinkID(parts[1])

	profileLink, response, err := iamIdentityClient.GetLink(getLinkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLink failed %s\n%s", err, response))
	}

	if err = d.Set("profile_id", getLinkOptions.ProfileID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting profile_id: %s", err))
	}
	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting cr_type: %s", err))
	}
	linkMap := resourceIBMIamTrustedProfileLinkCreateProfileLinkRequestLinkToMap(*profileLink.Link)
	if err = d.Set("link", []map[string]interface{}{linkMap}); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting link: %s", err))
	}
	if err = d.Set("name", profileLink.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("link_id", profileLink.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting id: %s", err))
	}
	if err = d.Set("entity_tag", profileLink.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(profileLink.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", flex.DateTimeToString(profileLink.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting modified_at: %s", err))
	}

	return nil
}

func resourceIBMIamTrustedProfileLinkCreateProfileLinkRequestLinkToMap(createProfileLinkRequestLink iamidentityv1.ProfileLinkLink) map[string]interface{} {
	createProfileLinkRequestLinkMap := map[string]interface{}{}

	createProfileLinkRequestLinkMap["crn"] = createProfileLinkRequestLink.CRN
	createProfileLinkRequestLinkMap["namespace"] = createProfileLinkRequestLink.Namespace
	if createProfileLinkRequestLink.Name != nil {
		createProfileLinkRequestLinkMap["name"] = createProfileLinkRequestLink.Name
	}

	return createProfileLinkRequestLinkMap
}

func resourceIBMIamTrustedProfileLinkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Invalid ID %s", err))
	}

	deleteLinkOptions := &iamidentityv1.DeleteLinkOptions{}

	deleteLinkOptions.SetProfileID(parts[0])
	deleteLinkOptions.SetLinkID(parts[1])

	response, err := iamIdentityClient.DeleteLink(deleteLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteLink failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
