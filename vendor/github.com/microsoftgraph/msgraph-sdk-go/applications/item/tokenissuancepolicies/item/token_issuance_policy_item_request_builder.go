package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i083d9ad14fa1e44957def936929cba5238351218462d4ad39a246ee0bc99ba4a "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenissuancepolicies/item/ref"
)

// TokenIssuancePolicyItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\tokenIssuancePolicies\{tokenIssuancePolicy-id}
type TokenIssuancePolicyItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// NewTokenIssuancePolicyItemRequestBuilderInternal instantiates a new TokenIssuancePolicyItemRequestBuilder and sets the default values.
func NewTokenIssuancePolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TokenIssuancePolicyItemRequestBuilder) {
    m := &TokenIssuancePolicyItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}/tokenIssuancePolicies/{tokenIssuancePolicy%2Did}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTokenIssuancePolicyItemRequestBuilder instantiates a new TokenIssuancePolicyItemRequestBuilder and sets the default values.
func NewTokenIssuancePolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TokenIssuancePolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTokenIssuancePolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of application entities.
func (m *TokenIssuancePolicyItemRequestBuilder) Ref()(*i083d9ad14fa1e44957def936929cba5238351218462d4ad39a246ee0bc99ba4a.RefRequestBuilder) {
    return i083d9ad14fa1e44957def936929cba5238351218462d4ad39a246ee0bc99ba4a.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
