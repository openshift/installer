package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i145fb7990a2e3fdab872136aac6bb8be775248efc831fc0ceac357c425bae97f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/decline"
    i263c92c411c89f2be6a8e54de5fc9a4dde982d0e6d8de040c71d4afe63a7748c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/cancel"
    i580c1b6081b469d14dea10fa0ff4e3fcf1e92e486940fe417a29f9f5671bbb53 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances"
    i896d4bd262b2cdabe891cd22c1b08cbcde12e3b522ff979daf20e82ad1ea711b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/snoozereminder"
    i9e33412fe1bdfa6c018b792dd1429c9276d49c06fe8c88117fede514ab522720 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/forward"
    ia02addaa9051893cf138a732042bf251456286968b685c328064520251cb0cf3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/extensions"
    ia3dd18af0ea5a7f297b2a2edda6570fc313fc5708c7e2634ba8259e0ae32b2a5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/tentativelyaccept"
    ia62de877959a200aca4887792a575daf1a8f28bf9cfcb7a1fff2251ff773c997 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/accept"
    ia93978b34743023ce443a9997d7630ed381895aed8d859d3f4ac245be62fb849 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/multivalueextendedproperties"
    iaae9964b69928703e801f04d74fc2a4f7d6498b08c64baf44a7eb2b78040bcfc "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/calendar"
    iae0f8cfe0219a22eb2de5c896f00609e7fb3fb5d12dfd6203fee1bdcc65508d7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/dismissreminder"
    ib55ea188c0bda46cccf674fa54b66612f2fa73c5e8d49e10718ad5878343eb9f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/attachments"
    id8bc39608b45df9d2acb9059f8be81c090091fe783571dbbd283d47b33d831d5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/singlevalueextendedproperties"
    i01f8c58bdadd020940a9f0c219ace3903f9dd0225f07ec8fe1fdfc9eb6c2c0f8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/multivalueextendedproperties/item"
    i061c7c8e6a1b9d55c8cc8301b7c4370b5310fd0b6c3a2c8c5217710a38ec869b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/extensions/item"
    i57a0e5d0a9a8aa62edc5ca53b30b8ad23c80c88a7af4e91c43ef9504c35fd4a4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/singlevalueextendedproperties/item"
    i66e4305766d125651f3fd13dcddaba80315e4ca5673cb2b5eda12fdabe6174fb "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/instances/item"
    ia402890f08866a24e37023bade2f4fd90979ed0d599b4b006b83fd80bee80d01 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*ia62de877959a200aca4887792a575daf1a8f28bf9cfcb7a1fff2251ff773c997.AcceptRequestBuilder) {
    return ia62de877959a200aca4887792a575daf1a8f28bf9cfcb7a1fff2251ff773c997.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ib55ea188c0bda46cccf674fa54b66612f2fa73c5e8d49e10718ad5878343eb9f.AttachmentsRequestBuilder) {
    return ib55ea188c0bda46cccf674fa54b66612f2fa73c5e8d49e10718ad5878343eb9f.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ia402890f08866a24e37023bade2f4fd90979ed0d599b4b006b83fd80bee80d01.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ia402890f08866a24e37023bade2f4fd90979ed0d599b4b006b83fd80bee80d01.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*iaae9964b69928703e801f04d74fc2a4f7d6498b08c64baf44a7eb2b78040bcfc.CalendarRequestBuilder) {
    return iaae9964b69928703e801f04d74fc2a4f7d6498b08c64baf44a7eb2b78040bcfc.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i263c92c411c89f2be6a8e54de5fc9a4dde982d0e6d8de040c71d4afe63a7748c.CancelRequestBuilder) {
    return i263c92c411c89f2be6a8e54de5fc9a4dde982d0e6d8de040c71d4afe63a7748c.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendar/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i145fb7990a2e3fdab872136aac6bb8be775248efc831fc0ceac357c425bae97f.DeclineRequestBuilder) {
    return i145fb7990a2e3fdab872136aac6bb8be775248efc831fc0ceac357c425bae97f.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*iae0f8cfe0219a22eb2de5c896f00609e7fb3fb5d12dfd6203fee1bdcc65508d7.DismissReminderRequestBuilder) {
    return iae0f8cfe0219a22eb2de5c896f00609e7fb3fb5d12dfd6203fee1bdcc65508d7.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*ia02addaa9051893cf138a732042bf251456286968b685c328064520251cb0cf3.ExtensionsRequestBuilder) {
    return ia02addaa9051893cf138a732042bf251456286968b685c328064520251cb0cf3.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i061c7c8e6a1b9d55c8cc8301b7c4370b5310fd0b6c3a2c8c5217710a38ec869b.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i061c7c8e6a1b9d55c8cc8301b7c4370b5310fd0b6c3a2c8c5217710a38ec869b.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i9e33412fe1bdfa6c018b792dd1429c9276d49c06fe8c88117fede514ab522720.ForwardRequestBuilder) {
    return i9e33412fe1bdfa6c018b792dd1429c9276d49c06fe8c88117fede514ab522720.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*i580c1b6081b469d14dea10fa0ff4e3fcf1e92e486940fe417a29f9f5671bbb53.InstancesRequestBuilder) {
    return i580c1b6081b469d14dea10fa0ff4e3fcf1e92e486940fe417a29f9f5671bbb53.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i66e4305766d125651f3fd13dcddaba80315e4ca5673cb2b5eda12fdabe6174fb.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i66e4305766d125651f3fd13dcddaba80315e4ca5673cb2b5eda12fdabe6174fb.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ia93978b34743023ce443a9997d7630ed381895aed8d859d3f4ac245be62fb849.MultiValueExtendedPropertiesRequestBuilder) {
    return ia93978b34743023ce443a9997d7630ed381895aed8d859d3f4ac245be62fb849.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i01f8c58bdadd020940a9f0c219ace3903f9dd0225f07ec8fe1fdfc9eb6c2c0f8.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i01f8c58bdadd020940a9f0c219ace3903f9dd0225f07ec8fe1fdfc9eb6c2c0f8.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*id8bc39608b45df9d2acb9059f8be81c090091fe783571dbbd283d47b33d831d5.SingleValueExtendedPropertiesRequestBuilder) {
    return id8bc39608b45df9d2acb9059f8be81c090091fe783571dbbd283d47b33d831d5.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i57a0e5d0a9a8aa62edc5ca53b30b8ad23c80c88a7af4e91c43ef9504c35fd4a4.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i57a0e5d0a9a8aa62edc5ca53b30b8ad23c80c88a7af4e91c43ef9504c35fd4a4.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i896d4bd262b2cdabe891cd22c1b08cbcde12e3b522ff979daf20e82ad1ea711b.SnoozeReminderRequestBuilder) {
    return i896d4bd262b2cdabe891cd22c1b08cbcde12e3b522ff979daf20e82ad1ea711b.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ia3dd18af0ea5a7f297b2a2edda6570fc313fc5708c7e2634ba8259e0ae32b2a5.TentativelyAcceptRequestBuilder) {
    return ia3dd18af0ea5a7f297b2a2edda6570fc313fc5708c7e2634ba8259e0ae32b2a5.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
