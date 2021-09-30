package instance

import (
	"fmt"
	"github.com/IBM-Cloud/power-go-client/errors"
	"log"
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
		return nil, fmt.Errorf(errors.CreateCloneOperationFailed, err)
	}
	return resp.Payload, nil
}

// Get status of a clone request
func (f *IBMPICloneVolumeClient) Get(powerinstanceid, clonetaskid string, timeout time.Duration) (*models.CloneTaskStatus, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesClonetasksGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithCloneTaskID(clonetaskid)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesClonetasksGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.GetCloneOperationFailed, powerinstanceid, err)
	}
	return resp.Payload, nil
}

// Create a volumes clone V2  Version = This is the prepare operation
func (f *IBMPICloneVolumeClient) CreateV2Clone(powerinstanceid string, cloneparams *p_cloud_volumes.PcloudV2VolumesclonePostParams) (*models.VolumesClone, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesclonePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(cloneparams.Body)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesclonePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.PrepareCloneOperationFailed, err)
	}
	return resp.Payload, nil

}

// Get a list of volume-clones request for a cloud instance
func (f *IBMPICloneVolumeClient) GetV2Clones(powerinstanceid, query_filter string) (*models.VolumesClones, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneGetallParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithFilter(&query_filter)
	log.Printf("the query filter is %s", query_filter)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.GetCloneOperationFailed, powerinstanceid, err)
	}
	return resp.Payload, nil
}

// Delete a volume- clone request
func (f *IBMPICloneVolumeClient) DeleteClone(cloneParams *p_cloud_volumes.PcloudV2VolumescloneDeleteParams, id, cloudinstance string, timeout time.Duration) (models.Object, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(cloudinstance).WithVolumesCloneID(id)

	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneDelete(params, ibmpisession.NewAuth(f.session, cloudinstance))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.DeleteCloneOperationFailed, err)
	}
	return resp.Payload, nil
}

// Initiate the start clone request

func (f *IBMPICloneVolumeClient) StartClone(powerinstanceid, volumeCloneID string, timeout time.Duration) (*models.VolumesClone, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneStartPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumesCloneID(volumeCloneID)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneStartPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.StartCloneOperationFailed, err)
	}
	return resp.Payload, nil
}

// Initiate the execute action for a clone request
func (f *IBMPICloneVolumeClient) PrepareClone(powerinstanceid, volumeCloneID string, timeout time.Duration) (*models.VolumesClone, error) {
	params := p_cloud_volumes.NewPcloudV2VolumescloneExecutePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumesCloneID(volumeCloneID)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneExecutePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.PrepareCloneOperationFailed, err)
	}
	return resp.Payload, nil
}

// Get V2Clone Task Status

func (f *IBMPICloneVolumeClient) GetV2CloneStatus(powerinstanceid, clone_name string) (*models.VolumesCloneDetail, error) {
	log.Printf("starting the clone status get operation..")
	params := p_cloud_volumes.NewPcloudV2VolumescloneGetParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumesCloneID(clone_name)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumescloneGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf(errors.GetCloneOperationFailed, powerinstanceid, err)
	}
	return resp.Payload, nil
}
