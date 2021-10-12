package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDmsEnterpriseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseInstanceCreate,
		Read:   resourceAlicloudDmsEnterpriseInstanceRead,
		Update: resourceAlicloudDmsEnterpriseInstanceUpdate,
		Delete: resourceAlicloudDmsEnterpriseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_link_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"database_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dba_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dba_nick_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dba_uid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ddl_online": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"ecs_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ecs_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_alias"},
			},
			"instance_alias": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_alias' has been deprecated from version 1.100.0. Use 'instance_name' instead.",
				ConflictsWith: []string{"instance_name"},
			},
			"instance_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"query_timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"safe_rule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"safe_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_test": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'state' has been deprecated from version 1.100.0. Use 'status' instead.",
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"use_dsql": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDmsEnterpriseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RegisterInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("data_link_name"); ok {
		request["DataLinkName"] = v
	}

	request["DatabasePassword"] = d.Get("database_password")
	request["DatabaseUser"] = d.Get("database_user")
	request["DbaUid"] = d.Get("dba_uid")
	if v, ok := d.GetOk("ddl_online"); ok {
		request["DdlOnline"] = v
	}

	if v, ok := d.GetOk("ecs_instance_id"); ok {
		request["EcsInstanceId"] = v
	}

	if v, ok := d.GetOk("ecs_region"); ok {
		request["EcsRegion"] = v
	}

	request["EnvType"] = d.Get("env_type")
	request["ExportTimeout"] = d.Get("export_timeout")
	request["Host"] = d.Get("host")
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceAlias"] = v
	} else if v, ok := d.GetOk("instance_alias"); ok {
		request["InstanceAlias"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "instance_alias" or "instance_name" must be set one!`))
	}

	request["InstanceSource"] = d.Get("instance_source")
	request["InstanceType"] = d.Get("instance_type")
	request["NetworkType"] = d.Get("network_type")
	request["Port"] = d.Get("port")
	request["QueryTimeout"] = d.Get("query_timeout")
	request["SafeRule"] = d.Get("safe_rule")
	if v, ok := d.GetOk("sid"); ok {
		request["Sid"] = v
	}

	if v, ok := d.GetOkExists("skip_test"); ok {
		request["SkipTest"] = v
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}

	if v, ok := d.GetOk("use_dsql"); ok {
		request["UseDsql"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	wait := incrementalWait(3*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"RegisterInstanceFailure"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Host"], ":", request["Port"]))

	return resourceAlicloudDmsEnterpriseInstanceUpdate(d, meta)
}
func resourceAlicloudDmsEnterpriseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	object, err := dms_enterpriseService.DescribeDmsEnterpriseInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_instance dms_enterpriseService.DescribeDmsEnterpriseInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host", parts[0])
	d.Set("port", formatInt(parts[1]))
	d.Set("data_link_name", object["DataLinkName"])
	d.Set("database_user", object["DatabaseUser"])
	d.Set("dba_id", object["DbaId"])
	d.Set("ddl_online", formatInt(object["DdlOnline"]))
	d.Set("ecs_instance_id", object["EcsInstanceId"])
	d.Set("ecs_region", object["EcsRegion"])
	d.Set("env_type", object["EnvType"])
	d.Set("export_timeout", formatInt(object["ExportTimeout"]))
	d.Set("instance_id", object["InstanceId"])
	d.Set("instance_name", object["InstanceAlias"])
	d.Set("instance_alias", object["InstanceAlias"])
	d.Set("instance_source", object["InstanceSource"])
	d.Set("instance_type", object["InstanceType"])
	d.Set("query_timeout", formatInt(object["QueryTimeout"]))
	d.Set("safe_rule_id", object["SafeRuleId"])
	d.Set("sid", object["Sid"])
	d.Set("status", object["State"])
	d.Set("state", object["State"])
	d.Set("use_dsql", formatInt(object["UseDsql"]))
	d.Set("vpc_id", object["VpcId"])
	d.Set("dba_nick_name", object["DbaNickName"])

	return nil
}
func resourceAlicloudDmsEnterpriseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"Host": parts[0],
		"Port": parts[1],
	}
	if !d.IsNewResource() && d.HasChange("database_password") {
		update = true
	}
	request["DatabasePassword"] = d.Get("database_password")
	if !d.IsNewResource() && d.HasChange("database_user") {
		update = true
	}
	request["DatabaseUser"] = d.Get("database_user")
	if d.HasChange("dba_id") {
		update = true
	}
	request["DbaId"] = d.Get("dba_id")
	if !d.IsNewResource() && d.HasChange("env_type") {
		update = true
	}
	request["EnvType"] = d.Get("env_type")
	if !d.IsNewResource() && d.HasChange("export_timeout") {
		update = true
	}
	request["ExportTimeout"] = d.Get("export_timeout")
	if d.HasChange("instance_id") {
		update = true
	}
	request["InstanceId"] = d.Get("instance_id")
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceAlias"] = d.Get("instance_name")
	}
	if !d.IsNewResource() && d.HasChange("instance_alias") {
		update = true
		request["InstanceAlias"] = d.Get("instance_alias")
	}
	if request["InstanceAlias"] == nil {
		request["InstanceAlias"] = d.Get("instance_name")
	}
	if !d.IsNewResource() && d.HasChange("instance_source") {
		update = true
	}
	request["InstanceSource"] = d.Get("instance_source")
	if !d.IsNewResource() && d.HasChange("instance_type") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if !d.IsNewResource() && d.HasChange("query_timeout") {
		update = true
	}
	request["QueryTimeout"] = d.Get("query_timeout")
	if d.HasChange("safe_rule_id") {
		update = true
	}
	request["SafeRuleId"] = d.Get("safe_rule_id")
	if !d.IsNewResource() && d.HasChange("data_link_name") {
		update = true
		request["DataLinkName"] = d.Get("data_link_name")
	}
	if !d.IsNewResource() && d.HasChange("ddl_online") {
		update = true
		request["DdlOnline"] = d.Get("ddl_online")
	}
	if !d.IsNewResource() && d.HasChange("ecs_instance_id") {
		update = true
		request["EcsInstanceId"] = d.Get("ecs_instance_id")
	}
	if !d.IsNewResource() && d.HasChange("ecs_region") {
		update = true
		request["EcsRegion"] = d.Get("ecs_region")
	}
	if !d.IsNewResource() && d.HasChange("sid") {
		update = true
		request["Sid"] = d.Get("sid")
	}
	if !d.IsNewResource() && d.HasChange("use_dsql") {
		update = true
		request["UseDsql"] = d.Get("use_dsql")
	}
	if !d.IsNewResource() && d.HasChange("vpc_id") {
		update = true
		request["VpcId"] = d.Get("vpc_id")
	}
	if update {
		if _, ok := d.GetOkExists("skip_test"); ok {
			request["SkipTest"] = d.Get("skip_test")
		}
		if _, ok := d.GetOk("tid"); ok {
			request["Tid"] = d.Get("tid")
		}
		action := "UpdateInstance"
		conn, err := client.NewDmsenterpriseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudDmsEnterpriseInstanceRead(d, meta)
}
func resourceAlicloudDmsEnterpriseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteInstance"
	var response map[string]interface{}
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Host": parts[0],
		"Port": parts[1],
	}

	if v, ok := d.GetOk("sid"); ok {
		request["Sid"] = v
	}
	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
