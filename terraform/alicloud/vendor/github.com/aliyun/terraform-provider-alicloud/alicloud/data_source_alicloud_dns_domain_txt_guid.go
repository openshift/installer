package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDnsDomainTxtGuid() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsDomainTxtGuidRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADD_SUB_DOMAIN", "RETRIEVAL"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh", "jp"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudDnsDomainTxtGuidRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateGetTxtRecordForVerifyRequest()
	request.RegionId = client.RegionId
	domainName := d.Get("domain_name").(string)
	request.DomainName = domainName
	request.Type = d.Get("type").(string)
	if lang, ok := d.GetOk("lang"); ok {
		request.Lang = lang.(string)
	}

	raw, err := client.WithDnsClient(func(dnsclient *alidns.Client) (i interface{}, err error) {
		return dnsclient.GetTxtRecordForVerify(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dns_domain_txt_guid", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*alidns.GetTxtRecordForVerifyResponse)

	rr := response.RR
	value := response.Value

	s := map[string]interface{}{
		"rr":    rr,
		"value": value,
	}

	ids := []string{domainName, rr, value}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rr", rr); err != nil {
		return WrapError(err)
	}

	if err := d.Set("value", value); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
