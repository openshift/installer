// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISSecurityGroupTarget() *schema.Resource {
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

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this security group target",
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
		if start != "" {
			listSecurityGroupTargetsOptions.Start = &start
		}
		groups, response, err := sess.ListSecurityGroupTargets(listSecurityGroupTargetsOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Getting InstanceGroup Managers %s\n%s", err, response)
		}
		if *groups.TotalCount == int64(0) {
			break
		}

		start = flex.GetNext(groups.Next)
		allrecs = append(allrecs, groups.Targets...)

		if start == "" {
			break
		}

	}

	for _, securityGroupTargetReferenceIntf := range allrecs {
		securityGroupTargetReference := securityGroupTargetReferenceIntf.(*vpcv1.SecurityGroupTargetReference)
		if *securityGroupTargetReference.Name == name {
			d.Set("target", *securityGroupTargetReference.ID)
			d.Set("crn", securityGroupTargetReference.CRN)
			// d.Set("resource_type", *securityGroupTargetReference.ResourceType)
			if securityGroupTargetReference.Deleted != nil {
				d.Set("more_info", *securityGroupTargetReference.Deleted.MoreInfo)
			}
			if securityGroupTargetReference != nil && securityGroupTargetReference.ResourceType != nil {
				d.Set("resource_type", *securityGroupTargetReference.ResourceType)
			}
			d.SetId(fmt.Sprintf("%s/%s", securityGroupID, *securityGroupTargetReference.ID))
			return nil
		}
	}
	return fmt.Errorf("Security Group Target %s not found", name)
}
