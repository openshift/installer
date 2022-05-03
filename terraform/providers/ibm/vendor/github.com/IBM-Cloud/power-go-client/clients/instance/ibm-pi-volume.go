package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/errors"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIVolumeClient ..
type IBMPIVolumeClient struct {
	IBMPIClient
}

// NewIBMPIVolumeClient ...
func NewIBMPIVolumeClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIVolumeClient {
	return &IBMPIVolumeClient{
		*NewIBMPIClient(ctx, sess, cloudInstanceID),
	}
}

//Get information about a single volume only
func (f *IBMPIVolumeClient) Get(id string) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesGetParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithVolumeID(id)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf(errors.GetVolumeOperationFailed, id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get Volume %s", id)
	}
	return resp.Payload, nil
}

// GetAll volumes
func (f *IBMPIVolumeClient) GetAll() (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get all Volumes for Cloud Instance %s: %w", f.cloudInstanceID, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all Volumes for Cloud Instance %s", f.cloudInstanceID)
	}
	return resp.Payload, nil
}

// GetAll volumes
func (f *IBMPIVolumeClient) GetAllAffinityVolumes(affinity string) (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithAffinity(&affinity)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get all Volumes with affinity %s for Cloud Instance %s: %w", affinity, f.cloudInstanceID, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all Volumes with affinity %s for Cloud Instance %s", affinity, f.cloudInstanceID)
	}
	return resp.Payload, nil
}

//CreateVolumeV2 ...
func (f *IBMPIVolumeClient) CreateVolumeV2(body *models.MultiVolumesCreate) (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudV2VolumesPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithBody(body)
	resp, err := f.session.Power.PCloudVolumes.PcloudV2VolumesPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf(errors.CreateVolumeV2OperationFailed, *body.Name, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Create Volume v2")
	}
	return resp.Payload, nil
}

// CreateVolume ...
func (f *IBMPIVolumeClient) CreateVolume(body *models.CreateDataVolume) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithBody(body)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf(errors.CreateVolumeOperationFailed, *body.Name, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Create Volume")
	}
	return resp.Payload, nil
}

// UpdateVolume ...
func (f *IBMPIVolumeClient) UpdateVolume(id string, body *models.UpdateVolume) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesPutParams().
		WithContext(f.ctx).WithTimeout(helpers.PIUpdateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithVolumeID(id).
		WithBody(body)
	resp, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesPut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf(errors.UpdateVolumeOperationFailed, id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Update Volume %s", id)
	}
	return resp.Payload, nil
}

// DeleteVolume ...
func (f *IBMPIVolumeClient) DeleteVolume(id string) error {
	params := p_cloud_volumes.NewPcloudCloudinstancesVolumesDeleteParams().
		WithContext(f.ctx).WithTimeout(helpers.PIDeleteTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithVolumeID(id)
	_, err := f.session.Power.PCloudVolumes.PcloudCloudinstancesVolumesDelete(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf(errors.DeleteVolumeOperationFailed, id, err)
	}
	return nil
}

// Attach a volume
func (f *IBMPIVolumeClient) Attach(id, volumename string) error {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithVolumeID(volumename)
	_, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf(errors.AttachVolumeOperationFailed, volumename, err)
	}
	return nil
}

//Detach a volume
func (f *IBMPIVolumeClient) Detach(id, volumename string) error {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesDeleteParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithVolumeID(volumename)
	_, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesDelete(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf(errors.DetachVolumeOperationFailed, volumename, err)
	}
	return nil
}

// GetAll volumes part of an instance
func (f *IBMPIVolumeClient) GetAllInstanceVolumes(id string) (*models.Volumes, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get all Volumes for PI Instance %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all Volumes for PI Instance %s", id)
	}
	return resp.Payload, nil
}

// SetBootVolume as the boot volume - PUT Operation
func (f *IBMPIVolumeClient) SetBootVolume(id, volumename string) error {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesSetbootPutParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithVolumeID(volumename)
	_, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesSetbootPut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to set the boot volume %s for instance %s", volumename, id)
	}
	return nil
}

// CheckVolumeAttach if the volume is attached to the instance
func (f *IBMPIVolumeClient) CheckVolumeAttach(id, volumeID string) (*models.Volume, error) {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesGetParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithVolumeID(volumeID)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to validate that the volume %s is attached to the pvminstance %s: %w", volumeID, id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to validate that the volume %s is attached to the pvminstance %s", volumeID, id)
	}
	return resp.Payload, nil
}

// UpdateVolumeAttach if the volume is attached to the instance
func (f *IBMPIVolumeClient) UpdateVolumeAttach(id, volumeID string, body *models.PVMInstanceVolumeUpdate) error {
	params := p_cloud_volumes.NewPcloudPvminstancesVolumesPutParams().
		WithContext(f.ctx).WithTimeout(helpers.PIUpdateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithVolumeID(volumeID).WithBody(body)
	resp, err := f.session.Power.PCloudVolumes.PcloudPvminstancesVolumesPut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to validate that the volume %s is attached to the pvminstance %s: %w", volumeID, id, err)
	}
	if resp == nil || resp.Payload == nil {
		return fmt.Errorf("failed to validate that the volume %s is attached to the pvminstance %s", volumeID, id)
	}
	return nil
}
