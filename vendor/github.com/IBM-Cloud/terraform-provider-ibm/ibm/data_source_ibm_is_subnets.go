// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSubnets                = "subnets"
	isSubnetResourceGroupID  = "resource_group"
	isSubnetRoutingTableName = "routing_table_name"
)

func dataSourceIBMISSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSubnetsRead,

		Schema: map[string]*schema.Schema{
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

	for {
		if start != "" {
			options.Start = &start
		}
		subnets, response, err := sess.ListSubnets(options)
		if err != nil {
			return fmt.Errorf("Error Fetching subnets %s\n%s", err, response)
		}
		start = GetNext(subnets.Next)
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
