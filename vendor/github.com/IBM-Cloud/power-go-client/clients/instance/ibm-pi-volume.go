package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIVolumeClient ..
type IBMPIVolumeClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

const (

	// Timeouts for power
	postTimeOut   = 30 * time.Second
	getTimeOut    = 60 * time.Second
	deleteTimeOut = 30 * time.Second
)

// NewIBMPIVolumeClient ...
func NewIBMPIVolumeClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIVolumeClient {
	return &IBMPIVolumeClient{
		sess, powerinstanceid,
	}
}

//Get information about a single volume only
func (f *IBMPIVolumeClient) Get(id, powerinstanceid string, timeout time.Duration) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumeID(id)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PI Volume %s :%s", id, err)
	}
	return resp.Payload, nil
}

//CreateVolumeV2 ...
func (f *IBMPIVolumeClient) CreateVolumeV2(createVolDefs *p_cloud_volumes.PcloudV2VolumesPostParams, powerinstanceid string, timeout time.Duration) (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(createVolDefs.Body)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Create PI Volume %s :%s", *createVolDefs.Body.Name, err)
	}
	return resp.Payload, nil
}

// CreateVolume ...
func (f *IBMPIVolumeClient) CreateVolume(createVolDefs *p_cloud_volumes.PcloudCloudinstancesVolumesPostParams, powerinstanceid string, timeout time.Duration) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(createVolDefs.Body)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Create PI Instance Volume %s :%s", *createVolDefs.Body.Name, err)
	}
	return resp.Payload, nil
}

// UpdateVolume ...
func (f *IBMPIVolumeClient) UpdateVolume(updateVolDefs *p_cloud_volumes.PcloudCloudinstancesVolumesPutParams, volumeid, powerinstanceid string, timeout time.Duration) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPutParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(updateVolDefs.Body).WithVolumeID(volumeid)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Update PI Instance Volume %s :%s", volumeid, err)
	}
	return resp.Payload, nil
}

// DeleteVolume ...
func (f *IBMPIVolumeClient) DeleteVolume(id string, powerinstanceid string, timeout time.Duration) error {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumeID(id)
	_, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf("Failed to Delete PI Instance Volume %s :%s", id, err)
	}
	return nil
}

//Create ..
// TO be Deprecated
func (f *IBMPIVolumeClient) Create(volumename string, volumesize float64, volumetype string, volumeshareable bool, powerinstanceid string, timeout time.Duration) (*models.Volume, error) {

	var body = models.CreateDataVolume{
		Name:      &volumename,
		Size:      &volumesize,
		DiskType:  volumetype,
		Shareable: &volumeshareable,
	}

	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(&body)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Create PI Instance Volume %s :%s", volumename, err)
	}
	return resp.Payload, nil
}

// Delete ...
func (f *IBMPIVolumeClient) Delete(id string, powerinstanceid string, timeout time.Duration) error {
	//var cloudinstanceid = f.session.PowerServiceInstance
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumeID(id)
	_, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf("Failed to Delete PI Instance Volume %s :%s", id, err)
	}
	return nil
}

// Update ...
func (f *IBMPIVolumeClient) Update(id, volumename string, volumesize float64, volumeshare bool, powerinstanceid string, timeout time.Duration) (*models.Volume, error) {

	var patchbody = models.UpdateVolume{}
	if &volumename != nil {
		patchbody.Name = &volumename
	}
	if &volumesize != nil {
		patchbody.Size = volumesize
	}
	if &volumeshare != nil {
		patchbody.Shareable = &volumeshare
	}

	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPutParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithVolumeID(id).WithBody(&patchbody)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Update PI Instance Volume %s :%s", id, err)
	}
	return resp.Payload, nil
}

// Attach a volume
func (f *IBMPIVolumeClient) Attach(id, volumename string, powerinstanceid string, timeout time.Duration) (models.Object, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id).WithVolumeID(volumename)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Attach PI Instance Volume %s :%s", id, err)
	}
	return resp.Payload, nil

}

//Detach a volume
func (f *IBMPIVolumeClient) Detach(id, volumename string, powerinstanceid string, timeout time.Duration) (models.Object, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesDeleteParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id).WithVolumeID(volumename)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		//return nil, errors.ToError(err)
		return nil, fmt.Errorf("Failed to detach the volume [%s ] for pvm instance with id [%s]: %s", volumename, id, err)
	}
	return resp.Payload, nil

}

// GetAll volumes part of an instance
func (f *IBMPIVolumeClient) GetAll(id, cloudInstanceID string, timeout time.Duration) (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesGetallParamsWithTimeout(helpers.PIGetTimeOut).WithPvmInstanceID(id).WithCloudInstanceID(cloudInstanceID)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesGetall(params, ibmpisession.NewAuth(f.session, cloudInstanceID))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get all PI Instance Volumes %s :%s", id, err)
	}
	return resp.Payload, nil

}

// SetBootVolume as the boot volume - PUT Operation
func (f *IBMPIVolumeClient) SetBootVolume(id, volumename, cloudInstanceID string, timeout time.Duration) (models.Object, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesSetbootPutParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(cloudInstanceID).WithPvmInstanceID(id).WithVolumeID(volumename)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesSetbootPut(params, ibmpisession.NewAuth(f.session, cloudInstanceID))
	if err != nil || resp == nil || resp.Payload == nil {
		//return nil, errors.ToError(err)
		return nil, fmt.Errorf("Failed to set the boot volume %s for cloud instance id [%s] ", volumename, cloudInstanceID)
	}
	return resp.Payload, nil
}

// CheckVolumeAttach if the volume is attached to the instance
func (f *IBMPIVolumeClient) CheckVolumeAttach(cloudInstanceID, pvmInstanceID, volumeID string, timeout time.Duration) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(cloudInstanceID).WithPvmInstanceID(pvmInstanceID).WithVolumeID(volumeID)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesGet(params, ibmpisession.NewAuth(f.session, cloudInstanceID))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to validate that the volume [%s] is attached to the pvminstance [%s]: %s", volumeID, pvmInstanceID, err)
	}
	return resp.Payload, nil
}
