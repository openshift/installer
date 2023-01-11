package appliesto

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i01b7156a2542039015b03b4c145895442723775d50760476d6a858e7a38c0229 "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item/appliesto/getavailableextensionproperties"
    i291ff03133de15263bbbb9161ec52a9fd3780244fb39767c57a0d3b944a877bb "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item/appliesto/count"
    i5fcc1cbba281bf089e86d1f0fd053b72254efdc576b56f07f5f474773e8f4e0b "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item/appliesto/getbyids"
    i9be808779624fafbf1ab48016999de8af7668ec9c3d03a9ba21e6fa0fb30642c "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item/appliesto/ref"
    ibb1e01f7e6a4c02e7b91f4bf4b3d9806fdaa9391df8c9da069d523afca17ae59 "github.com/microsoftgraph/msgraph-sdk-go/policies/featurerolloutpolicies/item/appliesto/validateproperties"
)

// AppliesToRequestBuilder provides operations to manage the appliesTo property of the microsoft.graph.featureRolloutPolicy entity.
type AppliesToRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppliesToRequestBuilderGetQueryParameters nullable. Specifies a list of directoryObjects that feature is enabled for.
type AppliesToRequestBuilderGetQueryParameters struct {
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
// AppliesToRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AppliesToRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AppliesToRequestBuilderGetQueryParameters
}
// AppliesToRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AppliesToRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAppliesToRequestBuilderInternal instantiates a new AppliesToRequestBuilder and sets the default values.
func NewAppliesToRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AppliesToRequestBuilder) {
    m := &AppliesToRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/policies/featureRolloutPolicies/{featureRolloutPolicy%2Did}/appliesTo{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewAppliesToRequestBuilder instantiates a new AppliesToRequestBuilder and sets the default values.
func NewAppliesToRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AppliesToRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAppliesToRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *AppliesToRequestBuilder) Count()(*i291ff03133de15263bbbb9161ec52a9fd3780244fb39767c57a0d3b944a877bb.CountRequestBuilder) {
    return i291ff03133de15263bbbb9161ec52a9fd3780244fb39767c57a0d3b944a877bb.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation nullable. Specifies a list of directoryObjects that feature is enabled for.
func (m *AppliesToRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *AppliesToRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create new navigation property to appliesTo for policies
func (m *AppliesToRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *AppliesToRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get nullable. Specifies a list of directoryObjects that feature is enabled for.
func (m *AppliesToRequestBuilder) Get(ctx context.Context, requestConfiguration *AppliesToRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *AppliesToRequestBuilder) GetAvailableExtensionProperties()(*i01b7156a2542039015b03b4c145895442723775d50760476d6a858e7a38c0229.GetAvailableExtensionPropertiesRequestBuilder) {
    return i01b7156a2542039015b03b4c145895442723775d50760476d6a858e7a38c0229.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *AppliesToRequestBuilder) GetByIds()(*i5fcc1cbba281bf089e86d1f0fd053b72254efdc576b56f07f5f474773e8f4e0b.GetByIdsRequestBuilder) {
    return i5fcc1cbba281bf089e86d1f0fd053b72254efdc576b56f07f5f474773e8f4e0b.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create new navigation property to appliesTo for policies
func (m *AppliesToRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *AppliesToRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
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
// Ref provides operations to manage the collection of policyRoot entities.
func (m *AppliesToRequestBuilder) Ref()(*i9be808779624fafbf1ab48016999de8af7668ec9c3d03a9ba21e6fa0fb30642c.RefRequestBuilder) {
    return i9be808779624fafbf1ab48016999de8af7668ec9c3d03a9ba21e6fa0fb30642c.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *AppliesToRequestBuilder) ValidateProperties()(*ibb1e01f7e6a4c02e7b91f4bf4b3d9806fdaa9391df8c9da069d523afca17ae59.ValidatePropertiesRequestBuilder) {
    return ibb1e01f7e6a4c02e7b91f4bf4b3d9806fdaa9391df8c9da069d523afca17ae59.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
