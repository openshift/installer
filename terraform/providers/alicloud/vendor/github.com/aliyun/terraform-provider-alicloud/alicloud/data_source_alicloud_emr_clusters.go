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

func dataSourceAlicloudEmrClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_type_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MANUAL", "ON-DEMAND"}, false),
			},
			"default_status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"deposit_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HALF_MANAGED", "MANAGED"}, false),
			},
			"is_desc": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"machine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ECS",
				ValidateFunc: validation.StringInSlice([]string{"DOCKER", "ECS", "ECS_FROM_ECM_HOSTPOOL", "PYHSICAL_MACHINE"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ABNORMAL", "CREATE_FAILED", "CREATING", "IDLE", "RELEASED", "RELEASE_FAILED", "RELEASING", "RUNNING", "WAIT_FOR_PAY"}, false),
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
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
				Default:  50,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zk_links": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"link": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"auto_scaling_allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling_by_load_allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_scaling_spot_with_limit_allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bootstrap_action_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"arg": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"bootstrap_failed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
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
						"create_resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deposit_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eas_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_uncompleted_order": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extra_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"high_availability_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"host_group_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"band_width": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"charge_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_core": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disk_capacity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disk_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_group_change_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_group_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory_capacity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"node_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"nodes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"create_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_infos": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"device": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"disk_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"disk_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"size": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"type": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"emr_expired_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"expired_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"inner_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"instance_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pub_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"support_ipv6": {
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
									"period": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"host_pool_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hp_biz_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hp_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_meta_db": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta_store_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relate_cluster_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resize_disk_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"running_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"software_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"emr_ver": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"softwares": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"display_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"only_display": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"start_tpe": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stop_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"user_defined_emr_ecs_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
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

func dataSourceAlicloudEmrClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListClusters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("cluster_type_list"); ok {
		request["ClusterTypeList"] = v
	}

	if v, ok := d.GetOk("create_type"); ok {
		request["CreateType"] = v
	}
	if v, ok := d.GetOkExists("default_status"); ok {
		request["DefaultStatus"] = v
	}
	if v, ok := d.GetOk("deposit_type"); ok {
		request["DepositType"] = v
	}

	if v, ok := d.GetOkExists("is_desc"); ok {
		request["IsDesc"] = v
	}
	if v, ok := d.GetOk("machine_type"); ok {
		request["MachineType"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status_list"); ok {
		request["StatusList"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
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

	var response map[string]interface{}
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Clusters.ClusterInfo", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Clusters.ClusterInfo", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if clusterNameRegex != nil && !clusterNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
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
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["Id"]),
			"cluster_id":            fmt.Sprint(object["Id"]),
			"cluster_name":          object["Name"],
			"create_resource":       object["CreateResource"],
			"create_time":           fmt.Sprint(object["CreateTime"]),
			"deposit_type":          object["DepositType"],
			"expired_time":          fmt.Sprint(object["ExpiredTime"]),
			"machine_type":          object["MachineType"],
			"meta_store_type":       object["MetaStoreType"],
			"payment_type":          object["ChargeType"],
			"period":                formatInt(object["Period"]),
			"running_time":          formatInt(object["RunningTime"]),
			"status":                object["Status"],
			"type":                  object["Type"],
			"has_uncompleted_order": object["HasUncompletedOrder"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, object["Name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["Id"])
		emrService := EmrService{client}
		getResp, err := emrService.DescribeClusterBasicInfo(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["auto_scaling_allowed"] = getResp["AutoScalingAllowed"]
		mapping["auto_scaling_by_load_allowed"] = getResp["AutoScalingByLoadAllowed"]
		mapping["auto_scaling_spot_with_limit_allowed"] = getResp["AutoScalingSpotWithLimitAllowed"]

		bootstrapAction := make([]map[string]interface{}, 0)
		if bootstrapActionList, ok := getResp["BootstrapActionList"].(map[string]interface{})["BootstrapAction"].([]interface{}); ok {
			for _, v := range bootstrapActionList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"arg":  m1["Arg"],
						"name": m1["Name"],
						"path": m1["Path"],
					}
					bootstrapAction = append(bootstrapAction, temp1)
				}
			}
		}
		mapping["bootstrap_action_list"] = bootstrapAction

		softwareInfo := make([]map[string]interface{}, 0)
		if softwareInfoList, ok := getResp["SoftwareInfo"].(map[string]interface{}); ok {

			temp1 := map[string]interface{}{
				"emr_ver":      softwareInfoList["EmrVer"],
				"cluster_type": softwareInfoList["ClusterType"],
			}

			if softwareInfoList["Softwares"] != nil {
				softwareMaps := make([]map[string]interface{}, 0)
				for _, softwareValue := range softwareInfoList["Softwares"].(map[string]interface{})["Software"].([]interface{}) {
					software := softwareValue.(map[string]interface{})
					softwareMap := map[string]interface{}{
						"start_tpe":    software["StartTpe"],
						"version":      software["Version"],
						"display_name": software["DisplayName"],
						"only_display": software["OnlyDisplay"],
						"name":         software["Name"],
					}
					softwareMaps = append(softwareMaps, softwareMap)
				}
				temp1["softwares"] = softwareMaps
			}
			softwareInfo = append(softwareInfo, temp1)
		}
		mapping["software_info"] = softwareInfo

		hostPoolInfo := make([]map[string]interface{}, 0)
		if hostPoolInfoList, ok := getResp["HostPoolInfo"].(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"hp_biz_id": hostPoolInfoList["HpBizId"],
				"hp_name":   hostPoolInfoList["HpName"],
			}
			hostPoolInfo = append(hostPoolInfo, temp1)
		}
		mapping["host_pool_info"] = hostPoolInfo

		relateClusterInfo := make([]map[string]interface{}, 0)
		if relateClusterInfoList, ok := getResp["RelateClusterInfo"].(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"cluster_id":   relateClusterInfoList["ClusterId"],
				"cluster_name": relateClusterInfoList["ClusterName"],
				"cluster_type": relateClusterInfoList["ClusterType"],
				"status":       relateClusterInfoList["Status"],
			}
			relateClusterInfo = append(relateClusterInfo, temp1)
		}
		mapping["relate_cluster_info"] = relateClusterInfo

		mapping["bootstrap_failed"] = getResp["BootstrapFailed"]

		mapping["create_type"] = getResp["CreateType"]
		mapping["eas_enable"] = getResp["EasEnable"]
		mapping["extra_info"] = getResp["ExtraInfo"]
		mapping["image_id"] = getResp["ImageId"]
		mapping["local_meta_db"] = getResp["LocalMetaDb"]
		mapping["net_type"] = getResp["NetType"]
		mapping["security_group_id"] = getResp["SecurityGroupId"]
		mapping["security_group_name"] = getResp["SecurityGroupName"]
		mapping["start_time"] = getResp["StartTime"]
		mapping["stop_time"] = getResp["StopTime"]
		mapping["user_defined_emr_ecs_role"] = getResp["UserDefinedEmrEcsRole"]
		mapping["user_id"] = getResp["UserId"]
		mapping["vswitch_id"] = getResp["VSwitchId"]
		mapping["vpc_id"] = getResp["VpcId"]
		mapping["zone_id"] = getResp["ZoneId"]
		emrService = EmrService{client}
		getResp1, err := emrService.DataSourceDescribeEmrCluster(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["access_info"] = getResp1["AccessInfo"]
		accessInfos := make([]map[string]interface{}, 0)
		if accessInfosList, ok := getResp1["AccessInfo"].([]interface{}); ok {
			for _, v := range accessInfosList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{}
					if m1["ZKLinks"] != nil {
						accessInfosMaps := make([]map[string]interface{}, 0)
						for _, accessInfosValue := range m1["ZKLinks"].([]interface{}) {
							accessInfo := accessInfosValue.(map[string]interface{})
							accessInfosMap := map[string]interface{}{
								"link": accessInfo["Link"],
								"port": accessInfo["Port"],
							}
							accessInfosMaps = append(accessInfosMaps, accessInfosMap)
						}
						temp1["zk_links"] = accessInfosMaps
					}
					accessInfos = append(accessInfos, temp1)
				}
			}
		}
		mapping["access_info"] = accessInfos

		mapping["auto_scaling_enable"] = getResp1["AutoScalingEnable"]
		mapping["high_availability_enable"] = getResp1["HighAvailabilityEnable"]

		hostGroup := make([]map[string]interface{}, 0)
		if hostGroupList, ok := getResp1["HostGroupList"].(map[string]interface{})["HostGroup"].([]interface{}); ok {
			for _, v := range hostGroupList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"band_width":             m1["BandWidth"],
						"charge_type":            m1["ChargeType"],
						"cpu_core":               formatInt(m1["CpuCore"]),
						"disk_capacity":          formatInt(m1["DiskCapacity"]),
						"disk_count":             formatInt(m1["DiskCount"]),
						"disk_type":              m1["DiskType"],
						"host_group_change_type": m1["HostGroupChangeType"],
						"host_group_id":          m1["HostGroupId"],
						"host_group_name":        m1["HostGroupName"],
						"host_group_type":        m1["HostGroupType"],
						"instance_type":          m1["InstanceType"],
						"memory_capacity":        formatInt(m1["MemoryCapacity"]),
						"node_count":             formatInt(m1["NodeCount"]),
						"period":                 m1["Period"],
					}
					if m1["Nodes"] != nil {
						nodesMaps := make([]map[string]interface{}, 0)
						for _, nodesValue := range m1["Nodes"].(map[string]interface{})["Node"].([]interface{}) {
							nodes := nodesValue.(map[string]interface{})
							nodesMap := map[string]interface{}{
								"create_time":      nodes["CreateTime"],
								"status":           nodes["Status"],
								"emr_expired_time": nodes["EmrExpiredTime"],
								"support_ipv6":     nodes["SupportIpV6"],
								"zone_id":          nodes["ZoneId"],
								"inner_ip":         nodes["InnerIp"],
								"instance_id":      nodes["InstanceId"],
								"expired_time":     nodes["ExpiredTime"],
								"pub_ip":           nodes["PubIp"],
							}

							if nodes["DiskInfos"] != nil {
								deskInfoMaps := make([]map[string]interface{}, 0)
								for _, diskInfosValue := range nodes["DiskInfos"].(map[string]interface{})["DiskInfo"].([]interface{}) {
									disks := diskInfosValue.(map[string]interface{})
									disksMap := map[string]interface{}{
										"device":    disks["Device"],
										"disk_id":   disks["DiskId"],
										"disk_name": disks["DiskName"],
										"size":      disks["Size"],
										"type":      disks["Type"],
									}
									deskInfoMaps = append(deskInfoMaps, disksMap)
								}
								nodesMap["disk_infos"] = deskInfoMaps
							}
							nodesMaps = append(nodesMaps, nodesMap)
						}
						temp1["nodes"] = nodesMaps
					}
					hostGroup = append(hostGroup, temp1)
				}
			}
		}
		mapping["host_group_list"] = hostGroup
		mapping["resize_disk_enable"] = getResp1["ResizeDiskEnable"]

		emrService = EmrService{client}
		getResp2, err := emrService.ListTagResources(id, "cluster")
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(getResp2)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
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
	if err := d.Set("total_count", formatInt(response["TotalCount"])); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
