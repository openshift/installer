package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudKvstoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKvstoreInstanceCreate,
		Read:   resourceAlicloudKvstoreInstanceRead,
		Update: resourceAlicloudKvstoreInstanceUpdate,
		Delete: resourceAlicloudKvstoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: redisPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: redisPostPaidAndRenewDiffSuppressFunc,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"business_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"parameters"},
			},
			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"coupon_no": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "youhuiquan_promotion_option_id_for_blank",
			},
			"db_instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_name"},
			},
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_name' has been deprecated from version 1.101.0. Use 'db_instance_name' instead.",
				ConflictsWith: []string{"db_instance_name"},
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_backup_log": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Default:      0,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.8", "4.0", "5.0", "6.0"}, false),
				Default:      "5.0",
			},
			"force_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"global_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"global_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_release_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Memcache", "Redis"}, false),
				Default:      "Redis",
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE", "double", "single"}, false),
				Deprecated:   "Field 'node_type' has been deprecated from version 1.120.1",
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
				Default:      "UPGRADE",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"instance_charge_type"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Deprecated:    "Field 'instance_charge_type' has been deprecated from version 1.101.0. Use 'payment_type' instead.",
				ConflictsWith: []string{"payment_type"},
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "12", "2", "24", "3", "36", "4", "5", "6", "7", "8", "9"}, false),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"private_connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_connection_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"restore_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ssl_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disable", "Enable", "Update"}, false),
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: redisSecurityGroupIdDiffSuppressFunc,
			},
			"security_ip_group_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ip_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"srcdb_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_auth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Close", "Open"}, false),
				Default:      "Open",
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'availability_zone' has been deprecated from version 1.101.0. Use 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"enable_public": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'enable_public' has been deprecated from version 1.101.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"connection_string_prefix": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'connection_string_prefix' has been deprecated from version 1.101.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"connection_string": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'connection_string' has been deprecated from version 1.101.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					return hashcode.String(
						v.(map[string]interface{})["name"].(string) + "|" + v.(map[string]interface{})["value"].(string))
				},
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'parameters' has been deprecated from version 1.101.0. Use 'config' instead.",
				ConflictsWith: []string{"config"},
			},
		},
	}
}

func resourceAlicloudKvstoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}

	request := r_kvstore.CreateCreateInstanceRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = convertBoolToString(v.(bool))
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		request.AutoRenewPeriod = convertIntergerToString(v.(int))
	}

	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request.AutoUseCoupon = convertBoolToString(v.(bool))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		request.BackupId = v.(string)
	}

	if v, ok := d.GetOk("business_info"); ok {
		request.BusinessInfo = v.(string)
	}

	if v, ok := d.GetOk("capacity"); ok {
		request.Capacity = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("coupon_no"); ok {
		request.CouponNo = v.(string)
	}

	if v, ok := d.GetOk("db_instance_name"); ok {
		request.InstanceName = v.(string)
	} else if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = v.(string)
	}

	if v, ok := d.GetOk("dedicated_host_group_id"); ok {
		request.DedicatedHostGroupId = v.(string)
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = v.(string)
	}

	if v, ok := d.GetOkExists("global_instance"); ok {
		request.GlobalInstance = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("global_instance_id"); ok {
		request.GlobalInstanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request.InstanceClass = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = v.(string)
	}

	request.NetworkType = "CLASSIC"
	if v, ok := d.GetOk("node_type"); ok {
		request.NodeType = v.(string)
	}

	request.Password = d.Get("password").(string)
	if request.Password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.Password = decryptResp
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request.ChargeType = v.(string)
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request.ChargeType = v.(string)
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = v.(string)
	}

	if v, ok := d.GetOk("private_ip"); ok {
		request.PrivateIpAddress = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("restore_time"); ok {
		request.RestoreTime = v.(string)
	}

	if v, ok := d.GetOk("srcdb_instance_id"); ok {
		request.SrcDBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request.ZoneId = v.(string)
	}

	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request.SecondaryZoneId = v.(string)
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request.NetworkType = "VPC"
		request.VpcId = vsw.VpcId
		request.VSwitchId = vswitchId
		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		}
	}
	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.CreateInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*r_kvstore.CreateInstanceResponse)
	d.SetId(fmt.Sprintf("%v", response.InstanceId))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 300*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudKvstoreInstanceUpdate(d, meta)
}
func resourceAlicloudKvstoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_instance r_kvstoreService.DescribeKvstoreInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("enable_public", false)
	d.Set("connection_string", "")
	net, _ := r_kvstoreService.DescribeKvstoreConnection(d.Id())
	for _, instanceNetInfo := range net {
		if instanceNetInfo.DBInstanceNetType == "0" {
			d.Set("enable_public", true)
			d.Set("connection_string", instanceNetInfo.ConnectionString)
		}
		if instanceNetInfo.DBInstanceNetType == "2" {
			d.Set("private_connection_port", instanceNetInfo.Port)
		}
	}

	d.Set("bandwidth", object.Bandwidth)
	d.Set("capacity", object.Capacity)
	d.Set("config", object.Config)
	d.Set("connection_domain", object.ConnectionDomain)
	d.Set("db_instance_name", object.InstanceName)
	d.Set("instance_name", object.InstanceName)
	d.Set("end_time", object.EndTime)
	d.Set("engine_version", object.EngineVersion)
	d.Set("instance_class", object.InstanceClass)
	d.Set("instance_release_protection", object.InstanceReleaseProtection)
	d.Set("instance_type", object.InstanceType)
	d.Set("maintain_end_time", object.MaintainEndTime)
	d.Set("maintain_start_time", object.MaintainStartTime)
	d.Set("node_type", object.NodeType)
	d.Set("payment_type", object.ChargeType)
	d.Set("instance_charge_type", object.ChargeType)

	d.Set("private_ip", object.PrivateIp)
	d.Set("qps", object.QPS)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("status", object.InstanceStatus)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		if !ignoredTags(t.Key, t.Value) {
			tags[t.Key] = t.Value
		}
	}
	d.Set("tags", tags)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("vpc_auth_mode", object.VpcAuthMode)
	d.Set("zone_id", object.ZoneId)
	d.Set("availability_zone", object.ZoneId)
	d.Set("secondary_zone_id", object.SecondaryZoneId)
	describeBackupPolicyObject, err := r_kvstoreService.DescribeBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_period", strings.Split(describeBackupPolicyObject.PreferredBackupPeriod, ","))
	d.Set("backup_time", describeBackupPolicyObject.PreferredBackupTime)

	if _, ok := d.GetOk("security_group_id"); ok {

		describeSecurityGroupConfigurationObject, err := r_kvstoreService.DescribeSecurityGroupConfiguration(d.Id())
		if err != nil {
			return WrapError(err)
		}
		sgs := make([]string, 0)
		for _, sg := range describeSecurityGroupConfigurationObject.EcsSecurityGroupRelation {
			sgs = append(sgs, sg.SecurityGroupId)
		}
		d.Set("security_group_id", strings.Join(sgs, ","))
	}
	if _, ok := d.GetOk("ssl_enable"); ok {

		describeInstanceSSLObject, err := r_kvstoreService.DescribeInstanceSSL(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("ssl_enable", describeInstanceSSLObject.SSLEnabled)
	}
	if _, ok := d.GetOk("security_ips"); ok {

		describeSecurityIpsObject, err := r_kvstoreService.DescribeSecurityIps(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("security_ip_group_attribute", describeSecurityIpsObject.SecurityIpGroupAttribute)
		d.Set("security_ip_group_name", describeSecurityIpsObject.SecurityIpGroupName)
		d.Set("security_ips", strings.Split(describeSecurityIpsObject.SecurityIpList, ","))
	}
	if object.ChargeType == string(PrePaid) {

		describeInstanceAutoRenewalAttributeObject, err := r_kvstoreService.DescribeInstanceAutoRenewalAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
		autoRenew, err := strconv.ParseBool(describeInstanceAutoRenewalAttributeObject.AutoRenew)
		if err != nil {
			// invalid request response
			return WrapError(err)
		}
		d.Set("auto_renew", autoRenew)
		d.Set("auto_renew_period", describeInstanceAutoRenewalAttributeObject.Duration)
	}
	//refresh parameters
	if err = refreshParameters(d, meta); err != nil {
		return WrapError(err)
	}
	return nil
}
func resourceAlicloudKvstoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	d.Partial(true)

	if d.HasChange("payment_type") {
		object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("payment_type").(string)
		if object.ChargeType != target {
			if target == "PrePaid" {
				request := r_kvstore.CreateTransformToPrePaidRequest()
				request.InstanceId = d.Id()
				if v, ok := d.GetOk("period"); ok {
					if v, err := strconv.Atoi(v.(string)); err == nil {
						request.Period = requests.NewInteger(v)
					} else {
						return WrapError(err)
					}
				}
				if v, ok := d.GetOk("auto_renew"); ok {
					request.AutoPay = requests.NewBoolean(v.(bool))
				}
				raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
					return r_kvstoreClient.TransformToPrePaid(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("payment_type")
		}
	}
	if d.HasChange("tags") {
		if err := r_kvstoreService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("config") {
		request := r_kvstore.CreateModifyInstanceConfigRequest()
		request.InstanceId = d.Id()
		respJson, err := convertMaptoJsonString(d.Get("config").(map[string]interface{}))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		request.Config = respJson
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceConfig(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("config")
	}
	if d.HasChange("vpc_auth_mode") && d.Get("vswitch_id") != "" {
		request := r_kvstore.CreateModifyInstanceVpcAuthModeRequest()
		request.InstanceId = d.Id()
		request.VpcAuthMode = d.Get("vpc_auth_mode").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceVpcAuthMode(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("vpc_auth_mode")
	}
	if d.HasChange("security_group_id") {
		request := r_kvstore.CreateModifySecurityGroupConfigurationRequest()
		request.DBInstanceId = d.Id()
		request.SecurityGroupId = d.Get("security_group_id").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifySecurityGroupConfiguration(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_id")
	}
	update := false
	request := r_kvstore.CreateModifyInstanceAutoRenewalAttributeRequest()
	request.DBInstanceId = d.Id()
	if !d.IsNewResource() && d.HasChange("auto_renew") {
		update = true
		request.AutoRenew = convertBoolToString(d.Get("auto_renew").(bool))
	}
	if !d.IsNewResource() && d.HasChange("auto_renew_period") {
		update = true
		request.Duration = convertIntergerToString(d.Get("auto_renew_period").(int))
	}
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceAutoRenewalAttribute(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}
	update = false
	modifyInstanceMaintainTimeReq := r_kvstore.CreateModifyInstanceMaintainTimeRequest()
	modifyInstanceMaintainTimeReq.InstanceId = d.Id()
	if d.HasChange("maintain_end_time") {
		update = true
	}
	modifyInstanceMaintainTimeReq.MaintainEndTime = d.Get("maintain_end_time").(string)
	if d.HasChange("maintain_start_time") {
		update = true
	}
	modifyInstanceMaintainTimeReq.MaintainStartTime = d.Get("maintain_start_time").(string)
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceMaintainTime(modifyInstanceMaintainTimeReq)
		})
		addDebug(modifyInstanceMaintainTimeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceMaintainTimeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	update = false
	modifyInstanceMajorVersionReq := r_kvstore.CreateModifyInstanceMajorVersionRequest()
	modifyInstanceMajorVersionReq.InstanceId = d.Id()
	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
	}
	modifyInstanceMajorVersionReq.MajorVersion = d.Get("engine_version").(string)
	modifyInstanceMajorVersionReq.EffectiveTime = "0"
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceMajorVersion(modifyInstanceMajorVersionReq)
		})
		addDebug(modifyInstanceMajorVersionReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceMajorVersionReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 300*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("engine_version")
	}
	if d.HasChange("ssl_enable") {
		request := r_kvstore.CreateModifyInstanceSSLRequest()
		request.InstanceId = d.Id()
		request.SSLEnabled = d.Get("ssl_enable").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceSSL(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil && !IsExpectedErrors(err, []string{"SSLDisableStateExistsFault"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_enable")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := r_kvstore.CreateModifyResourceGroupRequest()
		request.InstanceId = d.Id()
		request.RegionId = client.RegionId
		request.ResourceGroupId = d.Get("resource_group_id").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyResourceGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	migrateToOtherZoneReq := r_kvstore.CreateMigrateToOtherZoneRequest()
	migrateToOtherZoneReq.DBInstanceId = d.Id()
	if !d.IsNewResource() && (d.HasChange("zone_id") || d.HasChange("availability_zone")) {
		update = true
	}
	if v, ok := d.GetOk("zone_id"); ok {
		migrateToOtherZoneReq.ZoneId = v.(string)
	} else {
		migrateToOtherZoneReq.ZoneId = d.Get("availability_zone").(string)
	}
	migrateToOtherZoneReq.EffectiveTime = "0"
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
		migrateToOtherZoneReq.VSwitchId = d.Get("vswitch_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("secondary_zone_id") {
		update = true
		migrateToOtherZoneReq.SecondaryZoneId = d.Get("secondary_zone_id").(string)
	}
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.MigrateToOtherZone(migrateToOtherZoneReq)
		})
		addDebug(migrateToOtherZoneReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), migrateToOtherZoneReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 300*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("availability_zone")
		d.SetPartial("zone_id")
		d.SetPartial("vswitch_id")
	}
	update = false
	modifyBackupPolicyReq := r_kvstore.CreateModifyBackupPolicyRequest()
	modifyBackupPolicyReq.InstanceId = d.Id()
	if d.HasChange("backup_period") {
		update = true
	}
	modifyBackupPolicyReq.PreferredBackupPeriod = convertListToCommaSeparate(d.Get("backup_period").(*schema.Set).List())
	if d.HasChange("backup_time") {
		update = true
	}
	modifyBackupPolicyReq.PreferredBackupTime = d.Get("backup_time").(string)
	if update {
		if _, ok := d.GetOk("enable_backup_log"); ok {
			modifyBackupPolicyReq.EnableBackupLog = requests.NewInteger(d.Get("enable_backup_log").(int))
		}
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyBackupPolicy(modifyBackupPolicyReq)
		})
		addDebug(modifyBackupPolicyReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyBackupPolicyReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("backup_period")
		d.SetPartial("backup_time")
	}
	update = false
	modifyInstanceAttributeReq := r_kvstore.CreateModifyInstanceAttributeRequest()
	modifyInstanceAttributeReq.InstanceId = d.Id()
	if !d.IsNewResource() && d.HasChange("db_instance_name") {
		update = true
		modifyInstanceAttributeReq.InstanceName = d.Get("db_instance_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		modifyInstanceAttributeReq.InstanceName = d.Get("instance_name").(string)
	}
	if d.HasChange("instance_release_protection") {
		update = true
		modifyInstanceAttributeReq.InstanceReleaseProtection = requests.NewBoolean(d.Get("instance_release_protection").(bool))
	}
	if !d.IsNewResource() && (d.HasChange("password") || d.HasChange("kms_encrypted_password")) {
		update = true
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			modifyInstanceAttributeReq.NewPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			modifyInstanceAttributeReq.NewPassword = decryptResp
		}
	}
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyInstanceAttribute(modifyInstanceAttributeReq)
		})
		addDebug(modifyInstanceAttributeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceAttributeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_name")
		d.SetPartial("db_instance_name")
		d.SetPartial("instance_release_protection")
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("password")
	}
	update = false
	modifyDBInstanceConnectionStringReq := r_kvstore.CreateModifyDBInstanceConnectionStringRequest()
	modifyDBInstanceConnectionStringReq.DBInstanceId = d.Id()
	if d.HasChange("private_connection_prefix") {
		update = true
		modifyDBInstanceConnectionStringReq.NewConnectionString = d.Get("private_connection_prefix").(string)
	}
	if d.HasChange("private_connection_port") {
		update = true
		modifyDBInstanceConnectionStringReq.Port = d.Get("private_connection_port").(string)
	}
	modifyDBInstanceConnectionStringReq.IPType = "Private"
	if update {
		object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
		modifyDBInstanceConnectionStringReq.CurrentConnectionString = object.ConnectionDomain
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyDBInstanceConnectionString(modifyDBInstanceConnectionStringReq)
		})
		addDebug(modifyDBInstanceConnectionStringReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyDBInstanceConnectionStringReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("private_connection_prefix")
		d.SetPartial("private_connection_port")
	}
	update = false
	modifySecurityIpsReq := r_kvstore.CreateModifySecurityIpsRequest()
	modifySecurityIpsReq.InstanceId = d.Id()
	if d.HasChange("security_ips") {
		update = true
	}
	modifySecurityIpsReq.SecurityIps = convertListToCommaSeparate(d.Get("security_ips").(*schema.Set).List())
	if d.HasChange("security_ip_group_attribute") {
		update = true
		modifySecurityIpsReq.SecurityIpGroupAttribute = d.Get("security_ip_group_attribute").(string)
	}
	if d.HasChange("security_ip_group_name") {
		update = true
		modifySecurityIpsReq.SecurityIpGroupName = d.Get("security_ip_group_name").(string)
	}
	if update {
		if _, ok := d.GetOk("modify_mode"); ok {
			modifySecurityIpsReq.ModifyMode = convertModifyModeRequest(d.Get("modify_mode").(int))
		}
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifySecurityIps(modifySecurityIpsReq)
		})
		addDebug(modifySecurityIpsReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifySecurityIpsReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_ips")
		d.SetPartial("security_ip_group_attribute")
		d.SetPartial("security_ip_group_name")
	}
	update = false
	modifyInstanceSpecReq := r_kvstore.CreateModifyInstanceSpecRequest()
	modifyInstanceSpecReq.InstanceId = d.Id()
	modifyInstanceSpecReq.AutoPay = requests.NewBoolean(true)
	modifyInstanceSpecReq.EffectiveTime = "0"
	if !d.IsNewResource() && d.HasChange("engine_version") {

		modifyInstanceSpecReq.MajorVersion = d.Get("engine_version").(string)
	}
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
		modifyInstanceSpecReq.InstanceClass = d.Get("instance_class").(string)
	}
	if update {
		if _, ok := d.GetOk("business_info"); ok {
			modifyInstanceSpecReq.BusinessInfo = d.Get("business_info").(string)
		}
		if _, ok := d.GetOk("coupon_no"); ok {
			modifyInstanceSpecReq.CouponNo = d.Get("coupon_no").(string)
		}
		if _, ok := d.GetOkExists("force_upgrade"); ok {
			modifyInstanceSpecReq.ForceUpgrade = requests.NewBoolean(d.Get("force_upgrade").(bool))
		}
		if _, ok := d.GetOk("order_type"); ok {
			modifyInstanceSpecReq.OrderType = d.Get("order_type").(string)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			modifyInstanceSpecReq.ClientToken = buildClientToken(modifyInstanceSpecReq.GetActionName())
			args := *modifyInstanceSpecReq
			raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
				return r_kvstoreClient.ModifyInstanceSpec(&args)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"MissingRedisUsedmemoryUnsupportPerfItem"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceSpecReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 360*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("engine_version")
		d.SetPartial("instance_class")
	}
	if d.HasChange("parameters") {
		request := r_kvstore.CreateModifyInstanceConfigRequest()
		request.InstanceId = d.Id()
		config := make(map[string]interface{})
		documented := d.Get("parameters").(*schema.Set).List()
		if len(documented) > 0 {
			for _, i := range documented {
				key := i.(map[string]interface{})["name"].(string)
				value := i.(map[string]interface{})["value"]
				config[key] = value
			}
			cfg, _ := convertMaptoJsonString(config)
			request.Config = cfg

			raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
				return r_kvstoreClient.ModifyInstanceConfig(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("parameters")
	}
	if d.HasChange("enable_public") {
		prefix := d.Get("connection_string_prefix").(string)
		port := fmt.Sprintf("%s", d.Get("port"))
		target := d.Get("enable_public").(bool)

		if target {
			request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
			request.InstanceId = d.Id()
			request.ConnectionStringPrefix = prefix
			request.Port = port

			raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.AllocateInstancePublicConnection(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if !target && d.Get("connection_string") != "" {
			request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
			request.InstanceId = d.Id()
			request.CurrentConnectionString = d.Get("connection_string").(string)

			raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.ReleaseInstancePublicConnection(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("enable_public")
	}
	d.Partial(false)
	return resourceAlicloudKvstoreInstanceRead(d, meta)
}
func resourceAlicloudKvstoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := r_kvstore.CreateDeleteInstanceRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("global_instance_id"); ok {
		request.GlobalInstanceId = v.(string)
	}
	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DeleteInstance(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertModifyModeRequest(input int) string {
	switch input {
	case 0:
		return "Cover"
	case 1:
		return "Append"
	case 2:
		return "Delete"
	}
	return ""
}

func refreshParameters(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}

	var param []map[string]interface{}
	_, ok := d.GetOk("parameters")
	if !ok {
		d.Set("parameters", param)
		return nil
	}

	object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(object.Config), &m)
	if err != nil {
		return WrapError(err)
	}

	for k, v := range m {
		parameter := map[string]interface{}{
			"name":  k,
			"value": v,
		}
		param = append(param, parameter)
	}
	d.Set("parameters", param)
	return nil
}
