// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vmware

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vmware-go-sdk/vmwarev1"
)

func DataSourceIbmVmaasVdc() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmVmaasVdcRead,

		Schema: map[string]*schema.Schema{
			"vmaas_vdc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique ID for a specified virtual data center.",
			},
			"accept_language": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Language.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of this virtual data center (VDC).",
			},
			"provisioned_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is provisioned and available to use.",
			},
			"cpu": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The vCPU usage limit on the virtual data center (VDC). Supported for VDCs deployed on a multitenant Cloud Director site. This property is applicable when the resource pool type is reserved.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique ID for the virtual data center (VDC) in IBM Cloud.",
			},
			"deleted_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is deleted.",
			},
			"director_site": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Director site in which to deploy the virtual data center (VDC).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for the Cloud Director site.",
						},
						"pvdc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource pool within the Director Site in which to deploy the virtual data center (VDC).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for the resource pool.",
									},
									"provider_type": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Determines how resources are made available to the virtual data center (VDC). Required for VDCs deployed on a multitenant Cloud Director site.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the resource pool type.",
												},
											},
										},
									},
								},
							},
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the VMware Cloud Director tenant portal where this virtual data center (VDC) can be managed.",
						},
					},
				},
			},
			"edges": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VMware NSX-T networking edges deployed on the virtual data center (VDC). NSX-T edges are used for bridging virtualization networking to the physical public-internet and IBM private networking.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for the edge.",
						},
						"public_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The public IP addresses assigned to the edge.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"size": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The size of the edge.The size can be specified only for performance edges. Larger sizes require more capacity from the Cloud Director site in which the virtual data center (VDC) was created to be deployed.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Determines the state of the edge.",
						},
						"transit_gateways": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connected IBM Transit Gateways.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for an IBM Transit Gateway.",
									},
									"connections": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IBM Transit Gateway connections.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The autogenerated name for this connection.",
												},
												"transit_gateway_connection_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name of the connection created on the IBM Transit Gateway.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Determines the state of the connection.",
												},
												"local_gateway_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Local gateway IP address for the connection.",
												},
												"remote_gateway_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Remote gateway IP address for the connection.",
												},
												"local_tunnel_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Local tunnel IP address for the connection.",
												},
												"remote_tunnel_ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Remote tunnel IP address for the connection.",
												},
												"local_bgp_asn": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Local network BGP ASN for the connection.",
												},
												"remote_bgp_asn": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remote network BGP ASN for the connection.",
												},
												"network_account_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the account that owns the connected network.",
												},
												"network_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of the network that is connected through this connection. Only \"unbound_gre_tunnel\" is supported.",
												},
												"base_network_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of the network that the unbound GRE tunnel is targeting. Only \"classic\" is supported.",
												},
												"zone": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The location of the connection.",
												},
											},
										},
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Determines the state of the IBM Transit Gateway based on its connections.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of edge to be deployed.Efficiency edges allow for multiple VDCs to share some edge resources. Performance edges do not share resources between VDCs.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edge version.",
						},
					},
				},
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about why the request to create the virtual data center (VDC) cannot be completed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An error code specific to the error encountered.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A message that describes why the error ocurred.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL that links to a page with more information about this error.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human readable ID for the virtual data center (VDC).",
			},
			"ordered_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time that the virtual data center (VDC) is ordered.",
			},
			"org_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the VMware Cloud Director organization that contains this virtual data center (VDC). VMware Cloud Director organizations are used to create strong boundaries between VDCs. There is a complete isolation of user administration, networking, workloads, and VMware Cloud Director catalogs between different Director organizations.",
			},
			"ram": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The RAM usage limit on the virtual data center (VDC) in GB (1024^3 bytes). Supported for VDCs deployed on a multitenant Cloud Director site. This property is applicable when the resource pool type is reserved.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines the state of the virtual data center.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether this virtual data center is in a single-tenant or multitenant Cloud Director site.",
			},
			"fast_provisioning_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether this virtual data center has fast provisioning enabled or not.",
			},
			"rhel_byol": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the RHEL VMs will be using the license from IBM or the customer will use their own license (BYOL).",
			},
			"windows_byol": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the Microsoft Windows VMs will be using the license from IBM or the customer will use their own license (BYOL).",
			},
		},
	}
}

func dataSourceIbmVmaasVdcRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getVdcOptions := &vmwarev1.GetVdcOptions{}

	getVdcOptions.SetID(d.Get("vmaas_vdc_id").(string))
	if _, ok := d.GetOk("accept_language"); ok {
		getVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	vDC, response, err := vmwareClient.GetVdcWithContext(context, getVdcOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVdcWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVdcWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getVdcOptions.ID))

	if err = d.Set("href", vDC.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("provisioned_at", flex.DateTimeToString(vDC.ProvisionedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting provisioned_at: %s", err))
	}

	if err = d.Set("cpu", flex.IntValue(vDC.Cpu)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cpu: %s", err))
	}

	if err = d.Set("crn", vDC.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("deleted_at", flex.DateTimeToString(vDC.DeletedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting deleted_at: %s", err))
	}

	directorSite := []map[string]interface{}{}
	if vDC.DirectorSite != nil {
		modelMap, err := dataSourceIbmVmaasVdcVDCDirectorSiteToMap(vDC.DirectorSite)
		if err != nil {
			return diag.FromErr(err)
		}
		directorSite = append(directorSite, modelMap)
	}
	if err = d.Set("director_site", directorSite); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting director_site %s", err))
	}

	edges := []map[string]interface{}{}
	if vDC.Edges != nil {
		for _, modelItem := range vDC.Edges {
			modelMap, err := dataSourceIbmVmaasVdcEdgeToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			edges = append(edges, modelMap)
		}
	}
	if err = d.Set("edges", edges); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting edges %s", err))
	}

	statusReasons := []map[string]interface{}{}
	if vDC.StatusReasons != nil {
		for _, modelItem := range vDC.StatusReasons {
			modelMap, err := dataSourceIbmVmaasVdcStatusReasonToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			statusReasons = append(statusReasons, modelMap)
		}
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_reasons %s", err))
	}

	if err = d.Set("name", vDC.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("ordered_at", flex.DateTimeToString(vDC.OrderedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ordered_at: %s", err))
	}

	if err = d.Set("org_name", vDC.OrgName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting org_name: %s", err))
	}

	if err = d.Set("ram", flex.IntValue(vDC.Ram)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ram: %s", err))
	}

	if err = d.Set("status", vDC.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if err = d.Set("type", vDC.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("fast_provisioning_enabled", vDC.FastProvisioningEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting fast_provisioning_enabled: %s", err))
	}

	if err = d.Set("rhel_byol", vDC.RhelByol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rhel_byol: %s", err))
	}

	if err = d.Set("windows_byol", vDC.WindowsByol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting windows_byol: %s", err))
	}

	return nil
}

func dataSourceIbmVmaasVdcVDCDirectorSiteToMap(model *vmwarev1.VDCDirectorSite) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	pvdcMap, err := dataSourceIbmVmaasVdcDirectorSitePVDCToMap(model.Pvdc)
	if err != nil {
		return modelMap, err
	}
	modelMap["pvdc"] = []map[string]interface{}{pvdcMap}
	modelMap["url"] = model.URL
	return modelMap, nil
}

func dataSourceIbmVmaasVdcDirectorSitePVDCToMap(model *vmwarev1.DirectorSitePVDC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	if model.ProviderType != nil {
		providerTypeMap, err := dataSourceIbmVmaasVdcVDCProviderTypeToMap(model.ProviderType)
		if err != nil {
			return modelMap, err
		}
		modelMap["provider_type"] = []map[string]interface{}{providerTypeMap}
	}
	return modelMap, nil
}

func dataSourceIbmVmaasVdcVDCProviderTypeToMap(model *vmwarev1.VDCProviderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmVmaasVdcEdgeToMap(model *vmwarev1.Edge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["public_ips"] = model.PublicIps
	if model.Size != nil {
		modelMap["size"] = model.Size
	}
	modelMap["status"] = model.Status
	transitGateways := []map[string]interface{}{}
	for _, transitGatewaysItem := range model.TransitGateways {
		transitGatewaysItemMap, err := dataSourceIbmVmaasVdcTransitGatewayToMap(&transitGatewaysItem)
		if err != nil {
			return modelMap, err
		}
		transitGateways = append(transitGateways, transitGatewaysItemMap)
	}
	modelMap["transit_gateways"] = transitGateways
	modelMap["type"] = model.Type
	modelMap["version"] = model.Version
	return modelMap, nil
}

func dataSourceIbmVmaasVdcTransitGatewayToMap(model *vmwarev1.TransitGateway) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	connections := []map[string]interface{}{}
	for _, connectionsItem := range model.Connections {
		connectionsItemMap, err := dataSourceIbmVmaasVdcTransitGatewayConnectionToMap(&connectionsItem)
		if err != nil {
			return modelMap, err
		}
		connections = append(connections, connectionsItemMap)
	}
	modelMap["connections"] = connections
	modelMap["status"] = model.Status
	return modelMap, nil
}

func dataSourceIbmVmaasVdcTransitGatewayConnectionToMap(model *vmwarev1.TransitGatewayConnection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.TransitGatewayConnectionName != nil {
		modelMap["transit_gateway_connection_name"] = model.TransitGatewayConnectionName
	}
	modelMap["status"] = model.Status
	if model.LocalGatewayIp != nil {
		modelMap["local_gateway_ip"] = model.LocalGatewayIp
	}
	if model.RemoteGatewayIp != nil {
		modelMap["remote_gateway_ip"] = model.RemoteGatewayIp
	}
	if model.LocalTunnelIp != nil {
		modelMap["local_tunnel_ip"] = model.LocalTunnelIp
	}
	if model.RemoteTunnelIp != nil {
		modelMap["remote_tunnel_ip"] = model.RemoteTunnelIp
	}
	if model.LocalBgpAsn != nil {
		modelMap["local_bgp_asn"] = flex.IntValue(model.LocalBgpAsn)
	}
	if model.RemoteBgpAsn != nil {
		modelMap["remote_bgp_asn"] = flex.IntValue(model.RemoteBgpAsn)
	}
	modelMap["network_account_id"] = model.NetworkAccountID
	modelMap["network_type"] = model.NetworkType
	modelMap["base_network_type"] = model.BaseNetworkType
	modelMap["zone"] = model.Zone
	return modelMap, nil
}

func dataSourceIbmVmaasVdcStatusReasonToMap(model *vmwarev1.StatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}
