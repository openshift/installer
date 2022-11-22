package primarychannel

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1e4013461c5123188ab7ccd33dd297c0d5a7865ba9fea8e6c138306141eeec4a "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/completemigration"
    i26a22db42badb207f51d6bdd9193eccfdc7cad972de53c3620609a905e198f7a "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/removeemail"
    i4785692a9d61d4a8638d06fe97e1daa4feb18f7b743e7ae34ff3f07f3d9425b4 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/messages"
    i6093a5b7b609cfa10a7aee9a44ffd2db461382a56d5f3817bbb86584ec07b126 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/sharedwithteams"
    i76bbd884af07122b67f18d59d9b49c81200bf91786320bfb62dc2c5cce7d254e "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    i977f2dd26da5a2e1880df87bd5372aa7a19a14868a426678151d974f360c2a01 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/members"
    i9cea999ed45be284a8003fe37addd7bac8685c5b60e557f3adc8dc8bcf2d27b5 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/provisionemail"
    ib05faad6285a5d11c891d13f0450d43dfbe0ddb57be2fa94177f219559ef0527 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/tabs"
    id169406704fae3a188bcd329b3cbae975ed48842ae13c31a260f221158f8516a "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/filesfolder"
    i0193f7b86687af396beaaecb295b77e66ded553cf49127c0d06bc79b38a4d93f "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/tabs/item"
    ib00fe6f5242d43b4e0cf713930bcfb1828c83941d27088c19568b60fa6e113d8 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/messages/item"
    ie734300e5bedd84acd036f8be309a4094f66ee9a58c5ec5b270079e8efc64674 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/members/item"
    if0ea137f3825c0b8efcf099d1486febfa9d0dd12f02b49b849bd50671902b8d8 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel/sharedwithteams/item"
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
func (m *PrimaryChannelRequestBuilder) CompleteMigration()(*i1e4013461c5123188ab7ccd33dd297c0d5a7865ba9fea8e6c138306141eeec4a.CompleteMigrationRequestBuilder) {
    return i1e4013461c5123188ab7ccd33dd297c0d5a7865ba9fea8e6c138306141eeec4a.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewPrimaryChannelRequestBuilderInternal instantiates a new PrimaryChannelRequestBuilder and sets the default values.
func NewPrimaryChannelRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrimaryChannelRequestBuilder) {
    m := &PrimaryChannelRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/teams/{team%2Did}/primaryChannel{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property primaryChannel for teams
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
// CreatePatchRequestInformation update the navigation property primaryChannel in teams
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
// Delete delete navigation property primaryChannel for teams
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
func (m *PrimaryChannelRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*i76bbd884af07122b67f18d59d9b49c81200bf91786320bfb62dc2c5cce7d254e.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return i76bbd884af07122b67f18d59d9b49c81200bf91786320bfb62dc2c5cce7d254e.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) FilesFolder()(*id169406704fae3a188bcd329b3cbae975ed48842ae13c31a260f221158f8516a.FilesFolderRequestBuilder) {
    return id169406704fae3a188bcd329b3cbae975ed48842ae13c31a260f221158f8516a.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *PrimaryChannelRequestBuilder) Members()(*i977f2dd26da5a2e1880df87bd5372aa7a19a14868a426678151d974f360c2a01.MembersRequestBuilder) {
    return i977f2dd26da5a2e1880df87bd5372aa7a19a14868a426678151d974f360c2a01.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MembersById(id string)(*ie734300e5bedd84acd036f8be309a4094f66ee9a58c5ec5b270079e8efc64674.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return ie734300e5bedd84acd036f8be309a4094f66ee9a58c5ec5b270079e8efc64674.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Messages()(*i4785692a9d61d4a8638d06fe97e1daa4feb18f7b743e7ae34ff3f07f3d9425b4.MessagesRequestBuilder) {
    return i4785692a9d61d4a8638d06fe97e1daa4feb18f7b743e7ae34ff3f07f3d9425b4.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MessagesById(id string)(*ib00fe6f5242d43b4e0cf713930bcfb1828c83941d27088c19568b60fa6e113d8.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return ib00fe6f5242d43b4e0cf713930bcfb1828c83941d27088c19568b60fa6e113d8.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property primaryChannel in teams
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
func (m *PrimaryChannelRequestBuilder) ProvisionEmail()(*i9cea999ed45be284a8003fe37addd7bac8685c5b60e557f3adc8dc8bcf2d27b5.ProvisionEmailRequestBuilder) {
    return i9cea999ed45be284a8003fe37addd7bac8685c5b60e557f3adc8dc8bcf2d27b5.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *PrimaryChannelRequestBuilder) RemoveEmail()(*i26a22db42badb207f51d6bdd9193eccfdc7cad972de53c3620609a905e198f7a.RemoveEmailRequestBuilder) {
    return i26a22db42badb207f51d6bdd9193eccfdc7cad972de53c3620609a905e198f7a.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeams()(*i6093a5b7b609cfa10a7aee9a44ffd2db461382a56d5f3817bbb86584ec07b126.SharedWithTeamsRequestBuilder) {
    return i6093a5b7b609cfa10a7aee9a44ffd2db461382a56d5f3817bbb86584ec07b126.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeamsById(id string)(*if0ea137f3825c0b8efcf099d1486febfa9d0dd12f02b49b849bd50671902b8d8.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return if0ea137f3825c0b8efcf099d1486febfa9d0dd12f02b49b849bd50671902b8d8.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Tabs()(*ib05faad6285a5d11c891d13f0450d43dfbe0ddb57be2fa94177f219559ef0527.TabsRequestBuilder) {
    return ib05faad6285a5d11c891d13f0450d43dfbe0ddb57be2fa94177f219559ef0527.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) TabsById(id string)(*i0193f7b86687af396beaaecb295b77e66ded553cf49127c0d06bc79b38a4d93f.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i0193f7b86687af396beaaecb295b77e66ded553cf49127c0d06bc79b38a4d93f.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
