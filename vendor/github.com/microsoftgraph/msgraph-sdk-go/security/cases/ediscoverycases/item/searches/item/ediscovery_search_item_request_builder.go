package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
    i011377c6234c0d9b7ec8331fb54d711964dbfce60a852ebcd6bcb67cfb0f7201 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/noncustodialsources"
    i108f500ba0337f68c69d081fdc1baf1da67981a029cc46cd55bcbb59d2f5e132 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/custodiansources"
    i230698fd5f86ad1b935ee0c129ab3e4e19a8b52675e16ef480c1cb8ac1c8ef9b "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/lastestimatestatisticsoperation"
    i83739c129d620b1070fe64795ee0e2a738218f6c3ec0b7a008e1adf890343b07 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/additionalsources"
    i93ba3bfd10e83ad08be0c2a6eb7231d725df6f0cde5f2f67d84601f198b47e97 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/estimatestatistics"
    ia5644c6e0fd2c736f8c83e14b34df3ff00a6361d19807e5ecc0a54794e3ca63d "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/purgedata"
    iecd6a3b6025ea7091d425fac6cf5d255595c827b86484d73aacf3742528b5cb7 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/addtoreviewsetoperation"
    i8f302d797cc13f4a2ecbd36d19a3a180ea9fc8bb3596fb79074e6ecd02dc6bd4 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/additionalsources/item"
    idec6a28edc7173c367f0bd14b204505e7b6645a0381c4ee75a61bb11c7fc79f0 "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/noncustodialsources/item"
    if56cdee80c7867e870947a6c362dab6e3bb804f5bf5cc92a2a195f902db1785c "github.com/microsoftgraph/msgraph-sdk-go/security/cases/ediscoverycases/item/searches/item/custodiansources/item"
)

