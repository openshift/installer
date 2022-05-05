package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIInstanceClient ...
type IBMPIInstanceClient struct {
	IBMPIClient
}

// NewIBMPIInstanceClient ...
func NewIBMPIInstanceClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIInstanceClient {
	return &IBMPIInstanceClient{
		*NewIBMPIClient(ctx, sess, cloudInstanceID),
	}
}

//Get information about a single pvm only
func (f *IBMPIInstanceClient) Get(id string) (*models.PVMInstance, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesGetParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get PVM Instance %s :%w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get PVM Instance %s", id)
	}
	return resp.Payload, nil
}

// GetAll Information about all the PVM Instances for a Client
func (f *IBMPIInstanceClient) GetAll() (*models.PVMInstances, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get all PVM Instances of Power Instance %s :%w", f.cloudInstanceID, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all PVM Instances of Power Instance %s", f.cloudInstanceID)
	}
	return resp.Payload, nil
}

//Create ...
func (f *IBMPIInstanceClient) Create(body *models.PVMInstanceCreate) (*models.PVMInstanceList, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithBody(body)
	postok, postcreated, postAccepted, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Create PVM Instance :%w", err)
	}
	if postok != nil && len(postok.Payload) > 0 {
		return &postok.Payload, nil
	}
	if postcreated != nil && len(postcreated.Payload) > 0 {
		return &postcreated.Payload, nil
	}
	if postAccepted != nil && len(postAccepted.Payload) > 0 {
		return &postAccepted.Payload, nil
	}
	return nil, fmt.Errorf("failed to Create PVM Instance")
}

// Delete PVM Instances
func (f *IBMPIInstanceClient) Delete(id string) error {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesDeleteParams().
		WithContext(f.ctx).WithTimeout(helpers.PIDeleteTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	_, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesDelete(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to Delete PVM Instance %s :%w", id, err)
	}
	return nil
}

// Update PVM Instances
func (f *IBMPIInstanceClient) Update(id string, body *models.PVMInstanceUpdate) (*models.PVMInstanceUpdateResponse, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesPutParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).WithBody(body)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesPut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Update PVM Instance %s :%w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Update PVM Instance %s", id)
	}
	return resp.Payload, nil
}

// Action PVM Instances Operations
func (f *IBMPIInstanceClient) Action(id string, body *models.PVMInstanceAction) error {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesActionPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	_, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesActionPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to perform Action on PVM Instance %s :%w", id, err)
	}
	return nil

}

// PostConsoleURL Generate the Console URL
func (f *IBMPIInstanceClient) PostConsoleURL(id string) (*models.PVMInstanceConsole, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesConsolePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	postok, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesConsolePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Generate the Console URL PVM Instance %s :%w", id, err)
	}
	if postok == nil || postok.Payload == nil {
		return nil, fmt.Errorf("failed to Generate the Console URL PVM Instance %s", id)
	}
	return postok.Payload, nil
}

// List the available console languages for an instance
func (f *IBMPIInstanceClient) GetConsoleLanguages(id string) (*models.ConsoleLanguages, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesConsoleGetParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesConsoleGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get Console Languages for PVM Instance %s :%w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get Console Languages for PVM Instance %s", id)
	}
	return resp.Payload, nil
}

// List the available console languages for an instance
func (f *IBMPIInstanceClient) UpdateConsoleLanguage(id string, body *models.ConsoleLanguage) (*models.ConsoleLanguage, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesConsolePutParams().
		WithContext(f.ctx).WithTimeout(helpers.PIUpdateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesConsolePut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Update Console Language for PVM Instance %s :%w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Update Console Language for PVM Instance %s", id)
	}
	return resp.Payload, nil
}

// CaptureInstanceToImageCatalog Captures an instance
func (f *IBMPIInstanceClient) CaptureInstanceToImageCatalog(id string, body *models.PVMInstanceCapture) error {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesCapturePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	_, _, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesCapturePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to Capture the PVM Instance %s: %w", id, err)
	}
	return nil

}

//CaptureInstanceToImageCatalog Captures V2
func (f *IBMPIInstanceClient) CaptureInstanceToImageCatalogV2(id string, body *models.PVMInstanceCapture) (*models.JobReference, error) {
	params := p_cloud_p_vm_instances.NewPcloudV2PvminstancesCapturePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudV2PvminstancesCapturePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Capture the PVM Instance %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Capture the PVM Instance %s", id)
	}
	return resp.Payload, nil
}

// CreatePvmSnapShot Create a snapshot of the instance
func (f *IBMPIInstanceClient) CreatePvmSnapShot(id string, body *models.SnapshotCreate) (*models.SnapshotCreateResponse, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	snapshotpostaccepted, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesSnapshotsPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Create the snapshot for the pvminstance %s: %w", id, err)
	}
	if snapshotpostaccepted == nil || snapshotpostaccepted.Payload == nil {
		return nil, fmt.Errorf("failed to Create the snapshot for the pvminstance %s", id)
	}
	return snapshotpostaccepted.Payload, nil
}

// CreateClone ...
func (f *IBMPIInstanceClient) CreateClone(id string, body *models.PVMInstanceClone) (*models.PVMInstance, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesClonePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	clonePost, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesClonePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to create the clone of the pvm instance %s: %w", id, err)
	}
	if clonePost == nil || clonePost.Payload == nil {
		return nil, fmt.Errorf("failed to create the clone of the pvm instance %s", id)
	}
	return clonePost.Payload, nil
}

// GetSnapShotVM Get information about the snapshots for a vm
func (f *IBMPIInstanceClient) GetSnapShotVM(id string) (*models.Snapshots, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesSnapshotsGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get the snapshot for the pvminstance %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get the snapshot for the pvminstance %s", id)
	}
	return resp.Payload, nil

}

// RestoreSnapShotVM Restore a snapshot
func (f *IBMPIInstanceClient) RestoreSnapShotVM(id, snapshotid, restoreAction string, body *models.SnapshotRestore) (*models.Snapshot, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsRestorePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithSnapshotID(snapshotid).WithRestoreFailAction(&restoreAction).
		WithBody(body)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesSnapshotsRestorePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to restrore the snapshot for the pvminstance %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to restrore the snapshot for the pvminstance %s", id)
	}
	return resp.Payload, nil
}

// AddNetwork Add a network to the instance
func (f *IBMPIInstanceClient) AddNetwork(id string, body *models.PVMInstanceAddNetwork) (*models.PVMInstanceNetwork, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesNetworksPostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesNetworksPost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to attach the network to the pvminstanceid %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to attach the network to the pvminstanceid %s", id)
	}
	return resp.Payload, nil
}

// Delete a network from an instance
func (f *IBMPIInstanceClient) DeleteNetwork(id string, body *models.PVMInstanceRemoveNetwork) error {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesNetworksDeleteParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(id).
		WithBody(body)
	_, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesNetworksDelete(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to delete the network to the pvminstanceid %s: %w", id, err)
	}
	return nil
}
