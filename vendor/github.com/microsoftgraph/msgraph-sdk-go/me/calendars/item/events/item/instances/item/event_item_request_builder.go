package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i224a43a6f6c22839104e3af8e09c25d6ee5653d0bf039c8cc955d19054af9239 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/calendar"
    i3b73bf1d367d8d58b43286e8815e65b652bdd9632a043961e96cdd4f21a4a456 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/attachments"
    i62e84548f3d8b139eec410acddc1a9ec548d38f53c3adf0accda54573c280237 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/accept"
    i6a76f7e1f0e5de6b644f66577ae554f00c63758d9839dbac747e2a975de2c062 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/extensions"
    i708657df116b55afae5c575793f00bc62083de220954f2b30070d94329ad4c0b "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/cancel"
    i78fa73479b6675db5fd8062a94af6edb28ce40ca11518c2cffa6b9e6892a94df "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/tentativelyaccept"
    ia2d1e74d0c475d24b54809e71da0e1ba8c25369b3fc198918003686d5c13a3b8 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/forward"
    iba5221e61e3b660e493309a9e0dc226fb7a9236ea0382bc762f0b58feb7753c5 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/singlevalueextendedproperties"
    ic1bd788757d48993d62ddaea6350b29916d32a63078388e6faeb72dc45a05d7f "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/multivalueextendedproperties"
    id303083c5a91f2be0053b04d5e3553e1dfe384bb4839cb4b08f8664a109d216f "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/dismissreminder"
    id888541b36ea9db1d12cf8e62e960c429547d5e0db473794ec8fc9bcde6b3301 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/decline"
    idf7d1d2c96845dbdba55aa4449dfe108112b1f0b7ff81ef7da936f0a0f979647 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/snoozereminder"
    i122f3095d5e5a9401b8d3bf34749578d52cc4f098effb0dbeab87f3b7c74cb97 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/singlevalueextendedproperties/item"
    i67c4f1ef75a41a5c6ae05d3b85e11b1c45c47f9e8bd5361a63b25adb073a82d3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/multivalueextendedproperties/item"
    ia34394981640465160cc339a2f60b56da7d55e5f7dd248039a8790e1081043f6 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/attachments/item"
    ic06690588d7937c8b2811db52aac78b63b86909ef55badd8436b8736211d0c1f "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events/item/instances/item/extensions/item"
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
func (m *EventItemRequestBuilder) Accept()(*i62e84548f3d8b139eec410acddc1a9ec548d38f53c3adf0accda54573c280237.AcceptRequestBuilder) {
    return i62e84548f3d8b139eec410acddc1a9ec548d38f53c3adf0accda54573c280237.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i3b73bf1d367d8d58b43286e8815e65b652bdd9632a043961e96cdd4f21a4a456.AttachmentsRequestBuilder) {
    return i3b73bf1d367d8d58b43286e8815e65b652bdd9632a043961e96cdd4f21a4a456.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ia34394981640465160cc339a2f60b56da7d55e5f7dd248039a8790e1081043f6.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ia34394981640465160cc339a2f60b56da7d55e5f7dd248039a8790e1081043f6.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i224a43a6f6c22839104e3af8e09c25d6ee5653d0bf039c8cc955d19054af9239.CalendarRequestBuilder) {
    return i224a43a6f6c22839104e3af8e09c25d6ee5653d0bf039c8cc955d19054af9239.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i708657df116b55afae5c575793f00bc62083de220954f2b30070d94329ad4c0b.CancelRequestBuilder) {
    return i708657df116b55afae5c575793f00bc62083de220954f2b30070d94329ad4c0b.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendars/{calendar%2Did}/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*id888541b36ea9db1d12cf8e62e960c429547d5e0db473794ec8fc9bcde6b3301.DeclineRequestBuilder) {
    return id888541b36ea9db1d12cf8e62e960c429547d5e0db473794ec8fc9bcde6b3301.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*id303083c5a91f2be0053b04d5e3553e1dfe384bb4839cb4b08f8664a109d216f.DismissReminderRequestBuilder) {
    return id303083c5a91f2be0053b04d5e3553e1dfe384bb4839cb4b08f8664a109d216f.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i6a76f7e1f0e5de6b644f66577ae554f00c63758d9839dbac747e2a975de2c062.ExtensionsRequestBuilder) {
    return i6a76f7e1f0e5de6b644f66577ae554f00c63758d9839dbac747e2a975de2c062.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ic06690588d7937c8b2811db52aac78b63b86909ef55badd8436b8736211d0c1f.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ic06690588d7937c8b2811db52aac78b63b86909ef55badd8436b8736211d0c1f.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ia2d1e74d0c475d24b54809e71da0e1ba8c25369b3fc198918003686d5c13a3b8.ForwardRequestBuilder) {
    return ia2d1e74d0c475d24b54809e71da0e1ba8c25369b3fc198918003686d5c13a3b8.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ic1bd788757d48993d62ddaea6350b29916d32a63078388e6faeb72dc45a05d7f.MultiValueExtendedPropertiesRequestBuilder) {
    return ic1bd788757d48993d62ddaea6350b29916d32a63078388e6faeb72dc45a05d7f.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i67c4f1ef75a41a5c6ae05d3b85e11b1c45c47f9e8bd5361a63b25adb073a82d3.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i67c4f1ef75a41a5c6ae05d3b85e11b1c45c47f9e8bd5361a63b25adb073a82d3.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*iba5221e61e3b660e493309a9e0dc226fb7a9236ea0382bc762f0b58feb7753c5.SingleValueExtendedPropertiesRequestBuilder) {
    return iba5221e61e3b660e493309a9e0dc226fb7a9236ea0382bc762f0b58feb7753c5.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i122f3095d5e5a9401b8d3bf34749578d52cc4f098effb0dbeab87f3b7c74cb97.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i122f3095d5e5a9401b8d3bf34749578d52cc4f098effb0dbeab87f3b7c74cb97.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*idf7d1d2c96845dbdba55aa4449dfe108112b1f0b7ff81ef7da936f0a0f979647.SnoozeReminderRequestBuilder) {
    return idf7d1d2c96845dbdba55aa4449dfe108112b1f0b7ff81ef7da936f0a0f979647.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i78fa73479b6675db5fd8062a94af6edb28ce40ca11518c2cffa6b9e6892a94df.TentativelyAcceptRequestBuilder) {
    return i78fa73479b6675db5fd8062a94af6edb28ce40ca11518c2cffa6b9e6892a94df.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
