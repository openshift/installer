package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0f91cc3524d407ce127cf52c636ea3d70741ec93c24f733e934c355ef01e9673 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/forward"
    i208443f24ffa0e9d495211e2061b0c16c1cc4e5aebf7a138ffbe44f915a17e0b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/calendar"
    i2122e565a8aa4c89017cd0e7f55859c05eae01b9b38ffdc65ae46e71e9707ee4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/accept"
    i341a1e7e328f323fd2b5b915fd2dbad7995b2bf8f3418f1e5b5e4ac55bb815c4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/decline"
    i3bb2079b15e2279e42f683d8f17113622024239edddadace92a48a4f793d9522 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/tentativelyaccept"
    i85b3ecb1df5dabd5d3ad4c4e6a54e2d1660771b26010ff83b38fa7db3649b996 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/dismissreminder"
    i9407c9bcb1ee71f799ed9d8c791c167b82daf93c1270274e376f805cc9f76d5f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/extensions"
    i994428626bd985a642668dee4c95629ec8caffe540062bae0249e43b54dfc697 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances"
    i9eddc52577a5555a33aafc56b6ec4c751b6c5272a204d1db9d90ab5f60ff902e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/singlevalueextendedproperties"
    ia5d34a53a5d690fc00e53f31e3ea500bbc666ac906da17540530fa20f6d74e70 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/attachments"
    iad6745b576f0cdfa153c080e6ed44bd19ba5fa213eacc4c29dc1e699984f0e3c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/multivalueextendedproperties"
    ibcb96cadb1d3953182dae076a8d0d2f4d74b8522757c415c12ef16502beb9e53 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/cancel"
    ie4ee11a97dc68135822d054e4649c9f9aa8b85b739ea56737a94898a6de75e4f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/snoozereminder"
    i36231df696b83109272fd3bb1ff522ae5b1c7890f3c08c8d848fe767a18cd47a "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/extensions/item"
    i4b0486fcdf6c612baf1d1ebbb06ea4f753f858e537fec0fc098234b5a3bf66f9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/attachments/item"
    i5df7d90865ed74cfa89b9670f8906bc275af7e716c852c939338962b443c4fdb "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item"
    i7df1508e513db666e7370befb3f751526bf7cdf8a085672e191438aec048b505 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/multivalueextendedproperties/item"
    i8bde4d3140c0c2f25ed25759ec43c0e794a2a480b2fc8d25b5f99b0a48d2c100 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/singlevalueextendedproperties/item"
)

