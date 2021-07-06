// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISSecurityGroupTarget() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceIBMISSecurityGroupTargetRead,

		Schema: map[string]*schema.Schema{

			"security_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security group id",
			},

			"target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group target identifier",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func dataSourceIBMISSecurityGroupTargetRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	securityGroupID := d.Get("security_group").(string)
	name := d.Get("name").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.SecurityGroupTargetReferenceIntf{}

	for {
		listSecurityGroupTargetsOptions := sess.NewListSecurityGroupTargetsOptions(securityGroupID)

		groups, response, err := sess.ListSecurityGroupTargets(listSecurityGroupTargetsOptions)
		if err != nil {
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

	for _, securityGroupTargetReferenceIntf := range allrecs {
		securityGroupTargetReference := securityGroupTargetReferenceIntf.(*vpcv1.SecurityGroupTargetReference)
		if *securityGroupTargetReference.Name == name {
			d.Set("target", *securityGroupTargetReference.ID)
			d.Set("resource_type", *securityGroupTargetReference.ResourceType)
			if securityGroupTargetReference.Deleted != nil {
				d.Set("more_info", *securityGroupTargetReference.Deleted.MoreInfo)
			}
			d.SetId(fmt.Sprintf("%s/%s", securityGroupID, *securityGroupTargetReference.ID))
			return nil
		}
	}
	return fmt.Errorf("Security Group Target %s not found", name)
}
