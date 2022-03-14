package alicloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAdbDbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAdbDbClustersRead,
		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"descriptions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{"ClassChanging", "Creating", "Deleting", "NetAddressCreating", "NetAddressDeleting", "Preparing", "Restoring", "Running"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"commodity_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compute_resource": {
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
						"db_cluster_category": {
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
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_node_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"elastic_io_resource": {
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
						"executor_count": {
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
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rds_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
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
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
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
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudAdbDbClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["DBClusterDescription"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["DBClusterStatus"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag.*"] = tags
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}
	var objects []interface{}
	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_db_clusters", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBCluster", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if descriptionRegex != nil {
				if !descriptionRegex.MatchString(fmt.Sprint(item["DBClusterDescription"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBClusterId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	descriptions := make([]string, 0)
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"commodity_code":          object["CommodityCode"],
			"compute_resource":        object["ComputeResource"],
			"connection_string":       object["ConnectionString"],
			"create_time":             object["CreateTime"],
			"db_cluster_category":     convertAdbDBClusterCategoryResponse(object["Category"].(string)),
			"id":                      fmt.Sprint(object["DBClusterId"]),
			"db_cluster_id":           fmt.Sprint(object["DBClusterId"]),
			"db_cluster_network_type": object["DBClusterNetworkType"],
			"network_type":            object["DBClusterNetworkType"],
			"db_cluster_type":         object["DBClusterType"],
			"db_node_class":           object["DBNodeClass"],
			"db_node_count":           object["DBNodeCount"],
			"db_node_storage":         object["DBNodeStorage"],
			"description":             object["DBClusterDescription"],
			"disk_type":               object["DiskType"],
			"dts_job_id":              object["DtsJobId"],
			"elastic_io_resource":     formatInt(object["ElasticIOResource"]),
			"engine":                  object["Engine"],
			"executor_count":          object["ExecutorCount"],
			"expire_time":             object["ExpireTime"],
			"expired":                 fmt.Sprint(object["Expired"]),
			"lock_mode":               object["LockMode"],
			"lock_reason":             object["LockReason"],
			"payment_type":            convertAdbDBClusterPaymentTypeResponse(object["PayType"].(string)),
			"charge_type":             getChargeType(object["PayType"].(string)),
			"port":                    formatInt(object["Port"]),
			"rds_instance_id":         object["RdsInstanceId"],
			"region_id":               object["RegionId"],
			"resource_group_id":       object["ResourceGroupId"],
			"status":                  object["DBClusterStatus"],
			"storage_resource":        object["StorageResource"],
			"vpc_cloud_instance_id":   object["VPCCloudInstanceId"],
			"vpc_id":                  object["VPCId"],
			"vswitch_id":              object["VSwitchId"],
			"zone_id":                 object["ZoneId"],
			"mode":                    object["Mode"],
			"db_cluster_version":      object["DBVersion"],
		}
		descriptions = append(descriptions, object["DBClusterDescription"].(string))

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
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["DBClusterId"]))
			s = append(s, mapping)
			continue
		}

		adbService := AdbService{client}
		id := fmt.Sprint(object["DBClusterId"])
		if object["PayType"].(string) == string(Prepaid) {
			getResp, err := adbService.DescribeAutoRenewAttribute(id)
			if err != nil {
				return WrapError(err)
			}
			mapping["auto_renew_period"] = getResp["Duration"]
			mapping["renewal_status"] = getResp["RenewalStatus"]
		}
		getResp1, err := adbService.DescribeDBClusterAccessWhiteList(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["security_ips"] = strings.Split(getResp1["SecurityIPList"].(string), ",")

		getResp2, err := adbService.DescribeAdbDbCluster(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["engine_version"] = getResp2["EngineVersion"]
		mapping["maintain_time"] = getResp2["MaintainTime"]

		ids = append(ids, fmt.Sprint(object["DBClusterId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", formatInt(response["TotalCount"])); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
