package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1bf8c2f95e6e929ffe140e396a108a74fd704da511c9355ecf989be5baee6ae7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    i52508fdeacdff04c86ef6ed98d578ffd2c77d62a7edfb160198322b9a67e8e0a "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/tabs"
    i793c23296f16d9688a86912d15fe2bf5d6212cd55a992ac7d34ce9ea17c9d8cb "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/provisionemail"
    i99487f33782211bac82941f8c49dfba40147399ed21891f1ad74a96382992ae3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/completemigration"
    ia0a1c012adec3a5f091e102240b8c385a20d774609da05ca9b6a1fa6f9c9d76d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/sharedwithteams"
    ib0c11d6a85aba358cbcfc13c9be3902096670bb120397dd26d3faa62652e2524 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/filesfolder"
    ic775254656b39a9ea58d56de7e743ecea61bb0555fd277cbb8e893134cdc519e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/members"
    ie5ca968c3c0027bf6b121a128bc8b2a3b8c58e61b6f1a6ccaa044003449557cf "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/removeemail"
    if0b23a5506431d362cbdca425249e3b64c69fb21596d086069afe81af11c676d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/messages"
    i51bfbcbb97372642de8f2ed3fe79c93955eafb6f6ba3bfe2532d49160430b57d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/members/item"
    i7aaac21246fdaa52fcd18af794bbc834c1e127facfe4e64e697ea1931cbe2496 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/tabs/item"
    ia6c1a592dc0444820daee3899e4be3c57cdf1079f6c5d508460448df15f9926d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/messages/item"
    icbecf33726874f634ca1f90f1559300d9f98b33b3db2ad035c454fb48a1bd846 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item/sharedwithteams/item"
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
func (m *ChannelItemRequestBuilder) CompleteMigration()(*i99487f33782211bac82941f8c49dfba40147399ed21891f1ad74a96382992ae3.CompleteMigrationRequestBuilder) {
    return i99487f33782211bac82941f8c49dfba40147399ed21891f1ad74a96382992ae3.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewChannelItemRequestBuilderInternal instantiates a new ChannelItemRequestBuilder and sets the default values.
func NewChannelItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChannelItemRequestBuilder) {
    m := &ChannelItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}/channels/{channel%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property channels for users
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
// CreatePatchRequestInformation update the navigation property channels in users
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
// Delete delete navigation property channels for users
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
func (m *ChannelItemRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*i1bf8c2f95e6e929ffe140e396a108a74fd704da511c9355ecf989be5baee6ae7.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return i1bf8c2f95e6e929ffe140e396a108a74fd704da511c9355ecf989be5baee6ae7.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) FilesFolder()(*ib0c11d6a85aba358cbcfc13c9be3902096670bb120397dd26d3faa62652e2524.FilesFolderRequestBuilder) {
    return ib0c11d6a85aba358cbcfc13c9be3902096670bb120397dd26d3faa62652e2524.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *ChannelItemRequestBuilder) Members()(*ic775254656b39a9ea58d56de7e743ecea61bb0555fd277cbb8e893134cdc519e.MembersRequestBuilder) {
    return ic775254656b39a9ea58d56de7e743ecea61bb0555fd277cbb8e893134cdc519e.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MembersById(id string)(*i51bfbcbb97372642de8f2ed3fe79c93955eafb6f6ba3bfe2532d49160430b57d.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i51bfbcbb97372642de8f2ed3fe79c93955eafb6f6ba3bfe2532d49160430b57d.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Messages()(*if0b23a5506431d362cbdca425249e3b64c69fb21596d086069afe81af11c676d.MessagesRequestBuilder) {
    return if0b23a5506431d362cbdca425249e3b64c69fb21596d086069afe81af11c676d.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MessagesById(id string)(*ia6c1a592dc0444820daee3899e4be3c57cdf1079f6c5d508460448df15f9926d.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return ia6c1a592dc0444820daee3899e4be3c57cdf1079f6c5d508460448df15f9926d.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property channels in users
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
func (m *ChannelItemRequestBuilder) ProvisionEmail()(*i793c23296f16d9688a86912d15fe2bf5d6212cd55a992ac7d34ce9ea17c9d8cb.ProvisionEmailRequestBuilder) {
    return i793c23296f16d9688a86912d15fe2bf5d6212cd55a992ac7d34ce9ea17c9d8cb.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *ChannelItemRequestBuilder) RemoveEmail()(*ie5ca968c3c0027bf6b121a128bc8b2a3b8c58e61b6f1a6ccaa044003449557cf.RemoveEmailRequestBuilder) {
    return ie5ca968c3c0027bf6b121a128bc8b2a3b8c58e61b6f1a6ccaa044003449557cf.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeams()(*ia0a1c012adec3a5f091e102240b8c385a20d774609da05ca9b6a1fa6f9c9d76d.SharedWithTeamsRequestBuilder) {
    return ia0a1c012adec3a5f091e102240b8c385a20d774609da05ca9b6a1fa6f9c9d76d.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeamsById(id string)(*icbecf33726874f634ca1f90f1559300d9f98b33b3db2ad035c454fb48a1bd846.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return icbecf33726874f634ca1f90f1559300d9f98b33b3db2ad035c454fb48a1bd846.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Tabs()(*i52508fdeacdff04c86ef6ed98d578ffd2c77d62a7edfb160198322b9a67e8e0a.TabsRequestBuilder) {
    return i52508fdeacdff04c86ef6ed98d578ffd2c77d62a7edfb160198322b9a67e8e0a.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) TabsById(id string)(*i7aaac21246fdaa52fcd18af794bbc834c1e127facfe4e64e697ea1931cbe2496.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i7aaac21246fdaa52fcd18af794bbc834c1e127facfe4e64e697ea1931cbe2496.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
