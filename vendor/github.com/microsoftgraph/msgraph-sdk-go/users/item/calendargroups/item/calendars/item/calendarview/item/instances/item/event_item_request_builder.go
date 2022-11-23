package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i129d2b2b7386bbe656bf60028d4f5bea68f14c7529c9ce22900eb9d173d5607f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/decline"
    i5912a36466869b15f1760f5481d7f99751bf0fa1994c8956a162b1c16c782d6d "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/dismissreminder"
    i65cd74bb06865abc6771396e5525947cec5078cfca4ee99362efe3a5e866cc94 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/extensions"
    i7559f3b6e95df335c7bfc1089a4378d407a819b47827e4a3f3152d67a87136d4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/forward"
    i799e30b06602a7df5a17d8afc28cc6a8b103b8165d900f9fce26ba3de4cc0844 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/accept"
    ia1d9d5b7a2a2cfd4026a0df2895f07012da0011fb2fb585c22567a2e4f44ec89 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/tentativelyaccept"
    ia3f03df8544ee71d749c96e3fe2255384fa21dc18cdf1722846800477b7f21e6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/cancel"
    iaf8768970d08e349ba7031c622571c57023cb2a20d30187c3bb368fc6b5539a3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/snoozereminder"
    ic1edefa996b4d1844c76fb2f84a777aa0a65bfd901eccb28186e34ad7e97efc8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/calendar"
    ie6c0dac54e647e42c0d889e3dcfdc2545b660a93ff304ee74b07c315d5fb95ab "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/attachments"
    if75be8d844c5daa95fb1d84905bae6ec1f712cd6ca806ac56b050d068c17c8d7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties"
    ifaac39ed7e498fa2b4b90c7e2365897c36e7d994bfa5612dfcde2a6969648a71 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties"
    i20539128e96a2f01c57fe0224117c0d8ace248f888d69a0c3f769a7f1c7c9b09 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/extensions/item"
    i2a03976fdfd37c9d174c31049686743531c8cc24f38cc9e07945e7c1a7a82806 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties/item"
    i7b23954f818a9b74508bc0a02417108455b8da36b94135440ff2f4b2dac36f46 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/attachments/item"
    i9025c7adc003d2d7747eef3ddca89aaeca8f0b6628378ab2b78bd5962a67b1b6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i799e30b06602a7df5a17d8afc28cc6a8b103b8165d900f9fce26ba3de4cc0844.AcceptRequestBuilder) {
    return i799e30b06602a7df5a17d8afc28cc6a8b103b8165d900f9fce26ba3de4cc0844.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ie6c0dac54e647e42c0d889e3dcfdc2545b660a93ff304ee74b07c315d5fb95ab.AttachmentsRequestBuilder) {
    return ie6c0dac54e647e42c0d889e3dcfdc2545b660a93ff304ee74b07c315d5fb95ab.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i7b23954f818a9b74508bc0a02417108455b8da36b94135440ff2f4b2dac36f46.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i7b23954f818a9b74508bc0a02417108455b8da36b94135440ff2f4b2dac36f46.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ic1edefa996b4d1844c76fb2f84a777aa0a65bfd901eccb28186e34ad7e97efc8.CalendarRequestBuilder) {
    return ic1edefa996b4d1844c76fb2f84a777aa0a65bfd901eccb28186e34ad7e97efc8.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ia3f03df8544ee71d749c96e3fe2255384fa21dc18cdf1722846800477b7f21e6.CancelRequestBuilder) {
    return ia3f03df8544ee71d749c96e3fe2255384fa21dc18cdf1722846800477b7f21e6.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i129d2b2b7386bbe656bf60028d4f5bea68f14c7529c9ce22900eb9d173d5607f.DeclineRequestBuilder) {
    return i129d2b2b7386bbe656bf60028d4f5bea68f14c7529c9ce22900eb9d173d5607f.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i5912a36466869b15f1760f5481d7f99751bf0fa1994c8956a162b1c16c782d6d.DismissReminderRequestBuilder) {
    return i5912a36466869b15f1760f5481d7f99751bf0fa1994c8956a162b1c16c782d6d.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i65cd74bb06865abc6771396e5525947cec5078cfca4ee99362efe3a5e866cc94.ExtensionsRequestBuilder) {
    return i65cd74bb06865abc6771396e5525947cec5078cfca4ee99362efe3a5e866cc94.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i20539128e96a2f01c57fe0224117c0d8ace248f888d69a0c3f769a7f1c7c9b09.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i20539128e96a2f01c57fe0224117c0d8ace248f888d69a0c3f769a7f1c7c9b09.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i7559f3b6e95df335c7bfc1089a4378d407a819b47827e4a3f3152d67a87136d4.ForwardRequestBuilder) {
    return i7559f3b6e95df335c7bfc1089a4378d407a819b47827e4a3f3152d67a87136d4.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ifaac39ed7e498fa2b4b90c7e2365897c36e7d994bfa5612dfcde2a6969648a71.MultiValueExtendedPropertiesRequestBuilder) {
    return ifaac39ed7e498fa2b4b90c7e2365897c36e7d994bfa5612dfcde2a6969648a71.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i9025c7adc003d2d7747eef3ddca89aaeca8f0b6628378ab2b78bd5962a67b1b6.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i9025c7adc003d2d7747eef3ddca89aaeca8f0b6628378ab2b78bd5962a67b1b6.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*if75be8d844c5daa95fb1d84905bae6ec1f712cd6ca806ac56b050d068c17c8d7.SingleValueExtendedPropertiesRequestBuilder) {
    return if75be8d844c5daa95fb1d84905bae6ec1f712cd6ca806ac56b050d068c17c8d7.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i2a03976fdfd37c9d174c31049686743531c8cc24f38cc9e07945e7c1a7a82806.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i2a03976fdfd37c9d174c31049686743531c8cc24f38cc9e07945e7c1a7a82806.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*iaf8768970d08e349ba7031c622571c57023cb2a20d30187c3bb368fc6b5539a3.SnoozeReminderRequestBuilder) {
    return iaf8768970d08e349ba7031c622571c57023cb2a20d30187c3bb368fc6b5539a3.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ia1d9d5b7a2a2cfd4026a0df2895f07012da0011fb2fb585c22567a2e4f44ec89.TentativelyAcceptRequestBuilder) {
    return ia1d9d5b7a2a2cfd4026a0df2895f07012da0011fb2fb585c22567a2e4f44ec89.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
