package deleteditems

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i334d8fedfd505c81d9078410e8949db770559d503313bdc074b7c3aa3d1eb225 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/group"
    i3d5ada12c07368d2eb27040bf4a7e70eba678ee6db4ee1f0ac93aa0a6f7d60c0 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/getbyids"
    icba28028f1df1a97032518c6af5fa4b34f37fee861857d5d9ace496a2421a662 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/getavailableextensionproperties"
    id560e6c50a2ad80dbc730a694f2b7cf69b16304069e2f56be17a61dc1f271f23 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/user"
    id64a33d4a1456f0a7c3b2dc31e4984ae757e62be968137695545f0dd4f5be2c2 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/count"
    id8f04b482d6f4071942a84bb0fbf25521758c3efd03dbf5c599ed4766bd33c07 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/application"
    ie70b17c2909fb54ff8c49ac02a9dcc9a2528984c5ea5deb722364e3745c1eb69 "github.com/microsoftgraph/msgraph-sdk-go/directory/deleteditems/validateproperties"
)

// DeletedItemsRequestBuilder provides operations to manage the deletedItems property of the microsoft.graph.directory entity.
type DeletedItemsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DeletedItemsRequestBuilderGetQueryParameters recently deleted items. Read-only. Nullable.
type DeletedItemsRequestBuilderGetQueryParameters struct {
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
// DeletedItemsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeletedItemsRequestBuilderGetQueryParameters
}
// DeletedItemsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedItemsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Application casts the previous resource to application.
func (m *DeletedItemsRequestBuilder) Application()(*id8f04b482d6f4071942a84bb0fbf25521758c3efd03dbf5c599ed4766bd33c07.ApplicationRequestBuilder) {
    return id8f04b482d6f4071942a84bb0fbf25521758c3efd03dbf5c599ed4766bd33c07.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDeletedItemsRequestBuilderInternal instantiates a new DeletedItemsRequestBuilder and sets the default values.
func NewDeletedItemsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsRequestBuilder) {
    m := &DeletedItemsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directory/deletedItems{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDeletedItemsRequestBuilder instantiates a new DeletedItemsRequestBuilder and sets the default values.
func NewDeletedItemsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedItemsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedItemsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *DeletedItemsRequestBuilder) Count()(*id64a33d4a1456f0a7c3b2dc31e4984ae757e62be968137695545f0dd4f5be2c2.CountRequestBuilder) {
    return id64a33d4a1456f0a7c3b2dc31e4984ae757e62be968137695545f0dd4f5be2c2.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation recently deleted items. Read-only. Nullable.
func (m *DeletedItemsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DeletedItemsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create new navigation property to deletedItems for directory
func (m *DeletedItemsRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *DeletedItemsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get recently deleted items. Read-only. Nullable.
func (m *DeletedItemsRequestBuilder) Get(ctx context.Context, requestConfiguration *DeletedItemsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *DeletedItemsRequestBuilder) GetAvailableExtensionProperties()(*icba28028f1df1a97032518c6af5fa4b34f37fee861857d5d9ace496a2421a662.GetAvailableExtensionPropertiesRequestBuilder) {
    return icba28028f1df1a97032518c6af5fa4b34f37fee861857d5d9ace496a2421a662.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *DeletedItemsRequestBuilder) GetByIds()(*i3d5ada12c07368d2eb27040bf4a7e70eba678ee6db4ee1f0ac93aa0a6f7d60c0.GetByIdsRequestBuilder) {
    return i3d5ada12c07368d2eb27040bf4a7e70eba678ee6db4ee1f0ac93aa0a6f7d60c0.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Group casts the previous resource to group.
func (m *DeletedItemsRequestBuilder) Group()(*i334d8fedfd505c81d9078410e8949db770559d503313bdc074b7c3aa3d1eb225.GroupRequestBuilder) {
    return i334d8fedfd505c81d9078410e8949db770559d503313bdc074b7c3aa3d1eb225.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create new navigation property to deletedItems for directory
func (m *DeletedItemsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, requestConfiguration *DeletedItemsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
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
// User casts the previous resource to user.
func (m *DeletedItemsRequestBuilder) User()(*id560e6c50a2ad80dbc730a694f2b7cf69b16304069e2f56be17a61dc1f271f23.UserRequestBuilder) {
    return id560e6c50a2ad80dbc730a694f2b7cf69b16304069e2f56be17a61dc1f271f23.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *DeletedItemsRequestBuilder) ValidateProperties()(*ie70b17c2909fb54ff8c49ac02a9dcc9a2528984c5ea5deb722364e3745c1eb69.ValidatePropertiesRequestBuilder) {
    return ie70b17c2909fb54ff8c49ac02a9dcc9a2528984c5ea5deb722364e3745c1eb69.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
