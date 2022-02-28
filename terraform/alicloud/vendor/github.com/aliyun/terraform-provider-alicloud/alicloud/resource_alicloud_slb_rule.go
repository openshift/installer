package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbRuleCreate,
		Read:   resourceAliyunSlbRuleRead,
		Update: resourceAliyunSlbRuleUpdate,
		Delete: resourceAliyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Required:     true,
				ForceNew:     true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-slb-rule",
			},

			"listener_sync": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
				Optional:     true,
				Default:      string(OnFlag),
			},
			"scheduler": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"wrr", "wlc", "rr"}, false),
				Optional:         true,
				Default:          WRRScheduler,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cookie": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 200),
				Optional:         true,
				DiffSuppressFunc: slbRuleCookieDiffSuppressFunc,
			},
			"cookie_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 86400),
				Optional:         true,
				DiffSuppressFunc: slbRuleCookieTimeoutDiffSuppressFunc,
			},
			"health_check": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			"health_check_http_code": {
				Type: schema.TypeString,
				ValidateFunc: validateAllowedSplitStringValue([]string{
					string(HTTP_2XX), string(HTTP_3XX), string(HTTP_4XX), string(HTTP_5XX)}, ","),
				Optional:         true,
				Default:          HTTP_2XX,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 50),
				Optional:         true,
				Default:          2,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_domain": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 80),
				Optional:         true,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_uri": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringLenBetween(1, 80),
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type: schema.TypeInt,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 65535),
					validation.IntInSlice([]int{-520})),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"health_check_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 300),
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: slbRuleHealthCheckDiffSuppressFunc,
			},
			//http & https
			"sticky_session": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"on", "off"}, false),
				Optional:         true,
				Default:          OffFlag,
				DiffSuppressFunc: slbRuleListenerSyncDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(InsertStickySessionType),
					string(ServerStickySessionType)}, false),
				Optional:         true,
				DiffSuppressFunc: slbRuleStickySessionTypeDiffSuppressFunc,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slb_id := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	name := strings.Trim(d.Get("name").(string), " ")
	group_id := strings.Trim(d.Get("server_group_id").(string), " ")

	var domain, url, rule string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		url = v.(string)
	}

	if domain == "" && url == "" {
		return WrapError(Error("At least one 'domain' or 'url' must be set."))
	} else if domain == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Url':'%s','VServerGroupId':'%s'}]", name, url, group_id)
	} else if url == "" {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','VServerGroupId':'%s'}]", name, domain, group_id)
	} else {
		rule = fmt.Sprintf("[{'RuleName':'%s','Domain':'%s','Url':'%s','VServerGroupId':'%s'}]", name, domain, url, group_id)
	}

	request := slb.CreateCreateRulesRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = slb_id
	request.ListenerPort = requests.NewInteger(port)
	request.RuleList = rule
	var raw interface{}
	var err error
	if err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateRules(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"BackendServer.configuring", "OperationFailed.ListenerStatusNotSupport"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*slb.CreateRulesResponse)
	d.SetId(response.Rules.Rule[0].RuleId)

	return resourceAliyunSlbRuleUpdate(d, meta)
}

func resourceAliyunSlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbRule(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.RuleName)
	d.Set("load_balancer_id", object.LoadBalancerId)
	if port, err := strconv.Atoi(object.ListenerPort); err != nil {
		return WrapError(err)
	} else {
		d.Set("frontend_port", port)
	}
	d.Set("domain", object.Domain)
	d.Set("url", object.Url)
	d.Set("server_group_id", object.VServerGroupId)
	d.Set("sticky_session", object.StickySession)
	d.Set("sticky_session_type", object.StickySessionType)
	d.Set("unhealthy_threshold", object.UnhealthyThreshold)
	d.Set("healthy_threshold", object.HealthyThreshold)
	d.Set("health_check_timeout", object.HealthCheckTimeout)
	d.Set("health_check_connect_port", object.HealthCheckConnectPort)
	d.Set("health_check_uri", object.HealthCheckURI)
	d.Set("health_check", object.HealthCheck)
	d.Set("health_check_http_code", object.HealthCheckHttpCode)
	d.Set("health_check_interval", object.HealthCheckInterval)
	d.Set("scheduler", object.Scheduler)
	d.Set("listener_sync", object.ListenerSync)
	d.Set("cookie_timeout", object.CookieTimeout)
	d.Set("cookie", object.Cookie)
	d.Set("health_check_domain", object.HealthCheckDomain)
	return nil
}

func resourceAliyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	update := false
	fullUpdate := false
	request := slb.CreateSetRuleRequest()
	request.RuleId = d.Id()
	if listenerSync, ok := d.GetOk("listener_sync"); ok && listenerSync == string(OffFlag) {
		if stickySession := d.Get("sticky_session"); stickySession == string(OnFlag) {
			if _, ok := d.GetOk("sticky_session_type"); !ok {
				return WrapError(Error(`'sticky_session_type': required field is not set when the sticky_session is 'on'.`))
			}
		}
		if stickySessionType := d.Get("sticky_session_type"); stickySessionType == string(InsertStickySessionType) {
			if _, ok := d.GetOk("cookie_timeout"); !ok {
				return WrapError(Error(`'cookie_timeout': required field is not set when the sticky_session_type is 'insert'.`))
			}
		}
		if stickySessionType := d.Get("sticky_session_type"); stickySessionType == string(ServerStickySessionType) {
			if _, ok := d.GetOk("cookie"); !ok {
				return WrapError(Error(`'cookie': required field is not set when the sticky_session_type is 'server'.`))
			}
		}
	}
	if d.HasChange("server_group_id") {
		request.VServerGroupId = d.Get("server_group_id").(string)
		update = true
	}

	if d.HasChange("name") {
		request.RuleName = d.Get("name").(string)
		update = true
	}

	fullUpdate = d.HasChange("listener_sync") || d.HasChange("scheduler") || d.HasChange("cookie") || d.HasChange("cookie_timeout") || d.HasChange("health_check") || d.HasChange("health_check_http_code") ||
		d.HasChange("health_check_interval") || d.HasChange("health_check_domain") || d.HasChange("health_check_uri") || d.HasChange("health_check_connect_port") || d.HasChange("health_check_timeout") ||
		d.HasChange("healthy_threshold") || d.HasChange("unhealthy_threshold") || d.HasChange("sticky_session") || d.HasChange("sticky_session_type")

	if fullUpdate {
		request.ListenerSync = d.Get("listener_sync").(string)
		if listenerSync, ok := d.GetOk("listener_sync"); ok && listenerSync == string(OffFlag) {
			request.Scheduler = d.Get("scheduler").(string)
			request.HealthCheck = d.Get("health_check").(string)
			request.StickySession = d.Get("sticky_session").(string)
			if request.HealthCheck == string(OnFlag) {
				request.HealthCheckTimeout = requests.NewInteger(d.Get("health_check_timeout").(int))
				request.HealthCheckURI = d.Get("health_check_uri").(string)
				request.HealthyThreshold = requests.NewInteger(d.Get("healthy_threshold").(int))
				request.UnhealthyThreshold = requests.NewInteger(d.Get("unhealthy_threshold").(int))
				request.HealthCheckInterval = requests.NewInteger(d.Get("health_check_interval").(int))
				request.HealthCheckHttpCode = d.Get("health_check_http_code").(string)
				if healthCheckDomain, ok := d.GetOk("health_check_domain"); ok {
					request.HealthCheckDomain = healthCheckDomain.(string)
				}
				if healthCheckConnectPort, ok := d.GetOk("health_check_connect_port"); ok {
					request.HealthCheckConnectPort = requests.NewInteger(healthCheckConnectPort.(int))
				}
			}
			if request.StickySession == string(OnFlag) {
				request.StickySessionType = d.Get("sticky_session_type").(string)
				if request.StickySessionType == string(InsertStickySessionType) {
					request.CookieTimeout = requests.NewInteger(d.Get("cookie_timeout").(int))
				}
				if request.StickySessionType == string(ServerStickySessionType) {
					request.Cookie = d.Get("cookie").(string)
				}
			}
		}

	}
	if update || fullUpdate {
		client := meta.(*connectivity.AliyunClient)
		request.RegionId = client.RegionId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetRule(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunSlbRuleRead(d, meta)
}

func resourceAliyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current rule's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the rule.", lbId))
		}
	}

	request := slb.CreateDeleteRulesRequest()
	request.RegionId = client.RegionId
	request.RuleIds = fmt.Sprintf("['%s']", d.Id())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteRules(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.ListenerStatusNotSupport"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRuleId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(slbService.WaitForSlbRule(d.Id(), Deleted, DefaultTimeoutMedium))

}
