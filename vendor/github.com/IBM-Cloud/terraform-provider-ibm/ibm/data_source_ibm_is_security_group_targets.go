// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISSecurityGroupTargets() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceIBMISSecurityGroupTargetsRead,

		Schema: map[string]*schema.Schema{

			"security_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security group id",
			},

			"targets": {
				Type:        schema.TypeList,
				Description: "List of targets",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group target identifier",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group target name",
						},

						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource Type",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about deleted resources",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISSecurityGroupTargetsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	securityGroupID := d.Get("security_group").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.SecurityGroupTargetReferenceIntf{}

	for {
		listSecurityGroupTargetsOptions := sess.NewListSecurityGroupTargetsOptions(securityGroupID)

		groups, response, err := sess.ListSecurityGroupTargets(listSecurityGroupTargetsOptions)
		if err != nil || groups == nil {
			return fmt.Errorf("Error Getting InstanceGroup Managers %s\n%s", err, response)
		}
		if *groups.TotalCount == int64(0) {
			break
		}

		start = GetNext(groups.Next)
		allrecs = append(allrecs, groups.Targets...)

		if start == "" {
			break
		}

	}

	targets := make([]map[string]interface{}, 0)
	for _, securityGroupTargetReferenceIntf := range allrecs {
		securityGroupTargetReference := securityGroupTargetReferenceIntf.(*vpcv1.SecurityGroupTargetReference)
		tr := map[string]interface{}{
			"name":          *securityGroupTargetReference.Name,
			"target":        *securityGroupTargetReference.ID,
			"resource_type": *securityGroupTargetReference.ResourceType,
		}
		if securityGroupTargetReference.Deleted != nil {
			tr["more_info"] = *securityGroupTargetReference.Deleted.MoreInfo
		}
		targets = append(targets, tr)
	}
	d.Set("targets", targets)
	d.SetId(securityGroupID)
	return nil
}
