package agentpools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AbortLatestOperationOperationResponse struct {
	HttpResponse *http.Response
}

// AbortLatestOperation ...
func (c AgentPoolsClient) AbortLatestOperation(ctx context.Context, id AgentPoolId) (result AbortLatestOperationOperationResponse, err error) {
	req, err := c.preparerForAbortLatestOperation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "AbortLatestOperation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "AbortLatestOperation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAbortLatestOperation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "AbortLatestOperation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAbortLatestOperation prepares the AbortLatestOperation request.
func (c AgentPoolsClient) preparerForAbortLatestOperation(ctx context.Context, id AgentPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/abort", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAbortLatestOperation handles the response to the AbortLatestOperation request. The method always
// closes the http.Response Body.
func (c AgentPoolsClient) responderForAbortLatestOperation(resp *http.Response) (result AbortLatestOperationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
