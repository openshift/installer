package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCallerIdentityRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resp, err := client.GetCallerIdentity()
	if err != nil {
		return err
	}
	d.SetId(resp.PrincipalId)
	if err := d.Set("account_id", resp.AccountId); err != nil {
		return err
	}
	if err := d.Set("arn", resp.Arn); err != nil {
		return err
	}
	if err := d.Set("identity_type", resp.IdentityType); err != nil {
		return err
	}
	return nil
}
