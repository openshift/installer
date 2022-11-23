package root

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0033efa979ecc8ac80fdb29d960015d241a7cd255681fd7cbf71d410f9604ad4 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/restore"
    i032fb091fa50ded67bd8413cc227e15a89212f22579d11ba5fcda786303b3d88 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/versions"
    i03c9455ec93075afc95f623da889cf2de63e90a3d4f5f5e8ec43ed1d37289360 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/checkin"
    i03ddf9ed2e8d64e93fc44a7258e3589c38d6f455c1358c4d75d1614935a98e46 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/children"
    i1009a2ff258bf0b86c5aa9b40c49762ddd855f338f1947b3a8f825c3d8524721 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/content"
    i27c4eb9aab96ee3918e3877e6df557c98aee08cbeebedc02e45682b3af84cf0a "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/preview"
    i404bc6886d1c9123a28011fced42055517ed97d7f967f514fb167d219b2c0db0 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/listitem"
    i5dc3a384aa3ae3c4baf5dbb42b7537ab2ed308f452fd7110946ba1eab1240e05 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/checkout"
    i71cf1b728ee8302f69703a7e6a0e0731e7b4dbecd838960e2ffda2156d16db08 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/deltawithtoken"
    i72d8e675471adc6caea657caad244220759b7b1d1af78e568f69f1b0488ae9c2 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/copy"
    i7669f738788c84c7ecf0251b8089bc20ff3a6e5b410b8045267fa73a34277705 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/analytics"
    i79c6c3c119b0fd26cac016493adf8593a63ae4fab341dd8853915459798f15db "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i7b2a9d9e2b194b78db259bec4bcb10cb8c8747663a268008a135977ab274f5ba "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/getactivitiesbyinterval"
    i7e2854ef00272a1811fdcaeaf5c86f99d0ec13adf274e99eac8180b08d5f8ae1 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/createlink"
    i80b1e91bf7b013cdf8dbdf9b8d8771f6a1b6e697578eefaecea07952add6388c "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/permissions"
    ia80c39dd409571ac746317206288180c5750a50b838f9b4840df1960308280e9 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/unfollow"
    ib327b2142df568febe244d5e9882de94f5049f86c795d0adc6dc9a062df1985d "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/invite"
    ic4f6f4ea374f94b486e173595e0a59f167b50e3195232b90a72d9d013abfa7d2 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/validatepermission"
    id1bc6c7776430f10716cbe5e70e0fffffd8064c286891200f2a817b8e144fe9e "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/follow"
    idcdd5295fe11eb0b451e619b74b2ec775e69aa4c0031109a4ad4ab63a2d346c2 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/delta"
    ief1f6f4bf8be57145cbf8d0f439184c8d23abbbcc4b38699c6c809db84be6beb "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/searchwithq"
    ief5ea07106b57464e0d9313bab8da87f6427ea4ce6317fef25f88f0475243573 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/thumbnails"
    if705254df95007eb4d749c340e906c098a5584daa3d253a27dad196a9e68aea6 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/createuploadsession"
    ifce4a3fdda3e56bdf741ded90a576a9c0e9f6480576f685c4c8e5eceb911a516 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/subscriptions"
    i689f37510f9fad8e0be6242df80f8b79fff25b25c310cfaf55e7b93c51c780af "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/subscriptions/item"
    ia16000547e3bbe576a5a946472bea8c98307843b62bbdb98e519f2ae8166f2c2 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/versions/item"
    ie2d9fc455cd20ad77bcd5ffeb7d799dfb360cc67d9482d6709d02097def6a274 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/children/item"
    ieeb69da49c1200fd91a53dded9fdde76e190b3cb09517e3aab7c813ac910de6d "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/permissions/item"
    if641b83445d13c546e420bdc7e9b23cb86e18512438d6cde4308543c0d7c8cb4 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root/thumbnails/item"
)

