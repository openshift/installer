package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudServiceMeshServiceMeshes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceMeshServiceMeshesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"running", "initial"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meshes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"intranet_api_server_endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"intranet_pilot_endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_api_server_endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_pilot_endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_server_loadbalancer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"api_server_public_eip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"pilot_public_eip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"pilot_public_loadbalancer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"edition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mesh_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_log": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"audit": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"project": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"customized_zipkin": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_locality_lb": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"include_ip_ranges": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kiali": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"opa": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"limit_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"limit_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"log_level": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"outbound_traffic_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pilot": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http10_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"trace_sampling": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"prometheus": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"external_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"use_external": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"proxy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_domain": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"limit_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"limit_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"sidecar_injector": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_injection_policy_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"enable_namespaces_by_default": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"init_cni_configuration": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"exclude_namespaces": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"limit_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"limit_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_cpu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"request_memory": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sidecar_injector_webhook_as_yaml": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"telemetry": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tracing": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"network": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitche_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_mesh_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_mesh_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudServiceMeshServiceMeshesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeServiceMeshes"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var serviceMeshNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		serviceMeshNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewServicemeshClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-01-11"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_mesh_service_meshes", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ServiceMeshes", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ServiceMeshes", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if serviceMeshNameRegex != nil && !serviceMeshNameRegex.MatchString(fmt.Sprint(item["ServiceMeshInfo"].(map[string]interface{})["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			id := item["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"]
			if _, ok := idsMap[fmt.Sprint(id)]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != item["ServiceMeshInfo"].(map[string]interface{})["State"].(string) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"clusters": object["Clusters"],
		}
		if objectArg, ok := object["ServiceMeshInfo"].(map[string]interface{}); ok {
			mapping["create_time"] = objectArg["CreationTime"]
			mapping["error_message"] = objectArg["ErrorMessage"]
			mapping["id"] = fmt.Sprint(objectArg["ServiceMeshId"])
			mapping["service_mesh_id"] = fmt.Sprint(objectArg["ServiceMeshId"])
			mapping["service_mesh_name"] = objectArg["Name"]
			mapping["status"] = objectArg["State"]
			mapping["version"] = objectArg["Version"]
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ServiceMeshInfo"].(map[string]interface{})["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"])
		servicemeshService := ServicemeshService{client}
		getResp, err := servicemeshService.DescribeServiceMeshServiceMesh(id)
		if err != nil {
			return WrapError(err)
		}
		endpointsSli := make([]map[string]interface{}, 0)
		if endpoints, ok := getResp["Endpoints"]; ok {
			endpointsMap := make(map[string]interface{})
			if endpointsArg, ok := endpoints.(map[string]interface{}); ok && len(endpointsArg) > 0 {
				endpointsMap["intranet_pilot_endpoint"] = endpointsArg["IntranetPilotEndpoint"]
				endpointsMap["public_pilot_endpoint"] = endpointsArg["PublicPilotEndpoint"]
				endpointsMap["intranet_api_server_endpoint"] = endpointsArg["IntranetApiServerEndpoint"]
				endpointsMap["public_api_server_endpoint"] = endpointsArg["PublicApiServerEndpoint"]
				endpointsSli = append(endpointsSli, endpointsMap)
				mapping["endpoints"] = endpointsSli
			}
		}
		if spec, ok := getResp["Spec"]; ok {
			if specArg, ok := spec.(map[string]interface{}); ok && len(specArg) > 0 {
				meshConfigSli := make([]map[string]interface{}, 0)
				if meshConfig, ok := specArg["MeshConfig"]; ok {
					meshConfigMap := make(map[string]interface{})
					if meshConfigArg, ok := meshConfig.(map[string]interface{}); ok && len(meshConfigArg) > 0 {
						accessLogSli := make([]map[string]interface{}, 0)
						if accessLog, ok := meshConfigArg["AccessLog"]; ok {
							if accessLogArg, ok := accessLog.(map[string]interface{}); ok && len(accessLogArg) > 0 {
								accessLogMap := make(map[string]interface{})
								accessLogMap["enabled"] = accessLogArg["Enabled"]
								accessLogSli = append(accessLogSli, accessLogMap)
							}
						}
						meshConfigMap["access_log"] = accessLogSli
						auditSli := make([]map[string]interface{}, 0)
						if audit, ok := meshConfigArg["Audit"]; ok {
							auditMap := make(map[string]interface{})
							if auditArg, ok := audit.(map[string]interface{}); ok && len(auditArg) > 0 {
								auditMap["enabled"] = auditArg["Enabled"]
								auditMap["project"] = auditArg["Project"]
								auditSli = append(auditSli, auditMap)
								meshConfigMap["audit"] = auditSli
							}
						}
						mapping["edition"] = meshConfigArg["Profile"]
						meshConfigMap["customized_zipkin"] = meshConfigArg["CustomizedZipkin"]
						meshConfigMap["enable_locality_lb"] = meshConfigArg["EnableLocalityLB"]
						meshConfigMap["include_ip_ranges"] = meshConfigArg["IncludeIPRanges"]

						kialiSli := make([]map[string]interface{}, 0)
						if kiali, ok := meshConfigArg["Kiali"]; ok {
							if kialiArg, ok := kiali.(map[string]interface{}); ok && len(kialiArg) > 0 {
								kialiMap := make(map[string]interface{})
								kialiMap["enabled"] = kialiArg["Enabled"]
								if v, ok := kialiArg["Url"]; ok {
									kialiMap["url"] = v
								}
								kialiSli = append(kialiSli, kialiMap)
							}
						}
						meshConfigMap["kiali"] = kialiSli

						opaSli := make([]map[string]interface{}, 0)
						if opa, ok := meshConfigArg["OPA"]; ok {
							opaMap := make(map[string]interface{})
							if opaArg, ok := opa.(map[string]interface{}); ok && len(opaArg) > 0 {
								opaMap["enabled"] = opaArg["Enabled"]
								opaMap["limit_cpu"] = opaArg["LimitCPU"]
								opaMap["limit_memory"] = opaArg["LimitMemory"]
								opaMap["log_level"] = opaArg["LogLevel"]
								opaMap["request_cpu"] = opaArg["RequestCPU"]
								opaMap["request_memory"] = opaArg["RequestMemory"]
							}
							opaSli = append(opaSli, opaMap)
						}
						meshConfigMap["opa"] = opaSli
						meshConfigMap["outbound_traffic_policy"] = meshConfigArg["OutboundTrafficPolicy"]

						pilotSli := make([]map[string]interface{}, 0)
						if pilot := meshConfigArg["Pilot"]; ok {
							if pilotArg, ok := pilot.(map[string]interface{}); ok && len(pilotArg) > 0 {
								pilotMap := make(map[string]interface{})
								pilotMap["http10_enabled"] = pilotArg["Http10Enabled"]
								pilotMap["trace_sampling"] = pilotArg["TraceSampling"]
								pilotSli = append(pilotSli, pilotMap)
							}
						}
						meshConfigMap["pilot"] = pilotSli

						prometheusSli := make([]map[string]interface{}, 0)
						if prometheus, ok := meshConfigArg["Prometheus"]; ok {
							if prometheusArg, ok := prometheus.(map[string]interface{}); ok && len(prometheusArg) > 0 {
								prometheusMap := make(map[string]interface{})
								prometheusMap["external_url"] = prometheusArg["ExternalUrl"]
								prometheusMap["use_external"] = prometheusArg["UseExternal"]
								prometheusSli = append(prometheusSli, prometheusMap)
							}
						}
						meshConfigMap["prometheus"] = prometheusSli

						proxySli := make([]map[string]interface{}, 0)
						if proxy, ok := meshConfigArg["Proxy"]; ok {
							if proxyArg, ok := proxy.(map[string]interface{}); ok && len(proxyArg) > 0 {
								proxyMap := make(map[string]interface{})
								proxyMap["cluster_domain"] = proxyArg["ClusterDomain"]
								proxyMap["limit_cpu"] = proxyArg["LimitCPU"]
								proxyMap["limit_memory"] = proxyArg["LimitMemory"]
								proxyMap["request_cpu"] = proxyArg["RequestCPU"]
								proxyMap["request_memory"] = proxyArg["RequestMemory"]
								proxySli = append(proxySli, proxyMap)
							}
						}
						meshConfigMap["proxy"] = proxySli

						sidecarInjectorSli := make([]map[string]interface{}, 0)
						if sidecarInjector, ok := meshConfigArg["SidecarInjector"]; ok {
							if sidecarInjectorArg, ok := sidecarInjector.(map[string]interface{}); ok && len(sidecarInjectorArg) > 0 {
								sidecarInjectorMap := make(map[string]interface{})
								sidecarInjectorMap["auto_injection_policy_enabled"] = sidecarInjectorArg["AutoInjectionPolicyEnabled"]
								sidecarInjectorMap["enable_namespaces_by_default"] = sidecarInjectorArg["EnableNamespacesByDefault"]

								initCNIConfigurationSli := make([]map[string]interface{}, 0)
								if initCNIConfiguration, ok := sidecarInjectorArg["InitCNIConfiguration"]; ok {
									if initCNIConfigurationArg, ok := initCNIConfiguration.(map[string]interface{}); ok && len(initCNIConfigurationArg) > 0 {
										initCNIConfigurationMap := make(map[string]interface{})
										initCNIConfigurationMap["enabled"] = initCNIConfigurationArg["Enabled"]
										initCNIConfigurationMap["exclude_namespaces"] = initCNIConfigurationArg["ExcludeNamespaces"]
										initCNIConfigurationSli = append(initCNIConfigurationSli, initCNIConfigurationMap)
									}
								}
								sidecarInjectorMap["init_cni_configuration"] = initCNIConfigurationSli
								sidecarInjectorMap["limit_cpu"] = sidecarInjectorArg["LimitCPU"]
								sidecarInjectorMap["limit_memory"] = sidecarInjectorArg["LimitMemory"]
								sidecarInjectorMap["request_cpu"] = sidecarInjectorArg["RequestCPU"]
								sidecarInjectorMap["request_memory"] = sidecarInjectorArg["RequestMemory"]
								if v, ok := sidecarInjectorArg["SidecarInjectorWebhookAsYaml"]; ok {
									sidecarInjectorMap["sidecar_injector_webhook_as_yaml"] = v
								}
								sidecarInjectorSli = append(sidecarInjectorSli, sidecarInjectorMap)
							}
						}
						meshConfigMap["sidecar_injector"] = sidecarInjectorSli
						meshConfigMap["telemetry"] = meshConfigArg["Telemetry"]
						meshConfigMap["tracing"] = meshConfigArg["Tracing"]
						meshConfigSli = append(meshConfigSli, meshConfigMap)
					}
				}
				mapping["mesh_config"] = meshConfigSli

				loadBalancerSli := make([]map[string]interface{}, 0)
				if loadBalancer, ok := specArg["LoadBalancer"]; ok {
					if loadBalancerArg, ok := loadBalancer.(map[string]interface{}); ok && len(loadBalancerArg) > 0 {
						loadBalancerMap := make(map[string]interface{})
						loadBalancerMap["pilot_public_eip"] = loadBalancerArg["PilotPublicEip"]
						loadBalancerMap["pilot_public_loadbalancer_id"] = loadBalancerArg["PilotPublicLoadbalancerId"]
						loadBalancerMap["api_server_loadbalancer_id"] = loadBalancerArg["ApiServerLoadbalancerId"]
						loadBalancerMap["api_server_public_eip"] = loadBalancerArg["ApiServerPublicEip"]
						loadBalancerSli = append(loadBalancerSli, loadBalancerMap)
					}
				}
				mapping["load_balancer"] = loadBalancerSli

				networkSli := make([]map[string]interface{}, 0)
				if network, ok := specArg["Network"]; ok {
					if networkArg, ok := network.(map[string]interface{}); ok && len(networkArg) > 0 {
						networkMap := make(map[string]interface{})
						networkMap["vswitche_list"] = networkArg["VSwitches"]
						networkMap["vpc_id"] = networkArg["VpcId"]
						networkMap["security_group_id"] = networkArg["SecurityGroupId"]
						networkSli = append(networkSli, networkMap)
					}
				}
				mapping["network"] = networkSli
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("meshes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
