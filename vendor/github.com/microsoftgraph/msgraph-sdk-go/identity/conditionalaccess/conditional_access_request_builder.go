package conditionalaccess

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i031da455aa638a5f89c85fdb0d3e8882b40650ad3a1e93d7b33885e282048fa5 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/authenticationcontextclassreferences"
    i27d3a1bae6f4577a4de2876ba58fad41be78be972958441916bfc0b15d805933 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/policies"
    i6bf5a84cd951bc7d61e669aafc7a4d48a4738a1af268d702686266e96e8d6a44 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/templates"
    ie167212f1f78fb7bf5e6a91e4a3df660b8a5499daeed0e885580a0dde04d710e "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/namedlocations"
    i574b61c7cf81bce24f5146815b9c92a366c8a4de82277e47a1fc7817eec71339 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/templates/item"
    i652ce65e13fdceb3d1ed10efc4a8f6c0d893a47d549accb475d48e05eedf3251 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/namedlocations/item"
    i91801d979eb0028b76d066e7195d63b5ed00d456ccdceaa9f6bb9a091401a467 "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/policies/item"
    ie664735c68cf367c24c6ad8cc3ee7e735baa549391a4d0ba0d619f7bf748294d "github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/authenticationcontextclassreferences/item"
)

// ConditionalAccessRequestBuilder provides operations to manage the conditionalAccess property of the microsoft.graph.identityContainer entity.
type ConditionalAccessRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ConditionalAccessRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ConditionalAccessRequestBuilderGetQueryParameters the entry point for the Conditional Access (CA) object model.
type ConditionalAccessRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ConditionalAccessRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ConditionalAccessRequestBuilderGetQueryParameters
}
// ConditionalAccessRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConditionalAccessRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AuthenticationContextClassReferences provides operations to manage the authenticationContextClassReferences property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) AuthenticationContextClassReferences()(*i031da455aa638a5f89c85fdb0d3e8882b40650ad3a1e93d7b33885e282048fa5.AuthenticationContextClassReferencesRequestBuilder) {
    return i031da455aa638a5f89c85fdb0d3e8882b40650ad3a1e93d7b33885e282048fa5.NewAuthenticationContextClassReferencesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthenticationContextClassReferencesById provides operations to manage the authenticationContextClassReferences property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) AuthenticationContextClassReferencesById(id string)(*ie664735c68cf367c24c6ad8cc3ee7e735baa549391a4d0ba0d619f7bf748294d.AuthenticationContextClassReferenceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["authenticationContextClassReference%2Did"] = id
    }
    return ie664735c68cf367c24c6ad8cc3ee7e735baa549391a4d0ba0d619f7bf748294d.NewAuthenticationContextClassReferenceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewConditionalAccessRequestBuilderInternal instantiates a new ConditionalAccessRequestBuilder and sets the default values.
func NewConditionalAccessRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessRequestBuilder) {
    m := &ConditionalAccessRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/identity/conditionalAccess{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewConditionalAccessRequestBuilder instantiates a new ConditionalAccessRequestBuilder and sets the default values.
func NewConditionalAccessRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConditionalAccessRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewConditionalAccessRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property conditionalAccess for identity
func (m *ConditionalAccessRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ConditionalAccessRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the entry point for the Conditional Access (CA) object model.
func (m *ConditionalAccessRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ConditionalAccessRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property conditionalAccess in identity
func (m *ConditionalAccessRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable, requestConfiguration *ConditionalAccessRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property conditionalAccess for identity
func (m *ConditionalAccessRequestBuilder) Delete(ctx context.Context, requestConfiguration *ConditionalAccessRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the entry point for the Conditional Access (CA) object model.
func (m *ConditionalAccessRequestBuilder) Get(ctx context.Context, requestConfiguration *ConditionalAccessRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateConditionalAccessRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable), nil
}
// NamedLocations provides operations to manage the namedLocations property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) NamedLocations()(*ie167212f1f78fb7bf5e6a91e4a3df660b8a5499daeed0e885580a0dde04d710e.NamedLocationsRequestBuilder) {
    return ie167212f1f78fb7bf5e6a91e4a3df660b8a5499daeed0e885580a0dde04d710e.NewNamedLocationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NamedLocationsById provides operations to manage the namedLocations property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) NamedLocationsById(id string)(*i652ce65e13fdceb3d1ed10efc4a8f6c0d893a47d549accb475d48e05eedf3251.NamedLocationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["namedLocation%2Did"] = id
    }
    return i652ce65e13fdceb3d1ed10efc4a8f6c0d893a47d549accb475d48e05eedf3251.NewNamedLocationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property conditionalAccess in identity
func (m *ConditionalAccessRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable, requestConfiguration *ConditionalAccessRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateConditionalAccessRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ConditionalAccessRootable), nil
}
// Policies provides operations to manage the policies property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) Policies()(*i27d3a1bae6f4577a4de2876ba58fad41be78be972958441916bfc0b15d805933.PoliciesRequestBuilder) {
    return i27d3a1bae6f4577a4de2876ba58fad41be78be972958441916bfc0b15d805933.NewPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PoliciesById provides operations to manage the policies property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) PoliciesById(id string)(*i91801d979eb0028b76d066e7195d63b5ed00d456ccdceaa9f6bb9a091401a467.ConditionalAccessPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conditionalAccessPolicy%2Did"] = id
    }
    return i91801d979eb0028b76d066e7195d63b5ed00d456ccdceaa9f6bb9a091401a467.NewConditionalAccessPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Templates provides operations to manage the templates property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) Templates()(*i6bf5a84cd951bc7d61e669aafc7a4d48a4738a1af268d702686266e96e8d6a44.TemplatesRequestBuilder) {
    return i6bf5a84cd951bc7d61e669aafc7a4d48a4738a1af268d702686266e96e8d6a44.NewTemplatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TemplatesById provides operations to manage the templates property of the microsoft.graph.conditionalAccessRoot entity.
func (m *ConditionalAccessRequestBuilder) TemplatesById(id string)(*i574b61c7cf81bce24f5146815b9c92a366c8a4de82277e47a1fc7817eec71339.ConditionalAccessTemplateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conditionalAccessTemplate%2Did"] = id
    }
    return i574b61c7cf81bce24f5146815b9c92a366c8a4de82277e47a1fc7817eec71339.NewConditionalAccessTemplateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