// RootRequestBuilder provides operations to manage the root property of the microsoft.graph.drive entity.
type RootRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// RootRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RootRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// RootRequestBuilderGetQueryParameters retrieve the metadata for a driveItem in a drive by file system path or ID.`item-id` is the ID of a driveItem. It may also be the unique ID of a SharePoint list item.
type RootRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// RootRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RootRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *RootRequestBuilderGetQueryParameters
}
// RootRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RootRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Analytics provides operations to manage the analytics property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Analytics()(*i7669f738788c84c7ecf0251b8089bc20ff3a6e5b410b8045267fa73a34277705.AnalyticsRequestBuilder) {
    return i7669f738788c84c7ecf0251b8089bc20ff3a6e5b410b8045267fa73a34277705.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *RootRequestBuilder) Checkin()(*i03c9455ec93075afc95f623da889cf2de63e90a3d4f5f5e8ec43ed1d37289360.CheckinRequestBuilder) {
    return i03c9455ec93075afc95f623da889cf2de63e90a3d4f5f5e8ec43ed1d37289360.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *RootRequestBuilder) Checkout()(*i5dc3a384aa3ae3c4baf5dbb42b7537ab2ed308f452fd7110946ba1eab1240e05.CheckoutRequestBuilder) {
    return i5dc3a384aa3ae3c4baf5dbb42b7537ab2ed308f452fd7110946ba1eab1240e05.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Children()(*i03ddf9ed2e8d64e93fc44a7258e3589c38d6f455c1358c4d75d1614935a98e46.ChildrenRequestBuilder) {
    return i03ddf9ed2e8d64e93fc44a7258e3589c38d6f455c1358c4d75d1614935a98e46.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ChildrenById(id string)(*ie2d9fc455cd20ad77bcd5ffeb7d799dfb360cc67d9482d6709d02097def6a274.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ie2d9fc455cd20ad77bcd5ffeb7d799dfb360cc67d9482d6709d02097def6a274.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewRootRequestBuilderInternal instantiates a new RootRequestBuilder and sets the default values.
func NewRootRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RootRequestBuilder) {
    m := &RootRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/drives/{drive%2Did}/root{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewRootRequestBuilder instantiates a new RootRequestBuilder and sets the default values.
func NewRootRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RootRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRootRequestBuilderInternal(urlParams, requestAdapter)
}
// Content provides operations to manage the media for the user entity.
func (m *RootRequestBuilder) Content()(*i1009a2ff258bf0b86c5aa9b40c49762ddd855f338f1947b3a8f825c3d8524721.ContentRequestBuilder) {
    return i1009a2ff258bf0b86c5aa9b40c49762ddd855f338f1947b3a8f825c3d8524721.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *RootRequestBuilder) Copy()(*i72d8e675471adc6caea657caad244220759b7b1d1af78e568f69f1b0488ae9c2.CopyRequestBuilder) {
    return i72d8e675471adc6caea657caad244220759b7b1d1af78e568f69f1b0488ae9c2.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property root for me
func (m *RootRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *RootRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation retrieve the metadata for a driveItem in a drive by file system path or ID.`item-id` is the ID of a driveItem. It may also be the unique ID of a SharePoint list item.
func (m *RootRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *RootRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *RootRequestBuilder) CreateLink()(*i7e2854ef00272a1811fdcaeaf5c86f99d0ec13adf274e99eac8180b08d5f8ae1.CreateLinkRequestBuilder) {
    return i7e2854ef00272a1811fdcaeaf5c86f99d0ec13adf274e99eac8180b08d5f8ae1.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property root in me
func (m *RootRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *RootRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *RootRequestBuilder) CreateUploadSession()(*if705254df95007eb4d749c340e906c098a5584daa3d253a27dad196a9e68aea6.CreateUploadSessionRequestBuilder) {
    return if705254df95007eb4d749c340e906c098a5584daa3d253a27dad196a9e68aea6.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property root for me
func (m *RootRequestBuilder) Delete(ctx context.Context, requestConfiguration *RootRequestBuilderDeleteRequestConfiguration)(error) {
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
func (m *RootRequestBuilder) Delta()(*idcdd5295fe11eb0b451e619b74b2ec775e69aa4c0031109a4ad4ab63a2d346c2.DeltaRequestBuilder) {
    return idcdd5295fe11eb0b451e619b74b2ec775e69aa4c0031109a4ad4ab63a2d346c2.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *RootRequestBuilder) DeltaWithToken(token *string)(*i71cf1b728ee8302f69703a7e6a0e0731e7b4dbecd838960e2ffda2156d16db08.DeltaWithTokenRequestBuilder) {
    return i71cf1b728ee8302f69703a7e6a0e0731e7b4dbecd838960e2ffda2156d16db08.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *RootRequestBuilder) Follow()(*id1bc6c7776430f10716cbe5e70e0fffffd8064c286891200f2a817b8e144fe9e.FollowRequestBuilder) {
    return id1bc6c7776430f10716cbe5e70e0fffffd8064c286891200f2a817b8e144fe9e.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get retrieve the metadata for a driveItem in a drive by file system path or ID.`item-id` is the ID of a driveItem. It may also be the unique ID of a SharePoint list item.
func (m *RootRequestBuilder) Get(ctx context.Context, requestConfiguration *RootRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
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
func (m *RootRequestBuilder) GetActivitiesByInterval()(*i7b2a9d9e2b194b78db259bec4bcb10cb8c8747663a268008a135977ab274f5ba.GetActivitiesByIntervalRequestBuilder) {
    return i7b2a9d9e2b194b78db259bec4bcb10cb8c8747663a268008a135977ab274f5ba.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *RootRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i79c6c3c119b0fd26cac016493adf8593a63ae4fab341dd8853915459798f15db.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i79c6c3c119b0fd26cac016493adf8593a63ae4fab341dd8853915459798f15db.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *RootRequestBuilder) Invite()(*ib327b2142df568febe244d5e9882de94f5049f86c795d0adc6dc9a062df1985d.InviteRequestBuilder) {
    return ib327b2142df568febe244d5e9882de94f5049f86c795d0adc6dc9a062df1985d.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ListItem()(*i404bc6886d1c9123a28011fced42055517ed97d7f967f514fb167d219b2c0db0.ListItemRequestBuilder) {
    return i404bc6886d1c9123a28011fced42055517ed97d7f967f514fb167d219b2c0db0.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property root in me
func (m *RootRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, requestConfiguration *RootRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DriveItemable, error) {
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
func (m *RootRequestBuilder) Permissions()(*i80b1e91bf7b013cdf8dbdf9b8d8771f6a1b6e697578eefaecea07952add6388c.PermissionsRequestBuilder) {
    return i80b1e91bf7b013cdf8dbdf9b8d8771f6a1b6e697578eefaecea07952add6388c.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) PermissionsById(id string)(*ieeb69da49c1200fd91a53dded9fdde76e190b3cb09517e3aab7c813ac910de6d.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return ieeb69da49c1200fd91a53dded9fdde76e190b3cb09517e3aab7c813ac910de6d.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *RootRequestBuilder) Preview()(*i27c4eb9aab96ee3918e3877e6df557c98aee08cbeebedc02e45682b3af84cf0a.PreviewRequestBuilder) {
    return i27c4eb9aab96ee3918e3877e6df557c98aee08cbeebedc02e45682b3af84cf0a.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *RootRequestBuilder) Restore()(*i0033efa979ecc8ac80fdb29d960015d241a7cd255681fd7cbf71d410f9604ad4.RestoreRequestBuilder) {
    return i0033efa979ecc8ac80fdb29d960015d241a7cd255681fd7cbf71d410f9604ad4.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *RootRequestBuilder) SearchWithQ(q *string)(*ief1f6f4bf8be57145cbf8d0f439184c8d23abbbcc4b38699c6c809db84be6beb.SearchWithQRequestBuilder) {
    return ief1f6f4bf8be57145cbf8d0f439184c8d23abbbcc4b38699c6c809db84be6beb.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Subscriptions()(*ifce4a3fdda3e56bdf741ded90a576a9c0e9f6480576f685c4c8e5eceb911a516.SubscriptionsRequestBuilder) {
    return ifce4a3fdda3e56bdf741ded90a576a9c0e9f6480576f685c4c8e5eceb911a516.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) SubscriptionsById(id string)(*i689f37510f9fad8e0be6242df80f8b79fff25b25c310cfaf55e7b93c51c780af.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return i689f37510f9fad8e0be6242df80f8b79fff25b25c310cfaf55e7b93c51c780af.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Thumbnails()(*ief5ea07106b57464e0d9313bab8da87f6427ea4ce6317fef25f88f0475243573.ThumbnailsRequestBuilder) {
    return ief5ea07106b57464e0d9313bab8da87f6427ea4ce6317fef25f88f0475243573.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ThumbnailsById(id string)(*if641b83445d13c546e420bdc7e9b23cb86e18512438d6cde4308543c0d7c8cb4.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return if641b83445d13c546e420bdc7e9b23cb86e18512438d6cde4308543c0d7c8cb4.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *RootRequestBuilder) Unfollow()(*ia80c39dd409571ac746317206288180c5750a50b838f9b4840df1960308280e9.UnfollowRequestBuilder) {
    return ia80c39dd409571ac746317206288180c5750a50b838f9b4840df1960308280e9.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *RootRequestBuilder) ValidatePermission()(*ic4f6f4ea374f94b486e173595e0a59f167b50e3195232b90a72d9d013abfa7d2.ValidatePermissionRequestBuilder) {
    return ic4f6f4ea374f94b486e173595e0a59f167b50e3195232b90a72d9d013abfa7d2.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Versions()(*i032fb091fa50ded67bd8413cc227e15a89212f22579d11ba5fcda786303b3d88.VersionsRequestBuilder) {
    return i032fb091fa50ded67bd8413cc227e15a89212f22579d11ba5fcda786303b3d88.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) VersionsById(id string)(*ia16000547e3bbe576a5a946472bea8c98307843b62bbdb98e519f2ae8166f2c2.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return ia16000547e3bbe576a5a946472bea8c98307843b62bbdb98e519f2ae8166f2c2.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
