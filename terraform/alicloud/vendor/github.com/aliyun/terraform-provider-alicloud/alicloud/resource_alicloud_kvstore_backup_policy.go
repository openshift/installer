package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudKVStoreBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreBackupPolicyCreate,
		Read:   resourceAlicloudKVStoreBackupPolicyRead,
		Update: resourceAlicloudKVStoreBackupPolicyUpdate,
		Delete: resourceAlicloudKVStoreBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},
			"backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKVStoreBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("instance_id").(string))

	return resourceAlicloudKVStoreBackupPolicyUpdate(d, meta)
}

func resourceAlicloudKVStoreBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeKVstoreBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("backup_time", object.PreferredBackupTime)
	d.Set("backup_period", strings.Split(object.PreferredBackupPeriod, ","))

	return nil
}

func resourceAlicloudKVStoreBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	if d.HasChange("backup_time") || d.HasChange("backup_period") {
		client := meta.(*connectivity.AliyunClient)
		kvstoreService := KvstoreService{client}

		request := r_kvstore.CreateModifyBackupPolicyRequest()
		request.RegionId = client.RegionId
		request.InstanceId = d.Id()
		request.PreferredBackupTime = d.Get("backup_time").(string)
		periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
		request.PreferredBackupPeriod = strings.Join(periodList, COMMA_SEPARATED)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyBackupPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		// There is a random error and need waiting some seconds to ensure the update is success
		_, err = kvstoreService.DescribeKVstoreBackupPolicy(d.Id())
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudKVStoreBackupPolicyRead(d, meta)
}

func resourceAlicloudKVStoreBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// In case of a delete we are resetting to default values which is Monday - Sunday each 3am-4am
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	request := r_kvstore.CreateModifyBackupPolicyRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()

	request.PreferredBackupTime = "01:00Z-02:00Z"
	request.PreferredBackupPeriod = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
