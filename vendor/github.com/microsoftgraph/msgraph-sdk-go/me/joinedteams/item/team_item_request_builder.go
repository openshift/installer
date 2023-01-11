package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i04a148e32be31a86cd21b897c8b55a4508d63dbae47a02eaeb122662dbd2ff9b "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/template"
    i1e6cfda39bf71ff6c49bf501ccb56055e496f12dd4a459a4c9f832017b3ed495 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/incomingchannels"
    i5cc953f63726e53531e5d00c609bd14469e1ce23944c22b6aea0986c18fb0043 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/primarychannel"
    i6bdaa0e5d28581a8a1e2d235e3e781f04d9c598ddc6abab4f4f4482948b402c4 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/clone"
    i7a26ea73a074ce439c29db16635d2ace903fc13c34496531fb310f143a72ad46 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/tags"
    i85d270c8e10c233ea48f496d72490ff86a9da2a715ff50b98d9bb49990543826 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/completemigration"
    i8a1cdbeac728d5d9d3409d0d7085c53384ad37435e0292d966ed94bbc4155a05 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/group"
    i9fa0e9d329dc2b42ce0cc0330991bb8f8e864efaaef5061789d895e28321a6b2 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/operations"
    iaaf6381de60bdbb052122028b1440af6c015d4e8c24cf971f1c3968b5a4f3c6b "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/archive"
    ic08a09e622b3f3279dad3fb1dc0d9adf50886e9a08f52205e44e27965d46190a "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/allchannels"
    ica5ec7aca95820534d087da722646ac0c0eb14565fc4c028141ff7080240b58d "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/photo"
    icf925d6e8373dd15bb408b246595c1d2598e1881b555a62a02a76300fc5b7cd2 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/installedapps"
    id463d65124ba412b3980ec713bebe8eb90e4a925515f5f993cd79d5b01b70907 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels"
    ie8c9d4ba2aa39351cc5580d10d45b8717b93a5d75b46c5be678a969e9175bcc6 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/sendactivitynotification"
    if5e128abdd6cb9303edd6d603805cc08bcb3a30eb8ff5fabd66714b2846349ee "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/unarchive"
    ifc923048eff969ef232f17cdaf6c11e18676c5c9e2918bc19001d30cbdc4b5c1 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/members"
    ifd53534f50d40567e607c2213e794582e29aa46a0c07e2d406db231a42a0140a "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/schedule"
    i03b2a35397701f29dd92bcf0dd549e252ab75a1bd539b1b1ca11a58738ba2c51 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/members/item"
    i0c17287a3757a64da956e62a7d60ea201a306f2234f98185836c4750f17a3b18 "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/allchannels/item"
    i289f2fd32e9a16e441f2a6b0e403b7f5f23d0c30924e6e4401734b1469fd06ce "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/tags/item"
    i77d75bcb3611254a7012d361f29a795e840693affe1de5861c58d310feff299c "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/operations/item"
    i986134cfb00ff93c88ffb51965b534e9da908849cf680fc5e4ff907a20066f7d "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/incomingchannels/item"
    id12792e5e45c250abd3ecec9bbc023cd6dc63edea08efc7907dcad0da46a8a0d "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/channels/item"
    ifaff06990f5e17731e64f75f4cb575f1ada76645d84341378dccf3aef4f66fab "github.com/microsoftgraph/msgraph-sdk-go/me/joinedteams/item/installedapps/item"
)

