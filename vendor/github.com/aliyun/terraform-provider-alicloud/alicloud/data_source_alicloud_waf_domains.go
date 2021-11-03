package alicloud

import (
	"fmt"
	"regexp"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudWafDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudWafDomainsRead,
		Schema: map[string]*schema.Schema{
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http2_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"http_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"http_to_user_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"https_redirect": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_access_product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancing": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_headers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"read_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"write_time": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudWafDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDomainNames"
	request := make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	var domainNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		domainNameRegex = r
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
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_waf_domains", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range response["DomainNames"].([]interface{}) {
		if domainNameRegex != nil {
			if !domainNameRegex.MatchString(object.(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[object.(string)]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object),
			"domain_name": object,
			"domain":      object,
		}
		ids = append(ids, fmt.Sprint(object.(string)))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object.(string))
			s = append(s, mapping)
			continue
		}

		waf_openapiService := Waf_openapiService{client}
		id := fmt.Sprint(request["InstanceId"], ":", object.(string))
		getResp, err := waf_openapiService.DescribeWafDomain(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["cluster_type"] = convertClusterTypeResponse(formatInt(getResp["ClusterType"]))
		mapping["cname"] = getResp["Cname"]
		mapping["connection_time"] = getResp["ConnectionTime"]
		mapping["http2_port"] = convertJsonStringToStringList(getResp["Http2Port"])
		mapping["http_port"] = convertJsonStringToStringList(getResp["HttpPort"])
		mapping["http_to_user_ip"] = convertHttpToUserIpResponse(formatInt(getResp["HttpToUserIp"]))
		mapping["https_port"] = convertJsonStringToStringList(getResp["HttpsPort"])
		mapping["https_redirect"] = convertHttpsRedirectResponse(formatInt(getResp["HttpsRedirect"]))
		mapping["is_access_product"] = convertIsAccessProductResponse(formatInt(getResp["IsAccessProduct"]))
		mapping["load_balancing"] = convertLoadBalancingResponse(formatInt(getResp["LoadBalancing"]))
		if v, ok := getResp["LogHeaders"].([]interface{}); ok {
			logHeaders := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				logHeaders = append(logHeaders, map[string]interface{}{
					"key":   item["k"].(string),
					"value": item["v"].(string),
				})
			}
			mapping["log_headers"] = logHeaders
		}
		mapping["read_time"] = getResp["ReadTime"]
		mapping["resource_group_id"] = getResp["ResourceGroupId"]
		mapping["source_ips"] = getResp["SourceIps"]
		mapping["version"] = getResp["Version"]
		mapping["write_time"] = getResp["WriteTime"]
		names = append(names, object.(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("domains", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
