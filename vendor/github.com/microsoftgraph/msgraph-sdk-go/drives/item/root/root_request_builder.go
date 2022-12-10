package root

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0c83e896ce3f14167434af67b0c833c0001d446b958794d29b20f00cefdba1ab "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/unfollow"
    i10fb790bbe09603533347e125df89fb9a33ef78b650cc7274d8bc7cbf78b1e56 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/invite"
    i1c05393909ced8c970ac8182fc4a6c5cc419883d0d1475ce0773c5cf783ba989 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/listitem"
    i353c71cf4e7f00f79bca7c43c481c2b7fccfaf2bc763be7f85e457bcdb2a832b "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/content"
    i35ec8d39bca25c3225ac281fac7242810cb642a0a53e14a875e33c70b99c171e "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/analytics"
    i3c24e0d0e594021b9a5040d170802c6fe01c9e39c6fb7c59681d77d731fe9b4c "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/createuploadsession"
    i3e0c800f29f310a253f02b35214628e7f3e1d58029de18307deadf71c019a972 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/getactivitiesbyinterval"
    i3ed0881d10cd2f4686e9a20ef9e3df745f32c4a5deaf241497268510a8dfb38f "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/createlink"
    i4fe456dbe4b9db0c319a13c072df6cd5de3fc0e9fdc876f33cdaecb702546d10 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i6332f0883e6ab0e66763410aa3cab286bc78b7b61f88fefedb1be3fb930832df "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/thumbnails"
    i6a155f272c5391daf639268076f7e78fb66358ef654a4f57db99ce00a32f3c11 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/children"
    i7d3f6d8421e7656fb91d2a5f58ce09709f147443a8c14bc065dc764c9b20e94a "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/follow"
    i8586e46c09b38f90b06f86eb42fd10319ce6acb9451aefe28f39853cfb5719b9 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/versions"
    i9d7e3dff92ecf83dceea8302ab5d7d9705c203886940133f82426f3f2a2eab3d "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/checkin"
    ia016f7159d32c9474ca21bb83155879bb92111457cd5a4a8da952cae87a9dd57 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/restore"
    ia05f97a343c7dbf32d5d42d3b714d8d4e180ebdebfee2e83c7a95ec1ed6e9e3e "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/permissions"
    ia8efe22714d926b99d7a1272ac32f8a7fc99f84889f31fc8917091d21f9413c6 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/delta"
    ia8f4a48393a3e5d61e77b4702b5c64c47a37ba6ed0a2853d3fee9568c80bcc1a "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/checkout"
    iabea5c49e48376850d4c62da6ee59d3cc84c4c9d26139e0ec1bfde329c0e3331 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/copy"
    ib668b4d5a43288034d2e7f9226cf4069e9c6e36c554f12e93013a1be6a8768ab "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/validatepermission"
    ic1c29d96fa432fa3653e15cdf1d85ce8c5c9934d509c7db7c45b20df335b2c83 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/deltawithtoken"
    ic1deb0ea5c033d7b6ab9e0618caeb86aae80c379b0f2803f804018ffeaa48d9c "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/preview"
    ic72fcde28010dbd679eab351c33b57107adc23c1cc48e905dd109128ec3bc4db "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/subscriptions"
    id4fb6c3815ee0c1e71dcaa77620e89ca2f6be8a9767a42e4e40ad3aa464742f8 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/searchwithq"
    i162b49d374c91a27c7d63a8c7103a36e56411b6b3e84fc64ecbf5d828ae8a8af "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/versions/item"
    i3d94da0dddfc16aca3f8f4350e9b0d9d0a38bc1ffb67f012ae9cb026450ac2a1 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/subscriptions/item"
    i896e3d0df31741011967f273d16af8fb1abe2a50bd0bb2f742e0aead26932993 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/thumbnails/item"
    i8ca583e1745c18b27ce93c86c752cf4c03f78604e1970006d1886875c01639c7 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/permissions/item"
    ic7d208ecab0cc1c68790901ba2ba97554ff0a0a86b46da4f6d1dba8ab740987f "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root/children/item"
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
func (m *RootRequestBuilder) Analytics()(*i35ec8d39bca25c3225ac281fac7242810cb642a0a53e14a875e33c70b99c171e.AnalyticsRequestBuilder) {
    return i35ec8d39bca25c3225ac281fac7242810cb642a0a53e14a875e33c70b99c171e.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *RootRequestBuilder) Checkin()(*i9d7e3dff92ecf83dceea8302ab5d7d9705c203886940133f82426f3f2a2eab3d.CheckinRequestBuilder) {
    return i9d7e3dff92ecf83dceea8302ab5d7d9705c203886940133f82426f3f2a2eab3d.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *RootRequestBuilder) Checkout()(*ia8f4a48393a3e5d61e77b4702b5c64c47a37ba6ed0a2853d3fee9568c80bcc1a.CheckoutRequestBuilder) {
    return ia8f4a48393a3e5d61e77b4702b5c64c47a37ba6ed0a2853d3fee9568c80bcc1a.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Children()(*i6a155f272c5391daf639268076f7e78fb66358ef654a4f57db99ce00a32f3c11.ChildrenRequestBuilder) {
    return i6a155f272c5391daf639268076f7e78fb66358ef654a4f57db99ce00a32f3c11.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ChildrenById(id string)(*ic7d208ecab0cc1c68790901ba2ba97554ff0a0a86b46da4f6d1dba8ab740987f.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ic7d208ecab0cc1c68790901ba2ba97554ff0a0a86b46da4f6d1dba8ab740987f.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewRootRequestBuilderInternal instantiates a new RootRequestBuilder and sets the default values.
func NewRootRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RootRequestBuilder) {
    m := &RootRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drives/{drive%2Did}/root{?%24select,%24expand}";
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
// Content provides operations to manage the media for the drive entity.
func (m *RootRequestBuilder) Content()(*i353c71cf4e7f00f79bca7c43c481c2b7fccfaf2bc763be7f85e457bcdb2a832b.ContentRequestBuilder) {
    return i353c71cf4e7f00f79bca7c43c481c2b7fccfaf2bc763be7f85e457bcdb2a832b.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *RootRequestBuilder) Copy()(*iabea5c49e48376850d4c62da6ee59d3cc84c4c9d26139e0ec1bfde329c0e3331.CopyRequestBuilder) {
    return iabea5c49e48376850d4c62da6ee59d3cc84c4c9d26139e0ec1bfde329c0e3331.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property root for drives
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
func (m *RootRequestBuilder) CreateLink()(*i3ed0881d10cd2f4686e9a20ef9e3df745f32c4a5deaf241497268510a8dfb38f.CreateLinkRequestBuilder) {
    return i3ed0881d10cd2f4686e9a20ef9e3df745f32c4a5deaf241497268510a8dfb38f.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property root in drives
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
func (m *RootRequestBuilder) CreateUploadSession()(*i3c24e0d0e594021b9a5040d170802c6fe01c9e39c6fb7c59681d77d731fe9b4c.CreateUploadSessionRequestBuilder) {
    return i3c24e0d0e594021b9a5040d170802c6fe01c9e39c6fb7c59681d77d731fe9b4c.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property root for drives
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
func (m *RootRequestBuilder) Delta()(*ia8efe22714d926b99d7a1272ac32f8a7fc99f84889f31fc8917091d21f9413c6.DeltaRequestBuilder) {
    return ia8efe22714d926b99d7a1272ac32f8a7fc99f84889f31fc8917091d21f9413c6.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *RootRequestBuilder) DeltaWithToken(token *string)(*ic1c29d96fa432fa3653e15cdf1d85ce8c5c9934d509c7db7c45b20df335b2c83.DeltaWithTokenRequestBuilder) {
    return ic1c29d96fa432fa3653e15cdf1d85ce8c5c9934d509c7db7c45b20df335b2c83.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *RootRequestBuilder) Follow()(*i7d3f6d8421e7656fb91d2a5f58ce09709f147443a8c14bc065dc764c9b20e94a.FollowRequestBuilder) {
    return i7d3f6d8421e7656fb91d2a5f58ce09709f147443a8c14bc065dc764c9b20e94a.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *RootRequestBuilder) GetActivitiesByInterval()(*i3e0c800f29f310a253f02b35214628e7f3e1d58029de18307deadf71c019a972.GetActivitiesByIntervalRequestBuilder) {
    return i3e0c800f29f310a253f02b35214628e7f3e1d58029de18307deadf71c019a972.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *RootRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i4fe456dbe4b9db0c319a13c072df6cd5de3fc0e9fdc876f33cdaecb702546d10.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i4fe456dbe4b9db0c319a13c072df6cd5de3fc0e9fdc876f33cdaecb702546d10.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *RootRequestBuilder) Invite()(*i10fb790bbe09603533347e125df89fb9a33ef78b650cc7274d8bc7cbf78b1e56.InviteRequestBuilder) {
    return i10fb790bbe09603533347e125df89fb9a33ef78b650cc7274d8bc7cbf78b1e56.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ListItem()(*i1c05393909ced8c970ac8182fc4a6c5cc419883d0d1475ce0773c5cf783ba989.ListItemRequestBuilder) {
    return i1c05393909ced8c970ac8182fc4a6c5cc419883d0d1475ce0773c5cf783ba989.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property root in drives
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
func (m *RootRequestBuilder) Permissions()(*ia05f97a343c7dbf32d5d42d3b714d8d4e180ebdebfee2e83c7a95ec1ed6e9e3e.PermissionsRequestBuilder) {
    return ia05f97a343c7dbf32d5d42d3b714d8d4e180ebdebfee2e83c7a95ec1ed6e9e3e.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) PermissionsById(id string)(*i8ca583e1745c18b27ce93c86c752cf4c03f78604e1970006d1886875c01639c7.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i8ca583e1745c18b27ce93c86c752cf4c03f78604e1970006d1886875c01639c7.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *RootRequestBuilder) Preview()(*ic1deb0ea5c033d7b6ab9e0618caeb86aae80c379b0f2803f804018ffeaa48d9c.PreviewRequestBuilder) {
    return ic1deb0ea5c033d7b6ab9e0618caeb86aae80c379b0f2803f804018ffeaa48d9c.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *RootRequestBuilder) Restore()(*ia016f7159d32c9474ca21bb83155879bb92111457cd5a4a8da952cae87a9dd57.RestoreRequestBuilder) {
    return ia016f7159d32c9474ca21bb83155879bb92111457cd5a4a8da952cae87a9dd57.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *RootRequestBuilder) SearchWithQ(q *string)(*id4fb6c3815ee0c1e71dcaa77620e89ca2f6be8a9767a42e4e40ad3aa464742f8.SearchWithQRequestBuilder) {
    return id4fb6c3815ee0c1e71dcaa77620e89ca2f6be8a9767a42e4e40ad3aa464742f8.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Subscriptions()(*ic72fcde28010dbd679eab351c33b57107adc23c1cc48e905dd109128ec3bc4db.SubscriptionsRequestBuilder) {
    return ic72fcde28010dbd679eab351c33b57107adc23c1cc48e905dd109128ec3bc4db.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) SubscriptionsById(id string)(*i3d94da0dddfc16aca3f8f4350e9b0d9d0a38bc1ffb67f012ae9cb026450ac2a1.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return i3d94da0dddfc16aca3f8f4350e9b0d9d0a38bc1ffb67f012ae9cb026450ac2a1.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Thumbnails()(*i6332f0883e6ab0e66763410aa3cab286bc78b7b61f88fefedb1be3fb930832df.ThumbnailsRequestBuilder) {
    return i6332f0883e6ab0e66763410aa3cab286bc78b7b61f88fefedb1be3fb930832df.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ThumbnailsById(id string)(*i896e3d0df31741011967f273d16af8fb1abe2a50bd0bb2f742e0aead26932993.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return i896e3d0df31741011967f273d16af8fb1abe2a50bd0bb2f742e0aead26932993.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *RootRequestBuilder) Unfollow()(*i0c83e896ce3f14167434af67b0c833c0001d446b958794d29b20f00cefdba1ab.UnfollowRequestBuilder) {
    return i0c83e896ce3f14167434af67b0c833c0001d446b958794d29b20f00cefdba1ab.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *RootRequestBuilder) ValidatePermission()(*ib668b4d5a43288034d2e7f9226cf4069e9c6e36c554f12e93013a1be6a8768ab.ValidatePermissionRequestBuilder) {
    return ib668b4d5a43288034d2e7f9226cf4069e9c6e36c554f12e93013a1be6a8768ab.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Versions()(*i8586e46c09b38f90b06f86eb42fd10319ce6acb9451aefe28f39853cfb5719b9.VersionsRequestBuilder) {
    return i8586e46c09b38f90b06f86eb42fd10319ce6acb9451aefe28f39853cfb5719b9.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) VersionsById(id string)(*i162b49d374c91a27c7d63a8c7103a36e56411b6b3e84fc64ecbf5d828ae8a8af.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return i162b49d374c91a27c7d63a8c7103a36e56411b6b3e84fc64ecbf5d828ae8a8af.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
