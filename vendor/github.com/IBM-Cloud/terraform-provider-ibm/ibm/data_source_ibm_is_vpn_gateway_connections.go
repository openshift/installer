// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isvpnGatewayConnections  = "connections"
	isVPNGatewayID           = "vpn_gateway"
	isVPNGatewayConnectionID = "id"
)

func dataSourceIBMISVPNGatewayConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewayConnectionsRead,

		Schema: map[string]*schema.Schema{

			isVPNGatewayID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier ",
			},

			isvpnGatewayConnections: {
				Type:        schema.TypeList,
				Description: "Collection of VPN Gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						isVPNGatewayConnectionAdminAuthenticationmode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The authentication mode",
						},
						isVPNGatewayConnectionCreatedat: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this VPN gateway connection was created",
						},
						isVPNGatewayConnectionAdminStateup: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "VPN gateway connection admin state",
						},
						isVPNGatewayConnectionDeadPeerDetectionAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action detection for dead peer detection action",
						},
						isVPNGatewayConnectionDeadPeerDetectionInterval: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval for dead peer detection interval",
						},
						isVPNGatewayConnectionDeadPeerDetectionTimeout: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout for dead peer detection",
						},
						isVPNGatewayConnectionID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN gateway connection",
						},

						isVPNGatewayConnectionIKEPolicy: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN gateway connection IKE Policy",
						},
						isVPNGatewayConnectionIPSECPolicy: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP security policy for vpn gateway connection",
						},
						isVPNGatewayConnectionMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode of the VPN gateway",
						},
						isVPNGatewayConnectionName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN Gateway connection name",
						},
						isVPNGatewayConnectionPeerAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN gateway connection peer address",
						},
						isVPNGatewayConnectionResourcetype: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
						isVPNGatewayConnectionStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN gateway connection status",
						},

						isVPNGatewayConnectionTunnels: {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							MinItems:    0,
							Description: "The VPN tunnel configuration for this VPN gateway connection (in static route mode)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address of the VPN gateway member in which the tunnel resides",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the VPN Tunnel",
									},
								},
							},
						},
						isVPNGatewayConnectionLocalCIDRS: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "VPN gateway connection local CIDRs",
						},

						isVPNGatewayConnectionPeerCIDRS: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "VPN gateway connection peer CIDRs",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMVPNGatewayConnectionsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	vpngatewayID := d.Get(isVPNGatewayID).(string)
	listvpnGWConnectionOptions := sess.NewListVPNGatewayConnectionsOptions(vpngatewayID)

	availableVPNGatewayConnections, detail, err := sess.ListVPNGatewayConnections(listvpnGWConnectionOptions)
	if err != nil {
		return fmt.Errorf("Error reading list of VPN Gateway Connections:%s\n%s", err, detail)
	}
	vpngatewayconnections := make([]map[string]interface{}, 0)
	for _, instance := range availableVPNGatewayConnections.Connections {
		gatewayconnection := map[string]interface{}{}
		data := instance.(*vpcv1.VPNGatewayConnection)
		gatewayconnection[isVPNGatewayConnectionAdminAuthenticationmode] = *data.AuthenticationMode
		gatewayconnection[isVPNGatewayConnectionCreatedat] = data.CreatedAt.String()
		gatewayconnection[isVPNGatewayConnectionAdminStateup] = *data.AdminStateUp
		gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionAction] = *data.DeadPeerDetection.Action
		gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionInterval] = *data.DeadPeerDetection.Interval
		gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionTimeout] = *data.DeadPeerDetection.Timeout
		gatewayconnection[isVPNGatewayConnectionID] = *data.ID

		if data.IkePolicy != nil {
			gatewayconnection[isVPNGatewayConnectionIKEPolicy] = *data.IkePolicy.ID
		}
		if data.IpsecPolicy != nil {
			gatewayconnection[isVPNGatewayConnectionIPSECPolicy] = *data.IpsecPolicy.ID
		}
		if data.LocalCIDRs != nil {
			gatewayconnection[isVPNGatewayConnectionLocalCIDRS] = flattenStringList(data.LocalCIDRs)
		}
		if data.PeerCIDRs != nil {
			gatewayconnection[isVPNGatewayConnectionPeerCIDRS] = flattenStringList(data.PeerCIDRs)
		}
		gatewayconnection[isVPNGatewayConnectionMode] = *data.Mode
		gatewayconnection[isVPNGatewayConnectionName] = *data.Name
		gatewayconnection[isVPNGatewayConnectionPeerAddress] = *data.PeerAddress
		gatewayconnection[isVPNGatewayConnectionResourcetype] = *data.ResourceType
		gatewayconnection[isVPNGatewayConnectionStatus] = *data.Status
		//if data.Tunnels != nil {
		if len(data.Tunnels) > 0 {
			vpcTunnelsList := make([]map[string]interface{}, 0)
			for _, vpcTunnel := range data.Tunnels {
				currentTunnel := map[string]interface{}{}
				if vpcTunnel.PublicIP != nil {
					if vpcTunnel.PublicIP != nil {
						currentTunnel["address"] = *vpcTunnel.PublicIP.Address
					}
					if vpcTunnel.Status != nil {
						currentTunnel["status"] = *vpcTunnel.Status
					}
					vpcTunnelsList = append(vpcTunnelsList, currentTunnel)
				}
			}
			gatewayconnection[isVPNGatewayConnectionTunnels] = vpcTunnelsList
		}

		vpngatewayconnections = append(vpngatewayconnections, gatewayconnection)
	}

	d.SetId(dataSourceIBMVPNGatewayConnectionsID(d))
	d.Set(isvpnGatewayConnections, vpngatewayconnections)
	return nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewayConnectionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
