package ownedobjects

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i4df0a049a4392d6e4a658bc89720a18c16df1e9e765f88ad93fd42da67e4f3c4 "github.com/microsoftgraph/msgraph-sdk-go/me/ownedobjects/serviceprincipal"
    i6276b1ac50166701f0778eb9a93a719714d4f2d6873e92e5b355f33d4e580448 "github.com/microsoftgraph/msgraph-sdk-go/me/ownedobjects/application"
    ib2a122a45f2f9ecdd4dec70c33255e57920c1fe41fe02b4b01be0195ffea3106 "github.com/microsoftgraph/msgraph-sdk-go/me/ownedobjects/count"
    ibaa1dc10d052b4b8688aceb90d308f0ac89f2de381f3ad51f64596a7fcedab6c "github.com/microsoftgraph/msgraph-sdk-go/me/ownedobjects/group"
)

// OwnedObjectsRequestBuilder provides operations to manage the ownedObjects property of the microsoft.graph.user entity.
type OwnedObjectsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OwnedObjectsRequestBuilderGetQueryParameters directory objects that are owned by the user. Read-only. Nullable. Supports $expand.
type OwnedObjectsRequestBuilderGetQueryParameters struct {
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
// OwnedObjectsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OwnedObjectsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OwnedObjectsRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *OwnedObjectsRequestBuilder) Application()(*i6276b1ac50166701f0778eb9a93a719714d4f2d6873e92e5b355f33d4e580448.ApplicationRequestBuilder) {
    return i6276b1ac50166701f0778eb9a93a719714d4f2d6873e92e5b355f33d4e580448.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOwnedObjectsRequestBuilderInternal instantiates a new OwnedObjectsRequestBuilder and sets the default values.
func NewOwnedObjectsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedObjectsRequestBuilder) {
    m := &OwnedObjectsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/ownedObjects{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOwnedObjectsRequestBuilder instantiates a new OwnedObjectsRequestBuilder and sets the default values.
func NewOwnedObjectsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedObjectsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOwnedObjectsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *OwnedObjectsRequestBuilder) Count()(*ib2a122a45f2f9ecdd4dec70c33255e57920c1fe41fe02b4b01be0195ffea3106.CountRequestBuilder) {
    return ib2a122a45f2f9ecdd4dec70c33255e57920c1fe41fe02b4b01be0195ffea3106.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation directory objects that are owned by the user. Read-only. Nullable. Supports $expand.
func (m *OwnedObjectsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OwnedObjectsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get directory objects that are owned by the user. Read-only. Nullable. Supports $expand.
func (m *OwnedObjectsRequestBuilder) Get(ctx context.Context, requestConfiguration *OwnedObjectsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
// Group casts the previous resource to group.
func (m *OwnedObjectsRequestBuilder) Group()(*ibaa1dc10d052b4b8688aceb90d308f0ac89f2de381f3ad51f64596a7fcedab6c.GroupRequestBuilder) {
    return ibaa1dc10d052b4b8688aceb90d308f0ac89f2de381f3ad51f64596a7fcedab6c.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *OwnedObjectsRequestBuilder) ServicePrincipal()(*i4df0a049a4392d6e4a658bc89720a18c16df1e9e765f88ad93fd42da67e4f3c4.ServicePrincipalRequestBuilder) {
    return i4df0a049a4392d6e4a658bc89720a18c16df1e9e765f88ad93fd42da67e4f3c4.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
