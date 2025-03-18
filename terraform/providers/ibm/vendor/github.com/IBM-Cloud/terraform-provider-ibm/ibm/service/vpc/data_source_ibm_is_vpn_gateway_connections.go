// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isvpnGatewayConnections  = "connections"
	isVPNGatewayID           = "vpn_gateway"
	isVPNGatewayConnectionID = "id"
)

func DataSourceIBMISVPNGatewayConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewayConnectionsRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the vpn gateway connection",
			},
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
						"distribute_traffic": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the traffic is distributed between the `up` tunnels of the VPN gateway connection when the VPC route's next hop is a VPN connection. If `false`, the traffic is only routed through the `up` tunnel with the lower `public_ip` address.",
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
							Deprecated:  "peer_address is deprecated, use peer instead",
						},
						// new breaking change
						"establish_mode": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The establish mode of the VPN gateway connection:- `bidirectional`: Either side of the VPN gateway can initiate IKE protocol   negotiations or rekeying processes.- `peer_only`: Only the peer can initiate IKE protocol negotiations for this VPN gateway   connection. Additionally, the peer is responsible for initiating the rekeying process   after the connection is established. If rekeying does not occur, the VPN gateway   connection will be brought down after its lifetime expires.",
						},
						"routing_protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing protocols for this VPN gateway connection.",
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
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "VPN gateway connection local CIDRs",
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
									"cidrs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The peer CIDRs for this resource.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
								},
							},
						},
						"psk": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The preshared key.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The href of the vpn gateway connection.",
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
						isVPNGatewayConnectionStatusreasons: {
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
							Deprecated:  "local_cidrs is deprecated, use local instead",
						},

						isVPNGatewayConnectionPeerCIDRS: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "VPN gateway connection peer CIDRs",
							Deprecated:  "peer_cidrs is deprecated, use peer instead",
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
	if statusIntf, ok := d.GetOk("status"); ok {
		status := statusIntf.(string)
		listvpnGWConnectionOptions.Status = &status
	}
	availableVPNGatewayConnections, detail, err := sess.ListVPNGatewayConnections(listvpnGWConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error reading list of VPN Gateway Connections:%s\n%s", err, detail)
	}
	vpngatewayconnections := make([]map[string]interface{}, 0)
	for _, instance := range availableVPNGatewayConnections.Connections {
		gatewayconnection, err := getvpnGatewayConnectionIntfData(instance)
		if err != nil {
			return err
		}
		vpngatewayconnections = append(vpngatewayconnections, gatewayconnection)
	}

	d.SetId(dataSourceIBMVPNGatewayConnectionsID(d))
	d.Set(isvpnGatewayConnections, vpngatewayconnections)
	return nil
}

