package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i115424fb947e5363ae5a2dbcd21aedc22f38dffe0cc73f1ef73b25306b7603d1 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/snoozereminder"
    i4ecf9972e71ceed18714c2f2e8df53347f0e3eb946fa092de94977ea4fed1ceb "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/decline"
    i5121a60846fa1d710d99f63cb11af175a28ec5079443325e118d875f5c00ff43 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/forward"
    i6fcb984531f9e619175ab5500d10abe7f35a19a0222db8119e4558b14210ee7e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/accept"
    i7ed3cb920739c72fa84a300b2103c5b0d218576cca12359b662685935f0eaf1e "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/calendar"
    i896eff89b7229688e940bcb6969ce2d2c0bdb24746c8c9588b2fd26db63028d7 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/attachments"
    ib0e713445b5a213b2e1bc107fd4c5258be4b72d7d08cae4177f0f252ecc674e9 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/singlevalueextendedproperties"
    ibf55e0cbfd799c432674550c326d27458a9d46ad44fb6847daaa94f3085266af "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/cancel"
    ic320338b3b6969f4fefa9b2514352265bf35f14d06f238916bcfa74353b23abe "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/multivalueextendedproperties"
    icedd85fbcd7041ab524d00f58e3fc8cd64b80a6040a18f66af1097654d14e25a "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/extensions"
    if70024e0e615a818d126650626f99631a3f58b7c2c14bc1a5a482345208428b3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/dismissreminder"
    ife044224236e6edea2f74297331ca97a9b058bbe292a86fada70da05707d7460 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/tentativelyaccept"
    i1f1103b1f10d19af70cdb5d44f1c4993d2eb251eb003bda9353d58a870ade6da "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/extensions/item"
    i2879a2ee532e317841d70b52cee430d07751f97864ca66e3c10e502ec78ce73a "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/multivalueextendedproperties/item"
    i2a39f49bad40b8b085e6b6ade7998f88875082e8dcb30589596db0523283fe95 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/singlevalueextendedproperties/item"
    i6b2bff425bc13e13fd61c4d59ba9a0d4738b430aa1f1e13e97659caefc058e2d "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item/instances/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*i6fcb984531f9e619175ab5500d10abe7f35a19a0222db8119e4558b14210ee7e.AcceptRequestBuilder) {
    return i6fcb984531f9e619175ab5500d10abe7f35a19a0222db8119e4558b14210ee7e.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i896eff89b7229688e940bcb6969ce2d2c0bdb24746c8c9588b2fd26db63028d7.AttachmentsRequestBuilder) {
    return i896eff89b7229688e940bcb6969ce2d2c0bdb24746c8c9588b2fd26db63028d7.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i6b2bff425bc13e13fd61c4d59ba9a0d4738b430aa1f1e13e97659caefc058e2d.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i6b2bff425bc13e13fd61c4d59ba9a0d4738b430aa1f1e13e97659caefc058e2d.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i7ed3cb920739c72fa84a300b2103c5b0d218576cca12359b662685935f0eaf1e.CalendarRequestBuilder) {
    return i7ed3cb920739c72fa84a300b2103c5b0d218576cca12359b662685935f0eaf1e.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*ibf55e0cbfd799c432674550c326d27458a9d46ad44fb6847daaa94f3085266af.CancelRequestBuilder) {
    return ibf55e0cbfd799c432674550c326d27458a9d46ad44fb6847daaa94f3085266af.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendar/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i4ecf9972e71ceed18714c2f2e8df53347f0e3eb946fa092de94977ea4fed1ceb.DeclineRequestBuilder) {
    return i4ecf9972e71ceed18714c2f2e8df53347f0e3eb946fa092de94977ea4fed1ceb.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*if70024e0e615a818d126650626f99631a3f58b7c2c14bc1a5a482345208428b3.DismissReminderRequestBuilder) {
    return if70024e0e615a818d126650626f99631a3f58b7c2c14bc1a5a482345208428b3.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*icedd85fbcd7041ab524d00f58e3fc8cd64b80a6040a18f66af1097654d14e25a.ExtensionsRequestBuilder) {
    return icedd85fbcd7041ab524d00f58e3fc8cd64b80a6040a18f66af1097654d14e25a.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i1f1103b1f10d19af70cdb5d44f1c4993d2eb251eb003bda9353d58a870ade6da.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i1f1103b1f10d19af70cdb5d44f1c4993d2eb251eb003bda9353d58a870ade6da.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i5121a60846fa1d710d99f63cb11af175a28ec5079443325e118d875f5c00ff43.ForwardRequestBuilder) {
    return i5121a60846fa1d710d99f63cb11af175a28ec5079443325e118d875f5c00ff43.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ic320338b3b6969f4fefa9b2514352265bf35f14d06f238916bcfa74353b23abe.MultiValueExtendedPropertiesRequestBuilder) {
    return ic320338b3b6969f4fefa9b2514352265bf35f14d06f238916bcfa74353b23abe.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i2879a2ee532e317841d70b52cee430d07751f97864ca66e3c10e502ec78ce73a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i2879a2ee532e317841d70b52cee430d07751f97864ca66e3c10e502ec78ce73a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*ib0e713445b5a213b2e1bc107fd4c5258be4b72d7d08cae4177f0f252ecc674e9.SingleValueExtendedPropertiesRequestBuilder) {
    return ib0e713445b5a213b2e1bc107fd4c5258be4b72d7d08cae4177f0f252ecc674e9.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i2a39f49bad40b8b085e6b6ade7998f88875082e8dcb30589596db0523283fe95.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i2a39f49bad40b8b085e6b6ade7998f88875082e8dcb30589596db0523283fe95.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*i115424fb947e5363ae5a2dbcd21aedc22f38dffe0cc73f1ef73b25306b7603d1.SnoozeReminderRequestBuilder) {
    return i115424fb947e5363ae5a2dbcd21aedc22f38dffe0cc73f1ef73b25306b7603d1.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ife044224236e6edea2f74297331ca97a9b058bbe292a86fada70da05707d7460.TentativelyAcceptRequestBuilder) {
    return ife044224236e6edea2f74297331ca97a9b058bbe292a86fada70da05707d7460.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
