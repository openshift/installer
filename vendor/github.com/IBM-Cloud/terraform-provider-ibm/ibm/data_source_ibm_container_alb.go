// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerALB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerALBRead,

		Schema: map[string]*schema.Schema{
			"alb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ALB ID",
			},
			"alb_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB type",
			},
			"cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster id",
			},
			"user_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP assigned by the user",
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "set to true if ALB needs to be enabled",
			},
			"disable_deployment": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set to true if ALB needs to be disabled",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB name",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB zone",
			},
		},
	}
}

func dataSourceIBMContainerALBRead(d *schema.ResourceData, meta interface{}) error {
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	albID := d.Get("alb_id").(string)

	albAPI := albClient.Albs()
	targetEnv, err := getAlbTargetHeader(d, meta)
	if err != nil {
		return err
	}
	albConfig, err := albAPI.GetALB(albID, targetEnv)
	if err != nil {
		return err
	}

	d.SetId(albID)
	d.Set("alb_type", &albConfig.ALBType)
	d.Set("cluster", &albConfig.ClusterID)
	d.Set("name", &albConfig.Name)
	d.Set("enable", &albConfig.Enable)
	d.Set("disable_deployment", &albConfig.DisableDeployment)
	d.Set("replicas", &albConfig.NumOfInstances)
	d.Set("resize", &albConfig.Resize)
	d.Set("user_ip", &albConfig.ALBIP)
	d.Set("zone", &albConfig.Zone)
	return nil
}
