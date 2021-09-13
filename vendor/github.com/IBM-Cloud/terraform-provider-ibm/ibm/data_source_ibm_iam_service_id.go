// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
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
							Deprecated:  "bound_to attribute in service_ids list has been deprecated",
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

	name := d.Get("name").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	iamClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	start := ""
	allrecs := []iamidentityv1.ServiceID{}
	var pg int64 = 100
	for {
		listServiceIDOptions := iamidentityv1.ListServiceIdsOptions{
			AccountID: &userDetails.userAccount,
			Pagesize:  &pg,
			Name:      &name,
		}
		if start != "" {
			listServiceIDOptions.Pagetoken = &start
		}

		serviceIDs, resp, err := iamClient.ListServiceIds(&listServiceIDOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing Service Ids %s %s", err, resp)
		}
		start = GetNextIAM(serviceIDs.Next)
		allrecs = append(allrecs, serviceIDs.Serviceids...)
		if start == "" {
			break
		}
	}
	if len(allrecs) == 0 {
		return fmt.Errorf("[ERROR] No serviceID found with name [%s]", name)

	}

	serviceIDListMap := make([]map[string]interface{}, 0, len(allrecs))
	for _, serviceID := range allrecs {
		l := map[string]interface{}{
			"id": serviceID.ID,
			// "bound_to":    serviceID.BoundTo,
			"version":     serviceID.EntityTag,
			"description": serviceID.Description,
			"crn":         serviceID.CRN,
			"locked":      serviceID.Locked,
			"iam_id":      serviceID.IamID,
		}
		serviceIDListMap = append(serviceIDListMap, l)
	}
	d.SetId(name)
	d.Set("service_ids", serviceIDListMap)
	return nil
}
