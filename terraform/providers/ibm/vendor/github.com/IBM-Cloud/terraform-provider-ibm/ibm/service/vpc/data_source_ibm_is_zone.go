// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isZoneName   = "name"
	isZoneRegion = "region"
	isZoneStatus = "status"

	isZoneDataCenter    = "data_center"
	isZoneUniversalName = "universal_name"
)

func DataSourceIBMISZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISZoneRead,

		Schema: map[string]*schema.Schema{

			isZoneName: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneRegion: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			isZoneDataCenter: {
				Type:     schema.TypeString,
				Computed: true,
			},
			isZoneUniversalName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMISZoneRead(d *schema.ResourceData, meta interface{}) error {
	regionName := d.Get(isZoneRegion).(string)
	zoneName := d.Get(isZoneName).(string)
	return zoneGet(d, meta, regionName, zoneName)
}

func zoneGet(d *schema.ResourceData, meta interface{}, regionName, zoneName string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getRegionZoneOptions := &vpcv1.GetRegionZoneOptions{
		RegionName: &regionName,
		Name:       &zoneName,
	}
	zone, _, err := sess.GetRegionZone(getRegionZoneOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from region name + zone name.
	id := fmt.Sprintf("%s.%s", *zone.Region.Name, *zone.Name)
	d.SetId(id)
	d.Set(isZoneName, *zone.Name)
	d.Set(isZoneRegion, *zone.Region.Name)
	d.Set(isZoneStatus, *zone.Status)
	if zone.DataCenter != nil {
		d.Set(isZoneDataCenter, *zone.DataCenter)
	}
	if zone.UniversalName != nil {
		d.Set(isZoneUniversalName, *zone.UniversalName)
	}
	return nil
}
