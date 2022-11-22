package primarychannel

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i23c2230a9d26cc65f1377304e7e2458d1952291a8ec639fd039abd93267cf215 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/messages"
    i3772d1545a94e182a501f85827b513f4692061bccbe8b23eb73c53a831f72d14 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/removeemail"
    i42121f7ad8826a000f65df3b00059d6fc81d8554c4cbbff019865a887d20a56a "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/provisionemail"
    i57e14fb40bbf624605d5248ef685ec9b4257446262e6f8bd78b5c13c71e1e3df "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/tabs"
    i75cf55986f617df38b649b4ebcf62bc92728e8a45cbbb6b95f7040252e76e32c "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/completemigration"
    i7674f91698779c7d9f326e96636581b5648ecc350b417dfca4417e91d69fa246 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/sharedwithteams"
    i7dcb24e562a6855f9b4081d07d2b6077073bb1827c1841433eb6644b41bca8f2 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/filesfolder"
    i8cdf95eab0d1df8385b6f1bdd3b098788d83540da047f10671bda06132b1adb3 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/members"
    ic01dbb97d6e3853c23435148feebfecb4638fffb80dcdcccee1e88a8a42f9470 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/doesuserhaveaccessuseriduseridtenantidtenantiduserprincipalnameuserprincipalname"
    i11b5e61a109dfce2102d2d81a8a621a28b579c197305053d9f61fd8e3a46a788 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/members/item"
    i7dfaad961bf25e32c3ce15ea57029e5e024ff862e31e23a16254e10f24113084 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/messages/item"
    ibeddd7f15f21eea6c1989a2c2ecf37d69b81be13bb77e9a96fd8bb8e0a63c3f7 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/tabs/item"
    id353d0ea6d61580fb6f14bd570e7f3f563fcdca43d360caa85303cff4a184de8 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel/sharedwithteams/item"
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
func (m *PrimaryChannelRequestBuilder) CompleteMigration()(*i75cf55986f617df38b649b4ebcf62bc92728e8a45cbbb6b95f7040252e76e32c.CompleteMigrationRequestBuilder) {
    return i75cf55986f617df38b649b4ebcf62bc92728e8a45cbbb6b95f7040252e76e32c.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewPrimaryChannelRequestBuilderInternal instantiates a new PrimaryChannelRequestBuilder and sets the default values.
func NewPrimaryChannelRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PrimaryChannelRequestBuilder) {
    m := &PrimaryChannelRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/joinedTeams/{team%2Did}/primaryChannel{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property primaryChannel for me
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
// CreatePatchRequestInformation update the navigation property primaryChannel in me
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
// Delete delete navigation property primaryChannel for me
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
func (m *PrimaryChannelRequestBuilder) DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalName()(*ic01dbb97d6e3853c23435148feebfecb4638fffb80dcdcccee1e88a8a42f9470.DoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilder) {
    return ic01dbb97d6e3853c23435148feebfecb4638fffb80dcdcccee1e88a8a42f9470.NewDoesUserHaveAccessuserIdUserIdTenantIdTenantIdUserPrincipalNameUserPrincipalNameRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FilesFolder provides operations to manage the filesFolder property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) FilesFolder()(*i7dcb24e562a6855f9b4081d07d2b6077073bb1827c1841433eb6644b41bca8f2.FilesFolderRequestBuilder) {
    return i7dcb24e562a6855f9b4081d07d2b6077073bb1827c1841433eb6644b41bca8f2.NewFilesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *PrimaryChannelRequestBuilder) Members()(*i8cdf95eab0d1df8385b6f1bdd3b098788d83540da047f10671bda06132b1adb3.MembersRequestBuilder) {
    return i8cdf95eab0d1df8385b6f1bdd3b098788d83540da047f10671bda06132b1adb3.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MembersById(id string)(*i11b5e61a109dfce2102d2d81a8a621a28b579c197305053d9f61fd8e3a46a788.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i11b5e61a109dfce2102d2d81a8a621a28b579c197305053d9f61fd8e3a46a788.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Messages()(*i23c2230a9d26cc65f1377304e7e2458d1952291a8ec639fd039abd93267cf215.MessagesRequestBuilder) {
    return i23c2230a9d26cc65f1377304e7e2458d1952291a8ec639fd039abd93267cf215.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) MessagesById(id string)(*i7dfaad961bf25e32c3ce15ea57029e5e024ff862e31e23a16254e10f24113084.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return i7dfaad961bf25e32c3ce15ea57029e5e024ff862e31e23a16254e10f24113084.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property primaryChannel in me
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
func (m *PrimaryChannelRequestBuilder) ProvisionEmail()(*i42121f7ad8826a000f65df3b00059d6fc81d8554c4cbbff019865a887d20a56a.ProvisionEmailRequestBuilder) {
    return i42121f7ad8826a000f65df3b00059d6fc81d8554c4cbbff019865a887d20a56a.NewProvisionEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RemoveEmail provides operations to call the removeEmail method.
func (m *PrimaryChannelRequestBuilder) RemoveEmail()(*i3772d1545a94e182a501f85827b513f4692061bccbe8b23eb73c53a831f72d14.RemoveEmailRequestBuilder) {
    return i3772d1545a94e182a501f85827b513f4692061bccbe8b23eb73c53a831f72d14.NewRemoveEmailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeams provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeams()(*i7674f91698779c7d9f326e96636581b5648ecc350b417dfca4417e91d69fa246.SharedWithTeamsRequestBuilder) {
    return i7674f91698779c7d9f326e96636581b5648ecc350b417dfca4417e91d69fa246.NewSharedWithTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharedWithTeamsById provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) SharedWithTeamsById(id string)(*id353d0ea6d61580fb6f14bd570e7f3f563fcdca43d360caa85303cff4a184de8.SharedWithChannelTeamInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedWithChannelTeamInfo%2Did"] = id
    }
    return id353d0ea6d61580fb6f14bd570e7f3f563fcdca43d360caa85303cff4a184de8.NewSharedWithChannelTeamInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) Tabs()(*i57e14fb40bbf624605d5248ef685ec9b4257446262e6f8bd78b5c13c71e1e3df.TabsRequestBuilder) {
    return i57e14fb40bbf624605d5248ef685ec9b4257446262e6f8bd78b5c13c71e1e3df.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.channel entity.
func (m *PrimaryChannelRequestBuilder) TabsById(id string)(*ibeddd7f15f21eea6c1989a2c2ecf37d69b81be13bb77e9a96fd8bb8e0a63c3f7.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return ibeddd7f15f21eea6c1989a2c2ecf37d69b81be13bb77e9a96fd8bb8e0a63c3f7.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
