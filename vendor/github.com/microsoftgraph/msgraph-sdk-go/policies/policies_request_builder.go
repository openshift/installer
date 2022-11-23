package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i15c1315dca3d448080a500f31a0343661829915a3d57e0bac2b5ef3d5b702cd7 "github.com/microsoftgraph/msgraph-sdk-go/policies/identitysecuritydefaultsenforcementpolicy"
    i15eb5011e1661e803ff4682ed565e8609c250e7ec6aaf6510e8d4cb37c79c25f "github.com/microsoftgraph/msgraph-sdk-go/policies/rolemanagementpolicyassignments"
    i50fbe096cb4bcbe54ce4584948f2fa87f8486532011e2c9024a169ba2a3dacde "github.com/microsoftgraph/msgraph-sdk-go/policies/activitybasedtimeoutpolicies"
    i598aadab7f9532dcd78b5237b03ae7e6ac9390b26cc04cb85038bcb797fcb4e3 "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies"
    i5eff42eff2580c09a93ee6ffbdd897f38083be436dd15c59a2e4891c29fd0b5e "github.com/microsoftgraph/msgraph-sdk-go/policies/authenticationmethodspolicy"
    i684bc90d26f7d7a7b87627223964f1cc0cab2499ef34d5bdbb89a5ac5cf49209 "github.com/microsoftgraph/msgraph-sdk-go/policies/authenticationflowspolicy"
    i6894c2598f4ac73a5643fa86cbbf110868e7746e8c9ed5da452b0170eebe03b3 "github.com/microsoftgraph/msgraph-sdk-go/policies/homerealmdiscoverypolicies"
    i7cd232f5a1288d29625a9b1dc442455d81df9bedea7cc94eef2584add7856e18 "github.com/microsoftgraph/msgraph-sdk-go/policies/authorizationpolicy"
    i800d2e5efbf3f2562fb570f6abc9bec962927b27f95f6e0025aa92bfe23f13ab "github.com/microsoftgraph/msgraph-sdk-go/policies/claimsmappingpolicies"
    i8cf0127b53d26923e1649ec43c875bc6f0f70ece0a4b9d3cf3377d9e0b4720cf "github.com/microsoftgraph/msgraph-sdk-go/policies/tokenissuancepolicies"
    i94e6a9fd02466104e1f244dfcb22cf331ccf308e2f911adea2fe5a4a87bf5ccc "github.com/microsoftgraph/msgraph-sdk-go/policies/permissiongrantpolicies"
    i966a3041c5833ccf56ce1da16c013fd76d376cfc8b939cc01e633d8b495cb208 "github.com/microsoftgraph/msgraph-sdk-go/policies/conditionalaccesspolicies"
    i992d94460d3d6690972cde6f4f03058259573cadd2ea4b4ed7196e5ddfc1e6e5 "github.com/microsoftgraph/msgraph-sdk-go/policies/rolemanagementpolicies"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ic72674a7b815c6dfc1c4fcd42f7ca9ca5d4b8ea0c1a65ea7317905471fa20a0f "github.com/microsoftgraph/msgraph-sdk-go/policies/adminconsentrequestpolicy"
    iebb445e03b5a1725abc955bb63d1192b435aa40b5fee4976b801435c3fc9ec06 "github.com/microsoftgraph/msgraph-sdk-go/policies/tokenlifetimepolicies"
    ifc08343a031e849b74e359a9f87e4aa9d0ae860bdb89c5f9131d85756a69debb "github.com/microsoftgraph/msgraph-sdk-go/policies/crosstenantaccesspolicy"
    i15708f9856db4407ccc78f2397a412ca44ce57d9fbe9dd3e17f7bf4747b8cf71 "github.com/microsoftgraph/msgraph-sdk-go/policies/rolemanagementpolicies/item"
    i569717b3ca8408936df0d0a4e253024a8a7b041522350cb7c8977c2c62a9b9d6 "github.com/microsoftgraph/msgraph-sdk-go/policies/tokenlifetimepolicies/item"
    i75b801e9d4a8913922e67e9e75dd091c9280212e449c4d8e76d13ab572b48c86 "github.com/microsoftgraph/msgraph-sdk-go/policies/activitybasedtimeoutpolicies/item"
    i77b5f2108467d08d93443762e93f7967aca61b1916748814eb2ddc49bc0fb0ed "github.com/microsoftgraph/msgraph-sdk-go/policies/tokenissuancepolicies/item"
    i843df2bef70b4bcdc0e002955adffd0d0793c6957377635528f401c17a961d02 "github.com/microsoftgraph/msgraph-sdk-go/policies/conditionalaccesspolicies/item"
    ia9d7ff4c255eb343aa34f45265c76078ea351f2ba2bad54efe9863a838dd38f0 "github.com/microsoftgraph/msgraph-sdk-go/policies/permissiongrantpolicies/item"
    iae405fd168f9b22628a993c1d320fde787ac74895fc32dab558860d3df5e3cf7 "github.com/microsoftgraph/msgraph-sdk-go/policies/homerealmdiscoverypolicies/item"
    ifd481b5a9f37ded8a4d12f7c0db8d1e5cc30754c49b7cbbab66e580c09455018 "github.com/microsoftgraph/msgraph-sdk-go/policies/claimsmappingpolicies/item"
    ifdb3aa1e2fe0663a1e8072ad2268b1766bb6cee12e54afb4c9bfe76480e4184a "github.com/microsoftgraph/msgraph-sdk-go/policies/rolemanagementpolicyassignments/item"
    iff6a6d8b424e8a919281280c87147525819afe04366318faafc92dbee9442123 "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item"
)

