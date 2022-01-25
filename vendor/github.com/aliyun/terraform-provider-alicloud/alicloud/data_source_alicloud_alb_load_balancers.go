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

func dataSourceAlicloudAlbLoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbLoadBalancersRead,
		Schema: map[string]*schema.Schema{
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Intranet", "Internet"}, false),
			},
			"load_balancer_bussiness_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Abnormal", "Normal"}, false),
				Deprecated:   "Field 'load_balancer_bussiness_status' has been deprecated from provider version 1.142.0 and it will be removed in the future version. Please use the new attribute 'load_balancer_business_status' instead.",
			},
			"load_balancer_business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Abnormal", "Normal"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"load_balancer_ids": {
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
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Configuring", "CreateFailed", "Inactive", "Provisioning"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"balancers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_project": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_store": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"address_allocated_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_protection_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enabled_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_billing_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pay_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"load_balancer_bussiness_status": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Field 'load_balancer_bussiness_status' has been deprecated from provider version 1.142.0 and it will be removed in the future version. Please use the new parameter 'load_balancer_business_status' instead.",
						},
						"load_balancer_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_edition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_operation_locks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lock_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lock_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"modification_protection_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reason": {
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
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"vswitch_id": {
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

func dataSourceAlicloudAlbLoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListLoadBalancers"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = v
	}

	if v, ok := d.GetOk("load_balancer_business_status"); ok {
		request["LoadBalancerBussinessStatus"] = v
	} else if v, ok := d.GetOk("load_balancer_bussiness_status"); ok {
		request["LoadBalancerBussinessStatus"] = v
	}

	if m, ok := d.GetOk("load_balancer_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("LoadBalancerIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerNames.1"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["LoadBalancerStatus"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
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
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcIds.*"] = v
	}
	if m, ok := d.GetOk("vpc_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("VpcIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var loadBalancerNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		loadBalancerNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_load_balancers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.LoadBalancers", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LoadBalancers", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if loadBalancerNameRegex != nil && !loadBalancerNameRegex.MatchString(fmt.Sprint(item["LoadBalancerName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["LoadBalancerId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
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
			"access_log_config":              object["AccessLogConfig"],
			"address_allocated_mode":         object["AddressAllocatedMode"],
			"address_type":                   object["AddressType"],
			"bandwidth_package_id":           object["BandwidthPackageId"],
			"create_time":                    object["CreateTime"],
			"deletion_protection_config":     object["DeletionProtectionConfig"],
			"dns_name":                       object["DNSName"],
			"load_balancer_business_status":  object["LoadBalancerBussinessStatus"],
			"load_balancer_bussiness_status": object["LoadBalancerBussinessStatus"],
			"load_balancer_edition":          object["LoadBalancerEdition"],
			"id":                             fmt.Sprint(object["LoadBalancerId"]),
			"load_balancer_id":               fmt.Sprint(object["LoadBalancerId"]),
			"load_balancer_name":             object["LoadBalancerName"],
			"modification_protection_config": object["ModificationProtectionConfig"],
			"resource_group_id":              object["ResourceGroupId"],
			"vpc_id":                         object["VpcId"],
		}

		loadBalancerOperationLocks := make([]map[string]interface{}, 0)
		if loadBalancerOperationLocksList, ok := object["LoadBalancerOperationLocks"].([]interface{}); ok {
			for _, v := range loadBalancerOperationLocksList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"lock_reason": m1["LockReason"],
						"lock_type":   m1["LockType"],
					}
					loadBalancerOperationLocks = append(loadBalancerOperationLocks, temp1)
				}
			}
		}
		mapping["load_balancer_operation_locks"] = loadBalancerOperationLocks

		accessLogConfigSli := make([]map[string]interface{}, 0)
		if len(object["AccessLogConfig"].(map[string]interface{})) > 0 {
			accessLogConfig := object["AccessLogConfig"]
			accessLogConfigMap := make(map[string]interface{})
			accessLogConfigMap["log_project"] = accessLogConfig.(map[string]interface{})["LogProject"]
			accessLogConfigMap["log_store"] = accessLogConfig.(map[string]interface{})["LogStore"]
			accessLogConfigSli = append(accessLogConfigSli, accessLogConfigMap)
		}
		mapping["access_log_config"] = accessLogConfigSli

		deletionProtectionConfigSli := make([]map[string]interface{}, 0)
		if len(object["DeletionProtectionConfig"].(map[string]interface{})) > 0 {
			deletionProtectionConfig := object["DeletionProtectionConfig"]
			deletionProtectionConfigMap := make(map[string]interface{})
			deletionProtectionConfigMap["enabled"] = deletionProtectionConfig.(map[string]interface{})["Enabled"]
			deletionProtectionConfigMap["enabled_time"] = deletionProtectionConfig.(map[string]interface{})["EnabledTime"]
			deletionProtectionConfigSli = append(deletionProtectionConfigSli, deletionProtectionConfigMap)
		}
		mapping["deletion_protection_config"] = deletionProtectionConfigSli

		modificationProtectionConfigSli := make([]map[string]interface{}, 0)
		if len(object["ModificationProtectionConfig"].(map[string]interface{})) > 0 {
			modificationProtectionConfig := object["ModificationProtectionConfig"]
			modificationProtectionConfigMap := make(map[string]interface{})
			modificationProtectionConfigMap["reason"] = modificationProtectionConfig.(map[string]interface{})["Reason"]
			modificationProtectionConfigMap["status"] = modificationProtectionConfig.(map[string]interface{})["Status"]
			modificationProtectionConfigSli = append(modificationProtectionConfigSli, modificationProtectionConfigMap)
		}
		mapping["modification_protection_config"] = modificationProtectionConfigSli

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags", object)
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
		names = append(names, object["LoadBalancerName"])

		id := fmt.Sprint(object["LoadBalancerId"])
		albService := AlbService{client}
		getResp, err := albService.DescribeAlbLoadBalancer(id)

		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getResp["LoadBalancerStatus"].(string) {
			continue
		}
		mapping["status"] = getResp["LoadBalancerStatus"]

		zoneMappings := make([]map[string]interface{}, 0)
		if zoneMappingsList, ok := getResp["ZoneMappings"].([]interface{}); ok {
			for _, v := range zoneMappingsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"vswitch_id": m1["VSwitchId"],
						"zone_id":    m1["ZoneId"],
					}
					if m1["LoadBalancerAddresses"] != nil {
						loadBalancerAddressesMaps := make([]map[string]interface{}, 0)
						for _, loadBalancerAddressesValue := range m1["LoadBalancerAddresses"].([]interface{}) {
							loadBalancerAddresses := loadBalancerAddressesValue.(map[string]interface{})
							loadBalancerAddressesMap := map[string]interface{}{
								"address": loadBalancerAddresses["Address"],
							}
							loadBalancerAddressesMaps = append(loadBalancerAddressesMaps, loadBalancerAddressesMap)
						}
						temp1["load_balancer_addresses"] = loadBalancerAddressesMaps
					}
					zoneMappings = append(zoneMappings, temp1)
				}
			}
		}
		mapping["zone_mappings"] = zoneMappings

		loadBalancerBillingConfigSli := make([]map[string]interface{}, 0)
		if object["LoadBalancerBillingConfig"] != nil && len(object["LoadBalancerBillingConfig"].(map[string]interface{})) > 0 {
			loadBalancerBillingConfig := object["LoadBalancerBillingConfig"]
			loadBalancerBillingConfigMap := make(map[string]interface{})
			loadBalancerBillingConfigMap["pay_type"] = convertAlbLoadBalancerPaymentTypeResponse(loadBalancerBillingConfig.(map[string]interface{})["PayType"].(string))
			loadBalancerBillingConfigSli = append(loadBalancerBillingConfigSli, loadBalancerBillingConfigMap)
		}
		mapping["load_balancer_billing_config"] = loadBalancerBillingConfigSli

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("balancers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
