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

type GetBackupSecurityPINRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DppBaseResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetBackupSecurityPINRequestsObjectsOperationResponse, error)
}

type GetBackupSecurityPINRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

func (r GetBackupSecurityPINRequestsObjectsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetBackupSecurityPINRequestsObjectsOperationResponse) LoadMore(ctx context.Context) (resp GetBackupSecurityPINRequestsObjectsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetBackupSecurityPINRequestsObjects ...
func (c ResourceGuardsClient) GetBackupSecurityPINRequestsObjects(ctx context.Context, id ResourceGuardId) (resp GetBackupSecurityPINRequestsObjectsOperationResponse, err error) {
	req, err := c.preparerForGetBackupSecurityPINRequestsObjects(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetBackupSecurityPINRequestsObjects(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGetBackupSecurityPINRequestsObjects prepares the GetBackupSecurityPINRequestsObjects request.
func (c ResourceGuardsClient) preparerForGetBackupSecurityPINRequestsObjects(ctx context.Context, id ResourceGuardId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getBackupSecurityPINRequests", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetBackupSecurityPINRequestsObjectsWithNextLink prepares the GetBackupSecurityPINRequestsObjects request with the given nextLink token.
func (c ResourceGuardsClient) preparerForGetBackupSecurityPINRequestsObjectsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetBackupSecurityPINRequestsObjects handles the response to the GetBackupSecurityPINRequestsObjects request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetBackupSecurityPINRequestsObjects(resp *http.Response) (result GetBackupSecurityPINRequestsObjectsOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetBackupSecurityPINRequestsObjectsOperationResponse, err error) {
			req, err := c.preparerForGetBackupSecurityPINRequestsObjectsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetBackupSecurityPINRequestsObjects(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetBackupSecurityPINRequestsObjects", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetBackupSecurityPINRequestsObjectsComplete retrieves all of the results into a single object
func (c ResourceGuardsClient) GetBackupSecurityPINRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetBackupSecurityPINRequestsObjectsCompleteResult, error) {
	return c.GetBackupSecurityPINRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetBackupSecurityPINRequestsObjectsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardsClient) GetBackupSecurityPINRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (resp GetBackupSecurityPINRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	page, err := c.GetBackupSecurityPINRequestsObjects(ctx, id)
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

	out := GetBackupSecurityPINRequestsObjectsCompleteResult{
		Items: items,
	}
	return out, nil
}
