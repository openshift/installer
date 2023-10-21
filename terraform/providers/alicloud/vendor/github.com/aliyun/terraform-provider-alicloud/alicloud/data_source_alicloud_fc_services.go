package alicloud

import (
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudFcServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcServicesRead,

		Schema: map[string]*schema.Schema{
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
			"services": {
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_access": {
							Type:     schema.TypeBool,
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
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"logstore": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
						"vpc_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"vswitch_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
						"nas_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"mount_points": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_addr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_dir": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var names []string
	var serviceMappings []map[string]interface{}
	nextToken := ""
	for {
		request := fc.NewListServicesInput()
		if nextToken != "" {
			request.NextToken = &nextToken
		}
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.ListServices(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_services", "ListServices", FcGoSdk)
		}
		addDebug("ListServices", raw, requestInfo, request)
		response, _ := raw.(*fc.ListServicesOutput)

		if response.Services == nil || len(response.Services) < 1 {
			break
		}
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
		for _, service := range response.Services {
			mapping := map[string]interface{}{
				"id":                     *service.ServiceID,
				"name":                   *service.ServiceName,
				"description":            *service.Description,
				"role":                   *service.Role,
				"internet_access":        *service.InternetAccess,
				"creation_time":          *service.CreatedTime,
				"last_modification_time": *service.LastModifiedTime,
			}

			var logConfigMappings []map[string]interface{}
			if service.LogConfig != nil {
				logConfigMappings = append(logConfigMappings, map[string]interface{}{
					"project":  *service.LogConfig.Project,
					"logstore": *service.LogConfig.Logstore,
				})
			}
			mapping["log_config"] = logConfigMappings

			var vpcConfigMappings []map[string]interface{}
			if service.VPCConfig != nil &&
				(service.VPCConfig.VPCID != nil || service.VPCConfig.SecurityGroupID != nil) {
				vpcConfigMappings = append(vpcConfigMappings, map[string]interface{}{
					"vpc_id":            *service.VPCConfig.VPCID,
					"vswitch_ids":       service.VPCConfig.VSwitchIDs,
					"security_group_id": *service.VPCConfig.SecurityGroupID,
				})
			}
			mapping["vpc_config"] = vpcConfigMappings

			var nasConfigMappings []map[string]interface{}
			if service.NASConfig != nil {
				nasConfig := map[string]interface{}{
					"user_id":  *service.NASConfig.UserID,
					"group_id": *service.NASConfig.GroupID,
				}
				var mountPoints []map[string]interface{}
				for _, v := range service.NASConfig.MountPoints {
					mountPoints = append(mountPoints, map[string]interface{}{
						"server_addr": v.ServerAddr,
						"mount_dir":   v.MountDir,
					})
				}
				nasConfig["mount_points"] = mountPoints
				nasConfigMappings = append(nasConfigMappings, nasConfig)
			}
			mapping["nas_config"] = nasConfigMappings

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
				if _, ok := idsMap[*service.ServiceID]; !ok {
					continue
				}
			}
			serviceMappings = append(serviceMappings, mapping)
			ids = append(ids, *service.ServiceID)
			names = append(names, *service.ServiceName)
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
	if err := d.Set("services", serviceMappings); err != nil {
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
		writeToFile(output.(string), serviceMappings)
	}
	return nil
}
