package workspaces

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func StatusDirectoryState(ctx context.Context, conn *workspaces.WorkSpaces, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDirectoryByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.State), nil
	}
}

// nosemgrep:ci.workspaces-in-func-name
func StatusWorkspaceState(ctx context.Context, conn *workspaces.WorkSpaces, workspaceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := conn.DescribeWorkspacesWithContext(ctx, &workspaces.DescribeWorkspacesInput{
			WorkspaceIds: aws.StringSlice([]string{workspaceID}),
		})
		if err != nil {
			return nil, workspaces.WorkspaceStateError, err
		}

		if len(output.Workspaces) == 0 {
			return output, workspaces.WorkspaceStateTerminated, nil
		}

		workspace := output.Workspaces[0]

		// https://docs.aws.amazon.com/workspaces/latest/api/API_TerminateWorkspaces.html
		// State TERMINATED is overridden with TERMINATING to catch up directory metadata clean up.
		if aws.StringValue(workspace.State) == workspaces.WorkspaceStateTerminated {
			return workspace, workspaces.WorkspaceStateTerminating, nil
		}

		return workspace, aws.StringValue(workspace.State), nil
	}
}
