package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i068c994d6e3103d7a8c451470fd8946c9a526bf7defe804a3414bec3f058c7c0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/attachments"
    i178988a2b9571cc52545bd2a40fe9246c6b63e31f9a8adbe3f7c47021e3ebfd5 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/snoozereminder"
    i39f7b9eb2753296d89e1366bb2fa48858540dd472648ec034a0c3a194183a61c "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/cancel"
    i65f4d175fc410ed923f2cd477c6e12f9b6cfed24669da0d6147a1eea2b7353b0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/forward"
    i6d2604c2a16b79e867c6950984ac155a11a951d6935ff922fb3d80299b8aae9c "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties"
    i732564651bd295ace901245d83f7cee6eecbd49650f303b4d6f278dce8d4196c "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/decline"
    i75b85d086dff2c57d9c90524d349608e34c224e7740e4a9a09c8e3892e4d650b "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/calendar"
    ia3e61a8f98e0b4ebf997055c903954f21cee726b10079e1c431a38baa7bb1de9 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/accept"
    iaa82cd4a7429733dee7556dd9022c3f856a8e77926ee4dba3e89fb673d6c269b "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/tentativelyaccept"
    ic38a64ab747dedc47441f8456aa350055b7b0ce2810f51c3c45312d8542baab0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/dismissreminder"
    ic7e25b5dfc5b0afcb9c1a08774f134a922789857f8a916436156eac48cd979cc "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties"
    icf3fe84f8e6c9727dc9e6c970c133185bbdca23532e728e26e585eacbd8d2d45 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/extensions"
    i34702bf27bd058d7b5e871f07e6514ca3dd947d728fc016f03eb36f170b53074 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/attachments/item"
    iaaa7634fa7b22059d9b7249765ef9f412f206b4d530675317f19af70e389cf34 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties/item"
    iac1bea5f44cb60b84d138bb975b473a16544fe8096eae5c117e8a7ac282abf8f "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/extensions/item"
    iea9fc93fa7aea2d4879fe2764aa490f2a684a91d6a20a48d909c8cfdcfe9d479 "github.com/microsoftgraph/msgraph-sdk-go/me/calendargroups/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*ia3e61a8f98e0b4ebf997055c903954f21cee726b10079e1c431a38baa7bb1de9.AcceptRequestBuilder) {
    return ia3e61a8f98e0b4ebf997055c903954f21cee726b10079e1c431a38baa7bb1de9.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i068c994d6e3103d7a8c451470fd8946c9a526bf7defe804a3414bec3f058c7c0.AttachmentsRequestBuilder) {
    return i068c994d6e3103d7a8c451470fd8946c9a526bf7defe804a3414bec3f058c7c0.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i34702bf27bd058d7b5e871f07e6514ca3dd947d728fc016f03eb36f170b53074.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i34702bf27bd058d7b5e871f07e6514ca3dd947d728fc016f03eb36f170b53074.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i75b85d086dff2c57d9c90524d349608e34c224e7740e4a9a09c8e3892e4d650b.CalendarRequestBuilder) {
    return i75b85d086dff2c57d9c90524d349608e34c224e7740e4a9a09c8e3892e4d650b.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i39f7b9eb2753296d89e1366bb2fa48858540dd472648ec034a0c3a194183a61c.CancelRequestBuilder) {
    return i39f7b9eb2753296d89e1366bb2fa48858540dd472648ec034a0c3a194183a61c.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendarGroups/{calendarGroup%2Did}/calendars/{calendar%2Did}/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i732564651bd295ace901245d83f7cee6eecbd49650f303b4d6f278dce8d4196c.DeclineRequestBuilder) {
    return i732564651bd295ace901245d83f7cee6eecbd49650f303b4d6f278dce8d4196c.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*ic38a64ab747dedc47441f8456aa350055b7b0ce2810f51c3c45312d8542baab0.DismissReminderRequestBuilder) {
    return ic38a64ab747dedc47441f8456aa350055b7b0ce2810f51c3c45312d8542baab0.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*icf3fe84f8e6c9727dc9e6c970c133185bbdca23532e728e26e585eacbd8d2d45.ExtensionsRequestBuilder) {
    return icf3fe84f8e6c9727dc9e6c970c133185bbdca23532e728e26e585eacbd8d2d45.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*iac1bea5f44cb60b84d138bb975b473a16544fe8096eae5c117e8a7ac282abf8f.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return iac1bea5f44cb60b84d138bb975b473a16544fe8096eae5c117e8a7ac282abf8f.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i65f4d175fc410ed923f2cd477c6e12f9b6cfed24669da0d6147a1eea2b7353b0.ForwardRequestBuilder) {
    return i65f4d175fc410ed923f2cd477c6e12f9b6cfed24669da0d6147a1eea2b7353b0.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ic7e25b5dfc5b0afcb9c1a08774f134a922789857f8a916436156eac48cd979cc.MultiValueExtendedPropertiesRequestBuilder) {
    return ic7e25b5dfc5b0afcb9c1a08774f134a922789857f8a916436156eac48cd979cc.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*iea9fc93fa7aea2d4879fe2764aa490f2a684a91d6a20a48d909c8cfdcfe9d479.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return iea9fc93fa7aea2d4879fe2764aa490f2a684a91d6a20a48d909c8cfdcfe9d479.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i6d2604c2a16b79e867c6950984ac155a11a951d6935ff922fb3d80299b8aae9c.SingleValueExtendedPropertiesRequestBuilder) {
    return i6d2604c2a16b79e867c6950984ac155a11a951d6935ff922fb3d80299b8aae9c.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*iaaa7634fa7b22059d9b7249765ef9f412f206b4d530675317f19af70e389cf34.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return iaaa7634fa7b22059d9b7249765ef9f412f206b4d530675317f19af70e389cf34.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i178988a2b9571cc52545bd2a40fe9246c6b63e31f9a8adbe3f7c47021e3ebfd5.SnoozeReminderRequestBuilder) {
    return i178988a2b9571cc52545bd2a40fe9246c6b63e31f9a8adbe3f7c47021e3ebfd5.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*iaa82cd4a7429733dee7556dd9022c3f856a8e77926ee4dba3e89fb673d6c269b.TentativelyAcceptRequestBuilder) {
    return iaa82cd4a7429733dee7556dd9022c3f856a8e77926ee4dba3e89fb673d6c269b.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
