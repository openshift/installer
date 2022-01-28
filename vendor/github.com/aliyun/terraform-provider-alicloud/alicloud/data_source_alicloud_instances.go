package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				//must contain a valid status, expected Creating, Starting, Running, Stopping, Stopped
				ValidateFunc: validation.StringInSlice([]string{
					string(Running),
					string(Stopped),
					string(Creating),
					string(Starting),
					string(Stopping),
				}, false),
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							//Set:      imageDiskDeviceMappingHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"category": {
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
						"tags": tagsSchema(),
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request.InstanceIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeInstancesTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	var allInstances []ecs.Instance

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request.PageNumber = requests.NewInteger(v.(int))
	} else {
		request.PageNumber = requests.NewInteger(1)
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request.PageSize = requests.NewInteger(v.(int))
	} else {
		request.PageSize = requests.NewInteger(PageSizeLarge)
	}
	var response *ecs.DescribeInstancesResponse
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeInstancesResponse)
		if isPagingRequest(d) {
			allInstances = response.Instances.Instance
			break
		}
		if len(response.Instances.Instance) < 1 {
			break
		}

		allInstances = append(allInstances, response.Instances.Instance...)
		pageSize, err := strconv.Atoi(string(request.PageSize))
		if err != nil {
			return WrapError(err)
		}

		if len(response.Instances.Instance) < pageSize {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredInstancesTemp []ecs.Instance

	nameRegex, ok := d.GetOk("name_regex")
	imageId, okImg := d.GetOk("image_id")
	if (ok && nameRegex.(string) != "") || (okImg && imageId.(string) != "") {
		var instanceNameRegex *regexp.Regexp
		if nameRegex != "" {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			instanceNameRegex = r
		}
		for _, inst := range allInstances {
			if instanceNameRegex != nil && !instanceNameRegex.MatchString(inst.InstanceName) {
				continue
			}
			if imageId.(string) != "" && inst.ImageId != imageId.(string) {
				continue
			}
			filteredInstancesTemp = append(filteredInstancesTemp, inst)
		}
	} else {
		filteredInstancesTemp = allInstances
	}
	// Filter by ram role name and fetch the instance role name
	instanceIds := make([]string, 0)
	for _, inst := range filteredInstancesTemp {
		if inst.InstanceNetworkType == "classic" {
			continue
		}
		instanceIds = append(instanceIds, inst.InstanceId)
	}
	instanceRoleNameMap := make(map[string]string)
	for index := 0; index < len(instanceIds); index += 100 {
		// DescribeInstanceRamRole parameter InstanceIds supports at most 100 items once
		request := ecs.CreateDescribeInstanceRamRoleRequest()
		request.InstanceIds = convertListToJsonString(convertListStringToListInterface(instanceIds[index:IntMin(index+100, len(instanceIds))]))
		request.RamRoleName = d.Get("ram_role_name").(string)
		request.PageSize = requests.NewInteger(PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)
		for {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeInstanceRamRole(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
			if len(response.InstanceRamRoleSets.InstanceRamRoleSet) < 1 {
				break
			}
			for _, role := range response.InstanceRamRoleSets.InstanceRamRoleSet {
				instanceRoleNameMap[role.InstanceId] = role.RamRoleName
			}

			if len(response.InstanceRamRoleSets.InstanceRamRoleSet) < PageSizeLarge {
				break
			}

			if page, err := getNextpageNumber(request.PageNumber); err != nil {
				return WrapError(err)
			} else {
				request.PageNumber = page
			}
		}
	}
	instanceDiskMappings, err := getInstanceDisksMappings(instanceRoleNameMap, meta)
	if err != nil {
		return WrapError(err)
	}

	return instancessDescriptionAttributes(d, filteredInstancesTemp, instanceRoleNameMap, instanceDiskMappings, meta, response.TotalCount)
}

// populate the numerous fields that the instance description returns.
func instancessDescriptionAttributes(d *schema.ResourceData, instances []ecs.Instance, instanceRoleNameMap map[string]string, instanceDisksMap map[string][]map[string]interface{}, meta interface{}, totalCount int) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, inst := range instances {
		// if instance can not in instanceRoleNameMap, it should be removed.
		if _, ok := instanceRoleNameMap[inst.InstanceId]; !ok {
			continue
		}
		mapping := map[string]interface{}{
			"id":                         inst.InstanceId,
			"region_id":                  inst.RegionId,
			"availability_zone":          inst.ZoneId,
			"status":                     inst.Status,
			"name":                       inst.InstanceName,
			"instance_type":              inst.InstanceType,
			"vpc_id":                     inst.VpcAttributes.VpcId,
			"vswitch_id":                 inst.VpcAttributes.VSwitchId,
			"image_id":                   inst.ImageId,
			"description":                inst.Description,
			"security_groups":            inst.SecurityGroupIds.SecurityGroupId,
			"resource_group_id":          inst.ResourceGroupId,
			"eip":                        inst.EipAddress.IpAddress,
			"key_name":                   inst.KeyPairName,
			"ram_role_name":              instanceRoleNameMap[inst.InstanceId],
			"spot_strategy":              inst.SpotStrategy,
			"creation_time":              inst.CreationTime,
			"instance_charge_type":       inst.InstanceChargeType,
			"internet_charge_type":       inst.InternetChargeType,
			"internet_max_bandwidth_out": inst.InternetMaxBandwidthOut,
			// Complex types get their own functions
			"disk_device_mappings": instanceDisksMap[inst.InstanceId],
			"tags":                 ecsService.tagsToMap(inst.Tags.Tag),
		}
		if len(inst.InnerIpAddress.IpAddress) > 0 {
			mapping["private_ip"] = inst.InnerIpAddress.IpAddress[0]
		} else {
			mapping["private_ip"] = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			mapping["public_ip"] = inst.PublicIpAddress.IpAddress[0]
		} else {
			mapping["public_ip"] = inst.VpcAttributes.NatIpAddress
		}

		ids = append(ids, inst.InstanceId)
		names = append(names, inst.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("names", names)
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", totalCount); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

//Returns a mapping of instance disks
func getInstanceDisksMappings(instanceMap map[string]string, meta interface{}) (map[string][]map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDescribeDisksRequest()
	request.PageSize = requests.NewInteger(PageSizeXLarge)
	request.PageNumber = requests.NewInteger(1)
	instanceDisks := make(map[string][]map[string]interface{})
	var allDisks []ecs.Disk
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			return instanceDisks, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeDisksResponse)

		if response == nil || len(response.Disks.Disk) < 1 {
			break
		}

		allDisks = append(allDisks, response.Disks.Disk...)

		if len(response.Disks.Disk) < PageSizeXLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return instanceDisks, WrapError(err)
		}
		request.PageNumber = page
	}
	for _, disk := range allDisks {
		if _, ok := instanceMap[disk.InstanceId]; !ok {
			continue
		}
		mapping := map[string]interface{}{
			"device":   disk.Device,
			"size":     disk.Size,
			"category": disk.Category,
			"type":     disk.Type,
		}
		instanceDisks[disk.InstanceId] = append(instanceDisks[disk.InstanceId], mapping)
	}

	return instanceDisks, nil
}
