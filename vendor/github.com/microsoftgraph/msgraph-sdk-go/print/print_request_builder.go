package print

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i0bda270a8db970823b8449367f0a1b6167aeb67bb1a0b116061e10c72e9e8d5d "github.com/microsoftgraph/msgraph-sdk-go/print/printers"
    i0f704d64ae3c276c8ff23cafb44315ef4067a668629360d124bdabfea38afdd3 "github.com/microsoftgraph/msgraph-sdk-go/print/connectors"
    i54feb3a1bee5e2d3891bf4b38d7ad46df95f464ec5502fc27ce12645532dcf04 "github.com/microsoftgraph/msgraph-sdk-go/print/shares"
    i802e345f3915c26a53a0e538ad96acb01b0ad515a9510581f00bcedb2e2d1266 "github.com/microsoftgraph/msgraph-sdk-go/print/services"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ib24f76d25d0e8f20d062be56f8b94808153dc20498d344b76af7867075ea8ad9 "github.com/microsoftgraph/msgraph-sdk-go/print/operations"
    ied8956e290d0c898cb2abe2b85d9373567ead7ae2e05a776597718dd3f1709e8 "github.com/microsoftgraph/msgraph-sdk-go/print/taskdefinitions"
    i0b3afc35ae08ae48e886b1ea1ae2fe47da46f82bd161cb48f9ff550676f78301 "github.com/microsoftgraph/msgraph-sdk-go/print/services/item"
    i7b3d94bc34459d3d2fa55e5dc4096ab7096e743ffc216c2d05b8f451f231fde1 "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item"
    i827522a32e9dfc8181a7f7793aec64fdf0814658430464c0b64af19c8e4e9a96 "github.com/microsoftgraph/msgraph-sdk-go/print/taskdefinitions/item"
    i8e18f5c596955d2371ff489e959f3488755280f5b7014ba5a0f5b95b7dc44bf2 "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item"
    id0e61828b4320e56c027152ab07b04414264a89aa169cd832394397ada28defd "github.com/microsoftgraph/msgraph-sdk-go/print/operations/item"
    iebd42d165d4ff6aca1eec85e8351f67b97a4208ea00a5bc8ac67f4f1f191b677 "github.com/microsoftgraph/msgraph-sdk-go/print/connectors/item"
)

// PrintRequestBuilder provides operations to manage the print singleton.
type PrintRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PrintRequestBuilderGetQueryParameters get print
type PrintRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PrintRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrintRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PrintRequestBuilderGetQueryParameters
}
// PrintRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrintRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Connectors provides operations to manage the connectors property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) Connectors()(*i0f704d64ae3c276c8ff23cafb44315ef4067a668629360d124bdabfea38afdd3.ConnectorsRequestBuilder) {
    return i0f704d64ae3c276c8ff23cafb44315ef4067a668629360d124bdabfea38afdd3.NewConnectorsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ConnectorsById provides operations to manage the connectors property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) ConnectorsById(id string)(*iebd42d165d4ff6aca1eec85e8351f67b97a4208ea00a5bc8ac67f4f1f191b677.PrintConnectorItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printConnector%2Did"] = id
    }
    return iebd42d165d4ff6aca1eec85e8351f67b97a4208ea00a5bc8ac67f4f1f191b677.NewPrintConnectorItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPrintRequestBuilderInternal instantiates a new PrintRequestBuilder and sets the default values.
func NewPrintRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrintRequestBuilder) {
    m := &PrintRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/print{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPrintRequestBuilder instantiates a new PrintRequestBuilder and sets the default values.
func NewPrintRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrintRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPrintRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get print
func (m *PrintRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PrintRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update print
func (m *PrintRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable, requestConfiguration *PrintRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get print
func (m *PrintRequestBuilder) Get(ctx context.Context, requestConfiguration *PrintRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable), nil
}
// Operations provides operations to manage the operations property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) Operations()(*ib24f76d25d0e8f20d062be56f8b94808153dc20498d344b76af7867075ea8ad9.OperationsRequestBuilder) {
    return ib24f76d25d0e8f20d062be56f8b94808153dc20498d344b76af7867075ea8ad9.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) OperationsById(id string)(*id0e61828b4320e56c027152ab07b04414264a89aa169cd832394397ada28defd.PrintOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printOperation%2Did"] = id
    }
    return id0e61828b4320e56c027152ab07b04414264a89aa169cd832394397ada28defd.NewPrintOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update print
func (m *PrintRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable, requestConfiguration *PrintRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printable), nil
}
// Printers provides operations to manage the printers property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) Printers()(*i0bda270a8db970823b8449367f0a1b6167aeb67bb1a0b116061e10c72e9e8d5d.PrintersRequestBuilder) {
    return i0bda270a8db970823b8449367f0a1b6167aeb67bb1a0b116061e10c72e9e8d5d.NewPrintersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PrintersById provides operations to manage the printers property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) PrintersById(id string)(*i7b3d94bc34459d3d2fa55e5dc4096ab7096e743ffc216c2d05b8f451f231fde1.PrinterItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printer%2Did"] = id
    }
    return i7b3d94bc34459d3d2fa55e5dc4096ab7096e743ffc216c2d05b8f451f231fde1.NewPrinterItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Services provides operations to manage the services property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) Services()(*i802e345f3915c26a53a0e538ad96acb01b0ad515a9510581f00bcedb2e2d1266.ServicesRequestBuilder) {
    return i802e345f3915c26a53a0e538ad96acb01b0ad515a9510581f00bcedb2e2d1266.NewServicesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicesById provides operations to manage the services property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) ServicesById(id string)(*i0b3afc35ae08ae48e886b1ea1ae2fe47da46f82bd161cb48f9ff550676f78301.PrintServiceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printService%2Did"] = id
    }
    return i0b3afc35ae08ae48e886b1ea1ae2fe47da46f82bd161cb48f9ff550676f78301.NewPrintServiceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Shares provides operations to manage the shares property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) Shares()(*i54feb3a1bee5e2d3891bf4b38d7ad46df95f464ec5502fc27ce12645532dcf04.SharesRequestBuilder) {
    return i54feb3a1bee5e2d3891bf4b38d7ad46df95f464ec5502fc27ce12645532dcf04.NewSharesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharesById provides operations to manage the shares property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) SharesById(id string)(*i8e18f5c596955d2371ff489e959f3488755280f5b7014ba5a0f5b95b7dc44bf2.PrinterShareItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printerShare%2Did"] = id
    }
    return i8e18f5c596955d2371ff489e959f3488755280f5b7014ba5a0f5b95b7dc44bf2.NewPrinterShareItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TaskDefinitions provides operations to manage the taskDefinitions property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) TaskDefinitions()(*ied8956e290d0c898cb2abe2b85d9373567ead7ae2e05a776597718dd3f1709e8.TaskDefinitionsRequestBuilder) {
    return ied8956e290d0c898cb2abe2b85d9373567ead7ae2e05a776597718dd3f1709e8.NewTaskDefinitionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TaskDefinitionsById provides operations to manage the taskDefinitions property of the microsoft.graph.print entity.
func (m *PrintRequestBuilder) TaskDefinitionsById(id string)(*i827522a32e9dfc8181a7f7793aec64fdf0814658430464c0b64af19c8e4e9a96.PrintTaskDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printTaskDefinition%2Did"] = id
    }
    return i827522a32e9dfc8181a7f7793aec64fdf0814658430464c0b64af19c8e4e9a96.NewPrintTaskDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
