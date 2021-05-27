// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSubnets = "subnets"
)

func dataSourceIBMISSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSubnetsRead,

		Schema: map[string]*schema.Schema{

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
						"ipv6_cidr_block": {
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
						"resource_group": {
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
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicSubnetList(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := subnetList(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicSubnetList(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.Subnet{}
	for {
		options := &vpcclassicv1.ListSubnetsOptions{}
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
		subnetsInfo = append(subnetsInfo, l)
	}
	d.SetId(dataSourceIBMISSubnetsID(d))
	d.Set(isSubnets, subnetsInfo)
	return nil
}

func subnetList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.Subnet{}
	for {
		options := &vpcv1.ListSubnetsOptions{}
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
