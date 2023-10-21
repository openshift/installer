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

func dataSourceAlicloudConfigAggregators() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigAggregatorsRead,
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
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "Normal"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aggregators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_type": {
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
						"aggregator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudConfigAggregatorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAggregators"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var aggregatorNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		aggregatorNameRegex = r
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
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-09-07"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_aggregators", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if aggregatorNameRegex != nil {
				if !aggregatorNameRegex.MatchString(fmt.Sprint(item["AggregatorName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AggregatorId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != convertConfigAggregatorStatusResponse(formatInt(item["AggregatorStatus"])) {
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
			"account_id":      fmt.Sprint(object["AccountId"]),
			"id":              fmt.Sprint(object["AggregatorId"]),
			"aggregator_id":   fmt.Sprint(object["AggregatorId"]),
			"aggregator_name": object["AggregatorName"],
			"aggregator_type": object["AggregatorType"],
			"description":     object["Description"],
			"status":          convertConfigAggregatorStatusResponse(formatInt(object["AggregatorStatus"])),
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, object["AggregatorName"])
			s = append(s, mapping)
			continue
		}

		configService := ConfigService{client}
		id := fmt.Sprint(object["AggregatorId"])
		getResp, err := configService.DescribeConfigAggregator(id)
		if err != nil {
			return WrapError(err)
		}

		aggregatorAccounts := make([]map[string]interface{}, 0)
		if aggregatorAccountsList, ok := getResp["AggregatorAccounts"].([]interface{}); ok {
			for _, v := range aggregatorAccountsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"account_id":   fmt.Sprint(m1["AccountId"]),
						"account_name": m1["AccountName"],
						"account_type": m1["AccountType"],
					}
					aggregatorAccounts = append(aggregatorAccounts, temp1)
				}
			}
		}
		mapping["aggregator_accounts"] = aggregatorAccounts
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AggregatorName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("aggregators", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
