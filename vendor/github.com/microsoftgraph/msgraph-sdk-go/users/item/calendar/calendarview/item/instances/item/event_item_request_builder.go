package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0a1772f845f9085861aa090a493458a98894090017e54a50f4c5b88de6644f90 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/singlevalueextendedproperties"
    i12642507ce20354eb16e1d03db64db137eee76ee94c176fcde74028377188b9f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/tentativelyaccept"
    i267fa55853fc556031a9e8954e6d6ca7189170e57c14c5e9becdde2fc1a25509 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/forward"
    i2db7a98d47f916ba53c004f1824aef73f2e99fd54aa97c680370fe090f7e39ad "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/extensions"
    i32a121e8104eb28a076666d209446bbafc37ad8d98c7a39c5f3bec7e5b981a17 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/decline"
    i4aff33fe86c5b97a31c51b5cc8f9b25789985c6f10f08b5f1b8ef661e90c5631 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/attachments"
    i51e9808b3bbcc42c101ad287b0e4d35e9d3d6e9e98cdb407094e22377d4a5c60 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/dismissreminder"
    i5750fd2e6ce08fd287cd06791c29fcb169a7c93950f74de254caae9571cf1627 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/calendar"
    i611d094d48a3d02bdba0eb293c70d1357016e16620d57617b01667f9128608f9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/multivalueextendedproperties"
    i74a6f8db4dc95067036c531df6f30ff434c1c7836b5225ed2fac1b9c6f36627b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/snoozereminder"
    i772f4b3da338df2e528317bf41b414aae3b767a697ec743d0992e790dffc7192 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/accept"
    ibda755edec3f52d7eb9620d5e348523a4620dd11e6f74a5359583a040965cd29 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/cancel"
    i7e15d65fb8896ad37f5d5bafc592ec1b831d58049ce192cc0743eabb2e44c602 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/multivalueextendedproperties/item"
    i854847a82f4ba11c9ef275cfd1248e266e65cad2a886be229151764e47b5704e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/attachments/item"
    ic03d56835c6ce122a368ce4889338baf6dc5aea7daa3939cc717f452afcdfba8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/extensions/item"
    ica43d1587570c8b62a74eb2fa6ad4efe2a21fa8bb7d34c1336464e09ef54f656 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item/singlevalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i772f4b3da338df2e528317bf41b414aae3b767a697ec743d0992e790dffc7192.AcceptRequestBuilder) {
    return i772f4b3da338df2e528317bf41b414aae3b767a697ec743d0992e790dffc7192.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i4aff33fe86c5b97a31c51b5cc8f9b25789985c6f10f08b5f1b8ef661e90c5631.AttachmentsRequestBuilder) {
    return i4aff33fe86c5b97a31c51b5cc8f9b25789985c6f10f08b5f1b8ef661e90c5631.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i854847a82f4ba11c9ef275cfd1248e266e65cad2a886be229151764e47b5704e.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i854847a82f4ba11c9ef275cfd1248e266e65cad2a886be229151764e47b5704e.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i5750fd2e6ce08fd287cd06791c29fcb169a7c93950f74de254caae9571cf1627.CalendarRequestBuilder) {
    return i5750fd2e6ce08fd287cd06791c29fcb169a7c93950f74de254caae9571cf1627.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ibda755edec3f52d7eb9620d5e348523a4620dd11e6f74a5359583a040965cd29.CancelRequestBuilder) {
    return ibda755edec3f52d7eb9620d5e348523a4620dd11e6f74a5359583a040965cd29.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendar/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i32a121e8104eb28a076666d209446bbafc37ad8d98c7a39c5f3bec7e5b981a17.DeclineRequestBuilder) {
    return i32a121e8104eb28a076666d209446bbafc37ad8d98c7a39c5f3bec7e5b981a17.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i51e9808b3bbcc42c101ad287b0e4d35e9d3d6e9e98cdb407094e22377d4a5c60.DismissReminderRequestBuilder) {
    return i51e9808b3bbcc42c101ad287b0e4d35e9d3d6e9e98cdb407094e22377d4a5c60.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i2db7a98d47f916ba53c004f1824aef73f2e99fd54aa97c680370fe090f7e39ad.ExtensionsRequestBuilder) {
    return i2db7a98d47f916ba53c004f1824aef73f2e99fd54aa97c680370fe090f7e39ad.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ic03d56835c6ce122a368ce4889338baf6dc5aea7daa3939cc717f452afcdfba8.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ic03d56835c6ce122a368ce4889338baf6dc5aea7daa3939cc717f452afcdfba8.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i267fa55853fc556031a9e8954e6d6ca7189170e57c14c5e9becdde2fc1a25509.ForwardRequestBuilder) {
    return i267fa55853fc556031a9e8954e6d6ca7189170e57c14c5e9becdde2fc1a25509.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i611d094d48a3d02bdba0eb293c70d1357016e16620d57617b01667f9128608f9.MultiValueExtendedPropertiesRequestBuilder) {
    return i611d094d48a3d02bdba0eb293c70d1357016e16620d57617b01667f9128608f9.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i7e15d65fb8896ad37f5d5bafc592ec1b831d58049ce192cc0743eabb2e44c602.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i7e15d65fb8896ad37f5d5bafc592ec1b831d58049ce192cc0743eabb2e44c602.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i0a1772f845f9085861aa090a493458a98894090017e54a50f4c5b88de6644f90.SingleValueExtendedPropertiesRequestBuilder) {
    return i0a1772f845f9085861aa090a493458a98894090017e54a50f4c5b88de6644f90.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ica43d1587570c8b62a74eb2fa6ad4efe2a21fa8bb7d34c1336464e09ef54f656.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ica43d1587570c8b62a74eb2fa6ad4efe2a21fa8bb7d34c1336464e09ef54f656.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i74a6f8db4dc95067036c531df6f30ff434c1c7836b5225ed2fac1b9c6f36627b.SnoozeReminderRequestBuilder) {
    return i74a6f8db4dc95067036c531df6f30ff434c1c7836b5225ed2fac1b9c6f36627b.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i12642507ce20354eb16e1d03db64db137eee76ee94c176fcde74028377188b9f.TentativelyAcceptRequestBuilder) {
    return i12642507ce20354eb16e1d03db64db137eee76ee94c176fcde74028377188b9f.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
