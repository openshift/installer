package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclCreate,
		Read:   resourceAlicloudSlbAclRead,
		Update: resourceAlicloudSlbAclUpdate,
		Delete: resourceAlicloudSlbAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipv4",
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudSlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateCreateAccessControlListRequest()
	request.RegionId = client.RegionId
	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}
	request.AclName = strings.TrimSpace(d.Get("name").(string))
	request.AddressIPVersion = d.Get("ip_version").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateAccessControlList(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_acl", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.CreateAccessControlListResponse)

	d.SetId(response.AclId)
	return resourceAlicloudSlbAclUpdate(d, meta)
}

func resourceAlicloudSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	tags, err := slbService.DescribeTags(d.Id(), nil, TagResourceAcl)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", slbService.tagsToMap(tags))

	object, err := slbService.DescribeSlbAcl(d.Id())
	if err != nil {
		if IsExpectedErrors(err, []string{"AclNotExist"}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.AclName)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("ip_version", object.AddressIPVersion)

	if err := d.Set("entry_list", slbService.FlattenSlbAclEntryMappings(object.AclEntrys.AclEntry)); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	d.Partial(true)

	if err := slbService.setInstanceTags(d, TagResourceAcl); err != nil {
		return WrapError(err)
	}

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetAccessControlListAttributeRequest()
		request.RegionId = client.RegionId
		request.AclId = d.Id()
		request.AclName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetAccessControlListAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			if err := slbService.SlbRemoveAccessControlListEntry(remove, d.Id()); err != nil {
				return WrapError(err)
			}
		}

		if len(add) > 0 {
			if err := slbService.SlbAddAccessControlListEntry(add, d.Id()); err != nil {
				return WrapError(err)
			}
		}
		d.SetPartial("entry_list")
	}

	d.Partial(false)

	return resourceAlicloudSlbAclRead(d, meta)
}

func resourceAlicloudSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteAccessControlListRequest()
	request.RegionId = client.RegionId
	request.AclId = d.Id()
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AclInUsed"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)

		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if !IsExpectedErrors(err, []string{"AclNotExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return WrapError(slbService.WaitForSlbAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
