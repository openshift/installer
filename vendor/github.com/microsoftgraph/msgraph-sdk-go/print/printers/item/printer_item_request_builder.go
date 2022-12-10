package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05f5154d5dc8e1096ef3ba18fb5bdb890712250391686b0ffd44f043e05db1fb "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/tasktriggers"
    i6ef3a73ca24cdfe4f01546d744e3ad5d5f17082ec1283caf5b2176d392cb9cd9 "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/connectors"
    i74b470ddc01590a04dfc5ba245b7ccb094a2f7b819d0a4ea94773869da1de8bc "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/shares"
    i95e425711e2e59b85807bfbd36bd18a7893a62957fccbb0dcc17647c9c3f5e90 "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/restorefactorydefaults"
    i02acec2b11a647c1fc0601a6c9f6283bf36064facb17fdc38860e7fc02ca0db0 "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/connectors/item"
    i652aa29a7f8828f592dc5bf3e9ada513fbcda2a67c9d060f001cc99e28d86159 "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/shares/item"
    i8a35ad54635f07016096adf7157efc92ede0057cd4a790d06a32332ce94f22cc "github.com/microsoftgraph/msgraph-sdk-go/print/printers/item/tasktriggers/item"
)

// PrinterItemRequestBuilder provides operations to manage the printers property of the microsoft.graph.print entity.
type PrinterItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PrinterItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// PrinterItemRequestBuilderGetQueryParameters the list of printers registered in the tenant.
type PrinterItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PrinterItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PrinterItemRequestBuilderGetQueryParameters
}
// PrinterItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Connectors provides operations to manage the connectors property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) Connectors()(*i6ef3a73ca24cdfe4f01546d744e3ad5d5f17082ec1283caf5b2176d392cb9cd9.ConnectorsRequestBuilder) {
    return i6ef3a73ca24cdfe4f01546d744e3ad5d5f17082ec1283caf5b2176d392cb9cd9.NewConnectorsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ConnectorsById provides operations to manage the connectors property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) ConnectorsById(id string)(*i02acec2b11a647c1fc0601a6c9f6283bf36064facb17fdc38860e7fc02ca0db0.PrintConnectorItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printConnector%2Did"] = id
    }
    return i02acec2b11a647c1fc0601a6c9f6283bf36064facb17fdc38860e7fc02ca0db0.NewPrintConnectorItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPrinterItemRequestBuilderInternal instantiates a new PrinterItemRequestBuilder and sets the default values.
func NewPrinterItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrinterItemRequestBuilder) {
    m := &PrinterItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/print/printers/{printer%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPrinterItemRequestBuilder instantiates a new PrinterItemRequestBuilder and sets the default values.
func NewPrinterItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrinterItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPrinterItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property printers for print
func (m *PrinterItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *PrinterItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the list of printers registered in the tenant.
func (m *PrinterItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PrinterItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property printers in print
func (m *PrinterItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable, requestConfiguration *PrinterItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property printers for print
func (m *PrinterItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *PrinterItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the list of printers registered in the tenant.
func (m *PrinterItemRequestBuilder) Get(ctx context.Context, requestConfiguration *PrinterItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrinterFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable), nil
}
// Patch update the navigation property printers in print
func (m *PrinterItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable, requestConfiguration *PrinterItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrinterFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Printerable), nil
}
// RestoreFactoryDefaults provides operations to call the restoreFactoryDefaults method.
func (m *PrinterItemRequestBuilder) RestoreFactoryDefaults()(*i95e425711e2e59b85807bfbd36bd18a7893a62957fccbb0dcc17647c9c3f5e90.RestoreFactoryDefaultsRequestBuilder) {
    return i95e425711e2e59b85807bfbd36bd18a7893a62957fccbb0dcc17647c9c3f5e90.NewRestoreFactoryDefaultsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Shares provides operations to manage the shares property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) Shares()(*i74b470ddc01590a04dfc5ba245b7ccb094a2f7b819d0a4ea94773869da1de8bc.SharesRequestBuilder) {
    return i74b470ddc01590a04dfc5ba245b7ccb094a2f7b819d0a4ea94773869da1de8bc.NewSharesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharesById provides operations to manage the shares property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) SharesById(id string)(*i652aa29a7f8828f592dc5bf3e9ada513fbcda2a67c9d060f001cc99e28d86159.PrinterShareItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printerShare%2Did"] = id
    }
    return i652aa29a7f8828f592dc5bf3e9ada513fbcda2a67c9d060f001cc99e28d86159.NewPrinterShareItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TaskTriggers provides operations to manage the taskTriggers property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) TaskTriggers()(*i05f5154d5dc8e1096ef3ba18fb5bdb890712250391686b0ffd44f043e05db1fb.TaskTriggersRequestBuilder) {
    return i05f5154d5dc8e1096ef3ba18fb5bdb890712250391686b0ffd44f043e05db1fb.NewTaskTriggersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TaskTriggersById provides operations to manage the taskTriggers property of the microsoft.graph.printer entity.
func (m *PrinterItemRequestBuilder) TaskTriggersById(id string)(*i8a35ad54635f07016096adf7157efc92ede0057cd4a790d06a32332ce94f22cc.PrintTaskTriggerItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printTaskTrigger%2Did"] = id
    }
    return i8a35ad54635f07016096adf7157efc92ede0057cd4a790d06a32332ce94f22cc.NewPrintTaskTriggerItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
