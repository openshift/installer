// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

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
			"peer_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the peer VPN gateway.",
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
			},
			"peer_cidrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The peer CIDRs for this resource.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	vpnGatewayConnection := &vpcv1.VPNGatewayConnection{}

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
				return diag.FromErr(fmt.Errorf("Error reading list of VPN Gateways:%s\n%s", err, detail))
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
			return diag.FromErr(fmt.Errorf("Error reading list of VPN Gateway Connections:%s\n%s", err, detail))
		}

		vpn_gateway_conn_found := false
		for _, connectionItem := range availableVPNGatewayConnections.Connections {
			connection := connectionItem.(*vpcv1.VPNGatewayConnection)
			if *connection.Name == vpn_gateway_connection_name {
				vpnGatewayConnection = connection
				vpn_gateway_conn_found = true
				break
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
		if err != nil || vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection) == nil {
			log.Printf("[DEBUG] GetVPNGatewayConnectionWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetVPNGatewayConnectionWithContext failed %s\n%s", err, response))
		}
		vpnGatewayConnection = vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection)
	}

	d.SetId(fmt.Sprintf("%s/%s", vpn_gateway_id, *vpnGatewayConnection.ID))

	if err = d.Set("admin_state_up", vpnGatewayConnection.AdminStateUp); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting admin_state_up: %s", err))
	}
	if err = d.Set("authentication_mode", vpnGatewayConnection.AuthenticationMode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting authentication_mode: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayConnection.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if vpnGatewayConnection.DeadPeerDetection != nil {
		err = d.Set("dead_peer_detection", dataSourceVPNGatewayConnectionFlattenDeadPeerDetection(*vpnGatewayConnection.DeadPeerDetection))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting dead_peer_detection %s", err))
		}
	}
	if err = d.Set("href", vpnGatewayConnection.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if vpnGatewayConnection.IkePolicy != nil {
		err = d.Set("ike_policy", dataSourceVPNGatewayConnectionFlattenIkePolicy(*vpnGatewayConnection.IkePolicy))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting ike_policy %s", err))
		}
	}

	if vpnGatewayConnection.IpsecPolicy != nil {
		err = d.Set("ipsec_policy", dataSourceVPNGatewayConnectionFlattenIpsecPolicy(*vpnGatewayConnection.IpsecPolicy))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting ipsec_policy %s", err))
		}
	}
	if err = d.Set("mode", vpnGatewayConnection.Mode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mode: %s", err))
	}
	if err = d.Set("name", vpnGatewayConnection.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("peer_address", vpnGatewayConnection.PeerAddress); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting peer_address: %s", err))
	}
	if err = d.Set("psk", vpnGatewayConnection.Psk); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting psk: %s", err))
	}
	if err = d.Set("resource_type", vpnGatewayConnection.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("status", vpnGatewayConnection.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("routing_protocol", vpnGatewayConnection.RoutingProtocol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting routing_protocol: %s", err))
	}

	if vpnGatewayConnection.Tunnels != nil {
		err = d.Set("tunnels", dataSourceVPNGatewayConnectionFlattenTunnels(vpnGatewayConnection.Tunnels))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tunnels %s", err))
		}
	}

	if len(vpnGatewayConnection.LocalCIDRs) > 0 {
		err = d.Set("local_cidrs", vpnGatewayConnection.LocalCIDRs)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting local CIDRs %s", err))
		}
	}

	if len(vpnGatewayConnection.PeerCIDRs) > 0 {
		err = d.Set("peer_cidrs", vpnGatewayConnection.PeerCIDRs)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting Peer CIDRs %s", err))
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

func dataSourceVPNGatewayConnectionIkePolicyDeletedToMap(deletedItem vpcv1.IkePolicyReferenceDeleted) (deletedMap map[string]interface{}) {
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

func dataSourceVPNGatewayConnectionIpsecPolicyDeletedToMap(deletedItem vpcv1.IPsecPolicyReferenceDeleted) (deletedMap map[string]interface{}) {
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

func dataSourceVPNGatewayConnectionTunnelsPublicIPToMap(publicIPItem vpcv1.IP) (publicIPMap map[string]interface{}) {
	publicIPMap = map[string]interface{}{}

	if publicIPItem.Address != nil {
		publicIPMap["address"] = publicIPItem.Address
	}

	return publicIPMap
}
