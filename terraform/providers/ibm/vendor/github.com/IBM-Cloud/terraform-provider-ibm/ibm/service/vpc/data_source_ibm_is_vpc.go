// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVPCRead,

		Schema: map[string]*schema.Schema{
			isVPCDefaultNetworkACL: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isVPCClassicAccess: {
				Type:     schema.TypeBool,
				Computed: true,
			},

			isVPCDefaultRoutingTable: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default routing table associated with VPC",
			},

			isVPCName: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isVPCName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_subnet", isVPCName),
			},
			"identifier": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{isVPCName, "identifier"},
				Optional:     true,
			},

			isVPCDefaultNetworkACLName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Network ACL name",
			},

			isVPCDefaultSecurityGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default security group name",
			},

			isVPCDefaultSecurityGroupCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default security group CRN",
			},

			isVPCDefaultNetworkACLCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Network ACL CRN",
			},

			isVPCDefaultRoutingTableName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default routing table name",
			},

			isVPCResourceGroup: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isVPCStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isVPCDefaultSecurityGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group associated with VPC",
			},

			isVPCTags: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      flex.ResourceIBMVPCHash,
			},

			isVPCAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},

			isVPCCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			cseSourceAddresses: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud service endpoint IP Address",
						},

						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location info of CSE Address",
						},
					},
				},
			},

			isVPCDns: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The DNS configuration for this VPC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPCDnsEnableHub: &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this VPC is enabled as a DNS name resolution hub.",
						},
						isVPCDnsResolutionBindingCount: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of DNS resolution bindings for this VPC.",
						},
						isVPCDnsResolver: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DNS resolver configuration for the VPC.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVPCDnsResolverServers: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The DNS servers for this VPC. The servers are populated:- by the system when `dns.resolver.type` is `system`- using the DNS servers in `dns.resolver.vpc` when `dns.resolver.type` is `delegated`- using `dns.resolver.manual_servers` when the `dns.resolver.type` is `manual`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVPCDnsResolverManualServersAddress: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												isVPCDnsResolverManualServersZoneAffinity: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Zone name, if present, DHCP configuration for this zone will have this DNS server listed first.",
												},
											},
										},
									},
									isVPCDnsResolverType: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the DNS resolver used for the VPC.- `delegated`: DNS server addresses are provided by the DNS resolver of the VPC               specified in `dns.resolver.vpc`.- `manual`: DNS server addresses are specified in `dns.resolver.manual_servers`.- `system`: DNS server addresses are provided by the system.",
									},
									isVPCDnsResolverVpc: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The VPC whose DNS resolver provides the DNS server addresses for this VPC.The VPC may be remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVPCDnsResolverVpcCrn: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this VPC.",
												},
												isVPCDnsResolverVpcDeleted: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolverVpcDeletedMoreInfo: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												isVPCDnsResolverVpcHref: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this VPC.",
												},
												isVPCDnsResolverVpcId: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this VPC.",
												},
												isVPCDnsResolverVpcName: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this VPC. The name is unique across all VPCs in the region.",
												},
												isVPCDnsResolverVpcRemote: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVPCDnsResolverVpcRemoteAccount: &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The unique identifier for this account.",
																		},
																		isVPCDnsResolverResourceType: &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The resource type.",
																		},
																	},
																},
															},
															isVPCDnsResolverVpcRemoteRegion: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region name. If present, this property indicates that the referenced resource is remote to this region, and identifies the native region.",
															},
														},
													},
												},
												isVPCDnsResolverResourceType: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									isVPCDnsResolverManualServers: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The manually specified DNS servers for this VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVPCDnsResolverManualServersAddress: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												isVPCDnsResolverManualServersZoneAffinity: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the zone. If present, DHCP configuration for this zone will have this DNS server listed first.",
												},
											},
										},
									},
									isVPCDnsResolverConfiguration: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration of the system DNS resolver for this VPC.- `custom_resolver`: A custom DNS resolver is configured for this VPC.- `private_resolver`: A private DNS resolver is configured for this VPC. Applicable when  the VPC has either or both of the following:    - at least one endpoint gateway residing in it    - a [DNS Services](https://cloud.ibm.com/docs/dns-svcs) private zone configured for it- `default`: The provider default DNS resolvers are configured for this VPC.  This system DNS resolver configuration is used when the VPC has:  - no custom DNS resolver configured for it, and  - no endpoint gateways residing in it, and  - no [DNS Services](https://cloud.ibm.com/docs/dns-svcs) private zone configured for it.",
									},
								},
							},
						},
					},
				},
			},

			"health_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `health_state` (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},
			"health_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
			},
			isVPCSecurityGroupList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPCSecurityGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group name",
						},

						isVPCSecurityGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group id",
						},

						isSecurityGroupRules: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Security Rules",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									isVPCSecurityGroupRuleID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule ID",
									},

									isVPCSecurityGroupRuleDirection: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Direction of traffic to enforce, either inbound or outbound",
									},

									isVPCSecurityGroupRuleIPVersion: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP version: ipv4",
									},

									isVPCSecurityGroupRuleRemote: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
									},

									isVPCSecurityGroupRuleType: {
										Type:     schema.TypeInt,
										Computed: true,
									},

									isVPCSecurityGroupRuleCode: {
										Type:     schema.TypeInt,
										Computed: true,
									},

									isVPCSecurityGroupRulePortMin: {
										Type:     schema.TypeInt,
										Computed: true,
									},

									isVPCSecurityGroupRulePortMax: {
										Type:     schema.TypeInt,
										Computed: true,
									},

									isVPCSecurityGroupRuleProtocol: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			subnetsList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subent name",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet ID",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet status",
						},

						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet location",
						},

						totalIPV4AddressCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total IPv4 address count in the subnet",
						},

						availableIPV4AddressCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Available IPv4 address count in the subnet",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMISVpcValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPCName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "is",
			CloudDataRange:             []string{"service:vpc", "resolved_to:id"}})

	ibmISVpcDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_vpc", Schema: validateSchema}
	return &ibmISVpcDataSourceValidator
}

func dataSourceIBMISVPCRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get(isVPCName).(string)
	id := d.Get("identifier").(string)
	err := vpcGetByNameOrId(d, meta, name, id)
	if err != nil {
		return err
	}
	return nil
}

func vpcGetByNameOrId(d *schema.ResourceData, meta interface{}, name, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	flag := false
	if id != "" {
		getVpcsOptions := &vpcv1.GetVPCOptions{
			ID: &id,
		}
		vpcGet, response, err := sess.GetVPC(getVpcsOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching vpc %s\n%s", err, response)
		}
		flag = true
		setVpcDetails(d, vpcGet, meta, sess)
	} else {
		start := ""
		allrecs := []vpcv1.VPC{}
		for {
			listVpcsOptions := &vpcv1.ListVpcsOptions{}
			if start != "" {
				listVpcsOptions.Start = &start
			}
			vpcs, response, err := sess.ListVpcs(listVpcsOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Fetching vpcs %s\n%s", err, response)
			}
			start = flex.GetNext(vpcs.Next)
			allrecs = append(allrecs, vpcs.Vpcs...)
			if start == "" {
				break
			}
		}
		for _, v := range allrecs {
			if *v.Name == name {
				flag = true
				setVpcDetails(d, &v, meta, sess)
			}
		}
	}
	if !flag {
		return fmt.Errorf("[ERROR] No VPC found with name %s", name)
	}
	return nil
}

