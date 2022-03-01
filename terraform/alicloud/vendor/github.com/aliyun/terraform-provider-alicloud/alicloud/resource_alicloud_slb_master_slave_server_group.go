package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSlbMasterSlaveServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbMasterSlaveServerGroupCreate,
		Read:   resourceAliyunSlbMasterSlaveServerGroupRead,
		Update: resourceAliyunSlbMasterSlaveServerGroupUpdate,
		Delete: resourceAliyunSlbMasterSlaveServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(ECS),
							ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
						},
						"server_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Master", "Slave"}, false),
						},
						"is_backup": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1}),
							Removed:      "Field 'is_backup' has been removed from provider version 1.63.0.",
						},
					},
				},
				MaxItems: 2,
				MinItems: 2,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunSlbMasterSlaveServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateCreateMasterSlaveServerGroupRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	if v, ok := d.GetOk("name"); ok {
		request.MasterSlaveServerGroupName = v.(string)
	}
	if v, ok := d.GetOk("servers"); ok {
		request.MasterSlaveBackendServers = expandMasterSlaveBackendServersToString(v.(*schema.Set).List())
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateMasterSlaveServerGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_master_slave_server_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateMasterSlaveServerGroupResponse)
	d.SetId(response.MasterSlaveServerGroupId)

	return resourceAliyunSlbMasterSlaveServerGroupRead(d, meta)
}

func resourceAliyunSlbMasterSlaveServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbMasterSlaveServerGroup(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.MasterSlaveServerGroupName)
	d.Set("load_balancer_id", object.LoadBalancerId)

	servers := make([]map[string]interface{}, 0)

	for _, server := range object.MasterSlaveBackendServers.MasterSlaveBackendServer {
		s := map[string]interface{}{
			"server_id":   server.ServerId,
			"port":        server.Port,
			"weight":      server.Weight,
			"type":        server.Type,
			"server_type": server.ServerType,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunSlbMasterSlaveServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAliyunSlbMasterSlaveServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
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
			return WrapError(fmt.Errorf("Current master-slave server group's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteMasterSlaveServerGroupRequest()
	request.RegionId = client.RegionId
	request.MasterSlaveServerGroupId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteMasterSlaveServerGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RspoolVipExist"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbMasterSlaveServerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
