// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isvpnGateways            = "vpn_gateways"
	isVPNGatewayResourceType = "resource_type"
	isVPNGatewayCrn          = "crn"
)

func DataSourceIBMISVPNGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewaysRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this vpn gateway belongs to",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The mode of this vpn gateway.",
			},
			isvpnGateways: {
				Type:        schema.TypeList,
				Description: "Collection of VPN Gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPNGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN Gateway instance name",
						},
						isVPNGatewayCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this VPN gateway was created",
						},
						isVPNGatewayCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN",
						},
						isVPNGatewayMembers: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of VPN gateway members",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The public IP address assigned to the VPN gateway member",
									},

									"private_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private IP address assigned to the VPN gateway member",
									},
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The high availability role assigned to the VPN gateway member",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the VPN gateway member",
									},
								},
							},
						},

						isVPNGatewayResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},

						isVPNGatewayStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN gateway",
						},
						isVPNGatewayHealthState: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
						},
						isVPNGatewayHealthReasons: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						isVPNGatewayLifecycleState: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the VPN route.",
						},
						isVPNGatewayLifecycleReasons: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current lifecycle_state (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						isVPNGatewaySubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPNGateway subnet info",
						},
						isVPNGatewayResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "resource group identifiers ",
						},
						isVPNGatewayMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " VPN gateway mode(policy/route) ",
						},
						"vpc": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VPC for the VPN Gateway",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this VPC.",
									},
								},
							},
						},
						isVPNGatewayTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "VPN Gateway tags list",
						},
						isVPNGatewayAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMVPNGatewaysRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listvpnGWOptions := sess.NewListVPNGatewaysOptions()
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listvpnGWOptions.ResourceGroupID = &resGroup
	}
	if modeIntf, ok := d.GetOk("mode"); ok {
		mode := modeIntf.(string)
		listvpnGWOptions.Mode = &mode
	}
	start := ""
	allrecs := []vpcv1.VPNGatewayIntf{}
	for {
		if start != "" {
			listvpnGWOptions.Start = &start
		}
		availableVPNGateways, detail, err := sess.ListVPNGateways(listvpnGWOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error reading list of VPN Gateways:%s\n%s", err, detail)
		}
		start = flex.GetNext(availableVPNGateways.Next)
		allrecs = append(allrecs, availableVPNGateways.VPNGateways...)
		if start == "" {
			break
		}
	}

	vpngateways := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		gateway := map[string]interface{}{}
		data := instance.(*vpcv1.VPNGateway)
		gateway[isVPNGatewayName] = *data.Name
		gateway[isVPNGatewayCreatedAt] = data.CreatedAt.String()
		gateway[isVPNGatewayResourceType] = *data.ResourceType
		gateway[isVPNGatewayHealthState] = *data.HealthState
		gateway[isVPNGatewayHealthReasons] = resourceVPNGatewayRouteFlattenHealthReasons(data.HealthReasons)
		gateway[isVPNGatewayLifecycleState] = *data.LifecycleState
		gateway[isVPNGatewayLifecycleReasons] = resourceVPNGatewayFlattenLifecycleReasons(data.LifecycleReasons)
		gateway[isVPNGatewayMode] = *data.Mode
		gateway[isVPNGatewayResourceGroup] = *data.ResourceGroup.ID
		gateway[isVPNGatewaySubnet] = *data.Subnet.ID
		gateway[isVPNGatewayCrn] = *data.CRN
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *data.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
		}
		gateway[isVPNGatewayTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *data.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource VPC VPN Gateway (%s) access tags: %s", d.Id(), err)
		}
		gateway[isVPNGatewayAccessTags] = accesstags
		if data.Members != nil {
			vpcMembersIpsList := make([]map[string]interface{}, 0)
			for _, memberIP := range data.Members {
				currentMemberIP := map[string]interface{}{}
				if memberIP.PublicIP != nil {
					currentMemberIP["address"] = *memberIP.PublicIP.Address
					currentMemberIP["role"] = *memberIP.Role
					vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
				}
				if memberIP.PrivateIP != nil && memberIP.PrivateIP.Address != nil {
					currentMemberIP["private_address"] = *memberIP.PrivateIP.Address
				}
			}
			gateway[isVPNGatewayMembers] = vpcMembersIpsList
		}

		if data.VPC != nil {
			vpcList := []map[string]interface{}{}
			vpcList = append(vpcList, dataSourceVPNServerCollectionVPNGatewayVpcReferenceToMap(data.VPC))
			gateway["vpc"] = vpcList
		}

		vpngateways = append(vpngateways, gateway)
	}

	d.SetId(dataSourceIBMVPNGatewaysID(d))
	d.Set(isvpnGateways, vpngateways)
	return nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVPNServerCollectionVPNGatewayVpcReferenceToMap(vpcsItem *vpcv1.VPCReference) (vpcsMap map[string]interface{}) {
	vpcsMap = map[string]interface{}{}

	if vpcsItem.CRN != nil {
		vpcsMap["crn"] = vpcsItem.CRN
	}
	if vpcsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNGatewayCollectionVpcsDeletedToMap(*vpcsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcsMap["deleted"] = deletedList
	}
	if vpcsItem.Href != nil {
		vpcsMap["href"] = vpcsItem.Href
	}
	if vpcsItem.ID != nil {
		vpcsMap["id"] = vpcsItem.ID
	}
	if vpcsItem.Name != nil {
		vpcsMap["name"] = vpcsItem.Name
	}

	return vpcsMap
}

func dataSourceVPNGatewayCollectionVpcsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
