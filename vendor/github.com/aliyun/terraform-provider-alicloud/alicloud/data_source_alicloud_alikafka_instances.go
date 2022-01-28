package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deploy_type": {
							Type:     schema.TypeInt,
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
						"io_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eip_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"topic_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"paid_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"msg_retain": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ssl_end_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgrade_service_detail_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current2_open_source_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"allowed_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deploy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allowed_ip_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"port_range": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"internet_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allowed_ip_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"port_range": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"domain_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_domain_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sasl_domain_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetInstanceList"
	request := make(map[string]interface{})
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId

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
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var response map[string]interface{}
	pageNo, pageSize := 1, PageSizeLarge
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceList.InstanceVO", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList.InstanceVO", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < pageSize {
			break
		}
		pageNo++
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)

	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		paidType := PostPaid
		if object["PaidType"] == 0 {
			paidType = PrePaid
		}

		mapping := map[string]interface{}{
			"id":                   object["InstanceId"],
			"name":                 object["Name"],
			"create_time":          object["CreateTime"],
			"service_status":       object["ServiceStatus"],
			"deploy_type":          object["DeployType"],
			"vpc_id":               object["VpcId"],
			"vswitch_id":           object["VSwitchId"],
			"io_max":               object["IoMax"],
			"eip_max":              object["EipMax"],
			"disk_type":            object["DiskType"],
			"disk_size":            object["DiskSize"],
			"topic_quota":          object["TopicNumLimit"],
			"paid_type":            paidType,
			"service_version":      object["UpgradeServiceDetailInfo"].(map[string]interface{})["Current2OpenSourceVersion"],
			"spec_type":            object["SpecType"],
			"zone_id":              object["ZoneId"],
			"end_point":            object["EndPoint"],
			"security_group":       object["SecurityGroup"],
			"config":               object["AllConfig"],
			"expired_time":         object["ExpiredTime"],
			"msg_retain":           object["MsgRetain"],
			"ssl_end_point":        object["SslEndPoint"],
			"domain_endpoint":      object["DomainEndpoint"],
			"ssl_domain_endpoint":  object["SslDomainEndpoint"],
			"sasl_domain_endpoint": object["SaslDomainEndpoint"],
		}
		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.TagVO", object)
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

		DetailInfoMaps := make([]map[string]interface{}, 0)
		if _, ok := object["UpgradeServiceDetailInfo"].(map[string]interface{}); ok {
			UpgradeServiceDetailInfoMap := map[string]interface{}{}
			UpgradeServiceDetailInfoMap["current2_open_source_version"] = object["UpgradeServiceDetailInfo"].(map[string]interface{})["Current2OpenSourceVersion"]
			DetailInfoMaps = append(DetailInfoMaps, UpgradeServiceDetailInfoMap)
		}
		mapping["upgrade_service_detail_info"] = DetailInfoMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, mapping["name"])
		id := fmt.Sprint(object["InstanceId"])

		AlikaService := AlikafkaService{client}
		if d.Get("enable_details").(bool) {
			getResp, err := AlikaService.GetAllowedIpList(id)
			if err != nil {
				return WrapError(err)
			}

			allowedListMaps := make([]map[string]interface{}, 0)
			if defaultActionsList, ok := getResp["AllowedList"].(map[string]interface{}); ok {
				defaultActionsMap := map[string]interface{}{}
				defaultActionsMap["deploy_type"] = defaultActionsList["DeployType"]
				if forwardGroupConfigArg, ok := defaultActionsList["VpcList"].([]interface{}); ok {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					for _, serverGroupTuples := range forwardGroupConfigArg {
						serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
						serverGroupTuplesMap := map[string]interface{}{}
						serverGroupTuplesMap["port_range"] = serverGroupTuplesArg["PortRange"]
						serverGroupTuplesMap["allowed_ip_list"] = serverGroupTuplesArg["AllowedIpList"]
						serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
					}
					defaultActionsMap["vpc_list"] = serverGroupTuplesMaps
				}

				if forwardGroupConfigArg, ok := defaultActionsList["InternetList"].([]interface{}); ok {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					for _, serverGroupTuples := range forwardGroupConfigArg {
						serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
						serverGroupTuplesMap := map[string]interface{}{}
						serverGroupTuplesMap["port_range"] = serverGroupTuplesArg["PortRange"]
						serverGroupTuplesMap["allowed_ip_list"] = serverGroupTuplesArg["AllowedIpList"]
						serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
					}
					defaultActionsMap["internet_list"] = serverGroupTuplesMaps
				}

				allowedListMaps = append(allowedListMaps, defaultActionsMap)

			}
			mapping["allowed_list"] = allowedListMaps
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
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