func getvpnGatewayConnectionIntfData(vpnGatewayConnectionIntf vpcv1.VPNGatewayConnectionIntf) (map[string]interface{}, error) {
	gatewayconnection := map[string]interface{}{}
	switch reflect.TypeOf(vpnGatewayConnectionIntf).String() {
	case "*vpcv1.VPNGatewayConnection":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection)
			gatewayconnection["id"] = vpnGatewayConnection.ID
			gatewayconnection["admin_state_up"] = vpnGatewayConnection.AdminStateUp
			gatewayconnection["authentication_mode"] = vpnGatewayConnection.AuthenticationMode
			gatewayconnection["created_at"] = flex.DateTimeToString(vpnGatewayConnection.CreatedAt)

			if vpnGatewayConnection.DeadPeerDetection != nil {
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionAction] = vpnGatewayConnection.DeadPeerDetection.Action
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionInterval] = vpnGatewayConnection.DeadPeerDetection.Interval
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionTimeout] = vpnGatewayConnection.DeadPeerDetection.Timeout
			}
			gatewayconnection["href"] = vpnGatewayConnection.Href

			if vpnGatewayConnection.IkePolicy != nil {
				gatewayconnection["ike_policy"] = vpnGatewayConnection.IkePolicy.ID
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				gatewayconnection["ipsec_policy"] = vpnGatewayConnection.IpsecPolicy.ID
			}
			gatewayconnection["mode"] = vpnGatewayConnection.Mode
			gatewayconnection["name"] = vpnGatewayConnection.Name
			gatewayconnection["distribute_traffic"] = vpnGatewayConnection.DistributeTraffic

			// breaking changes
			gatewayconnection["establish_mode"] = vpnGatewayConnection.EstablishMode
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return gatewayconnection, err
				}
				local = append(local, modelMap)
			}
			gatewayconnection["local"] = local

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return gatewayconnection, err
				}
				peer = append(peer, modelMap)
			}
			gatewayconnection["peer"] = peer
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				gatewayconnection["peer_address"] = peer.Address
			}
			gatewayconnection["psk"] = vpnGatewayConnection.Psk
			gatewayconnection["resource_type"] = vpnGatewayConnection.ResourceType
			gatewayconnection["status"] = vpnGatewayConnection.Status
			gatewayconnection["status_reasons"] = resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)
			gatewayconnection["routing_protocol"] = vpnGatewayConnection.RoutingProtocol

			if vpnGatewayConnection.Tunnels != nil {
				gatewayconnection["tunnels"] = dataSourceVPNGatewayConnectionsFlattenTunnels(vpnGatewayConnection.Tunnels)
			}
		}
	case "*vpcv1.VPNGatewayConnectionRouteMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionRouteMode)
			gatewayconnection["id"] = vpnGatewayConnection.ID
			gatewayconnection["admin_state_up"] = vpnGatewayConnection.AdminStateUp
			gatewayconnection["authentication_mode"] = vpnGatewayConnection.AuthenticationMode
			gatewayconnection["created_at"] = flex.DateTimeToString(vpnGatewayConnection.CreatedAt)

			if vpnGatewayConnection.DeadPeerDetection != nil {
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionAction] = vpnGatewayConnection.DeadPeerDetection.Action
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionInterval] = vpnGatewayConnection.DeadPeerDetection.Interval
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionTimeout] = vpnGatewayConnection.DeadPeerDetection.Timeout
			}
			gatewayconnection["href"] = vpnGatewayConnection.Href
			if vpnGatewayConnection.IkePolicy != nil {
				gatewayconnection["ike_policy"] = vpnGatewayConnection.IkePolicy.ID
			}
			gatewayconnection["distribute_traffic"] = vpnGatewayConnection.DistributeTraffic

			if vpnGatewayConnection.IpsecPolicy != nil {
				gatewayconnection["ipsec_policy"] = vpnGatewayConnection.IpsecPolicy.ID
			}
			gatewayconnection["mode"] = vpnGatewayConnection.Mode
			gatewayconnection["name"] = vpnGatewayConnection.Name

			// breaking changes
			gatewayconnection["establish_mode"] = vpnGatewayConnection.EstablishMode
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return gatewayconnection, err
				}
				local = append(local, modelMap)
			}
			gatewayconnection["local"] = local

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return gatewayconnection, err
				}
				peer = append(peer, modelMap)
			}
			gatewayconnection["peer"] = peer
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				gatewayconnection["peer_address"] = peer.Address
			}
			gatewayconnection["psk"] = vpnGatewayConnection.Psk
			gatewayconnection["resource_type"] = vpnGatewayConnection.ResourceType
			gatewayconnection["status"] = vpnGatewayConnection.Status
			gatewayconnection["status_reasons"] = resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)
			gatewayconnection["routing_protocol"] = vpnGatewayConnection.RoutingProtocol

			if vpnGatewayConnection.Tunnels != nil {
				gatewayconnection["tunnels"] = dataSourceVPNGatewayConnectionsFlattenTunnels(vpnGatewayConnection.Tunnels)
			}
		}
	case "*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode)
			gatewayconnection["id"] = vpnGatewayConnection.ID
			gatewayconnection["admin_state_up"] = vpnGatewayConnection.AdminStateUp
			gatewayconnection["authentication_mode"] = vpnGatewayConnection.AuthenticationMode
			gatewayconnection["created_at"] = flex.DateTimeToString(vpnGatewayConnection.CreatedAt)

			if vpnGatewayConnection.DeadPeerDetection != nil {
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionAction] = vpnGatewayConnection.DeadPeerDetection.Action
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionInterval] = vpnGatewayConnection.DeadPeerDetection.Interval
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionTimeout] = vpnGatewayConnection.DeadPeerDetection.Timeout
			}
			gatewayconnection["distribute_traffic"] = vpnGatewayConnection.DistributeTraffic
			gatewayconnection["href"] = vpnGatewayConnection.Href
			if vpnGatewayConnection.IkePolicy != nil {
				gatewayconnection["ike_policy"] = vpnGatewayConnection.IkePolicy.ID
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				gatewayconnection["ipsec_policy"] = vpnGatewayConnection.IpsecPolicy.ID
			}
			gatewayconnection["mode"] = vpnGatewayConnection.Mode
			gatewayconnection["name"] = vpnGatewayConnection.Name

			// breaking changes
			gatewayconnection["establish_mode"] = vpnGatewayConnection.EstablishMode
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return gatewayconnection, err
				}
				local = append(local, modelMap)
			}
			gatewayconnection["local"] = local

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionStaticRouteModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return gatewayconnection, err
				}
				peer = append(peer, modelMap)
			}
			gatewayconnection["peer"] = peer
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
				gatewayconnection["peer_address"] = peer.Address
			}
			gatewayconnection["psk"] = vpnGatewayConnection.Psk
			gatewayconnection["resource_type"] = vpnGatewayConnection.ResourceType
			gatewayconnection["status"] = vpnGatewayConnection.Status
			gatewayconnection["status_reasons"] = resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)
			gatewayconnection["routing_protocol"] = vpnGatewayConnection.RoutingProtocol

			if vpnGatewayConnection.Tunnels != nil {
				gatewayconnection["tunnels"] = dataSourceVPNGatewayConnectionsFlattenTunnels(vpnGatewayConnection.Tunnels)
			}
		}
	case "*vpcv1.VPNGatewayConnectionPolicyMode":
		{
			vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnectionPolicyMode)
			gatewayconnection["id"] = vpnGatewayConnection.ID
			gatewayconnection["admin_state_up"] = vpnGatewayConnection.AdminStateUp
			gatewayconnection["authentication_mode"] = vpnGatewayConnection.AuthenticationMode
			gatewayconnection["created_at"] = flex.DateTimeToString(vpnGatewayConnection.CreatedAt)

			if vpnGatewayConnection.DeadPeerDetection != nil {
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionAction] = vpnGatewayConnection.DeadPeerDetection.Action
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionInterval] = vpnGatewayConnection.DeadPeerDetection.Interval
				gatewayconnection[isVPNGatewayConnectionDeadPeerDetectionTimeout] = vpnGatewayConnection.DeadPeerDetection.Timeout
			}
			gatewayconnection["href"] = vpnGatewayConnection.Href
			if vpnGatewayConnection.IkePolicy != nil {
				gatewayconnection["ike_policy"] = vpnGatewayConnection.IkePolicy.ID
			}

			if vpnGatewayConnection.IpsecPolicy != nil {
				gatewayconnection["ipsec_policy"] = vpnGatewayConnection.IpsecPolicy.ID
			}
			gatewayconnection["mode"] = vpnGatewayConnection.Mode
			gatewayconnection["name"] = vpnGatewayConnection.Name

			// breaking changes
			gatewayconnection["establish_mode"] = vpnGatewayConnection.EstablishMode
			local := []map[string]interface{}{}
			if vpnGatewayConnection.Local != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModeLocalToMap(vpnGatewayConnection.Local)
				if err != nil {
					return gatewayconnection, err
				}
				local = append(local, modelMap)
			}
			gatewayconnection["local"] = local

			peer := []map[string]interface{}{}
			if vpnGatewayConnection.Peer != nil {
				modelMap, err := dataSourceIBMIsVPNGatewayConnectionVPNGatewayConnectionPolicyModePeerToMap(vpnGatewayConnection.Peer)
				if err != nil {
					return gatewayconnection, err
				}
				peer = append(peer, modelMap)
			}
			gatewayconnection["peer"] = peer
			// Deprecated
			if vpnGatewayConnection.Peer != nil {
				peer := vpnGatewayConnection.Peer.(*vpcv1.VPNGatewayConnectionPolicyModePeer)
				gatewayconnection["peer_address"] = peer.Address
				if len(peer.CIDRs) > 0 {
					gatewayconnection["peer_cidrs"] = peer.CIDRs
				}
			}
			gatewayconnection["psk"] = vpnGatewayConnection.Psk
			gatewayconnection["resource_type"] = vpnGatewayConnection.ResourceType
			gatewayconnection["status"] = vpnGatewayConnection.Status
			gatewayconnection["status_reasons"] = resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)
			// Deprecated
			if vpnGatewayConnection.Local != nil {
				local := vpnGatewayConnection.Local
				if len(local.CIDRs) > 0 {
					gatewayconnection["local_cidrs"] = local.CIDRs
				}
			}

		}
	}
	return gatewayconnection, nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewayConnectionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModeLocalToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModeLocal) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentities := []map[string]interface{}{}
	for _, ikeIdentitiesItem := range model.IkeIdentities {
		ikeIdentitiesItemMap, err := dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityToMap(ikeIdentitiesItem)
		if err != nil {
			return modelMap, err
		}
		ikeIdentities = append(ikeIdentities, ikeIdentitiesItemMap)
	}
	modelMap["ike_identities"] = ikeIdentities
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityToMap(model vpcv1.VPNGatewayConnectionIkeIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdnToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostnameToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4ToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyIDToMap(model.(*vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID))
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

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdnToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityFqdn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostnameToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityHostname) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4ToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityIPv4) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyIDToMap(model *vpcv1.VPNGatewayConnectionIkeIdentityVPNGatewayConnectionIkeIdentityKeyID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = string(*model.Value)
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModePeerToMap(model vpcv1.VPNGatewayConnectionStaticRouteModePeerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddressToMap(model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn); ok {
		return dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdnToMap(model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn))
	} else if _, ok := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VPNGatewayConnectionStaticRouteModePeer)
		ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
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

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddressToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["address"] = model.Address
	return modelMap, nil
}

func dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdnToMap(model *vpcv1.VPNGatewayConnectionStaticRouteModePeerVPNGatewayConnectionPeerByFqdn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	ikeIdentityMap, err := dataSourceIBMIsVPNGatewayConnectionsVPNGatewayConnectionIkeIdentityToMap(model.IkeIdentity)
	if err != nil {
		return modelMap, err
	}
	modelMap["ike_identity"] = []map[string]interface{}{ikeIdentityMap}
	modelMap["type"] = model.Type
	modelMap["fqdn"] = model.Fqdn
	return modelMap, nil
}
func dataSourceVPNGatewayConnectionsFlattenTunnels(result []vpcv1.VPNGatewayConnectionStaticRouteModeTunnel) (tunnels []map[string]interface{}) {
	for _, tunnelsItem := range result {
		tunnels = append(tunnels, dataSourceVPNGatewayConnectionsTunnelsToMap(tunnelsItem))
	}

	return tunnels
}

func dataSourceVPNGatewayConnectionsTunnelsToMap(tunnelsItem vpcv1.VPNGatewayConnectionStaticRouteModeTunnel) (tunnelsMap map[string]interface{}) {
	tunnelsMap = map[string]interface{}{}

	if tunnelsItem.PublicIP != nil {
		tunnelsMap["address"] = tunnelsItem.PublicIP.Address
	}
	if tunnelsItem.Status != nil {
		tunnelsMap["status"] = tunnelsItem.Status
	}

	return tunnelsMap
}
