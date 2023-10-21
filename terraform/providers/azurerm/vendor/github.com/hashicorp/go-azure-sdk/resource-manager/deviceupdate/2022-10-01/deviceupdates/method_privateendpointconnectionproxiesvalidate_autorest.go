package deviceupdates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProxiesValidateOperationResponse struct {
	HttpResponse *http.Response
}

// PrivateEndpointConnectionProxiesValidate ...
func (c DeviceupdatesClient) PrivateEndpointConnectionProxiesValidate(ctx context.Context, id PrivateEndpointConnectionProxyId, input PrivateEndpointConnectionProxy) (result PrivateEndpointConnectionProxiesValidateOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionProxiesValidate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesValidate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesValidate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionProxiesValidate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesValidate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionProxiesValidate prepares the PrivateEndpointConnectionProxiesValidate request.
func (c DeviceupdatesClient) preparerForPrivateEndpointConnectionProxiesValidate(ctx context.Context, id PrivateEndpointConnectionProxyId, input PrivateEndpointConnectionProxy) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateEndpointConnectionProxiesValidate handles the response to the PrivateEndpointConnectionProxiesValidate request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForPrivateEndpointConnectionProxiesValidate(resp *http.Response) (result PrivateEndpointConnectionProxiesValidateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
