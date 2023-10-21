package configurationassignments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListParentOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListConfigurationAssignmentsResult
}

// ListParent ...
func (c ConfigurationAssignmentsClient) ListParent(ctx context.Context, id ResourceGroupProviderId) (result ListParentOperationResponse, err error) {
	req, err := c.preparerForListParent(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "ListParent", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "ListParent", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListParent(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationassignments.ConfigurationAssignmentsClient", "ListParent", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListParent prepares the ListParent request.
func (c ConfigurationAssignmentsClient) preparerForListParent(ctx context.Context, id ResourceGroupProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Maintenance/configurationAssignments", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListParent handles the response to the ListParent request. The method always
// closes the http.Response Body.
func (c ConfigurationAssignmentsClient) responderForListParent(resp *http.Response) (result ListParentOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
