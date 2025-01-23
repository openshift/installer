// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isReservedIPLimit  = "limit"
	isReservedIPSort   = "sort"
	isReservedIPs      = "reserved_ips"
	isReservedIPsCount = "total_count"
)

func DataSourceIBMISReservedIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSdataSourceIBMISReservedIPsRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet identifier.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isReservedIPs: {
				Type:        schema.TypeList,
				Description: "Collection of reserved IPs in this subnet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservedIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						isReservedIPAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If reserved ip shall be deleted automatically",
						},
						isReservedIPCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the reserved IP was created.",
						},
						isReservedIPhref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						isReservedIPID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP",
						},
						isReservedIPLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the reserved IP",
						},
						isReservedIPName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this reserved IP.",
						},
						isReservedIPOwner: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
						},
						isReservedIPType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						isReservedIPTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserved IP target id",
						},
						isReservedIPTargetCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn for target.",
						},
						"target_reference": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The target this reserved IP is bound to.If absent, this reserved IP is provider-owned or unbound.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this endpoint gateway.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this endpoint gateway.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this endpoint gateway.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this endpoint gateway. The name is unique across all endpoint gateways in the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
					},
				},
			},
			isReservedIPsCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages",
			},
		},
	}
}

func dataSdataSourceIBMISReservedIPsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnetID := d.Get(isSubNetID).(string)

	// Flatten all the reserved IPs
	start := ""
	allrecs := []vpcv1.ReservedIP{}
	for {
		options := &vpcv1.ListSubnetReservedIpsOptions{SubnetID: &subnetID}

		if start != "" {
			options.Start = &start
		}

		result, response, err := sess.ListSubnetReservedIps(options)
		if err != nil || response == nil || result == nil {
			return fmt.Errorf("[ERROR] Error fetching reserved ips %s\n%s", err, response)
		}
		start = flex.GetNext(result.Next)
		allrecs = append(allrecs, result.ReservedIps...)
		if start == "" {
			break
		}
	}

	// Now store all the reserved IP info with their response tags
	reservedIPs := []map[string]interface{}{}
	for _, data := range allrecs {
		ipsOutput := map[string]interface{}{}
		ipsOutput[isReservedIPAddress] = *data.Address
		ipsOutput[isReservedIPAutoDelete] = *data.AutoDelete
		ipsOutput[isReservedIPCreatedAt] = (*data.CreatedAt).String()
		ipsOutput[isReservedIPhref] = *data.Href
		ipsOutput[isReservedIPID] = *data.ID
		ipsOutput[isReservedIPLifecycleState] = data.LifecycleState
		ipsOutput[isReservedIPName] = *data.Name
		ipsOutput[isReservedIPOwner] = *data.Owner
		ipsOutput[isReservedIPType] = *data.ResourceType
		target := []map[string]interface{}{}
		if data.Target != nil {
			modelMap, err := dataSourceIBMIsReservedIPReservedIPTargetToMap(data.Target)
			if err != nil {
				return err
			}
			target = append(target, modelMap)
		}
		ipsOutput["target_reference"] = target
		if len(target) > 0 {
			ipsOutput[isReservedIPTarget] = target[0]["id"]
			ipsOutput[isReservedIPTargetCrn] = target[0]["crn"]
		}

		reservedIPs = append(reservedIPs, ipsOutput)
	}

	d.SetId(time.Now().UTC().String()) // This is not any reserved ip or subnet id but state id
	d.Set(isReservedIPs, reservedIPs)
	d.Set(isReservedIPsCount, len(reservedIPs))
	d.Set(isSubNetID, subnetID)
	return nil
}

func dataSourceIBMIsReservedIPReservedIPTargetToMap(model vpcv1.ReservedIPTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ReservedIPTargetEndpointGatewayReference); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetEndpointGatewayReferenceToMap(model.(*vpcv1.ReservedIPTargetEndpointGatewayReference))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetVirtualNetworkInterfaceReferenceReservedIPTargetContext); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetVirtualNetworkInterfaceReferenceReservedIPTargetContextToMap(model.(*vpcv1.ReservedIPTargetVirtualNetworkInterfaceReferenceReservedIPTargetContext))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetNetworkInterfaceReferenceTargetContextToMap(model.(*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetBareMetalServerNetworkInterfaceReferenceTargetContext); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetBareMetalServerNetworkInterfaceReferenceTargetContextToMap(model.(*vpcv1.ReservedIPTargetBareMetalServerNetworkInterfaceReferenceTargetContext))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetLoadBalancerReference); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetLoadBalancerReferenceToMap(model.(*vpcv1.ReservedIPTargetLoadBalancerReference))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetVPNGatewayReference); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetVPNGatewayReferenceToMap(model.(*vpcv1.ReservedIPTargetVPNGatewayReference))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetVPNServerReference); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetVPNServerReferenceToMap(model.(*vpcv1.ReservedIPTargetVPNServerReference))
	} else if _, ok := model.(*vpcv1.ReservedIPTargetGenericResourceReference); ok {
		return dataSourceIBMIsReservedIPReservedIPTargetGenericResourceReferenceToMap(model.(*vpcv1.ReservedIPTargetGenericResourceReference))
	} else if _, ok := model.(*vpcv1.ReservedIPTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ReservedIPTarget)
		if model.CRN != nil {
			modelMap["crn"] = model.CRN
		}
		if model.Deleted != nil {
			deletedMap, err := dataSourceIBMIsReservedIPEndpointGatewayReferenceDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = model.Href
		}
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ReservedIPTargetIntf subtype encountered")
	}
}

func dataSourceIBMIsReservedIPEndpointGatewayReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetEndpointGatewayReferenceToMap(model *vpcv1.ReservedIPTargetEndpointGatewayReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPEndpointGatewayReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetVirtualNetworkInterfaceReferenceReservedIPTargetContextToMap(model *vpcv1.ReservedIPTargetVirtualNetworkInterfaceReferenceReservedIPTargetContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetNetworkInterfaceReferenceTargetContextToMap(model *vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPNetworkInterfaceReferenceTargetContextDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPNetworkInterfaceReferenceTargetContextDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetBareMetalServerNetworkInterfaceReferenceTargetContextToMap(model *vpcv1.ReservedIPTargetBareMetalServerNetworkInterfaceReferenceTargetContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPBareMetalServerNetworkInterfaceReferenceTargetContextDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPBareMetalServerNetworkInterfaceReferenceTargetContextDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetLoadBalancerReferenceToMap(model *vpcv1.ReservedIPTargetLoadBalancerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPLoadBalancerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPLoadBalancerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetVPNGatewayReferenceToMap(model *vpcv1.ReservedIPTargetVPNGatewayReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPVPNGatewayReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPVPNGatewayReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetVPNServerReferenceToMap(model *vpcv1.ReservedIPTargetVPNServerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPVPNServerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPVPNServerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsReservedIPReservedIPTargetGenericResourceReferenceToMap(model *vpcv1.ReservedIPTargetGenericResourceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsReservedIPGenericResourceReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsReservedIPGenericResourceReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
