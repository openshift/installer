package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i35c871e364b7bd9c669482c3e558d704d96c2c6922e224cf70fc6ba4cb816548 "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item/allowedusers"
    i48812370840add29fbda20a153052c7c465785cc5280697e04043880e76211ef "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item/printer"
    if71bc0def059e5cd43ad758ac90e0476d545be67abc269b44c62b0e319bd0c1c "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item/allowedgroups"
    i19852dfc496eca0fe63108cfac8774cb2809469b264ed6de9f0834e636ab428c "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item/allowedgroups/item"
    i30908241df49cb55a70218cf04d4d0e1e3dfdde514367f11bb0ada2960507068 "github.com/microsoftgraph/msgraph-sdk-go/print/shares/item/allowedusers/item"
)

// PrinterShareItemRequestBuilder provides operations to manage the shares property of the microsoft.graph.print entity.
type PrinterShareItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PrinterShareItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterShareItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// PrinterShareItemRequestBuilderGetQueryParameters the list of printer shares registered in the tenant.
type PrinterShareItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PrinterShareItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterShareItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PrinterShareItemRequestBuilderGetQueryParameters
}
// PrinterShareItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrinterShareItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllowedGroups provides operations to manage the allowedGroups property of the microsoft.graph.printerShare entity.
func (m *PrinterShareItemRequestBuilder) AllowedGroups()(*if71bc0def059e5cd43ad758ac90e0476d545be67abc269b44c62b0e319bd0c1c.AllowedGroupsRequestBuilder) {
    return if71bc0def059e5cd43ad758ac90e0476d545be67abc269b44c62b0e319bd0c1c.NewAllowedGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllowedGroupsById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.print.shares.item.allowedGroups.item collection
func (m *PrinterShareItemRequestBuilder) AllowedGroupsById(id string)(*i19852dfc496eca0fe63108cfac8774cb2809469b264ed6de9f0834e636ab428c.GroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["group%2Did"] = id
    }
    return i19852dfc496eca0fe63108cfac8774cb2809469b264ed6de9f0834e636ab428c.NewGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AllowedUsers provides operations to manage the allowedUsers property of the microsoft.graph.printerShare entity.
func (m *PrinterShareItemRequestBuilder) AllowedUsers()(*i35c871e364b7bd9c669482c3e558d704d96c2c6922e224cf70fc6ba4cb816548.AllowedUsersRequestBuilder) {
    return i35c871e364b7bd9c669482c3e558d704d96c2c6922e224cf70fc6ba4cb816548.NewAllowedUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllowedUsersById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.print.shares.item.allowedUsers.item collection
func (m *PrinterShareItemRequestBuilder) AllowedUsersById(id string)(*i30908241df49cb55a70218cf04d4d0e1e3dfdde514367f11bb0ada2960507068.UserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["user%2Did"] = id
    }
    return i30908241df49cb55a70218cf04d4d0e1e3dfdde514367f11bb0ada2960507068.NewUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPrinterShareItemRequestBuilderInternal instantiates a new PrinterShareItemRequestBuilder and sets the default values.
func NewPrinterShareItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrinterShareItemRequestBuilder) {
    m := &PrinterShareItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/print/shares/{printerShare%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPrinterShareItemRequestBuilder instantiates a new PrinterShareItemRequestBuilder and sets the default values.
func NewPrinterShareItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrinterShareItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPrinterShareItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property shares for print
func (m *PrinterShareItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *PrinterShareItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the list of printer shares registered in the tenant.
func (m *PrinterShareItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PrinterShareItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property shares in print
func (m *PrinterShareItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable, requestConfiguration *PrinterShareItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property shares for print
func (m *PrinterShareItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *PrinterShareItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the list of printer shares registered in the tenant.
func (m *PrinterShareItemRequestBuilder) Get(ctx context.Context, requestConfiguration *PrinterShareItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrinterShareFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable), nil
}
// Patch update the navigation property shares in print
func (m *PrinterShareItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable, requestConfiguration *PrinterShareItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrinterShareFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrinterShareable), nil
}
// Printer provides operations to manage the printer property of the microsoft.graph.printerShare entity.
func (m *PrinterShareItemRequestBuilder) Printer()(*i48812370840add29fbda20a153052c7c465785cc5280697e04043880e76211ef.PrinterRequestBuilder) {
    return i48812370840add29fbda20a153052c7c465785cc5280697e04043880e76211ef.NewPrinterRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
