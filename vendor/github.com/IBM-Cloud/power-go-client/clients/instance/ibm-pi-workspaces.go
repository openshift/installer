package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/workspaces"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIWorkspacesClient
type IBMPIWorkspacesClient struct {
	IBMPIClient
}

// NewIBMPIWorkspacesClient
func NewIBMPIWorkspacesClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIWorkspacesClient {
	return &IBMPIWorkspacesClient{
		*NewIBMPIClient(ctx, sess, cloudInstanceID),
	}
}

// Get a workspace
func (f *IBMPIWorkspacesClient) Get(cloudInstanceID string) (*models.Workspace, error) {
	if f.session.IsOnPrem() {
		return nil, fmt.Errorf("operation not supported in satellite location, check documentation")
	}
	params := workspaces.NewV1WorkspacesGetParams().WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).WithWorkspaceID(cloudInstanceID)
	resp, err := f.session.Power.Workspaces.V1WorkspacesGet(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf(errors.GetWorkspaceOperationFailed, f.cloudInstanceID, err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get Workspace %s", f.cloudInstanceID)
	}
	return resp.Payload, nil
}

// Get all workspaces
func (f *IBMPIWorkspacesClient) GetAll() (*models.Workspaces, error) {
	if f.session.IsOnPrem() {
		return nil, fmt.Errorf("operation not supported in satellite location, check documentation")
	}
	params := workspaces.NewV1WorkspacesGetallParams().WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut)
	resp, err := f.session.Power.Workspaces.V1WorkspacesGetall(params, f.session.AuthInfo(f.cloudInstanceID))
	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf("failed to Get all Workspaces: %w", err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all Workspaces")
	}
	return resp.Payload, nil
}
