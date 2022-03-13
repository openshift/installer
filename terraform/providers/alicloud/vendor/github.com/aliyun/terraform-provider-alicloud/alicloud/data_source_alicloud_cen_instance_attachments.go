package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenInstanceAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenInstanceAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"child_instance_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"child_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "VBR", "CCN"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Aetaching", "Attached", "Attaching"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_attach_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_owner_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"child_instance_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"child_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenInstanceAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	if v, ok := d.GetOk("child_instance_region_id"); ok {
		request.ChildInstanceRegionId = v.(string)
	}
	if v, ok := d.GetOk("child_instance_type"); ok {
		request.ChildInstanceType = v.(string)
	}
	request.CenId = d.Get("instance_id").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.ChildInstance
	status, statusOk := d.GetOk("status")
	var response *cbn.DescribeCenAttachedChildInstancesResponse
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_instance_attachments", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*cbn.DescribeCenAttachedChildInstancesResponse)

		for _, item := range response.ChildInstances.ChildInstance {
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.ChildInstances.ChildInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                         fmt.Sprintf("%v:%v:%v:%v", request.CenId, object.ChildInstanceId, object.ChildInstanceType, object.ChildInstanceRegionId),
			"child_instance_attach_time": object.ChildInstanceAttachTime,
			"child_instance_id":          object.ChildInstanceId,
			"child_instance_owner_id":    object.ChildInstanceOwnerId,
			"child_instance_region_id":   object.ChildInstanceRegionId,
			"child_instance_type":        object.ChildInstanceType,
			"instance_id":                request.CenId,
			"status":                     object.Status,
		}
		ids = append(ids, fmt.Sprintf("%v:%v:%v:%v", request.CenId, object.ChildInstanceId, object.ChildInstanceType, object.ChildInstanceRegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
