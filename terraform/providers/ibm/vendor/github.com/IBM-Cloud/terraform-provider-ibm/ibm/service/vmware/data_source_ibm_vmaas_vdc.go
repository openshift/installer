// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package vmware

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
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
						"private_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The private IP addresses assigned to the edge.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"private_only": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the edge is private only. The default value is True for a private Cloud Director site and False for a public Cloud Director site.",
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
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region where the IBM Transit Gateway is deployed.",
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
			"org_href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the organization that owns the VDC.",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_vmaas_vdc", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVdcOptions := &vmwarev1.GetVdcOptions{}

	getVdcOptions.SetID(d.Get("vmaas_vdc_id").(string))
	if _, ok := d.GetOk("accept_language"); ok {
		getVdcOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	vDC, _, err := vmwareClient.GetVdcWithContext(context, getVdcOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVdcWithContext failed: %s", err.Error()), "(Data) ibm_vmaas_vdc", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getVdcOptions.ID)

	if err = d.Set("href", vDC.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-href").GetDiag()
	}

	if !core.IsNil(vDC.ProvisionedAt) {
		if err = d.Set("provisioned_at", flex.DateTimeToString(vDC.ProvisionedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting provisioned_at: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-provisioned_at").GetDiag()
		}
	}

	if !core.IsNil(vDC.Cpu) {
		if err = d.Set("cpu", flex.IntValue(vDC.Cpu)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cpu: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-cpu").GetDiag()
		}
	}

	if err = d.Set("crn", vDC.Crn); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-crn").GetDiag()
	}

	if !core.IsNil(vDC.DeletedAt) {
		if err = d.Set("deleted_at", flex.DateTimeToString(vDC.DeletedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting deleted_at: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-deleted_at").GetDiag()
		}
	}

	directorSite := []map[string]interface{}{}
	directorSiteMap, err := DataSourceIbmVmaasVdcVDCDirectorSiteToMap(vDC.DirectorSite)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_vmaas_vdc", "read", "director_site-to-map").GetDiag()
	}
	directorSite = append(directorSite, directorSiteMap)
	if err = d.Set("director_site", directorSite); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting director_site: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-director_site").GetDiag()
	}

	edges := []map[string]interface{}{}
	for _, edgesItem := range vDC.Edges {
		edgesItemMap, err := DataSourceIbmVmaasVdcEdgeToMap(&edgesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_vmaas_vdc", "read", "edges-to-map").GetDiag()
		}
		edges = append(edges, edgesItemMap)
	}
	if err = d.Set("edges", edges); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting edges: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-edges").GetDiag()
	}

	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range vDC.StatusReasons {
		statusReasonsItemMap, err := DataSourceIbmVmaasVdcStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_vmaas_vdc", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_reasons: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-status_reasons").GetDiag()
	}

	if err = d.Set("name", vDC.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-name").GetDiag()
	}

	if err = d.Set("ordered_at", flex.DateTimeToString(vDC.OrderedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ordered_at: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-ordered_at").GetDiag()
	}

	if err = d.Set("org_href", vDC.OrgHref); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting org_href: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-org_href").GetDiag()
	}

	if err = d.Set("org_name", vDC.OrgName); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting org_name: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-org_name").GetDiag()
	}

	if !core.IsNil(vDC.Ram) {
		if err = d.Set("ram", flex.IntValue(vDC.Ram)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ram: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-ram").GetDiag()
		}
	}

	if err = d.Set("status", vDC.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-status").GetDiag()
	}

	if err = d.Set("type", vDC.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-type").GetDiag()
	}

	if err = d.Set("fast_provisioning_enabled", vDC.FastProvisioningEnabled); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting fast_provisioning_enabled: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-fast_provisioning_enabled").GetDiag()
	}

	if err = d.Set("rhel_byol", vDC.RhelByol); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting rhel_byol: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-rhel_byol").GetDiag()
	}

	if err = d.Set("windows_byol", vDC.WindowsByol); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting windows_byol: %s", err), "(Data) ibm_vmaas_vdc", "read", "set-windows_byol").GetDiag()
	}

	return nil
}

