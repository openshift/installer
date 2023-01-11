package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2aa695d2031f70f5d124f226c0bb7d90cba2b281788a0e5f062cfd23464f568e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/dismissreminder"
    i32d20a3dbf1acc1aeb627a8d900b6fc7fc488631ad29c0fe0952e66026ab99f2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/forward"
    i33da0e0ab9fc76274bc798c0179ec96b44cd3ca8ae698eb71bb470d7d3a64266 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/singlevalueextendedproperties"
    i4aa0979a3421851170b7f9c219354310292ae88cecddd0c34707076ae4c39872 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/snoozereminder"
    i50529b46cb9815e1ee6989b50a5e7c9d65cdbb639fbeeb17a1745097ca1de34a "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/multivalueextendedproperties"
    i626438065d3e2e4196a8edc03c59cccec617b67714a70388b865384f45bf7d54 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/extensions"
    i76b4b348bf50e98fa2fe202aa8497f6cc42d90db777ae33eef2bad180ba1997e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/decline"
    i77ccc0b59e6aff950dfcfe7893300dcac0af0f4c6a63312c005fcd867b4af5c4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/cancel"
    i8409d063ec88dfe68141ac2e114f2ce751bd923950f2b07ebd6d6538402649d7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/attachments"
    ibbc9746075ab51b54d510dd1431f7223ed53e3c1ead9dbfbf8d4b1c10db05783 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/tentativelyaccept"
    ic2ec03aae4a44329c27dd64dfc4e4a30cf275384f3242568f08e627453d33cf9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/calendar"
    ie4edeb5b8615215dcf77c191289d9633b816d905d2859e386e9bbeaa680a7924 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/accept"
    i14989f66f870841793770a7d9450c5e1f60ce6f731702c9e59616eb27ef280ec "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/extensions/item"
    i44fd1ffc04685be592ee41c9e9af8f5115b183df895b51a20476f788d5cdcff2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/attachments/item"
    i62922412829845e613386fdd7e95855dafb210457259ca4401e86b26f35228bf "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/multivalueextendedproperties/item"
    icbe569d9ab31fcdf58f7ea64bcbe8538e08f365719a1773e15b606d4c87a7a33 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item/singlevalueextendedproperties/item"
)

// EventItemRequestBuilder provides operations to manage the instances property of the microsoft.graph.event entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderGetQueryParameters the occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
type EventItemRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
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
func (m *EventItemRequestBuilder) Accept()(*ie4edeb5b8615215dcf77c191289d9633b816d905d2859e386e9bbeaa680a7924.AcceptRequestBuilder) {
    return ie4edeb5b8615215dcf77c191289d9633b816d905d2859e386e9bbeaa680a7924.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i8409d063ec88dfe68141ac2e114f2ce751bd923950f2b07ebd6d6538402649d7.AttachmentsRequestBuilder) {
    return i8409d063ec88dfe68141ac2e114f2ce751bd923950f2b07ebd6d6538402649d7.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i44fd1ffc04685be592ee41c9e9af8f5115b183df895b51a20476f788d5cdcff2.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i44fd1ffc04685be592ee41c9e9af8f5115b183df895b51a20476f788d5cdcff2.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ic2ec03aae4a44329c27dd64dfc4e4a30cf275384f3242568f08e627453d33cf9.CalendarRequestBuilder) {
    return ic2ec03aae4a44329c27dd64dfc4e4a30cf275384f3242568f08e627453d33cf9.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i77ccc0b59e6aff950dfcfe7893300dcac0af0f4c6a63312c005fcd867b4af5c4.CancelRequestBuilder) {
    return i77ccc0b59e6aff950dfcfe7893300dcac0af0f4c6a63312c005fcd867b4af5c4.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendar/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
// CreateGetRequestInformation the occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
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
func (m *EventItemRequestBuilder) Decline()(*i76b4b348bf50e98fa2fe202aa8497f6cc42d90db777ae33eef2bad180ba1997e.DeclineRequestBuilder) {
    return i76b4b348bf50e98fa2fe202aa8497f6cc42d90db777ae33eef2bad180ba1997e.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i2aa695d2031f70f5d124f226c0bb7d90cba2b281788a0e5f062cfd23464f568e.DismissReminderRequestBuilder) {
    return i2aa695d2031f70f5d124f226c0bb7d90cba2b281788a0e5f062cfd23464f568e.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i626438065d3e2e4196a8edc03c59cccec617b67714a70388b865384f45bf7d54.ExtensionsRequestBuilder) {
    return i626438065d3e2e4196a8edc03c59cccec617b67714a70388b865384f45bf7d54.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i14989f66f870841793770a7d9450c5e1f60ce6f731702c9e59616eb27ef280ec.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i14989f66f870841793770a7d9450c5e1f60ce6f731702c9e59616eb27ef280ec.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i32d20a3dbf1acc1aeb627a8d900b6fc7fc488631ad29c0fe0952e66026ab99f2.ForwardRequestBuilder) {
    return i32d20a3dbf1acc1aeb627a8d900b6fc7fc488631ad29c0fe0952e66026ab99f2.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
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
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i50529b46cb9815e1ee6989b50a5e7c9d65cdbb639fbeeb17a1745097ca1de34a.MultiValueExtendedPropertiesRequestBuilder) {
    return i50529b46cb9815e1ee6989b50a5e7c9d65cdbb639fbeeb17a1745097ca1de34a.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i62922412829845e613386fdd7e95855dafb210457259ca4401e86b26f35228bf.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i62922412829845e613386fdd7e95855dafb210457259ca4401e86b26f35228bf.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i33da0e0ab9fc76274bc798c0179ec96b44cd3ca8ae698eb71bb470d7d3a64266.SingleValueExtendedPropertiesRequestBuilder) {
    return i33da0e0ab9fc76274bc798c0179ec96b44cd3ca8ae698eb71bb470d7d3a64266.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*icbe569d9ab31fcdf58f7ea64bcbe8538e08f365719a1773e15b606d4c87a7a33.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return icbe569d9ab31fcdf58f7ea64bcbe8538e08f365719a1773e15b606d4c87a7a33.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i4aa0979a3421851170b7f9c219354310292ae88cecddd0c34707076ae4c39872.SnoozeReminderRequestBuilder) {
    return i4aa0979a3421851170b7f9c219354310292ae88cecddd0c34707076ae4c39872.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ibbc9746075ab51b54d510dd1431f7223ed53e3c1ead9dbfbf8d4b1c10db05783.TentativelyAcceptRequestBuilder) {
    return ibbc9746075ab51b54d510dd1431f7223ed53e3c1ead9dbfbf8d4b1c10db05783.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
