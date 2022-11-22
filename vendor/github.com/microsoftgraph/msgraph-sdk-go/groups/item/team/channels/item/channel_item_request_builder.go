package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i11526ee5699d623863f00b0b685d457eb13b19842baadbeddfd36eecf73dcdf5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/completemigration"
    i1370e74741c1f731fcda20e6e79b93c6cbee203716569746c1e43924df754b98 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/filesfolder"
    i5f0644cbd9043ef3ecfed8730c6a012f1f93d0bbc2a78c33eccf27e35ec1a2b7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    iac6867405139fcb2f6241d662e12bba0c0fbc752185e7378ef357512b9d5fd59 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/sharedwithteams"
    iaf3119de5714740cbdb24901fa13a2700a0dc0f152b343a2c356536580c6e10b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/tabs"
    ib9ab15dbda262e6b76749343df5d5f49e644a1926521ae149443532a84b18125 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/removeemail"
    icb337b1772f31d93683bb7dfb1edf3c4a97b38d5449f9ca882a52e68e02d5e27 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/members"
    iece274c950455f68d3b9284c5b46b89d9f9f8aa9c5482a9ba6d18dde70bd04f8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/messages"
    if559d86b89add997d564b55f90c9800ec81efa7e7c483735cee45eedd11460bd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/provisionemail"
    i045b00b932c574167b22217b96e68e85ba6f525a94bcd3b95118701b46c60f03 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/tabs/item"
    ib543babcb7b560b08aa30ad60d5db8d65b281cf0de36ac91c72b9eb5099e576b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/members/item"
    ic2fc4726b6ef1949b7220307cf2677355049193c9fd3815132e291979132f29b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/messages/item"
    ie56d1c765e4668273ba532044619456ce7199ab6bb9d66e68a89ff5aec1c0bad "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item/sharedwithteams/item"
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
func (m *ChannelItemRequestBuilder) CompleteMigration()(*i11526ee5699d623863f00b0b685d457eb13b19842baadbeddfd36eecf73dcdf5.CompleteMigrationRequestBuilder) {
    return i11526ee5699d623863f00b0b685d457eb13b19842baadbeddfd36eecf73dcdf5.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewChannelItemRequestBuilderInternal instantiates a new ChannelItemRequestBuilder and sets the default values.
func NewChannelItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChannelItemRequestBuilder) {
    m := &ChannelItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/team/channels/{channel%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property channels for groups
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
// CreatePatchRequestInformation update the navigation property channels in groups
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
// Delete delete navigation property channels for groups
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
func (m *ChannelItemRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*i5f0644cbd9043ef3ecfed8730c6a012f1f93d0bbc2a78c33eccf27e35ec1a2b7.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return i5f0644cbd9043ef3ecfed8730c6a012f1f93d0bbc2a78c33eccf27e35ec1a2b7.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) FilesFolder()(*i1370e74741c1f731fcda20e6e79b93c6cbee203716569746c1e43924df754b98.FilesFolderRequestBuilder) {
    return i1370e74741c1f731fcda20e6e79b93c6cbee203716569746c1e43924df754b98.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *ChannelItemRequestBuilder) Members()(*icb337b1772f31d93683bb7dfb1edf3c4a97b38d5449f9ca882a52e68e02d5e27.MembersRequestBuilder) {
    return icb337b1772f31d93683bb7dfb1edf3c4a97b38d5449f9ca882a52e68e02d5e27.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MembersById(id string)(*ib543babcb7b560b08aa30ad60d5db8d65b281cf0de36ac91c72b9eb5099e576b.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return ib543babcb7b560b08aa30ad60d5db8d65b281cf0de36ac91c72b9eb5099e576b.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Messages()(*iece274c950455f68d3b9284c5b46b89d9f9f8aa9c5482a9ba6d18dde70bd04f8.MessagesRequestBuilder) {
    return iece274c950455f68d3b9284c5b46b89d9f9f8aa9c5482a9ba6d18dde70bd04f8.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) MessagesById(id string)(*ic2fc4726b6ef1949b7220307cf2677355049193c9fd3815132e291979132f29b.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return ic2fc4726b6ef1949b7220307cf2677355049193c9fd3815132e291979132f29b.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property channels in groups
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
func (m *ChannelItemRequestBuilder) ProvisionEmail()(*if559d86b89add997d564b55f90c9800ec81efa7e7c483735cee45eedd11460bd.ProvisionEmailRequestBuilder) {
    return if559d86b89add997d564b55f90c9800ec81efa7e7c483735cee45eedd11460bd.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *ChannelItemRequestBuilder) RemoveEmail()(*ib9ab15dbda262e6b76749343df5d5f49e644a1926521ae149443532a84b18125.RemoveEmailRequestBuilder) {
    return ib9ab15dbda262e6b76749343df5d5f49e644a1926521ae149443532a84b18125.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeams()(*iac6867405139fcb2f6241d662e12bba0c0fbc752185e7378ef357512b9d5fd59.SharedWithTeamsRequestBuilder) {
    return iac6867405139fcb2f6241d662e12bba0c0fbc752185e7378ef357512b9d5fd59.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) SharedWithTeamsById(id string)(*ie56d1c765e4668273ba532044619456ce7199ab6bb9d66e68a89ff5aec1c0bad.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return ie56d1c765e4668273ba532044619456ce7199ab6bb9d66e68a89ff5aec1c0bad.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) Tabs()(*iaf3119de5714740cbdb24901fa13a2700a0dc0f152b343a2c356536580c6e10b.TabsRequestBuilder) {
    return iaf3119de5714740cbdb24901fa13a2700a0dc0f152b343a2c356536580c6e10b.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *ChannelItemRequestBuilder) TabsById(id string)(*i045b00b932c574167b22217b96e68e85ba6f525a94bcd3b95118701b46c60f03.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i045b00b932c574167b22217b96e68e85ba6f525a94bcd3b95118701b46c60f03.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
