// Copyright 2024 Google LLC. All Rights Reserved.
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

func DCLInterconnectAttachmentSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/InterconnectAttachment",
			Description: "The Compute InterconnectAttachment resource",
			StructName:  "InterconnectAttachment",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a InterconnectAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "interconnectAttachment",
						Required:    true,
						Description: "A full instance of a InterconnectAttachment",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a InterconnectAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "interconnectAttachment",
						Required:    true,
						Description: "A full instance of a InterconnectAttachment",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a InterconnectAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "interconnectAttachment",
						Required:    true,
						Description: "A full instance of a InterconnectAttachment",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all InterconnectAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "region",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many InterconnectAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "region",
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
				"InterconnectAttachment": &dcl.Component{
					Title: "InterconnectAttachment",
					ID:    "projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}",
					Locations: []string{
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"region",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"adminEnabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "AdminEnabled",
								Description: "Determines whether this Attachment will carry packets. Not present for PARTNER_PROVIDER.",
							},
							"bandwidth": &dcl.Property{
								Type:        "string",
								GoName:      "Bandwidth",
								GoType:      "InterconnectAttachmentBandwidthEnum",
								Description: "Provisioned bandwidth capacity for the interconnect attachment. For attachments of type DEDICATED, the user can set the bandwidth. For attachments of type PARTNER, the Google Partner that is operating the interconnect must set the bandwidth. Output only for PARTNER type, mutable for PARTNER_PROVIDER and DEDICATED, and can take one of the following values: - BPS_50M: 50 Mbit/s - BPS_100M: 100 Mbit/s - BPS_200M: 200 Mbit/s - BPS_300M: 300 Mbit/s - BPS_400M: 400 Mbit/s - BPS_500M: 500 Mbit/s - BPS_1G: 1 Gbit/s - BPS_2G: 2 Gbit/s - BPS_5G: 5 Gbit/s - BPS_10G: 10 Gbit/s - BPS_20G: 20 Gbit/s - BPS_50G: 50 Gbit/s",
								Enum: []string{
									"BPS_50M",
									"BPS_100M",
									"BPS_200M",
									"BPS_300M",
									"BPS_400M",
									"BPS_500M",
									"BPS_1G",
									"BPS_2G",
									"BPS_5G",
									"BPS_10G",
									"BPS_20G",
									"BPS_50G",
								},
							},
							"candidateSubnets": &dcl.Property{
								Type:        "array",
								GoName:      "CandidateSubnets",
								Description: "Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc). Google will attempt to select an unused /29 from the supplied candidate prefix(es). The request will fail if all possible /29s are in use on Google's edge. If not supplied, Google will randomly select an unused /29 from all of link-local space.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"cloudRouterIPAddress": &dcl.Property{
								Type:        "string",
								GoName:      "CloudRouterIPAddress",
								ReadOnly:    true,
								Description: "IPv4 address + prefix length to be configured on Cloud Router Interface for this interconnect attachment.",
								Immutable:   true,
							},
							"customerRouterIPAddress": &dcl.Property{
								Type:        "string",
								GoName:      "CustomerRouterIPAddress",
								ReadOnly:    true,
								Description: "IPv4 address + prefix length to be configured on the customer router subinterface for this interconnect attachment.",
								Immutable:   true,
							},
							"dataplaneVersion": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "DataplaneVersion",
								Description: "Dataplane version for this InterconnectAttachment. This field is only present for Dataplane version 2 and higher. Absence of this field in the API output indicates that the Dataplane is version 1.",
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource.",
							},
							"edgeAvailabilityDomain": &dcl.Property{
								Type:        "string",
								GoName:      "EdgeAvailabilityDomain",
								GoType:      "InterconnectAttachmentEdgeAvailabilityDomainEnum",
								Description: "Desired availability domain for the attachment. Only available for type PARTNER, at creation time, and can take one of the following values: - AVAILABILITY_DOMAIN_ANY - AVAILABILITY_DOMAIN_1 - AVAILABILITY_DOMAIN_2 For improved reliability, customers should configure a pair of attachments, one per availability domain. The selected availability domain will be provided to the Partner via the pairing key, so that the provisioned circuit will lie in the specified domain. If not specified, the value will default to AVAILABILITY_DOMAIN_ANY.",
								Enum: []string{
									"AVAILABILITY_DOMAIN_ANY",
									"AVAILABILITY_DOMAIN_1",
									"AVAILABILITY_DOMAIN_2",
								},
							},
							"encryption": &dcl.Property{
								Type:        "string",
								GoName:      "Encryption",
								GoType:      "InterconnectAttachmentEncryptionEnum",
								Description: "Indicates the user-supplied encryption option of this VLAN attachment (interconnectAttachment). Can only be specified at attachment creation for PARTNER or DEDICATED attachments. Possible values are: - `NONE` - This is the default value, which means that the VLAN attachment carries unencrypted traffic. VMs are able to send traffic to, or receive traffic from, such a VLAN attachment. - `IPSEC` - The VLAN attachment carries only encrypted traffic that is encrypted by an IPsec device, such as an HA VPN gateway or third-party IPsec VPN. VMs cannot directly send traffic to, or receive traffic from, such a VLAN attachment. To use _IPsec-encrypted Cloud Interconnect_, the VLAN attachment must be created with this option. Not currently available publicly.",
								Enum: []string{
									"NONE",
									"IPSEC",
								},
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "The unique identifier for the resource. This identifier is defined by the server.",
								Immutable:   true,
							},
							"interconnect": &dcl.Property{
								Type:        "string",
								GoName:      "Interconnect",
								Description: "URL of the underlying Interconnect object that this attachment's traffic will traverse through.",
							},
							"ipsecInternalAddresses": &dcl.Property{
								Type:        "array",
								GoName:      "IpsecInternalAddresses",
								Description: "A list of URLs of addresses that have been reserved for the VLAN attachment. Used only for the VLAN attachment that has the encryption option as IPSEC. The addresses must be regional internal IP address ranges. When creating an HA VPN gateway over the VLAN attachment, if the attachment is configured to use a regional internal IP address, then the VPN gateway's IP address is allocated from the IP address range specified here. For example, if the HA VPN gateway's interface 0 is paired to this VLAN attachment, then a regional internal IP address for the VPN gateway interface 0 will be allocated from the IP address specified for this VLAN attachment. If this field is not specified when creating the VLAN attachment, then later on when creating an HA VPN gateway on this VLAN attachment, the HA VPN gateway's IP address is allocated from the regional external IP address pool. Not currently available publicly.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"mtu": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Mtu",
								Description: "Maximum Transmission Unit (MTU), in bytes, of packets passing through this interconnect attachment. Only 1440 and 1500 are allowed. If not specified, the value will default to 1440.",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
							},
							"operationalStatus": &dcl.Property{
								Type:        "string",
								GoName:      "OperationalStatus",
								GoType:      "InterconnectAttachmentOperationalStatusEnum",
								ReadOnly:    true,
								Description: "The current status of whether or not this interconnect attachment is functional, which can take one of the following values: - OS_ACTIVE: The attachment has been turned up and is ready to use. - OS_UNPROVISIONED: The attachment is not ready to use yet, because turnup is not complete.",
								Immutable:   true,
								Enum: []string{
									"OS_ACTIVE",
									"OS_UNPROVISIONED",
								},
							},
							"pairingKey": &dcl.Property{
								Type:        "string",
								GoName:      "PairingKey",
								Description: "The opaque identifier of an PARTNER attachment used to initiate provisioning with a selected partner. Of the form \"XXXXX/region/domain\"",
							},
							"partnerAsn": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "PartnerAsn",
								Description: "Optional BGP ASN for the router supplied by a Layer 3 Partner if they configured BGP on behalf of the customer. Output only for PARTNER type, input only for PARTNER_PROVIDER, not available for DEDICATED.",
							},
							"partnerMetadata": &dcl.Property{
								Type:        "object",
								GoName:      "PartnerMetadata",
								GoType:      "InterconnectAttachmentPartnerMetadata",
								Description: "Informational metadata about Partner attachments from Partners to display to customers. Output only for for PARTNER type, mutable for PARTNER_PROVIDER, not available for DEDICATED.",
								Properties: map[string]*dcl.Property{
									"interconnectName": &dcl.Property{
										Type:        "string",
										GoName:      "InterconnectName",
										Description: "Plain text name of the Interconnect this attachment is connected to, as displayed in the Partner's portal. For instance \"Chicago 1\". This value may be validated to match approved Partner values.",
									},
									"partnerName": &dcl.Property{
										Type:        "string",
										GoName:      "PartnerName",
										Description: "Plain text name of the Partner providing this attachment. This value may be validated to match approved Partner values.",
									},
									"portalUrl": &dcl.Property{
										Type:        "string",
										GoName:      "PortalUrl",
										Description: "URL of the Partner's portal for this Attachment. Partners may customise this to be a deep link to the specific resource on the Partner portal. This value may be validated to match approved Partner values.",
									},
								},
							},
							"privateInterconnectInfo": &dcl.Property{
								Type:        "object",
								GoName:      "PrivateInterconnectInfo",
								GoType:      "InterconnectAttachmentPrivateInterconnectInfo",
								ReadOnly:    true,
								Description: "Information specific to an InterconnectAttachment. This property is populated if the interconnect that this is attached to is of type DEDICATED.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"tag8021q": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Tag8021q",
										ReadOnly:    true,
										Description: "802.1q encapsulation tag to be used for traffic between Google and the customer, going to and from this network and region.",
										Immutable:   true,
									},
								},
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
								Parameter: true,
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								Description: "URL of the region where the regional interconnect attachment resides. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body.",
								Immutable:   true,
							},
							"router": &dcl.Property{
								Type:        "string",
								GoName:      "Router",
								Description: "URL of the Cloud Router to be used for dynamic routing. This router must be in the same region as this InterconnectAttachment. The InterconnectAttachment will automatically connect the Interconnect to the network & region within which the Cloud Router is configured.",
							},
							"satisfiesPzs": &dcl.Property{
								Type:        "boolean",
								GoName:      "SatisfiesPzs",
								ReadOnly:    true,
								Description: "Set to true if the resource satisfies the zone separation organization policy constraints and false otherwise. Defaults to false if the field is not present.",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "Server-defined URL for the resource.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "InterconnectAttachmentStateEnum",
								ReadOnly:    true,
								Description: "The current state of this attachment's functionality. Enum values ACTIVE and UNPROVISIONED are shared by DEDICATED/PRIVATE, PARTNER, and PARTNER_PROVIDER interconnect attachments, while enum values PENDING_PARTNER, PARTNER_REQUEST_RECEIVED, and PENDING_CUSTOMER are used for only PARTNER and PARTNER_PROVIDER interconnect attachments. This state can take one of the following values: - ACTIVE: The attachment has been turned up and is ready to use. - UNPROVISIONED: The attachment is not ready to use yet, because turnup is not complete. - PENDING_PARTNER: A newly-created PARTNER attachment that has not yet been configured on the Partner side. - PARTNER_REQUEST_RECEIVED: A PARTNER attachment is in the process of provisioning after a PARTNER_PROVIDER attachment was created that references it. - PENDING_CUSTOMER: A PARTNER or PARTNER_PROVIDER attachment that is waiting for a customer to activate it. - DEFUNCT: The attachment was deleted externally and is no longer functional. This could be because the associated Interconnect was removed, or because the other side of a Partner attachment was deleted. Possible values: DEPRECATED, OBSOLETE, DELETED, ACTIVE",
								Immutable:   true,
								Enum: []string{
									"DEPRECATED",
									"OBSOLETE",
									"DELETED",
									"ACTIVE",
								},
							},
							"type": &dcl.Property{
								Type:        "string",
								GoName:      "Type",
								GoType:      "InterconnectAttachmentTypeEnum",
								Description: "The type of interconnect attachment this is, which can take one of the following values: - DEDICATED: an attachment to a Dedicated Interconnect. - PARTNER: an attachment to a Partner Interconnect, created by the customer. - PARTNER_PROVIDER: an attachment to a Partner Interconnect, created by the partner. Possible values: PATH, OTHER, PARAMETER",
								Enum: []string{
									"PATH",
									"OTHER",
									"PARAMETER",
								},
							},
							"vlanTag8021q": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "VlanTag8021q",
								Description: "The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094. Only specified at creation time.",
							},
						},
					},
				},
			},
		},
	}
}
