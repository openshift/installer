// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISVPNGatewayConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayConnectionRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_name", "vpn_gateway"},
				Description:  "The VPN gateway identifier.",
			},
			"vpn_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_name", "vpn_gateway"},
				Description:  "The VPN gateway name.",
			},
			"vpn_gateway_connection": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_connection", "vpn_gateway_connection_name"},
				Description:  "The VPN gateway connection identifier.",
			},
			"vpn_gateway_connection_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_connection", "vpn_gateway_connection_name"},
				Description:  "The VPN gateway connection name.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to false, the VPN gateway connection is shut down.",
			},
			"authentication_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication mode. Only `psk` is currently supported.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this VPN gateway connection was created.",
			},
			"dead_peer_detection": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Dead Peer Detection settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dead Peer Detection actions.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dead Peer Detection interval in seconds.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dead Peer Detection timeout in seconds. Must be at least the interval.",
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPN connection's canonical URL.",
			},
			"ike_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IKE policy. If absent, [auto-negotiation isused](https://cloud.ibm.com/docs/vpc?topic=vpc-using-vpn&interface=ui#ike-auto-negotiation-phase-1).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
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
							Description: "The IKE policy's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this IKE policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this IKE policy.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"ipsec_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IPsec policy. If absent, [auto-negotiation isused](https://cloud.ibm.com/docs/vpc?topic=vpc-using-vpn&interface=ui#ipsec-auto-negotiation-phase-2).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
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
							Description: "The IPsec policy's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this IPsec policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this IPsec policy.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mode of the VPN gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this VPN gateway connection.",
			},

			// new breaking changes
			"establish_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The establish mode of the VPN gateway connection:- `bidirectional`: Either side of the VPN gateway can initiate IKE protocol   negotiations or rekeying processes.- `peer_only`: Only the peer can initiate IKE protocol negotiations for this VPN gateway   connection. Additionally, the peer is responsible for initiating the rekeying process   after the connection is established. If rekeying does not occur, the VPN gateway   connection will be brought down after its lifetime expires.",
			},
			"local": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_identities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The local IKE identities.A VPN gateway in static route mode consists of two members in active-active mode. The first identity applies to the first member, and the second identity applies to the second member.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IKE identity type.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the backup policy on which the unexpected property value was encountered.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IKE identity FQDN value.",
									},
								},
							},
						},
						"cidrs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The local CIDRs for this resource.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"peer": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_identity": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The peer IKE identity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IKE identity type.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the backup policy on which the unexpected property value was encountered.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IKE identity FQDN value.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates whether `peer.address` or `peer.fqdn` is used.",
						},
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the peer VPN gateway for this connection.",
						},
						"fqdn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The FQDN of the peer VPN gateway for this connection.",
						},
						"cidrs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The peer CIDRs for this resource.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"peer_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the peer VPN gateway.",
				Deprecated:  "peer_address is deprecated, use peer instead",
			},
			"psk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The preshared key.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of a VPN gateway connection.",
			},
			"status_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"routing_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing protocols are disabled for this VPN gateway connection.",
			},
			"tunnels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPN tunnel configuration for this VPN gateway connection (in static route mode).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the VPN gateway member in which the tunnel resides.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN Tunnel.",
						},
					},
				},
			},
			"local_cidrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The local CIDRs for this resource.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated: "local_cidrs is deprecated, use local instead",
			},
			"peer_cidrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The peer CIDRs for this resource.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated: "peer_cidrs is deprecated, use peer instead",
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayConnectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	vpn_gateway_id := d.Get("vpn_gateway").(string)
	vpn_gateway_name := d.Get("vpn_gateway_name").(string)
	vpn_gateway_connection := d.Get("vpn_gateway_connection").(string)
	vpn_gateway_connection_name := d.Get("vpn_gateway_connection_name").(string)

	var vpnGatewayConnection vpcv1.VPNGatewayConnectionIntf

	if vpn_gateway_name != "" {
		listvpnGWOptions := vpcClient.NewListVPNGatewaysOptions()

		start := ""
		allrecs := []vpcv1.VPNGatewayIntf{}
		for {
			if start != "" {
				listvpnGWOptions.Start = &start
			}
			availableVPNGateways, detail, err := vpcClient.ListVPNGatewaysWithContext(context, listvpnGWOptions)
			if err != nil || availableVPNGateways == nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error reading list of VPN Gateways:%s\n%s", err, detail))
			}
			start = flex.GetNext(availableVPNGateways.Next)
			allrecs = append(allrecs, availableVPNGateways.VPNGateways...)
			if start == "" {
				break
			}
		}
		vpn_gateway_found := false
		for _, vpnGatewayIntfItem := range allrecs {
			if *vpnGatewayIntfItem.(*vpcv1.VPNGateway).Name == vpn_gateway_name {
				vpnGateway := vpnGatewayIntfItem.(*vpcv1.VPNGateway)
				vpn_gateway_id = *vpnGateway.ID
				vpn_gateway_found = true
				break
			}
		}
		if !vpn_gateway_found {
			log.Printf("[DEBUG] No vpn gateway and connection found with given name %s", vpn_gateway_name)
			return diag.FromErr(fmt.Errorf("No vpn gateway and connection found with given name %s", vpn_gateway_name))
		}
	}

	if vpn_gateway_connection_name != "" {
		listvpnGWConnectionOptions := vpcClient.NewListVPNGatewayConnectionsOptions(vpn_gateway_id)

		availableVPNGatewayConnections, detail, err := vpcClient.ListVPNGatewayConnections(listvpnGWConnectionOptions)
		if err != nil || availableVPNGatewayConnections == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error reading list of VPN Gateway Connections:%s\n%s", err, detail))
		}

		vpn_gateway_conn_found := false
		for _, connectionItem := range availableVPNGatewayConnections.Connections {
			switch reflect.TypeOf(connectionItem).String() {
			case "*vpcv1.VPNGatewayConnection":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnection)
					if *connection.Name == vpn_gateway_connection_name {
						vpnGatewayConnection = connectionItem
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionRouteMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionRouteMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpnGatewayConnection = connectionItem
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpnGatewayConnection = connectionItem
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionPolicyMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionPolicyMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpnGatewayConnection = connectionItem
						vpn_gateway_conn_found = true
						break
					}
				}
			}

		}
		if !vpn_gateway_conn_found {
			return diag.FromErr(fmt.Errorf("VPN gateway connection %s not found", vpn_gateway_connection_name))
		}
	} else if vpn_gateway_connection != "" {
		getVPNGatewayConnectionOptions := &vpcv1.GetVPNGatewayConnectionOptions{}

		getVPNGatewayConnectionOptions.SetVPNGatewayID(vpn_gateway_id)
		getVPNGatewayConnectionOptions.SetID(vpn_gateway_connection)

		vpnGatewayConnectionIntf, response, err := vpcClient.GetVPNGatewayConnectionWithContext(context, getVPNGatewayConnectionOptions)
		if err != nil || vpnGatewayConnectionIntf == nil {
			log.Printf("[DEBUG] GetVPNGatewayConnectionWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetVPNGatewayConnectionWithContext failed %s\n%s", err, response))
		}
		vpnGatewayConnection = vpnGatewayConnectionIntf
	}

	setvpnGatewayConnectionIntfDatasourceData(d, vpn_gateway_id, vpnGatewayConnection)
	return nil
}

