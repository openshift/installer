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

func dataSourceAlicloudScdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudScdnDomainsRead,
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"check_failed", "checking", "configure_failed", "configuring", "offline", "online"}, false),
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
						"cert_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cert_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_pub": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
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
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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

func dataSourceAlicloudScdnDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeScdnUserDomains"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["DomainStatus"] = v
	}
	request["PageSize"] = PageSizeXLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
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
	conn, err := client.NewScdnClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_scdn_domains", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if domainNameRegex != nil && !domainNameRegex.MatchString(fmt.Sprint(item["DomainName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DomainName"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cname":             object["Cname"],
			"create_time":       object["GmtCreated"],
			"description":       object["Description"],
			"id":                fmt.Sprint(object["DomainName"]),
			"domain_name":       fmt.Sprint(object["DomainName"]),
			"gmt_modified":      object["GmtModified"],
			"resource_group_id": object["ResourceGroupId"],
			"status":            object["DomainStatus"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DomainName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DomainName"])
		scdnService := ScdnService{client}
		getResp1, err := scdnService.DescribeScdnDomainCertificateInfo(id)
		if err != nil {
			return WrapError(err)
		}

		if certInfosMap, ok := getResp1["CertInfos"].(map[string]interface{}); ok && certInfosMap != nil {
			if certInfoList, ok := certInfosMap["CertInfo"]; ok && certInfoList != nil {
				certInfosMaps := make([]map[string]interface{}, 0)
				for _, certInfoListItem := range certInfoList.([]interface{}) {
					if certInfoListItemMap, ok := certInfoListItem.(map[string]interface{}); ok {
						certInfoListItemMap["cert_name"] = certInfoListItemMap["CertName"]
						certInfoListItemMap["cert_type"] = certInfoListItemMap["CertType"]
						certInfoListItemMap["ssl_protocol"] = certInfoListItemMap["SslProtocol"]
						certInfoListItemMap["ssl_pub"] = certInfoListItemMap["SslPub"]
						certInfosMaps = append(certInfosMaps, certInfoListItemMap)
					}
				}
				mapping["cert_infos"] = certInfosMaps
			}
		}
		getResp4, err := scdnService.DescribeScdnDomain(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp4["Sources"].(map[string]interface{})["Source"].([]interface{}); ok {
			source := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"content":  item["Content"],
					"enabled":  item["Enabled"],
					"port":     item["Port"],
					"priority": item["Priority"],
					"type":     item["Type"],
				}

				source = append(source, temp)
			}
			mapping["sources"] = source
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

	if err := d.Set("domains", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
