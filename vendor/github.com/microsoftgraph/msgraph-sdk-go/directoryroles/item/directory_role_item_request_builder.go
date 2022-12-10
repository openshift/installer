package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0b751179e4031051436a326b2702d4abde4bb19175d1daaf554c829a18c29177 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/checkmemberobjects"
    i718fb7a72b38652faad082566b31b033df2c9f8d04cc0af99ee1e607f87e8d80 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/scopedmembers"
    iaeb18d19aec9445180dbe3b2ba9b76bbfbcf9731777c886b4f0fa5501b97e58c "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/checkmembergroups"
    ibcb647a3bce4367afc71872d9ca7c9cc64edcc4e60fe4831e8fcab3de590a294 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/restore"
    id49876e5f70243d663a4f93ec6d77d7f89f0c91f1eda2daad48f4f037781d7b0 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/getmemberobjects"
    id63def69f45b85c0f48fc2880b865fc141de2bfccf786d2d6a9813abb7b316e9 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/getmembergroups"
    if1385e9a44279e575f64436ffd1ff5c93bc83fbaa5a4051782729fbe1d3bdb3d "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members"
    i4dc00e738a43cdc3dbc464cba96748a3f629a6e7f36a76cbb15bf379872ebba9 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item"
    ic73a5746a64610f345c78c4114780b56c25d27fdb8734dbd3b6500ebf1ad197c "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/scopedmembers/item"
)

// DirectoryRoleItemRequestBuilder provides operations to manage the collection of directoryRole entities.
type DirectoryRoleItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryRoleItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRoleItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DirectoryRoleItemRequestBuilderGetQueryParameters retrieve the properties of a directoryRole object. The role must be activated in tenant for a successful response. You can use both the object ID and template ID of the **directoryRole** with this API. The template ID of a built-in role is immutable and can be seen in the role description on the Azure portal. For details, see Role template IDs.
type DirectoryRoleItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DirectoryRoleItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRoleItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DirectoryRoleItemRequestBuilderGetQueryParameters
}
// DirectoryRoleItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRoleItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
func (m *DirectoryRoleItemRequestBuilder) CheckMemberGroups()(*iaeb18d19aec9445180dbe3b2ba9b76bbfbcf9731777c886b4f0fa5501b97e58c.CheckMemberGroupsRequestBuilder) {
    return iaeb18d19aec9445180dbe3b2ba9b76bbfbcf9731777c886b4f0fa5501b97e58c.NewCheckMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
func (m *DirectoryRoleItemRequestBuilder) CheckMemberObjects()(*i0b751179e4031051436a326b2702d4abde4bb19175d1daaf554c829a18c29177.CheckMemberObjectsRequestBuilder) {
    return i0b751179e4031051436a326b2702d4abde4bb19175d1daaf554c829a18c29177.NewCheckMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryRoleItemRequestBuilderInternal instantiates a new DirectoryRoleItemRequestBuilder and sets the default values.
func NewDirectoryRoleItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryRoleItemRequestBuilder) {
    m := &DirectoryRoleItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directoryRoles/{directoryRole%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDirectoryRoleItemRequestBuilder instantiates a new DirectoryRoleItemRequestBuilder and sets the default values.
func NewDirectoryRoleItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryRoleItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDirectoryRoleItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete entity from directoryRoles
func (m *DirectoryRoleItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *DirectoryRoleItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation retrieve the properties of a directoryRole object. The role must be activated in tenant for a successful response. You can use both the object ID and template ID of the **directoryRole** with this API. The template ID of a built-in role is immutable and can be seen in the role description on the Azure portal. For details, see Role template IDs.
func (m *DirectoryRoleItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DirectoryRoleItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update entity in directoryRoles
func (m *DirectoryRoleItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable, requestConfiguration *DirectoryRoleItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete entity from directoryRoles
func (m *DirectoryRoleItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DirectoryRoleItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the properties of a directoryRole object. The role must be activated in tenant for a successful response. You can use both the object ID and template ID of the **directoryRole** with this API. The template ID of a built-in role is immutable and can be seen in the role description on the Azure portal. For details, see Role template IDs.
func (m *DirectoryRoleItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DirectoryRoleItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryRoleFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable), nil
}
// GetMemberGroups provides operations to call the getMemberGroups method.
func (m *DirectoryRoleItemRequestBuilder) GetMemberGroups()(*id63def69f45b85c0f48fc2880b865fc141de2bfccf786d2d6a9813abb7b316e9.GetMemberGroupsRequestBuilder) {
    return id63def69f45b85c0f48fc2880b865fc141de2bfccf786d2d6a9813abb7b316e9.NewGetMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetMemberObjects provides operations to call the getMemberObjects method.
func (m *DirectoryRoleItemRequestBuilder) GetMemberObjects()(*id49876e5f70243d663a4f93ec6d77d7f89f0c91f1eda2daad48f4f037781d7b0.GetMemberObjectsRequestBuilder) {
    return id49876e5f70243d663a4f93ec6d77d7f89f0c91f1eda2daad48f4f037781d7b0.NewGetMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.directoryRole entity.
func (m *DirectoryRoleItemRequestBuilder) Members()(*if1385e9a44279e575f64436ffd1ff5c93bc83fbaa5a4051782729fbe1d3bdb3d.MembersRequestBuilder) {
    return if1385e9a44279e575f64436ffd1ff5c93bc83fbaa5a4051782729fbe1d3bdb3d.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.directoryRoles.item.members.item collection
func (m *DirectoryRoleItemRequestBuilder) MembersById(id string)(*i4dc00e738a43cdc3dbc464cba96748a3f629a6e7f36a76cbb15bf379872ebba9.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i4dc00e738a43cdc3dbc464cba96748a3f629a6e7f36a76cbb15bf379872ebba9.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update entity in directoryRoles
func (m *DirectoryRoleItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable, requestConfiguration *DirectoryRoleItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryRoleFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryRoleable), nil
}
// Restore provides operations to call the restore method.
func (m *DirectoryRoleItemRequestBuilder) Restore()(*ibcb647a3bce4367afc71872d9ca7c9cc64edcc4e60fe4831e8fcab3de590a294.RestoreRequestBuilder) {
    return ibcb647a3bce4367afc71872d9ca7c9cc64edcc4e60fe4831e8fcab3de590a294.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ScopedMembers provides operations to manage the scopedMembers property of the microsoft.graph.directoryRole entity.
func (m *DirectoryRoleItemRequestBuilder) ScopedMembers()(*i718fb7a72b38652faad082566b31b033df2c9f8d04cc0af99ee1e607f87e8d80.ScopedMembersRequestBuilder) {
    return i718fb7a72b38652faad082566b31b033df2c9f8d04cc0af99ee1e607f87e8d80.NewScopedMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ScopedMembersById provides operations to manage the scopedMembers property of the microsoft.graph.directoryRole entity.
func (m *DirectoryRoleItemRequestBuilder) ScopedMembersById(id string)(*ic73a5746a64610f345c78c4114780b56c25d27fdb8734dbd3b6500ebf1ad197c.ScopedRoleMembershipItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["scopedRoleMembership%2Did"] = id
    }
    return ic73a5746a64610f345c78c4114780b56c25d27fdb8734dbd3b6500ebf1ad197c.NewScopedRoleMembershipItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
