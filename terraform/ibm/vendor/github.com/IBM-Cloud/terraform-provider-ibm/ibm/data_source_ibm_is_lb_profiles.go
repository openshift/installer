// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLbsProfiles = "lb_profiles"
)

func dataSourceIBMISLbProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISLbProfilesRead,

		Schema: map[string]*schema.Schema{

			isLbsProfiles: {
				Type:        schema.TypeList,
				Description: "Collection of load balancer profile collectors",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this load balancer profile",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this load balancer profile",
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this load balancer profile belongs to",
						},
						"route_mode_supported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The route mode support for a load balancer with this profile depends on its configuration",
						},
						"route_mode_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route mode type for this load balancer profile, one of [fixed, dependent]",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISLbProfilesRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	start := ""
	allrecs := []vpcv1.LoadBalancerProfile{}
	for {
		listOptions := &vpcv1.ListLoadBalancerProfilesOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		profileCollectors, response, err := sess.ListLoadBalancerProfiles(listOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Load Balancer Profiles for VPC %s\n%s", err, response)
		}
		start = GetNext(profileCollectors.Next)
		allrecs = append(allrecs, profileCollectors.Profiles...)
		if start == "" {
			break
		}
	}
	lbprofilesInfo := make([]map[string]interface{}, 0)
	for _, profileCollector := range allrecs {

		l := map[string]interface{}{
			"name":   *profileCollector.Name,
			"href":   *profileCollector.Href,
			"family": *profileCollector.Family,
		}
		if profileCollector.RouteModeSupported != nil {
			routeMode := profileCollector.RouteModeSupported
			switch reflect.TypeOf(routeMode).String() {
			case "*vpcv1.LoadBalancerProfileRouteModeSupportedFixed":
				{
					rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupportedFixed)
					l["route_mode_supported"] = rms.Value
					l["route_mode_type"] = rms.Type
				}
			case "*vpcv1.LoadBalancerProfileRouteModeSupportedDependent":
				{
					rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupportedDependent)
					if rms.Type != nil {
						l["route_mode_type"] = *rms.Type
					}
				}
			case "*vpcv1.LoadBalancerProfileRouteModeSupported":
				{
					rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupported)
					if rms.Type != nil {
						l["route_mode_type"] = *rms.Type
					}
					if rms.Value != nil {
						l["route_mode_supported"] = *rms.Value
					}
				}
			}
		}
		lbprofilesInfo = append(lbprofilesInfo, l)
	}
	d.SetId(dataSourceIBMISLbProfilesID(d))
	d.Set(isLbsProfiles, lbprofilesInfo)
	return nil
}

// dataSourceIBMISLbProfilesID returns a reasonable ID for a profileCollector list.
func dataSourceIBMISLbProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
