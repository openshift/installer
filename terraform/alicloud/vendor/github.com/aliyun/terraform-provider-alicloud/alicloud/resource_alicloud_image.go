package alicloud

import (
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageCreate,
		Read:   resourceAliCloudImageRead,
		Update: resourceAliCloudImageUpdate,
		Delete: resourceAliCloudImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "x86_64",
				ValidateFunc: validation.StringInSlice([]string{
					"x86_64",
					"i386",
				}, false),
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"disk_device_mapping", "snapshot_id"},
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id", "disk_device_mapping"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'name' has been deprecated from version 1.69.0. Use `image_name` instead.",
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Ubuntu",
				ValidateFunc: validation.StringInSlice([]string{
					"CentOS",
					"Ubuntu",
					"SUSE",
					"OpenSUSE",
					"RedHat",
					"Debian",
					"CoreOS",
					"Aliyun Linux",
					"Windows Server 2003",
					"Windows Server 2008",
					"Windows Server 2012",
					"Windows 7",
					"Customized Linux",
					"Others Linux",
				}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_device_mapping": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_id", "snapshot_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),
			// Not the public attribute and it used to automatically delete dependence snapshots while deleting the image.
			// Available in 1.136.0
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}
func resourceAliCloudImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	// Make sure the instance status is Running or Stopped
	if v, ok := d.GetOk("instance_id"); ok {
		instance, err := ecsService.DescribeInstance(v.(string))
		if err != nil {
			return WrapError(err)
		}
		status := Status(instance.Status)
		if status != Running && status != Stopped {
			return WrapError(Error("You must make sure that the status of the specified instance is Running or Stopped. "))
		}
	}

	// The snapshot cannot be a snapshot created before July 15, 2013 (inclusive)
	if snapshotId, ok := d.GetOk("snapshot_id"); ok {
		snapshot, err := ecsService.DescribeSnapshot(snapshotId.(string))
		if err != nil {
			return WrapError(err)
		}
		snapshotCreationTime, err := time.Parse("2006-01-02T15:04:05Z", snapshot.CreationTime)
		if err != nil {
			return WrapErrorf(err, IdMsg, snapshotId)
		}
		deadlineTime, _ := time.Parse("2006-01-02T15:04:05Z", "2013-07-16T00:00:00Z")
		if deadlineTime.After(snapshotCreationTime) {
			return WrapError(Error("the specified snapshot cannot be created on or before July 15, 2013."))
		}
	}
	request := ecs.CreateCreateImageRequest()
	request.RegionId = client.RegionId
	if instanceId, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = instanceId.(string)
	}
	if value, ok := d.GetOk("disk_device_mapping"); ok {
		diskDeviceMappings := value.([]interface{})
		if diskDeviceMappings != nil && len(diskDeviceMappings) > 0 {
			mappings := make([]ecs.CreateImageDiskDeviceMapping, 0, len(diskDeviceMappings))
			for _, diskDeviceMapping := range diskDeviceMappings {
				mapping := diskDeviceMapping.(map[string]interface{})
				deviceMapping := ecs.CreateImageDiskDeviceMapping{
					SnapshotId: mapping["snapshot_id"].(string),
					Size:       strconv.Itoa(mapping["size"].(int)),
					DiskType:   mapping["disk_type"].(string),
					Device:     mapping["device"].(string),
				}
				mappings = append(mappings, deviceMapping)
			}
			request.DiskDeviceMapping = &mappings
		}
	}

	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		imageTags := make([]ecs.CreateImageTag, 0, len(tags))
		for k, v := range tags {
			imageTag := ecs.CreateImageTag{
				Key:   k,
				Value: v.(string),
			}
			imageTags = append(imageTags, imageTag)
		}
		request.Tag = &imageTags
	}
	if snapshotId, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = snapshotId.(string)
	}
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.Platform = d.Get("platform").(string)
	request.ImageName = d.Get("image_name").(string)
	request.Description = d.Get("description").(string)
	request.Architecture = d.Get("architecture").(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateImage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
				time.Sleep(time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.CreateImageResponse)
		d.SetId(response.ImageId)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("resource_group_id") {
		action := "JoinResourceGroup"
		request := map[string]interface{}{
			"ResourceType":    "image",
			"ResourceId":      d.Id(),
			"RegionId":        client.RegionId,
			"ResourceGroupId": d.Get("resource_group_id"),
		}
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("resource_group_id")
	}
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_image ecsService.DescribeImageById Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("platform", object.Platform)
	d.Set("image_name", object.ImageName)
	d.Set("name", object.ImageName)
	d.Set("description", object.Description)
	d.Set("architecture", object.Architecture)
	d.Set("disk_device_mapping", FlattenImageDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping))
	tags, err := ecsService.ListTagResources(d.Id(), "image")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}
	return WrapError(err)
}

func resourceAliCloudImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	return ecsService.deleteImage(d)
}

func FlattenImageDiskDeviceMappings(list []ecs.DiskDeviceMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		size, _ := strconv.Atoi(i.Size)
		l := map[string]interface{}{
			"device":      i.Device,
			"size":        size,
			"snapshot_id": i.SnapshotId,
			"disk_type":   i.Type,
		}
		result = append(result, l)
	}

	return result
}
