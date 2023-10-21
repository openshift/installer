package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUpdateSummaryOperationResponse struct {
	HttpResponse *http.Response
	Model        *UpdateSummary
}

// GetUpdateSummary ...
func (c DevicesClient) GetUpdateSummary(ctx context.Context, id DataBoxEdgeDeviceId) (result GetUpdateSummaryOperationResponse, err error) {
	req, err := c.preparerForGetUpdateSummary(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetUpdateSummary", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetUpdateSummary", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetUpdateSummary(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetUpdateSummary", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetUpdateSummary prepares the GetUpdateSummary request.
func (c DevicesClient) preparerForGetUpdateSummary(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateSummary/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetUpdateSummary handles the response to the GetUpdateSummary request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForGetUpdateSummary(resp *http.Response) (result GetUpdateSummaryOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
