package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i07aefebef2360a05f1f978a4dacfbfc958f6263953281cb4e4810c24e8fdbe51 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/forward"
    i130329fd346118a25e215e9ae08e6acafc2fb545bcc2544394704d98040e0fd5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/tentativelyaccept"
    i1b6d6c642ae53b6d75418c14f94432b149eb82d231e9ef4719cabf3dd36e4a84 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/decline"
    i47c119d74e0c4d1b7732433a7add8a684909bad78a2169bdadc752f8052b2cde "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/attachments"
    i8bc1b83864742db92b87f3ce691ba7194b9e4dfcdee44bc1f989db1fd24e0d9c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/dismissreminder"
    i8daf6ca6792aac2f5e22afadd880b9416a73f6d66f3d3a4cb4ce49b702beabff "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances"
    iaafeb0437eaa3b26607893fda20fbdae77ac9d7e2333ed4d5b73883b77acff0a "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/snoozereminder"
    iabd200b2c68fcfd56b3a70da55accf4de8cff2e02b0264e62b3ec575b67993fd "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/cancel"
    ib1691f681b565035438b2134f4d2fff4a5be77dd06fbf2cefb004ae233672e01 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/calendar"
    ic2ca36017987ae7a5daa23b834522d6cb1b251d72fe6d72b471f8ca4ed2d2264 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/multivalueextendedproperties"
    ic9d89564024358edace4430f9c467c2085e6b36f70aa7b6fb0605f47881321b4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/extensions"
    icfb49c8c2637abda0a6e0b8db57c259abb3d4ec46f18de110de6cc049d986341 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/singlevalueextendedproperties"
    if35aa5926f23e100256e1aa6716e7d17471936c8947905b607e1035b11cc8f34 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/accept"
    i065657a3daf360b62ebf3c6e101d6bd29d67ca95848bd319d1ffeb09214c3cd5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/attachments/item"
    i49c7e84876761d3a9df15aa2f74f21992b0ba54e8eea9060140c4135b4b76920 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/singlevalueextendedproperties/item"
    i9facca5460444460406e305a714abf1b6c6c7e0be4c7bedda91eb9b394bd9ede "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/multivalueextendedproperties/item"
    ibe0f4b62b27e66ba9582cbb8b856bbb6bd0e3e8a30c85f69943b3fc01aac856e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/instances/item"
    ibf212a370d44f80b700c5b5b75239c5e02548b6abbd3484f7872aa4674bf95c4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item/extensions/item"
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
func (m *EventItemRequestBuilder) Accept()(*if35aa5926f23e100256e1aa6716e7d17471936c8947905b607e1035b11cc8f34.AcceptRequestBuilder) {
    return if35aa5926f23e100256e1aa6716e7d17471936c8947905b607e1035b11cc8f34.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i47c119d74e0c4d1b7732433a7add8a684909bad78a2169bdadc752f8052b2cde.AttachmentsRequestBuilder) {
    return i47c119d74e0c4d1b7732433a7add8a684909bad78a2169bdadc752f8052b2cde.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i065657a3daf360b62ebf3c6e101d6bd29d67ca95848bd319d1ffeb09214c3cd5.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i065657a3daf360b62ebf3c6e101d6bd29d67ca95848bd319d1ffeb09214c3cd5.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ib1691f681b565035438b2134f4d2fff4a5be77dd06fbf2cefb004ae233672e01.CalendarRequestBuilder) {
    return ib1691f681b565035438b2134f4d2fff4a5be77dd06fbf2cefb004ae233672e01.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*iabd200b2c68fcfd56b3a70da55accf4de8cff2e02b0264e62b3ec575b67993fd.CancelRequestBuilder) {
    return iabd200b2c68fcfd56b3a70da55accf4de8cff2e02b0264e62b3ec575b67993fd.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendar/events/{event%2Did}{?%24select}";
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
// CreateDeleteRequestInformation delete navigation property events for users
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
// CreatePatchRequestInformation update the navigation property events in users
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
func (m *EventItemRequestBuilder) Decline()(*i1b6d6c642ae53b6d75418c14f94432b149eb82d231e9ef4719cabf3dd36e4a84.DeclineRequestBuilder) {
    return i1b6d6c642ae53b6d75418c14f94432b149eb82d231e9ef4719cabf3dd36e4a84.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property events for users
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
func (m *EventItemRequestBuilder) DismissReminder()(*i8bc1b83864742db92b87f3ce691ba7194b9e4dfcdee44bc1f989db1fd24e0d9c.DismissReminderRequestBuilder) {
    return i8bc1b83864742db92b87f3ce691ba7194b9e4dfcdee44bc1f989db1fd24e0d9c.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*ic9d89564024358edace4430f9c467c2085e6b36f70aa7b6fb0605f47881321b4.ExtensionsRequestBuilder) {
    return ic9d89564024358edace4430f9c467c2085e6b36f70aa7b6fb0605f47881321b4.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ibf212a370d44f80b700c5b5b75239c5e02548b6abbd3484f7872aa4674bf95c4.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ibf212a370d44f80b700c5b5b75239c5e02548b6abbd3484f7872aa4674bf95c4.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i07aefebef2360a05f1f978a4dacfbfc958f6263953281cb4e4810c24e8fdbe51.ForwardRequestBuilder) {
    return i07aefebef2360a05f1f978a4dacfbfc958f6263953281cb4e4810c24e8fdbe51.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*i8daf6ca6792aac2f5e22afadd880b9416a73f6d66f3d3a4cb4ce49b702beabff.InstancesRequestBuilder) {
    return i8daf6ca6792aac2f5e22afadd880b9416a73f6d66f3d3a4cb4ce49b702beabff.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*ibe0f4b62b27e66ba9582cbb8b856bbb6bd0e3e8a30c85f69943b3fc01aac856e.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return ibe0f4b62b27e66ba9582cbb8b856bbb6bd0e3e8a30c85f69943b3fc01aac856e.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ic2ca36017987ae7a5daa23b834522d6cb1b251d72fe6d72b471f8ca4ed2d2264.MultiValueExtendedPropertiesRequestBuilder) {
    return ic2ca36017987ae7a5daa23b834522d6cb1b251d72fe6d72b471f8ca4ed2d2264.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i9facca5460444460406e305a714abf1b6c6c7e0be4c7bedda91eb9b394bd9ede.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i9facca5460444460406e305a714abf1b6c6c7e0be4c7bedda91eb9b394bd9ede.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property events in users
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
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*icfb49c8c2637abda0a6e0b8db57c259abb3d4ec46f18de110de6cc049d986341.SingleValueExtendedPropertiesRequestBuilder) {
    return icfb49c8c2637abda0a6e0b8db57c259abb3d4ec46f18de110de6cc049d986341.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i49c7e84876761d3a9df15aa2f74f21992b0ba54e8eea9060140c4135b4b76920.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i49c7e84876761d3a9df15aa2f74f21992b0ba54e8eea9060140c4135b4b76920.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*iaafeb0437eaa3b26607893fda20fbdae77ac9d7e2333ed4d5b73883b77acff0a.SnoozeReminderRequestBuilder) {
    return iaafeb0437eaa3b26607893fda20fbdae77ac9d7e2333ed4d5b73883b77acff0a.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i130329fd346118a25e215e9ae08e6acafc2fb545bcc2544394704d98040e0fd5.TentativelyAcceptRequestBuilder) {
    return i130329fd346118a25e215e9ae08e6acafc2fb545bcc2544394704d98040e0fd5.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
