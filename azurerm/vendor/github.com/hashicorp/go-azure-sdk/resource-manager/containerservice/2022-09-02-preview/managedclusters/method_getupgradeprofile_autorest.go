package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUpgradeProfileOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedClusterUpgradeProfile
}

// GetUpgradeProfile ...
func (c ManagedClustersClient) GetUpgradeProfile(ctx context.Context, id ManagedClusterId) (result GetUpgradeProfileOperationResponse, err error) {
	req, err := c.preparerForGetUpgradeProfile(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetUpgradeProfile", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetUpgradeProfile", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetUpgradeProfile(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetUpgradeProfile", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetUpgradeProfile prepares the GetUpgradeProfile request.
func (c ManagedClustersClient) preparerForGetUpgradeProfile(ctx context.Context, id ManagedClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/upgradeProfiles/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetUpgradeProfile handles the response to the GetUpgradeProfile request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetUpgradeProfile(resp *http.Response) (result GetUpgradeProfileOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
