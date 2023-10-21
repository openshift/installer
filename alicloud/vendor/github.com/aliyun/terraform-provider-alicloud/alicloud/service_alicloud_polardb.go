package alicloud

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type PolarDBService struct {
	client *connectivity.AliyunClient
}

func (s *PolarDBService) DescribePolarDBCluster(id string) (instance *polardb.DBCluster, err error) {
	request := polardb.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusters(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClustersResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *PolarDBService) DescribePolarDBClusterAttribute(id string) (instance *polardb.DescribeDBClusterAttributeResponse, err error) {
	request := polardb.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClusterAttributeResponse)

	return response, nil
}

func (s *PolarDBService) DescribePolarDBAutoRenewAttribute(id string) (instance *polardb.AutoRenewAttribute, err error) {
	request := polardb.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *PolarDBService) DescribeParameters(id string) (ds *polardb.DescribeDBClusterParametersResponse, err error) {
	request := polardb.CreateDescribeDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterParameters(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterParametersResponse)
	return response, err
}

func (s *PolarDBService) GrantPolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateGrantAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName
	request.AccountPrivilege = parts[2]

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDBService) RevokePolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateRevokeAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.RevokeAccountPrivilege(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDBService) WaitForPolarDBAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			return WrapError(err)
		}

		exist := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}

	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		ready := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName && dp.AccountPrivilege == parts[2] {
					ready = true
					break
				}
			}
		}
		if status == Deleted && !ready {
			break
		}
		if status != Deleted && ready {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribePolarDBAccountPrivilege(id string) (account *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccountPrivilege", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.Accounts[0], nil
}

func (s *PolarDBService) WaitForPolarDBConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *PolarDBService) WaitPolardbEndpointConfigEffect(id string, item map[string]string, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		effected := true
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapError(err)
		}

		var endpoint polardb.DBEndpoint
		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					endpoint = o
					break
				}
			}
		}
		if value, ok := item["Nodes"]; ok {
			if endpoint.Nodes != value {
				effected = false
			}
		}
		if value, ok := item["ReadWriteMode"]; ok {
			if endpoint.ReadWriteMode != value {
				effected = false
			}
		}
		if value, ok := item["AutoAddNewNodes"]; ok {
			if endpoint.AutoAddNewNodes != value {
				effected = false
			}
		}
		if value, ok := item["EndpointConfig"]; ok {
			expectConfig := make(map[string]string)
			actualConfig := make(map[string]string)
			err = json.Unmarshal([]byte(value), &expectConfig)
			err = json.Unmarshal([]byte(endpoint.EndpointConfig), &actualConfig)
			for k, v := range expectConfig {
				if subVal, ok := actualConfig[k]; !ok || subVal != v {
					effected = false
					break
				}
			}
		}
		if effected {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, endpoint, item, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *PolarDBService) WaitForPolarDBEndpoints(d *schema.ResourceData, status Status, endpointIds *schema.Set, timeout int) (string, error) {
	var dbEndpointId string
	if d.Id() != "" {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return "", WrapError(err)
		}
		dbEndpointId = parts[1]
	}
	dbClusterId := d.Get("db_cluster_id").(string)
	endpointType := d.Get("endpoint_type").(string)

	newEndpoint := make(map[string]string)
	newEndpoint["endpoint_type"] = endpointType

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		endpoints, err := s.DescribePolarDBInstanceNetInfo(dbClusterId)
		if err != nil {
			return "", WrapError(err)
		}

		var deleted bool
		deleted = true
		for _, value := range endpoints {
			if status == Deleted {
				if dbEndpointId == value.DBEndpointId {
					deleted = false
				}
				continue
			}
			if !endpointIds.Contains(value.DBEndpointId) && value.EndpointType == endpointType {
				return value.DBEndpointId, nil
			}
		}
		if status == Deleted && deleted {
			return "", nil
		}

		if time.Now().After(deadline) {
			return "", WrapErrorf(err, WaitTimeoutMsg, dbClusterId, GetFunc(1), timeout, endpoints, newEndpoint, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return "", nil
}

func (s *PolarDBService) DescribePolarDBConnection(id string) (*polardb.Address, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if NotFoundError(err) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return nil, WrapError(err)
		}

		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					for _, p := range o.AddressItems {
						if p.NetType == "Public" {
							return &p, nil
						}
					}
				}
			}
		}
		time.Sleep(DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("DBConnection", id)), NotFoundMsg, ProviderERROR)
}

