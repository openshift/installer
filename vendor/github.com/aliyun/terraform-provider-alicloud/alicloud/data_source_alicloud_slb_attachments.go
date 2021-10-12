package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSlbAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"slb_attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	instanceIdsMap := make(map[string]string)
	if v, ok := d.GetOk("instance_ids"); ok {
		for _, vv := range v.([]interface{}) {
			instanceIdsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_attachments", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	var filteredBackendServersTemp []slb.BackendServerInDescribeLoadBalancerAttribute
	if len(instanceIdsMap) > 0 {
		for _, backendServer := range response.BackendServers.BackendServer {
			if len(instanceIdsMap) > 0 {
				if _, ok := instanceIdsMap[backendServer.ServerId]; !ok {
					continue
				}
			}

			filteredBackendServersTemp = append(filteredBackendServersTemp, backendServer)
		}
	} else {
		filteredBackendServersTemp = response.BackendServers.BackendServer
	}

	return slbAttachmentsDescriptionAttributes(d, filteredBackendServersTemp)
}

func slbAttachmentsDescriptionAttributes(d *schema.ResourceData, backendServers []slb.BackendServerInDescribeLoadBalancerAttribute) error {
	var ids []string
	var s []map[string]interface{}

	for _, backendServer := range backendServers {
		mapping := map[string]interface{}{
			"instance_id": backendServer.ServerId,
			"weight":      backendServer.Weight,
		}

		ids = append(ids, backendServer.ServerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_attachments", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
