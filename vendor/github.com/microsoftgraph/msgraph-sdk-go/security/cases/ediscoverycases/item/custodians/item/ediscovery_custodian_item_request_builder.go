package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
    i13941c9d88af85554fed7be8defe9bd91ce9902e55573f88d799747c5dfd8e82 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/unifiedgroupsources"
    i1869c24cbd3274a8c33be4fda4e1f763877159294011e7ec482fcd566cd75e99 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/sitesources"
    i69a932198cf358e277c15d3468ebfab21a130c7751f6cd66dd20b707d5947ff1 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/lastindexoperation"
    i7536d78f675815361616e6c4cb30cd9dfd237d77a976d90cb47312d48369a76c "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/usersources"
    i799a1eb30a03251c27ea2ae2e887e353346da13e01be4a22edf504ec083d743a "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/removehold"
    ib2a8ef42b6524c07a09d3fc949f5eb5e275eadabe56c240c55d51c4012ba09e3 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/applyhold"
    ic1a98148b3e81a067532df947756c9d8a1c79ede5edf2b16d24786aedce309f8 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/updateindex"
    ic27cbd57655ffa65ff97a5af26e217f5cd8317344562330e109f18842fb6ab06 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/activate"
    ie5c3b211ff8a54e5ba14f56cd47e975ea605505714a1375d2f1c3c7eac30af85 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/release"
    i0a99804f8b27b23874e992f201aa8d77adbddb992dd89e3c57c520e042787d87 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/sitesources/item"
    ia1a7339eb6913976688accec03653670e1614c197f13df632862da97431dc874 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/unifiedgroupsources/item"
    ia64d457cd1a449247c12dc66ae6bc9a7df87f7374519a0072f8379482c539fb7 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/custodians/item/usersources/item"
)

