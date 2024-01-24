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

func dataSourceAlicloudDdoscooDomainResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdoscooDomainResourcesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"query_domain_pattern": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"fuzzy", "exact"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cc_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cc_rule_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cc_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http2_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"https_ext": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"policy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"proxy_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proxy_ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"proxy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"real_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"rs_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ssl_ciphers": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_protocols": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDdoscooDomainResourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDomainResource"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		request["InstanceIds"] = v
	}
	if v, ok := d.GetOk("query_domain_pattern"); ok {
		request["QueryDomainPattern"] = v
	}
	request["PageSize"] = PageSizeSmall
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
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_domain_resources", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.WebRules", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.WebRules", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Domain"])]; !ok {
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
			"black_list":      object["BlackList"],
			"cc_enabled":      object["CcEnabled"],
			"cc_rule_enabled": object["CcRuleEnabled"],
			"cc_template":     object["CcTemplate"],
			"cert_name":       object["CertName"],
			"id":              fmt.Sprint(object["Domain"]),
			"domain":          fmt.Sprint(object["Domain"]),
			"http2_enable":    object["Http2Enable"],
			"https_ext":       object["HttpsExt"],
			"instance_ids":    object["InstanceIds"],
			"policy_mode":     object["PolicyMode"],
			"proxy_enabled":   object["ProxyEnabled"],
			"real_servers":    object["RealServers"],
			"rs_type":         formatInt(object["RsType"]),
			"ssl_ciphers":     object["SslCiphers"],
			"ssl_protocols":   object["SslProtocols"],
			"white_list":      object["WhiteList"],
		}

		proxyTypes := make([]map[string]interface{}, 0)
		if proxyTypesList, ok := object["ProxyTypes"].([]interface{}); ok {
			for _, v := range proxyTypesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"proxy_ports": m1["ProxyPorts"],
						"proxy_type":  m1["ProxyType"],
					}
					proxyTypes = append(proxyTypes, temp1)
				}
			}
		}
		mapping["proxy_types"] = proxyTypes
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("resources", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
