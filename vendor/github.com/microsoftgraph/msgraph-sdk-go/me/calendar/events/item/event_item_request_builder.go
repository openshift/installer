package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i083a9367fb14f0da75caea7bc857922fe10d91693dde0f11517ab6420e557018 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/attachments"
    i0f781326a67fbc9776be4d1ccbf43639b9bf8077faa925a7b0f90f0e03141c07 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/snoozereminder"
    i10a043ae8449d6b8061165bb67946e604bc2614783ade4a4557a2bb8d2a5a2b0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/decline"
    i339d2d1a30983a01dd7de00efdd9e7d51cf0424dd9c4c96d335f9901fcff9763 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances"
    i49ad5e4dbfe1e924deb723580385a06e87f447b9ed73e648b64a934117f1c784 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/forward"
    i62d05692265cf61165d6bb2a1bb3a36d03ea2604586e3ebd550fae433333488e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/singlevalueextendedproperties"
    i69bb0e8065fcc365cd51c4f9cb100c9289f8c378fe585dd0126d7eaee8953622 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/calendar"
    i87334fbab17098d26d733b7c6d4d979d97f28091d5d1b3c41c2e4adcd9d2f2a3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/accept"
    id1d8fbd767f32fa77bb0a1e1f8b79c51f1a58b7dae1ccebcd3a3001af55d4f44 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/extensions"
    ie04ca0f3f9823daf845f8b3b165c34feaffc35ed4189263a607834108a606dbe "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/cancel"
    ie6641980e1225311bdf1c499d054bc5d07aa799f49902ef90e71b01800c2a555 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/dismissreminder"
    iee14c2abd2614ebab15ceaf1135456d72c4e1cd953dd5d7b3347f819c29c6002 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/tentativelyaccept"
    if78feb360a058ff87bf5a02baba6f7d4fb28a9e5f93fc081392b457ad6fae3eb "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/multivalueextendedproperties"
    i3f305339dd659581dddbb6081f2d1981394ba5f409b3fff8d5affc26648441c4 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/singlevalueextendedproperties/item"
    i658437a8657f51783d49c87367e2ae5b495c006ac48f6a7010d75098566455a0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/multivalueextendedproperties/item"
    i868bcbc69b2c3a35a368f852b47fffbb78632a8f15944de468f0e130a0242300 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/attachments/item"
    i955a97236a9f337897b15ed49defa3aa5c5d7910ae15f5b5efe5fbae5e75cb05 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/extensions/item"
    ic8bb7a27f91fd567d4417f19bf0aacb6595974ef7d4a830d5f1842796aab82f7 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item"
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
func (m *EventItemRequestBuilder) Accept()(*i87334fbab17098d26d733b7c6d4d979d97f28091d5d1b3c41c2e4adcd9d2f2a3.AcceptRequestBuilder) {
    return i87334fbab17098d26d733b7c6d4d979d97f28091d5d1b3c41c2e4adcd9d2f2a3.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i083a9367fb14f0da75caea7bc857922fe10d91693dde0f11517ab6420e557018.AttachmentsRequestBuilder) {
    return i083a9367fb14f0da75caea7bc857922fe10d91693dde0f11517ab6420e557018.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i868bcbc69b2c3a35a368f852b47fffbb78632a8f15944de468f0e130a0242300.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i868bcbc69b2c3a35a368f852b47fffbb78632a8f15944de468f0e130a0242300.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i69bb0e8065fcc365cd51c4f9cb100c9289f8c378fe585dd0126d7eaee8953622.CalendarRequestBuilder) {
    return i69bb0e8065fcc365cd51c4f9cb100c9289f8c378fe585dd0126d7eaee8953622.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ie04ca0f3f9823daf845f8b3b165c34feaffc35ed4189263a607834108a606dbe.CancelRequestBuilder) {
    return ie04ca0f3f9823daf845f8b3b165c34feaffc35ed4189263a607834108a606dbe.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendar/events/{event%2Did}{?%24select}";
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
// CreateDeleteRequestInformation delete navigation property events for me
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
// CreatePatchRequestInformation update the navigation property events in me
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
func (m *EventItemRequestBuilder) Decline()(*i10a043ae8449d6b8061165bb67946e604bc2614783ade4a4557a2bb8d2a5a2b0.DeclineRequestBuilder) {
    return i10a043ae8449d6b8061165bb67946e604bc2614783ade4a4557a2bb8d2a5a2b0.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property events for me
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
func (m *EventItemRequestBuilder) DismissReminder()(*ie6641980e1225311bdf1c499d054bc5d07aa799f49902ef90e71b01800c2a555.DismissReminderRequestBuilder) {
    return ie6641980e1225311bdf1c499d054bc5d07aa799f49902ef90e71b01800c2a555.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*id1d8fbd767f32fa77bb0a1e1f8b79c51f1a58b7dae1ccebcd3a3001af55d4f44.ExtensionsRequestBuilder) {
    return id1d8fbd767f32fa77bb0a1e1f8b79c51f1a58b7dae1ccebcd3a3001af55d4f44.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i955a97236a9f337897b15ed49defa3aa5c5d7910ae15f5b5efe5fbae5e75cb05.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i955a97236a9f337897b15ed49defa3aa5c5d7910ae15f5b5efe5fbae5e75cb05.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i49ad5e4dbfe1e924deb723580385a06e87f447b9ed73e648b64a934117f1c784.ForwardRequestBuilder) {
    return i49ad5e4dbfe1e924deb723580385a06e87f447b9ed73e648b64a934117f1c784.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*i339d2d1a30983a01dd7de00efdd9e7d51cf0424dd9c4c96d335f9901fcff9763.InstancesRequestBuilder) {
    return i339d2d1a30983a01dd7de00efdd9e7d51cf0424dd9c4c96d335f9901fcff9763.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*ic8bb7a27f91fd567d4417f19bf0aacb6595974ef7d4a830d5f1842796aab82f7.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return ic8bb7a27f91fd567d4417f19bf0aacb6595974ef7d4a830d5f1842796aab82f7.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*if78feb360a058ff87bf5a02baba6f7d4fb28a9e5f93fc081392b457ad6fae3eb.MultiValueExtendedPropertiesRequestBuilder) {
    return if78feb360a058ff87bf5a02baba6f7d4fb28a9e5f93fc081392b457ad6fae3eb.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i658437a8657f51783d49c87367e2ae5b495c006ac48f6a7010d75098566455a0.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i658437a8657f51783d49c87367e2ae5b495c006ac48f6a7010d75098566455a0.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property events in me
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
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i62d05692265cf61165d6bb2a1bb3a36d03ea2604586e3ebd550fae433333488e.SingleValueExtendedPropertiesRequestBuilder) {
    return i62d05692265cf61165d6bb2a1bb3a36d03ea2604586e3ebd550fae433333488e.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i3f305339dd659581dddbb6081f2d1981394ba5f409b3fff8d5affc26648441c4.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i3f305339dd659581dddbb6081f2d1981394ba5f409b3fff8d5affc26648441c4.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i0f781326a67fbc9776be4d1ccbf43639b9bf8077faa925a7b0f90f0e03141c07.SnoozeReminderRequestBuilder) {
    return i0f781326a67fbc9776be4d1ccbf43639b9bf8077faa925a7b0f90f0e03141c07.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*iee14c2abd2614ebab15ceaf1135456d72c4e1cd953dd5d7b3347f819c29c6002.TentativelyAcceptRequestBuilder) {
    return iee14c2abd2614ebab15ceaf1135456d72c4e1cd953dd5d7b3347f819c29c6002.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
