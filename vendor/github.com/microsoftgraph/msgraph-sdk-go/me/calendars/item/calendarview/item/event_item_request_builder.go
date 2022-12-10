package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0f0ce6d736b6f504b471613c5ac15739a55461fedfe57875a5750f9b53c0f226 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/decline"
    i1a18a6b6a39e454f6c8d96179b980f6da06555b91b37580c9ec8ba0305e6de5f "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/extensions"
    i1c6a239e026b8cb9867b574711b33427f6c5744c2f14261d944366070fc353e5 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/calendar"
    i1ea6a33ab7759770447ac7a4d978f10017eacb31c6f0871731eeb148649fb5bc "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/tentativelyaccept"
    i26977c51a41815fe0d7ff9cbafcc8af4a348722bd293191d771b43d03047b305 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/singlevalueextendedproperties"
    i566f9c16a6fe85b63ce8dbb3f2fea3693448cd4a91cd02b2e3a7e71a63267860 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/instances"
    i77eb73a898b6f8d95bb82e9c2efcf03f741c8f2f7e7176dde95d01eed01893c0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/attachments"
    i78a836c35965e9975e5e331baffe1a63b8525eef66bb1b8645b6e1e39525bef3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/accept"
    i829fbd366d6ff12b081486076f81163508113569bf85669ed18cd726a41affc8 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/forward"
    ib0d55d703ff4d551e433c43443d5dd4e5ee73676034a7c9cb395c174c3e3c434 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/cancel"
    ie6bafb4b17c2a9ac6f86bde54240348db1f716dae601c3aaa3ba5246a106e2c7 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/multivalueextendedproperties"
    ie7a0d6568144c2179c8661e11538ef71f2732b97fd1f89e98eaa928afb3f952e "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/dismissreminder"
    ifc94e234d0cd341ab21ed1133e622aa778dde1f8670fa77373fecf95bb710d09 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/snoozereminder"
    i435b79a3073f9e3b67954c2f479ceba4b84b060beffbf4eddfdcd11b9c0afdd6 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/multivalueextendedproperties/item"
    i5c4cc58062b30213bec3fe8c70fe2c7eba1480b941dd3c248011fe0fb3c4dcf9 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/extensions/item"
    i9d4b2dde6d1a46763d115eb03242b93910eb5891036290fb5e8a49322ead62da "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/instances/item"
    i9e8c599582a73f84f65086c644266f1ecd869813360b75664d307bffb48d75bc "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/attachments/item"
    if591e3f43c598d182e086c7492b69164db7e06e4a23e37dd47fab1c1dfa49207 "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/calendarview/item/singlevalueextendedproperties/item"
)

