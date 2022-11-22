package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
    i0dee13fbbfab71a50e5566f2661d1d03f919ad275cd6c579617271e19406c70a "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/close"
    i67c4f3f8e8190f015213b2d3b265df2dc0747be0479d12355ccc5e0d7d8b8c58 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches"
    i717b9e1b2a0f68a8ec70146fa3075f72b30547f0b5fcbf201a3785ead600d336 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/operations"
    i787ee6e6d3c618ac8b391b3b1878ff0a2bd1e97f1b0d46ead96271c983110e48 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/tags"
    i8c393405892dbb22dfdddd56df47ec954195c0c6da1b1ae9c48e52cb7411d7ea "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/reviewsets"
    i8e1fe4855aca99eeadcd9321ed9443a1ac8bfd382b0c3c9963910ac5b6cdaf07 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/reopen"
    i9ef2d21b8fc69ba564a09403cd053d4d0735d936a826c7d2d564f1bb53c69740 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians"
    icac6ae8eb04885268777e9e86122e96896ad10b410f9bb104f3cad628eaea4e3 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/settings"
    id9d040632dbcbd497f9825c9a1eaf38b82b7529df01d632005643e53e9dcfec4 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/noncustodialdatasources"
    i0e99de5ab7cd5bd59a2a5c7d8039ba39dca2c25d48df8731bc7a3ede70007471 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/noncustodialdatasources/item"
    i26ebee56f309e3e6f8bf82c3cafa9eb54ee770f662b556d38343cd3e2d5a4f06 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item"
    i93ce34d72c05108db6d8c1ffc6122f45be91bc22789ccfccce5344ddbd4c8f99 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/tags/item"
    i97075ede7eb61e89f265b8ed8ec733f32d0abd1e07e5332d37dd11cba36b083a "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item"
    iba98a13db08b2d0b6857419935b53645fd30ba64a9165c5ea608a98c1d80c1fd "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/operations/item"
    id47c4ed6a666e7b18c72bc095983bb1369ae6f56786b1801b7900b6fca3206c3 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/reviewsets/item"
)

