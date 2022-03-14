package alicloud

import (
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipAssociationCreate,
		Read:   resourceAliyunEipAssociationRead,
		Update: resourceAliyunEipAssociationUpdate,
		Delete: resourceAliyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateEipAddressRequest()
	request.RegionId = client.RegionId
	request.AllocationId = Trim(d.Get("allocation_id").(string))
	request.InstanceId = Trim(d.Get("instance_id").(string))
	request.InstanceType = EcsInstance
	// There is a product api bug about clientToken and after fixed , the clientToken will be opened again.
	//request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(request.InstanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(request.InstanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	if privateIPAddress, ok := d.GetOk("private_ip_address"); ok {
		request.PrivateIpAddress = privateIPAddress.(string)
	}
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_association", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := vpcService.WaitForEip(request.AllocationId, InUse, 60); err != nil {
		return WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if request.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(request.AllocationId + ":" + request.InstanceId)

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object.InstanceId)
	d.Set("allocation_id", object.AllocationId)
	d.Set("instance_type", object.InstanceType)
	d.Set("force", d.Get("force").(bool))
	return nil
}

func resourceAliyunEipAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] The update method is used to ensure that the force parameter does not need to add forcenew.")
	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.RegionId = client.RegionId
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.InstanceType = EcsInstance
	request.Force = requests.NewBoolean(d.Get("force").(bool))
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectHaVipStatus", "TaskConflict",
				"InvalidIpStatus.HasBeenUsedBySnatTable", "InvalidIpStatus.HasBeenUsedByForwardEntry", "InvalidStatus.EniStatusNotSupport"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEipAssociation(d.Id(), Available, DefaultTimeoutMedium))
}
