package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_s_a_p"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

/*  ChangeLog

2020-June-05 : Added the timeout variable to the clients since a lot of the SB / Powervc calls are timing out.

*/

// IBMPIInstanceClient ...
type IBMPIInstanceClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPIInstanceClient ...
func NewIBMPIInstanceClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIInstanceClient {
	return &IBMPIInstanceClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

//Get information about a single pvm only
func (f *IBMPIInstanceClient) Get(id, powerinstanceid string, timeout time.Duration) (*models.PVMInstance, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PVM Instance %s :%s", id, err)
	}
	return resp.Payload, nil
}

// GetAll Information about all the PVM Instances for a Client
func (f *IBMPIInstanceClient) GetAll(powerinstanceid string, timeout time.Duration) (*models.PVMInstances, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesGetallParamsWithTimeout(getTimeOut).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get all PVM Instances of Power Instance %s :%s", powerinstanceid, err)
	}
	return resp.Payload, nil
}

//Create ...
func (f *IBMPIInstanceClient) Create(powerdef *p_cloud_p_vm_instances.PcloudPvminstancesPostParams, powerinstanceid string, timeout time.Duration) (*models.PVMInstanceList, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(powerdef.Body)
	postok, postcreated, postAccepted, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf("Failed to Create PVM Instance :%s", err)
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
	return nil, nil
}

// Delete PVM Instances
func (f *IBMPIInstanceClient) Delete(id, powerinstanceid string, timeout time.Duration) error {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id)
	_, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return fmt.Errorf("Failed to Delete PVM Instance %s :%s", id, err)
	}

	return nil
}

// Update PVM Instances
func (f *IBMPIInstanceClient) Update(id, powerinstanceid string, powerupdateparams *p_cloud_p_vm_instances.PcloudPvminstancesPutParams, timeout time.Duration) (*models.PVMInstanceUpdateResponse, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesPutParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id).WithBody(powerupdateparams.Body)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Update PVM Instance %s :%s", id, err)
	}
	return resp.Payload, nil
}

// Action PVM Instances Operations
func (f *IBMPIInstanceClient) Action(poweractionparams *p_cloud_p_vm_instances.PcloudPvminstancesActionPostParams, id, powerinstanceid string, timeout time.Duration) (models.Object, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesActionPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id).WithBody(poweractionparams.Body)
	postok, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesActionPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Action PVM Instance :%s", err)
	}

	return postok.Payload, nil

}

// PostConsoleURL Generate the Console URL
func (f *IBMPIInstanceClient) PostConsoleURL(id, powerinstanceid string, timeout time.Duration) (models.Object, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesConsolePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id)
	postok, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesConsolePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Generate the Console URL PVM Instance:%s", err)
	}
	return postok.Payload, nil
}

// CaptureInstanceToImageCatalog Captures an instance
func (f *IBMPIInstanceClient) CaptureInstanceToImageCatalog(id, powerinstanceid string, picaptureparams *p_cloud_p_vm_instances.PcloudPvminstancesCapturePostParams, timeout time.Duration) (models.Object, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesCapturePostParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(id).WithBody(picaptureparams.Body)
	postok, _, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesCapturePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Generate the Console URL PVM Instance:%s", err)
	}
	return postok.Payload, nil

}

// CreatePvmSnapShot Create a snapshot of the instance
func (f *IBMPIInstanceClient) CreatePvmSnapShot(snapshotdef *p_cloud_p_vm_instances.PcloudPvminstancesSnapshotsPostParams, pvminstanceid, powerinstanceid string, timeout time.Duration) (*models.SnapshotCreateResponse, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsPostParamsWithTimeout(helpers.PICreateTimeOut).WithPvmInstanceID(pvminstanceid).WithCloudInstanceID(powerinstanceid).WithBody(snapshotdef.Body)
	snapshotpostaccepted, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesSnapshotsPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || snapshotpostaccepted == nil {
		return nil, fmt.Errorf("Failed to Create the snapshot %s for the pvminstance : %s", pvminstanceid, err)
	}
	return snapshotpostaccepted.Payload, nil
}

