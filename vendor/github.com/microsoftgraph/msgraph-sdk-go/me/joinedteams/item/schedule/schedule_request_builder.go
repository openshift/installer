package schedule

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i075694b2a8efa878fe00f87889862d739294455f6efff8199d2d989f3ebf4cba "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/schedulinggroups"
    i26a34e0f10f26b3bb561d271a4595092c5f4cd605c28eae42c233eb6f4658521 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/share"
    i2ad47c4606f7a2a188f343b2ac1db602b5e9e1dbfe7bfbea053a304d630990a6 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timeoffrequests"
    i363cad643714e467edd342c9bf7adf7984461e48c9f613f331bdbda9805c7c29 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/openshiftchangerequests"
    i49bb3bef27839d6b6c772c52798a99b9ec445098226c0138d02c56b11d27fc35 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timesoff"
    i56ad5fc9582390c984982fd18a9da92f7672dd0720777db0bb4e02d9a66e3478 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/swapshiftschangerequests"
    i8282dff0ec1ee0d108c1eb2b9bd96541dc02ca1291abc179f3125208e078cd76 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/shifts"
    ic60cd14fb26ecb0b8abe2f8d972cf6c98b37478730d7982231e50648bb1cf343 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/offershiftrequests"
    id52f0a18e9b113d2d9a076f88fc4daac41e43ad2dd5d9b1d9b75035c42631fdb "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/openshifts"
    ie5d5271becf3f8b7fe1400f8f22fed56b444f4b4d65ab9bddf1e300e584c7e39 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timeoffreasons"
    i22d9083c01b60c357dfa3ef52a74e75974fece0d9cc374707c62582e6c670d78 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/offershiftrequests/item"
    i325e1dbc08b3c7de4ec1e0544ea652801f2a7e0e8f122804a01b02f15024054d "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timeoffreasons/item"
    i53445d1cde7ea9f4c78f439e5c3177819401b1e542ff41fcfe9d701d56688231 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timeoffrequests/item"
    i568663956637218756e049423f943bfd975e2fbcc942fda0b43fe2d336b24986 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/swapshiftschangerequests/item"
    i6eb8e2f139966e506839381b3b22a154b0cb6ca4a1e8db1c00104fe5676b5221 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/openshiftchangerequests/item"
    ia0480667d398c5699f3671622951a91abc2c43918021edfd9867308a881e5256 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/shifts/item"
    ib8462250f78a8053c2a681770d8f3127c44cb3b9f10de3c8eb69e7de6e4578f2 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/schedulinggroups/item"
    ibcd0f63ee32639148598b0d6925bc07eed0f7e18d3004b7f1cb46fc2e6e26562 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/openshifts/item"
    id54e49d27fe39eb2b19a73b7f41eb9c65318265f3e9be65ced863c7266f4755b "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule/timesoff/item"
)

