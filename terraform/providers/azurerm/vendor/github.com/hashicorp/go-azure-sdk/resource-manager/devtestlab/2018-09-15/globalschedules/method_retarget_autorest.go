package globalschedules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetargetOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Retarget ...
func (c GlobalSchedulesClient) Retarget(ctx context.Context, id ScheduleId, input RetargetScheduleProperties) (result RetargetOperationResponse, err error) {
	req, err := c.preparerForRetarget(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "globalschedules.GlobalSchedulesClient", "Retarget", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRetarget(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "globalschedules.GlobalSchedulesClient", "Retarget", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RetargetThenPoll performs Retarget then polls until it's completed
func (c GlobalSchedulesClient) RetargetThenPoll(ctx context.Context, id ScheduleId, input RetargetScheduleProperties) error {
	result, err := c.Retarget(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Retarget: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Retarget: %+v", err)
	}

	return nil
}

// preparerForRetarget prepares the Retarget request.
func (c GlobalSchedulesClient) preparerForRetarget(ctx context.Context, id ScheduleId, input RetargetScheduleProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/retarget", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRetarget sends the Retarget request. The method will close the
// http.Response Body if it receives an error.
func (c GlobalSchedulesClient) senderForRetarget(ctx context.Context, req *http.Request) (future RetargetOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
