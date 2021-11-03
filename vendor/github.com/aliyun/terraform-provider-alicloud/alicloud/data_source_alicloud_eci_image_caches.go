package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEciImageCaches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEciImageCachesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"image": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"image_cache_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Failed", "Preparing", "Ready"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"caches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"first_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"expire_date_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_cache_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_cache_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
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

func dataSourceAlicloudEciImageCachesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := eci.CreateDescribeImageCachesRequest()
	if v, ok := d.GetOk("image"); ok {
		request.Image = v.(string)
	}
	if v, ok := d.GetOk("image_cache_name"); ok {
		request.ImageCacheName = v.(string)
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = v.(string)
	}
	var objects []eci.DescribeImageCachesImageCache0
	var imageCacheNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		imageCacheNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response *eci.DescribeImageCachesResponse
	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eci_image_caches", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*eci.DescribeImageCachesResponse)

	for _, item := range response.ImageCaches {
		if imageCacheNameRegex != nil {
			if !imageCacheNameRegex.MatchString(item.ImageCacheName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[item.ImageCacheId]; !ok {
				continue
			}
		}
		if statusOk && status != "" && status != item.Status {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"container_group_id": object.ContainerGroupId,
			"expire_date_time":   object.ExpireDateTime,
			"id":                 object.ImageCacheId,
			"image_cache_id":     object.ImageCacheId,
			"image_cache_name":   object.ImageCacheName,
			"images":             object.Images,
			"progress":           object.Progress,
			"snapshot_id":        object.SnapshotId,
			"status":             object.Status,
		}

		var eventsList []map[string]interface{}
		for _, v := range object.Events {
			eventsList = append(eventsList, map[string]interface{}{
				"count":           v.Count,
				"first_timestamp": v.FirstTimestamp,
				"last_timestamp":  v.LastTimestamp,
				"message":         v.Message,
				"name":            v.Name,
				"type":            v.Type,
			})
		}
		mapping["events"] = eventsList
		ids[i] = object.ImageCacheId
		names = append(names, object.ImageCacheName)
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("caches", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
