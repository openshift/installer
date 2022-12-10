package sites

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i04168ffe2353eecc18b03924033a8b3a497518b96de485de3c54646fd2088af1 "github.com/microsoftgraph/msgraph-sdk-go/sites/add"
    i36e45f2c65a6887462977eb125acfe848013e7ff15c056bbdef95a2654f52bba "github.com/microsoftgraph/msgraph-sdk-go/sites/remove"
    i699f99a29a98391477e726e777013f008f0b2380c685707edf69f79319d49705 "github.com/microsoftgraph/msgraph-sdk-go/sites/count"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SitesRequestBuilder provides operations to manage the collection of site entities.
type SitesRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// SitesRequestBuilderGetQueryParameters search across a SharePoint tenant for [sites][] that match keywords provided. The only property that works for sorting is **createdDateTime**. The search filter is a free text search that uses multiple properties when retrieving the search results.
type SitesRequestBuilderGetQueryParameters struct {
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
// SitesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SitesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SitesRequestBuilderGetQueryParameters
}
// Add provides operations to call the add method.
func (m *SitesRequestBuilder) Add()(*i04168ffe2353eecc18b03924033a8b3a497518b96de485de3c54646fd2088af1.AddRequestBuilder) {
    return i04168ffe2353eecc18b03924033a8b3a497518b96de485de3c54646fd2088af1.NewAddRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewSitesRequestBuilderInternal instantiates a new SitesRequestBuilder and sets the default values.
func NewSitesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SitesRequestBuilder) {
    m := &SitesRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/sites{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewSitesRequestBuilder instantiates a new SitesRequestBuilder and sets the default values.
func NewSitesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SitesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSitesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *SitesRequestBuilder) Count()(*i699f99a29a98391477e726e777013f008f0b2380c685707edf69f79319d49705.CountRequestBuilder) {
    return i699f99a29a98391477e726e777013f008f0b2380c685707edf69f79319d49705.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation search across a SharePoint tenant for [sites][] that match keywords provided. The only property that works for sorting is **createdDateTime**. The search filter is a free text search that uses multiple properties when retrieving the search results.
func (m *SitesRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *SitesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get search across a SharePoint tenant for [sites][] that match keywords provided. The only property that works for sorting is **createdDateTime**. The search filter is a free text search that uses multiple properties when retrieving the search results.
func (m *SitesRequestBuilder) Get(ctx context.Context, requestConfiguration *SitesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SiteCollectionResponseable), nil
}
// Remove provides operations to call the remove method.
func (m *SitesRequestBuilder) Remove()(*i36e45f2c65a6887462977eb125acfe848013e7ff15c056bbdef95a2654f52bba.RemoveRequestBuilder) {
    return i36e45f2c65a6887462977eb125acfe848013e7ff15c056bbdef95a2654f52bba.NewRemoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
