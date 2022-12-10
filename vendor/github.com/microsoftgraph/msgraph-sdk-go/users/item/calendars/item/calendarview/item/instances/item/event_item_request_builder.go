package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i272324389d6bf1f60ccbfa77144be57c1d70c0718d79da10d8dee55e5f62919e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/cancel"
    i2ffebe5209de4076e8a82adab1a4bfca6c05cdf6f208c455125b48c57db65904 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties"
    i44867788db34306f0996cf494a44a02c2fdfeba2f1d248dfbf3719ef279f5d35 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/decline"
    i48b536acd4f01ac6a9b0c428b9de43a1cbde1b708b519c38213a57d87b96504f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties"
    i7cfc21431c7f3715849334d5beda9eb12748eea231818a18c33f9cd7311e2ea7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/calendar"
    i7d3117e3afb0fbd44802c1ad6d61175bbd7b795db5acdd75c10d1c1da8aed689 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/extensions"
    i81701a116c93dee173023918621f350692ddfe7873f6de037bccfa9bc683dc45 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/accept"
    i83473d2a0bc99449ff6981736c3c17fbee0979720e884cddb635db5c58ed629c "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/forward"
    i9bf36797cd9f9784c47c527ad0e9d3362b3dc78d433e24e67f96bcd93f8c8dd2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/dismissreminder"
    iba1c2f9e853ce3684320c2cf66b1bf622a32492dfdb465c4f5900f2609cac816 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/attachments"
    ic2b75b590eb1afa68b5b139fda2fa1f65d5573ce52b3716a67f1e95fa70569e8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/snoozereminder"
    idcd892355bb9b0a3694ef802b374e6be28dd95de1898988e0f5500a51feedf6f "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/tentativelyaccept"
    i0a8a1ba6dbe110e4b0d2005064e6c047ae4063f6e259f00465d057bc61bf5204 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/singlevalueextendedproperties/item"
    i7d0e215fee21f7bd421540d9e6919033aaaf98c664446238d006381d73e70691 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/attachments/item"
    icc005ad8da2b3b97a3de6350a707a53f032a558960b40cd2afcc1c28e246ada6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/extensions/item"
    ifcf0b211e141030ced11d35ff3b92d777831766af0ce38e3192dd97c536b3221 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item/calendarview/item/instances/item/multivalueextendedproperties/item"
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
func (m *EventItemRequestBuilder) Accept()(*i81701a116c93dee173023918621f350692ddfe7873f6de037bccfa9bc683dc45.AcceptRequestBuilder) {
    return i81701a116c93dee173023918621f350692ddfe7873f6de037bccfa9bc683dc45.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*iba1c2f9e853ce3684320c2cf66b1bf622a32492dfdb465c4f5900f2609cac816.AttachmentsRequestBuilder) {
    return iba1c2f9e853ce3684320c2cf66b1bf622a32492dfdb465c4f5900f2609cac816.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i7d0e215fee21f7bd421540d9e6919033aaaf98c664446238d006381d73e70691.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i7d0e215fee21f7bd421540d9e6919033aaaf98c664446238d006381d73e70691.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i7cfc21431c7f3715849334d5beda9eb12748eea231818a18c33f9cd7311e2ea7.CalendarRequestBuilder) {
    return i7cfc21431c7f3715849334d5beda9eb12748eea231818a18c33f9cd7311e2ea7.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i272324389d6bf1f60ccbfa77144be57c1d70c0718d79da10d8dee55e5f62919e.CancelRequestBuilder) {
    return i272324389d6bf1f60ccbfa77144be57c1d70c0718d79da10d8dee55e5f62919e.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendars/{calendar%2Did}/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i44867788db34306f0996cf494a44a02c2fdfeba2f1d248dfbf3719ef279f5d35.DeclineRequestBuilder) {
    return i44867788db34306f0996cf494a44a02c2fdfeba2f1d248dfbf3719ef279f5d35.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i9bf36797cd9f9784c47c527ad0e9d3362b3dc78d433e24e67f96bcd93f8c8dd2.DismissReminderRequestBuilder) {
    return i9bf36797cd9f9784c47c527ad0e9d3362b3dc78d433e24e67f96bcd93f8c8dd2.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i7d3117e3afb0fbd44802c1ad6d61175bbd7b795db5acdd75c10d1c1da8aed689.ExtensionsRequestBuilder) {
    return i7d3117e3afb0fbd44802c1ad6d61175bbd7b795db5acdd75c10d1c1da8aed689.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*icc005ad8da2b3b97a3de6350a707a53f032a558960b40cd2afcc1c28e246ada6.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return icc005ad8da2b3b97a3de6350a707a53f032a558960b40cd2afcc1c28e246ada6.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i83473d2a0bc99449ff6981736c3c17fbee0979720e884cddb635db5c58ed629c.ForwardRequestBuilder) {
    return i83473d2a0bc99449ff6981736c3c17fbee0979720e884cddb635db5c58ed629c.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i2ffebe5209de4076e8a82adab1a4bfca6c05cdf6f208c455125b48c57db65904.MultiValueExtendedPropertiesRequestBuilder) {
    return i2ffebe5209de4076e8a82adab1a4bfca6c05cdf6f208c455125b48c57db65904.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ifcf0b211e141030ced11d35ff3b92d777831766af0ce38e3192dd97c536b3221.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ifcf0b211e141030ced11d35ff3b92d777831766af0ce38e3192dd97c536b3221.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i48b536acd4f01ac6a9b0c428b9de43a1cbde1b708b519c38213a57d87b96504f.SingleValueExtendedPropertiesRequestBuilder) {
    return i48b536acd4f01ac6a9b0c428b9de43a1cbde1b708b519c38213a57d87b96504f.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i0a8a1ba6dbe110e4b0d2005064e6c047ae4063f6e259f00465d057bc61bf5204.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i0a8a1ba6dbe110e4b0d2005064e6c047ae4063f6e259f00465d057bc61bf5204.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ic2b75b590eb1afa68b5b139fda2fa1f65d5573ce52b3716a67f1e95fa70569e8.SnoozeReminderRequestBuilder) {
    return ic2b75b590eb1afa68b5b139fda2fa1f65d5573ce52b3716a67f1e95fa70569e8.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*idcd892355bb9b0a3694ef802b374e6be28dd95de1898988e0f5500a51feedf6f.TentativelyAcceptRequestBuilder) {
    return idcd892355bb9b0a3694ef802b374e6be28dd95de1898988e0f5500a51feedf6f.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
