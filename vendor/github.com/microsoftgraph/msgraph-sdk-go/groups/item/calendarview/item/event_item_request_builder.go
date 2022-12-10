package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i24e24e2af499471f8e4f73b6d94d99165b1e84d2131623dcc1a383199e8f050d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/singlevalueextendedproperties"
    i2d1fe0bd4017d6a5f945a9ceefc6fde8147a060a766d7001c9b0a50c6ba6ee00 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/multivalueextendedproperties"
    i44dd09532dcb7a7d8ed2a9efa5056d3e20b725877822e9b1c509e378c8db9c90 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/snoozereminder"
    i493deeccf2b51ea98ffd5d94edc505e26b0559809a7e344b863c50a129ced7fe "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/tentativelyaccept"
    i555550f048433b0f866e19c810c2cac5a25998a168b6c16dc1a7f6a70665fd92 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/calendar"
    i61b71d3da7a06d0f54197cc290f668d03c0ff14f20103e48cab218d31dab0913 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/extensions"
    i6c90a1bd7a4f0cece06e6b34e378ee70ca26356395f3e3815e647b33c612c4eb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/forward"
    i6d4c88f7616ab1a1de0f125a057a6e9383f4b21373f7c695fd9c6db145fa1754 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances"
    i7a68ed9d6c87ce91e310c0b60439fa878a2f3f2ce780223e121fe57306de4cf0 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/attachments"
    i89ff886574ba825b3220369ac4d9bdb3e03e16efbd810974d1a8b0c6a7d97cf8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/decline"
    i945af8cb490de0be43df42f623102f1a698bf30d1e58d7b44500a64140273ec6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/accept"
    idf3581de0ced5dd508d8f36e7490c56e109571dd437697fe9c1edb7bce3dac89 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/dismissreminder"
    ie208e488c4a77bd04b8849d39ad4a5b1c796289043cd257b42ba1a96585cd84e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/cancel"
    i243da75f0745ff07c19b77e029d0409d2627daa65cc8fff1c986fb76b2b765a9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/attachments/item"
    i64313c875409e11e8ad455ef0785e8f2d79c5d511bf6a34bdeec84f3da8c6795 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/singlevalueextendedproperties/item"
    i66a62ea9b41daca5ad35a026f8fdd67e1aa3e96483f9cb8758db312e258ae78f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/instances/item"
    i71e82e42e7ad7a048127c6e49b504d33eb9ea367867057c7a9bf3cb88ef2c551 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/extensions/item"
    ib73431f50c914fa465532d76bca0db38ab141314b5872ddb060382316b1988dd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendarview/item/multivalueextendedproperties/item"
)