func setVpcDetails(d *schema.ResourceData, vpc *vpcv1.VPC, meta interface{}, sess *vpcv1.VpcV1) error {
	if vpc != nil {
		d.SetId(*vpc.ID)
		d.Set("identifier", *vpc.ID)
		d.Set(isVPCName, *vpc.Name)
		d.Set(isVPCClassicAccess, *vpc.ClassicAccess)
		d.Set(isVPCStatus, *vpc.Status)
		if vpc.ResourceGroup != nil {
			d.Set(isVPCResourceGroup, *vpc.ResourceGroup.ID)
		}
		if vpc.DefaultNetworkACL != nil {
			d.Set(isVPCDefaultNetworkACLName, *vpc.DefaultNetworkACL.Name)
			d.Set(isVPCDefaultNetworkACL, *vpc.DefaultNetworkACL.ID)
			d.Set(isVPCDefaultNetworkACLCRN, vpc.DefaultNetworkACL.CRN)
		} else {
			d.Set(isVPCDefaultNetworkACL, nil)
		}
		if vpc.DefaultRoutingTable != nil {
			d.Set(isVPCDefaultRoutingTableName, *vpc.DefaultRoutingTable.Name)
			d.Set(isVPCDefaultRoutingTable, *vpc.DefaultRoutingTable.ID)
		}
		if vpc.DefaultSecurityGroup != nil {
			d.Set(isVPCDefaultSecurityGroupName, *vpc.DefaultSecurityGroup.Name)
			d.Set(isVPCDefaultSecurityGroup, *vpc.DefaultSecurityGroup.ID)
			d.Set(isVPCDefaultSecurityGroupCRN, vpc.DefaultSecurityGroup.CRN)
		} else {
			d.Set(isVPCDefaultSecurityGroup, nil)
		}
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *vpc.CRN, "", isVPCUserTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of vpc (%s) tags : %s", d.Id(), err)
		}
		d.Set(isVPCTags, tags)
		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpc.CRN, "", isVPCAccessTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of vpc (%s) access tags: %s", d.Id(), err)
		}
		d.Set(isVPCAccessTags, accesstags)
		d.Set(isVPCCRN, *vpc.CRN)

		healthReasons := []map[string]interface{}{}
		if vpc.HealthReasons != nil {
			for _, modelItem := range vpc.HealthReasons {
				modelMap, err := dataSourceIBMIsVPCVPCHealthReasonToMap(&modelItem)
				if err != nil {
					return err
				}
				healthReasons = append(healthReasons, modelMap)
			}
		}
		if err = d.Set("health_reasons", healthReasons); err != nil {
			return fmt.Errorf("[ERROR] Error setting health_reasons %s", err)
		}

		if err = d.Set("health_state", vpc.HealthState); err != nil {
			return fmt.Errorf("[ERROR] Error setting health_state: %s", err)
		}

		controller, err := flex.GetBaseController(meta)
		if err != nil {
			return err
		}
		d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/vpcs")
		d.Set(flex.ResourceName, *vpc.Name)
		d.Set(flex.ResourceCRN, *vpc.CRN)
		d.Set(flex.ResourceStatus, *vpc.Status)
		if vpc.ResourceGroup != nil {
			d.Set(flex.ResourceGroupName, *vpc.ResourceGroup.Name)
		}
		//set the cse ip addresses info
		if vpc.CseSourceIps != nil {
			cseSourceIpsList := make([]map[string]interface{}, 0)
			for _, sourceIP := range vpc.CseSourceIps {
				currentCseSourceIp := map[string]interface{}{}
				if sourceIP.IP != nil {
					currentCseSourceIp["address"] = *sourceIP.IP.Address
					currentCseSourceIp["zone_name"] = *sourceIP.Zone.Name
					cseSourceIpsList = append(cseSourceIpsList, currentCseSourceIp)
				}
			}
			d.Set(cseSourceAddresses, cseSourceIpsList)
		}

		// adding pagination support for subnets inside vpc

		startSub := ""
		allrecsSub := []vpcv1.Subnet{}
		options := &vpcv1.ListSubnetsOptions{}

		for {
			if startSub != "" {
				options.Start = &startSub
			}
			s, response, err := sess.ListSubnets(options)
			if err != nil {
				return fmt.Errorf("[ERROR] Error fetching subnets %s\n%s", err, response)
			}
			startSub = flex.GetNext(s.Next)
			allrecsSub = append(allrecsSub, s.Subnets...)
			if startSub == "" {
				break
			}
		}
		if err == nil {
			subnetsInfo := make([]map[string]interface{}, 0)
			for _, subnet := range allrecsSub {
				if *subnet.VPC.ID == d.Id() {
					l := map[string]interface{}{
						"name":                    *subnet.Name,
						"id":                      *subnet.ID,
						"status":                  *subnet.Status,
						"zone":                    *subnet.Zone.Name,
						totalIPV4AddressCount:     *subnet.TotalIpv4AddressCount,
						availableIPV4AddressCount: *subnet.AvailableIpv4AddressCount,
					}
					subnetsInfo = append(subnetsInfo, l)
				}
			}
			d.Set(subnetsList, subnetsInfo)
		}

		// adding pagination support for sg inside vpc

		startSg := ""
		allrecsSg := []vpcv1.SecurityGroup{}

		for {
			vpcId := d.Id()
			listSgOptions := &vpcv1.ListSecurityGroupsOptions{
				VPCID: &vpcId,
			}
			if startSg != "" {
				listSgOptions.Start = &startSg
			}
			sgs, response, err := sess.ListSecurityGroups(listSgOptions)
			if err != nil || sgs == nil {
				return fmt.Errorf("[ERROR] Error fetching Security Groups %s\n%s", err, response)
			}
			if *sgs.TotalCount == int64(0) {
				break
			}
			startSg = flex.GetNext(sgs.Next)
			allrecsSg = append(allrecsSg, sgs.SecurityGroups...)

			if startSg == "" {
				break
			}

		}

		securityGroupList := make([]map[string]interface{}, 0)

		for _, group := range allrecsSg {
			g := make(map[string]interface{})

			g[isVPCSecurityGroupName] = *group.Name
			g[isVPCSecurityGroupID] = *group.ID

			rules := make([]map[string]interface{}, 0)
			for _, sgrule := range group.Rules {
				switch reflect.TypeOf(sgrule).String() {
				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
						r := make(map[string]interface{})
						if rule.Code != nil {
							r[isVPCSecurityGroupRuleCode] = int(*rule.Code)
						}
						if rule.Type != nil {
							r[isVPCSecurityGroupRuleType] = int(*rule.Type)
						}
						r[isVPCSecurityGroupRuleDirection] = *rule.Direction
						r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
						}
						r[isVPCSecurityGroupRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
								}
							}
						}
						rules = append(rules, r)
					}

				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
						r := make(map[string]interface{})
						r[isVPCSecurityGroupRuleDirection] = *rule.Direction
						r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
						}
						r[isVPCSecurityGroupRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
								}
							}
						}
						rules = append(rules, r)
					}

				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
						r := make(map[string]interface{})
						r[isVPCSecurityGroupRuleDirection] = *rule.Direction
						r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
						if rule.PortMin != nil {
							r[isVPCSecurityGroupRulePortMin] = int(*rule.PortMin)
						}
						if rule.PortMax != nil {
							r[isVPCSecurityGroupRulePortMax] = int(*rule.PortMax)
						}
						r[isVPCSecurityGroupRuleID] = *rule.ID
						if rule.Protocol != nil {
							r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
						}

						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
								}
							}
						}
						rules = append(rules, r)
					}
				}
			}
			g[isVPCSgRules] = rules
			securityGroupList = append(securityGroupList, g)
		}
		d.Set(isVPCSecurityGroupList, securityGroupList)
		if !core.IsNil(vpc.Dns) {
			dnsMap, err := dataSourceIBMIsVPCVpcdnsToMap(vpc.Dns)
			if err != nil {
				return err
			}
			if err = d.Set(isVPCDns, []map[string]interface{}{dnsMap}); err != nil {
				return fmt.Errorf("[ERROR] Error setting dns: %s", err)
			}
		}
		return nil
	}
	return nil
}

