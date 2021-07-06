package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPICloneVolumeClient ...
type IBMPICloneVolumeClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPICloneVolumeClient ...
func NewIBMPICloneVolumeClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPICloneVolumeClient {
	return &IBMPICloneVolumeClient{
		sess, powerinstanceid,
	}
}

//Create a clone volume using V2 of the API - This creates a clone
func (f *IBMPICloneVolumeClient) Create(cloneParams *p_cloud_volumes.PcloudV2VolumesClonePostParams, timeout time.Duration) (*models.CloneTaskReference, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesClonePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(cloneParams.CloudInstanceID).WithBody(cloneParams.Body)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesClonePost(params, ibmpisession.NewAuth(f.session, cloneParams.CloudInstanceID))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform Create operation... %s", err)
	}
	return resp.Payload, nil
}

//DeleteClone Deletes a clone
func (f *IBMPICloneVolumeClient) DeleteClone(cloneParams *p_cloud_volumes.PcloudV2VolumescloneDeleteParams, id, cloudinstance string, timeout time.Duration) (models.Object, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(cloudinstance).WithVolumesCloneID(id)

	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneDelete(params, ibmpisession.NewAuth(f.session, cloudinstance))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform Delete operation... %s", err)
	}
	return resp.Payload, nil
}

// Cancel a Clone

// Get status of a clone request
func (f *IBMPICloneVolumeClient) Get(powerinstanceid, clonetaskid string, timeout time.Duration) (*models.CloneTaskStatus, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesClonetasksGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithCloneTaskID(clonetaskid)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesClonetasksGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform the get operation for clones... %s", err)
	}
	return resp.Payload, nil
}

//StartClone ...
func (f *IBMPICloneVolumeClient) StartClone(powerinstanceid, volumeCloneID string, timeout time.Duration) (*models.VolumesClone, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneStartPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumesCloneID(volumeCloneID)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneStartPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform the start operation for clones... %s", err)
	}
	return resp.Payload, nil
}
