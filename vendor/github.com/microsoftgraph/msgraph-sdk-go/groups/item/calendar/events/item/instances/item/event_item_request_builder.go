package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0323a63dffab877f26880ff9d898dd683b137f110e681d7d437282a535f11c87 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/multivalueextendedproperties"
    i053391bc6adcb52e7a06235af432dde4b91f64c74f6aad7c530db1a2ac8d9fd5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/tentativelyaccept"
    i1eb540a156784a1c1466788d5998a6259a5bf9f8a6ad1d2f5c5350dd62845d31 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/dismissreminder"
    i374a3886dd17efd341a3266630c9c0193af40ed92cb05ffed90478816503f921 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/extensions"
    i3e367fd62461b1e59d434b1a9a1979a1913780b0085f47513db6136e5be05f43 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/accept"
    i7723196f3625220ead5e66ac8e7564fd63f8a974c7523942da7e11c443dce12c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/forward"
    i838bdc040b67a67f42d3a319019f7d751e6f2f638f8176040470892eda9fbefb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/snoozereminder"
    i8fad80f883c81dbb1610d3530af0b03447f5b648bdebd725cd0179a10144e8bf "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/calendar"
    icabac1b170b5a175bf40ceb126e85ca0f20d6ee5bb8f8ab53809cd2e31a78331 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/attachments"
    idebd0716278787fa5fab28a133715f58f2552b8fe8e2f7a9071417bf972eadeb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/singlevalueextendedproperties"
    if0a1b5804d4380cb388c51cff9421ac2b14ff858b9ac7ce1dd5e7b91d97b1860 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/decline"
    if13bf5db4c861d2c48743a82cb0fd3db30f4e2e2265b5c93b6a107b69b25b609 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/cancel"
    i1c8db7e05868a0cafb415b29138075693c961fa5e33b194b9814d37594776fa3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/attachments/item"
    id22b22085a212d960602b7515034a17a3a20f6b6d057247a2fdc61ffd7c670c7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/multivalueextendedproperties/item"
    id90eb8213af7914ce85043314a795aaf85858ffbbb8b82006be605ab9841293e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/extensions/item"
    iefc15682fcc08d84e696ba4641824094ecee2f76501173d5589e0b22c0d722f7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item/instances/item/singlevalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i3e367fd62461b1e59d434b1a9a1979a1913780b0085f47513db6136e5be05f43.AcceptRequestBuilder) {
    return i3e367fd62461b1e59d434b1a9a1979a1913780b0085f47513db6136e5be05f43.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*icabac1b170b5a175bf40ceb126e85ca0f20d6ee5bb8f8ab53809cd2e31a78331.AttachmentsRequestBuilder) {
    return icabac1b170b5a175bf40ceb126e85ca0f20d6ee5bb8f8ab53809cd2e31a78331.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i1c8db7e05868a0cafb415b29138075693c961fa5e33b194b9814d37594776fa3.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i1c8db7e05868a0cafb415b29138075693c961fa5e33b194b9814d37594776fa3.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i8fad80f883c81dbb1610d3530af0b03447f5b648bdebd725cd0179a10144e8bf.CalendarRequestBuilder) {
    return i8fad80f883c81dbb1610d3530af0b03447f5b648bdebd725cd0179a10144e8bf.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*if13bf5db4c861d2c48743a82cb0fd3db30f4e2e2265b5c93b6a107b69b25b609.CancelRequestBuilder) {
    return if13bf5db4c861d2c48743a82cb0fd3db30f4e2e2265b5c93b6a107b69b25b609.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/calendar/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*if0a1b5804d4380cb388c51cff9421ac2b14ff858b9ac7ce1dd5e7b91d97b1860.DeclineRequestBuilder) {
    return if0a1b5804d4380cb388c51cff9421ac2b14ff858b9ac7ce1dd5e7b91d97b1860.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i1eb540a156784a1c1466788d5998a6259a5bf9f8a6ad1d2f5c5350dd62845d31.DismissReminderRequestBuilder) {
    return i1eb540a156784a1c1466788d5998a6259a5bf9f8a6ad1d2f5c5350dd62845d31.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i374a3886dd17efd341a3266630c9c0193af40ed92cb05ffed90478816503f921.ExtensionsRequestBuilder) {
    return i374a3886dd17efd341a3266630c9c0193af40ed92cb05ffed90478816503f921.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*id90eb8213af7914ce85043314a795aaf85858ffbbb8b82006be605ab9841293e.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return id90eb8213af7914ce85043314a795aaf85858ffbbb8b82006be605ab9841293e.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i7723196f3625220ead5e66ac8e7564fd63f8a974c7523942da7e11c443dce12c.ForwardRequestBuilder) {
    return i7723196f3625220ead5e66ac8e7564fd63f8a974c7523942da7e11c443dce12c.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i0323a63dffab877f26880ff9d898dd683b137f110e681d7d437282a535f11c87.MultiValueExtendedPropertiesRequestBuilder) {
    return i0323a63dffab877f26880ff9d898dd683b137f110e681d7d437282a535f11c87.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*id22b22085a212d960602b7515034a17a3a20f6b6d057247a2fdc61ffd7c670c7.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return id22b22085a212d960602b7515034a17a3a20f6b6d057247a2fdc61ffd7c670c7.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*idebd0716278787fa5fab28a133715f58f2552b8fe8e2f7a9071417bf972eadeb.SingleValueExtendedPropertiesRequestBuilder) {
    return idebd0716278787fa5fab28a133715f58f2552b8fe8e2f7a9071417bf972eadeb.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*iefc15682fcc08d84e696ba4641824094ecee2f76501173d5589e0b22c0d722f7.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return iefc15682fcc08d84e696ba4641824094ecee2f76501173d5589e0b22c0d722f7.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i838bdc040b67a67f42d3a319019f7d751e6f2f638f8176040470892eda9fbefb.SnoozeReminderRequestBuilder) {
    return i838bdc040b67a67f42d3a319019f7d751e6f2f638f8176040470892eda9fbefb.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i053391bc6adcb52e7a06235af432dde4b91f64c74f6aad7c530db1a2ac8d9fd5.TentativelyAcceptRequestBuilder) {
    return i053391bc6adcb52e7a06235af432dde4b91f64c74f6aad7c530db1a2ac8d9fd5.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
