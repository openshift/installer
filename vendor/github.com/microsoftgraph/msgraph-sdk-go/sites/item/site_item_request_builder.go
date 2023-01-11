package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05bca8d67f58c6d9aa648acf49fac24e4d5361e451fb444c3868fc782faa3a3a "github.com/microsoftgraph/msgraph-sdk-go/sites/item/items"
    i0c26d90ad9ba1101a52431870fbb281c6a4f6b751415aa8750d1294df6fc9e2c "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes"
    i10cc803300c5d47903ff6a5a1eb5b34159f7eb7927f500c05b63f005fcd7a3b7 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/onenote"
    i1f30e998cc7f16739a41575c86813834123cba8d5d18b460b177e4d1033ca60d "github.com/microsoftgraph/msgraph-sdk-go/sites/item/drive"
    i339cee83997cfcbf197752f780095780d8a49c408ea09318bf2a269c76c80974 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/analytics"
    i367b3d6454445f7875bc897535d25edb404746299a8839bb9f87747de65f43fd "github.com/microsoftgraph/msgraph-sdk-go/sites/item/externalcolumns"
    i3ca97aef887804bcab28c6fc0da670ecf38816601cb616e2dd9231cb1cc3a18b "github.com/microsoftgraph/msgraph-sdk-go/sites/item/getapplicablecontenttypesforlistwithlistid"
    i3cd7ab6f3d06321d73344309d80fb03ef4dcccd5c3c20fc9e878db6dc8e33e51 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i4ffbbec24206153c780835b7a43a03ed6b772efe99d5161b43d3a4ba10302a0a "github.com/microsoftgraph/msgraph-sdk-go/sites/item/getbypathwithpath"
    i55cb8b4eba146099911d28cd6cbf1c2648018a786384c9ec7c40a8f9baf109f0 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/lists"
    i56d3808dd8956f908cd3292b0ec9af6c42c18061cbc621b53136a2fa051b6542 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/termstore"
    i5cf45aa21d0025c798a2868eb9d4091cba588f6ea645d747c04aba2a4591fe1b "github.com/microsoftgraph/msgraph-sdk-go/sites/item/permissions"
    i6dae79969c0cc01a1ae54240ed306a0bd1ad56e07a0515dc28cfab53e29bb580 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/sites"
    i8c3520c945a3499d24747186abcb37377ea0d8afa3d1bfa0cd8a7e6fcf9b2a83 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/termstores"
    i939b9e822b399f2c722f94a48564ad5eeba57ffb61b7a071a9ff1916605e56e8 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/getactivitiesbyinterval"
    ibce5e54be30545b916933fc1be11ed771de0b52c7e3fe58aea612d980394b5b5 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/operations"
    ic3e3392ffd4cfbeca5a615ab4821d60c9132bc0218c910f3bbb9967204677d15 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/drives"
    ifdd24c47f65a3313024bac4f3f3e730f886d4041edbb2510eda841cd16478038 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/columns"
    i209a587ff4beb5ee9c360ac0d6e4a3c7f62c54965cf6114ee02077fff23801c5 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/columns/item"
    i25b238e1bb04065b78bbf7b8af856ab5c032ddb3bf1f7e6a6dd2660de6f8217a "github.com/microsoftgraph/msgraph-sdk-go/sites/item/sites/item"
    i78f58c637039463a95013be84966d9b8b0483bfd65aa77130e076a7f7381cfd6 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/lists/item"
    i9bdc46bee3db40ed2ad5cb8d27cc47b07556e9f9bec575d1a403c32a5cce68ae "github.com/microsoftgraph/msgraph-sdk-go/sites/item/permissions/item"
    ia4734d1ec9f14b1703e1b22f66af2a51aea0b15668f690fb492176e062cae20e "github.com/microsoftgraph/msgraph-sdk-go/sites/item/termstores/item"
    ia5bb104926264b975572fae9aaa50ca289316f014c3b08eb8446548e5146dedb "github.com/microsoftgraph/msgraph-sdk-go/sites/item/drives/item"
    ib0bdd659558cc945118844ad67edcde17d87d06cec11209f3e35fd2264b9821c "github.com/microsoftgraph/msgraph-sdk-go/sites/item/externalcolumns/item"
    id1dc8bf51296de1e26b04da095d413b0cf7ad29e21afa5f4cd9ff1ede4a46807 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/operations/item"
    id5c99b5d474798a011325ba892c9f2ef398ac5ca9126710bc745e020c22b60fd "github.com/microsoftgraph/msgraph-sdk-go/sites/item/items/item"
    if642ca98ef1cfea5e51e56e879548fde41cf1297f3ce162d09a6224e659bc1f8 "github.com/microsoftgraph/msgraph-sdk-go/sites/item/contenttypes/item"
)

