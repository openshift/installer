package permissiongrants

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i210b5891202f914911898d53aea1de0175a6d9a7339fedd712be86e69f61d0a8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/permissiongrants/getavailableextensionproperties"
    i75b4656732491cc699ae2839299632804556c722ec2e0aeff4d4ce49b7d174c5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/permissiongrants/getbyids"
    i7d58d163b6cd2b0c9e4028e1172c49a5605cd2aa9d321bba95295cd9b0dbaedc "github.com/microsoftgraph/msgraph-sdk-go/groups/item/permissiongrants/count"
    id34941651d03077081aceaccd59b616bb77a4640c063a17f50ecac1d42881f20 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/permissiongrants/validateproperties"
)

// PermissionGrantsRequestBuilder provides operations to manage the permissionGrants property of the microsoft.graph.group entity.
type PermissionGrantsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PermissionGrantsRequestBuilderGetQueryParameters list all resource-specific permission grants on the group. This list specifies the Azure AD apps that have access to the **group**, along with the corresponding kind of resource-specific access that each app has.
type PermissionGrantsRequestBuilderGetQueryParameters struct {
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
// PermissionGrantsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PermissionGrantsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PermissionGrantsRequestBuilderGetQueryParameters
}
// PermissionGrantsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PermissionGrantsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewPermissionGrantsRequestBuilderInternal instantiates a new PermissionGrantsRequestBuilder and sets the default values.
func NewPermissionGrantsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PermissionGrantsRequestBuilder) {
    m := &PermissionGrantsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/permissionGrants{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPermissionGrantsRequestBuilder instantiates a new PermissionGrantsRequestBuilder and sets the default values.
func NewPermissionGrantsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PermissionGrantsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPermissionGrantsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *PermissionGrantsRequestBuilder) Count()(*i7d58d163b6cd2b0c9e4028e1172c49a5605cd2aa9d321bba95295cd9b0dbaedc.CountRequestBuilder) {
    return i7d58d163b6cd2b0c9e4028e1172c49a5605cd2aa9d321bba95295cd9b0dbaedc.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation list all resource-specific permission grants on the group. This list specifies the Azure AD apps that have access to the **group**, along with the corresponding kind of resource-specific access that each app has.
func (m *PermissionGrantsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PermissionGrantsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create new navigation property to permissionGrants for groups
func (m *PermissionGrantsRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantable, requestConfiguration *PermissionGrantsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get list all resource-specific permission grants on the group. This list specifies the Azure AD apps that have access to the **group**, along with the corresponding kind of resource-specific access that each app has.
func (m *PermissionGrantsRequestBuilder) Get(ctx context.Context, requestConfiguration *PermissionGrantsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateResourceSpecificPermissionGrantCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantCollectionResponseable), nil
}
// GetAvailableExtensionProperties provides operations to call the getAvailableExtensionProperties method.
func (m *PermissionGrantsRequestBuilder) GetAvailableExtensionProperties()(*i210b5891202f914911898d53aea1de0175a6d9a7339fedd712be86e69f61d0a8.GetAvailableExtensionPropertiesRequestBuilder) {
    return i210b5891202f914911898d53aea1de0175a6d9a7339fedd712be86e69f61d0a8.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *PermissionGrantsRequestBuilder) GetByIds()(*i75b4656732491cc699ae2839299632804556c722ec2e0aeff4d4ce49b7d174c5.GetByIdsRequestBuilder) {
    return i75b4656732491cc699ae2839299632804556c722ec2e0aeff4d4ce49b7d174c5.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create new navigation property to permissionGrants for groups
func (m *PermissionGrantsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantable, requestConfiguration *PermissionGrantsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateResourceSpecificPermissionGrantFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ResourceSpecificPermissionGrantable), nil
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *PermissionGrantsRequestBuilder) ValidateProperties()(*id34941651d03077081aceaccd59b616bb77a4640c063a17f50ecac1d42881f20.ValidatePropertiesRequestBuilder) {
    return id34941651d03077081aceaccd59b616bb77a4640c063a17f50ecac1d42881f20.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
