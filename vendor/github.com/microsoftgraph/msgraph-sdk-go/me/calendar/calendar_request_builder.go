package calendar

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i321470654b3ef34bacc414654086f423258f1193a26b214185a5cc558e2bb5e2 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview"
    i6af69cbc9d5ee4e56d0394b5dbb9ba47c157c0ee557e79646739cd31671483d3 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/multivalueextendedproperties"
    i8106511447607356fcb7172494684f7bb195617a2f3cfb8090371e4dca9494cf "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/singlevalueextendedproperties"
    ia93089358fcf18f5dfb97b53182f1524884b6a134ef09311582420c51b266f63 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/allowedcalendarsharingroleswithuser"
    ib746c7bd792e0607afacbee87215b032700dfb09c38134acd3e5367b09bd1eef "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/getschedule"
    ic3f5d8e3fd6959673a012943934769f057abadba324812644982eeef5f407def "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events"
    iefe8e12aae13619f806b3f0379ff65f77fd7287cd264dc16f120a83b59db3999 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarpermissions"
    i209c23501c843b5dd2a14fe5d4f0a449bc55aac8f20ad31a9d86e7cb75577e28 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/multivalueextendedproperties/item"
    i79296554e9c2c5b59da0a63741ace146ef92bbb0bc577e1866c1d097111c377c "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/singlevalueextendedproperties/item"
    i96d50ecda61b15c98effcc511b1eff207b365d1e3b05ef2f704013788d04ccae "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/events/item"
    i9e50aa45f8b8ff846d0226fbc640efa378ef39164e74364f6e3b11c0735dfc91 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarview/item"
    if33c45189b52ce3c1a36912a56beedbaacd4a8273d6116ef5f83bd582e5a8506 "github.com/microsoftgraph/msgraph-sdk-go/me/calendar/calendarpermissions/item"
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
func (m *CalendarRequestBuilder) AllowedCalendarSharingRolesWithUser(user *string)(*ia93089358fcf18f5dfb97b53182f1524884b6a134ef09311582420c51b266f63.AllowedCalendarSharingRolesWithUserRequestBuilder) {
    return ia93089358fcf18f5dfb97b53182f1524884b6a134ef09311582420c51b266f63.NewAllowedCalendarSharingRolesWithUserRequestBuilderInternal(m.pathParameters, m.requestAdapter, user);
}
// CalendarPermissions provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissions()(*iefe8e12aae13619f806b3f0379ff65f77fd7287cd264dc16f120a83b59db3999.CalendarPermissionsRequestBuilder) {
    return iefe8e12aae13619f806b3f0379ff65f77fd7287cd264dc16f120a83b59db3999.NewCalendarPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarPermissionsById provides operations to manage the calendarPermissions property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarPermissionsById(id string)(*if33c45189b52ce3c1a36912a56beedbaacd4a8273d6116ef5f83bd582e5a8506.CalendarPermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["calendarPermission%2Did"] = id
    }
    return if33c45189b52ce3c1a36912a56beedbaacd4a8273d6116ef5f83bd582e5a8506.NewCalendarPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CalendarView provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarView()(*i321470654b3ef34bacc414654086f423258f1193a26b214185a5cc558e2bb5e2.CalendarViewRequestBuilder) {
    return i321470654b3ef34bacc414654086f423258f1193a26b214185a5cc558e2bb5e2.NewCalendarViewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarViewById provides operations to manage the calendarView property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) CalendarViewById(id string)(*i9e50aa45f8b8ff846d0226fbc640efa378ef39164e74364f6e3b11c0735dfc91.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return i9e50aa45f8b8ff846d0226fbc640efa378ef39164e74364f6e3b11c0735dfc91.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewCalendarRequestBuilderInternal instantiates a new CalendarRequestBuilder and sets the default values.
func NewCalendarRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CalendarRequestBuilder) {
    m := &CalendarRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/calendar{?%24select}";
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
func (m *CalendarRequestBuilder) Events()(*ic3f5d8e3fd6959673a012943934769f057abadba324812644982eeef5f407def.EventsRequestBuilder) {
    return ic3f5d8e3fd6959673a012943934769f057abadba324812644982eeef5f407def.NewEventsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// EventsById provides operations to manage the events property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) EventsById(id string)(*i96d50ecda61b15c98effcc511b1eff207b365d1e3b05ef2f704013788d04ccae.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return i96d50ecda61b15c98effcc511b1eff207b365d1e3b05ef2f704013788d04ccae.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *CalendarRequestBuilder) GetSchedule()(*ib746c7bd792e0607afacbee87215b032700dfb09c38134acd3e5367b09bd1eef.GetScheduleRequestBuilder) {
    return ib746c7bd792e0607afacbee87215b032700dfb09c38134acd3e5367b09bd1eef.NewGetScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedProperties()(*i6af69cbc9d5ee4e56d0394b5dbb9ba47c157c0ee557e79646739cd31671483d3.MultiValueExtendedPropertiesRequestBuilder) {
    return i6af69cbc9d5ee4e56d0394b5dbb9ba47c157c0ee557e79646739cd31671483d3.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i209c23501c843b5dd2a14fe5d4f0a449bc55aac8f20ad31a9d86e7cb75577e28.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i209c23501c843b5dd2a14fe5d4f0a449bc55aac8f20ad31a9d86e7cb75577e28.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *CalendarRequestBuilder) SingleValueExtendedProperties()(*i8106511447607356fcb7172494684f7bb195617a2f3cfb8090371e4dca9494cf.SingleValueExtendedPropertiesRequestBuilder) {
    return i8106511447607356fcb7172494684f7bb195617a2f3cfb8090371e4dca9494cf.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.calendar entity.
func (m *CalendarRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i79296554e9c2c5b59da0a63741ace146ef92bbb0bc577e1866c1d097111c377c.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i79296554e9c2c5b59da0a63741ace146ef92bbb0bc577e1866c1d097111c377c.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
