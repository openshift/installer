package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i02e61d7c24506eb3092637591a2fe281450dbe31d56751354fa7a05f72a89a6d "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/disablelostmode"
    i15a1a10a0c024f930afef3b930ac6a0a48505b8f50dc765c9208ba2816bfb8fb "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/resetpasscode"
    i35e5de1b6d1217fb83cdfe34b93910dd53c07bcbaed491f2d9dd30bf097a2e27 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/cleanwindowsdevice"
    i3f5657959eef872a8c3b583ceb3e90ed3871b1ab719ab6bf86fbd61313c45635 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/users"
    i4be1f6e23df553f3d25539b06b6e72e653c54673be6c256636259b2dfc46b028 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/deleteuserfromsharedappledevice"
    i5626258f553fa5ab744196016cddf3bdb99eca3ad27f9716b75a1816ddcac4e4 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/windowsdefenderupdatesignatures"
    i579420ac2331f8db64f4efb8f6b2f9b7e59da9e53c2f60c7189463b2f6729a66 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/retire"
    i61111df9932978d316b7c2d5bf743126bd3dfc765577371722fdfd9540e546c6 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/deviceconfigurationstates"
    i73fe8c13857767fd45ca26708823e85bca6083a2a2ebdb0dda0d4fe4e91cdf17 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/updatewindowsdeviceaccount"
    i957d1faa5f562ab89577eb46e79203c7e50075bcfa0818779d21777aac9c4774 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/locatedevice"
    ib2246236511e9e6bd5f149ec83d2768718af4bee5cf6b48497a15ba9859004da "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/devicecompliancepolicystates"
    ib4e430bec489c09a72c51ad3159e98695c04be6d6cbd03769b4074193f098ce8 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/rebootnow"
    ib6fe4058a826f6d827fec2a61b53730458811701046b26ceb4044feeb7ddf7d1 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/recoverpasscode"
    ib975ed5a3045ece0bc94c5dd2f4c5a4ea484dc80b1830219bb365eea8fd973c0 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/logoutsharedappledeviceactiveuser"
    ic1c10c22035a19935689244731094da0d49c68bb4707c06ce0ecf896730268c2 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/windowsdefenderscan"
    ic2df49236f6d10a558da4728bb3e2ce53874fd50cfc571bf9e691b66d446221c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/syncdevice"
    ic61092c249f07b9e411d18b3b8310e71210bdf5c5de775c13112cab083c5d082 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/wipe"
    icc2068c1350cb53b049a71e87eb1cf48456fa9e9b48d284b8abaa38d4207ffc2 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/shutdown"
    id96b7ae916d8c8e465b8a821c2fe41d4d020ec0a7610953b779987f93a48d88e "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/devicecategory"
    ie6232acc82dadd162e8e29367892609c026901eba827c7fdbb14386e5d36b99d "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/remotelock"
    if1e16de23b7ea9301d8d5208c98e249a57e7e1e4df4785a249e053121429dd1e "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/requestremoteassistance"
    if245b21847517f2bdf21fefeaac5356812b1d5b2cad61e07f1e97c61a72adfeb "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/bypassactivationlock"
    i257ff0bc433efcedf2c2cd195e3b06a2fbc6131ac66e39a2a9444f64b497cbf5 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/deviceconfigurationstates/item"
    i30c735a31ceee62863c65dfaa2404ac107c81c4dfd9ccc96d5eee45ca24f04f2 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/manageddevices/item/devicecompliancepolicystates/item"
)

