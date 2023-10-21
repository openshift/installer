package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSaeApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeApplicationCreate,
		Read:   resourceAlicloudSaeApplicationRead,
		Update: resourceAlicloudSaeApplicationUpdate,
		Delete: resourceAlicloudSaeApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_config": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_enable_application_scaling_rule": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"batch_wait_time": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"change_order_desc": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_args": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_map_mount_desc": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1000, 16000, 2000, 32000, 4000, 500, 8000}),
			},
			"custom_host_alias": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"edas_container_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ahas": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enable_grey_tag_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"envs": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jar_start_args": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jar_start_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"liveness": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1024, 131072, 16384, 2048, 32768, 4096, 65536, 8192}),
			},
			"min_ready_instances": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"mount_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nas_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_ak_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oss_ak_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oss_mount_descs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FatJar", "Image", "War"}, false),
			},
			"package_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"php_arms_config_location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"php_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"php_config_location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"post_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_stop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"readiness": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_configs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"RUNNING", "STOPPED", "UNKNOWN"}, false),
			},
			"termination_grace_period_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tomcat_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"war_start_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"intranet": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
							Optional:     true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"internet": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
							Optional:     true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"internet_slb_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"intranet_slb_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intranet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSaeApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/app/createApplication"
	request := make(map[string]*string)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request["AppName"] = StringPointer(d.Get("app_name").(string))
	request["PackageType"] = StringPointer(d.Get("package_type").(string))
	request["Replicas"] = StringPointer(strconv.Itoa(d.Get("replicas").(int)))
	if v, ok := d.GetOk("app_description"); ok {
		request["AppDescription"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOkExists("auto_config"); ok {
		request["AutoConfig"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("command"); ok {
		request["Command"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("command_args"); ok {
		request["CommandArgs"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("config_map_mount_desc"); ok {
		request["ConfigMapMountDesc"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("cpu"); ok {
		request["Cpu"] = StringPointer(strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("custom_host_alias"); ok {
		request["CustomHostAlias"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOkExists("deploy"); ok {
		request["Deploy"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("edas_container_version"); ok {
		request["EdasContainerVersion"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("envs"); ok {
		request["Envs"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("image_url"); ok {
		request["ImageUrl"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("jar_start_args"); ok {
		request["JarStartArgs"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("jar_start_options"); ok {
		request["JarStartOptions"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("jdk"); ok {
		request["Jdk"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("liveness"); ok {
		request["Liveness"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("memory"); ok {
		request["Memory"] = StringPointer(strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("mount_desc"); ok {
		request["MountDesc"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("mount_host"); ok {
		request["MountHost"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("namespace_id"); ok {
		request["NamespaceId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("nas_id"); ok {
		request["NasId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("oss_ak_id"); ok {
		request["OssAkId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("oss_ak_secret"); ok {
		request["OssAkSecret"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("oss_mount_descs"); ok {
		request["OssMountDescs"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("package_url"); ok {
		request["PackageUrl"] = StringPointer(v.(string))
	}
	request["PackageVersion"] = StringPointer(strconv.FormatInt(time.Now().Unix(), 10))
	if v, ok := d.GetOk("php_arms_config_location"); ok {
		request["PhpArmsConfigLocation"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("php_config"); ok {
		request["PhpConfig"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("php_config_location"); ok {
		request["PhpConfigLocation"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("post_start"); ok {
		request["PostStart"] = StringPointer(v.(string))
	}
	if d.HasChange("pre_stop") {
		request["PreStop"] = StringPointer(d.Get("pre_stop").(string))
	}
	if v, ok := d.GetOk("pre_stop"); ok {
		request["PreStop"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("readiness"); ok {
		request["Readiness"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("sls_configs"); ok {
		request["SlsConfigs"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		request["TerminationGracePeriodSeconds"] = StringPointer(strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("timezone"); ok {
		request["Timezone"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("tomcat_config"); ok {
		request["TomcatConfig"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("war_start_options"); ok {
		request["WarStartOptions"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("web_container"); ok {
		request["WebContainer"] = StringPointer(v.(string))
	}
	if v, exist := d.GetOk("vpc_id"); exist {
		request["VpcId"] = StringPointer(v.(string))
	}
	if v, exist := d.GetOk("vswitch_id"); exist {
		request["VSwitchId"] = StringPointer(v.(string))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_application", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["AppId"]))

	return resourceAlicloudSaeApplicationUpdate(d, meta)
}
func resourceAlicloudSaeApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_application saeService.DescribeSaeApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("app_description", object["AppDescription"])
	d.Set("app_name", object["AppName"])
	d.Set("command", object["Command"])
	d.Set("command_args", object["CommandArgs"])

	d.Set("config_map_mount_desc", object["ConfigMapMountDesc"])
	if v, ok := object["Cpu"]; ok && fmt.Sprint(v) != "0" {
		d.Set("cpu", formatInt(v))
	}
	d.Set("custom_host_alias", object["CustomHostAlias"])
	d.Set("edas_container_version", object["EdasContainerVersion"])
	d.Set("envs", object["Envs"])
	d.Set("image_url", object["ImageUrl"])
	d.Set("jar_start_args", object["JarStartArgs"])
	d.Set("jar_start_options", object["JarStartOptions"])
	d.Set("jdk", object["Jdk"])
	d.Set("liveness", object["Liveness"])
	if v, ok := object["Memory"]; ok && fmt.Sprint(v) != "0" {
		d.Set("memory", formatInt(v))
	}
	if v, ok := object["MinReadyInstances"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_ready_instances", formatInt(v))
	}

	if v, ok := object["MountDesc"].([]interface{}); ok {
		mountDesc, err := convertListObjectToCommaSeparate(v)
		if err != nil {
			return WrapError(err)
		}
		d.Set("mount_desc", mountDesc)
	}

	d.Set("mount_host", object["MountHost"])
	d.Set("namespace_id", object["NamespaceId"])
	d.Set("nas_id", object["NasId"])
	d.Set("oss_ak_id", object["OssAkId"])
	d.Set("oss_ak_secret", object["OssAkSecret"])
	d.Set("oss_mount_descs", object["OssMountDescs"])
	d.Set("package_type", object["PackageType"])
	d.Set("package_url", object["PackageUrl"])
	d.Set("package_version", object["PackageVersion"])
	d.Set("php_arms_config_location", object["PhpArmsConfigLocation"])
	d.Set("php_config", object["PhpConfig"])
	d.Set("php_config_location", object["PhpConfigLocation"])
	d.Set("post_start", object["PostStart"])
	d.Set("pre_stop", object["PreStop"])
	d.Set("readiness", object["Readiness"])
	if v, ok := object["Replicas"]; ok && fmt.Sprint(v) != "0" {
		d.Set("replicas", formatInt(v))
	}
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("sls_configs", object["SlsConfigs"])
	if v, ok := object["TerminationGracePeriodSeconds"]; ok && fmt.Sprint(v) != "0" {
		d.Set("termination_grace_period_seconds", formatInt(v))
	}
	d.Set("timezone", object["Timezone"])
	d.Set("tomcat_config", object["TomcatConfig"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("war_start_options", object["WarStartOptions"])
	d.Set("web_container", object["WebContainer"])
	describeApplicationStatusObject, err := saeService.DescribeApplicationStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("status", describeApplicationStatusObject["CurrentStatus"])

	describeApplicationSlbObject, err := saeService.DescribeApplicationSlb(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("internet_ip", describeApplicationSlbObject["InternetIp"])
	d.Set("intranet_ip", describeApplicationSlbObject["IntranetIp"])
	d.Set("intranet_slb_id", describeApplicationSlbObject["IntranetSlbId"])
	d.Set("internet_slb_id", describeApplicationSlbObject["InternetSlbId"])
	intranetArray := make([]interface{}, 0)
	if v, ok := describeApplicationSlbObject["Intranet"]; ok {
		for _, intranet := range v.([]interface{}) {
			intranetObject := intranet.(map[string]interface{})
			intranetObj := map[string]interface{}{
				"https_cert_id": intranetObject["HttpsCertId"],
				"protocol":      intranetObject["Protocol"],
				"target_port":   intranetObject["TargetPort"],
				"port":          intranetObject["Port"],
			}
			intranetArray = append(intranetArray, intranetObj)
		}
	}
	d.Set("intranet", intranetArray)

	internetArray := make([]interface{}, 0)
	if v, ok := describeApplicationSlbObject["Internet"]; ok {
		for _, internet := range v.([]interface{}) {
			internetObject := internet.(map[string]interface{})
			internetObj := map[string]interface{}{
				"https_cert_id": internetObject["HttpsCertId"],
				"protocol":      internetObject["Protocol"],
				"target_port":   internetObject["TargetPort"],
				"port":          internetObject["Port"],
			}
			internetArray = append(internetArray, internetObj)
		}
	}
	d.Set("internet", internetArray)

	return nil
}
func resourceAlicloudSaeApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	saeService := SaeService{client}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	request := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	update = false
	if d.HasChange("replicas") && !d.IsNewResource() {
		update = true
	}
	request["Replicas"] = StringPointer(strconv.Itoa(d.Get("replicas").(int)))
	if update {
		action := "/pop/v1/sam/app/rescaleApplication"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
	}
	update = false
	//DeployApplication
	request = map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	if d.HasChange("command") {
		update = true
	}
	if v, ok := d.GetOk("command"); ok {
		request["Command"] = StringPointer(v.(string))
	}

	if d.HasChange("command_args") {
		update = true
	}
	if v, ok := d.GetOk("command_args"); ok {
		request["CommandArgs"] = StringPointer(v.(string))
	}

	if d.HasChange("config_map_mount_desc") {
		update = true
	}
	if v, ok := d.GetOk("config_map_mount_desc"); ok {
		request["ConfigMapMountDesc"] = StringPointer(v.(string))
	}

	if d.HasChange("custom_host_alias") {
		update = true
	}
	if v, ok := d.GetOk("custom_host_alias"); ok {
		request["CustomHostAlias"] = StringPointer(v.(string))
	}

	if d.HasChange("edas_container_version") {
		update = true
	}
	if v, ok := d.GetOk("edas_container_version"); ok {
		request["EdasContainerVersion"] = StringPointer(v.(string))
	}

	if d.HasChange("envs") {
		update = true
	}
	if v, ok := d.GetOk("envs"); ok {
		request["Envs"] = StringPointer(v.(string))
	}

	if d.HasChange("image_url") {
		update = true
	}
	if v, ok := d.GetOk("image_url"); ok {
		request["ImageUrl"] = StringPointer(v.(string))
	}

	if d.HasChange("jar_start_args") {
		update = true
	}
	if v, ok := d.GetOk("jar_start_args"); ok {
		request["JarStartArgs"] = StringPointer(v.(string))
	}

	if d.HasChange("jar_start_options") {
		update = true
	}
	if v, ok := d.GetOk("jar_start_options"); ok {
		request["JarStartOptions"] = StringPointer(v.(string))
	}

	if d.HasChange("jdk") {
		update = true
	}
	if v, ok := d.GetOk("jdk"); ok {
		request["Jdk"] = StringPointer(v.(string))
	}

	if d.HasChange("liveness") {
		update = true
	}
	if v, ok := d.GetOk("liveness"); ok {
		request["Liveness"] = StringPointer(v.(string))
	}

	if d.HasChange("mount_desc") {
		update = true
	}
	if v, ok := d.GetOk("mount_desc"); ok {
		request["MountDesc"] = StringPointer(v.(string))
	}

	if d.HasChange("mount_host") {
		update = true
	}
	if v, ok := d.GetOk("mount_host"); ok {
		request["MountHost"] = StringPointer(v.(string))
	}

	if d.HasChange("nas_id") {
		update = true
	}
	if v, ok := d.GetOk("nas_id"); ok {
		request["NasId"] = StringPointer(v.(string))
	}

	if d.HasChange("oss_ak_id") {
		update = true
	}
	if v, ok := d.GetOk("oss_ak_id"); ok {
		request["OssAkId"] = StringPointer(v.(string))
	}

	if d.HasChange("oss_ak_secret") {
		update = true
	}
	if v, ok := d.GetOk("oss_ak_secret"); ok {
		request["OssAkSecret"] = StringPointer(v.(string))
	}

	if d.HasChange("oss_mount_descs") {
		update = true
	}
	if v, ok := d.GetOk("oss_mount_descs"); ok {
		request["OssMountDescs"] = StringPointer(v.(string))
	}

	if d.HasChange("package_url") {
		update = true
	}
	if v, ok := d.GetOk("package_url"); ok {
		request["PackageUrl"] = StringPointer(v.(string))
	}

	if d.HasChange("php_arms_config_location") {
		update = true
	}
	if v, ok := d.GetOk("php_arms_config_location"); ok {
		request["PhpArmsConfigLocation"] = StringPointer(v.(string))
	}

	if d.HasChange("php_config") {
		update = true
	}
	if v, ok := d.GetOk("php_config"); ok {
		request["PhpConfig"] = StringPointer(v.(string))
	}

	if d.HasChange("php_config_location") {
		update = true
	}
	if v, ok := d.GetOk("php_config_location"); ok {
		request["PhpConfigLocation"] = StringPointer(v.(string))
	}

	if d.HasChange("post_start") {
		update = true
	}
	if v, ok := d.GetOk("post_start"); ok {
		request["PostStart"] = StringPointer(v.(string))
	}

	if d.HasChange("pre_stop") {
		update = true
	}
	if v, ok := d.GetOk("pre_stop"); ok {
		request["PreStop"] = StringPointer(v.(string))
	}

	if d.HasChange("readiness") {
		update = true
	}
	if v, ok := d.GetOk("readiness"); ok {
		request["Readiness"] = StringPointer(v.(string))
	}

	if d.HasChange("sls_configs") {
		update = true
	}
	if v, ok := d.GetOk("sls_configs"); ok {
		request["SlsConfigs"] = StringPointer(v.(string))
	}

	if d.HasChange("termination_grace_period_seconds") {
		update = true
	}
	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		request["TerminationGracePeriodSeconds"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if d.HasChange("timezone") {
		update = true
	}
	if v, ok := d.GetOk("timezone"); ok {
		request["Timezone"] = StringPointer(v.(string))
	}

	if d.HasChange("tomcat_config") {
		update = true
	}
	if v, ok := d.GetOk("tomcat_config"); ok {
		request["TomcatConfig"] = StringPointer(v.(string))
	}

	if d.HasChange("war_start_options") {
		update = true
	}
	if v, ok := d.GetOk("war_start_options"); ok {
		request["WarStartOptions"] = StringPointer(v.(string))
	}

	if d.HasChange("web_container") {
		update = true
	}
	if v, ok := d.GetOk("web_container"); ok {
		request["WebContainer"] = StringPointer(v.(string))
	}

	if d.HasChange("auto_enable_application_scaling_rule") {
		update = true
	}
	if v, ok := d.GetOk("auto_enable_application_scaling_rule"); ok {
		request["AutoEnableApplicationScalingRule"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	if d.HasChange("min_ready_instances") {
		update = true
	}
	if v, ok := d.GetOk("min_ready_instances"); ok {
		request["MinReadyInstances"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("batch_wait_time"); ok {
		request["BatchWaitTime"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if d.HasChange("change_order_desc") {
		update = true
	}
	if v, ok := d.GetOk("change_order_desc"); ok {
		request["ChangeOrderDesc"] = StringPointer(v.(string))
	}

	if d.HasChange("enable_ahas") {
		update = true
	}
	if v, ok := d.GetOk("enable_ahas"); ok {
		request["EnableAhas"] = StringPointer(v.(string))
	}

	if d.HasChange("enable_grey_tag_route") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_grey_tag_route"); ok {
		request["EnableGreyTagRoute"] = StringPointer(v.(string))
	}

	if d.HasChange("update_strategy") {
		update = true
	}
	if v, ok := d.GetOk("update_strategy"); ok {
		request["UpdateStrategy"] = StringPointer(v.(string))
	}

	if update {
		action := "/pop/v1/sam/app/deployApplication"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		request["PackageVersion"] = StringPointer(strconv.FormatInt(time.Now().Unix(), 10))
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		d.SetPartial("command")
		d.SetPartial("command_args")
		d.SetPartial("config_map_mount_desc")
		d.SetPartial("custom_host_alias")
		d.SetPartial("edas_container_version")
		d.SetPartial("envs")
		d.SetPartial("image_url")
		d.SetPartial("jar_start_args")
		d.SetPartial("jar_start_options")
		d.SetPartial("jdk")
		d.SetPartial("liveness")
		d.SetPartial("mount_desc")
		d.SetPartial("mount_host")
		d.SetPartial("nas_id")
		d.SetPartial("oss_ak_id")
		d.SetPartial("oss_ak_secret")
		d.SetPartial("oss_mount_descs")
		d.SetPartial("package_url")
		d.SetPartial("package_version")
		d.SetPartial("php_arms_config_location")
		d.SetPartial("php_config")
		d.SetPartial("php_config_location")
		d.SetPartial("post_start")
		d.SetPartial("pre_stop")
		d.SetPartial("readiness")
		d.SetPartial("sls_configs")
		d.SetPartial("min_ready_instances")
		d.SetPartial("auto_enable_application_scaling_rule")
		d.SetPartial("termination_grace_period_seconds")
		d.SetPartial("timezone")
		d.SetPartial("tomcat_config")
		d.SetPartial("war_start_options")
		d.SetPartial("web_container")
	}

	// update SLB
	update = false
	if err := saeService.UpdateSlb(d); err != nil {
		return WrapError(err)
	}

	// 【Exists】update security_group_id
	if d.HasChange("security_group_id") {
		d.Partial(true)
		request := map[string]*string{
			"AppId": StringPointer(d.Id()),
		}
		if v, ok := d.GetOk("security_group_id"); ok {
			request["SecurityGroupId"] = StringPointer(v.(string))
		}
		action := "/pop/v1/sam/app/updateAppSecurityGroup"
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "PUT "+action, response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "PUT "+action, response))
		}
		d.SetPartial("security_group_id")
	}

	// 【Exists】update rescaleApplicationVertically（CPU+Memory）
	if (d.HasChange("cpu") || d.HasChange("memory")) && !d.IsNewResource() {
		d.Partial(true)
		request := map[string]*string{
			"AppId": StringPointer(d.Id()),
		}
		if v, ok := d.GetOk("cpu"); ok {
			request["Cpu"] = StringPointer(strconv.Itoa(v.(int)))
		}
		if v, ok := d.GetOk("memory"); ok {
			request["Memory"] = StringPointer(strconv.Itoa(v.(int)))
		}
		action := "/pop/v1/sam/app/rescaleApplicationVertically"
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		d.SetPartial("cpu")
		d.SetPartial("memory")
	}
	//	【Exists】update status
	if d.HasChange("status") {
		d.Partial(true)
		object, err := saeService.DescribeApplicationStatus(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["CurrentStatus"].(string) != target {
			if target == "RUNNING" {
				request := map[string]*string{
					"AppId": StringPointer(d.Id()),
				}
				action := "/pop/v1/sam/app/startApplication"
				conn, err := client.NewServerlessClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
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
				if respBody, isExist := response["body"]; isExist {
					response = respBody.(map[string]interface{})
				} else {
					return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
				}
				if fmt.Sprint(response["Success"]) == "false" {
					return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
				}
			}
			if target == "STOPPED" {
				request := map[string]*string{
					"AppId": StringPointer(d.Id()),
				}
				action := "/pop/v1/sam/app/stopApplication"
				conn, err := client.NewServerlessClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.InvalidStatus"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
				}
				if respBody, isExist := response["body"]; isExist {
					response = respBody.(map[string]interface{})
				} else {
					return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
				}
				if fmt.Sprint(response["Success"]) == "false" {
					return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudSaeApplicationRead(d, meta)
}
func resourceAlicloudSaeApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/app/deleteApplication"
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}

	action = "/pop/v1/sam/app/describeApplicationConfig"
	request = map[string]*string{
		"AppId": StringPointer(d.Id()),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait = incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if response != nil {
			err = fmt.Errorf("application have not been destroyed yet")
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func convertListObjectToCommaSeparate(configured []interface{}) (string, error) {
	if len(configured) < 1 {
		return "", nil
	}
	result := "["
	for i, v := range configured {
		rail := ","
		if i == len(configured)-1 {
			rail = ""
		}
		vv, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		result += string(vv) + rail
	}
	return result + "]", nil
}
