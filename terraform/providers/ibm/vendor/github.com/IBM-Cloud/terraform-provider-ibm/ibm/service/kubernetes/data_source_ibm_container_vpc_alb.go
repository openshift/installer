// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerVPCClusterALB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerVpcALBRead,
		Schema: map[string]*schema.Schema{
			"alb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ALB ID",
			},
			"alb_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_deployment": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resize": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
			},
		},
	}
}

func dataSourceIBMContainerVpcALBRead(d *schema.ResourceData, meta interface{}) error {
	albClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	albID := d.Get("alb_id").(string)
	albAPI := albClient.Albs()
	targetEnv, _ := getVpcClusterTargetHeader(d)
	albConfig, err := albAPI.GetAlb(albID, targetEnv)
	if err != nil {
		return err
	}

	d.Set("alb_type", albConfig.AlbType)
	d.Set("cluster", albConfig.Cluster)
	d.Set("name", albConfig.Name)
	d.Set("enable", albConfig.Enable)
	d.Set("disable_deployment", albConfig.DisableDeployment)
	d.Set("resize", albConfig.Resize)
	d.Set("zone", albConfig.ZoneAlb)
	d.Set("status", albConfig.Status)
	d.Set("state", albConfig.State)
	d.Set("load_balancer_hostname", albConfig.LoadBalancerHostname)
	d.SetId(albID)
	return nil
}
