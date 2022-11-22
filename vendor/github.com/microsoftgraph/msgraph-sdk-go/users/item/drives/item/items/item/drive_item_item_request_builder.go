package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i122eddd996f16b226c48c377d84c38e03669c5a21c3ee2707526b309df2c302d "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/restore"
    i256c3973c754c32db343a69a5559b486c82f4ea275d1e6fd7bb118a5a50e60af "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/thumbnails"
    i2ce61a24701dcd9d8631bf4f41aea243edb561f20f6457d8dbffe771cfeda605 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/unfollow"
    i2f4222c1ed68afebf22bf0da4254fdf6af0f49a584d8f64c30d84619628c5702 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/createuploadsession"
    i3aaa21b8eba72994150daba77d9145bb89dbc0362a7056b1e55dcd1a378cf29f "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/preview"
    i3d788c112ff8e894a71cea0e2d9b345f866197fb045501f32a77b990aae0bfb0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i4ad95a6c78d2bbc0d551f7f83feff56a755fe274173046a052c8f16440e2285d "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/analytics"
    i579e57cbfb4d8a5e7cc8d94b574194fc4f4bd6e2d79312f1e1157dd947b65ba3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/checkin"
    i6144d1beb34ba8a816c7cda0a208dcbbb5369dbe1cdfa63b8497560473fe56da "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/checkout"
    i62abe1cd2ad3f6c1b7996efc49d1e4d60bef31a3791278f28040c4154802a94a "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/delta"
    i6d37f4d3da93fc052cf83cefb3a3fe9a2e878b6566d73a2c1fc3071c8c7a2f33 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/subscriptions"
    i83008558442e64c5bf0a1d83617b26864b62c96dbd9c076a5f42b595dce1422c "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/validatepermission"
    i8990c7872938b401500cc1bcda9ebf1b6b6e0a83ed110ee83cf0c057aba93c2c "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/listitem"
    i89c6e89124333138b9b97beae347587f5f59a6e35a85975682a5cb2c642ba3c9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/permissions"
    i964db02df3900361c40d2725931853cf67f5d0cf823e0f303de7f1ff59cc5b74 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/versions"
    ia4abe4a62946631fa5ea220db620b019268baefc42cafcb20beb7c7b445ea1a7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/searchwithq"
    ia6e4a32ab59e8a0dac20d9c8c675d1e9769dce33343b20563e4f38d128f1c87b "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/deltawithtoken"
    iaa7f1f161a46ee80ab49f872d06a2f7ce23c8f9b318af39a83d20351f9f40117 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/invite"
    ib765f5039e518fa8697650fac3d5a116e8e30a2ed1d8c8078f20263cd324b937 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/getactivitiesbyinterval"
    ib9cce61d016695c26fef8b858bf13cb530681491bb6ed31832d66e9416b81fab "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/copy"
    ibd487c2b87dca3c99ca7ed3a08867086ddeee165038524010541c989b2bbfda5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/children"
    ibe31a294249e6c09f9e5d7e969376f9cacacfdc79edfebcb3ad110b0b1201974 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/content"
    ic6e54ecc33d74165b9e0937cf9a8079a78c785bc408d793b6a557e14d0746f01 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/follow"
    if893ac1d6b2c84a0589a854c08136626d29bbac6dfa19fdf7f6fbc0661e6dbaf "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/createlink"
    i425f10a4034e47270e96c2be09315414c174bf653b7b820427cd68b4fd005678 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/children/item"
    i618b1a4f48a82b0f679af0dabcedafb590ed0ddb82ade1d6b7303352f8ba3a33 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/thumbnails/item"
    i79508a4e60d6ec9761fd29cdf217891cdb087b26380379217248e3e3d3f0754e "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/versions/item"
    ibbf87a9d0822f48c3ac14f2e636b7123f483b764d9f9915d2f4e696e3368faf6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/permissions/item"
    iee83a084ce5d1fc2c76bad7adfa9a9ed98b719098045e1c3d5837aca5bf78bb2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item/subscriptions/item"
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
func (m *DriveItemItemRequestBuilder) Analytics()(*i4ad95a6c78d2bbc0d551f7f83feff56a755fe274173046a052c8f16440e2285d.AnalyticsRequestBuilder) {
    return i4ad95a6c78d2bbc0d551f7f83feff56a755fe274173046a052c8f16440e2285d.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *DriveItemItemRequestBuilder) Checkin()(*i579e57cbfb4d8a5e7cc8d94b574194fc4f4bd6e2d79312f1e1157dd947b65ba3.CheckinRequestBuilder) {
    return i579e57cbfb4d8a5e7cc8d94b574194fc4f4bd6e2d79312f1e1157dd947b65ba3.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *DriveItemItemRequestBuilder) Checkout()(*i6144d1beb34ba8a816c7cda0a208dcbbb5369dbe1cdfa63b8497560473fe56da.CheckoutRequestBuilder) {
    return i6144d1beb34ba8a816c7cda0a208dcbbb5369dbe1cdfa63b8497560473fe56da.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Children()(*ibd487c2b87dca3c99ca7ed3a08867086ddeee165038524010541c989b2bbfda5.ChildrenRequestBuilder) {
    return ibd487c2b87dca3c99ca7ed3a08867086ddeee165038524010541c989b2bbfda5.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ChildrenById(id string)(*i425f10a4034e47270e96c2be09315414c174bf653b7b820427cd68b4fd005678.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did1"] = id
    }
    return i425f10a4034e47270e96c2be09315414c174bf653b7b820427cd68b4fd005678.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemItemRequestBuilderInternal instantiates a new DriveItemItemRequestBuilder and sets the default values.
func NewDriveItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemItemRequestBuilder) {
    m := &DriveItemItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/drives/{drive%2Did}/items/{driveItem%2Did}{?%24select,%24expand}";
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
// Content provides operations to manage the media for the user entity.
func (m *DriveItemItemRequestBuilder) Content()(*ibe31a294249e6c09f9e5d7e969376f9cacacfdc79edfebcb3ad110b0b1201974.ContentRequestBuilder) {
    return ibe31a294249e6c09f9e5d7e969376f9cacacfdc79edfebcb3ad110b0b1201974.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *DriveItemItemRequestBuilder) Copy()(*ib9cce61d016695c26fef8b858bf13cb530681491bb6ed31832d66e9416b81fab.CopyRequestBuilder) {
    return ib9cce61d016695c26fef8b858bf13cb530681491bb6ed31832d66e9416b81fab.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property items for users
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
func (m *DriveItemItemRequestBuilder) CreateLink()(*if893ac1d6b2c84a0589a854c08136626d29bbac6dfa19fdf7f6fbc0661e6dbaf.CreateLinkRequestBuilder) {
    return if893ac1d6b2c84a0589a854c08136626d29bbac6dfa19fdf7f6fbc0661e6dbaf.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property items in users
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
func (m *DriveItemItemRequestBuilder) CreateUploadSession()(*i2f4222c1ed68afebf22bf0da4254fdf6af0f49a584d8f64c30d84619628c5702.CreateUploadSessionRequestBuilder) {
    return i2f4222c1ed68afebf22bf0da4254fdf6af0f49a584d8f64c30d84619628c5702.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property items for users
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
func (m *DriveItemItemRequestBuilder) Delta()(*i62abe1cd2ad3f6c1b7996efc49d1e4d60bef31a3791278f28040c4154802a94a.DeltaRequestBuilder) {
    return i62abe1cd2ad3f6c1b7996efc49d1e4d60bef31a3791278f28040c4154802a94a.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *DriveItemItemRequestBuilder) DeltaWithToken(token *string)(*ia6e4a32ab59e8a0dac20d9c8c675d1e9769dce33343b20563e4f38d128f1c87b.DeltaWithTokenRequestBuilder) {
    return ia6e4a32ab59e8a0dac20d9c8c675d1e9769dce33343b20563e4f38d128f1c87b.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *DriveItemItemRequestBuilder) Follow()(*ic6e54ecc33d74165b9e0937cf9a8079a78c785bc408d793b6a557e14d0746f01.FollowRequestBuilder) {
    return ic6e54ecc33d74165b9e0937cf9a8079a78c785bc408d793b6a557e14d0746f01.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *DriveItemItemRequestBuilder) GetActivitiesByInterval()(*ib765f5039e518fa8697650fac3d5a116e8e30a2ed1d8c8078f20263cd324b937.GetActivitiesByIntervalRequestBuilder) {
    return ib765f5039e518fa8697650fac3d5a116e8e30a2ed1d8c8078f20263cd324b937.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *DriveItemItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i3d788c112ff8e894a71cea0e2d9b345f866197fb045501f32a77b990aae0bfb0.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i3d788c112ff8e894a71cea0e2d9b345f866197fb045501f32a77b990aae0bfb0.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *DriveItemItemRequestBuilder) Invite()(*iaa7f1f161a46ee80ab49f872d06a2f7ce23c8f9b318af39a83d20351f9f40117.InviteRequestBuilder) {
    return iaa7f1f161a46ee80ab49f872d06a2f7ce23c8f9b318af39a83d20351f9f40117.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ListItem()(*i8990c7872938b401500cc1bcda9ebf1b6b6e0a83ed110ee83cf0c057aba93c2c.ListItemRequestBuilder) {
    return i8990c7872938b401500cc1bcda9ebf1b6b6e0a83ed110ee83cf0c057aba93c2c.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property items in users
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
func (m *DriveItemItemRequestBuilder) Permissions()(*i89c6e89124333138b9b97beae347587f5f59a6e35a85975682a5cb2c642ba3c9.PermissionsRequestBuilder) {
    return i89c6e89124333138b9b97beae347587f5f59a6e35a85975682a5cb2c642ba3c9.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) PermissionsById(id string)(*ibbf87a9d0822f48c3ac14f2e636b7123f483b764d9f9915d2f4e696e3368faf6.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return ibbf87a9d0822f48c3ac14f2e636b7123f483b764d9f9915d2f4e696e3368faf6.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *DriveItemItemRequestBuilder) Preview()(*i3aaa21b8eba72994150daba77d9145bb89dbc0362a7056b1e55dcd1a378cf29f.PreviewRequestBuilder) {
    return i3aaa21b8eba72994150daba77d9145bb89dbc0362a7056b1e55dcd1a378cf29f.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *DriveItemItemRequestBuilder) Restore()(*i122eddd996f16b226c48c377d84c38e03669c5a21c3ee2707526b309df2c302d.RestoreRequestBuilder) {
    return i122eddd996f16b226c48c377d84c38e03669c5a21c3ee2707526b309df2c302d.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemItemRequestBuilder) SearchWithQ(q *string)(*ia4abe4a62946631fa5ea220db620b019268baefc42cafcb20beb7c7b445ea1a7.SearchWithQRequestBuilder) {
    return ia4abe4a62946631fa5ea220db620b019268baefc42cafcb20beb7c7b445ea1a7.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Subscriptions()(*i6d37f4d3da93fc052cf83cefb3a3fe9a2e878b6566d73a2c1fc3071c8c7a2f33.SubscriptionsRequestBuilder) {
    return i6d37f4d3da93fc052cf83cefb3a3fe9a2e878b6566d73a2c1fc3071c8c7a2f33.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) SubscriptionsById(id string)(*iee83a084ce5d1fc2c76bad7adfa9a9ed98b719098045e1c3d5837aca5bf78bb2.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return iee83a084ce5d1fc2c76bad7adfa9a9ed98b719098045e1c3d5837aca5bf78bb2.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Thumbnails()(*i256c3973c754c32db343a69a5559b486c82f4ea275d1e6fd7bb118a5a50e60af.ThumbnailsRequestBuilder) {
    return i256c3973c754c32db343a69a5559b486c82f4ea275d1e6fd7bb118a5a50e60af.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ThumbnailsById(id string)(*i618b1a4f48a82b0f679af0dabcedafb590ed0ddb82ade1d6b7303352f8ba3a33.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return i618b1a4f48a82b0f679af0dabcedafb590ed0ddb82ade1d6b7303352f8ba3a33.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *DriveItemItemRequestBuilder) Unfollow()(*i2ce61a24701dcd9d8631bf4f41aea243edb561f20f6457d8dbffe771cfeda605.UnfollowRequestBuilder) {
    return i2ce61a24701dcd9d8631bf4f41aea243edb561f20f6457d8dbffe771cfeda605.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *DriveItemItemRequestBuilder) ValidatePermission()(*i83008558442e64c5bf0a1d83617b26864b62c96dbd9c076a5f42b595dce1422c.ValidatePermissionRequestBuilder) {
    return i83008558442e64c5bf0a1d83617b26864b62c96dbd9c076a5f42b595dce1422c.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Versions()(*i964db02df3900361c40d2725931853cf67f5d0cf823e0f303de7f1ff59cc5b74.VersionsRequestBuilder) {
    return i964db02df3900361c40d2725931853cf67f5d0cf823e0f303de7f1ff59cc5b74.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) VersionsById(id string)(*i79508a4e60d6ec9761fd29cdf217891cdb087b26380379217248e3e3d3f0754e.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return i79508a4e60d6ec9761fd29cdf217891cdb087b26380379217248e3e3d3f0754e.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