func dataSourceIBMIsVPCVpcdnsToMap(model *vpcv1.Vpcdns) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["enable_hub"] = model.EnableHub
	modelMap["resolution_binding_count"] = flex.IntValue(model.ResolutionBindingCount)
	resolverMap, err := dataSourceIBMIsVPCVpcdnsResolverToMap(model.Resolver)
	if err != nil {
		return modelMap, err
	}
	modelMap["resolver"] = []map[string]interface{}{resolverMap}
	return modelMap, nil
}

func dataSourceIBMIsVPCVpcdnsResolverToMap(model vpcv1.VpcdnsResolverIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VpcdnsResolverTypeDelegated); ok {
		return dataSourceIBMIsVPCVpcdnsResolverTypeDelegatedToMap(model.(*vpcv1.VpcdnsResolverTypeDelegated))
	} else if _, ok := model.(*vpcv1.VpcdnsResolverTypeManual); ok {
		return dataSourceIBMIsVPCVpcdnsResolverTypeManualToMap(model.(*vpcv1.VpcdnsResolverTypeManual))
	} else if _, ok := model.(*vpcv1.VpcdnsResolverTypeSystem); ok {
		return dataSourceIBMIsVPCVpcdnsResolverTypeSystemToMap(model.(*vpcv1.VpcdnsResolverTypeSystem))
	} else if _, ok := model.(*vpcv1.VpcdnsResolver); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VpcdnsResolver)
		servers := []map[string]interface{}{}
		for _, serversItem := range model.Servers {
			serversItemMap, err := resourceIBMIsVPCDnsServerToMap(&serversItem)
			if err != nil {
				return modelMap, err
			}
			servers = append(servers, serversItemMap)
		}
		modelMap["servers"] = servers
		modelMap["type"] = model.Type
		if model.VPC != nil {
			vpcMap, err := dataSourceIBMIsVPCVPCReferenceDnsResolverContextToMap(model.VPC)
			if err != nil {
				return modelMap, err
			}
			modelMap["vpc"] = []map[string]interface{}{vpcMap}
		}

		if model.ManualServers != nil {
			manualServers := []map[string]interface{}{}
			for _, manualServersItem := range model.ManualServers {
				manualServersItemMap, err := resourceIBMIsVPCDnsServerToMap(&manualServersItem)
				if err != nil {
					return modelMap, err
				}
				manualServers = append(manualServers, manualServersItemMap)
			}
			modelMap["manual_servers"] = manualServers
		}
		if model.Configuration != nil {
			modelMap["configuration"] = model.Configuration
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VpcdnsResolverIntf subtype encountered")
	}
}

