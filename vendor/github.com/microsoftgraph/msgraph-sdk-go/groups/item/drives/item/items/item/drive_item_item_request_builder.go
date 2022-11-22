package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05344efea1d4be12319bb7c5b025931ba7fddb71fbb3bd3b02fd86b2045334bb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/searchwithq"
    i0d458b9aeee02ffeb65912917c6004b91b08588d3d5bb5fc6170d2a65cbedbad "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/versions"
    i1aedb2f6be136754e3647b75293759292ae86e7310b6e3787f647681d024843b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/checkin"
    i1ed2465c0650789a0c3468d361f04ea3245520b6432f7943f40a6ee9676c6acf "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/listitem"
    i23651271b73ff8571d8f44bbd1125267d399139c8c72a7c44f316b3945664ead "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/getactivitiesbyinterval"
    i24f0fbf5a7fa0c4033c74de4fd4521cbcff7b002c33de034fb040ba8ade48785 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/thumbnails"
    i278382eff19a4f6576ab54b504ed3a92af853342cf6d967e355d77b87ae5682e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/analytics"
    i3fd24d91ff9a68cd8ab305bc77c0e3a39dfa186a15011df11a4b5fc502dec13c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/checkout"
    i4036cf103711e119d410b51381a3753e915edb2604ae85b0bc6a6bacbdc6bb4e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/permissions"
    i42fa059243d4d6f94a8d8af65655d8380e7295023bac261d178a44d8de892879 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/createuploadsession"
    i4b8bbe9046c6571e0165aec127cfeb3543bcdf885ff70152863ddf40e7d0d2c1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/deltawithtoken"
    i678dc3bfe43b0f4ba842f2afe235234266d0e069c4762ad06375c7259ea088e1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/delta"
    i78fc85cea17bb32600e14078c87c3a54df2c187baf5debea9e79c13dde469c56 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/createlink"
    i8750b9a7d5153cfba7feee458bc38687ae022fe35f765ca82b44d4c8a5f9b75d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/validatepermission"
    i8be9918444eb916cc3b16cb4512b4d12db1c08417ad7ac322e7012d482ece5a8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/children"
    ia3eb3e928c003d78a015e8909eccce1b1e621a0dce53019f0a349300b40486ff "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/invite"
    iaa1dd4dbfef1600abbc07910a8f9ebdee0b731f1f52de14cd5c5a409c903b514 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/restore"
    ic5b5a9c76b232e9f8309205eb2439d8c1a441bb1845d3820b279a86439929575 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/preview"
    icfd37dc8fd55b0377aa523f6ea282c008004d48bda1d53f9d09743b07c36bf38 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    id2af1f91c111b065cd00ac5fd2ae94d9a7deed35ac0f9b58aa0af63831ae7540 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/subscriptions"
    id7b09db828cbb5e8824f7e9cf9185a1352772de9559fdcf810f71b9b224d80ca "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/unfollow"
    iddc9789ef72548a22e35a0f02875ce215d7550c8e58bc2eec1ed3ff5ea63b800 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/copy"
    idf9137d634486c6e4658161bdc130d7295285da61f119f49788b296545899115 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/follow"
    iebeaca5006903bd5bb81da4af078255c578cf1de4c76a9091cd5f04314c57b3c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/content"
    i062cd71ffa356d3cb382ab1456aa8574c689d6f8c9f7af3b69ce17e7674545fd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/permissions/item"
    i17c6bf71a98febcbc5031e9fa22b3ff715c3bc4d9d36ca3e7e0eda5b23026a9f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/thumbnails/item"
    i2fb6c5629ec400178c97aa38cc849d45535ca7b9b401ff03636d95de92a25651 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/children/item"
    ic5a933d516483ca370b2accb708efb949aaef36f165b4030548eb3650108c39a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/versions/item"
    ie5031af962b4b44302cad16972559e77bfba62775a5268d679d5453d0f1659ae "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item/subscriptions/item"
)

