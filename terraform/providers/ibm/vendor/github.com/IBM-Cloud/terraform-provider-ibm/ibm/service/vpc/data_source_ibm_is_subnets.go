// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSubnets                = "subnets"
	isSubnetResourceGroupID  = "resource_group"
	isSubnetRoutingTableName = "routing_table_name"
	isSubnetResourceZone     = "zone"
	isSubnetResourceVpc      = "vpc"
	isSubnetResourceVpcCrn   = "vpc_crn"
	isSubnetResourceVpcName  = "vpc_name"
)

func DataSourceIBMISSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSubnetsRead,

		Schema: map[string]*schema.Schema{
			isSubnetResourceVpc: {
				Type:        schema.TypeString,
				Description: "ID of the VPC",
				Optional:    true,
			},
			isSubnetResourceVpcName: {
				Type:        schema.TypeString,
				Description: "Name of the VPC",
				Optional:    true,
			},
			isSubnetResourceVpcCrn: {
				Type:        schema.TypeString,
				Description: "CRN of the VPC",
				Optional:    true,
			},
			isSubnetResourceZone: {
				Type:        schema.TypeString,
				Description: "Name of the Zone ",
				Optional:    true,
			},
			isSubnetResourceGroupID: {
				Type:        schema.TypeString,
				Description: "Resource Group ID",
				Optional:    true,
			},

			isSubnetRoutingTableName: {
				Type:        schema.TypeString,
				Description: "Name of the routing table",
				Optional:    true,
			},

			isSubnetRoutingTableID: {
				Type:        schema.TypeString,
				Description: "ID of the routing table",
				Optional:    true,
			},

			isSubnets: {
				Type:        schema.TypeList,
				Description: "List of subnets",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"crn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_ipv4_address_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isSubnetResourceGroupID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_ipv4_address_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routing_table": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The routing table for this subnet",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this routing table.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this routing table.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this routing table.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The crn for this routing table.",
									},
									"resource_type": {
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
		},
	}
}

func dataSourceIBMISSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	err := subnetList(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func subnetList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.Subnet{}

	var resourceGroup string
	if v, ok := d.GetOk(isSubnetResourceGroupID); ok {
		resourceGroup = v.(string)
	}

	var routingTable string
	if v, ok := d.GetOk(isSubnetRoutingTableID); ok {
		routingTable = v.(string)
	}

	var resourceTableName string
	if v, ok := d.GetOk(isSubnetRoutingTableName); ok {
		resourceTableName = v.(string)
	}

	var zone string
	if v, ok := d.GetOk(isSubnetResourceZone); ok {
		zone = v.(string)
	}

	var vpc string
	if v, ok := d.GetOk(isSubnetResourceVpc); ok {
		vpc = v.(string)
	}

	var vpcName string
	if v, ok := d.GetOk(isSubnetResourceVpcName); ok {
		vpcName = v.(string)
	}

	var vpcCrn string
	if v, ok := d.GetOk(isSubnetResourceVpcCrn); ok {
		vpcCrn = v.(string)
	}

	for {
		options := &vpcv1.ListSubnetsOptions{}
		if resourceGroup != "" {
			options.SetResourceGroupID(resourceGroup)
		}
		if routingTable != "" {
			options.SetRoutingTableID(routingTable)
		}
		if resourceTableName != "" {
			options.SetRoutingTableName(resourceTableName)
		}
		if zone != "" {
			options.SetZoneName(zone)
		}
		if vpc != "" {
			options.SetVPCID(vpc)
		}
		if vpcName != "" {
			options.SetVPCName(vpcName)
		}
		if vpcCrn != "" {
			options.SetVPCCRN(vpcCrn)
		}
		if start != "" {
			options.Start = &start
		}
		subnets, response, err := sess.ListSubnets(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching subnets %s\n%s", err, response)
		}
		start = flex.GetNext(subnets.Next)
		allrecs = append(allrecs, subnets.Subnets...)
		if start == "" {
			break
		}
	}
	subnetsInfo := make([]map[string]interface{}, 0)
	for _, subnet := range allrecs {

		var aac string = strconv.FormatInt(*subnet.AvailableIpv4AddressCount, 10)
		var tac string = strconv.FormatInt(*subnet.TotalIpv4AddressCount, 10)
		l := map[string]interface{}{
			"name":                         *subnet.Name,
			"id":                           *subnet.ID,
			"status":                       *subnet.Status,
			"crn":                          *subnet.CRN,
			"ipv4_cidr_block":              *subnet.Ipv4CIDRBlock,
			"available_ipv4_address_count": aac,
			"network_acl":                  *subnet.NetworkACL.Name,
			"total_ipv4_address_count":     tac,
			"vpc":                          *subnet.VPC.ID,
			"zone":                         *subnet.Zone.Name,
		}
		if subnet.PublicGateway != nil {
			l["public_gateway"] = *subnet.PublicGateway.ID
		}
		if subnet.ResourceGroup != nil {
			l["resource_group"] = *subnet.ResourceGroup.ID
		}
		if subnet.RoutingTable != nil {
			l["routing_table"] = dataSourceSubnetFlattenroutingTable(*subnet.RoutingTable)
		}
		subnetsInfo = append(subnetsInfo, l)
	}
	d.SetId(dataSourceIBMISSubnetsID(d))
	d.Set(isSubnets, subnetsInfo)
	return nil
}

// dataSourceIBMISSubnetsId returns a reasonable ID for a subnet list.
func dataSourceIBMISSubnetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
