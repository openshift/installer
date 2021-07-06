// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISInstanceGroupManagerPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceGroupManagerPoliciesRead,

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

			"instance_group_manager_policies": {
				Type:        schema.TypeList,
				Description: "List of instance group manager policies",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instance group manager policy.",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceIBMISInstanceGroupManagerPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupManagerID := d.Get("instance_group_manager").(string)
	instanceGroupID := d.Get("instance_group").(string)

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

	policies := make([]map[string]interface{}, 0)
	for _, data := range allrecs {
		instanceGroupManagerPolicy := data.(*vpcv1.InstanceGroupManagerPolicy)
		policy := map[string]interface{}{
			"id":           fmt.Sprintf("%s/%s/%s", instanceGroupID, instanceGroupManagerID, *instanceGroupManagerPolicy.ID),
			"name":         *instanceGroupManagerPolicy.Name,
			"metric_value": *instanceGroupManagerPolicy.MetricValue,
			"metric_type":  *instanceGroupManagerPolicy.MetricType,
			"policy_type":  *instanceGroupManagerPolicy.PolicyType,
			"policy_id":    *instanceGroupManagerPolicy.ID,
		}
		policies = append(policies, policy)
	}
	d.Set("instance_group_manager_policies", policies)
	d.SetId(dataSourceIBMISInstanceGroupManagerPoliciesID(d))
	return nil
}

// dataSourceIBMISInstanceGroupManagerPoliciesID returns a reasonable ID for a instance group manager policies list.
func dataSourceIBMISInstanceGroupManagerPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
