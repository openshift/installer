package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1bf7e274532170f1bce30fa85891cd1060aaf4432c850c3e6f600c74225358dd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/publish"
    i4fd99e30bd7e1ec5146299a811b6699b7a48158aff453a0ca332de49c3afb49e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columnlinks"
    i5980d4c6089817dc360b0282a3de8f025a947b489914b96333ae7bb03e59b88a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/base"
    i8d5d0e8815604e277a1a9310a20dd65c29f73cdd754a3ce39aac678b60860ca1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/basetypes"
    i965a40c98e9504af392fa55c0e043f32c195f17afac01638dd363f4492bea4e9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/ispublished"
    ibbce353b9cf3173510a3ab5b4232bab32639d0afdeb30cb6fa7be20ae2f9db7e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/copytodefaultcontentlocation"
    ic897225cd3a6c5e1636ce8c17f93a0c8b2150f74edeb8ac1991c2c3c584636f6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columnpositions"
    id483c2708d28edda81621fb56e0347a6ae606808dcbb1ae7f49bb596e0150406 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columns"
    id4f65ae8ee4c57e2f7ab4326616b2d20b8bfcc68216f9b5f5cdbeaf6a6099b33 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/associatewithhubsites"
    iee2737a01517790f43d05681c34a6491d2832aa82ce6073946c0001bc2a818f9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/unpublish"
    i1149c500ddccfb7c4db5f677250807d83df1d279a796576eb9613fedd3b2578f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columnlinks/item"
    i93827d02ff2765cecf70f15c1b4ff586989d30c2b0051a346ee15cd23c2ba844 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/basetypes/item"
    icf1aa423f6e90fa0f5e91966c342d736b4bf0d2e71c543272f6a32ecd56802ce "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columnpositions/item"
    ie281909f6939213db2ea9fee97bafef5874b4d351854a3f510c03c428183b137 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item/contenttypes/item/columns/item"
)

