package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbListenerCreate,
		Read:   resourceAlicloudAlbListenerRead,
		Update: resourceAlicloudAlbListenerUpdate,
		Delete: resourceAlicloudAlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_log_record_customized_headers_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"access_log_tracing_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tracing_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tracing_sample": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10000),
						},
						"tracing_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"acl_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_relations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"acl_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"acl_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"White", "Black"}, false),
						},
					},
				},
			},
			"default_actions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"forward_group_config": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_id": {
													Type:     schema.TypeString,
													Required: true,
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gzip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"http2_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("listener_protocol"); ok && v.(string) == "HTTPS" {
						return false
					}
					return true
				},
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"listener_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([^\x00-\xff]|[\w.,;/@-]){2,256}$`), "\t\nThe description of the listener.\n\nThe description must be 2 to 256 characters in length. The name can contain only the characters in the following string: /^([^\\x00-\\xff]|[\\w.,;/@-]){2,256}$/."),
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"listener_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "QUIC"}, false),
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quic_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quic_listener_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"quic_upgrade_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if v, ok := d.GetOk("listener_protocol"); ok && v.(string) == "HTTPS" {
									return false
								}
								return true
							},
						},
					},
				},
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 180),
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("listener_protocol"); ok && v.(string) == "HTTPS" {
						return false
					}
					return true
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"xforwarded_for_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("listener_protocol"); ok && v.(string) == "HTTPS" {
						return false
					}
					return true
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"xforwardedforclientcert_issuerdnalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcert_issuerdnenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertclientverifyalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertclientverifyenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertfingerprintalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertfingerprintenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertsubjectdnalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertsubjectdnenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientsrcportenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforprotoenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforslbidenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforslbportenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateListener"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("certificates"); ok {
		certificatesMaps := make([]map[string]interface{}, 0)
		for _, certificates := range v.(*schema.Set).List() {
			certificatesArg := certificates.(map[string]interface{})
			certificatesMap := map[string]interface{}{}
			certificatesMap["CertificateId"] = certificatesArg["certificate_id"]
			certificatesMaps = append(certificatesMaps, certificatesMap)
		}
		request["Certificates"] = certificatesMaps
	}
	defaultActionsMaps := make([]map[string]interface{}, 0)
	for _, defaultActions := range d.Get("default_actions").(*schema.Set).List() {
		defaultActionsArg := defaultActions.(map[string]interface{})
		defaultActionsMap := map[string]interface{}{}
		defaultActionsMap["Type"] = defaultActionsArg["type"]
		forwardGroupConfigMap := map[string]interface{}{}
		for _, forwardGroupConfig := range defaultActionsArg["forward_group_config"].(*schema.Set).List() {
			forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
			serverGroupTuplesMaps := make([]map[string]interface{}, 0)
			for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
				serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
				serverGroupTuplesMap := map[string]interface{}{}
				serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
				serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
			}
			forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
		}

		defaultActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		defaultActionsMaps = append(defaultActionsMaps, defaultActionsMap)
	}

	request["DefaultActions"] = defaultActionsMaps
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		request["GzipEnabled"] = v
	}
	if v, ok := d.GetOkExists("http2_enabled"); ok {
		request["Http2Enabled"] = v
	}
	if v, ok := d.GetOk("idle_timeout"); ok {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOk("listener_description"); ok {
		request["ListenerDescription"] = v
	}
	request["ListenerPort"] = d.Get("listener_port")
	request["ListenerProtocol"] = d.Get("listener_protocol")
	request["LoadBalancerId"] = d.Get("load_balancer_id")
	if v, ok := d.GetOk("request_timeout"); ok {
		request["RequestTimeout"] = v
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}
	if v, ok := d.GetOk("xforwarded_for_config"); ok {
		xforwardedForConfigMap := map[string]interface{}{}
		for _, xforwardedForConfig := range v.(*schema.Set).List() {
			xforwardedForConfigArg := xforwardedForConfig.(map[string]interface{})
			xforwardedForConfigMap["XForwardedForClientCertIssuerDNAlias"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnalias"]
			xforwardedForConfigMap["XForwardedForClientCertIssuerDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnenabled"]
			xforwardedForConfigMap["XForwardedForClientCertClientVerifyAlias"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyalias"]
			xforwardedForConfigMap["XForwardedForClientCertClientVerifyEnabled"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyenabled"]
			xforwardedForConfigMap["XForwardedForClientCertFingerprintAlias"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintalias"]
			xforwardedForConfigMap["XForwardedForClientCertFingerprintEnabled"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintenabled"]
			xforwardedForConfigMap["XForwardedForClientCertSubjectDNAlias"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnalias"]
			xforwardedForConfigMap["XForwardedForClientCertSubjectDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnenabled"]
			xforwardedForConfigMap["XForwardedForClientSrcPortEnabled"] = xforwardedForConfigArg["xforwardedforclientsrcportenabled"]
			xforwardedForConfigMap["XForwardedForEnabled"] = xforwardedForConfigArg["xforwardedforenabled"]
			xforwardedForConfigMap["XForwardedForProtoEnabled"] = xforwardedForConfigArg["xforwardedforprotoenabled"]
			xforwardedForConfigMap["XForwardedForSLBIdEnabled"] = xforwardedForConfigArg["xforwardedforslbidenabled"]
			xforwardedForConfigMap["XForwardedForSLBPortEnabled"] = xforwardedForConfigArg["xforwardedforslbportenabled"]
		}

		request["XForwardedForConfig"] = xforwardedForConfigMap
	}
	request["ClientToken"] = buildClientToken("CreateListener")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))
	albService := AlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlbListenerUpdate(d, meta)
}

func resourceAlicloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbListener(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_listener albService.DescribeAlbListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_log_record_customized_headers_enabled", object["LogConfig"].(map[string]interface{})["AccessLogRecordCustomizedHeadersEnabled"])

	accessLogTracingConfigSli := make([]map[string]interface{}, 0)
	if accessLogTracingConfig, ok := object["LogConfig"].(map[string]interface{})["AccessLogTracingConfig"]; ok && len(accessLogTracingConfig.(map[string]interface{})) > 0 {
		accessLogTracingConfigMap := make(map[string]interface{})
		accessLogTracingConfigMap["tracing_enabled"] = accessLogTracingConfig.(map[string]interface{})["TracingEnabled"]
		accessLogTracingConfigMap["tracing_sample"] = accessLogTracingConfig.(map[string]interface{})["TracingSample"]
		accessLogTracingConfigMap["tracing_type"] = accessLogTracingConfig.(map[string]interface{})["TracingType"]
		accessLogTracingConfigSli = append(accessLogTracingConfigSli, accessLogTracingConfigMap)
	}
	d.Set("access_log_tracing_config", accessLogTracingConfigSli)

	aclConfigSli := make([]map[string]interface{}, 0)
	if aclConfig, ok := object["AclConfig"]; ok && len(aclConfig.(map[string]interface{})) > 0 {
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
	d.Set("acl_config", aclConfigSli)
	if certificatesList, ok := object["Certificates"]; ok && certificatesList != nil {
		certificatesMaps := make([]map[string]interface{}, 0)
		for _, certificatesListItem := range certificatesList.([]interface{}) {
			if certificatesListItemArg, ok := certificatesListItem.(map[string]interface{}); ok {
				certificatesListItemMap := map[string]interface{}{}
				certificatesListItemMap["certificate_id"] = certificatesListItemArg["CertificateId"]
				certificatesMaps = append(certificatesMaps, certificatesListItemMap)
			}
		}

		d.Set("certificates", certificatesMaps)
	}

	if defaultActionsList, ok := object["DefaultActions"]; ok && defaultActionsList != nil {
		defaultActionsMaps := make([]map[string]interface{}, 0)
		for _, defaultActions := range defaultActionsList.([]interface{}) {
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

		d.Set("default_actions", defaultActionsMaps)
	}

	d.Set("gzip_enabled", object["GzipEnabled"])
	d.Set("http2_enabled", object["Http2Enabled"])
	if v, ok := object["IdleTimeout"]; ok && fmt.Sprint(v) != "0" {
		d.Set("idle_timeout", formatInt(v))
	}
	d.Set("listener_description", object["ListenerDescription"])
	if v, ok := object["ListenerPort"]; ok && fmt.Sprint(v) != "0" {
		d.Set("listener_port", formatInt(v))
	}
	d.Set("listener_protocol", object["ListenerProtocol"])
	d.Set("load_balancer_id", object["LoadBalancerId"])
	d.Set("status", object["ListenerStatus"])

	if quicConfig, ok := object["QuicConfig"]; ok {
		quicConfigSli := make([]map[string]interface{}, 0)
		if len(quicConfig.(map[string]interface{})) > 0 {
			quicConfigMap := make(map[string]interface{})
			quicConfigMap["quic_listener_id"] = quicConfig.(map[string]interface{})["QuicListenerId"]
			quicConfigMap["quic_upgrade_enabled"] = quicConfig.(map[string]interface{})["QuicUpgradeEnabled"]
			quicConfigSli = append(quicConfigSli, quicConfigMap)
		}
		d.Set("quic_config", quicConfigSli)
	}

	if v, ok := object["RequestTimeout"]; ok && fmt.Sprint(v) != "0" {
		d.Set("request_timeout", formatInt(v))
	}
	d.Set("security_policy_id", object["SecurityPolicyId"])

	if xforwardedForConfig, ok := object["XForwardedForConfig"]; ok && len(xforwardedForConfig.(map[string]interface{})) > 0 {
		xforwardedForConfigSli := make([]map[string]interface{}, 0)
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
		d.Set("xforwarded_for_config", xforwardedForConfigSli)
	}
	return nil
}

func resourceAlicloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ListenerId": d.Id(),
	}
	if d.HasChange("access_log_record_customized_headers_enabled") || d.IsNewResource() {
		if v, ok := d.GetOkExists("access_log_record_customized_headers_enabled"); ok {
			update = true
			request["AccessLogRecordCustomizedHeadersEnabled"] = v
		}
	}
	if d.HasChange("access_log_tracing_config") {
		if v, ok := d.GetOk("access_log_tracing_config"); ok {
			update = true
			accessLogTracingConfigMap := map[string]interface{}{}
			for _, certificates := range v.(*schema.Set).List() {
				certificatesArg := certificates.(map[string]interface{})
				accessLogTracingConfigMap["TracingEnabled"] = certificatesArg["tracing_enabled"]
				accessLogTracingConfigMap["TracingSample"] = certificatesArg["tracing_sample"]
				accessLogTracingConfigMap["TracingType"] = certificatesArg["tracing_type"]
			}
			request["AccessLogTracingConfig"] = accessLogTracingConfigMap
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateListenerLogConfig"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateListenerLogConfig")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("access_log_record_customized_headers_enabled")
		d.SetPartial("access_log_tracing_config")
	}
	update = false

	updateListenerAttributeReq := map[string]interface{}{
		"ListenerId": d.Id(),
	}
	if d.HasChange("certificates") {
		update = true
		if v, ok := d.GetOk("certificates"); ok {
			certificatesMaps := make([]map[string]interface{}, 0)
			for _, certificates := range v.(*schema.Set).List() {
				certificatesArg := certificates.(map[string]interface{})
				certificatesMap := map[string]interface{}{}
				certificatesMap["CertificateId"] = certificatesArg["certificate_id"]
				certificatesMaps = append(certificatesMaps, certificatesMap)
			}
			updateListenerAttributeReq["Certificates"] = certificatesMaps
		}
	}
	if !d.IsNewResource() && d.HasChange("default_actions") {
		update = true
		defaultActionsMaps := make([]map[string]interface{}, 0)
		for _, defaultActions := range d.Get("default_actions").(*schema.Set).List() {
			defaultActionsArg := defaultActions.(map[string]interface{})
			defaultActionsMap := map[string]interface{}{}
			defaultActionsMap["Type"] = defaultActionsArg["type"]
			forwardGroupConfigMap := map[string]interface{}{}
			for _, forwardGroupConfig := range defaultActionsArg["forward_group_config"].(*schema.Set).List() {
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				serverGroupTuplesMaps := make([]map[string]interface{}, 0)
				for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
					serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
					serverGroupTuplesMap := map[string]interface{}{}
					serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
				forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
			}

			defaultActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
			defaultActionsMaps = append(defaultActionsMaps, defaultActionsMap)
		}
		updateListenerAttributeReq["DefaultActions"] = defaultActionsMaps
	}
	if !d.IsNewResource() && d.HasChange("server_group_id") {
		update = true
		updateListenerAttributeReq["DefaultActions.*.ForwardGroupConfig.ServerGroupTuples.*.ServerGroupId"] = d.Get("server_group_id")
	}
	if !d.IsNewResource() && d.HasChange("type") {
		update = true
		updateListenerAttributeReq["DefaultActions.*.Type"] = d.Get("type")
	}
	if !d.IsNewResource() && d.HasChange("gzip_enabled") {
		update = true
		if v, ok := d.GetOkExists("gzip_enabled"); ok {
			updateListenerAttributeReq["GzipEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("http2_enabled") {
		update = true
		if v, ok := d.GetOkExists("http2_enabled"); ok {
			updateListenerAttributeReq["Http2Enabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("idle_timeout") {
		update = true
		if v, ok := d.GetOk("idle_timeout"); ok {
			updateListenerAttributeReq["IdleTimeout"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("listener_description") {
		update = true
		if v, ok := d.GetOk("listener_description"); ok {
			updateListenerAttributeReq["ListenerDescription"] = v
		}
	}
	if d.HasChange("quic_config") {
		if v, ok := d.GetOk("quic_config"); ok {
			update = true
			quicConfigMap := map[string]interface{}{}

			for _, quicConfig := range v.(*schema.Set).List() {
				quicConfigArg := quicConfig.(map[string]interface{})
				quicConfigMap["QuicListenerId"] = quicConfigArg["quic_listener_id"]
				quicConfigMap["QuicUpgradeEnabled"] = quicConfigArg["quic_upgrade_enabled"]
			}

			updateListenerAttributeReq["QuicConfig"] = quicConfigMap
		}
	}
	if !d.IsNewResource() && d.HasChange("request_timeout") {
		update = true
		if v, ok := d.GetOk("request_timeout"); ok {
			updateListenerAttributeReq["RequestTimeout"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("security_policy_id") {
		update = true
		if v, ok := d.GetOk("security_policy_id"); ok {
			updateListenerAttributeReq["SecurityPolicyId"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("xforwarded_for_config") {
		update = true
		if v, ok := d.GetOk("xforwarded_for_config"); ok {
			xforwardedForConfigMap := map[string]interface{}{}
			for _, xforwardedForConfig := range v.(*schema.Set).List() {
				xforwardedForConfigArg := xforwardedForConfig.(map[string]interface{})
				xforwardedForConfigMap["XForwardedForClientCertIssuerDNAlias"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnalias"]
				xforwardedForConfigMap["XForwardedForClientCertIssuerDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnenabled"]
				xforwardedForConfigMap["XForwardedForClientCertClientVerifyAlias"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyalias"]
				xforwardedForConfigMap["XForwardedForClientCertClientVerifyEnabled"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyenabled"]
				xforwardedForConfigMap["XForwardedForClientCertFingerprintAlias"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintalias"]
				xforwardedForConfigMap["XForwardedForClientCertFingerprintEnabled"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintenabled"]
				xforwardedForConfigMap["XForwardedForClientCertSubjectDNAlias"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnalias"]
				xforwardedForConfigMap["XForwardedForClientCertSubjectDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnenabled"]
				xforwardedForConfigMap["XForwardedForClientSrcPortEnabled"] = xforwardedForConfigArg["xforwardedforclientsrcportenabled"]
				xforwardedForConfigMap["XForwardedForEnabled"] = xforwardedForConfigArg["xforwardedforenabled"]
				xforwardedForConfigMap["XForwardedForProtoEnabled"] = xforwardedForConfigArg["xforwardedforprotoenabled"]
				xforwardedForConfigMap["XForwardedForSLBIdEnabled"] = xforwardedForConfigArg["xforwardedforslbidenabled"]
				xforwardedForConfigMap["XForwardedForSLBPortEnabled"] = xforwardedForConfigArg["xforwardedforslbportenabled"]
			}

			updateListenerAttributeReq["XForwardedForConfig"] = xforwardedForConfigMap
		}
	}

	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateListenerAttributeReq["DryRun"] = v
		}
		action := "UpdateListenerAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		updateListenerAttributeReq["ClientToken"] = buildClientToken("UpdateListenerAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateListenerAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateListenerAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running", "Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("certificates")
		d.SetPartial("default_actions")
		d.SetPartial("server_group_id")
		d.SetPartial("type")
		d.SetPartial("gzip_enabled")
		d.SetPartial("http2_enabled")
		d.SetPartial("idle_timeout")
		d.SetPartial("listener_description")
		d.SetPartial("quic_config")
		d.SetPartial("request_timeout")
		d.SetPartial("security_policy_id")
		d.SetPartial("xforwarded_for_config")
	}
	if d.HasChange("status") {
		object, err := albService.DescribeAlbListener(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["ListenerStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"ListenerId": d.Id(),
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StartListener"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("StartListener")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"ListenerId": d.Id(),
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StopListener"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("StopListener")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	if d.HasChange("acl_config") {
		oldAssociateAcls, newAssociateVpcs := d.GetChange("acl_config")
		oldAssociateAclsSet := oldAssociateAcls.(*schema.Set)
		newAssociateAclsSet := newAssociateVpcs.(*schema.Set)
		removed := oldAssociateAclsSet.Difference(newAssociateAclsSet)
		added := newAssociateAclsSet.Difference(oldAssociateAclsSet)
		if removed.Len() > 0 {
			action := "DissociateAclsFromListener"
			dissociateAclsFromListenerReq := map[string]interface{}{
				"ListenerId": d.Id(),
			}
			dissociateAclsFromListenerReq["ClientToken"] = buildClientToken("DissociateAclsFromListener")
			associateAclIds := make([]string, 0)
			for _, aclConfig := range removed.List() {
				if aclRelationsMaps, ok := aclConfig.(map[string]interface{})["acl_relations"]; ok {
					for _, aclRelationsMap := range aclRelationsMaps.(*schema.Set).List() {
						associateAclIds = append(associateAclIds, aclRelationsMap.(map[string]interface{})["acl_id"].(string))
					}
				}
			}
			dissociateAclsFromListenerReq["AclIds"] = associateAclIds
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, dissociateAclsFromListenerReq, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, dissociateAclsFromListenerReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if added.Len() > 0 {
			action := "AssociateAclsWithListener"
			associateAclsWithListenerReq := map[string]interface{}{
				"ListenerId": d.Id(),
			}
			associateAclsWithListenerReq["ClientToken"] = buildClientToken("AssociateAclsWithListener")
			associateAclIds := make([]string, 0)
			for _, aclConfig := range added.List() {
				if aclRelationsMaps, ok := aclConfig.(map[string]interface{})["acl_relations"]; ok {
					for _, aclRelationsMap := range aclRelationsMaps.(*schema.Set).List() {
						associateAclIds = append(associateAclIds, aclRelationsMap.(map[string]interface{})["acl_id"].(string))
					}
				}
				associateAclsWithListenerReq["AclType"] = aclConfig.(map[string]interface{})["acl_type"]
			}
			associateAclsWithListenerReq["AclIds"] = associateAclIds
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, associateAclsWithListenerReq, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, associateAclsWithListenerReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	return resourceAlicloudAlbListenerRead(d, meta)
}

func resourceAlicloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteListener"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ListenerId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteListener")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "Throttling", "-22031"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Listener", "ResourceNotFound.LoadBalancer"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
