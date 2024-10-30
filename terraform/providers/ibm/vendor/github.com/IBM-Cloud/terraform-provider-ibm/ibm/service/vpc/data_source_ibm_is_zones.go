// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isZoneNames = "zones"
	isZonesInfo = "zone_info"
)

func DataSourceIBMISZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISZonesRead,

		Schema: map[string]*schema.Schema{

			isZoneRegion: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneStatus: {
				Type:     schema.TypeString,
				Optional: true,
			},

			isZoneNames: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			isZonesInfo: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zones information in the region",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isZoneName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isZoneUniversalName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isZoneDataCenter: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isZoneStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISZonesRead(d *schema.ResourceData, meta interface{}) error {

	regionName := d.Get(isZoneRegion).(string)
	return zonesList(d, meta, regionName)
}

func zonesList(d *schema.ResourceData, meta interface{}, regionName string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listRegionZonesOptions := &vpcv1.ListRegionZonesOptions{
		RegionName: &regionName,
	}
	availableZones, _, err := sess.ListRegionZones(listRegionZonesOptions)
	if err != nil {
		return err
	}
	names := make([]string, 0)
	status := d.Get(isZoneStatus).(string)
	zonesList := make([]map[string]interface{}, 0)
	for _, zone := range availableZones.Zones {
		zoneInfo := map[string]interface{}{}
		if status == "" || *zone.Status == status {
			names = append(names, *zone.Name)
			zoneInfo[isZoneName] = *zone.Name
			zoneInfo[isZoneStatus] = *zone.Status
			if zone.DataCenter != nil {
				zoneInfo[isZoneDataCenter] = *zone.DataCenter
			}
			if zone.UniversalName != nil {
				zoneInfo[isZoneUniversalName] = *zone.UniversalName
			}
		}
		zonesList = append(zonesList, zoneInfo)
	}
	d.SetId(dataSourceIBMISZonesId(d))
	d.Set(isZoneNames, names)
	d.Set(isZonesInfo, zonesList)
	return nil
}

// dataSourceIBMISZonesId returns a reasonable ID for a zone list.
func dataSourceIBMISZonesId(d *schema.ResourceData) string {
	// Our zone list is not guaranteed to be stable because the content
	// of the list can vary between two calls if any of the following
	// events occur between calls:
	// - a zone is added to our region
	// - a zone is dropped from our region
	// - we are using the status filter and the status of one or more
	//   zones changes between calls.
	//
	// For simplicity we are using a timestamp for the required terraform id.
	// If we find through usage that this choice is too ephemeral for our users
	// then we can change this function to use a more stable id, perhaps
	// composed from a hash of the list contents. But, for now, a timestamp
	// is good enough.
	return time.Now().UTC().String()
}
