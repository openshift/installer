package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05aa6bb50d61c9a91685ccd79551fe45c7fd5a6f41b2f9052961feb2dc394eb0 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/associatewithhubsites"
    i0cf1daad8519139f1c1dddb4eb3f61c0e4ec6e71759a998ee5a46fd4b2468890 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columns"
    i2bd83908df726274d2c1ed9ae647d9d262c28eb0503cd31251764fddce791655 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columnpositions"
    i2e3ff723a01a654065641276ad9c59facdf69c3d223421c5e324df2133e20d2d "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columnlinks"
    i43a81088f78dc3b0d7c4c4386ce5e5fe9e9adeeeed6c450fa8819629e5edd1d2 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/publish"
    i44ccd150eeae5d04e5afc727774f113c375506bb3a712e0dc54818ed2d2a7173 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/ispublished"
    i5e2cbca2e66013cc685d3b90b6f6620ddfc6a3ec8bd3560353ccf1bc0243fec6 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/unpublish"
    ib8daa131ed0a2d3b42d5cabdcf429cf168fa29e9581c9b560445404fc3855531 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/basetypes"
    id257be66b68c6ec4f94dd1b16a346a3628e3e2ff5283493c88ea924e3651e546 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/base"
    ie4cb5930f587e11ddf957e2716c252f45c2bae1f0cbb3eca003f83d9478ad44d "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/copytodefaultcontentlocation"
    i4e76b560d7ab063fee73b44a70701c5d35c8cd952639ebdad225d0fd1dd440b6 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columns/item"
    i8031712eaa64ad6bcb9a08e5abc085cfafe42aef18be156e924188e8a86b8403 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/basetypes/item"
    i9d5cdbf190630a4b53f62d7eb59d02a008abe57fc0685c8324eec7d57bd1e98d "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columnpositions/item"
    if03de5dd741d9ed2d25876eea430ca53f3feda55c4d1a8e9ac5b875bac63413a "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item/columnlinks/item"
)

// ContentTypeItemRequestBuilder provides operations to manage the contentTypes property of the microsoft.graph.site entity.
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
// ContentTypeItemRequestBuilderGetQueryParameters the collection of content types defined for this site.
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
func (m *ContentTypeItemRequestBuilder) AssociateWithHubSites()(*i05aa6bb50d61c9a91685ccd79551fe45c7fd5a6f41b2f9052961feb2dc394eb0.AssociateWithHubSitesRequestBuilder) {
    return i05aa6bb50d61c9a91685ccd79551fe45c7fd5a6f41b2f9052961feb2dc394eb0.NewAssociateWithHubSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Base provides operations to manage the base property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Base()(*id257be66b68c6ec4f94dd1b16a346a3628e3e2ff5283493c88ea924e3651e546.BaseRequestBuilder) {
    return id257be66b68c6ec4f94dd1b16a346a3628e3e2ff5283493c88ea924e3651e546.NewBaseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypes provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypes()(*ib8daa131ed0a2d3b42d5cabdcf429cf168fa29e9581c9b560445404fc3855531.BaseTypesRequestBuilder) {
    return ib8daa131ed0a2d3b42d5cabdcf429cf168fa29e9581c9b560445404fc3855531.NewBaseTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypesById provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypesById(id string)(*i8031712eaa64ad6bcb9a08e5abc085cfafe42aef18be156e924188e8a86b8403.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did1"] = id
    }
    return i8031712eaa64ad6bcb9a08e5abc085cfafe42aef18be156e924188e8a86b8403.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnLinks provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinks()(*i2e3ff723a01a654065641276ad9c59facdf69c3d223421c5e324df2133e20d2d.ColumnLinksRequestBuilder) {
    return i2e3ff723a01a654065641276ad9c59facdf69c3d223421c5e324df2133e20d2d.NewColumnLinksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnLinksById provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinksById(id string)(*if03de5dd741d9ed2d25876eea430ca53f3feda55c4d1a8e9ac5b875bac63413a.ColumnLinkItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnLink%2Did"] = id
    }
    return if03de5dd741d9ed2d25876eea430ca53f3feda55c4d1a8e9ac5b875bac63413a.NewColumnLinkItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnPositions provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositions()(*i2bd83908df726274d2c1ed9ae647d9d262c28eb0503cd31251764fddce791655.ColumnPositionsRequestBuilder) {
    return i2bd83908df726274d2c1ed9ae647d9d262c28eb0503cd31251764fddce791655.NewColumnPositionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnPositionsById provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositionsById(id string)(*i9d5cdbf190630a4b53f62d7eb59d02a008abe57fc0685c8324eec7d57bd1e98d.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i9d5cdbf190630a4b53f62d7eb59d02a008abe57fc0685c8324eec7d57bd1e98d.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Columns()(*i0cf1daad8519139f1c1dddb4eb3f61c0e4ec6e71759a998ee5a46fd4b2468890.ColumnsRequestBuilder) {
    return i0cf1daad8519139f1c1dddb4eb3f61c0e4ec6e71759a998ee5a46fd4b2468890.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnsById(id string)(*i4e76b560d7ab063fee73b44a70701c5d35c8cd952639ebdad225d0fd1dd440b6.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i4e76b560d7ab063fee73b44a70701c5d35c8cd952639ebdad225d0fd1dd440b6.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewContentTypeItemRequestBuilderInternal instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    m := &ContentTypeItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/sites/{site%2Did}/contentTypes/{contentType%2Did}{?%24select,%24expand}";
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
func (m *ContentTypeItemRequestBuilder) CopyToDefaultContentLocation()(*ie4cb5930f587e11ddf957e2716c252f45c2bae1f0cbb3eca003f83d9478ad44d.CopyToDefaultContentLocationRequestBuilder) {
    return ie4cb5930f587e11ddf957e2716c252f45c2bae1f0cbb3eca003f83d9478ad44d.NewCopyToDefaultContentLocationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property contentTypes for sites
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
// CreateGetRequestInformation the collection of content types defined for this site.
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
// CreatePatchRequestInformation update the navigation property contentTypes in sites
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
// Delete delete navigation property contentTypes for sites
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
// Get the collection of content types defined for this site.
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
func (m *ContentTypeItemRequestBuilder) IsPublished()(*i44ccd150eeae5d04e5afc727774f113c375506bb3a712e0dc54818ed2d2a7173.IsPublishedRequestBuilder) {
    return i44ccd150eeae5d04e5afc727774f113c375506bb3a712e0dc54818ed2d2a7173.NewIsPublishedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property contentTypes in sites
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
func (m *ContentTypeItemRequestBuilder) Publish()(*i43a81088f78dc3b0d7c4c4386ce5e5fe9e9adeeeed6c450fa8819629e5edd1d2.PublishRequestBuilder) {
    return i43a81088f78dc3b0d7c4c4386ce5e5fe9e9adeeeed6c450fa8819629e5edd1d2.NewPublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unpublish provides operations to call the unpublish method.
func (m *ContentTypeItemRequestBuilder) Unpublish()(*i5e2cbca2e66013cc685d3b90b6f6620ddfc6a3ec8bd3560353ccf1bc0243fec6.UnpublishRequestBuilder) {
    return i5e2cbca2e66013cc685d3b90b6f6620ddfc6a3ec8bd3560353ccf1bc0243fec6.NewUnpublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
