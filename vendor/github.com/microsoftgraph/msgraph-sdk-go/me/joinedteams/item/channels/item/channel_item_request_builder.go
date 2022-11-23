package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1135d2333f45686e0e95fc542221d60e6b9c4e28f8a0a333bab31e3d74638da0 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    i253a46a809c827c7298a764f8c82f551eb9c16db15abb743fffe2eba96700378 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/filesfolder"
    i2ad3e3c936b6898b388b6d3cc1c2364d7851f04fd4217471b556b15e62ffff65 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/provisionemail"
    i40fe54e22de0311e21811c5781924675aa201c68840b6af48bd6b80c4ccfc083 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/messages"
    i8162c3ce6781bd606b762c0a0a2d4db93f251d69b22fccbc03e039e4c8dc2c76 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/sharedwithteams"
    iaf082408b1fd7fd4e8e3c6c1befccb1388486f663eb9ad351a7f8f61d4cf7be6 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/completemigration"
    ie326949afcb08e6a5436a367b925b7f68ec414a2c72fce4b6164a4b239264d52 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/tabs"
    ief0cacbd8be774a1502106b209117815bb3293f805cb67532513c01f04e36513 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/members"
    iefaaa288cea77216f47ed2ce48474ac53a86f88b057403e15cbe3142e3187502 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/removeemail"
    i2da7ef6ccb4db37d9d9fe954811ee5ddad198aa369fa874d4f6099c7f62e0bb6 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/messages/item"
    i3a81824602dda13e4a5b75e20ede850d39a8762eb95c95f58b595895ef98be06 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/tabs/item"
    iacd3c3e82b4dcdba0a068ae5f345aa156aee3f7cf8b623e086211d0dc51afdef "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/sharedwithteams/item"
    if78b603399676ac0a458d83dbd0e2eb877a642ddfd394dd4905b49017bf3160d "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item/members/item"
)

// ChannelItemRequestBuilder provides operations to manage the channels property of the microsoft.graph.team entity.
type ChannelItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ChannelItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChannelItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ChannelItemRequestBuilderGetQueryParameters the collection of channels and messages associated with the team.
type ChannelItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ChannelItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChannelItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ChannelItemRequestBuilderGetQueryParameters
}
// ChannelItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChannelItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *ChannelItemRequestBuilder) CompleteMigration()(*iaf082408b1fd7fd4e8e3c6c1befccb1388486f663eb9ad351a7f8f61d4cf7be6.CompleteMigrationRequestBuilder) {
    return iaf082408b1fd7fd4e8e3c6c1befccb1388486f663eb9ad351a7f8f61d4cf7be6.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewChannelItemRequestBuilderInternal instantiates a new ChannelItemRequestBuilder and sets the default values.
func NewChannelItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChannelItemRequestBuilder) {
    m := &ChannelItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/joinedTeams/{team%2Did}/channels/{channel%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewChannelItemRequestBuilder instantiates a new ChannelItemRequestBuilder and sets the default values.
func NewChannelItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChannelItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewChannelItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property channels for me
func (m *ChannelItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ChannelItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the collection of channels and messages associated with the team.
func (m *ChannelItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ChannelItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property channels in me
func (m *ChannelItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, requestConfiguration *ChannelItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property channels for me
func (m *ChannelItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ChannelItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName provides operations to call the doesUserHaveAccess method.
func (m *ChannelItemRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*i1135d2333f45686e0e95fc542221d60e6b9c4e28f8a0a333bab31e3d74638da0.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return i1135d2333f45686e0e95fc542221d60e6b9c4e28f8a0a333bab31e3d74638da0.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) FilesFolder()(*i253a46a809c827c7298a764f8c82f551eb9c16db15abb743fffe2eba96700378.FilesFolderRequestBuilder) {
    return i253a46a809c827c7298a764f8c82f551eb9c16db15abb743fffe2eba96700378.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the collection of channels and messages associated with the team.
func (m *ChannelItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ChannelItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChannelFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable), nil
}
// Members provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Members()(*ief0cacbd8be774a1502106b209117815bb3293f805cb67532513c01f04e36513.MembersRequestBuilder) {
    return ief0cacbd8be774a1502106b209117815bb3293f805cb67532513c01f04e36513.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MembersById(id string)(*if78b603399676ac0a458d83dbd0e2eb877a642ddfd394dd4905b49017bf3160d.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return if78b603399676ac0a458d83dbd0e2eb877a642ddfd394dd4905b49017bf3160d.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Messages()(*i40fe54e22de0311e21811c5781924675aa201c68840b6af48bd6b80c4ccfc083.MessagesRequestBuilder) {
    return i40fe54e22de0311e21811c5781924675aa201c68840b6af48bd6b80c4ccfc083.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MessagesById(id string)(*i2da7ef6ccb4db37d9d9fe954811ee5ddad198aa369fa874d4f6099c7f62e0bb6.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return i2da7ef6ccb4db37d9d9fe954811ee5ddad198aa369fa874d4f6099c7f62e0bb6.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property channels in me
func (m *ChannelItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, requestConfiguration *ChannelItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChannelFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable), nil
}
// ProvisionEmail provides operations to call the provisionEmail method.
func (m *ChannelItemRequestBuilder) ProvisionEmail()(*i2ad3e3c936b6898b388b6d3cc1c2364d7851f04fd4217471b556b15e62ffff65.ProvisionEmailRequestBuilder) {
    return i2ad3e3c936b6898b388b6d3cc1c2364d7851f04fd4217471b556b15e62ffff65.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *ChannelItemRequestBuilder) RemoveEmail()(*iefaaa288cea77216f47ed2ce48474ac53a86f88b057403e15cbe3142e3187502.RemoveEmailRequestBuilder) {
    return iefaaa288cea77216f47ed2ce48474ac53a86f88b057403e15cbe3142e3187502.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeams()(*i8162c3ce6781bd606b762c0a0a2d4db93f251d69b22fccbc03e039e4c8dc2c76.SharedWithTeamsRequestBuilder) {
    return i8162c3ce6781bd606b762c0a0a2d4db93f251d69b22fccbc03e039e4c8dc2c76.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeamsById(id string)(*iacd3c3e82b4dcdba0a068ae5f345aa156aee3f7cf8b623e086211d0dc51afdef.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return iacd3c3e82b4dcdba0a068ae5f345aa156aee3f7cf8b623e086211d0dc51afdef.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Tabs()(*ie326949afcb08e6a5436a367b925b7f68ec414a2c72fce4b6164a4b239264d52.TabsRequestBuilder) {
    return ie326949afcb08e6a5436a367b925b7f68ec414a2c72fce4b6164a4b239264d52.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) TabsById(id string)(*i3a81824602dda13e4a5b75e20ede850d39a8762eb95c95f58b595895ef98be06.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i3a81824602dda13e4a5b75e20ede850d39a8762eb95c95f58b595895ef98be06.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