func dataSourceIBMIsVPCVpcdnsResolverTypeManualToMap(model *vpcv1.VpcdnsResolverTypeManual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	servers := []map[string]interface{}{}
	for _, serversItem := range model.Servers {
		serversItemMap, err := resourceIBMIsVPCDnsServerToMap(&serversItem)
		if err != nil {
			return modelMap, err
		}
		servers = append(servers, serversItemMap)
	}
	modelMap["servers"] = servers
	manualServers := []map[string]interface{}{}
	for _, manualServersItem := range model.ManualServers {
		manualServersItemMap, err := resourceIBMIsVPCDnsServerToMap(&manualServersItem)
		if err != nil {
			return modelMap, err
		}
		manualServers = append(manualServers, manualServersItemMap)
	}
	modelMap["manual_servers"] = manualServers
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsVPCVpcdnsResolverTypeDelegatedToMap(model *vpcv1.VpcdnsResolverTypeDelegated) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	servers := []map[string]interface{}{}
	for _, serversItem := range model.Servers {
		serversItemMap, err := resourceIBMIsVPCDnsServerToMap(&serversItem)
		if err != nil {
			return modelMap, err
		}
		servers = append(servers, serversItemMap)
	}
	modelMap["servers"] = servers
	modelMap["type"] = model.Type
	if model.VPC != nil {
		vpcMap, err := dataSourceIBMIsVPCVPCReferenceDnsResolverContextToMap(model.VPC)
		if err != nil {
			return modelMap, err
		}
		modelMap["vpc"] = []map[string]interface{}{vpcMap}
	}
	return modelMap, nil
}
func dataSourceIBMIsVPCVPCHealthReasonToMap(model *vpcv1.VPCHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIsVPCVpcdnsResolverTypeSystemToMap(model *vpcv1.VpcdnsResolverTypeSystem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	servers := []map[string]interface{}{}
	for _, serversItem := range model.Servers {
		serversItemMap, err := resourceIBMIsVPCDnsServerToMap(&serversItem)
		if err != nil {
			return modelMap, err
		}
		servers = append(servers, serversItemMap)
	}
	modelMap["servers"] = servers
	modelMap["configuration"] = model.Configuration
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsVPCVPCReferenceDnsResolverContextToMap(model *vpcv1.VPCReferenceDnsResolverContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsVPCVPCReferenceDnsResolverContextDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := dataSourceIBMIsVPCVPCRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsVPCVPCReferenceDnsResolverContextDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsVPCVPCRemoteToMap(model *vpcv1.VPCRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := dataSourceIBMIsVPCAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := dataSourceIBMIsVPCRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func dataSourceIBMIsVPCAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsVPCRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}