// EdiscoveryCaseItemRequestBuilder provides operations to manage the ediscoveryCases property of the microsoft.graph.security.casesRoot entity.
type EdiscoveryCaseItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EdiscoveryCaseItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCaseItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EdiscoveryCaseItemRequestBuilderGetQueryParameters get ediscoveryCases from security
type EdiscoveryCaseItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EdiscoveryCaseItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCaseItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EdiscoveryCaseItemRequestBuilderGetQueryParameters
}
// EdiscoveryCaseItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCaseItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Close provides operations to call the close method.
func (m *EdiscoveryCaseItemRequestBuilder) Close()(*i0dee13fbbfab71a50e5566f2661d1d03f919ad275cd6c579617271e19406c70a.CloseRequestBuilder) {
    return i0dee13fbbfab71a50e5566f2661d1d03f919ad275cd6c579617271e19406c70a.NewCloseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEdiscoveryCaseItemRequestBuilderInternal instantiates a new EdiscoveryCaseItemRequestBuilder and sets the default values.
func NewEdiscoveryCaseItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoveryCaseItemRequestBuilder) {
    m := &EdiscoveryCaseItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEdiscoveryCaseItemRequestBuilder instantiates a new EdiscoveryCaseItemRequestBuilder and sets the default values.
func NewEdiscoveryCaseItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoveryCaseItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEdiscoveryCaseItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property ediscoveryCases for security
func (m *EdiscoveryCaseItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EdiscoveryCaseItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get ediscoveryCases from security
func (m *EdiscoveryCaseItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EdiscoveryCaseItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property ediscoveryCases in security
func (m *EdiscoveryCaseItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable, requestConfiguration *EdiscoveryCaseItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Custodians provides operations to manage the custodians property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) Custodians()(*i9ef2d21b8fc69ba564a09403cd053d4d0735d936a826c7d2d564f1bb53c69740.CustodiansRequestBuilder) {
    return i9ef2d21b8fc69ba564a09403cd053d4d0735d936a826c7d2d564f1bb53c69740.NewCustodiansRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CustodiansById provides operations to manage the custodians property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) CustodiansById(id string)(*i26ebee56f309e3e6f8bf82c3cafa9eb54ee770f662b556d38343cd3e2d5a4f06.EdiscoveryCustodianItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoveryCustodian%2Did"] = id
    }
    return i26ebee56f309e3e6f8bf82c3cafa9eb54ee770f662b556d38343cd3e2d5a4f06.NewEdiscoveryCustodianItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Delete delete navigation property ediscoveryCases for security
func (m *EdiscoveryCaseItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EdiscoveryCaseItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get ediscoveryCases from security
func (m *EdiscoveryCaseItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EdiscoveryCaseItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoveryCaseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable), nil
}
// NoncustodialDataSources provides operations to manage the noncustodialDataSources property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) NoncustodialDataSources()(*id9d040632dbcbd497f9825c9a1eaf38b82b7529df01d632005643e53e9dcfec4.NoncustodialDataSourcesRequestBuilder) {
    return id9d040632dbcbd497f9825c9a1eaf38b82b7529df01d632005643e53e9dcfec4.NewNoncustodialDataSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NoncustodialDataSourcesById provides operations to manage the noncustodialDataSources property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) NoncustodialDataSourcesById(id string)(*i0e99de5ab7cd5bd59a2a5c7d8039ba39dca2c25d48df8731bc7a3ede70007471.EdiscoveryNoncustodialDataSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoveryNoncustodialDataSource%2Did"] = id
    }
    return i0e99de5ab7cd5bd59a2a5c7d8039ba39dca2c25d48df8731bc7a3ede70007471.NewEdiscoveryNoncustodialDataSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) Operations()(*i717b9e1b2a0f68a8ec70146fa3075f72b30547f0b5fcbf201a3785ead600d336.OperationsRequestBuilder) {
    return i717b9e1b2a0f68a8ec70146fa3075f72b30547f0b5fcbf201a3785ead600d336.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) OperationsById(id string)(*iba98a13db08b2d0b6857419935b53645fd30ba64a9165c5ea608a98c1d80c1fd.CaseOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["caseOperation%2Did"] = id
    }
    return iba98a13db08b2d0b6857419935b53645fd30ba64a9165c5ea608a98c1d80c1fd.NewCaseOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property ediscoveryCases in security
func (m *EdiscoveryCaseItemRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable, requestConfiguration *EdiscoveryCaseItemRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoveryCaseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCaseable), nil
}
// Reopen provides operations to call the reopen method.
func (m *EdiscoveryCaseItemRequestBuilder) Reopen()(*i8e1fe4855aca99eeadcd9321ed9443a1ac8bfd382b0c3c9963910ac5b6cdaf07.ReopenRequestBuilder) {
    return i8e1fe4855aca99eeadcd9321ed9443a1ac8bfd382b0c3c9963910ac5b6cdaf07.NewReopenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReviewSets provides operations to manage the reviewSets property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) ReviewSets()(*i8c393405892dbb22dfdddd56df47ec954195c0c6da1b1ae9c48e52cb7411d7ea.ReviewSetsRequestBuilder) {
    return i8c393405892dbb22dfdddd56df47ec954195c0c6da1b1ae9c48e52cb7411d7ea.NewReviewSetsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReviewSetsById provides operations to manage the reviewSets property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) ReviewSetsById(id string)(*id47c4ed6a666e7b18c72bc095983bb1369ae6f56786b1801b7900b6fca3206c3.EdiscoveryReviewSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoveryReviewSet%2Did"] = id
    }
    return id47c4ed6a666e7b18c72bc095983bb1369ae6f56786b1801b7900b6fca3206c3.NewEdiscoveryReviewSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Searches provides operations to manage the searches property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) Searches()(*i67c4f3f8e8190f015213b2d3b265df2dc0747be0479d12355ccc5e0d7d8b8c58.SearchesRequestBuilder) {
    return i67c4f3f8e8190f015213b2d3b265df2dc0747be0479d12355ccc5e0d7d8b8c58.NewSearchesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchesById provides operations to manage the searches property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) SearchesById(id string)(*i97075ede7eb61e89f265b8ed8ec733f32d0abd1e07e5332d37dd11cba36b083a.EdiscoverySearchItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoverySearch%2Did"] = id
    }
    return i97075ede7eb61e89f265b8ed8ec733f32d0abd1e07e5332d37dd11cba36b083a.NewEdiscoverySearchItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Settings provides operations to manage the settings property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) Settings()(*icac6ae8eb04885268777e9e86122e96896ad10b410f9bb104f3cad628eaea4e3.SettingsRequestBuilder) {
    return icac6ae8eb04885268777e9e86122e96896ad10b410f9bb104f3cad628eaea4e3.NewSettingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tags provides operations to manage the tags property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) Tags()(*i787ee6e6d3c618ac8b391b3b1878ff0a2bd1e97f1b0d46ead96271c983110e48.TagsRequestBuilder) {
    return i787ee6e6d3c618ac8b391b3b1878ff0a2bd1e97f1b0d46ead96271c983110e48.NewTagsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TagsById provides operations to manage the tags property of the microsoft.graph.security.ediscoveryCase entity.
func (m *EdiscoveryCaseItemRequestBuilder) TagsById(id string)(*i93ce34d72c05108db6d8c1ffc6122f45be91bc22789ccfccce5344ddbd4c8f99.EdiscoveryReviewTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoveryReviewTag%2Did"] = id
    }
    return i93ce34d72c05108db6d8c1ffc6122f45be91bc22789ccfccce5344ddbd4c8f99.NewEdiscoveryReviewTagItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
