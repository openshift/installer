package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i13f9a52f0faf606b01bfa3955e81ee30242ef21fa0101aa087fbede7948f1b27 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/resetpasscode"
    i1cea6eb12ff2031dedeb2ca8e7862c20336acfdade144ef517a8e5f962fde1f6 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/cleanwindowsdevice"
    i1e119e421440451e7a53df7e5ca93960ef8e6779e31ca34c2612ddf6e4f5328e "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/windowsdefenderscan"
    i227d0481e3969744ed37e693056ba632a2619393c65c499dec2bb5818ee14670 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/disablelostmode"
    i2e55400e3a6b1518964ef2369bbc57f4e9a9c4a253afb8263d0c7dd315d4c634 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/locatedevice"
    i2f0f1986c49cb55444d8f3078b1690316ddcf576fc6f835186d573e09b7ecc57 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/windowsdefenderupdatesignatures"
    i48553f07197f253783753adbaf7d41660e2a12cf4e334a15ce0a642355263b5e "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/retire"
    i4910726fdcc07627f867a8093698926b8b1ebb256d1ec4e04b21c86a7ee852eb "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/logoutsharedappledeviceactiveuser"
    i4d4ffdac6d6f354742f9f83e0d75cd1f102d1dd45a13d732b9a0c59657c5b53e "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/requestremoteassistance"
    i62df4e156ab0e11abc901553dd3d14a35d775f80cadfec4b46046b962d16f910 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/wipe"
    i674214359fe29326a03866451c37e7e463c6460897eff0b1713788bcd244994d "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/remotelock"
    i6c92c5c6ee7fb6a63e6a0166f2ec31aa4b47828c1a5ff577350d48023fb89b7a "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/rebootnow"
    i83acc3dd442c42abdbabb9f4fb5b7f362e8c2133c09e2ab7c8978a58dcaba49a "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/devicecompliancepolicystates"
    i93be5ba1e2ed1dcc8ae36f63d85f3edb6388941e09d0e839b2f1a5f66d4919f4 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/updatewindowsdeviceaccount"
    i97119b774733e620c79e77fe13b9bb32fe545bb2ed890cab4c2a70c7c3e2a442 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/bypassactivationlock"
    ia9eca82f666e7935a7636c1a6840da82b973be896a102904ee6f6cb28bc3f48a "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/deviceconfigurationstates"
    ibcabf3b7cd8bdd83cd21114c5e5f0274cd32297b1b565fc77144db1020c9d980 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/users"
    ibece143defa494d6c2d968eb6a807859f4c7077b7f099cd2c37993aec1843123 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/syncdevice"
    ic7e7f352418887585e7f8a86e825eb5a78d5d8b112d8eb01d1373aec67367a26 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/devicecategory"
    iceb6710bce94d2cc31507d20a496d678d25042f6be1b3ab09b691e833abcccc9 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/recoverpasscode"
    id8aabeb9ac25ed8c755eee299962060c1b32b52e8491a650fa00ecfd2bcb5f28 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/shutdown"
    id98973b9a80b6d1c37c9c0f1386e5d814f433103bc9e0e7fba3ee12787a54dd8 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/deleteuserfromsharedappledevice"
    i04ce841e69bdde74a0a5ab2a5f5b64b7e52f890d6dce32e3f002ee145b68f803 "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/devicecompliancepolicystates/item"
    i635fc9cf6abf1cd6f28dfa5885addac5c0d71ed597869597de4dc6569fcace0a "github.com/microsoftgraph/msgraph-sdk-go/me/manageddevices/item/deviceconfigurationstates/item"
)

