// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package compute

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLVpnTunnelSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/VpnTunnel",
			Description: "The Compute VpnTunnel resource",
			StructName:  "VpnTunnel",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a VpnTunnel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "vpnTunnel",
						Required:    true,
						Description: "A full instance of a VpnTunnel",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a VpnTunnel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "vpnTunnel",
						Required:    true,
						Description: "A full instance of a VpnTunnel",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a VpnTunnel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "vpnTunnel",
						Required:    true,
						Description: "A full instance of a VpnTunnel",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all VpnTunnel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many VpnTunnel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"VpnTunnel": &dcl.Component{
					Title:           "VpnTunnel",
					ID:              "projects/{{project}}/regions/{{location}}/vpnTunnels/{{name}}",
					UsesStateHint:   true,
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"sharedSecret",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource.",
								Immutable:   true,
							},
							"detailedStatus": &dcl.Property{
								Type:        "string",
								GoName:      "DetailedStatus",
								ReadOnly:    true,
								Description: "Detailed status message for the VPN tunnel.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "The unique identifier for the resource. This identifier is defined by the server.",
								Immutable:   true,
							},
							"ikeVersion": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "IkeVersion",
								Description: "IKE protocol version to use when establishing the VPN tunnel with the peer VPN gateway. Acceptable IKE versions are 1 or 2. The default version is 2.",
								Immutable:   true,
								Default:     2,
							},
							"localTrafficSelector": &dcl.Property{
								Type:          "array",
								GoName:        "LocalTrafficSelector",
								Description:   "Local traffic selector to use when establishing the VPN tunnel with the peer VPN gateway. The value should be a CIDR formatted string, for example: 192.168.0.0/16. The ranges must be disjoint. Only IPv4 is supported.",
								Immutable:     true,
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "set",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "Name of the region where the VPN tunnel resides.",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt) Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
							"peerExternalGateway": &dcl.Property{
								Type:        "string",
								GoName:      "PeerExternalGateway",
								Description: "URL of the peer side external VPN gateway to which this VPN tunnel is connected. Provided by the client when the VPN tunnel is created. This field is exclusive with the field peerGcpGateway.",
								Immutable:   true,
							},
							"peerExternalGatewayInterface": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "PeerExternalGatewayInterface",
								Description: "The interface ID of the external VPN gateway to which this VPN tunnel is connected. Provided by the client when the VPN tunnel is created.",
								Immutable:   true,
							},
							"peerGcpGateway": &dcl.Property{
								Type:        "string",
								GoName:      "PeerGcpGateway",
								Description: "URL of the peer side HA GCP VPN gateway to which this VPN tunnel is connected. Provided by the client when the VPN tunnel is created. This field can be used when creating highly available VPN from VPC network to VPC network, the field is exclusive with the field peerExternalGateway. If provided, the VPN tunnel will automatically use the same vpnGatewayInterface ID in the peer GCP VPN gateway.",
								Immutable:   true,
							},
							"peerIP": &dcl.Property{
								Type:          "string",
								GoName:        "PeerIP",
								Description:   "IP address of the peer VPN gateway. Only IPv4 is supported.",
								Immutable:     true,
								ServerDefault: true,
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"remoteTrafficSelector": &dcl.Property{
								Type:          "array",
								GoName:        "RemoteTrafficSelector",
								Description:   "Remote traffic selectors to use when establishing the VPN tunnel with the peer VPN gateway. The value should be a CIDR formatted string, for example: 192.168.0.0/16. The ranges should be disjoint. Only IPv4 is supported.",
								Immutable:     true,
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "set",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"router": &dcl.Property{
								Type:        "string",
								GoName:      "Router",
								Description: "URL of the router resource to be used for dynamic routing.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/Router",
										Field:    "selfLink",
									},
								},
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "Server-defined URL for the resource.",
								Immutable:   true,
							},
							"sharedSecret": &dcl.Property{
								Type:        "string",
								GoName:      "SharedSecret",
								Description: "Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway.",
								Immutable:   true,
								Sensitive:   true,
								Unreadable:  true,
							},
							"sharedSecretHash": &dcl.Property{
								Type:        "string",
								GoName:      "SharedSecretHash",
								ReadOnly:    true,
								Description: "Hash of the shared secret.",
								Immutable:   true,
							},
							"status": &dcl.Property{
								Type:        "string",
								GoName:      "Status",
								GoType:      "VpnTunnelStatusEnum",
								ReadOnly:    true,
								Description: "The status of the VPN tunnel, which can be one of the following:  * PROVISIONING: Resource is being allocated for the VPN tunnel.  * WAITING_FOR_FULL_CONFIG: Waiting to receive all VPN-related configs from   the user. Network, TargetVpnGateway, VpnTunnel, ForwardingRule, and Route   resources are needed to setup the VPN tunnel.  * FIRST_HANDSHAKE: Successful first handshake with the peer VPN.  * ESTABLISHED: Secure session is successfully established with the peer VPN.  * NETWORK_ERROR: Deprecated, replaced by NO_INCOMING_PACKETS  * AUTHORIZATION_ERROR: Auth error (for example, bad shared secret).  * NEGOTIATION_FAILURE: Handshake failed.  * DEPROVISIONING: Resources are being deallocated for the VPN tunnel.  * FAILED: Tunnel creation has failed and the tunnel is not ready to be used.  * NO_INCOMING_PACKETS: No incoming packets from peer.  * REJECTED: Tunnel configuration was rejected, can be result of being blocklisted.  * ALLOCATING_RESOURCES: Cloud VPN is in the process of allocating all required resources.  * STOPPED: Tunnel is stopped due to its Forwarding Rules being deleted for Classic VPN tunnels or the project is in frozen state.  * PEER_IDENTITY_MISMATCH: Peer identity does not match peer IP, probably behind NAT.  * TS_NARROWING_NOT_ALLOWED: Traffic selector narrowing not allowed for an HA-VPN tunnel.",
								Immutable:   true,
								Enum: []string{
									"PROVISIONING",
									"WAITING_FOR_FULL_CONFIG",
									"FIRST_HANDSHAKE",
									"ESTABLISHED",
									"NO_INCOMING_PACKETS",
									"AUTHORIZATION_ERROR",
									"NEGOTIATION_FAILURE",
									"DEPROVISIONING",
									"FAILED",
									"REJECTED",
									"ALLOCATING_RESOURCES",
									"STOPPED",
									"PEER_IDENTITY_MISMATCH",
									"TS_NARROWING_NOT_ALLOWED",
								},
							},
							"targetVpnGateway": &dcl.Property{
								Type:        "string",
								GoName:      "TargetVpnGateway",
								Description: "URL of the Target VPN gateway with which this VPN tunnel is associated. Provided by the client when the VPN tunnel is created.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/TargetVpnGateway",
										Field:    "selfLink",
									},
								},
							},
							"vpnGateway": &dcl.Property{
								Type:        "string",
								GoName:      "VpnGateway",
								Description: "URL of the VPN gateway with which this VPN tunnel is associated. Provided by the client when the VPN tunnel is created. This must be used (instead of target_vpn_gateway) if a High Availability VPN gateway resource is created.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/VpnGateway",
										Field:    "selfLink",
									},
								},
							},
							"vpnGatewayInterface": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "VpnGatewayInterface",
								Description: "The interface ID of the VPN gateway with which this VPN tunnel is associated.",
								Immutable:   true,
								SendEmpty:   true,
							},
						},
					},
				},
			},
		},
	}
}
