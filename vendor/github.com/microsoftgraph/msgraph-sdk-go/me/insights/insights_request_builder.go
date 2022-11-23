package insights

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i29d2c155f681081063d0b2c0ab92796f0e2b78c469f71af45a16d8240b348ebf "github.com/microsoftgraph/msgraph-sdk-go/me/insights/shared"
    i8893ecb7e9d38ec8961e185f29391303f74779021473d1c37caeafd1abebef9a "github.com/microsoftgraph/msgraph-sdk-go/me/insights/trending"
    i99fad623e5b2b4a7ae85326fa692745befba2f6b4fb179d4572c3e7219fbf3ac "github.com/microsoftgraph/msgraph-sdk-go/me/insights/used"
    i7eb746d3d91c620af79236bba2af4b708c6b65d5f4ec8f9cc42e5d4305f1c396 "github.com/microsoftgraph/msgraph-sdk-go/me/insights/shared/item"
    ibf98f8f6ecbe2d059f97c4d3970656db9710f9bafa095bf25ac791007a80e997 "github.com/microsoftgraph/msgraph-sdk-go/me/insights/trending/item"
    ic47fb69b7ecc13ce7f81973ffdde08efe7b70445311d1e83964064ef79497098 "github.com/microsoftgraph/msgraph-sdk-go/me/insights/used/item"
)

// InsightsRequestBuilder provides operations to manage the insights property of the microsoft.graph.user entity.
type InsightsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// InsightsRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type InsightsRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// InsightsRequestBuilderGetQueryParameters get insights from me
type InsightsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// InsightsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type InsightsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *InsightsRequestBuilderGetQueryParameters
}
// InsightsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type InsightsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewInsightsRequestBuilderInternal instantiates a new InsightsRequestBuilder and sets the default values.
func NewInsightsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InsightsRequestBuilder) {
    m := &InsightsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/insights{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewInsightsRequestBuilder instantiates a new InsightsRequestBuilder and sets the default values.
func NewInsightsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InsightsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewInsightsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property insights for me
func (m *InsightsRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *InsightsRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get insights from me
func (m *InsightsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *InsightsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property insights in me
func (m *InsightsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable, requestConfiguration *InsightsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property insights for me
func (m *InsightsRequestBuilder) Delete(ctx context.Context, requestConfiguration *InsightsRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get insights from me
func (m *InsightsRequestBuilder) Get(ctx context.Context, requestConfiguration *InsightsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOfficeGraphInsightsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable), nil
}
// Patch update the navigation property insights in me
func (m *InsightsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable, requestConfiguration *InsightsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOfficeGraphInsightsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OfficeGraphInsightsable), nil
}
// Shared provides operations to manage the shared property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) Shared()(*i29d2c155f681081063d0b2c0ab92796f0e2b78c469f71af45a16d8240b348ebf.SharedRequestBuilder) {
    return i29d2c155f681081063d0b2c0ab92796f0e2b78c469f71af45a16d8240b348ebf.NewSharedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedById provides operations to manage the shared property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) SharedById(id string)(*i7eb746d3d91c620af79236bba2af4b708c6b65d5f4ec8f9cc42e5d4305f1c396.SharedInsightItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedInsight%2Did"] = id
    }
    return i7eb746d3d91c620af79236bba2af4b708c6b65d5f4ec8f9cc42e5d4305f1c396.NewSharedInsightItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Trending provides operations to manage the trending property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) Trending()(*i8893ecb7e9d38ec8961e185f29391303f74779021473d1c37caeafd1abebef9a.TrendingRequestBuilder) {
    return i8893ecb7e9d38ec8961e185f29391303f74779021473d1c37caeafd1abebef9a.NewTrendingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TrendingById provides operations to manage the trending property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) TrendingById(id string)(*ibf98f8f6ecbe2d059f97c4d3970656db9710f9bafa095bf25ac791007a80e997.TrendingItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["trending%2Did"] = id
    }
    return ibf98f8f6ecbe2d059f97c4d3970656db9710f9bafa095bf25ac791007a80e997.NewTrendingItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Used provides operations to manage the used property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) Used()(*i99fad623e5b2b4a7ae85326fa692745befba2f6b4fb179d4572c3e7219fbf3ac.UsedRequestBuilder) {
    return i99fad623e5b2b4a7ae85326fa692745befba2f6b4fb179d4572c3e7219fbf3ac.NewUsedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UsedById provides operations to manage the used property of the microsoft.graph.officeGraphInsights entity.
func (m *InsightsRequestBuilder) UsedById(id string)(*ic47fb69b7ecc13ce7f81973ffdde08efe7b70445311d1e83964064ef79497098.UsedInsightItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["usedInsight%2Did"] = id
    }
    return ic47fb69b7ecc13ce7f81973ffdde08efe7b70445311d1e83964064ef79497098.NewUsedInsightItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
