package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogStoreCreate,
		Read:   resourceAlicloudLogStoreRead,
		Update: resourceAlicloudLogStoreUpdate,
		Delete: resourceAlicloudLogStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 3650),
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "" {
						return false
					}
					return true
				},
			},
			"shards": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"max_split_shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 64),
			},
			"append_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_web_tracking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encrypt_conf": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"encrypt_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "default",
							ValidateFunc: validation.StringInSlice([]string{"default", "m4"}, false),
						},
						"user_cmk_info": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmk_key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"arn": {
										Type:     schema.TypeString,
										Required: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logstore := &sls.LogStore{
		Name:          d.Get("name").(string),
		TTL:           d.Get("retention_period").(int),
		ShardCount:    d.Get("shard_count").(int),
		WebTracking:   d.Get("enable_web_tracking").(bool),
		AutoSplit:     d.Get("auto_split").(bool),
		MaxSplitShard: d.Get("max_split_shard_count").(int),
		AppendMeta:    d.Get("append_meta").(bool),
	}
	if encrypt := buildEncrypt(d); encrypt != nil {
		logstore.EncryptConf = encrypt
	}
	var requestinfo *sls.Client
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.CreateLogStoreV2(d.Get("project").(string), logstore)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateLogStoreV2", raw, requestinfo, map[string]interface{}{
				"project":  d.Get("project").(string),
				"logstore": logstore,
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "CreateLogStoreV2", AliyunLogGoSdkERROR)
	}
	// Wait for the store to be available
	time.Sleep(60 * time.Second)
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogStoreUpdate(d, meta)
}

func resourceAlicloudLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogStore(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("retention_period", object.TTL)
	d.Set("shard_count", object.ShardCount)
	var shards []*sls.Shard
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		shards, err = object.ListShards()
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListShards", shards)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "ListShards", AliyunLogGoSdkERROR)
	}
	var shardList []map[string]interface{}
	for _, s := range shards {
		mapping := map[string]interface{}{
			"id":        s.ShardID,
			"status":    s.Status,
			"begin_key": s.InclusiveBeginKey,
			"end_key":   s.ExclusiveBeginKey,
		}
		shardList = append(shardList, mapping)
	}
	d.Set("shards", shardList)
	d.Set("append_meta", object.AppendMeta)
	d.Set("auto_split", object.AutoSplit)
	d.Set("enable_web_tracking", object.WebTracking)
	d.Set("max_split_shard_count", object.MaxSplitShard)
	if encrypt := object.EncryptConf; encrypt != nil {
		encryptMap := map[string]interface{}{
			"enable":       encrypt.Enable,
			"encrypt_type": encrypt.EncryptType,
		}
		if userCmkInfo := encrypt.UserCmkInfo; userCmkInfo != nil {
			userCmkInfoMap := map[string]interface{}{
				"cmk_key_id": userCmkInfo.CmkKeyId,
				"arn":        userCmkInfo.Arn,
				"region_id":  userCmkInfo.RegionId,
			}
			encryptMap["user_cmk_info"] = []map[string]interface{}{userCmkInfoMap}
		}
		if err := d.Set("encrypt_conf", []map[string]interface{}{encryptMap}); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	if d.IsNewResource() {
		return resourceAlicloudLogStoreRead(d, meta)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	if d.HasChange("retention_period") {
		update = true
		d.SetPartial("retention_period")
	}
	if d.HasChange("max_split_shard_count") {
		update = true
		d.SetPartial("max_split_shard_count")
	}
	if d.HasChange("enable_web_tracking") {
		update = true
		d.SetPartial("enable_web_tracking")
	}
	if d.HasChange("append_meta") {
		update = true
		d.SetPartial("append_meta")
	}
	if d.HasChange("auto_split") {
		update = true
		d.SetPartial("auto_split")
	}

	if update {
		store, err := logService.DescribeLogStore(d.Id())
		if err != nil {
			return WrapError(err)
		}
		store.MaxSplitShard = d.Get("max_split_shard_count").(int)
		store.TTL = d.Get("retention_period").(int)
		store.WebTracking = d.Get("enable_web_tracking").(bool)
		store.AppendMeta = d.Get("append_meta").(bool)
		store.AutoSplit = d.Get("auto_split").(bool)
		var requestInfo *sls.Client
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateLogStoreV2(parts[0], store)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogStoreV2", AliyunLogGoSdkERROR)
		}
		if debugOn() {
			addDebug("UpdateLogStoreV2", raw, requestInfo, map[string]interface{}{
				"project":  parts[0],
				"logstore": store,
			})
		}
	}
	d.Partial(false)

	return resourceAlicloudLogStoreRead(d, meta)
}

func resourceAlicloudLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var raw interface{}
	var requestInfo *sls.Client
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteLogStore(parts[0], parts[1])
		})
		if err != nil {
			if code, ok := err.(*sls.Error); ok {
				if "LogStoreNotExist" == code.Code {
					return nil
				}
			}
			if IsExpectedErrors(err, []string{"ProjectForbidden", LogClientTimeout}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("DeleteLogStore", raw, requestInfo, map[string]interface{}{
		"project":  parts[0],
		"logstore": parts[1],
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "DeleteLogStore", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogStore(d.Id(), Deleted, DefaultTimeout))
}

func buildEncrypt(d *schema.ResourceData) *sls.EncryptConf {
	var encryptConf *sls.EncryptConf
	if field, ok := d.GetOk("encrypt_conf"); ok {
		encryptConf = new(sls.EncryptConf)
		value := field.(*schema.Set).List()[0].(map[string]interface{})
		encryptConf.Enable = value["enable"].(bool)
		encryptConf.EncryptType = value["encrypt_type"].(string)
		cmkInfo := value["user_cmk_info"].(*schema.Set).List()
		if len(cmkInfo) > 0 {
			cmk := cmkInfo[0].(map[string]interface{})
			encryptConf.UserCmkInfo = &sls.EncryptUserCmkConf{
				CmkKeyId: cmk["cmk_key_id"].(string),
				Arn:      cmk["arn"].(string),
				RegionId: cmk["region_id"].(string),
			}
		}
	}
	return encryptConf
}
