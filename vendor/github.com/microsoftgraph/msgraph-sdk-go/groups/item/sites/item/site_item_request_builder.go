package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i19d3b905956d0d5ab3e405e912ce5512047d11c0368cfb7b3532e4e86e5a07af "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/termstores"
    i23ed04a56449d98fcd568705a0cf04a1b544ee9821f049f84ae0c25b1d8a1f76 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/items"
    i433a1560954e7a6c21be030d3739b5c8e7a3a62b25192a404e6284478dbdf361 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/columns"
    i4530be9f2649e2da7651adc83045d48e03ce2d02fe5bbfac583977b13834180c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists"
    i5ff933dfab28b9bc528e49bf8331befadb5004877c8387d60f5e76079d8ba1a6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/externalcolumns"
    i70eefd6cd3bf5a38a63273bd720e33374f89e72e616be1625b28a3121ac42d7e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/getactivitiesbyinterval"
    i7bf3dbc49a37a69a92d9c0b2f5894dbd2387cd5cf2c265de0a9f1307dedf7007 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/onenote"
    i9619fa512c16d940bb8a99f87dcbb11e238865522dd9f8aa27e8eb6839d210c8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/sites"
    iabe6a78fc3f592fa955e512cf883dcc83f9a1c269265fa5904300f3f3c66d279 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/termstore"
    ibcb8a11a874b2cf73574fef772774539d24c59dca1bdbe499cc36aa8f140653b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    ibe39270ba4e624ff4fe340aa2bd6d3631b4dd3b549fc6c5dbafcf7756c2cae95 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/permissions"
    ica4ea80f6f6f649de989441fa87ef61c9a4e8922191481297e55c79e72481491 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/operations"
    id7f3688a19ad77598caac412c4d52e7de1d4e1280b3475235ac9cd72911b8374 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/getbypathwithpath"
    ida62d1e3e0c3a4155187312adb160423df6f1b84474a7c37dac6b38ee0d57c30 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/contenttypes"
    ie2e0fba360015570188988834e1d1616938882b51276657b518960b0529c8395 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/drive"
    ie48cea0f35857c6390d80a8f786c40458fa4f1b697426e534a5847736966c883 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/getapplicablecontenttypesforlistwithlistid"
    ie977d634115f6edbcba02945f5f558e1f1060c3078fc1043cfc6a70da5598048 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/analytics"
    if423f3b786c94937b52a1c87021c6cf25ba0bba2b0405608a8d1629f108d2664 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/drives"
    i20d71e187e2869b1bf133549cf92b03dbdd4fcd7a6312d7e4cfab185397ec504 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/drives/item"
    i3d15a6f13c74ebb900c3f4d7e8b8cf4dc95de03aed0a0c4b8a2a4e6ed582e97e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/columns/item"
    i3eb82c303a8d1465dafb4c1bc15b34ff40013f94aac87f058c363a7aba92e4e8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/externalcolumns/item"
    i41ed64eca1f91fe55946dd46e765c7ec66bd75b7a8e2a22aa52490bb8b7322de "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/operations/item"
    i542a41cec057b321ef1e41187fa571edf863979b732f2c19975560ced8e76f8b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/permissions/item"
    ia63bedd538371b0186a5b6eee2c0d77c2ae0f708dcde2fcdeeebb39fe26bd8b7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/termstores/item"
    iba2aa0555939eae7c495b9747ca02d399d5a9d9768108b185661cc038711fdbc "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/items/item"
    ic3cb2a24c56a0f014badb537b2e05366face417e10e73eda79e507ce267e5183 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/contenttypes/item"
    icb7ef07a689e39fb1ed0e6596004ecf47a3815ad87f19c96bff8fb7ad90304b9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/sites/item"
    iea900cf48c5d3da0d6a702e72cfb722443dcaaa70a27924a2c1a8b756a55448d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/sites/item/lists/item"
)

