package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0175c2914b12b61cc774e65aa639603895b1ae1642fd0d5746778bab3a505278 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/restore"
    i05d5f2d2f7d8643c0cc65d242a1e8e39efdfff6868938e24cc4d55d7dc247a81 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/thumbnails"
    i0cb44ac3ff5b85c2f1592086ca294cb612434d9f0b242efaa25ca47763aeef53 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/subscriptions"
    i1028a583cd3165679a2742667fe557347c15e838469016c886ba9642f536e3c5 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/analytics"
    i1279c7a2c940534be4779e183d6614a0624954cae0ba0fd7b18a7d74097cd266 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/createlink"
    i18a280d283ce2e40c61dc3a789783cd9b87e7d03b58573398c49a171ab7a1061 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/children"
    i48430594ddfc5f5ffb366a62303fd243e98398831f6294f94fbe58871c0afb9d "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/checkout"
    i4dacd278c290b484ea557826599ab3b681b64884fd13f4110cdbc03d72fd179e "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/unfollow"
    i52a4c1c4396bb6dfa2711ab3a6e974807f91b1191aa227cc5eaafc77c74b3b5a "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/preview"
    i5425ede49864b60be1b85cf869bde8f098a0f47f21b91eadb584348c6a59870d "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/getactivitiesbyinterval"
    i619c92bc6529cafebebfcd7f39e4ed704fee779152dbc90514101c724e7ae835 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/copy"
    i6afd67fc45a6c6ed94f8b5901f953074f83bb810cbb940ad413154c343a3a948 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/validatepermission"
    i8e060e70f61b546a5a24f8f386ceec3e29344eda06dbd47476538c8c557debac "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/checkin"
    i9bb06ffe49645f7e2e889f53f7bc1b0ac72ccb5e2106edc587893cd72d3288d2 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/getactivitiesbyintervalwithstartdatetimewithenddatetimewithinterval"
    iae2098c15020f00a005771b8a0dfdf0cc230504133ee3754b7e1f29877241492 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/searchwithq"
    ibb991d85daa24e5224252c1e7c3c5839e8297cb86bc8339eedab2a007f24a616 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/deltawithtoken"
    icc851f6701cc3e379edf8f1c8a5a20a493cce7efe46c146a987b9411a8c91b78 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/listitem"
    id49bde74d8f16edb48b5c947d716968caa13db975be8099d653271f59f39deab "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/delta"
    id5401dd754f1beaff524a09c4501c5fceb1e5bdd9e70bbe3d95723c2f329c5da "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/content"
    id7d7999927066bc35d0394f9f642e8385955712b09c4a555591918b8084870a2 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/follow"
    id9669323a3b4752f5b4c8d5f78acfd3b98dfdc33cd348c62aa7c25cc35e019f1 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/versions"
    ie33d67c99c5a2f3be2428b388d567f731532c5ddf60174b4f8c6daeff9fbfaae "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/createuploadsession"
    if1296dc61ef5c60fc9a2b02825626c2b7d61e21c6cdf5462ccaef6d25386fd23 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/invite"
    ifda175920025647e9407948753f11cc7a3fabff7fd6e426cf29edc97c408e5c2 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/permissions"
    i2ab5710457ce9e7de18564cd6ca79881ae7e842006bdf4476c908f1264282bdd "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/subscriptions/item"
    i90ba6214fcbef64f1de16cc5f1976222c710318406c262c19cfa0401c934bdd3 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/children/item"
    ib57b57d66e5436184bc3dd19450c0882627b04c9aa4d6e54d68ff2c33c761387 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/versions/item"
    ic92493f36906ee5f9721f2db0a85d2445ba0fba9bae564c8d6aaae67916102a7 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/thumbnails/item"
    id6e448f4be60986224ef6568b4a9b7d39a61886fbd386e3281d8989508828a70 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item/permissions/item"
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
func (m *DriveItemItemRequestBuilder) Analytics()(*i1028a583cd3165679a2742667fe557347c15e838469016c886ba9642f536e3c5.AnalyticsRequestBuilder) {
    return i1028a583cd3165679a2742667fe557347c15e838469016c886ba9642f536e3c5.NewAnalyticsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkin provides operations to call the checkin method.
func (m *DriveItemItemRequestBuilder) Checkin()(*i8e060e70f61b546a5a24f8f386ceec3e29344eda06dbd47476538c8c557debac.CheckinRequestBuilder) {
    return i8e060e70f61b546a5a24f8f386ceec3e29344eda06dbd47476538c8c557debac.NewCheckinRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Checkout provides operations to call the checkout method.
func (m *DriveItemItemRequestBuilder) Checkout()(*i48430594ddfc5f5ffb366a62303fd243e98398831f6294f94fbe58871c0afb9d.CheckoutRequestBuilder) {
    return i48430594ddfc5f5ffb366a62303fd243e98398831f6294f94fbe58871c0afb9d.NewCheckoutRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Children provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Children()(*i18a280d283ce2e40c61dc3a789783cd9b87e7d03b58573398c49a171ab7a1061.ChildrenRequestBuilder) {
    return i18a280d283ce2e40c61dc3a789783cd9b87e7d03b58573398c49a171ab7a1061.NewChildrenRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildrenById provides operations to manage the children property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ChildrenById(id string)(*i90ba6214fcbef64f1de16cc5f1976222c710318406c262c19cfa0401c934bdd3.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did1"] = id
    }
    return i90ba6214fcbef64f1de16cc5f1976222c710318406c262c19cfa0401c934bdd3.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemItemRequestBuilderInternal instantiates a new DriveItemItemRequestBuilder and sets the default values.
func NewDriveItemItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemItemRequestBuilder) {
    m := &DriveItemItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}{?%24select,%24expand}";
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
// Content provides operations to manage the media for the drive entity.
func (m *DriveItemItemRequestBuilder) Content()(*id5401dd754f1beaff524a09c4501c5fceb1e5bdd9e70bbe3d95723c2f329c5da.ContentRequestBuilder) {
    return id5401dd754f1beaff524a09c4501c5fceb1e5bdd9e70bbe3d95723c2f329c5da.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *DriveItemItemRequestBuilder) Copy()(*i619c92bc6529cafebebfcd7f39e4ed704fee779152dbc90514101c724e7ae835.CopyRequestBuilder) {
    return i619c92bc6529cafebebfcd7f39e4ed704fee779152dbc90514101c724e7ae835.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property items for drives
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
func (m *DriveItemItemRequestBuilder) CreateLink()(*i1279c7a2c940534be4779e183d6614a0624954cae0ba0fd7b18a7d74097cd266.CreateLinkRequestBuilder) {
    return i1279c7a2c940534be4779e183d6614a0624954cae0ba0fd7b18a7d74097cd266.NewCreateLinkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatePatchRequestInformation update the navigation property items in drives
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
func (m *DriveItemItemRequestBuilder) CreateUploadSession()(*ie33d67c99c5a2f3be2428b388d567f731532c5ddf60174b4f8c6daeff9fbfaae.CreateUploadSessionRequestBuilder) {
    return ie33d67c99c5a2f3be2428b388d567f731532c5ddf60174b4f8c6daeff9fbfaae.NewCreateUploadSessionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property items for drives
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
func (m *DriveItemItemRequestBuilder) Delta()(*id49bde74d8f16edb48b5c947d716968caa13db975be8099d653271f59f39deab.DeltaRequestBuilder) {
    return id49bde74d8f16edb48b5c947d716968caa13db975be8099d653271f59f39deab.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeltaWithToken provides operations to call the delta method.
func (m *DriveItemItemRequestBuilder) DeltaWithToken(token *string)(*ibb991d85daa24e5224252c1e7c3c5839e8297cb86bc8339eedab2a007f24a616.DeltaWithTokenRequestBuilder) {
    return ibb991d85daa24e5224252c1e7c3c5839e8297cb86bc8339eedab2a007f24a616.NewDeltaWithTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, token);
}
// Follow provides operations to call the follow method.
func (m *DriveItemItemRequestBuilder) Follow()(*id7d7999927066bc35d0394f9f642e8385955712b09c4a555591918b8084870a2.FollowRequestBuilder) {
    return id7d7999927066bc35d0394f9f642e8385955712b09c4a555591918b8084870a2.NewFollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *DriveItemItemRequestBuilder) GetActivitiesByInterval()(*i5425ede49864b60be1b85cf869bde8f098a0f47f21b91eadb584348c6a59870d.GetActivitiesByIntervalRequestBuilder) {
    return i5425ede49864b60be1b85cf869bde8f098a0f47f21b91eadb584348c6a59870d.NewGetActivitiesByIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval provides operations to call the getActivitiesByInterval method.
func (m *DriveItemItemRequestBuilder) GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithInterval(endDateTime *string, interval *string, startDateTime *string)(*i9bb06ffe49645f7e2e889f53f7bc1b0ac72ccb5e2106edc587893cd72d3288d2.GetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilder) {
    return i9bb06ffe49645f7e2e889f53f7bc1b0ac72ccb5e2106edc587893cd72d3288d2.NewGetActivitiesByIntervalWithStartDateTimeWithEndDateTimeWithIntervalRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, interval, startDateTime);
}
// Invite provides operations to call the invite method.
func (m *DriveItemItemRequestBuilder) Invite()(*if1296dc61ef5c60fc9a2b02825626c2b7d61e21c6cdf5462ccaef6d25386fd23.InviteRequestBuilder) {
    return if1296dc61ef5c60fc9a2b02825626c2b7d61e21c6cdf5462ccaef6d25386fd23.NewInviteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ListItem provides operations to manage the listItem property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ListItem()(*icc851f6701cc3e379edf8f1c8a5a20a493cce7efe46c146a987b9411a8c91b78.ListItemRequestBuilder) {
    return icc851f6701cc3e379edf8f1c8a5a20a493cce7efe46c146a987b9411a8c91b78.NewListItemRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property items in drives
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
func (m *DriveItemItemRequestBuilder) Permissions()(*ifda175920025647e9407948753f11cc7a3fabff7fd6e426cf29edc97c408e5c2.PermissionsRequestBuilder) {
    return ifda175920025647e9407948753f11cc7a3fabff7fd6e426cf29edc97c408e5c2.NewPermissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionsById provides operations to manage the permissions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) PermissionsById(id string)(*id6e448f4be60986224ef6568b4a9b7d39a61886fbd386e3281d8989508828a70.PermissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["permission%2Did"] = id
    }
    return id6e448f4be60986224ef6568b4a9b7d39a61886fbd386e3281d8989508828a70.NewPermissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Preview provides operations to call the preview method.
func (m *DriveItemItemRequestBuilder) Preview()(*i52a4c1c4396bb6dfa2711ab3a6e974807f91b1191aa227cc5eaafc77c74b3b5a.PreviewRequestBuilder) {
    return i52a4c1c4396bb6dfa2711ab3a6e974807f91b1191aa227cc5eaafc77c74b3b5a.NewPreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *DriveItemItemRequestBuilder) Restore()(*i0175c2914b12b61cc774e65aa639603895b1ae1642fd0d5746778bab3a505278.RestoreRequestBuilder) {
    return i0175c2914b12b61cc774e65aa639603895b1ae1642fd0d5746778bab3a505278.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemItemRequestBuilder) SearchWithQ(q *string)(*iae2098c15020f00a005771b8a0dfdf0cc230504133ee3754b7e1f29877241492.SearchWithQRequestBuilder) {
    return iae2098c15020f00a005771b8a0dfdf0cc230504133ee3754b7e1f29877241492.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Subscriptions()(*i0cb44ac3ff5b85c2f1592086ca294cb612434d9f0b242efaa25ca47763aeef53.SubscriptionsRequestBuilder) {
    return i0cb44ac3ff5b85c2f1592086ca294cb612434d9f0b242efaa25ca47763aeef53.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) SubscriptionsById(id string)(*i2ab5710457ce9e7de18564cd6ca79881ae7e842006bdf4476c908f1264282bdd.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return i2ab5710457ce9e7de18564cd6ca79881ae7e842006bdf4476c908f1264282bdd.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Thumbnails provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Thumbnails()(*i05d5f2d2f7d8643c0cc65d242a1e8e39efdfff6868938e24cc4d55d7dc247a81.ThumbnailsRequestBuilder) {
    return i05d5f2d2f7d8643c0cc65d242a1e8e39efdfff6868938e24cc4d55d7dc247a81.NewThumbnailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ThumbnailsById provides operations to manage the thumbnails property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) ThumbnailsById(id string)(*ic92493f36906ee5f9721f2db0a85d2445ba0fba9bae564c8d6aaae67916102a7.ThumbnailSetItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["thumbnailSet%2Did"] = id
    }
    return ic92493f36906ee5f9721f2db0a85d2445ba0fba9bae564c8d6aaae67916102a7.NewThumbnailSetItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unfollow provides operations to call the unfollow method.
func (m *DriveItemItemRequestBuilder) Unfollow()(*i4dacd278c290b484ea557826599ab3b681b64884fd13f4110cdbc03d72fd179e.UnfollowRequestBuilder) {
    return i4dacd278c290b484ea557826599ab3b681b64884fd13f4110cdbc03d72fd179e.NewUnfollowRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ValidatePermission provides operations to call the validatePermission method.
func (m *DriveItemItemRequestBuilder) ValidatePermission()(*i6afd67fc45a6c6ed94f8b5901f953074f83bb810cbb940ad413154c343a3a948.ValidatePermissionRequestBuilder) {
    return i6afd67fc45a6c6ed94f8b5901f953074f83bb810cbb940ad413154c343a3a948.NewValidatePermissionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Versions provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) Versions()(*id9669323a3b4752f5b4c8d5f78acfd3b98dfdc33cd348c62aa7c25cc35e019f1.VersionsRequestBuilder) {
    return id9669323a3b4752f5b4c8d5f78acfd3b98dfdc33cd348c62aa7c25cc35e019f1.NewVersionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// VersionsById provides operations to manage the versions property of the microsoft.graph.driveItem entity.
func (m *DriveItemItemRequestBuilder) VersionsById(id string)(*ib57b57d66e5436184bc3dd19450c0882627b04c9aa4d6e54d68ff2c33c761387.DriveItemVersionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItemVersion%2Did"] = id
    }
    return ib57b57d66e5436184bc3dd19450c0882627b04c9aa4d6e54d68ff2c33c761387.NewDriveItemVersionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
