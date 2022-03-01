package alicloud

import (
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudImagesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"owners": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// must contain a valid Image owner, expected ImageOwnerSystem, ImageOwnerSelf, ImageOwnerOthers, ImageOwnerMarketplace, ImageOwnerDefault
				ValidateFunc: validation.StringInSlice([]string{"system", "self", "others", "marketplace", ""}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Available",
				ValidateFunc: validation.StringInSlice([]string{"Available", "Creating", "Waiting", "UnAvailable", "CreateFailed", "Deprecated"}, false),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_support_io_optimized": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_support_cloud_init": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_family": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"instance", "none"}, false),
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"windows", "linux"}, false),
			},
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"i386", "x86_64"}, false),
			},
			"action_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CreateEcs", "CreateOS"}, false),
			},
			"tags": tagsSchema(),
			// Computed values.
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_owner_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_name_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// Complex computed values
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							//Set:      imageDiskDeviceMappingHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_self_shared": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_subscribed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_copied": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_io_optimized": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"image_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"usage": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

// dataSourceAlicloudImagesDescriptionRead performs the Alicloud Image lookup.
func dataSourceAlicloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	owners, ownersOk := d.GetOk("owners")
	mostRecent, mostRecentOk := d.GetOk("most_recent")

	if nameRegexOk == false && ownersOk == false && mostRecentOk == false {
		return WrapError(Error("One of name_regex, owners or most_recent must be assigned"))
	}

	request := ecs.CreateDescribeImagesRequest()
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeXLarge)

	if ownersOk {
		request.ImageOwnerAlias = owners.(string)
	}

	if status, ok := d.GetOk("status"); ok && status.(string) != "" {
		request.Status = status.(string)
	}

	if v, ok := d.GetOk("image_id"); ok && v.(string) != "" {
		request.ImageId = v.(string)
	}

	if v, ok := d.GetOk("image_name"); ok && v.(string) != "" {
		request.ImageName = v.(string)
	}

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		request.SnapshotId = v.(string)
	}

	if v, ok := d.GetOk("image_family"); ok && v.(string) != "" {
		request.ImageFamily = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok && v.(string) != "" {
		request.InstanceType = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("usage"); ok && v.(string) != "" {
		request.Usage = v.(string)
	}

	if v, ok := d.GetOk("architecture"); ok && v.(string) != "" {
		request.Architecture = v.(string)
	}

	if v, ok := d.GetOk("os_type"); ok && v.(string) != "" {
		request.OSType = v.(string)
	}

	if v, ok := d.GetOk("action_type"); ok && v.(string) != "" {
		request.ActionType = v.(string)
	}

	if v, ok := d.GetOk("is_support_io_optimized"); ok {
		request.IsSupportIoOptimized = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("is_support_cloud_init"); ok {
		request.IsSupportCloudinit = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request.DryRun = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("tags"); ok {
		var reqTags []ecs.DescribeImagesTag
		for k, v := range v.(map[string]interface{}) {
			reqTags = append(reqTags, ecs.DescribeImagesTag{
				Key:   k,
				Value: v.(string),
			})
		}
		request.Tag = &reqTags
	}

	var allImages []ecs.Image

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeImages(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_images", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*ecs.DescribeImagesResponse)
		if response == nil || len(response.Images.Image) < 1 {
			break
		}

		allImages = append(allImages, response.Images.Image...)

		if len(response.Images.Image) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredImages []ecs.Image
	if nameRegexOk {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return WrapError(err)
		}
		for _, image := range allImages {
			// Check for a very rare case where the response would include no
			// image name. No name means nothing to attempt a match against,
			// therefore we are skipping such image.
			if image.ImageName == "" {
				log.Printf("[WARN] Unable to find Image name to match against "+
					"for image ID %q, nothing to do.",
					image.ImageId)
				continue
			}
			if r.MatchString(image.ImageName) {
				filteredImages = append(filteredImages, image)
			}
		}
	} else {
		filteredImages = allImages[:]
	}

	var images []ecs.Image

	if len(filteredImages) > 1 && mostRecent.(bool) {
		// Query returned single result.
		images = append(images, mostRecentImage(filteredImages))
	} else {
		images = filteredImages
	}

	return imagesDescriptionAttributes(d, images, meta)
}

// populate the numerous fields that the image description returns.
func imagesDescriptionAttributes(d *schema.ResourceData, images []ecs.Image, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, image := range images {
		mapping := map[string]interface{}{
			"id":                      image.ImageId,
			"architecture":            image.Architecture,
			"creation_time":           image.CreationTime,
			"description":             image.Description,
			"image_id":                image.ImageId,
			"image_owner_alias":       image.ImageOwnerAlias,
			"os_name":                 image.OSName,
			"os_name_en":              image.OSNameEn,
			"os_type":                 image.OSType,
			"name":                    image.ImageName,
			"platform":                image.Platform,
			"status":                  image.Status,
			"state":                   image.Status,
			"size":                    image.Size,
			"is_self_shared":          image.IsSelfShared,
			"is_subscribed":           image.IsSubscribed,
			"is_copied":               image.IsCopied,
			"is_support_io_optimized": image.IsSupportIoOptimized,
			"image_version":           image.ImageVersion,
			"progress":                image.Progress,
			"usage":                   image.Usage,
			"product_code":            image.ProductCode,

			// Complex types get their own functions
			"disk_device_mappings": imageDiskDeviceMappings(image.DiskDeviceMappings.DiskDeviceMapping),
			"tags":                 imageTagsMappings(d, image.ImageId, meta),
		}

		ids = append(ids, image.ImageId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("images", s); err != nil {
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

//Find most recent image
type imageSort []ecs.Image

func (a imageSort) Len() int {
	return len(a)
}
func (a imageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a imageSort) Less(i, j int) bool {
	itime, _ := time.Parse(time.RFC3339, a[i].CreationTime)
	jtime, _ := time.Parse(time.RFC3339, a[j].CreationTime)
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []ecs.Image) ecs.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

// Returns a set of disk device mappings.
func imageDiskDeviceMappings(m []ecs.DiskDeviceMapping) []map[string]interface{} {
	var s []map[string]interface{}

	for _, v := range m {
		mapping := map[string]interface{}{
			"device":      v.Device,
			"size":        v.Size,
			"snapshot_id": v.SnapshotId,
		}

		s = append(s, mapping)
	}

	return s
}

//Returns a mapping of image tags
func imageTagsMappings(d *schema.ResourceData, imageId string, meta interface{}) map[string]string {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	tags, err := ecsService.DescribeTags(imageId, TagResourceImage)

	if err != nil {
		return nil
	}

	return ecsTagsToMap(tags)
}
