package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	Hb_LAUNCHING            = "LAUNCHING"
	Hb_CREATING             = "CREATING"
	Hb_ACTIVATION           = "ACTIVATION"
	Hb_DELETING             = "DELETING"
	Hb_CREATE_FAILED        = "CREATE_FAILED"
	Hb_NODE_RESIZING        = "HBASE_SCALE_OUT"
	Hb_NODE_RESIZING_FAILED = "NODE_RESIZE_FAILED"
	Hb_DISK_RESIZING        = "HBASE_EXPANDING"
	Hb_DISK_RESIZE_FAILED   = "DISK_RESIZING_FAILED"
	Hb_LEVEL_MODIFY         = "INSTANCE_LEVEL_MODIFY"
	Hb_LEVEL_MODIFY_FAILED  = "INSTANCE_LEVEL_MODIFY_FAILED"
	Hb_HBASE_COLD_EXPANDING = "HBASE_COLD_EXPANDING"
)

type HBaseService struct {
	client *connectivity.AliyunClient
}

func (s *HBaseService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := hbase.CreateUnTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.UnTagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := hbase.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *HBaseService) diffTags(oldTags, newTags []hbase.TagResourcesTag) ([]hbase.TagResourcesTag, []hbase.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []hbase.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *HBaseService) tagsFromMap(m map[string]interface{}) []hbase.TagResourcesTag {
	result := make([]hbase.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, hbase.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *HBaseService) tagsToMap(tags []hbase.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *HBaseService) ignoreTag(t hbase.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *HBaseService) DescribeHBaseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbaseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstance"

	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Hbase:Instance", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

//pop has limit, support next.
func (s *HBaseService) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, _ := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) DescribeSecurityGroups(id string) (object hbase.DescribeSecurityGroupsResponse, err error) {
	request := hbase.CreateDescribeSecurityGroupsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeSecurityGroups(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeSecurityGroupsResponse)
	return *response, nil
}

func (s *HBaseService) HBaseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *HBaseService) ModifyClusterDeletionProtection(clusterId string, protection bool) error {
	request := hbase.CreateModifyClusterDeletionProtectionRequest()
	request.ClusterId = clusterId
	request.Protection = requests.NewBoolean(protection)
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.ModifyClusterDeletionProtection(request)
	})
	if err != nil {
		return WrapErrorf(err, clusterId+" modifyClusterDeletionProtection failed")
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *HBaseService) DescribeEndpoints(id string) (object hbase.DescribeEndpointsResponse, err error) {
	request := hbase.CreateDescribeEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeEndpoints(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeEndpointsResponse)
	return *response, nil
}

func (s *HBaseService) DescribeClusterConnection(id string) (object hbase.DescribeClusterConnectionResponse, err error) {
	request := hbase.CreateDescribeClusterConnectionRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeClusterConnection(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeClusterConnectionResponse)
	return *response, nil
}
