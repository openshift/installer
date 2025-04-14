// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
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

		return nil
	}
	return nil
}
