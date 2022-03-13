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

func dataSourceAlicloudDirectMailDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDirectMailDomainsRead,
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
			"key_word": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(4, 50),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4"}, false),
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
						"cname_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_confirm_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_mx": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_spf": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_txt": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icp_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mx_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mx_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spf_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spf_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tl_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tracef_record": {
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

func dataSourceAlicloudDirectMailDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryDomainByParam"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNo"] = 1
	if v, ok := d.GetOk("key_word"); ok {
		request["KeyWord"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
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
	conn, err := client.NewDmClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_direct_mail_domains", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.data.domain", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.domain", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if domainNameRegex != nil && !domainNameRegex.MatchString(fmt.Sprint(item["DomainName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DomainId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cname_auth_status": object["CnameAuthStatus"],
			"create_time":       object["CreateTime"],
			"id":                fmt.Sprint(object["DomainId"]),
			"domain_id":         fmt.Sprint(object["DomainId"]),
			"domain_name":       object["DomainName"],
			"icp_status":        object["IcpStatus"],
			"mx_auth_status":    object["MxAuthStatus"],
			"spf_auth_status":   object["SpfAuthStatus"],
			"status":            object["DomainStatus"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DomainName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DomainId"])
		dmService := DmService{client}
		getResp, err := dmService.DescribeDirectMailDomain(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["cname_confirm_status"] = getResp["CnameConfirmStatus"]
		mapping["cname_record"] = getResp["CnameRecord"]
		mapping["default_domain"] = getResp["DefaultDomain"]
		mapping["dns_mx"] = getResp["DnsMx"]
		mapping["dns_spf"] = getResp["DnsSpf"]
		mapping["dns_txt"] = getResp["DnsTxt"]
		mapping["domain_type"] = getResp["DomainType"]
		mapping["mx_record"] = getResp["MxRecord"]
		mapping["spf_record"] = getResp["SpfRecord"]
		mapping["tl_domain_name"] = getResp["TlDomainName"]
		mapping["tracef_record"] = getResp["TracefRecord"]
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
