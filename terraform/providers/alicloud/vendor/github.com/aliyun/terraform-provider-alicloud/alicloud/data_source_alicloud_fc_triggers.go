package alicloud

import (
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudFcTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcTriggersRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
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
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"triggers": {
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
						"source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invocation_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcTriggersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)
	functionName := d.Get("function_name").(string)

	var ids []string
	var names []string
	var triggerMappings []map[string]interface{}
	nextToken := ""
	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		request := fc.NewListTriggersInput(serviceName, functionName)
		if nextToken != "" {
			request.NextToken = &nextToken
		}
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.ListTriggers(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_triggers", "ListTriggers", FcGoSdk)
		}
		addDebug("ListTriggers", raw, requestInfo, request)
		response, _ := raw.(*fc.ListTriggersOutput)

		if len(response.Triggers) < 1 {
			break
		}

		for _, trigger := range response.Triggers {
			var sourceARN, invocationRole string
			if trigger.SourceARN != nil {
				sourceARN = *trigger.SourceARN
			}
			if trigger.InvocationRole != nil {
				invocationRole = *trigger.InvocationRole
			}
			mapping := map[string]interface{}{
				"id":                     *trigger.TriggerID,
				"name":                   *trigger.TriggerName,
				"source_arn":             sourceARN,
				"type":                   *trigger.TriggerType,
				"invocation_role":        invocationRole,
				"config":                 string(trigger.RawTriggerConfig),
				"creation_time":          *trigger.CreatedTime,
				"last_modification_time": *trigger.LastModifiedTime,
			}

			nameRegex, ok := d.GetOk("name_regex")
			if ok && nameRegex.(string) != "" {
				var r *regexp.Regexp
				if nameRegex != "" {
					r, err = regexp.Compile(nameRegex.(string))
					if err != nil {
						return WrapError(err)
					}
				}
				if r != nil && !r.MatchString(mapping["name"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[*trigger.TriggerID]; !ok {
					continue
				}
			}
			triggerMappings = append(triggerMappings, mapping)
			ids = append(ids, *trigger.TriggerID)
			names = append(names, *trigger.TriggerName)
		}

		nextToken = ""
		if response.NextToken != nil {
			nextToken = *response.NextToken
		}
		if nextToken == "" {
			break
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("triggers", triggerMappings); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), triggerMappings)
	}
	return nil
}
