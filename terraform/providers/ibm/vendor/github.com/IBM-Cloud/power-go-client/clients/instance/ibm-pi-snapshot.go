package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_snapshots"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPISnapshotClient ...
type IBMPISnapshotClient struct {
	IBMPIClient
}

// NewIBMPISnapshotClient ...
func NewIBMPISnapshotClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPISnapshotClient {
	return &IBMPISnapshotClient{
		*NewIBMPIClient(ctx, sess, cloudInstanceID),
	}
}

//Get information about a single snapshot only
func (f *IBMPISnapshotClient) Get(id string) (*models.Snapshot, error) {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsGetParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithSnapshotID(id)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get PI Snapshot %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get PI Snapshot %s", id)
	}
	return resp.Payload, nil
}

// Delete ...
func (f *IBMPISnapshotClient) Delete(id string) error {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsDeleteParams().
		WithContext(f.ctx).WithTimeout(helpers.PIDeleteTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithSnapshotID(id)
	_, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsDelete(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return fmt.Errorf("failed to Delete PI Snapshot %s: %w", id, err)
	}
	return nil
}

// Update ...
func (f *IBMPISnapshotClient) Update(id string, body *models.SnapshotUpdate) (models.Object, error) {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsPutParams().
		WithContext(f.ctx).WithTimeout(helpers.PIUpdateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithSnapshotID(id).
		WithBody(body)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsPut(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Update PI Snapshot %s: %w", id, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Update PI Snapshot %s", id)
	}
	return resp.Payload, nil
}

// GetAll snapshots
func (f *IBMPISnapshotClient) GetAll() (*models.Snapshots, error) {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsGetallParams().
		WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(f.cloudInstanceID)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to Get all PI Snapshots: %w", err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all PI Snapshots")
	}
	return resp.Payload, nil
}

// Create or Restore a Snapshot
func (f *IBMPISnapshotClient) Create(instanceID, snapshotID, restoreFailAction string) (*models.Snapshot, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsRestorePostParams().
		WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut).
		WithCloudInstanceID(f.cloudInstanceID).WithPvmInstanceID(instanceID).
		WithSnapshotID(snapshotID).WithRestoreFailAction(&restoreFailAction)
	resp, err := f.session.Power.PCloudpVMInstances.PcloudPvminstancesSnapshotsRestorePost(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, fmt.Errorf("failed to restore PI Snapshot %s of the instance %s: %w", snapshotID, instanceID, err)
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to restore PI Snapshot %s of the instance %s", snapshotID, instanceID)
	}
	return resp.Payload, nil
}
