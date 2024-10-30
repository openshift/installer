// Copyright IBM Corp. 2017, 2021, 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dlGateways         = "gateways"
	dlGatewaysId       = "id"
	dlSpecificPrefixes = "specific_prefixes"
)

func DataSourceIBMDLGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLGatewaysRead,
		Schema: map[string]*schema.Schema{
			dlGateways: {
				Type:        schema.TypeList,
				Description: "Collection of direct link gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlGatewaysId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the data source gateways",
						},
						dlAsPrepends: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of AS Prepend configuration information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlCreatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time AS Prepend was created",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time AS Prepend was created",
									},
									dlLength: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of times the ASN to appended to the AS Path",
									},
									dlPolicy: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route type this AS Prepend applies to",
									},
									dlPrefix: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Comma separated list of prefixes this AS Prepend applies to. Maximum of 10 prefixes. If not specified, this AS Prepend applies to all prefixes.",
									},
									dlSpecificPrefixes: {
										Type:        schema.TypeList,
										Description: "Array of prefixes this AS Prepend applies to",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									dlUpdatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time AS Prepend was updated",
									},
								},
							},
						},
						dlDefault_export_route_filter: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default directional route filter action that applies to routes that do not match any directional route filters",
						},
						dlDefault_import_route_filter: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default directional route filter action that applies to routes that do not match any directional route filters",
						},
						dlAuthenticationKey: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP MD5 authentication key",
						},
						dlBfdInterval: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BFD Interval",
						},
						dlBfdMultiplier: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BFD Multiplier",
						},
						dlBfdStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BFD Status",
						},
						dlBfdStatusUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BFD Status",
						},
						dlBgpAsn: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP ASN",
						},
						dlBgpBaseCidr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP base CIDR",
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
						dlBgpStatusUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Date and time BGP status was updated",
						},
						dlCompletionNoticeRejectReason: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason for completion notice rejection",
						},
						dlConnectionMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of services this Gateway is attached to. Mode transit means this Gateway will be attached to Transit Gateway Service and direct means this Gateway will be attached to vpc or classic connection",
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
						dlCrossConnectRouter: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cross connect router",
						},
						dlGlobal: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Gateways with global routing (true) can connect to networks outside their associated region",
						},
						dlLinkStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway link status",
						},
						dlLinkStatusUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Date and time Link status was updated",
						},
						dlLocationDisplayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway location long name",
						},
						dlLocationName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway location",
						},
						dlMetered: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Metered billing option",
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
						dlMacSecConfig: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "MACsec configuration information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlActive: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate whether MACsec protection should be active (true) or inactive (false) for this MACsec enabled gateway",
									},
									dlActiveCak: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Active connectivity association key.",
									},
									dlPrimaryCak: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Desired primary connectivity association key. Keys for a MACsec configuration must have names with an even number of characters from [0-9a-fA-F]",
									},
									dlFallbackCak: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Fallback connectivity association key. Keys used for MACsec configuration must have names with an even number of characters from [0-9a-fA-F]",
									},
									dlSakExpiryTime: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Secure Association Key (SAK) expiry time in seconds",
									},
									dlSecurityPolicy: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Packets without MACsec headers are not dropped when security_policy is should_secure.",
									},
									dlWindowSize: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Replay protection window size",
									},
									dlCipherSuite: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SAK cipher suite",
									},
									dlConfidentialityOffset: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Confidentiality Offset",
									},
									dlCryptographicAlgorithm: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cryptographic Algorithm",
									},
									dlKeyServerPriority: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Key Server Priority",
									},
									dlMacSecConfigStatus: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current status of MACsec on the device for this gateway",
									},
								},
							},
						},
						dlChangeRequest: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changes pending approval for provider managed Direct Link Connect gateways",
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

func dataSourceIBMDLGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	listGatewaysOptionsModel := &directlinkv1.ListGatewaysOptions{}
	listGateways, response, err := directLink.ListGateways(listGatewaysOptionsModel)
	if err != nil {
		log.Println("[WARN] Error listing dl Gateway", response, err)
		return err
	}
	gateways := make([]map[string]interface{}, 0)
	for _, gwIntf := range listGateways.Gateways {

		gateway := map[string]interface{}{}
		instance := gwIntf.(*directlinkv1.GatewayCollectionGatewaysItem)

		if instance.ID != nil {
			gateway["id"] = *instance.ID
		}
		if instance.Name != nil {
			gateway[dlName] = *instance.Name
		}
		if instance.Crn != nil {
			gateway[dlCrn] = *instance.Crn
		}
		if instance.BgpAsn != nil {
			gateway[dlBgpAsn] = *instance.BgpAsn
		}
		if instance.BgpIbmCidr != nil {
			gateway[dlBgpIbmCidr] = *instance.BgpIbmCidr
		}
		if instance.BgpIbmAsn != nil {
			gateway[dlBgpIbmAsn] = *instance.BgpIbmAsn
		}
		if instance.Metered != nil {
			gateway[dlMetered] = *instance.Metered
		}
		if instance.CrossConnectRouter != nil {
			gateway[dlCrossConnectRouter] = *instance.CrossConnectRouter
		}
		if instance.BgpBaseCidr != nil {
			gateway[dlBgpBaseCidr] = *instance.BgpBaseCidr
		}
		if instance.BgpCerCidr != nil {
			gateway[dlBgpCerCidr] = *instance.BgpCerCidr
		}

		if instance.ProviderApiManaged != nil {
			gateway[dlProviderAPIManaged] = *instance.ProviderApiManaged
		}
		if instance.Type != nil {
			gateway[dlType] = *instance.Type
		}
		if instance.SpeedMbps != nil {
			gateway[dlSpeedMbps] = *instance.SpeedMbps
		}
		if instance.OperationalStatus != nil {
			gateway[dlOperationalStatus] = *instance.OperationalStatus
		}
		if instance.BgpStatus != nil {
			gateway[dlBgpStatus] = *instance.BgpStatus
		}
		if instance.BgpStatusUpdatedAt != nil {
			gateway[dlBgpStatusUpdatedAt] = instance.BgpStatusUpdatedAt.String()
		}
		if instance.LocationName != nil {
			gateway[dlLocationName] = *instance.LocationName
		}
		if instance.LocationDisplayName != nil {
			gateway[dlLocationDisplayName] = *instance.LocationDisplayName
		}
		if instance.Vlan != nil {
			gateway[dlVlan] = *instance.Vlan
		}
		if instance.Global != nil {
			gateway[dlGlobal] = *instance.Global
		}
		if instance.Port != nil {
			gateway[dlPort] = *instance.Port.ID
		}
		if instance.LinkStatus != nil {
			gateway[dlLinkStatus] = *instance.LinkStatus
		}
		if instance.LinkStatusUpdatedAt != nil {
			gateway[dlLinkStatusUpdatedAt] = instance.LinkStatusUpdatedAt.String()
		}
		if instance.CreatedAt != nil {
			gateway[dlCreatedAt] = instance.CreatedAt.String()
		}
		if instance.ResourceGroup != nil {
			rg := instance.ResourceGroup
			gateway[dlResourceGroup] = *rg.ID
		}
		if instance.DefaultExportRouteFilter != nil {
			gateway[dlDefault_export_route_filter] = *instance.DefaultExportRouteFilter
		}
		if instance.DefaultImportRouteFilter != nil {
			gateway[dlDefault_import_route_filter] = *instance.DefaultImportRouteFilter
		}

		//Show the BFD Config parameters if set
		if instance.BfdConfig != nil {
			if instance.BfdConfig.Interval != nil {
				gateway[dlBfdInterval] = *instance.BfdConfig.Interval
			}

			if instance.BfdConfig.Multiplier != nil {
				gateway[dlBfdMultiplier] = *instance.BfdConfig.Multiplier
			}

			if instance.BfdConfig.BfdStatus != nil {
				gateway[dlBfdStatus] = *instance.BfdConfig.BfdStatus
			}

			if instance.BfdConfig.BfdStatusUpdatedAt != nil {
				gateway[dlBfdStatusUpdatedAt] = instance.BfdConfig.BfdStatusUpdatedAt.String()
			}
		}

		asPrependList := make([]map[string]interface{}, 0)
		if len(instance.AsPrepends) > 0 {
			for _, asPrepend := range instance.AsPrepends {
				asPrependItem := map[string]interface{}{}
				asPrependItem[dlResourceId] = asPrepend.ID
				asPrependItem[dlLength] = asPrepend.Length
				asPrependItem[dlPrefix] = asPrepend.Prefix
				asPrependItem[dlSpecificPrefixes] = asPrepend.SpecificPrefixes
				asPrependItem[dlPolicy] = asPrepend.Policy
				asPrependItem[dlCreatedAt] = asPrepend.CreatedAt.String()
				asPrependItem[dlUpdatedAt] = asPrepend.UpdatedAt.String()

				asPrependList = append(asPrependList, asPrependItem)
			}

		}
		gateway[dlAsPrepends] = asPrependList

		dtype := *instance.Type
		if dtype == "dedicated" {
			if instance.MacsecConfig != nil {
				macsecList := make([]map[string]interface{}, 0)
				currentMacSec := map[string]interface{}{}
				// Construct an instance of the GatewayMacsecConfigTemplate model
				gatewayMacsecConfigTemplateModel := instance.MacsecConfig
				if gatewayMacsecConfigTemplateModel.Active != nil {
					currentMacSec[dlActive] = *gatewayMacsecConfigTemplateModel.Active
				}
				if gatewayMacsecConfigTemplateModel.ActiveCak != nil {
					if gatewayMacsecConfigTemplateModel.ActiveCak.Crn != nil {
						currentMacSec[dlActiveCak] = *gatewayMacsecConfigTemplateModel.ActiveCak.Crn
					}
				}
				if gatewayMacsecConfigTemplateModel.PrimaryCak != nil {
					currentMacSec[dlPrimaryCak] = *gatewayMacsecConfigTemplateModel.PrimaryCak.Crn
				}
				if gatewayMacsecConfigTemplateModel.FallbackCak != nil {
					if gatewayMacsecConfigTemplateModel.FallbackCak.Crn != nil {
						currentMacSec[dlFallbackCak] = *gatewayMacsecConfigTemplateModel.FallbackCak.Crn
					}
				}
				if gatewayMacsecConfigTemplateModel.SakExpiryTime != nil {
					currentMacSec[dlSakExpiryTime] = *gatewayMacsecConfigTemplateModel.SakExpiryTime
				}
				if gatewayMacsecConfigTemplateModel.SecurityPolicy != nil {
					currentMacSec[dlSecurityPolicy] = *gatewayMacsecConfigTemplateModel.SecurityPolicy
				}
				if gatewayMacsecConfigTemplateModel.WindowSize != nil {
					currentMacSec[dlWindowSize] = *gatewayMacsecConfigTemplateModel.WindowSize
				}
				if gatewayMacsecConfigTemplateModel.CipherSuite != nil {
					currentMacSec[dlCipherSuite] = *gatewayMacsecConfigTemplateModel.CipherSuite
				}
				if gatewayMacsecConfigTemplateModel.ConfidentialityOffset != nil {
					currentMacSec[dlConfidentialityOffset] = *gatewayMacsecConfigTemplateModel.ConfidentialityOffset
				}
				if gatewayMacsecConfigTemplateModel.CryptographicAlgorithm != nil {
					currentMacSec[dlCryptographicAlgorithm] = *gatewayMacsecConfigTemplateModel.CryptographicAlgorithm
				}
				if gatewayMacsecConfigTemplateModel.KeyServerPriority != nil {
					currentMacSec[dlKeyServerPriority] = *gatewayMacsecConfigTemplateModel.KeyServerPriority
				}
				if gatewayMacsecConfigTemplateModel.Status != nil {
					currentMacSec[dlMacSecConfigStatus] = *gatewayMacsecConfigTemplateModel.Status
				}
				macsecList = append(macsecList, currentMacSec)
				gateway[dlMacSecConfig] = macsecList
			}
		}
		if instance.ChangeRequest != nil {
			gatewayChangeRequestIntf := instance.ChangeRequest
			gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
			gateway[dlChangeRequest] = *gatewayChangeRequest.Type
		}

		if instance.AuthenticationKey != nil {
			gateway[dlAuthenticationKey] = *instance.AuthenticationKey.Crn
		}

		if instance.ConnectionMode != nil {
			gateway[dlConnectionMode] = *instance.ConnectionMode
		}

		gateways = append(gateways, gateway)
	}
	d.SetId(dataSourceIBMDLGatewaysID(d))
	d.Set(dlGateways, gateways)
	return nil
}

// dataSourceIBMDLGatewaysID returns a reasonable ID for a direct link gateways list.
func dataSourceIBMDLGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
