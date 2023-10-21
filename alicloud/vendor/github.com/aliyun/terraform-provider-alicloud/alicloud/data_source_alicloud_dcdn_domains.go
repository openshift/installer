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

func dataSourceAlicloudDcdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDcdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"change_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"change_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"check_domain_show": {
				Type:     schema.TypeBool,
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
			"domain_search_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_token": {
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
						"cert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
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
						"ssl_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_pub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
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
									"weight": {
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

func dataSourceAlicloudDcdnDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDcdnUserDomains"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("change_end_time"); ok {
		request["ChangeEndTime"] = v
	}
	if v, ok := d.GetOk("change_start_time"); ok {
		request["ChangeStartTime"] = v
	}
	if v, ok := d.GetOkExists("check_domain_show"); ok {
		request["CheckDomainShow"] = v
	}
	if v, ok := d.GetOk("domain_search_type"); ok {
		request["DomainSearchType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("security_token"); ok {
		request["SecurityToken"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["DomainStatus"] = v
	}
	request["PageSize"] = PageSizeLarge
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
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if domainNameRegex != nil {
				if !domainNameRegex.MatchString(item["DomainName"].(string)) {
					continue
				}
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
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cname":             object["Cname"],
			"id":                fmt.Sprint(object["DomainName"]),
			"domain_name":       fmt.Sprint(object["DomainName"]),
			"gmt_modified":      object["GmtModified"],
			"resource_group_id": object["ResourceGroupId"],
			"ssl_protocol":      object["SSLProtocol"],
			"status":            object["DomainStatus"],
		}
		if v, ok := object["Sources"].(map[string]interface{})["Source"].([]interface{}); ok {
			source := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				source = append(source, map[string]interface{}{
					"content":  item["Content"],
					"port":     item["Port"],
					"priority": item["Priority"],
					"type":     item["Type"],
					"weight":   item["Weight"],
				})
			}
			mapping["sources"] = source
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["DomainName"]))
			names = append(names, object["DomainName"].(string))
			s = append(s, mapping)
			continue
		}

		dcdnService := DcdnService{client}
		id := fmt.Sprint(object["DomainName"])
		getResp, err := dcdnService.DescribeDcdnDomainCertificateInfo(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["cert_name"] = getResp["CertName"]
		mapping["ssl_pub"] = getResp["SSLPub"]
		getResp1, err := dcdnService.DescribeDcdnDomain(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["description"] = getResp1["Description"]
		mapping["scope"] = getResp1["Scope"]

		ids = append(ids, fmt.Sprint(object["DomainName"]))
		names = append(names, object["DomainName"].(string))
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
