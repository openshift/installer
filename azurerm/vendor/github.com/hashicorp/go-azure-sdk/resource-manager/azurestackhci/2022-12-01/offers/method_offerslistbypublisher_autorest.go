package offers

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

type OffersListByPublisherOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Offer

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (OffersListByPublisherOperationResponse, error)
}

type OffersListByPublisherCompleteResult struct {
	Items []Offer
}

func (r OffersListByPublisherOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r OffersListByPublisherOperationResponse) LoadMore(ctx context.Context) (resp OffersListByPublisherOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type OffersListByPublisherOperationOptions struct {
	Expand *string
}

func DefaultOffersListByPublisherOperationOptions() OffersListByPublisherOperationOptions {
	return OffersListByPublisherOperationOptions{}
}

func (o OffersListByPublisherOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o OffersListByPublisherOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// OffersListByPublisher ...
func (c OffersClient) OffersListByPublisher(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions) (resp OffersListByPublisherOperationResponse, err error) {
	req, err := c.preparerForOffersListByPublisher(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForOffersListByPublisher(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForOffersListByPublisher prepares the OffersListByPublisher request.
func (c OffersClient) preparerForOffersListByPublisher(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/offers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForOffersListByPublisherWithNextLink prepares the OffersListByPublisher request with the given nextLink token.
func (c OffersClient) preparerForOffersListByPublisherWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForOffersListByPublisher handles the response to the OffersListByPublisher request. The method always
// closes the http.Response Body.
func (c OffersClient) responderForOffersListByPublisher(resp *http.Response) (result OffersListByPublisherOperationResponse, err error) {
	type page struct {
		Values   []Offer `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result OffersListByPublisherOperationResponse, err error) {
			req, err := c.preparerForOffersListByPublisherWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForOffersListByPublisher(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByPublisher", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// OffersListByPublisherComplete retrieves all of the results into a single object
func (c OffersClient) OffersListByPublisherComplete(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions) (OffersListByPublisherCompleteResult, error) {
	return c.OffersListByPublisherCompleteMatchingPredicate(ctx, id, options, OfferOperationPredicate{})
}

// OffersListByPublisherCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OffersClient) OffersListByPublisherCompleteMatchingPredicate(ctx context.Context, id PublisherId, options OffersListByPublisherOperationOptions, predicate OfferOperationPredicate) (resp OffersListByPublisherCompleteResult, err error) {
	items := make([]Offer, 0)

	page, err := c.OffersListByPublisher(ctx, id, options)
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

	out := OffersListByPublisherCompleteResult{
		Items: items,
	}
	return out, nil
}
