package resourceguards

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

type GetDeleteProtectedItemRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DppBaseResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetDeleteProtectedItemRequestsObjectsOperationResponse, error)
}

type GetDeleteProtectedItemRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

func (r GetDeleteProtectedItemRequestsObjectsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetDeleteProtectedItemRequestsObjectsOperationResponse) LoadMore(ctx context.Context) (resp GetDeleteProtectedItemRequestsObjectsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetDeleteProtectedItemRequestsObjects ...
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjects(ctx context.Context, id ResourceGuardId) (resp GetDeleteProtectedItemRequestsObjectsOperationResponse, err error) {
	req, err := c.preparerForGetDeleteProtectedItemRequestsObjects(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetDeleteProtectedItemRequestsObjects(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGetDeleteProtectedItemRequestsObjects prepares the GetDeleteProtectedItemRequestsObjects request.
func (c ResourceGuardsClient) preparerForGetDeleteProtectedItemRequestsObjects(ctx context.Context, id ResourceGuardId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/deleteProtectedItemRequests", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetDeleteProtectedItemRequestsObjectsWithNextLink prepares the GetDeleteProtectedItemRequestsObjects request with the given nextLink token.
func (c ResourceGuardsClient) preparerForGetDeleteProtectedItemRequestsObjectsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetDeleteProtectedItemRequestsObjects handles the response to the GetDeleteProtectedItemRequestsObjects request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDeleteProtectedItemRequestsObjects(resp *http.Response) (result GetDeleteProtectedItemRequestsObjectsOperationResponse, err error) {
	type page struct {
		Values   []DppBaseResource `json:"value"`
		NextLink *string           `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetDeleteProtectedItemRequestsObjectsOperationResponse, err error) {
			req, err := c.preparerForGetDeleteProtectedItemRequestsObjectsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetDeleteProtectedItemRequestsObjects(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDeleteProtectedItemRequestsObjects", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetDeleteProtectedItemRequestsObjectsComplete retrieves all of the results into a single object
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetDeleteProtectedItemRequestsObjectsCompleteResult, error) {
	return c.GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (resp GetDeleteProtectedItemRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	page, err := c.GetDeleteProtectedItemRequestsObjects(ctx, id)
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

	out := GetDeleteProtectedItemRequestsObjectsCompleteResult{
		Items: items,
	}
	return out, nil
}
