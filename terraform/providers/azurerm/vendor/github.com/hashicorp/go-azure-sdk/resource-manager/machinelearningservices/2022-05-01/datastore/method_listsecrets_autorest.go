package datastore

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSecretsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatastoreSecrets
}

// ListSecrets ...
func (c DatastoreClient) ListSecrets(ctx context.Context, id DataStoreId) (result ListSecretsOperationResponse, err error) {
	req, err := c.preparerForListSecrets(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datastore.DatastoreClient", "ListSecrets", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datastore.DatastoreClient", "ListSecrets", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSecrets(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datastore.DatastoreClient", "ListSecrets", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSecrets prepares the ListSecrets request.
func (c DatastoreClient) preparerForListSecrets(ctx context.Context, id DataStoreId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listSecrets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSecrets handles the response to the ListSecrets request. The method always
// closes the http.Response Body.
func (c DatastoreClient) responderForListSecrets(resp *http.Response) (result ListSecretsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("reading response body for DatastoreSecrets: %+v", err)
	}
	model, err := unmarshalDatastoreSecretsImplementation(b)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
