package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenRouteMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenRouteMapCreate,
		Read:   resourceAlicloudCenRouteMapRead,
		Update: resourceAlicloudCenRouteMapUpdate,
		Delete: resourceAlicloudCenRouteMapDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"as_path_match_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Include", "Complete"}, false),
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_match_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Include", "Complete"}, false),
			},
			"community_match_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Include", "Complete", "Contain"}, false),
			},
			"community_operate_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Additive", "Replace"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_child_instance_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_cidr_blocks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_instance_ids_reverse_match": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"destination_route_table_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"map_result": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Permit", "Deny"}, false),
			},
			"match_asns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"match_community_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"next_priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"operate_community_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"prepend_as_path": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"route_map_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"route_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_child_instance_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_instance_ids_reverse_match": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"source_region_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_route_table_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transmit_direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RegionIn", "RegionOut"}, false),
			},
		},
	}
}

func resourceAlicloudCenRouteMapCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateCreateCenRouteMapRequest()
	if v, ok := d.GetOk("as_path_match_mode"); ok {
		request.AsPathMatchMode = v.(string)
	}
	request.CenId = d.Get("cen_id").(string)
	request.CenRegionId = d.Get("cen_region_id").(string)
	if v, ok := d.GetOk("cidr_match_mode"); ok {
		request.CidrMatchMode = v.(string)
	}
	if v, ok := d.GetOk("community_match_mode"); ok {
		request.CommunityMatchMode = v.(string)
	}
	if v, ok := d.GetOk("community_operate_mode"); ok {
		request.CommunityOperateMode = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("destination_child_instance_types"); ok {
		destinationChildInstanceTypes := expandStringList(v.(*schema.Set).List())
		request.DestinationChildInstanceTypes = &destinationChildInstanceTypes
	}
	if v, ok := d.GetOk("destination_cidr_blocks"); ok {
		destinationCidrBlocks := expandStringList(v.(*schema.Set).List())
		request.DestinationCidrBlocks = &destinationCidrBlocks
	}
	if v, ok := d.GetOk("destination_instance_ids"); ok {
		destinationInstanceIds := expandStringList(v.(*schema.Set).List())
		request.DestinationInstanceIds = &destinationInstanceIds
	}
	if v, ok := d.GetOkExists("destination_instance_ids_reverse_match"); ok {
		request.DestinationInstanceIdsReverseMatch = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("destination_route_table_ids"); ok {
		destinationRouteTableIds := expandStringList(v.(*schema.Set).List())
		request.DestinationRouteTableIds = &destinationRouteTableIds
	}
	request.MapResult = d.Get("map_result").(string)
	if v, ok := d.GetOk("match_asns"); ok {
		matchAsns := expandStringList(v.(*schema.Set).List())
		request.MatchAsns = &matchAsns
	}
	if v, ok := d.GetOk("match_community_set"); ok {
		matchCommunitySet := expandStringList(v.(*schema.Set).List())
		request.MatchCommunitySet = &matchCommunitySet
	}
	if v, ok := d.GetOk("next_priority"); ok {
		request.NextPriority = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("operate_community_set"); ok {
		operateCommunitySet := expandStringList(v.(*schema.Set).List())
		request.OperateCommunitySet = &operateCommunitySet
	}
	if v, ok := d.GetOk("preference"); ok {
		request.Preference = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("prepend_as_path"); ok {
		prependAsPath := expandStringList(v.(*schema.Set).List())
		request.PrependAsPath = &prependAsPath
	}
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	if v, ok := d.GetOk("route_types"); ok {
		routeTypes := expandStringList(v.(*schema.Set).List())
		request.RouteTypes = &routeTypes
	}
	if v, ok := d.GetOk("source_child_instance_types"); ok {
		sourceChildInstanceTypes := expandStringList(v.(*schema.Set).List())
		request.SourceChildInstanceTypes = &sourceChildInstanceTypes
	}
	if v, ok := d.GetOk("source_instance_ids"); ok {
		sourceInstanceIds := expandStringList(v.(*schema.Set).List())
		request.SourceInstanceIds = &sourceInstanceIds
	}
	if v, ok := d.GetOkExists("source_instance_ids_reverse_match"); ok {
		request.SourceInstanceIdsReverseMatch = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("source_region_ids"); ok {
		sourceRegionIds := expandStringList(v.(*schema.Set).List())
		request.SourceRegionIds = &sourceRegionIds
	}
	if v, ok := d.GetOk("source_route_table_ids"); ok {
		sourceRouteTableIds := expandStringList(v.(*schema.Set).List())
		request.SourceRouteTableIds = &sourceRouteTableIds
	}
	request.TransmitDirection = d.Get("transmit_direction").(string)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCenRouteMap(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.CreateCenRouteMapResponse)
		d.SetId(fmt.Sprintf("%v:%v", request.CenId, response.RouteMapId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_route_map", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenRouteMapStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenRouteMapRead(d, meta)
}
func resourceAlicloudCenRouteMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenRouteMap(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cen_id", parts[0])
	d.Set("as_path_match_mode", object.AsPathMatchMode)
	d.Set("cen_region_id", object.CenRegionId)
	d.Set("cidr_match_mode", object.CidrMatchMode)
	d.Set("community_match_mode", object.CommunityMatchMode)
	d.Set("community_operate_mode", object.CommunityOperateMode)
	d.Set("description", object.Description)
	d.Set("destination_child_instance_types", object.DestinationChildInstanceTypes.DestinationChildInstanceType)
	d.Set("destination_cidr_blocks", object.DestinationCidrBlocks.DestinationCidrBlock)
	d.Set("destination_instance_ids", object.DestinationInstanceIds.DestinationInstanceId)
	d.Set("destination_instance_ids_reverse_match", object.DestinationInstanceIdsReverseMatch)
	d.Set("destination_route_table_ids", object.DestinationRouteTableIds.DestinationRouteTableId)
	d.Set("map_result", object.MapResult)
	d.Set("match_asns", object.MatchAsns.MatchAsn)
	d.Set("match_community_set", object.MatchCommunitySet.MatchCommunity)
	d.Set("next_priority", object.NextPriority)
	d.Set("operate_community_set", object.OperateCommunitySet.OperateCommunity)
	d.Set("preference", object.Preference)
	d.Set("prepend_as_path", object.PrependAsPath.AsPath)
	d.Set("priority", object.Priority)
	d.Set("route_types", object.RouteTypes.RouteType)
	d.Set("source_child_instance_types", object.SourceChildInstanceTypes.SourceChildInstanceType)
	d.Set("source_instance_ids", object.SourceInstanceIds.SourceInstanceId)
	d.Set("source_instance_ids_reverse_match", object.SourceInstanceIdsReverseMatch)
	d.Set("source_region_ids", object.SourceRegionIds.SourceRegionId)
	d.Set("source_route_table_ids", object.SourceRouteTableIds.SourceRouteTableId)
	d.Set("status", object.Status)
	d.Set("transmit_direction", object.TransmitDirection)
	return nil
}
func resourceAlicloudCenRouteMapUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := cbn.CreateModifyCenRouteMapRequest()
	request.CenId = parts[0]
	request.RouteMapId = parts[1]
	if d.HasChange("cen_region_id") {
		update = true
	}
	request.CenRegionId = d.Get("cen_region_id").(string)
	if d.HasChange("map_result") {
		update = true
	}
	request.MapResult = d.Get("map_result").(string)
	if d.HasChange("priority") {
		update = true
	}
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	if d.HasChange("as_path_match_mode") {
		update = true
		request.AsPathMatchMode = d.Get("as_path_match_mode").(string)
	}
	if d.HasChange("cidr_match_mode") {
		update = true
		request.CidrMatchMode = d.Get("cidr_match_mode").(string)
	}
	if d.HasChange("community_match_mode") {
		update = true
		request.CommunityMatchMode = d.Get("community_match_mode").(string)
	}
	if d.HasChange("community_operate_mode") {
		update = true
		request.CommunityOperateMode = d.Get("community_operate_mode").(string)
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if d.HasChange("destination_child_instance_types") {
		update = true
		destinationChildInstanceTypes := expandStringList(d.Get("destination_child_instance_types").(*schema.Set).List())
		request.DestinationChildInstanceTypes = &destinationChildInstanceTypes

	}
	if d.HasChange("destination_cidr_blocks") {
		update = true
		destinationCidrBlocks := expandStringList(d.Get("destination_cidr_blocks").(*schema.Set).List())
		request.DestinationCidrBlocks = &destinationCidrBlocks

	}
	if d.HasChange("destination_instance_ids") {
		update = true
		destinationInstanceIds := expandStringList(d.Get("destination_instance_ids").(*schema.Set).List())
		request.DestinationInstanceIds = &destinationInstanceIds

	}
	if d.HasChange("destination_instance_ids_reverse_match") {
		update = true
		request.DestinationInstanceIdsReverseMatch = requests.NewBoolean(d.Get("destination_instance_ids_reverse_match").(bool))
	}
	if d.HasChange("destination_route_table_ids") {
		update = true
		destinationRouteTableIds := expandStringList(d.Get("destination_route_table_ids").(*schema.Set).List())
		request.DestinationRouteTableIds = &destinationRouteTableIds

	}
	if d.HasChange("match_asns") {
		update = true
		matchAsns := expandStringList(d.Get("match_asns").(*schema.Set).List())
		request.MatchAsns = &matchAsns

	}
	if d.HasChange("match_community_set") {
		update = true
		matchCommunitySet := expandStringList(d.Get("match_community_set").(*schema.Set).List())
		request.MatchCommunitySet = &matchCommunitySet

	}
	if d.HasChange("next_priority") {
		update = true
		request.NextPriority = requests.NewInteger(d.Get("next_priority").(int))
	}
	if d.HasChange("operate_community_set") {
		update = true
		operateCommunitySet := expandStringList(d.Get("operate_community_set").(*schema.Set).List())
		request.OperateCommunitySet = &operateCommunitySet

	}
	if d.HasChange("preference") {
		update = true
		request.Preference = requests.NewInteger(d.Get("preference").(int))
	}
	if d.HasChange("prepend_as_path") {
		update = true
		prependAsPath := expandStringList(d.Get("prepend_as_path").(*schema.Set).List())
		request.PrependAsPath = &prependAsPath

	}
	if d.HasChange("route_types") {
		update = true
		routeTypes := expandStringList(d.Get("route_types").(*schema.Set).List())
		request.RouteTypes = &routeTypes

	}
	if d.HasChange("source_child_instance_types") {
		update = true
		sourceChildInstanceTypes := expandStringList(d.Get("source_child_instance_types").(*schema.Set).List())
		request.SourceChildInstanceTypes = &sourceChildInstanceTypes

	}
	if d.HasChange("source_instance_ids") {
		update = true
		sourceInstanceIds := expandStringList(d.Get("source_instance_ids").(*schema.Set).List())
		request.SourceInstanceIds = &sourceInstanceIds

	}
	if d.HasChange("source_instance_ids_reverse_match") {
		update = true
		request.SourceInstanceIdsReverseMatch = requests.NewBoolean(d.Get("source_instance_ids_reverse_match").(bool))
	}
	if d.HasChange("source_region_ids") {
		update = true
		sourceRegionIds := expandStringList(d.Get("source_region_ids").(*schema.Set).List())
		request.SourceRegionIds = &sourceRegionIds

	}
	if d.HasChange("source_route_table_ids") {
		update = true
		sourceRouteTableIds := expandStringList(d.Get("source_route_table_ids").(*schema.Set).List())
		request.SourceRouteTableIds = &sourceRouteTableIds

	}
	if update {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenRouteMap(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCenRouteMapRead(d, meta)
}
func resourceAlicloudCenRouteMapDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := cbn.CreateDeleteCenRouteMapRequest()
	request.CenId = parts[0]
	request.RouteMapId = parts[1]
	request.CenRegionId = d.Get("cen_region_id").(string)
	err = resource.Retry(300*time.Second, func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenRouteMap(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.RouteMap", "Throttling.User"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
