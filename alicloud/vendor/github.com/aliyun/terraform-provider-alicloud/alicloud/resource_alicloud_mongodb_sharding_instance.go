package alicloud

import (
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

func resourceAlicloudMongoDBShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBShardingInstanceCreate,
		Read:   resourceAlicloudMongoDBShardingInstanceRead,
		Update: resourceAlicloudMongoDBShardingInstanceUpdate,
		Delete: resourceAlicloudMongoDBShardingInstanceDelete,
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
				Computed:     true,
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
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" || d.Get("engine_version").(string) < "4.0"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled"}, false),
				Optional:     true,
				ForceNew:     true,
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
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"readonly_replicas": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 5),
							Computed:     true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},

			"mongo_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"tags": tagsSchema(),
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"config_server_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildMongoDBShardingCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateShardingDBInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := dds.CreateCreateShardingDBInstanceRequest()

	request.RegionId = string(client.Region)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
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

	shardList, ok := d.GetOk("shard_list")
	if ok {
		replicaSets := []dds.CreateShardingDBInstanceReplicaSet{}
		for _, rew := range shardList.([]interface{}) {
			item := rew.(map[string]interface{})
			readonlyReplicas := 0
			if item["readonly_replicas"] != nil {
				readonlyReplicas = item["readonly_replicas"].(int)
			}
			class := item["node_class"].(string)
			nodeStorage := item["node_storage"].(int)
			replicaSets = append(replicaSets, dds.CreateShardingDBInstanceReplicaSet{strconv.Itoa(readonlyReplicas), strconv.Itoa(nodeStorage), class})
		}
		request.ReplicaSet = &replicaSets
	}

	mongoList, ok := d.GetOk("mongo_list")
	if ok {
		mongos := []dds.CreateShardingDBInstanceMongos{}
		for _, rew := range mongoList.([]interface{}) {
			item := rew.(map[string]interface{})
			class := item["node_class"].(string)
			mongos = append(mongos, dds.CreateShardingDBInstanceMongos{class})
		}
		request.Mongos = &mongos
	}

	request.ConfigServer = &[]dds.CreateShardingDBInstanceConfigServer{{"20", "dds.cs.mid"}}

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
				return nil, WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in multi the zone " + request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch " + vsw.VSwitchId + " isn't in the zone " + request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if ok && PayType(request.ChargeType) == PrePaid {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List()), COMMA_SEPARATED)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request.AutoRenew = strconv.FormatBool(v.(bool))
	}

	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceAlicloudMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBShardingCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateShardingDBInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*dds.CreateShardingDBInstanceResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(response.DBInstanceId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBShardingInstanceUpdate(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_time", backupPolicy.PreferredBackupTime)
	d.Set("backup_period", strings.Split(backupPolicy.PreferredBackupPeriod, ","))
	d.Set("retention_period", backupPolicy.BackupRetentionPeriod)

	d.Set("name", instance.DBInstanceDescription)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("storage_engine", instance.StorageEngine)
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

	mongosList := []map[string]interface{}{}
	for _, item := range instance.MongosList.MongosAttribute {
		mongo := map[string]interface{}{
			"node_class":     item.NodeClass,
			"node_id":        item.NodeId,
			"port":           item.Port,
			"connect_string": item.ConnectSting,
		}
		mongosList = append(mongosList, mongo)
	}
	err = d.Set("mongo_list", mongosList)
	if err != nil {
		return WrapError(err)
	}

	shardList := []map[string]interface{}{}
	for _, item := range instance.ShardList.ShardAttribute {
		shard := map[string]interface{}{
			"node_id":           item.NodeId,
			"node_storage":      item.NodeStorage,
			"node_class":        item.NodeClass,
			"readonly_replicas": item.ReadonlyReplicas,
		}
		shardList = append(shardList, shard)
	}
	err = d.Set("shard_list", shardList)
	if err != nil {
		return WrapError(err)
	}
	tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo.TDEStatus)

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

	d.Set("tags", ddsService.tagsInAttributeToMap(instance.Tags.Tag))

	configServerSets := make([]map[string]interface{}, 0)
	for _, v := range instance.ConfigserverList.ConfigserverAttribute {
		configServerSets = append(configServerSets, map[string]interface{}{
			"max_iops":         v.MaxIOPS,
			"connect_string":   v.ConnectString,
			"node_class":       v.NodeClass,
			"max_connections":  v.MaxConnections,
			"port":             v.Port,
			"node_description": v.NodeDescription,
			"node_id":          v.NodeId,
			"node_storage":     v.NodeStorage,
		})
	}
	err = d.Set("config_server_list", configServerSets)

	return nil
}

func resourceAlicloudMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	d.Partial(true)

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
		return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
	}

	if d.HasChange("shard_list") {
		state, diff := d.GetChange("shard_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeShard, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("shard_list")
	}

	if d.HasChange("mongo_list") {
		state, diff := d.GetChange("mongo_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d, MongoDBShardingNodeMongos, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("mongo_list")
	}

	if d.HasChange("name") {
		request := dds.CreateModifyDBInstanceDescriptionRequest()
		request.RegionId = client.RegionId
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
		d.SetPartial("account_password")
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
	d.Partial(false)
	return resourceAlicloudMongoDBShardingInstanceRead(d, meta)
}

func resourceAlicloudMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		log.Printf("[WARN] Cannot destroy resourceAlicloudMongoDBShardingInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDeleteDBInstanceRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
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
	return nil
}
