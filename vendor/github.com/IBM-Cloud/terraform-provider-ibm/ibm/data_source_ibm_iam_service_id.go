// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMIAMServiceID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMServiceIDRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the serviceID",
				Type:        schema.TypeString,
				Required:    true,
			},

			"service_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"bound_to": {
							Description: "bound to of the serviceID",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"crn": {
							Description: "CRN of the serviceID",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"description": {
							Description: "description of the serviceID",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"version": {
							Description: "Version of the serviceID",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"locked": {
							Description: "lock state of the serviceID",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"iam_id": {
							Description: "The IAM ID of the serviceID",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMServiceIDRead(d *schema.ResourceData, meta interface{}) error {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	boundTo := crn.New(userDetails.cloudName, userDetails.cloudType)
	boundTo.ScopeType = crn.ScopeAccount
	boundTo.Scope = userDetails.userAccount

	serviceIDS, err := iamClient.ServiceIds().FindByName(boundTo.String(), name)
	if err != nil {
		return err
	}

	if len(serviceIDS) == 0 {
		return fmt.Errorf("No serviceID found with name [%s]", name)

	}

	serviceIDListMap := make([]map[string]interface{}, 0, len(serviceIDS))
	for _, serviceID := range serviceIDS {
		l := map[string]interface{}{
			"id":          serviceID.UUID,
			"bound_to":    serviceID.BoundTo,
			"version":     serviceID.Version,
			"description": serviceID.Description,
			"crn":         serviceID.CRN,
			"locked":      serviceID.Locked,
			"iam_id":      serviceID.IAMID,
		}
		serviceIDListMap = append(serviceIDListMap, l)
	}
	d.SetId(name)
	d.Set("service_ids", serviceIDListMap)
	return nil
}
