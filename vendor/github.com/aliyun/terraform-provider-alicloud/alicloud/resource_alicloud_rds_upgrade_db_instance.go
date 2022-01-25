package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRdsUpgradeDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsUpgradeDbInstanceCreate,
		Read:   resourceAlicloudRdsUpgradeDbInstanceRead,
		Update: resourceAlicloudRdsUpgradeDbInstanceUpdate,
		Delete: resourceAlicloudRdsUpgradeDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cert", "perfer", "verify-ca", "verify-full"}, false),
				Computed:     true,
			},
			"auto_upgrade_minor_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Auto", "Manual"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine").(string) != "MySQL"
				},
			},
			"ca_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"aliyun", "custom"}, false),
				Computed:     true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ca_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert_revocation_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_crl_enabled": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"db_instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_essd2", "cloud_essd3", "cloud_ssd", "local_ssd"}, false),
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Auto", "Down", "TempUpgrade", "Up"}, false),
			},
			"effective_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ha_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"RPO", "RTO"}, false),
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "VPC"}, false),
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"released_keep_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replication_acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cert", "perfer", "verify-ca", "verify-full"}, false),
				Computed:     true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"server_cert": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_biz": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Optional:     true,
				Computed:     true,
			},
			"switch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Async", "Semi-sync", "Sync"}, false),
			},
			"tde_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id_slave_1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"switch_time_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Immediate", "MaintainTime"}, false),
			},
			"switch_over": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"collect_stat_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Before", "After"}, false),
			},
			"target_major_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsUpgradeDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	conn, err := client.NewRdsClient()
	action := "UpgradeDBInstanceMajorVersionPrecheck"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Get("source_db_instance_id"),
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("target_major_version"); ok && v.(string) != "" {
		request["TargetMajorVersion"] = v
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, rdsService.RdsUpgradeMajorVersionRefreshFunc(d.Get("source_db_instance_id").(string), formatInt(response["TaskId"]), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	action = "UpgradeDBInstanceMajorVersion"
	request = make(map[string]interface{})
	request["SourceIp"] = client.SourceIp

	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("db_instance_class"); ok {
		request["DBInstanceClass"] = v
	}
	if v, ok := d.GetOk("db_instance_storage"); ok {
		request["DBInstanceStorage"] = v
	}
	request["PayType"] = convertRdsInstancePaymentTypeRequest(d.Get("payment_type"))
	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}
	if v, ok := d.GetOk("switch_time_mode"); ok {
		request["SwitchTimeMode"] = v
	}
	if v, ok := d.GetOk("switch_over"); ok {
		request["SwitchOver"] = v
	}
	if v, ok := d.GetOk("collect_stat_mode"); ok {
		request["CollectStatMode"] = v
	}
	if v, ok := d.GetOk("target_major_version"); ok {
		request["TargetMajorVersion"] = v
	}
	request["DBInstanceId"] = d.Get("source_db_instance_id")
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	request["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("zone_id_slave_1"); ok {
		request["ZoneIdSlave1"] = v
	}
	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_upgrade_db_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	// wait instance status change from Creating to running
	stateConf = BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRdsUpgradeDbInstanceUpdate(d, meta)
}
func resourceAlicloudRdsUpgradeDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeRdsCloneDbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_upgrade_db_instance rdsService.DescribeRdsCloneDbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auto_upgrade_minor_version", object["AutoUpgradeMinorVersion"])
	d.Set("db_instance_class", object["DBInstanceClass"])
	d.Set("db_instance_description", object["DBInstanceDescription"])
	if v, ok := object["DBInstanceStorage"]; ok && fmt.Sprint(v) != "0" {
		d.Set("db_instance_storage", formatInt(v))
	}
	d.Set("db_instance_storage_type", object["DBInstanceStorageType"])
	d.Set("dedicated_host_group_id", object["DedicatedHostGroupId"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("instance_network_type", object["InstanceNetworkType"])
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	if len(object["SlaveZones"].(map[string]interface{})["SlaveZone"].([]interface{})) > 0 {
		d.Set("zone_id_slave_1", object["SlaveZones"].(map[string]interface{})["SlaveZone"].([]interface{})[0].(map[string]interface{})["ZoneId"])
	}
	d.Set("payment_type", convertRdsInstancePaymentTypeResponse(object["PayType"]))
	d.Set("port", object["Port"])
	d.Set("connection_string", object["ConnectionString"])

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return WrapError(err)
	}
	describeDBInstanceHAConfigObject, err := rdsService.DescribeDBInstanceHAConfig(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("sync_mode", describeDBInstanceHAConfigObject["SyncMode"])
	d.Set("ha_mode", describeDBInstanceHAConfigObject["HAMode"])
	dbInstanceIpArrayName := "default"
	describeDBInstanceIPArrayListObject, err := rdsService.GetSecurityIps(d.Id(), dbInstanceIpArrayName)
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", describeDBInstanceIPArrayListObject)
	describeDBInstanceNetInfoObject, err := rdsService.DescribeDBInstanceNetInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("private_ip_address", describeDBInstanceNetInfoObject[0].(map[string]interface{})["IPAddress"])

	describeDBInstanceSSLObject, err := rdsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("acl", describeDBInstanceSSLObject["ACL"])
	d.Set("ca_type", describeDBInstanceSSLObject["CAType"])
	d.Set("client_ca_cert", describeDBInstanceSSLObject["ClientCACert"])
	d.Set("client_cert_revocation_list", describeDBInstanceSSLObject["ClientCertRevocationList"])
	d.Set("replication_acl", describeDBInstanceSSLObject["ReplicationACL"])
	d.Set("server_cert", describeDBInstanceSSLObject["ServerCert"])
	d.Set("server_key", describeDBInstanceSSLObject["ServerKey"])
	if v, ok := describeDBInstanceSSLObject["SSLEnabled"]; ok && v.(string) != "" {
		sslEnabled := 0
		if v == "on" {
			sslEnabled = 1
		}
		d.Set("ssl_enabled'", sslEnabled)
	}
	return nil
}
func resourceAlicloudRdsUpgradeDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("auto_upgrade_minor_version") {
		update = true
	}
	if v, ok := d.GetOk("auto_upgrade_minor_version"); ok {
		request["AutoUpgradeMinorVersion"] = v
	}
	if update {
		action := "ModifyDBInstanceAutoUpgradeMinorVersion"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyDBInstanceAutoUpgradeMinorVersion")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_upgrade_minor_version")
	}
	update = false
	modifyDBInstanceDescriptionReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("db_instance_description") {
		update = true
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		modifyDBInstanceDescriptionReq["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceDescriptionReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceDescriptionReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_instance_description")
	}
	update = false
	modifyDBInstanceMaintainTimeReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("maintain_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_time"); ok {
		modifyDBInstanceMaintainTimeReq["MaintainTime"] = v
	}
	if update {
		action := "ModifyDBInstanceMaintainTime"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyDBInstanceMaintainTime")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceMaintainTimeReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceMaintainTimeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_time")
	}
	update = false
	modifyDBInstanceHAConfigReq := map[string]interface{}{
		"DbInstanceId": d.Id(),
	}
	if d.HasChange("sync_mode") {
		update = true
	}
	if d.HasChange("ha_mode") {
		update = true
	}
	if update {
		if v, ok := d.GetOk("sync_mode"); ok {
			modifyDBInstanceHAConfigReq["SyncMode"] = v
		}
		if v, ok := d.GetOk("ha_mode"); ok {
			modifyDBInstanceHAConfigReq["HAMode"] = v
		}
		action := "ModifyDBInstanceHAConfig"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceHAConfigReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceHAConfigReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("sync_mode")
	}
	update = false
	switchDBInstanceVpcReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("private_ip_address") {
		update = true
		if v, ok := d.GetOk("private_ip_address"); ok {
			switchDBInstanceVpcReq["PrivateIpAddress"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("vpc_id") {
		update = true
		if v, ok := d.GetOk("vpc_id"); ok {
			switchDBInstanceVpcReq["VPCId"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
		if v, ok := d.GetOk("vswitch_id"); ok {
			switchDBInstanceVpcReq["VSwitchId"] = v
		}
	}
	if update {
		action := "SwitchDBInstanceVpc"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, switchDBInstanceVpcReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, switchDBInstanceVpcReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("private_ip_address")
		d.SetPartial("vpc_id")
		d.SetPartial("vswitch_id")
	}
	update = false
	modifyDBInstanceConnectionStringReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("port") {
		update = true
	}
	if d.HasChange("connection_string_prefix") {
		update = true
	}
	if v, ok := d.GetOk("connection_string"); ok {
		modifyDBInstanceConnectionStringReq["CurrentConnectionString"] = v
	}
	if update {
		if v, ok := d.GetOk("connection_string_prefix"); ok {
			modifyDBInstanceConnectionStringReq["ConnectionStringPrefix"] = v
		}
		if v, ok := d.GetOk("port"); ok {
			modifyDBInstanceConnectionStringReq["Port"] = v
		}
		action := "ModifyDBInstanceConnectionString"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceConnectionStringReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceConnectionStringReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("connection_string")
	}
	update = false
	modifyDBInstanceTDEReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("encryption_key") {
		update = true
		if v, ok := d.GetOk("encryption_key"); ok {
			modifyDBInstanceTDEReq["EncryptionKey"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("tde_status"); ok {
			modifyDBInstanceTDEReq["TDEStatus"] = v
		}
		if v, ok := d.GetOk("certificate"); ok {
			modifyDBInstanceTDEReq["Certificate"] = v
		}
		if v, ok := d.GetOk("db_name"); ok {
			modifyDBInstanceTDEReq["DBName"] = v
		}
		if v, ok := d.GetOk("password"); ok {
			modifyDBInstanceTDEReq["PassWord"] = v
		}
		if v, ok := d.GetOk("private_key"); ok {
			modifyDBInstanceTDEReq["PrivateKey"] = v
		}
		if v, ok := d.GetOk("role_arn"); ok {
			modifyDBInstanceTDEReq["RoleArn"] = v
		}
		action := "ModifyDBInstanceTDE"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceTDEReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceTDEReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("encryption_key")
	}
	update = false
	modifySecurityIpsReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("security_ips") {
		update = true
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}
		modifySecurityIpsReq["SecurityIps"] = ipstr
	}
	if update {
		action := "ModifySecurityIps"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifySecurityIpsReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySecurityIpsReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("security_ips")
	}
	update = false
	modifyDBInstanceSSLReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if v, ok := d.GetOk("connection_string"); ok {
		modifyDBInstanceSSLReq["ConnectionString"] = v
	}
	if d.HasChange("acl") {
		update = true
		if v, ok := d.GetOk("acl"); ok {
			modifyDBInstanceSSLReq["ACL"] = v
		}
	}
	if d.HasChange("ca_type") {
		update = true
		if v, ok := d.GetOk("ca_type"); ok {
			modifyDBInstanceSSLReq["CAType"] = v
		}
	}
	if d.HasChange("client_ca_cert") {
		update = true
		if v, ok := d.GetOk("client_ca_cert"); ok {
			modifyDBInstanceSSLReq["ClientCACert"] = v
		}
	}
	if d.HasChange("client_cert_revocation_list") {
		update = true
		if v, ok := d.GetOk("client_cert_revocation_list"); ok {
			modifyDBInstanceSSLReq["ClientCertRevocationList"] = v
		}
	}
	if d.HasChange("replication_acl") {
		update = true
		if v, ok := d.GetOk("replication_acl"); ok {
			modifyDBInstanceSSLReq["ReplicationACL"] = v
		}
	}
	if d.HasChange("server_cert") {
		update = true
		if v, ok := d.GetOk("server_cert"); ok {
			modifyDBInstanceSSLReq["ServerCert"] = v
		}
	}
	if d.HasChange("server_key") {
		update = true
		if v, ok := d.GetOk("server_key"); ok {
			modifyDBInstanceSSLReq["ServerKey"] = v
		}
	}
	if d.HasChange("client_ca_enabled") {
		update = true
		if v, ok := d.GetOk("client_ca_enabled"); ok {
			modifyDBInstanceSSLReq["ClientCAEnabled"] = v
		}
	}
	if d.HasChange("client_crl_enabled") {
		update = true
		if v, ok := d.GetOk("client_crl_enabled"); ok {
			modifyDBInstanceSSLReq["ClientCrlEnabled"] = v
		}
	}
	if d.HasChange("ssl_enabled") {
		update = true
		if v, ok := d.GetOk("ssl_enabled"); ok {
			modifyDBInstanceSSLReq["SSLEnabled"] = v
		}
	}
	if update {
		action := "ModifyDBInstanceSSL"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceSSLReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceSSLReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("connection_string")
		d.SetPartial("acl")
		d.SetPartial("ca_type")
		d.SetPartial("client_ca_cert")
		d.SetPartial("client_cert_revocation_list")
		d.SetPartial("replication_acl")
		d.SetPartial("server_cert")
		d.SetPartial("server_key")
		d.SetPartial("ssl_enabled")
	}
	update = false
	modifyDBInstanceSpecReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	if v, ok := d.GetOk("payment_type"); ok {
		modifyDBInstanceSpecReq["PayType"] = convertRdsInstancePaymentTypeRequest(v)
	}
	if !d.IsNewResource() && d.HasChange("db_instance_class") {
		update = true
		if v, ok := d.GetOk("db_instance_class"); ok {
			modifyDBInstanceSpecReq["DBInstanceClass"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("db_instance_storage") {
		update = true
		if v, ok := d.GetOk("db_instance_storage"); ok {
			modifyDBInstanceSpecReq["DBInstanceStorage"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("db_instance_storage_type") {
		update = true
		modifyDBInstanceSpecReq["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_group_id") {
		update = true
		if v, ok := d.GetOk("dedicated_host_group_id"); ok {
			modifyDBInstanceSpecReq["DedicatedHostGroupId"] = v
		}
	}
	if d.HasChange("engine_version") {
		update = true
		if v, ok := d.GetOk("engine_version"); ok {
			modifyDBInstanceSpecReq["EngineVersion"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("zone_id") {
		update = true
		if v, ok := d.GetOk("zone_id"); ok {
			modifyDBInstanceSpecReq["ZoneId"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("direction"); ok {
			modifyDBInstanceSpecReq["Direction"] = v
		}
		if v, ok := d.GetOk("effective_time"); ok {
			modifyDBInstanceSpecReq["EffectiveTime"] = v
		}
		if v, ok := d.GetOk("resource_group_id"); ok {
			modifyDBInstanceSpecReq["ResourceGroupId"] = v
		}
		if v, ok := d.GetOk("source_biz"); ok {
			modifyDBInstanceSpecReq["SourceBiz"] = v
		}
		if v, ok := d.GetOk("switch_time"); ok {
			modifyDBInstanceSpecReq["SwitchTime"] = v
		}
		action := "ModifyDBInstanceSpec"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, modifyDBInstanceSpecReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceSpecReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait instance status change from Creating to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_instance_class")
		d.SetPartial("db_instance_storage")
		d.SetPartial("db_instance_storage_type")
		d.SetPartial("dedicated_host_group_id")
		d.SetPartial("engine_version")
		d.SetPartial("zone_id")
	}
	d.Partial(false)
	return resourceAlicloudRdsUpgradeDbInstanceRead(d, meta)
}
func resourceAlicloudRdsUpgradeDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     client.RegionId,
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("released_keep_policy"); ok {
		request["ReleasedKeepPolicy"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