func (s *PolarDBService) DescribePolarDBInstanceNetInfo(id string) ([]polardb.DBEndpoint, error) {

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterEndpoints(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
	if len(response.Items) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBInstanceNetInfo", id)), NotFoundMsg, ProviderERROR)
	}

	return response.Items, nil
}

func (s *PolarDBService) DescribePolarDBClusterEndpoint(id string) (*polardb.DBEndpoint, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = dbClusterId
	request.DBEndpointId = dbEndpointId

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterEndpoints(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
	if len(response.Items) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBEndpoint", dbEndpointId)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items[0], nil
}

func (s *PolarDBService) DescribePolarDBClusterSSL(d *schema.ResourceData) (ssl *polardb.DescribeDBClusterSSLResponse, err error) {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return nil, WrapError(err)
	}
	dbClusterId := parts[0]

	request := polardb.CreateDescribeDBClusterSSLRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = dbClusterId
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterSSL(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, _ := raw.(*polardb.DescribeDBClusterSSLResponse)

	return response, nil
}

func (s *PolarDBService) DescribePolarDBDatabase(id string) (ds *polardb.Database, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	dbName := parts[1]
	request := polardb.CreateDescribeDatabasesRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = dbName

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDatabases(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDatabasesResponse)
	if len(response.Databases.Database) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBDatabase", dbName)), NotFoundMsg, ProviderERROR)
	}
	ds = &response.Databases.Database[0]
	return ds, nil
}

