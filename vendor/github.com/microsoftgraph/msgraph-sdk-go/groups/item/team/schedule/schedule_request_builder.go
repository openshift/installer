package schedule

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i072c72d2684659fe7cbdb65aef47d6c1918f631d9ab0628e1f00b14b8b8c2188 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/shifts"
    i252a222eb77235b52b18b0c9f8303553c657839a4db9ead11772772fd8630fe4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/openshiftchangerequests"
    i2547a7f1d67b3d63ac7cd6ad852788ce52638547a6ec23de203e52fa44a92944 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/share"
    i3b0fd31e4b254703a84b00af06b87d47018b7f09e267fbbdec09a339330613fc "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timesoff"
    i4ee57924c99393c68749c89c3f52497a0a762f96e8c85319a78b103466ae6d15 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/swapshiftschangerequests"
    i6d9dcfcd3f7fb04246276019aa4e2779a6d73d2fdefb23aeb6b7727980ae72a8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/openshifts"
    i98ee724eb951b642ca62569abf1d3cda0e9ca342b1c35e6bb93d1ab0b5946191 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/offershiftrequests"
    ib957f55cc6bceca14111f9d522afca85a7a87ee4297cb1a52c18263455a1be9f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timeoffreasons"
    ibe75ffd15e1edfe2fb53fc58a7e53cc582e467c94a026b855b863e17e863e4e2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/schedulinggroups"
    ie9c0fe875266534e995b493c7e216891e7e63ab8cb66e0a1bb94a0ed4bd476f0 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timeoffrequests"
    i25835bf2ea4b2783655a613062849327eae2c462a3fb0bb246ab4e81c0aeefcd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/schedulinggroups/item"
    i4371af37049a130938818ac064e9a7234fa37759b90ddd60c724a9f2e40312aa "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/openshifts/item"
    i49f769d63eb0f564f5e3fd6dcde8d8b66f068577823c2849fb223a24af8c4b6d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/shifts/item"
    i4d1975d833f0871642b081140a8181ebd5b9b14b52c8daf8c17e9fc462496d93 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timeoffrequests/item"
    i7c1a8cd3fa594c1fc7433210be87d46a669f79ffa033b2ed121559b1a7836b36 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/openshiftchangerequests/item"
    i8d73b6115390f96f3f9effc6c115d9a8c01e4bccf3a6ffb524fe8f9d6d189521 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/swapshiftschangerequests/item"
    i91916952245ecadf3a132c25d6712618462caf09e7999a937a6bdd80260050a1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timesoff/item"
    iad1beb593b32b2b10705cc5c3c2e390baf898c1d2ec6abdecf140dd092a393c5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/timeoffreasons/item"
    iea0766a6107f3958d23f6d55165a212b7c18810c581726007daff6658ae36796 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule/offershiftrequests/item"
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
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/team/schedule{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property schedule for groups
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
// CreatePutRequestInformation update the navigation property schedule in groups
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
// Delete delete navigation property schedule for groups
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
func (m *ScheduleRequestBuilder) OfferShiftRequests()(*i98ee724eb951b642ca62569abf1d3cda0e9ca342b1c35e6bb93d1ab0b5946191.OfferShiftRequestsRequestBuilder) {
    return i98ee724eb951b642ca62569abf1d3cda0e9ca342b1c35e6bb93d1ab0b5946191.NewOfferShiftRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OfferShiftRequestsById provides operations to manage the offerShiftRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OfferShiftRequestsById(id string)(*iea0766a6107f3958d23f6d55165a212b7c18810c581726007daff6658ae36796.OfferShiftRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["offerShiftRequest%2Did"] = id
    }
    return iea0766a6107f3958d23f6d55165a212b7c18810c581726007daff6658ae36796.NewOfferShiftRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShiftChangeRequests provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequests()(*i252a222eb77235b52b18b0c9f8303553c657839a4db9ead11772772fd8630fe4.OpenShiftChangeRequestsRequestBuilder) {
    return i252a222eb77235b52b18b0c9f8303553c657839a4db9ead11772772fd8630fe4.NewOpenShiftChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftChangeRequestsById provides operations to manage the openShiftChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftChangeRequestsById(id string)(*i7c1a8cd3fa594c1fc7433210be87d46a669f79ffa033b2ed121559b1a7836b36.OpenShiftChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShiftChangeRequest%2Did"] = id
    }
    return i7c1a8cd3fa594c1fc7433210be87d46a669f79ffa033b2ed121559b1a7836b36.NewOpenShiftChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OpenShifts provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShifts()(*i6d9dcfcd3f7fb04246276019aa4e2779a6d73d2fdefb23aeb6b7727980ae72a8.OpenShiftsRequestBuilder) {
    return i6d9dcfcd3f7fb04246276019aa4e2779a6d73d2fdefb23aeb6b7727980ae72a8.NewOpenShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OpenShiftsById provides operations to manage the openShifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) OpenShiftsById(id string)(*i4371af37049a130938818ac064e9a7234fa37759b90ddd60c724a9f2e40312aa.OpenShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["openShift%2Did"] = id
    }
    return i4371af37049a130938818ac064e9a7234fa37759b90ddd60c724a9f2e40312aa.NewOpenShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Put update the navigation property schedule in groups
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
func (m *ScheduleRequestBuilder) SchedulingGroups()(*ibe75ffd15e1edfe2fb53fc58a7e53cc582e467c94a026b855b863e17e863e4e2.SchedulingGroupsRequestBuilder) {
    return ibe75ffd15e1edfe2fb53fc58a7e53cc582e467c94a026b855b863e17e863e4e2.NewSchedulingGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchedulingGroupsById provides operations to manage the schedulingGroups property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SchedulingGroupsById(id string)(*i25835bf2ea4b2783655a613062849327eae2c462a3fb0bb246ab4e81c0aeefcd.SchedulingGroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["schedulingGroup%2Did"] = id
    }
    return i25835bf2ea4b2783655a613062849327eae2c462a3fb0bb246ab4e81c0aeefcd.NewSchedulingGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Share provides operations to call the share method.
func (m *ScheduleRequestBuilder) Share()(*i2547a7f1d67b3d63ac7cd6ad852788ce52638547a6ec23de203e52fa44a92944.ShareRequestBuilder) {
    return i2547a7f1d67b3d63ac7cd6ad852788ce52638547a6ec23de203e52fa44a92944.NewShareRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Shifts provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) Shifts()(*i072c72d2684659fe7cbdb65aef47d6c1918f631d9ab0628e1f00b14b8b8c2188.ShiftsRequestBuilder) {
    return i072c72d2684659fe7cbdb65aef47d6c1918f631d9ab0628e1f00b14b8b8c2188.NewShiftsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShiftsById provides operations to manage the shifts property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) ShiftsById(id string)(*i49f769d63eb0f564f5e3fd6dcde8d8b66f068577823c2849fb223a24af8c4b6d.ShiftItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["shift%2Did"] = id
    }
    return i49f769d63eb0f564f5e3fd6dcde8d8b66f068577823c2849fb223a24af8c4b6d.NewShiftItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SwapShiftsChangeRequests provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequests()(*i4ee57924c99393c68749c89c3f52497a0a762f96e8c85319a78b103466ae6d15.SwapShiftsChangeRequestsRequestBuilder) {
    return i4ee57924c99393c68749c89c3f52497a0a762f96e8c85319a78b103466ae6d15.NewSwapShiftsChangeRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SwapShiftsChangeRequestsById provides operations to manage the swapShiftsChangeRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) SwapShiftsChangeRequestsById(id string)(*i8d73b6115390f96f3f9effc6c115d9a8c01e4bccf3a6ffb524fe8f9d6d189521.SwapShiftsChangeRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["swapShiftsChangeRequest%2Did"] = id
    }
    return i8d73b6115390f96f3f9effc6c115d9a8c01e4bccf3a6ffb524fe8f9d6d189521.NewSwapShiftsChangeRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffReasons provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasons()(*ib957f55cc6bceca14111f9d522afca85a7a87ee4297cb1a52c18263455a1be9f.TimeOffReasonsRequestBuilder) {
    return ib957f55cc6bceca14111f9d522afca85a7a87ee4297cb1a52c18263455a1be9f.NewTimeOffReasonsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffReasonsById provides operations to manage the timeOffReasons property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffReasonsById(id string)(*iad1beb593b32b2b10705cc5c3c2e390baf898c1d2ec6abdecf140dd092a393c5.TimeOffReasonItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffReason%2Did"] = id
    }
    return iad1beb593b32b2b10705cc5c3c2e390baf898c1d2ec6abdecf140dd092a393c5.NewTimeOffReasonItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimeOffRequests provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequests()(*ie9c0fe875266534e995b493c7e216891e7e63ab8cb66e0a1bb94a0ed4bd476f0.TimeOffRequestsRequestBuilder) {
    return ie9c0fe875266534e995b493c7e216891e7e63ab8cb66e0a1bb94a0ed4bd476f0.NewTimeOffRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimeOffRequestsById provides operations to manage the timeOffRequests property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimeOffRequestsById(id string)(*i4d1975d833f0871642b081140a8181ebd5b9b14b52c8daf8c17e9fc462496d93.TimeOffRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOffRequest%2Did"] = id
    }
    return i4d1975d833f0871642b081140a8181ebd5b9b14b52c8daf8c17e9fc462496d93.NewTimeOffRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TimesOff provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOff()(*i3b0fd31e4b254703a84b00af06b87d47018b7f09e267fbbdec09a339330613fc.TimesOffRequestBuilder) {
    return i3b0fd31e4b254703a84b00af06b87d47018b7f09e267fbbdec09a339330613fc.NewTimesOffRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TimesOffById provides operations to manage the timesOff property of the microsoft.graph.schedule entity.
func (m *ScheduleRequestBuilder) TimesOffById(id string)(*i91916952245ecadf3a132c25d6712618462caf09e7999a937a6bdd80260050a1.TimeOffItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["timeOff%2Did"] = id
    }
    return i91916952245ecadf3a132c25d6712618462caf09e7999a937a6bdd80260050a1.NewTimeOffItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
