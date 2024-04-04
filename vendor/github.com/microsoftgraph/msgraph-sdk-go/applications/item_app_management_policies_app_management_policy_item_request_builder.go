package applications

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\appManagementPolicies\{appManagementPolicy-id}
type ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal instantiates a new AppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    m := &ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}/appManagementPolicies/{appManagementPolicy%2Did}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams
    m.requestAdapter = requestAdapter
    return m
}
// NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder instantiates a new AppManagementPolicyItemRequestBuilder and sets the default values.
func NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemAppManagementPoliciesAppManagementPolicyItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of application entities.
func (m *ItemAppManagementPoliciesAppManagementPolicyItemRequestBuilder) Ref()(*ItemAppManagementPoliciesItemRefRequestBuilder) {
    return NewItemAppManagementPoliciesItemRefRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
