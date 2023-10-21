package synthetics

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	canaryCreatedTimeout = 5 * time.Minute
	canaryRunningTimeout = 5 * time.Minute
	canaryStoppedTimeout = 5 * time.Minute
	canaryDeletedTimeout = 5 * time.Minute
)

func waitCanaryReady(ctx context.Context, conn *synthetics.Synthetics, name string) (*synthetics.Canary, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{synthetics.CanaryStateCreating, synthetics.CanaryStateUpdating},
		Target:  []string{synthetics.CanaryStateReady},
		Refresh: statusCanaryState(ctx, conn, name),
		Timeout: canaryCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*synthetics.Canary); ok {
		if status := output.Status; aws.StringValue(status.State) == synthetics.CanaryStateError {
			tfresource.SetLastError(err, fmt.Errorf("%s: %s", aws.StringValue(status.StateReasonCode), aws.StringValue(status.StateReason)))
		}

		return output, err
	}

	return nil, err
}

func waitCanaryStopped(ctx context.Context, conn *synthetics.Synthetics, name string) (*synthetics.Canary, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			synthetics.CanaryStateStopping,
			synthetics.CanaryStateUpdating,
			synthetics.CanaryStateRunning,
			synthetics.CanaryStateReady,
			synthetics.CanaryStateStarting,
		},
		Target:  []string{synthetics.CanaryStateStopped},
		Refresh: statusCanaryState(ctx, conn, name),
		Timeout: canaryStoppedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*synthetics.Canary); ok {
		if status := output.Status; aws.StringValue(status.State) == synthetics.CanaryStateError {
			tfresource.SetLastError(err, fmt.Errorf("%s: %s", aws.StringValue(status.StateReasonCode), aws.StringValue(status.StateReason)))
		}

		return output, err
	}

	return nil, err
}

func waitCanaryRunning(ctx context.Context, conn *synthetics.Synthetics, name string) (*synthetics.Canary, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			synthetics.CanaryStateStarting,
			synthetics.CanaryStateUpdating,
			synthetics.CanaryStateStarting,
			synthetics.CanaryStateReady,
		},
		Target:  []string{synthetics.CanaryStateRunning},
		Refresh: statusCanaryState(ctx, conn, name),
		Timeout: canaryRunningTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*synthetics.Canary); ok {
		if status := output.Status; aws.StringValue(status.State) == synthetics.CanaryStateError {
			tfresource.SetLastError(err, fmt.Errorf("%s: %s", aws.StringValue(status.StateReasonCode), aws.StringValue(status.StateReason)))
		}

		return output, err
	}

	return nil, err
}

func waitCanaryDeleted(ctx context.Context, conn *synthetics.Synthetics, name string) (*synthetics.Canary, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{synthetics.CanaryStateDeleting, synthetics.CanaryStateStopped},
		Target:  []string{},
		Refresh: statusCanaryState(ctx, conn, name),
		Timeout: canaryDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*synthetics.Canary); ok {
		if status := output.Status; aws.StringValue(status.State) == synthetics.CanaryStateError {
			tfresource.SetLastError(err, fmt.Errorf("%s: %s", aws.StringValue(status.StateReasonCode), aws.StringValue(status.StateReason)))
		}

		return output, err
	}

	return nil, err
}