// PoliciesRequestBuilder provides operations to manage the policyRoot singleton.
type PoliciesRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PoliciesRequestBuilderGetQueryParameters get policies
type PoliciesRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PoliciesRequestBuilderGetQueryParameters
}
// PoliciesRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PoliciesRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ActivityBasedTimeoutPolicies provides operations to manage the activityBasedTimeoutPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ActivityBasedTimeoutPolicies()(*i50fbe096cb4bcbe54ce4584948f2fa87f8486532011e2c9024a169ba2a3dacde.ActivityBasedTimeoutPoliciesRequestBuilder) {
    return i50fbe096cb4bcbe54ce4584948f2fa87f8486532011e2c9024a169ba2a3dacde.NewActivityBasedTimeoutPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ActivityBasedTimeoutPoliciesById provides operations to manage the activityBasedTimeoutPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ActivityBasedTimeoutPoliciesById(id string)(*i75b801e9d4a8913922e67e9e75dd091c9280212e449c4d8e76d13ab572b48c86.ActivityBasedTimeoutPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["activityBasedTimeoutPolicy%2Did"] = id
    }
    return i75b801e9d4a8913922e67e9e75dd091c9280212e449c4d8e76d13ab572b48c86.NewActivityBasedTimeoutPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AdminConsentRequestPolicy provides operations to manage the adminConsentRequestPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) AdminConsentRequestPolicy()(*ic72674a7b815c6dfc1c4fcd42f7ca9ca5d4b8ea0c1a65ea7317905471fa20a0f.AdminConsentRequestPolicyRequestBuilder) {
    return ic72674a7b815c6dfc1c4fcd42f7ca9ca5d4b8ea0c1a65ea7317905471fa20a0f.NewAdminConsentRequestPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthenticationFlowsPolicy provides operations to manage the authenticationFlowsPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) AuthenticationFlowsPolicy()(*i684bc90d26f7d7a7b87627223964f1cc0cab2499ef34d5bdbb89a5ac5cf49209.AuthenticationFlowsPolicyRequestBuilder) {
    return i684bc90d26f7d7a7b87627223964f1cc0cab2499ef34d5bdbb89a5ac5cf49209.NewAuthenticationFlowsPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthenticationMethodsPolicy provides operations to manage the authenticationMethodsPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) AuthenticationMethodsPolicy()(*i5eff42eff2580c09a93ee6ffbdd897f38083be436dd15c59a2e4891c29fd0b5e.AuthenticationMethodsPolicyRequestBuilder) {
    return i5eff42eff2580c09a93ee6ffbdd897f38083be436dd15c59a2e4891c29fd0b5e.NewAuthenticationMethodsPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthorizationPolicy provides operations to manage the authorizationPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) AuthorizationPolicy()(*i7cd232f5a1288d29625a9b1dc442455d81df9bedea7cc94eef2584add7856e18.AuthorizationPolicyRequestBuilder) {
    return i7cd232f5a1288d29625a9b1dc442455d81df9bedea7cc94eef2584add7856e18.NewAuthorizationPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClaimsMappingPolicies provides operations to manage the claimsMappingPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ClaimsMappingPolicies()(*i800d2e5efbf3f2562fb570f6abc9bec962927b27f95f6e0025aa92bfe23f13ab.ClaimsMappingPoliciesRequestBuilder) {
    return i800d2e5efbf3f2562fb570f6abc9bec962927b27f95f6e0025aa92bfe23f13ab.NewClaimsMappingPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClaimsMappingPoliciesById provides operations to manage the claimsMappingPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ClaimsMappingPoliciesById(id string)(*ifd481b5a9f37ded8a4d12f7c0db8d1e5cc30754c49b7cbbab66e580c09455018.ClaimsMappingPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["claimsMappingPolicy%2Did"] = id
    }
    return ifd481b5a9f37ded8a4d12f7c0db8d1e5cc30754c49b7cbbab66e580c09455018.NewClaimsMappingPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ConditionalAccessPolicies provides operations to manage the conditionalAccessPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ConditionalAccessPolicies()(*i966a3041c5833ccf56ce1da16c013fd76d376cfc8b939cc01e633d8b495cb208.ConditionalAccessPoliciesRequestBuilder) {
    return i966a3041c5833ccf56ce1da16c013fd76d376cfc8b939cc01e633d8b495cb208.NewConditionalAccessPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ConditionalAccessPoliciesById provides operations to manage the conditionalAccessPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) ConditionalAccessPoliciesById(id string)(*i843df2bef70b4bcdc0e002955adffd0d0793c6957377635528f401c17a961d02.ConditionalAccessPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conditionalAccessPolicy%2Did"] = id
    }
    return i843df2bef70b4bcdc0e002955adffd0d0793c6957377635528f401c17a961d02.NewConditionalAccessPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPoliciesRequestBuilderInternal instantiates a new PoliciesRequestBuilder and sets the default values.
func NewPoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PoliciesRequestBuilder) {
    m := &PoliciesRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/policies{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPoliciesRequestBuilder instantiates a new PoliciesRequestBuilder and sets the default values.
func NewPoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get policies
func (m *PoliciesRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update policies
func (m *PoliciesRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, requestConfiguration *PoliciesRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CrossTenantAccessPolicy provides operations to manage the crossTenantAccessPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) CrossTenantAccessPolicy()(*ifc08343a031e849b74e359a9f87e4aa9d0ae860bdb89c5f9131d85756a69debb.CrossTenantAccessPolicyRequestBuilder) {
    return ifc08343a031e849b74e359a9f87e4aa9d0ae860bdb89c5f9131d85756a69debb.NewCrossTenantAccessPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FeatureRolloutPolicies provides operations to manage the featureRolloutPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) FeatureRolloutPolicies()(*i598aadab7f9532dcd78b5237b03ae7e6ac9390b26cc04cb85038bcb797fcb4e3.FeatureRolloutPoliciesRequestBuilder) {
    return i598aadab7f9532dcd78b5237b03ae7e6ac9390b26cc04cb85038bcb797fcb4e3.NewFeatureRolloutPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FeatureRolloutPoliciesById provides operations to manage the featureRolloutPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) FeatureRolloutPoliciesById(id string)(*iff6a6d8b424e8a919281280c87147525819afe04366318faafc92dbee9442123.FeatureRolloutPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["featureRolloutPolicy%2Did"] = id
    }
    return iff6a6d8b424e8a919281280c87147525819afe04366318faafc92dbee9442123.NewFeatureRolloutPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get get policies
func (m *PoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *PoliciesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePolicyRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable), nil
}
// HomeRealmDiscoveryPolicies provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) HomeRealmDiscoveryPolicies()(*i6894c2598f4ac73a5643fa86cbbf110868e7746e8c9ed5da452b0170eebe03b3.HomeRealmDiscoveryPoliciesRequestBuilder) {
    return i6894c2598f4ac73a5643fa86cbbf110868e7746e8c9ed5da452b0170eebe03b3.NewHomeRealmDiscoveryPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// HomeRealmDiscoveryPoliciesById provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) HomeRealmDiscoveryPoliciesById(id string)(*iae405fd168f9b22628a993c1d320fde787ac74895fc32dab558860d3df5e3cf7.HomeRealmDiscoveryPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["homeRealmDiscoveryPolicy%2Did"] = id
    }
    return iae405fd168f9b22628a993c1d320fde787ac74895fc32dab558860d3df5e3cf7.NewHomeRealmDiscoveryPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// IdentitySecurityDefaultsEnforcementPolicy provides operations to manage the identitySecurityDefaultsEnforcementPolicy property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) IdentitySecurityDefaultsEnforcementPolicy()(*i15c1315dca3d448080a500f31a0343661829915a3d57e0bac2b5ef3d5b702cd7.IdentitySecurityDefaultsEnforcementPolicyRequestBuilder) {
    return i15c1315dca3d448080a500f31a0343661829915a3d57e0bac2b5ef3d5b702cd7.NewIdentitySecurityDefaultsEnforcementPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update policies
func (m *PoliciesRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, requestConfiguration *PoliciesRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePolicyRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable), nil
}
// PermissionGrantPolicies provides operations to manage the permissionGrantPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) PermissionGrantPolicies()(*i94e6a9fd02466104e1f244dfcb22cf331ccf308e2f911adea2fe5a4a87bf5ccc.PermissionGrantPoliciesRequestBuilder) {
    return i94e6a9fd02466104e1f244dfcb22cf331ccf308e2f911adea2fe5a4a87bf5ccc.NewPermissionGrantPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionGrantPoliciesById provides operations to manage the permissionGrantPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) PermissionGrantPoliciesById(id string)(*ia9d7ff4c255eb343aa34f45265c76078ea351f2ba2bad54efe9863a838dd38f0.PermissionGrantPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permissionGrantPolicy%2Did"] = id
    }
    return ia9d7ff4c255eb343aa34f45265c76078ea351f2ba2bad54efe9863a838dd38f0.NewPermissionGrantPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleManagementPolicies provides operations to manage the roleManagementPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) RoleManagementPolicies()(*i992d94460d3d6690972cde6f4f03058259573cadd2ea4b4ed7196e5ddfc1e6e5.RoleManagementPoliciesRequestBuilder) {
    return i992d94460d3d6690972cde6f4f03058259573cadd2ea4b4ed7196e5ddfc1e6e5.NewRoleManagementPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleManagementPoliciesById provides operations to manage the roleManagementPolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) RoleManagementPoliciesById(id string)(*i15708f9856db4407ccc78f2397a412ca44ce57d9fbe9dd3e17f7bf4747b8cf71.UnifiedRoleManagementPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleManagementPolicy%2Did"] = id
    }
    return i15708f9856db4407ccc78f2397a412ca44ce57d9fbe9dd3e17f7bf4747b8cf71.NewUnifiedRoleManagementPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleManagementPolicyAssignments provides operations to manage the roleManagementPolicyAssignments property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) RoleManagementPolicyAssignments()(*i15eb5011e1661e803ff4682ed565e8609c250e7ec6aaf6510e8d4cb37c79c25f.RoleManagementPolicyAssignmentsRequestBuilder) {
    return i15eb5011e1661e803ff4682ed565e8609c250e7ec6aaf6510e8d4cb37c79c25f.NewRoleManagementPolicyAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleManagementPolicyAssignmentsById provides operations to manage the roleManagementPolicyAssignments property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) RoleManagementPolicyAssignmentsById(id string)(*ifdb3aa1e2fe0663a1e8072ad2268b1766bb6cee12e54afb4c9bfe76480e4184a.UnifiedRoleManagementPolicyAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleManagementPolicyAssignment%2Did"] = id
    }
    return ifdb3aa1e2fe0663a1e8072ad2268b1766bb6cee12e54afb4c9bfe76480e4184a.NewUnifiedRoleManagementPolicyAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TokenIssuancePolicies provides operations to manage the tokenIssuancePolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) TokenIssuancePolicies()(*i8cf0127b53d26923e1649ec43c875bc6f0f70ece0a4b9d3cf3377d9e0b4720cf.TokenIssuancePoliciesRequestBuilder) {
    return i8cf0127b53d26923e1649ec43c875bc6f0f70ece0a4b9d3cf3377d9e0b4720cf.NewTokenIssuancePoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TokenIssuancePoliciesById provides operations to manage the tokenIssuancePolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) TokenIssuancePoliciesById(id string)(*i77b5f2108467d08d93443762e93f7967aca61b1916748814eb2ddc49bc0fb0ed.TokenIssuancePolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["tokenIssuancePolicy%2Did"] = id
    }
    return i77b5f2108467d08d93443762e93f7967aca61b1916748814eb2ddc49bc0fb0ed.NewTokenIssuancePolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TokenLifetimePolicies provides operations to manage the tokenLifetimePolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) TokenLifetimePolicies()(*iebb445e03b5a1725abc955bb63d1192b435aa40b5fee4976b801435c3fc9ec06.TokenLifetimePoliciesRequestBuilder) {
    return iebb445e03b5a1725abc955bb63d1192b435aa40b5fee4976b801435c3fc9ec06.NewTokenLifetimePoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TokenLifetimePoliciesById provides operations to manage the tokenLifetimePolicies property of the microsoft.graph.policyRoot entity.
func (m *PoliciesRequestBuilder) TokenLifetimePoliciesById(id string)(*i569717b3ca8408936df0d0a4e253024a8a7b041522350cb7c8977c2c62a9b9d6.TokenLifetimePolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["tokenLifetimePolicy%2Did"] = id
    }
    return i569717b3ca8408936df0d0a4e253024a8a7b041522350cb7c8977c2c62a9b9d6.NewTokenLifetimePolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