// ManagedDeviceItemRequestBuilder provides operations to manage the managedDevices property of the microsoft.graph.user entity.
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
// ManagedDeviceItemRequestBuilderGetQueryParameters the managed devices associated with the user.
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
func (m *ManagedDeviceItemRequestBuilder) BypassActivationLock()(*i97119b774733e620c79e77fe13b9bb32fe545bb2ed890cab4c2a70c7c3e2a442.BypassActivationLockRequestBuilder) {
    return i97119b774733e620c79e77fe13b9bb32fe545bb2ed890cab4c2a70c7c3e2a442.NewBypassActivationLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CleanWindowsDevice provides operations to call the cleanWindowsDevice method.
func (m *ManagedDeviceItemRequestBuilder) CleanWindowsDevice()(*i1cea6eb12ff2031dedeb2ca8e7862c20336acfdade144ef517a8e5f962fde1f6.CleanWindowsDeviceRequestBuilder) {
    return i1cea6eb12ff2031dedeb2ca8e7862c20336acfdade144ef517a8e5f962fde1f6.NewCleanWindowsDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewManagedDeviceItemRequestBuilderInternal instantiates a new ManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDeviceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceItemRequestBuilder) {
    m := &ManagedDeviceItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/managedDevices/{managedDevice%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property managedDevices for me
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
// CreateGetRequestInformation the managed devices associated with the user.
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
// CreatePatchRequestInformation update the navigation property managedDevices in me
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
// Delete delete navigation property managedDevices for me
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
func (m *ManagedDeviceItemRequestBuilder) DeleteUserFromSharedAppleDevice()(*id98973b9a80b6d1c37c9c0f1386e5d814f433103bc9e0e7fba3ee12787a54dd8.DeleteUserFromSharedAppleDeviceRequestBuilder) {
    return id98973b9a80b6d1c37c9c0f1386e5d814f433103bc9e0e7fba3ee12787a54dd8.NewDeleteUserFromSharedAppleDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCategory provides operations to manage the deviceCategory property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCategory()(*ic7e7f352418887585e7f8a86e825eb5a78d5d8b112d8eb01d1373aec67367a26.DeviceCategoryRequestBuilder) {
    return ic7e7f352418887585e7f8a86e825eb5a78d5d8b112d8eb01d1373aec67367a26.NewDeviceCategoryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStates provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStates()(*i83acc3dd442c42abdbabb9f4fb5b7f362e8c2133c09e2ab7c8978a58dcaba49a.DeviceCompliancePolicyStatesRequestBuilder) {
    return i83acc3dd442c42abdbabb9f4fb5b7f362e8c2133c09e2ab7c8978a58dcaba49a.NewDeviceCompliancePolicyStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStatesById provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStatesById(id string)(*i04ce841e69bdde74a0a5ab2a5f5b64b7e52f890d6dce32e3f002ee145b68f803.DeviceCompliancePolicyStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceCompliancePolicyState%2Did"] = id
    }
    return i04ce841e69bdde74a0a5ab2a5f5b64b7e52f890d6dce32e3f002ee145b68f803.NewDeviceCompliancePolicyStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceConfigurationStates provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStates()(*ia9eca82f666e7935a7636c1a6840da82b973be896a102904ee6f6cb28bc3f48a.DeviceConfigurationStatesRequestBuilder) {
    return ia9eca82f666e7935a7636c1a6840da82b973be896a102904ee6f6cb28bc3f48a.NewDeviceConfigurationStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceConfigurationStatesById provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStatesById(id string)(*i635fc9cf6abf1cd6f28dfa5885addac5c0d71ed597869597de4dc6569fcace0a.DeviceConfigurationStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationState%2Did"] = id
    }
    return i635fc9cf6abf1cd6f28dfa5885addac5c0d71ed597869597de4dc6569fcace0a.NewDeviceConfigurationStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DisableLostMode provides operations to call the disableLostMode method.
func (m *ManagedDeviceItemRequestBuilder) DisableLostMode()(*i227d0481e3969744ed37e693056ba632a2619393c65c499dec2bb5818ee14670.DisableLostModeRequestBuilder) {
    return i227d0481e3969744ed37e693056ba632a2619393c65c499dec2bb5818ee14670.NewDisableLostModeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the managed devices associated with the user.
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
func (m *ManagedDeviceItemRequestBuilder) LocateDevice()(*i2e55400e3a6b1518964ef2369bbc57f4e9a9c4a253afb8263d0c7dd315d4c634.LocateDeviceRequestBuilder) {
    return i2e55400e3a6b1518964ef2369bbc57f4e9a9c4a253afb8263d0c7dd315d4c634.NewLocateDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LogoutSharedAppleDeviceActiveUser provides operations to call the logoutSharedAppleDeviceActiveUser method.
func (m *ManagedDeviceItemRequestBuilder) LogoutSharedAppleDeviceActiveUser()(*i4910726fdcc07627f867a8093698926b8b1ebb256d1ec4e04b21c86a7ee852eb.LogoutSharedAppleDeviceActiveUserRequestBuilder) {
    return i4910726fdcc07627f867a8093698926b8b1ebb256d1ec4e04b21c86a7ee852eb.NewLogoutSharedAppleDeviceActiveUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property managedDevices in me
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
func (m *ManagedDeviceItemRequestBuilder) RebootNow()(*i6c92c5c6ee7fb6a63e6a0166f2ec31aa4b47828c1a5ff577350d48023fb89b7a.RebootNowRequestBuilder) {
    return i6c92c5c6ee7fb6a63e6a0166f2ec31aa4b47828c1a5ff577350d48023fb89b7a.NewRebootNowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RecoverPasscode provides operations to call the recoverPasscode method.
func (m *ManagedDeviceItemRequestBuilder) RecoverPasscode()(*iceb6710bce94d2cc31507d20a496d678d25042f6be1b3ab09b691e833abcccc9.RecoverPasscodeRequestBuilder) {
    return iceb6710bce94d2cc31507d20a496d678d25042f6be1b3ab09b691e833abcccc9.NewRecoverPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoteLock provides operations to call the remoteLock method.
func (m *ManagedDeviceItemRequestBuilder) RemoteLock()(*i674214359fe29326a03866451c37e7e463c6460897eff0b1713788bcd244994d.RemoteLockRequestBuilder) {
    return i674214359fe29326a03866451c37e7e463c6460897eff0b1713788bcd244994d.NewRemoteLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RequestRemoteAssistance provides operations to call the requestRemoteAssistance method.
func (m *ManagedDeviceItemRequestBuilder) RequestRemoteAssistance()(*i4d4ffdac6d6f354742f9f83e0d75cd1f102d1dd45a13d732b9a0c59657c5b53e.RequestRemoteAssistanceRequestBuilder) {
    return i4d4ffdac6d6f354742f9f83e0d75cd1f102d1dd45a13d732b9a0c59657c5b53e.NewRequestRemoteAssistanceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResetPasscode provides operations to call the resetPasscode method.
func (m *ManagedDeviceItemRequestBuilder) ResetPasscode()(*i13f9a52f0faf606b01bfa3955e81ee30242ef21fa0101aa087fbede7948f1b27.ResetPasscodeRequestBuilder) {
    return i13f9a52f0faf606b01bfa3955e81ee30242ef21fa0101aa087fbede7948f1b27.NewResetPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Retire provides operations to call the retire method.
func (m *ManagedDeviceItemRequestBuilder) Retire()(*i48553f07197f253783753adbaf7d41660e2a12cf4e334a15ce0a642355263b5e.RetireRequestBuilder) {
    return i48553f07197f253783753adbaf7d41660e2a12cf4e334a15ce0a642355263b5e.NewRetireRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShutDown provides operations to call the shutDown method.
func (m *ManagedDeviceItemRequestBuilder) ShutDown()(*id8aabeb9ac25ed8c755eee299962060c1b32b52e8491a650fa00ecfd2bcb5f28.ShutDownRequestBuilder) {
    return id8aabeb9ac25ed8c755eee299962060c1b32b52e8491a650fa00ecfd2bcb5f28.NewShutDownRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SyncDevice provides operations to call the syncDevice method.
func (m *ManagedDeviceItemRequestBuilder) SyncDevice()(*ibece143defa494d6c2d968eb6a807859f4c7077b7f099cd2c37993aec1843123.SyncDeviceRequestBuilder) {
    return ibece143defa494d6c2d968eb6a807859f4c7077b7f099cd2c37993aec1843123.NewSyncDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UpdateWindowsDeviceAccount provides operations to call the updateWindowsDeviceAccount method.
func (m *ManagedDeviceItemRequestBuilder) UpdateWindowsDeviceAccount()(*i93be5ba1e2ed1dcc8ae36f63d85f3edb6388941e09d0e839b2f1a5f66d4919f4.UpdateWindowsDeviceAccountRequestBuilder) {
    return i93be5ba1e2ed1dcc8ae36f63d85f3edb6388941e09d0e839b2f1a5f66d4919f4.NewUpdateWindowsDeviceAccountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Users provides operations to manage the users property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) Users()(*ibcabf3b7cd8bdd83cd21114c5e5f0274cd32297b1b565fc77144db1020c9d980.UsersRequestBuilder) {
    return ibcabf3b7cd8bdd83cd21114c5e5f0274cd32297b1b565fc77144db1020c9d980.NewUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderScan provides operations to call the windowsDefenderScan method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderScan()(*i1e119e421440451e7a53df7e5ca93960ef8e6779e31ca34c2612ddf6e4f5328e.WindowsDefenderScanRequestBuilder) {
    return i1e119e421440451e7a53df7e5ca93960ef8e6779e31ca34c2612ddf6e4f5328e.NewWindowsDefenderScanRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderUpdateSignatures provides operations to call the windowsDefenderUpdateSignatures method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderUpdateSignatures()(*i2f0f1986c49cb55444d8f3078b1690316ddcf576fc6f835186d573e09b7ecc57.WindowsDefenderUpdateSignaturesRequestBuilder) {
    return i2f0f1986c49cb55444d8f3078b1690316ddcf576fc6f835186d573e09b7ecc57.NewWindowsDefenderUpdateSignaturesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Wipe provides operations to call the wipe method.
func (m *ManagedDeviceItemRequestBuilder) Wipe()(*i62df4e156ab0e11abc901553dd3d14a35d775f80cadfec4b46046b962d16f910.WipeRequestBuilder) {
    return i62df4e156ab0e11abc901553dd3d14a35d775f80cadfec4b46046b962d16f910.NewWipeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
