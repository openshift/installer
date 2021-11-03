package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEdasSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasSlbAttachmentCreate,
		Read:   resourceAlicloudEdasSlbAttachmentRead,
		Delete: resourceAlicloudEdasSlbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slb_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slb_ip": {
				Type:         schema.TypeString,
				ValidateFunc: validation.SingleIP(),
				Required:     true,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
				ForceNew:     true,
			},
			"vserver_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slb_status": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	slbId := d.Get("slb_id").(string)

	request := edas.CreateBindSlbRequest()
	request.RegionId = client.RegionId
	request.Type = d.Get("type").(string)
	request.AppId = appId
	request.SlbId = slbId
	request.SlbIp = d.Get("slb_ip").(string)
	if v, ok := d.GetOk("listener_port"); ok {
		request.ListenerPort = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("vserver_group_id"); ok {
		request.VServerGroupId = v.(string)
	}

	if err := edasService.SyncResource("slb"); err != nil {
		return err
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.BindSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.BindSlbResponse)
	if response.Code != 200 {
		return WrapError(Error("bind slb failed for " + response.Message))
	}
	d.SetId(appId + ":" + slbId)
	return resourceAlicloudEdasInstanceApplicationAttachmentRead(d, meta)
}

func resourceAlicloudEdasSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	id := d.Id()
	strs, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}

	regionId := client.RegionId
	slbId := strs[1]
	appId := strs[0]

	rq := edas.CreateGetApplicationRequest()
	rq.RegionId = regionId
	rq.AppId = appId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(rq)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", rq.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(rq.GetActionName(), raw, rq.RoaRequest, rq)

	rs := raw.(*edas.GetApplicationResponse)
	if rs.Applcation.SlbId != slbId && rs.Applcation.ExtSlbId != slbId {
		return WrapError(Error("can not find slb:" + slbId))
	}

	request := edas.CreateListSlbRequest()
	request.RegionId = regionId

	raw, err = edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.ListSlbResponse)
	if response.Code != 200 {
		return WrapError(Error("List Slb failed for " + response.Message))
	}

	for _, slb := range response.SlbList.SlbEntity {
		if slb.SlbId == slbId {
			d.Set("slb_status", slb.SlbStatus)
			d.Set("vswitch_id", slb.VswitchId)
			return nil
		}
	}

	return nil
}

func resourceAlicloudEdasSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateUnbindSlbRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Get("app_id").(string)
	request.SlbId = d.Get("slb_id").(string)
	request.Type = d.Get("type").(string)

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.UnbindSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.UnbindSlbResponse)
	if response.Code != 200 {
		return WrapError(Error("unbind slb failed," + response.Message))
	}

	return nil
}
