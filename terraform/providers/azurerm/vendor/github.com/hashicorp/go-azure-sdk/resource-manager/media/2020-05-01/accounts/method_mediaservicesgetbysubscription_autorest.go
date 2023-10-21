package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesGetBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaService
}

// MediaservicesGetBySubscription ...
func (c AccountsClient) MediaservicesGetBySubscription(ctx context.Context, id MediaServiceId) (result MediaservicesGetBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForMediaservicesGetBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGetBySubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGetBySubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesGetBySubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGetBySubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesGetBySubscription prepares the MediaservicesGetBySubscription request.
func (c AccountsClient) preparerForMediaservicesGetBySubscription(ctx context.Context, id MediaServiceId) (*http.Request, error) {
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

// responderForMediaservicesGetBySubscription handles the response to the MediaservicesGetBySubscription request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesGetBySubscription(resp *http.Response) (result MediaservicesGetBySubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
