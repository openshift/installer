package root

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i02452b02d09bf832000dca9e33aca96ccaf26f7b61841c0c6c476a7d12f7eca7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/checkin"
    i025b69ef2d2ddbfb7e9807a38f690bdf4cd63d6712a67206da809a0cdfc42189 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/thumbnails"
    i0dae1ce8f5993df8d5b23e2ff35f73b9208a960111e9d84bb59ce507c2a43b23 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/checkout"
    i118ae752cb14a143b3827b518f2d208f8022e14c050d29068b8e5fb83442f32d "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/delta"
    i135e884c820565705b38a640ffdf31280e134de90528d9785ec0e75376321a61 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/analytics"
    i2ea67eb04bbbfa82d90dd88dddc4c42ccdbb62c2e552c27c0091b1608040bcf8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/listitem"
    i3d0e4f8aa6ec904f8c27942ce3497aadbb1cd99228b4c60d2a161cc1cd3859cd "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/children"
    i55fa05beaa401ce538b7b8e5032419e0945c3be38d72f54e68c69874fa1d1cbb "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/createlink"
    i593e2f9c5ea33ff66149a808910d4ae3e5ae255503b54f3742437caf64e884f8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/copy"
    i5bbb0ccbc4c3843f9a789010d1b81a6c7dc71a17cc52f5e0d68f386d2057dbc9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/subscriptions"
    i766af667e01e173e5625ba5550a8532b78b66aeb04193536e22716027e5989fe "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/preview"
    i7b6ee7d5b6663b375ed921bc6ba9ad187a415e8a469daf7f11aa380553d07065 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/versions"
    i8acc4d3596ccc8e653935bec6d055cc4adf3955ec6cd7df214397d0354d0dcd3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/invite"
    i8d536d7b0adff285ec8d565fcfda8f8b51f664ac66fcc1466ee5dd701e61677a "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i90812a66c9855e30e4b6cca138d506affed6dd51888715cd07f5e980a860bebc "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/follow"
    iaa1d84731a4b8785f6dacbbb5439c2508806285148dd6a0c01c1be4f2b099002 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/content"
    ib1efd02c4c2751e576dd0de1c915b8e6db635eae24dd01117248f2e015d27ec6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/permissions"
    ib9006619a5801080c2f277b9a7f4c490ddb97e0ccf598880d9a1352d1c7761ae "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/validatepermission"
    ibf0f337b4c15da08148eb948869b44de35c3789dee48fbc9aea5721c7d98fa8a "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/createuploadsession"
    idf0b3a3e5899ec908b59474e7501b458c62e53de0501bb85594046515f15f2ff "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/unfollow"
    ie3fb0038bbda0f8f65b4e0c77d4cd4306b20cf2f887b782c207b18256226b5eb "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/getactivitiesbyinterval"
    ie5b7fc64a3bffb4b52a76cf59bee4382bfa061320905c4bfedd0c17eb9aa3984 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/searchwithq"
    if916bbe2a3dc128eaa8201e2fc9de95a3ed150dd25f862ba6637e16faf747d86 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/restore"
    ifb73e41cc51e69e24e449a47f621bf9ba0b65e48d6f61e225f736a939d09b985 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/deltawithtoken"
    i67ec6eaee9045e07acc845ccdf57c93da38efc06e18c1bb8dfd55e85d7b6320d "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/subscriptions/item"
    i8dac7fcd490022c229aca597d1ef2174bb968ffc2cd99e554ec494acf608d88c "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/thumbnails/item"
    iccb7f2d5d33cfb54ca0c62427ff964d27d77254ebd63aa68529b5eed5a42266d "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/children/item"
    id0c42f467a4df86619438eccced198bfbc3edae46d9a7697c9794a1f605d4d11 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/permissions/item"
    id5e75288d7196634c44036e792f83a221d2e22ef830f4c48a8c0e09650983d14 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root/versions/item"
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
func (m *RootRequestBuilder) Analytics()(*i135e884c820565705b38a640ffdf31280e134de90528d9785ec0e75376321a61.AnalyticsRequestBuilder) {
    return i135e884c820565705b38a640ffdf31280e134de90528d9785ec0e75376321a61.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *RootRequestBuilder) Checkin()(*i02452b02d09bf832000dca9e33aca96ccaf26f7b61841c0c6c476a7d12f7eca7.CheckinRequestBuilder) {
    return i02452b02d09bf832000dca9e33aca96ccaf26f7b61841c0c6c476a7d12f7eca7.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *RootRequestBuilder) Checkout()(*i0dae1ce8f5993df8d5b23e2ff35f73b9208a960111e9d84bb59ce507c2a43b23.CheckoutRequestBuilder) {
    return i0dae1ce8f5993df8d5b23e2ff35f73b9208a960111e9d84bb59ce507c2a43b23.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Children()(*i3d0e4f8aa6ec904f8c27942ce3497aadbb1cd99228b4c60d2a161cc1cd3859cd.ChildrenRequestBuilder) {
    return i3d0e4f8aa6ec904f8c27942ce3497aadbb1cd99228b4c60d2a161cc1cd3859cd.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ChildrenById(id string)(*iccb7f2d5d33cfb54ca0c62427ff964d27d77254ebd63aa68529b5eed5a42266d.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return iccb7f2d5d33cfb54ca0c62427ff964d27d77254ebd63aa68529b5eed5a42266d.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewRootRequestBuilderInternal instantiates a new RootRequestBuilder and sets the default values.
func NewRootRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RootRequestBuilder) {
    m := &RootRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/drives/{drive%2Did}/root{?%24select,%24expand}";
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
func (m *RootRequestBuilder) Content()(*iaa1d84731a4b8785f6dacbbb5439c2508806285148dd6a0c01c1be4f2b099002.ContentRequestBuilder) {
    return iaa1d84731a4b8785f6dacbbb5439c2508806285148dd6a0c01c1be4f2b099002.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *RootRequestBuilder) Copy()(*i593e2f9c5ea33ff66149a808910d4ae3e5ae255503b54f3742437caf64e884f8.CopyRequestBuilder) {
    return i593e2f9c5ea33ff66149a808910d4ae3e5ae255503b54f3742437caf64e884f8.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property root for users
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
func (m *RootRequestBuilder) CreateLink()(*i55fa05beaa401ce538b7b8e5032419e0945c3be38d72f54e68c69874fa1d1cbb.CreateLinkRequestBuilder) {
    return i55fa05beaa401ce538b7b8e5032419e0945c3be38d72f54e68c69874fa1d1cbb.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property root in users
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
func (m *RootRequestBuilder) CreateUploadSession()(*ibf0f337b4c15da08148eb948869b44de35c3789dee48fbc9aea5721c7d98fa8a.CreateUploadSessionRequestBuilder) {
    return ibf0f337b4c15da08148eb948869b44de35c3789dee48fbc9aea5721c7d98fa8a.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property root for users
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
func (m *RootRequestBuilder) Delta()(*i118ae752cb14a143b3827b518f2d208f8022e14c050d29068b8e5fb83442f32d.DeltaRequestBuilder) {
    return i118ae752cb14a143b3827b518f2d208f8022e14c050d29068b8e5fb83442f32d.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *RootRequestBuilder) DeltaWithToken(token *string)(*ifb73e41cc51e69e24e449a47f621bf9ba0b65e48d6f61e225f736a939d09b985.DeltaWithTokenRequestBuilder) {
    return ifb73e41cc51e69e24e449a47f621bf9ba0b65e48d6f61e225f736a939d09b985.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *RootRequestBuilder) Follow()(*i90812a66c9855e30e4b6cca138d506affed6dd51888715cd07f5e980a860bebc.FollowRequestBuilder) {
    return i90812a66c9855e30e4b6cca138d506affed6dd51888715cd07f5e980a860bebc.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *RootRequestBuilder) GetActivitiesByInterval()(*ie3fb0038bbda0f8f65b4e0c77d4cd4306b20cf2f887b782c207b18256226b5eb.GetActivitiesByIntervalRequestBuilder) {
    return ie3fb0038bbda0f8f65b4e0c77d4cd4306b20cf2f887b782c207b18256226b5eb.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *RootRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i8d536d7b0adff285ec8d565fcfda8f8b51f664ac66fcc1466ee5dd701e61677a.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i8d536d7b0adff285ec8d565fcfda8f8b51f664ac66fcc1466ee5dd701e61677a.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *RootRequestBuilder) Invite()(*i8acc4d3596ccc8e653935bec6d055cc4adf3955ec6cd7df214397d0354d0dcd3.InviteRequestBuilder) {
    return i8acc4d3596ccc8e653935bec6d055cc4adf3955ec6cd7df214397d0354d0dcd3.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ListItem()(*i2ea67eb04bbbfa82d90dd88dddc4c42ccdbb62c2e552c27c0091b1608040bcf8.ListItemRequestBuilder) {
    return i2ea67eb04bbbfa82d90dd88dddc4c42ccdbb62c2e552c27c0091b1608040bcf8.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property root in users
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
func (m *RootRequestBuilder) Permissions()(*ib1efd02c4c2751e576dd0de1c915b8e6db635eae24dd01117248f2e015d27ec6.PermissionsRequestBuilder) {
    return ib1efd02c4c2751e576dd0de1c915b8e6db635eae24dd01117248f2e015d27ec6.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) PermissionsById(id string)(*id0c42f467a4df86619438eccced198bfbc3edae46d9a7697c9794a1f605d4d11.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return id0c42f467a4df86619438eccced198bfbc3edae46d9a7697c9794a1f605d4d11.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *RootRequestBuilder) Preview()(*i766af667e01e173e5625ba5550a8532b78b66aeb04193536e22716027e5989fe.PreviewRequestBuilder) {
    return i766af667e01e173e5625ba5550a8532b78b66aeb04193536e22716027e5989fe.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *RootRequestBuilder) Restore()(*if916bbe2a3dc128eaa8201e2fc9de95a3ed150dd25f862ba6637e16faf747d86.RestoreRequestBuilder) {
    return if916bbe2a3dc128eaa8201e2fc9de95a3ed150dd25f862ba6637e16faf747d86.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *RootRequestBuilder) SearchWithQ(q *string)(*ie5b7fc64a3bffb4b52a76cf59bee4382bfa061320905c4bfedd0c17eb9aa3984.SearchWithQRequestBuilder) {
    return ie5b7fc64a3bffb4b52a76cf59bee4382bfa061320905c4bfedd0c17eb9aa3984.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Subscriptions()(*i5bbb0ccbc4c3843f9a789010d1b81a6c7dc71a17cc52f5e0d68f386d2057dbc9.SubscriptionsRequestBuilder) {
    return i5bbb0ccbc4c3843f9a789010d1b81a6c7dc71a17cc52f5e0d68f386d2057dbc9.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) SubscriptionsById(id string)(*i67ec6eaee9045e07acc845ccdf57c93da38efc06e18c1bb8dfd55e85d7b6320d.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return i67ec6eaee9045e07acc845ccdf57c93da38efc06e18c1bb8dfd55e85d7b6320d.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Thumbnails()(*i025b69ef2d2ddbfb7e9807a38f690bdf4cd63d6712a67206da809a0cdfc42189.ThumbnailsRequestBuilder) {
    return i025b69ef2d2ddbfb7e9807a38f690bdf4cd63d6712a67206da809a0cdfc42189.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ThumbnailsById(id string)(*i8dac7fcd490022c229aca597d1ef2174bb968ffc2cd99e554ec494acf608d88c.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return i8dac7fcd490022c229aca597d1ef2174bb968ffc2cd99e554ec494acf608d88c.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *RootRequestBuilder) Unfollow()(*idf0b3a3e5899ec908b59474e7501b458c62e53de0501bb85594046515f15f2ff.UnfollowRequestBuilder) {
    return idf0b3a3e5899ec908b59474e7501b458c62e53de0501bb85594046515f15f2ff.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *RootRequestBuilder) ValidatePermission()(*ib9006619a5801080c2f277b9a7f4c490ddb97e0ccf598880d9a1352d1c7761ae.ValidatePermissionRequestBuilder) {
    return ib9006619a5801080c2f277b9a7f4c490ddb97e0ccf598880d9a1352d1c7761ae.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Versions()(*i7b6ee7d5b6663b375ed921bc6ba9ad187a415e8a469daf7f11aa380553d07065.VersionsRequestBuilder) {
    return i7b6ee7d5b6663b375ed921bc6ba9ad187a415e8a469daf7f11aa380553d07065.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) VersionsById(id string)(*id5e75288d7196634c44036e792f83a221d2e22ef830f4c48a8c0e09650983d14.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return id5e75288d7196634c44036e792f83a221d2e22ef830f4c48a8c0e09650983d14.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
