package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbServerGroupCreate,
		Read:   resourceAlicloudAlbServerGroupRead,
		Update: resourceAlicloudAlbServerGroupUpdate,
		Delete: resourceAlicloudAlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_http_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
						},
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 50),
						},
						"health_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD"}, false),
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"health_check_codes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTPS", "HTTP"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Sch", "Wlc", "Wrr"}, false),
			},
			"server_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"server_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Ecs", "Eni", "Eci"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sticky_session_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cookie_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.Any(validation.IntInSlice([]int{0}), validation.IntBetween(1, 86400)),
						},
						"sticky_session_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sticky_session_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"Insert", "Server"}, false),
						},
					},
				},
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudAlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServerGroup"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOk("server_group_name"); ok {
		request["ServerGroupName"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("health_check_config"); ok {

		healthCheckConfig := make(map[string]interface{})
		for _, healthCheckConfigArgs := range v.(*schema.Set).List() {
			healthCheckConfigArg := healthCheckConfigArgs.(map[string]interface{})

			healthCheckConfig["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]

			if healthCheckConfig["HealthCheckEnabled"] == true {
				healthCheckConfig["HealthCheckCodes"] = healthCheckConfigArg["health_check_codes"]
				healthCheckConfig["HealthCheckConnectPort"] = healthCheckConfigArg["health_check_connect_port"]
				healthCheckConfig["HealthCheckHost"] = healthCheckConfigArg["health_check_host"]
				healthCheckConfig["HealthCheckHttpVersion"] = healthCheckConfigArg["health_check_http_version"]
				healthCheckConfig["HealthCheckInterval"] = healthCheckConfigArg["health_check_interval"]
				healthCheckConfig["HealthCheckMethod"] = healthCheckConfigArg["health_check_method"]
				healthCheckConfig["HealthCheckPath"] = healthCheckConfigArg["health_check_path"]
				healthCheckConfig["HealthCheckProtocol"] = healthCheckConfigArg["health_check_protocol"]
				healthCheckConfig["HealthCheckTimeout"] = healthCheckConfigArg["health_check_timeout"]
				healthCheckConfig["HealthyThreshold"] = healthCheckConfigArg["healthy_threshold"]
				healthCheckConfig["UnhealthyThreshold"] = healthCheckConfigArg["unhealthy_threshold"]
			}
		}
		request["HealthCheckConfig"] = healthCheckConfig
	}
	if v, ok := d.GetOk("sticky_session_config"); ok {

		stickySessionConfig := make(map[string]interface{})
		for _, stickySessionConfigArgs := range v.(*schema.Set).List() {
			stickySessionConfigArg := stickySessionConfigArgs.(map[string]interface{})
			stickySessionConfig["StickySessionEnabled"] = stickySessionConfigArg["sticky_session_enabled"]
			if stickySessionConfig["StickySessionEnabled"] == true {
				stickySessionConfig["StickySessionType"] = stickySessionConfigArg["sticky_session_type"]
				if stickySessionConfig["StickySessionType"] == "Server" {
					stickySessionConfig["Cookie"] = stickySessionConfigArg["cookie"]
				}
				if stickySessionConfig["StickySessionType"] == "Insert" {
					stickySessionConfig["CookieTimeout"] = stickySessionConfigArg["cookie_timeout"]
				}
			}
		}

		request["StickySessionConfig"] = stickySessionConfig
	}

	request["ClientToken"] = buildClientToken("CreateServerGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_server_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))
	albService := AlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlbServerGroupUpdate(d, meta)
}
func resourceAlicloudAlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbServerGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_server_group albService.DescribeAlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	healthCheckConfigSli := make([]map[string]interface{}, 0)
	if len(object["HealthCheckConfig"].(map[string]interface{})) > 0 {
		healthCheckConfig := object["HealthCheckConfig"]
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckConfigMap["health_check_codes"] = healthCheckConfig.(map[string]interface{})["HealthCheckCodes"]
		healthCheckConfigMap["health_check_connect_port"] = formatInt(healthCheckConfig.(map[string]interface{})["HealthCheckConnectPort"])
		healthCheckConfigMap["health_check_enabled"] = healthCheckConfig.(map[string]interface{})["HealthCheckEnabled"]
		healthCheckConfigMap["health_check_host"] = healthCheckConfig.(map[string]interface{})["HealthCheckHost"]
		healthCheckConfigMap["health_check_http_version"] = healthCheckConfig.(map[string]interface{})["HealthCheckHttpVersion"]
		healthCheckConfigMap["health_check_interval"] = formatInt(healthCheckConfig.(map[string]interface{})["HealthCheckInterval"])
		healthCheckConfigMap["health_check_method"] = healthCheckConfig.(map[string]interface{})["HealthCheckMethod"]
		healthCheckConfigMap["health_check_path"] = healthCheckConfig.(map[string]interface{})["HealthCheckPath"]
		healthCheckConfigMap["health_check_protocol"] = healthCheckConfig.(map[string]interface{})["HealthCheckProtocol"]
		healthCheckConfigMap["health_check_timeout"] = formatInt(healthCheckConfig.(map[string]interface{})["HealthCheckTimeout"])
		healthCheckConfigMap["healthy_threshold"] = formatInt(healthCheckConfig.(map[string]interface{})["HealthyThreshold"])
		healthCheckConfigMap["unhealthy_threshold"] = formatInt(healthCheckConfig.(map[string]interface{})["UnhealthyThreshold"])
		healthCheckConfigSli = append(healthCheckConfigSli, healthCheckConfigMap)
	}
	d.Set("health_check_config", healthCheckConfigSli)
	d.Set("protocol", object["Protocol"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("scheduler", object["Scheduler"])
	d.Set("server_group_name", object["ServerGroupName"])
	d.Set("status", object["ServerGroupStatus"])

	stickySessionConfigSli := make([]map[string]interface{}, 0)
	if len(object["StickySessionConfig"].(map[string]interface{})) > 0 {
		stickySessionConfig := object["StickySessionConfig"]
		stickySessionConfigMap := make(map[string]interface{})
		stickySessionConfigMap["cookie"] = stickySessionConfig.(map[string]interface{})["Cookie"]
		stickySessionConfigMap["cookie_timeout"] = stickySessionConfig.(map[string]interface{})["CookieTimeout"]
		stickySessionConfigMap["sticky_session_enabled"] = stickySessionConfig.(map[string]interface{})["StickySessionEnabled"]
		stickySessionConfigMap["sticky_session_type"] = stickySessionConfig.(map[string]interface{})["StickySessionType"]
		stickySessionConfigSli = append(stickySessionConfigSli, stickySessionConfigMap)
	}
	d.Set("sticky_session_config", stickySessionConfigSli)
	d.Set("vpc_id", object["VpcId"])
	serversList, err := albService.ListServerGroupServers(d.Id())
	serversMaps := make([]map[string]interface{}, 0)
	if serversList != nil {
		for _, serversListItem := range serversList {
			if serversListItemMap, ok := serversListItem.(map[string]interface{}); ok {
				serversMap := make(map[string]interface{}, 0)
				serversMap["description"] = serversListItemMap["Description"]
				serversMap["port"] = serversListItemMap["Port"]
				serversMap["server_id"] = serversListItemMap["ServerId"]
				serversMap["server_ip"] = serversListItemMap["ServerIp"]
				serversMap["server_type"] = serversListItemMap["ServerType"]
				serversMap["status"] = serversListItemMap["Status"]
				serversMap["weight"] = serversListItemMap["Weight"]
				serversMaps = append(serversMaps, serversMap)
			}
		}
	}
	d.Set("servers", serversMaps)

	listTagResourcesObject, err := albService.ListTagResources(d.Id(), "servergroup")
	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}
func resourceAlicloudAlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}
	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ResourceType"] = "servergroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	update = false
	if d.HasChange("servers") {

		removed, added := d.GetChange("servers")
		removeServersFromServerGroupReq := map[string]interface{}{
			"ServerGroupId": d.Id(),
		}

		removeServersMaps := make([]map[string]interface{}, 0)
		for _, servers := range removed.(*schema.Set).List() {
			update = true
			serversMap := map[string]interface{}{}
			serversArg := servers.(map[string]interface{})
			serversMap["Description"] = serversArg["description"]
			serversMap["Port"] = serversArg["port"]
			serversMap["ServerId"] = serversArg["server_id"]
			serversMap["ServerIp"] = serversArg["server_ip"]
			serversMap["ServerType"] = serversArg["server_type"]
			serversMap["Weight"] = serversArg["weight"]
			removeServersMaps = append(removeServersMaps, serversMap)
		}
		removeServersFromServerGroupReq["Servers"] = removeServersMaps

		if update {

			action := "RemoveServersFromServerGroup"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}

			removeServersFromServerGroupReq["ClientToken"] = buildClientToken(action)
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, removeServersFromServerGroupReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeServersFromServerGroupReq)

			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		update = false
		addServersToServerGroupReq := map[string]interface{}{
			"ServerGroupId": d.Id(),
		}
		addServersMaps := make([]map[string]interface{}, 0)
		for _, servers := range added.(*schema.Set).List() {
			update = true
			serversArg := servers.(map[string]interface{})
			serversMap := map[string]interface{}{}
			serversMap["Description"] = serversArg["description"]
			serversMap["Port"] = serversArg["port"]
			serversMap["ServerId"] = serversArg["server_id"]
			serversMap["ServerIp"] = serversArg["server_ip"]
			serversMap["ServerType"] = serversArg["server_type"]
			serversMap["Weight"] = serversArg["weight"]
			addServersMaps = append(addServersMaps, serversMap)
		}
		addServersToServerGroupReq["Servers"] = addServersMaps

		if update {

			action := "AddServersToServerGroup"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			addServersToServerGroupReq["ClientToken"] = buildClientToken(action)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, addServersToServerGroupReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addServersToServerGroupReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

			d.SetPartial("servers")
		}
	}

	update = false
	updateServerGroupAttributeReq := map[string]interface{}{
		"ServerGroupId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		if v, ok := d.GetOk("scheduler"); ok {
			updateServerGroupAttributeReq["Scheduler"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
		if v, ok := d.GetOk("server_group_name"); ok {
			updateServerGroupAttributeReq["ServerGroupName"] = v
		}
	}
	if d.HasChange("health_check_config") {
		update = true
		if v, ok := d.GetOk("health_check_config"); ok {
			healthCheckConfig := make(map[string]interface{})
			for _, healthCheckConfigArgs := range v.(*schema.Set).List() {
				healthCheckConfigArg := healthCheckConfigArgs.(map[string]interface{})
				healthCheckConfig["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]
				if healthCheckConfig["HealthCheckEnabled"] == true {
					healthCheckConfig["HealthCheckCodes"] = healthCheckConfigArg["health_check_codes"]
					healthCheckConfig["HealthCheckConnectPort"] = healthCheckConfigArg["health_check_connect_port"]
					healthCheckConfig["HealthCheckHost"] = healthCheckConfigArg["health_check_host"]
					healthCheckConfig["HealthCheckHttpVersion"] = healthCheckConfigArg["health_check_http_version"]
					healthCheckConfig["HealthCheckInterval"] = healthCheckConfigArg["health_check_interval"]
					healthCheckConfig["HealthCheckMethod"] = healthCheckConfigArg["health_check_method"]
					healthCheckConfig["HealthCheckPath"] = healthCheckConfigArg["health_check_path"]
					healthCheckConfig["HealthCheckProtocol"] = healthCheckConfigArg["health_check_protocol"]
					healthCheckConfig["HealthCheckTimeout"] = healthCheckConfigArg["health_check_timeout"]
					healthCheckConfig["HealthyThreshold"] = healthCheckConfigArg["healthy_threshold"]
					healthCheckConfig["UnhealthyThreshold"] = healthCheckConfigArg["unhealthy_threshold"]
				}
			}
			updateServerGroupAttributeReq["HealthCheckConfig"] = healthCheckConfig
		}

	}
	if d.HasChange("sticky_session_config") {
		update = true
		if v, ok := d.GetOk("sticky_session_config"); ok {
			stickySessionConfig := make(map[string]interface{})
			for _, stickySessionConfigArgs := range v.(*schema.Set).List() {
				stickySessionConfigArg := stickySessionConfigArgs.(map[string]interface{})
				stickySessionConfig["StickySessionEnabled"] = stickySessionConfigArg["sticky_session_enabled"]
				if stickySessionConfig["StickySessionEnabled"] == true {
					stickySessionConfig["StickySessionType"] = stickySessionConfigArg["sticky_session_type"]
					if stickySessionConfig["StickySessionType"] == "Server" {
						stickySessionConfig["Cookie"] = stickySessionConfigArg["cookie"]
					}
					if stickySessionConfig["StickySessionType"] == "Insert" {
						stickySessionConfig["CookieTimeout"] = stickySessionConfigArg["cookie_timeout"]
					}
				}
			}
			updateServerGroupAttributeReq["StickySessionConfig"] = stickySessionConfig
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateServerGroupAttributeReq["DryRun"] = v
		}
		action := "UpdateServerGroupAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateServerGroupAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateServerGroupAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup", "ResourceNotFound.ServerGroup"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateServerGroupAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("scheduler")
		d.SetPartial("server_group_name")
	}
	d.Partial(false)
	return resourceAlicloudAlbServerGroupRead(d, meta)
}
func resourceAlicloudAlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServerGroup"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServerGroupId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteServerGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup", "ResourceInUse.ServerGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.ServerGroup"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
