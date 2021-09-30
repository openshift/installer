package instance

import (
	"fmt"
	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_placement_groups"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

//IBMPIPlacementGroupClient ...
type IBMPIPlacementGroupClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPIImageClient ...
func NewIBMPIPlacementGroupClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIPlacementGroupClient {
	return &IBMPIPlacementGroupClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

// Get PI Placementgroup
func (f *IBMPIPlacementGroupClient) Get(id, powerinstanceid string) (*models.PlacementGroup, error) {

	params := p_cloud_placement_groups.NewPcloudPlacementgroupsGetParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid).WithPlacementGroupID(id)
	resp, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf(errors.GetPlacementGroupOperationFailed, id, err)
	}
	return resp.Payload, nil
}

// GEt All placement groups

func (f *IBMPIPlacementGroupClient) GetAll(powerinstanceid string) (*models.PlacementGroups, error) {
	params := p_cloud_placement_groups.NewPcloudPlacementgroupsGetallParamsWithTimeout(helpers.PIGetTimeOut).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf(errors.GetPlacementGroupOperationFailed, powerinstanceid, err)
	}
	return resp.Payload, nil
}

//Create the placement group
func (f *IBMPIPlacementGroupClient) Create(powerdef *p_cloud_placement_groups.PcloudPlacementgroupsPostParams, powerinstanceid string) (*models.PlacementGroup, error) {

	params := p_cloud_placement_groups.NewPcloudPlacementgroupsPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithBody(powerdef.Body)
	result, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || result == nil || result.Payload == nil {
		return nil, fmt.Errorf(errors.CreatePlacementGroupOperationFailed, powerinstanceid, err)
	}
	return result.Payload, nil

}

// Delete Placement Group
func (f *IBMPIPlacementGroupClient) Delete(id string, powerinstanceid string) error {
	params := p_cloud_placement_groups.NewPcloudPlacementgroupsDeleteParamsWithTimeout(helpers.PIDeleteTimeOut).WithCloudInstanceID(powerinstanceid).WithPlacementGroupID(id)
	_, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf(errors.DeletePlacementGroupOperationFailed, id, err)
	}
	return nil
}

// Adding a member to a  Placement Group
func (f *IBMPIPlacementGroupClient) Update(placementdef *p_cloud_placement_groups.PcloudPlacementgroupsMembersPostParams, placementgroupid, powerinstanceid string) (*models.PlacementGroup, error) {

	params := p_cloud_placement_groups.NewPcloudPlacementgroupsMembersPostParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPlacementGroupID(placementgroupid).WithBody(placementdef.Body)
	resp, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsMembersPost(params, ibmpisession.NewAuth(f.session, f.powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf(errors.UpdatePlacementGroupOperationFailed, powerinstanceid, placementgroupid, err)
	}
	return resp.Payload, nil
}

// Delete Member from Placement Group
func (f *IBMPIPlacementGroupClient) DeleteMember(placementdef *p_cloud_placement_groups.PcloudPlacementgroupsMembersPostParams, placementgroupid, powerinstanceid string) (*models.PlacementGroup, error) {

	params := p_cloud_placement_groups.NewPcloudPlacementgroupsMembersDeleteParamsWithTimeout(helpers.PICreateTimeOut).WithCloudInstanceID(powerinstanceid).WithPlacementGroupID(placementgroupid).WithBody(placementdef.Body)
	resp, err := f.session.Power.PCloudPlacementGroups.PcloudPlacementgroupsMembersDelete(params, ibmpisession.NewAuth(f.session, f.powerinstanceid))

	if err != nil {
		return nil, fmt.Errorf(errors.DeleteMemberPlacementGroupOperationFailed, powerinstanceid, placementgroupid, err)
	}
	return resp.Payload, nil
}
