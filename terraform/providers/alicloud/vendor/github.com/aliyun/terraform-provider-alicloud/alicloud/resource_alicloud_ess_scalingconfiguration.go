package alicloud

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingConfigurationCreate,
		Read:   resourceAliyunEssScalingConfigurationRead,
		Update: resourceAliyunEssScalingConfigurationUpdate,
		Delete: resourceAliyunEssScalingConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
				ConflictsWith: []string{"instance_types"},
			},
			"instance_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"instance_type"},
				MaxItems:      int(MaxScalingConfigurationInstanceTypes),
			},
			"io_optimized": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute io_optimized has been deprecated on instance resource. All the launched alicloud instances will be IO optimized. Suggest to remove it from your template.",
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"security_group_ids"},
			},
			"security_group_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"security_group_id"},
				Optional:      true,
				MaxItems:      16,
			},
			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PayByBandwidth,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"credit_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(CreditSpecificationStandard),
					string(CreditSpecificationUnlimited),
				}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_ssd", "cloud_essd", "cloud_efficiency"}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(20, 500),
			},
			"system_disk_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_disk": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Attribute device has been deprecated on disk attachment resource. Suggest to remove it from your template.",
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_ids' has been deprecated from provider version 1.6.0. New resource 'alicloud_ess_attachment' replaces it.",
			},

			"substitute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ESS-Instance",
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("password_inherit").(bool)
				},
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_encrypted_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("password_inherit").(bool) || d.Get("password").(string) != ""
				},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"system_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	// Ensure instance_type is generation three
	client := meta.(*connectivity.AliyunClient)
	request, err := buildAlicloudEssScalingConfigurationArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	request.IoOptimized = string(IOOptimized)
	if d.Get("is_outdated").(bool) == true {
		request.IoOptimized = string(NoneOptimized)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingConfiguration(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling, "IncorrectScalingGroupStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateScalingConfigurationResponse)
		d.SetId(response.ScalingConfigurationId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalingconfiguration", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunEssScalingConfigurationUpdate(d, meta)
}

func resourceAliyunEssScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	if d.HasChange("active") {
		c, err := essService.DescribeEssScalingConfiguration(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}

		if d.Get("active").(bool) {
			if c.LifecycleState == string(Inactive) {

				err := essService.ActiveEssScalingConfiguration(c.ScalingGroupId, d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		} else {
			if c.LifecycleState == string(Active) {
				_, err := activeSubstituteScalingConfiguration(d, meta)
				if err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("active")
	}

	if err := enableEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}

	if err := modifyEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}

	d.Partial(false)

	return resourceAliyunEssScalingConfigurationRead(d, meta)
}

func modifyEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyScalingConfigurationRequest()
	request.ScalingConfigurationId = d.Id()

	if d.HasChange("override") {
		request.Override = requests.NewBoolean(d.Get("override").(bool))
		d.SetPartial("override")
	}

	if d.HasChange("password_inherit") {
		request.PasswordInherit = requests.NewBoolean(d.Get("password_inherit").(bool))
		d.SetPartial("password_inherit")
	}
	if d.HasChange("image_id") || d.Get("override").(bool) {
		request.ImageId = d.Get("image_id").(string)
		d.SetPartial("image_id")
	}

	if d.HasChange("image_name") || d.Get("override").(bool) {
		request.ImageName = d.Get("image_name").(string)
		d.SetPartial("image_name")
	}

	hasChangeInstanceType := d.HasChange("instance_type")
	hasChangeInstanceTypes := d.HasChange("instance_types")
	if hasChangeInstanceType || hasChangeInstanceTypes || d.Get("override").(bool) {
		instanceType := d.Get("instance_type").(string)
		instanceTypes := d.Get("instance_types").([]interface{})
		if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
			return fmt.Errorf("instance_type or instance_types must be assigned")
		}
		types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
		if instanceTypes != nil && len(instanceTypes) > 0 {
			types = expandStringList(instanceTypes)
		}
		if instanceType != "" {
			types = append(types, instanceType)
		}
		request.InstanceTypes = &types
	}

	hasChangeSecurityGroupId := d.HasChange("security_group_id")
	hasChangeSecurityGroupIds := d.HasChange("security_group_ids")
	if hasChangeSecurityGroupId || hasChangeSecurityGroupIds || d.Get("override").(bool) {
		securityGroupId := d.Get("security_group_id").(string)
		securityGroupIds := d.Get("security_group_ids").([]interface{})
		if securityGroupId == "" && (securityGroupIds == nil || len(securityGroupIds) == 0) {
			return fmt.Errorf("securityGroupId or securityGroupIds must be assigned")
		}
		if securityGroupIds != nil && len(securityGroupIds) > 0 {
			sgs := expandStringList(securityGroupIds)
			request.SecurityGroupIds = &sgs
		}

		if securityGroupId != "" {
			request.SecurityGroupId = securityGroupId
		}
	}

	if d.HasChange("scaling_configuration_name") {
		request.ScalingConfigurationName = d.Get("scaling_configuration_name").(string)
		d.SetPartial("scaling_configuration_name")
	}

	if d.HasChange("internet_charge_type") {
		request.InternetChargeType = d.Get("internet_charge_type").(string)
		d.SetPartial("internet_charge_type")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
		d.SetPartial("internet_max_bandwidth_out")
	}

	if d.HasChange("credit_specification") {
		request.CreditSpecification = d.Get("credit_specification").(string)
		d.SetPartial("credit_specification")
	}

	if d.HasChange("system_disk_category") {
		request.SystemDiskCategory = d.Get("system_disk_category").(string)
		d.SetPartial("system_disk_category")
	}

	if d.HasChange("system_disk_size") {
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		d.SetPartial("system_disk_size")
	}

	if d.HasChange("system_disk_name") {
		request.SystemDiskDiskName = d.Get("system_disk_name").(string)
		d.SetPartial("system_disk_name")
	}

	if d.HasChange("system_disk_description") {
		request.SystemDiskDescription = d.Get("system_disk_description").(string)
		d.SetPartial("system_disk_description")
	}

	if d.HasChange("system_disk_auto_snapshot_policy_id") {
		request.SystemDiskAutoSnapshotPolicyId = d.Get("system_disk_auto_snapshot_policy_id").(string)
		d.SetPartial("system_disk_auto_snapshot_policy_id")
	}

	if d.HasChange("system_disk_performance_level") {
		request.SystemDiskPerformanceLevel = d.Get("system_disk_performance_level").(string)
		d.SetPartial("system_disk_performance_level")
	}

	if d.HasChange("resource_group_id") {
		request.ResourceGroupId = d.Get("resource_group_id").(string)
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				request.UserData = v.(string)
			} else {
				request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
			}
		}
		d.SetPartial("user_data")
	}

	if d.HasChange("role_name") {
		request.RamRoleName = d.Get("role_name").(string)
		d.SetPartial("role_name")
	}

	if d.HasChange("key_name") {
		request.KeyPairName = d.Get("key_name").(string)
		d.SetPartial("key_name")
	}

	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		d.SetPartial("instance_name")
	}

	if d.HasChange("host_name") {
		request.HostName = d.Get("host_name").(string)
		d.SetPartial("host_name")
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := "{"
			for key, value := range v.(map[string]interface{}) {
				tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
			}
			request.Tags = strings.TrimSuffix(tags, ",") + "}"
		}
		d.SetPartial("tags")
	}

	if d.HasChange("data_disk") {
		dds, ok := d.GetOk("data_disk")
		if ok {
			disks := dds.([]interface{})
			createDataDisks := make([]ess.ModifyScalingConfigurationDataDisk, 0, len(disks))
			for _, e := range disks {
				pack := e.(map[string]interface{})
				dataDisk := ess.ModifyScalingConfigurationDataDisk{
					Size:                 strconv.Itoa(pack["size"].(int)),
					Category:             pack["category"].(string),
					SnapshotId:           pack["snapshot_id"].(string),
					DeleteWithInstance:   strconv.FormatBool(pack["delete_with_instance"].(bool)),
					Device:               pack["device"].(string),
					Encrypted:            strconv.FormatBool(pack["encrypted"].(bool)),
					KMSKeyId:             pack["kms_key_id"].(string),
					DiskName:             pack["name"].(string),
					Description:          pack["description"].(string),
					AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
					PerformanceLevel:     pack["performance_level"].(string),
				}
				createDataDisks = append(createDataDisks, dataDisk)
			}
			request.DataDisk = &createDataDisks
		}
		d.SetPartial("data_disk")
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingConfiguration(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func enableEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if d.HasChange("enable") {
		sgId := d.Get("scaling_group_id").(string)
		group, err := essService.DescribeEssScalingGroup(sgId)
		if err != nil {
			return WrapError(err)
		}

		if d.Get("enable").(bool) {
			if group.LifecycleState == string(Inactive) {

				object, err := essService.DescribeEssScalingConfifurations(sgId)

				if err != nil {
					return WrapError(err)
				}
				activeConfig := ""
				var csIds []string
				for _, c := range object {
					csIds = append(csIds, c.ScalingConfigurationId)
					if c.LifecycleState == string(Active) {
						activeConfig = c.ScalingConfigurationId
					}
				}

				if activeConfig == "" {
					return WrapError(Error("Please active a scaling configuration before enabling scaling group %s. Its all scaling configuration are %s.",
						sgId, strings.Join(csIds, ",")))
				}

				request := ess.CreateEnableScalingGroupRequest()
				request.RegionId = client.RegionId
				request.ScalingGroupId = sgId
				request.ActiveScalingConfigurationId = activeConfig

				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Active, DefaultTimeout); err != nil {
					return WrapError(err)
				}

				d.SetPartial("scaling_configuration_id")
			}
		} else {
			if group.LifecycleState == string(Active) {
				request := ess.CreateDisableScalingGroupRequest()
				request.RegionId = client.RegionId
				request.ScalingGroupId = sgId
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DisableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Inactive, DefaultTimeout); err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("enable")
	}

	return nil
}

func resourceAliyunEssScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}
	object, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("active", object.LifecycleState == string(Active))
	d.Set("image_id", object.ImageId)
	d.Set("image_name", object.ImageName)
	d.Set("scaling_configuration_name", object.ScalingConfigurationName)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("internet_max_bandwidth_in", object.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", object.InternetMaxBandwidthOut)
	d.Set("credit_specification", object.CreditSpecification)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("system_disk_name", object.SystemDiskName)
	d.Set("system_disk_description", object.SystemDiskDescription)
	d.Set("system_disk_auto_snapshot_policy_id", object.SystemDiskAutoSnapshotPolicyId)
	d.Set("system_disk_performance_level", object.SystemDiskPerformanceLevel)
	d.Set("data_disk", essService.flattenDataDiskMappings(object.DataDisks.DataDisk))
	d.Set("role_name", object.RamRoleName)
	d.Set("key_name", object.KeyPairName)
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("tags", essTagsToMap(object.Tags.Tag))
	d.Set("instance_name", object.InstanceName)
	d.Set("override", d.Get("override").(bool))
	d.Set("password_inherit", object.PasswordInherit)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("host_name", object.HostName)

	if sg, ok := d.GetOk("security_group_id"); ok && sg.(string) != "" {
		d.Set("security_group_id", object.SecurityGroupId)
	}
	if sgs, ok := d.GetOk("security_group_ids"); ok && len(sgs.([]interface{})) > 0 {
		d.Set("security_group_ids", object.SecurityGroupIds.SecurityGroupId)
	}
	if instanceType, ok := d.GetOk("instance_type"); ok && instanceType.(string) != "" {
		d.Set("instance_type", object.InstanceType)
	}
	if instanceTypes, ok := d.GetOk("instance_types"); ok && len(instanceTypes.([]interface{})) > 0 {
		d.Set("instance_types", object.InstanceTypes.InstanceType)
	}
	userData := d.Get("user_data")
	if userData.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(userData.(string))
		if base64DecodeError == nil {
			d.Set("user_data", object.UserData)
		} else {
			d.Set("user_data", userDataHashSum(object.UserData))
		}
	} else {
		d.Set("user_data", userDataHashSum(object.UserData))
	}
	return nil
}

func resourceAliyunEssScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	object, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = object.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return nil
	} else if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
		if d.Get("force_delete").(bool) {
			request := ess.CreateDeleteScalingGroupRequest()
			request.ScalingGroupId = object.ScalingGroupId
			request.ForceDelete = requests.NewBoolean(true)
			request.RegionId = client.RegionId
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DeleteScalingGroup(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
					return nil
				}
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
		}
		return WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), object.ScalingGroupId))
	}

	deleteScalingConfigurationRequest := ess.CreateDeleteScalingConfigurationRequest()
	deleteScalingConfigurationRequest.ScalingConfigurationId = d.Id()

	rawDeleteScalingConfiguration, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingConfiguration(deleteScalingConfigurationRequest)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound", "InvalidScalingConfigurationId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), rawDeleteScalingConfiguration, request.RpcRequest, request)

	return WrapError(essService.WaitForScalingConfiguration(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudEssScalingConfigurationArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingConfigurationRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	//ecsService := EcsService{client}
	//zoneId, validZones, _, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	//if err != nil {
	//	return nil, WrapError(err)
	//}

	request := ess.CreateCreateScalingConfigurationRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	request.ImageId = d.Get("image_id").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.PasswordInherit = requests.NewBoolean(d.Get("password_inherit").(bool))

	securityGroupId := d.Get("security_group_id").(string)
	securityGroupIds := d.Get("security_group_ids").([]interface{})

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password != "" {
		request.Password = password
	} else if kmsPassword != "" {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return nil, WrapError(err)
		}
		request.Password = decryptResp
	}

	if securityGroupId == "" && (securityGroupIds == nil || len(securityGroupIds) == 0) {
		return nil, WrapError(Error("security_group_id or security_group_ids must be assigned"))
	}

	if securityGroupIds != nil && len(securityGroupIds) > 0 {
		sgs := expandStringList(securityGroupIds)
		request.SecurityGroupIds = &sgs
	}

	if securityGroupId != "" {
		request.SecurityGroupId = securityGroupId
	}

	types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
	instanceType := d.Get("instance_type").(string)
	instanceTypes := d.Get("instance_types").([]interface{})
	if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
		return nil, WrapError(Error("instance_type or instance_types must be assigned"))
	}

	if instanceTypes != nil && len(instanceTypes) > 0 {
		types = expandStringList(instanceTypes)
	}

	if instanceType != "" {
		types = append(types, instanceType)
	}
	//for _, v := range types {
	//	if err := ecsService.InstanceTypeValidation(v, zoneId, validZones); err != nil {
	//		return nil, WrapError(err)
	//	}
	//}
	request.InstanceTypes = &types

	if v := d.Get("scaling_configuration_name").(string); v != "" {
		request.ScalingConfigurationName = v
	}

	if v := d.Get("image_name").(string); v != "" {
		request.ImageName = v
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		request.InternetChargeType = v
	}

	if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
		request.InternetMaxBandwidthIn = requests.NewInteger(v)
	}

	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))

	if v := d.Get("credit_specification").(string); v != "" {
		request.CreditSpecification = v
	}

	if v := d.Get("system_disk_category").(string); v != "" {
		request.SystemDiskCategory = v
	}

	if v := d.Get("system_disk_size").(int); v != 0 {
		request.SystemDiskSize = requests.NewInteger(v)
	}

	if v := d.Get("system_disk_name").(string); v != "" {
		request.SystemDiskDiskName = v
	}

	if v := d.Get("system_disk_description").(string); v != "" {
		request.SystemDiskDescription = v
	}

	if v := d.Get("system_disk_auto_snapshot_policy_id").(string); v != "" {
		request.SystemDiskAutoSnapshotPolicyId = v
	}

	if v := d.Get("system_disk_performance_level").(string); v != "" {
		request.SystemDiskPerformanceLevel = v
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}

	dds, ok := d.GetOk("data_disk")
	if ok {
		disks := dds.([]interface{})
		createDataDisks := make([]ess.CreateScalingConfigurationDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := ess.CreateScalingConfigurationDataDisk{
				Size:                 strconv.Itoa(pack["size"].(int)),
				Category:             pack["category"].(string),
				SnapshotId:           pack["snapshot_id"].(string),
				DeleteWithInstance:   strconv.FormatBool(pack["delete_with_instance"].(bool)),
				Device:               pack["device"].(string),
				Encrypted:            strconv.FormatBool(pack["encrypted"].(bool)),
				KMSKeyId:             pack["kms_key_id"].(string),
				DiskName:             pack["name"].(string),
				Description:          pack["description"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				PerformanceLevel:     pack["performance_level"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		request.DataDisk = &createDataDisks
	}

	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		request.RamRoleName = v.(string)
	}

	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		request.KeyPairName = v.(string)
	}

	if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			request.UserData = v.(string)
		} else {
			request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := "{"
		for key, value := range v.(map[string]interface{}) {
			tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
		}
		request.Tags = strings.TrimSuffix(tags, ",") + "}"
	}

	if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
		request.InstanceName = v.(string)
	}

	if v, ok := d.GetOk("host_name"); ok && v.(string) != "" {
		request.HostName = v.(string)
	}

	return request, nil
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) (configures []ess.ScalingConfiguration, err error) {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	substituteId, ok := d.GetOk("substitute")

	c, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		err = WrapError(err)
		return
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return
	}

	if !ok || substituteId.(string) == "" {

		if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
			return configures, WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId))
		}

		var configs []string
		for _, cc := range response.ScalingConfigurations.ScalingConfiguration {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return configures, WrapError(Error("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ",")))

	}

	err = essService.ActiveEssScalingConfiguration(c.ScalingGroupId, substituteId.(string))
	if err != nil {
		return configures, WrapError(Error("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substituteId.(string)))
	}

	return response.ScalingConfigurations.ScalingConfiguration, nil
}
