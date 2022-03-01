package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMongodbServerlessInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbServerlessInstanceCreate,
		Read:   resourceAlicloudMongodbServerlessInstanceRead,
		Update: resourceAlicloudMongodbServerlessInstanceUpdate,
		Delete: resourceAlicloudMongodbServerlessInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z!#$%^&*()_+-=]{8,32}`), "account_password must consist of uppercase letters, lowercase letters, numbers, and special characters"),
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"capacity_unit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 8000),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"MongoDB"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"4.2"}, false),
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
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"period_price_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Day", "Month"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_ip_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_ip_group_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ip_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMongodbServerlessInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServerlessDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["DBInstanceStorage"] = d.Get("db_instance_storage")
	request["CapacityUnit"] = d.Get("capacity_unit")
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}
	request["EngineVersion"] = d.Get("engine_version")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_price_type"); ok {
		request["PeriodPriceType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("storage_engine"); ok {
		request["StorageEngine"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	request["VpcId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateServerlessDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_serverless_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	MongoDBService := MongoDBService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, MongoDBService.MongodbServerlessInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMongodbServerlessInstanceUpdate(d, meta)
}
func resourceAlicloudMongodbServerlessInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongoDBService := MongoDBService{client}
	object, err := mongoDBService.DescribeMongodbServerlessInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_serverless_instance MongoDBService.DescribeMongodbServerlessInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("capacity_unit", formatInt(object["CapacityUnit"]))
	d.Set("db_instance_description", object["DBInstanceDescription"])
	if v, ok := object["DBInstanceStorage"]; ok && fmt.Sprint(v) != "0" {
		d.Set("db_instance_storage", formatInt(v))
	}
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("maintain_end_time", object["MaintainEndTime"])
	d.Set("maintain_start_time", object["MaintainStartTime"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["DBInstanceStatus"])
	d.Set("storage_engine", convertMongodbServerlessInstanceStorageEngineResponse(object["StorageEngine"].(string)))

	listTagResourcesObject, err := mongoDBService.ListTagResources(d.Id(), "INSTANCE")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("vpc_id", object["VPCId"])
	securityIpGroupsObject, err := mongoDBService.DescribeSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if securityIpGroupsMap, ok := securityIpGroupsObject["SecurityIpGroups"].(map[string]interface{}); ok && securityIpGroupsMap != nil {
		if securityIpGroup, ok := securityIpGroupsMap["SecurityIpGroup"]; ok && securityIpGroup != nil {
			securityIpGroupMaps := make([]map[string]interface{}, 0)
			for _, iPArrayListItem := range securityIpGroup.([]interface{}) {
				if v, ok := iPArrayListItem.(map[string]interface{}); ok {
					if v["SecurityIpGroupName"].(string) == "default" {
						continue
					}
					iPArrayListItemMap := make(map[string]interface{})
					iPArrayListItemMap["security_ip_group_attribute"] = v["SecurityIpGroupAttribute"]
					iPArrayListItemMap["security_ip_group_name"] = v["SecurityIpGroupName"]
					iPArrayListItemMap["security_ip_list"] = v["SecurityIpList"]
					securityIpGroupMaps = append(securityIpGroupMaps, iPArrayListItemMap)
				}
			}
			err = d.Set("security_ip_groups", securityIpGroupMaps)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
func resourceAlicloudMongodbServerlessInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	var response map[string]interface{}
	d.Partial(true)
	if d.HasChange("tags") {
		if err := MongoDBService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}
		if v, ok := d.GetOk("db_instance_description"); ok {
			request["DBInstanceDescription"] = v
		}
		action := "ModifyDBInstanceDescription"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("db_instance_description")
	}
	if d.HasChange("maintain_end_time") || d.HasChange("maintain_start_time") {
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}
		request["MaintainStartTime"] = d.Get("maintain_start_time")
		request["MaintainEndTime"] = d.Get("maintain_end_time")
		action := "ModifyDBInstanceMaintainTime"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	if d.HasChange("security_ip_groups") {
		oraw, nraw := d.GetChange("security_ip_groups")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			removeSecurityIpsReq := map[string]interface{}{
				"DBInstanceId": d.Id(),
				"ModifyMode":   "Delete",
			}
			for _, whiteList := range remove {
				whiteListArg := whiteList.(map[string]interface{})
				removeSecurityIpsReq["SecurityIpGroupAttribute"] = whiteListArg["security_ip_group_attribute"]
				removeSecurityIpsReq["SecurityIpGroupName"] = whiteListArg["security_ip_group_name"]
				removeSecurityIpsReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifySecurityIps"
				conn, err := client.NewDdsClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, removeSecurityIpsReq, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removeSecurityIpsReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		if len(create) > 0 {
			createSecurityIpsReq := map[string]interface{}{
				"DBInstanceId": d.Id(),
				"ModifyMode":   "Append",
			}
			for _, whiteList := range create {
				whiteListArg := whiteList.(map[string]interface{})
				createSecurityIpsReq["SecurityIpGroupAttribute"] = whiteListArg["security_ip_group_attribute"]
				createSecurityIpsReq["SecurityIpGroupName"] = whiteListArg["security_ip_group_name"]
				createSecurityIpsReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifySecurityIps"
				conn, err := client.NewDdsClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, createSecurityIpsReq, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, createSecurityIpsReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		d.SetPartial("security_ip_groups")
	}
	if !d.IsNewResource() && (d.HasChange("db_instance_storage") || d.HasChange("capacity_unit")) {
		request := map[string]interface{}{
			"DBInstanceId": d.Id(),
		}
		if v, ok := d.GetOk("capacity_unit"); ok {
			request["DBInstanceClass"] = v
		}
		if v, ok := d.GetOk("db_instance_storage"); ok {
			request["DBInstanceStorage"] = v
		}
		action := "ModifyDBInstanceSpec"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, MongoDBService.MongodbServerlessInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("capacity_unit")
		d.SetPartial("db_instance_storage")
	}
	d.Partial(false)

	return resourceAlicloudMongodbServerlessInstanceRead(d, meta)
}
func resourceAlicloudMongodbServerlessInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource Alicloud Resource Mongodb Serverless Subscription Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertMongodbServerlessInstancePayTypeResponse(source string) string {
	switch source {
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertMongodbServerlessInstanceStorageEngineResponse(source string) string {
	switch source {
	case "wiredTiger":
		return "WiredTiger"
	}
	return source
}
