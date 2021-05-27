package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_snapshots"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPISnapshotClient ...
type IBMPISnapshotClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPISnapshotClient ...
func NewIBMPISnapshotClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPISnapshotClient {
	return &IBMPISnapshotClient{
		sess, powerinstanceid,
	}
}

//Get information about a single snapshot only
func (f *IBMPISnapshotClient) Get(id, powerinstanceid string, timeout time.Duration) (*models.Snapshot, error) {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithSnapshotID(id)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PI Snapshot %s :%s", id, err)
	}
	return resp.Payload, nil
}

// Delete ...
func (f *IBMPISnapshotClient) Delete(id string, powerinstanceid string, timeout time.Duration) error {
	//var cloudinstanceid = f.session.PowerServiceInstance
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(powerinstanceid).WithSnapshotID(id)
	_, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf("Failed to Delete PI Snapshot %s :%s", id, err)
	}
	return nil
}

// Update ...
func (f *IBMPISnapshotClient) Update(id, powerinstanceid string, snapshotdef *models.SnapshotUpdate, timeout time.Duration) (models.Object, error) {

	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsPutParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithSnapshotID(id).WithBody(snapshotdef)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf("Failed to Update PI Snapshot %s :%s", id, err)
	}
	return resp.Payload, nil
}

// GetAll snapshots part of an instance
func (f *IBMPISnapshotClient) GetAll(id, powerinstanceid string, timeout time.Duration) (*models.Snapshots, error) {
	params := p_cloud_snapshots.NewPcloudCloudinstancesSnapshotsGetallParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudSnapshots.PcloudCloudinstancesSnapshotsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Get all  PI Snapshot %s :%s", id, err)
	}
	return resp.Payload, nil

}

// Create or Restore a Snapshot
func (f *IBMPISnapshotClient) Create(pvminstanceid, powerinstanceid, snapshotid, restorefailAction string, timeout time.Duration) (*models.Snapshot, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsRestorePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithCloudInstanceID(pvminstanceid).WithSnapshotID(snapshotid).WithRestoreFailAction(&restorefailAction)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesSnapshotsRestorePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to restore  PI Snapshot %s of the instance %s :%s", snapshotid, pvminstanceid, err)
	}
	return resp.Payload, nil
}
