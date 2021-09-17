// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isRegionEndpoint = "endpoint"
	isRegionName     = "name"
	isRegionStatus   = "status"
)

func dataSourceIBMISRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISRegionRead,

		Schema: map[string]*schema.Schema{

			isRegionEndpoint: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isRegionName: {
				Type:     schema.TypeString,
				Required: true,
			},

			isRegionStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMISRegionRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	return regionGet(d, meta, name)
}

func regionGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getRegionOptions := &vpcv1.GetRegionOptions{
		Name: &name,
	}
	region, _, err := sess.GetRegion(getRegionOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from region name.
	d.SetId(*region.Name)
	d.Set(isRegionEndpoint, *region.Endpoint)
	d.Set(isRegionName, *region.Name)
	d.Set(isRegionStatus, *region.Status)
	return nil
}
