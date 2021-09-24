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

func dataSourceAlicloudEciContainerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEciContainerGroupsRead,
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
			"container_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{"Pending", "Running", "ScheduleFailed", "Scheduling", "Succeeded"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"with_event": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"args": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"environment_vars": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_pull_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"ready": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"restart_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_mounts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"working_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"discount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dns_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name_servers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"searches": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"eci_security_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sysctls": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"eni_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"first_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
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
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"init_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"args": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"environment_vars": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_pull_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"ready": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"restart_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_mounts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"working_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"succeeded_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_file_volume_config_file_to_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"disk_volume_disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_volume_fs_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flex_volume_driver": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flex_volume_fs_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flex_volume_options": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nfs_volume_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nfs_volume_read_only": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nfs_volume_server": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
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
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
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

func dataSourceAlicloudEciContainerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeContainerGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("container_group_name"); ok {
		request["ContainerGroupName"] = v
	}
	if v, ok := d.GetOk("limit"); ok {
		request["Limit"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOkExists("with_event"); ok {
		request["WithEvent"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var containerGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		containerGroupNameRegex = r
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
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eci_container_groups", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.ContainerGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ContainerGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if containerGroupNameRegex != nil {
				if !containerGroupNameRegex.MatchString(fmt.Sprint(item["ContainerGroupName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ContainerGroupId"])]; !ok {
					continue
				}
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
			"id":                   fmt.Sprint(object["ContainerGroupId"]),
			"container_group_id":   fmt.Sprint(object["ContainerGroupId"]),
			"container_group_name": object["ContainerGroupName"],
			"cpu":                  object["Cpu"],
			"discount":             formatInt(object["Discount"]),
			"eni_instance_id":      object["EniInstanceId"],
			"expired_time":         object["ExpiredTime"],
			"failed_time":          object["FailedTime"],
			"instance_type":        object["InstanceType"],
			"internet_ip":          object["InternetIp"],
			"intranet_ip":          object["IntranetIp"],
			"ipv6_address":         object["Ipv6Address"],
			"memory":               object["Memory"],
			"ram_role_name":        object["RamRoleName"],
			"resource_group_id":    object["ResourceGroupId"],
			"security_group_id":    object["SecurityGroupId"],
			"status":               object["Status"],
			"succeeded_time":       object["SucceededTime"],
			"vswitch_id":           object["VSwitchId"],
			"vpc_id":               object["VpcId"],
			"zone_id":              object["ZoneId"],
		}

		containers := make([]map[string]interface{}, 0)
		if containersList, ok := object["Containers"].([]interface{}); ok {
			for _, v := range containersList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"args":              m1["Args"],
						"commands":          m1["Commands"],
						"cpu":               m1["Cpu"],
						"gpu":               m1["Gpu"],
						"image":             m1["Image"],
						"image_pull_policy": m1["ImagePullPolicy"],
						"memory":            m1["Memory"],
						"name":              m1["Name"],
						"ready":             m1["Ready"],
						"restart_count":     m1["RestartCount"],
						"working_dir":       m1["WorkingDir"],
					}
					if m1["EnvironmentVars"] != nil {
						environmentVarsMaps := make([]map[string]interface{}, 0)
						for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
							environmentVars := environmentVarsValue.(map[string]interface{})
							environmentVarsMap := map[string]interface{}{
								"key":   environmentVars["Key"],
								"value": environmentVars["Value"],
							}
							environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
						}
						temp1["environment_vars"] = environmentVarsMaps
					}
					if m1["Ports"] != nil {
						portsMaps := make([]map[string]interface{}, 0)
						for _, portsValue := range m1["Ports"].([]interface{}) {
							ports := portsValue.(map[string]interface{})
							portsMap := map[string]interface{}{
								"port":     ports["Port"],
								"protocol": ports["Protocol"],
							}
							portsMaps = append(portsMaps, portsMap)
						}
						temp1["ports"] = portsMaps
					}
					if m1["VolumeMounts"] != nil {
						volumeMountsMaps := make([]map[string]interface{}, 0)
						for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
							volumeMounts := volumeMountsValue.(map[string]interface{})
							volumeMountsMap := map[string]interface{}{
								"mount_path": volumeMounts["MountPath"],
								"name":       volumeMounts["Name"],
								"read_only":  volumeMounts["ReadOnly"],
							}
							volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
						}
						temp1["volume_mounts"] = volumeMountsMaps
					}
					containers = append(containers, temp1)
				}
			}
		}
		mapping["containers"] = containers

		events := make([]map[string]interface{}, 0)
		if eventsList, ok := object["Events"].([]interface{}); ok {
			for _, v := range eventsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"count":           m1["Count"],
						"first_timestamp": m1["FirstTimestamp"],
						"last_timestamp":  m1["LastTimestamp"],
						"message":         m1["Message"],
						"name":            m1["Name"],
						"reason":          m1["Reason"],
						"type":            m1["Type"],
					}
					events = append(events, temp1)
				}
			}
		}
		mapping["events"] = events

		hostAliases := make([]map[string]interface{}, 0)
		if hostAliasesList, ok := object["HostAliases"].([]interface{}); ok {
			for _, v := range hostAliasesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"hostnames": m1["Hostnames"],
						"ip":        m1["Ip"],
					}
					hostAliases = append(hostAliases, temp1)
				}
			}
		}
		mapping["host_aliases"] = hostAliases

		initContainers := make([]map[string]interface{}, 0)
		if initContainersList, ok := object["InitContainers"].([]interface{}); ok {
			for _, v := range initContainersList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"args":              m1["Args"],
						"commands":          m1["Command"],
						"cpu":               m1["Cpu"],
						"gpu":               m1["Gpu"],
						"image":             m1["Image"],
						"image_pull_policy": m1["ImagePullPolicy"],
						"memory":            m1["Memory"],
						"name":              m1["Name"],
						"ready":             m1["Ready"],
						"restart_count":     m1["RestartCount"],
						"working_dir":       m1["WorkingDir"],
					}
					if m1["EnvironmentVars"] != nil {
						environmentVarsMaps := make([]map[string]interface{}, 0)
						for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
							environmentVars := environmentVarsValue.(map[string]interface{})
							environmentVarsMap := map[string]interface{}{
								"key":   environmentVars["Key"],
								"value": environmentVars["Value"],
							}
							environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
						}
						temp1["environment_vars"] = environmentVarsMaps
					}
					if m1["Ports"] != nil {
						portsMaps := make([]map[string]interface{}, 0)
						for _, portsValue := range m1["Ports"].([]interface{}) {
							ports := portsValue.(map[string]interface{})
							portsMap := map[string]interface{}{
								"port":     ports["Port"],
								"protocol": ports["Protocol"],
							}
							portsMaps = append(portsMaps, portsMap)
						}
						temp1["ports"] = portsMaps
					}
					if m1["VolumeMounts"] != nil {
						volumeMountsMaps := make([]map[string]interface{}, 0)
						for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
							volumeMounts := volumeMountsValue.(map[string]interface{})
							volumeMountsMap := map[string]interface{}{
								"mount_path": volumeMounts["MountPath"],
								"name":       volumeMounts["Name"],
								"read_only":  volumeMounts["ReadOnly"],
							}
							volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
						}
						temp1["volume_mounts"] = volumeMountsMaps
					}
					initContainers = append(initContainers, temp1)
				}
			}
		}
		mapping["init_containers"] = initContainers

		volumes := make([]map[string]interface{}, 0)
		if volumesList, ok := object["Volumes"].([]interface{}); ok {
			for _, v := range volumesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"disk_volume_disk_id":  m1["DiskVolumeDiskId"],
						"disk_volume_fs_type":  m1["DiskVolumeFsType"],
						"flex_volume_driver":   m1["FlexVolumeDriver"],
						"flex_volume_fs_type":  m1["FlexVolumeFsType"],
						"flex_volume_options":  m1["FlexVolumeOptions"],
						"nfs_volume_path":      m1["NFSVolumePath"],
						"nfs_volume_read_only": m1["NFSVolumeReadOnly"],
						"nfs_volume_server":    m1["NFSVolumeServer"],
						"name":                 m1["Name"],
						"type":                 m1["Type"],
					}
					if m1["ConfigFileVolumeConfigFileToPaths"] != nil {
						configFileVolumeConfigFileToPathsMaps := make([]map[string]interface{}, 0)
						for _, configFileVolumeConfigFileToPathsValue := range m1["ConfigFileVolumeConfigFileToPaths"].([]interface{}) {
							configFileVolumeConfigFileToPaths := configFileVolumeConfigFileToPathsValue.(map[string]interface{})
							configFileVolumeConfigFileToPathsMap := map[string]interface{}{
								"content": configFileVolumeConfigFileToPaths["Content"],
								"path":    configFileVolumeConfigFileToPaths["Path"],
							}
							configFileVolumeConfigFileToPathsMaps = append(configFileVolumeConfigFileToPathsMaps, configFileVolumeConfigFileToPathsMap)
						}
						temp1["config_file_volume_config_file_to_paths"] = configFileVolumeConfigFileToPathsMaps
					}
					volumes = append(volumes, temp1)
				}
			}
		}
		mapping["volumes"] = volumes

		eciSecurityContextSli := make([]map[string]interface{}, 0)
		if len(object["EciSecurityContext"].(map[string]interface{})) > 0 {
			eciSecurityContext := object["EciSecurityContext"]
			eciSecurityContextMap := make(map[string]interface{})

			sysctlsSli := make([]map[string]interface{}, 0)
			if len(eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{})) > 0 {
				for _, sysctls := range eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{}) {
					sysctlsMap := make(map[string]interface{})
					sysctlsMap["name"] = sysctls.(map[string]interface{})["Name"]
					sysctlsMap["value"] = sysctls.(map[string]interface{})["Value"]
					sysctlsSli = append(sysctlsSli, sysctlsMap)
				}
			}
			eciSecurityContextMap["sysctls"] = sysctlsSli
			eciSecurityContextSli = append(eciSecurityContextSli, eciSecurityContextMap)
		}
		mapping["eci_security_context"] = eciSecurityContextSli

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ContainerGroupId"]))
			names = append(names, object["ContainerGroupName"])
			s = append(s, mapping)
			continue
		}

		eciService := EciService{client}
		id := fmt.Sprint(object["ContainerGroupId"])
		getResp, err := eciService.DescribeEciContainerGroup(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["restart_policy"] = getResp["RestartPolicy"]
		ids = append(ids, fmt.Sprint(object["ContainerGroupId"]))
		names = append(names, object["ContainerGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
