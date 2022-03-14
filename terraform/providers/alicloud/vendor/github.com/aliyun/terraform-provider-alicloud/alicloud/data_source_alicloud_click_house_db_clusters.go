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

func dataSourceAlicloudClickHouseDbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudClickHouseDbClustersRead,
		Schema: map[string]*schema.Schema{
			"db_cluster_description": {
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
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "Restarting", "Preparing", "Running"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ali_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commodity_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_access_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_cluster_ip_array_attribute": {
										Type:     schema.TypeString,
										Computed: true,
										Removed:  "Field 'db_cluster_ip_array_attribute' has been removed from provider",
									},
									"db_cluster_ip_array_name": {
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
						"db_cluster_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_type": {
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
						"encryption_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption_type": {
							Type:     schema.TypeString,
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
						"is_expired": {
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"scale_out_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"progress": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ratio": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_backup": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"support_https_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"support_mysql_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
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
						"control_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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

func dataSourceAlicloudClickHouseDbClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["DBClusterStatus"] = v
	}
	request["RegionId"] = client.RegionId
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
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_click_house_db_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DBClusters.DBCluster", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBClusters.DBCluster", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBClusterId"])]; !ok {
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
			"ali_uid":                 object["AliUid"],
			"bid":                     object["Bid"],
			"category":                object["Category"],
			"commodity_code":          object["CommodityCode"],
			"connection_string":       object["ConnectionString"],
			"create_time":             object["CreateTime"],
			"db_cluster_description":  object["DBClusterDescription"],
			"id":                      fmt.Sprint(object["DBClusterId"]),
			"db_cluster_id":           fmt.Sprint(object["DBClusterId"]),
			"db_cluster_network_type": object["DBClusterNetworkType"],
			"db_node_class":           object["DBNodeClass"],
			"db_node_count":           fmt.Sprint(object["DBNodeCount"]),
			"db_node_storage":         fmt.Sprint(object["DBNodeStorage"]),
			"expire_time":             object["ExpireTime"],
			"lock_mode":               object["LockMode"],
			"lock_reason":             object["LockReason"],
			"payment_type":            convertClickHouseDbClusterPaymentTypeResponse(object["PayType"].(string)),
			"port":                    formatInt(object["Port"]),
			"storage_type":            object["StorageType"],
			"vswitch_id":              object["VSwitchId"],
			"vpc_cloud_instance_id":   object["VpcCloudInstanceId"],
			"vpc_id":                  object["VpcId"],
			"zone_id":                 object["ZoneId"],
			"control_version":         object["ControlVersion"],
			"status":                  object["DBClusterStatus"],
		}
		if val, exist := object["IsExpired"]; exist {
			object["is_expired"] = val.(bool)
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DBClusterId"])
		clickhouseService := ClickhouseService{client}
		getResp1, err := clickhouseService.DescribeDBClusterAccessWhiteList(id)
		if err != nil {
			return WrapError(err)
		}

		iPArray := make([]map[string]interface{}, 0)
		if iPArrayList, ok := getResp1["DBClusterAccessWhiteList"].(map[string]interface{})["IPArray"].([]interface{}); ok {
			for _, v := range iPArrayList {
				if m1, ok := v.(map[string]interface{}); ok {
					if m1["DBClusterIPArrayName"].(string) == "default" || m1["DBClusterIPArrayName"].(string) == "dms" {
						continue
					}
					temp1 := map[string]interface{}{
						"db_cluster_ip_array_attribute": m1["DBClusterIPArrayAttribute"],
						"db_cluster_ip_array_name":      m1["DBClusterIPArrayName"],
						"security_ip_list":              m1["SecurityIPList"],
					}
					iPArray = append(iPArray, temp1)
				}
			}
		}
		mapping["db_cluster_access_white_list"] = iPArray
		getResp, err := clickhouseService.DescribeClickHouseDbCluster(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["db_cluster_type"] = getResp["DBClusterType"]
		mapping["encryption_key"] = getResp["EncryptionKey"]
		mapping["encryption_type"] = getResp["EncryptionType"]
		mapping["engine"] = getResp["Engine"]
		mapping["engine_version"] = getResp["EngineVersion"]
		mapping["maintain_time"] = getResp["MaintainTime"]
		mapping["public_connection_string"] = getResp["PublicConnectionString"]
		mapping["public_port"] = getResp["PublicPort"]
		if v, ok := getResp["SupportBackup"]; ok && fmt.Sprint(v) != "0" {
			mapping["support_backup"] = formatInt(v)
		}
		mapping["support_https_port"] = getResp["SupportHttpsPort"]
		mapping["support_mysql_port"] = getResp["SupportMysqlPort"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
