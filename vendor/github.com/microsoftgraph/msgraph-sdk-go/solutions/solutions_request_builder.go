package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i6d5b6940d5e24f8aee3eb7eef5499c5a97c9c71485b644e5b94efadbc56609b0 "github.com/microsoftgraph/msgraph-sdk-go/solutions/bookingcurrencies"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ic2077c84f4cea40593f46a097344f757e71b10d42626a06348d8a6856a4a6408 "github.com/microsoftgraph/msgraph-sdk-go/solutions/bookingbusinesses"
    i30929f1a7f5c848dd5ee1e9d3e9b1dc2926ba5d92f9f5ca889f41edee949d39a "github.com/microsoftgraph/msgraph-sdk-go/solutions/bookingcurrencies/item"
    i98532a0d43958ac46a84969d106d7877583a8e83cd026d67fbc826e462b1e8c7 "github.com/microsoftgraph/msgraph-sdk-go/solutions/bookingbusinesses/item"
)

// SolutionsRequestBuilder provides operations to manage the solutionsRoot singleton.
type SolutionsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// SolutionsRequestBuilderGetQueryParameters get solutions
type SolutionsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SolutionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SolutionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SolutionsRequestBuilderGetQueryParameters
}
// SolutionsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SolutionsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BookingBusinesses provides operations to manage the bookingBusinesses property of the microsoft.graph.solutionsRoot entity.
func (m *SolutionsRequestBuilder) BookingBusinesses()(*ic2077c84f4cea40593f46a097344f757e71b10d42626a06348d8a6856a4a6408.BookingBusinessesRequestBuilder) {
    return ic2077c84f4cea40593f46a097344f757e71b10d42626a06348d8a6856a4a6408.NewBookingBusinessesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BookingBusinessesById provides operations to manage the bookingBusinesses property of the microsoft.graph.solutionsRoot entity.
func (m *SolutionsRequestBuilder) BookingBusinessesById(id string)(*i98532a0d43958ac46a84969d106d7877583a8e83cd026d67fbc826e462b1e8c7.BookingBusinessItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["bookingBusiness%2Did"] = id
    }
    return i98532a0d43958ac46a84969d106d7877583a8e83cd026d67fbc826e462b1e8c7.NewBookingBusinessItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// BookingCurrencies provides operations to manage the bookingCurrencies property of the microsoft.graph.solutionsRoot entity.
func (m *SolutionsRequestBuilder) BookingCurrencies()(*i6d5b6940d5e24f8aee3eb7eef5499c5a97c9c71485b644e5b94efadbc56609b0.BookingCurrenciesRequestBuilder) {
    return i6d5b6940d5e24f8aee3eb7eef5499c5a97c9c71485b644e5b94efadbc56609b0.NewBookingCurrenciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BookingCurrenciesById provides operations to manage the bookingCurrencies property of the microsoft.graph.solutionsRoot entity.
func (m *SolutionsRequestBuilder) BookingCurrenciesById(id string)(*i30929f1a7f5c848dd5ee1e9d3e9b1dc2926ba5d92f9f5ca889f41edee949d39a.BookingCurrencyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["bookingCurrency%2Did"] = id
    }
    return i30929f1a7f5c848dd5ee1e9d3e9b1dc2926ba5d92f9f5ca889f41edee949d39a.NewBookingCurrencyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewSolutionsRequestBuilderInternal instantiates a new SolutionsRequestBuilder and sets the default values.
func NewSolutionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SolutionsRequestBuilder) {
    m := &SolutionsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/solutions{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewSolutionsRequestBuilder instantiates a new SolutionsRequestBuilder and sets the default values.
func NewSolutionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SolutionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSolutionsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get solutions
func (m *SolutionsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *SolutionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update solutions
func (m *SolutionsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable, requestConfiguration *SolutionsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get solutions
func (m *SolutionsRequestBuilder) Get(ctx context.Context, requestConfiguration *SolutionsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSolutionsRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable), nil
}
// Patch update solutions
func (m *SolutionsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable, requestConfiguration *SolutionsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSolutionsRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SolutionsRootable), nil
}
