package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2981f58c7f06c0386f9a4e727317aac9d4c17f231776484a8690114b12a85e5e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/dismissreminder"
    i3317acf916015a6da98992a383137812a5961e5c81e803d22522b2484c48033d "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/forward"
    i38436f0a4e92197268eee818a3e5019cb401034642c8e29cda6af642314b18ce "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/decline"
    i3a260f80b9a8a1f84911d79d4cdcf7bc1d18069773a7f6200a5f6a416868ab04 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/extensions"
    i3e7e30189bb9604f17498df3ab4b1e61dc383ec2de779cbb344374170c352485 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/snoozereminder"
    i5e3379b62cbc68e92960e6e91cd5dc54475bf7d326ab6332654cc7d7fdebca1e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/cancel"
    i73a38c0e9446fa2d7f78ea0d1b6d29befa145927b44439b28f9d2f22a32bcff7 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/attachments"
    i857d339808a1eb72877ecb77512cc129ce7eb1d489c756636e0d8640404d5eac "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/multivalueextendedproperties"
    id2c4b5c449138e176f6eb07afd0ed786b9fec3efba8c9c0dc5c7128b53d7043a "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/accept"
    ie01b50d6f430eb576e8d08605a1c6be014f5239a0ad0a26f95521f85fb8de7e5 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/calendar"
    if3b11f993193c9a0f6e47f8102337e61cb11c136e848eac16c9fa4a1805915dc "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/singlevalueextendedproperties"
    if5e5d38f1d503ab7b83c4d2fdf8c74d829bd0dd8de7f8a8cdd7b612a99a2b9a2 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/tentativelyaccept"
    i07201606642086d089115851b2b58182e0b540badcccd5a11c61214c00edbf81 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/extensions/item"
    i3e804013a81aecff2c19c493c5d680862564cbd38a48d47aeac83d3712b0f930 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/singlevalueextendedproperties/item"
    ie56b9658d8470399c10c76bcf6ba4e2c225096cc3580b7ae93778faa0c4c3153 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/multivalueextendedproperties/item"
    if4a028fefbf5e7f4c6969694da62a99987871155bdaae6ae5f256353567dba98 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*id2c4b5c449138e176f6eb07afd0ed786b9fec3efba8c9c0dc5c7128b53d7043a.AcceptRequestBuilder) {
    return id2c4b5c449138e176f6eb07afd0ed786b9fec3efba8c9c0dc5c7128b53d7043a.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i73a38c0e9446fa2d7f78ea0d1b6d29befa145927b44439b28f9d2f22a32bcff7.AttachmentsRequestBuilder) {
    return i73a38c0e9446fa2d7f78ea0d1b6d29befa145927b44439b28f9d2f22a32bcff7.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*if4a028fefbf5e7f4c6969694da62a99987871155bdaae6ae5f256353567dba98.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return if4a028fefbf5e7f4c6969694da62a99987871155bdaae6ae5f256353567dba98.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ie01b50d6f430eb576e8d08605a1c6be014f5239a0ad0a26f95521f85fb8de7e5.CalendarRequestBuilder) {
    return ie01b50d6f430eb576e8d08605a1c6be014f5239a0ad0a26f95521f85fb8de7e5.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i5e3379b62cbc68e92960e6e91cd5dc54475bf7d326ab6332654cc7d7fdebca1e.CancelRequestBuilder) {
    return i5e3379b62cbc68e92960e6e91cd5dc54475bf7d326ab6332654cc7d7fdebca1e.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendar/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i38436f0a4e92197268eee818a3e5019cb401034642c8e29cda6af642314b18ce.DeclineRequestBuilder) {
    return i38436f0a4e92197268eee818a3e5019cb401034642c8e29cda6af642314b18ce.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i2981f58c7f06c0386f9a4e727317aac9d4c17f231776484a8690114b12a85e5e.DismissReminderRequestBuilder) {
    return i2981f58c7f06c0386f9a4e727317aac9d4c17f231776484a8690114b12a85e5e.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i3a260f80b9a8a1f84911d79d4cdcf7bc1d18069773a7f6200a5f6a416868ab04.ExtensionsRequestBuilder) {
    return i3a260f80b9a8a1f84911d79d4cdcf7bc1d18069773a7f6200a5f6a416868ab04.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i07201606642086d089115851b2b58182e0b540badcccd5a11c61214c00edbf81.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i07201606642086d089115851b2b58182e0b540badcccd5a11c61214c00edbf81.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i3317acf916015a6da98992a383137812a5961e5c81e803d22522b2484c48033d.ForwardRequestBuilder) {
    return i3317acf916015a6da98992a383137812a5961e5c81e803d22522b2484c48033d.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i857d339808a1eb72877ecb77512cc129ce7eb1d489c756636e0d8640404d5eac.MultiValueExtendedPropertiesRequestBuilder) {
    return i857d339808a1eb72877ecb77512cc129ce7eb1d489c756636e0d8640404d5eac.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ie56b9658d8470399c10c76bcf6ba4e2c225096cc3580b7ae93778faa0c4c3153.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ie56b9658d8470399c10c76bcf6ba4e2c225096cc3580b7ae93778faa0c4c3153.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*if3b11f993193c9a0f6e47f8102337e61cb11c136e848eac16c9fa4a1805915dc.SingleValueExtendedPropertiesRequestBuilder) {
    return if3b11f993193c9a0f6e47f8102337e61cb11c136e848eac16c9fa4a1805915dc.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i3e804013a81aecff2c19c493c5d680862564cbd38a48d47aeac83d3712b0f930.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i3e804013a81aecff2c19c493c5d680862564cbd38a48d47aeac83d3712b0f930.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i3e7e30189bb9604f17498df3ab4b1e61dc383ec2de779cbb344374170c352485.SnoozeReminderRequestBuilder) {
    return i3e7e30189bb9604f17498df3ab4b1e61dc383ec2de779cbb344374170c352485.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*if5e5d38f1d503ab7b83c4d2fdf8c74d829bd0dd8de7f8a8cdd7b612a99a2b9a2.TentativelyAcceptRequestBuilder) {
    return if5e5d38f1d503ab7b83c4d2fdf8c74d829bd0dd8de7f8a8cdd7b612a99a2b9a2.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
