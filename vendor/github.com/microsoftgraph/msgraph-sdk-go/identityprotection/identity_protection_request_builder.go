package identityprotection

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i734828c35d0209cc42d94bd3da5e79564f27866a2e7dd505cec26e26172ddd00 "github.com/microsoftgraph/msgraph-sdk-go/identityprotection/riskdetections"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    if659cd8a8fa7e51b07cf2aa293599e63216814779670d31cfb9fa4a2facb3f81 "github.com/microsoftgraph/msgraph-sdk-go/identityprotection/riskyusers"
    i4fd6061886571b3d64bdb991095eb35ae40fd05f187e9565e0e84610c3889be0 "github.com/microsoftgraph/msgraph-sdk-go/identityprotection/riskyusers/item"
    i919d73984bc501fa850d0ae7fc63eed2d0af650bd7a764d682b904e3619f0c40 "github.com/microsoftgraph/msgraph-sdk-go/identityprotection/riskdetections/item"
)

// IdentityProtectionRequestBuilder provides operations to manage the identityProtectionRoot singleton.
type IdentityProtectionRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// IdentityProtectionRequestBuilderGetQueryParameters get identityProtection
type IdentityProtectionRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// IdentityProtectionRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type IdentityProtectionRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *IdentityProtectionRequestBuilderGetQueryParameters
}
// IdentityProtectionRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type IdentityProtectionRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewIdentityProtectionRequestBuilderInternal instantiates a new IdentityProtectionRequestBuilder and sets the default values.
func NewIdentityProtectionRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*IdentityProtectionRequestBuilder) {
    m := &IdentityProtectionRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/identityProtection{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewIdentityProtectionRequestBuilder instantiates a new IdentityProtectionRequestBuilder and sets the default values.
func NewIdentityProtectionRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*IdentityProtectionRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewIdentityProtectionRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get identityProtection
func (m *IdentityProtectionRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *IdentityProtectionRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update identityProtection
func (m *IdentityProtectionRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable, requestConfiguration *IdentityProtectionRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get identityProtection
func (m *IdentityProtectionRequestBuilder) Get(ctx context.Context, requestConfiguration *IdentityProtectionRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentityProtectionRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable), nil
}
// Patch update identityProtection
func (m *IdentityProtectionRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable, requestConfiguration *IdentityProtectionRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentityProtectionRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProtectionRootable), nil
}
// RiskDetections provides operations to manage the riskDetections property of the microsoft.graph.identityProtectionRoot entity.
func (m *IdentityProtectionRequestBuilder) RiskDetections()(*i734828c35d0209cc42d94bd3da5e79564f27866a2e7dd505cec26e26172ddd00.RiskDetectionsRequestBuilder) {
    return i734828c35d0209cc42d94bd3da5e79564f27866a2e7dd505cec26e26172ddd00.NewRiskDetectionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RiskDetectionsById provides operations to manage the riskDetections property of the microsoft.graph.identityProtectionRoot entity.
func (m *IdentityProtectionRequestBuilder) RiskDetectionsById(id string)(*i919d73984bc501fa850d0ae7fc63eed2d0af650bd7a764d682b904e3619f0c40.RiskDetectionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["riskDetection%2Did"] = id
    }
    return i919d73984bc501fa850d0ae7fc63eed2d0af650bd7a764d682b904e3619f0c40.NewRiskDetectionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RiskyUsers provides operations to manage the riskyUsers property of the microsoft.graph.identityProtectionRoot entity.
func (m *IdentityProtectionRequestBuilder) RiskyUsers()(*if659cd8a8fa7e51b07cf2aa293599e63216814779670d31cfb9fa4a2facb3f81.RiskyUsersRequestBuilder) {
    return if659cd8a8fa7e51b07cf2aa293599e63216814779670d31cfb9fa4a2facb3f81.NewRiskyUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RiskyUsersById provides operations to manage the riskyUsers property of the microsoft.graph.identityProtectionRoot entity.
func (m *IdentityProtectionRequestBuilder) RiskyUsersById(id string)(*i4fd6061886571b3d64bdb991095eb35ae40fd05f187e9565e0e84610c3889be0.RiskyUserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["riskyUser%2Did"] = id
    }
    return i4fd6061886571b3d64bdb991095eb35ae40fd05f187e9565e0e84610c3889be0.NewRiskyUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
