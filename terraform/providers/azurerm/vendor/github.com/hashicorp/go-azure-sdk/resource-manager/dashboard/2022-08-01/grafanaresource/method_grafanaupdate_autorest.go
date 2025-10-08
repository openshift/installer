package grafanaresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedGrafana
}

// GrafanaUpdate ...
func (c GrafanaResourceClient) GrafanaUpdate(ctx context.Context, id GrafanaId, input ManagedGrafanaUpdateParameters) (result GrafanaUpdateOperationResponse, err error) {
	req, err := c.preparerForGrafanaUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGrafanaUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGrafanaUpdate prepares the GrafanaUpdate request.
func (c GrafanaResourceClient) preparerForGrafanaUpdate(ctx context.Context, id GrafanaId, input ManagedGrafanaUpdateParameters) (*http.Request, error) {
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

// responderForGrafanaUpdate handles the response to the GrafanaUpdate request. The method always
// closes the http.Response Body.
func (c GrafanaResourceClient) responderForGrafanaUpdate(resp *http.Response) (result GrafanaUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
