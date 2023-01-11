package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0fdd6e00c3a70ace39a5146d6c55b15af0de25c040237603b09168231c256f9b "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/deviceconfigurationstates"
    i14175b6b023a6eb71ef6ba830a21cb6cd13dd609e7dc15bd405158e0b3df79ae "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/requestremoteassistance"
    i1754afc016e1fdc501acb882736ad04e30bc5650f53d3c074ab8ac8a51c6fa8f "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/devicecompliancepolicystates"
    i2b3c6b981da171994b1b46a7294676fb6bdc0c93bb24e519dd1306595ead2c80 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/updatewindowsdeviceaccount"
    i30838d0e3fd5b58dc055f8efb1004348c1be948971f57117d40f2bbf1b4535f0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/cleanwindowsdevice"
    i322fafcb7ae8cab29e30bcc25707daf94d942483b0b2026f21e1a15ec8be0f92 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/shutdown"
    i4b83748504956b892ea1a7f149037fca181d069c7d1673044fada70771e33093 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/retire"
    i5889a0f91c17b2ff6be4d4d7f398164e86344dbddaf46d10d8dbf893bd1fdcae "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/deleteuserfromsharedappledevice"
    i8cb1d9fa2d4d44aff3eee0f4bd2fa9c55841709e6d7067c757f8da63d9f42acb "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/logoutsharedappledeviceactiveuser"
    i94ba78e1d9aea49972909e0476fb3f1e5a65ad70fb7b6b0bfa5b59fcc655e084 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/syncdevice"
    i992fc981e39a2a010a5298e4dd852c4b0b9d2285b09665f5884b4fd2dcca552e "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/remotelock"
    i9daab2b0b5282b9d8617504ab740c1b4eeb654e8993c10684911865beb7f6513 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/devicecategory"
    ia76862a49870b0532e8c4b374ba7384c2818ea7114ca76a78ad21e98b0c92501 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/locatedevice"
    iad89fd53db2b505c2f511349b0beb9bdda197521ce058a20d54d9f5a00f41230 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/wipe"
    ib6f3d272b640361ab70ab361d75bbbcd3009458c8853a93188c3a1a2a636b702 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/users"
    ic645e4837ba91a567171d12c28ba82eb15a0c5ba22dfc84e46d5174771406adc "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/windowsdefenderscan"
    ic9510153a0140fff3d042d83894479343cdd52e080d68ae1293f3eb456c459e3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/resetpasscode"
    icb15920be1a61530593e35306720a47cfee2b31808ee7b11ef7d482703361e75 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/rebootnow"
    id09c4b3a73bddb54267ba24e9836a171941a344b098f8a53e9044afc145d5939 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/disablelostmode"
    if07c243615b7cd8f004a1303c7440d0fb98ff8cb9456cb22dee43fbf230ad815 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/bypassactivationlock"
    if36ea65c1b212f583ceb77d3c50dc0ebe8e096afabd7f9662ae9ee0019217b7a "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/recoverpasscode"
    if54975d74aac049350e7bc12c684f45645594faffab68df3d530a91cd5877663 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/windowsdefenderupdatesignatures"
    i5b5d439175020246e151a9558d6d568e6bc4bd9fdd376e082964b229d2b71629 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/devicecompliancepolicystates/item"
    i5e9268af6ea2183cbe612b8d18a9aeec3e11e6ad80e059611e3cafbb64907758 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item/deviceconfigurationstates/item"
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
func (m *ManagedDeviceItemRequestBuilder) BypassActivationLock()(*if07c243615b7cd8f004a1303c7440d0fb98ff8cb9456cb22dee43fbf230ad815.BypassActivationLockRequestBuilder) {
    return if07c243615b7cd8f004a1303c7440d0fb98ff8cb9456cb22dee43fbf230ad815.NewBypassActivationLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CleanWindowsDevice provides operations to call the cleanWindowsDevice method.
func (m *ManagedDeviceItemRequestBuilder) CleanWindowsDevice()(*i30838d0e3fd5b58dc055f8efb1004348c1be948971f57117d40f2bbf1b4535f0.CleanWindowsDeviceRequestBuilder) {
    return i30838d0e3fd5b58dc055f8efb1004348c1be948971f57117d40f2bbf1b4535f0.NewCleanWindowsDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewManagedDeviceItemRequestBuilderInternal instantiates a new ManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDeviceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDeviceItemRequestBuilder) {
    m := &ManagedDeviceItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/managedDevices/{managedDevice%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property managedDevices for users
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
// CreatePatchRequestInformation update the navigation property managedDevices in users
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
// Delete delete navigation property managedDevices for users
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
func (m *ManagedDeviceItemRequestBuilder) DeleteUserFromSharedAppleDevice()(*i5889a0f91c17b2ff6be4d4d7f398164e86344dbddaf46d10d8dbf893bd1fdcae.DeleteUserFromSharedAppleDeviceRequestBuilder) {
    return i5889a0f91c17b2ff6be4d4d7f398164e86344dbddaf46d10d8dbf893bd1fdcae.NewDeleteUserFromSharedAppleDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCategory provides operations to manage the deviceCategory property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCategory()(*i9daab2b0b5282b9d8617504ab740c1b4eeb654e8993c10684911865beb7f6513.DeviceCategoryRequestBuilder) {
    return i9daab2b0b5282b9d8617504ab740c1b4eeb654e8993c10684911865beb7f6513.NewDeviceCategoryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStates provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStates()(*i1754afc016e1fdc501acb882736ad04e30bc5650f53d3c074ab8ac8a51c6fa8f.DeviceCompliancePolicyStatesRequestBuilder) {
    return i1754afc016e1fdc501acb882736ad04e30bc5650f53d3c074ab8ac8a51c6fa8f.NewDeviceCompliancePolicyStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceCompliancePolicyStatesById provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStatesById(id string)(*i5b5d439175020246e151a9558d6d568e6bc4bd9fdd376e082964b229d2b71629.DeviceCompliancePolicyStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceCompliancePolicyState%2Did"] = id
    }
    return i5b5d439175020246e151a9558d6d568e6bc4bd9fdd376e082964b229d2b71629.NewDeviceCompliancePolicyStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceConfigurationStates provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStates()(*i0fdd6e00c3a70ace39a5146d6c55b15af0de25c040237603b09168231c256f9b.DeviceConfigurationStatesRequestBuilder) {
    return i0fdd6e00c3a70ace39a5146d6c55b15af0de25c040237603b09168231c256f9b.NewDeviceConfigurationStatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceConfigurationStatesById provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) DeviceConfigurationStatesById(id string)(*i5e9268af6ea2183cbe612b8d18a9aeec3e11e6ad80e059611e3cafbb64907758.DeviceConfigurationStateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceConfigurationState%2Did"] = id
    }
    return i5e9268af6ea2183cbe612b8d18a9aeec3e11e6ad80e059611e3cafbb64907758.NewDeviceConfigurationStateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DisableLostMode provides operations to call the disableLostMode method.
func (m *ManagedDeviceItemRequestBuilder) DisableLostMode()(*id09c4b3a73bddb54267ba24e9836a171941a344b098f8a53e9044afc145d5939.DisableLostModeRequestBuilder) {
    return id09c4b3a73bddb54267ba24e9836a171941a344b098f8a53e9044afc145d5939.NewDisableLostModeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *ManagedDeviceItemRequestBuilder) LocateDevice()(*ia76862a49870b0532e8c4b374ba7384c2818ea7114ca76a78ad21e98b0c92501.LocateDeviceRequestBuilder) {
    return ia76862a49870b0532e8c4b374ba7384c2818ea7114ca76a78ad21e98b0c92501.NewLocateDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LogoutSharedAppleDeviceActiveUser provides operations to call the logoutSharedAppleDeviceActiveUser method.
func (m *ManagedDeviceItemRequestBuilder) LogoutSharedAppleDeviceActiveUser()(*i8cb1d9fa2d4d44aff3eee0f4bd2fa9c55841709e6d7067c757f8da63d9f42acb.LogoutSharedAppleDeviceActiveUserRequestBuilder) {
    return i8cb1d9fa2d4d44aff3eee0f4bd2fa9c55841709e6d7067c757f8da63d9f42acb.NewLogoutSharedAppleDeviceActiveUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property managedDevices in users
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
func (m *ManagedDeviceItemRequestBuilder) RebootNow()(*icb15920be1a61530593e35306720a47cfee2b31808ee7b11ef7d482703361e75.RebootNowRequestBuilder) {
    return icb15920be1a61530593e35306720a47cfee2b31808ee7b11ef7d482703361e75.NewRebootNowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RecoverPasscode provides operations to call the recoverPasscode method.
func (m *ManagedDeviceItemRequestBuilder) RecoverPasscode()(*if36ea65c1b212f583ceb77d3c50dc0ebe8e096afabd7f9662ae9ee0019217b7a.RecoverPasscodeRequestBuilder) {
    return if36ea65c1b212f583ceb77d3c50dc0ebe8e096afabd7f9662ae9ee0019217b7a.NewRecoverPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoteLock provides operations to call the remoteLock method.
func (m *ManagedDeviceItemRequestBuilder) RemoteLock()(*i992fc981e39a2a010a5298e4dd852c4b0b9d2285b09665f5884b4fd2dcca552e.RemoteLockRequestBuilder) {
    return i992fc981e39a2a010a5298e4dd852c4b0b9d2285b09665f5884b4fd2dcca552e.NewRemoteLockRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RequestRemoteAssistance provides operations to call the requestRemoteAssistance method.
func (m *ManagedDeviceItemRequestBuilder) RequestRemoteAssistance()(*i14175b6b023a6eb71ef6ba830a21cb6cd13dd609e7dc15bd405158e0b3df79ae.RequestRemoteAssistanceRequestBuilder) {
    return i14175b6b023a6eb71ef6ba830a21cb6cd13dd609e7dc15bd405158e0b3df79ae.NewRequestRemoteAssistanceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResetPasscode provides operations to call the resetPasscode method.
func (m *ManagedDeviceItemRequestBuilder) ResetPasscode()(*ic9510153a0140fff3d042d83894479343cdd52e080d68ae1293f3eb456c459e3.ResetPasscodeRequestBuilder) {
    return ic9510153a0140fff3d042d83894479343cdd52e080d68ae1293f3eb456c459e3.NewResetPasscodeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Retire provides operations to call the retire method.
func (m *ManagedDeviceItemRequestBuilder) Retire()(*i4b83748504956b892ea1a7f149037fca181d069c7d1673044fada70771e33093.RetireRequestBuilder) {
    return i4b83748504956b892ea1a7f149037fca181d069c7d1673044fada70771e33093.NewRetireRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ShutDown provides operations to call the shutDown method.
func (m *ManagedDeviceItemRequestBuilder) ShutDown()(*i322fafcb7ae8cab29e30bcc25707daf94d942483b0b2026f21e1a15ec8be0f92.ShutDownRequestBuilder) {
    return i322fafcb7ae8cab29e30bcc25707daf94d942483b0b2026f21e1a15ec8be0f92.NewShutDownRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SyncDevice provides operations to call the syncDevice method.
func (m *ManagedDeviceItemRequestBuilder) SyncDevice()(*i94ba78e1d9aea49972909e0476fb3f1e5a65ad70fb7b6b0bfa5b59fcc655e084.SyncDeviceRequestBuilder) {
    return i94ba78e1d9aea49972909e0476fb3f1e5a65ad70fb7b6b0bfa5b59fcc655e084.NewSyncDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UpdateWindowsDeviceAccount provides operations to call the updateWindowsDeviceAccount method.
func (m *ManagedDeviceItemRequestBuilder) UpdateWindowsDeviceAccount()(*i2b3c6b981da171994b1b46a7294676fb6bdc0c93bb24e519dd1306595ead2c80.UpdateWindowsDeviceAccountRequestBuilder) {
    return i2b3c6b981da171994b1b46a7294676fb6bdc0c93bb24e519dd1306595ead2c80.NewUpdateWindowsDeviceAccountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Users provides operations to manage the users property of the microsoft.graph.managedDevice entity.
func (m *ManagedDeviceItemRequestBuilder) Users()(*ib6f3d272b640361ab70ab361d75bbbcd3009458c8853a93188c3a1a2a636b702.UsersRequestBuilder) {
    return ib6f3d272b640361ab70ab361d75bbbcd3009458c8853a93188c3a1a2a636b702.NewUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderScan provides operations to call the windowsDefenderScan method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderScan()(*ic645e4837ba91a567171d12c28ba82eb15a0c5ba22dfc84e46d5174771406adc.WindowsDefenderScanRequestBuilder) {
    return ic645e4837ba91a567171d12c28ba82eb15a0c5ba22dfc84e46d5174771406adc.NewWindowsDefenderScanRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WindowsDefenderUpdateSignatures provides operations to call the windowsDefenderUpdateSignatures method.
func (m *ManagedDeviceItemRequestBuilder) WindowsDefenderUpdateSignatures()(*if54975d74aac049350e7bc12c684f45645594faffab68df3d530a91cd5877663.WindowsDefenderUpdateSignaturesRequestBuilder) {
    return if54975d74aac049350e7bc12c684f45645594faffab68df3d530a91cd5877663.NewWindowsDefenderUpdateSignaturesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Wipe provides operations to call the wipe method.
func (m *ManagedDeviceItemRequestBuilder) Wipe()(*iad89fd53db2b505c2f511349b0beb9bdda197521ce058a20d54d9f5a00f41230.WipeRequestBuilder) {
    return iad89fd53db2b505c2f511349b0beb9bdda197521ce058a20d54d9f5a00f41230.NewWipeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
