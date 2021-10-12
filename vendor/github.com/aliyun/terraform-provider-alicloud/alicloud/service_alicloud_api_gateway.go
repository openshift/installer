package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type CloudApiService struct {
	client *connectivity.AliyunClient
}

func (s *CloudApiService) DescribeApiGatewayGroup(id string) (*cloudapi.DescribeApiGroupResponse, error) {
	apiGroup := &cloudapi.DescribeApiGroupResponse{}
	request := cloudapi.CreateDescribeApiGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupId = id

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApiGroup"}) {
			return apiGroup, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return apiGroup, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	apiGroup, _ = raw.(*cloudapi.DescribeApiGroupResponse)
	if apiGroup.GroupId == "" {
		return apiGroup, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return apiGroup, nil
}
func (s *CloudApiService) WaitForApiGatewayGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(object.GroupId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, string(object.GroupId), id, ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApp(id string) (*cloudapi.DescribeAppResponse, error) {
	app := &cloudapi.DescribeAppResponse{}
	request := cloudapi.CreateDescribeAppRequest()
	request.RegionId = s.client.RegionId
	request.AppId = requests.Integer(id)

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApp(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return app, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return app, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	app, _ = raw.(*cloudapi.DescribeAppResponse)
	return app, nil
}

func (s *CloudApiService) WaitForApiGatewayApp(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApp(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if fmt.Sprint(object.AppId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object.AppId), id, ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApi(id string) (*cloudapi.DescribeApiResponse, error) {
	api := &cloudapi.DescribeApiResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return api, WrapError(err)
	}
	request := cloudapi.CreateDescribeApiRequest()
	request.RegionId = s.client.RegionId
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
			return api, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return api, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	api, _ = raw.(*cloudapi.DescribeApiResponse)
	if api.ApiId == "" {
		return api, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayApi", id)), NotFoundMsg, ProviderERROR)
	}
	return api, nil
}

func (s *CloudApiService) WaitForApiGatewayApi(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApi(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		respId := fmt.Sprintf("%s%s%s", object.GroupId, COLON_SEPARATED, object.ApiId)
		if respId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, respId, id, ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayAppAttachment(id string) (*cloudapi.AuthorizedApp, error) {
	app := &cloudapi.AuthorizedApp{}
	request := cloudapi.CreateDescribeAuthorizedAppsRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return app, WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.StageName = parts[3]
	appId, _ := strconv.ParseInt(parts[2], 10, 64)

	var allApps []cloudapi.AuthorizedApp

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAuthorizedApps(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
				return app, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return app, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeAuthorizedAppsResponse)

		allApps = append(allApps, response.AuthorizedApps.AuthorizedApp...)

		if len(allApps) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return app, WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredAppsTemp []cloudapi.AuthorizedApp
	for _, app := range allApps {
		if app.AppId == appId {
			filteredAppsTemp = append(filteredAppsTemp, app)
		}
	}

	if len(filteredAppsTemp) < 1 {
		return app, WrapErrorf(Error(GetNotFoundMessage("ApigatewayAppAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	app = &filteredAppsTemp[0]
	return app, nil
}

func (s *CloudApiService) DescribeApiGatewayVpcAccess(id string) (*cloudapi.VpcAccessAttribute, error) {
	vpc := &cloudapi.VpcAccessAttribute{}
	request := cloudapi.CreateDescribeVpcAccessesRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return vpc, WrapError(err)
	}
	var allVpcs []cloudapi.VpcAccessAttribute

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeVpcAccesses(request)
		})
		if err != nil {
			return vpc, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeVpcAccessesResponse)

		allVpcs = append(allVpcs, response.VpcAccessAttributes.VpcAccessAttribute...)

		if len(allVpcs) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return vpc, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredVpcsTemp []cloudapi.VpcAccessAttribute
	for _, vpc := range allVpcs {
		iPort, _ := strconv.Atoi(parts[3])
		if vpc.Port == iPort && vpc.InstanceId == parts[2] && vpc.VpcId == parts[1] && vpc.Name == parts[0] {
			filteredVpcsTemp = append(filteredVpcsTemp, vpc)
		}
	}

	if len(filteredVpcsTemp) < 1 {
		return vpc, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayVpcAccess", id)), NotFoundMsg, ProviderERROR)
	}
	return &filteredVpcsTemp[0], nil
}

func (s *CloudApiService) WaitForApiGatewayAppAttachment(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	appIds := parts[2]
	for {
		object, err := s.DescribeApiGatewayAppAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.FormatInt(object.AppId, 10) == appIds && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatInt(object.AppId, 10), appIds, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CloudApiService) DescribeDeployedApi(id string, stageName string) (*cloudapi.DescribeDeployedApiResponse, error) {
	api := &cloudapi.DescribeDeployedApiResponse{}
	request := cloudapi.CreateDescribeDeployedApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return api, WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeDeployedApi(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return api, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return api, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	api, _ = raw.(*cloudapi.DescribeDeployedApiResponse)
	return api, nil
}

func (s *CloudApiService) DeployedApi(id string, stageName string) (err error) {
	request := cloudapi.CreateDeployApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName
	request.Description = DeployCommonDescription

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeployApi(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) AbolishApi(id string, stageName string) (err error) {
	request := cloudapi.CreateAbolishApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.AbolishApi(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []cloudapi.TagResource, err error) {
	request := cloudapi.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []cloudapi.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, cloudapi.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}
	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cloudapi.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *CloudApiService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := cloudapi.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(resourceType)
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cloudapi.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(resourceType)
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")

	return nil
}

func (s *CloudApiService) tagsToMap(tags []cloudapi.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *CloudApiService) ignoreTag(t cloudapi.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *CloudApiService) diffTags(oldTags, newTags []cloudapi.TagResourcesTag) ([]cloudapi.TagResourcesTag, []cloudapi.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []cloudapi.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *CloudApiService) tagsFromMap(m map[string]interface{}) []cloudapi.TagResourcesTag {
	result := make([]cloudapi.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, cloudapi.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}
