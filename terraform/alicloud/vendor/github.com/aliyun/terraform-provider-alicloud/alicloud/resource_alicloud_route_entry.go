package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteEntryCreate,
		Read:   resourceAliyunRouteEntryRead,
		Delete: resourceAliyunRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute router_id has been deprecated and suggest removing it from your template.",
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidrblock": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
		},
	}
}

func resourceAliyunRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	rtId := d.Get("route_table_id").(string)
	cidr := d.Get("destination_cidrblock").(string)
	nt := d.Get("nexthop_type").(string)
	ni := d.Get("nexthop_id").(string)

	table, err := vpcService.QueryRouteTableById(rtId)

	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateCreateRouteEntryRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = rtId
	request.DestinationCidrBlock = cidr
	request.NextHopType = nt
	request.NextHopId = ni
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RouteEntryName = d.Get("name").(string)

	// retry 10 min to create lots of entries concurrently
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
			return resource.NonRetryableError(err)
		}
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateRouteEntry(&args)
		})
		if err != nil {
			// Route Entry does not support concurrence when creating or deleting it;
			// Route Entry does not support creating or deleting within 5 seconds frequently
			// It must ensure all the route entries, vpc, vswitches' status must be available before creating or deleting route entry.
			if IsExpectedErrors(err, []string{"TaskConflict", "IncorrectRouteEntryStatus", Throttling, "IncorrectVpcStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"RouterEntryConflict.Duplicated"}) {
			en, err := vpcService.DescribeRouteEntry(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)
			if err != nil {
				return WrapError(err)
			}
			return WrapError(Error("The route entry %s has already existed. "+
				"Please import it using ID '%s:%s:%s:%s:%s' or specify a new 'destination_cidrblock' and try again.",
				en.DestinationCidrBlock, en.RouteTableId, table.VRouterId, en.DestinationCidrBlock, en.NextHopType, ni))
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	// route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id

	d.SetId(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)

	if err := vpcService.WaitForRouteEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunRouteEntryRead(d, meta)
}

func resourceAliyunRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	object, err := vpcService.DescribeRouteEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("router_id", parts[1])
	d.Set("route_table_id", object.RouteTableId)
	d.Set("destination_cidrblock", object.DestinationCidrBlock)
	d.Set("nexthop_type", object.NextHopType)
	d.Set("nexthop_id", object.InstanceId)
	d.Set("name", object.RouteEntryName)
	return nil
}

func resourceAliyunRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAliyunRouteEntryDeleteArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	request.RegionId = client.RegionId
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 5)
	rtId := parts[0]
	if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	retryTimes := 7
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectVpcStatus", "TaskConflict", "IncorrectRouteEntryStatus", "Forbbiden", "UnknownError"}) {
				// Route Entry does not support creating or deleting within 5 seconds frequently
				time.Sleep(time.Duration(retryTimes) * time.Second)
				retryTimes += 7
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRouteEntry.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForRouteEntry(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunRouteEntryDeleteArgs(d *schema.ResourceData, meta interface{}) (*vpc.DeleteRouteEntryRequest, error) {

	request := vpc.CreateDeleteRouteEntryRequest()
	request.RouteTableId = d.Get("route_table_id").(string)
	request.DestinationCidrBlock = d.Get("destination_cidrblock").(string)

	if v := d.Get("destination_cidrblock").(string); v != "" {
		request.DestinationCidrBlock = v
	}

	if v := d.Get("nexthop_id").(string); v != "" {
		request.NextHopId = v
	}

	return request, nil
}