func DataSourceIbmVmaasVdcVDCDirectorSiteToMap(model *vmwarev1.VDCDirectorSite) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	pvdcMap, err := DataSourceIbmVmaasVdcDirectorSitePVDCToMap(model.Pvdc)
	if err != nil {
		return modelMap, err
	}
	modelMap["pvdc"] = []map[string]interface{}{pvdcMap}
	modelMap["url"] = *model.URL
	return modelMap, nil
}

func DataSourceIbmVmaasVdcDirectorSitePVDCToMap(model *vmwarev1.DirectorSitePVDC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	if model.ProviderType != nil {
		providerTypeMap, err := DataSourceIbmVmaasVdcVDCProviderTypeToMap(model.ProviderType)
		if err != nil {
			return modelMap, err
		}
		modelMap["provider_type"] = []map[string]interface{}{providerTypeMap}
	}
	return modelMap, nil
}

func DataSourceIbmVmaasVdcVDCProviderTypeToMap(model *vmwarev1.VDCProviderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIbmVmaasVdcEdgeToMap(model *vmwarev1.Edge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["public_ips"] = model.PublicIps
	modelMap["private_ips"] = model.PrivateIps
	if model.PrivateOnly != nil {
		modelMap["private_only"] = *model.PrivateOnly
	}
	modelMap["size"] = *model.Size
	modelMap["status"] = *model.Status
	transitGateways := []map[string]interface{}{}
	for _, transitGatewaysItem := range model.TransitGateways {
		transitGatewaysItemMap, err := DataSourceIbmVmaasVdcTransitGatewayToMap(&transitGatewaysItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		transitGateways = append(transitGateways, transitGatewaysItemMap)
	}
	modelMap["transit_gateways"] = transitGateways
	modelMap["type"] = *model.Type
	modelMap["version"] = *model.Version
	return modelMap, nil
}

func DataSourceIbmVmaasVdcTransitGatewayToMap(model *vmwarev1.TransitGateway) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	connections := []map[string]interface{}{}
	for _, connectionsItem := range model.Connections {
		connectionsItemMap, err := DataSourceIbmVmaasVdcTransitGatewayConnectionToMap(&connectionsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		connections = append(connections, connectionsItemMap)
	}
	modelMap["connections"] = connections
	modelMap["status"] = *model.Status
	modelMap["region"] = *model.Region
	return modelMap, nil
}

func DataSourceIbmVmaasVdcTransitGatewayConnectionToMap(model *vmwarev1.TransitGatewayConnection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	if model.TransitGatewayConnectionName != nil {
		modelMap["transit_gateway_connection_name"] = *model.TransitGatewayConnectionName
	}
	modelMap["status"] = *model.Status
	if model.LocalGatewayIp != nil {
		modelMap["local_gateway_ip"] = *model.LocalGatewayIp
	}
	if model.RemoteGatewayIp != nil {
		modelMap["remote_gateway_ip"] = *model.RemoteGatewayIp
	}
	if model.LocalTunnelIp != nil {
		modelMap["local_tunnel_ip"] = *model.LocalTunnelIp
	}
	if model.RemoteTunnelIp != nil {
		modelMap["remote_tunnel_ip"] = *model.RemoteTunnelIp
	}
	if model.LocalBgpAsn != nil {
		modelMap["local_bgp_asn"] = flex.IntValue(model.LocalBgpAsn)
	}
	if model.RemoteBgpAsn != nil {
		modelMap["remote_bgp_asn"] = flex.IntValue(model.RemoteBgpAsn)
	}
	modelMap["network_account_id"] = *model.NetworkAccountID
	modelMap["network_type"] = *model.NetworkType
	modelMap["base_network_type"] = *model.BaseNetworkType
	modelMap["zone"] = *model.Zone
	return modelMap, nil
}

func DataSourceIbmVmaasVdcStatusReasonToMap(model *vmwarev1.StatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
