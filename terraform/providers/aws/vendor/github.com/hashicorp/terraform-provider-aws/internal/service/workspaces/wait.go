package workspaces

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	DirectoryDeregisterInvalidResourceStateTimeout = 2 * time.Minute
	DirectoryRegisterInvalidResourceStateTimeout   = 2 * time.Minute

	// Maximum amount of time to wait for a Directory to return Registered
	DirectoryRegisteredTimeout = 10 * time.Minute

	// Maximum amount of time to wait for a Directory to return Deregistered
	DirectoryDeregisteredTimeout = 10 * time.Minute

	// Maximum amount of time to wait for a WorkSpace to return Available
	WorkspaceAvailableTimeout = 30 * time.Minute

	// Maximum amount of time to wait for a WorkSpace while returning Updating
	WorkspaceUpdatingTimeout = 10 * time.Minute

	// Amount of time to delay before checking WorkSpace when updating
	WorkspaceUpdatingDelay = 1 * time.Minute

	// Maximum amount of time to wait for a WorkSpace to return Terminated
	WorkspaceTerminatedTimeout = 10 * time.Minute
)

func WaitDirectoryRegistered(ctx context.Context, conn *workspaces.WorkSpaces, directoryID string) (*workspaces.WorkspaceDirectory, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{workspaces.WorkspaceDirectoryStateRegistering},
		Target:  []string{workspaces.WorkspaceDirectoryStateRegistered},
		Refresh: StatusDirectoryState(ctx, conn, directoryID),
		Timeout: DirectoryRegisteredTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*workspaces.WorkspaceDirectory); ok {
		return v, err
	}

	return nil, err
}

func WaitDirectoryDeregistered(ctx context.Context, conn *workspaces.WorkSpaces, directoryID string) (*workspaces.WorkspaceDirectory, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceDirectoryStateRegistering,
			workspaces.WorkspaceDirectoryStateRegistered,
			workspaces.WorkspaceDirectoryStateDeregistering,
		},
		Target:  []string{},
		Refresh: StatusDirectoryState(ctx, conn, directoryID),
		Timeout: DirectoryDeregisteredTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*workspaces.WorkspaceDirectory); ok {
		return v, err
	}

	return nil, err
}

func WaitWorkspaceAvailable(ctx context.Context, conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStatePending,
			workspaces.WorkspaceStateStarting,
		},
		Target:  []string{workspaces.WorkspaceStateAvailable},
		Refresh: StatusWorkspaceState(ctx, conn, workspaceID),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}

func WaitWorkspaceTerminated(ctx context.Context, conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStatePending,
			workspaces.WorkspaceStateAvailable,
			workspaces.WorkspaceStateImpaired,
			workspaces.WorkspaceStateUnhealthy,
			workspaces.WorkspaceStateRebooting,
			workspaces.WorkspaceStateStarting,
			workspaces.WorkspaceStateRebuilding,
			workspaces.WorkspaceStateRestoring,
			workspaces.WorkspaceStateMaintenance,
			workspaces.WorkspaceStateAdminMaintenance,
			workspaces.WorkspaceStateSuspended,
			workspaces.WorkspaceStateUpdating,
			workspaces.WorkspaceStateStopping,
			workspaces.WorkspaceStateStopped,
			workspaces.WorkspaceStateTerminating,
			workspaces.WorkspaceStateError,
		},
		Target:  []string{workspaces.WorkspaceStateTerminated},
		Refresh: StatusWorkspaceState(ctx, conn, workspaceID),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}

func WaitWorkspaceUpdated(ctx context.Context, conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	// OperationInProgressException: The properties of this WorkSpace are currently under modification. Please try again in a moment.
	// AWS Workspaces service doesn't change instance status to "Updating" during property modification. Respective AWS Support feature request has been created. Meanwhile, artificial delay is placed here as a workaround.
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStateUpdating,
		},
		Target: []string{
			workspaces.WorkspaceStateAvailable,
			workspaces.WorkspaceStateStopped,
		},
		Refresh: StatusWorkspaceState(ctx, conn, workspaceID),
		Delay:   WorkspaceUpdatingDelay,
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}
