// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileIdentities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileIdentitiesRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entity tag of the profile identities response.",
			},
			"identities": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of identities.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity.",
						},
						"identifier": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn' it uses account id contained in the CRN.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the identity.",
						},
						"accounts": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Only valid for the type user. Accounts from which a user can assume the trusted profile.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the identity that can assume the trusted profile. This is optional field for all the types of identities. When this field is not set for the identity type 'serviceid' then the description of the service id is used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfileIdentitiesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileIdentitiesOptions := &iamidentityv1.GetProfileIdentitiesOptions{}

	getProfileIdentitiesOptions.SetProfileID(d.Get("profile_id").(string))

	profileIdentitiesResponse, response, err := iamIdentityClient.GetProfileIdentitiesWithContext(context, getProfileIdentitiesOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProfileIdentitiesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfileIdentitiesWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIamTrustedProfileIdentitiesID(d))

	if err = d.Set("entity_tag", profileIdentitiesResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	identities := []map[string]interface{}{}
	if profileIdentitiesResponse.Identities != nil {
		for _, modelItem := range profileIdentitiesResponse.Identities {
			modelMap, err := dataSourceIBMIamTrustedProfileIdentitiesProfileIdentityResponseToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			identities = append(identities, modelMap)
		}
	}
	if err = d.Set("identities", identities); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting identities %s", err))
	}

	return nil
}

// dataSourceIBMIamTrustedProfileIdentitiesID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfileIdentitiesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIamTrustedProfileIdentitiesProfileIdentityResponseToMap(model *iamidentityv1.ProfileIdentityResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["identifier"] = model.Identifier
	modelMap["type"] = model.Type
	if model.Accounts != nil {
		modelMap["accounts"] = model.Accounts
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}