// SiteItemRequestBuilder provides operations to manage the collection of site entities.
type SiteItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// SiteItemRequestBuilderGetQueryParameters retrieve properties and relationships for a [site][] resource.A **site** resource represents a team site in SharePoint.
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
func (m *SiteItemRequestBuilder) Analytics()(*i339cee83997cfcbf197752f780095780d8a49c408ea09318bf2a269c76c80974.AnalyticsRequestBuilder) {
    return i339cee83997cfcbf197752f780095780d8a49c408ea09318bf2a269c76c80974.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Columns()(*ifdd24c47f65a3313024bac4f3f3e730f886d4041edbb2510eda841cd16478038.ColumnsRequestBuilder) {
    return ifdd24c47f65a3313024bac4f3f3e730f886d4041edbb2510eda841cd16478038.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ColumnsById(id string)(*i209a587ff4beb5ee9c360ac0d6e4a3c7f62c54965cf6114ee02077fff23801c5.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i209a587ff4beb5ee9c360ac0d6e4a3c7f62c54965cf6114ee02077fff23801c5.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewSiteItemRequestBuilderInternal instantiates a new SiteItemRequestBuilder and sets the default values.
func NewSiteItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SiteItemRequestBuilder) {
    m := &SiteItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/sites/{site%2Did}{?%24select,%24expand}";
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
func (m *SiteItemRequestBuilder) ContentTypes()(*i0c26d90ad9ba1101a52431870fbb281c6a4f6b751415aa8750d1294df6fc9e2c.ContentTypesRequestBuilder) {
    return i0c26d90ad9ba1101a52431870fbb281c6a4f6b751415aa8750d1294df6fc9e2c.NewContentTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContentTypesById provides operations to manage the contentTypes property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ContentTypesById(id string)(*if642ca98ef1cfea5e51e56e879548fde41cf1297f3ce162d09a6224e659bc1f8.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did"] = id
    }
    return if642ca98ef1cfea5e51e56e879548fde41cf1297f3ce162d09a6224e659bc1f8.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CreateGetRequestInformation retrieve properties and relationships for a [site][] resource.A **site** resource represents a team site in SharePoint.
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
// CreatePatchRequestInformation update entity in sites
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
func (m *SiteItemRequestBuilder) Drive()(*i1f30e998cc7f16739a41575c86813834123cba8d5d18b460b177e4d1033ca60d.DriveRequestBuilder) {
    return i1f30e998cc7f16739a41575c86813834123cba8d5d18b460b177e4d1033ca60d.NewDriveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Drives provides operations to manage the drives property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Drives()(*ic3e3392ffd4cfbeca5a615ab4821d60c9132bc0218c910f3bbb9967204677d15.DrivesRequestBuilder) {
    return ic3e3392ffd4cfbeca5a615ab4821d60c9132bc0218c910f3bbb9967204677d15.NewDrivesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DrivesById provides operations to manage the drives property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) DrivesById(id string)(*ia5bb104926264b975572fae9aaa50ca289316f014c3b08eb8446548e5146dedb.DriveItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["drive%2Did"] = id
    }
    return ia5bb104926264b975572fae9aaa50ca289316f014c3b08eb8446548e5146dedb.NewDriveItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ExternalColumns provides operations to manage the externalColumns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ExternalColumns()(*i367b3d6454445f7875bc897535d25edb404746299a8839bb9f87747de65f43fd.ExternalColumnsRequestBuilder) {
    return i367b3d6454445f7875bc897535d25edb404746299a8839bb9f87747de65f43fd.NewExternalColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExternalColumnsById provides operations to manage the externalColumns property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ExternalColumnsById(id string)(*ib0bdd659558cc945118844ad67edcde17d87d06cec11209f3e35fd2264b9821c.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return ib0bdd659558cc945118844ad67edcde17d87d06cec11209f3e35fd2264b9821c.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get retrieve properties and relationships for a [site][] resource.A **site** resource represents a team site in SharePoint.
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
func (m *SiteItemRequestBuilder) GetActivitiesByInterval()(*i939b9e822b399f2c722f94a48564ad5eeba57ffb61b7a071a9ff1916605e56e8.GetActivitiesByIntervalRequestBuilder) {
    return i939b9e822b399f2c722f94a48564ad5eeba57ffb61b7a071a9ff1916605e56e8.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *SiteItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i3cd7ab6f3d06321d73344309d80fb03ef4dcccd5c3c20fc9e878db6dc8e33e51.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i3cd7ab6f3d06321d73344309d80fb03ef4dcccd5c3c20fc9e878db6dc8e33e51.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// GetApplicableContentTypesForListWithListId provides operations to call the getApplicableContentTypesForList method.
func (m *SiteItemRequestBuilder) GetApplicableContentTypesForListWithListId(listId *string)(*i3ca97aef887804bcab28c6fc0da670ecf38816601cb616e2dd9231cb1cc3a18b.GetApplicableContentTypesForListWithListIdRequestBuilder) {
    return i3ca97aef887804bcab28c6fc0da670ecf38816601cb616e2dd9231cb1cc3a18b.NewGetApplicableContentTypesForListWithListIdRequestBuilderInternal(m.pathParameters, m.requestAdapter, listId);
}
// GetByPathWithPath provides operations to call the getByPath method.
func (m *SiteItemRequestBuilder) GetByPathWithPath(path *string)(*i4ffbbec24206153c780835b7a43a03ed6b772efe99d5161b43d3a4ba10302a0a.GetByPathWithPathRequestBuilder) {
    return i4ffbbec24206153c780835b7a43a03ed6b772efe99d5161b43d3a4ba10302a0a.NewGetByPathWithPathRequestBuilderInternal(m.pathParameters, m.requestAdapter, path);
}
// Items provides operations to manage the items property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Items()(*i05bca8d67f58c6d9aa648acf49fac24e4d5361e451fb444c3868fc782faa3a3a.ItemsRequestBuilder) {
    return i05bca8d67f58c6d9aa648acf49fac24e4d5361e451fb444c3868fc782faa3a3a.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ItemsById(id string)(*id5c99b5d474798a011325ba892c9f2ef398ac5ca9126710bc745e020c22b60fd.BaseItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["baseItem%2Did"] = id
    }
    return id5c99b5d474798a011325ba892c9f2ef398ac5ca9126710bc745e020c22b60fd.NewBaseItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Lists provides operations to manage the lists property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Lists()(*i55cb8b4eba146099911d28cd6cbf1c2648018a786384c9ec7c40a8f9baf109f0.ListsRequestBuilder) {
    return i55cb8b4eba146099911d28cd6cbf1c2648018a786384c9ec7c40a8f9baf109f0.NewListsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListsById provides operations to manage the lists property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) ListsById(id string)(*i78f58c637039463a95013be84966d9b8b0483bfd65aa77130e076a7f7381cfd6.ListItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["list%2Did"] = id
    }
    return i78f58c637039463a95013be84966d9b8b0483bfd65aa77130e076a7f7381cfd6.NewListItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Onenote()(*i10cc803300c5d47903ff6a5a1eb5b34159f7eb7927f500c05b63f005fcd7a3b7.OnenoteRequestBuilder) {
    return i10cc803300c5d47903ff6a5a1eb5b34159f7eb7927f500c05b63f005fcd7a3b7.NewOnenoteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Operations()(*ibce5e54be30545b916933fc1be11ed771de0b52c7e3fe58aea612d980394b5b5.OperationsRequestBuilder) {
    return ibce5e54be30545b916933fc1be11ed771de0b52c7e3fe58aea612d980394b5b5.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) OperationsById(id string)(*id1dc8bf51296de1e26b04da095d413b0cf7ad29e21afa5f4cd9ff1ede4a46807.RichLongRunningOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["richLongRunningOperation%2Did"] = id
    }
    return id1dc8bf51296de1e26b04da095d413b0cf7ad29e21afa5f4cd9ff1ede4a46807.NewRichLongRunningOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update entity in sites
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
func (m *SiteItemRequestBuilder) Permissions()(*i5cf45aa21d0025c798a2868eb9d4091cba588f6ea645d747c04aba2a4591fe1b.PermissionsRequestBuilder) {
    return i5cf45aa21d0025c798a2868eb9d4091cba588f6ea645d747c04aba2a4591fe1b.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) PermissionsById(id string)(*i9bdc46bee3db40ed2ad5cb8d27cc47b07556e9f9bec575d1a403c32a5cce68ae.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i9bdc46bee3db40ed2ad5cb8d27cc47b07556e9f9bec575d1a403c32a5cce68ae.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Sites provides operations to manage the sites property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) Sites()(*i6dae79969c0cc01a1ae54240ed306a0bd1ad56e07a0515dc28cfab53e29bb580.SitesRequestBuilder) {
    return i6dae79969c0cc01a1ae54240ed306a0bd1ad56e07a0515dc28cfab53e29bb580.NewSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SitesById provides operations to manage the sites property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) SitesById(id string)(*i25b238e1bb04065b78bbf7b8af856ab5c032ddb3bf1f7e6a6dd2660de6f8217a.SiteItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["site%2Did1"] = id
    }
    return i25b238e1bb04065b78bbf7b8af856ab5c032ddb3bf1f7e6a6dd2660de6f8217a.NewSiteItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TermStore provides operations to manage the termStore property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStore()(*i56d3808dd8956f908cd3292b0ec9af6c42c18061cbc621b53136a2fa051b6542.TermStoreRequestBuilder) {
    return i56d3808dd8956f908cd3292b0ec9af6c42c18061cbc621b53136a2fa051b6542.NewTermStoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TermStores provides operations to manage the termStores property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStores()(*i8c3520c945a3499d24747186abcb37377ea0d8afa3d1bfa0cd8a7e6fcf9b2a83.TermStoresRequestBuilder) {
    return i8c3520c945a3499d24747186abcb37377ea0d8afa3d1bfa0cd8a7e6fcf9b2a83.NewTermStoresRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TermStoresById provides operations to manage the termStores property of the microsoft.graph.site entity.
func (m *SiteItemRequestBuilder) TermStoresById(id string)(*ia4734d1ec9f14b1703e1b22f66af2a51aea0b15668f690fb492176e062cae20e.StoreItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["store%2Did"] = id
    }
    return ia4734d1ec9f14b1703e1b22f66af2a51aea0b15668f690fb492176e062cae20e.NewStoreItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
