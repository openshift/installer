// Copyright IBM Corp. 2017, 2021, 2024  All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	tgConnName        = "name"
	tgConnections     = "connections"
	ID                = "id"
	tgBaseNetworkType = "base_network_type"
)

func DataSourceIBMTransitGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayRead,

		Schema: map[string]*schema.Schema{
			tgName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Transit Gateway identifier",
				ValidateFunc: validate.InvokeValidator("ibm_tg_gateway", tgName),
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
						tgBaseConnectionId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgBaseNetworkType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgLocalBgpAsn: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						tgLocalGatewayIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgLocalTunnelIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgRemoteBgpAsn: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						tgRemoteGatewayIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgRemoteTunnelIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgZone: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgMtu: {
							Type:     schema.TypeInt,
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
						tgDefaultPrefixFilter: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_tg_connection_prefix_filter", tgAction),
							Description:  "Whether to permit or deny the prefix filter",
						},
						tgrGREtunnels: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of GRE tunnels for a transit gateway redundant GRE tunnel connection. This field is required for 'redundant_gre' connections",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									tgconTunnelName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this  rGRE tunnel.",
									},
									tgLocalGatewayIp: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The local gateway IP address.",
									},
									tgLocalTunnelIp: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The local tunnel IP address.",
									},
									tgRemoteGatewayIp: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The remote gateway IP address.",
									},
									tgRemoteTunnelIp: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The remote tunnel IP address.",
									},
									tgZone: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Location of GRE tunnel.",
									},
									tgCreatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that this connection was created",
									},
									tgUpdatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that this connection was last updated",
									},
									tgGreTunnelStatus: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "What is the current configuration state of this connection. Possible values: [attached,failed,pending,deleting,detaching,detached]",
									},
									tgGreTunnelId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Transit Gateway Connection identifier",
									},
									tgMtu: {
										Type:     schema.TypeInt,
										Computed: true,
									},
									tgLocalBgpAsn: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The local network BGP ASN. ",
									},
									tgRemoteBgpAsn: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The remote network BGP ASN.",
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

func dataSourceIBMTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	listTransitGatewaysOptionsModel := &transitgatewayapisv1.ListTransitGatewaysOptions{}
	listTransitGateways, response, err := client.ListTransitGateways(listTransitGatewaysOptionsModel)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while listing transit gateways %s\n%s", err, response)
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
		return fmt.Errorf("[ERROR] Couldn't find any gateway with the specified name: (%s)", gwName)
	}

	return dataSourceIBMTransitGatewayConnectionsRead(d, meta)

}
func dataSourceIBMTransitGatewayConnectionsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	startSub := ""
	listTransitGatewayConnectionsOptions := &transitgatewayapisv1.ListTransitGatewayConnectionsOptions{}
	tgGatewayId := d.Id()
	log.Println("tgGatewayId: ", tgGatewayId)

	listTransitGatewayConnectionsOptions.SetTransitGatewayID(tgGatewayId)
	connections := make([]map[string]interface{}, 0)
	for {

		if startSub != "" {
			listTransitGatewayConnectionsOptions.Start = &startSub
		}
		listTGConnections, response, err := client.ListTransitGatewayConnections(listTransitGatewayConnectionsOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while listing transit gateway connections %s\n%s", err, response)
		}
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
			if instance.BaseNetworkType != nil {
				tgConn[tgBaseNetworkType] = *instance.BaseNetworkType
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

			if *instance.NetworkType != "redundant_gre" {

				if instance.BaseConnectionID != nil {
					tgConn[tgBaseConnectionId] = *instance.BaseConnectionID
				}
				if instance.LocalBgpAsn != nil {
					tgConn[tgLocalBgpAsn] = *instance.LocalBgpAsn
				}
				if instance.LocalGatewayIp != nil {
					tgConn[tgLocalGatewayIp] = *instance.LocalGatewayIp
				}
				if instance.LocalTunnelIp != nil {
					tgConn[tgLocalTunnelIp] = *instance.LocalTunnelIp
				}
				if instance.RemoteBgpAsn != nil {
					tgConn[tgRemoteBgpAsn] = *instance.RemoteBgpAsn
				}
				if instance.RemoteGatewayIp != nil {
					tgConn[tgRemoteGatewayIp] = *instance.RemoteGatewayIp
				}
				if instance.RemoteTunnelIp != nil {
					tgConn[tgRemoteTunnelIp] = *instance.RemoteTunnelIp
				}
				if instance.Zone != nil {
					tgConn[tgZone] = *instance.Zone.Name
				}
				if instance.Mtu != nil {
					tgConn[tgMtu] = *instance.Mtu
				}
				if instance.PrefixFiltersDefault != nil {
					tgConn[tgDefaultPrefixFilter] = *instance.PrefixFiltersDefault
				}
			}

			if instance.Tunnels != nil {
				// read the tunnels
				rGREtunnels := make([]map[string]interface{}, 0)
				for _, rGREtunnel := range instance.Tunnels {

					tunnel := map[string]interface{}{}
					if rGREtunnel.ID != nil {
						tunnel[tgGreTunnelId] = *rGREtunnel.ID
					}
					if rGREtunnel.LocalGatewayIp != nil {
						tunnel[tgLocalGatewayIp] = *rGREtunnel.LocalGatewayIp
					}
					if rGREtunnel.LocalTunnelIp != nil {
						tunnel[tgLocalTunnelIp] = *rGREtunnel.LocalTunnelIp
					}
					if rGREtunnel.RemoteGatewayIp != nil {
						tunnel[tgRemoteGatewayIp] = *rGREtunnel.RemoteGatewayIp
					}
					if rGREtunnel.RemoteTunnelIp != nil {
						tunnel[tgRemoteTunnelIp] = *rGREtunnel.RemoteTunnelIp
					}
					if rGREtunnel.Mtu != nil {
						tunnel[tgMtu] = *rGREtunnel.Mtu
					}
					if rGREtunnel.RemoteBgpAsn != nil {
						tunnel[tgRemoteBgpAsn] = *rGREtunnel.RemoteBgpAsn
					}
					if rGREtunnel.Name != nil {
						tunnel[tgconTunnelName] = *rGREtunnel.Name
					}
					if rGREtunnel.Zone.Name != nil {
						tunnel[tgZone] = *rGREtunnel.Zone.Name
					}
					if rGREtunnel.LocalBgpAsn != nil {
						tunnel[tgLocalBgpAsn] = *rGREtunnel.LocalBgpAsn
					}
					if rGREtunnel.Status != nil {
						tunnel[tgGreTunnelStatus] = *rGREtunnel.Status
					}
					if rGREtunnel.CreatedAt != nil {
						tunnel[tgCreatedAt] = rGREtunnel.CreatedAt.String()
					}
					if rGREtunnel.UpdatedAt != nil {
						tunnel[tgUpdatedAt] = rGREtunnel.UpdatedAt.String()
					}
					rGREtunnels = append(rGREtunnels, tunnel)
				}
				if len(rGREtunnels) > 0 {
					tgConn[tgrGREtunnels] = rGREtunnels
				}
			}
			connections = append(connections, tgConn)
		}
		startSub = flex.GetNext(listTGConnections.Next)
		if startSub == "" {
			break
		}
	}
	d.Set(tgConnections, connections)
	return nil

}
