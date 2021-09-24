package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceAttachmentCreate,
		Read:   resourceAlicloudCenInstanceAttachmentRead,
		Update: resourceAlicloudCenInstanceAttachmentUpdate,
		Delete: resourceAlicloudCenInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"child_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"child_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "VBR", "CCN"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateAttachCenChildInstanceRequest()
	request.ChildInstanceId = d.Get("child_instance_id").(string)
	if v, ok := d.GetOk("child_instance_owner_id"); ok {
		request.ChildInstanceOwnerId = requests.NewInteger(v.(int))
	}

	request.ChildInstanceRegionId = d.Get("child_instance_region_id").(string)
	request.ChildInstanceType = d.Get("child_instance_type").(string)
	request.CenId = d.Get("instance_id").(string)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AttachCenChildInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "InvalidOperation.ChildInstanceStatus", "Operation.Blocking", "OperationFailed.InvalidVpcStatus", "Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(fmt.Sprintf("%v:%v:%v:%v", request.CenId, request.ChildInstanceId, request.ChildInstanceType, request.ChildInstanceRegionId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_instance_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenInstanceAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenInstanceAttachmentRead(d, meta)
}
func resourceAlicloudCenInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	if len(strings.Split(d.Id(), ":")) == 2 {
		childType, _ := GetCenChildInstanceType(d.Get("child_instance_id").(string))
		d.SetId(fmt.Sprintf("%v:%v:%v", d.Id(), childType, d.Get("child_instance_region_id").(string)))
	}
	object, err := cbnService.DescribeCenInstanceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_instance_attachment cbnService.DescribeCenInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("child_instance_id", parts[1])
	d.Set("child_instance_region_id", parts[3])
	d.Set("child_instance_type", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("child_instance_owner_id", object.ChildInstanceOwnerId)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudCenInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudCenInstanceAttachmentRead(d, meta)
}
func resourceAlicloudCenInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) == 2 {
		childType, _ := GetCenChildInstanceType(d.Get("child_instance_id").(string))
		d.SetId(fmt.Sprintf("%v:%v:%v", d.Id(), childType, d.Get("child_instance_region_id").(string)))
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	cbnService := CbnService{client}
	request := cbn.CreateDetachCenChildInstanceRequest()
	request.ChildInstanceId = parts[1]
	request.ChildInstanceRegionId = parts[3]
	request.ChildInstanceType = parts[2]
	request.CenId = parts[0]
	if v, ok := d.GetOk("cen_owner_id"); ok {
		request.CenOwnerId = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("child_instance_owner_id"); ok {
		request.ChildInstanceOwnerId = requests.NewInteger(v.(int))
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DetachCenChildInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Bid.NotFound", "ParameterInstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenInstanceAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
