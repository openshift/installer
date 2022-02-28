package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMongoDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBInstanceCreate,
		Read:   resourceAlicloudMongoDBInstanceRead,
		Update: resourceAlicloudMongoDBInstanceUpdate,
		Delete: resourceAlicloudMongoDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(10, 2000),
				Required:     true,
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 5, 7}),
				Optional:     true,
				Computed:     true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger", "RocksDB"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
			},
			//Computed
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" || d.Get("engine_version").(string) < "4.0"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled"}, false),
				Optional:     true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"replica_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica_set_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
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
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func buildMongoDBCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	request := dds.CreateCreateDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
	request.DBInstanceStorage = requests.NewInteger(d.Get("db_instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("db_instance_class").(string))
	request.DBInstanceDescription = d.Get("name").(string)

	request.AccountPassword = d.Get("account_password").(string)
	if request.AccountPassword == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, WrapError(err)
			}
			request.AccountPassword = decryptResp
		}
	}

	request.ZoneId = d.Get("zone_id").(string)
	request.StorageEngine = d.Get("storage_engine").(string)

	if replication_factor, ok := d.GetOk("replication_factor"); ok {
		request.ReplicationFactor = strconv.Itoa(replication_factor.(int))
	}

	request.NetworkType = string(Classic)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(fmt.Errorf("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if PayType(request.ChargeType) == PrePaid && ok {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		request.AutoRenew = strconv.FormatBool(v.(bool))
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceAlicloudMongoDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*dds.CreateDBInstanceResponse)
	d.SetId(response.DBInstanceId)
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBInstanceUpdate(d, meta)
}

func resourceAlicloudMongoDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	instance, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if len(instance.ReplicaSets.ReplicaSet) > 0 {
		replicaSets := make([]map[string]interface{}, 0)
		for _, v := range instance.ReplicaSets.ReplicaSet {
			replicaSets = append(replicaSets, map[string]interface{}{
				"vswitch_id":            v.VSwitchId,
				"connection_port":       v.ConnectionPort,
				"replica_set_role":      v.ReplicaSetRole,
				"connection_domain":     v.ConnectionDomain,
				"vpc_cloud_instance_id": v.VPCCloudInstanceId,
				"network_type":          v.NetworkType,
				"vpc_id":                v.VPCId,
			})
		}
		d.Set("replica_sets", replicaSets)
	}
	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_time", backupPolicy.PreferredBackupTime)
	d.Set("backup_period", strings.Split(backupPolicy.PreferredBackupPeriod, ","))
	d.Set("retention_period", backupPolicy.BackupRetentionPeriod)

	ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", ips)

	groupIp, err := ddsService.DescribeMongoDBSecurityGroupId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(groupIp.Items.RdsEcsSecurityGroupRel) > 0 {
		d.Set("security_group_id", groupIp.Items.RdsEcsSecurityGroupRel[0].SecurityGroupId)
	}

	d.Set("name", instance.DBInstanceDescription)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("db_instance_class", instance.DBInstanceClass)
	d.Set("db_instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.ChargeType)
	if instance.ChargeType == "PrePaid" {
		period, err := computePeriodByUnit(instance.CreationTime, instance.ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("storage_engine", instance.StorageEngine)
	d.Set("maintain_start_time", instance.MaintainStartTime)
	d.Set("maintain_end_time", instance.MaintainEndTime)
	d.Set("replica_set_name", instance.ReplicaSetName)

	sslAction, err := ddsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction.SSLStatus)
	if v, err := strconv.Atoi(instance.ReplicationFactor); err != nil {
		log.Println(WrapError(err))
	} else {
		d.Set("replication_factor", v)
	}
	if instance.ReplicationFactor != "" && instance.ReplicationFactor != "1" {
		tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("tde_Status", tdeInfo.TDEStatus)
	}

	d.Set("tags", ddsService.tagsInAttributeToMap(instance.Tags.Tag))
	return nil
}

func resourceAlicloudMongoDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	d.Partial(true)

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		prePaidRequest := dds.CreateTransformToPrePaidRequest()
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.AutoPay = requests.NewBoolean(true)
		prePaidRequest.Period = requests.NewInteger(d.Get("period").(int))
		if v, ok := d.GetOk("auto_renew"); ok {
			prePaidRequest.AutoRenew = strconv.FormatBool(v.(bool))
		}
		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.TransformToPrePaid(prePaidRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("period")
	}

	if d.HasChange("backup_time") || d.HasChange("backup_period") {
		if err := ddsService.MotifyMongoDBBackupPolicy(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("backup_time")
		d.SetPartial("backup_period")
	}

	if d.HasChange("tde_status") {
		request := dds.CreateModifyDBInstanceTDERequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.TDEStatus = d.Get("tde_status").(string)
		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("tde_status")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		request := dds.CreateModifyDBInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	if d.HasChange("security_group_id") {
		request := dds.CreateModifySecurityGroupConfigurationRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()
		request.SecurityGroupId = d.Get("security_group_id").(string)

		wait := incrementalWait(2*time.Second, 3*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
				return client.ModifySecurityGroupConfiguration(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{"InstanceStatusInvalid"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("security_group_id")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudMongoDBInstanceRead(d, meta)
	}

	if d.HasChange("name") {
		request := dds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("name").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceDescription(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}

	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := ddsService.ModifyMongoDBSecurityIps(d, ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		var accountPassword string
		if accountPassword = d.Get("account_password").(string); accountPassword != "" {
			d.SetPartial("account_password")
		} else if kmsPassword := d.Get("kms_encrypted_password").(string); kmsPassword != "" {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			accountPassword = decryptResp
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
		}

		err := ddsService.ResetAccountPassword(d, accountPassword)
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("ssl_action") {
		request := dds.CreateModifyDBInstanceSSLRequest()
		request.DBInstanceId = d.Id()
		request.RegionId = client.RegionId
		request.SSLAction = d.Get("ssl_action").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSSL(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("ssl_action")
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("db_instance_storage") ||
		d.HasChange("db_instance_class") ||
		d.HasChange("replication_factor") {

		request := dds.CreateModifyDBInstanceSpecRequest()
		request.DBInstanceId = d.Id()

		if d.Get("instance_charge_type").(string) == "PrePaid" {
			if v, ok := d.GetOk("order_type"); ok {
				request.OrderType = v.(string)
			}
		}
		request.DBInstanceClass = d.Get("db_instance_class").(string)
		request.DBInstanceStorage = strconv.Itoa(d.Get("db_instance_storage").(int))
		request.ReplicationFactor = strconv.Itoa(d.Get("replication_factor").(int))

		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSpec(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("db_instance_class")
		d.SetPartial("db_instance_storage")
		d.SetPartial("replication_factor")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAlicloudMongoDBInstanceRead(d, meta)
}

func resourceAlicloudMongoDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// Pre paid instance can not be release.
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		log.Printf("[WARN] Cannot destroy resourceAlicloudMongoDBInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
