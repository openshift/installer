// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isPublicGateways = "public_gateways"
)

func dataSourceIBMISPublicGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISPublicGatewaysRead,

		Schema: map[string]*schema.Schema{
			isPublicGateways: {
				Type:        schema.TypeList,
				Description: "List of public gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public gateway id",
						},
						isPublicGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public gateway Name",
						},

						isPublicGatewayFloatingIP: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Public gateway floating IP",
						},

						isPublicGatewayStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public gateway instance status",
						},

						isPublicGatewayResourceGroup: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Public gateway resource group info",
						},

						isPublicGatewayVPC: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public gateway VPC info",
						},

						isPublicGatewayZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public gateway zone info",
						},

						isPublicGatewayTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         resourceIBMVPCHash,
							Description: "Service tags for the public gateway instance",
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
					},
				},
			},
		},
	}
}

func dataSourceIBMISPublicGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	err := publicGatewaysGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func publicGatewaysGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	rgroup := ""
	if rg, ok := d.GetOk(isPublicGatewayResourceGroup); ok {
		rgroup = rg.(string)
	}
	start := ""
	allrecs := []vpcv1.PublicGateway{}
	for {
		listPublicGatewaysOptions := &vpcv1.ListPublicGatewaysOptions{}
		if start != "" {
			listPublicGatewaysOptions.Start = &start
		}
		if rgroup != "" {
			listPublicGatewaysOptions.ResourceGroupID = &rgroup
		}
		publicgws, response, err := sess.ListPublicGateways(listPublicGatewaysOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching public gateways %s\n%s", err, response)
		}
		start = GetNext(publicgws.Next)
		allrecs = append(allrecs, publicgws.PublicGateways...)
		if start == "" {
			break
		}
	}
	publicgwInfo := make([]map[string]interface{}, 0)
	for _, publicgw := range allrecs {
		id := *publicgw.ID
		l := map[string]interface{}{
			"id":                  id,
			isPublicGatewayName:   *publicgw.Name,
			isPublicGatewayStatus: *publicgw.Status,
			isPublicGatewayZone:   *publicgw.Zone.Name,
			isPublicGatewayVPC:    *publicgw.VPC.ID,

			ResourceName:   *publicgw.Name,
			ResourceCRN:    *publicgw.CRN,
			ResourceStatus: *publicgw.Status,
		}
		if publicgw.FloatingIP != nil {
			floatIP := map[string]interface{}{
				"id":                             *publicgw.FloatingIP.ID,
				isPublicGatewayFloatingIPAddress: *publicgw.FloatingIP.Address,
			}
			l[isPublicGatewayFloatingIP] = floatIP
		}
		tags, err := GetTagsUsingCRN(meta, *publicgw.CRN)
		if err != nil {
			log.Printf(
				"Error on get of vpc public gateway (%s) tags: %s", *publicgw.ID, err)
		}
		l[isPublicGatewayTags] = tags
		controller, err := getBaseController(meta)
		if err != nil {
			return err
		}
		l[ResourceControllerURL] = controller + "/vpc-ext/network/publicGateways"
		if publicgw.ResourceGroup != nil {
			l[isPublicGatewayResourceGroup] = *publicgw.ResourceGroup.ID
			l[ResourceGroupName] = *publicgw.ResourceGroup.Name
		}
		publicgwInfo = append(publicgwInfo, l)
	}
	d.SetId(dataSourceIBMISPublicGatewaysID(d))
	d.Set(isPublicGateways, publicgwInfo)
	return nil
}

// dataSourceIBMISPublicGatewaysID returns a reasonable ID for a Public Gateway list.
func dataSourceIBMISPublicGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
