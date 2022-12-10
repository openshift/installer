package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0a0fc3d004d78576d043eb521ad26f53fe8214afd2b3ebe973fe93e60f680f71 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/attachments"
    i16240d0de448bc4c6216c08277550221d94bc0f48ba2accd218749ad7c953187 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/extensions"
    i3c7b0004be1fa25f87efde7027fd8e976d9cf7bf986843ed38a31557bf5e7a95 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/accept"
    i49c9c1594c0f4696377dfecc53ccf788c8e8b710f8313547cbcd0af990f5bbf6 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/cancel"
    i68b73555ae5e3f902578ad78535a4acc022d87da89f72cc819d660c671ed5ece "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/tentativelyaccept"
    i8bfe1260d815ef118a4accfe78279062de1c84617b6bf2fd2fcbee8aaa6e4a3e "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/multivalueextendedproperties"
    ib177b4b7402e2bd1b4c939d9165f6cc6ad4c888362c100157bbf3567e762c3a8 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/singlevalueextendedproperties"
    ib7d6bec077bf9e6ed260a5dcae2c2e771c594d05995ed7713b4f6d55c140e8e9 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/snoozereminder"
    ic273ea006cf9b2cbb523ae048eff477ad9414b6846e6486ee9ecb96844555f5e "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/forward"
    idc8f27b80133ba40f0ebae7ae73b0927fa97d40f63c920f7718b0b1bb3646f33 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/decline"
    idcb9651181c4cfa39068cf83f36f38021b6d83f2785c96c836454a8066720a49 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/dismissreminder"
    if25488a44284a9c63a32442f9d9a88bb047100fbcbc4cbe3daa8228fb39d90d7 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/calendar"
    i684e98bc756b097dd50799c2b776635f96b124dbc95dd1f25f6d308ac7ecfa6f "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/attachments/item"
    ia70c4b879b44928811313bb8380c609fb4f17a2429974bfc650c78042dd445a2 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/extensions/item"
    ib4ab2585af40e6299bd85e8d11fe39d4c9e2950cba8aaab47b74a63390ade002 "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/singlevalueextendedproperties/item"
    ifb082cc9281770b0f14bb66061bb2f300737211dde2040e7b7d2bb48d23b51ae "github.com/microsoftgraph/msgraph-sdk-go/me/events/item/instances/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i3c7b0004be1fa25f87efde7027fd8e976d9cf7bf986843ed38a31557bf5e7a95.AcceptRequestBuilder) {
    return i3c7b0004be1fa25f87efde7027fd8e976d9cf7bf986843ed38a31557bf5e7a95.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i0a0fc3d004d78576d043eb521ad26f53fe8214afd2b3ebe973fe93e60f680f71.AttachmentsRequestBuilder) {
    return i0a0fc3d004d78576d043eb521ad26f53fe8214afd2b3ebe973fe93e60f680f71.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i684e98bc756b097dd50799c2b776635f96b124dbc95dd1f25f6d308ac7ecfa6f.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i684e98bc756b097dd50799c2b776635f96b124dbc95dd1f25f6d308ac7ecfa6f.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*if25488a44284a9c63a32442f9d9a88bb047100fbcbc4cbe3daa8228fb39d90d7.CalendarRequestBuilder) {
    return if25488a44284a9c63a32442f9d9a88bb047100fbcbc4cbe3daa8228fb39d90d7.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i49c9c1594c0f4696377dfecc53ccf788c8e8b710f8313547cbcd0af990f5bbf6.CancelRequestBuilder) {
    return i49c9c1594c0f4696377dfecc53ccf788c8e8b710f8313547cbcd0af990f5bbf6.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/events/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*idc8f27b80133ba40f0ebae7ae73b0927fa97d40f63c920f7718b0b1bb3646f33.DeclineRequestBuilder) {
    return idc8f27b80133ba40f0ebae7ae73b0927fa97d40f63c920f7718b0b1bb3646f33.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*idcb9651181c4cfa39068cf83f36f38021b6d83f2785c96c836454a8066720a49.DismissReminderRequestBuilder) {
    return idcb9651181c4cfa39068cf83f36f38021b6d83f2785c96c836454a8066720a49.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i16240d0de448bc4c6216c08277550221d94bc0f48ba2accd218749ad7c953187.ExtensionsRequestBuilder) {
    return i16240d0de448bc4c6216c08277550221d94bc0f48ba2accd218749ad7c953187.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*ia70c4b879b44928811313bb8380c609fb4f17a2429974bfc650c78042dd445a2.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ia70c4b879b44928811313bb8380c609fb4f17a2429974bfc650c78042dd445a2.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ic273ea006cf9b2cbb523ae048eff477ad9414b6846e6486ee9ecb96844555f5e.ForwardRequestBuilder) {
    return ic273ea006cf9b2cbb523ae048eff477ad9414b6846e6486ee9ecb96844555f5e.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i8bfe1260d815ef118a4accfe78279062de1c84617b6bf2fd2fcbee8aaa6e4a3e.MultiValueExtendedPropertiesRequestBuilder) {
    return i8bfe1260d815ef118a4accfe78279062de1c84617b6bf2fd2fcbee8aaa6e4a3e.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ifb082cc9281770b0f14bb66061bb2f300737211dde2040e7b7d2bb48d23b51ae.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ifb082cc9281770b0f14bb66061bb2f300737211dde2040e7b7d2bb48d23b51ae.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*ib177b4b7402e2bd1b4c939d9165f6cc6ad4c888362c100157bbf3567e762c3a8.SingleValueExtendedPropertiesRequestBuilder) {
    return ib177b4b7402e2bd1b4c939d9165f6cc6ad4c888362c100157bbf3567e762c3a8.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ib4ab2585af40e6299bd85e8d11fe39d4c9e2950cba8aaab47b74a63390ade002.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ib4ab2585af40e6299bd85e8d11fe39d4c9e2950cba8aaab47b74a63390ade002.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ib7d6bec077bf9e6ed260a5dcae2c2e771c594d05995ed7713b4f6d55c140e8e9.SnoozeReminderRequestBuilder) {
    return ib7d6bec077bf9e6ed260a5dcae2c2e771c594d05995ed7713b4f6d55c140e8e9.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i68b73555ae5e3f902578ad78535a4acc022d87da89f72cc819d660c671ed5ece.TentativelyAcceptRequestBuilder) {
    return i68b73555ae5e3f902578ad78535a4acc022d87da89f72cc819d660c671ed5ece.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
