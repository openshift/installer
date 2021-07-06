// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isZoneName   = "name"
	isZoneRegion = "region"
	isZoneStatus = "status"
)

func dataSourceIBMISZone() *schema.Resource {
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
		},
	}
}

func dataSourceIBMISZoneRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	regionName := d.Get(isZoneRegion).(string)
	zoneName := d.Get(isZoneName).(string)
	if userDetails.generation == 1 {
		err := classicZoneGet(d, meta, regionName, zoneName)
		if err != nil {
			return err
		}
	} else {
		err := zoneGet(d, meta, regionName, zoneName)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicZoneGet(d *schema.ResourceData, meta interface{}, regionName, zoneName string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getRegionZoneOptions := &vpcclassicv1.GetRegionZoneOptions{
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
	return nil
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
	return nil
}