// EventItemRequestBuilder provides operations to manage the calendarView property of the microsoft.graph.group entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderGetQueryParameters the calendar view for the calendar. Read-only.
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
func (m *EventItemRequestBuilder) Accept()(*i945af8cb490de0be43df42f623102f1a698bf30d1e58d7b44500a64140273ec6.AcceptRequestBuilder) {
    return i945af8cb490de0be43df42f623102f1a698bf30d1e58d7b44500a64140273ec6.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i7a68ed9d6c87ce91e310c0b60439fa878a2f3f2ce780223e121fe57306de4cf0.AttachmentsRequestBuilder) {
    return i7a68ed9d6c87ce91e310c0b60439fa878a2f3f2ce780223e121fe57306de4cf0.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i243da75f0745ff07c19b77e029d0409d2627daa65cc8fff1c986fb76b2b765a9.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i243da75f0745ff07c19b77e029d0409d2627daa65cc8fff1c986fb76b2b765a9.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i555550f048433b0f866e19c810c2cac5a25998a168b6c16dc1a7f6a70665fd92.CalendarRequestBuilder) {
    return i555550f048433b0f866e19c810c2cac5a25998a168b6c16dc1a7f6a70665fd92.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ie208e488c4a77bd04b8849d39ad4a5b1c796289043cd257b42ba1a96585cd84e.CancelRequestBuilder) {
    return ie208e488c4a77bd04b8849d39ad4a5b1c796289043cd257b42ba1a96585cd84e.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
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
// CreateGetRequestInformation the calendar view for the calendar. Read-only.
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
func (m *EventItemRequestBuilder) Decline()(*i89ff886574ba825b3220369ac4d9bdb3e03e16efbd810974d1a8b0c6a7d97cf8.DeclineRequestBuilder) {
    return i89ff886574ba825b3220369ac4d9bdb3e03e16efbd810974d1a8b0c6a7d97cf8.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*idf3581de0ced5dd508d8f36e7490c56e109571dd437697fe9c1edb7bce3dac89.DismissReminderRequestBuilder) {
    return idf3581de0ced5dd508d8f36e7490c56e109571dd437697fe9c1edb7bce3dac89.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i61b71d3da7a06d0f54197cc290f668d03c0ff14f20103e48cab218d31dab0913.ExtensionsRequestBuilder) {
    return i61b71d3da7a06d0f54197cc290f668d03c0ff14f20103e48cab218d31dab0913.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i71e82e42e7ad7a048127c6e49b504d33eb9ea367867057c7a9bf3cb88ef2c551.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i71e82e42e7ad7a048127c6e49b504d33eb9ea367867057c7a9bf3cb88ef2c551.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i6c90a1bd7a4f0cece06e6b34e378ee70ca26356395f3e3815e647b33c612c4eb.ForwardRequestBuilder) {
    return i6c90a1bd7a4f0cece06e6b34e378ee70ca26356395f3e3815e647b33c612c4eb.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the calendar view for the calendar. Read-only.
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
func (m *EventItemRequestBuilder) Instances()(*i6d4c88f7616ab1a1de0f125a057a6e9383f4b21373f7c695fd9c6db145fa1754.InstancesRequestBuilder) {
    return i6d4c88f7616ab1a1de0f125a057a6e9383f4b21373f7c695fd9c6db145fa1754.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i66a62ea9b41daca5ad35a026f8fdd67e1aa3e96483f9cb8758db312e258ae78f.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i66a62ea9b41daca5ad35a026f8fdd67e1aa3e96483f9cb8758db312e258ae78f.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i2d1fe0bd4017d6a5f945a9ceefc6fde8147a060a766d7001c9b0a50c6ba6ee00.MultiValueExtendedPropertiesRequestBuilder) {
    return i2d1fe0bd4017d6a5f945a9ceefc6fde8147a060a766d7001c9b0a50c6ba6ee00.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ib73431f50c914fa465532d76bca0db38ab141314b5872ddb060382316b1988dd.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ib73431f50c914fa465532d76bca0db38ab141314b5872ddb060382316b1988dd.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i24e24e2af499471f8e4f73b6d94d99165b1e84d2131623dcc1a383199e8f050d.SingleValueExtendedPropertiesRequestBuilder) {
    return i24e24e2af499471f8e4f73b6d94d99165b1e84d2131623dcc1a383199e8f050d.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i64313c875409e11e8ad455ef0785e8f2d79c5d511bf6a34bdeec84f3da8c6795.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i64313c875409e11e8ad455ef0785e8f2d79c5d511bf6a34bdeec84f3da8c6795.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i44dd09532dcb7a7d8ed2a9efa5056d3e20b725877822e9b1c509e378c8db9c90.SnoozeReminderRequestBuilder) {
    return i44dd09532dcb7a7d8ed2a9efa5056d3e20b725877822e9b1c509e378c8db9c90.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i493deeccf2b51ea98ffd5d94edc505e26b0559809a7e344b863c50a129ced7fe.TentativelyAcceptRequestBuilder) {
    return i493deeccf2b51ea98ffd5d94edc505e26b0559809a7e344b863c50a129ced7fe.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
