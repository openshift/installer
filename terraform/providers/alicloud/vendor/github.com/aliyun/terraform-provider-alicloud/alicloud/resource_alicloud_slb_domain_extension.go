package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbDomainExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbDomainExtensionCreate,
		Read:   resourceAliyunSlbDomainExtensionRead,
		Update: resourceAliyunSlbDomainExtensionUpdate,
		Delete: resourceAliyunSlbDomainExtensionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Required:     true,
				ForceNew:     true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunSlbDomainExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateCreateDomainExtensionRequest()
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))
	request.Domain = d.Get("domain").(string)
	request.ServerCertificateId = d.Get("server_certificate_id").(string)

	var response *slb.CreateDomainExtensionResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateDomainExtension(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DomainExtensionProcessing"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*slb.CreateDomainExtensionResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_domain_extension", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.DomainExtensionId)
	return resourceAliyunSlbDomainExtensionRead(d, meta)
}

func resourceAliyunSlbDomainExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	domainExtension, err := slbService.DescribeDomainExtensionAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("id", domainExtension.DomainExtensionId)
	d.Set("load_balancer_id", domainExtension.LoadBalancerId)
	d.Set("domain", domainExtension.Domain)
	d.Set("server_certificate_id", domainExtension.ServerCertificateId)
	d.Set("frontend_port", domainExtension.ListenerPort)

	return nil
}

func resourceAliyunSlbDomainExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("server_certificate_id") {
		request := slb.CreateSetDomainExtensionAttributeRequest()
		request.DomainExtensionId = d.Id()
		request.ServerCertificateId = d.Get("server_certificate_id").(string)
		client := meta.(*connectivity.AliyunClient)
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetDomainExtensionAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"BackendServer.configuring", "DomainExtensionProcessing"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_domain_extension", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("server_certificate_id")
	}
	return resourceAliyunSlbDomainExtensionRead(d, meta)
}

func resourceAliyunSlbDomainExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return WrapError(fmt.Errorf("Current domain extension's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteDomainExtensionRequest()
	request.DomainExtensionId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteDomainExtension(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DomainExtensionProcessing", "InternalError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId"}) {
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)

	}
	return WrapError(slbService.WaitForSlbDomainExtension(d.Id(), Deleted, DefaultTimeoutMedium))
}
