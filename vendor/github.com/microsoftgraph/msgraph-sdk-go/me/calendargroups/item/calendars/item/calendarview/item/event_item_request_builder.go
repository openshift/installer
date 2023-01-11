package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i034f075a9742087cc62fa2fa3c2d823b0c34579b549387b0217d091d4a084aed "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/singlevalueextendedproperties"
    i0d1833fafe91c99d46532bdcc396354079ec6601492e051fcf683e98f02b4aee "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/decline"
    i322baabcc84e3997da553dda536125328af71f44c1c87052ffeafb14a6877592 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/cancel"
    i42aeafc43ff4f3886ddd2c06a4c17d4ba4cc3b7fd3216d6c353debc41b3d4d14 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/calendar"
    i554abe6eb0940df39daf173fba7aa2976e43e0e1cde74ba38fb1ad212c524b1e "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/dismissreminder"
    i7e9c02a874a50f511f5d63885525677f008606e724bca4866861a307c5f4886a "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/extensions"
    i85e1eadeddf7ec2c307967b39d84d66cd802ef3edd7cea67e3fcdd2e9efae123 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/attachments"
    i8bd23f6c7e0c42bea79bd661447dc93d8e0ada59e8bba7ad08a5783f58f6ad68 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/multivalueextendedproperties"
    i94dcdd3851da0f02fa640225ba331bea6617b62727e0061c6bb8e9ab386b8629 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/accept"
    ibcb3c8ba4cd7beda528fefdc4942827a7ca6e5c6b961928da585dfc095291fb4 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/forward"
    ic66f697cfef967cd4effb4c9d92b29e287a4def44f9b2aea82d334789318299f "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/snoozereminder"
    ie4a846e389b081fa82f31b8581ae8dfed0401f4a81a686e4cc965b82d9606851 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/tentativelyaccept"
    ieb9e6dcb9bcea88f684ee7112127bd4d1df6deb986fab653d792f9b1389cc795 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances"
    i4e4fe41f8ec4b0f84b1c21baa5931909bd57fcb9768d3be13b17583642944ef8 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/multivalueextendedproperties/item"
    i62e0076ba3224fb44507e4ef678645cbc8c53f273bb42881cc65121616fbd8f0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/attachments/item"
    ie57987dfe2c89ca06fcc0fc5f474e14183f8c3cec8825b4b733c1da6a5d7ea6a "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/extensions/item"
    if26639b4bc627d0cd7049c7f814bf3f0458856e94fc4fce0c5d95011495f7257 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item"
    if7192229b65a228ea27ede8842c1b7d3638ac8120f07a51c43c8174a1920dd14 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/singlevalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i94dcdd3851da0f02fa640225ba331bea6617b62727e0061c6bb8e9ab386b8629.AcceptRequestBuilder) {
    return i94dcdd3851da0f02fa640225ba331bea6617b62727e0061c6bb8e9ab386b8629.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i85e1eadeddf7ec2c307967b39d84d66cd802ef3edd7cea67e3fcdd2e9efae123.AttachmentsRequestBuilder) {
    return i85e1eadeddf7ec2c307967b39d84d66cd802ef3edd7cea67e3fcdd2e9efae123.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i62e0076ba3224fb44507e4ef678645cbc8c53f273bb42881cc65121616fbd8f0.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i62e0076ba3224fb44507e4ef678645cbc8c53f273bb42881cc65121616fbd8f0.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i42aeafc43ff4f3886ddd2c06a4c17d4ba4cc3b7fd3216d6c353debc41b3d4d14.CalendarRequestBuilder) {
    return i42aeafc43ff4f3886ddd2c06a4c17d4ba4cc3b7fd3216d6c353debc41b3d4d14.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i322baabcc84e3997da553dda536125328af71f44c1c87052ffeafb14a6877592.CancelRequestBuilder) {
    return i322baabcc84e3997da553dda536125328af71f44c1c87052ffeafb14a6877592.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/calendarView/{event%2Did}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i0d1833fafe91c99d46532bdcc396354079ec6601492e051fcf683e98f02b4aee.DeclineRequestBuilder) {
    return i0d1833fafe91c99d46532bdcc396354079ec6601492e051fcf683e98f02b4aee.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i554abe6eb0940df39daf173fba7aa2976e43e0e1cde74ba38fb1ad212c524b1e.DismissReminderRequestBuilder) {
    return i554abe6eb0940df39daf173fba7aa2976e43e0e1cde74ba38fb1ad212c524b1e.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i7e9c02a874a50f511f5d63885525677f008606e724bca4866861a307c5f4886a.ExtensionsRequestBuilder) {
    return i7e9c02a874a50f511f5d63885525677f008606e724bca4866861a307c5f4886a.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ie57987dfe2c89ca06fcc0fc5f474e14183f8c3cec8825b4b733c1da6a5d7ea6a.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ie57987dfe2c89ca06fcc0fc5f474e14183f8c3cec8825b4b733c1da6a5d7ea6a.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ibcb3c8ba4cd7beda528fefdc4942827a7ca6e5c6b961928da585dfc095291fb4.ForwardRequestBuilder) {
    return ibcb3c8ba4cd7beda528fefdc4942827a7ca6e5c6b961928da585dfc095291fb4.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*ieb9e6dcb9bcea88f684ee7112127bd4d1df6deb986fab653d792f9b1389cc795.InstancesRequestBuilder) {
    return ieb9e6dcb9bcea88f684ee7112127bd4d1df6deb986fab653d792f9b1389cc795.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*if26639b4bc627d0cd7049c7f814bf3f0458856e94fc4fce0c5d95011495f7257.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return if26639b4bc627d0cd7049c7f814bf3f0458856e94fc4fce0c5d95011495f7257.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i8bd23f6c7e0c42bea79bd661447dc93d8e0ada59e8bba7ad08a5783f58f6ad68.MultiValueExtendedPropertiesRequestBuilder) {
    return i8bd23f6c7e0c42bea79bd661447dc93d8e0ada59e8bba7ad08a5783f58f6ad68.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i4e4fe41f8ec4b0f84b1c21baa5931909bd57fcb9768d3be13b17583642944ef8.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i4e4fe41f8ec4b0f84b1c21baa5931909bd57fcb9768d3be13b17583642944ef8.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i034f075a9742087cc62fa2fa3c2d823b0c34579b549387b0217d091d4a084aed.SingleValueExtendedPropertiesRequestBuilder) {
    return i034f075a9742087cc62fa2fa3c2d823b0c34579b549387b0217d091d4a084aed.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*if7192229b65a228ea27ede8842c1b7d3638ac8120f07a51c43c8174a1920dd14.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return if7192229b65a228ea27ede8842c1b7d3638ac8120f07a51c43c8174a1920dd14.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ic66f697cfef967cd4effb4c9d92b29e287a4def44f9b2aea82d334789318299f.SnoozeReminderRequestBuilder) {
    return ic66f697cfef967cd4effb4c9d92b29e287a4def44f9b2aea82d334789318299f.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ie4a846e389b081fa82f31b8581ae8dfed0401f4a81a686e4cc965b82d9606851.TentativelyAcceptRequestBuilder) {
    return ie4a846e389b081fa82f31b8581ae8dfed0401f4a81a686e4cc965b82d9606851.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
