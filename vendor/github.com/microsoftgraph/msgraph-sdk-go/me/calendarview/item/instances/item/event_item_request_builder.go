package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i047511a86dd96f86390856dcdd303d1e4d376a1bde031c6f00cddbe54d732022 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/cancel"
    i1513186621a8d49c87229c6e5d07a633fa56b489d4246d85f233dac97f90f5c3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/multivalueextendedproperties"
    i1c9ac4c6283505a22162c7f8690d6acec224ca3b68f6506ada58952467ec4abd "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/extensions"
    i3d51a594762fc074bd2e62e1f226ab21b1fca1bd1c595a68fe36050700268eb2 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/tentativelyaccept"
    i4fa8ca9f998bd5bd29718077b44e5de2b34431251ee75be9d9b92b8cd57afa91 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/forward"
    i5a877e113fc35bfbd07e814c73eed69ca38344d277aeaaed4fae0c75df54dc77 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/attachments"
    i89c6c9192c1cb49df5ef36b41d2a7f1e93307dfcd982548893ddcd3cf40479be "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/dismissreminder"
    ia2e0ca70ae4e6d806f222d053141e3a24e5c2d07ef2558b3d8182d5515061330 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/accept"
    iaf925c938d533ffb8df602f79ec606ba61eeb812f89f6027b216a163fcffed15 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/decline"
    ic94adb8229abb79a5fcfd407ef062559e0a6b12b377a5a5d750cc8b0ba1c41f2 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/singlevalueextendedproperties"
    id09daa0d64b25ca646cfc102a53100dca3cd3bc9af3340a6e85dd4213f60cc98 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/snoozereminder"
    id3106db13673225b97d26dceb41ad56217d4923df53ed71a326ab20a7c9cb61d "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/calendar"
    i7a90d7b78e805e2b7f1f2b753e687c92db71b52a73f98a3df4aadfe124576c19 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/extensions/item"
    ic298ea33bf31e4be0fe59a4f20297a719e5271914cc2bf1653465e3af505f8f7 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/multivalueextendedproperties/item"
    ic7d6728a237849a0f5a5965e0501bda413d23001139934b60b5a74a0713e232a "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/singlevalueextendedproperties/item"
    iebeba9115c421074573ae4ec2fe90bf71ca0a26697cfea0123d7e422578a9b52 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*ia2e0ca70ae4e6d806f222d053141e3a24e5c2d07ef2558b3d8182d5515061330.AcceptRequestBuilder) {
    return ia2e0ca70ae4e6d806f222d053141e3a24e5c2d07ef2558b3d8182d5515061330.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i5a877e113fc35bfbd07e814c73eed69ca38344d277aeaaed4fae0c75df54dc77.AttachmentsRequestBuilder) {
    return i5a877e113fc35bfbd07e814c73eed69ca38344d277aeaaed4fae0c75df54dc77.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*iebeba9115c421074573ae4ec2fe90bf71ca0a26697cfea0123d7e422578a9b52.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return iebeba9115c421074573ae4ec2fe90bf71ca0a26697cfea0123d7e422578a9b52.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*id3106db13673225b97d26dceb41ad56217d4923df53ed71a326ab20a7c9cb61d.CalendarRequestBuilder) {
    return id3106db13673225b97d26dceb41ad56217d4923df53ed71a326ab20a7c9cb61d.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i047511a86dd96f86390856dcdd303d1e4d376a1bde031c6f00cddbe54d732022.CancelRequestBuilder) {
    return i047511a86dd96f86390856dcdd303d1e4d376a1bde031c6f00cddbe54d732022.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*iaf925c938d533ffb8df602f79ec606ba61eeb812f89f6027b216a163fcffed15.DeclineRequestBuilder) {
    return iaf925c938d533ffb8df602f79ec606ba61eeb812f89f6027b216a163fcffed15.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i89c6c9192c1cb49df5ef36b41d2a7f1e93307dfcd982548893ddcd3cf40479be.DismissReminderRequestBuilder) {
    return i89c6c9192c1cb49df5ef36b41d2a7f1e93307dfcd982548893ddcd3cf40479be.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i1c9ac4c6283505a22162c7f8690d6acec224ca3b68f6506ada58952467ec4abd.ExtensionsRequestBuilder) {
    return i1c9ac4c6283505a22162c7f8690d6acec224ca3b68f6506ada58952467ec4abd.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i7a90d7b78e805e2b7f1f2b753e687c92db71b52a73f98a3df4aadfe124576c19.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i7a90d7b78e805e2b7f1f2b753e687c92db71b52a73f98a3df4aadfe124576c19.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i4fa8ca9f998bd5bd29718077b44e5de2b34431251ee75be9d9b92b8cd57afa91.ForwardRequestBuilder) {
    return i4fa8ca9f998bd5bd29718077b44e5de2b34431251ee75be9d9b92b8cd57afa91.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i1513186621a8d49c87229c6e5d07a633fa56b489d4246d85f233dac97f90f5c3.MultiValueExtendedPropertiesRequestBuilder) {
    return i1513186621a8d49c87229c6e5d07a633fa56b489d4246d85f233dac97f90f5c3.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ic298ea33bf31e4be0fe59a4f20297a719e5271914cc2bf1653465e3af505f8f7.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ic298ea33bf31e4be0fe59a4f20297a719e5271914cc2bf1653465e3af505f8f7.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*ic94adb8229abb79a5fcfd407ef062559e0a6b12b377a5a5d750cc8b0ba1c41f2.SingleValueExtendedPropertiesRequestBuilder) {
    return ic94adb8229abb79a5fcfd407ef062559e0a6b12b377a5a5d750cc8b0ba1c41f2.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ic7d6728a237849a0f5a5965e0501bda413d23001139934b60b5a74a0713e232a.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ic7d6728a237849a0f5a5965e0501bda413d23001139934b60b5a74a0713e232a.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*id09daa0d64b25ca646cfc102a53100dca3cd3bc9af3340a6e85dd4213f60cc98.SnoozeReminderRequestBuilder) {
    return id09daa0d64b25ca646cfc102a53100dca3cd3bc9af3340a6e85dd4213f60cc98.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i3d51a594762fc074bd2e62e1f226ab21b1fca1bd1c595a68fe36050700268eb2.TentativelyAcceptRequestBuilder) {
    return i3d51a594762fc074bd2e62e1f226ab21b1fca1bd1c595a68fe36050700268eb2.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
