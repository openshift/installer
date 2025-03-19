// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLbsProfiles    = "lb_profiles"
	isLbsProfileName = "name"
)

func DataSourceIBMISLbProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISLbProfilesRead,

		Schema: map[string]*schema.Schema{
			isLbsProfileName: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer profile name.",
			},
			isLbsProfiles: {
				Type:        schema.TypeList,
				Description: "Collection of load balancer profile collectors",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBAccessModes: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The access mode for a load balancer with this profile",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for access mode",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Access modes for this profile",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Access modes for this profile",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"availability": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The availability mode for a load balancer with this profile",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of availability, one of [fixed, dependent]",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The availability of this load balancer, one of [subnet, region]. Applicable only if type is fixed",
									},
								},
							},
						},
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
						"instance_groups_supported": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance groups support for the load balancer with this profile",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of support for instance groups, one of [fixed, dependent]",
									},
									"value": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether Instance groups are supported for this profile. Applicable only if type is fixed",
									},
								},
							},
						},
						"source_ip_session_persistence_supported": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source IP session ip persistence support for a load balancer with this profile",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of support for session ip persistence, one of [fixed, dependent on configuration]",
									},
									"value": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether session ip persistence are supported for this profile. Applicable only if type is fixed",
									},
								},
							},
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
	if lbprofilenameok, ok := d.GetOk(isLbsProfileName); ok {
		lbprofilename := lbprofilenameok.(string)
		getLoadBalancerProfileOptions := &vpcv1.GetLoadBalancerProfileOptions{
			Name: &lbprofilename,
		}
		lbProfile, response, err := sess.GetLoadBalancerProfile(getLoadBalancerProfileOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Load Balancer Profile(%s) for VPC %s\n%s", lbprofilename, err, response)
		}
		allrecs = append(allrecs, *lbProfile)
	} else {
		for {
			listOptions := &vpcv1.ListLoadBalancerProfilesOptions{}
			if start != "" {
				listOptions.Start = &start
			}
			profileCollectors, response, err := sess.ListLoadBalancerProfiles(listOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Fetching Load Balancer Profiles for VPC %s\n%s", err, response)
			}
			start = flex.GetNext(profileCollectors.Next)
			allrecs = append(allrecs, profileCollectors.Profiles...)
			if start == "" {
				break
			}
		}
	}
	lbprofilesInfo := make([]map[string]interface{}, 0)
	for _, profileCollector := range allrecs {

		l := map[string]interface{}{
			"name":   *profileCollector.Name,
			"href":   *profileCollector.Href,
			"family": *profileCollector.Family,
		}
		if profileCollector.UDPSupported != nil {
			udpSupport := profileCollector.UDPSupported
			switch reflect.TypeOf(udpSupport).String() {
			case "*vpcv1.LoadBalancerProfileUDPSupportedFixed":
				{
					udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupportedFixed)
					l["udp_supported"] = udp.Value
					l["udp_supported_type"] = udp.Type
				}
			case "*vpcv1.LoadBalancerProfileUDPSupportedDependent":
				{
					udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupportedDependent)
					if udp.Type != nil {
						l["udp_supported_type"] = *udp.Type
					}
				}
			case "*vpcv1.LoadBalancerProfileUDPSupported":
				{
					udp := udpSupport.(*vpcv1.LoadBalancerProfileUDPSupported)
					if udp.Type != nil {
						l["udp_supported_type"] = *udp.Type
					}
					if udp.Value != nil {
						l["udp_supported"] = *udp.Value
					}
				}
			}
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

		if profileCollector.AccessModes != nil {
			accessModes := profileCollector.AccessModes
			AccessModesMap := map[string]interface{}{}
			AccessModesList := []map[string]interface{}{}
			if accessModes.Type != nil {
				AccessModesMap["type"] = *accessModes.Type
			}
			if len(accessModes.Values) > 0 {
				AccessModesMap["values"] = accessModes.Values
			}
			AccessModesList = append(AccessModesList, AccessModesMap)
			l[isLBAccessModes] = AccessModesList
		}
		if profileCollector.Availability != nil {
			availabilitySupport := profileCollector.Availability.(*vpcv1.LoadBalancerProfileAvailability)
			availabilitySupportMap := map[string]interface{}{}
			availabilitySupportList := []map[string]interface{}{}
			if availabilitySupport.Type != nil {
				availabilitySupportMap["type"] = *availabilitySupport.Type
			}
			if availabilitySupport.Value != nil {
				availabilitySupportMap["value"] = *availabilitySupport.Value
			}
			availabilitySupportList = append(availabilitySupportList, availabilitySupportMap)
			l["availability"] = availabilitySupportList
		}
		if profileCollector.InstanceGroupsSupported != nil {
			instanceGroupSupport := profileCollector.InstanceGroupsSupported.(*vpcv1.LoadBalancerProfileInstanceGroupsSupported)
			instanceGroupSupportMap := map[string]interface{}{}
			instanceGroupSupportList := []map[string]interface{}{}
			if instanceGroupSupport.Type != nil {
				instanceGroupSupportMap["type"] = *instanceGroupSupport.Type
			}
			if instanceGroupSupport.Value != nil {
				instanceGroupSupportMap["value"] = *instanceGroupSupport.Value
			}
			instanceGroupSupportList = append(instanceGroupSupportList, instanceGroupSupportMap)
			l["source_ip_session_persistence_supported"] = instanceGroupSupportList
		}
		if profileCollector.SourceIPSessionPersistenceSupported != nil {
			sourceIpPersistenceSupport := profileCollector.SourceIPSessionPersistenceSupported.(*vpcv1.LoadBalancerProfileSourceIPSessionPersistenceSupported)
			sourceIpPersistenceSupportMap := map[string]interface{}{}
			sourceIpPersistenceSupportList := []map[string]interface{}{}
			if sourceIpPersistenceSupport.Type != nil {
				sourceIpPersistenceSupportMap["type"] = *sourceIpPersistenceSupport.Type
			}
			if sourceIpPersistenceSupport.Value != nil {
				sourceIpPersistenceSupportMap["value"] = *sourceIpPersistenceSupport.Value
			}
			sourceIpPersistenceSupportList = append(sourceIpPersistenceSupportList, sourceIpPersistenceSupportMap)
			l["instance_groups_supported"] = sourceIpPersistenceSupportList
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
