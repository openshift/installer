package tokenissuancepolicies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0a7527c31d3fe38436ca4070a7467620ae738b216ecf91deddc0faf2a5967298 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenissuancepolicies/count"
    i76eae71f1d2dfce714e601f9c9310637a9b224b3dac53190f11a4812b4c6e5e2 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenissuancepolicies/ref"
)

// TokenIssuancePoliciesRequestBuilder provides operations to manage the tokenIssuancePolicies property of the microsoft.graph.application entity.
type TokenIssuancePoliciesRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TokenIssuancePoliciesRequestBuilderGetQueryParameters list the tokenIssuancePolicy objects that are assigned to an application.
type TokenIssuancePoliciesRequestBuilderGetQueryParameters struct {
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
// TokenIssuancePoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TokenIssuancePoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TokenIssuancePoliciesRequestBuilderGetQueryParameters
}
// NewTokenIssuancePoliciesRequestBuilderInternal instantiates a new TokenIssuancePoliciesRequestBuilder and sets the default values.
func NewTokenIssuancePoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TokenIssuancePoliciesRequestBuilder) {
    m := &TokenIssuancePoliciesRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}/tokenIssuancePolicies{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTokenIssuancePoliciesRequestBuilder instantiates a new TokenIssuancePoliciesRequestBuilder and sets the default values.
func NewTokenIssuancePoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TokenIssuancePoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTokenIssuancePoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *TokenIssuancePoliciesRequestBuilder) Count()(*i0a7527c31d3fe38436ca4070a7467620ae738b216ecf91deddc0faf2a5967298.CountRequestBuilder) {
    return i0a7527c31d3fe38436ca4070a7467620ae738b216ecf91deddc0faf2a5967298.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation list the tokenIssuancePolicy objects that are assigned to an application.
func (m *TokenIssuancePoliciesRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TokenIssuancePoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get list the tokenIssuancePolicy objects that are assigned to an application.
func (m *TokenIssuancePoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *TokenIssuancePoliciesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TokenIssuancePolicyCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTokenIssuancePolicyCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TokenIssuancePolicyCollectionResponseable), nil
}
// Ref provides operations to manage the collection of application entities.
func (m *TokenIssuancePoliciesRequestBuilder) Ref()(*i76eae71f1d2dfce714e601f9c9310637a9b224b3dac53190f11a4812b4c6e5e2.RefRequestBuilder) {
    return i76eae71f1d2dfce714e601f9c9310637a9b224b3dac53190f11a4812b4c6e5e2.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
