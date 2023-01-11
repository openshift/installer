package owneddevices

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2fa5782fd1e7f38ba8a5740355fbfa7e8d920abcc35255b8e04927467d7094a8 "github.com/microsoftgraph/msgraph-sdk-go/me/owneddevices/endpoint"
    i622bb27f3e3b795f5825bf562e76acd49951cfe27b98420952bed3a640f875f8 "github.com/microsoftgraph/msgraph-sdk-go/me/owneddevices/device"
    ia9847cef2752fa6526b8edfd0d9e0c8bc1bda772b96715c0026cacc9b8061e0a "github.com/microsoftgraph/msgraph-sdk-go/me/owneddevices/approleassignment"
    ibbf3bf8a68958b4caea8f9049039b0a3fea7c5502d8caf5b04faf3a68b819295 "github.com/microsoftgraph/msgraph-sdk-go/me/owneddevices/count"
)

// OwnedDevicesRequestBuilder provides operations to manage the ownedDevices property of the microsoft.graph.user entity.
type OwnedDevicesRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OwnedDevicesRequestBuilderGetQueryParameters devices that are owned by the user. Read-only. Nullable. Supports $expand.
type OwnedDevicesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// OwnedDevicesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OwnedDevicesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OwnedDevicesRequestBuilderGetQueryParameters
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *OwnedDevicesRequestBuilder) AppRoleAssignment()(*ia9847cef2752fa6526b8edfd0d9e0c8bc1bda772b96715c0026cacc9b8061e0a.AppRoleAssignmentRequestBuilder) {
    return ia9847cef2752fa6526b8edfd0d9e0c8bc1bda772b96715c0026cacc9b8061e0a.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOwnedDevicesRequestBuilderInternal instantiates a new OwnedDevicesRequestBuilder and sets the default values.
func NewOwnedDevicesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedDevicesRequestBuilder) {
    m := &OwnedDevicesRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/ownedDevices{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOwnedDevicesRequestBuilder instantiates a new OwnedDevicesRequestBuilder and sets the default values.
func NewOwnedDevicesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedDevicesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOwnedDevicesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *OwnedDevicesRequestBuilder) Count()(*ibbf3bf8a68958b4caea8f9049039b0a3fea7c5502d8caf5b04faf3a68b819295.CountRequestBuilder) {
    return ibbf3bf8a68958b4caea8f9049039b0a3fea7c5502d8caf5b04faf3a68b819295.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation devices that are owned by the user. Read-only. Nullable. Supports $expand.
func (m *OwnedDevicesRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OwnedDevicesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Device casts the previous resource to device.
func (m *OwnedDevicesRequestBuilder) Device()(*i622bb27f3e3b795f5825bf562e76acd49951cfe27b98420952bed3a640f875f8.DeviceRequestBuilder) {
    return i622bb27f3e3b795f5825bf562e76acd49951cfe27b98420952bed3a640f875f8.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Endpoint casts the previous resource to endpoint.
func (m *OwnedDevicesRequestBuilder) Endpoint()(*i2fa5782fd1e7f38ba8a5740355fbfa7e8d920abcc35255b8e04927467d7094a8.EndpointRequestBuilder) {
    return i2fa5782fd1e7f38ba8a5740355fbfa7e8d920abcc35255b8e04927467d7094a8.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get devices that are owned by the user. Read-only. Nullable. Supports $expand.
func (m *OwnedDevicesRequestBuilder) Get(ctx context.Context, requestConfiguration *OwnedDevicesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
