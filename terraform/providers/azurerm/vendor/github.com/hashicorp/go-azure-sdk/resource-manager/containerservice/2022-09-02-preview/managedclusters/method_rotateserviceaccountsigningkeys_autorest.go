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

type RotateServiceAccountSigningKeysOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RotateServiceAccountSigningKeys ...
func (c ManagedClustersClient) RotateServiceAccountSigningKeys(ctx context.Context, id ManagedClusterId) (result RotateServiceAccountSigningKeysOperationResponse, err error) {
	req, err := c.preparerForRotateServiceAccountSigningKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RotateServiceAccountSigningKeys", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRotateServiceAccountSigningKeys(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RotateServiceAccountSigningKeys", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RotateServiceAccountSigningKeysThenPoll performs RotateServiceAccountSigningKeys then polls until it's completed
func (c ManagedClustersClient) RotateServiceAccountSigningKeysThenPoll(ctx context.Context, id ManagedClusterId) error {
	result, err := c.RotateServiceAccountSigningKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RotateServiceAccountSigningKeys: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RotateServiceAccountSigningKeys: %+v", err)
	}

	return nil
}

// preparerForRotateServiceAccountSigningKeys prepares the RotateServiceAccountSigningKeys request.
func (c ManagedClustersClient) preparerForRotateServiceAccountSigningKeys(ctx context.Context, id ManagedClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rotateServiceAccountSigningKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRotateServiceAccountSigningKeys sends the RotateServiceAccountSigningKeys request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedClustersClient) senderForRotateServiceAccountSigningKeys(ctx context.Context, req *http.Request) (future RotateServiceAccountSigningKeysOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
