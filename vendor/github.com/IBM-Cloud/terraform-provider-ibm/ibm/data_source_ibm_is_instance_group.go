// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceGroupRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user-defined name for this instance group",
			},

			"instance_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance template ID",
			},

			"membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of instances in the instance group",
			},

			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group ID",
			},

			"subnets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "list of subnet IDs",
			},

			"application_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Used by the instance group when scaling up instances to supply the port for the load balancer pool member.",
			},

			"load_balancer_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "load balancer pool ID",
			},

			"managers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "list of Managers associated with instancegroup",
			},

			"vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vpc instance",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance group status - deleting, healthy, scaling, unhealthy",
			},
		},
	}
}

func dataSourceIBMISInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	name := d.Get("name")

	// Support for pagination
	start := ""
	allrecs := []vpcv1.InstanceGroup{}
	for {
		listInstanceGroupOptions := vpcv1.ListInstanceGroupsOptions{}
		if start != "" {
			listInstanceGroupOptions.Start = &start
		}
		instanceGroupsCollection, response, err := sess.ListInstanceGroups(&listInstanceGroupOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching InstanceGroups %s\n%s", err, response)
		}
		start = GetNext(instanceGroupsCollection.Next)
		allrecs = append(allrecs, instanceGroupsCollection.InstanceGroups...)

		if start == "" {
			break
		}

	}

	for _, instanceGroup := range allrecs {
		if *instanceGroup.Name == name {
			d.Set("name", *instanceGroup.Name)
			d.Set("instance_template", *instanceGroup.InstanceTemplate.ID)
			d.Set("membership_count", *instanceGroup.MembershipCount)
			d.Set("resource_group", *instanceGroup.ResourceGroup.ID)
			d.SetId(*instanceGroup.ID)
			if instanceGroup.ApplicationPort != nil {
				d.Set("application_port", *instanceGroup.ApplicationPort)
			}
			subnets := make([]string, 0)
			for i := 0; i < len(instanceGroup.Subnets); i++ {
				subnets = append(subnets, string(*(instanceGroup.Subnets[i].ID)))
			}
			if instanceGroup.LoadBalancerPool != nil {
				d.Set("load_balancer_pool", *instanceGroup.LoadBalancerPool.ID)
			}
			d.Set("subnets", subnets)
			managers := make([]string, 0)
			for i := 0; i < len(instanceGroup.Managers); i++ {
				managers = append(managers, string(*(instanceGroup.Managers[i].ID)))
			}
			d.Set("managers", managers)
			d.Set("vpc", *instanceGroup.VPC.ID)
			d.Set("status", *instanceGroup.Status)
			return nil
		}
	}
	return fmt.Errorf("Instance group %s not found", name)
}
