package errors

import (
	"encoding/json"
	"reflect"
	//"strconv"

	"github.com/IBM-Cloud/power-go-client/power/models"
)

// start of Placementgroup Messages

const GetPlacementGroupOperationFailed = "failed to perform Get Placement Group Operation for placement group  %s with error %v"
const CreatePlacementGroupOperationFailed = "failed to perform Create Placement Group Operation for powerinstanceid %s with error  %v"
const DeletePlacementGroupOperationFailed = "failed to perform Delete Placement Group Operation for placement group %s with error %v"
const UpdatePlacementGroupOperationFailed = "failed to perform Update Placement Group Operation for powerinstanceid  %s and placement group %s with error %v"
const DeleteMemberPlacementGroupOperationFailed = "failed to perform Delete Placement Group Operation for powerinstanceid  %s and placement group %s with error %v"

// start of Cloud Connection Messages

const GetCloudConnectionOperationFailed = "failed to perform Get Cloud Connections Operation for powerinstanceid %s with error %v"
const CreateCloudConnectionOperationFailed = "failed to perform Create Cloud Connection Operation for powerinstanceid %s with error %v"
const UpdateCloudConnectionOperationFailed = "failed to perform Update Cloud Connection Operation for cloudconnectionid  %s with error %v"
const DeleteCloudConnectionOperationFailed = "failed to perform Delete Cloud Connection Operation for cloudconnectionid  %s with error %v"

// start of System-Pools Messages
const GetSystemPoolsOperationFailed = "failed to perform Get System Pools Operation for powerinstanceid %s with error %v"

// start of Image Messages

const GetImageOperationFailed = "failed to perform Get Image Operation for image  %s with error %v"
const CreateImageOperationFailed = "failed to perform Create Image Operation for powerinstanceid %s with error  %v"

// Start of Network Messages
const GetNetworkOperationFailed = "failed to perform Get Network  Operation for Network id %s with error %v"
const CreateNetworkOperationFailed = "failed to perform Create Network Operation for Network %s with error %v"
const CreateNetworkPortOperationFailed = "failed to perform Create Network Port Operation for Network %s with error %v"
const AttachNetworkPortOperationFailed = "failed to perform Attach Network Port Operation for Port %s to Network %s with error %v"

// start of Volume Messages
const DeleteVolumeOperationFailed = "failed to perform Delete Volume Operation for volume %s with error %v"
const UpdateVolumeOperationFailed = "failed to perform Update Volume Operation for volume %s with error %v"
const GetVolumeOperationFailed = "failed to perform the Get Volume Operation for volume  %s with error %v"
const CreateVolumeOperationFailed = "failed to perform the Create volume Operation for volume  %s with  error %v"
const CreateVolumeV2OperationFailed = "failed to perform the Create volume Operation V2 for volume  %s  with error  %v"
const AttachVolumeOperationFailed = "failed to perform the Attach volume Operation for volume  %s  with error  %v"

// start of Clone Messages
const StartCloneOperationFailed = "failed to start the clone operation for %v"
const PrepareCloneOperationFailed = "failed to prepare the clone operation for %v"
const DeleteCloneOperationFailed = "failed to perform Delete operation %v"
const GetCloneOperationFailed = "failed to get the volumes-clone for the power instanceid  %s with error %v"
const CreateCloneOperationFailed = "failed to perform the create clone operation %v"

// start of Cloud Instance Messages
const GetCloudInstanceOperationFailed = "failed to Get Cloud Instance %s with error %v"
const UpdateCloudInstanceOperationFailed = "failed to update the Cloud instance %s with error %v"
const DeleteCloudInstanceOperationFailed = "failed to delete the Cloud instance %s with error %v"

// start of PI Key Messages
const GetPIKeyOperationFailed = "failed to Get PI Key %s with error %v"
const CreatePIKeyOperationFailed = "failed to Create PI Key %s with error %v"
const DeletePIKeyOperationFailed = "failed to Delete PI Key %s with error %v"

// ErrorTarget ...
type ErrorTarget struct {
	Name string
	Type string
}

// SingleError ...
type SingleError struct {
	Code     string
	Message  string
	MoreInfo string
	Target   ErrorTarget
}

// PowerError ...
type Error struct {
	Payload *models.Error
}

func (e Error) Error() string {
	b, _ := json.Marshal(e.Payload)
	return string(b)
}

// ToError ...
func ToError(err error) error {
	if err == nil {
		return nil
	}

	// check if its ours
	kind := reflect.TypeOf(err).Kind()
	if kind != reflect.Ptr {
		return err
	}

	// next follow pointer
	errstruct := reflect.TypeOf(err).Elem()
	if errstruct.Kind() != reflect.Struct {
		return err
	}

	n := errstruct.NumField()
	found := false
	for i := 0; i < n; i++ {
		if errstruct.Field(i).Name == "Payload" {
			found = true
			break
		}
	}

	if !found {
		return err
	}

	// check if a payload field exists
	payloadValue := reflect.ValueOf(err).Elem().FieldByName("Payload")
	if payloadValue.Interface() == nil {
		return err
	}

	payloadIntf := payloadValue.Elem().Interface()
	payload, parsed := payloadIntf.(models.Error)
	if !parsed {
		return err
	}

	var reterr = Error{
		Payload: &payload,
	}

	return reterr
}