func setvpnGatewayConnectionIntfDatasourceData(d *schema.ResourceData, vpn_gateway_id string, vpnGatewayConnectionIntf vpcv1.VPNGatewayConnectionIntf) error {
	var err error
	switch reflect.TypeOf(vpnGatewayConnectionIntf).String() {
	case "*vpcv1.VPNGatewayConnection":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection)
			d.SetId(fmt.Sprintf("%s/%s", vpn_gateway_id, *vpnGatewayConnection.ID))
			if err = d.Set("admin_state_up", vpnGatewayConnection.AdminStateUp); err != nil {
				return fmt.Errorf("[ERROR] Error setting admin_state_up: %s", err)
			}
			if err = d.Set("authentication_mode", vpnGatewayConnection.AuthenticationMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting authentication_mode: %s", err)
			}
			if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayConnection.CreatedAt)); err != nil {
				return fmt.Errorf("[ERROR] Error setting created_at: %s", err)
			}

			if vpnGatewayConnection.DeadPeerDetection != nil {
				err = d.Set("dead_peer_detection", dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(*vpnGatewayConnection.DeadPeerDetection))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting dead_peer_detection %s", err)
				}
			}
			if err = d.Set("href", vpnGatewayConnection.Href); err != nil {
				return fmt.Errorf("[ERROR] Error setting href: %s", err)
			}

			if vpnGatewayConnection.IkePolicy != nil {
				err = d.Set("ike_policy", dataSourceVPNGatewayConnectionFlattenIkePolicy(*vpnGatewayConnection.IkePolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ike_policy %s", err)
				}
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				err = d.Set("ipsec_policy", dataSourceVPNGatewayConnectionFlattenIpsecPolicy(*vpnGatewayConnection.IpsecPolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ipsec_policy %s", err)
				}
			}
			if err = d.Set("mode", vpnGatewayConnection.Mode); err != nil {
				return fmt.Errorf("[ERROR] Error setting mode: %s", err)
			}
			if err = d.Set("name", vpnGatewayConnection.Name); err != nil {
				return fmt.Errorf("[ERROR] Error setting name: %s", err)
			}

			// breaking changes
			if err = d.Set("establish_mode", vpnGatewayConnection.EstablishMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting establish_mode: %s", err)
			}
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return err
				}
				local = append(local, modelMap)
			}
			if err = d.Set("local", local); err != nil {
				return fmt.Errorf("[ERROR] Error setting local %s", err)
			}

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return err
				}
				peer = append(peer, modelMap)
			}
			if err = d.Set("peer", peer); err != nil {
				return fmt.Errorf("[ERROR] Error setting peer %s", err)
			}
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				if err = d.Set("peer_address", peer.Address); err != nil {
					return fmt.Errorf("[ERROR] Error setting peer_address: %s", err)
				}
			}
			if err = d.Set("psk", vpnGatewayConnection.Psk); err != nil {
				return fmt.Errorf("[ERROR] Error setting psk: %s", err)
			}
			if err = d.Set("resource_type", vpnGatewayConnection.ResourceType); err != nil {
				return fmt.Errorf("[ERROR] Error setting resource_type: %s", err)
			}
			if err = d.Set("status", vpnGatewayConnection.Status); err != nil {
				return fmt.Errorf("[ERROR] Error setting status: %s", err)
			}
			if err := d.Set("status_reasons", resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)); err != nil {
				return fmt.Errorf("[ERROR] Error setting status_reasons: %s", err)
			}
			if err = d.Set("routing_protocol", vpnGatewayConnection.RoutingProtocol); err != nil {
				return fmt.Errorf("[ERROR] Error setting routing_protocol: %s", err)
			}

			if vpnGatewayConnection.Tunnels != nil {
				err = d.Set("tunnels", dataSourceVPNGatewayConnectionFlattenTunnels(vpnGatewayConnection.Tunnels))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting tunnels %s", err)
				}
			}
		}
	case "*vpcv1.VPNGatewayConnectionRouteMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionRouteMode)
			d.SetId(fmt.Sprintf("%s/%s", vpn_gateway_id, *vpnGatewayConnection.ID))
			if err = d.Set("admin_state_up", vpnGatewayConnection.AdminStateUp); err != nil {
				return fmt.Errorf("[ERROR] Error setting admin_state_up: %s", err)
			}
			if err = d.Set("authentication_mode", vpnGatewayConnection.AuthenticationMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting authentication_mode: %s", err)
			}
			if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayConnection.CreatedAt)); err != nil {
				return fmt.Errorf("[ERROR] Error setting created_at: %s", err)
			}

			if vpnGatewayConnection.DeadPeerDetection != nil {
				err = d.Set("dead_peer_detection", dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(*vpnGatewayConnection.DeadPeerDetection))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting dead_peer_detection %s", err)
				}
			}
			if err = d.Set("href", vpnGatewayConnection.Href); err != nil {
				return fmt.Errorf("[ERROR] Error setting href: %s", err)
			}

			if vpnGatewayConnection.IkePolicy != nil {
				err = d.Set("ike_policy", dataSourceVPNGatewayConnectionFlattenIkePolicy(*vpnGatewayConnection.IkePolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ike_policy %s", err)
				}
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				err = d.Set("ipsec_policy", dataSourceVPNGatewayConnectionFlattenIpsecPolicy(*vpnGatewayConnection.IpsecPolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ipsec_policy %s", err)
				}
			}
			if err = d.Set("mode", vpnGatewayConnection.Mode); err != nil {
				return fmt.Errorf("[ERROR] Error setting mode: %s", err)
			}
			if err = d.Set("name", vpnGatewayConnection.Name); err != nil {
				return fmt.Errorf("[ERROR] Error setting name: %s", err)
			}

			// breaking changes
			if err = d.Set("establish_mode", vpnGatewayConnection.EstablishMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting establish_mode: %s", err)
			}
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return err
				}
				local = append(local, modelMap)
			}
			if err = d.Set("local", local); err != nil {
				return fmt.Errorf("[ERROR] Error setting local %s", err)
			}

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return err
				}
				peer = append(peer, modelMap)
			}
			if err = d.Set("peer", peer); err != nil {
				return fmt.Errorf("[ERROR] Error setting peer %s", err)
			}
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				if err = d.Set("peer_address", peer.Address); err != nil {
					return fmt.Errorf("[ERROR] Error setting peer_address: %s", err)
				}
			}
			if err = d.Set("psk", vpnGatewayConnection.Psk); err != nil {
				return fmt.Errorf("[ERROR] Error setting psk: %s", err)
			}
			if err = d.Set("resource_type", vpnGatewayConnection.ResourceType); err != nil {
				return fmt.Errorf("[ERROR] Error setting resource_type: %s", err)
			}
			if err = d.Set("status", vpnGatewayConnection.Status); err != nil {
				return fmt.Errorf("[ERROR] Error setting status: %s", err)
			}
			if err := d.Set("status_reasons", resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)); err != nil {
				return fmt.Errorf("[ERROR] Error setting status_reasons: %s", err)
			}
			if err = d.Set("routing_protocol", vpnGatewayConnection.RoutingProtocol); err != nil {
				return fmt.Errorf("[ERROR] Error setting routing_protocol: %s", err)
			}

			if vpnGatewayConnection.Tunnels != nil {
				err = d.Set("tunnels", dataSourceVPNGatewayConnectionFlattenTunnels(vpnGatewayConnection.Tunnels))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting tunnels %s", err)
				}
			}
		}
	case "*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode)
			d.SetId(fmt.Sprintf("%s/%s", vpn_gateway_id, *vpnGatewayConnection.ID))
			if err = d.Set("admin_state_up", vpnGatewayConnection.AdminStateUp); err != nil {
				return fmt.Errorf("[ERROR] Error setting admin_state_up: %s", err)
			}
			if err = d.Set("authentication_mode", vpnGatewayConnection.AuthenticationMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting authentication_mode: %s", err)
			}
			if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayConnection.CreatedAt)); err != nil {
				return fmt.Errorf("[ERROR] Error setting created_at: %s", err)
			}

			if vpnGatewayConnection.DeadPeerDetection != nil {
				err = d.Set("dead_peer_detection", dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(*vpnGatewayConnection.DeadPeerDetection))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting dead_peer_detection %s", err)
				}
			}
			if err = d.Set("href", vpnGatewayConnection.Href); err != nil {
				return fmt.Errorf("[ERROR] Error setting href: %s", err)
			}

			if vpnGatewayConnection.IkePolicy != nil {
				err = d.Set("ike_policy", dataSourceVPNGatewayConnectionFlattenIkePolicy(*vpnGatewayConnection.IkePolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ike_policy %s", err)
				}
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				err = d.Set("ipsec_policy", dataSourceVPNGatewayConnectionFlattenIpsecPolicy(*vpnGatewayConnection.IpsecPolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ipsec_policy %s", err)
				}
			}
			if err = d.Set("mode", vpnGatewayConnection.Mode); err != nil {
				return fmt.Errorf("[ERROR] Error setting mode: %s", err)
			}
			if err = d.Set("name", vpnGatewayConnection.Name); err != nil {
				return fmt.Errorf("[ERROR] Error setting name: %s", err)
			}

			// breaking changes
			if err = d.Set("establish_mode", vpnGatewayConnection.EstablishMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting establish_mode: %s", err)
			}
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return err
				}
				local = append(local, modelMap)
			}
			if err = d.Set("local", local); err != nil {
				return fmt.Errorf("[ERROR] Error setting local %s", err)
			}

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return err
				}
				peer = append(peer, modelMap)
			}
			if err = d.Set("peer", peer); err != nil {
				return fmt.Errorf("[ERROR] Error setting peer %s", err)
			}
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				if err = d.Set("peer_address", peer.Address); err != nil {
					return fmt.Errorf("[ERROR] Error setting peer_address: %s", err)
				}
			}
			if err = d.Set("psk", vpnGatewayConnection.Psk); err != nil {
				return fmt.Errorf("[ERROR] Error setting psk: %s", err)
			}
			if err = d.Set("resource_type", vpnGatewayConnection.ResourceType); err != nil {
				return fmt.Errorf("[ERROR] Error setting resource_type: %s", err)
			}
			if err = d.Set("status", vpnGatewayConnection.Status); err != nil {
				return fmt.Errorf("[ERROR] Error setting status: %s", err)
			}
			if err := d.Set("status_reasons", resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)); err != nil {
				return fmt.Errorf("[ERROR] Error setting status_reasons: %s", err)
			}
			if err = d.Set("routing_protocol", vpnGatewayConnection.RoutingProtocol); err != nil {
				return fmt.Errorf("[ERROR] Error setting routing_protocol: %s", err)
			}

			if vpnGatewayConnection.Tunnels != nil {
				err = d.Set("tunnels", dataSourceVPNGatewayConnectionFlattenTunnels(vpnGatewayConnection.Tunnels))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting tunnels %s", err)
				}
			}
		}
	case "*vpcv1.VPNGatewayConnectionPolicyMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionPolicyMode)
			d.SetId(fmt.Sprintf("%s/%s", vpn_gateway_id, *vpnGatewayConnection.ID))
			if err = d.Set("admin_state_up", vpnGatewayConnection.AdminStateUp); err != nil {
				return fmt.Errorf("[ERROR] Error setting admin_state_up: %s", err)
			}
			if err = d.Set("authentication_mode", vpnGatewayConnection.AuthenticationMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting authentication_mode: %s", err)
			}
			if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayConnection.CreatedAt)); err != nil {
				return fmt.Errorf("[ERROR] Error setting created_at: %s", err)
			}

			if vpnGatewayConnection.DeadPeerDetection != nil {
				err = d.Set("dead_peer_detection", dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(*vpnGatewayConnection.DeadPeerDetection))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting dead_peer_detection %s", err)
				}
			}
			if err = d.Set("href", vpnGatewayConnection.Href); err != nil {
				return fmt.Errorf("[ERROR] Error setting href: %s", err)
			}

			if vpnGatewayConnection.IkePolicy != nil {
				err = d.Set("ike_policy", dataSourceVPNGatewayConnectionFlattenIkePolicy(*vpnGatewayConnection.IkePolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ike_policy %s", err)
				}
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				err = d.Set("ipsec_policy", dataSourceVPNGatewayConnectionFlattenIpsecPolicy(*vpnGatewayConnection.IpsecPolicy))
				if err != nil {
					return fmt.Errorf("[ERROR] Error setting ipsec_policy %s", err)
				}
			}
			if err = d.Set("mode", vpnGatewayConnection.Mode); err != nil {
				return fmt.Errorf("[ERROR] Error setting mode: %s", err)
			}
			if err = d.Set("name", vpnGatewayConnection.Name); err != nil {
				return fmt.Errorf("[ERROR] Error setting name: %s", err)
			}

			// breaking changes
			if err = d.Set("establish_mode", vpnGatewayConnection.EstablishMode); err != nil {
				return fmt.Errorf("[ERROR] Error setting establish_mode: %s", err)
			}
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return err
				}
				local = append(local, modelMap)
			}
			if err = d.Set("local", local); err != nil {
				return fmt.Errorf("[ERROR] Error setting local %s", err)
			}

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return err
				}
				peer = append(peer, modelMap)
			}
			if err = d.Set("peer", peer); err != nil {
				return fmt.Errorf("[ERROR] Error setting peer %s", err)
			}
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionPolicyModePeer)
				if err = d.Set("peer_address", peer.Address); err != nil {
					return fmt.Errorf("[ERROR] Error setting peer_address: %s", err)
				}
				if len(peer.CIDRs) > 0 {
					err = d.Set("peer_cidrs", peer.CIDRs)
					if err != nil {
						return fmt.Errorf("[ERROR] Error setting Peer CIDRs %s", err)
					}
				}
			}
			if err = d.Set("psk", vpnGatewayConnection.Psk); err != nil {
				return fmt.Errorf("[ERROR] Error setting psk: %s", err)
			}
			if err = d.Set("resource_type", vpnGatewayConnection.ResourceType); err != nil {
				return fmt.Errorf("[ERROR] Error setting resource_type: %s", err)
			}
			if err = d.Set("status", vpnGatewayConnection.Status); err != nil {
				return fmt.Errorf("[ERROR] Error setting status: %s", err)
			}
			if err := d.Set("status_reasons", resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)); err != nil {
				return fmt.Errorf("[ERROR] Error setting status_reasons: %s", err)
			}
			// Deprecated
			if vpnGatewayConnection.Local != nil {
				local := vpnGatewayConnection.Local
				if len(local.CIDRs) > 0 {
					err = d.Set("local_cidrs", local.CIDRs)
					if err != nil {
						return fmt.Errorf("[ERROR] Error setting local CIDRs %s", err)
					}
				}
			}

		}
	}
	return nil
}

func dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(result vpcv1.VPNGatewayConnectionDpd) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNGatewayConnectionDeadPeerDetectionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNGatewayConnectionDeadPeerDetectionToMap(deadPeerDetectionItem vpcv1.VPNGatewayConnectionDpd) (deadPeerDetectionMap map[string]interface{}) {
	deadPeerDetectionMap = map[string]interface{}{}

	if deadPeerDetectionItem.Action != nil {
		deadPeerDetectionMap["action"] = deadPeerDetectionItem.Action
	}
	if deadPeerDetectionItem.Interval != nil {
		deadPeerDetectionMap["interval"] = deadPeerDetectionItem.Interval
	}
	if deadPeerDetectionItem.Timeout != nil {
		deadPeerDetectionMap["timeout"] = deadPeerDetectionItem.Timeout
	}

	return deadPeerDetectionMap
}

func dataSourceVPNGatewayConnectionFlattenIkePolicy(result vpcv1.IkePolicyReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNGatewayConnectionIkePolicyToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNGatewayConnectionIkePolicyToMap(ikePolicyItem vpcv1.IkePolicyReference) (ikePolicyMap map[string]interface{}) {
	ikePolicyMap = map[string]interface{}{}

	if ikePolicyItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNGatewayConnectionIkePolicyDeletedToMap(*ikePolicyItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		ikePolicyMap["deleted"] = deletedList
	}
	if ikePolicyItem.Href != nil {
		ikePolicyMap["href"] = ikePolicyItem.Href
	}
	if ikePolicyItem.ID != nil {
		ikePolicyMap["id"] = ikePolicyItem.ID
	}
	if ikePolicyItem.Name != nil {
		ikePolicyMap["name"] = ikePolicyItem.Name
	}
	if ikePolicyItem.ResourceType != nil {
		ikePolicyMap["resource_type"] = ikePolicyItem.ResourceType
	}

	return ikePolicyMap
}

func dataSourceVPNGatewayConnectionIkePolicyDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNGatewayConnectionFlattenIpsecPolicy(result vpcv1.IPsecPolicyReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNGatewayConnectionIpsecPolicyToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNGatewayConnectionIpsecPolicyToMap(ipsecPolicyItem vpcv1.IPsecPolicyReference) (ipsecPolicyMap map[string]interface{}) {
	ipsecPolicyMap = map[string]interface{}{}

	if ipsecPolicyItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNGatewayConnectionIpsecPolicyDeletedToMap(*ipsecPolicyItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		ipsecPolicyMap["deleted"] = deletedList
	}
	if ipsecPolicyItem.Href != nil {
		ipsecPolicyMap["href"] = ipsecPolicyItem.Href
	}
	if ipsecPolicyItem.ID != nil {
		ipsecPolicyMap["id"] = ipsecPolicyItem.ID
	}
	if ipsecPolicyItem.Name != nil {
		ipsecPolicyMap["name"] = ipsecPolicyItem.Name
	}
	if ipsecPolicyItem.ResourceType != nil {
		ipsecPolicyMap["resource_type"] = ipsecPolicyItem.ResourceType
	}

	return ipsecPolicyMap
}

func dataSourceVPNGatewayConnectionIpsecPolicyDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNGatewayConnectionFlattenTunnels(result []vpcv1.VPNGatewayConnectionStaticRouteModeTunnel) (tunnels []map[string]interface{}) {
	for _, tunnelsItem := range result {
		tunnels = append(tunnels, dataSourceVPNGatewayConnectionTunnelsToMap(tunnelsItem))
	}

	return tunnels
}

func dataSourceVPNGatewayConnectionTunnelsToMap(tunnelsItem vpcv1.VPNGatewayConnectionStaticRouteModeTunnel) (tunnelsMap map[string]interface{}) {
	tunnelsMap = map[string]interface{}{}

	if tunnelsItem.PublicIP != nil {
		tunnelsMap["public_ip_address"] = tunnelsItem.PublicIP.Address
	}
	if tunnelsItem.Status != nil {
		tunnelsMap["status"] = tunnelsItem.Status
	}

	return tunnelsMap
}

// helper methods

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStatusReasonToMap(model *vpcv1.VPNGatewayConnectionStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModeLocal) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentities := []map[string]interface{}{}
	for _, ikeIdentitiesItem := range model.IkeIdentities {
		ikeIdentitiesItemMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(ikeIdentitiesItem)
		if err != nil {
			return modelMap, err
		}
		ikeIdentities = append(ikeIdentities, ikeIdentitiesItemMap)
	}
	modelMap["ike_identities"] = ikeIdentities
	return modelMap, nil
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModeLocalToMap(model *vpcv1.VPNGatewayConnectionPolicyModeLocal) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentities := []map[string]interface{}{}
	for _, ikeIdentitiesItem := range model.IkeIdentities {
		ikeIdentitiesItemMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(ikeIdentitiesItem)
		if err != nil {
			return modelMap, err
		}
		ikeIdentities = append(ikeIdentities, ikeIdentitiesItemMap)
	}
	modelMap["ike_identities"] = ikeIdentities
	modelMap["cidrs"] = model.CIDRs
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model vpcv1.VPNGatewayConnectionIkeIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdnToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostnameToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4ToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyIDToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VPNGatewayConnectionIkeIdentity)
		modelMap["type"] = model.Type
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VPNGatewayConnectionIkeIdentityIntf subtype encountered")
	}
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdnToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostnameToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4ToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyIDToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = string(*model.Value)
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(model vpcv1.VPNGatewayConnectionStaticRouteModePeerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddressToMap(model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdnToMap(model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
		ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
		if err != nil {
			return modelMap, err
		}
		modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
		modelMap["type"] = model.Type
		if model.Address != nil {
			modelMap["address"] = model.Address
		}
		if model.Fqdn != nil {
			modelMap["fqdn"] = model.Fqdn
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VPNGatewayConnectionStaticRouteModePeerIntf subtype encountered")
	}
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerToMap(model vpcv1.VPNGatewayConnectionPolicyModePeerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByAddress); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByAddressToMap(model.(*vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByAddress))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByFqdn); ok {
		return dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByFqdnToMap(model.(*vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByFqdn))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionPolicyModePeer); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VPNGatewayConnectionPolicyModePeer)
		ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
		if err != nil {
			return modelMap, err
		}
		modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
		modelMap["type"] = model.Type
		if model.Address != nil {
			modelMap["address"] = model.Address
		}
		if model.Fqdn != nil {
			modelMap["fqdn"] = model.Fqdn
		}
		if model.CIDRs != nil {
			modelMap["cidrs"] = model.CIDRs
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VPNGatewayConnectionPolicyModePeerIntf subtype encountered")
	}
}

func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddressToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["address"] = model.Address
	return modelMap, nil
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByAddressToMap(model *vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["address"] = model.Address
	if model.CIDRs != nil {
		modelMap["cidrs"] = model.CIDRs
	}
	return modelMap, nil
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdnToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["fqdn"] = model.Fqdn
	return modelMap, nil
}
func dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByFqdnToMap(model *vpcv1.VPNGatewayConnectionPolicyModePeerVPNGatewayConnectionPeerByFqdn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["fqdn"] = model.Fqdn
	if model.CIDRs != nil {
		modelMap["cidrs"] = model.CIDRs
	}
	return modelMap, nil
}

// PrettyPrint print pretty.
func PrettifyPrint(result interface{}) string {
	output, err := json.MarshalIndent(result, "", "    ")

	if err == nil {
		return fmt.Sprintf("%v", string(output))
	}
	return string(output)
}
