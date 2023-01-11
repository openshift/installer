package externalsponsors

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i41016cceeff842ed9a006c661821c6b48641a95c280b56f894c38b38b4f338a8 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item/externalsponsors/validateproperties"
    i4efeeddbbfa57da208dfa7b863ad9067cfc0c816854dd71cb5a6ca38c8a184b6 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item/externalsponsors/getbyids"
    i5780e326fb3af6db2938c8ca9c539e53f905ef00146f0fb0131233caaa460917 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item/externalsponsors/getavailableextensionproperties"
    ia251373d77add160c436f9217f5f9433a4c524c145be9e94c9cb03181c9a1c45 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item/externalsponsors/count"
    id167466353efea4b0c493edfe4558f5b2e21849a0a97e4983e926a18deca523b "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item/externalsponsors/ref"
)

// ExternalSponsorsRequestBuilder provides operations to manage the externalSponsors property of the microsoft.graph.connectedOrganization entity.
type ExternalSponsorsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ExternalSponsorsRequestBuilderGetQueryParameters retrieve a list of a connectedOrganization's external sponsors.  The external sponsors are a set of users who can approve requests on behalf of other users from that connected organization.
type ExternalSponsorsRequestBuilderGetQueryParameters struct {
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
// ExternalSponsorsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ExternalSponsorsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ExternalSponsorsRequestBuilderGetQueryParameters
}
// ExternalSponsorsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ExternalSponsorsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewExternalSponsorsRequestBuilderInternal instantiates a new ExternalSponsorsRequestBuilder and sets the default values.
func NewExternalSponsorsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ExternalSponsorsRequestBuilder) {
    m := &ExternalSponsorsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/identityGovernance/entitlementManagement/connectedOrganizations/{connectedOrganization%2Did}/externalSponsors{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewExternalSponsorsRequestBuilder instantiates a new ExternalSponsorsRequestBuilder and sets the default values.
func NewExternalSponsorsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ExternalSponsorsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewExternalSponsorsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *ExternalSponsorsRequestBuilder) Count()(*ia251373d77add160c436f9217f5f9433a4c524c145be9e94c9cb03181c9a1c45.CountRequestBuilder) {
    return ia251373d77add160c436f9217f5f9433a4c524c145be9e94c9cb03181c9a1c45.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation retrieve a list of a connectedOrganization's external sponsors.  The external sponsors are a set of users who can approve requests on behalf of other users from that connected organization.
func (m *ExternalSponsorsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ExternalSponsorsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create new navigation property to externalSponsors for identityGovernance
func (m *ExternalSponsorsRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *ExternalSponsorsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Get retrieve a list of a connectedOrganization's external sponsors.  The external sponsors are a set of users who can approve requests on behalf of other users from that connected organization.
func (m *ExternalSponsorsRequestBuilder) Get(ctx context.Context, requestConfiguration *ExternalSponsorsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
// GetAvailableExtensionProperties provides operations to call the getAvailableExtensionProperties method.
func (m *ExternalSponsorsRequestBuilder) GetAvailableExtensionProperties()(*i5780e326fb3af6db2938c8ca9c539e53f905ef00146f0fb0131233caaa460917.GetAvailableExtensionPropertiesRequestBuilder) {
    return i5780e326fb3af6db2938c8ca9c539e53f905ef00146f0fb0131233caaa460917.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *ExternalSponsorsRequestBuilder) GetByIds()(*i4efeeddbbfa57da208dfa7b863ad9067cfc0c816854dd71cb5a6ca38c8a184b6.GetByIdsRequestBuilder) {
    return i4efeeddbbfa57da208dfa7b863ad9067cfc0c816854dd71cb5a6ca38c8a184b6.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create new navigation property to externalSponsors for identityGovernance
func (m *ExternalSponsorsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *ExternalSponsorsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable), nil
}
// Ref provides operations to manage the collection of identityGovernance entities.
func (m *ExternalSponsorsRequestBuilder) Ref()(*id167466353efea4b0c493edfe4558f5b2e21849a0a97e4983e926a18deca523b.RefRequestBuilder) {
    return id167466353efea4b0c493edfe4558f5b2e21849a0a97e4983e926a18deca523b.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *ExternalSponsorsRequestBuilder) ValidateProperties()(*i41016cceeff842ed9a006c661821c6b48641a95c280b56f894c38b38b4f338a8.ValidatePropertiesRequestBuilder) {
    return i41016cceeff842ed9a006c661821c6b48641a95c280b56f894c38b38b4f338a8.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
