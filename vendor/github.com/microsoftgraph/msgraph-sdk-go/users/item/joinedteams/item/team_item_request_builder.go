package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2600ad0538dd095947610852034037aed3bb130d1317dc4866865ae049bf71de "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/completemigration"
    i5187612c80915156278ab470c00ab17b5872de60712ebc154fe8b5a67d372c0b "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/photo"
    i5363721350f81d4226008b2a0abc5fa76fabff34763dc7de5066c621e0c0f6be "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/operations"
    i539671039e993ec9795209a4951310067d75618d6c4524b564660e6d7f6dec35 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/template"
    i66dfe7a7c3387fa20a0d979d3bfba6fd40a47368ae24042ec7def4c570e011bb "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/archive"
    i6ea3db8c980a5b6b549dcba6aaa6b387b8c6b1cb073e6178a10627e6fb06a389 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/primarychannel"
    i81bd7dd61df3d95dbc4620ad1479451afe1469239e95f230e4f09a46b785da75 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/schedule"
    i8f24546b063616adb950ffa448f634e2722a9ef55b17354a587109842784dc87 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/clone"
    i95a7fa386ac768e6cf2f520a17857a2cdde3f6c2daf61c03b795675a411b4a10 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/unarchive"
    ib1deb2599b9a847d5026bfa2b1b44e59f2a829010f536668074417c636ce36db "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels"
    ic07997e8175d4979dbf6a5fdc175868456c576f787269936469ca3cbe98bf07e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/allchannels"
    ic44d4f0e5acd07f5481d8f8b84a79cab11dfbee059cf6770fd4e2b23db0c1a93 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/group"
    ic914f26c893a65c2d7a0c81250eacc4c0f9d860f64ab1b1cf1c364c6743699a8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/installedapps"
    icb354bee8864e0c1131259c04983ee9ce8d3bf6fded1f9c90997e6b8b476d0ff "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/tags"
    iccac5382005a5ca2eac80d35c6684f5a5d2a95fd6a785277685abb4bca00f086 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/incomingchannels"
    ie7bd69bc163dd5d2c6dbbe5978d8c7eab2d6a90de4ac1a9c005cb73477cec6bb "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/sendactivitynotification"
    ief730b7100e12b3385f77770a9d755bebd575023d52fdbb9414e8c86230b027d "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/members"
    i317663812c69a7c8ef4e5e9c721877170ab1423cef01deb27befa0de404a522b "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/allchannels/item"
    i6d038dcd104603d7281b4eac363006d79b67e92eca3e8c364b49d02b84649d79 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/channels/item"
    i98b219e6f0357bddb879690365c172b8f5fa9226021c11592a8062b6d8a9415e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/members/item"
    ic36d8fd750746e0ef3bba4a0ac0d5869787d2d7347af1b8a0ff929a9af209bf4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/incomingchannels/item"
    ie5af489748d2fc36f15958b988ad08901566118f084112092561d1d1b446b768 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/installedapps/item"
    ie78978d46a3704a561b7ecbd9e66a5ee9be50b99960e26229ec758c80ee98271 "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/operations/item"
    if61ac464c1526ee232b6bb1f4cbe350432b784a34a1cfcb34d9d61824fecb42e "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item/tags/item"
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
// TeamItemRequestBuilderGetQueryParameters get joinedTeams from users
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
func (m *TeamItemRequestBuilder) AllChannels()(*ic07997e8175d4979dbf6a5fdc175868456c576f787269936469ca3cbe98bf07e.AllChannelsRequestBuilder) {
    return ic07997e8175d4979dbf6a5fdc175868456c576f787269936469ca3cbe98bf07e.NewAllChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllChannelsById provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) AllChannelsById(id string)(*i317663812c69a7c8ef4e5e9c721877170ab1423cef01deb27befa0de404a522b.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i317663812c69a7c8ef4e5e9c721877170ab1423cef01deb27befa0de404a522b.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Archive provides operations to call the archive method.
func (m *TeamItemRequestBuilder) Archive()(*i66dfe7a7c3387fa20a0d979d3bfba6fd40a47368ae24042ec7def4c570e011bb.ArchiveRequestBuilder) {
    return i66dfe7a7c3387fa20a0d979d3bfba6fd40a47368ae24042ec7def4c570e011bb.NewArchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Channels()(*ib1deb2599b9a847d5026bfa2b1b44e59f2a829010f536668074417c636ce36db.ChannelsRequestBuilder) {
    return ib1deb2599b9a847d5026bfa2b1b44e59f2a829010f536668074417c636ce36db.NewChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChannelsById provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) ChannelsById(id string)(*i6d038dcd104603d7281b4eac363006d79b67e92eca3e8c364b49d02b84649d79.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i6d038dcd104603d7281b4eac363006d79b67e92eca3e8c364b49d02b84649d79.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Clone provides operations to call the clone method.
func (m *TeamItemRequestBuilder) Clone()(*i8f24546b063616adb950ffa448f634e2722a9ef55b17354a587109842784dc87.CloneRequestBuilder) {
    return i8f24546b063616adb950ffa448f634e2722a9ef55b17354a587109842784dc87.NewCloneRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *TeamItemRequestBuilder) CompleteMigration()(*i2600ad0538dd095947610852034037aed3bb130d1317dc4866865ae049bf71de.CompleteMigrationRequestBuilder) {
    return i2600ad0538dd095947610852034037aed3bb130d1317dc4866865ae049bf71de.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTeamItemRequestBuilderInternal instantiates a new TeamItemRequestBuilder and sets the default values.
func NewTeamItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamItemRequestBuilder) {
    m := &TeamItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property joinedTeams for users
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
// CreateGetRequestInformation get joinedTeams from users
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
// CreatePatchRequestInformation update the navigation property joinedTeams in users
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
// Delete delete navigation property joinedTeams for users
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
// Get get joinedTeams from users
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
func (m *TeamItemRequestBuilder) Group()(*ic44d4f0e5acd07f5481d8f8b84a79cab11dfbee059cf6770fd4e2b23db0c1a93.GroupRequestBuilder) {
    return ic44d4f0e5acd07f5481d8f8b84a79cab11dfbee059cf6770fd4e2b23db0c1a93.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannels()(*iccac5382005a5ca2eac80d35c6684f5a5d2a95fd6a785277685abb4bca00f086.IncomingChannelsRequestBuilder) {
    return iccac5382005a5ca2eac80d35c6684f5a5d2a95fd6a785277685abb4bca00f086.NewIncomingChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannelsById provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannelsById(id string)(*ic36d8fd750746e0ef3bba4a0ac0d5869787d2d7347af1b8a0ff929a9af209bf4.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return ic36d8fd750746e0ef3bba4a0ac0d5869787d2d7347af1b8a0ff929a9af209bf4.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledApps()(*ic914f26c893a65c2d7a0c81250eacc4c0f9d860f64ab1b1cf1c364c6743699a8.InstalledAppsRequestBuilder) {
    return ic914f26c893a65c2d7a0c81250eacc4c0f9d860f64ab1b1cf1c364c6743699a8.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledAppsById(id string)(*ie5af489748d2fc36f15958b988ad08901566118f084112092561d1d1b446b768.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return ie5af489748d2fc36f15958b988ad08901566118f084112092561d1d1b446b768.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Members()(*ief730b7100e12b3385f77770a9d755bebd575023d52fdbb9414e8c86230b027d.MembersRequestBuilder) {
    return ief730b7100e12b3385f77770a9d755bebd575023d52fdbb9414e8c86230b027d.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) MembersById(id string)(*i98b219e6f0357bddb879690365c172b8f5fa9226021c11592a8062b6d8a9415e.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i98b219e6f0357bddb879690365c172b8f5fa9226021c11592a8062b6d8a9415e.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Operations()(*i5363721350f81d4226008b2a0abc5fa76fabff34763dc7de5066c621e0c0f6be.OperationsRequestBuilder) {
    return i5363721350f81d4226008b2a0abc5fa76fabff34763dc7de5066c621e0c0f6be.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) OperationsById(id string)(*ie78978d46a3704a561b7ecbd9e66a5ee9be50b99960e26229ec758c80ee98271.TeamsAsyncOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAsyncOperation%2Did"] = id
    }
    return ie78978d46a3704a561b7ecbd9e66a5ee9be50b99960e26229ec758c80ee98271.NewTeamsAsyncOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property joinedTeams in users
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
func (m *TeamItemRequestBuilder) Photo()(*i5187612c80915156278ab470c00ab17b5872de60712ebc154fe8b5a67d372c0b.PhotoRequestBuilder) {
    return i5187612c80915156278ab470c00ab17b5872de60712ebc154fe8b5a67d372c0b.NewPhotoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) PrimaryChannel()(*i6ea3db8c980a5b6b549dcba6aaa6b387b8c6b1cb073e6178a10627e6fb06a389.PrimaryChannelRequestBuilder) {
    return i6ea3db8c980a5b6b549dcba6aaa6b387b8c6b1cb073e6178a10627e6fb06a389.NewPrimaryChannelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Schedule()(*i81bd7dd61df3d95dbc4620ad1479451afe1469239e95f230e4f09a46b785da75.ScheduleRequestBuilder) {
    return i81bd7dd61df3d95dbc4620ad1479451afe1469239e95f230e4f09a46b785da75.NewScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *TeamItemRequestBuilder) SendActivityNotification()(*ie7bd69bc163dd5d2c6dbbe5978d8c7eab2d6a90de4ac1a9c005cb73477cec6bb.SendActivityNotificationRequestBuilder) {
    return ie7bd69bc163dd5d2c6dbbe5978d8c7eab2d6a90de4ac1a9c005cb73477cec6bb.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Tags()(*icb354bee8864e0c1131259c04983ee9ce8d3bf6fded1f9c90997e6b8b476d0ff.TagsRequestBuilder) {
    return icb354bee8864e0c1131259c04983ee9ce8d3bf6fded1f9c90997e6b8b476d0ff.NewTagsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TagsById provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) TagsById(id string)(*if61ac464c1526ee232b6bb1f4cbe350432b784a34a1cfcb34d9d61824fecb42e.TeamworkTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamworkTag%2Did"] = id
    }
    return if61ac464c1526ee232b6bb1f4cbe350432b784a34a1cfcb34d9d61824fecb42e.NewTeamworkTagItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Template()(*i539671039e993ec9795209a4951310067d75618d6c4524b564660e6d7f6dec35.TemplateRequestBuilder) {
    return i539671039e993ec9795209a4951310067d75618d6c4524b564660e6d7f6dec35.NewTemplateRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unarchive provides operations to call the unarchive method.
func (m *TeamItemRequestBuilder) Unarchive()(*i95a7fa386ac768e6cf2f520a17857a2cdde3f6c2daf61c03b795675a411b4a10.UnarchiveRequestBuilder) {
    return i95a7fa386ac768e6cf2f520a17857a2cdde3f6c2daf61c03b795675a411b4a10.NewUnarchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
