package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0608c7e63daed1d8c25dae69476d30ed0865692d2083143f0218ba545bd1e6fe "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/forward"
    i07d0e4dc59220db4d7a12d39275288711f40897a0177e9303992283ab77defff "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/singlevalueextendedproperties"
    i600ccf8ecb63c2306bea2313987324f661d6a99b86ce8c4d346b598ccb947723 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/accept"
    i6ea1b94176d49fe77dc1401acdc47611ad8a52c9c6679fb9a0f4042e4a8b1e2b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/attachments"
    i6f6a379fccd6d7ecab5d208cda6abdeea04b51f72052df61fbd909540c1d7785 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/calendar"
    i80f3a65475ab0b52025d7da02f157b571a71ec1653a35767a6445d3396413731 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/cancel"
    i88a6407fe07fdde4b4f64edd95fdd05984b8f248ff13190c22eecee6198b3935 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/extensions"
    i90956c2f657489b41c77cca666d102dea0bdf9d73b17f4a49c73f6e2ed3271e9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/decline"
    i98290a55afd1d1f0fd9f466411a505e295f1873fb97266f9cc506b0d4d6917ac "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/multivalueextendedproperties"
    ia39929ce6cbf360020596a4720c2bdc855842e132ab54471ed8c11ddcf661102 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/dismissreminder"
    iac5a3fc9fdba162b0d168cefc298dc8c085d69c261727fb186025595df5972f8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/tentativelyaccept"
    if29ba39feb86bf644a6ce867550b8dc648df46f2908d6b2ac6dcd043bc00fb6e "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/snoozereminder"
    i0ac354fe96550f13333780b4988ce2c8fbb97ccb242f8a0b579a9d16173fcce3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/extensions/item"
    i2ca41bd2f6db760c9dda4d70fef3da36f4fff8e684fcc1a720ef69b56eb396cf "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/singlevalueextendedproperties/item"
    i6aa496b042042eb3121230fe915280bf61f9463e3012cdaddc86bdfc32da2186 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/multivalueextendedproperties/item"
    ib7e466ee049923144d5c5bd395e80154c60e05395739ccdb4e370be08c2b4127 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item/instances/item/attachments/item"
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
func (m *EventItemRequestBuilder) Accept()(*i600ccf8ecb63c2306bea2313987324f661d6a99b86ce8c4d346b598ccb947723.AcceptRequestBuilder) {
    return i600ccf8ecb63c2306bea2313987324f661d6a99b86ce8c4d346b598ccb947723.NewAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Attachments()(*i6ea1b94176d49fe77dc1401acdc47611ad8a52c9c6679fb9a0f4042e4a8b1e2b.AttachmentsRequestBuilder) {
    return i6ea1b94176d49fe77dc1401acdc47611ad8a52c9c6679fb9a0f4042e4a8b1e2b.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) AttachmentsById(id string)(*ib7e466ee049923144d5c5bd395e80154c60e05395739ccdb4e370be08c2b4127.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ib7e466ee049923144d5c5bd395e80154c60e05395739ccdb4e370be08c2b4127.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Calendar()(*i6f6a379fccd6d7ecab5d208cda6abdeea04b51f72052df61fbd909540c1d7785.CalendarRequestBuilder) {
    return i6f6a379fccd6d7ecab5d208cda6abdeea04b51f72052df61fbd909540c1d7785.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cancel provides operations to call the cancel method.
func (m *EventItemRequestBuilder) Cancel()(*i80f3a65475ab0b52025d7da02f157b571a71ec1653a35767a6445d3396413731.CancelRequestBuilder) {
    return i80f3a65475ab0b52025d7da02f157b571a71ec1653a35767a6445d3396413731.NewCancelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEventItemRequestBuilderInternal instantiates a new EventItemRequestBuilder and sets the default values.
func NewEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EventItemRequestBuilder) {
    m := &EventItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendarView/{event%2Did}/instances/{event%2Did1}{?%24select}";
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
func (m *EventItemRequestBuilder) Decline()(*i90956c2f657489b41c77cca666d102dea0bdf9d73b17f4a49c73f6e2ed3271e9.DeclineRequestBuilder) {
    return i90956c2f657489b41c77cca666d102dea0bdf9d73b17f4a49c73f6e2ed3271e9.NewDeclineRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DismissReminder provides operations to call the dismissReminder method.
func (m *EventItemRequestBuilder) DismissReminder()(*ia39929ce6cbf360020596a4720c2bdc855842e132ab54471ed8c11ddcf661102.DismissReminderRequestBuilder) {
    return ia39929ce6cbf360020596a4720c2bdc855842e132ab54471ed8c11ddcf661102.NewDismissReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) Extensions()(*i88a6407fe07fdde4b4f64edd95fdd05984b8f248ff13190c22eecee6198b3935.ExtensionsRequestBuilder) {
    return i88a6407fe07fdde4b4f64edd95fdd05984b8f248ff13190c22eecee6198b3935.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) ExtensionsById(id string)(*i0ac354fe96550f13333780b4988ce2c8fbb97ccb242f8a0b579a9d16173fcce3.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i0ac354fe96550f13333780b4988ce2c8fbb97ccb242f8a0b579a9d16173fcce3.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *EventItemRequestBuilder) Forward()(*i0608c7e63daed1d8c25dae69476d30ed0865692d2083143f0218ba545bd1e6fe.ForwardRequestBuilder) {
    return i0608c7e63daed1d8c25dae69476d30ed0865692d2083143f0218ba545bd1e6fe.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *EventItemRequestBuilder) MultiValueExtendedProperties()(*i98290a55afd1d1f0fd9f466411a505e295f1873fb97266f9cc506b0d4d6917ac.MultiValueExtendedPropertiesRequestBuilder) {
    return i98290a55afd1d1f0fd9f466411a505e295f1873fb97266f9cc506b0d4d6917ac.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i6aa496b042042eb3121230fe915280bf61f9463e3012cdaddc86bdfc32da2186.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i6aa496b042042eb3121230fe915280bf61f9463e3012cdaddc86bdfc32da2186.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedProperties()(*i07d0e4dc59220db4d7a12d39275288711f40897a0177e9303992283ab77defff.SingleValueExtendedPropertiesRequestBuilder) {
    return i07d0e4dc59220db4d7a12d39275288711f40897a0177e9303992283ab77defff.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.event entity.
func (m *EventItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i2ca41bd2f6db760c9dda4d70fef3da36f4fff8e684fcc1a720ef69b56eb396cf.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i2ca41bd2f6db760c9dda4d70fef3da36f4fff8e684fcc1a720ef69b56eb396cf.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SnoozeReminder provides operations to call the snoozeReminder method.
func (m *EventItemRequestBuilder) SnoozeReminder()(*if29ba39feb86bf644a6ce867550b8dc648df46f2908d6b2ac6dcd043bc00fb6e.SnoozeReminderRequestBuilder) {
    return if29ba39feb86bf644a6ce867550b8dc648df46f2908d6b2ac6dcd043bc00fb6e.NewSnoozeReminderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TentativelyAccept provides operations to call the tentativelyAccept method.
func (m *EventItemRequestBuilder) TentativelyAccept()(*iac5a3fc9fdba162b0d168cefc298dc8c085d69c261727fb186025595df5972f8.TentativelyAcceptRequestBuilder) {
    return iac5a3fc9fdba162b0d168cefc298dc8c085d69c261727fb186025595df5972f8.NewTentativelyAcceptRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
