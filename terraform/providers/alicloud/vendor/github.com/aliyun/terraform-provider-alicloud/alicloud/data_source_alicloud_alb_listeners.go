package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbListenersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"listener_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "QUIC"}, false),
			},
			"load_balancer_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Provisioning", "Running", "Configuring", "Stopped"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log_record_customized_headers_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"access_log_tracing_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tracing_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tracing_sample": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"tracing_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"acl_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"acl_relations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"acl_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"acl_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"certificates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"forward_group_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_tuples": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"server_group_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"gzip_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"http2_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"idle_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_results": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quic_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"quic_listener_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quic_upgrade_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"request_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xforwarded_for_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"xforwardedforclientcert_issuerdnalias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"xforwardedforclientcert_issuerdnenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforclientcertclientverifyalias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"xforwardedforclientcertclientverifyenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforclientcertfingerprintalias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"xforwardedforclientcertfingerprintenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforclientcertsubjectdnalias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"xforwardedforclientcertsubjectdnenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforclientsrcportenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforprotoenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforslbidenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"xforwardedforslbportenabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
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

func dataSourceAlicloudAlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListListeners"
	request := make(map[string]interface{})
	if m, ok := d.GetOk("listener_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("ListenerIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("listener_protocol"); ok {
		request["ListenerProtocol"] = v
	}
	if m, ok := d.GetOk("load_balancer_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("LoadBalancerIds.%d", k+1)] = v.(string)
		}
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}

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
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_listeners", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Listeners", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Listeners", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ListenerId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["ListenerStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"max_results": fmt.Sprint(response["MaxResults"]),
			"next_token":  response["NextToken"],
		}
		mapping["access_log_record_customized_headers_enabled"] = object["LogConfig.AccessLogRecordCustomizedHeadersEnabled"]
		mapping["gzip_enabled"] = object["GzipEnabled"]
		mapping["http2_enabled"] = object["Http2Enabled"]
		mapping["idle_timeout"] = object["IdleTimeout"]
		mapping["listener_description"] = object["ListenerDescription"]
		mapping["id"] = fmt.Sprint(object["ListenerId"])
		mapping["listener_id"] = object["ListenerId"]
		mapping["listener_port"] = object["ListenerPort"]
		mapping["listener_protocol"] = object["ListenerProtocol"]
		mapping["load_balancer_id"] = object["LoadBalancerId"]
		mapping["request_timeout"] = object["RequestTimeout"]
		mapping["security_policy_id"] = object["SecurityPolicyId"]
		mapping["status"] = object["ListenerStatus"]
		xforwardedForConfigSli := make([]map[string]interface{}, 0)
		if xforwardedForConfig, ok := object["XForwardedForConfig"]; ok && len(xforwardedForConfig.(map[string]interface{})) > 0 {
			xforwardedForConfigMap := make(map[string]interface{})
			xforwardedForConfigMap["xforwardedforclientcert_issuerdnalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertIssuerDNAlias"]
			xforwardedForConfigMap["xforwardedforclientcert_issuerdnenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertIssuerDNEnabled"]
			xforwardedForConfigMap["xforwardedforclientcertclientverifyalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertClientVerifyAlias"]
			xforwardedForConfigMap["xforwardedforclientcertclientverifyenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertClientVerifyEnabled"]
			xforwardedForConfigMap["xforwardedforclientcertfingerprintalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertFingerprintAlias"]
			xforwardedForConfigMap["xforwardedforclientcertfingerprintenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertFingerprintEnabled"]
			xforwardedForConfigMap["xforwardedforclientcertsubjectdnalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertSubjectDNAlias"]
			xforwardedForConfigMap["xforwardedforclientcertsubjectdnenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertSubjectDNEnabled"]
			xforwardedForConfigMap["xforwardedforclientsrcportenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientSrcPortEnabled"]
			xforwardedForConfigMap["xforwardedforenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForEnabled"]
			xforwardedForConfigMap["xforwardedforprotoenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForProtoEnabled"]
			xforwardedForConfigMap["xforwardedforslbidenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForSLBIdEnabled"]
			xforwardedForConfigMap["xforwardedforslbportenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForSLBPortEnabled"]
			xforwardedForConfigSli = append(xforwardedForConfigSli, xforwardedForConfigMap)
		}
		mapping["xforwarded_for_config"] = xforwardedForConfigSli

		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ListenerId"])
		albService := AlbService{client}
		getResp, err := albService.DescribeAlbListener(id)
		if err != nil {
			return WrapError(err)
		}
		certificatesMaps := make([]map[string]interface{}, 0)
		if certificatesList, ok := getResp["Certificates"]; ok && certificatesList != nil {
			for _, certificatesListItem := range certificatesList.([]interface{}) {
				if certificatesListItemArg, ok := certificatesListItem.(map[string]interface{}); ok {
					certificatesListItemMap := map[string]interface{}{}
					certificatesListItemMap["certificate_id"] = certificatesListItemArg["CertificateId"]
					certificatesMaps = append(certificatesMaps, certificatesListItemMap)
				}
			}
		}
		mapping["certificates"] = certificatesMaps
		defaultActionsMaps := make([]map[string]interface{}, 0)
		if defaultActionsList, ok := getResp["DefaultActions"].([]interface{}); ok {
			for _, defaultActions := range defaultActionsList {
				defaultActionsArg := defaultActions.(map[string]interface{})
				defaultActionsMap := map[string]interface{}{}
				defaultActionsMap["type"] = defaultActionsArg["Type"]
				forwardGroupConfig := defaultActionsArg["ForwardGroupConfig"]
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				serverGroupTuplesMaps := make([]map[string]interface{}, 0)
				for _, serverGroupTuples := range forwardGroupConfigArg["ServerGroupTuples"].([]interface{}) {
					serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
					serverGroupTuplesMap := map[string]interface{}{}
					serverGroupTuplesMap["server_group_id"] = serverGroupTuplesArg["ServerGroupId"]
					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
				forwardGroupConfigMaps := make([]map[string]interface{}, 0)
				forwardGroupConfigMap := map[string]interface{}{}
				forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
				forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
				defaultActionsMap["forward_group_config"] = forwardGroupConfigMaps
				defaultActionsMaps = append(defaultActionsMaps, defaultActionsMap)
			}
		}
		mapping["default_actions"] = defaultActionsMaps
		aclConfigSli := make([]map[string]interface{}, 0)
		if aclConfig, ok := getResp["AclConfig"]; ok && len(aclConfig.(map[string]interface{})) > 0 {
			aclConfigMap := make(map[string]interface{})
			aclRelationsSli := make([]map[string]interface{}, 0)
			if v, ok := aclConfig.(map[string]interface{})["AclRelations"]; ok && len(v.([]interface{})) > 0 {
				for _, aclRelations := range v.([]interface{}) {
					aclRelationsMap := make(map[string]interface{})
					aclRelationsMap["acl_id"] = aclRelations.(map[string]interface{})["AclId"]
					aclRelationsMap["status"] = aclRelations.(map[string]interface{})["Status"]
					aclRelationsSli = append(aclRelationsSli, aclRelationsMap)
				}
			}
			aclConfigMap["acl_relations"] = aclRelationsSli
			aclConfigMap["acl_type"] = aclConfig.(map[string]interface{})["AclType"]
			aclConfigSli = append(aclConfigSli, aclConfigMap)
		}
		mapping["acl_config"] = aclConfigSli

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("listeners", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
