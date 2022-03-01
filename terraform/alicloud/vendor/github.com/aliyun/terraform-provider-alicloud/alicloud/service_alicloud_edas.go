package alicloud

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EdasService struct {
	client *connectivity.AliyunClient
}

type Hook struct {
	Exec      *Exec      `json:"exec,omitempty"`
	HttpGet   *HttpGet   `json:"httpGet,omitempty"`
	TcpSocket *TcpSocket `json:"tcpSocket,omitempty"`
}

type Exec struct {
	Command []string `json:"command"`
}

type HttpGet struct {
	Path        string       `json:"path"`
	Port        int          `json:"port"`
	Scheme      string       `json:"scheme"`
	HttpHeaders []HttpHeader `json:"httpHeaders"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TcpSocket struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Prober struct {
	FailureThreshold    int `json:"failureThreshold"`
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	SuccessThreshold    int `json:"successThreshold"`
	TimeoutSeconds      int `json:"timeoutSeconds"`
	Hook                `json:",inline"`
}

func (e *EdasService) GetChangeOrderStatus(id string) (info *edas.ChangeOrderInfo, err error) {
	request := edas.CreateGetChangeOrderInfoRequest()
	request.RegionId = e.client.RegionId
	request.ChangeOrderId = id

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetChangeOrderInfo(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return info, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return info, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.GetChangeOrderInfoResponse)
	return &rsp.ChangeOrderInfo, nil

}

func (e *EdasService) GetDeployGroup(appId, groupId string) (groupInfo *edas.DeployGroup, err error) {
	request := edas.CreateListDeployGroupRequest()
	request.RegionId = e.client.RegionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	if err != nil {
		return groupInfo, WrapErrorf(err, DefaultErrorMsg, appId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.ListDeployGroupResponse)
	if rsp.Code != 200 {
		return groupInfo, Error("get deploy group failed for " + rsp.Message)
	}
	for _, group := range rsp.DeployGroupList.DeployGroup {
		if group.GroupId == groupId {
			return &group, nil
		}
	}
	return groupInfo, nil
}

func (e *EdasService) EdasChangeOrderStatusRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := e.GetChangeOrderStatus(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if strconv.Itoa(object.Status) == failState {
				return object, strconv.Itoa(object.Status), WrapError(Error(FailedToReachTargetStatus, strconv.Itoa(object.Status)))
			}
		}

		return object, strconv.Itoa(object.Status), nil
	}
}

func (e *EdasService) SyncResource(resourceType string) error {
	request := edas.CreateSynchronizeResourceRequest()
	request.RegionId = e.client.RegionId
	request.Type = resourceType

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.SynchronizeResource(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "sync resource", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.SynchronizeResourceResponse)
	if rsp.Code != 200 || !rsp.Success {
		return WrapError(Error("sync resource failed for " + rsp.Message))
	}

	return nil
}

func (e *EdasService) CheckEcsStatus(instanceIds string, count int) error {
	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = e.client.RegionId
	request.Status = "Running"
	request.PageSize = requests.NewInteger(100)
	request.InstanceIds = instanceIds

	raw, err := e.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, instanceIds, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	rsp := raw.(*ecs.DescribeInstancesResponse)

	if len(rsp.Instances.Instance) != count {
		return WrapErrorf(Error("not enough instances"), DefaultErrorMsg, instanceIds, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

func (e *EdasService) GetLastPackgeVersion(appId, groupId string) (string, error) {
	var versionId string
	request := edas.CreateQueryApplicationStatusRequest()
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application_package_version", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.QueryApplicationStatusResponse)

	if response.Code != 200 {
		return "", WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, group := range response.AppInfo.GroupList.Group {
		if group.GroupId == groupId {
			versionId = group.PackageVersionId
		}
	}

	rq := edas.CreateListHistoryDeployVersionRequest()
	rq.AppId = appId

	raw, err = e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListHistoryDeployVersion(rq)
	})
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application_package_version_list", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp, _ := raw.(*edas.ListHistoryDeployVersionResponse)

	if rsp.Code != 200 {
		return "", WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, version := range rsp.PackageVersionList.PackageVersion {
		if version.Id == versionId {
			return version.PackageVersion, nil
		}
	}

	return "", nil
}

func (e *EdasService) DescribeEdasApplication(appId string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	regionId := e.client.RegionId

	request := edas.CreateGetApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})
	if err != nil {
		return application, WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetApplicationResponse)
	if response.Code != 200 {
		return application, WrapError(Error("get application error :" + response.Message))
	}

	v := response.Applcation

	return &v, nil
}

func (e *EdasService) DescribeEdasCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	regionId := e.client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return cluster, WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		return cluster, WrapError(Error("create cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
}

func (e *EdasService) DescribeEdasDeployGroup(id string) (*edas.DeployGroup, error) {
	group := &edas.DeployGroup{}
	regionId := e.client.RegionId

	strs := strings.Split(id, ":")

	request := edas.CreateListDeployGroupRequest()
	request.RegionId = regionId
	request.AppId = strs[0]

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	if err != nil {
		return group, WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_deploy_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListDeployGroupResponse)
	if response.Code != 200 {
		return group, WrapError(Error("create cluster failed for " + response.Message))
	}

	for _, v := range response.DeployGroupList.DeployGroup {
		if v.ClusterName == strs[1] {
			return &v, nil
		}
	}

	return group, nil
}

func (e *EdasService) DescribeEdasInstanceClusterAttachment(id string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasCluster(v[0])
	if err != nil {
		return cluster, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationDeployment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationScale(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasSlbAttachment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

type CommandArg struct {
	Argument string `json:"argument" xml:"argument"`
}

func (e *EdasService) GetK8sCommandArgs(args []interface{}) (string, error) {
	aString := make([]CommandArg, 0)
	for _, v := range args {
		aString = append(aString, CommandArg{Argument: v.(string)})
	}
	b, err := json.Marshal(aString)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) GetK8sCommandArgsForDeploy(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

type K8sEnv struct {
	Name  string `json:"name" xml:"name"`
	Value string `json:"value" xml:"value"`
}

func (e *EdasService) GetK8sEnvs(envs map[string]interface{}) (string, error) {
	k8sEnvs := make([]K8sEnv, 0)
	for n, v := range envs {
		k8sEnvs = append(k8sEnvs, K8sEnv{Name: n, Value: v.(string)})
	}

	b, err := json.Marshal(k8sEnvs)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) QueryK8sAppPackageType(appId string) (string, error) {
	request := edas.CreateGetApplicationRequest()
	request.RegionId = e.client.RegionId
	request.AppId = appId
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})

	if err != nil {
		return "", WrapError(err)
	}

	response, _ := raw.(*edas.GetApplicationResponse)
	if response.Code != 200 {
		return "", WrapError(Error("get application for appId:" + appId + " failed:" + response.Message))
	}
	if len(response.Applcation.ApplicationType) > 0 {
		return response.Applcation.ApplicationType, nil
	}
	return "", WrapError(Error("not package type for appId:" + appId))
}

func (e *EdasService) DescribeEdasK8sCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	regionId := e.client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return cluster, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		if strings.Contains(response.Message, "does not exist") {
			return cluster, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return cluster, WrapError(Error("create k8s cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
}

func (e *EdasService) DescribeEdasK8sApplication(appId string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	regionId := e.client.RegionId

	request := edas.CreateGetK8sApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetK8sApplication(request)
	})
	if err != nil {
		return application, WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetK8sApplicationResponse)
	if response.Code != 200 {
		if strings.Contains(response.Message, "does not exist") {
			return application, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return application, WrapError(Error("get k8s application error :" + response.Message))
	}

	v := response.Applcation

	return &v, nil
}

func (e *EdasService) PreStopEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) PostStartEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) LivenessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}

func (e *EdasService) ReadinessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}
