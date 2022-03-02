package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSaeApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeApplicationsRead,
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"field_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"appName", "appIds", "slbIps", "instanceIps"}, false),
			},
			"field_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"order_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"running", "instances"}, false),
			},
			"reverse": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RUNNING", "STOPPED", "UNKNOWN"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_args": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_map_mount_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_origin_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_host_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"edas_container_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"envs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jar_start_args": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jar_start_options": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jdk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"liveness": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_ready_instances": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mount_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nas_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_ak_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_ak_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_mount_descs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"php_arms_config_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"php_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"php_config_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"post_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pre_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readiness": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_configs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"termination_grace_period_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"acr_assume_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tomcat_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"war_start_options": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_container": {
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

func dataSourceAlicloudSaeApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/pop/v1/sam/app/listApplications"
	request := make(map[string]*string)
	if v, ok := d.GetOk("app_name"); ok {
		request["AppName"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("field_type"); ok {
		request["FieldType"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("field_value"); ok {
		request["FieldValue"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("namespace_id"); ok {
		request["NamespaceId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		request["OrderBy"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOkExists("reverse"); ok {
		request["Reverse"] = StringPointer(v.(string))
	}
	request["PageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	request["CurrentPage"] = StringPointer("1")
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
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_applications", "GET "+action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
		}
		resp, err := jsonpath.Get("$.Data.Applications", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Applications", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AppId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		currentPage, err := strconv.Atoi(*request["CurrentPage"])
		if err != nil {
			return WrapError(err)
		}
		request["CurrentPage"] = StringPointer(strconv.Itoa(currentPage + 1))
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"app_description": object["AppDescription"],
			"app_name":        object["AppName"],
			"id":              fmt.Sprint(object["AppId"]),
			"application_id":  fmt.Sprint(object["AppId"]),
			"namespace_id":    object["NamespaceId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["AppId"])
		saeService := SaeService{client}
		getResp, err := saeService.DescribeSaeApplication(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["acr_assume_role_arn"] = getResp["AcrAssumeRoleArn"]
		mapping["command"] = getResp["Command"]
		mapping["command_args"] = getResp["CommandArgs"]
		config_map_mount_desc, _ := convertArrayObjectToJsonString(getResp["ConfigMapMountDesc"])
		mapping["config_map_mount_desc"] = config_map_mount_desc

		if v, ok := getResp["Cpu"]; ok && fmt.Sprint(v) != "0" {
			mapping["cpu"] = formatInt(v)
		}
		mapping["custom_host_alias"] = getResp["CustomHostAlias"]
		mapping["edas_container_version"] = getResp["EdasContainerVersion"]
		mapping["envs"] = getResp["Envs"]
		mapping["image_url"] = getResp["ImageUrl"]
		mapping["jar_start_args"] = getResp["JarStartArgs"]
		mapping["jar_start_options"] = getResp["JarStartOptions"]
		mapping["jdk"] = getResp["Jdk"]
		mapping["liveness"] = getResp["Liveness"]
		if v, ok := getResp["Memory"]; ok && fmt.Sprint(v) != "0" {
			mapping["memory"] = formatInt(v)
		}
		if v, ok := getResp["MinReadyInstances"]; ok && fmt.Sprint(v) != "0" {
			mapping["min_ready_instances"] = formatInt(v)
		}
		mapping["mount_desc"] = getResp["MountDesc"]
		mapping["mount_host"] = getResp["MountHost"]
		mapping["nas_id"] = getResp["NasId"]
		mapping["oss_ak_id"] = getResp["OssAkId"]
		mapping["oss_ak_secret"] = getResp["OssAkSecret"]
		mapping["oss_mount_descs"] = getResp["OssMountDescs"]
		mapping["package_type"] = getResp["PackageType"]
		mapping["package_url"] = getResp["PackageUrl"]
		mapping["package_version"] = getResp["PackageVersion"]
		mapping["php_arms_config_location"] = getResp["PhpArmsConfigLocation"]
		mapping["php_config"] = getResp["PhpConfig"]
		mapping["php_config_location"] = getResp["PhpConfigLocation"]
		mapping["post_start"] = getResp["PostStart"]
		mapping["pre_stop"] = getResp["PreStop"]
		mapping["readiness"] = getResp["Readiness"]
		if v, ok := getResp["Replicas"]; ok && fmt.Sprint(v) != "0" {
			mapping["replicas"] = formatInt(v)
		}
		mapping["security_group_id"] = getResp["SecurityGroupId"]
		mapping["sls_configs"] = getResp["SlsConfigs"]
		if v, ok := getResp["TerminationGracePeriodSeconds"]; ok && fmt.Sprint(v) != "0" {
			mapping["termination_grace_period_seconds"] = formatInt(v)
		}
		mapping["timezone"] = getResp["Timezone"]
		mapping["tomcat_config"] = getResp["TomcatConfig"]
		mapping["vswitch_id"] = getResp["VSwitchId"]
		mapping["vpc_id"] = getResp["VpcId"]
		mapping["war_start_options"] = getResp["WarStartOptions"]
		mapping["web_container"] = getResp["WebContainer"]
		saeService = SaeService{client}
		if imageUrl, exist := getResp["ImageUrl"]; exist {
			applicationImageResp, err := saeService.DescribeApplicationImage(id, imageUrl.(string))
			if err != nil {
				return WrapError(err)
			}
			mapping["region_id"] = applicationImageResp["RegionId"]
			mapping["repo_name"] = applicationImageResp["RepoName"]
			mapping["repo_namespace"] = applicationImageResp["RepoNamespace"]
			mapping["repo_origin_type"] = applicationImageResp["RepoOriginType"]
		}

		getRespStatus, err := saeService.DescribeApplicationStatus(id)
		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getRespStatus["CurrentStatus"].(string) {
			continue
		}
		mapping["create_time"] = getRespStatus["CreateTime"]
		mapping["status"] = getRespStatus["CurrentStatus"]

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("applications", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
