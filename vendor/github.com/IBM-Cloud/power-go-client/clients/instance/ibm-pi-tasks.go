package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_tasks"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPITaskClient ...
type IBMPITaskClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPITaskClient ...
func NewIBMPITaskClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPITaskClient {
	return &IBMPITaskClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

// Get ...
func (f *IBMPITaskClient) Get(id, powerinstanceid string) (*models.Task, error) {
	params := p_cloud_tasks.NewPcloudTasksGetParamsWithTimeout(postTimeOut).WithTaskID(id)
	resp, err := f.session.Power.PCloudTasks.PcloudTasksGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to get the task id ... %s", err)
	}
	return resp.Payload, nil
}

// Delete ...
func (f *IBMPITaskClient) Delete(id, powerinstanceid string) (models.Object, error) {

	params := p_cloud_tasks.NewPcloudTasksDeleteParamsWithTimeout(postTimeOut).WithTaskID(id)
	resp, err := f.session.Power.PCloudTasks.PcloudTasksDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to delete the task id ... %s", err)
	}
	return resp.Payload, nil
}