// CreateClone ...
func (f *IBMPIInstanceClient) CreateClone(clonedef *p_cloud_p_vm_instances.PcloudPvminstancesClonePostParams, pvminstanceid, powerinstanceid string) (*models.PVMInstance, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesClonePostParamsWithTimeout(helpers.PICreateTimeOut).WithPvmInstanceID(pvminstanceid).WithCloudInstanceID(powerinstanceid).WithBody(clonedef.Body)
	clonePost, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesClonePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to create the clone of the pvm instance %s", err)
	}
	return clonePost.Payload, nil
}

// GetSnapShotVM Get information about the snapshots for a vm
func (f *IBMPIInstanceClient) GetSnapShotVM(powerinstanceid, pvminstanceid string, timeout time.Duration) (*models.Snapshots, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsGetallParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(pvminstanceid)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesSnapshotsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get the snapshot for the pvminstance [%s]: %s", pvminstanceid, err)
	}
	return resp.Payload, nil

}

// RestoreSnapShotVM Restore a snapshot
func (f *IBMPIInstanceClient) RestoreSnapShotVM(powerinstanceid, pvminstanceid, snapshotid, restoreAction string, restoreparams *p_cloud_p_vm_instances.PcloudPvminstancesSnapshotsRestorePostParams, timeout time.Duration) (*models.Snapshot, error) {
	params := p_cloud_p_vm_instances.NewPcloudPvminstancesSnapshotsRestorePostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(pvminstanceid).WithSnapshotID(snapshotid).WithRestoreFailAction(&restoreAction).WithBody(restoreparams.Body)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesSnapshotsRestorePost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to restrore the snapshot for the pvminstance [%s]: %s", pvminstanceid, err)
	}
	return resp.Payload, nil
}

// AddNetwork Add a network to the instance
func (f *IBMPIInstanceClient) AddNetwork(powerinstanceid, pvminstanceid string, networkdef *p_cloud_p_vm_instances.PcloudPvminstancesNetworksPostParams, timeout time.Duration) (*models.PVMInstanceNetwork, error) {

	params := p_cloud_p_vm_instances.NewPcloudPvminstancesNetworksPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPvmInstanceID(pvminstanceid).WithBody(networkdef.Body)
	resp, err := f.session.Power.PCloudPVMInstances.PcloudPvminstancesNetworksPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload.NetworkID == "" {
		return nil, fmt.Errorf("Failed to attach the network to the pvminstanceid %s : %s", pvminstanceid, err)
	}
	return resp.Payload, nil
}

// Delete a network from an instance

// CreateSAP Create SAP Systems
func (f *IBMPIInstanceClient) CreateSAP(powerdef *p_cloud_s_a_p.PcloudSapPostParams, powerinstanceid string, timeout time.Duration) (*models.PVMInstanceList, error) {

	params := p_cloud_s_a_p.NewPcloudSapPostParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithBody(powerdef.Body)
	postok, postcreated, postAccepted, err := f.session.Power.PCloudSAP.PcloudSapPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf("Failed to create sap to the power instance %s : %s", powerinstanceid, err)
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

	//return &postok.Payload, nil
	return nil, nil
}

// GetSAPProfiles Get All SAP Profiles
func (f *IBMPIInstanceClient) GetSAPProfiles(powerinstanceid string) (*models.SAPProfiles, error) {

	params := p_cloud_s_a_p.NewPcloudSapGetallParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudSAP.PcloudSapGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to get sap profiles to the power instance %s : %s", powerinstanceid, err)
	}
	return resp.Payload, nil
}

// GetSap Get an SAP profile
func (f *IBMPIInstanceClient) GetSap(powerinstanceid, sapprofileID string) (*models.SAPProfile, error) {
	params := p_cloud_s_a_p.NewPcloudSapGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithSapProfileID(sapprofileID)
	resp, err := f.session.Power.PCloudSAP.PcloudSapGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to get sap profile %s to the power instance %s : %s", sapprofileID, powerinstanceid, err)
	}
	return resp.Payload, nil

}
