// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		ReadContext: dataSourceIBMISZoneRead,

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
func dataSourceIBMISZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	regionName := d.Get(isZoneRegion).(string)
	zoneName := d.Get(isZoneName).(string)
	return zoneGet(ctx, d, meta, regionName, zoneName)
}

func zoneGet(ctx context.Context, d *schema.ResourceData, meta interface{}, regionName, zoneName string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_zone", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getRegionZoneOptions := &vpcv1.GetRegionZoneOptions{
		RegionName: &regionName,
		Name:       &zoneName,
	}
	zone, _, err := sess.GetRegionZoneWithContext(ctx, getRegionZoneOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRegionZoneWithContext failed: %s", err.Error()), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	// For lack of anything better, compose our id from region name + zone name.
	id := fmt.Sprintf("%s.%s", *zone.Region.Name, *zone.Name)
	d.SetId(id)
	if !core.IsNil(zone.Name) {
		if err = d.Set(isZoneName, zone.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_zone", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(zone.Region) {
		if err = d.Set(isZoneRegion, *zone.Region.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_is_zone", "read", "set-region").GetDiag()
		}
	}
	if err = d.Set("status", zone.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_zone", "read", "set-status").GetDiag()
	}

	if !core.IsNil(zone.DataCenter) {
		if err = d.Set("data_center", zone.DataCenter); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting data_center: %s", err), "(Data) ibm_is_zone", "read", "set-data_center").GetDiag()
		}
	}

	if !core.IsNil(zone.UniversalName) {
		if err = d.Set("universal_name", zone.UniversalName); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting universal_name: %s", err), "(Data) ibm_is_zone", "read", "set-universal_name").GetDiag()
		}
	}
	return nil
}
