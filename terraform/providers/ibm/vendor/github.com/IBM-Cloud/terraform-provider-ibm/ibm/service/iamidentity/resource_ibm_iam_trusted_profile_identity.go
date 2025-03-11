// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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

func ResourceIBMIamTrustedProfileIdentity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamTrustedProfileIdentityCreate,
		ReadContext:   resourceIBMIamTrustedProfileIdentityRead,
		DeleteContext: resourceIBMIamTrustedProfileIdentityDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the trusted profile.",
			},
			"identity_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_identity", "identity_type"),
				Description:  "Type of the identity.",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn' it uses account id contained in the CRN.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_identity", "type"),
				Description:  "Type of the identity.",
			},
			"accounts": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Only valid for the type user. Accounts from which a user can assume the trusted profile.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the identity that can assume the trusted profile. This is optional field for all the types of identities. When this field is not set for the identity type 'serviceid' then the description of the service id is used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.",
			},
		},
	}
}

func ResourceIBMIamTrustedProfileIdentityValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identity_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "crn, serviceid, user",
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "crn, serviceid, user",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_identity", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIamTrustedProfileIdentityCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	setProfileIdentityOptions := &iamidentityv1.SetProfileIdentityOptions{}

	setProfileIdentityOptions.SetProfileID(d.Get("profile_id").(string))
	setProfileIdentityOptions.SetIdentityType(d.Get("identity_type").(string))
	setProfileIdentityOptions.SetIdentifier(d.Get("identifier").(string))
	setProfileIdentityOptions.SetType(d.Get("type").(string))
	if _, ok := d.GetOk("accounts"); ok {
		var accounts []string
		for _, v := range d.Get("accounts").([]interface{}) {
			accountsItem := v.(string)
			accounts = append(accounts, accountsItem)
		}
		setProfileIdentityOptions.SetAccounts(accounts)
	}
	if _, ok := d.GetOk("description"); ok {
		setProfileIdentityOptions.SetDescription(d.Get("description").(string))
	}

	profileIdentityResponse, response, err := iamIdentityClient.SetProfileIdentityWithContext(context, setProfileIdentityOptions)
	if err != nil {
		log.Printf("[DEBUG] SetProfileIdentityWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("SetProfileIdentityWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s|%s|%s", *setProfileIdentityOptions.ProfileID, *setProfileIdentityOptions.IdentityType, *profileIdentityResponse.Identifier))

	return resourceIBMIamTrustedProfileIdentityRead(context, d, meta)
}

func resourceIBMIamTrustedProfileIdentityRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileIdentityOptions := &iamidentityv1.GetProfileIdentityOptions{}

	parts_to_use, original_err := flex.SepIdParts(d.Id(), "|")
	if original_err != nil {
		parts, err := flex.SepIdParts(d.Id(), "/") // compatability - can be removed in future release
		if err != nil {
			return diag.FromErr(original_err)
		}
		parts_to_use = parts
	}

	getProfileIdentityOptions.SetProfileID(parts_to_use[0])
	getProfileIdentityOptions.SetIdentityType(parts_to_use[1])
	getProfileIdentityOptions.SetIdentifierID(parts_to_use[2])

	profileIdentityResponse, response, err := iamIdentityClient.GetProfileIdentityWithContext(context, getProfileIdentityOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProfileIdentityWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfileIdentityWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s|%s|%s", *getProfileIdentityOptions.ProfileID, *getProfileIdentityOptions.IdentityType, *getProfileIdentityOptions.IdentifierID))

	if err = d.Set("profile_id", getProfileIdentityOptions.ProfileID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile_id: %s", err))
	}
	if err = d.Set("identity_type", getProfileIdentityOptions.IdentityType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting identity_type: %s", err))
	}
	if err = d.Set("identifier", profileIdentityResponse.Identifier); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting identifier: %s", err))
	}
	if err = d.Set("type", profileIdentityResponse.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if !core.IsNil(profileIdentityResponse.Accounts) {
		if err = d.Set("accounts", profileIdentityResponse.Accounts); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting accounts: %s", err))
		}
	}
	if !core.IsNil(profileIdentityResponse.Description) {
		if err = d.Set("description", profileIdentityResponse.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}

	return nil
}

func resourceIBMIamTrustedProfileIdentityDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProfileIdentityOptions := &iamidentityv1.DeleteProfileIdentityOptions{}

	parts_to_use, original_err := flex.SepIdParts(d.Id(), "|")
	if original_err != nil {
		parts, err := flex.SepIdParts(d.Id(), "/") // compatability - remove in future release
		if err != nil {
			return diag.FromErr(original_err)
		}
		parts_to_use = parts
	}

	deleteProfileIdentityOptions.SetProfileID(parts_to_use[0])
	deleteProfileIdentityOptions.SetIdentityType(parts_to_use[1])
	deleteProfileIdentityOptions.SetIdentifierID(parts_to_use[2])

	response, err := iamIdentityClient.DeleteProfileIdentityWithContext(context, deleteProfileIdentityOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProfileIdentityWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProfileIdentityWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