// EventItemRequestBuilder provides operations to manage the calendarView property of the microsoft.graph.user entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderGetQueryParameters the calendar view for the calendar. Read-only. Nullable.
type EventItemRequestBuilderGetQueryParameters struct {
    // The end date and time of the time range, represented in ISO 8601 format. For example, 2019-11-08T20:00:00-08:00
    EndDateTime *string
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // The start date and time of the time range, represented in ISO 8601 format. For example, 2019-11-08T19:00:00-08:00
    StartDateTime *string
}
// EventItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EventItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EventItemRequestBuilderGetQueryParameters
}
// Accept provides operations to call the accept method.
func (m *EventItemRequestBuilder) Accept()(*i2122e565a8aa4c89017cd0e7f55859c05eae01b9b38ffdc65ae46e71e9707ee4.AcceptRequestBuilder) {
    return i2122e565a8aa4c89017cd0e7f55859c05eae01b9b38ffdc65ae46e71e9707ee4.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ia5d34a53a5d690fc00e53f31e3ea500bbc666ac906da17540530fa20f6d74e70.AttachmentsRequestBuilder) {
    return ia5d34a53a5d690fc00e53f31e3ea500bbc666ac906da17540530fa20f6d74e70.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i4b0486fcdf6c612baf1d1ebbb06ea4f753f858e537fec0fc098234b5a3bf66f9.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i4b0486fcdf6c612baf1d1ebbb06ea4f753f858e537fec0fc098234b5a3bf66f9.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i208443f24ffa0e9d495211e2061b0c16c1cc4e5aebf7a138ffbe44f915a17e0b.CalendarRequestBuilder) {
    return i208443f24ffa0e9d495211e2061b0c16c1cc4e5aebf7a138ffbe44f915a17e0b.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ibcb96cadb1d3953182dae076a8d0d2f4d74b8522757c415c12ef16502beb9e53.CancelRequestBuilder) {
    return ibcb96cadb1d3953182dae076a8d0d2f4d74b8522757c415c12ef16502beb9e53.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEventItemRequestBuilder instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEventItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation the calendar view for the calendar. Read-only. Nullable.
func (m *EventItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EventItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Decline provides operations to call the decline method.
func (m *EventItemRequestBuilder) Decline()(*i341a1e7e328f323fd2b5b915fd2dbad7995b2bf8f3418f1e5b5e4ac55bb815c4.DeclineRequestBuilder) {
    return i341a1e7e328f323fd2b5b915fd2dbad7995b2bf8f3418f1e5b5e4ac55bb815c4.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i85b3ecb1df5dabd5d3ad4c4e6a54e2d1660771b26010ff83b38fa7db3649b996.DismissReminderRequestBuilder) {
    return i85b3ecb1df5dabd5d3ad4c4e6a54e2d1660771b26010ff83b38fa7db3649b996.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i9407c9bcb1ee71f799ed9d8c791c167b82daf93c1270274e376f805cc9f76d5f.ExtensionsRequestBuilder) {
    return i9407c9bcb1ee71f799ed9d8c791c167b82daf93c1270274e376f805cc9f76d5f.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i36231df696b83109272fd3bb1ff522ae5b1c7890f3c08c8d848fe767a18cd47a.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i36231df696b83109272fd3bb1ff522ae5b1c7890f3c08c8d848fe767a18cd47a.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i0f91cc3524d407ce127cf52c636ea3d70741ec93c24f733e934c355ef01e9673.ForwardRequestBuilder) {
    return i0f91cc3524d407ce127cf52c636ea3d70741ec93c24f733e934c355ef01e9673.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the calendar view for the calendar. Read-only. Nullable.
func (m *EventItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EventItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Eventable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEventFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Eventable), nil
}
// Instances provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Instances()(*i994428626bd985a642668dee4c95629ec8caffe540062bae0249e43b54dfc697.InstancesRequestBuilder) {
    return i994428626bd985a642668dee4c95629ec8caffe540062bae0249e43b54dfc697.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i5df7d90865ed74cfa89b9670f8906bc275af7e716c852c939338962b443c4fdb.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i5df7d90865ed74cfa89b9670f8906bc275af7e716c852c939338962b443c4fdb.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*iad6745b576f0cdfa153c080e6ed44bd19ba5fa213eacc4c29dc1e699984f0e3c.MultiValueExtendedPropertiesRequestBuilder) {
    return iad6745b576f0cdfa153c080e6ed44bd19ba5fa213eacc4c29dc1e699984f0e3c.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i7df1508e513db666e7370befb3f751526bf7cdf8a085672e191438aec048b505.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i7df1508e513db666e7370befb3f751526bf7cdf8a085672e191438aec048b505.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i9eddc52577a5555a33aafc56b6ec4c751b6c5272a204d1db9d90ab5f60ff902e.SingleValueExtendedPropertiesRequestBuilder) {
    return i9eddc52577a5555a33aafc56b6ec4c751b6c5272a204d1db9d90ab5f60ff902e.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i8bde4d3140c0c2f25ed25759ec43c0e794a2a480b2fc8d25b5f99b0a48d2c100.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i8bde4d3140c0c2f25ed25759ec43c0e794a2a480b2fc8d25b5f99b0a48d2c100.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ie4ee11a97dc68135822d054e4649c9f9aa8b85b739ea56737a94898a6de75e4f.SnoozeReminderRequestBuilder) {
    return ie4ee11a97dc68135822d054e4649c9f9aa8b85b739ea56737a94898a6de75e4f.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i3bb2079b15e2279e42f683d8f17113622024239edddadace92a48a4f793d9522.TentativelyAcceptRequestBuilder) {
    return i3bb2079b15e2279e42f683d8f17113622024239edddadace92a48a4f793d9522.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
