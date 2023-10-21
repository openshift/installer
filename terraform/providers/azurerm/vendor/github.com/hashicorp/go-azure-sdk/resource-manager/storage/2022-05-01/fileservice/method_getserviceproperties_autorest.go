package fileservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetServicePropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *FileServiceProperties
}

// GetServiceProperties ...
func (c FileServiceClient) GetServiceProperties(ctx context.Context, id StorageAccountId) (result GetServicePropertiesOperationResponse, err error) {
	req, err := c.preparerForGetServiceProperties(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileservice.FileServiceClient", "GetServiceProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileservice.FileServiceClient", "GetServiceProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetServiceProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileservice.FileServiceClient", "GetServiceProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetServiceProperties prepares the GetServiceProperties request.
func (c FileServiceClient) preparerForGetServiceProperties(ctx context.Context, id StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/fileServices/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetServiceProperties handles the response to the GetServiceProperties request. The method always
// closes the http.Response Body.
func (c FileServiceClient) responderForGetServiceProperties(resp *http.Response) (result GetServicePropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
