// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISPublicGateway() *schema.Resource {
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
	}
}

func dataSourceIBMISPublicGatewayRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	name := d.Get(isPublicGatewayName).(string)
	if userDetails.generation == 1 {
		err := classicPublicGatewayGet(d, meta, name)
		if err != nil {
			return err
		}
	} else {
		err := publicGatewayGet(d, meta, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicPublicGatewayGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.PublicGateway{}
	for {
		listPublicGatewaysOptions := &vpcclassicv1.ListPublicGatewaysOptions{}
		if start != "" {
			listPublicGatewaysOptions.Start = &start
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
			tags, err := GetTagsUsingCRN(meta, *publicgw.CRN)
			if err != nil {
				log.Printf(
					"Error on get of vpc public gateway (%s) tags: %s", *publicgw.ID, err)
			}
			d.Set(isPublicGatewayTags, tags)
			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/network/publicGateways")
			d.Set(ResourceName, *publicgw.Name)
			d.Set(ResourceCRN, *publicgw.CRN)
			d.Set(ResourceStatus, *publicgw.Status)
			return nil
		}
	}
	return fmt.Errorf("No Public Gateway found with name %s", name)
}

func publicGatewayGet(d *schema.ResourceData, meta interface{}, name string) error {
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
			tags, err := GetTagsUsingCRN(meta, *publicgw.CRN)
			if err != nil {
				log.Printf(
					"Error on get of vpc public gateway (%s) tags: %s", *publicgw.ID, err)
			}
			d.Set(isPublicGatewayTags, tags)
			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc-ext/network/publicGateways")
			d.Set(ResourceName, *publicgw.Name)
			d.Set(ResourceCRN, *publicgw.CRN)
			d.Set(ResourceStatus, *publicgw.Status)
			if publicgw.ResourceGroup != nil {
				d.Set(isPublicGatewayResourceGroup, *publicgw.ResourceGroup.ID)
				d.Set(ResourceGroupName, *publicgw.ResourceGroup.Name)
			}
			return nil
		}
	}
	return fmt.Errorf("No Public gateway found with name %s", name)
}
