package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type EmrService struct {
	client *connectivity.AliyunClient
}

func (s *EmrService) DescribeEmrCluster(id string) (*emr.DescribeClusterV2Response, error) {
	response := &emr.DescribeClusterV2Response{}
	request := emr.CreateDescribeClusterV2Request()
	request.Id = id

	raw, err := s.client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.DescribeClusterV2(request)
	})

	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest)
	response, _ = raw.(*emr.DescribeClusterV2Response)
	if response.ClusterInfo.Status == "RELEASED" {
		return response, WrapErrorf(Error(GetNotFoundMessage("EmrCluster", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *EmrService) WaitForEmrCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.ClusterInfo.Id == id && status != Deleted {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterInfo.Id, id, ProviderERROR)
		}
	}
	return nil
}

func (s *EmrService) EmrClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.ClusterInfo.Status == failState {
				return object, object.ClusterInfo.Status, WrapError(Error(FailedToReachTargetStatus, object.ClusterInfo.Status))
			}
		}

		return object, object.ClusterInfo.Status, nil
	}
}

func (s *EmrService) setEmrClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := emr.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(TagResourceCluster)
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := emr.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(TagResourceCluster)
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *EmrService) DescribeEmrClusterTags(resourceId string, resourceType TagResourceType) (tags []emr.TagResource, err error) {
	request := emr.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*emr.ListTagResourcesResponse)
	tags = response.TagResources.TagResource

	return
}

func (s *EmrService) diffTags(oldTags, newTags []emr.TagResourcesTag) ([]emr.TagResourcesTag, []emr.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []emr.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *EmrService) tagsFromMap(m map[string]interface{}) []emr.TagResourcesTag {
	result := make([]emr.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, emr.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *EmrService) tagsToMap(tags []emr.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.TagKey] = t.TagValue
	}
	return result
}