func (s *PolarDBService) WaitForPolarDBDatabase(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBDatabase(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if status != Deleted && object != nil && object.DBName == parts[1] {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBName, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccount(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribePolarDBAccount(id string) (ds *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccount", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.Accounts[0], nil
}

// WaitForInstance waits for instance to given status
func (s *PolarDBService) WaitForPolarDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBConnectionPrefix(id, prefix string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBConnection(id)
		if err != nil {
			return WrapError(err)
		}
		parts := strings.Split(object.ConnectionString, ".")
		if prefix == parts[0] {
			break
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, parts[0], prefix, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) RefreshEndpointConfig(d *schema.ResourceData) error {
	var config map[string]interface{}
	config = make(map[string]interface{}, 0)
	documented, ok := d.GetOk("endpoint_config")
	if !ok {
		d.Set("endpoint_config", config)
		return nil
	}
	object, err := s.DescribePolarDBClusterEndpoint(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var endpointConfig = make(map[string]interface{})
	err = json.Unmarshal([]byte(object.EndpointConfig), &endpointConfig)
	if err != nil {
		return WrapError(err)
	}

	for k, v := range documented.(map[string]interface{}) {
		if _, ok := endpointConfig[k]; ok {
			config[k] = v
		}
	}
	if err := d.Set("endpoint_config", config); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) RefreshParameters(d *schema.ResourceData) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk("parameters")
	if !ok {
		d.Set("parameters", param)
		return nil
	}
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if err := d.Set("parameters", param); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) ModifyParameters(d *schema.ResourceData) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]string)
	allConfig := make(map[string]string)
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterParameters(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		if err := s.WaitForPolarDBParameter(d.Id(), DefaultTimeoutMedium, allConfig); err != nil {
			return WrapError(err)
		}
	}
	d.SetPartial("parameters")
	return nil
}

func (s *PolarDBService) setClusterTags(d *schema.ResourceData) error {
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
			request := polardb.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := polardb.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
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

func (s *PolarDBService) diffTags(oldTags, newTags []polardb.TagResourcesTag) ([]polardb.TagResourcesTag, []polardb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []polardb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *PolarDBService) tagsToMap(tags []polardb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *PolarDBService) tagsFromMap(m map[string]interface{}) []polardb.TagResourcesTag {
	result := make([]polardb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, polardb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *PolarDBService) ignoreTag(t polardb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagValue)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *PolarDBService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []polardb.TagResource, err error) {
	request := polardb.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

// WaitForCluster waits for cluster to given status
func (s *PolarDBService) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribeDBSecurityIps(clusterId string, dbClusterIPArrayName string) (ips []string, err error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.DBClusterIPArray {
		if ip.DBClusterIPArrayName == dbClusterIPArrayName {
			ipstr += separator + ip.SecurityIps
			separator = COMMA_SEPARATED
		}
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *PolarDBService) ModifyDBSecurityIps(clusterId, ips string) error {

	request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterAccessWhitelist(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) DescribeBackupPolicy(id string) (policy *polardb.DescribeBackupPolicyResponse, err error) {

	request := polardb.CreateDescribeBackupPolicyRequest()
	request.DBClusterId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.DescribeBackupPolicy(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return policy, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return raw.(*polardb.DescribeBackupPolicyResponse), nil
}

func (s *PolarDBService) ModifyDBBackupPolicy(clusterId, backupTime, backupPeriod string) error {

	request := polardb.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime

	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.ModifyBackupPolicy(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) DescribeDBAuditLogCollectorStatus(id string) (collectorStatus string, err error) {
	request := polardb.CreateDescribeDBClusterAuditLogCollectorRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.DescribeDBClusterAuditLogCollector(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return "", WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return collectorStatus, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*polardb.DescribeDBClusterAuditLogCollectorResponse)

	return response.CollectorStatus, nil
}

func (s *PolarDBService) PolarDBClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, WrapError(Error(FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *PolarDBService) WaitForPolarDBParameter(clusterId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeParameters(clusterId)
		if err != nil {
			return WrapError(err)
		}

		var actuals = make(map[string]string)
		for _, i := range object.RunningParameters.Parameter {
			actuals[i.ParameterName] = i.ParameterValue
		}

		match := true

		got_value := ""
		expected_value := ""

		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					got_value = actual
					expected_value = expect
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, clusterId, GetFunc(1), timeout, got_value, expected_value, ProviderERROR)
		}
	}
	return nil
}

func (s *PolarDBService) DescribeDBClusterTDE(id string) (map[string]interface{}, error) {
	action := "DescribeDBClusterTDE"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewPolarDBClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *PolarDBService) WaitForPolarDBTDEStatus(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Minute)
	for {
		object, err := s.DescribeDBClusterTDE(id)
		if err != nil {
			return WrapError(err)
		}
		if object["TDEStatus"] == status {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBNodeClass(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribePolarDBCluster(id)
		if err != nil {
			return WrapError(err)
		}
		clusterAttribute, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			return WrapError(err)
		}
		if len(clusters.DBNodes.DBNode) == len(clusterAttribute.DBNodes) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, ProviderERROR)
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBPayType(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribePolarDBCluster(id)
		if err != nil {
			return WrapError(err)
		}
		if clusters.PayType == status {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, clusters, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *PolarDBService) PolarDBClusterTDEStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribeDBClusterTDE(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["TDEStatus"].(string) == failState {
				return object, object["TDEStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *PolarDBService) DescribeDBSecurityGroups(clusterId string) ([]string, error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	groups := make([]string, 0)
	dbClusterSecurityGroups := resp.DBClusterSecurityGroups.DBClusterSecurityGroup
	for _, group := range dbClusterSecurityGroups {

		groups = append(groups, group.SecurityGroupId)
	}
	return groups, nil
}

func (s *PolarDBService) ModifyDBAccessWhitelistSecurityIps(d *schema.ResourceData) error {
	if l, ok := d.GetOk("db_cluster_ip_array"); ok {
		for _, e := range l.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			//ips expand string list
			ipList := expandStringList(pack["security_ips"].(*schema.Set).List())
			ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
			// default disable connect from outside
			if ipstr == "" {
				ipstr = LOCAL_HOST_IP
			}
			request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
			request.RegionId = s.client.RegionId
			request.DBClusterId = d.Id()
			request.SecurityIps = ipstr
			request.DBClusterIPArrayName = pack["db_cluster_ip_array_name"].(string)
			request.ModifyMode = pack["modify_mode"].(string)
			raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBClusterAccessWhitelist(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err := s.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
		}
	}
	d.SetPartial("db_cluster_ip_array")
	return nil
}

func (s *PolarDBService) DBClusterIPArrays(d *schema.ResourceData, attribute string, dbClusterIPArrayName string) error {
	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	dbClusterIPArray := response.Items.DBClusterIPArray
	var dbClusterIPArrays = make([]map[string]interface{}, 0, len(dbClusterIPArray))
	for _, i := range dbClusterIPArray {
		if i.DBClusterIPArrayName == dbClusterIPArrayName {
			l := map[string]interface{}{
				"db_cluster_ip_array_name": i.DBClusterIPArrayName,
				"security_ips":             convertPolarDBIpsSetToString(i.SecurityIps),
			}
			dbClusterIPArrays = append(dbClusterIPArrays, l)
		}
	}
	if err := d.Set(attribute, dbClusterIPArrays); err != nil {
		return WrapError(err)
	}

	return nil
}

func convertPolarDBIpsSetToString(sourceIps string) []string {
	ipsMap := make(map[string]string)

	for _, ip := range strings.Split(sourceIps, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}
	var ips []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			ips = append(ips, key)
		}
	}
	return ips
}
