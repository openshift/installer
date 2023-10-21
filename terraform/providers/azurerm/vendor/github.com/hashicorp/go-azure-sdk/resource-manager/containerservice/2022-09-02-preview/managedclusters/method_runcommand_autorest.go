package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RunCommand ...
func (c ManagedClustersClient) RunCommand(ctx context.Context, id ManagedClusterId, input RunCommandRequest) (result RunCommandOperationResponse, err error) {
	req, err := c.preparerForRunCommand(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RunCommand", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRunCommand(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RunCommand", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RunCommandThenPoll performs RunCommand then polls until it's completed
func (c ManagedClustersClient) RunCommandThenPoll(ctx context.Context, id ManagedClusterId, input RunCommandRequest) error {
	result, err := c.RunCommand(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RunCommand: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RunCommand: %+v", err)
	}

	return nil
}

// preparerForRunCommand prepares the RunCommand request.
func (c ManagedClustersClient) preparerForRunCommand(ctx context.Context, id ManagedClusterId, input RunCommandRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/runCommand", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRunCommand sends the RunCommand request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedClustersClient) senderForRunCommand(ctx context.Context, req *http.Request) (future RunCommandOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
