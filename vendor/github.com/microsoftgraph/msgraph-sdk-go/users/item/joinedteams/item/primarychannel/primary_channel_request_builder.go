package primarychannel

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i088b238319b5250557b39864d8d3c74d555659e596e1673f7a94fbcde0b2ed1a "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/sharedwithteams"
    i1cf55ca0a825cd1cac192481d048117f2a39d3c38bf116a271f891df4ec37b4e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/messages"
    i1edd6a970c3b4ab42c23d1b6eb53d6e64e98e80a0842bbf3bb7cba83b5d464e4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/filesfolder"
    i83ea2bc60b68a36c45c77a2f1476931aed5417ffec59dd4772298a0544d7d58f "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/removeemail"
    i9b14eb39a7e5c1ffdc2c1d4786456b4d00217974730db20bea548f66f1476654 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    ib09a59de5c12db445e6f36aaba2d77a43c3728a772fd214f17547ae07e2614b4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/members"
    ib50348ec811d0ba0748f574b275ea09412ffa3079660cb0588275df9d54ccb92 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/tabs"
    ie7936a806322eabc26d4bac3d1049d6a6e65c4b59f4d8899adbeca160c339f01 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/completemigration"
    ifd249b5e7bc45bf262c727daac53c1d3c3d82a0be934c29c9acab3f62f5eb9b8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/provisionemail"
    i0ca9484c89452a4bc89c2adb6fa50859ce27953ce6eeff8b05452f93cd422d0b "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/tabs/item"
    i3b1baf61208ab9a11e7db26e1306bcb1b020ac68643127b2262519ea13d8cf2b "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/messages/item"
    idf42e0e625ec60873dcf4117181f59264c1d666a648124f9f973ec7af2e1ad2d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/sharedwithteams/item"
    if0a6e5cbcdc9ab20e7dd380367a323df63997762be3102745e3275cb41762f3e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel/members/item"
)

// PrimaryChannelRequestBuilder provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
type PrimaryChannelRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PrimaryChannelRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrimaryChannelRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// PrimaryChannelRequestBuilderGetQueryParameters get the default channel, **General**, of a team.
type PrimaryChannelRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PrimaryChannelRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrimaryChannelRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PrimaryChannelRequestBuilderGetQueryParameters
}
// PrimaryChannelRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PrimaryChannelRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *PrimaryChannelRequestBuilder) CompleteMigration()(*ie7936a806322eabc26d4bac3d1049d6a6e65c4b59f4d8899adbeca160c339f01.CompleteMigrationRequestBuilder) {
    return ie7936a806322eabc26d4bac3d1049d6a6e65c4b59f4d8899adbeca160c339f01.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewPrimaryChannelRequestBuilderInternal instantiates a new PrimaryChannelRequestBuilder and sets the default values.
func NewPrimaryChannelRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrimaryChannelRequestBuilder) {
    m := &PrimaryChannelRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}/primaryChannel{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPrimaryChannelRequestBuilder instantiates a new PrimaryChannelRequestBuilder and sets the default values.
func NewPrimaryChannelRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrimaryChannelRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPrimaryChannelRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property primaryChannel for users
func (m *PrimaryChannelRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *PrimaryChannelRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get the default channel, **General**, of a team.
func (m *PrimaryChannelRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PrimaryChannelRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property primaryChannel in users
func (m *PrimaryChannelRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, requestConfiguration *PrimaryChannelRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property primaryChannel for users
func (m *PrimaryChannelRequestBuilder) Delete(ctx context.Context, requestConfiguration *PrimaryChannelRequestBuilderDeleteRequestConfiguration)(error) {
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
func (m *PrimaryChannelRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*i9b14eb39a7e5c1ffdc2c1d4786456b4d00217974730db20bea548f66f1476654.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return i9b14eb39a7e5c1ffdc2c1d4786456b4d00217974730db20bea548f66f1476654.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) FilesFolder()(*i1edd6a970c3b4ab42c23d1b6eb53d6e64e98e80a0842bbf3bb7cba83b5d464e4.FilesFolderRequestBuilder) {
    return i1edd6a970c3b4ab42c23d1b6eb53d6e64e98e80a0842bbf3bb7cba83b5d464e4.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get the default channel, **General**, of a team.
func (m *PrimaryChannelRequestBuilder) Get(ctx context.Context, requestConfiguration *PrimaryChannelRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, error) {
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
func (m *PrimaryChannelRequestBuilder) Members()(*ib09a59de5c12db445e6f36aaba2d77a43c3728a772fd214f17547ae07e2614b4.MembersRequestBuilder) {
    return ib09a59de5c12db445e6f36aaba2d77a43c3728a772fd214f17547ae07e2614b4.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MembersById(id string)(*if0a6e5cbcdc9ab20e7dd380367a323df63997762be3102745e3275cb41762f3e.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return if0a6e5cbcdc9ab20e7dd380367a323df63997762be3102745e3275cb41762f3e.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Messages()(*i1cf55ca0a825cd1cac192481d048117f2a39d3c38bf116a271f891df4ec37b4e.MessagesRequestBuilder) {
    return i1cf55ca0a825cd1cac192481d048117f2a39d3c38bf116a271f891df4ec37b4e.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MessagesById(id string)(*i3b1baf61208ab9a11e7db26e1306bcb1b020ac68643127b2262519ea13d8cf2b.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return i3b1baf61208ab9a11e7db26e1306bcb1b020ac68643127b2262519ea13d8cf2b.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property primaryChannel in users
func (m *PrimaryChannelRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, requestConfiguration *PrimaryChannelRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Channelable, error) {
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
func (m *PrimaryChannelRequestBuilder) ProvisionEmail()(*ifd249b5e7bc45bf262c727daac53c1d3c3d82a0be934c29c9acab3f62f5eb9b8.ProvisionEmailRequestBuilder) {
    return ifd249b5e7bc45bf262c727daac53c1d3c3d82a0be934c29c9acab3f62f5eb9b8.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *PrimaryChannelRequestBuilder) RemoveEmail()(*i83ea2bc60b68a36c45c77a2f1476931aed5417ffec59dd4772298a0544d7d58f.RemoveEmailRequestBuilder) {
    return i83ea2bc60b68a36c45c77a2f1476931aed5417ffec59dd4772298a0544d7d58f.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeams()(*i088b238319b5250557b39864d8d3c74d555659e596e1673f7a94fbcde0b2ed1a.SharedWithTeamsRequestBuilder) {
    return i088b238319b5250557b39864d8d3c74d555659e596e1673f7a94fbcde0b2ed1a.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeamsById(id string)(*idf42e0e625ec60873dcf4117181f59264c1d666a648124f9f973ec7af2e1ad2d.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return idf42e0e625ec60873dcf4117181f59264c1d666a648124f9f973ec7af2e1ad2d.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Tabs()(*ib50348ec811d0ba0748f574b275ea09412ffa3079660cb0588275df9d54ccb92.TabsRequestBuilder) {
    return ib50348ec811d0ba0748f574b275ea09412ffa3079660cb0588275df9d54ccb92.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) TabsById(id string)(*i0ca9484c89452a4bc89c2adb6fa50859ce27953ce6eeff8b05452f93cd422d0b.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i0ca9484c89452a4bc89c2adb6fa50859ce27953ce6eeff8b05452f93cd422d0b.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
