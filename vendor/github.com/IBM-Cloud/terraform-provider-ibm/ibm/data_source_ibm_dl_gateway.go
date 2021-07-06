// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlGateway        = "gateway"
	dlSecurityPolicy = "security_policy"
	dlActiveCak      = "active_cak"
)

func dataSourceIBMDLGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLGatewayRead,
		Schema: map[string]*schema.Schema{
			dlName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique user-defined name for this gateway",
				ValidateFunc: InvokeValidator("ibm_dl_gateway", dlName),
			},

			dlGatewaysVirtualConnections: {
				Type:        schema.TypeList,
				Description: "Collection of direct link gateway virtual connections",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlVCCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time resource was created",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual connection",
						},
						dlVCStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the virtual connection.Possible values: [pending,attached,approval_pending,rejected,expired,deleting,detached_by_network_pending,detached_by_network]",
						},
						dlVCNetworkAccount: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "For virtual connections across two different IBM Cloud Accounts network_account indicates the account that owns the target network.",
						},
						dlVCNetworkId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the target network. For type=vpc virtual connections this is the CRN of the target VPC. This field does not apply to type=classic connections.",
						},
						dlVCType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of virtual connection. (classic,vpc)",
						},
						dlVCName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this virtual connection. Virtualconnection names are unique within a gateway. This is the name of thevirtual connection itself, the network being connected may have its ownname attribute",
						},
					},
				},
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
							Description: "Desired primary connectivity association key.",
						},
						dlFallbackCak: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fallback connectivity association key.",
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
			dlCompletionNoticeRejectReason: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason for completion notice rejection",
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
			dlVlan: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VLAN allocated for this gateway",
			},
		},
	}
}

func dataSourceIBMDLGatewayVirtualConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := meta.(ClientSession).DirectlinkV1API()

	if err != nil {
		return err
	}
	listVcOptions := &directlinkv1.ListGatewayVirtualConnectionsOptions{}
	dlGatewayId := d.Id()
	listVcOptions.SetGatewayID(dlGatewayId)
	listGatewayVirtualConnections, response, err := directLink.ListGatewayVirtualConnections(listVcOptions)
	if err != nil {
		return fmt.Errorf("Error while listing directlink gateway's virtual connections XXX %s\n%s", err, response)
	}
	gatewayVCs := make([]map[string]interface{}, 0)
	for _, instance := range listGatewayVirtualConnections.VirtualConnections {
		gatewayVC := map[string]interface{}{}

		if instance.ID != nil {
			gatewayVC[ID] = *instance.ID
		}
		if instance.Name != nil {
			gatewayVC[dlVCName] = *instance.Name
		}
		if instance.Type != nil {
			gatewayVC[dlVCType] = *instance.Type
		}
		if instance.NetworkAccount != nil {
			gatewayVC[dlVCNetworkAccount] = *instance.NetworkAccount
		}
		if instance.NetworkID != nil {
			gatewayVC[dlVCNetworkId] = *instance.NetworkID
		}
		if instance.CreatedAt != nil {
			gatewayVC[dlVCCreatedAt] = instance.CreatedAt.String()

		}
		if instance.Status != nil {
			gatewayVC[dlVCStatus] = *instance.Status
		}

		gatewayVCs = append(gatewayVCs, gatewayVC)
	}
	d.SetId(dlGatewayId)

	d.Set(dlGatewaysVirtualConnections, gatewayVCs)
	return nil
}
func dataSourceIBMDLGatewayRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	dlGatewayName := d.Get(dlName).(string)

	if err != nil {
		return err
	}
	listGatewaysOptionsModel := &directlinkv1.ListGatewaysOptions{}
	listGateways, response, err := directLink.ListGateways(listGatewaysOptionsModel)
	if err != nil {
		log.Println("[WARN] Error listing dl Gateway", response, err)
		return err
	}
	var found bool

	for _, instance := range listGateways.Gateways {

		if *instance.Name == dlGatewayName {
			found = true
			if instance.ID != nil {
				d.SetId(*instance.ID)
			}
			if instance.Name != nil {
				d.Set(dlName, *instance.Name)
			}
			if instance.Crn != nil {
				d.Set(dlCrn, *instance.Crn)
			}
			if instance.BgpAsn != nil {
				d.Set(dlBgpAsn, *instance.BgpAsn)
			}
			if instance.BgpIbmCidr != nil {
				d.Set(dlBgpIbmCidr, *instance.BgpIbmCidr)
			}
			if instance.BgpIbmAsn != nil {
				d.Set(dlBgpIbmAsn, *instance.BgpIbmAsn)
			}
			if instance.Metered != nil {
				d.Set(dlMetered, *instance.Metered)
			}
			if instance.CrossConnectRouter != nil {
				d.Set(dlCrossConnectRouter, *instance.CrossConnectRouter)
			}
			if instance.BgpBaseCidr != nil {
				d.Set(dlBgpBaseCidr, *instance.BgpBaseCidr)
			}
			if instance.BgpCerCidr != nil {
				d.Set(dlBgpCerCidr, *instance.BgpCerCidr)
			}

			if instance.ProviderApiManaged != nil {
				d.Set(dlProviderAPIManaged, *instance.ProviderApiManaged)
			}
			if instance.Type != nil {
				d.Set(dlType, *instance.Type)
			}
			if instance.SpeedMbps != nil {
				d.Set(dlSpeedMbps, *instance.SpeedMbps)
			}
			if instance.OperationalStatus != nil {
				d.Set(dlOperationalStatus, *instance.OperationalStatus)
			}
			if instance.BgpStatus != nil {
				d.Set(dlBgpStatus, *instance.BgpStatus)
			}
			if instance.LocationName != nil {
				d.Set(dlLocationName, *instance.LocationName)
			}
			if instance.LocationDisplayName != nil {
				d.Set(dlLocationDisplayName, *instance.LocationDisplayName)
			}
			if instance.Vlan != nil {
				d.Set(dlVlan, *instance.Vlan)
			}
			if instance.Global != nil {
				d.Set(dlGlobal, *instance.Global)
			}
			if instance.Port != nil {
				d.Set(dlPort, *instance.Port.ID)
			}
			if instance.LinkStatus != nil {
				d.Set(dlLinkStatus, *instance.LinkStatus)
			}
			if instance.CreatedAt != nil {
				d.Set(dlCreatedAt, instance.CreatedAt.String())
			}
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
					d.Set(dlMacSecConfig, macsecList)
				}
			}
			if instance.ChangeRequest != nil {
				gatewayChangeRequestIntf := instance.ChangeRequest
				gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
				d.Set(dlChangeRequest, *gatewayChangeRequest.Type)
			}
			if instance.ResourceGroup != nil {
				rg := instance.ResourceGroup
				d.Set(dlResourceGroup, *rg.ID)
			}

		}
	}

	if !found {
		return fmt.Errorf(
			"Error Gateway with name  (%s) not found ", dlGatewayName)
	}
	return dataSourceIBMDLGatewayVirtualConnectionsRead(d, meta)
}
