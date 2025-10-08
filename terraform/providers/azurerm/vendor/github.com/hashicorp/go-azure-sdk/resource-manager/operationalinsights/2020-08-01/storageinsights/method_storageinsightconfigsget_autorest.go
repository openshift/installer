package storageinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightConfigsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageInsight
}

// StorageInsightConfigsGet ...
func (c StorageInsightsClient) StorageInsightConfigsGet(ctx context.Context, id StorageInsightConfigId) (result StorageInsightConfigsGetOperationResponse, err error) {
	req, err := c.preparerForStorageInsightConfigsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStorageInsightConfigsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStorageInsightConfigsGet prepares the StorageInsightConfigsGet request.
func (c StorageInsightsClient) preparerForStorageInsightConfigsGet(ctx context.Context, id StorageInsightConfigId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStorageInsightConfigsGet handles the response to the StorageInsightConfigsGet request. The method always
// closes the http.Response Body.
func (c StorageInsightsClient) responderForStorageInsightConfigsGet(resp *http.Response) (result StorageInsightConfigsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
