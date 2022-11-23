package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1398a2b6b7b6fa6075615f0b696c5d5d393449822e80a67ff6e1f1b006570589 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/cancel"
    i23879bc3fa3602cbaf93cee98173f588cf737baf3fa364fa5d532b4a7ccffb6f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances"
    i4b518553c867ab557e4a9014d64bff52b4a5614abc228a534b622097747f58c4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/extensions"
    i855851d197827c1ee6475b539166033ff1c226b94affc6b7876227a822cefe56 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/singlevalueextendedproperties"
    i91e19aea3690b9482c589e1b4fffc198beb93c8882e7ad5dc30a2c61b975241d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/forward"
    i967611b2c88418ecd27c9ba55d8e1a0fb980af96f161a7941b0122d0af98e20b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/multivalueextendedproperties"
    i973bdbcda3dd1f6f25a1482be6ac8c9929a425f352242a5f72ad56ed2695a133 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/dismissreminder"
    i9d4cd76f12610ad565059847d3a79b13115ec8cc1967d272b7f9e33b62908d70 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/decline"
    ia9eabac6ecd410545dc4495208341bd846651d6673db813cabee1323b7ccebb1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/attachments"
    iaf6ca7f88694008763f7c64dc5f9e9a4696d364070e649a47f06c9cbb9720f79 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/accept"
    ic09438710245b2671bf04abe5db2656009ca7799d1d80105aa23a063cf2f4e40 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/tentativelyaccept"
    idaa6a05522f747ce1a1a1c6a68f0551d5b3d05627b4ce5ddef9c85b80d4684b2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/calendar"
    ie274417b453065e3b9c1c84afef158b9149fe29aa54c90fc8fe2a8361a6a6880 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/snoozereminder"
    i3924f06647673effb94eed5dd5d2d7fe2b5b2f58ac82d0f02eab48c1219f1c07 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/multivalueextendedproperties/item"
    i5042ede149a551ff0f8fb6fd788882bbcfd795931aefebb153289527246e5e4d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/instances/item"
    i68e56f812a05f1c50a7fc8625d20c43951316b2122068afdbe3b9ec586bd6fe2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/extensions/item"
    i9fcaa4fe7aa3132d8c368ba6abc3e2992d0ff4b8a67614921b91e92aab8cab92 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/singlevalueextendedproperties/item"
    ia3fe1bf9b023096d20271ede0e7f809b9ea0d5dc59475200a395629f1ecb917c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/events/item/attachments/item"
)

// EventItemRequestBuilder provides operations to manage the events property of the microsoft.graph.group entity.
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
// EventItemRequestBuilderGetQueryParameters the group's calendar events.
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
func (m *EventItemRequestBuilder) Accept()(*iaf6ca7f88694008763f7c64dc5f9e9a4696d364070e649a47f06c9cbb9720f79.AcceptRequestBuilder) {
    return iaf6ca7f88694008763f7c64dc5f9e9a4696d364070e649a47f06c9cbb9720f79.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*ia9eabac6ecd410545dc4495208341bd846651d6673db813cabee1323b7ccebb1.AttachmentsRequestBuilder) {
    return ia9eabac6ecd410545dc4495208341bd846651d6673db813cabee1323b7ccebb1.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ia3fe1bf9b023096d20271ede0e7f809b9ea0d5dc59475200a395629f1ecb917c.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ia3fe1bf9b023096d20271ede0e7f809b9ea0d5dc59475200a395629f1ecb917c.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*idaa6a05522f747ce1a1a1c6a68f0551d5b3d05627b4ce5ddef9c85b80d4684b2.CalendarRequestBuilder) {
    return idaa6a05522f747ce1a1a1c6a68f0551d5b3d05627b4ce5ddef9c85b80d4684b2.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i1398a2b6b7b6fa6075615f0b696c5d5d393449822e80a67ff6e1f1b006570589.CancelRequestBuilder) {
    return i1398a2b6b7b6fa6075615f0b696c5d5d393449822e80a67ff6e1f1b006570589.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/events/{event%2Did}{?%24select}";
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
// CreateGetRequestInformation the group's calendar events.
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
func (m *EventItemRequestBuilder) Decline()(*i9d4cd76f12610ad565059847d3a79b13115ec8cc1967d272b7f9e33b62908d70.DeclineRequestBuilder) {
    return i9d4cd76f12610ad565059847d3a79b13115ec8cc1967d272b7f9e33b62908d70.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) DismissReminder()(*i973bdbcda3dd1f6f25a1482be6ac8c9929a425f352242a5f72ad56ed2695a133.DismissReminderRequestBuilder) {
    return i973bdbcda3dd1f6f25a1482be6ac8c9929a425f352242a5f72ad56ed2695a133.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i4b518553c867ab557e4a9014d64bff52b4a5614abc228a534b622097747f58c4.ExtensionsRequestBuilder) {
    return i4b518553c867ab557e4a9014d64bff52b4a5614abc228a534b622097747f58c4.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i68e56f812a05f1c50a7fc8625d20c43951316b2122068afdbe3b9ec586bd6fe2.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i68e56f812a05f1c50a7fc8625d20c43951316b2122068afdbe3b9ec586bd6fe2.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i91e19aea3690b9482c589e1b4fffc198beb93c8882e7ad5dc30a2c61b975241d.ForwardRequestBuilder) {
    return i91e19aea3690b9482c589e1b4fffc198beb93c8882e7ad5dc30a2c61b975241d.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the group's calendar events.
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
func (m *EventItemRequestBuilder) Instances()(*i23879bc3fa3602cbaf93cee98173f588cf737baf3fa364fa5d532b4a7ccffb6f.InstancesRequestBuilder) {
    return i23879bc3fa3602cbaf93cee98173f588cf737baf3fa364fa5d532b4a7ccffb6f.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i5042ede149a551ff0f8fb6fd788882bbcfd795931aefebb153289527246e5e4d.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i5042ede149a551ff0f8fb6fd788882bbcfd795931aefebb153289527246e5e4d.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i967611b2c88418ecd27c9ba55d8e1a0fb980af96f161a7941b0122d0af98e20b.MultiValueExtendedPropertiesRequestBuilder) {
    return i967611b2c88418ecd27c9ba55d8e1a0fb980af96f161a7941b0122d0af98e20b.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i3924f06647673effb94eed5dd5d2d7fe2b5b2f58ac82d0f02eab48c1219f1c07.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i3924f06647673effb94eed5dd5d2d7fe2b5b2f58ac82d0f02eab48c1219f1c07.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i855851d197827c1ee6475b539166033ff1c226b94affc6b7876227a822cefe56.SingleValueExtendedPropertiesRequestBuilder) {
    return i855851d197827c1ee6475b539166033ff1c226b94affc6b7876227a822cefe56.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i9fcaa4fe7aa3132d8c368ba6abc3e2992d0ff4b8a67614921b91e92aab8cab92.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i9fcaa4fe7aa3132d8c368ba6abc3e2992d0ff4b8a67614921b91e92aab8cab92.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ie274417b453065e3b9c1c84afef158b9149fe29aa54c90fc8fe2a8361a6a6880.SnoozeReminderRequestBuilder) {
    return ie274417b453065e3b9c1c84afef158b9149fe29aa54c90fc8fe2a8361a6a6880.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*ic09438710245b2671bf04abe5db2656009ca7799d1d80105aa23a063cf2f4e40.TentativelyAcceptRequestBuilder) {
    return ic09438710245b2671bf04abe5db2656009ca7799d1d80105aa23a063cf2f4e40.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
