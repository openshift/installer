package queueservice

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

type QueueListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ListQueue

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueueListOperationResponse, error)
}

type QueueListCompleteResult struct {
	Items []ListQueue
}

func (r QueueListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueueListOperationResponse) LoadMore(ctx context.Context) (resp QueueListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type QueueListOperationOptions struct {
	Filter      *string
	Maxpagesize *string
}

func DefaultQueueListOperationOptions() QueueListOperationOptions {
	return QueueListOperationOptions{}
}

func (o QueueListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o QueueListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Maxpagesize != nil {
		out["$maxpagesize"] = *o.Maxpagesize
	}

	return out
}

// QueueList ...
func (c QueueServiceClient) QueueList(ctx context.Context, id StorageAccountId, options QueueListOperationOptions) (resp QueueListOperationResponse, err error) {
	req, err := c.preparerForQueueList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueueList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForQueueList prepares the QueueList request.
func (c QueueServiceClient) preparerForQueueList(ctx context.Context, id StorageAccountId, options QueueListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/queueServices/default/queues", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForQueueListWithNextLink prepares the QueueList request with the given nextLink token.
func (c QueueServiceClient) preparerForQueueListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueueList handles the response to the QueueList request. The method always
// closes the http.Response Body.
func (c QueueServiceClient) responderForQueueList(resp *http.Response) (result QueueListOperationResponse, err error) {
	type page struct {
		Values   []ListQueue `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueueListOperationResponse, err error) {
			req, err := c.preparerForQueueListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueueList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// QueueListComplete retrieves all of the results into a single object
func (c QueueServiceClient) QueueListComplete(ctx context.Context, id StorageAccountId, options QueueListOperationOptions) (QueueListCompleteResult, error) {
	return c.QueueListCompleteMatchingPredicate(ctx, id, options, ListQueueOperationPredicate{})
}

// QueueListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c QueueServiceClient) QueueListCompleteMatchingPredicate(ctx context.Context, id StorageAccountId, options QueueListOperationOptions, predicate ListQueueOperationPredicate) (resp QueueListCompleteResult, err error) {
	items := make([]ListQueue, 0)

	page, err := c.QueueList(ctx, id, options)
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

	out := QueueListCompleteResult{
		Items: items,
	}
	return out, nil
}