// ContentTypeItemRequestBuilder provides operations to manage the contentTypes property of the microsoft.graph.list entity.
type ContentTypeItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ContentTypeItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ContentTypeItemRequestBuilderGetQueryParameters the collection of content types present in this list.
type ContentTypeItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ContentTypeItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ContentTypeItemRequestBuilderGetQueryParameters
}
// ContentTypeItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AssociateWithHubSites provides operations to call the associateWithHubSites method.
func (m *ContentTypeItemRequestBuilder) AssociateWithHubSites()(*id4f65ae8ee4c57e2f7ab4326616b2d20b8bfcc68216f9b5f5cdbeaf6a6099b33.AssociateWithHubSitesRequestBuilder) {
    return id4f65ae8ee4c57e2f7ab4326616b2d20b8bfcc68216f9b5f5cdbeaf6a6099b33.NewAssociateWithHubSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Base provides operations to manage the base property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Base()(*i5980d4c6089817dc360b0282a3de8f025a947b489914b96333ae7bb03e59b88a.BaseRequestBuilder) {
    return i5980d4c6089817dc360b0282a3de8f025a947b489914b96333ae7bb03e59b88a.NewBaseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypes provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypes()(*i8d5d0e8815604e277a1a9310a20dd65c29f73cdd754a3ce39aac678b60860ca1.BaseTypesRequestBuilder) {
    return i8d5d0e8815604e277a1a9310a20dd65c29f73cdd754a3ce39aac678b60860ca1.NewBaseTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypesById provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypesById(id string)(*i93827d02ff2765cecf70f15c1b4ff586989d30c2b0051a346ee15cd23c2ba844.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did1"] = id
    }
    return i93827d02ff2765cecf70f15c1b4ff586989d30c2b0051a346ee15cd23c2ba844.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnLinks provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinks()(*i4fd99e30bd7e1ec5146299a811b6699b7a48158aff453a0ca332de49c3afb49e.ColumnLinksRequestBuilder) {
    return i4fd99e30bd7e1ec5146299a811b6699b7a48158aff453a0ca332de49c3afb49e.NewColumnLinksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnLinksById provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinksById(id string)(*i1149c500ddccfb7c4db5f677250807d83df1d279a796576eb9613fedd3b2578f.ColumnLinkItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnLink%2Did"] = id
    }
    return i1149c500ddccfb7c4db5f677250807d83df1d279a796576eb9613fedd3b2578f.NewColumnLinkItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnPositions provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositions()(*ic897225cd3a6c5e1636ce8c17f93a0c8b2150f74edeb8ac1991c2c3c584636f6.ColumnPositionsRequestBuilder) {
    return ic897225cd3a6c5e1636ce8c17f93a0c8b2150f74edeb8ac1991c2c3c584636f6.NewColumnPositionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnPositionsById provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositionsById(id string)(*icf1aa423f6e90fa0f5e91966c342d736b4bf0d2e71c543272f6a32ecd56802ce.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return icf1aa423f6e90fa0f5e91966c342d736b4bf0d2e71c543272f6a32ecd56802ce.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Columns()(*id483c2708d28edda81621fb56e0347a6ae606808dcbb1ae7f49bb596e0150406.ColumnsRequestBuilder) {
    return id483c2708d28edda81621fb56e0347a6ae606808dcbb1ae7f49bb596e0150406.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnsById(id string)(*ie281909f6939213db2ea9fee97bafef5874b4d351854a3f510c03c428183b137.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return ie281909f6939213db2ea9fee97bafef5874b4d351854a3f510c03c428183b137.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewContentTypeItemRequestBuilderInternal instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    m := &ContentTypeItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}/lists/{list%2Did}/contentTypes/{contentType%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewContentTypeItemRequestBuilder instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewContentTypeItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CopyToDefaultContentLocation provides operations to call the copyToDefaultContentLocation method.
func (m *ContentTypeItemRequestBuilder) CopyToDefaultContentLocation()(*ibbce353b9cf3173510a3ab5b4232bab32639d0afdeb30cb6fa7be20ae2f9db7e.CopyToDefaultContentLocationRequestBuilder) {
    return ibbce353b9cf3173510a3ab5b4232bab32639d0afdeb30cb6fa7be20ae2f9db7e.NewCopyToDefaultContentLocationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property contentTypes for groups
func (m *ContentTypeItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the collection of content types present in this list.
func (m *ContentTypeItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property contentTypes in groups
func (m *ContentTypeItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, requestConfiguration *ContentTypeItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property contentTypes for groups
func (m *ContentTypeItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the collection of content types present in this list.
func (m *ContentTypeItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateContentTypeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable), nil
}
// IsPublished provides operations to call the isPublished method.
func (m *ContentTypeItemRequestBuilder) IsPublished()(*i965a40c98e9504af392fa55c0e043f32c195f17afac01638dd363f4492bea4e9.IsPublishedRequestBuilder) {
    return i965a40c98e9504af392fa55c0e043f32c195f17afac01638dd363f4492bea4e9.NewIsPublishedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property contentTypes in groups
func (m *ContentTypeItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, requestConfiguration *ContentTypeItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateContentTypeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable), nil
}
// Publish provides operations to call the publish method.
func (m *ContentTypeItemRequestBuilder) Publish()(*i1bf7e274532170f1bce30fa85891cd1060aaf4432c850c3e6f600c74225358dd.PublishRequestBuilder) {
    return i1bf7e274532170f1bce30fa85891cd1060aaf4432c850c3e6f600c74225358dd.NewPublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unpublish provides operations to call the unpublish method.
func (m *ContentTypeItemRequestBuilder) Unpublish()(*iee2737a01517790f43d05681c34a6491d2832aa82ce6073946c0001bc2a818f9.UnpublishRequestBuilder) {
    return iee2737a01517790f43d05681c34a6491d2832aa82ce6073946c0001bc2a818f9.NewUnpublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
