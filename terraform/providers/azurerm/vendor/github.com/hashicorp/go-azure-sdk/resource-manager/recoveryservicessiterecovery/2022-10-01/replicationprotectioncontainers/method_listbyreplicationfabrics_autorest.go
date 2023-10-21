package replicationprotectioncontainers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByReplicationFabricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ProtectionContainer

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByReplicationFabricsOperationResponse, error)
}

type ListByReplicationFabricsCompleteResult struct {
	Items []ProtectionContainer
}

func (r ListByReplicationFabricsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByReplicationFabricsOperationResponse) LoadMore(ctx context.Context) (resp ListByReplicationFabricsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByReplicationFabrics ...
func (c ReplicationProtectionContainersClient) ListByReplicationFabrics(ctx context.Context, id ReplicationFabricId) (resp ListByReplicationFabricsOperationResponse, err error) {
	req, err := c.preparerForListByReplicationFabrics(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByReplicationFabrics(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByReplicationFabrics prepares the ListByReplicationFabrics request.
func (c ReplicationProtectionContainersClient) preparerForListByReplicationFabrics(ctx context.Context, id ReplicationFabricId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/replicationProtectionContainers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByReplicationFabricsWithNextLink prepares the ListByReplicationFabrics request with the given nextLink token.
func (c ReplicationProtectionContainersClient) preparerForListByReplicationFabricsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByReplicationFabrics handles the response to the ListByReplicationFabrics request. The method always
// closes the http.Response Body.
func (c ReplicationProtectionContainersClient) responderForListByReplicationFabrics(resp *http.Response) (result ListByReplicationFabricsOperationResponse, err error) {
	type page struct {
		Values   []ProtectionContainer `json:"value"`
		NextLink *string               `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByReplicationFabricsOperationResponse, err error) {
			req, err := c.preparerForListByReplicationFabricsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByReplicationFabrics(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "ListByReplicationFabrics", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByReplicationFabricsComplete retrieves all of the results into a single object
func (c ReplicationProtectionContainersClient) ListByReplicationFabricsComplete(ctx context.Context, id ReplicationFabricId) (ListByReplicationFabricsCompleteResult, error) {
	return c.ListByReplicationFabricsCompleteMatchingPredicate(ctx, id, ProtectionContainerOperationPredicate{})
}

// ListByReplicationFabricsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ReplicationProtectionContainersClient) ListByReplicationFabricsCompleteMatchingPredicate(ctx context.Context, id ReplicationFabricId, predicate ProtectionContainerOperationPredicate) (resp ListByReplicationFabricsCompleteResult, err error) {
	items := make([]ProtectionContainer, 0)

	page, err := c.ListByReplicationFabrics(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := ListByReplicationFabricsCompleteResult{
		Items: items,
	}
	return out, nil
}
