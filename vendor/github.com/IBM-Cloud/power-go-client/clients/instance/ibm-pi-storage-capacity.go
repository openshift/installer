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

//GetAll information about all the storage pools
func (f *IBMPIStorageCapacityClient) GetAll(powerinstanceid string, timeout time.Duration) (*models.StoragePoolsCapacity, error) {
	params := p_cloud_storage_capacity.NewPcloudStoragecapacityPoolsGetallParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudStorageCapacity.PcloudStoragecapacityPoolsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to get all storage pools %s", err)
	}
	return resp.Payload, nil
}
