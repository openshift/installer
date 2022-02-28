package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamAccountAliasRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateGetAccountAliasRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_account_alias", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.GetAccountAliasResponse)
	d.SetId(response.AccountAlias)
	d.Set("account_alias", response.AccountAlias)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		s := map[string]interface{}{"account_alias": response.AccountAlias}
		writeToFile(output.(string), s)
	}
	return nil
}
