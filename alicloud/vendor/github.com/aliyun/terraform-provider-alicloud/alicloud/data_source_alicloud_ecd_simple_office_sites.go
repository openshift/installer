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

func dataSourceAlicloudEcdSimpleOfficeSites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdSimpleOfficeSitesRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGISTERED", "REGISTERING"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sites": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:       schema.TypeInt,
							Computed:   true,
							Deprecated: "Field 'bandwidth' has been deprecated from provider version 1.142.0.",
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_vpc_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dns_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_admin_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_cross_desktop_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_internet_access": {
							Type:       schema.TypeBool,
							Computed:   true,
							Deprecated: "Field 'enable_internet_access' has been deprecated from provider version 1.142.0.",
						},
						"file_system_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"mfa_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"simple_office_site_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sso_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sso_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_dns_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"sub_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trust_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdSimpleOfficeSitesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeOfficeSites"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var simpleOfficeSiteNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		simpleOfficeSiteNameRegex = r
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
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_simple_office_sites", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.OfficeSites", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.OfficeSites", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if simpleOfficeSiteNameRegex != nil && !simpleOfficeSiteNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["OfficeSiteId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
			// todo:  bandwidth depends on network_package resource， you can find it in alicloud_ecd_network_packages
			//"bandwidth":                   formatInt(object["Bandwidth"]),
			//"enable_internet_access":      object["EnableInternetAccess"],
			"cen_id":                      object["CenId"],
			"cidr_block":                  object["CidrBlock"],
			"create_time":                 object["CreationTime"],
			"custom_security_group_id":    object["CustomSecurityGroupId"],
			"desktop_access_type":         convertDesktopAccessType(object["DesktopAccessType"].(string)),
			"desktop_vpc_endpoint":        object["DesktopVpcEndpoint"],
			"dns_address":                 object["DnsAddress"],
			"dns_user_name":               object["DnsUserName"],
			"domain_name":                 object["DomainName"],
			"domain_password":             object["DomainPassword"],
			"domain_user_name":            object["DomainUserName"],
			"enable_admin_access":         object["EnableAdminAccess"],
			"enable_cross_desktop_access": object["EnableCrossDesktopAccess"],
			"enable_internet_access":      object["EnableInternetAccess"],
			"file_system_ids":             object["FileSystemIds"],
			"mfa_enabled":                 object["MfaEnabled"],
			"network_package_id":          object["NetworkPackageId"],
			"id":                          object["OfficeSiteId"],
			"office_site_id":              object["OfficeSiteId"],
			"office_site_type":            object["OfficeSiteType"],
			"simple_office_site_name":     object["Name"],
			"sso_enabled":                 object["SsoEnabled"],
			"sso_status":                  response["SsoStatus"],
			"status":                      object["Status"],
			"sub_dns_address":             object["SubDnsAddress"],
			"sub_domain_name":             object["SubDomainName"],
			"trust_password":              object["TrustPassword"],
			"vswitch_ids":                 object["VSwitchIds"],
			"vpc_id":                      object["VpcId"],
		}

		ids = append(ids, fmt.Sprint(object["OfficeSiteId"]))
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

	if err := d.Set("sites", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
