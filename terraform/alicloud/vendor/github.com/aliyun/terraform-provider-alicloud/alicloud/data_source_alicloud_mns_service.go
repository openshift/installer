package alicloud

import (
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMnsService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMnsServiceRead,

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
func dataSourceAlicloudMnsServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("MnsServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.WithMnsClient(func(client *ali_mns.MNSClient) (interface{}, error) {
			return ali_mns.NewAccountManager(*client).OpenService()
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug("OpenService", response, nil)
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.Opened"}) {
			d.SetId("MnsServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mns_service", "OpenService", AlibabaCloudSdkGoERROR)
	}
	d.SetId("MnsServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
