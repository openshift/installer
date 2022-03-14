package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEciImageCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEciImageCacheCreate,
		Read:   resourceAlicloudEciImageCacheRead,
		Delete: resourceAlicloudEciImageCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"container_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_cache_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_cache_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"image_registry_credential": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				ForceNew: true,
			},
			"images": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEciImageCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}

	request := eci.CreateCreateImageCacheRequest()
	if v, ok := d.GetOk("eip_instance_id"); ok {
		request.EipInstanceId = v.(string)
	}
	request.ImageCacheName = d.Get("image_cache_name").(string)
	if v, ok := d.GetOk("image_cache_size"); ok {
		request.ImageCacheSize = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("image_registry_credential"); ok {
		imageRegistryCredential := []eci.CreateImageCacheImageRegistryCredential{}
		for _, e := range v.(*schema.Set).List() {
			password := e.(map[string]interface{})["password"]
			server := e.(map[string]interface{})["server"]
			userName := e.(map[string]interface{})["user_name"]
			imageRegistryCredential = append(imageRegistryCredential, eci.CreateImageCacheImageRegistryCredential{
				Password: password.(string),
				Server:   server.(string),
				UserName: userName.(string),
			})
		}
		request.ImageRegistryCredential = &imageRegistryCredential
	}
	image := expandStringList(d.Get("images").(*schema.Set).List())
	request.Image = image
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("retention_days"); ok {
		request.RetentionDays = requests.NewInteger(v.(int))
	}
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}

	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.CreateImageCache(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eci_image_cache", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*eci.CreateImageCacheResponse)
	d.SetId(fmt.Sprintf("%v", response.ImageCacheId))
	stateConf := BuildStateConf([]string{"Preparing", "Creating"}, []string{"Ready"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, eciService.EciImageCacheStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEciImageCacheRead(d, meta)
}
func resourceAlicloudEciImageCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	object, err := eciService.DescribeEciImageCache(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("container_group_id", object.ContainerGroupId)
	d.Set("image_cache_name", object.ImageCacheName)
	d.Set("images", object.Images)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudEciImageCacheDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := eci.CreateDeleteImageCacheRequest()
	request.ImageCacheId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DeleteImageCache(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
