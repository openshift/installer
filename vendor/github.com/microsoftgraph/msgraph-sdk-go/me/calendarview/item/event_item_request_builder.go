package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1d0689ad330d2744c1fb18da8874f5ab0d3d9f9e089b981c1847e0732b28d4b6 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/tentativelyaccept"
    i465e55e32341d082edd97b4cbbe30b74404cd0753d5a94c441f6e442632d344d "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/cancel"
    i53f43a3531f80c6e0571cd6533e13dfb0a9bec5d78a4451ab57d06d69569e313 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/decline"
    i68fc9e57c0f6013ec42006755700065554fc5bf5e29fcd5025a538b57570720b "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/singlevalueextendedproperties"
    i6a090f392ca1e1ce7d5306cbfcc6ecedfebdc0920501b06deec7800a7c0a2a4d "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/extensions"
    i7ebbe5973f1cfdd153c7211a7d729e6acf754214d73692548caaa452fc437931 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/dismissreminder"
    i890bc6825225252197ba0329e5325181337dcc5018722ea3f24133acc97d663d "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/attachments"
    ib37f35045d733a190b9416b6893f854438f6a18b0b560d7fd6225d7b528fc019 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances"
    ib4cb01573cfaac49f1987ed83b0866022e257847a6e6a595526d0a91e42cd048 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/snoozereminder"
    ibd0b52fb44db2ad950860bec1e2b4ed68ce5d6349861ec6e194d282bde7bcc51 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/forward"
    ibe852a23e0398849c98cf64aa13ee2a52ee98610fea6b7e377c6a7dd1c1628c0 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/multivalueextendedproperties"
    ie12260a766bb83c9b88880512999d5f33b231a5677bf79c4d07097242d01255a "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/calendar"
    ie711528247ed53f1948a3ebd7b38785aa216df02cdea65cc2a13d9fb60e90e54 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/accept"
    i624e81c2db0b3bfef96d02d14ccc5aafca1e94da30a6ef51f34764348903e10e "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/attachments/item"
    i6d746bffa422f6cab9bce504d82f353fa1e96f30d9135162a2867c0b6ee56409 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/extensions/item"
    i981aed02b0a2dbdb188fef5063ef9909c9da485cac50836f2c3d65d618577a23 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/singlevalueextendedproperties/item"
    ic2fee68681a1462f085ba98b6ba6395da8df913e922d9bff315377a4a254da2e "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/instances/item"
    icbdb2876a9fc3cde1981cefae88c00544c63c09c8176a43b5750b723fee1cfa3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendarview/item/multivalueextendedproperties/item"
)

