package root

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i089dbefd5043e5207e6e2000a1bf061c66b399725298944660bf5cd253b7264f "github.com/microsoftgraph/msgraph-sdk-go/drive/root/listitem"
    i0a99a8fb5d7b1c4484fa76b06399453c096f3cf9cfeabf8424f11f92fd5a2c44 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/createuploadsession"
    i0e77718e2d8ed0fb79e161140808b50d563606786049bdea23582bd151ec244a "github.com/microsoftgraph/msgraph-sdk-go/drive/root/createlink"
    i36527fa4119014183fa4364bccc2e86f46102925dbb6d147430dbc9f5b435d7d "github.com/microsoftgraph/msgraph-sdk-go/drive/root/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    i3a9efcafad0a4d278675735198209554258a87b732c52719994c901d8fed19c1 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/deltawithtoken"
    i6b4e0044e217f76a6f04af27156aad54c99ecb8702fa8501cf102dcf06a6ef21 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/invite"
    i70e5eddc8eb57c1b8081c8f33ff2e6347843570fcedd75546786b72ab61917a5 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/checkout"
    i7411a7f72e668a7ff18a176006bed801955df8a77e55f0225b8480fca0ed334f "github.com/microsoftgraph/msgraph-sdk-go/drive/root/restore"
    i7ab3c96c86b3b38921e0b8f81dcd5b872cc990c5cf8809f28bea1fb4dcf1558e "github.com/microsoftgraph/msgraph-sdk-go/drive/root/permissions"
    i87e0f9e6aaa1c9caa4c22be1bb55dea82446b84482de27f4512e49f0a4483fed "github.com/microsoftgraph/msgraph-sdk-go/drive/root/thumbnails"
    i8bcf5f285bb3bc5ee6cab0bfec528bd18ccc1540ba888eadeb410b0c568c7c46 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/follow"
    i949a469670415d616140fbcd2d21fb2698aa23a0aa4925da17b5adc3e9bb0b66 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/preview"
    i96d93f3c07afb20ba2c213bfcc5ee660bdd3eb0ba3d5da3a5ccc49f9c680afe2 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/unfollow"
    i9cb9cab89d4c3e68cb0d46e2040e4006c8ff6c3dcbbb3212185458679195d5fc "github.com/microsoftgraph/msgraph-sdk-go/drive/root/delta"
    i9fc974917f60b393695b0f954f0bf74a22c0e71c8f699ae2319dba8c253fbcca "github.com/microsoftgraph/msgraph-sdk-go/drive/root/getactivitiesbyinterval"
    ia29571508baf63a873910b5dae4f5a73689ef617c6a0a050e827a3e260b72e75 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/versions"
    ia929efb0838f1045c0b71fc7e927522ade38d9a38eaf5c4a09af5e7b841887c6 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/children"
    icf65e813432b1534e2e77beccb7ec7a8e1b99ec53da6c4cec46646028f46ba24 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/validatepermission"
    id0ce06da0b25054b4dbfde94257d62b553c93d1753f5b9cb07ff4be8788b8b85 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/copy"
    id724096a8c3b90aa72504dcd2761cc99b632bad9d405f048a071a5f697d2195f "github.com/microsoftgraph/msgraph-sdk-go/drive/root/searchwithq"
    idb5c2e560e3daac4a3ed8f616c4dd2ffd7defc33875f4bd870963f083b1c03a8 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/content"
    idec39c75646fe0f514f66c621871a98ad8f286b849dea224ba584455e0d2e34e "github.com/microsoftgraph/msgraph-sdk-go/drive/root/analytics"
    ie8ffa2b378f01cc070adb481bd1045a5f9af75b0337e47c56ebcfe149d370bb3 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/checkin"
    if9e05debf18ae7a89d29e3f92e3d20b4bd85e1de8b0abc6c5d1464f8a2df79ae "github.com/microsoftgraph/msgraph-sdk-go/drive/root/subscriptions"
    i451951d9cef8df3c8b2a8be8726dc1285485859b8f467b1880fe9484d1f86b30 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/thumbnails/item"
    i56826de7fe9540a3d79c02bb8e18e83cbe6733df09aac1d6cf30342c155743a3 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/children/item"
    i5a7c91cef86c9d63b5934d1ca4411013dd41d7a5772f7cc81ff281eff124ef60 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/versions/item"
    i83a35aa1860d0b10df007e13131be0541e05e6b90b022f9bdfbeb8e94bcc3562 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/permissions/item"
    iabb32685aec6f73966e56b90a18908830f75cc93fae47ddcb640f1ffc9a3be02 "github.com/microsoftgraph/msgraph-sdk-go/drive/root/subscriptions/item"
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
func (m *RootRequestBuilder) Analytics()(*idec39c75646fe0f514f66c621871a98ad8f286b849dea224ba584455e0d2e34e.AnalyticsRequestBuilder) {
    return idec39c75646fe0f514f66c621871a98ad8f286b849dea224ba584455e0d2e34e.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *RootRequestBuilder) Checkin()(*ie8ffa2b378f01cc070adb481bd1045a5f9af75b0337e47c56ebcfe149d370bb3.CheckinRequestBuilder) {
    return ie8ffa2b378f01cc070adb481bd1045a5f9af75b0337e47c56ebcfe149d370bb3.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *RootRequestBuilder) Checkout()(*i70e5eddc8eb57c1b8081c8f33ff2e6347843570fcedd75546786b72ab61917a5.CheckoutRequestBuilder) {
    return i70e5eddc8eb57c1b8081c8f33ff2e6347843570fcedd75546786b72ab61917a5.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Children()(*ia929efb0838f1045c0b71fc7e927522ade38d9a38eaf5c4a09af5e7b841887c6.ChildrenRequestBuilder) {
    return ia929efb0838f1045c0b71fc7e927522ade38d9a38eaf5c4a09af5e7b841887c6.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ChildrenById(id string)(*i56826de7fe9540a3d79c02bb8e18e83cbe6733df09aac1d6cf30342c155743a3.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i56826de7fe9540a3d79c02bb8e18e83cbe6733df09aac1d6cf30342c155743a3.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewRootRequestBuilderInternal instantiates a new RootRequestBuilder and sets the default values.
func NewRootRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RootRequestBuilder) {
    m := &RootRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drive/root{?%24select,%24expand}";
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
func (m *RootRequestBuilder) Content()(*idb5c2e560e3daac4a3ed8f616c4dd2ffd7defc33875f4bd870963f083b1c03a8.ContentRequestBuilder) {
    return idb5c2e560e3daac4a3ed8f616c4dd2ffd7defc33875f4bd870963f083b1c03a8.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *RootRequestBuilder) Copy()(*id0ce06da0b25054b4dbfde94257d62b553c93d1753f5b9cb07ff4be8788b8b85.CopyRequestBuilder) {
    return id0ce06da0b25054b4dbfde94257d62b553c93d1753f5b9cb07ff4be8788b8b85.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property root for drive
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
func (m *RootRequestBuilder) CreateLink()(*i0e77718e2d8ed0fb79e161140808b50d563606786049bdea23582bd151ec244a.CreateLinkRequestBuilder) {
    return i0e77718e2d8ed0fb79e161140808b50d563606786049bdea23582bd151ec244a.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property root in drive
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
func (m *RootRequestBuilder) CreateUploadSession()(*i0a99a8fb5d7b1c4484fa76b06399453c096f3cf9cfeabf8424f11f92fd5a2c44.CreateUploadSessionRequestBuilder) {
    return i0a99a8fb5d7b1c4484fa76b06399453c096f3cf9cfeabf8424f11f92fd5a2c44.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property root for drive
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
func (m *RootRequestBuilder) Delta()(*i9cb9cab89d4c3e68cb0d46e2040e4006c8ff6c3dcbbb3212185458679195d5fc.DeltaRequestBuilder) {
    return i9cb9cab89d4c3e68cb0d46e2040e4006c8ff6c3dcbbb3212185458679195d5fc.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *RootRequestBuilder) DeltaWithToken(token *string)(*i3a9efcafad0a4d278675735198209554258a87b732c52719994c901d8fed19c1.DeltaWithTokenRequestBuilder) {
    return i3a9efcafad0a4d278675735198209554258a87b732c52719994c901d8fed19c1.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *RootRequestBuilder) Follow()(*i8bcf5f285bb3bc5ee6cab0bfec528bd18ccc1540ba888eadeb410b0c568c7c46.FollowRequestBuilder) {
    return i8bcf5f285bb3bc5ee6cab0bfec528bd18ccc1540ba888eadeb410b0c568c7c46.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *RootRequestBuilder) GetActivitiesByInterval()(*i9fc974917f60b393695b0f954f0bf74a22c0e71c8f699ae2319dba8c253fbcca.GetActivitiesByIntervalRequestBuilder) {
    return i9fc974917f60b393695b0f954f0bf74a22c0e71c8f699ae2319dba8c253fbcca.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *RootRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i36527fa4119014183fa4364bccc2e86f46102925dbb6d147430dbc9f5b435d7d.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i36527fa4119014183fa4364bccc2e86f46102925dbb6d147430dbc9f5b435d7d.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *RootRequestBuilder) Invite()(*i6b4e0044e217f76a6f04af27156aad54c99ecb8702fa8501cf102dcf06a6ef21.InviteRequestBuilder) {
    return i6b4e0044e217f76a6f04af27156aad54c99ecb8702fa8501cf102dcf06a6ef21.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ListItem()(*i089dbefd5043e5207e6e2000a1bf061c66b399725298944660bf5cd253b7264f.ListItemRequestBuilder) {
    return i089dbefd5043e5207e6e2000a1bf061c66b399725298944660bf5cd253b7264f.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property root in drive
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
func (m *RootRequestBuilder) Permissions()(*i7ab3c96c86b3b38921e0b8f81dcd5b872cc990c5cf8809f28bea1fb4dcf1558e.PermissionsRequestBuilder) {
    return i7ab3c96c86b3b38921e0b8f81dcd5b872cc990c5cf8809f28bea1fb4dcf1558e.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) PermissionsById(id string)(*i83a35aa1860d0b10df007e13131be0541e05e6b90b022f9bdfbeb8e94bcc3562.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return i83a35aa1860d0b10df007e13131be0541e05e6b90b022f9bdfbeb8e94bcc3562.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *RootRequestBuilder) Preview()(*i949a469670415d616140fbcd2d21fb2698aa23a0aa4925da17b5adc3e9bb0b66.PreviewRequestBuilder) {
    return i949a469670415d616140fbcd2d21fb2698aa23a0aa4925da17b5adc3e9bb0b66.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *RootRequestBuilder) Restore()(*i7411a7f72e668a7ff18a176006bed801955df8a77e55f0225b8480fca0ed334f.RestoreRequestBuilder) {
    return i7411a7f72e668a7ff18a176006bed801955df8a77e55f0225b8480fca0ed334f.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *RootRequestBuilder) SearchWithQ(q *string)(*id724096a8c3b90aa72504dcd2761cc99b632bad9d405f048a071a5f697d2195f.SearchWithQRequestBuilder) {
    return id724096a8c3b90aa72504dcd2761cc99b632bad9d405f048a071a5f697d2195f.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Subscriptions()(*if9e05debf18ae7a89d29e3f92e3d20b4bd85e1de8b0abc6c5d1464f8a2df79ae.SubscriptionsRequestBuilder) {
    return if9e05debf18ae7a89d29e3f92e3d20b4bd85e1de8b0abc6c5d1464f8a2df79ae.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) SubscriptionsById(id string)(*iabb32685aec6f73966e56b90a18908830f75cc93fae47ddcb640f1ffc9a3be02.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return iabb32685aec6f73966e56b90a18908830f75cc93fae47ddcb640f1ffc9a3be02.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Thumbnails()(*i87e0f9e6aaa1c9caa4c22be1bb55dea82446b84482de27f4512e49f0a4483fed.ThumbnailsRequestBuilder) {
    return i87e0f9e6aaa1c9caa4c22be1bb55dea82446b84482de27f4512e49f0a4483fed.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) ThumbnailsById(id string)(*i451951d9cef8df3c8b2a8be8726dc1285485859b8f467b1880fe9484d1f86b30.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return i451951d9cef8df3c8b2a8be8726dc1285485859b8f467b1880fe9484d1f86b30.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *RootRequestBuilder) Unfollow()(*i96d93f3c07afb20ba2c213bfcc5ee660bdd3eb0ba3d5da3a5ccc49f9c680afe2.UnfollowRequestBuilder) {
    return i96d93f3c07afb20ba2c213bfcc5ee660bdd3eb0ba3d5da3a5ccc49f9c680afe2.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *RootRequestBuilder) ValidatePermission()(*icf65e813432b1534e2e77beccb7ec7a8e1b99ec53da6c4cec46646028f46ba24.ValidatePermissionRequestBuilder) {
    return icf65e813432b1534e2e77beccb7ec7a8e1b99ec53da6c4cec46646028f46ba24.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) Versions()(*ia29571508baf63a873910b5dae4f5a73689ef617c6a0a050e827a3e260b72e75.VersionsRequestBuilder) {
    return ia29571508baf63a873910b5dae4f5a73689ef617c6a0a050e827a3e260b72e75.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *RootRequestBuilder) VersionsById(id string)(*i5a7c91cef86c9d63b5934d1ca4411013dd41d7a5772f7cc81ff281eff124ef60.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return i5a7c91cef86c9d63b5934d1ca4411013dd41d7a5772f7cc81ff281eff124ef60.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