// TeamItemRequestBuilder provides operations to manage the joinedTeams property of the microsoft.graph.user entity.
type TeamItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TeamItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TeamItemRequestBuilderGetQueryParameters get joinedTeams from me
type TeamItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TeamItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TeamItemRequestBuilderGetQueryParameters
}
// TeamItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllChannels provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) AllChannels()(*ic08a09e622b3f3279dad3fb1dc0d9adf50886e9a08f52205e44e27965d46190a.AllChannelsRequestBuilder) {
    return ic08a09e622b3f3279dad3fb1dc0d9adf50886e9a08f52205e44e27965d46190a.NewAllChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllChannelsById provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) AllChannelsById(id string)(*i0c17287a3757a64da956e62a7d60ea201a306f2234f98185836c4750f17a3b18.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i0c17287a3757a64da956e62a7d60ea201a306f2234f98185836c4750f17a3b18.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Archive provides operations to call the archive method.
func (m *TeamItemRequestBuilder) Archive()(*iaaf6381de60bdbb052122028b1440af6c015d4e8c24cf971f1c3968b5a4f3c6b.ArchiveRequestBuilder) {
    return iaaf6381de60bdbb052122028b1440af6c015d4e8c24cf971f1c3968b5a4f3c6b.NewArchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Channels()(*id463d65124ba412b3980ec713bebe8eb90e4a925515f5f993cd79d5b01b70907.ChannelsRequestBuilder) {
    return id463d65124ba412b3980ec713bebe8eb90e4a925515f5f993cd79d5b01b70907.NewChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChannelsById provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) ChannelsById(id string)(*id12792e5e45c250abd3ecec9bbc023cd6dc63edea08efc7907dcad0da46a8a0d.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return id12792e5e45c250abd3ecec9bbc023cd6dc63edea08efc7907dcad0da46a8a0d.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Clone provides operations to call the clone method.
func (m *TeamItemRequestBuilder) Clone()(*i6bdaa0e5d28581a8a1e2d235e3e781f04d9c598ddc6abab4f4f4482948b402c4.CloneRequestBuilder) {
    return i6bdaa0e5d28581a8a1e2d235e3e781f04d9c598ddc6abab4f4f4482948b402c4.NewCloneRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *TeamItemRequestBuilder) CompleteMigration()(*i85d270c8e10c233ea48f496d72490ff86a9da2a715ff50b98d9bb49990543826.CompleteMigrationRequestBuilder) {
    return i85d270c8e10c233ea48f496d72490ff86a9da2a715ff50b98d9bb49990543826.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTeamItemRequestBuilderInternal instantiates a new TeamItemRequestBuilder and sets the default values.
func NewTeamItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamItemRequestBuilder) {
    m := &TeamItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/joinedTeams/{team%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTeamItemRequestBuilder instantiates a new TeamItemRequestBuilder and sets the default values.
func NewTeamItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTeamItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property joinedTeams for me
func (m *TeamItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *TeamItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get joinedTeams from me
func (m *TeamItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TeamItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property joinedTeams in me
func (m *TeamItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property joinedTeams for me
func (m *TeamItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TeamItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get joinedTeams from me
func (m *TeamItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TeamItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// Group provides operations to manage the group property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Group()(*i8a1cdbeac728d5d9d3409d0d7085c53384ad37435e0292d966ed94bbc4155a05.GroupRequestBuilder) {
    return i8a1cdbeac728d5d9d3409d0d7085c53384ad37435e0292d966ed94bbc4155a05.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannels()(*i1e6cfda39bf71ff6c49bf501ccb56055e496f12dd4a459a4c9f832017b3ed495.IncomingChannelsRequestBuilder) {
    return i1e6cfda39bf71ff6c49bf501ccb56055e496f12dd4a459a4c9f832017b3ed495.NewIncomingChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannelsById provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannelsById(id string)(*i986134cfb00ff93c88ffb51965b534e9da908849cf680fc5e4ff907a20066f7d.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i986134cfb00ff93c88ffb51965b534e9da908849cf680fc5e4ff907a20066f7d.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledApps()(*icf925d6e8373dd15bb408b246595c1d2598e1881b555a62a02a76300fc5b7cd2.InstalledAppsRequestBuilder) {
    return icf925d6e8373dd15bb408b246595c1d2598e1881b555a62a02a76300fc5b7cd2.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledAppsById(id string)(*ifaff06990f5e17731e64f75f4cb575f1ada76645d84341378dccf3aef4f66fab.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return ifaff06990f5e17731e64f75f4cb575f1ada76645d84341378dccf3aef4f66fab.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Members()(*ifc923048eff969ef232f17cdaf6c11e18676c5c9e2918bc19001d30cbdc4b5c1.MembersRequestBuilder) {
    return ifc923048eff969ef232f17cdaf6c11e18676c5c9e2918bc19001d30cbdc4b5c1.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) MembersById(id string)(*i03b2a35397701f29dd92bcf0dd549e252ab75a1bd539b1b1ca11a58738ba2c51.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i03b2a35397701f29dd92bcf0dd549e252ab75a1bd539b1b1ca11a58738ba2c51.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Operations()(*i9fa0e9d329dc2b42ce0cc0330991bb8f8e864efaaef5061789d895e28321a6b2.OperationsRequestBuilder) {
    return i9fa0e9d329dc2b42ce0cc0330991bb8f8e864efaaef5061789d895e28321a6b2.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) OperationsById(id string)(*i77d75bcb3611254a7012d361f29a795e840693affe1de5861c58d310feff299c.TeamsAsyncOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAsyncOperation%2Did"] = id
    }
    return i77d75bcb3611254a7012d361f29a795e840693affe1de5861c58d310feff299c.NewTeamsAsyncOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property joinedTeams in me
func (m *TeamItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// Photo provides operations to manage the photo property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Photo()(*ica5ec7aca95820534d087da722646ac0c0eb14565fc4c028141ff7080240b58d.PhotoRequestBuilder) {
    return ica5ec7aca95820534d087da722646ac0c0eb14565fc4c028141ff7080240b58d.NewPhotoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) PrimaryChannel()(*i5cc953f63726e53531e5d00c609bd14469e1ce23944c22b6aea0986c18fb0043.PrimaryChannelRequestBuilder) {
    return i5cc953f63726e53531e5d00c609bd14469e1ce23944c22b6aea0986c18fb0043.NewPrimaryChannelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Schedule()(*ifd53534f50d40567e607c2213e794582e29aa46a0c07e2d406db231a42a0140a.ScheduleRequestBuilder) {
    return ifd53534f50d40567e607c2213e794582e29aa46a0c07e2d406db231a42a0140a.NewScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *TeamItemRequestBuilder) SendActivityNotification()(*ie8c9d4ba2aa39351cc5580d10d45b8717b93a5d75b46c5be678a969e9175bcc6.SendActivityNotificationRequestBuilder) {
    return ie8c9d4ba2aa39351cc5580d10d45b8717b93a5d75b46c5be678a969e9175bcc6.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Tags()(*i7a26ea73a074ce439c29db16635d2ace903fc13c34496531fb310f143a72ad46.TagsRequestBuilder) {
    return i7a26ea73a074ce439c29db16635d2ace903fc13c34496531fb310f143a72ad46.NewTagsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TagsById provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) TagsById(id string)(*i289f2fd32e9a16e441f2a6b0e403b7f5f23d0c30924e6e4401734b1469fd06ce.TeamworkTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamworkTag%2Did"] = id
    }
    return i289f2fd32e9a16e441f2a6b0e403b7f5f23d0c30924e6e4401734b1469fd06ce.NewTeamworkTagItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Template()(*i04a148e32be31a86cd21b897c8b55a4508d63dbae47a02eaeb122662dbd2ff9b.TemplateRequestBuilder) {
    return i04a148e32be31a86cd21b897c8b55a4508d63dbae47a02eaeb122662dbd2ff9b.NewTemplateRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unarchive provides operations to call the unarchive method.
func (m *TeamItemRequestBuilder) Unarchive()(*if5e128abdd6cb9303edd6d603805cc08bcb3a30eb8ff5fabd66714b2846349ee.UnarchiveRequestBuilder) {
    return if5e128abdd6cb9303edd6d603805cc08bcb3a30eb8ff5fabd66714b2846349ee.NewUnarchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
