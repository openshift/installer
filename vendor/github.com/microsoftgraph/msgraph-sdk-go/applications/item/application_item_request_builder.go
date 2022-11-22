package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i02dc3442cb07e14dd6564850055582daa12c5271c29402a2eacef3296ce0336a "github.com/microsoftgraph/msgraph-sdk-go/applications/item/federatedidentitycredentials"
    i0b84a273ded1e6669385f95709d9a7cf84a87df5fd73dc8c7af3710571de2128 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/unsetverifiedpublisher"
    i1a8360833e23698ad94af084d2df35cce4a1972916936379ed6648bfa6ba57ec "github.com/microsoftgraph/msgraph-sdk-go/applications/item/logo"
    i1afcec2f462ce9653f6d9c178d0d0542f5684412e62323ab5a3f979e1b79b5b8 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenissuancepolicies"
    i1cf4d6bdfc37eea0b8dc06a0fbf17ea8c1b58d8dc24649421cd74e454d266909 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/addkey"
    i402e9b9506b4f160b32ff917d1b2fe43f4d4615296648051dc3eb89140278233 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/removekey"
    i486b52d4dad87ceeac08bd843c5504a58e010e9931346242f546e8106ff87250 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenlifetimepolicies"
    i5a085ba4baea1613f7766106526258f582e52faf4f832b1ce2655815416aa722 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/removepassword"
    i6f17a12e2de6cae814f546d592ddc360f0666722e94dabfd2392eaa13b36740a "github.com/microsoftgraph/msgraph-sdk-go/applications/item/setverifiedpublisher"
    i72720b8b61b527b526dacbdb4f1b2ebc44378e69b3a31f054e52851fc51aee5a "github.com/microsoftgraph/msgraph-sdk-go/applications/item/checkmembergroups"
    i8644212e8e36b861035674e39cad93b1868aca40e9292c7019cd29ae25503583 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/createdonbehalfof"
    i8f079434e6086cc56be6ed0afce670e9aa0f7ff3ac0107e2f5d9ac8d358fe2dc "github.com/microsoftgraph/msgraph-sdk-go/applications/item/getmembergroups"
    i9565e5e51a86644270ab0a3f267a0bbce14c338625e43342715917d0c6d5eec7 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners"
    i9dc2689198e150b8c4d03044f025827c3242c48569acd497be6fec72d51f6797 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/homerealmdiscoverypolicies"
    ia607556161094bfdb9a3407887dfa9b53a9d7a0857a0fd091dab27dd31a9459a "github.com/microsoftgraph/msgraph-sdk-go/applications/item/checkmemberobjects"
    iafba803a7908a654425c76ab1e010310bb6222243e602ddaab50bc32848dac3c "github.com/microsoftgraph/msgraph-sdk-go/applications/item/getmemberobjects"
    ib5b5951df8af6837ed6ea3344f5688d5f2483d29ea1c165d2d219ec5edf941da "github.com/microsoftgraph/msgraph-sdk-go/applications/item/restore"
    id3ff5385c48c1f8ddd4e5f53452e89d8e30d638f1b87dc3265bf1f8adb5d3078 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/addpassword"
    ida306ab63547733b6e18148e7f72a169be3442196d415312a52cbd9978ab6961 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/extensionproperties"
    i07ab470de2cf9d615c15406822a9ab9154633fa9ed5641ce1d97bbfbb1863560 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/federatedidentitycredentials/item"
    i2af67e5b951d866dcb390e8cdd53a703b8d1613ab54b5b0b111652d63e7426e3 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/homerealmdiscoverypolicies/item"
    i97c12a9c6c5e3dabf4cee7a203250908050ed917a949921ccf7e4f24cf095d25 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenlifetimepolicies/item"
    id5f0be733c94801a6e075b5b4e84614f608f9c7a8ece0f49a59084244eff9f56 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item"
    id868d42838b59366f189543d6b65b78e527af17acc61cbc962cae87c3c961e40 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/tokenissuancepolicies/item"
    ie9db136445b8a8e3d0b048cc9910b564635f65b63edab67c97c083c411010b77 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/extensionproperties/item"
)

