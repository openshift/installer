package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbAttachmentCreate,
		Read:   resourceAliyunSlbAttachmentRead,
		Update: resourceAliyunSlbAttachmentUpdate,
		Delete: resourceAliyunSlbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(0, 100),
			},

			"server_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ecs",
				ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
			},

			"backend_servers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Get("load_balancer_id").(string))
	if err != nil {
		return WrapError(err)
	}
	d.SetId(object.LoadBalancerId)

	return resourceAliyunSlbAttachmentUpdate(d, meta)
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	servers := object.BackendServers.BackendServer
	instanceIds := make([]string, 0, len(servers))
	var weight int
	var serverType string
	if len(servers) > 0 {
		weight = servers[0].Weight
		serverType = servers[0].Type
		for _, e := range servers {
			instanceIds = append(instanceIds, e.ServerId)
		}
	}

	d.Set("load_balancer_id", object.LoadBalancerId)
	d.Set("instance_ids", instanceIds)
	d.Set("weight", weight)
	d.Set("server_type", serverType)
	d.Set("backend_servers", strings.Join(instanceIds, ","))

	return nil
}

func resourceAliyunSlbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	update := false
	weight := d.Get("weight").(int)
	oldServerType, serverType := d.GetChange("server_type")

	if d.HasChange("server_type") {
		update = true
		d.SetPartial("server_type")
	}
	if d.HasChange("weight") {
		update = true
		d.SetPartial("weight")
	}
	if d.HasChange("instance_ids") {
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(add) > 0 {
			request := slb.CreateAddBackendServersRequest()
			request.RegionId = client.RegionId
			request.LoadBalancerId = d.Id()
			request.BackendServers = expandBackendServersToString(ns.Difference(os).List(), weight, serverType.(string))
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.AddBackendServers(request)
				})
				if err != nil {
					if IsExpectedErrors(err, SlbIsBusy) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
		if len(remove) > 0 {
			request := slb.CreateRemoveBackendServersRequest()
			request.RegionId = client.RegionId
			request.LoadBalancerId = d.Id()
			request.BackendServers = expandBackendServersToString(os.Difference(ns).List(), weight, oldServerType.(string))
			if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
					return slbClient.RemoveBackendServers(request)
				})
				if err != nil {
					if IsExpectedErrors(err, SlbIsBusy) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		if len(add) < 1 && len(remove) < 1 {
			update = true
		}
		d.SetPartial("instance_ids")
	}

	if update {
		request := slb.CreateSetBackendServersRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.BackendServers = expandBackendServersToString(d.Get("instance_ids").(*schema.Set).List(), weight, serverType.(string))
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetBackendServers(request)
			})
			if err != nil {
				if IsExpectedErrors(err, SlbIsBusy) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliyunSlbAttachmentRead(d, meta)

}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbInstance, err := slbService.DescribeSlb(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", d.Id()))
		}
	}

	instanceSet := d.Get("instance_ids").(*schema.Set)
	weight := d.Get("weight").(int)
	serverType := d.Get("server_type").(string)
	if len(instanceSet.List()) > 0 {
		request := slb.CreateRemoveBackendServersRequest()
		request.RegionId = client.RegionId
		request.LoadBalancerId = d.Id()
		request.BackendServers = expandBackendServersToString(d.Get("instance_ids").(*schema.Set).List(), weight, serverType)
		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.RemoveBackendServers(request)
			})
			if err != nil {
				if IsExpectedErrors(err, SlbIsBusy) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return WrapError(slbService.WaitSlbAttribute(d.Id(), instanceSet, DefaultTimeout))
}
