package nginxdeployment

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

type DeploymentsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeploymentsCreateOrUpdate ...
func (c NginxDeploymentClient) DeploymentsCreateOrUpdate(ctx context.Context, id NginxDeploymentId, input NginxDeployment) (result DeploymentsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForDeploymentsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeploymentsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeploymentsCreateOrUpdateThenPoll performs DeploymentsCreateOrUpdate then polls until it's completed
func (c NginxDeploymentClient) DeploymentsCreateOrUpdateThenPoll(ctx context.Context, id NginxDeploymentId, input NginxDeployment) error {
	result, err := c.DeploymentsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DeploymentsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeploymentsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForDeploymentsCreateOrUpdate prepares the DeploymentsCreateOrUpdate request.
func (c NginxDeploymentClient) preparerForDeploymentsCreateOrUpdate(ctx context.Context, id NginxDeploymentId, input NginxDeployment) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDeploymentsCreateOrUpdate sends the DeploymentsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c NginxDeploymentClient) senderForDeploymentsCreateOrUpdate(ctx context.Context, req *http.Request) (future DeploymentsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
