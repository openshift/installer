package alertprocessingrules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRulesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AlertProcessingRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AlertProcessingRulesListBySubscriptionOperationResponse, error)
}

type AlertProcessingRulesListBySubscriptionCompleteResult struct {
	Items []AlertProcessingRule
}

func (r AlertProcessingRulesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AlertProcessingRulesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp AlertProcessingRulesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AlertProcessingRulesListBySubscription ...
func (c AlertProcessingRulesClient) AlertProcessingRulesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp AlertProcessingRulesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAlertProcessingRulesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAlertProcessingRulesListBySubscription prepares the AlertProcessingRulesListBySubscription request.
func (c AlertProcessingRulesClient) preparerForAlertProcessingRulesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.AlertsManagement/actionRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAlertProcessingRulesListBySubscriptionWithNextLink prepares the AlertProcessingRulesListBySubscription request with the given nextLink token.
func (c AlertProcessingRulesClient) preparerForAlertProcessingRulesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAlertProcessingRulesListBySubscription handles the response to the AlertProcessingRulesListBySubscription request. The method always
// closes the http.Response Body.
func (c AlertProcessingRulesClient) responderForAlertProcessingRulesListBySubscription(resp *http.Response) (result AlertProcessingRulesListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []AlertProcessingRule `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AlertProcessingRulesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForAlertProcessingRulesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAlertProcessingRulesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AlertProcessingRulesListBySubscriptionComplete retrieves all of the results into a single object
func (c AlertProcessingRulesClient) AlertProcessingRulesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (AlertProcessingRulesListBySubscriptionCompleteResult, error) {
	return c.AlertProcessingRulesListBySubscriptionCompleteMatchingPredicate(ctx, id, AlertProcessingRuleOperationPredicate{})
}

// AlertProcessingRulesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AlertProcessingRulesClient) AlertProcessingRulesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AlertProcessingRuleOperationPredicate) (resp AlertProcessingRulesListBySubscriptionCompleteResult, err error) {
	items := make([]AlertProcessingRule, 0)

	page, err := c.AlertProcessingRulesListBySubscription(ctx, id)
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

	out := AlertProcessingRulesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
