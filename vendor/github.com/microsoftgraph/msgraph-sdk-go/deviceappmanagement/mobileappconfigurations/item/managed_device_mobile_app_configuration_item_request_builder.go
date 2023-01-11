package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i20a49af14049cce13bbcbf77639ba345f23a979ba735144f96caf7cc46ac52d4 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/userstatuses"
    i6bbd09e352a0e87235e9c9775845f897603974356216a5236acfbea5f16155b0 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/assignments"
    i8d9da7a408610e6b2da98617dc04ec0862fb37929855ca63ba17948a3e7c6869 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/devicestatussummary"
    ia2a6efbe6de075696db4a436ca6bc9745df0eb596dd5de9281c86bf2bebcd940 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/assign"
    id36e539f78070cadf5b4f1fcefe8eeff25266761b7e3541b45788010e5cdfb9b "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/userstatussummary"
    iebd1b47b377fdbd0e32c261a4ebc852aebb3266e58bde42f93309b4e2edc6d16 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/devicestatuses"
    i5644f6cefb6e37ba63ec8582d3a158ec48759832b4146d5ffd2d350288b5e5ca "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/devicestatuses/item"
    i72a4318a7e8f28adce62f7520318129090631e70fa23a9371561c34fac1f0582 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/userstatuses/item"
    i88ff03a861b0c98567d5cfa6d93d8d9d0de34bb73fdd948e55765d9858a275c3 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement/mobileappconfigurations/item/assignments/item"
)

// ManagedDeviceMobileAppConfigurationItemRequestBuilder provides operations to manage the mobileAppConfigurations property of the microsoft.graph.deviceAppManagement entity.
type ManagedDeviceMobileAppConfigurationItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ManagedDeviceMobileAppConfigurationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceMobileAppConfigurationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ManagedDeviceMobileAppConfigurationItemRequestBuilderGetQueryParameters the Managed Device Mobile Application Configurations.
type ManagedDeviceMobileAppConfigurationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ManagedDeviceMobileAppConfigurationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceMobileAppConfigurationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ManagedDeviceMobileAppConfigurationItemRequestBuilderGetQueryParameters
}
// ManagedDeviceMobileAppConfigurationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceMobileAppConfigurationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Assign provides operations to call the assign method.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) Assign()(*ia2a6efbe6de075696db4a436ca6bc9745df0eb596dd5de9281c86bf2bebcd940.AssignRequestBuilder) {
    return ia2a6efbe6de075696db4a436ca6bc9745df0eb596dd5de9281c86bf2bebcd940.NewAssignRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Assignments provides operations to manage the assignments property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) Assignments()(*i6bbd09e352a0e87235e9c9775845f897603974356216a5236acfbea5f16155b0.AssignmentsRequestBuilder) {
    return i6bbd09e352a0e87235e9c9775845f897603974356216a5236acfbea5f16155b0.NewAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentsById provides operations to manage the assignments property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) AssignmentsById(id string)(*i88ff03a861b0c98567d5cfa6d93d8d9d0de34bb73fdd948e55765d9858a275c3.ManagedDeviceMobileAppConfigurationAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["managedDeviceMobileAppConfigurationAssignment%2Did"] = id
    }
    return i88ff03a861b0c98567d5cfa6d93d8d9d0de34bb73fdd948e55765d9858a275c3.NewManagedDeviceMobileAppConfigurationAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewManagedDeviceMobileAppConfigurationItemRequestBuilderInternal instantiates a new ManagedDeviceMobileAppConfigurationItemRequestBuilder and sets the default values.
func NewManagedDeviceMobileAppConfigurationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceMobileAppConfigurationItemRequestBuilder) {
    m := &ManagedDeviceMobileAppConfigurationItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/deviceAppManagement/mobileAppConfigurations/{managedDeviceMobileAppConfiguration%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewManagedDeviceMobileAppConfigurationItemRequestBuilder instantiates a new ManagedDeviceMobileAppConfigurationItemRequestBuilder and sets the default values.
func NewManagedDeviceMobileAppConfigurationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceMobileAppConfigurationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewManagedDeviceMobileAppConfigurationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property mobileAppConfigurations for deviceAppManagement
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the Managed Device Mobile Application Configurations.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property mobileAppConfigurations in deviceAppManagement
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property mobileAppConfigurations for deviceAppManagement
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DeviceStatuses provides operations to manage the deviceStatuses property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) DeviceStatuses()(*iebd1b47b377fdbd0e32c261a4ebc852aebb3266e58bde42f93309b4e2edc6d16.DeviceStatusesRequestBuilder) {
    return iebd1b47b377fdbd0e32c261a4ebc852aebb3266e58bde42f93309b4e2edc6d16.NewDeviceStatusesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceStatusesById provides operations to manage the deviceStatuses property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) DeviceStatusesById(id string)(*i5644f6cefb6e37ba63ec8582d3a158ec48759832b4146d5ffd2d350288b5e5ca.ManagedDeviceMobileAppConfigurationDeviceStatusItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["managedDeviceMobileAppConfigurationDeviceStatus%2Did"] = id
    }
    return i5644f6cefb6e37ba63ec8582d3a158ec48759832b4146d5ffd2d350288b5e5ca.NewManagedDeviceMobileAppConfigurationDeviceStatusItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceStatusSummary provides operations to manage the deviceStatusSummary property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) DeviceStatusSummary()(*i8d9da7a408610e6b2da98617dc04ec0862fb37929855ca63ba17948a3e7c6869.DeviceStatusSummaryRequestBuilder) {
    return i8d9da7a408610e6b2da98617dc04ec0862fb37929855ca63ba17948a3e7c6869.NewDeviceStatusSummaryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the Managed Device Mobile Application Configurations.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable), nil
}
// Patch update the navigation property mobileAppConfigurations in deviceAppManagement
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable, requestConfiguration *ManagedDeviceMobileAppConfigurationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceMobileAppConfigurationable), nil
}
// UserStatuses provides operations to manage the userStatuses property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) UserStatuses()(*i20a49af14049cce13bbcbf77639ba345f23a979ba735144f96caf7cc46ac52d4.UserStatusesRequestBuilder) {
    return i20a49af14049cce13bbcbf77639ba345f23a979ba735144f96caf7cc46ac52d4.NewUserStatusesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UserStatusesById provides operations to manage the userStatuses property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) UserStatusesById(id string)(*i72a4318a7e8f28adce62f7520318129090631e70fa23a9371561c34fac1f0582.ManagedDeviceMobileAppConfigurationUserStatusItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["managedDeviceMobileAppConfigurationUserStatus%2Did"] = id
    }
    return i72a4318a7e8f28adce62f7520318129090631e70fa23a9371561c34fac1f0582.NewManagedDeviceMobileAppConfigurationUserStatusItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UserStatusSummary provides operations to manage the userStatusSummary property of the microsoft.graph.managedDeviceMobileAppConfiguration entity.
func (m *ManagedDeviceMobileAppConfigurationItemRequestBuilder) UserStatusSummary()(*id36e539f78070cadf5b4f1fcefe8eeff25266761b7e3541b45788010e5cdfb9b.UserStatusSummaryRequestBuilder) {
    return id36e539f78070cadf5b4f1fcefe8eeff25266761b7e3541b45788010e5cdfb9b.NewUserStatusSummaryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
