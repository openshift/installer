package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i12c00cf6d17f76414eb8bf4af116ae0f94a375b16a3fdcc2a520db9d8b51bb69 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/calendar"
    i382bc5ab005174ca0a157f3ddbbe59b93465e03639aea98bec85f98a02704a20 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/extensions"
    i52fe9bd763a706a85cde883eea661538b98807c8ed4896d571ad9f373de4a24b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/dismissreminder"
    i8cb640b8e824bfdf805685b8399a0ddb642f5f041f6f7a6d1a2adbbd8b04242c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/attachments"
    i91a25c006065bdbda1ec6dbad739849f30b2ed0ea0a864cf75a845796dfb463a "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/cancel"
    ia41b7c24cc3f3baeaa7815e53e95f6fede450c797c61686c04ae5509cff1fee9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/accept"
    iac5010b0952b1db93baa23920b116d827c3b2ec1b7e9168259f8668f64c64732 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/tentativelyaccept"
    ib2891d93c6ded684ce7abfd2909cc1a20fc2090c50db20bbe6da0375eb2f3fe7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/instances"
    ib2f603875392f9c2b340caee4eab840215f32e0520cb72749ad970c7f46e8ced "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/decline"
    id148b2204d6e284392369bbf295c7ad603ca498161a3166593a744b1705d35d7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/multivalueextendedproperties"
    ieabc03565c45dc93b6dd73031a0f62fcc81744682614225dee84662b38c937e3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/forward"
    if1d75ca20ba55d1dcc9dfcad50836d4235106da1e2a01e8cbd7d0a07f344df20 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/singlevalueextendedproperties"
    ife09ca7c64720802742bbf52d953fd6747f5910742fbabcc2a740709adb9c2cc "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/snoozereminder"
    i8fcec19e7c3ce8aa06ca5657386121d5140cfc8138e4f515339a4475a80f52be "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/instances/item"
    iabd4185ab5ac976d8418c3afea4e9d5b739795c1459c53296cbe08e4cdb2a44a "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/extensions/item"
    id1c9a5e9fa77008270b7af81a4a4663bc2e9e1151e2512f7f7783693e4a62854 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/singlevalueextendedproperties/item"
    ie0f0bdcc9602484baf79d027f6483fa9fbafd07fe65976cf1a96cc9ced359802 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/attachments/item"
    ie456cc01b9ad3635053091de1d8edb69f6dfa30b76d79790c6668bf7307efe8e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/events/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*ia41b7c24cc3f3baeaa7815e53e95f6fede450c797c61686c04ae5509cff1fee9.AcceptRequestBuilder) {
    return ia41b7c24cc3f3baeaa7815e53e95f6fede450c797c61686c04ae5509cff1fee9.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i8cb640b8e824bfdf805685b8399a0ddb642f5f041f6f7a6d1a2adbbd8b04242c.AttachmentsRequestBuilder) {
    return i8cb640b8e824bfdf805685b8399a0ddb642f5f041f6f7a6d1a2adbbd8b04242c.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ie0f0bdcc9602484baf79d027f6483fa9fbafd07fe65976cf1a96cc9ced359802.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ie0f0bdcc9602484baf79d027f6483fa9fbafd07fe65976cf1a96cc9ced359802.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i12c00cf6d17f76414eb8bf4af116ae0f94a375b16a3fdcc2a520db9d8b51bb69.CalendarRequestBuilder) {
    return i12c00cf6d17f76414eb8bf4af116ae0f94a375b16a3fdcc2a520db9d8b51bb69.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i91a25c006065bdbda1ec6dbad739849f30b2ed0ea0a864cf75a845796dfb463a.CancelRequestBuilder) {
    return i91a25c006065bdbda1ec6dbad739849f30b2ed0ea0a864cf75a845796dfb463a.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendars/{calendar%2Did}/events/{event%2Did}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*ib2f603875392f9c2b340caee4eab840215f32e0520cb72749ad970c7f46e8ced.DeclineRequestBuilder) {
    return ib2f603875392f9c2b340caee4eab840215f32e0520cb72749ad970c7f46e8ced.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) DismissReminder()(*i52fe9bd763a706a85cde883eea661538b98807c8ed4896d571ad9f373de4a24b.DismissReminderRequestBuilder) {
    return i52fe9bd763a706a85cde883eea661538b98807c8ed4896d571ad9f373de4a24b.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i382bc5ab005174ca0a157f3ddbbe59b93465e03639aea98bec85f98a02704a20.ExtensionsRequestBuilder) {
    return i382bc5ab005174ca0a157f3ddbbe59b93465e03639aea98bec85f98a02704a20.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*iabd4185ab5ac976d8418c3afea4e9d5b739795c1459c53296cbe08e4cdb2a44a.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return iabd4185ab5ac976d8418c3afea4e9d5b739795c1459c53296cbe08e4cdb2a44a.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ieabc03565c45dc93b6dd73031a0f62fcc81744682614225dee84662b38c937e3.ForwardRequestBuilder) {
    return ieabc03565c45dc93b6dd73031a0f62fcc81744682614225dee84662b38c937e3.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) Instances()(*ib2891d93c6ded684ce7abfd2909cc1a20fc2090c50db20bbe6da0375eb2f3fe7.InstancesRequestBuilder) {
    return ib2891d93c6ded684ce7abfd2909cc1a20fc2090c50db20bbe6da0375eb2f3fe7.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*i8fcec19e7c3ce8aa06ca5657386121d5140cfc8138e4f515339a4475a80f52be.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return i8fcec19e7c3ce8aa06ca5657386121d5140cfc8138e4f515339a4475a80f52be.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*id148b2204d6e284392369bbf295c7ad603ca498161a3166593a744b1705d35d7.MultiValueExtendedPropertiesRequestBuilder) {
    return id148b2204d6e284392369bbf295c7ad603ca498161a3166593a744b1705d35d7.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ie456cc01b9ad3635053091de1d8edb69f6dfa30b76d79790c6668bf7307efe8e.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ie456cc01b9ad3635053091de1d8edb69f6dfa30b76d79790c6668bf7307efe8e.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*if1d75ca20ba55d1dcc9dfcad50836d4235106da1e2a01e8cbd7d0a07f344df20.SingleValueExtendedPropertiesRequestBuilder) {
    return if1d75ca20ba55d1dcc9dfcad50836d4235106da1e2a01e8cbd7d0a07f344df20.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*id1c9a5e9fa77008270b7af81a4a4663bc2e9e1151e2512f7f7783693e4a62854.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return id1c9a5e9fa77008270b7af81a4a4663bc2e9e1151e2512f7f7783693e4a62854.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ife09ca7c64720802742bbf52d953fd6747f5910742fbabcc2a740709adb9c2cc.SnoozeReminderRequestBuilder) {
    return ife09ca7c64720802742bbf52d953fd6747f5910742fbabcc2a740709adb9c2cc.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*iac5010b0952b1db93baa23920b116d827c3b2ec1b7e9168259f8668f64c64732.TentativelyAcceptRequestBuilder) {
    return iac5010b0952b1db93baa23920b116d827c3b2ec1b7e9168259f8668f64c64732.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
