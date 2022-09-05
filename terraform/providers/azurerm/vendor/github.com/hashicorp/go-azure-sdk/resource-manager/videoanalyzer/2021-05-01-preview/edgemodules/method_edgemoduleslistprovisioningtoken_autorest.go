package edgemodules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesListProvisioningTokenOperationResponse struct {
	HttpResponse *http.Response
	Model        *EdgeModuleProvisioningToken
}

// EdgeModulesListProvisioningToken ...
func (c EdgeModulesClient) EdgeModulesListProvisioningToken(ctx context.Context, id EdgeModuleId, input ListProvisioningTokenInput) (result EdgeModulesListProvisioningTokenOperationResponse, err error) {
	req, err := c.preparerForEdgeModulesListProvisioningToken(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "edgemodules.EdgeModulesClient", "EdgeModulesListProvisioningToken", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "edgemodules.EdgeModulesClient", "EdgeModulesListProvisioningToken", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEdgeModulesListProvisioningToken(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "edgemodules.EdgeModulesClient", "EdgeModulesListProvisioningToken", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEdgeModulesListProvisioningToken prepares the EdgeModulesListProvisioningToken request.
func (c EdgeModulesClient) preparerForEdgeModulesListProvisioningToken(ctx context.Context, id EdgeModuleId, input ListProvisioningTokenInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listProvisioningToken", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEdgeModulesListProvisioningToken handles the response to the EdgeModulesListProvisioningToken request. The method always
// closes the http.Response Body.
func (c EdgeModulesClient) responderForEdgeModulesListProvisioningToken(resp *http.Response) (result EdgeModulesListProvisioningTokenOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
