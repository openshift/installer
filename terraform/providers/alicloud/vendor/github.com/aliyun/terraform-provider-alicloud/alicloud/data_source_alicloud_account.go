package alicloud

import (
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccountRead,

		Schema: map[string]*schema.Schema{
			// Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudAccountRead(d *schema.ResourceData, meta interface{}) error {
	accountId, err := meta.(*connectivity.AliyunClient).AccountId()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] alicloud_account - account ID found: %#v", accountId)

	d.SetId(accountId)

	return nil
}
