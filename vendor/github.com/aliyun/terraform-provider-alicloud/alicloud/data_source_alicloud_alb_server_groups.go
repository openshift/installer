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

func dataSourceAlicloudAlbServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbServerGroupsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
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
			"server_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"server_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Configuring", "Provisioning"}, false),
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_check_codes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"health_check_connect_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"health_check_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_http_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"sticky_session_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cookie": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sticky_session_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"sticky_session_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_id": {
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

func dataSourceAlicloudAlbServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListServerGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if m, ok := d.GetOk("server_group_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("ServerGroupIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("server_group_name"); ok {
		request["ServerGroupNames.1"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tagsFromMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var serverGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		serverGroupNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_server_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ServerGroups", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ServerGroups", response)
		}
		result, _ := resp.([]interface{})
		if resp != nil {
			for _, v := range result {
				item := v.(map[string]interface{})
				if serverGroupNameRegex != nil && !serverGroupNameRegex.MatchString(fmt.Sprint(item["ServerGroupName"])) {
					continue
				}
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(item["ServerGroupId"])]; !ok {
						continue
					}
				}
				if statusOk && status.(string) != "" && status.(string) != item["ServerGroupStatus"].(string) {
					continue
				}
				objects = append(objects, item)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"protocol":              object["Protocol"],
			"scheduler":             object["Scheduler"],
			"id":                    fmt.Sprint(object["ServerGroupId"]),
			"server_group_id":       fmt.Sprint(object["ServerGroupId"]),
			"server_group_name":     object["ServerGroupName"],
			"status":                object["ServerGroupStatus"],
			"sticky_session_config": object["StickySessionConfig"],
			"vpc_id":                object["VpcId"],
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
		mapping["health_check_config"] = healthCheckConfigSli

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
		mapping["sticky_session_config"] = stickySessionConfigSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ServerGroupName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["ServerGroupId"])
		albService := AlbService{client}
		serversList, err := albService.ListServerGroupServers(id)
		if err != nil {
			return WrapError(err)
		}

		servers := make([]map[string]interface{}, 0)
		for _, v := range serversList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description": m1["Description"],
					"port":        formatInt(m1["Port"]),
					"server_id":   m1["ServerId"],
					"server_ip":   m1["ServerIp"],
					"server_type": m1["ServerType"],
					"status":      m1["Status"],
					"weight":      formatInt(m1["Weight"]),
				}
				servers = append(servers, temp1)
			}
		}
		mapping["servers"] = servers

		tagResp, _ := albService.ListTagResources(id, "servergroup")
		mapping["tags"] = tagsToMap(tagResp)

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
