package outlook

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2d6c0e9680c3153270f78e152e34f25ac0be843ae551634f2ef6b57d69f361ea "github.com/microsoftgraph/msgraph-sdk-go/me/outlook/supportedtimezoneswithtimezonestandard"
    i4763caa20937545d0f4b4141b0840589720e6ef619737a593826e0ad669072e6 "github.com/microsoftgraph/msgraph-sdk-go/me/outlook/supportedtimezones"
    icc3e92b7111767420ec0be836502793f8a368fa468daf7fdf8e65cc66f0e087f "github.com/microsoftgraph/msgraph-sdk-go/me/outlook/mastercategories"
    iefb515d7c64e31fd65a25d5ca6cc93380fa4b1c860b5fda47701b19a3b2226f8 "github.com/microsoftgraph/msgraph-sdk-go/me/outlook/supportedlanguages"
    ifa2143032fa73c15f7bf7cc659bb2bb69f28a378f9584ec6f7ccf4d2b3d29868 "github.com/microsoftgraph/msgraph-sdk-go/me/outlook/mastercategories/item"
)

// OutlookRequestBuilder provides operations to manage the outlook property of the microsoft.graph.user entity.
type OutlookRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OutlookRequestBuilderGetQueryParameters get outlook from me
type OutlookRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// OutlookRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OutlookRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OutlookRequestBuilderGetQueryParameters
}
// NewOutlookRequestBuilderInternal instantiates a new OutlookRequestBuilder and sets the default values.
func NewOutlookRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OutlookRequestBuilder) {
    m := &OutlookRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/outlook{?%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOutlookRequestBuilder instantiates a new OutlookRequestBuilder and sets the default values.
func NewOutlookRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OutlookRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOutlookRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get outlook from me
func (m *OutlookRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OutlookRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get outlook from me
func (m *OutlookRequestBuilder) Get(ctx context.Context, requestConfiguration *OutlookRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OutlookUserable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOutlookUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OutlookUserable), nil
}
// MasterCategories provides operations to manage the masterCategories property of the microsoft.graph.outlookUser entity.
func (m *OutlookRequestBuilder) MasterCategories()(*icc3e92b7111767420ec0be836502793f8a368fa468daf7fdf8e65cc66f0e087f.MasterCategoriesRequestBuilder) {
    return icc3e92b7111767420ec0be836502793f8a368fa468daf7fdf8e65cc66f0e087f.NewMasterCategoriesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MasterCategoriesById provides operations to manage the masterCategories property of the microsoft.graph.outlookUser entity.
func (m *OutlookRequestBuilder) MasterCategoriesById(id string)(*ifa2143032fa73c15f7bf7cc659bb2bb69f28a378f9584ec6f7ccf4d2b3d29868.OutlookCategoryItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["outlookCategory%2Did"] = id
    }
    return ifa2143032fa73c15f7bf7cc659bb2bb69f28a378f9584ec6f7ccf4d2b3d29868.NewOutlookCategoryItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SupportedLanguages provides operations to call the supportedLanguages method.
func (m *OutlookRequestBuilder) SupportedLanguages()(*iefb515d7c64e31fd65a25d5ca6cc93380fa4b1c860b5fda47701b19a3b2226f8.SupportedLanguagesRequestBuilder) {
    return iefb515d7c64e31fd65a25d5ca6cc93380fa4b1c860b5fda47701b19a3b2226f8.NewSupportedLanguagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SupportedTimeZones provides operations to call the supportedTimeZones method.
func (m *OutlookRequestBuilder) SupportedTimeZones()(*i4763caa20937545d0f4b4141b0840589720e6ef619737a593826e0ad669072e6.SupportedTimeZonesRequestBuilder) {
    return i4763caa20937545d0f4b4141b0840589720e6ef619737a593826e0ad669072e6.NewSupportedTimeZonesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SupportedTimeZonesWithTimeZoneStandard provides operations to call the supportedTimeZones method.
func (m *OutlookRequestBuilder) SupportedTimeZonesWithTimeZoneStandard(timeZoneStandard *string)(*i2d6c0e9680c3153270f78e152e34f25ac0be843ae551634f2ef6b57d69f361ea.SupportedTimeZonesWithTimeZoneStandardRequestBuilder) {
    return i2d6c0e9680c3153270f78e152e34f25ac0be843ae551634f2ef6b57d69f361ea.NewSupportedTimeZonesWithTimeZoneStandardRequestBuilderInternal(m.pathParameters, m.requestAdapter, timeZoneStandard);
}
