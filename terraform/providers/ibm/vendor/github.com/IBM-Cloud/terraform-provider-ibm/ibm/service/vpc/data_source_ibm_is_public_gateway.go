// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISPublicGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISPublicGatewayRead,

		Schema: map[string]*schema.Schema{
			isPublicGatewayName: {
				Type:        schema.TypeString,
				Required:    true,
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
				Set:         flex.ResourceIBMVPCHash,
				Description: "Service tags for the public gateway instance",
			},

			isPublicGatewayAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
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

			isPublicGatewayCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
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
		},
	}
}

func dataSourceIBMISPublicGatewayRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	name := d.Get(isPublicGatewayName).(string)
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
			return fmt.Errorf("[ERROR] Error Fetching public gateways %s\n%s", err, response)
		}
		start = flex.GetNext(publicgws.Next)
		allrecs = append(allrecs, publicgws.PublicGateways...)
		if start == "" {
			break
		}
	}
	for _, publicgw := range allrecs {
		if *publicgw.Name == name {
			d.SetId(*publicgw.ID)
			d.Set(isPublicGatewayName, *publicgw.Name)
			if publicgw.FloatingIP != nil {
				floatIP := map[string]interface{}{
					"id":                             *publicgw.FloatingIP.ID,
					isPublicGatewayFloatingIPAddress: *publicgw.FloatingIP.Address,
				}
				d.Set(isPublicGatewayFloatingIP, floatIP)

			}
			d.Set(isPublicGatewayStatus, *publicgw.Status)
			d.Set(isPublicGatewayZone, *publicgw.Zone.Name)
			d.Set(isPublicGatewayVPC, *publicgw.VPC.ID)
			tags, err := flex.GetGlobalTagsUsingCRN(meta, *publicgw.CRN, "", isUserTagType)
			if err != nil {
				log.Printf(
					"Error on get of vpc public gateway (%s) tags: %s", *publicgw.ID, err)
			}
			d.Set(isPublicGatewayTags, tags)

			accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *publicgw.CRN, "", isAccessTagType)
			if err != nil {
				log.Printf(
					"Error on get of vpc public gateway (%s) access tags: %s", d.Id(), err)
			}

			d.Set(isPublicGatewayAccessTags, accesstags)

			controller, err := flex.GetBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/publicGateways")
			d.Set(flex.ResourceName, *publicgw.Name)
			d.Set(flex.ResourceCRN, *publicgw.CRN)
			d.Set(isPublicGatewayCRN, *publicgw.CRN)
			d.Set(flex.ResourceStatus, *publicgw.Status)
			if publicgw.ResourceGroup != nil {
				d.Set(isPublicGatewayResourceGroup, *publicgw.ResourceGroup.ID)
				d.Set(flex.ResourceGroupName, *publicgw.ResourceGroup.Name)
			}
			return nil
		}
	}
	return fmt.Errorf("[ERROR] No Public gateway found with name %s", name)
}