// EventItemRequestBuilder provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderGetQueryParameters the calendar view for the calendar. Navigation property. Read-only.
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
func (m *EventItemRequestBuilder) Accept()(*i78a836c35965e9975e5e331baffe1a63b8525eef66bb1b8645b6e1e39525bef3.AcceptRequestBuilder) {
    return i78a836c35965e9975e5e331baffe1a63b8525eef66bb1b8645b6e1e39525bef3.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i77eb73a898b6f8d95bb82e9c2efcf03f741c8f2f7e7176dde95d01eed01893c0.AttachmentsRequestBuilder) {
    return i77eb73a898b6f8d95bb82e9c2efcf03f741c8f2f7e7176dde95d01eed01893c0.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i9e8c599582a73f84f65086c644266f1ecd869813360b75664d307bffb48d75bc.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i9e8c599582a73f84f65086c644266f1ecd869813360b75664d307bffb48d75bc.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i1c6a239e026b8cb9867b574711b33427f6c5744c2f14261d944366070fc353e5.CalendarRequestBuilder) {
    return i1c6a239e026b8cb9867b574711b33427f6c5744c2f14261d944366070fc353e5.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ib0d55d703ff4d551e433c43443d5dd4e5ee73676034a7c9cb395c174c3e3c434.CancelRequestBuilder) {
    return ib0d55d703ff4d551e433c43443d5dd4e5ee73676034a7c9cb395c174c3e3c434.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendars/{calendar%2Did}/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
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
// CreateGetRequestInformation the calendar view for the calendar. Navigation property. Read-only.
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
func (m *EventItemRequestBuilder) Decline()(*i0f0ce6d736b6f504b471613c5ac15739a55461fedfe57875a5750f9b53c0f226.DeclineRequestBuilder) {
    return i0f0ce6d736b6f504b471613c5ac15739a55461fedfe57875a5750f9b53c0f226.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*ie7a0d6568144c2179c8661e11538ef71f2732b97fd1f89e98eaa928afb3f952e.DismissReminderRequestBuilder) {
    return ie7a0d6568144c2179c8661e11538ef71f2732b97fd1f89e98eaa928afb3f952e.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i1a18a6b6a39e454f6c8d96179b980f6da06555b91b37580c9ec8ba0305e6de5f.ExtensionsRequestBuilder) {
    return i1a18a6b6a39e454f6c8d96179b980f6da06555b91b37580c9ec8ba0305e6de5f.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i5c4cc58062b30213bec3fe8c70fe2c7eba1480b941dd3c248011fe0fb3c4dcf9.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i5c4cc58062b30213bec3fe8c70fe2c7eba1480b941dd3c248011fe0fb3c4dcf9.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i829fbd366d6ff12b081486076f81163508113569bf85669ed18cd726a41affc8.ForwardRequestBuilder) {
    return i829fbd366d6ff12b081486076f81163508113569bf85669ed18cd726a41affc8.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the calendar view for the calendar. Navigation property. Read-only.
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
func (m *EventItemRequestBuilder) Instances()(*i566f9c16a6fe85b63ce8dbb3f2fea3693448cd4a91cd02b2e3a7e71a63267860.InstancesRequestBuilder) {
    return i566f9c16a6fe85b63ce8dbb3f2fea3693448cd4a91cd02b2e3a7e71a63267860.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i9d4b2dde6d1a46763d115eb03242b93910eb5891036290fb5e8a49322ead62da.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i9d4b2dde6d1a46763d115eb03242b93910eb5891036290fb5e8a49322ead62da.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ie6bafb4b17c2a9ac6f86bde54240348db1f716dae601c3aaa3ba5246a106e2c7.MultiValueExtendedPropertiesRequestBuilder) {
    return ie6bafb4b17c2a9ac6f86bde54240348db1f716dae601c3aaa3ba5246a106e2c7.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i435b79a3073f9e3b67954c2f479ceba4b84b060beffbf4eddfdcd11b9c0afdd6.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i435b79a3073f9e3b67954c2f479ceba4b84b060beffbf4eddfdcd11b9c0afdd6.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i26977c51a41815fe0d7ff9cbafcc8af4a348722bd293191d771b43d03047b305.SingleValueExtendedPropertiesRequestBuilder) {
    return i26977c51a41815fe0d7ff9cbafcc8af4a348722bd293191d771b43d03047b305.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*if591e3f43c598d182e086c7492b69164db7e06e4a23e37dd47fab1c1dfa49207.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return if591e3f43c598d182e086c7492b69164db7e06e4a23e37dd47fab1c1dfa49207.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ifc94e234d0cd341ab21ed1133e622aa778dde1f8670fa77373fecf95bb710d09.SnoozeReminderRequestBuilder) {
    return ifc94e234d0cd341ab21ed1133e622aa778dde1f8670fa77373fecf95bb710d09.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i1ea6a33ab7759770447ac7a4d978f10017eacb31c6f0871731eeb148649fb5bc.TentativelyAcceptRequestBuilder) {
    return i1ea6a33ab7759770447ac7a4d978f10017eacb31c6f0871731eeb148649fb5bc.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
