// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISInstanceGroupManagerPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceGroupManagerPolicyRead,

		Schema: map[string]*schema.Schema{

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"instance_group_manager": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance group manager ID",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the instance group manager policy",
			},

			"metric_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of metric to be evaluated",
			},

			"metric_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The metric value to be evaluated",
			},

			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Policy for the Instance Group",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy ID",
			},
		},
	}
}

func dataSourceIBMISInstanceGroupManagerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupManagerID := d.Get("instance_group_manager").(string)
	instanceGroupID := d.Get("instance_group").(string)
	policyName := d.Get("name").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.InstanceGroupManagerPolicyIntf{}

	for {
		listInstanceGroupManagerPoliciesOptions := vpcv1.ListInstanceGroupManagerPoliciesOptions{
			InstanceGroupID:        &instanceGroupID,
			InstanceGroupManagerID: &instanceGroupManagerID,
		}

		instanceGroupManagerPolicyCollection, response, err := sess.ListInstanceGroupManagerPolicies(&listInstanceGroupManagerPoliciesOptions)
		if err != nil {
			return fmt.Errorf("Error Getting InstanceGroup Manager Policies %s\n%s", err, response)
		}
		start = GetNext(instanceGroupManagerPolicyCollection.Next)
		allrecs = append(allrecs, instanceGroupManagerPolicyCollection.Policies...)
		if start == "" {
			break
		}
	}

	for _, data := range allrecs {
		instanceGroupManagerPolicy := data.(*vpcv1.InstanceGroupManagerPolicy)
		if policyName == *instanceGroupManagerPolicy.Name {
			d.SetId(fmt.Sprintf("%s/%s/%s", instanceGroupID, instanceGroupManagerID, *instanceGroupManagerPolicy.ID))
			d.Set("policy_id", *instanceGroupManagerPolicy.ID)
			d.Set("metric_value", *instanceGroupManagerPolicy.MetricValue)
			d.Set("metric_type", *instanceGroupManagerPolicy.MetricType)
			d.Set("policy_type", *instanceGroupManagerPolicy.PolicyType)
			return nil
		}
	}
	return fmt.Errorf("Instance group manager policy %s not found", policyName)
}
