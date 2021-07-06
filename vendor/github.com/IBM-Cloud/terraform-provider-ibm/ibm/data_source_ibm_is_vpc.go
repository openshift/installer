// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISVPC() *schema.Resource {
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
				Required:     true,
				ValidateFunc: InvokeDataSourceValidator("ibm_is_subnet", isVPCName),
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
				Set:      resourceIBMVPCHash,
			},

			isVPCCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
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
							Required:    true,
							Description: "Security group id",
							ForceNew:    true,
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
										Description: "IP version: ipv4 or ipv6",
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

func dataSourceIBMISVpcValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPCName,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString})

	ibmISVpcDataSourceValidator := ResourceValidator{ResourceName: "ibm_is_vpc", Schema: validateSchema}
	return &ibmISVpcDataSourceValidator
}

func dataSourceIBMISVPCRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	name := d.Get(isVPCName).(string)
	if userDetails.generation == 1 {
		err := classicVpcGetByName(d, meta, name)
		if err != nil {
			return err
		}
	} else {
		err := vpcGetByName(d, meta, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicVpcGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.VPC{}
	for {
		listVpcsOptions := &vpcclassicv1.ListVpcsOptions{}
		if start != "" {
			listVpcsOptions.Start = &start
		}
		vpcs, response, err := sess.ListVpcs(listVpcsOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching vpcs %s\n%s", err, response)
		}
		start = GetNext(vpcs.Next)
		allrecs = append(allrecs, vpcs.Vpcs...)
		if start == "" {
			break
		}
	}
	for _, vpc := range allrecs {
		if *vpc.Name == name {
			d.SetId(*vpc.ID)
			d.Set(isVPCName, *vpc.Name)
			d.Set(isVPCClassicAccess, *vpc.ClassicAccess)
			d.Set(isVPCStatus, *vpc.Status)
			d.Set(isVPCResourceGroup, *vpc.ResourceGroup.ID)
			if vpc.DefaultNetworkACL != nil {
				d.Set(isVPCDefaultNetworkACL, *vpc.DefaultNetworkACL.ID)
			} else {
				d.Set(isVPCDefaultNetworkACL, nil)
			}
			if vpc.DefaultSecurityGroup != nil {
				d.Set(isVPCDefaultSecurityGroup, *vpc.DefaultSecurityGroup.ID)
			} else {
				d.Set(isVPCDefaultSecurityGroup, nil)
			}
			tags, err := GetTagsUsingCRN(meta, *vpc.CRN)
			if err != nil {
				log.Printf(
					"An error occured during reading of vpc (%s) tags : %s", d.Id(), err)
			}
			d.Set(isVPCTags, tags)
			d.Set(isVPCCRN, *vpc.CRN)

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/network/vpcs")
			d.Set(ResourceName, *vpc.Name)
			d.Set(ResourceCRN, *vpc.CRN)
			d.Set(ResourceStatus, *vpc.Status)
			if vpc.ResourceGroup != nil {
				d.Set(ResourceGroupName, *vpc.ResourceGroup.ID)
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
			options := &vpcclassicv1.ListSubnetsOptions{}
			s, response, err := sess.ListSubnets(options)
			if err != nil {
				log.Printf("Error Fetching subnets %s\n%s", err, response)
			} else {
				subnetsInfo := make([]map[string]interface{}, 0)
				for _, subnet := range s.Subnets {
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

			//Set Security group list

			listSgOptions := &vpcclassicv1.ListSecurityGroupsOptions{}
			sgs, _, err := sess.ListSecurityGroups(listSgOptions)
			if err != nil {
				return err
			}

			securityGroupList := make([]map[string]interface{}, 0)

			for _, group := range sgs.SecurityGroups {
				if *group.VPC.ID == d.Id() {
					g := make(map[string]interface{})

					g[isVPCSecurityGroupName] = *group.Name
					g[isVPCSecurityGroupID] = *group.ID

					rules := make([]map[string]interface{}, 0)
					for _, sgrule := range group.Rules {
						switch reflect.TypeOf(sgrule).String() {
						case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
							{
								rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
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
								remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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

						case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
							{
								rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
								r := make(map[string]interface{})
								r[isVPCSecurityGroupRuleDirection] = *rule.Direction
								r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
								if rule.Protocol != nil {
									r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
								}
								r[isVPCSecurityGroupRuleID] = *rule.ID
								remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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

						case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
							{
								rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
								r := make(map[string]interface{})
								r[isVPCSecurityGroupRuleDirection] = *rule.Direction
								r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
								if rule.PortMin != nil {
									r[isVPCSecurityGroupRulePortMin] = int(*rule.PortMin)
								}
								if rule.PortMax != nil {
									r[isVPCSecurityGroupRulePortMax] = int(*rule.PortMax)
								}

								if rule.Protocol != nil {
									r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
								}

								r[isVPCSecurityGroupRuleID] = *rule.ID
								remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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
			}

			d.Set(isVPCSecurityGroupList, securityGroupList)
			return nil
		}
	}
	return fmt.Errorf("No VPC found with name %s", name)
}
func vpcGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.VPC{}
	for {
		listVpcsOptions := &vpcv1.ListVpcsOptions{}
		if start != "" {
			listVpcsOptions.Start = &start
		}
		vpcs, response, err := sess.ListVpcs(listVpcsOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching vpcs %s\n%s", err, response)
		}
		start = GetNext(vpcs.Next)
		allrecs = append(allrecs, vpcs.Vpcs...)
		if start == "" {
			break
		}
	}
	for _, vpc := range allrecs {
		if *vpc.Name == name {
			d.SetId(*vpc.ID)
			d.Set(isVPCName, *vpc.Name)
			d.Set(isVPCClassicAccess, *vpc.ClassicAccess)
			d.Set(isVPCStatus, *vpc.Status)
			d.Set(isVPCResourceGroup, *vpc.ResourceGroup.ID)
			d.Set(isVPCDefaultNetworkACLName, *vpc.DefaultNetworkACL.Name)
			d.Set(isVPCDefaultRoutingTableName, *vpc.DefaultRoutingTable.Name)
			d.Set(isVPCDefaultSecurityGroupName, *vpc.DefaultSecurityGroup.Name)
			if vpc.DefaultNetworkACL != nil {
				d.Set(isVPCDefaultNetworkACL, *vpc.DefaultNetworkACL.ID)
			} else {
				d.Set(isVPCDefaultNetworkACL, nil)
			}
			if vpc.DefaultRoutingTable != nil {
				d.Set(isVPCDefaultRoutingTable, *vpc.DefaultRoutingTable.ID)
			}
			if vpc.DefaultSecurityGroup != nil {
				d.Set(isVPCDefaultSecurityGroup, *vpc.DefaultSecurityGroup.ID)
			} else {
				d.Set(isVPCDefaultSecurityGroup, nil)
			}
			tags, err := GetTagsUsingCRN(meta, *vpc.CRN)
			if err != nil {
				log.Printf(
					"An error occured during reading of vpc (%s) tags : %s", d.Id(), err)
			}
			d.Set(isVPCTags, tags)
			d.Set(isVPCCRN, *vpc.CRN)

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc-ext/network/vpcs")
			d.Set(ResourceName, *vpc.Name)
			d.Set(ResourceCRN, *vpc.CRN)
			d.Set(ResourceStatus, *vpc.Status)
			if vpc.ResourceGroup != nil {
				d.Set(ResourceGroupName, *vpc.ResourceGroup.Name)
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
			options := &vpcv1.ListSubnetsOptions{}
			s, response, err := sess.ListSubnets(options)
			if err != nil {
				log.Printf("Error Fetching subnets %s\n%s", err, response)
			} else {
				subnetsInfo := make([]map[string]interface{}, 0)
				for _, subnet := range s.Subnets {
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

			listSgOptions := &vpcv1.ListSecurityGroupsOptions{}
			sgs, _, err := sess.ListSecurityGroups(listSgOptions)
			if err != nil {
				return err
			}

			securityGroupList := make([]map[string]interface{}, 0)

			for _, group := range sgs.SecurityGroups {
				if *group.VPC.ID == d.Id() {
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
			}

			d.Set(isVPCSecurityGroupList, securityGroupList)

			return nil
		}
	}
	return fmt.Errorf("No VPC found with name %s", name)
}
