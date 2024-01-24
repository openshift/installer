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

func dataSourceAlicloudHbrEcsBackupClients() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrEcsBackupClientsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVATED", "DEACTIVATED", "INSTALLING", "INSTALL_FAILED", "NOT_INSTALLED", "REGISTERED", "STOPPED", "UNINSTALLING", "UNINSTALL_FAILED", "UNKNOWN", "UPGRADE_FAILED", "UPGRADING"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clients": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arch_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_proxy_setting": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_backup_client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_heart_beat_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_client_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_cpu_core": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_worker": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ipv4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrEcsBackupClientsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackupClients"
	request := make(map[string]interface{})
	request["ClientType"] = "ECS_CLIENT"
	if v, ok := d.GetOk("instance_ids"); ok {
		request["InstanceIds"] = convertListToJsonString(v.([]interface{}))
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
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
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_ecs_backup_clients", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Clients", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Clients", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ClientId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"arch_type":            object["ArchType"],
			"backup_status":        object["BackupStatus"],
			"client_type":          object["ClientType"],
			"client_version":       object["ClientVersion"],
			"create_time":          fmt.Sprint(object["CreatedTime"]),
			"data_network_type":    object["Settings"].(map[string]interface{})["DataNetworkType"],
			"data_proxy_setting":   object["Settings"].(map[string]interface{})["DataProxySetting"],
			"id":                   fmt.Sprint(object["ClientId"]),
			"ecs_backup_client_id": fmt.Sprint(object["ClientId"]),
			"hostname":             object["Hostname"],
			"instance_id":          object["InstanceId"],
			"instance_name":        object["InstanceName"],
			"last_heart_beat_time": fmt.Sprint(object["LastHeartBeatTime"]),
			"max_client_version":   object["MaxClientVersion"],
			"max_cpu_core":         object["Settings"].(map[string]interface{})["MaxCpuCore"],
			"max_worker":           object["Settings"].(map[string]interface{})["MaxWorker"],
			"os_type":              object["OsType"],
			"private_ipv4":         object["PrivateIpV4"],
			"proxy_host":           object["Settings"].(map[string]interface{})["ProxyHost"],
			"proxy_password":       object["Settings"].(map[string]interface{})["ProxyPassword"],
			"proxy_port":           fmt.Sprint(object["Settings"].(map[string]interface{})["ProxyPort"]),
			"proxy_user":           object["Settings"].(map[string]interface{})["ProxyUser"],
			"status":               object["Status"],
			"updated_time":         fmt.Sprint(object["UpdatedTime"]),
			"use_https":            object["Settings"].(map[string]interface{})["UseHttps"],
			"zone_id":              object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clients", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
