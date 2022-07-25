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

func dataSourceAlicloudGraphDatabaseDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGraphDatabaseDbInstancesRead,
		Schema: map[string]*schema.Schema{
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Running", "Deleting", "DBInstanceClassChanging", "NetAddressCreating", "NetAddressDeleting"}, false),
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
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_minor_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_cpu": {
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
						"db_instance_memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_minor_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"public_connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"read_only_db_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_ip_array": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_instance_ip_array_attribute": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"db_instance_ip_array_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"security_ips": {
										Type:     schema.TypeString,
										Optional: true,
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

func dataSourceAlicloudGraphDatabaseDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["DBInstanceStatus"] = v
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
	conn, err := client.NewGdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_graph_database_db_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"])]; !ok {
					continue
				}
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
			"create_time":             object["CreateTime"],
			"db_instance_category":    object["Category"],
			"db_instance_description": object["DBInstanceDescription"],
			"id":                      fmt.Sprint(object["DBInstanceId"]),
			"db_instance_id":          fmt.Sprint(object["DBInstanceId"]),
			"db_instance_type":        object["DBInstanceType"],
			"db_node_class":           object["DBNodeClass"],
			"db_node_count":           fmt.Sprint(object["DBNodeCount"]),
			"db_node_storage":         fmt.Sprint(object["DBNodeStorage"]),
			"db_version":              object["DBVersion"],
			"expire_time":             object["ExpireTime"],
			"expired":                 fmt.Sprint(object["Expired"]),
			"lock_mode":               object["LockMode"],
			"lock_reason":             object["LockReason"],
			"master_db_instance_id":   object["MasterInstanceId"],
			"payment_type":            convertGraphDatabaseDbInstancePaymentTypeResponse(object["PayType"]),
			"status":                  object["DBInstanceStatus"],
			"vswitch_id":              object["VSwitchId"],
			"vpc_id":                  object["VpcId"],
			"zone_id":                 object["ZoneId"],
		}
		if v := object["ReadOnlyDBInstanceIds"].(map[string]interface{})["ReadOnlyDBInstanceId"].([]interface{}); len(v) > 0 {
			mapping["read_only_db_instance_ids"] = v[0].(map[string]interface{})["DBInstanceId"]
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DBInstanceId"])
		gdbService := GdbService{client}
		getResp, err := gdbService.DescribeGraphDatabaseDbInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["connection_string"] = getResp["ConnectionString"]
		mapping["current_minor_version"] = getResp["CurrentMinorVersion"]
		mapping["db_instance_cpu"] = getResp["DBInstanceCPU"]
		mapping["db_instance_memory"] = getResp["DBInstanceMemory"]
		mapping["db_instance_network_type"] = getResp["DBInstanceNetworkType"]
		mapping["db_instance_storage_type"] = getResp["DBInstanceStorageType"]
		mapping["latest_minor_version"] = getResp["LatestMinorVersion"]
		mapping["maintain_time"] = getResp["MaintainTime"]
		if DBInstanceIPArray, ok := getResp["DBInstanceIPArray"]; ok {
			DBInstanceIPArrayAry, ok := DBInstanceIPArray.([]interface{})
			if ok && len(DBInstanceIPArrayAry) > 0 {
				DBInstanceIPArraySli := make([]map[string]interface{}, 0)
				for _, DBInstanceIPArrayArg := range DBInstanceIPArrayAry {
					DBInstanceIPArrayMap := make(map[string]interface{})
					DBInstanceIPArrayMap["security_ips"] = DBInstanceIPArrayArg.(map[string]interface{})["SecurityIps"]
					DBInstanceIPArrayMap["db_instance_ip_array_name"] = DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayName"]
					if v, ok := DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayAttribute"]; ok {
						DBInstanceIPArrayMap["db_instance_ip_array_attribute"] = v
					}
					DBInstanceIPArraySli = append(DBInstanceIPArraySli, DBInstanceIPArrayMap)
				}
				mapping["db_instance_ip_array"] = DBInstanceIPArraySli
			}
		}
		if v, ok := getResp["Port"]; ok && v.(string) != "" {
			mapping["port"] = formatInt(v)
		}
		mapping["public_connection_string"] = getResp["PublicConnectionString"]
		if v, ok := getResp["PublicPort"]; ok && fmt.Sprint(v) != "0" {
			mapping["public_port"] = formatInt(v)
		}
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