// EdiscoveryCustodianItemRequestBuilder provides operations to manage the custodians property of the microsoft.graph.security.ediscoveryCase entity.
type EdiscoveryCustodianItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EdiscoveryCustodianItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCustodianItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EdiscoveryCustodianItemRequestBuilderGetQueryParameters returns a list of case ediscoveryCustodian objects for this case.
type EdiscoveryCustodianItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EdiscoveryCustodianItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCustodianItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EdiscoveryCustodianItemRequestBuilderGetQueryParameters
}
// EdiscoveryCustodianItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoveryCustodianItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Activate provides operations to call the activate method.
func (m *EdiscoveryCustodianItemRequestBuilder) Activate()(*ic27cbd57655ffa65ff97a5af26e217f5cd8317344562330e109f18842fb6ab06.ActivateRequestBuilder) {
    return ic27cbd57655ffa65ff97a5af26e217f5cd8317344562330e109f18842fb6ab06.NewActivateRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ApplyHold provides operations to call the applyHold method.
func (m *EdiscoveryCustodianItemRequestBuilder) ApplyHold()(*ib2a8ef42b6524c07a09d3fc949f5eb5e275eadabe56c240c55d51c4012ba09e3.ApplyHoldRequestBuilder) {
    return ib2a8ef42b6524c07a09d3fc949f5eb5e275eadabe56c240c55d51c4012ba09e3.NewApplyHoldRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEdiscoveryCustodianItemRequestBuilderInternal instantiates a new EdiscoveryCustodianItemRequestBuilder and sets the default values.
func NewEdiscoveryCustodianItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoveryCustodianItemRequestBuilder) {
    m := &EdiscoveryCustodianItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}/custodians/{ediscoveryCustodian%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEdiscoveryCustodianItemRequestBuilder instantiates a new EdiscoveryCustodianItemRequestBuilder and sets the default values.
func NewEdiscoveryCustodianItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoveryCustodianItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEdiscoveryCustodianItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property custodians for security
func (m *EdiscoveryCustodianItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EdiscoveryCustodianItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation returns a list of case ediscoveryCustodian objects for this case.
func (m *EdiscoveryCustodianItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EdiscoveryCustodianItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property custodians in security
func (m *EdiscoveryCustodianItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable, requestConfiguration *EdiscoveryCustodianItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property custodians for security
func (m *EdiscoveryCustodianItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EdiscoveryCustodianItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get returns a list of case ediscoveryCustodian objects for this case.
func (m *EdiscoveryCustodianItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EdiscoveryCustodianItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoveryCustodianFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable), nil
}
// LastIndexOperation provides operations to manage the lastIndexOperation property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) LastIndexOperation()(*i69a932198cf358e277c15d3468ebfab21a130c7751f6cd66dd20b707d5947ff1.LastIndexOperationRequestBuilder) {
    return i69a932198cf358e277c15d3468ebfab21a130c7751f6cd66dd20b707d5947ff1.NewLastIndexOperationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property custodians in security
func (m *EdiscoveryCustodianItemRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable, requestConfiguration *EdiscoveryCustodianItemRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoveryCustodianFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryCustodianable), nil
}
// Release provides operations to call the release method.
func (m *EdiscoveryCustodianItemRequestBuilder) Release()(*ie5c3b211ff8a54e5ba14f56cd47e975ea605505714a1375d2f1c3c7eac30af85.ReleaseRequestBuilder) {
    return ie5c3b211ff8a54e5ba14f56cd47e975ea605505714a1375d2f1c3c7eac30af85.NewReleaseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveHold provides operations to call the removeHold method.
func (m *EdiscoveryCustodianItemRequestBuilder) RemoveHold()(*i799a1eb30a03251c27ea2ae2e887e353346da13e01be4a22edf504ec083d743a.RemoveHoldRequestBuilder) {
    return i799a1eb30a03251c27ea2ae2e887e353346da13e01be4a22edf504ec083d743a.NewRemoveHoldRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SiteSources provides operations to manage the siteSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) SiteSources()(*i1869c24cbd3274a8c33be4fda4e1f763877159294011e7ec482fcd566cd75e99.SiteSourcesRequestBuilder) {
    return i1869c24cbd3274a8c33be4fda4e1f763877159294011e7ec482fcd566cd75e99.NewSiteSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SiteSourcesById provides operations to manage the siteSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) SiteSourcesById(id string)(*i0a99804f8b27b23874e992f201aa8d77adbddb992dd89e3c57c520e042787d87.SiteSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["siteSource%2Did"] = id
    }
    return i0a99804f8b27b23874e992f201aa8d77adbddb992dd89e3c57c520e042787d87.NewSiteSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UnifiedGroupSources provides operations to manage the unifiedGroupSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) UnifiedGroupSources()(*i13941c9d88af85554fed7be8defe9bd91ce9902e55573f88d799747c5dfd8e82.UnifiedGroupSourcesRequestBuilder) {
    return i13941c9d88af85554fed7be8defe9bd91ce9902e55573f88d799747c5dfd8e82.NewUnifiedGroupSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UnifiedGroupSourcesById provides operations to manage the unifiedGroupSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) UnifiedGroupSourcesById(id string)(*ia1a7339eb6913976688accec03653670e1614c197f13df632862da97431dc874.UnifiedGroupSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedGroupSource%2Did"] = id
    }
    return ia1a7339eb6913976688accec03653670e1614c197f13df632862da97431dc874.NewUnifiedGroupSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UpdateIndex provides operations to call the updateIndex method.
func (m *EdiscoveryCustodianItemRequestBuilder) UpdateIndex()(*ic1a98148b3e81a067532df947756c9d8a1c79ede5edf2b16d24786aedce309f8.UpdateIndexRequestBuilder) {
    return ic1a98148b3e81a067532df947756c9d8a1c79ede5edf2b16d24786aedce309f8.NewUpdateIndexRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UserSources provides operations to manage the userSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) UserSources()(*i7536d78f675815361616e6c4cb30cd9dfd237d77a976d90cb47312d48369a76c.UserSourcesRequestBuilder) {
    return i7536d78f675815361616e6c4cb30cd9dfd237d77a976d90cb47312d48369a76c.NewUserSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UserSourcesById provides operations to manage the userSources property of the microsoft.graph.security.ediscoveryCustodian entity.
func (m *EdiscoveryCustodianItemRequestBuilder) UserSourcesById(id string)(*ia64d457cd1a449247c12dc66ae6bc9a7df87f7374519a0072f8379482c539fb7.UserSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["userSource%2Did"] = id
    }
    return ia64d457cd1a449247c12dc66ae6bc9a7df87f7374519a0072f8379482c539fb7.NewUserSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
