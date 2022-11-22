package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i13d9d853f6f596efbe21ea18f40d0592475a8ac990e8a695d81a6fbba14f4d26 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/cancel"
    i1a10f47c94aaf40994f428a074851ae6a632b047055a7f1cf3bdd2b7c7b6c0a5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/singlevalueextendedproperties"
    i1b64d030401b44fa9d8e6c1050f609bba7b33042a59456f5e5f4bf592c0d9b93 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/snoozereminder"
    i2e1f9fb6aaa8abf5f8df5e11ba377fe0143a5e43f512624a8d9cd6fb27421a78 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/dismissreminder"
    i91354a288a328a774d7e175c52b2d23e26aa9ea2c90072fa527d0a581ca4bc8e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/tentativelyaccept"
    ia088d5c6b1f3a9f9a2aeaad3b76325dbf4d601ee30542e1676570dbf971be3f8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/decline"
    iaf599fb4f2ae1ccff34614f5a03bcc714f19cbff2c171838f87122ed74aecb31 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances"
    ica7ba27254da644cf113708d63a508c5069debc7ee258b273bee580a0147682d "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/multivalueextendedproperties"
    idd8083a89a62c4a90c13d7095b22cbdaa8ecd14901ed17edda30c0d6194c4c4f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/calendar"
    ie21405766b71afd03cdca486d941e721d602f890de39f16cfb24032bae39294b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/extensions"
    ie8336fa6b5a4ea27b8de9a04a40237f8ebccb1dfe6f97770ce15b89d6c3e2f95 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/attachments"
    ie8d2a5a2eab16f2219d42c95df3d68b757fb3773211e5b6ae612f6646a67a1b1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/forward"
    if3456a6c7e36e0e413ac6ae411ab193961324f5aa72c986631a16de60564c13e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/accept"
    i089a2accd06c2c34e85f4442bcfe40065651bb4f57e23851aa07d1ea2becdc21 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/extensions/item"
    i0be4955da7d4a2edcd0b70a9d0b170efd62f033d7459723d8439132a58f846c9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item"
    iadb079fbc3d0d06ba36c2c0c66d5f3073b14e52cd15b1db7a9bd6befec932fa3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/multivalueextendedproperties/item"
    ib451c2fc927d457e02148536325b3839f57ef91bec0997dd64fb4888db642574 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/singlevalueextendedproperties/item"
    iee3b48332651e2bc68d677bf62b811a958d2bac5d0f1786cad3f5eece9d248f8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*if3456a6c7e36e0e413ac6ae411ab193961324f5aa72c986631a16de60564c13e.AcceptRequestBuilder) {
    return if3456a6c7e36e0e413ac6ae411ab193961324f5aa72c986631a16de60564c13e.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ie8336fa6b5a4ea27b8de9a04a40237f8ebccb1dfe6f97770ce15b89d6c3e2f95.AttachmentsRequestBuilder) {
    return ie8336fa6b5a4ea27b8de9a04a40237f8ebccb1dfe6f97770ce15b89d6c3e2f95.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*iee3b48332651e2bc68d677bf62b811a958d2bac5d0f1786cad3f5eece9d248f8.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return iee3b48332651e2bc68d677bf62b811a958d2bac5d0f1786cad3f5eece9d248f8.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*idd8083a89a62c4a90c13d7095b22cbdaa8ecd14901ed17edda30c0d6194c4c4f.CalendarRequestBuilder) {
    return idd8083a89a62c4a90c13d7095b22cbdaa8ecd14901ed17edda30c0d6194c4c4f.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i13d9d853f6f596efbe21ea18f40d0592475a8ac990e8a695d81a6fbba14f4d26.CancelRequestBuilder) {
    return i13d9d853f6f596efbe21ea18f40d0592475a8ac990e8a695d81a6fbba14f4d26.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/calendarView/{event%2Did}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*ia088d5c6b1f3a9f9a2aeaad3b76325dbf4d601ee30542e1676570dbf971be3f8.DeclineRequestBuilder) {
    return ia088d5c6b1f3a9f9a2aeaad3b76325dbf4d601ee30542e1676570dbf971be3f8.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i2e1f9fb6aaa8abf5f8df5e11ba377fe0143a5e43f512624a8d9cd6fb27421a78.DismissReminderRequestBuilder) {
    return i2e1f9fb6aaa8abf5f8df5e11ba377fe0143a5e43f512624a8d9cd6fb27421a78.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*ie21405766b71afd03cdca486d941e721d602f890de39f16cfb24032bae39294b.ExtensionsRequestBuilder) {
    return ie21405766b71afd03cdca486d941e721d602f890de39f16cfb24032bae39294b.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i089a2accd06c2c34e85f4442bcfe40065651bb4f57e23851aa07d1ea2becdc21.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i089a2accd06c2c34e85f4442bcfe40065651bb4f57e23851aa07d1ea2becdc21.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ie8d2a5a2eab16f2219d42c95df3d68b757fb3773211e5b6ae612f6646a67a1b1.ForwardRequestBuilder) {
    return ie8d2a5a2eab16f2219d42c95df3d68b757fb3773211e5b6ae612f6646a67a1b1.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*iaf599fb4f2ae1ccff34614f5a03bcc714f19cbff2c171838f87122ed74aecb31.InstancesRequestBuilder) {
    return iaf599fb4f2ae1ccff34614f5a03bcc714f19cbff2c171838f87122ed74aecb31.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i0be4955da7d4a2edcd0b70a9d0b170efd62f033d7459723d8439132a58f846c9.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i0be4955da7d4a2edcd0b70a9d0b170efd62f033d7459723d8439132a58f846c9.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ica7ba27254da644cf113708d63a508c5069debc7ee258b273bee580a0147682d.MultiValueExtendedPropertiesRequestBuilder) {
    return ica7ba27254da644cf113708d63a508c5069debc7ee258b273bee580a0147682d.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*iadb079fbc3d0d06ba36c2c0c66d5f3073b14e52cd15b1db7a9bd6befec932fa3.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return iadb079fbc3d0d06ba36c2c0c66d5f3073b14e52cd15b1db7a9bd6befec932fa3.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i1a10f47c94aaf40994f428a074851ae6a632b047055a7f1cf3bdd2b7c7b6c0a5.SingleValueExtendedPropertiesRequestBuilder) {
    return i1a10f47c94aaf40994f428a074851ae6a632b047055a7f1cf3bdd2b7c7b6c0a5.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ib451c2fc927d457e02148536325b3839f57ef91bec0997dd64fb4888db642574.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ib451c2fc927d457e02148536325b3839f57ef91bec0997dd64fb4888db642574.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i1b64d030401b44fa9d8e6c1050f609bba7b33042a59456f5e5f4bf592c0d9b93.SnoozeReminderRequestBuilder) {
    return i1b64d030401b44fa9d8e6c1050f609bba7b33042a59456f5e5f4bf592c0d9b93.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i91354a288a328a774d7e175c52b2d23e26aa9ea2c90072fa527d0a581ca4bc8e.TentativelyAcceptRequestBuilder) {
    return i91354a288a328a774d7e175c52b2d23e26aa9ea2c90072fa527d0a581ca4bc8e.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
