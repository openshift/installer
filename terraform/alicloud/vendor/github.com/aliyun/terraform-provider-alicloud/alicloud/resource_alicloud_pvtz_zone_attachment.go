package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPvtzZoneAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneAttachmentCreate,
		Read:   resourceAlicloudPvtzZoneAttachmentRead,
		Update: resourceAlicloudPvtzZoneAttachmentUpdate,
		Delete: resourceAlicloudPvtzZoneAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("vpcs") {
				d.SetNewComputed("vpc_ids")
			} else if d.HasChange("vpc_ids") {
				d.SetNewComputed("vpcs")
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpcs": {
				Type:          schema.TypeSet,
				Computed:      true,
				ConflictsWith: []string{"vpc_ids"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vpcs"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudPvtzZoneAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "BindZoneVpc"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}

	vpcIdMap := make(map[string]string)
	if v, ok := d.GetOk("vpcs"); ok {
		vpcs := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, j := range v.(*schema.Set).List() {
			vpcs[i] = make(map[string]interface{})
			regionId := j.(map[string]interface{})["region_id"]
			if regionId == "" {
				regionId = client.RegionId
			}
			vpcs[i]["RegionId"] = regionId
			vpcIdMap[j.(map[string]interface{})["vpc_id"].(string)] = j.(map[string]interface{})["vpc_id"].(string)
			vpcs[i]["VpcId"] = j.(map[string]interface{})["vpc_id"]
		}
		request["Vpcs"] = vpcs

	} else {
		vpcIds := d.Get("vpc_ids").(*schema.Set).List()
		vpcs := make([]map[string]interface{}, len(vpcIds))
		for i, j := range vpcIds {
			vpcs[i] = make(map[string]interface{})
			vpcs[i]["RegionId"] = client.RegionId
			vpcs[i]["VpcId"] = j.(string)
			vpcIdMap[j.(string)] = j.(string)
		}
		request["Vpcs"] = vpcs
	}

	request["ZoneId"] = d.Get("zone_id")
	wait := incrementalWait(3*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy", "Throttling.User", "ServiceUnavailable", "Zone.NotExists"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ZoneId"]))
	pvtzService := PvtzService{client}
	if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPvtzZoneAttachmentRead(d, meta)
}
func resourceAlicloudPvtzZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZoneAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_private_zone_zone_attachment pvtzService.DescribePvtzZoneAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("zone_id", d.Id())
	vpcIds := make([]string, 0)

	vpc := make([]map[string]interface{}, 0)
	if vpcList, ok := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{}); ok {
		for _, v := range vpcList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"region_id": m1["RegionId"],
					"vpc_id":    m1["VpcId"],
				}
				vpcIds = append(vpcIds, m1["VpcId"].(string))
				vpc = append(vpc, temp1)

			}
		}
	}
	if err := d.Set("vpcs", vpc); err != nil {
		return WrapError(err)
	}
	if err := d.Set("vpc_ids", vpcIds); err != nil {
		return WrapError(err)
	}
	return nil
}
func resourceAlicloudPvtzZoneAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	if d.HasChange("vpcs") || d.HasChange("vpc_ids") {
		request := map[string]interface{}{
			"ZoneId": d.Id(),
		}
		vpcIdMap := make(map[string]string)
		vpcs := make([]map[string]interface{}, len(d.Get("vpcs").(*schema.Set).List()))
		if d.HasChange("vpcs") && len(vpcs) != 0 {
			for i, j := range d.Get("vpcs").(*schema.Set).List() {
				vpcs[i] = make(map[string]interface{})
				regionId := j.(map[string]interface{})["region_id"]
				if regionId == "" {
					regionId = client.RegionId
				}
				vpcs[i]["RegionId"] = regionId
				vpcs[i]["VpcId"] = j.(map[string]interface{})["vpc_id"]
				vpcIdMap[j.(map[string]interface{})["vpc_id"].(string)] = j.(map[string]interface{})["vpc_id"].(string)
			}
		} else {
			vpcs = make([]map[string]interface{}, len(d.Get("vpc_ids").(*schema.Set).List()))
			for i, j := range d.Get("vpc_ids").(*schema.Set).List() {
				vpcs[i] = make(map[string]interface{})
				vpcs[i]["RegionId"] = client.RegionId
				vpcs[i]["VpcId"] = j.(string)
				vpcIdMap[j.(string)] = j.(string)
			}
		}
		request["Vpcs"] = vpcs
		if _, ok := d.GetOk("lang"); ok {
			request["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			request["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "BindZoneVpc"
		conn, err := client.NewPvtzClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 2*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy", "Throttling.User", "ServiceUnavailable", "Zone.NotExists"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		pvtzService := PvtzService{client}
		if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudPvtzZoneAttachmentRead(d, meta)
}
func resourceAlicloudPvtzZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "BindZoneVpc"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ZoneId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}
	vpcs := make([]map[string]interface{}, 0)
	request["Vpcs"] = vpcs
	wait := incrementalWait(3*time.Second, 2*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	pvtzService := PvtzService{client}
	return WrapError(pvtzService.WaitForPvtzZoneAttachment(d.Id(), Deleted, DefaultTimeout))
}