// EdiscoverySearchItemRequestBuilder provides operations to manage the searches property of the microsoft.graph.security.ediscoveryCase entity.
type EdiscoverySearchItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EdiscoverySearchItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoverySearchItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EdiscoverySearchItemRequestBuilderGetQueryParameters returns a list of eDiscoverySearch objects associated with this case.
type EdiscoverySearchItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EdiscoverySearchItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoverySearchItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EdiscoverySearchItemRequestBuilderGetQueryParameters
}
// EdiscoverySearchItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EdiscoverySearchItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AdditionalSources provides operations to manage the additionalSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) AdditionalSources()(*i83739c129d620b1070fe64795ee0e2a738218f6c3ec0b7a008e1adf890343b07.AdditionalSourcesRequestBuilder) {
    return i83739c129d620b1070fe64795ee0e2a738218f6c3ec0b7a008e1adf890343b07.NewAdditionalSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AdditionalSourcesById provides operations to manage the additionalSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) AdditionalSourcesById(id string)(*i8f302d797cc13f4a2ecbd36d19a3a180ea9fc8bb3596fb79074e6ecd02dc6bd4.DataSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["dataSource%2Did"] = id
    }
    return i8f302d797cc13f4a2ecbd36d19a3a180ea9fc8bb3596fb79074e6ecd02dc6bd4.NewDataSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AddToReviewSetOperation provides operations to manage the addToReviewSetOperation property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) AddToReviewSetOperation()(*iecd6a3b6025ea7091d425fac6cf5d255595c827b86484d73aacf3742528b5cb7.AddToReviewSetOperationRequestBuilder) {
    return iecd6a3b6025ea7091d425fac6cf5d255595c827b86484d73aacf3742528b5cb7.NewAddToReviewSetOperationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewEdiscoverySearchItemRequestBuilderInternal instantiates a new EdiscoverySearchItemRequestBuilder and sets the default values.
func NewEdiscoverySearchItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoverySearchItemRequestBuilder) {
    m := &EdiscoverySearchItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}/searches/{ediscoverySearch%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEdiscoverySearchItemRequestBuilder instantiates a new EdiscoverySearchItemRequestBuilder and sets the default values.
func NewEdiscoverySearchItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EdiscoverySearchItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEdiscoverySearchItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property searches for security
func (m *EdiscoverySearchItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EdiscoverySearchItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation returns a list of eDiscoverySearch objects associated with this case.
func (m *EdiscoverySearchItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EdiscoverySearchItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property searches in security
func (m *EdiscoverySearchItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable, requestConfiguration *EdiscoverySearchItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CustodianSources provides operations to manage the custodianSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) CustodianSources()(*i108f500ba0337f68c69d081fdc1baf1da67981a029cc46cd55bcbb59d2f5e132.CustodianSourcesRequestBuilder) {
    return i108f500ba0337f68c69d081fdc1baf1da67981a029cc46cd55bcbb59d2f5e132.NewCustodianSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CustodianSourcesById provides operations to manage the custodianSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) CustodianSourcesById(id string)(*if56cdee80c7867e870947a6c362dab6e3bb804f5bf5cc92a2a195f902db1785c.DataSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["dataSource%2Did"] = id
    }
    return if56cdee80c7867e870947a6c362dab6e3bb804f5bf5cc92a2a195f902db1785c.NewDataSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Delete delete navigation property searches for security
func (m *EdiscoverySearchItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EdiscoverySearchItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// EstimateStatistics provides operations to call the estimateStatistics method.
func (m *EdiscoverySearchItemRequestBuilder) EstimateStatistics()(*i93ba3bfd10e83ad08be0c2a6eb7231d725df6f0cde5f2f67d84601f198b47e97.EstimateStatisticsRequestBuilder) {
    return i93ba3bfd10e83ad08be0c2a6eb7231d725df6f0cde5f2f67d84601f198b47e97.NewEstimateStatisticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get returns a list of eDiscoverySearch objects associated with this case.
func (m *EdiscoverySearchItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EdiscoverySearchItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoverySearchFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable), nil
}
// LastEstimateStatisticsOperation provides operations to manage the lastEstimateStatisticsOperation property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) LastEstimateStatisticsOperation()(*i230698fd5f86ad1b935ee0c129ab3e4e19a8b52675e16ef480c1cb8ac1c8ef9b.LastEstimateStatisticsOperationRequestBuilder) {
    return i230698fd5f86ad1b935ee0c129ab3e4e19a8b52675e16ef480c1cb8ac1c8ef9b.NewLastEstimateStatisticsOperationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NoncustodialSources provides operations to manage the noncustodialSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) NoncustodialSources()(*i011377c6234c0d9b7ec8331fb54d711964dbfce60a852ebcd6bcb67cfb0f7201.NoncustodialSourcesRequestBuilder) {
    return i011377c6234c0d9b7ec8331fb54d711964dbfce60a852ebcd6bcb67cfb0f7201.NewNoncustodialSourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NoncustodialSourcesById provides operations to manage the noncustodialSources property of the microsoft.graph.security.ediscoverySearch entity.
func (m *EdiscoverySearchItemRequestBuilder) NoncustodialSourcesById(id string)(*idec6a28edc7173c367f0bd14b204505e7b6645a0381c4ee75a61bb11c7fc79f0.EdiscoveryNoncustodialDataSourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["ediscoveryNoncustodialDataSource%2Did"] = id
    }
    return idec6a28edc7173c367f0bd14b204505e7b6645a0381c4ee75a61bb11c7fc79f0.NewEdiscoveryNoncustodialDataSourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property searches in security
func (m *EdiscoverySearchItemRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable, requestConfiguration *EdiscoverySearchItemRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoverySearchFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoverySearchable), nil
}
// PurgeData provides operations to call the purgeData method.
func (m *EdiscoverySearchItemRequestBuilder) PurgeData()(*ia5644c6e0fd2c736f8c83e14b34df3ff00a6361d19807e5ecc0a54794e3ca63d.PurgeDataRequestBuilder) {
    return ia5644c6e0fd2c736f8c83e14b34df3ff00a6361d19807e5ecc0a54794e3ca63d.NewPurgeDataRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
