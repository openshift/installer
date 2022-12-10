package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i055ad2f09c6431b30e6d3528a3dd11b667dcf169a3943e92ebbdc9f40c102910 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/directreports"
    i15c3213b1d34c79d320391089a05683e39eb479cb5d4b030a5b53be13ac05712 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/getmembergroups"
    i4830f34009f28f830e2d6c1e84624d4bbfc5a6d035e6f1097558162d8a77592d "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof"
    i56f89ff65e41535b67a714e8fdea36f3ad34961ef4ad2e4c91a286c72ea09b44 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/restore"
    i6a99a368f031e6c1f338c5fc4bcbdad85d242df8d3d5dfe1b7b0475df468b3dd "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/manager"
    i79406b1f55eac2f939b44b7a8f42f4c1fcc6945c0eb2e6ff7fb5ee4c06796b4d "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof"
    ib08660edc0c07cb1bfb954212b3de0dfac964791cee0a520dea1ff1fc58ce970 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/checkmemberobjects"
    id7379fe1d66c0786481cb810e571e9cf8ca3a1db19fb8596ae3ab26b168ebd4b "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/checkmembergroups"
    ifff3f2949bd5013659a326ca1cde357b44996411a8f259036c89d4ee51051f34 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/getmemberobjects"
    i10b5e6a502bf4e0294bb5060a3060bddfc4c78df75ba9fbb4e0b2c938e4dbcda "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item"
    i334f74d58761225d35a308cd57cea9169133186ab28a294f59b928e984bfa1c4 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/directreports/item"
    if2622610c8389e13476d20088b38942700590b9eea0b48a2c9847786361c95cd "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/item"
)

// OrgContactItemRequestBuilder provides operations to manage the collection of orgContact entities.
type OrgContactItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OrgContactItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OrgContactItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// OrgContactItemRequestBuilderGetQueryParameters get the properties and relationships of an organizational contact.
type OrgContactItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// OrgContactItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OrgContactItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OrgContactItemRequestBuilderGetQueryParameters
}
// OrgContactItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OrgContactItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
func (m *OrgContactItemRequestBuilder) CheckMemberGroups()(*id7379fe1d66c0786481cb810e571e9cf8ca3a1db19fb8596ae3ab26b168ebd4b.CheckMemberGroupsRequestBuilder) {
    return id7379fe1d66c0786481cb810e571e9cf8ca3a1db19fb8596ae3ab26b168ebd4b.NewCheckMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
func (m *OrgContactItemRequestBuilder) CheckMemberObjects()(*ib08660edc0c07cb1bfb954212b3de0dfac964791cee0a520dea1ff1fc58ce970.CheckMemberObjectsRequestBuilder) {
    return ib08660edc0c07cb1bfb954212b3de0dfac964791cee0a520dea1ff1fc58ce970.NewCheckMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOrgContactItemRequestBuilderInternal instantiates a new OrgContactItemRequestBuilder and sets the default values.
func NewOrgContactItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OrgContactItemRequestBuilder) {
    m := &OrgContactItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/contacts/{orgContact%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOrgContactItemRequestBuilder instantiates a new OrgContactItemRequestBuilder and sets the default values.
func NewOrgContactItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OrgContactItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOrgContactItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete entity from contacts
func (m *OrgContactItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *OrgContactItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get the properties and relationships of an organizational contact.
func (m *OrgContactItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OrgContactItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update entity in contacts
func (m *OrgContactItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, requestConfiguration *OrgContactItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete entity from contacts
func (m *OrgContactItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *OrgContactItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DirectReports provides operations to manage the directReports property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) DirectReports()(*i055ad2f09c6431b30e6d3528a3dd11b667dcf169a3943e92ebbdc9f40c102910.DirectReportsRequestBuilder) {
    return i055ad2f09c6431b30e6d3528a3dd11b667dcf169a3943e92ebbdc9f40c102910.NewDirectReportsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectReportsById provides operations to manage the directReports property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) DirectReportsById(id string)(*i334f74d58761225d35a308cd57cea9169133186ab28a294f59b928e984bfa1c4.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i334f74d58761225d35a308cd57cea9169133186ab28a294f59b928e984bfa1c4.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get get the properties and relationships of an organizational contact.
func (m *OrgContactItemRequestBuilder) Get(ctx context.Context, requestConfiguration *OrgContactItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOrgContactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable), nil
}
// GetMemberGroups provides operations to call the getMemberGroups method.
func (m *OrgContactItemRequestBuilder) GetMemberGroups()(*i15c3213b1d34c79d320391089a05683e39eb479cb5d4b030a5b53be13ac05712.GetMemberGroupsRequestBuilder) {
    return i15c3213b1d34c79d320391089a05683e39eb479cb5d4b030a5b53be13ac05712.NewGetMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetMemberObjects provides operations to call the getMemberObjects method.
func (m *OrgContactItemRequestBuilder) GetMemberObjects()(*ifff3f2949bd5013659a326ca1cde357b44996411a8f259036c89d4ee51051f34.GetMemberObjectsRequestBuilder) {
    return ifff3f2949bd5013659a326ca1cde357b44996411a8f259036c89d4ee51051f34.NewGetMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Manager provides operations to manage the manager property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) Manager()(*i6a99a368f031e6c1f338c5fc4bcbdad85d242df8d3d5dfe1b7b0475df468b3dd.ManagerRequestBuilder) {
    return i6a99a368f031e6c1f338c5fc4bcbdad85d242df8d3d5dfe1b7b0475df468b3dd.NewManagerRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MemberOf provides operations to manage the memberOf property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) MemberOf()(*i4830f34009f28f830e2d6c1e84624d4bbfc5a6d035e6f1097558162d8a77592d.MemberOfRequestBuilder) {
    return i4830f34009f28f830e2d6c1e84624d4bbfc5a6d035e6f1097558162d8a77592d.NewMemberOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MemberOfById provides operations to manage the memberOf property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) MemberOfById(id string)(*if2622610c8389e13476d20088b38942700590b9eea0b48a2c9847786361c95cd.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return if2622610c8389e13476d20088b38942700590b9eea0b48a2c9847786361c95cd.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update entity in contacts
func (m *OrgContactItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, requestConfiguration *OrgContactItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOrgContactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable), nil
}
// Restore provides operations to call the restore method.
func (m *OrgContactItemRequestBuilder) Restore()(*i56f89ff65e41535b67a714e8fdea36f3ad34961ef4ad2e4c91a286c72ea09b44.RestoreRequestBuilder) {
    return i56f89ff65e41535b67a714e8fdea36f3ad34961ef4ad2e4c91a286c72ea09b44.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TransitiveMemberOf provides operations to manage the transitiveMemberOf property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) TransitiveMemberOf()(*i79406b1f55eac2f939b44b7a8f42f4c1fcc6945c0eb2e6ff7fb5ee4c06796b4d.TransitiveMemberOfRequestBuilder) {
    return i79406b1f55eac2f939b44b7a8f42f4c1fcc6945c0eb2e6ff7fb5ee4c06796b4d.NewTransitiveMemberOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TransitiveMemberOfById provides operations to manage the transitiveMemberOf property of the microsoft.graph.orgContact entity.
func (m *OrgContactItemRequestBuilder) TransitiveMemberOfById(id string)(*i10b5e6a502bf4e0294bb5060a3060bddfc4c78df75ba9fbb4e0b2c938e4dbcda.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i10b5e6a502bf4e0294bb5060a3060bddfc4c78df75ba9fbb4e0b2c938e4dbcda.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
