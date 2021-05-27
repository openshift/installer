// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	tgConnName    = "name"
	tgConnections = "connections"
)

func dataSourceIBMTransitGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayRead,

		Schema: map[string]*schema.Schema{
			tgName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Transit Gateway identifier",
				ValidateFunc: InvokeValidator("ibm_tg_gateway", tgName),
			},
			tgCrn: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgLocation: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgGlobal: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			tgStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgUpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgResourceGroup: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgConnections: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateway connections",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgNetworkAccountID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the account which owns the network that is being connected. Generally only used if the network is in a different account than the gateway.",
						},
						tgNetworkId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConnName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgNetworkType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConectionCreatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConnectionStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgUpdatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	listTransitGatewaysOptionsModel := &transitgatewayapisv1.ListTransitGatewaysOptions{}
	listTransitGateways, response, err := client.ListTransitGateways(listTransitGatewaysOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while listing transit gateways %s\n%s", err, response)
	}

	gwName := d.Get(tgName).(string)
	var foundGateway bool
	for _, tgw := range listTransitGateways.TransitGateways {

		if *tgw.Name == gwName {
			d.SetId(*tgw.ID)
			d.Set(tgCrn, tgw.Crn)
			d.Set(tgName, tgw.Name)
			d.Set(tgLocation, tgw.Location)
			d.Set(tgCreatedAt, tgw.CreatedAt.String())

			if tgw.UpdatedAt != nil {
				d.Set(tgUpdatedAt, tgw.UpdatedAt.String())
			}
			d.Set(tgGlobal, tgw.Global)
			d.Set(tgStatus, tgw.Status)

			if tgw.ResourceGroup != nil {
				rg := tgw.ResourceGroup
				d.Set(tgResourceGroup, *rg.ID)
			}
			foundGateway = true
		}
	}

	if !foundGateway {
		return fmt.Errorf(
			"Couldn't find any gateway with the specified name: (%s)", gwName)
	}

	return dataSourceIBMTransitGatewayConnectionsRead(d, meta)

}
func dataSourceIBMTransitGatewayConnectionsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	listTransitGatewayConnectionsOptions := &transitgatewayapisv1.ListTransitGatewayConnectionsOptions{}
	tgGatewayId := d.Id()
	log.Println("tgGatewayId: ", tgGatewayId)

	listTransitGatewayConnectionsOptions.SetTransitGatewayID(tgGatewayId)
	listTGConnections, response, err := client.ListTransitGatewayConnections(listTransitGatewayConnectionsOptions)
	if err != nil {
		return fmt.Errorf("Error while listing transit gateway connections %s\n%s", err, response)
	}
	connections := make([]map[string]interface{}, 0)

	for _, instance := range listTGConnections.Connections {
		tgConn := map[string]interface{}{}

		if instance.ID != nil {
			tgConn[ID] = *instance.ID
		}
		if instance.Name != nil {
			tgConn[tgConnName] = *instance.Name
		}
		if instance.NetworkType != nil {
			tgConn[tgNetworkType] = *instance.NetworkType
		}

		if instance.NetworkID != nil {
			tgConn[tgNetworkId] = *instance.NetworkID
		}
		if instance.NetworkAccountID != nil {
			tgConn[tgNetworkAccountID] = *instance.NetworkAccountID
		}

		if instance.CreatedAt != nil {
			tgConn[tgConectionCreatedAt] = instance.CreatedAt.String()

		}
		if instance.UpdatedAt != nil {
			tgConn[tgUpdatedAt] = instance.UpdatedAt.String()

		}
		if instance.Status != nil {
			tgConn[tgConnectionStatus] = *instance.Status
		}

		connections = append(connections, tgConn)

	}
	d.Set(tgConnections, connections)

	return nil

}
