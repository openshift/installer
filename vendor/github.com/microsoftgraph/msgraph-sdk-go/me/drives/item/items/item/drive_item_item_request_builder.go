package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05c0ec45567fb3d8897550eb8ff2422453da4e6f6d020b16f038801c76613cb0 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/checkin"
    i19d6788f5023c539ef47ca9c3c8aca68e43729ccb50978593ba361436fcbc357 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/unfollow"
    i30a62b46006f460f830724a34185eabe6addc39b02c7369c642d0f0daa0e55b1 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/delta"
    i347de3ec81ee72cada957f9274a2049b35b0a5adeeea551ab42adea127d3020a "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/versions"
    i4d19529841f348c677157e7e33d1ed6edd4f81dc72950b477198cfc74d7e39f9 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/createlink"
    i5304ce2548f901afe419817e19c86188c2a91f438d46d99102bb441eeff4f0b3 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/content"
    i55bff457bac0cc82d6d6986c1093837e445c6daf83712c030ddc359db974fb0b "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i5911408c9fb6d0b7b5506ed43692ffa0e9ef7ea4aa6bb400252473412d679648 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/subscriptions"
    i5aab00d764a5f9ad9061b7231fa5a867864ca548458690794621d6663cc7494f "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/follow"
    i6310746a06c75ef864413806f6415f3cf09d4e2c5251a2476ed6aca8fb80a235 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/getactivitiesbyinterval"
    i66dc379faf68048034fe69a71cd7f87234a72e24bc6b264f77d8f56825ced3ea "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/restore"
    i697fe94c20cf8aa9d9ca8670aa74113ca512fc7e962733fbc87546eebac86d99 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/thumbnails"
    i746d7e15a556f6f04633a940990ca8cd5733284011c498aa2f3efda3d02ec850 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/listitem"
    i9388a70340cfd5e47ba875db1c2e399c8aae526124bab493eb81968bf695efc9 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/createuploadsession"
    i988afeb970a0ba81d62fe3aa5c8bd88ff30fb4b675b544bbaf5d48276346e0d1 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/invite"
    iae6a563f2ec29b53d4ce491d2675d9c4d905991ad21aca971dea26f08acd1255 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/validatepermission"
    ibaffe2d11367d67304962b88ba02d998906cb5eb0f54d8ddbb61bee317476e87 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/analytics"
    ibe02574c8f615741fc45ce6218a6cadf885d6fe1213b9a2728f6c2560f64bc8b "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/children"
    ic6086ed46c44902cd86e4fbb614c18d42889165eb4282e4c5cad4224a15bec58 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/deltawithtoken"
    ica0c5577a659c2db73294766edd0634bcfb6b46cd86fba6ef3695da20154b5fb "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/checkout"
    id9c3ac74cd32445dd9c8e64a18843db0fb93f188536cbbff19bcc3146673aa02 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/preview"
    if6a7a6c5900d9d753285a6913fb9e665cb1e4d74b8249bd99a4c3db927fb81a4 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/searchwithq"
    ifab0ee1ab8b0ac70ace5f3a9a8a6be20f9f7bcaa0a4a597c12b0193c63a01bfc "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/permissions"
    ife7f55f714b5be81f74d010999842a4ebbe531af12f5e0c350b453baff8d123c "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/copy"
    i41715f4164032be0448074d9b13e0a1f3deea014767e99c2401d7c0c2be460fa "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/permissions/item"
    i43d5a4f8297869aa1d898f64c47af83c1f2f56e21d6cbe49173818f05cdb50e6 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/versions/item"
    i6074b6299f116dcb8c8e4a9411037fe415721114dd52ebb0cc1f250577d97877 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/subscriptions/item"
    ibbd03e76ea11da931ba857caee75fe5915f32400c809adf19ab65c9715024303 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/thumbnails/item"
    ic1d7a0ae14cc9523d32ba0aac77590208358c9eb1f33aa8aae89b53a722a2f34 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item/children/item"
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
func (m *DriveItemItemRequestBuilder) Analytics()(*ibaffe2d11367d67304962b88ba02d998906cb5eb0f54d8ddbb61bee317476e87.AnalyticsRequestBuilder) {
    return ibaffe2d11367d67304962b88ba02d998906cb5eb0f54d8ddbb61bee317476e87.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *DriveItemItemRequestBuilder) Checkin()(*i05c0ec45567fb3d8897550eb8ff2422453da4e6f6d020b16f038801c76613cb0.CheckinRequestBuilder) {
    return i05c0ec45567fb3d8897550eb8ff2422453da4e6f6d020b16f038801c76613cb0.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *DriveItemItemRequestBuilder) Checkout()(*ica0c5577a659c2db73294766edd0634bcfb6b46cd86fba6ef3695da20154b5fb.CheckoutRequestBuilder) {
    return ica0c5577a659c2db73294766edd0634bcfb6b46cd86fba6ef3695da20154b5fb.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Children()(*ibe02574c8f615741fc45ce6218a6cadf885d6fe1213b9a2728f6c2560f64bc8b.ChildrenRequestBuilder) {
    return ibe02574c8f615741fc45ce6218a6cadf885d6fe1213b9a2728f6c2560f64bc8b.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ChildrenById(id string)(*ic1d7a0ae14cc9523d32ba0aac77590208358c9eb1f33aa8aae89b53a722a2f34.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did1"] = id
    }
    return ic1d7a0ae14cc9523d32ba0aac77590208358c9eb1f33aa8aae89b53a722a2f34.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemItemRequestBuilderInternal instantiates a new DriveItemItemRequestBuilder and sets the default values.
func NewDriveItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemItemRequestBuilder) {
    m := &DriveItemItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/drives/{drive%2Did}/items/{driveItem%2Did}{?%24select,%24expand}";
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
func (m *DriveItemItemRequestBuilder) Content()(*i5304ce2548f901afe419817e19c86188c2a91f438d46d99102bb441eeff4f0b3.ContentRequestBuilder) {
    return i5304ce2548f901afe419817e19c86188c2a91f438d46d99102bb441eeff4f0b3.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *DriveItemItemRequestBuilder) Copy()(*ife7f55f714b5be81f74d010999842a4ebbe531af12f5e0c350b453baff8d123c.CopyRequestBuilder) {
    return ife7f55f714b5be81f74d010999842a4ebbe531af12f5e0c350b453baff8d123c.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property items for me
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
func (m *DriveItemItemRequestBuilder) CreateLink()(*i4d19529841f348c677157e7e33d1ed6edd4f81dc72950b477198cfc74d7e39f9.CreateLinkRequestBuilder) {
    return i4d19529841f348c677157e7e33d1ed6edd4f81dc72950b477198cfc74d7e39f9.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property items in me
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
func (m *DriveItemItemRequestBuilder) CreateUploadSession()(*i9388a70340cfd5e47ba875db1c2e399c8aae526124bab493eb81968bf695efc9.CreateUploadSessionRequestBuilder) {
    return i9388a70340cfd5e47ba875db1c2e399c8aae526124bab493eb81968bf695efc9.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property items for me
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
func (m *DriveItemItemRequestBuilder) Delta()(*i30a62b46006f460f830724a34185eabe6addc39b02c7369c642d0f0daa0e55b1.DeltaRequestBuilder) {
    return i30a62b46006f460f830724a34185eabe6addc39b02c7369c642d0f0daa0e55b1.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *DriveItemItemRequestBuilder) DeltaWithToken(token *string)(*ic6086ed46c44902cd86e4fbb614c18d42889165eb4282e4c5cad4224a15bec58.DeltaWithTokenRequestBuilder) {
    return ic6086ed46c44902cd86e4fbb614c18d42889165eb4282e4c5cad4224a15bec58.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *DriveItemItemRequestBuilder) Follow()(*i5aab00d764a5f9ad9061b7231fa5a867864ca548458690794621d6663cc7494f.FollowRequestBuilder) {
    return i5aab00d764a5f9ad9061b7231fa5a867864ca548458690794621d6663cc7494f.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *DriveItemItemRequestBuilder) GetActivitiesByInterval()(*i6310746a06c75ef864413806f6415f3cf09d4e2c5251a2476ed6aca8fb80a235.GetActivitiesByIntervalRequestBuilder) {
    return i6310746a06c75ef864413806f6415f3cf09d4e2c5251a2476ed6aca8fb80a235.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *DriveItemItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i55bff457bac0cc82d6d6986c1093837e445c6daf83712c030ddc359db974fb0b.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i55bff457bac0cc82d6d6986c1093837e445c6daf83712c030ddc359db974fb0b.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *DriveItemItemRequestBuilder) Invite()(*i988afeb970a0ba81d62fe3aa5c8bd88ff30fb4b675b544bbaf5d48276346e0d1.InviteRequestBuilder) {
    return i988afeb970a0ba81d62fe3aa5c8bd88ff30fb4b675b544bbaf5d48276346e0d1.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ListItem()(*i746d7e15a556f6f04633a940990ca8cd5733284011c498aa2f3efda3d02ec850.ListItemRequestBuilder) {
    return i746d7e15a556f6f04633a940990ca8cd5733284011c498aa2f3efda3d02ec850.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property items in me
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
func (m *DriveItemItemRequestBuilder) Permissions()(*ifab0ee1ab8b0ac70ace5f3a9a8a6be20f9f7bcaa0a4a597c12b0193c63a01bfc.PermissionsRequestBuilder) {
    return ifab0ee1ab8b0ac70ace5f3a9a8a6be20f9f7bcaa0a4a597c12b0193c63a01bfc.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) PermissionsById(id string)(*i41715f4164032be0448074d9b13e0a1f3deea014767e99c2401d7c0c2be460fa.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i41715f4164032be0448074d9b13e0a1f3deea014767e99c2401d7c0c2be460fa.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *DriveItemItemRequestBuilder) Preview()(*id9c3ac74cd32445dd9c8e64a18843db0fb93f188536cbbff19bcc3146673aa02.PreviewRequestBuilder) {
    return id9c3ac74cd32445dd9c8e64a18843db0fb93f188536cbbff19bcc3146673aa02.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *DriveItemItemRequestBuilder) Restore()(*i66dc379faf68048034fe69a71cd7f87234a72e24bc6b264f77d8f56825ced3ea.RestoreRequestBuilder) {
    return i66dc379faf68048034fe69a71cd7f87234a72e24bc6b264f77d8f56825ced3ea.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemItemRequestBuilder) SearchWithQ(q *string)(*if6a7a6c5900d9d753285a6913fb9e665cb1e4d74b8249bd99a4c3db927fb81a4.SearchWithQRequestBuilder) {
    return if6a7a6c5900d9d753285a6913fb9e665cb1e4d74b8249bd99a4c3db927fb81a4.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Subscriptions()(*i5911408c9fb6d0b7b5506ed43692ffa0e9ef7ea4aa6bb400252473412d679648.SubscriptionsRequestBuilder) {
    return i5911408c9fb6d0b7b5506ed43692ffa0e9ef7ea4aa6bb400252473412d679648.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) SubscriptionsById(id string)(*i6074b6299f116dcb8c8e4a9411037fe415721114dd52ebb0cc1f250577d97877.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return i6074b6299f116dcb8c8e4a9411037fe415721114dd52ebb0cc1f250577d97877.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Thumbnails()(*i697fe94c20cf8aa9d9ca8670aa74113ca512fc7e962733fbc87546eebac86d99.ThumbnailsRequestBuilder) {
    return i697fe94c20cf8aa9d9ca8670aa74113ca512fc7e962733fbc87546eebac86d99.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ThumbnailsById(id string)(*ibbd03e76ea11da931ba857caee75fe5915f32400c809adf19ab65c9715024303.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return ibbd03e76ea11da931ba857caee75fe5915f32400c809adf19ab65c9715024303.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *DriveItemItemRequestBuilder) Unfollow()(*i19d6788f5023c539ef47ca9c3c8aca68e43729ccb50978593ba361436fcbc357.UnfollowRequestBuilder) {
    return i19d6788f5023c539ef47ca9c3c8aca68e43729ccb50978593ba361436fcbc357.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *DriveItemItemRequestBuilder) ValidatePermission()(*iae6a563f2ec29b53d4ce491d2675d9c4d905991ad21aca971dea26f08acd1255.ValidatePermissionRequestBuilder) {
    return iae6a563f2ec29b53d4ce491d2675d9c4d905991ad21aca971dea26f08acd1255.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Versions()(*i347de3ec81ee72cada957f9274a2049b35b0a5adeeea551ab42adea127d3020a.VersionsRequestBuilder) {
    return i347de3ec81ee72cada957f9274a2049b35b0a5adeeea551ab42adea127d3020a.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) VersionsById(id string)(*i43d5a4f8297869aa1d898f64c47af83c1f2f56e21d6cbe49173818f05cdb50e6.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return i43d5a4f8297869aa1d898f64c47af83c1f2f56e21d6cbe49173818f05cdb50e6.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
