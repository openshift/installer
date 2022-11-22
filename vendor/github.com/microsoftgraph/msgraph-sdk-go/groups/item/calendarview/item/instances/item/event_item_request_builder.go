package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i3cc34074d4799ae365015afe8a944936199187ea89a89982620bc4d5270ee20a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/accept"
    i4acb1eeebf60790e92c699b157738449aea78a952be2463ce026b8040d7e5789 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/multivalueextendedproperties"
    i578c3c432a4af3223ca10fc3355e1e507450b8186f1aafc298ff748751814861 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/attachments"
    i6974512b0df45bd6478c5bd76ba677d7332b9275200ed141bb1c82d44ac9d26f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/singlevalueextendedproperties"
    i756e61eebdd85c15621a1de3592845102e21d2510bcf4f47c8b6f633b3d7cc8e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/calendar"
    i7758ded31927896a0324ac6ba89d90109741d24d1a194a275042f6f228cdc4be "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/dismissreminder"
    i849f6cd4968ad54bff7c4acc3d2c814119ed8035747edd7da13892771dd1f6cf "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/forward"
    i8c1de96d676c31bb3baa4d5578220df8cc694fcb1499fa3ced34cc06d55588c9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/snoozereminder"
    ib55ccb1e43bdb65ed253495385ecec98286a535ace3067311f181b27a9a373d7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/extensions"
    ib5a7878855ffd0e4fbdf84af0ae8b0f0a7974f30f055f4ced1d57fdb861e0fa5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/tentativelyaccept"
    ib9699fbc3be709b24c33827610ccf455d295816e912a53e123836e2c220e6877 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/decline"
    ic8a6a36a3fed4f835a1259b874a6c34cb4c966ccbabfa5b23cfa9244dbbb2151 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/cancel"
    i006d6330e3c4b5ce635e3dd7b2de257d6731702cbef84980ddf3887ec16eb40e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/singlevalueextendedproperties/item"
    i538b156a8aee4235e0fe1fbffcecd293c195c55b10dc0e4aaa92e8a97447ea26 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/attachments/item"
    iabde4e47b53f2370aeb27d4c48fa6dd8b902a360e08ca7a4f145a8b7d3e0fd2e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/extensions/item"
    ifd4e95066ee67088ba4bf159e28dd71f75754bc7edc71bf65b717d7f84239541 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i3cc34074d4799ae365015afe8a944936199187ea89a89982620bc4d5270ee20a.AcceptRequestBuilder) {
    return i3cc34074d4799ae365015afe8a944936199187ea89a89982620bc4d5270ee20a.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i578c3c432a4af3223ca10fc3355e1e507450b8186f1aafc298ff748751814861.AttachmentsRequestBuilder) {
    return i578c3c432a4af3223ca10fc3355e1e507450b8186f1aafc298ff748751814861.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i538b156a8aee4235e0fe1fbffcecd293c195c55b10dc0e4aaa92e8a97447ea26.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i538b156a8aee4235e0fe1fbffcecd293c195c55b10dc0e4aaa92e8a97447ea26.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i756e61eebdd85c15621a1de3592845102e21d2510bcf4f47c8b6f633b3d7cc8e.CalendarRequestBuilder) {
    return i756e61eebdd85c15621a1de3592845102e21d2510bcf4f47c8b6f633b3d7cc8e.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ic8a6a36a3fed4f835a1259b874a6c34cb4c966ccbabfa5b23cfa9244dbbb2151.CancelRequestBuilder) {
    return ic8a6a36a3fed4f835a1259b874a6c34cb4c966ccbabfa5b23cfa9244dbbb2151.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*ib9699fbc3be709b24c33827610ccf455d295816e912a53e123836e2c220e6877.DeclineRequestBuilder) {
    return ib9699fbc3be709b24c33827610ccf455d295816e912a53e123836e2c220e6877.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i7758ded31927896a0324ac6ba89d90109741d24d1a194a275042f6f228cdc4be.DismissReminderRequestBuilder) {
    return i7758ded31927896a0324ac6ba89d90109741d24d1a194a275042f6f228cdc4be.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*ib55ccb1e43bdb65ed253495385ecec98286a535ace3067311f181b27a9a373d7.ExtensionsRequestBuilder) {
    return ib55ccb1e43bdb65ed253495385ecec98286a535ace3067311f181b27a9a373d7.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*iabde4e47b53f2370aeb27d4c48fa6dd8b902a360e08ca7a4f145a8b7d3e0fd2e.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return iabde4e47b53f2370aeb27d4c48fa6dd8b902a360e08ca7a4f145a8b7d3e0fd2e.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i849f6cd4968ad54bff7c4acc3d2c814119ed8035747edd7da13892771dd1f6cf.ForwardRequestBuilder) {
    return i849f6cd4968ad54bff7c4acc3d2c814119ed8035747edd7da13892771dd1f6cf.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i4acb1eeebf60790e92c699b157738449aea78a952be2463ce026b8040d7e5789.MultiValueExtendedPropertiesRequestBuilder) {
    return i4acb1eeebf60790e92c699b157738449aea78a952be2463ce026b8040d7e5789.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ifd4e95066ee67088ba4bf159e28dd71f75754bc7edc71bf65b717d7f84239541.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ifd4e95066ee67088ba4bf159e28dd71f75754bc7edc71bf65b717d7f84239541.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i6974512b0df45bd6478c5bd76ba677d7332b9275200ed141bb1c82d44ac9d26f.SingleValueExtendedPropertiesRequestBuilder) {
    return i6974512b0df45bd6478c5bd76ba677d7332b9275200ed141bb1c82d44ac9d26f.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i006d6330e3c4b5ce635e3dd7b2de257d6731702cbef84980ddf3887ec16eb40e.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i006d6330e3c4b5ce635e3dd7b2de257d6731702cbef84980ddf3887ec16eb40e.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i8c1de96d676c31bb3baa4d5578220df8cc694fcb1499fa3ced34cc06d55588c9.SnoozeReminderRequestBuilder) {
    return i8c1de96d676c31bb3baa4d5578220df8cc694fcb1499fa3ced34cc06d55588c9.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ib5a7878855ffd0e4fbdf84af0ae8b0f0a7974f30f055f4ced1d57fdb861e0fa5.TentativelyAcceptRequestBuilder) {
    return ib5a7878855ffd0e4fbdf84af0ae8b0f0a7974f30f055f4ced1d57fdb861e0fa5.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
