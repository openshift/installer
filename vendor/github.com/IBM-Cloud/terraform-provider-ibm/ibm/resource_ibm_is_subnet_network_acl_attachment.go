// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSubnetID     = "subnet"
	isNetworkACLID = "network_acl"
)

func resourceIBMISSubnetNetworkACLAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISSubnetNetworkACLAttachmentCreate,
		Read:     resourceIBMISSubnetNetworkACLAttachmentRead,
		Update:   resourceIBMISSubnetNetworkACLAttachmentUpdate,
		Delete:   resourceIBMISSubnetNetworkACLAttachmentDelete,
		Exists:   resourceIBMISSubnetNetworkACLAttachmentExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isSubnetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier",
			},

			isNetworkACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of network ACL",
			},

			isNetworkACLName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL name",
			},

			isNetworkACLVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL VPC",
			},

			isNetworkACLResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group ID for the network ACL",
			},

			isNetworkACLRules: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this Network ACL rule",
						},
						isNetworkACLRuleName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this rule",
						},
						isNetworkACLRuleAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to allow or deny matching traffic",
						},
						isNetworkACLRuleIPVersion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version for this rule",
						},
						isNetworkACLRuleSource: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source CIDR block",
						},
						isNetworkACLRuleDestination: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination CIDR block",
						},
						isNetworkACLRuleDirection: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direction of traffic to enforce, either inbound or outbound",
						},
						isNetworkACLRuleICMP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRuleICMPCode: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic code to allow",
									},
									isNetworkACLRuleICMPType: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic type to allow",
									},
								},
							},
						},

						isNetworkACLRuleTCP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of TCP destination port range",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of TCP destination port range",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of TCP source port range",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of TCP source port range",
									},
								},
							},
						},

						isNetworkACLRuleUDP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of UDP destination port range",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of UDP destination port range",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of UDP source port range",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of UDP source port range",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIBMISSubnetNetworkACLAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnet := d.Get(isSubnetID).(string)
	networkACL := d.Get(isNetworkACLID).(string)

	// Construct an instance of the NetworkACLIdentityByID model
	networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
	networkACLIdentityModel.ID = &networkACL

	// Construct an instance of the ReplaceSubnetNetworkACLOptions model
	replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
	replaceSubnetNetworkACLOptionsModel.ID = &subnet
	replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
	resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

	if err != nil {
		log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
		return fmt.Errorf("Error while attaching a network ACL to a subnet %s\n%s", err, response)
	}
	d.SetId(subnet)
	log.Printf("[INFO] Network ACL : %s", *resultACL.ID)
	log.Printf("[INFO] Subnet ID : %s", subnet)

	return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
}

func resourceIBMISSubnetNetworkACLAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSubnetNetworkACLOptionsModel := &vpcv1.GetSubnetNetworkACLOptions{
		ID: &id,
	}
	nwacl, response, err := sess.GetSubnetNetworkACL(getSubnetNetworkACLOptionsModel)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting subnet's (%s) attached network ACL: %s\n%s", id, err, response)
	}
	d.Set(isNetworkACLName, *nwacl.Name)
	d.Set(isNetworkACLVPC, *nwacl.VPC.ID)
	if nwacl.ResourceGroup != nil {
		d.Set(isNetworkACLResourceGroup, *nwacl.ResourceGroup.ID)
	}

	rules := make([]interface{}, 0)
	if len(nwacl.Rules) > 0 {
		for _, rulex := range nwacl.Rules {
			log.Println("[DEBUG] Type of the Rule", reflect.TypeOf(rulex))
			rule := make(map[string]interface{})
			switch reflect.TypeOf(rulex).String() {
			case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp":
				{
					rulex := rulex.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp)
					rule[isNetworkACLRuleID] = *rulex.ID
					rule[isNetworkACLRuleName] = *rulex.Name
					rule[isNetworkACLRuleAction] = *rulex.Action
					rule[isNetworkACLRuleIPVersion] = *rulex.IPVersion
					rule[isNetworkACLRuleSource] = *rulex.Source
					rule[isNetworkACLRuleDestination] = *rulex.Destination
					rule[isNetworkACLRuleDirection] = *rulex.Direction
					rule[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
					rule[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
					icmp := make([]map[string]int, 1, 1)
					if rulex.Code != nil && rulex.Type != nil {
						icmp[0] = map[string]int{
							isNetworkACLRuleICMPCode: int(*rulex.Code),
							isNetworkACLRuleICMPType: int(*rulex.Code),
						}
					}
					rule[isNetworkACLRuleICMP] = icmp
				}
			case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp":
				{
					rulex := rulex.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp)
					rule[isNetworkACLRuleID] = *rulex.ID
					rule[isNetworkACLRuleName] = *rulex.Name
					rule[isNetworkACLRuleAction] = *rulex.Action
					rule[isNetworkACLRuleIPVersion] = *rulex.IPVersion
					rule[isNetworkACLRuleSource] = *rulex.Source
					rule[isNetworkACLRuleDestination] = *rulex.Destination
					rule[isNetworkACLRuleDirection] = *rulex.Direction
					if *rulex.Protocol == "tcp" {
						rule[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
						rule[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
						tcp := make([]map[string]int, 1, 1)
						tcp[0] = map[string]int{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						}
						tcp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
						tcp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
						rule[isNetworkACLRuleTCP] = tcp
					} else if *rulex.Protocol == "udp" {
						rule[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
						rule[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
						udp := make([]map[string]int, 1, 1)
						udp[0] = map[string]int{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						}
						udp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
						udp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
						rule[isNetworkACLRuleUDP] = udp
					}
				}
			case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll":
				{
					rulex := rulex.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll)
					rule[isNetworkACLRuleID] = *rulex.ID
					rule[isNetworkACLRuleName] = *rulex.Name
					rule[isNetworkACLRuleAction] = *rulex.Action
					rule[isNetworkACLRuleIPVersion] = *rulex.IPVersion
					rule[isNetworkACLRuleSource] = *rulex.Source
					rule[isNetworkACLRuleDestination] = *rulex.Destination
					rule[isNetworkACLRuleDirection] = *rulex.Direction
					rule[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
					rule[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
					rule[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
				}
			}
			rules = append(rules, rule)
		}
	}
	d.Set(isNetworkACLRules, rules)
	return nil
}

func resourceIBMISSubnetNetworkACLAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isNetworkACLID) {
		subnet := d.Get(isSubnetID).(string)
		networkACL := d.Get(isNetworkACLID).(string)

		// Construct an instance of the NetworkACLIdentityByID model
		networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
		networkACLIdentityModel.ID = &networkACL

		// Construct an instance of the ReplaceSubnetNetworkACLOptions model
		replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
		replaceSubnetNetworkACLOptionsModel.ID = &subnet
		replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
		resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
			return fmt.Errorf("Error while attaching a network ACL to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with Network ACL : %s", subnet, *resultACL.ID)

		d.SetId(subnet)
		return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
	}

	return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
}

func resourceIBMISSubnetNetworkACLAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	// Set the subnet with VPC default network ACL
	getSubnetOptions := &vpcv1.GetSubnetOptions{
		ID: &id,
	}
	subnet, response, err := sess.GetSubnet(getSubnetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Subnet (%s): %s\n%s", id, err, response)
	}
	// Fetch VPC
	vpcID := *subnet.VPC.ID

	getvpcOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getvpcOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting VPC : %s\n%s", err, response)
	}

	// Fetch default network ACL
	if vpc.DefaultNetworkACL != nil {
		log.Printf("[DEBUG] vpc default network acl is not null :%s", *vpc.DefaultNetworkACL.ID)
		// Construct an instance of the NetworkACLIdentityByID model
		networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
		networkACLIdentityModel.ID = vpc.DefaultNetworkACL.ID

		// Construct an instance of the ReplaceSubnetNetworkACLOptions model
		replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
		replaceSubnetNetworkACLOptionsModel.ID = &id
		replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
		resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
			return fmt.Errorf("Error while attaching a network ACL to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with VPC default Network ACL : %s", id, *resultACL.ID)
	} else {
		log.Printf("[DEBUG] vpc default network acl is  null")
	}

	d.SetId("")
	return nil
}

func resourceIBMISSubnetNetworkACLAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSubnetNetworkACLOptionsModel := &vpcv1.GetSubnetNetworkACLOptions{
		ID: &id,
	}
	_, response, err := sess.GetSubnetNetworkACL(getSubnetNetworkACLOptionsModel)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting subnet's attached network ACL: %s\n%s", err, response)
	}
	return true, nil
}