// DriveItemItemRequestBuilder provides operations to manage the items property of the microsoft.graph.drive entity.
type DriveItemItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DriveItemItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DriveItemItemRequestBuilderGetQueryParameters all items contained in the drive. Read-only. Nullable.
type DriveItemItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DriveItemItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DriveItemItemRequestBuilderGetQueryParameters
}
// DriveItemItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Analytics()(*i278382eff19a4f6576ab54b504ed3a92af853342cf6d967e355d77b87ae5682e.AnalyticsRequestBuilder) {
    return i278382eff19a4f6576ab54b504ed3a92af853342cf6d967e355d77b87ae5682e.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *DriveItemItemRequestBuilder) Checkin()(*i1aedb2f6be136754e3647b75293759292ae86e7310b6e3787f647681d024843b.CheckinRequestBuilder) {
    return i1aedb2f6be136754e3647b75293759292ae86e7310b6e3787f647681d024843b.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *DriveItemItemRequestBuilder) Checkout()(*i3fd24d91ff9a68cd8ab305bc77c0e3a39dfa186a15011df11a4b5fc502dec13c.CheckoutRequestBuilder) {
    return i3fd24d91ff9a68cd8ab305bc77c0e3a39dfa186a15011df11a4b5fc502dec13c.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Children()(*i8be9918444eb916cc3b16cb4512b4d12db1c08417ad7ac322e7012d482ece5a8.ChildrenRequestBuilder) {
    return i8be9918444eb916cc3b16cb4512b4d12db1c08417ad7ac322e7012d482ece5a8.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ChildrenById(id string)(*i2fb6c5629ec400178c97aa38cc849d45535ca7b9b401ff03636d95de92a25651.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did1"] = id
    }
    return i2fb6c5629ec400178c97aa38cc849d45535ca7b9b401ff03636d95de92a25651.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemItemRequestBuilderInternal instantiates a new DriveItemItemRequestBuilder and sets the default values.
func NewDriveItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemItemRequestBuilder) {
    m := &DriveItemItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/drives/{drive%2Did}/items/{driveItem%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDriveItemItemRequestBuilder instantiates a new DriveItemItemRequestBuilder and sets the default values.
func NewDriveItemItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDriveItemItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Content provides operations to manage the media for the group entity.
func (m *DriveItemItemRequestBuilder) Content()(*iebeaca5006903bd5bb81da4af078255c578cf1de4c76a9091cd5f04314c57b3c.ContentRequestBuilder) {
    return iebeaca5006903bd5bb81da4af078255c578cf1de4c76a9091cd5f04314c57b3c.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *DriveItemItemRequestBuilder) Copy()(*iddc9789ef72548a22e35a0f02875ce215d7550c8e58bc2eec1ed3ff5ea63b800.CopyRequestBuilder) {
    return iddc9789ef72548a22e35a0f02875ce215d7550c8e58bc2eec1ed3ff5ea63b800.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property items for groups
func (m *DriveItemItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *DriveItemItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation all items contained in the drive. Read-only. Nullable.
func (m *DriveItemItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DriveItemItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateLink provides operations to call the createLink method.
func (m *DriveItemItemRequestBuilder) CreateLink()(*i78fc85cea17bb32600e14078c87c3a54df2c187baf5debea9e79c13dde469c56.CreateLinkRequestBuilder) {
    return i78fc85cea17bb32600e14078c87c3a54df2c187baf5debea9e79c13dde469c56.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property items in groups
func (m *DriveItemItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *DriveItemItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateUploadSession provides operations to call the createUploadSession method.
func (m *DriveItemItemRequestBuilder) CreateUploadSession()(*i42fa059243d4d6f94a8d8af65655d8380e7295023bac261d178a44d8de892879.CreateUploadSessionRequestBuilder) {
    return i42fa059243d4d6f94a8d8af65655d8380e7295023bac261d178a44d8de892879.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property items for groups
func (m *DriveItemItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DriveItemItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Delta provides operations to call the delta method.
func (m *DriveItemItemRequestBuilder) Delta()(*i678dc3bfe43b0f4ba842f2afe235234266d0e069c4762ad06375c7259ea088e1.DeltaRequestBuilder) {
    return i678dc3bfe43b0f4ba842f2afe235234266d0e069c4762ad06375c7259ea088e1.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *DriveItemItemRequestBuilder) DeltaWithToken(token *string)(*i4b8bbe9046c6571e0165aec127cfeb3543bcdf885ff70152863ddf40e7d0d2c1.DeltaWithTokenRequestBuilder) {
    return i4b8bbe9046c6571e0165aec127cfeb3543bcdf885ff70152863ddf40e7d0d2c1.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *DriveItemItemRequestBuilder) Follow()(*idf9137d634486c6e4658161bdc130d7295285da61f119f49788b296545899115.FollowRequestBuilder) {
    return idf9137d634486c6e4658161bdc130d7295285da61f119f49788b296545899115.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get all items contained in the drive. Read-only. Nullable.
func (m *DriveItemItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DriveItemItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable), nil
}
// GetActivitiesByInterval provides operations to call the getActivitiesByInterval method.
func (m *DriveItemItemRequestBuilder) GetActivitiesByInterval()(*i23651271b73ff8571d8f44bbd1125267d399139c8c72a7c44f316b3945664ead.GetActivitiesByIntervalRequestBuilder) {
    return i23651271b73ff8571d8f44bbd1125267d399139c8c72a7c44f316b3945664ead.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *DriveItemItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*icfd37dc8fd55b0377aa523f6ea282c008004d48bda1d53f9d09743b07c36bf38.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return icfd37dc8fd55b0377aa523f6ea282c008004d48bda1d53f9d09743b07c36bf38.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *DriveItemItemRequestBuilder) Invite()(*ia3eb3e928c003d78a015e8909eccce1b1e621a0dce53019f0a349300b40486ff.InviteRequestBuilder) {
    return ia3eb3e928c003d78a015e8909eccce1b1e621a0dce53019f0a349300b40486ff.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ListItem()(*i1ed2465c0650789a0c3468d361f04ea3245520b6432f7943f40a6ee9676c6acf.ListItemRequestBuilder) {
    return i1ed2465c0650789a0c3468d361f04ea3245520b6432f7943f40a6ee9676c6acf.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property items in groups
func (m *DriveItemItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *DriveItemItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable), nil
}
// Permissions provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Permissions()(*i4036cf103711e119d410b51381a3753e915edb2604ae85b0bc6a6bacbdc6bb4e.PermissionsRequestBuilder) {
    return i4036cf103711e119d410b51381a3753e915edb2604ae85b0bc6a6bacbdc6bb4e.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) PermissionsById(id string)(*i062cd71ffa356d3cb382ab1456aa8574c689d6f8c9f7af3b69ce17e7674545fd.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i062cd71ffa356d3cb382ab1456aa8574c689d6f8c9f7af3b69ce17e7674545fd.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *DriveItemItemRequestBuilder) Preview()(*ic5b5a9c76b232e9f8309205eb2439d8c1a441bb1845d3820b279a86439929575.PreviewRequestBuilder) {
    return ic5b5a9c76b232e9f8309205eb2439d8c1a441bb1845d3820b279a86439929575.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *DriveItemItemRequestBuilder) Restore()(*iaa1dd4dbfef1600abbc07910a8f9ebdee0b731f1f52de14cd5c5a409c903b514.RestoreRequestBuilder) {
    return iaa1dd4dbfef1600abbc07910a8f9ebdee0b731f1f52de14cd5c5a409c903b514.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemItemRequestBuilder) SearchWithQ(q *string)(*i05344efea1d4be12319bb7c5b025931ba7fddb71fbb3bd3b02fd86b2045334bb.SearchWithQRequestBuilder) {
    return i05344efea1d4be12319bb7c5b025931ba7fddb71fbb3bd3b02fd86b2045334bb.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Subscriptions()(*id2af1f91c111b065cd00ac5fd2ae94d9a7deed35ac0f9b58aa0af63831ae7540.SubscriptionsRequestBuilder) {
    return id2af1f91c111b065cd00ac5fd2ae94d9a7deed35ac0f9b58aa0af63831ae7540.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) SubscriptionsById(id string)(*ie5031af962b4b44302cad16972559e77bfba62775a5268d679d5453d0f1659ae.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return ie5031af962b4b44302cad16972559e77bfba62775a5268d679d5453d0f1659ae.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Thumbnails()(*i24f0fbf5a7fa0c4033c74de4fd4521cbcff7b002c33de034fb040ba8ade48785.ThumbnailsRequestBuilder) {
    return i24f0fbf5a7fa0c4033c74de4fd4521cbcff7b002c33de034fb040ba8ade48785.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ThumbnailsById(id string)(*i17c6bf71a98febcbc5031e9fa22b3ff715c3bc4d9d36ca3e7e0eda5b23026a9f.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return i17c6bf71a98febcbc5031e9fa22b3ff715c3bc4d9d36ca3e7e0eda5b23026a9f.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *DriveItemItemRequestBuilder) Unfollow()(*id7b09db828cbb5e8824f7e9cf9185a1352772de9559fdcf810f71b9b224d80ca.UnfollowRequestBuilder) {
    return id7b09db828cbb5e8824f7e9cf9185a1352772de9559fdcf810f71b9b224d80ca.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *DriveItemItemRequestBuilder) ValidatePermission()(*i8750b9a7d5153cfba7feee458bc38687ae022fe35f765ca82b44d4c8a5f9b75d.ValidatePermissionRequestBuilder) {
    return i8750b9a7d5153cfba7feee458bc38687ae022fe35f765ca82b44d4c8a5f9b75d.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Versions()(*i0d458b9aeee02ffeb65912917c6004b91b08588d3d5bb5fc6170d2a65cbedbad.VersionsRequestBuilder) {
    return i0d458b9aeee02ffeb65912917c6004b91b08588d3d5bb5fc6170d2a65cbedbad.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) VersionsById(id string)(*ic5a933d516483ca370b2accb708efb949aaef36f165b4030548eb3650108c39a.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return ic5a933d516483ca370b2accb708efb949aaef36f165b4030548eb3650108c39a.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
