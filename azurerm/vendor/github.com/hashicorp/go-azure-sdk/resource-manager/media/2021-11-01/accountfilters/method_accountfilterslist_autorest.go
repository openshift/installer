package accountfilters

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

type AccountFiltersListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AccountFilter

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccountFiltersListOperationResponse, error)
}

type AccountFiltersListCompleteResult struct {
	Items []AccountFilter
}

func (r AccountFiltersListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccountFiltersListOperationResponse) LoadMore(ctx context.Context) (resp AccountFiltersListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AccountFiltersList ...
func (c AccountFiltersClient) AccountFiltersList(ctx context.Context, id MediaServiceId) (resp AccountFiltersListOperationResponse, err error) {
	req, err := c.preparerForAccountFiltersList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccountFiltersList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAccountFiltersList prepares the AccountFiltersList request.
func (c AccountFiltersClient) preparerForAccountFiltersList(ctx context.Context, id MediaServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/accountFilters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAccountFiltersListWithNextLink prepares the AccountFiltersList request with the given nextLink token.
func (c AccountFiltersClient) preparerForAccountFiltersListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccountFiltersList handles the response to the AccountFiltersList request. The method always
// closes the http.Response Body.
func (c AccountFiltersClient) responderForAccountFiltersList(resp *http.Response) (result AccountFiltersListOperationResponse, err error) {
	type page struct {
		Values   []AccountFilter `json:"value"`
		NextLink *string         `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccountFiltersListOperationResponse, err error) {
			req, err := c.preparerForAccountFiltersListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccountFiltersList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AccountFiltersListComplete retrieves all of the results into a single object
func (c AccountFiltersClient) AccountFiltersListComplete(ctx context.Context, id MediaServiceId) (AccountFiltersListCompleteResult, error) {
	return c.AccountFiltersListCompleteMatchingPredicate(ctx, id, AccountFilterOperationPredicate{})
}

// AccountFiltersListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AccountFiltersClient) AccountFiltersListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, predicate AccountFilterOperationPredicate) (resp AccountFiltersListCompleteResult, err error) {
	items := make([]AccountFilter, 0)

	page, err := c.AccountFiltersList(ctx, id)
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

	out := AccountFiltersListCompleteResult{
		Items: items,
	}
	return out, nil
}
