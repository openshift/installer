package calendar

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i07b8994e6a1fa340814b90f5b658702fed4412e7ee820e23e3e990f4afd9eb73 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/getschedule"
    i520ec78d0b8b7722563a2934d91014ff877091dfa7315c9686d8905e132ae5b3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/multivalueextendedproperties"
    i94b8da2154e4a088b032a31d42d1f7f3035dc1e0cc99bff3d7db6d04b78bcaf9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/allowedcalendarsharingroleswithuser"
    ib399bd7a52c229f4262ed1659c9cfbc736f51e55f5b34d90bd5973e1cd56e052 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/singlevalueextendedproperties"
    ic4fbedba8b71c96a8c1cc23f5f68b98bdb74a545c81a35239a1ae6a4bddc6749 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarpermissions"
    idffc983be1e870a2a4117185ece24daafd4929a0264366c338f3760c1423b076 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events"
    if44e7d37c091c522009df8d2b8b30d67287c1346420ed710b1f9997ffa334461 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview"
    i2133c4e03374cde57e17f2aa66857c6ba629b1bbfa7b89e0585f17c88b3470d0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarpermissions/item"
    i5db0ccf3677a8d13461e0b10ef2980107c6cb5223fecf31a0484b6d4eac391d7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/multivalueextendedproperties/item"
    i631cbd9bc2cfc12ac467e801a8ff1de1f272c89e5691913fac02728b1daca8e7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/singlevalueextendedproperties/item"
    i7da62e1f0e8b5076b925f58a24cd8bcdf2c18c876c5e3fac9e0cce1708c8d669 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/calendarview/item"
    iad6ec6cb8dcd07c8cb3bd03b2bc97831694f469367bb8b0cf70c13010c71fe91 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events/item"
)

// CalendarRequestBuilder provides operations to manage the calendar property of the microsoft.graph.user entity.
type CalendarRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// CalendarRequestBuilderGetQueryParameters get the properties and relationships of a calendar object. The calendar can be one for a user, or the default calendar of a Microsoft 365 group. There are two scenarios where an app can get another user's calendar:
type CalendarRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CalendarRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CalendarRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CalendarRequestBuilderGetQueryParameters
}
// CalendarRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CalendarRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllowedCalendarSharingRolesWithUser provides operations to call the allowedCalendarSharingRoles method.
func (m *CalendarRequestBuilder) AllowedCalendarSharingRolesWithUser(user *string)(*i94b8da2154e4a088b032a31d42d1f7f3035dc1e0cc99bff3d7db6d04b78bcaf9.AllowedCalendarSharingRolesWithUserRequestBuilder) {
    return i94b8da2154e4a088b032a31d42d1f7f3035dc1e0cc99bff3d7db6d04b78bcaf9.NewAllowedCalendarSharingRolesWithUserRequestBuilderInternal(m.pathParameters, m.requestAdapter, user);
}
// CalendarPermissions provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissions()(*ic4fbedba8b71c96a8c1cc23f5f68b98bdb74a545c81a35239a1ae6a4bddc6749.CalendarPermissionsRequestBuilder) {
    return ic4fbedba8b71c96a8c1cc23f5f68b98bdb74a545c81a35239a1ae6a4bddc6749.NewCalendarPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarPermissionsById provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissionsById(id string)(*i2133c4e03374cde57e17f2aa66857c6ba629b1bbfa7b89e0585f17c88b3470d0.CalendarPermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["calendarPermission%2Did"] = id
    }
    return i2133c4e03374cde57e17f2aa66857c6ba629b1bbfa7b89e0585f17c88b3470d0.NewCalendarPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CalendarView provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarView()(*if44e7d37c091c522009df8d2b8b30d67287c1346420ed710b1f9997ffa334461.CalendarViewRequestBuilder) {
    return if44e7d37c091c522009df8d2b8b30d67287c1346420ed710b1f9997ffa334461.NewCalendarViewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarViewById provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarViewById(id string)(*i7da62e1f0e8b5076b925f58a24cd8bcdf2c18c876c5e3fac9e0cce1708c8d669.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return i7da62e1f0e8b5076b925f58a24cd8bcdf2c18c876c5e3fac9e0cce1708c8d669.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewCalendarRequestBuilderInternal instantiates a new CalendarRequestBuilder and sets the default values.
func NewCalendarRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CalendarRequestBuilder) {
    m := &CalendarRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/calendar{?%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewCalendarRequestBuilder instantiates a new CalendarRequestBuilder and sets the default values.
func NewCalendarRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CalendarRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCalendarRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get the properties and relationships of a calendar object. The calendar can be one for a user, or the default calendar of a Microsoft 365 group. There are two scenarios where an app can get another user's calendar:
func (m *CalendarRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *CalendarRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the properties of a calendar object. The calendar can be one for a user, or the default calendar of a Microsoft 365 group.
func (m *CalendarRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable, requestConfiguration *CalendarRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Events provides operations to manage the events property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) Events()(*idffc983be1e870a2a4117185ece24daafd4929a0264366c338f3760c1423b076.EventsRequestBuilder) {
    return idffc983be1e870a2a4117185ece24daafd4929a0264366c338f3760c1423b076.NewEventsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// EventsById provides operations to manage the events property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) EventsById(id string)(*iad6ec6cb8dcd07c8cb3bd03b2bc97831694f469367bb8b0cf70c13010c71fe91.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return iad6ec6cb8dcd07c8cb3bd03b2bc97831694f469367bb8b0cf70c13010c71fe91.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get get the properties and relationships of a calendar object. The calendar can be one for a user, or the default calendar of a Microsoft 365 group. There are two scenarios where an app can get another user's calendar:
func (m *CalendarRequestBuilder) Get(ctx context.Context, requestConfiguration *CalendarRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCalendarFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable), nil
}
// GetSchedule provides operations to call the getSchedule method.
func (m *CalendarRequestBuilder) GetSchedule()(*i07b8994e6a1fa340814b90f5b658702fed4412e7ee820e23e3e990f4afd9eb73.GetScheduleRequestBuilder) {
    return i07b8994e6a1fa340814b90f5b658702fed4412e7ee820e23e3e990f4afd9eb73.NewGetScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedProperties()(*i520ec78d0b8b7722563a2934d91014ff877091dfa7315c9686d8905e132ae5b3.MultiValueExtendedPropertiesRequestBuilder) {
    return i520ec78d0b8b7722563a2934d91014ff877091dfa7315c9686d8905e132ae5b3.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i5db0ccf3677a8d13461e0b10ef2980107c6cb5223fecf31a0484b6d4eac391d7.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i5db0ccf3677a8d13461e0b10ef2980107c6cb5223fecf31a0484b6d4eac391d7.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the properties of a calendar object. The calendar can be one for a user, or the default calendar of a Microsoft 365 group.
func (m *CalendarRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable, requestConfiguration *CalendarRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCalendarFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Calendarable), nil
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) SingleValueExtendedProperties()(*ib399bd7a52c229f4262ed1659c9cfbc736f51e55f5b34d90bd5973e1cd56e052.SingleValueExtendedPropertiesRequestBuilder) {
    return ib399bd7a52c229f4262ed1659c9cfbc736f51e55f5b34d90bd5973e1cd56e052.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i631cbd9bc2cfc12ac467e801a8ff1de1f272c89e5691913fac02728b1daca8e7.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i631cbd9bc2cfc12ac467e801a8ff1de1f272c89e5691913fac02728b1daca8e7.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
