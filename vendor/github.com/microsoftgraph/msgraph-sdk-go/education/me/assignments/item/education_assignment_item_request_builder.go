package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i608fa3df84918be40f4d961440af7a1df9cdc70cdcbcd01332c72421f74df332 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/categories"
    i6182b98a227da91892039c997d1742545d169f3ded56acfb9ebddc2b8754acb3 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/rubric"
    i710174f1350381dcd776625bf1b593be962f1d2c9b1e368ae0bd178d3495728f "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/submissions"
    i74b91723ec2102e06ba61f2eec38b8a338e711092e99e7a6121ae3287f6f5788 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/publish"
    i874a2b579f5eb892e97d9c0c157220e5118544d0dca2de101640f4b07a1c3fe3 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/resources"
    ib7f8e3ce04c81d083f8ef638afdca65de5db8ffabbc758013b5a49cae3ee262d "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/setupfeedbackresourcesfolder"
    ie33a7d66ff2752d51b9a63d7758fc6c00305b565d547fa3e4e42f07c7a26bdec "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/setupresourcesfolder"
    i0a97b1183e4efc9029dbfdb7442c76c2d4fc13a243fc9089be3106c63fd1b7f2 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/submissions/item"
    i26284194e56af98cb4a7e2e347a5485f31cf2c535beac0528c6c3661c312a162 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/categories/item"
    i7b969696e68843bc2d2d30122cafa6f7eaf8a9843bf04feaf8693f3f794c0d94 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item/resources/item"
)

// EducationAssignmentItemRequestBuilder provides operations to manage the assignments property of the microsoft.graph.educationUser entity.
type EducationAssignmentItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EducationAssignmentItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationAssignmentItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EducationAssignmentItemRequestBuilderGetQueryParameters assignments belonging to the user.
type EducationAssignmentItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EducationAssignmentItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationAssignmentItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EducationAssignmentItemRequestBuilderGetQueryParameters
}
// EducationAssignmentItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationAssignmentItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Categories provides operations to manage the categories property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) Categories()(*i608fa3df84918be40f4d961440af7a1df9cdc70cdcbcd01332c72421f74df332.CategoriesRequestBuilder) {
    return i608fa3df84918be40f4d961440af7a1df9cdc70cdcbcd01332c72421f74df332.NewCategoriesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CategoriesById gets an item from the github.com/microsoftgraph/msgraph-sdk-go/.education.me.assignments.item.categories.item collection
func (m *EducationAssignmentItemRequestBuilder) CategoriesById(id string)(*i26284194e56af98cb4a7e2e347a5485f31cf2c535beac0528c6c3661c312a162.EducationCategoryItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationCategory%2Did"] = id
    }
    return i26284194e56af98cb4a7e2e347a5485f31cf2c535beac0528c6c3661c312a162.NewEducationCategoryItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewEducationAssignmentItemRequestBuilderInternal instantiates a new EducationAssignmentItemRequestBuilder and sets the default values.
func NewEducationAssignmentItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationAssignmentItemRequestBuilder) {
    m := &EducationAssignmentItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education/me/assignments/{educationAssignment%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEducationAssignmentItemRequestBuilder instantiates a new EducationAssignmentItemRequestBuilder and sets the default values.
func NewEducationAssignmentItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationAssignmentItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEducationAssignmentItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property assignments for education
func (m *EducationAssignmentItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EducationAssignmentItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation assignments belonging to the user.
func (m *EducationAssignmentItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EducationAssignmentItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property assignments in education
func (m *EducationAssignmentItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable, requestConfiguration *EducationAssignmentItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property assignments for education
func (m *EducationAssignmentItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EducationAssignmentItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get assignments belonging to the user.
func (m *EducationAssignmentItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EducationAssignmentItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationAssignmentFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable), nil
}
// Patch update the navigation property assignments in education
func (m *EducationAssignmentItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable, requestConfiguration *EducationAssignmentItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationAssignmentFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationAssignmentable), nil
}
// Publish provides operations to call the publish method.
func (m *EducationAssignmentItemRequestBuilder) Publish()(*i74b91723ec2102e06ba61f2eec38b8a338e711092e99e7a6121ae3287f6f5788.PublishRequestBuilder) {
    return i74b91723ec2102e06ba61f2eec38b8a338e711092e99e7a6121ae3287f6f5788.NewPublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Resources provides operations to manage the resources property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) Resources()(*i874a2b579f5eb892e97d9c0c157220e5118544d0dca2de101640f4b07a1c3fe3.ResourcesRequestBuilder) {
    return i874a2b579f5eb892e97d9c0c157220e5118544d0dca2de101640f4b07a1c3fe3.NewResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResourcesById provides operations to manage the resources property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) ResourcesById(id string)(*i7b969696e68843bc2d2d30122cafa6f7eaf8a9843bf04feaf8693f3f794c0d94.EducationAssignmentResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationAssignmentResource%2Did"] = id
    }
    return i7b969696e68843bc2d2d30122cafa6f7eaf8a9843bf04feaf8693f3f794c0d94.NewEducationAssignmentResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Rubric provides operations to manage the rubric property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) Rubric()(*i6182b98a227da91892039c997d1742545d169f3ded56acfb9ebddc2b8754acb3.RubricRequestBuilder) {
    return i6182b98a227da91892039c997d1742545d169f3ded56acfb9ebddc2b8754acb3.NewRubricRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetUpFeedbackResourcesFolder provides operations to call the setUpFeedbackResourcesFolder method.
func (m *EducationAssignmentItemRequestBuilder) SetUpFeedbackResourcesFolder()(*ib7f8e3ce04c81d083f8ef638afdca65de5db8ffabbc758013b5a49cae3ee262d.SetUpFeedbackResourcesFolderRequestBuilder) {
    return ib7f8e3ce04c81d083f8ef638afdca65de5db8ffabbc758013b5a49cae3ee262d.NewSetUpFeedbackResourcesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetUpResourcesFolder provides operations to call the setUpResourcesFolder method.
func (m *EducationAssignmentItemRequestBuilder) SetUpResourcesFolder()(*ie33a7d66ff2752d51b9a63d7758fc6c00305b565d547fa3e4e42f07c7a26bdec.SetUpResourcesFolderRequestBuilder) {
    return ie33a7d66ff2752d51b9a63d7758fc6c00305b565d547fa3e4e42f07c7a26bdec.NewSetUpResourcesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Submissions provides operations to manage the submissions property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) Submissions()(*i710174f1350381dcd776625bf1b593be962f1d2c9b1e368ae0bd178d3495728f.SubmissionsRequestBuilder) {
    return i710174f1350381dcd776625bf1b593be962f1d2c9b1e368ae0bd178d3495728f.NewSubmissionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubmissionsById provides operations to manage the submissions property of the microsoft.graph.educationAssignment entity.
func (m *EducationAssignmentItemRequestBuilder) SubmissionsById(id string)(*i0a97b1183e4efc9029dbfdb7442c76c2d4fc13a243fc9089be3106c63fd1b7f2.EducationSubmissionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSubmission%2Did"] = id
    }
    return i0a97b1183e4efc9029dbfdb7442c76c2d4fc13a243fc9089be3106c63fd1b7f2.NewEducationSubmissionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
