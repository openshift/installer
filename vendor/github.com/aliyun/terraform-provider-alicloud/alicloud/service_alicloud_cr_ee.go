package alicloud

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
)

func (c *CrService) ListCrEEInstances(pageNo int, pageSize int) (*cr_ee.ListInstanceResponse, error) {
	response := &cr_ee.ListInstanceResponse{}
	request := cr_ee.CreateListInstanceRequest()
	request.RegionId = c.client.RegionId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListInstance(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, "ListInstance", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListInstanceResponse)
	if !response.ListInstanceIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, "ListInstance", action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEEInstance(instanceId string) (*cr_ee.GetInstanceResponse, error) {
	response := &cr_ee.GetInstanceResponse{}
	request := cr_ee.CreateGetInstanceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	err := resource.Retry(6*time.Second, func() *resource.RetryError {
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
				time.Sleep(time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw, request.RpcRequest, request)
		response, _ = raw.(*cr_ee.GetInstanceResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	if !response.GetInstanceIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) GetCrEEInstanceUsage(instanceId string) (*cr_ee.GetInstanceUsageResponse, error) {
	response := &cr_ee.GetInstanceUsageResponse{}
	request := cr_ee.CreateGetInstanceUsageRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetInstanceUsage(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetInstanceUsageResponse)
	if !response.GetInstanceUsageIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) ListCrEEInstanceEndpoint(instanceId string) (*cr_ee.ListInstanceEndpointResponse, error) {
	response := &cr_ee.ListInstanceEndpointResponse{}
	request := cr_ee.CreateListInstanceEndpointRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListInstanceEndpoint(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.ListInstanceEndpointResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	if !response.ListInstanceEndpointIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) ListCrEENamespaces(instanceId string, pageNo int, pageSize int) (*cr_ee.ListNamespaceResponse, error) {
	response := &cr_ee.ListNamespaceResponse{}
	request := cr_ee.CreateListNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListNamespace(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListNamespaceResponse)
	if !response.ListNamespaceIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEENamespace(id string) (*cr_ee.GetNamespaceResponse, error) {
	response := &cr_ee.GetNamespaceResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	request := cr_ee.CreateGetNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetNamespace(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetNamespaceResponse)
	if !response.GetNamespaceIsSuccess {
		if response.Code == "NAMESPACE_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DeleteCrEENamespace(id string) (*cr_ee.DeleteNamespaceResponse, error) {
	response := &cr_ee.DeleteNamespaceResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	request := cr_ee.CreateDeleteNamespaceRequest()
	request.RegionId = c.client.RegionId

	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteNamespace(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteNamespaceResponse)
	if !response.DeleteNamespaceIsSuccess {
		if response.Code == "NAMESPACE_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) WaitForCrEENamespace(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		resp, err := c.DescribeCrEENamespace(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if resp.NamespaceName == parts[1] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, resp.NamespaceName, parts[1], ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *CrService) ListCrEERepos(instanceId string, namespace string, pageNo int, pageSize int) (*cr_ee.ListRepositoryResponse, error) {
	response := &cr_ee.ListRepositoryResponse{}
	request := cr_ee.CreateListRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoStatus = "ALL"
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepository(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", namespace), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListRepositoryResponse)
	if !response.ListRepositoryIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", namespace), action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEERepo(id string) (*cr_ee.GetRepositoryResponse, error) {
	response := &cr_ee.GetRepositoryResponse{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return response, WrapError(err)
	}
	request := cr_ee.CreateGetRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.RepoNamespaceName = parts[1]
	request.RepoName = parts[2]
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetRepository(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetRepositoryResponse)
	if !response.GetRepositoryIsSuccess {
		if response.Code == "REPO_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DeleteCrEERepo(id, repoId string) (*cr_ee.DeleteRepositoryResponse, error) {
	response := &cr_ee.DeleteRepositoryResponse{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return response, WrapError(err)
	}
	request := cr_ee.CreateDeleteRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.RepoId = repoId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteRepository(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteRepositoryResponse)
	if !response.DeleteRepositoryIsSuccess {
		if response.Code == "REPO_NOT_EXIST" {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(fmt.Errorf("%v", response), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) WaitForCrEERepo(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	for {
		resp, err := c.DescribeCrEERepo(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if resp.RepoName == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, resp.RepoName, parts[2], ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *CrService) ListCrEERepoTags(instanceId string, repoId string, pageNo int, pageSize int) (*cr_ee.ListRepoTagResponse, error) {
	response := &cr_ee.ListRepoTagResponse{}
	request := cr_ee.CreateListRepoTagRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoId = repoId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepoTag(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", repoId), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListRepoTagResponse)
	if !response.ListRepoTagIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", repoId), action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEESyncRule(id string) (*cr_ee.SyncRulesItem, error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}
	instanceId := parts[0]
	namespace := parts[1]
	syncRuleId := parts[2]

	pageNo := 1
	for {
		response := &cr_ee.ListRepoSyncRuleResponse{}
		request := cr_ee.CreateListRepoSyncRuleRequest()
		request.RegionId = c.client.RegionId
		request.InstanceId = instanceId
		request.NamespaceName = namespace
		request.PageNo = requests.NewInteger(pageNo)
		request.PageSize = requests.NewInteger(PageSizeLarge)
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListRepoSyncRule(request)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.ListRepoSyncRuleResponse)
		if !response.ListRepoSyncRuleIsSuccess {
			return nil, WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		for _, rule := range response.SyncRules {
			if rule.SyncRuleId == syncRuleId && rule.LocalInstanceId == instanceId {
				return &rule, nil
			}
		}

		if len(response.SyncRules) < PageSizeLarge {
			return nil, WrapErrorf(errors.New("sync rule not found"), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}

		pageNo++
	}
}
