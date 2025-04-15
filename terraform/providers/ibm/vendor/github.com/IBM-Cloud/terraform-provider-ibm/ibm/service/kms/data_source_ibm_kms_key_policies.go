// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"log"

	//kp "github.com/IBM/keyprotect-go-client"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMKMSkeyPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMKMSKeyPoliciesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key protect or hpcs instance GUID",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
				Default:      "public",
			},
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Key ID of the Key",
				ExactlyOneOf: []string{"key_id", "alias"},
			},
			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Alias of the Key",
				ExactlyOneOf: []string{"key_id", "alias"},
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Creates or updates one or more policies for the specified key",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rotation": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the key rotation with enabled and time interval in months, with a minimum of 1, and a maximum of 12",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The v4 UUID used to uniquely identify the policy resource, as specified by RFC 4122.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that created the policy.",
									},
									"creation_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date the policy was created. The date format follows RFC 3339.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that updated the policy.",
									},
									"last_update_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
									},
									"interval_month": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the key rotation time interval in months",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether the key rotation policy is enabled.",
									},
								},
							},
						},
						"dual_auth_delete": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data associated with the dual authorization delete policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The v4 UUID used to uniquely identify the policy resource, as specified by RFC 4122.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that created the policy.",
									},
									"creation_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date the policy was created. The date format follows RFC 3339.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that updated the policy.",
									},
									"last_update_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies the dual authorization policy on a single key.",
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

func dataSourceIBMKMSKeyPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	api, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	endpointType := d.Get("endpoint_type").(string)
	var id string
	if v, ok := d.GetOk("key_id"); ok {
		id = v.(string)
		d.Set("key_id", id)
	}
	if v, ok := d.GetOk("alias"); ok {
		id = v.(string)
		key, err := api.GetKey(context, id)
		if err != nil {
			return diag.Errorf("Failed to get Key: %s", err)
		}
		d.Set("alias", id)
		d.Set("key_id", key.ID)
	}
	policies, err := api.GetPolicies(context, id)
	if err != nil {
		return diag.Errorf("Failed to read policies: %s", err)
	}

	if len(policies) == 0 {
		log.Printf("No Policy Configurations read\n")
	} else {
		d.Set("policies", flex.FlattenKeyPolicies(policies))
	}
	d.SetId(instanceID)
	d.Set("instance_id", instanceID)
	d.Set("endpoint_type", endpointType)

	return nil
}
