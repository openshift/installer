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

func dataSourceAlicloudMongodbServerlessInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongodbServerlessInstancesRead,
		Schema: map[string]*schema.Schema{
			"db_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_instance_description": {
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
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "VPC"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "DBInstanceClassChanging", "DBInstanceNetTypeChanging", "Deleting", "EngineVersionUpgrading", "GuardSwitching", "HASwitching", "Importing", "ImportingFromOthers", "LinkSwitching", "MinorVersionUpgrading", "NodeCreating", "NodeDeleting", "Rebooting", "Restoring", "Running", "SSLModifying", "TDEModifying", "TempDBInstanceCreating", "Transing", "TransingToOthers", "released"}, false),
			},
			"tags": tagsSchema(),
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
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_unit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_release_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"db_instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_ip_group_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_ip_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_ip_list": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_auth_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
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

func dataSourceAlicloudMongodbServerlessInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_class"); ok {
		request["DBInstanceClass"] = v
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["DBInstanceType"] = "serverless"
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	status, statusOk := d.GetOk("status")
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
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
	var response map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mongodb_serverless_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBInstances.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["DBInstanceStatus"].(string) {
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
			"payment_type":            convertMongodbServerlessInstancePayTypeResponse(object["ChargeType"].(string)),
			"db_instance_class":       object["DBInstanceClass"],
			"db_instance_description": object["DBInstanceDescription"],
			"id":                      fmt.Sprint(object["DBInstanceId"]),
			"db_instance_id":          fmt.Sprint(object["DBInstanceId"]),
			"db_instance_storage":     formatInt(object["DBInstanceStorage"]),
			"engine":                  object["Engine"],
			"engine_version":          object["EngineVersion"],
			"expire_time":             object["ExpireTime"],
			"kind_code":               object["KindCode"],
			"lock_mode":               object["LockMode"],
			"network_type":            object["NetworkType"],
			"resource_group_id":       object["ResourceGroupId"],
			"status":                  object["DBInstanceStatus"],
			"vpc_auth_mode":           object["VpcAuthMode"],
			"zone_id":                 object["ZoneId"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DBInstanceId"])
		MongoDBService := MongoDBService{client}
		getResp, err := MongoDBService.DescribeMongodbServerlessInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["capacity_unit"] = formatInt(getResp["CapacityUnit"])
		mapping["db_instance_release_protection"] = getResp["DBInstanceReleaseProtection"]
		mapping["maintain_end_time"] = getResp["MaintainEndTime"]
		mapping["maintain_start_time"] = getResp["MaintainStartTime"]
		if v, ok := getResp["MaxConnections"]; ok && fmt.Sprint(v) != "0" {
			mapping["max_connections"] = formatInt(v)
		}
		if v, ok := getResp["MaxIOPS"]; ok && fmt.Sprint(v) != "0" {
			mapping["max_iops"] = formatInt(v)
		}
		mapping["protocol_type"] = getResp["ProtocolType"]
		mapping["storage_engine"] = convertMongodbServerlessInstanceStorageEngineResponse(getResp["StorageEngine"].(string))
		mapping["vswitch_id"] = getResp["VSwitchId"]
		mapping["vpc_id"] = getResp["VPCId"]
		getResp1, err := MongoDBService.DescribeSecurityIps(id)
		if err != nil {
			return WrapError(err)
		}
		iPArray := make([]map[string]interface{}, 0)
		if securityIpGroup, ok := getResp1["SecurityIpGroups"].(map[string]interface{})["SecurityIpGroup"].([]interface{}); ok {
			for _, v := range securityIpGroup {
				if m1, ok := v.(map[string]interface{}); ok {
					if m1["SecurityIpGroupName"].(string) == "default" {
						continue
					}
					temp1 := map[string]interface{}{
						"security_ip_group_attribute": m1["SecurityIpGroupAttribute"],
						"security_ip_group_name":      m1["SecurityIpGroupName"],
						"security_ip_list":            m1["SecurityIpList"],
					}
					iPArray = append(iPArray, temp1)
				}
			}
		}
		mapping["security_ip_groups"] = iPArray

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
