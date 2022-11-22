package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i373f90719cb3f0b97f21dfb295efc403c934571a604d4d6a7bd62ca30a44e015 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/decline"
    i54b4dc0b393454aea4bd0dd164a30a992216c43946031f9b90d97c00cd50f83b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/tentativelyaccept"
    i5f2e6f653d07ecb7b2247cf520e1af568cc5362ef0d2be82946ef4ad4acef11a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/calendar"
    i6405624f66ea71046ed4d92712a7fc0dc1c0773e306cf79641eff6a6341ddc4e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/dismissreminder"
    i7b4042292f5e408f0b303dd1f1bb04dd7eca9b9e8d2c9ad0f0dad072bd5bdf30 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/extensions"
    i860f704332d66191ef5006487d66838635ed277d6784f0cc821b0776c2480719 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/multivalueextendedproperties"
    i87f2def9d9f68286d61b7299f82fb7181145d355681c461945f2774639a001bb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/snoozereminder"
    i8c1ca72db80d40922fca8d9b05b8467812088b7c6db75f15d0bbcc2f134dfb3b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/cancel"
    iababe505c0a1c9ca0dda078443102a263523994dd5a228924133eea244c51070 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/accept"
    ibc1a8cb1d61cc40e3c09eab86a3716f5f5b1f6121f8f823947a66d8425ad65ff "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/forward"
    ic97821bf0d97a40cc49a9727bb9f16130d876883d0c0ecd4c81ff5dd54503f38 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/singlevalueextendedproperties"
    id13a32ef047da8a7524f21260998047193088e9e89e04696b106f076b087d8d1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/attachments"
    i3cec616c4914d987dbac48dbdab1ab50d8867952222702d410096231050b5cc7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/multivalueextendedproperties/item"
    i65d593b6ee680564ce3ed06f21c8a748e353302c7d22c7109dd7412c1e64fe50 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/singlevalueextendedproperties/item"
    ia794c2f35f8a3b32feef2f9b52d199f6dbf453fd3857ca899a220880b77b7095 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/extensions/item"
    idfb7558c4133b989c63f0f3673f9d1ae416736570bcce60a3290124917f9f625 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*iababe505c0a1c9ca0dda078443102a263523994dd5a228924133eea244c51070.AcceptRequestBuilder) {
    return iababe505c0a1c9ca0dda078443102a263523994dd5a228924133eea244c51070.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*id13a32ef047da8a7524f21260998047193088e9e89e04696b106f076b087d8d1.AttachmentsRequestBuilder) {
    return id13a32ef047da8a7524f21260998047193088e9e89e04696b106f076b087d8d1.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*idfb7558c4133b989c63f0f3673f9d1ae416736570bcce60a3290124917f9f625.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return idfb7558c4133b989c63f0f3673f9d1ae416736570bcce60a3290124917f9f625.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i5f2e6f653d07ecb7b2247cf520e1af568cc5362ef0d2be82946ef4ad4acef11a.CalendarRequestBuilder) {
    return i5f2e6f653d07ecb7b2247cf520e1af568cc5362ef0d2be82946ef4ad4acef11a.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i8c1ca72db80d40922fca8d9b05b8467812088b7c6db75f15d0bbcc2f134dfb3b.CancelRequestBuilder) {
    return i8c1ca72db80d40922fca8d9b05b8467812088b7c6db75f15d0bbcc2f134dfb3b.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i373f90719cb3f0b97f21dfb295efc403c934571a604d4d6a7bd62ca30a44e015.DeclineRequestBuilder) {
    return i373f90719cb3f0b97f21dfb295efc403c934571a604d4d6a7bd62ca30a44e015.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i6405624f66ea71046ed4d92712a7fc0dc1c0773e306cf79641eff6a6341ddc4e.DismissReminderRequestBuilder) {
    return i6405624f66ea71046ed4d92712a7fc0dc1c0773e306cf79641eff6a6341ddc4e.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i7b4042292f5e408f0b303dd1f1bb04dd7eca9b9e8d2c9ad0f0dad072bd5bdf30.ExtensionsRequestBuilder) {
    return i7b4042292f5e408f0b303dd1f1bb04dd7eca9b9e8d2c9ad0f0dad072bd5bdf30.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ia794c2f35f8a3b32feef2f9b52d199f6dbf453fd3857ca899a220880b77b7095.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ia794c2f35f8a3b32feef2f9b52d199f6dbf453fd3857ca899a220880b77b7095.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ibc1a8cb1d61cc40e3c09eab86a3716f5f5b1f6121f8f823947a66d8425ad65ff.ForwardRequestBuilder) {
    return ibc1a8cb1d61cc40e3c09eab86a3716f5f5b1f6121f8f823947a66d8425ad65ff.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i860f704332d66191ef5006487d66838635ed277d6784f0cc821b0776c2480719.MultiValueExtendedPropertiesRequestBuilder) {
    return i860f704332d66191ef5006487d66838635ed277d6784f0cc821b0776c2480719.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i3cec616c4914d987dbac48dbdab1ab50d8867952222702d410096231050b5cc7.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i3cec616c4914d987dbac48dbdab1ab50d8867952222702d410096231050b5cc7.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*ic97821bf0d97a40cc49a9727bb9f16130d876883d0c0ecd4c81ff5dd54503f38.SingleValueExtendedPropertiesRequestBuilder) {
    return ic97821bf0d97a40cc49a9727bb9f16130d876883d0c0ecd4c81ff5dd54503f38.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i65d593b6ee680564ce3ed06f21c8a748e353302c7d22c7109dd7412c1e64fe50.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i65d593b6ee680564ce3ed06f21c8a748e353302c7d22c7109dd7412c1e64fe50.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i87f2def9d9f68286d61b7299f82fb7181145d355681c461945f2774639a001bb.SnoozeReminderRequestBuilder) {
    return i87f2def9d9f68286d61b7299f82fb7181145d355681c461945f2774639a001bb.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i54b4dc0b393454aea4bd0dd164a30a992216c43946031f9b90d97c00cd50f83b.TentativelyAcceptRequestBuilder) {
    return i54b4dc0b393454aea4bd0dd164a30a992216c43946031f9b90d97c00cd50f83b.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
