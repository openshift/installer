package calendar

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i325e8c5a186a2ac5ba63357864499a2daeca2d6b1fe397a8e242931a81e90fd6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/singlevalueextendedproperties"
    i564e44237d1fd29264e29294fe61a277a2b68966ba3da977d4dc4188473efc6f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/allowedcalendarsharingroleswithuser"
    i59567e44c84457909a7acd2691855472309ae72d46905ad5b927e32ac1bc6775 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/calendarpermissions"
    i8a7b3ddf32454c04e11c42fafddbe8a65a1b8c29699e948e42fb5fb489efff7a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/getschedule"
    i91f19a09d718f883e167c3c1fea6e1dc0bafbd34f8ca677f95c72d7745bcd847 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/multivalueextendedproperties"
    ia4056a24d617675f13dc550d9eb8d60439d535c22994554b767157caae834324 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/calendarview"
    iae525d555fa999780edca2ea54ea681baf7b616326144106b66e173f4930b985 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events"
    i066e4b225a90f837574cba0d06582cb5234a89a58269aad5ccca1e0d592385ed "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/calendarview/item"
    i32f1c0e0a1af3233116ddc476f2459674db4e63222c4a5762ca7049d50e06a83 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/singlevalueextendedproperties/item"
    i84f844d52d5e497fca0d718f40dd524efb38a3a493f27e5b266d0f0b0d06bd58 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/calendarpermissions/item"
    ieb71583a78a26be4af7a9008915150796bd995e96567f7d812206ea5c47e34ee "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/events/item"
    iff61068eadc42d2652152ea70da0c96f72f2b3e78201eeb2e58a48230426fa5a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/calendar/multivalueextendedproperties/item"
)

// CalendarRequestBuilder provides operations to manage the calendar property of the microsoft.graph.group entity.
type CalendarRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// CalendarRequestBuilderGetQueryParameters the group's calendar. Read-only.
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
// AllowedCalendarSharingRolesWithUser provides operations to call the allowedCalendarSharingRoles method.
func (m *CalendarRequestBuilder) AllowedCalendarSharingRolesWithUser(user *string)(*i564e44237d1fd29264e29294fe61a277a2b68966ba3da977d4dc4188473efc6f.AllowedCalendarSharingRolesWithUserRequestBuilder) {
    return i564e44237d1fd29264e29294fe61a277a2b68966ba3da977d4dc4188473efc6f.NewAllowedCalendarSharingRolesWithUserRequestBuilderInternal(m.pathParameters, m.requestAdapter, user);
}
// CalendarPermissions provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissions()(*i59567e44c84457909a7acd2691855472309ae72d46905ad5b927e32ac1bc6775.CalendarPermissionsRequestBuilder) {
    return i59567e44c84457909a7acd2691855472309ae72d46905ad5b927e32ac1bc6775.NewCalendarPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarPermissionsById provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissionsById(id string)(*i84f844d52d5e497fca0d718f40dd524efb38a3a493f27e5b266d0f0b0d06bd58.CalendarPermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["calendarPermission%2Did"] = id
    }
    return i84f844d52d5e497fca0d718f40dd524efb38a3a493f27e5b266d0f0b0d06bd58.NewCalendarPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CalendarView provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarView()(*ia4056a24d617675f13dc550d9eb8d60439d535c22994554b767157caae834324.CalendarViewRequestBuilder) {
    return ia4056a24d617675f13dc550d9eb8d60439d535c22994554b767157caae834324.NewCalendarViewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarViewById provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarViewById(id string)(*i066e4b225a90f837574cba0d06582cb5234a89a58269aad5ccca1e0d592385ed.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return i066e4b225a90f837574cba0d06582cb5234a89a58269aad5ccca1e0d592385ed.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewCalendarRequestBuilderInternal instantiates a new CalendarRequestBuilder and sets the default values.
func NewCalendarRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CalendarRequestBuilder) {
    m := &CalendarRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/calendar{?%24select}";
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
// CreateGetRequestInformation the group's calendar. Read-only.
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
// Events provides operations to manage the events property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) Events()(*iae525d555fa999780edca2ea54ea681baf7b616326144106b66e173f4930b985.EventsRequestBuilder) {
    return iae525d555fa999780edca2ea54ea681baf7b616326144106b66e173f4930b985.NewEventsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// EventsById provides operations to manage the events property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) EventsById(id string)(*ieb71583a78a26be4af7a9008915150796bd995e96567f7d812206ea5c47e34ee.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return ieb71583a78a26be4af7a9008915150796bd995e96567f7d812206ea5c47e34ee.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get the group's calendar. Read-only.
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
func (m *CalendarRequestBuilder) GetSchedule()(*i8a7b3ddf32454c04e11c42fafddbe8a65a1b8c29699e948e42fb5fb489efff7a.GetScheduleRequestBuilder) {
    return i8a7b3ddf32454c04e11c42fafddbe8a65a1b8c29699e948e42fb5fb489efff7a.NewGetScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedProperties()(*i91f19a09d718f883e167c3c1fea6e1dc0bafbd34f8ca677f95c72d7745bcd847.MultiValueExtendedPropertiesRequestBuilder) {
    return i91f19a09d718f883e167c3c1fea6e1dc0bafbd34f8ca677f95c72d7745bcd847.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedPropertiesById(id string)(*iff61068eadc42d2652152ea70da0c96f72f2b3e78201eeb2e58a48230426fa5a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return iff61068eadc42d2652152ea70da0c96f72f2b3e78201eeb2e58a48230426fa5a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) SingleValueExtendedProperties()(*i325e8c5a186a2ac5ba63357864499a2daeca2d6b1fe397a8e242931a81e90fd6.SingleValueExtendedPropertiesRequestBuilder) {
    return i325e8c5a186a2ac5ba63357864499a2daeca2d6b1fe397a8e242931a81e90fd6.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i32f1c0e0a1af3233116ddc476f2459674db4e63222c4a5762ca7049d50e06a83.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i32f1c0e0a1af3233116ddc476f2459674db4e63222c4a5762ca7049d50e06a83.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
