package nutanix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fc "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixFCImagedClustersList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFCImagedClustersListRead,
		Schema: map[string]*schema.Schema{
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archived": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"imaged_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"archived": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cluster_external_ip": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"imaged_node_uuid_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"common_network_settings": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cvm_dns_servers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"hypervisor_dns_servers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cvm_ntp_servers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"hypervisor_ntp_servers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"storage_node_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"foundation_init_node_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"foundation_init_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blocks": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"block_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"nodes": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cvm_ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cvm_vlan_id": {
																Type:     schema.TypeInt,
																Computed: true,
																Optional: true,
															},
															"hardware_attributes_override": {
																Type:     schema.TypeMap,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"fc_imaged_node_uuid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hypervisor": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hypervisor_hostname": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"hypervisor_ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"image_now": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"ipmi_ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"ipv6_address": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"node_position": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"node_serial": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"clusters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_external_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cluster_init_now": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"cluster_init_successful": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"cluster_members": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"cluster_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cvm_dns_servers": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cvm_ntp_servers": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"redundancy_factor": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"timezone": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"cvm_gateway": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvm_netmask": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns_servers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hyperv_product_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hyperv_sku": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor_gateway": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor_iso_url": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"hypervisor_isos": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hypervisor_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sha256sum": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"hypervisor_netmask": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipmi_gateway": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipmi_netmask": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nos_package_url": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sha256sum": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"cluster_status": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_creation_started": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"intent_picked_up": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"imaging_stopped": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"node_progress_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"imaged_node_uuid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"imaging_stopped": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"intent_picked_up": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"percent_complete": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"message_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"aggregate_percent_complete": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"current_foundation_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_progress_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"percent_complete": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"message_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"foundation_session_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"destroyed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"imaged_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixFCImagedClustersListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral

	req := &fc.ImagedClustersListInput{}
	if len, lenok := d.GetOk("length"); lenok {
		req.Length = utils.IntPtr(len.(int))
	}
	if offset, offok := d.GetOk("offset"); offok {
		req.Offset = utils.IntPtr(offset.(int))
	}

	if filter, fok := d.GetOk("filters"); fok {
		filt := &fc.ImagedClustersListFilter{}
		filter := filter.([]interface{})[0].(map[string]interface{})
		filt.Archived = utils.BoolPtr(filter["archived"].(bool))
		req.Filters = filt
	}

	resp, err := conn.Service.ListImagedClusters(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	imagedClusters := flattenImagedClusters(resp.ImagedClusters)

	d.Set("imaged_clusters", imagedClusters)

	if resp.Metadata != nil {
		metalist := make([]map[string]interface{}, 0)
		meta := make(map[string]interface{})
		meta["length"] = (resp.Metadata.Length)
		meta["offset"] = (resp.Metadata.Offset)
		meta["total_matches"] = (resp.Metadata.TotalMatches)

		metalist = append(metalist, meta)
		d.Set("metadata", metalist)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenImagedClusters(imgCluster []*fc.ImagedClusterDetails) []map[string]interface{} {
	imgClusterList := make([]map[string]interface{}, len(imgCluster))
	if len(imgCluster) > 0 {
		for k, v := range imgCluster {
			imgClusterList[k] = map[string]interface{}{
				"current_time":              v.CurrentTime,
				"archived":                  v.Archived,
				"cluster_external_ip":       v.ClusterExternalIP,
				"imaged_node_uuid_list":     utils.StringValueSlice(v.ImagedNodeUUIDList),
				"common_network_settings":   flattenFCCommonNetworkSettings(v.CommonNetworkSettings),
				"storage_node_count":        v.StorageNodeCount,
				"redundancy_factor":         v.RedundancyFactor,
				"foundation_init_node_uuid": v.FoundationInitNodeUUID,
				"workflow_type":             v.WorkflowType,
				"cluster_name":              v.ClusterName,
				"foundation_init_config":    flattenFCFoundationInitConfig(v.FoundationInitConfig),
				"cluster_status":            flattenClusterStatus(v.ClusterStatus),
				"cluster_size":              v.ClusterSize,
				"destroyed":                 v.Destroyed,
				"created_timestamp":         v.CreatedTimestamp,
				"imaged_cluster_uuid":       v.ImagedClusterUUID,
			}
		}
	}
	return imgClusterList
}

func flattenFCCommonNetworkSettings(cnet *fc.CommonNetworkSettings) []interface{} {
	references := make([]interface{}, 0)
	if cnet != nil {
		reference := make(map[string]interface{})
		reference["cvm_dns_servers"] = utils.StringSlice(cnet.CvmDNSServers)
		reference["hypervisor_dns_servers"] = utils.StringSlice(cnet.HypervisorDNSServers)
		reference["cvm_ntp_servers"] = utils.StringSlice(cnet.CvmNtpServers)
		reference["hypervisor_ntp_servers"] = utils.StringSlice(cnet.HypervisorNtpServers)

		references = append(references, reference)
	}
	return references
}

func flattenClusterStatus(cs *fc.ClusterStatus) []interface{} {
	cstatus := make([]interface{}, 0)
	if cs != nil {
		csList := make(map[string]interface{})
		csList["intent_picked_up"] = utils.BoolValue(cs.IntentPickedUp)
		csList["cluster_creation_started"] = utils.BoolValue(cs.ClusterCreationStarted)
		csList["imaging_stopped"] = utils.BoolValue(cs.ImagingStopped)
		csList["aggregate_percent_complete"] = utils.Float64Value(cs.AggregatePercentComplete)
		csList["current_foundation_ip"] = utils.StringValue(cs.CurrentFoundationIP)
		csList["foundation_session_id"] = utils.StringValue(cs.FoundationSessionID)
		csList["node_progress_details"] = flattenNodeProgressDetails(cs.NodeProgressDetails)
		csList["cluster_progress_details"] = flattenClusterProgressDetails(cs.ClusterProgressDetails)

		cstatus = append(cstatus, csList)
	}
	return cstatus
}

func flattenNodeProgressDetails(np []*fc.NodeProgressDetail) []map[string]interface{} {
	npd := make([]map[string]interface{}, len(np))

	if len(np) > 0 {
		for k, v := range np {
			n := make(map[string]interface{})

			n["status"] = v.Status
			n["imaged_node_uuid"] = v.ImagedNodeUUID
			n["imaging_stopped"] = v.ImagingStopped
			n["intent_picked_up"] = v.IntentPickedUp
			n["percent_complete"] = v.PercentComplete
			n["message_list"] = utils.StringValueSlice(v.MessageList)

			npd[k] = n
		}
	}
	return npd
}

func flattenClusterProgressDetails(cp *fc.ClusterProgressDetails) []interface{} {
	cpDetails := make([]interface{}, 0)
	if cp != nil {
		cpd := make(map[string]interface{})
		cpd["cluster_name"] = utils.StringValue(cp.ClusterName)
		cpd["status"] = utils.StringValue(cp.Status)
		cpd["percent_complete"] = utils.Float64Value(cp.PercentComplete)
		cpd["message_list"] = utils.StringValueSlice(cp.MessageList)

		cpDetails = append(cpDetails, cpd)
	}
	return cpDetails
}

func flattenFCFoundationInitConfig(fci *fc.FoundationInitConfig) []interface{} {
	fciDetails := make([]interface{}, 0)
	if fci != nil {
		fcic := make(map[string]interface{})
		fcic["blocks"] = flattenFCBlock(fci.Blocks)
		fcic["clusters"] = flattenCluster(fci.Clusters)
		fcic["cvm_gateway"] = fci.CvmGateway
		fcic["cvm_netmask"] = fci.CvmNetmask
		fcic["dns_servers"] = fci.DNSServers
		fcic["hyperv_product_key"] = fci.HypervProductKey
		fcic["hyperv_sku"] = fci.HypervSku
		fcic["hypervisor_gateway"] = fci.HypervisorGateway
		fcic["hypervisor_iso_url"] = flattenHypervisorIsoURL(fci.HypervisorIsoURL)
		fcic["hypervisor_isos"] = flattenFCHypervisorIsos(fci.HypervisorIsos)
		fcic["hypervisor_netmask"] = fci.HypervisorNetmask
		fcic["ipmi_gateway"] = fci.IpmiGateway
		fcic["ipmi_netmask"] = fci.IpmiNetmask
		fcic["nos_package_url"] = flattenNosPackage(fci.NosPackageURL)

		fciDetails = append(fciDetails, fcic)
	}

	return fciDetails
}

func flattenFCBlock(fb []*fc.Blocks) []map[string]interface{} {
	res := make([]map[string]interface{}, len(fb))
	if len(fb) > 0 {
		for k, v := range fb {
			re := make(map[string]interface{})

			re["block_id"] = v.BlockID
			re["nodes"] = flattenNodes(v.Nodes)

			res[k] = re
		}
	}
	return res
}

func flattenNodes(fn []*fc.Nodes) []map[string]interface{} {
	res := make([]map[string]interface{}, len(fn))

	if len(fn) > 0 {
		for k, v := range fn {
			re := make(map[string]interface{})

			re["cvm_ip"] = v.CvmIP
			re["cvm_vlan_id"] = v.CvmVlanID
			re["fc_imaged_node_uuid"] = v.FcImagedNodeUUID
			re["hypervisor"] = v.Hypervisor
			re["hypervisor_hostname"] = v.HypervisorHostname
			re["hypervisor_ip"] = v.HypervisorIP
			re["image_now"] = v.ImageNow
			re["ipmi_ip"] = v.IpmiIP
			re["ipv6_address"] = v.IPv6Address
			re["node_position"] = v.NodePosition
			re["node_serial"] = v.NodeSerial
			re["hardware_attributes_override"] = flattenHardwareAttributes(v.HardwareAttributesOverride)

			res[k] = re
		}
	}
	return res
}

func flattenCluster(fc []*fc.Clusters) []map[string]interface{} {
	res := make([]map[string]interface{}, len(fc))

	if len(fc) > 0 {
		for k, v := range fc {
			re := make(map[string]interface{})

			re["cluster_external_ip"] = v.ClusterExternalIP
			re["cluster_init_now"] = v.ClusterInitNow
			re["cluster_init_successful"] = v.ClusterInitSuccessful
			re["cluster_members"] = utils.StringValueSlice(v.ClusterMembers)
			re["cluster_name"] = v.ClusterName
			re["cvm_dns_servers"] = v.CvmDNSServers
			re["cvm_ntp_servers"] = v.CvmNtpServers
			re["redundancy_factor"] = v.RedundancyFactor
			re["timezone"] = v.TimeZone

			res[k] = re
		}
	}
	return res
}

func flattenHypervisorIsoURL(hiso *fc.HypervisorIso) map[string]interface{} {
	return map[string]interface{}{
		"hypervisor_type": hiso.HypervisorType,
		"sha256sum":       hiso.Sha256sum,
		"url":             hiso.URL,
	}
}

func flattenFCHypervisorIsos(hisos []*fc.HypervisorIso) (res []map[string]interface{}) {
	if len(hisos) > 0 {
		for _, k := range hisos {
			re := make(map[string]interface{})

			re["hypervisor_type"] = k.HypervisorType
			re["sha256sum"] = k.Sha256sum
			re["url"] = k.URL

			res = append(res, re)
		}
	}
	return res
}

func flattenNosPackage(fnos *fc.NosPackageURL) (res []interface{}) {
	if fnos != nil {
		re := make(map[string]interface{})
		re["sha256sum"] = fnos.Sha256sum
		re["url"] = fnos.URL

		res = append(res, re)
	}
	return res
}

func flattenHardwareAttributes(pr map[string]interface{}) map[string]interface{} {
	mm := make(map[string]interface{})
	for k, v := range pr {
		mm[k] = fmt.Sprint(v)
	}
	return mm
}
