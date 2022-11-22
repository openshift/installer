package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    id596d9baf4431faf255161b3f48b6ab78efaf2b2a0636dee176263e72e58b5b0 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/claimsmappingpolicies/item/ref"
)

// ClaimsMappingPolicyItemRequestBuilder builds and executes requests for operations under \servicePrincipals\{servicePrincipal-id}\claimsMappingPolicies\{claimsMappingPolicy-id}
type ClaimsMappingPolicyItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// NewClaimsMappingPolicyItemRequestBuilderInternal instantiates a new ClaimsMappingPolicyItemRequestBuilder and sets the default values.
func NewClaimsMappingPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClaimsMappingPolicyItemRequestBuilder) {
    m := &ClaimsMappingPolicyItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/claimsMappingPolicies/{claimsMappingPolicy%2Did}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewClaimsMappingPolicyItemRequestBuilder instantiates a new ClaimsMappingPolicyItemRequestBuilder and sets the default values.
func NewClaimsMappingPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClaimsMappingPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClaimsMappingPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of servicePrincipal entities.
func (m *ClaimsMappingPolicyItemRequestBuilder) Ref()(*id596d9baf4431faf255161b3f48b6ab78efaf2b2a0636dee176263e72e58b5b0.RefRequestBuilder) {
    return id596d9baf4431faf255161b3f48b6ab78efaf2b2a0636dee176263e72e58b5b0.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
