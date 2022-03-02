package alicloud

import (
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudFcFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcFunctionsRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"functions": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"initializer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initialization_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"environment_variables": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"code_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_checksum": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ca_port": {
							Type:     schema.TypeInt,
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
						"custom_container_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"command": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"args": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)

	var ids []string
	var names []string
	var functionMappings []map[string]interface{}
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
		request := fc.NewListFunctionsInput(serviceName)
		if nextToken != "" {
			request.NextToken = &nextToken
		}
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.ListFunctions(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_functions", "ListFunctions", FcGoSdk)
		}
		addDebug("ListFunctions", raw, requestInfo, request)
		response, _ := raw.(*fc.ListFunctionsOutput)

		if response.Functions == nil || len(response.Functions) < 1 {
			break
		}

		for _, function := range response.Functions {
			mapping := map[string]interface{}{
				"id":                     *function.FunctionID,
				"name":                   *function.FunctionName,
				"description":            *function.Description,
				"runtime":                *function.Runtime,
				"handler":                *function.Handler,
				"timeout":                *function.Timeout,
				"memory_size":            *function.MemorySize,
				"code_size":              *function.CodeSize,
				"code_checksum":          *function.CodeChecksum,
				"creation_time":          *function.CreatedTime,
				"last_modification_time": *function.LastModifiedTime,
				"environment_variables":  function.EnvironmentVariables,
			}

			if function.Initializer != nil {
				mapping["initializer"] = *function.Initializer
			}

			if function.InitializationTimeout != nil {
				mapping["initialization_timeout"] = *function.InitializationTimeout
			}

			if function.InstanceConcurrency != nil {
				mapping["instance_concurrency"] = *function.InstanceConcurrency
			}

			if function.InstanceType != nil {
				mapping["instance_type"] = *function.InstanceType
			}

			if function.CAPort != nil {
				mapping["ca_port"] = *function.CAPort
			}

			if function.CustomContainerConfig != nil {
				var cfgList []map[string]interface{}
				cfg := map[string]interface{}{
					"image": *function.CustomContainerConfig.Image,
				}
				if function.CustomContainerConfig.Command != nil {
					cfg["command"] = *function.CustomContainerConfig.Command
				}
				if function.CustomContainerConfig.Args != nil {
					cfg["args"] = *function.CustomContainerConfig.Args
				}
				cfgList = append(cfgList, cfg)
				mapping["custom_container_config"] = cfgList
			}

			// Filter by function name.
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
			// Filter by function id.
			if len(idsMap) > 0 {
				if _, ok := idsMap[*function.FunctionID]; !ok {
					continue
				}
			}
			functionMappings = append(functionMappings, mapping)
			ids = append(ids, *function.FunctionID)
			names = append(names, *function.FunctionName)
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
	if err := d.Set("functions", functionMappings); err != nil {
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
		writeToFile(output.(string), functionMappings)
	}
	return nil
}
