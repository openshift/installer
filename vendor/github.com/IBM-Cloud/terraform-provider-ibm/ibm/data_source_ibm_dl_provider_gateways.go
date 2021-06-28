// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	dlProviderV2 "github.com/IBM/networking-go-sdk/directlinkproviderv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

const (
	dlProviderGateways   = "gateways"
	dlProviderGatewaysID = "id"
)

func dataSourceIBMDirectLinkProviderGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDirectLinkProviderGatewaysRead,

		Schema: map[string]*schema.Schema{

			dlProviderGateways: {

				Type:        schema.TypeList,
				Description: "Collection of direct link provider ports",
				Computed:    true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{
						dlProviderGatewaysID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the data source gateways",
						},
						dlBgpAsn: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP ASN",
						},

						dlBgpCerCidr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP customer edge router CIDR",
						},
						dlBgpIbmAsn: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IBM BGP ASN",
						},
						dlBgpIbmCidr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP IBM CIDR",
						},
						dlBgpStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway BGP status",
						},
						customerAccountID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Customer IBM Cloud account ID for the new gateway. A gateway object containing the pending create request will become available in the specified account.",
						},
						dlCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time resource was created",
						},
						dlCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN (Cloud Resource Name) of this gateway",
						},

						dlName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this gateway",
						},
						dlOperationalStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway operational status",
						},
						dlChangeRequest: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changes pending approval for provider managed Direct Link gateways",
						},
						dlPort: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway port",
						},
						dlProviderAPIManaged: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether gateway was created through a provider portal",
						},
						dlResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway resource group",
						},
						dlSpeedMbps: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Gateway speed in megabits per second",
						},
						dlType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway type",
						},
						dlVlan: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VLAN allocated for this gateway",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDirectLinkProviderGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	directLinkProvider, err := directlinkProviderClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []dlProviderV2.ProviderGateway{}
	for {
		listProviderGatewaysOptions := directLinkProvider.NewListProviderGatewaysOptions()
		if start != "" {
			listProviderGatewaysOptions.Start = &start
		}

		providerGateways, resp, err := directLinkProvider.ListProviderGateways(listProviderGatewaysOptions)
		if err != nil {
			log.Println("[WARN] Error listing dl provider gateways", providerGateways, resp, err)
			return err
		}
		start = GetNext(providerGateways.Next)
		allrecs = append(allrecs, providerGateways.Gateways...)
		if start == "" {
			break
		}
	}
	gatewayCollections := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		gatewayCollection := map[string]interface{}{}

		if instance.ID != nil {
			gatewayCollection[dlProviderGatewaysID] = *instance.ID
		}
		if instance.Name != nil {
			gatewayCollection[dlName] = *instance.Name
		}
		if instance.Crn != nil {
			gatewayCollection[dlCrn] = *instance.Crn
		}
		if instance.BgpAsn != nil {
			gatewayCollection[dlBgpAsn] = *instance.BgpAsn
		}
		if instance.BgpIbmCidr != nil {
			gatewayCollection[dlBgpIbmCidr] = *instance.BgpIbmCidr
		}
		if instance.BgpIbmAsn != nil {
			gatewayCollection[dlBgpIbmAsn] = *instance.BgpIbmAsn
		}

		if instance.BgpCerCidr != nil {
			gatewayCollection[dlBgpCerCidr] = *instance.BgpCerCidr
		}

		if instance.ProviderApiManaged != nil {
			gatewayCollection[dlProviderAPIManaged] = *instance.ProviderApiManaged
		}
		if instance.Type != nil {
			gatewayCollection[dlType] = *instance.Type
		}
		if instance.SpeedMbps != nil {
			gatewayCollection[dlSpeedMbps] = *instance.SpeedMbps
		}
		if instance.OperationalStatus != nil {
			gatewayCollection[dlOperationalStatus] = *instance.OperationalStatus
		}
		if instance.BgpStatus != nil {
			gatewayCollection[dlBgpStatus] = *instance.BgpStatus
		}
		if instance.Vlan != nil {
			gatewayCollection[dlVlan] = *instance.Vlan
		}

		if instance.Port != nil {
			gatewayCollection[dlPort] = *instance.Port.ID
		}

		if instance.CreatedAt != nil {
			gatewayCollection[dlCreatedAt] = instance.CreatedAt.String()
		}

		gatewayCollections = append(gatewayCollections, gatewayCollection)
	}
	d.SetId(dataSourceIBMDirectLinkProviderGatewaysReadID(d))
	d.Set(dlProviderGateways, gatewayCollections)
	return nil
}

func dataSourceIBMDirectLinkProviderGatewaysReadID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
