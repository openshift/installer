package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i57af6d665855a0506938ef0e172bba0a14925f6b6ca627e2a7f4b6b2b38fd3bc "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/forward"
    i5b955e22c663573118720c3bc93ddf49ea1d54fc25dbd270017f41a676a78121 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/calendar"
    i6886d0c1332bc9d55516ea947d572efa5b319910b95fbb0820df71ec1f978f5a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances"
    i860a97d3a0a455053fe8451f67950c18e9a2e873fe28c966f761d1441a1b3f8f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/cancel"
    i8c7ff56f4d95c55b53a19a6b026409669658ef2e87e0166571395d96a4f5dc9b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/accept"
    i970025de78959a55575efb35fa3a413a32ebb5ab4dd5a708532d02f291175272 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/extensions"
    ia5e7630df574c5850c03f09e0d5676dde3236e1bca7c9041678b2fa25ee07e57 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/decline"
    ia6e2cb4fce4bf7b4978f9f272441f1b3a147cf7e6cb2a0de353a804070b27112 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/tentativelyaccept"
    ib354b8e7fb23c71741722b1fefce34b3b747a56eef37f83a062018a757f1aa72 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/attachments"
    iba41db6ae0a52d498c587df3e5bfcd8faf2437bcf7d4e46fd758bb57a9b0793e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/singlevalueextendedproperties"
    ic243b07f463bb3f30c6cb6f2fd2d449e4e4cec617ebca670effd9fc299701f4c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/multivalueextendedproperties"
    id98aa3c369bca6039a3982c6e8eaad54a60e63bab85c3ca6d5ddcb7fef06e0ab "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/dismissreminder"
    if141e32134f129fe69c623ef5329c6c06edda15c3911a704e36f0879d7277b76 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/snoozereminder"
    i4d40e556c4fcba97401013dc12446fe5534169b8814c68f0f34bc69689136da0 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/singlevalueextendedproperties/item"
    i931f24b3eb23f7e9c765dc76b0676613a564dedcec2ade920230227f4d4da606 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/multivalueextendedproperties/item"
    ibfd185eccf652dada583d3ca36c3d4bfc599e1bd77b6d8088608b2683fa14880 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/extensions/item"
    ie19dffbaf9991ce0123621be61d8c60393863cee5c24c78b5a5311caf84c056f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/attachments/item"
    ie90f5d53aab3ae5d6e53998a425e8dac6aec14b0af9e64448368e99bd0616cb5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item"
)

