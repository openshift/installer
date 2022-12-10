package schedule

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i26e1b177111f8a00c755a9493873e39cc8e54d97a7ed846fbab5e2b98e6086aa "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/shifts"
    i303f2b59aaf94036e4c52e1f6eec495c231b2ea4269c6c620f1968cc7f37d29e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/offershiftrequests"
    i38b1d0cd4787f4c709deefdebc9011d65ea498ddd72bbe8c6d6753403d1da3a7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/openshifts"
    i4d899631f98079a2a1e192a51bd5f13a4ad4dcf317b82c1ae84986c572375f9d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timeoffreasons"
    i6977a5d84cfb8f3a7ca758621fa1892ac67095ec26831e278bead93844d09d6e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/schedulinggroups"
    i6a080e6c11a55185c1b16614943b747c47e5b588cd0dda9a243670d6689511f5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timeoffrequests"
    i76a90b8fe290b0eb35da52697b88cc104e65700992ff3f7f1ebf7568fc69de62 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/share"
    i86ba6ba234181696e23fb2cde502d1d1e9e082bb456412adf8467529ce308948 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/swapshiftschangerequests"
    i9ec7ccde7dc719a242dbb335c137486b82a7146d5e1fed1039608fadafbf62f6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/openshiftchangerequests"
    ib294d815ef40c29e26e10f3c1284fa16b88f1c5cd5c8c14103142614dd99c291 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timesoff"
    i0d1a68bf16e77d2487a937577799dd6167ffee937d965666d69f6ad7f3c09a50 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/schedulinggroups/item"
    i1081b52037899720d0759af03675121003d474a31d48434a51a13f00717d1de6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timeoffrequests/item"
    i25d31ae09339a87264a818f06856b020d1f7f458ad5bcd376f77616d64336b31 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/openshiftchangerequests/item"
    i544fd4f38bdf20701505814dc209e583f746b4b10fee5d775fea33949a2c1d8a "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timesoff/item"
    i623731e26f5962ffc40ed47d9738052818c87601af17f1ad3e5cf3bfa9ff0675 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/offershiftrequests/item"
    i75132cc5252e03a3d7b7df052ceea5cf0151431f6a6ebc0935d43698207d46b3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/timeoffreasons/item"
    i78242729c849fb24be50e8793396cac5a437ebdcd77f00949f3f62b2d7bab779 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/openshifts/item"
    if4d4fbf86081722c0a7ebc0e19f822a786b3897e217dc778edba5c1aa53a731c "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/shifts/item"
    if7aca6b7d8b5075b0972208f9d700383af1ec4e05db2a2e3d7672d44075b206f "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule/swapshiftschangerequests/item"
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
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}/schedule{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property schedule for users
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
// CreatePutRequestInformation update the navigation property schedule in users
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
// Delete delete navigation property schedule for users
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
func (m *ScheduleRequestBuilder) OfferShiftRequests()(*i303f2b59aaf94036e4c52e1f6eec495c231b2ea4269c6c620f1968cc7f37d29e.OfferShiftRequestsRequestBuilder) {
    return i303f2b59aaf94036e4c52e1f6eec495c231b2ea4269c6c620f1968cc7f37d29e.NewOfferShiftRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OfferShiftRequestsById provides operations to manage the offerShiftRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OfferShiftRequestsById(id string)(*i623731e26f5962ffc40ed47d9738052818c87601af17f1ad3e5cf3bfa9ff0675.OfferShiftRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["offerShiftRequest%2Did"] = id
    }
    return i623731e26f5962ffc40ed47d9738052818c87601af17f1ad3e5cf3bfa9ff0675.NewOfferShiftRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShiftChangeRequests provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequests()(*i9ec7ccde7dc719a242dbb335c137486b82a7146d5e1fed1039608fadafbf62f6.OpenShiftChangeRequestsRequestBuilder) {
    return i9ec7ccde7dc719a242dbb335c137486b82a7146d5e1fed1039608fadafbf62f6.NewOpenShiftChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftChangeRequestsById provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequestsById(id string)(*i25d31ae09339a87264a818f06856b020d1f7f458ad5bcd376f77616d64336b31.OpenShiftChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShiftChangeRequest%2Did"] = id
    }
    return i25d31ae09339a87264a818f06856b020d1f7f458ad5bcd376f77616d64336b31.NewOpenShiftChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShifts provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShifts()(*i38b1d0cd4787f4c709deefdebc9011d65ea498ddd72bbe8c6d6753403d1da3a7.OpenShiftsRequestBuilder) {
    return i38b1d0cd4787f4c709deefdebc9011d65ea498ddd72bbe8c6d6753403d1da3a7.NewOpenShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftsById provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftsById(id string)(*i78242729c849fb24be50e8793396cac5a437ebdcd77f00949f3f62b2d7bab779.OpenShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShift%2Did"] = id
    }
    return i78242729c849fb24be50e8793396cac5a437ebdcd77f00949f3f62b2d7bab779.NewOpenShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Put update the navigation property schedule in users
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
func (m *ScheduleRequestBuilder) SchedulingGroups()(*i6977a5d84cfb8f3a7ca758621fa1892ac67095ec26831e278bead93844d09d6e.SchedulingGroupsRequestBuilder) {
    return i6977a5d84cfb8f3a7ca758621fa1892ac67095ec26831e278bead93844d09d6e.NewSchedulingGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchedulingGroupsById provides operations to manage the schedulingGroups property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SchedulingGroupsById(id string)(*i0d1a68bf16e77d2487a937577799dd6167ffee937d965666d69f6ad7f3c09a50.SchedulingGroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["schedulingGroup%2Did"] = id
    }
    return i0d1a68bf16e77d2487a937577799dd6167ffee937d965666d69f6ad7f3c09a50.NewSchedulingGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Share provides operations to call the share method.
func (m *ScheduleRequestBuilder) Share()(*i76a90b8fe290b0eb35da52697b88cc104e65700992ff3f7f1ebf7568fc69de62.ShareRequestBuilder) {
    return i76a90b8fe290b0eb35da52697b88cc104e65700992ff3f7f1ebf7568fc69de62.NewShareRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Shifts provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) Shifts()(*i26e1b177111f8a00c755a9493873e39cc8e54d97a7ed846fbab5e2b98e6086aa.ShiftsRequestBuilder) {
    return i26e1b177111f8a00c755a9493873e39cc8e54d97a7ed846fbab5e2b98e6086aa.NewShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShiftsById provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) ShiftsById(id string)(*if4d4fbf86081722c0a7ebc0e19f822a786b3897e217dc778edba5c1aa53a731c.ShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["shift%2Did"] = id
    }
    return if4d4fbf86081722c0a7ebc0e19f822a786b3897e217dc778edba5c1aa53a731c.NewShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SwapShiftsChangeRequests provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequests()(*i86ba6ba234181696e23fb2cde502d1d1e9e082bb456412adf8467529ce308948.SwapShiftsChangeRequestsRequestBuilder) {
    return i86ba6ba234181696e23fb2cde502d1d1e9e082bb456412adf8467529ce308948.NewSwapShiftsChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SwapShiftsChangeRequestsById provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequestsById(id string)(*if7aca6b7d8b5075b0972208f9d700383af1ec4e05db2a2e3d7672d44075b206f.SwapShiftsChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["swapShiftsChangeRequest%2Did"] = id
    }
    return if7aca6b7d8b5075b0972208f9d700383af1ec4e05db2a2e3d7672d44075b206f.NewSwapShiftsChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffReasons provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasons()(*i4d899631f98079a2a1e192a51bd5f13a4ad4dcf317b82c1ae84986c572375f9d.TimeOffReasonsRequestBuilder) {
    return i4d899631f98079a2a1e192a51bd5f13a4ad4dcf317b82c1ae84986c572375f9d.NewTimeOffReasonsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffReasonsById provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasonsById(id string)(*i75132cc5252e03a3d7b7df052ceea5cf0151431f6a6ebc0935d43698207d46b3.TimeOffReasonItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffReason%2Did"] = id
    }
    return i75132cc5252e03a3d7b7df052ceea5cf0151431f6a6ebc0935d43698207d46b3.NewTimeOffReasonItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffRequests provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequests()(*i6a080e6c11a55185c1b16614943b747c47e5b588cd0dda9a243670d6689511f5.TimeOffRequestsRequestBuilder) {
    return i6a080e6c11a55185c1b16614943b747c47e5b588cd0dda9a243670d6689511f5.NewTimeOffRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffRequestsById provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequestsById(id string)(*i1081b52037899720d0759af03675121003d474a31d48434a51a13f00717d1de6.TimeOffRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffRequest%2Did"] = id
    }
    return i1081b52037899720d0759af03675121003d474a31d48434a51a13f00717d1de6.NewTimeOffRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimesOff provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOff()(*ib294d815ef40c29e26e10f3c1284fa16b88f1c5cd5c8c14103142614dd99c291.TimesOffRequestBuilder) {
    return ib294d815ef40c29e26e10f3c1284fa16b88f1c5cd5c8c14103142614dd99c291.NewTimesOffRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimesOffById provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOffById(id string)(*i544fd4f38bdf20701505814dc209e583f746b4b10fee5d775fea33949a2c1d8a.TimeOffItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOff%2Did"] = id
    }
    return i544fd4f38bdf20701505814dc209e583f746b4b10fee5d775fea33949a2c1d8a.NewTimeOffItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
