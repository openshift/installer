package blobcontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetLegalHoldOperationResponse struct {
	HttpResponse *http.Response
	Model        *LegalHold
}

// SetLegalHold ...
func (c BlobContainersClient) SetLegalHold(ctx context.Context, id ContainerId, input LegalHold) (result SetLegalHoldOperationResponse, err error) {
	req, err := c.preparerForSetLegalHold(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "SetLegalHold", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "SetLegalHold", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSetLegalHold(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "SetLegalHold", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSetLegalHold prepares the SetLegalHold request.
func (c BlobContainersClient) preparerForSetLegalHold(ctx context.Context, id ContainerId, input LegalHold) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/setLegalHold", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSetLegalHold handles the response to the SetLegalHold request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForSetLegalHold(resp *http.Response) (result SetLegalHoldOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
