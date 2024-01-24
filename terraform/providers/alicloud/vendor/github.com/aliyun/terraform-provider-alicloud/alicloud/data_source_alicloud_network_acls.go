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

func dataSourceAlicloudNetworkAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNetworkAclsRead,
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
			"network_acl_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Modifying"}, false),
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
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"egress_acl_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_entry_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ingress_acl_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_entry_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_acl_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAlicloudNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeNetworkAcls"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("network_acl_name"); ok {
		request["NetworkAclName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_id"); ok {
		request["ResourceId"] = v
	}
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var networkAclNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		networkAclNameRegex = r
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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_network_acls", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NetworkAcls.NetworkAcl", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkAcls.NetworkAcl", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if networkAclNameRegex != nil {
				if !networkAclNameRegex.MatchString(fmt.Sprint(item["NetworkAclName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["NetworkAclId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
			"description":      object["Description"],
			"id":               fmt.Sprint(object["NetworkAclId"]),
			"network_acl_id":   fmt.Sprint(object["NetworkAclId"]),
			"network_acl_name": object["NetworkAclName"],
			"status":           object["Status"],
			"vpc_id":           object["VpcId"],
		}

		egressAclEntry := make([]map[string]interface{}, 0)
		if egressAclEntryList, ok := object["EgressAclEntries"].(map[string]interface{})["EgressAclEntry"].([]interface{}); ok {
			for _, v := range egressAclEntryList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"description":            m1["Description"],
						"destination_cidr_ip":    m1["DestinationCidrIp"],
						"network_acl_entry_name": m1["NetworkAclEntryName"],
						"policy":                 m1["Policy"],
						"port":                   m1["Port"],
						"protocol":               m1["Protocol"],
					}
					egressAclEntry = append(egressAclEntry, temp1)
				}
			}
		}
		mapping["egress_acl_entries"] = egressAclEntry

		ingressAclEntry := make([]map[string]interface{}, 0)
		if ingressAclEntryList, ok := object["IngressAclEntries"].(map[string]interface{})["IngressAclEntry"].([]interface{}); ok {
			for _, v := range ingressAclEntryList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"description":            m1["Description"],
						"network_acl_entry_name": m1["NetworkAclEntryName"],
						"policy":                 m1["Policy"],
						"port":                   m1["Port"],
						"protocol":               m1["Protocol"],
						"source_cidr_ip":         m1["SourceCidrIp"],
					}
					ingressAclEntry = append(ingressAclEntry, temp1)
				}
			}
		}
		mapping["ingress_acl_entries"] = ingressAclEntry

		resourceMap := make([]map[string]interface{}, 0)
		if resourceMapList, ok := object["Resources"].(map[string]interface{})["Resource"].([]interface{}); ok {
			for _, v := range resourceMapList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"resource_id":   m1["ResourceId"],
						"resource_type": m1["ResourceType"],
						"status":        m1["Status"],
					}
					resourceMap = append(resourceMap, temp1)
				}
			}
		}
		mapping["resources"] = resourceMap
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["NetworkAclName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("acls", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
