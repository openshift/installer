package singlesignon

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatadogSingleSignOnResource
}

// ConfigurationsGet ...
func (c SingleSignOnClient) ConfigurationsGet(ctx context.Context, id SingleSignOnConfigurationId) (result ConfigurationsGetOperationResponse, err error) {
	req, err := c.preparerForConfigurationsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "singlesignon.SingleSignOnClient", "ConfigurationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "singlesignon.SingleSignOnClient", "ConfigurationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfigurationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "singlesignon.SingleSignOnClient", "ConfigurationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfigurationsGet prepares the ConfigurationsGet request.
func (c SingleSignOnClient) preparerForConfigurationsGet(ctx context.Context, id SingleSignOnConfigurationId) (*http.Request, error) {
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

// responderForConfigurationsGet handles the response to the ConfigurationsGet request. The method always
// closes the http.Response Body.
func (c SingleSignOnClient) responderForConfigurationsGet(resp *http.Response) (result ConfigurationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