// ManagedDeviceItemRequestBuilder provides operations to manage the managedDevices property of the microsoft.graph.deviceManagement entity.
type ManagedDeviceItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ManagedDeviceItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ManagedDeviceItemRequestBuilderGetQueryParameters the list of managed devices.
type ManagedDeviceItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ManagedDeviceItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ManagedDeviceItemRequestBuilderGetQueryParameters
}
// ManagedDeviceItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDeviceItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BypassActivationLock provides operations to call the bypassActivationLock method.
func (m *ManagedDeviceItemRequestBuilder) BypassActivationLock()(*if245b21847517f2bdf21fefeaac5356812b1d5b2cad61e07f1e97c61a72adfeb.BypassActivationLockRequestBuilder) {
    return if245b21847517f2bdf21fefeaac5356812b1d5b2cad61e07f1e97c61a72adfeb.NewBypassActivationLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CleanWindowsDevice provides operations to call the cleanWindowsDevice method.
func (m *ManagedDeviceItemRequestBuilder) CleanWindowsDevice()(*i35e5de1b6d1217fb83cdfe34b93910dd53c07bcbaed491f2d9dd30bf097a2e27.CleanWindowsDeviceRequestBuilder) {
    return i35e5de1b6d1217fb83cdfe34b93910dd53c07bcbaed491f2d9dd30bf097a2e27.NewCleanWindowsDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewManagedDeviceItemRequestBuilderInternal instantiates a new ManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDeviceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceItemRequestBuilder) {
    m := &ManagedDeviceItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/deviceManagement/managedDevices/{managedDevice%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewManagedDeviceItemRequestBuilder instantiates a new ManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDeviceItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewManagedDeviceItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property managedDevices for deviceManagement
func (m *ManagedDeviceItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ManagedDeviceItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the list of managed devices.
func (m *ManagedDeviceItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ManagedDeviceItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property managedDevices in deviceManagement
func (m *ManagedDeviceItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, requestConfiguration *ManagedDeviceItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property managedDevices for deviceManagement
func (m *ManagedDeviceItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ManagedDeviceItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DeleteUserFromSharedAppleDevice provides operations to call the deleteUserFromSharedAppleDevice method.
func (m *ManagedDeviceItemRequestBuilder) DeleteUserFromSharedAppleDevice()(*i4be1f6e23df553f3d25539b06b6e72e653c54673be6c256636259b2dfc46b028.DeleteUserFromSharedAppleDeviceRequestBuilder) {
    return i4be1f6e23df553f3d25539b06b6e72e653c54673be6c256636259b2dfc46b028.NewDeleteUserFromSharedAppleDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCategory provides operations to manage the deviceCategory property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCategory()(*id96b7ae916d8c8e465b8a821c2fe41d4d020ec0a7610953b779987f93a48d88e.DeviceCategoryRequestBuilder) {
    return id96b7ae916d8c8e465b8a821c2fe41d4d020ec0a7610953b779987f93a48d88e.NewDeviceCategoryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStates provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStates()(*ib2246236511e9e6bd5f149ec83d2768718af4bee5cf6b48497a15ba9859004da.DeviceCompliancePolicyStatesRequestBuilder) {
    return ib2246236511e9e6bd5f149ec83d2768718af4bee5cf6b48497a15ba9859004da.NewDeviceCompliancePolicyStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStatesById provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStatesById(id string)(*i30c735a31ceee62863c65dfaa2404ac107c81c4dfd9ccc96d5eee45ca24f04f2.DeviceCompliancePolicyStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceCompliancePolicyState%2Did"] = id
    }
    return i30c735a31ceee62863c65dfaa2404ac107c81c4dfd9ccc96d5eee45ca24f04f2.NewDeviceCompliancePolicyStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceConfigurationStates provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStates()(*i61111df9932978d316b7c2d5bf743126bd3dfc765577371722fdfd9540e546c6.DeviceConfigurationStatesRequestBuilder) {
    return i61111df9932978d316b7c2d5bf743126bd3dfc765577371722fdfd9540e546c6.NewDeviceConfigurationStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceConfigurationStatesById provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStatesById(id string)(*i257ff0bc433efcedf2c2cd195e3b06a2fbc6131ac66e39a2a9444f64b497cbf5.DeviceConfigurationStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationState%2Did"] = id
    }
    return i257ff0bc433efcedf2c2cd195e3b06a2fbc6131ac66e39a2a9444f64b497cbf5.NewDeviceConfigurationStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DisableLostMode provides operations to call the disableLostMode method.
func (m *ManagedDeviceItemRequestBuilder) DisableLostMode()(*i02e61d7c24506eb3092637591a2fe281450dbe31d56751354fa7a05f72a89a6d.DisableLostModeRequestBuilder) {
    return i02e61d7c24506eb3092637591a2fe281450dbe31d56751354fa7a05f72a89a6d.NewDisableLostModeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the list of managed devices.
func (m *ManagedDeviceItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ManagedDeviceItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable), nil
}
// LocateDevice provides operations to call the locateDevice method.
func (m *ManagedDeviceItemRequestBuilder) LocateDevice()(*i957d1faa5f562ab89577eb46e79203c7e50075bcfa0818779d21777aac9c4774.LocateDeviceRequestBuilder) {
    return i957d1faa5f562ab89577eb46e79203c7e50075bcfa0818779d21777aac9c4774.NewLocateDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LogoutSharedAppleDeviceActiveUser provides operations to call the logoutSharedAppleDeviceActiveUser method.
func (m *ManagedDeviceItemRequestBuilder) LogoutSharedAppleDeviceActiveUser()(*ib975ed5a3045ece0bc94c5dd2f4c5a4ea484dc80b1830219bb365eea8fd973c0.LogoutSharedAppleDeviceActiveUserRequestBuilder) {
    return ib975ed5a3045ece0bc94c5dd2f4c5a4ea484dc80b1830219bb365eea8fd973c0.NewLogoutSharedAppleDeviceActiveUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property managedDevices in deviceManagement
func (m *ManagedDeviceItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, requestConfiguration *ManagedDeviceItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable), nil
}
// RebootNow provides operations to call the rebootNow method.
func (m *ManagedDeviceItemRequestBuilder) RebootNow()(*ib4e430bec489c09a72c51ad3159e98695c04be6d6cbd03769b4074193f098ce8.RebootNowRequestBuilder) {
    return ib4e430bec489c09a72c51ad3159e98695c04be6d6cbd03769b4074193f098ce8.NewRebootNowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RecoverPasscode provides operations to call the recoverPasscode method.
func (m *ManagedDeviceItemRequestBuilder) RecoverPasscode()(*ib6fe4058a826f6d827fec2a61b53730458811701046b26ceb4044feeb7ddf7d1.RecoverPasscodeRequestBuilder) {
    return ib6fe4058a826f6d827fec2a61b53730458811701046b26ceb4044feeb7ddf7d1.NewRecoverPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoteLock provides operations to call the remoteLock method.
func (m *ManagedDeviceItemRequestBuilder) RemoteLock()(*ie6232acc82dadd162e8e29367892609c026901eba827c7fdbb14386e5d36b99d.RemoteLockRequestBuilder) {
    return ie6232acc82dadd162e8e29367892609c026901eba827c7fdbb14386e5d36b99d.NewRemoteLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RequestRemoteAssistance provides operations to call the requestRemoteAssistance method.
func (m *ManagedDeviceItemRequestBuilder) RequestRemoteAssistance()(*if1e16de23b7ea9301d8d5208c98e249a57e7e1e4df4785a249e053121429dd1e.RequestRemoteAssistanceRequestBuilder) {
    return if1e16de23b7ea9301d8d5208c98e249a57e7e1e4df4785a249e053121429dd1e.NewRequestRemoteAssistanceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResetPasscode provides operations to call the resetPasscode method.
func (m *ManagedDeviceItemRequestBuilder) ResetPasscode()(*i15a1a10a0c024f930afef3b930ac6a0a48505b8f50dc765c9208ba2816bfb8fb.ResetPasscodeRequestBuilder) {
    return i15a1a10a0c024f930afef3b930ac6a0a48505b8f50dc765c9208ba2816bfb8fb.NewResetPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Retire provides operations to call the retire method.
func (m *ManagedDeviceItemRequestBuilder) Retire()(*i579420ac2331f8db64f4efb8f6b2f9b7e59da9e53c2f60c7189463b2f6729a66.RetireRequestBuilder) {
    return i579420ac2331f8db64f4efb8f6b2f9b7e59da9e53c2f60c7189463b2f6729a66.NewRetireRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShutDown provides operations to call the shutDown method.
func (m *ManagedDeviceItemRequestBuilder) ShutDown()(*icc2068c1350cb53b049a71e87eb1cf48456fa9e9b48d284b8abaa38d4207ffc2.ShutDownRequestBuilder) {
    return icc2068c1350cb53b049a71e87eb1cf48456fa9e9b48d284b8abaa38d4207ffc2.NewShutDownRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SyncDevice provides operations to call the syncDevice method.
func (m *ManagedDeviceItemRequestBuilder) SyncDevice()(*ic2df49236f6d10a558da4728bb3e2ce53874fd50cfc571bf9e691b66d446221c.SyncDeviceRequestBuilder) {
    return ic2df49236f6d10a558da4728bb3e2ce53874fd50cfc571bf9e691b66d446221c.NewSyncDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UpdateWindowsDeviceAccount provides operations to call the updateWindowsDeviceAccount method.
func (m *ManagedDeviceItemRequestBuilder) UpdateWindowsDeviceAccount()(*i73fe8c13857767fd45ca26708823e85bca6083a2a2ebdb0dda0d4fe4e91cdf17.UpdateWindowsDeviceAccountRequestBuilder) {
    return i73fe8c13857767fd45ca26708823e85bca6083a2a2ebdb0dda0d4fe4e91cdf17.NewUpdateWindowsDeviceAccountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Users provides operations to manage the users property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) Users()(*i3f5657959eef872a8c3b583ceb3e90ed3871b1ab719ab6bf86fbd61313c45635.UsersRequestBuilder) {
    return i3f5657959eef872a8c3b583ceb3e90ed3871b1ab719ab6bf86fbd61313c45635.NewUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderScan provides operations to call the windowsDefenderScan method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderScan()(*ic1c10c22035a19935689244731094da0d49c68bb4707c06ce0ecf896730268c2.WindowsDefenderScanRequestBuilder) {
    return ic1c10c22035a19935689244731094da0d49c68bb4707c06ce0ecf896730268c2.NewWindowsDefenderScanRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderUpdateSignatures provides operations to call the windowsDefenderUpdateSignatures method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderUpdateSignatures()(*i5626258f553fa5ab744196016cddf3bdb99eca3ad27f9716b75a1816ddcac4e4.WindowsDefenderUpdateSignaturesRequestBuilder) {
    return i5626258f553fa5ab744196016cddf3bdb99eca3ad27f9716b75a1816ddcac4e4.NewWindowsDefenderUpdateSignaturesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Wipe provides operations to call the wipe method.
func (m *ManagedDeviceItemRequestBuilder) Wipe()(*ic61092c249f07b9e411d18b3b8310e71210bdf5c5de775c13112cab083c5d082.WipeRequestBuilder) {
    return ic61092c249f07b9e411d18b3b8310e71210bdf5c5de775c13112cab083c5d082.NewWipeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
