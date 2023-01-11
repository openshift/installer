package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i97b7f699ffeb18b6e5ae1e2e27f52abcd604d2c0f2ecd9cc8ed8216c23b05dbd "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/homerealmdiscoverypolicies/item/ref"
)

// HomeRealmDiscoveryPolicyItemRequestBuilder builds and executes requests for operations under \servicePrincipals\{servicePrincipal-id}\homeRealmDiscoveryPolicies\{homeRealmDiscoveryPolicy-id}
type HomeRealmDiscoveryPolicyItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// NewHomeRealmDiscoveryPolicyItemRequestBuilderInternal instantiates a new HomeRealmDiscoveryPolicyItemRequestBuilder and sets the default values.
func NewHomeRealmDiscoveryPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*HomeRealmDiscoveryPolicyItemRequestBuilder) {
    m := &HomeRealmDiscoveryPolicyItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/homeRealmDiscoveryPolicies/{homeRealmDiscoveryPolicy%2Did}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewHomeRealmDiscoveryPolicyItemRequestBuilder instantiates a new HomeRealmDiscoveryPolicyItemRequestBuilder and sets the default values.
func NewHomeRealmDiscoveryPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*HomeRealmDiscoveryPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewHomeRealmDiscoveryPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of servicePrincipal entities.
func (m *HomeRealmDiscoveryPolicyItemRequestBuilder) Ref()(*i97b7f699ffeb18b6e5ae1e2e27f52abcd604d2c0f2ecd9cc8ed8216c23b05dbd.RefRequestBuilder) {
    return i97b7f699ffeb18b6e5ae1e2e27f52abcd604d2c0f2ecd9cc8ed8216c23b05dbd.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
