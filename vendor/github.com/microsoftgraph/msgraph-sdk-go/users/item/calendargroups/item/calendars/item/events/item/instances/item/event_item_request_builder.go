package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0aefe50f098a449e69bee7ad796d3adaf5ca0225915758de5533b5af1cd5d3d4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/calendar"
    i275121e84c16a3befa700e0246c4d148291c4f5467bc1d593fe2186acd512bd7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/decline"
    i310314089d3422d43c9590127e71f71ad172e68540b68a482ebe308d1133dd25 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/singlevalueextendedproperties"
    i3ba75c90f6a2129e0dededbd49a857df83e8e5652e53088efa2e440031e55677 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/extensions"
    i6c5f60aa4919c599b5ca2ed4e78e6ceca68b0624738d1eade1a98c2d49222b3b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/forward"
    i7a2688b4dd22a5885f9cfdda0103f66e883120819be45676562c7d045c603b17 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/multivalueextendedproperties"
    i83a689f67f17f9e9c49b5376e4fb065afee84c3f7057a1ff28da735b7b9bc9e5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/dismissreminder"
    ia67d11714c9a47717b4e54a81cdbd069dd56de4ff701ebacb389de2d00bb8d9e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/snoozereminder"
    iadaf7a16adbbf6d8698ba597e063fa400def8ed175febd88751a03daf684121c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/attachments"
    iae6f0e1bcfce2bf343ad891175d3e665aa98c92b3b7d7dc8dd793365bd904343 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/tentativelyaccept"
    iae91fe3bafac33b450c88daa3d5f916463f4b49c314ecea0f5a6774536878bef "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/accept"
    idd3dcbcdf6aec68d0741cf89364a43bc5705b17b2fd86b48348315716154ff77 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/cancel"
    i3db4ccb0e3925e78597745943f5e74c735c59c99bc41b1043b1bd145762cd60d "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/attachments/item"
    i5826c74fe1163782016deeb01ad6e7c8edd6129aed51e020012eb5a0a2eceea2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/extensions/item"
    ib0dbf672e985e31b3cb13208e4540917f62ae0ea12baeb6cf13df0b957c1d7cb "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/multivalueextendedproperties/item"
    id611332bbbcace00f410ded8271d027e69d7b2f7f0d882ac5d3ade2aef3c63f6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/events/item/instances/item/singlevalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*iae91fe3bafac33b450c88daa3d5f916463f4b49c314ecea0f5a6774536878bef.AcceptRequestBuilder) {
    return iae91fe3bafac33b450c88daa3d5f916463f4b49c314ecea0f5a6774536878bef.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*iadaf7a16adbbf6d8698ba597e063fa400def8ed175febd88751a03daf684121c.AttachmentsRequestBuilder) {
    return iadaf7a16adbbf6d8698ba597e063fa400def8ed175febd88751a03daf684121c.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i3db4ccb0e3925e78597745943f5e74c735c59c99bc41b1043b1bd145762cd60d.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i3db4ccb0e3925e78597745943f5e74c735c59c99bc41b1043b1bd145762cd60d.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i0aefe50f098a449e69bee7ad796d3adaf5ca0225915758de5533b5af1cd5d3d4.CalendarRequestBuilder) {
    return i0aefe50f098a449e69bee7ad796d3adaf5ca0225915758de5533b5af1cd5d3d4.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*idd3dcbcdf6aec68d0741cf89364a43bc5705b17b2fd86b48348315716154ff77.CancelRequestBuilder) {
    return idd3dcbcdf6aec68d0741cf89364a43bc5705b17b2fd86b48348315716154ff77.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i275121e84c16a3befa700e0246c4d148291c4f5467bc1d593fe2186acd512bd7.DeclineRequestBuilder) {
    return i275121e84c16a3befa700e0246c4d148291c4f5467bc1d593fe2186acd512bd7.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i83a689f67f17f9e9c49b5376e4fb065afee84c3f7057a1ff28da735b7b9bc9e5.DismissReminderRequestBuilder) {
    return i83a689f67f17f9e9c49b5376e4fb065afee84c3f7057a1ff28da735b7b9bc9e5.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i3ba75c90f6a2129e0dededbd49a857df83e8e5652e53088efa2e440031e55677.ExtensionsRequestBuilder) {
    return i3ba75c90f6a2129e0dededbd49a857df83e8e5652e53088efa2e440031e55677.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i5826c74fe1163782016deeb01ad6e7c8edd6129aed51e020012eb5a0a2eceea2.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i5826c74fe1163782016deeb01ad6e7c8edd6129aed51e020012eb5a0a2eceea2.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i6c5f60aa4919c599b5ca2ed4e78e6ceca68b0624738d1eade1a98c2d49222b3b.ForwardRequestBuilder) {
    return i6c5f60aa4919c599b5ca2ed4e78e6ceca68b0624738d1eade1a98c2d49222b3b.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i7a2688b4dd22a5885f9cfdda0103f66e883120819be45676562c7d045c603b17.MultiValueExtendedPropertiesRequestBuilder) {
    return i7a2688b4dd22a5885f9cfdda0103f66e883120819be45676562c7d045c603b17.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ib0dbf672e985e31b3cb13208e4540917f62ae0ea12baeb6cf13df0b957c1d7cb.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ib0dbf672e985e31b3cb13208e4540917f62ae0ea12baeb6cf13df0b957c1d7cb.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i310314089d3422d43c9590127e71f71ad172e68540b68a482ebe308d1133dd25.SingleValueExtendedPropertiesRequestBuilder) {
    return i310314089d3422d43c9590127e71f71ad172e68540b68a482ebe308d1133dd25.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*id611332bbbcace00f410ded8271d027e69d7b2f7f0d882ac5d3ade2aef3c63f6.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return id611332bbbcace00f410ded8271d027e69d7b2f7f0d882ac5d3ade2aef3c63f6.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ia67d11714c9a47717b4e54a81cdbd069dd56de4ff701ebacb389de2d00bb8d9e.SnoozeReminderRequestBuilder) {
    return ia67d11714c9a47717b4e54a81cdbd069dd56de4ff701ebacb389de2d00bb8d9e.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*iae6f0e1bcfce2bf343ad891175d3e665aa98c92b3b7d7dc8dd793365bd904343.TentativelyAcceptRequestBuilder) {
    return iae6f0e1bcfce2bf343ad891175d3e665aa98c92b3b7d7dc8dd793365bd904343.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
