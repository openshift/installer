package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMseClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMseClustersRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"cluster_alias_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
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
			"request_pars": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DESTROY_FAILED", "DESTROY_ING", "DESTROY_SUCCESS", "INIT_FAILED", "INIT_ING", "INIT_SUCCESS", "INIT_TIME_OUT", "RESTART_FAILED", "RESTART_ING", "RESTART_SUCCESS", "SCALE_FAILED", "SCALE_ING", "SCALE_SUCCESS"}, false),
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
						"acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"init_cost_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_models": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"internet_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pod_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"single_tunnel_vip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pay_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pub_network_flow": {
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

func dataSourceAlicloudMseClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListClusters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_alias_name"); ok {
		request["ClusterAliasName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("request_pars"); ok {
		request["RequestPars"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 1
	var objects []map[string]interface{}
	var clusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		clusterNameRegex = r
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
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-05-31"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_clusters", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if clusterNameRegex != nil {
				if !clusterNameRegex.MatchString(fmt.Sprint(item["ClusterAliasName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["InitStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"app_version":      object["AppVersion"],
			"cluster_id":       object["ClusterId"],
			"cluster_name":     object["ClusterAliasName"],
			"cluster_type":     object["ClusterType"],
			"id":               fmt.Sprint(object["InstanceId"]),
			"instance_id":      fmt.Sprint(object["InstanceId"]),
			"internet_address": object["InternetAddress"],
			"internet_domain":  object["InternetDomain"],
			"intranet_address": object["IntranetAddress"],
			"intranet_domain":  object["IntranetDomain"],
			"status":           object["InitStatus"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["InstanceId"]))
			names = append(names, object["ClusterAliasName"])
			s = append(s, mapping)
			continue
		}

		mseService := MseService{client}
		id := fmt.Sprint(object["InstanceId"])
		getResp, err := mseService.DescribeMseCluster(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["acl_id"] = getResp["AclId"]
		mapping["cpu"] = getResp["Cpu"]
		mapping["health_status"] = getResp["HealthStatus"]
		mapping["init_cost_time"] = getResp["InitCostTime"]
		mapping["instance_count"] = getResp["InstanceCount"]

		instanceModels := make([]map[string]interface{}, 0)
		if instanceModelsList, ok := getResp["InstanceModels"].([]interface{}); ok {
			for _, v := range instanceModelsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"health_status":     m1["HealthStatus"],
						"instance_type":     m1["InstanceType"],
						"internet_ip":       m1["InternetIp"],
						"ip":                m1["Ip"],
						"pod_name":          m1["PodName"],
						"role":              m1["Role"],
						"single_tunnel_vip": m1["SingleTunnelVip"],
						"vip":               m1["Vip"],
					}
					instanceModels = append(instanceModels, temp1)
				}
			}
		}
		mapping["instance_models"] = instanceModels
		mapping["internet_port"] = getResp["InternetPort"]
		mapping["intranet_port"] = getResp["IntranetPort"]
		mapping["memory_capacity"] = getResp["MemoryCapacity"]
		mapping["pay_info"] = getResp["PayInfo"]
		mapping["pub_network_flow"] = getResp["PubNetworkFlow"]
		ids = append(ids, fmt.Sprint(object["InstanceId"]))
		names = append(names, object["ClusterAliasName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