// EventItemRequestBuilder provides operations to manage the calendarView property of the microsoft.graph.user entity.
type EventItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EventItemRequestBuilderGetQueryParameters the calendar view for the calendar. Read-only. Nullable.
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
func (m *EventItemRequestBuilder) Accept()(*ie711528247ed53f1948a3ebd7b38785aa216df02cdea65cc2a13d9fb60e90e54.AcceptRequestBuilder) {
    return ie711528247ed53f1948a3ebd7b38785aa216df02cdea65cc2a13d9fb60e90e54.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i890bc6825225252197ba0329e5325181337dcc5018722ea3f24133acc97d663d.AttachmentsRequestBuilder) {
    return i890bc6825225252197ba0329e5325181337dcc5018722ea3f24133acc97d663d.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*i624e81c2db0b3bfef96d02d14ccc5aafca1e94da30a6ef51f34764348903e10e.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i624e81c2db0b3bfef96d02d14ccc5aafca1e94da30a6ef51f34764348903e10e.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*ie12260a766bb83c9b88880512999d5f33b231a5677bf79c4d07097242d01255a.CalendarRequestBuilder) {
    return ie12260a766bb83c9b88880512999d5f33b231a5677bf79c4d07097242d01255a.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i465e55e32341d082edd97b4cbbe30b74404cd0753d5a94c441f6e442632d344d.CancelRequestBuilder) {
    return i465e55e32341d082edd97b4cbbe30b74404cd0753d5a94c441f6e442632d344d.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendarView/{event%2Did}{?startDateTime*,endDateTime*,%24select}";
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
// CreateGetRequestInformation the calendar view for the calendar. Read-only. Nullable.
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
func (m *EventItemRequestBuilder) Decline()(*i53f43a3531f80c6e0571cd6533e13dfb0a9bec5d78a4451ab57d06d69569e313.DeclineRequestBuilder) {
    return i53f43a3531f80c6e0571cd6533e13dfb0a9bec5d78a4451ab57d06d69569e313.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*i7ebbe5973f1cfdd153c7211a7d729e6acf754214d73692548caaa452fc437931.DismissReminderRequestBuilder) {
    return i7ebbe5973f1cfdd153c7211a7d729e6acf754214d73692548caaa452fc437931.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i6a090f392ca1e1ce7d5306cbfcc6ecedfebdc0920501b06deec7800a7c0a2a4d.ExtensionsRequestBuilder) {
    return i6a090f392ca1e1ce7d5306cbfcc6ecedfebdc0920501b06deec7800a7c0a2a4d.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i6d746bffa422f6cab9bce504d82f353fa1e96f30d9135162a2867c0b6ee56409.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i6d746bffa422f6cab9bce504d82f353fa1e96f30d9135162a2867c0b6ee56409.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*ibd0b52fb44db2ad950860bec1e2b4ed68ce5d6349861ec6e194d282bde7bcc51.ForwardRequestBuilder) {
    return ibd0b52fb44db2ad950860bec1e2b4ed68ce5d6349861ec6e194d282bde7bcc51.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the calendar view for the calendar. Read-only. Nullable.
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
func (m *EventItemRequestBuilder) Instances()(*ib37f35045d733a190b9416b6893f854438f6a18b0b560d7fd6225d7b528fc019.InstancesRequestBuilder) {
    return ib37f35045d733a190b9416b6893f854438f6a18b0b560d7fd6225d7b528fc019.NewInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstancesById provides operations to manage the instances property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) InstancesById(id string)(*ic2fee68681a1462f085ba98b6ba6395da8df913e922d9bff315377a4a254da2e.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did1"] = id
    }
    return ic2fee68681a1462f085ba98b6ba6395da8df913e922d9bff315377a4a254da2e.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*ibe852a23e0398849c98cf64aa13ee2a52ee98610fea6b7e377c6a7dd1c1628c0.MultiValueExtendedPropertiesRequestBuilder) {
    return ibe852a23e0398849c98cf64aa13ee2a52ee98610fea6b7e377c6a7dd1c1628c0.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*icbdb2876a9fc3cde1981cefae88c00544c63c09c8176a43b5750b723fee1cfa3.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return icbdb2876a9fc3cde1981cefae88c00544c63c09c8176a43b5750b723fee1cfa3.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i68fc9e57c0f6013ec42006755700065554fc5bf5e29fcd5025a538b57570720b.SingleValueExtendedPropertiesRequestBuilder) {
    return i68fc9e57c0f6013ec42006755700065554fc5bf5e29fcd5025a538b57570720b.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i981aed02b0a2dbdb188fef5063ef9909c9da485cac50836f2c3d65d618577a23.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i981aed02b0a2dbdb188fef5063ef9909c9da485cac50836f2c3d65d618577a23.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*ib4cb01573cfaac49f1987ed83b0866022e257847a6e6a595526d0a91e42cd048.SnoozeReminderRequestBuilder) {
    return ib4cb01573cfaac49f1987ed83b0866022e257847a6e6a595526d0a91e42cd048.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*i1d0689ad330d2744c1fb18da8874f5ab0d3d9f9e089b981c1847e0732b28d4b6.TentativelyAcceptRequestBuilder) {
    return i1d0689ad330d2744c1fb18da8874f5ab0d3d9f9e089b981c1847e0732b28d4b6.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