// ScheduleRequestBuilder provides operations to manage the schedule property of the microsoft.graph.team entity.
type ScheduleRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ScheduleRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ScheduleRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ScheduleRequestBuilderGetQueryParameters retrieve the properties and relationships of a schedule object. The schedule creation process conforms to the One API guideline for resource based long running operations (RELO).When clients use the PUT method, if the schedule is provisioned, the operation updates the schedule; otherwise, the operation starts the schedule provisioning process in the background. During schedule provisioning, clients can use the GET method to get the schedule and look at the `provisionStatus` property for the current state of the provisioning. If the provisioning failed, clients can get additional information from the `provisionStatusCode` property. Clients can also inspect the configuration of the schedule.
type ScheduleRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ScheduleRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ScheduleRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ScheduleRequestBuilderGetQueryParameters
}
// ScheduleRequestBuilderPutRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ScheduleRequestBuilderPutRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewScheduleRequestBuilderInternal instantiates a new ScheduleRequestBuilder and sets the default values.
func NewScheduleRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ScheduleRequestBuilder) {
    m := &ScheduleRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/joinedTeams/{team%2Did}/schedule{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewScheduleRequestBuilder instantiates a new ScheduleRequestBuilder and sets the default values.
func NewScheduleRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ScheduleRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewScheduleRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property schedule for me
func (m *ScheduleRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ScheduleRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreateGetRequestInformation retrieve the properties and relationships of a schedule object. The schedule creation process conforms to the One API guideline for resource based long running operations (RELO).When clients use the PUT method, if the schedule is provisioned, the operation updates the schedule; otherwise, the operation starts the schedule provisioning process in the background. During schedule provisioning, clients can use the GET method to get the schedule and look at the `provisionStatus` property for the current state of the provisioning. If the provisioning failed, clients can get additional information from the `provisionStatusCode` property. Clients can also inspect the configuration of the schedule.
func (m *ScheduleRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ScheduleRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET
    requestInfo.Headers["Accept"] = "application/json"
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreatePutRequestInformation update the navigation property schedule in me
func (m *ScheduleRequestBuilder) CreatePutRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable, requestConfiguration *ScheduleRequestBuilderPutRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PUT
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Delete delete navigation property schedule for me
func (m *ScheduleRequestBuilder) Delete(ctx context.Context, requestConfiguration *ScheduleRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.CreateDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContentAsync(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get retrieve the properties and relationships of a schedule object. The schedule creation process conforms to the One API guideline for resource based long running operations (RELO).When clients use the PUT method, if the schedule is provisioned, the operation updates the schedule; otherwise, the operation starts the schedule provisioning process in the background. During schedule provisioning, clients can use the GET method to get the schedule and look at the `provisionStatus` property for the current state of the provisioning. If the provisioning failed, clients can get additional information from the `provisionStatusCode` property. Clients can also inspect the configuration of the schedule.
func (m *ScheduleRequestBuilder) Get(ctx context.Context, requestConfiguration *ScheduleRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateScheduleFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable), nil
}
// OfferShiftRequests provides operations to manage the offerShiftRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OfferShiftRequests()(*ic60cd14fb26ecb0b8abe2f8d972cf6c98b37478730d7982231e50648bb1cf343.OfferShiftRequestsRequestBuilder) {
    return ic60cd14fb26ecb0b8abe2f8d972cf6c98b37478730d7982231e50648bb1cf343.NewOfferShiftRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OfferShiftRequestsById provides operations to manage the offerShiftRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OfferShiftRequestsById(id string)(*i22d9083c01b60c357dfa3ef52a74e75974fece0d9cc374707c62582e6c670d78.OfferShiftRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["offerShiftRequest%2Did"] = id
    }
    return i22d9083c01b60c357dfa3ef52a74e75974fece0d9cc374707c62582e6c670d78.NewOfferShiftRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShiftChangeRequests provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequests()(*i363cad643714e467edd342c9bf7adf7984461e48c9f613f331bdbda9805c7c29.OpenShiftChangeRequestsRequestBuilder) {
    return i363cad643714e467edd342c9bf7adf7984461e48c9f613f331bdbda9805c7c29.NewOpenShiftChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftChangeRequestsById provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequestsById(id string)(*i6eb8e2f139966e506839381b3b22a154b0cb6ca4a1e8db1c00104fe5676b5221.OpenShiftChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShiftChangeRequest%2Did"] = id
    }
    return i6eb8e2f139966e506839381b3b22a154b0cb6ca4a1e8db1c00104fe5676b5221.NewOpenShiftChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShifts provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShifts()(*id52f0a18e9b113d2d9a076f88fc4daac41e43ad2dd5d9b1d9b75035c42631fdb.OpenShiftsRequestBuilder) {
    return id52f0a18e9b113d2d9a076f88fc4daac41e43ad2dd5d9b1d9b75035c42631fdb.NewOpenShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftsById provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftsById(id string)(*ibcd0f63ee32639148598b0d6925bc07eed0f7e18d3004b7f1cb46fc2e6e26562.OpenShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShift%2Did"] = id
    }
    return ibcd0f63ee32639148598b0d6925bc07eed0f7e18d3004b7f1cb46fc2e6e26562.NewOpenShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Put update the navigation property schedule in me
func (m *ScheduleRequestBuilder) Put(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable, requestConfiguration *ScheduleRequestBuilderPutRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable, error) {
    requestInfo, err := m.CreatePutRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateScheduleFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Scheduleable), nil
}
// SchedulingGroups provides operations to manage the schedulingGroups property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SchedulingGroups()(*i075694b2a8efa878fe00f87889862d739294455f6efff8199d2d989f3ebf4cba.SchedulingGroupsRequestBuilder) {
    return i075694b2a8efa878fe00f87889862d739294455f6efff8199d2d989f3ebf4cba.NewSchedulingGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchedulingGroupsById provides operations to manage the schedulingGroups property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SchedulingGroupsById(id string)(*ib8462250f78a8053c2a681770d8f3127c44cb3b9f10de3c8eb69e7de6e4578f2.SchedulingGroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["schedulingGroup%2Did"] = id
    }
    return ib8462250f78a8053c2a681770d8f3127c44cb3b9f10de3c8eb69e7de6e4578f2.NewSchedulingGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Share provides operations to call the share method.
func (m *ScheduleRequestBuilder) Share()(*i26a34e0f10f26b3bb561d271a4595092c5f4cd605c28eae42c233eb6f4658521.ShareRequestBuilder) {
    return i26a34e0f10f26b3bb561d271a4595092c5f4cd605c28eae42c233eb6f4658521.NewShareRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Shifts provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) Shifts()(*i8282dff0ec1ee0d108c1eb2b9bd96541dc02ca1291abc179f3125208e078cd76.ShiftsRequestBuilder) {
    return i8282dff0ec1ee0d108c1eb2b9bd96541dc02ca1291abc179f3125208e078cd76.NewShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShiftsById provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) ShiftsById(id string)(*ia0480667d398c5699f3671622951a91abc2c43918021edfd9867308a881e5256.ShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["shift%2Did"] = id
    }
    return ia0480667d398c5699f3671622951a91abc2c43918021edfd9867308a881e5256.NewShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SwapShiftsChangeRequests provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequests()(*i56ad5fc9582390c984982fd18a9da92f7672dd0720777db0bb4e02d9a66e3478.SwapShiftsChangeRequestsRequestBuilder) {
    return i56ad5fc9582390c984982fd18a9da92f7672dd0720777db0bb4e02d9a66e3478.NewSwapShiftsChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SwapShiftsChangeRequestsById provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequestsById(id string)(*i568663956637218756e049423f943bfd975e2fbcc942fda0b43fe2d336b24986.SwapShiftsChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["swapShiftsChangeRequest%2Did"] = id
    }
    return i568663956637218756e049423f943bfd975e2fbcc942fda0b43fe2d336b24986.NewSwapShiftsChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffReasons provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasons()(*ie5d5271becf3f8b7fe1400f8f22fed56b444f4b4d65ab9bddf1e300e584c7e39.TimeOffReasonsRequestBuilder) {
    return ie5d5271becf3f8b7fe1400f8f22fed56b444f4b4d65ab9bddf1e300e584c7e39.NewTimeOffReasonsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffReasonsById provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasonsById(id string)(*i325e1dbc08b3c7de4ec1e0544ea652801f2a7e0e8f122804a01b02f15024054d.TimeOffReasonItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffReason%2Did"] = id
    }
    return i325e1dbc08b3c7de4ec1e0544ea652801f2a7e0e8f122804a01b02f15024054d.NewTimeOffReasonItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffRequests provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequests()(*i2ad47c4606f7a2a188f343b2ac1db602b5e9e1dbfe7bfbea053a304d630990a6.TimeOffRequestsRequestBuilder) {
    return i2ad47c4606f7a2a188f343b2ac1db602b5e9e1dbfe7bfbea053a304d630990a6.NewTimeOffRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffRequestsById provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequestsById(id string)(*i53445d1cde7ea9f4c78f439e5c3177819401b1e542ff41fcfe9d701d56688231.TimeOffRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffRequest%2Did"] = id
    }
    return i53445d1cde7ea9f4c78f439e5c3177819401b1e542ff41fcfe9d701d56688231.NewTimeOffRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimesOff provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOff()(*i49bb3bef27839d6b6c772c52798a99b9ec445098226c0138d02c56b11d27fc35.TimesOffRequestBuilder) {
    return i49bb3bef27839d6b6c772c52798a99b9ec445098226c0138d02c56b11d27fc35.NewTimesOffRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimesOffById provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOffById(id string)(*id54e49d27fe39eb2b19a73b7f41eb9c65318265f3e9be65ced863c7266f4755b.TimeOffItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOff%2Did"] = id
    }
    return id54e49d27fe39eb2b19a73b7f41eb9c65318265f3e9be65ced863c7266f4755b.NewTimeOffItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
