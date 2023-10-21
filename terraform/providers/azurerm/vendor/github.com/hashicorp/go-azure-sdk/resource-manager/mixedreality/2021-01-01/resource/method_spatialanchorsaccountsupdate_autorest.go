package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpatialAnchorsAccountsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *SpatialAnchorsAccount
}

// SpatialAnchorsAccountsUpdate ...
func (c ResourceClient) SpatialAnchorsAccountsUpdate(ctx context.Context, id SpatialAnchorsAccountId, input SpatialAnchorsAccount) (result SpatialAnchorsAccountsUpdateOperationResponse, err error) {
	req, err := c.preparerForSpatialAnchorsAccountsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSpatialAnchorsAccountsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSpatialAnchorsAccountsUpdate prepares the SpatialAnchorsAccountsUpdate request.
func (c ResourceClient) preparerForSpatialAnchorsAccountsUpdate(ctx context.Context, id SpatialAnchorsAccountId, input SpatialAnchorsAccount) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSpatialAnchorsAccountsUpdate handles the response to the SpatialAnchorsAccountsUpdate request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForSpatialAnchorsAccountsUpdate(resp *http.Response) (result SpatialAnchorsAccountsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
