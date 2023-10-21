package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSlbDomainExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbDomainExtensionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbDomainExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeDomainExtensionsRequest()
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeDomainExtensions(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_domain_extensions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*slb.DescribeDomainExtensionsResponse)
	var filteredDomainExtensionsTemp []slb.DomainExtension
	if len(idsMap) > 0 {
		for _, domainExtension := range response.DomainExtensions.DomainExtension {
			if len(idsMap) > 0 {
				if _, ok := idsMap[domainExtension.DomainExtensionId]; !ok {
					continue
				}
			}
			filteredDomainExtensionsTemp = append(filteredDomainExtensionsTemp, domainExtension)
		}
	} else {
		filteredDomainExtensionsTemp = response.DomainExtensions.DomainExtension
	}
	return slbDomainExtensionDescriptionAttributes(d, filteredDomainExtensionsTemp)
}

func slbDomainExtensionDescriptionAttributes(d *schema.ResourceData, domainExtensions []slb.DomainExtension) error {
	var ids []string
	var s []map[string]interface{}
	for _, domainExtension := range domainExtensions {
		mapping := map[string]interface{}{
			"id":                    domainExtension.DomainExtensionId,
			"domain":                domainExtension.Domain,
			"server_certificate_id": domainExtension.ServerCertificateId,
		}
		ids = append(ids, domainExtension.DomainExtensionId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("extensions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