// SiteItemRequestBuilder provides operations to manage the sites property of the microsoft.graph.group entity.
type SiteItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// SiteItemRequestBuilderGetQueryParameters the list of SharePoint sites in this group. Access the default site with /sites/root.
type SiteItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SiteItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SiteItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SiteItemRequestBuilderGetQueryParameters
}
// SiteItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SiteItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Analytics()(*ie977d634115f6edbcba02945f5f558e1f1060c3078fc1043cfc6a70da5598048.AnalyticsRequestBuilder) {
    return ie977d634115f6edbcba02945f5f558e1f1060c3078fc1043cfc6a70da5598048.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Columns()(*i433a1560954e7a6c21be030d3739b5c8e7a3a62b25192a404e6284478dbdf361.ColumnsRequestBuilder) {
    return i433a1560954e7a6c21be030d3739b5c8e7a3a62b25192a404e6284478dbdf361.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ColumnsById(id string)(*i3d15a6f13c74ebb900c3f4d7e8b8cf4dc95de03aed0a0c4b8a2a4e6ed582e97e.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i3d15a6f13c74ebb900c3f4d7e8b8cf4dc95de03aed0a0c4b8a2a4e6ed582e97e.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewSiteItemRequestBuilderInternal instantiates a new SiteItemRequestBuilder and sets the default values.
func NewSiteItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SiteItemRequestBuilder) {
    m := &SiteItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewSiteItemRequestBuilder instantiates a new SiteItemRequestBuilder and sets the default values.
func NewSiteItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SiteItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSiteItemRequestBuilderInternal(urlParams, requestAdapter)
}
// ContentTypes provides operations to manage the contentTypes property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ContentTypes()(*ida62d1e3e0c3a4155187312adb160423df6f1b84474a7c37dac6b38ee0d57c30.ContentTypesRequestBuilder) {
    return ida62d1e3e0c3a4155187312adb160423df6f1b84474a7c37dac6b38ee0d57c30.NewContentTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContentTypesById provides operations to manage the contentTypes property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ContentTypesById(id string)(*ic3cb2a24c56a0f014badb537b2e05366face417e10e73eda79e507ce267e5183.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did"] = id
    }
    return ic3cb2a24c56a0f014badb537b2e05366face417e10e73eda79e507ce267e5183.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CreateGetRequestInformation the list of SharePoint sites in this group. Access the default site with /sites/root.
func (m *SiteItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *SiteItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property sites in groups
func (m *SiteItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, requestConfiguration *SiteItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Drive provides operations to manage the drive property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Drive()(*ie2e0fba360015570188988834e1d1616938882b51276657b518960b0529c8395.DriveRequestBuilder) {
    return ie2e0fba360015570188988834e1d1616938882b51276657b518960b0529c8395.NewDriveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Drives provides operations to manage the drives property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Drives()(*if423f3b786c94937b52a1c87021c6cf25ba0bba2b0405608a8d1629f108d2664.DrivesRequestBuilder) {
    return if423f3b786c94937b52a1c87021c6cf25ba0bba2b0405608a8d1629f108d2664.NewDrivesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DrivesById provides operations to manage the drives property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) DrivesById(id string)(*i20d71e187e2869b1bf133549cf92b03dbdd4fcd7a6312d7e4cfab185397ec504.DriveItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["drive%2Did"] = id
    }
    return i20d71e187e2869b1bf133549cf92b03dbdd4fcd7a6312d7e4cfab185397ec504.NewDriveItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ExternalColumns provides operations to manage the externalColumns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ExternalColumns()(*i5ff933dfab28b9bc528e49bf8331befadb5004877c8387d60f5e76079d8ba1a6.ExternalColumnsRequestBuilder) {
    return i5ff933dfab28b9bc528e49bf8331befadb5004877c8387d60f5e76079d8ba1a6.NewExternalColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExternalColumnsById provides operations to manage the externalColumns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ExternalColumnsById(id string)(*i3eb82c303a8d1465dafb4c1bc15b34ff40013f94aac87f058c363a7aba92e4e8.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i3eb82c303a8d1465dafb4c1bc15b34ff40013f94aac87f058c363a7aba92e4e8.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get the list of SharePoint sites in this group. Access the default site with /sites/root.
func (m *SiteItemRequestBuilder) Get(ctx context.Context, requestConfiguration *SiteItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable), nil
}
// GetActivitiesByInterval provides operations to call the getActivitiesByInterval method.
func (m *SiteItemRequestBuilder) GetActivitiesByInterval()(*i70eefd6cd3bf5a38a63273bd720e33374f89e72e616be1625b28a3121ac42d7e.GetActivitiesByIntervalRequestBuilder) {
    return i70eefd6cd3bf5a38a63273bd720e33374f89e72e616be1625b28a3121ac42d7e.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *SiteItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*ibcb8a11a874b2cf73574fef772774539d24c59dca1bdbe499cc36aa8f140653b.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return ibcb8a11a874b2cf73574fef772774539d24c59dca1bdbe499cc36aa8f140653b.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// GetApplicableContentTypesForListWithListId provides operations to call the getApplicableContentTypesForList method.
func (m *SiteItemRequestBuilder) GetApplicableContentTypesForListWithListId(listId *string)(*ie48cea0f35857c6390d80a8f786c40458fa4f1b697426e534a5847736966c883.GetApplicableContentTypesForListWithListIdRequestBuilder) {
    return ie48cea0f35857c6390d80a8f786c40458fa4f1b697426e534a5847736966c883.NewGetApplicableContentTypesForListWithListIdRequestBuilderInternal(m.pathParameters, m.requestAdapter, listId);
}
// GetByPathWithPath provides operations to call the getByPath method.
func (m *SiteItemRequestBuilder) GetByPathWithPath(path *string)(*id7f3688a19ad77598caac412c4d52e7de1d4e1280b3475235ac9cd72911b8374.GetByPathWithPathRequestBuilder) {
    return id7f3688a19ad77598caac412c4d52e7de1d4e1280b3475235ac9cd72911b8374.NewGetByPathWithPathRequestBuilderInternal(m.pathParameters, m.requestAdapter, path);
}
// Items provides operations to manage the items property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Items()(*i23ed04a56449d98fcd568705a0cf04a1b544ee9821f049f84ae0c25b1d8a1f76.ItemsRequestBuilder) {
    return i23ed04a56449d98fcd568705a0cf04a1b544ee9821f049f84ae0c25b1d8a1f76.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ItemsById(id string)(*iba2aa0555939eae7c495b9747ca02d399d5a9d9768108b185661cc038711fdbc.BaseItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["baseItem%2Did"] = id
    }
    return iba2aa0555939eae7c495b9747ca02d399d5a9d9768108b185661cc038711fdbc.NewBaseItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Lists provides operations to manage the lists property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Lists()(*i4530be9f2649e2da7651adc83045d48e03ce2d02fe5bbfac583977b13834180c.ListsRequestBuilder) {
    return i4530be9f2649e2da7651adc83045d48e03ce2d02fe5bbfac583977b13834180c.NewListsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListsById provides operations to manage the lists property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ListsById(id string)(*iea900cf48c5d3da0d6a702e72cfb722443dcaaa70a27924a2c1a8b756a55448d.ListItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["list%2Did"] = id
    }
    return iea900cf48c5d3da0d6a702e72cfb722443dcaaa70a27924a2c1a8b756a55448d.NewListItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Onenote()(*i7bf3dbc49a37a69a92d9c0b2f5894dbd2387cd5cf2c265de0a9f1307dedf7007.OnenoteRequestBuilder) {
    return i7bf3dbc49a37a69a92d9c0b2f5894dbd2387cd5cf2c265de0a9f1307dedf7007.NewOnenoteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Operations()(*ica4ea80f6f6f649de989441fa87ef61c9a4e8922191481297e55c79e72481491.OperationsRequestBuilder) {
    return ica4ea80f6f6f649de989441fa87ef61c9a4e8922191481297e55c79e72481491.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) OperationsById(id string)(*i41ed64eca1f91fe55946dd46e765c7ec66bd75b7a8e2a22aa52490bb8b7322de.RichLongRunningOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["richLongRunningOperation%2Did"] = id
    }
    return i41ed64eca1f91fe55946dd46e765c7ec66bd75b7a8e2a22aa52490bb8b7322de.NewRichLongRunningOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property sites in groups
func (m *SiteItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, requestConfiguration *SiteItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSiteFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Siteable), nil
}
// Permissions provides operations to manage the permissions property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Permissions()(*ibe39270ba4e624ff4fe340aa2bd6d3631b4dd3b549fc6c5dbafcf7756c2cae95.PermissionsRequestBuilder) {
    return ibe39270ba4e624ff4fe340aa2bd6d3631b4dd3b549fc6c5dbafcf7756c2cae95.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) PermissionsById(id string)(*i542a41cec057b321ef1e41187fa571edf863979b732f2c19975560ced8e76f8b.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i542a41cec057b321ef1e41187fa571edf863979b732f2c19975560ced8e76f8b.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Sites provides operations to manage the sites property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Sites()(*i9619fa512c16d940bb8a99f87dcbb11e238865522dd9f8aa27e8eb6839d210c8.SitesRequestBuilder) {
    return i9619fa512c16d940bb8a99f87dcbb11e238865522dd9f8aa27e8eb6839d210c8.NewSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SitesById provides operations to manage the sites property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) SitesById(id string)(*icb7ef07a689e39fb1ed0e6596004ecf47a3815ad87f19c96bff8fb7ad90304b9.SiteItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["site%2Did1"] = id
    }
    return icb7ef07a689e39fb1ed0e6596004ecf47a3815ad87f19c96bff8fb7ad90304b9.NewSiteItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TermStore provides operations to manage the termStore property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStore()(*iabe6a78fc3f592fa955e512cf883dcc83f9a1c269265fa5904300f3f3c66d279.TermStoreRequestBuilder) {
    return iabe6a78fc3f592fa955e512cf883dcc83f9a1c269265fa5904300f3f3c66d279.NewTermStoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TermStores provides operations to manage the termStores property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStores()(*i19d3b905956d0d5ab3e405e912ce5512047d11c0368cfb7b3532e4e86e5a07af.TermStoresRequestBuilder) {
    return i19d3b905956d0d5ab3e405e912ce5512047d11c0368cfb7b3532e4e86e5a07af.NewTermStoresRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TermStoresById provides operations to manage the termStores property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStoresById(id string)(*ia63bedd538371b0186a5b6eee2c0d77c2ae0f708dcde2fcdeeebb39fe26bd8b7.StoreItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["store%2Did"] = id
    }
    return ia63bedd538371b0186a5b6eee2c0d77c2ae0f708dcde2fcdeeebb39fe26bd8b7.NewStoreItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
