package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDatabaseCreate,
		Read:   resourceAlicloudPolarDBDatabaseRead,
		Update: resourceAlicloudPolarDBDatabaseUpdate,
		Delete: resourceAlicloudPolarDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"db_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"character_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
				ForceNew: true,
			},

			"db_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPolarDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := polardb.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.DBName = d.Get("db_name").(string)
	request.CharacterSetName = d.Get("character_set_name").(string)

	if v, ok := d.GetOk("db_description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDatabase(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
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

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.DBName))

	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("db_cluster_id", parts[0])
	d.Set("db_name", object.DBName)
	d.Set("character_set_name", object.CharacterSetName)
	d.Set("db_description", object.DBDescription)

	return nil
}

func resourceAlicloudPolarDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("db_description") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyDBDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("db_description").(string)
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBDescription(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateDeleteDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := polarDBService.WaitForPolarDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteDatabase(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(polarDBService.WaitForPolarDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
