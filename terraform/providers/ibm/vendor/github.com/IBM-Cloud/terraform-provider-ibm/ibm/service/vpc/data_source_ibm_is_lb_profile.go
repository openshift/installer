// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISLbProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISLbProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
			"udp_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The UDP support for a load balancer with this profile",
			},
			"udp_supported_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UDP support type for a load balancer with this profile",
			},
		},
	}
}

func dataSourceIBMISLbProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbprofilename := d.Get(isLbsProfileName).(string)
	getLoadBalancerProfileOptions := &vpcv1.GetLoadBalancerProfileOptions{
		Name: &lbprofilename,
	}
	lbProfile, response, err := sess.GetLoadBalancerProfileWithContext(context, getLoadBalancerProfileOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Fetching Load Balancer Profile(%s) for VPC %s\n%s", lbprofilename, err, response))
	}

	d.Set("name", *lbProfile.Name)
	d.Set("href", *lbProfile.Href)
	d.Set("family", *lbProfile.Family)
	log.Printf("[INFO] lbprofile udp %v", lbProfile.UDPSupported)
	if lbProfile.UDPSupported != nil {
		udpSupport := lbProfile.UDPSupported
		log.Printf("[INFO] lbprofile udp %s", reflect.TypeOf(udpSupport).String())

		switch reflect.TypeOf(udpSupport).String() {
		case "*vpcv1.LoadBalancerProfileUDPSupportedFixed":
			{
				udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupportedFixed)
				d.Set("udp_supported", udp.Value)
				d.Set("udp_supported_type", udp.Type)
			}
		case "*vpcv1.LoadBalancerProfileUDPSupportedDependent":
			{
				udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupportedDependent)
				if udp.Type != nil {
					d.Set("udp_supported_type", *udp.Type)
				}
			}
		case "*vpcv1.LoadBalancerProfileUDPSupported":
			{
				udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupported)
				if udp.Type != nil {
					d.Set("udp_supported_type", *udp.Type)
				}
				if udp.Value != nil {
					d.Set("udp_supported", *udp.Value)
				}
			}
		}
	}
	if lbProfile.RouteModeSupported != nil {
		routeMode := lbProfile.RouteModeSupported
		switch reflect.TypeOf(routeMode).String() {
		case "*vpcv1.LoadBalancerProfileRouteModeSupportedFixed":
			{
				rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupportedFixed)
				d.Set("route_mode_supported", rms.Value)
				d.Set("route_mode_type", rms.Type)
			}
		case "*vpcv1.LoadBalancerProfileRouteModeSupportedDependent":
			{
				rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupportedDependent)
				if rms.Type != nil {
					d.Set("route_mode_type", *rms.Type)
				}
			}
		case "*vpcv1.LoadBalancerProfileRouteModeSupported":
			{
				rms := routeMode.(*vpcv1.LoadBalancerProfileRouteModeSupported)
				if rms.Type != nil {
					d.Set("route_mode_type", *rms.Type)
				}
				if rms.Value != nil {
					d.Set("route_mode_supported", *rms.Value)
				}
			}
		}
	}
	d.SetId(*lbProfile.Name)
	return nil
}
