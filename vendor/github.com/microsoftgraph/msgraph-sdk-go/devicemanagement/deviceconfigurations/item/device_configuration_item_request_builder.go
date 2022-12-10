package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1acaef5ee6ae1997bcbb07d275900c5d95eb14266a9c79c9ee811bf4c0edefb4 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/assignments"
    i1cbc049d4f1b0a0c7493469fdf2a72ede1929262173246b80d83717fb3509836 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/userstatusoverview"
    i2052556f4fb1ecc2b8520d93accde420e5a85783ea309802622a14a2965cac8c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/devicestatuses"
    i419f93c2382cc7ee59ef6aee59654c413d9feb258f5f656a8f44dc8fb4ad2b61 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/devicestatusoverview"
    i4cd13eb5be3949210a98b3e74a48e803b388b1057dc5b76fa296250220177115 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/assign"
    i5d3bb0910f46428b82e3ef51c35a9f795351e7ee61bfb08c89e5de44026ecaa2 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/getomasettingplaintextvaluewithsecretreferencevalueid"
    ib0c9f5570d692d204085c3256ffa0914246699c56424649fc82be68680943fb2 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/devicesettingstatesummaries"
    ie37508b74757cd3bd2b6dc98d9db2ba9651032937fe4385be4b4c611ade0f8fc "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/userstatuses"
    i13d5635acea3d1befd38097b223c984a262de351ddcc99002d56e0362ecf8c7c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/userstatuses/item"
    i37d5c64fb898f2ba64f90d5a72b1da2a1813227f63002283be7b64bced3b78de "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/devicesettingstatesummaries/item"
    i847fd9044721470530d680766ffdf92095b6c54818a689f742e81f854e71cbd3 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/assignments/item"
    ie1759fa3aebfd08bd0a0973770111e293804186cb2bccc14cfdf3f11a34c8f75 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/deviceconfigurations/item/devicestatuses/item"
)

// DeviceConfigurationItemRequestBuilder provides operations to manage the deviceConfigurations property of the microsoft.graph.deviceManagement entity.
type DeviceConfigurationItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DeviceConfigurationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceConfigurationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DeviceConfigurationItemRequestBuilderGetQueryParameters the device configurations.
type DeviceConfigurationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DeviceConfigurationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceConfigurationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeviceConfigurationItemRequestBuilderGetQueryParameters
}
// DeviceConfigurationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceConfigurationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Assign provides operations to call the assign method.
func (m *DeviceConfigurationItemRequestBuilder) Assign()(*i4cd13eb5be3949210a98b3e74a48e803b388b1057dc5b76fa296250220177115.AssignRequestBuilder) {
    return i4cd13eb5be3949210a98b3e74a48e803b388b1057dc5b76fa296250220177115.NewAssignRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Assignments provides operations to manage the assignments property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) Assignments()(*i1acaef5ee6ae1997bcbb07d275900c5d95eb14266a9c79c9ee811bf4c0edefb4.AssignmentsRequestBuilder) {
    return i1acaef5ee6ae1997bcbb07d275900c5d95eb14266a9c79c9ee811bf4c0edefb4.NewAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentsById provides operations to manage the assignments property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) AssignmentsById(id string)(*i847fd9044721470530d680766ffdf92095b6c54818a689f742e81f854e71cbd3.DeviceConfigurationAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationAssignment%2Did"] = id
    }
    return i847fd9044721470530d680766ffdf92095b6c54818a689f742e81f854e71cbd3.NewDeviceConfigurationAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDeviceConfigurationItemRequestBuilderInternal instantiates a new DeviceConfigurationItemRequestBuilder and sets the default values.
func NewDeviceConfigurationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceConfigurationItemRequestBuilder) {
    m := &DeviceConfigurationItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/deviceManagement/deviceConfigurations/{deviceConfiguration%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDeviceConfigurationItemRequestBuilder instantiates a new DeviceConfigurationItemRequestBuilder and sets the default values.
func NewDeviceConfigurationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceConfigurationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeviceConfigurationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property deviceConfigurations for deviceManagement
func (m *DeviceConfigurationItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *DeviceConfigurationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the device configurations.
func (m *DeviceConfigurationItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DeviceConfigurationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property deviceConfigurations in deviceManagement
func (m *DeviceConfigurationItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable, requestConfiguration *DeviceConfigurationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property deviceConfigurations for deviceManagement
func (m *DeviceConfigurationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DeviceConfigurationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DeviceSettingStateSummaries provides operations to manage the deviceSettingStateSummaries property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) DeviceSettingStateSummaries()(*ib0c9f5570d692d204085c3256ffa0914246699c56424649fc82be68680943fb2.DeviceSettingStateSummariesRequestBuilder) {
    return ib0c9f5570d692d204085c3256ffa0914246699c56424649fc82be68680943fb2.NewDeviceSettingStateSummariesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceSettingStateSummariesById provides operations to manage the deviceSettingStateSummaries property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) DeviceSettingStateSummariesById(id string)(*i37d5c64fb898f2ba64f90d5a72b1da2a1813227f63002283be7b64bced3b78de.SettingStateDeviceSummaryItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["settingStateDeviceSummary%2Did"] = id
    }
    return i37d5c64fb898f2ba64f90d5a72b1da2a1813227f63002283be7b64bced3b78de.NewSettingStateDeviceSummaryItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceStatuses provides operations to manage the deviceStatuses property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) DeviceStatuses()(*i2052556f4fb1ecc2b8520d93accde420e5a85783ea309802622a14a2965cac8c.DeviceStatusesRequestBuilder) {
    return i2052556f4fb1ecc2b8520d93accde420e5a85783ea309802622a14a2965cac8c.NewDeviceStatusesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceStatusesById provides operations to manage the deviceStatuses property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) DeviceStatusesById(id string)(*ie1759fa3aebfd08bd0a0973770111e293804186cb2bccc14cfdf3f11a34c8f75.DeviceConfigurationDeviceStatusItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationDeviceStatus%2Did"] = id
    }
    return ie1759fa3aebfd08bd0a0973770111e293804186cb2bccc14cfdf3f11a34c8f75.NewDeviceConfigurationDeviceStatusItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceStatusOverview provides operations to manage the deviceStatusOverview property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) DeviceStatusOverview()(*i419f93c2382cc7ee59ef6aee59654c413d9feb258f5f656a8f44dc8fb4ad2b61.DeviceStatusOverviewRequestBuilder) {
    return i419f93c2382cc7ee59ef6aee59654c413d9feb258f5f656a8f44dc8fb4ad2b61.NewDeviceStatusOverviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the device configurations.
func (m *DeviceConfigurationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DeviceConfigurationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceConfigurationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable), nil
}
// GetOmaSettingPlainTextValueWithSecretReferenceValueId provides operations to call the getOmaSettingPlainTextValue method.
func (m *DeviceConfigurationItemRequestBuilder) GetOmaSettingPlainTextValueWithSecretReferenceValueId(secretReferenceValueId *string)(*i5d3bb0910f46428b82e3ef51c35a9f795351e7ee61bfb08c89e5de44026ecaa2.GetOmaSettingPlainTextValueWithSecretReferenceValueIdRequestBuilder) {
    return i5d3bb0910f46428b82e3ef51c35a9f795351e7ee61bfb08c89e5de44026ecaa2.NewGetOmaSettingPlainTextValueWithSecretReferenceValueIdRequestBuilderInternal(m.pathParameters, m.requestAdapter, secretReferenceValueId);
}
// Patch update the navigation property deviceConfigurations in deviceManagement
func (m *DeviceConfigurationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable, requestConfiguration *DeviceConfigurationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceConfigurationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceConfigurationable), nil
}
// UserStatuses provides operations to manage the userStatuses property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) UserStatuses()(*ie37508b74757cd3bd2b6dc98d9db2ba9651032937fe4385be4b4c611ade0f8fc.UserStatusesRequestBuilder) {
    return ie37508b74757cd3bd2b6dc98d9db2ba9651032937fe4385be4b4c611ade0f8fc.NewUserStatusesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UserStatusesById provides operations to manage the userStatuses property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) UserStatusesById(id string)(*i13d5635acea3d1befd38097b223c984a262de351ddcc99002d56e0362ecf8c7c.DeviceConfigurationUserStatusItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationUserStatus%2Did"] = id
    }
    return i13d5635acea3d1befd38097b223c984a262de351ddcc99002d56e0362ecf8c7c.NewDeviceConfigurationUserStatusItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UserStatusOverview provides operations to manage the userStatusOverview property of the microsoft.graph.deviceConfiguration entity.
func (m *DeviceConfigurationItemRequestBuilder) UserStatusOverview()(*i1cbc049d4f1b0a0c7493469fdf2a72ede1929262173246b80d83717fb3509836.UserStatusOverviewRequestBuilder) {
    return i1cbc049d4f1b0a0c7493469fdf2a72ede1929262173246b80d83717fb3509836.NewUserStatusOverviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
