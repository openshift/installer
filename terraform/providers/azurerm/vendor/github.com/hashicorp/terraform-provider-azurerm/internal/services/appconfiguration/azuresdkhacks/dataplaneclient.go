package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
)

type DataPlaneClient struct {
	client *appconfiguration.BaseClient
}

func NewDataPlaneClient(client appconfiguration.BaseClient) DataPlaneClient {
	return DataPlaneClient{
		client: &client,
	}
}

func (c DataPlaneClient) GetKeyValuesComplete(ctx context.Context, key string, label string, after string, acceptDatetime string, selectParameter []string) (result KeyValueListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.GetKeyValues")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = c.GetKeyValues(ctx, key, label, after, acceptDatetime, selectParameter)
	return
}

func (c DataPlaneClient) GetKeyValues(ctx context.Context, key string, label string, after string, acceptDatetime string, selectParameter []string) (result KeyValueListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.GetKeyValues")
		defer func() {
			sc := -1
			if result.kvlr.Response.Response != nil {
				sc = result.kvlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = c.getKeyValuesNextResults
	req, err := c.client.GetKeyValuesPreparer(ctx, key, label, after, acceptDatetime, selectParameter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "GetKeyValues", nil, "Failure preparing request")
		return
	}

	resp, err := c.client.GetKeyValuesSender(req)
	if err != nil {
		result.kvlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "GetKeyValues", resp, "Failure sending request")
		return
	}

	result.kvlr, err = c.GetKeyValuesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "GetKeyValues", resp, "Failure responding to request")
		return
	}
	if result.kvlr.hasNextLink() && result.kvlr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

func (c DataPlaneClient) getKeyValuesNextResults(ctx context.Context, lastResults KeyValueListResult) (result KeyValueListResult, err error) {
	req, err := lastResults.keyValueListResultPreparer(ctx, c.client.Endpoint)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "getKeyValuesNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := c.client.GetKeyValuesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "getKeyValuesNextResults", resp, "Failure sending next results request")
	}
	result, err = c.GetKeyValuesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.BaseClient", "getKeyValuesNextResults", resp, "Failure responding to next results request")
	}
	return
}

func (c DataPlaneClient) GetKeyValuesResponder(resp *http.Response) (result KeyValueListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
