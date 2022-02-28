package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdBundles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdBundlesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"bundle_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"bundle_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SYSTEM", "CUSTOM"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bundles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bundle_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bundle_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bundle_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_type_attribute": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gpu_spec": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdBundlesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBundles"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("bundle_id"); ok {
		request["BundleId"] = v
	}
	if v, ok := d.GetOk("bundle_type"); ok {
		request["BundleType"] = v
	}

	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var bundleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		bundleNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_bundles", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Bundles", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Bundles", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if bundleNameRegex != nil && !bundleNameRegex.MatchString(fmt.Sprint(item["BundleName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["BundleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":           fmt.Sprint(object["BundleId"]),
			"bundle_id":    fmt.Sprint(object["BundleId"]),
			"bundle_name":  object["BundleName"],
			"bundle_type":  object["BundleType"],
			"description":  object["Description"],
			"desktop_type": object["DesktopType"],
			"image_id":     object["ImageId"],
			"os_type":      object["OsType"],
		}

		disks := make([]map[string]interface{}, 0)
		if disksList, ok := object["Disks"].([]interface{}); ok {
			for _, v := range disksList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"disk_size": m1["DiskSize"],
						"disk_type": m1["DiskType"],
					}
					disks = append(disks, temp1)
				}
			}
		}
		mapping["disks"] = disks

		desktopTypeAttributeSli := make([]map[string]interface{}, 0)
		if len(object["DesktopTypeAttribute"].(map[string]interface{})) > 0 {
			desktopTypeAttribute := object["DesktopTypeAttribute"]
			desktopTypeAttributeMap := make(map[string]interface{})
			desktopTypeAttributeMap["cpu_count"] = desktopTypeAttribute.(map[string]interface{})["CpuCount"]
			desktopTypeAttributeMap["gpu_count"] = desktopTypeAttribute.(map[string]interface{})["GpuCount"]
			desktopTypeAttributeMap["gpu_spec"] = desktopTypeAttribute.(map[string]interface{})["GpuSpec"]
			desktopTypeAttributeMap["memory_size"] = desktopTypeAttribute.(map[string]interface{})["MemorySize"]
			desktopTypeAttributeSli = append(desktopTypeAttributeSli, desktopTypeAttributeMap)
		}
		mapping["desktop_type_attribute"] = desktopTypeAttributeSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["BundleName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("bundles", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
