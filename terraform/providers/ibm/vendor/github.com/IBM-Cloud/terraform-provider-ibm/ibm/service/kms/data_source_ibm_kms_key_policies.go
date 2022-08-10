// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"log"
	"strings"

	//kp "github.com/IBM/keyprotect-go-client"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
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
				Optional:    true,
				Computed:    true,
				Description: "Creates or updates one or more policies for the specified key",
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rotation": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: []string{"policies.0.rotation", "policies.0.dual_auth_delete"},
							Description:  "Specifies the key rotation time interval in months, with a minimum of 1, and a maximum of 12",
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
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validate.ValidateAllowedRangeInt(1, 12),
										Description:  "Specifies the key rotation time interval in months",
									},
								},
							},
						},
						"dual_auth_delete": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: []string{"policies.0.rotation", "policies.0.dual_auth_delete"},
							Description:  "Data associated with the dual authorization delete policy.",
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
										Required:    true,
										Description: "If set to true, Key Protect enables a dual authorization policy on a single key.",
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
	api, err := meta.(conns.ClientSession).KeyManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	CrnInstanceID := strings.Split(instanceID, ":")
	if len(CrnInstanceID) > 3 {
		instanceID = CrnInstanceID[len(CrnInstanceID)-3]
	}
	endpointType := d.Get("endpoint_type").(string)

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return diag.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(api, endpointType, extensions)
	if err != nil {
		return diag.FromErr(err)
	}
	api.URL = URL
	api.Config.InstanceID = instanceID
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
