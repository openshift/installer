package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i112b1ed50d98e06e2afc6c27777a4413e4aec03373c2c045d8fab1d827230449 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/dismissreminder"
    i1bcbde3648b969716e475c91cdb32d8cdf62c76e85838d834083a236f2d38333 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/tentativelyaccept"
    i3c35eb67b420d24b8b792bd486949e59c9d8b4b7fb3f16f018266a3064470c46 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/snoozereminder"
    i4429fd37a56c4f85306e5b908f5976bb1569387d88606fdbc734b5b5afff1a5e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/multivalueextendedproperties"
    i4b3ac3d77873723261c8c94dc247eb0c4e33b3adc4a4b2528b3a8fc513f95250 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/cancel"
    i651ebf3ec997ed2debe308c2e0ed3a208f869f2a9c82390f706b6456e906e255 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances"
    i719ba32a86c43f7802b551a6d9320eb0b3e90c285148770a88c044f42b79b024 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/attachments"
    i7ded4b752eb3655e749d50f9eafabc0cbea03fdebdfb187b26ddd1566a7c4236 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/accept"
    i965695020684e888087236e60e5c0f7971367bc349f4a0088bd1af20ca363f05 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/forward"
    ia1eb05bd3c4f8481e8d0ac2eb5ff47aa090fb98c35c567df6e13064f70d0b3d0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/calendar"
    ia39bb4166c4ef2deee03aec7af70c9763aea4f1fc2b921860014e829fd638086 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/singlevalueextendedproperties"
    iece4056a8d35809aba6e4ff7ae53718558f0501f3bf175b276feb26afb59a056 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/extensions"
    iffae5ee2a2d9a81008b555644a80fc76abfbbd265c837cecf1568216e5db2972 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/decline"
    i0df45b242d0b94600d9cbfc3f4e57b9acc807e751ed7030c7c5d4fb2846927eb "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/instances/item"
    i1af302b703d5ff4044ea6d579861399d521115e720901fd67e398c46ac2c4dde "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/singlevalueextendedproperties/item"
    i20495d21133b26c7848a42901db05d39645c739c0ef4dfc27137db115c7ff68a "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/multivalueextendedproperties/item"
    i536a7b87db2cf3b5a9fc821e79fa11d85cd220fa19be5513d62932141a446474 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/extensions/item"
    i55f8e56a76f6d157ef3775e097b2a4e244e70c88edea43c5a8ba28bf62638a17 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*i7ded4b752eb3655e749d50f9eafabc0cbea03fdebdfb187b26ddd1566a7c4236.AcceptRequestBuilder) {
    return i7ded4b752eb3655e749d50f9eafabc0cbea03fdebdfb187b26ddd1566a7c4236.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i719ba32a86c43f7802b551a6d9320eb0b3e90c285148770a88c044f42b79b024.AttachmentsRequestBuilder) {
    return i719ba32a86c43f7802b551a6d9320eb0b3e90c285148770a88c044f42b79b024.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i55f8e56a76f6d157ef3775e097b2a4e244e70c88edea43c5a8ba28bf62638a17.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i55f8e56a76f6d157ef3775e097b2a4e244e70c88edea43c5a8ba28bf62638a17.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ia1eb05bd3c4f8481e8d0ac2eb5ff47aa090fb98c35c567df6e13064f70d0b3d0.CalendarRequestBuilder) {
    return ia1eb05bd3c4f8481e8d0ac2eb5ff47aa090fb98c35c567df6e13064f70d0b3d0.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i4b3ac3d77873723261c8c94dc247eb0c4e33b3adc4a4b2528b3a8fc513f95250.CancelRequestBuilder) {
    return i4b3ac3d77873723261c8c94dc247eb0c4e33b3adc4a4b2528b3a8fc513f95250.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendar/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*iffae5ee2a2d9a81008b555644a80fc76abfbbd265c837cecf1568216e5db2972.DeclineRequestBuilder) {
    return iffae5ee2a2d9a81008b555644a80fc76abfbbd265c837cecf1568216e5db2972.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i112b1ed50d98e06e2afc6c27777a4413e4aec03373c2c045d8fab1d827230449.DismissReminderRequestBuilder) {
    return i112b1ed50d98e06e2afc6c27777a4413e4aec03373c2c045d8fab1d827230449.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*iece4056a8d35809aba6e4ff7ae53718558f0501f3bf175b276feb26afb59a056.ExtensionsRequestBuilder) {
    return iece4056a8d35809aba6e4ff7ae53718558f0501f3bf175b276feb26afb59a056.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i536a7b87db2cf3b5a9fc821e79fa11d85cd220fa19be5513d62932141a446474.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i536a7b87db2cf3b5a9fc821e79fa11d85cd220fa19be5513d62932141a446474.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i965695020684e888087236e60e5c0f7971367bc349f4a0088bd1af20ca363f05.ForwardRequestBuilder) {
    return i965695020684e888087236e60e5c0f7971367bc349f4a0088bd1af20ca363f05.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*i651ebf3ec997ed2debe308c2e0ed3a208f869f2a9c82390f706b6456e906e255.InstancesRequestBuilder) {
    return i651ebf3ec997ed2debe308c2e0ed3a208f869f2a9c82390f706b6456e906e255.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i0df45b242d0b94600d9cbfc3f4e57b9acc807e751ed7030c7c5d4fb2846927eb.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i0df45b242d0b94600d9cbfc3f4e57b9acc807e751ed7030c7c5d4fb2846927eb.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i4429fd37a56c4f85306e5b908f5976bb1569387d88606fdbc734b5b5afff1a5e.MultiValueExtendedPropertiesRequestBuilder) {
    return i4429fd37a56c4f85306e5b908f5976bb1569387d88606fdbc734b5b5afff1a5e.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i20495d21133b26c7848a42901db05d39645c739c0ef4dfc27137db115c7ff68a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i20495d21133b26c7848a42901db05d39645c739c0ef4dfc27137db115c7ff68a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*ia39bb4166c4ef2deee03aec7af70c9763aea4f1fc2b921860014e829fd638086.SingleValueExtendedPropertiesRequestBuilder) {
    return ia39bb4166c4ef2deee03aec7af70c9763aea4f1fc2b921860014e829fd638086.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i1af302b703d5ff4044ea6d579861399d521115e720901fd67e398c46ac2c4dde.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i1af302b703d5ff4044ea6d579861399d521115e720901fd67e398c46ac2c4dde.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i3c35eb67b420d24b8b792bd486949e59c9d8b4b7fb3f16f018266a3064470c46.SnoozeReminderRequestBuilder) {
    return i3c35eb67b420d24b8b792bd486949e59c9d8b4b7fb3f16f018266a3064470c46.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i1bcbde3648b969716e475c91cdb32d8cdf62c76e85838d834083a236f2d38333.TentativelyAcceptRequestBuilder) {
    return i1bcbde3648b969716e475c91cdb32d8cdf62c76e85838d834083a236f2d38333.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
