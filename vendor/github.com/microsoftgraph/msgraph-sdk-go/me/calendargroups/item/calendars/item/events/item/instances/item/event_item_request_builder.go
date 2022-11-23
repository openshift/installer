package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1e9beaaa908cb69f7a15064712a334e45f37f65de40dfb1433c2677711205029 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/multivalueextendedproperties"
    i2b87a17c7347376924744cfb85f3b31622a6c63c8aceb78f038eef745d402c0c "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/extensions"
    i3312e04b15b398a35c0148087fe42d4779883f18181418b03190064a63911d15 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/decline"
    i35a887107143985495a91b95be79e2bacace109cfa36e36612ed39d286109d70 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/singlevalueextendedproperties"
    i3b4ca7974940ae3f515dbeb764f4aa86e2d1da374aa6418905cd563bd947b518 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/cancel"
    i60727935654a93268ee4639394116734dfb4ca9eed7a4b1d844dc5c8aab9da56 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/tentativelyaccept"
    i89b48cfea57c20e3f63d37262611306c9a5be2c23700f71dd45e9545a8847724 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/dismissreminder"
    i9931760ece29c07974ec54bff6c8e4f4c77244ee1ea4518576b391777dddc4e8 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/calendar"
    ib1cb79dbdd106a570acb347cdc29225e77d0dd6092ef90e813f83d298c7f8547 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/forward"
    ibbf7a203781518d029004cc83b027773b2f8cf8dd2e42ecdf53d53798cb92be8 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/accept"
    ieb3a436bb4866a12f9d68d298c21ced9d448dca16df2243c95a5d1afc8f91077 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/snoozereminder"
    iff5117267a51762cffe7680715ec5190d852a4f27f15f1a5bec7c3cb08069805 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/attachments"
    i2186e63cae725f83567a043a8eb82e808c5f5a204fb964c217d7adaaa7be8102 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/attachments/item"
    i525225bd2be0b8d861ed0c7ba4972b6ba48c03f453e0728d9543e48ffd8b7191 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/multivalueextendedproperties/item"
    i6ef2ced592c6090dd781593f3a59a7eed9230cdda657796584eb7301399c6a9f "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/singlevalueextendedproperties/item"
    ied5a4eff0211de3e6a71f8605b5ea559bf9c4a9b5b75aa4c11e4195b351da06c "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/events/item/instances/item/extensions/item"
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
func (m *EventItemRequestBuilder) Accept()(*ibbf7a203781518d029004cc83b027773b2f8cf8dd2e42ecdf53d53798cb92be8.AcceptRequestBuilder) {
    return ibbf7a203781518d029004cc83b027773b2f8cf8dd2e42ecdf53d53798cb92be8.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*iff5117267a51762cffe7680715ec5190d852a4f27f15f1a5bec7c3cb08069805.AttachmentsRequestBuilder) {
    return iff5117267a51762cffe7680715ec5190d852a4f27f15f1a5bec7c3cb08069805.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i2186e63cae725f83567a043a8eb82e808c5f5a204fb964c217d7adaaa7be8102.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i2186e63cae725f83567a043a8eb82e808c5f5a204fb964c217d7adaaa7be8102.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i9931760ece29c07974ec54bff6c8e4f4c77244ee1ea4518576b391777dddc4e8.CalendarRequestBuilder) {
    return i9931760ece29c07974ec54bff6c8e4f4c77244ee1ea4518576b391777dddc4e8.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i3b4ca7974940ae3f515dbeb764f4aa86e2d1da374aa6418905cd563bd947b518.CancelRequestBuilder) {
    return i3b4ca7974940ae3f515dbeb764f4aa86e2d1da374aa6418905cd563bd947b518.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i3312e04b15b398a35c0148087fe42d4779883f18181418b03190064a63911d15.DeclineRequestBuilder) {
    return i3312e04b15b398a35c0148087fe42d4779883f18181418b03190064a63911d15.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i89b48cfea57c20e3f63d37262611306c9a5be2c23700f71dd45e9545a8847724.DismissReminderRequestBuilder) {
    return i89b48cfea57c20e3f63d37262611306c9a5be2c23700f71dd45e9545a8847724.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i2b87a17c7347376924744cfb85f3b31622a6c63c8aceb78f038eef745d402c0c.ExtensionsRequestBuilder) {
    return i2b87a17c7347376924744cfb85f3b31622a6c63c8aceb78f038eef745d402c0c.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ied5a4eff0211de3e6a71f8605b5ea559bf9c4a9b5b75aa4c11e4195b351da06c.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ied5a4eff0211de3e6a71f8605b5ea559bf9c4a9b5b75aa4c11e4195b351da06c.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ib1cb79dbdd106a570acb347cdc29225e77d0dd6092ef90e813f83d298c7f8547.ForwardRequestBuilder) {
    return ib1cb79dbdd106a570acb347cdc29225e77d0dd6092ef90e813f83d298c7f8547.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i1e9beaaa908cb69f7a15064712a334e45f37f65de40dfb1433c2677711205029.MultiValueExtendedPropertiesRequestBuilder) {
    return i1e9beaaa908cb69f7a15064712a334e45f37f65de40dfb1433c2677711205029.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i525225bd2be0b8d861ed0c7ba4972b6ba48c03f453e0728d9543e48ffd8b7191.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i525225bd2be0b8d861ed0c7ba4972b6ba48c03f453e0728d9543e48ffd8b7191.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i35a887107143985495a91b95be79e2bacace109cfa36e36612ed39d286109d70.SingleValueExtendedPropertiesRequestBuilder) {
    return i35a887107143985495a91b95be79e2bacace109cfa36e36612ed39d286109d70.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i6ef2ced592c6090dd781593f3a59a7eed9230cdda657796584eb7301399c6a9f.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i6ef2ced592c6090dd781593f3a59a7eed9230cdda657796584eb7301399c6a9f.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ieb3a436bb4866a12f9d68d298c21ced9d448dca16df2243c95a5d1afc8f91077.SnoozeReminderRequestBuilder) {
    return ieb3a436bb4866a12f9d68d298c21ced9d448dca16df2243c95a5d1afc8f91077.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i60727935654a93268ee4639394116734dfb4ca9eed7a4b1d844dc5c8aab9da56.TentativelyAcceptRequestBuilder) {
    return i60727935654a93268ee4639394116734dfb4ca9eed7a4b1d844dc5c8aab9da56.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
