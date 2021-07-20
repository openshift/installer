package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_storage_capacity"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIStorageCapacityClient ..
type IBMPIStorageCapacityClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPIStorageCapacityClient ...
func NewIBMPIStorageCapacityClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIStorageCapacityClient {
	return &IBMPIStorageCapacityClient{
		sess, powerinstanceid,
	}
}

//Storage capacity for all available storage pools in a region
func (f *IBMPIStorageCapacityClient) GetAllStoragePools(powerinstanceid string, timeout time.Duration) (*models.StoragePoolsCapacity, error) {
	params := p_cloud_storage_capacity.NewPcloudStoragecapacityPoolsGetallParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudStorageCapacity.PcloudStoragecapacityPoolsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to get all storage pools %v", err)
	}
	return resp.Payload, nil
}

// Storage capacity for a storage pool in a region
func (f *IBMPIStorageCapacityClient) GetAvailableStoragePool(powerinstanceid, storagepool string, timeout time.Duration) (*models.StoragePoolCapacity, error) {
	params := p_cloud_storage_capacity.NewPcloudStoragecapacityPoolsGetParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithStoragePoolName(storagepool)
	resp, err := f.session.Power.PCloudStorageCapacity.PcloudStoragecapacityPoolsGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		fmt.Errorf("failed to get the capacity for storage pool %v", err)
	}
	return resp.Payload, nil
}

// Storage capacity for a storage type in a region
func (f *IBMPIStorageCapacityClient) GetAvailableStorageCapacity(powerinstanceid, storage_tier string, timeout time.Duration) (*models.StorageTypeCapacity, error) {
	params := p_cloud_storage_capacity.NewPcloudStoragecapacityTypesGetParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithStorageTypeName(storage_tier)
	resp, err := f.session.Power.PCloudStorageCapacity.PcloudStoragecapacityTypesGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		fmt.Errorf("failed to get the capacity for storage pool %v", err)
	}
	return resp.Payload, nil
}

// Storage capacity for all available storage types in a region
func (f *IBMPIStorageCapacityClient) GetAvailableStorageType(powerinstanceid string, timeout time.Duration) (*models.StorageTypesCapacity, error) {
	params := p_cloud_storage_capacity.NewPcloudStoragecapacityTypesGetallParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudStorageCapacity.PcloudStoragecapacityTypesGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		fmt.Errorf("failed to get the capacity for storage tiers %v", err)
	}
	return resp.Payload, nil
}
