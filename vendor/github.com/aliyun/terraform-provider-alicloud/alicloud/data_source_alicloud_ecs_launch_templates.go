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

func dataSourceAlicloudEcsLaunchTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsLaunchTemplatesRead,
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
			"launch_template_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_tags": tagsSchema(),
			"template_resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_release_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delete_with_instance": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"encrypted": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_version_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deployment_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_vm_os_config": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_owner_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_in": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"io_optimized": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_version_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"primary_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password_inherit": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
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
						"security_enhancement_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"spot_duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delete_with_instance": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"template_tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_description": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAlicloudEcsLaunchTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeLaunchTemplates"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("launch_template_name"); ok {
		request["LaunchTemplateName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("template_tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["TemplateTag"] = tags
	}
	if v, ok := d.GetOk("template_resource_group_id"); ok {
		request["TemplateResourceGroupId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var launchTemplateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		launchTemplateNameRegex = r
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
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_launch_templates", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.LaunchTemplateSets.LaunchTemplateSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LaunchTemplateSets.LaunchTemplateSet", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if launchTemplateNameRegex != nil {
				if !launchTemplateNameRegex.MatchString(fmt.Sprint(item["LaunchTemplateName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["LaunchTemplateId"])]; !ok {
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
			"created_by":             object["CreatedBy"],
			"default_version_number": object["DefaultVersionNumber"],
			"latest_version_number":  object["LatestVersionNumber"],
			"id":                     fmt.Sprint(object["LaunchTemplateId"]),
			"launch_template_id":     fmt.Sprint(object["LaunchTemplateId"]),
			"launch_template_name":   object["LaunchTemplateName"],
			"modified_time":          object["ModifiedTime"],
			"resource_group_id":      object["ResourceGroupId"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["template_tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["LaunchTemplateId"]))
			names = append(names, object["LaunchTemplateName"])
			s = append(s, mapping)
			continue
		}

		ecsService := EcsService{client}
		id := fmt.Sprint(object["LaunchTemplateId"])
		getResp, err := ecsService.DescribeLaunchTemplateVersions(id)
		if err != nil {
			return WrapError(err)
		}

		describeLaunchTemplateVersionsObject := make(map[string]interface{})
		for _, version := range getResp {
			if version.(map[string]interface{})["VersionNumber"] == object["LatestVersionNumber"] {
				describeLaunchTemplateVersionsObject = version.(map[string]interface{})
			}
		}

		mapping["auto_release_time"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["AutoReleaseTime"]

		dataDisk := make([]map[string]interface{}, 0)
		if dataDiskList, ok := describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["DataDisks"].(map[string]interface{})["DataDisk"].([]interface{}); ok {
			for _, v := range dataDiskList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"category":             m1["Category"],
						"delete_with_instance": m1["DeleteWithInstance"],
						"description":          m1["Description"],
						"encrypted":            m1["Encrypted"],
						"name":                 m1["DiskName"],
						"performance_level":    m1["PerformanceLevel"],
						"size":                 m1["Size"],
						"snapshot_id":          m1["SnapshotId"],
					}
					dataDisk = append(dataDisk, temp1)
				}
			}
		}
		mapping["data_disks"] = dataDisk
		mapping["deployment_set_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["DeploymentSetId"]
		mapping["description"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["Description"]
		mapping["enable_vm_os_config"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["EnableVmOsConfig"]
		mapping["host_name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["HostName"]
		mapping["image_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ImageId"]
		mapping["image_owner_alias"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ImageOwnerAlias"]
		mapping["instance_charge_type"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceChargeType"]
		mapping["instance_name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceName"]
		mapping["instance_type"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceType"]
		mapping["internet_charge_type"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetChargeType"]
		mapping["internet_max_bandwidth_in"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetMaxBandwidthIn"]
		mapping["internet_max_bandwidth_out"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetMaxBandwidthOut"]
		mapping["io_optimized"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["IoOptimized"]
		mapping["key_pair_name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["KeyPairName"]

		networkInterface := make([]map[string]interface{}, 0)
		if networkInterfaceList, ok := describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["NetworkInterfaces"].(map[string]interface{})["NetworkInterface"].([]interface{}); ok {
			for _, v := range networkInterfaceList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"description":       m1["Description"],
						"name":              m1["NetworkInterfaceName"],
						"primary_ip":        m1["PrimaryIpAddress"],
						"security_group_id": m1["SecurityGroupId"],
						"vswitch_id":        m1["VSwitchId"],
					}
					networkInterface = append(networkInterface, temp1)
				}
			}
		}
		mapping["network_interfaces"] = networkInterface
		mapping["network_type"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["NetworkType"]
		mapping["password_inherit"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["PasswordInherit"]
		mapping["period"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["Period"]
		mapping["private_ip_address"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["PrivateIpAddress"]
		mapping["ram_role_name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["RamRoleName"]
		mapping["security_enhancement_strategy"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityEnhancementStrategy"]
		mapping["security_group_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupId"]
		if describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupIds"] != nil {
			mapping["security_group_ids"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"]
		}
		mapping["spot_duration"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotDuration"]
		mapping["spot_price_limit"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotPriceLimit"]
		mapping["spot_strategy"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotStrategy"]
		systemDiskSli := make([]map[string]interface{}, 0)
		systemDiskMap := make(map[string]interface{})
		systemDiskMap["category"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Category"]
		systemDiskMap["delete_with_instance"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.DeleteWithInstance"]
		systemDiskMap["description"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Description"]
		systemDiskMap["iops"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Iops"]
		systemDiskMap["name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.DiskName"]
		systemDiskMap["performance_level"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.PerformanceLevel"]
		systemDiskMap["size"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Size"]
		systemDiskSli = append(systemDiskSli, systemDiskMap)
		mapping["system_disk"] = systemDiskSli
		mapping["user_data"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["UserData"]
		mapping["vswitch_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["VSwitchId"]
		mapping["version_description"] = describeLaunchTemplateVersionsObject["VersionDescription"]
		mapping["vpc_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["VpcId"]
		mapping["zone_id"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ZoneId"]
		ids = append(ids, fmt.Sprint(object["LaunchTemplateId"]))
		names = append(names, object["LaunchTemplateName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