// ApplicationItemRequestBuilder provides operations to manage the collection of application entities.
type ApplicationItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ApplicationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ApplicationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ApplicationItemRequestBuilderGetQueryParameters get the properties and relationships of an application object.
type ApplicationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ApplicationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ApplicationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ApplicationItemRequestBuilderGetQueryParameters
}
// ApplicationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ApplicationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AddKey provides operations to call the addKey method.
func (m *ApplicationItemRequestBuilder) AddKey()(*i1cf4d6bdfc37eea0b8dc06a0fbf17ea8c1b58d8dc24649421cd74e454d266909.AddKeyRequestBuilder) {
    return i1cf4d6bdfc37eea0b8dc06a0fbf17ea8c1b58d8dc24649421cd74e454d266909.NewAddKeyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AddPassword provides operations to call the addPassword method.
func (m *ApplicationItemRequestBuilder) AddPassword()(*id3ff5385c48c1f8ddd4e5f53452e89d8e30d638f1b87dc3265bf1f8adb5d3078.AddPasswordRequestBuilder) {
    return id3ff5385c48c1f8ddd4e5f53452e89d8e30d638f1b87dc3265bf1f8adb5d3078.NewAddPasswordRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
func (m *ApplicationItemRequestBuilder) CheckMemberGroups()(*i72720b8b61b527b526dacbdb4f1b2ebc44378e69b3a31f054e52851fc51aee5a.CheckMemberGroupsRequestBuilder) {
    return i72720b8b61b527b526dacbdb4f1b2ebc44378e69b3a31f054e52851fc51aee5a.NewCheckMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
func (m *ApplicationItemRequestBuilder) CheckMemberObjects()(*ia607556161094bfdb9a3407887dfa9b53a9d7a0857a0fd091dab27dd31a9459a.CheckMemberObjectsRequestBuilder) {
    return ia607556161094bfdb9a3407887dfa9b53a9d7a0857a0fd091dab27dd31a9459a.NewCheckMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewApplicationItemRequestBuilderInternal instantiates a new ApplicationItemRequestBuilder and sets the default values.
func NewApplicationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ApplicationItemRequestBuilder) {
    m := &ApplicationItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewApplicationItemRequestBuilder instantiates a new ApplicationItemRequestBuilder and sets the default values.
func NewApplicationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ApplicationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewApplicationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete an application object. When deleted, apps are moved to a temporary container and can be restored within 30 days. After that time, they are permanently deleted.
func (m *ApplicationItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ApplicationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatedOnBehalfOf provides operations to manage the createdOnBehalfOf property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) CreatedOnBehalfOf()(*i8644212e8e36b861035674e39cad93b1868aca40e9292c7019cd29ae25503583.CreatedOnBehalfOfRequestBuilder) {
    return i8644212e8e36b861035674e39cad93b1868aca40e9292c7019cd29ae25503583.NewCreatedOnBehalfOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get the properties and relationships of an application object.
func (m *ApplicationItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ApplicationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the properties of an application object.
func (m *ApplicationItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, requestConfiguration *ApplicationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete an application object. When deleted, apps are moved to a temporary container and can be restored within 30 days. After that time, they are permanently deleted.
func (m *ApplicationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ApplicationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// ExtensionProperties provides operations to manage the extensionProperties property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) ExtensionProperties()(*ida306ab63547733b6e18148e7f72a169be3442196d415312a52cbd9978ab6961.ExtensionPropertiesRequestBuilder) {
    return ida306ab63547733b6e18148e7f72a169be3442196d415312a52cbd9978ab6961.NewExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionPropertiesById provides operations to manage the extensionProperties property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) ExtensionPropertiesById(id string)(*ie9db136445b8a8e3d0b048cc9910b564635f65b63edab67c97c083c411010b77.ExtensionPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extensionProperty%2Did"] = id
    }
    return ie9db136445b8a8e3d0b048cc9910b564635f65b63edab67c97c083c411010b77.NewExtensionPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// FederatedIdentityCredentials provides operations to manage the federatedIdentityCredentials property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) FederatedIdentityCredentials()(*i02dc3442cb07e14dd6564850055582daa12c5271c29402a2eacef3296ce0336a.FederatedIdentityCredentialsRequestBuilder) {
    return i02dc3442cb07e14dd6564850055582daa12c5271c29402a2eacef3296ce0336a.NewFederatedIdentityCredentialsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FederatedIdentityCredentialsById provides operations to manage the federatedIdentityCredentials property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) FederatedIdentityCredentialsById(id string)(*i07ab470de2cf9d615c15406822a9ab9154633fa9ed5641ce1d97bbfbb1863560.FederatedIdentityCredentialItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["federatedIdentityCredential%2Did"] = id
    }
    return i07ab470de2cf9d615c15406822a9ab9154633fa9ed5641ce1d97bbfbb1863560.NewFederatedIdentityCredentialItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get get the properties and relationships of an application object.
func (m *ApplicationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ApplicationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable), nil
}
// GetMemberGroups provides operations to call the getMemberGroups method.
func (m *ApplicationItemRequestBuilder) GetMemberGroups()(*i8f079434e6086cc56be6ed0afce670e9aa0f7ff3ac0107e2f5d9ac8d358fe2dc.GetMemberGroupsRequestBuilder) {
    return i8f079434e6086cc56be6ed0afce670e9aa0f7ff3ac0107e2f5d9ac8d358fe2dc.NewGetMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetMemberObjects provides operations to call the getMemberObjects method.
func (m *ApplicationItemRequestBuilder) GetMemberObjects()(*iafba803a7908a654425c76ab1e010310bb6222243e602ddaab50bc32848dac3c.GetMemberObjectsRequestBuilder) {
    return iafba803a7908a654425c76ab1e010310bb6222243e602ddaab50bc32848dac3c.NewGetMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// HomeRealmDiscoveryPolicies provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) HomeRealmDiscoveryPolicies()(*i9dc2689198e150b8c4d03044f025827c3242c48569acd497be6fec72d51f6797.HomeRealmDiscoveryPoliciesRequestBuilder) {
    return i9dc2689198e150b8c4d03044f025827c3242c48569acd497be6fec72d51f6797.NewHomeRealmDiscoveryPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// HomeRealmDiscoveryPoliciesById provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) HomeRealmDiscoveryPoliciesById(id string)(*i2af67e5b951d866dcb390e8cdd53a703b8d1613ab54b5b0b111652d63e7426e3.HomeRealmDiscoveryPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["homeRealmDiscoveryPolicy%2Did"] = id
    }
    return i2af67e5b951d866dcb390e8cdd53a703b8d1613ab54b5b0b111652d63e7426e3.NewHomeRealmDiscoveryPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Logo provides operations to manage the media for the application entity.
func (m *ApplicationItemRequestBuilder) Logo()(*i1a8360833e23698ad94af084d2df35cce4a1972916936379ed6648bfa6ba57ec.LogoRequestBuilder) {
    return i1a8360833e23698ad94af084d2df35cce4a1972916936379ed6648bfa6ba57ec.NewLogoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Owners provides operations to manage the owners property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) Owners()(*i9565e5e51a86644270ab0a3f267a0bbce14c338625e43342715917d0c6d5eec7.OwnersRequestBuilder) {
    return i9565e5e51a86644270ab0a3f267a0bbce14c338625e43342715917d0c6d5eec7.NewOwnersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OwnersById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.applications.item.owners.item collection
func (m *ApplicationItemRequestBuilder) OwnersById(id string)(*id5f0be733c94801a6e075b5b4e84614f608f9c7a8ece0f49a59084244eff9f56.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return id5f0be733c94801a6e075b5b4e84614f608f9c7a8ece0f49a59084244eff9f56.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the properties of an application object.
func (m *ApplicationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, requestConfiguration *ApplicationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable), nil
}
// RemoveKey provides operations to call the removeKey method.
func (m *ApplicationItemRequestBuilder) RemoveKey()(*i402e9b9506b4f160b32ff917d1b2fe43f4d4615296648051dc3eb89140278233.RemoveKeyRequestBuilder) {
    return i402e9b9506b4f160b32ff917d1b2fe43f4d4615296648051dc3eb89140278233.NewRemoveKeyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemovePassword provides operations to call the removePassword method.
func (m *ApplicationItemRequestBuilder) RemovePassword()(*i5a085ba4baea1613f7766106526258f582e52faf4f832b1ce2655815416aa722.RemovePasswordRequestBuilder) {
    return i5a085ba4baea1613f7766106526258f582e52faf4f832b1ce2655815416aa722.NewRemovePasswordRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *ApplicationItemRequestBuilder) Restore()(*ib5b5951df8af6837ed6ea3344f5688d5f2483d29ea1c165d2d219ec5edf941da.RestoreRequestBuilder) {
    return ib5b5951df8af6837ed6ea3344f5688d5f2483d29ea1c165d2d219ec5edf941da.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetVerifiedPublisher provides operations to call the setVerifiedPublisher method.
func (m *ApplicationItemRequestBuilder) SetVerifiedPublisher()(*i6f17a12e2de6cae814f546d592ddc360f0666722e94dabfd2392eaa13b36740a.SetVerifiedPublisherRequestBuilder) {
    return i6f17a12e2de6cae814f546d592ddc360f0666722e94dabfd2392eaa13b36740a.NewSetVerifiedPublisherRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TokenIssuancePolicies provides operations to manage the tokenIssuancePolicies property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) TokenIssuancePolicies()(*i1afcec2f462ce9653f6d9c178d0d0542f5684412e62323ab5a3f979e1b79b5b8.TokenIssuancePoliciesRequestBuilder) {
    return i1afcec2f462ce9653f6d9c178d0d0542f5684412e62323ab5a3f979e1b79b5b8.NewTokenIssuancePoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TokenIssuancePoliciesById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.applications.item.tokenIssuancePolicies.item collection
func (m *ApplicationItemRequestBuilder) TokenIssuancePoliciesById(id string)(*id868d42838b59366f189543d6b65b78e527af17acc61cbc962cae87c3c961e40.TokenIssuancePolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["tokenIssuancePolicy%2Did"] = id
    }
    return id868d42838b59366f189543d6b65b78e527af17acc61cbc962cae87c3c961e40.NewTokenIssuancePolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TokenLifetimePolicies provides operations to manage the tokenLifetimePolicies property of the microsoft.graph.application entity.
func (m *ApplicationItemRequestBuilder) TokenLifetimePolicies()(*i486b52d4dad87ceeac08bd843c5504a58e010e9931346242f546e8106ff87250.TokenLifetimePoliciesRequestBuilder) {
    return i486b52d4dad87ceeac08bd843c5504a58e010e9931346242f546e8106ff87250.NewTokenLifetimePoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TokenLifetimePoliciesById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.applications.item.tokenLifetimePolicies.item collection
func (m *ApplicationItemRequestBuilder) TokenLifetimePoliciesById(id string)(*i97c12a9c6c5e3dabf4cee7a203250908050ed917a949921ccf7e4f24cf095d25.TokenLifetimePolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["tokenLifetimePolicy%2Did"] = id
    }
    return i97c12a9c6c5e3dabf4cee7a203250908050ed917a949921ccf7e4f24cf095d25.NewTokenLifetimePolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UnsetVerifiedPublisher provides operations to call the unsetVerifiedPublisher method.
func (m *ApplicationItemRequestBuilder) UnsetVerifiedPublisher()(*i0b84a273ded1e6669385f95709d9a7cf84a87df5fd73dc8c7af3710571de2128.UnsetVerifiedPublisherRequestBuilder) {
    return i0b84a273ded1e6669385f95709d9a7cf84a87df5fd73dc8c7af3710571de2128.NewUnsetVerifiedPublisherRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
