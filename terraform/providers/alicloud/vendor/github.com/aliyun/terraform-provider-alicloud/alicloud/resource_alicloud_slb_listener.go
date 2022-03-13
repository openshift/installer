package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbListenerCreate,
		Read:   resourceAliyunSlbListenerRead,
		Update: resourceAliyunSlbListenerUpdate,
		Delete: resourceAliyunSlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Required:     true,
				ForceNew:     true,
			},
			"lb_port": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'lb_port' has been deprecated, and using 'frontend_port' to replace.",
			},

			"backend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
				ForceNew:     true,
			},

			"instance_port": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'instance_port' has been deprecated, and using 'backend_port' to replace.",
			},

			"lb_protocol": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'lb_protocol' has been deprecated, and using 'protocol' to replace.",
			},

			"protocol": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp", "udp"}, false),
				Required:     true,
				ForceNew:     true,
			},

			"bandwidth": {
				Type: schema.TypeInt,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 1000),
					validation.IntInSlice([]int{-1})),
				Optional: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"wrr", "wlc", "rr", "sch", "tch", "qch"}, false),
				Optional:     true,
				Default:      WRRScheduler,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_slave_server_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acl_status": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
				Optional:     true,
				Default:      OffFlag,
			},
			"acl_type": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"black", "white"}, false),
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			"acl_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			//http & https
			"sticky_session": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OffFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(InsertStickySessionType),
					string(ServerStickySessionType)}, false),
				Optional:         true,
				DiffSuppressFunc: stickySessionTypeDiffSuppressFunc,
			},
			//http & https
			"cookie_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 86400),
				Optional:         true,
				DiffSuppressFunc: cookieTimeoutDiffSuppressFunc,
			},
			//http & https
			"cookie": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 200),
				Optional:         true,
				DiffSuppressFunc: cookieDiffSuppressFunc,
			},
			//tcp & udp
			"persistence_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 3600),
				Optional:         true,
				Default:          0,
				DiffSuppressFunc: tcpUdpDiffSuppressFunc,
			},
			//http & https
			"health_check": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//http & https
			"health_check_method": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"head", "get"}, false),
				Optional:     true,
				Computed:     true,
			},
			//tcp
			"health_check_type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(TCPHealthCheckType),
					string(HTTPHealthCheckType)}, false),
				Optional:         true,
				Default:          TCPHealthCheckType,
				DiffSuppressFunc: healthCheckTypeDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_domain": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringDoesNotMatch(regexp.MustCompile(`^\$_ip$`), "value '$_ip' has been deprecated, and empty string will replace it"),
				Optional:         true,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_uri": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 80),
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type: schema.TypeInt,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 65535),
					validation.IntInSlice([]int{-520})),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},

			"health_check_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 300),
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 50),
				Optional:         true,
				Default:          2,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_http_code": {
				Type: schema.TypeString,
				ValidateFunc: validateAllowedSplitStringValue([]string{
					string(HTTP_2XX), string(HTTP_3XX), string(HTTP_4XX), string(HTTP_5XX)}, ","),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//https
			"ssl_certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
				Deprecated:       "Field 'ssl_certificate_id' has been deprecated from 1.59.0 and using 'server_certificate_id' instead.",
			},
			"server_certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
			},

			"ca_certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
			},
			//http, https
			"gzip": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			"x_forwarded_for": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// At present, retrive client ip can not be modified, and it default to true.
						"retrive_client_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"retrive_slb_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retrive_slb_id": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retrive_slb_proto": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				MaxItems: 1,
			},
			//tcp
			"established_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(10, 900),
				Optional:         true,
				Default:          900,
				DiffSuppressFunc: establishedTimeoutDiffSuppressFunc,
			},

			//http & https
			"idle_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 60),
				Optional:         true,
				Default:          15,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},

			//http & https
			"request_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 180),
				Optional:         true,
				Default:          60,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},

			//https
			"enable_http2": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: httpsDiffSuppressFunc,
			},

			//https
			"tls_cipher_policy": {
				Type:             schema.TypeString,
				Default:          "tls_cipher_policy_1_0",
				ValidateFunc:     validation.StringInSlice([]string{"tls_cipher_policy_1_0", "tls_cipher_policy_1_1", "tls_cipher_policy_1_2", "tls_cipher_policy_1_2_strict"}, false),
				Optional:         true,
				DiffSuppressFunc: httpsDiffSuppressFunc,
			},
			"forward_port": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 65535),
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: forwardPortDiffSuppressFunc,
			},
			"listener_forward": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: httpDiffSuppressFunc,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliyunSlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	httpForward := false
	protocol := d.Get("protocol").(string)
	lbId := d.Get("load_balancer_id").(string)
	frontend := d.Get("frontend_port").(int)
	if listenerForward, ok := d.GetOk("listener_forward"); ok && listenerForward.(string) == string(OnFlag) {
		httpForward = true
	}
	request, err := buildListenerCommonArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.ApiName = fmt.Sprintf("CreateLoadBalancer%sListener", strings.ToUpper(protocol))

	if Protocol(protocol) == Http || Protocol(protocol) == Https {
		if httpForward {
			reqHttp, err := buildHttpForwardArgs(d, request)
			if err != nil {
				return WrapError(err)
			}
			request = reqHttp
		} else {
			reqHttp, err := buildHttpListenerArgs(d, request)
			if err != nil {
				return WrapError(err)
			}
			request = reqHttp
		}
		if Protocol(protocol) == Https {
			scId := d.Get("server_certificate_id").(string)
			if scId == "" {
				scId = d.Get("ssl_certificate_id").(string)
			}

			if scId == "" {
				return WrapError(Error(`'server_certificate_id': required field is not set when the protocol is 'https'.`))
			}
			request.QueryParams["ServerCertificateId"] = scId

			if caId, ok := d.GetOk("ca_certificate_id"); ok && caId.(string) != "" {
				request.QueryParams["CACertificateId"] = caId.(string)
			}
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_listener", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	d.SetId(lbId + ":" + protocol + ":" + strconv.Itoa(frontend))

	if err := slbService.WaitForSlbListener(d.Id(), Stopped, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	startLoadBalancerListenerRequest := slb.CreateStartLoadBalancerListenerRequest()
	startLoadBalancerListenerRequest.RegionId = client.RegionId
	startLoadBalancerListenerRequest.LoadBalancerId = lbId
	startLoadBalancerListenerRequest.ListenerPort = requests.NewInteger(frontend)
	startLoadBalancerListenerRequest.ListenerProtocol = protocol

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.StartLoadBalancerListener(startLoadBalancerListenerRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceIsConfiguring"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(startLoadBalancerListenerRequest.GetActionName(), raw, startLoadBalancerListenerRequest.RpcRequest, startLoadBalancerListenerRequest)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_listener", startLoadBalancerListenerRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err = slbService.WaitForSlbListener(d.Id(), Running, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	if httpForward {
		return resourceAliyunSlbListenerRead(d, meta)
	}
	return resourceAliyunSlbListenerUpdate(d, meta)
}

func resourceAliyunSlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("protocol", protocol)
	d.Set("load_balancer_id", lb_id)
	d.Set("frontend_port", port)
	d.SetId(lb_id + ":" + protocol + ":" + strconv.Itoa(port))
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := slbService.DescribeSlbListener(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			if IsExpectedErrors(err, SlbIsBusy) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}

		if port, ok := object["ListenerPort"]; ok && port.(float64) > 0 {
			readListener(d, object)
		} else {
			d.SetId("")
		}
		return nil
	})
}

func resourceAliyunSlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	proto := d.Get("protocol").(string)
	lb_id := d.Get("load_balancer_id").(string)
	frontend := d.Get("frontend_port").(int)
	d.SetId(lb_id + ":" + proto + ":" + strconv.Itoa(frontend))

	client := meta.(*connectivity.AliyunClient)
	protocol := Protocol(d.Get("protocol").(string))
	commonRequest, err := buildListenerCommonArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	commonRequest.ApiName = fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))

	update := false
	if d.HasChange("description") {
		update = true
	}
	if d.HasChange("scheduler") {
		commonRequest.QueryParams["Scheduler"] = d.Get("scheduler").(string)
		update = true
	}

	if d.HasChange("server_group_id") {
		serverGroupId := d.Get("server_group_id").(string)
		if serverGroupId != "" {
			commonRequest.QueryParams["VServerGroup"] = string(OnFlag)
			commonRequest.QueryParams["VServerGroupId"] = d.Get("server_group_id").(string)
		} else {
			commonRequest.QueryParams["VServerGroup"] = string(OffFlag)
		}
		update = true
	}

	if d.HasChange("master_slave_server_group_id") {
		serverGroupId := d.Get("master_slave_server_group_id").(string)
		if serverGroupId != "" {
			commonRequest.QueryParams["MasterSlaveServerGroup"] = string(OnFlag)
			commonRequest.QueryParams["MasterSlaveServerGroupId"] = d.Get("master_slave_server_group_id").(string)
		} else {
			commonRequest.QueryParams["MasterSlaveServerGroup"] = string(OffFlag)
		}
		update = true
	}

	if d.HasChange("bandwidth") {
		commonRequest.QueryParams["Bandwidth"] = strconv.Itoa(d.Get("bandwidth").(int))
		update = true
	}

	if d.HasChange("acl_status") {
		commonRequest.QueryParams["AclStatus"] = d.Get("acl_status").(string)
		update = true
	}

	if d.HasChange("acl_type") {
		commonRequest.QueryParams["AclType"] = d.Get("acl_type").(string)
		update = true
	}

	if d.HasChange("acl_id") {
		commonRequest.QueryParams["AclId"] = d.Get("acl_id").(string)
		update = true
	}

	httpArgs, err := buildHttpListenerArgs(d, commonRequest)
	if (protocol == Https || protocol == Http) && err != nil {
		return WrapError(err)
	}
	// http https
	if d.HasChange("sticky_session") {
		update = true
	}
	if d.HasChange("sticky_session_type") {
		update = true
	}
	if d.HasChange("cookie_timeout") {
		update = true
	}
	if d.HasChange("cookie") {
		update = true
	}
	if d.HasChange("health_check") {
		update = true
	}

	d.SetPartial("gzip")
	if d.Get("gzip").(bool) {
		httpArgs.QueryParams["Gzip"] = string(OnFlag)
	} else {
		httpArgs.QueryParams["Gzip"] = string(OffFlag)
	}

	d.SetPartial("x_forwarded_for")
	if len(d.Get("x_forwarded_for").([]interface{})) > 0 && (d.Get("protocol").(string) == "http" || d.Get("protocol").(string) == "https") {
		xff := d.Get("x_forwarded_for").([]interface{})[0].(map[string]interface{})
		if xff["retrive_slb_ip"].(bool) {
			httpArgs.QueryParams["XForwardedFor_SLBIP"] = string(OnFlag)
		} else {
			httpArgs.QueryParams["XForwardedFor_SLBIP"] = string(OffFlag)
		}
		if xff["retrive_slb_id"].(bool) {
			httpArgs.QueryParams["XForwardedFor_SLBID"] = string(OnFlag)
		} else {
			httpArgs.QueryParams["XForwardedFor_SLBID"] = string(OffFlag)
		}
		if xff["retrive_slb_proto"].(bool) {
			httpArgs.QueryParams["XForwardedFor_proto"] = string(OnFlag)
		} else {
			httpArgs.QueryParams["XForwardedFor_proto"] = string(OffFlag)
		}
	}

	if d.HasChange("gzip") || d.HasChange("x_forwarded_for") {
		update = true
	}

	// http https
	if d.HasChange("idle_timeout") {
		update = true
	}

	// http https
	if d.HasChange("request_timeout") {
		update = true
	}

	// http https
	if d.HasChange("health_check_method") {
		update = true
	}

	// http https tcp udp and health_check=on
	if d.HasChange("unhealthy_threshold") {
		commonRequest.QueryParams["UnhealthyThreshold"] = string(requests.NewInteger(d.Get("unhealthy_threshold").(int)))
		update = true
	}
	if d.HasChange("healthy_threshold") {
		commonRequest.QueryParams["HealthyThreshold"] = string(requests.NewInteger(d.Get("healthy_threshold").(int)))
		update = true
	}
	if d.HasChange("health_check_timeout") {
		commonRequest.QueryParams["HealthCheckConnectTimeout"] = string(requests.NewInteger(d.Get("health_check_timeout").(int)))
		update = true
	}
	if d.HasChange("health_check_interval") {
		commonRequest.QueryParams["HealthCheckInterval"] = string(requests.NewInteger(d.Get("health_check_interval").(int)))
		update = true
	}
	if d.HasChange("health_check_connect_port") {
		if port, ok := d.GetOk("health_check_connect_port"); ok {
			httpArgs.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
			commonRequest.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
			update = true
		}
	}

	// tcp and udp
	if d.HasChange("persistence_timeout") {
		commonRequest.QueryParams["PersistenceTimeout"] = string(requests.NewInteger(d.Get("persistence_timeout").(int)))
		update = true
	}

	tcpArgs := commonRequest
	udpArgs := commonRequest

	// http https tcp
	if d.HasChange("health_check_domain") {
		update = true
	}
	// todo: depends on Api fixing the default value
	if domain, ok := d.GetOk("health_check_domain"); ok {
		httpArgs.QueryParams["HealthCheckDomain"] = domain.(string)
		tcpArgs.QueryParams["HealthCheckDomain"] = domain.(string)
	}
	if d.HasChange("health_check_uri") {
		tcpArgs.QueryParams["HealthCheckURI"] = d.Get("health_check_uri").(string)
		update = true
	}
	if d.HasChange("health_check_http_code") {
		tcpArgs.QueryParams["HealthCheckHttpCode"] = getHealthCheckHttpCodeValue(d)
		update = true
	}

	// tcp
	if d.HasChange("health_check_type") {
		tcpArgs.QueryParams["HealthCheckType"] = d.Get("health_check_type").(string)
		update = true
	}

	// tcp
	if d.HasChange("established_timeout") {
		tcpArgs.QueryParams["EstablishedTimeout"] = string(requests.NewInteger(d.Get("established_timeout").(int)))
		update = true
	}

	// https
	httpsArgs := httpArgs
	if protocol == Https {
		scId := d.Get("server_certificate_id").(string)
		if scId == "" {
			scId = d.Get("ssl_certificate_id").(string)
		}
		if scId == "" {
			return WrapError(Error("'server_certificate_id': required field is not set when the protocol is 'https'."))
		}

		httpsArgs.QueryParams["ServerCertificateId"] = scId
		if d.HasChange("ssl_certificate_id") || d.HasChange("server_certificate_id") {
			update = true
		}

		if d.HasChange("enable_http2") {
			httpsArgs.QueryParams["EnableHttp2"] = d.Get("enable_http2").(string)
			update = true
		}

		if d.HasChange("tls_cipher_policy") {
			// spec changes check, can not be updated when load balancer instance is "Shared-Performance".
			slbService := SlbService{client}
			object, err := slbService.DescribeSlb(d.Get("load_balancer_id").(string))
			if err != nil {
				return WrapError(err)
			}
			spec := object.LoadBalancerSpec
			if spec == "" {
				if !d.IsNewResource() || string("tls_cipher_policy_1_0") != d.Get("tls_cipher_policy").(string) {
					return WrapError(Error("Currently the param \"tls_cipher_policy\" can not be updated when load balancer instance is \"Shared-Performance\"."))
				}
			} else {
				httpsArgs.QueryParams["TLSCipherPolicy"] = d.Get("tls_cipher_policy").(string)
				update = true
			}
		}

		if d.HasChange("ca_certificate_id") {
			httpsArgs.QueryParams["CACertificateId"] = d.Get("ca_certificate_id").(string)
			update = true
		}
	}

	if update {
		var request *requests.CommonRequest
		switch protocol {
		case Https:
			request = httpsArgs
		case Tcp:
			request = tcpArgs
		case Udp:
			request = udpArgs
		default:
			request = httpArgs
		}
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
	}

	d.Partial(false)

	return resourceAliyunSlbListenerRead(d, meta)
}

func resourceAliyunSlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	lbId, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("delete_protection_validation"); ok && v.(bool) {
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current listener's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the listener resource.", lbId))
		}
	}
	request := slb.CreateDeleteLoadBalancerListenerRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = lbId
	request.ListenerPort = requests.NewInteger(port)
	request.ListenerProtocol = protocol
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteLoadBalancerListener(request)
		})

		if err != nil {
			if IsExpectedErrors(err, SlbIsBusy) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(slbService.WaitForSlbListener(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildListenerCommonArgs(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request, err := slbService.BuildSlbCommonRequest()
	if err != nil {
		return request, WrapError(err)
	}
	request.RegionId = client.RegionId
	request.QueryParams["LoadBalancerId"] = d.Get("load_balancer_id").(string)
	request.QueryParams["ListenerPort"] = string(requests.NewInteger(d.Get("frontend_port").(int)))
	if backendServerPort, ok := d.GetOk("backend_port"); ok {
		request.QueryParams["BackendServerPort"] = string(requests.NewInteger(backendServerPort.(int)))
	}
	if bandWidth, ok := d.GetOk("bandwidth"); ok {
		request.QueryParams["Bandwidth"] = string(requests.NewInteger(bandWidth.(int)))
	}

	if groupId, ok := d.GetOk("server_group_id"); ok && groupId.(string) != "" {
		request.QueryParams["VServerGroupId"] = groupId.(string)
	}

	if groupId, ok := d.GetOk("master_slave_server_group_id"); ok && groupId.(string) != "" {
		request.QueryParams["MasterSlaveServerGroupId"] = groupId.(string)
	}
	// acl status
	if aclStatus, ok := d.GetOk("acl_status"); ok && aclStatus.(string) != "" {
		request.QueryParams["AclStatus"] = aclStatus.(string)
	}
	// acl type
	if aclType, ok := d.GetOk("acl_type"); ok && aclType.(string) != "" {
		request.QueryParams["AclType"] = aclType.(string)
	}
	// acl id
	if aclId, ok := d.GetOk("acl_id"); ok && aclId.(string) != "" {
		request.QueryParams["AclId"] = aclId.(string)
	}
	// description
	if description, ok := d.GetOk("description"); ok && description.(string) != "" {
		request.QueryParams["Description"] = description.(string)
	}
	// scheduler
	if scheduler, ok := d.GetOk("scheduler"); ok && scheduler.(string) != "" {
		request.QueryParams["Scheduler"] = scheduler.(string)
	}

	return request, nil

}
func buildHttpListenerArgs(d *schema.ResourceData, req *requests.CommonRequest) (*requests.CommonRequest, error) {
	stickySession := d.Get("sticky_session").(string)
	healthCheck := d.Get("health_check").(string)
	req.QueryParams["StickySession"] = stickySession
	req.QueryParams["HealthCheck"] = healthCheck
	req.QueryParams["RequestTimeout"] = string(requests.NewInteger(d.Get("request_timeout").(int)))
	req.QueryParams["IdleTimeout"] = string(requests.NewInteger(d.Get("idle_timeout").(int)))
	if stickySession == string(OnFlag) {
		sessionType, ok := d.GetOk("sticky_session_type")
		if !ok || sessionType.(string) == "" {
			return req, WrapError(Error("'sticky_session_type': required field is not set when the StickySession is %s.", OnFlag))
		} else {
			req.QueryParams["StickySessionType"] = sessionType.(string)

		}
		if sessionType.(string) == string(InsertStickySessionType) {
			if timeout, ok := d.GetOk("cookie_timeout"); !ok || timeout == 0 {
				return req, WrapError(Error("'cookie_timeout': required field is not set when the StickySession is %s and StickySessionType is %s.",
					OnFlag, InsertStickySessionType))
			} else {
				req.QueryParams["CookieTimeout"] = string(requests.NewInteger(timeout.(int)))
			}
		} else {
			if cookie, ok := d.GetOk("cookie"); !ok || cookie.(string) == "" {
				return req, WrapError(fmt.Errorf("'cookie': required field is not set when the StickySession is %s and StickySessionType is %s.",
					OnFlag, ServerStickySessionType))
			} else {
				req.QueryParams["Cookie"] = cookie.(string)
			}
		}
	}
	if healthCheck == string(OnFlag) {
		req.QueryParams["HealthCheckURI"] = d.Get("health_check_uri").(string)
		if port, ok := d.GetOk("health_check_connect_port"); ok {
			req.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
		}
		req.QueryParams["HealthyThreshold"] = string(requests.NewInteger(d.Get("healthy_threshold").(int)))
		req.QueryParams["UnhealthyThreshold"] = string(requests.NewInteger(d.Get("unhealthy_threshold").(int)))
		req.QueryParams["HealthCheckTimeout"] = string(requests.NewInteger(d.Get("health_check_timeout").(int)))
		req.QueryParams["HealthCheckInterval"] = string(requests.NewInteger(d.Get("health_check_interval").(int)))
		req.QueryParams["HealthCheckHttpCode"] = getHealthCheckHttpCodeValue(d)
		if d.Get("protocol").(string) == "http" || d.Get("protocol").(string) == "https" {
			if healthCheckMethod, ok := d.GetOk("health_check_method"); ok && healthCheckMethod.(string) != "" {
				req.QueryParams["HealthCheckMethod"] = healthCheckMethod.(string)
			}
		}
	}

	return req, nil
}

func buildHttpForwardArgs(d *schema.ResourceData, req *requests.CommonRequest) (*requests.CommonRequest, error) {
	stickySession := string(OffFlag)
	healthCheck := string(OffFlag)
	listenerForward := string(OnFlag)
	req.QueryParams["StickySession"] = stickySession
	req.QueryParams["HealthCheck"] = healthCheck
	req.QueryParams["ListenerForward"] = listenerForward
	/**
	if the user do not fill backend_port, give 80 to pass the SDK parameter check.
	*/
	if backEndServerPort, ok := d.GetOk("backend_port"); ok {
		req.QueryParams[""] = string(requests.NewInteger(backEndServerPort.(int)))
	} else {
		req.QueryParams["BackendServerPort"] = string("80")
	}
	if forwardPort, ok := d.GetOk("forward_port"); ok {
		req.QueryParams["ForwardPort"] = string(requests.NewInteger(forwardPort.(int)))
	}
	return req, nil
}

func parseListenerId(d *schema.ResourceData, meta interface{}) (string, string, int, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	parts, err := ParseSlbListenerId(d.Id())
	if err != nil {
		return "", "", 0, WrapError(err)
	}
	protocol := ""
	port := 0
	if len(parts) == 3 {
		protocol = parts[1]
		port, err = strconv.Atoi(parts[2])
	} else {
		if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
			protocol = v.(string)
		}
		port, err = strconv.Atoi(parts[1])
	}
	if err != nil {
		return "", "", 0, WrapError(err)
	}
	loadBalancer, err := slbService.DescribeSlb(parts[0])
	if err != nil {
		return "", "", 0, WrapError(err)
	}
	if protocol != "" {
		for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if portAndProtocol.ListenerPort == port && portAndProtocol.ListenerProtocol == protocol {
				return loadBalancer.LoadBalancerId, portAndProtocol.ListenerProtocol, port, nil
			}
		}
	} else {
		if len(loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol) > 1 {
			return "", "", 0, WrapError(Error("More than one listener was with with the same id: %s, please specify protocol.", d.Id()))
		}
		for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if portAndProtocol.ListenerPort == port {
				return loadBalancer.LoadBalancerId, portAndProtocol.ListenerProtocol, port, nil
			}
		}
	}
	return "", "", 0, GetNotFoundErrorFromString(GetNotFoundMessage("Listener", d.Id()))
}

