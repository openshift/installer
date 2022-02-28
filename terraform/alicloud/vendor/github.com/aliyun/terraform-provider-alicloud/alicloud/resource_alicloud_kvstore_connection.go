package alicloud

import (
	"fmt"
	"log"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudKvstoreConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKvstoreConnectionCreate,
		Read:   resourceAlicloudKvstoreConnectionRead,
		Update: resourceAlicloudKvstoreConnectionUpdate,
		Delete: resourceAlicloudKvstoreConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudKvstoreConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}

	request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
	request.ConnectionStringPrefix = d.Get("connection_string_prefix").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = d.Get("port").(string)
	request.RegionId = client.RegionId

	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.AllocateInstancePublicConnection(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", request.InstanceId))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudKvstoreConnectionUpdate(d, meta)
}
func resourceAlicloudKvstoreConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_connection r_kvstoreService.DescribeKvstoreConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	for _, instanceNetInfo := range object {
		if instanceNetInfo.DBInstanceNetType == "0" {
			d.Set("connection_string", instanceNetInfo.ConnectionString)
			d.Set("port", instanceNetInfo.Port)
		}
	}
	return nil
}
func resourceAlicloudKvstoreConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	update := false
	request := r_kvstore.CreateModifyDBInstanceConnectionStringRequest()
	request.DBInstanceId = d.Id()

	request.CurrentConnectionString = d.Get("connection_string").(string)
	if !d.IsNewResource() && d.HasChange("connection_string_prefix") {
		update = true
	}
	request.NewConnectionString = d.Get("connection_string_prefix").(string)
	request.IPType = "Public"
	if !d.IsNewResource() && d.HasChange("port") {
		update = true
		request.Port = d.Get("port").(string)
	}
	if update {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyDBInstanceConnectionString(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudKvstoreConnectionRead(d, meta)
}
func resourceAlicloudKvstoreConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
	request.InstanceId = d.Id()
	request.CurrentConnectionString = d.Get("connection_string").(string)
	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.ReleaseInstancePublicConnection(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
