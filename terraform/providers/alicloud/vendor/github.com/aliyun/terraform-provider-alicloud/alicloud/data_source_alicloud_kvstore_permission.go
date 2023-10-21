package alicloud

import (
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKVStorePermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStorePermissionRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudKVStorePermissionRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("KVStorePermissionHasNotBeenInitialize")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	request := r_kvstore.CreateInitializeKvstorePermissionRequest()
	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.InitializeKvstorePermission(request)
	})
	if err != nil {
		return WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId("KVStorePermissionHasBeenInitialize")
	d.Set("status", "Initialized")

	return nil
}