func readListener(d *schema.ResourceData, listener map[string]interface{}) {
	if val, ok := listener["BackendServerPort"]; ok {
		d.Set("backend_port", val.(float64))
	}
	if val, ok := listener["ListenerPort"]; ok {
		d.Set("frontend_port", val.(float64))
	}
	if val, ok := listener["Bandwidth"]; ok {
		d.Set("bandwidth", val.(float64))
	}
	if val, ok := listener["Scheduler"]; ok {
		d.Set("scheduler", val.(string))
	}
	if val, ok := listener["VServerGroupId"]; ok {
		d.Set("server_group_id", val.(string))
	}
	if val, ok := listener["MasterSlaveServerGroupId"]; ok {
		d.Set("master_slave_server_group_id", val.(string))
	}
	if val, ok := listener["AclStatus"]; ok {
		d.Set("acl_status", val.(string))
	}
	if val, ok := listener["AclType"]; ok {
		d.Set("acl_type", val.(string))
	}
	if val, ok := listener["AclId"]; ok {
		d.Set("acl_id", val.(string))
	}
	if val, ok := listener["HealthCheck"]; ok {
		d.Set("health_check", val.(string))
	}
	if val, ok := listener["StickySession"]; ok {
		d.Set("sticky_session", val.(string))
	}
	if val, ok := listener["StickySessionType"]; ok {
		d.Set("sticky_session_type", val.(string))
	}
	if val, ok := listener["CookieTimeout"]; ok {
		d.Set("cookie_timeout", val.(float64))
	}
	if val, ok := listener["Cookie"]; ok {
		d.Set("cookie", val.(string))
	}
	if val, ok := listener["PersistenceTimeout"]; ok {
		d.Set("persistence_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckType"]; ok {
		d.Set("health_check_type", val.(string))
	}
	if val, ok := listener["EstablishedTimeout"]; ok {
		d.Set("established_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckDomain"]; ok {
		d.Set("health_check_domain", val.(string))
	}
	if val, ok := listener["HealthCheckMethod"]; ok {
		d.Set("health_check_method", val.(string))
	}
	if val, ok := listener["HealthCheckConnectPort"]; ok {
		d.Set("health_check_connect_port", val.(float64))
	}
	if val, ok := listener["HealthCheckURI"]; ok {
		d.Set("health_check_uri", val.(string))
	}
	if val, ok := listener["HealthyThreshold"]; ok {
		d.Set("healthy_threshold", val.(float64))
	}
	if val, ok := listener["UnhealthyThreshold"]; ok {
		d.Set("unhealthy_threshold", val.(float64))
	}
	if val, ok := listener["HealthCheckTimeout"]; ok {
		d.Set("health_check_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckConnectTimeout"]; ok {
		d.Set("health_check_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckInterval"]; ok {
		d.Set("health_check_interval", val.(float64))
	}
	if val, ok := listener["HealthCheckHttpCode"]; ok {
		d.Set("health_check_http_code", val.(string))
	}
	if val, ok := listener["ServerCertificateId"]; ok {
		d.Set("ssl_certificate_id", val.(string))
		d.Set("server_certificate_id", val.(string))
	}
	if val, ok := listener["CACertificateId"]; ok {
		d.Set("ca_certificate_id", val.(string))
	}
	if val, ok := listener["EnableHttp2"]; ok {
		d.Set("enable_http2", val.(string))
	}

	if val, ok := listener["TLSCipherPolicy"]; ok {
		d.Set("tls_cipher_policy", val.(string))
	}

	if val, ok := listener["IdleTimeout"]; ok {
		d.Set("idle_timeout", val.(float64))
	}

	if val, ok := listener["RequestTimeout"]; ok {
		d.Set("request_timeout", val.(float64))
	}

	if val, ok := listener["Gzip"]; ok {
		d.Set("gzip", val.(string) == string(OnFlag))
	}
	if val, ok := listener["ListenerForward"]; ok {
		d.Set("listener_forward", val.(string))
	}
	if val, ok := listener["ForwardPort"]; ok {
		d.Set("forward_port", val.(float64))
	}
	xff := make(map[string]interface{})
	if val, ok := listener["XForwardedFor"]; ok {
		xff["retrive_client_ip"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_SLBIP"]; ok {
		xff["retrive_slb_ip"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_SLBID"]; ok {
		xff["retrive_slb_id"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_proto"]; ok {
		xff["retrive_slb_proto"] = val.(string) == string(OnFlag)
	}

	if len(xff) > 0 {
		d.Set("x_forwarded_for", []map[string]interface{}{xff})
	}
	if val, ok := listener["Description"]; ok {
		d.Set("description", val.(string))
	}

	return
}

func getHealthCheckHttpCodeValue(d *schema.ResourceData) string {
	if v, ok := d.GetOk("health_check_http_code"); ok && v.(string) != "" {
		return v.(string)
	}
	// After the version 1.144.0, removes health_check_http_code Default and using Computed instead.
	return string(HTTP_2XX)
}
