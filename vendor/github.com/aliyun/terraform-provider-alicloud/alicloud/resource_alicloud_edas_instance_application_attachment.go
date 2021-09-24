package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasInstanceApplicationAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasInstanceApplicationAttachmentCreate,
		Read:   resourceAlicloudEdasInstanceApplicationAttachmentRead,
		Delete: resourceAlicloudEdasInstanceApplicationAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ecc_info": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ecu_info": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasInstanceApplicationAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	ecuInfo := d.Get("ecu_info").([]interface{})
	aString := make([]string, len(ecuInfo))
	for i, v := range ecuInfo {
		aString[i] = v.(string)
	}

	request := edas.CreateScaleOutApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = appId
	request.DeployGroup = d.Get("deploy_group").(string)
	request.EcuInfo = strings.Join(aString, ",")

	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleOutApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_application_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ScaleOutApplicationResponse)
	changeOrderId = response.ChangeOrderId
	d.SetId(appId + ":" + strings.Join(aString, ","))
	if response.Code != 200 {
		return WrapError(Error("scaleOut application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasInstanceApplicationAttachmentRead(d, meta)
}

func resourceAlicloudEdasInstanceApplicationAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	appId := strs[0]
	regionId := client.RegionId
	ecuInfo := strs[1]
	aString := strings.Split(ecuInfo, ",")
	request := edas.CreateQueryApplicationStatusRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_application_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var eccs []string
	response := raw.(*edas.QueryApplicationStatusResponse)
	if response.Code != 200 {
		return WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}
	for _, ecc := range response.AppInfo.EccList.Ecc {
		for _, ecu := range aString {
			if ecu == ecc.EcuId {
				eccs = append(eccs, ecc.EccId)
			}
		}

	}

	d.Set("ecc_info", strings.Join(eccs, ","))
	return nil
}

func resourceAlicloudEdasInstanceApplicationAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateScaleInApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Get("app_id").(string)
	request.EccInfo = d.Get("ecc_info").(string)
	if v, ok := d.GetOk("force_status"); ok {
		request.ForceStatus = requests.NewBoolean(v.(bool))
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleInApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_application_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var changeOrderId string
	response, _ := raw.(*edas.ScaleInApplicationResponse)
	if response.Code != 200 {
		return WrapError(Error("scaleIn application failed for " + response.Message))
	}
	changeOrderId = response.ChangeOrderId

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}
