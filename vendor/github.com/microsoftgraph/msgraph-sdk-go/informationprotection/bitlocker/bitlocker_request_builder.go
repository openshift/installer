package bitlocker

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i5dc93f982e4ecf9567e43c10c4133c7c4df8d34016db35ea82573ec8759ce43d "github.com/microsoftgraph/msgraph-sdk-go/informationprotection/bitlocker/recoverykeys"
    i53a72d4c5b633ebaf2e92a3eae90f203f023e1af23ceaaa7b670a43df7cb21ed "github.com/microsoftgraph/msgraph-sdk-go/informationprotection/bitlocker/recoverykeys/item"
)

// BitlockerRequestBuilder provides operations to manage the bitlocker property of the microsoft.graph.informationProtection entity.
type BitlockerRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// BitlockerRequestBuilderGetQueryParameters get bitlocker from informationProtection
type BitlockerRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// BitlockerRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type BitlockerRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *BitlockerRequestBuilderGetQueryParameters
}
// NewBitlockerRequestBuilderInternal instantiates a new BitlockerRequestBuilder and sets the default values.
func NewBitlockerRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BitlockerRequestBuilder) {
    m := &BitlockerRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/informationProtection/bitlocker{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewBitlockerRequestBuilder instantiates a new BitlockerRequestBuilder and sets the default values.
func NewBitlockerRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*BitlockerRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewBitlockerRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get bitlocker from informationProtection
func (m *BitlockerRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *BitlockerRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get bitlocker from informationProtection
func (m *BitlockerRequestBuilder) Get(ctx context.Context, requestConfiguration *BitlockerRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Bitlockerable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateBitlockerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Bitlockerable), nil
}
// RecoveryKeys provides operations to manage the recoveryKeys property of the microsoft.graph.bitlocker entity.
func (m *BitlockerRequestBuilder) RecoveryKeys()(*i5dc93f982e4ecf9567e43c10c4133c7c4df8d34016db35ea82573ec8759ce43d.RecoveryKeysRequestBuilder) {
    return i5dc93f982e4ecf9567e43c10c4133c7c4df8d34016db35ea82573ec8759ce43d.NewRecoveryKeysRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RecoveryKeysById provides operations to manage the recoveryKeys property of the microsoft.graph.bitlocker entity.
func (m *BitlockerRequestBuilder) RecoveryKeysById(id string)(*i53a72d4c5b633ebaf2e92a3eae90f203f023e1af23ceaaa7b670a43df7cb21ed.BitlockerRecoveryKeyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["bitlockerRecoveryKey%2Did"] = id
    }
    return i53a72d4c5b633ebaf2e92a3eae90f203f023e1af23ceaaa7b670a43df7cb21ed.NewBitlockerRecoveryKeyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