// EventItemRequestBuilder provides operations to manage the events property of the microsoft.graph.calendar entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EventItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EventItemRequestBuilderGetQueryParameters the events in the calendar. Navigation property. Read-only.
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
// EventItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EventItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Accept provides operations to call the accept method.
func (m *EventItemRequestBuilder) Accept()(*i8c7ff56f4d95c55b53a19a6b026409669658ef2e87e0166571395d96a4f5dc9b.AcceptRequestBuilder) {
    return i8c7ff56f4d95c55b53a19a6b026409669658ef2e87e0166571395d96a4f5dc9b.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ib354b8e7fb23c71741722b1fefce34b3b747a56eef37f83a062018a757f1aa72.AttachmentsRequestBuilder) {
    return ib354b8e7fb23c71741722b1fefce34b3b747a56eef37f83a062018a757f1aa72.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ie19dffbaf9991ce0123621be61d8c60393863cee5c24c78b5a5311caf84c056f.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ie19dffbaf9991ce0123621be61d8c60393863cee5c24c78b5a5311caf84c056f.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i5b955e22c663573118720c3bc93ddf49ea1d54fc25dbd270017f41a676a78121.CalendarRequestBuilder) {
    return i5b955e22c663573118720c3bc93ddf49ea1d54fc25dbd270017f41a676a78121.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i860a97d3a0a455053fe8451f67950c18e9a2e873fe28c966f761d1441a1b3f8f.CancelRequestBuilder) {
    return i860a97d3a0a455053fe8451f67950c18e9a2e873fe28c966f761d1441a1b3f8f.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/calendar/events/{event%2Did}{?%24select}";
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
// CreateDeleteRequestInformation delete navigation property events for groups
func (m *EventItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EventItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreateGetRequestInformation the events in the calendar. Navigation property. Read-only.
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
// CreatePatchRequestInformation update the navigation property events in groups
func (m *EventItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Eventable, requestConfiguration *EventItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Decline provides operations to call the decline method.
func (m *EventItemRequestBuilder) Decline()(*ia5e7630df574c5850c03f09e0d5676dde3236e1bca7c9041678b2fa25ee07e57.DeclineRequestBuilder) {
    return ia5e7630df574c5850c03f09e0d5676dde3236e1bca7c9041678b2fa25ee07e57.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property events for groups
func (m *EventItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EventItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.CreateDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContentAsync(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*id98aa3c369bca6039a3982c6e8eaad54a60e63bab85c3ca6d5ddcb7fef06e0ab.DismissReminderRequestBuilder) {
    return id98aa3c369bca6039a3982c6e8eaad54a60e63bab85c3ca6d5ddcb7fef06e0ab.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i970025de78959a55575efb35fa3a413a32ebb5ab4dd5a708532d02f291175272.ExtensionsRequestBuilder) {
    return i970025de78959a55575efb35fa3a413a32ebb5ab4dd5a708532d02f291175272.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ibfd185eccf652dada583d3ca36c3d4bfc599e1bd77b6d8088608b2683fa14880.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ibfd185eccf652dada583d3ca36c3d4bfc599e1bd77b6d8088608b2683fa14880.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i57af6d665855a0506938ef0e172bba0a14925f6b6ca627e2a7f4b6b2b38fd3bc.ForwardRequestBuilder) {
    return i57af6d665855a0506938ef0e172bba0a14925f6b6ca627e2a7f4b6b2b38fd3bc.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the events in the calendar. Navigation property. Read-only.
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
func (m *EventItemRequestBuilder) Instances()(*i6886d0c1332bc9d55516ea947d572efa5b319910b95fbb0820df71ec1f978f5a.InstancesRequestBuilder) {
    return i6886d0c1332bc9d55516ea947d572efa5b319910b95fbb0820df71ec1f978f5a.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*ie90f5d53aab3ae5d6e53998a425e8dac6aec14b0af9e64448368e99bd0616cb5.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return ie90f5d53aab3ae5d6e53998a425e8dac6aec14b0af9e64448368e99bd0616cb5.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ic243b07f463bb3f30c6cb6f2fd2d449e4e4cec617ebca670effd9fc299701f4c.MultiValueExtendedPropertiesRequestBuilder) {
    return ic243b07f463bb3f30c6cb6f2fd2d449e4e4cec617ebca670effd9fc299701f4c.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i931f24b3eb23f7e9c765dc76b0676613a564dedcec2ade920230227f4d4da606.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i931f24b3eb23f7e9c765dc76b0676613a564dedcec2ade920230227f4d4da606.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property events in groups
func (m *EventItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Eventable, requestConfiguration *EventItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Eventable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
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
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*iba41db6ae0a52d498c587df3e5bfcd8faf2437bcf7d4e46fd758bb57a9b0793e.SingleValueExtendedPropertiesRequestBuilder) {
    return iba41db6ae0a52d498c587df3e5bfcd8faf2437bcf7d4e46fd758bb57a9b0793e.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i4d40e556c4fcba97401013dc12446fe5534169b8814c68f0f34bc69689136da0.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i4d40e556c4fcba97401013dc12446fe5534169b8814c68f0f34bc69689136da0.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*if141e32134f129fe69c623ef5329c6c06edda15c3911a704e36f0879d7277b76.SnoozeReminderRequestBuilder) {
    return if141e32134f129fe69c623ef5329c6c06edda15c3911a704e36f0879d7277b76.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ia6e2cb4fce4bf7b4978f9f272441f1b3a147cf7e6cb2a0de353a804070b27112.TentativelyAcceptRequestBuilder) {
    return ia6e2cb4fce4bf7b4978f9f272441f1b3a147cf7e6cb2a0de353a804070b27112.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
