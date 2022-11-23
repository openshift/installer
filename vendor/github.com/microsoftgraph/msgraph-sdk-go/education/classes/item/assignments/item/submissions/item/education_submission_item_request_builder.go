package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i3c6015cf8457d8d8ec01d019f3dbc202765651c2130e73641cf6b06caeb2d228 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/setupresourcesfolder"
    i4897ecfe261ff150792387d2c4b3c93a4a0c1fa8d17924ec2e00dd6cf21d022d "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/return_escaped"
    i6676fa5ab422c72e51792bf14b597f55c25d08008768ede327ecc59a581d7d7c "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/reassign"
    i685e1486255b2d524092d2de81f1a7d1d38632cbb99a29975e93cc433375ed77 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/submit"
    ic9e6f044df077c5c49203af4cbb1cb4bdd990da9a1e6a2c50a189335ae8c62fc "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/outcomes"
    id25fa828362bc547da297fba40d2cdbe268a1fd50d31a4ead0dd41ffa77fdfc6 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/resources"
    id45b66ebeec3910c7bc64c11a6669bb5fd10605cbdd08bde69bcaa49dd9ab2a7 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/unsubmit"
    id595ffd3ed593bffbb73e849bb31881e531d536f5c57d0b37c7933bb34ad1ad0 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/submittedresources"
    i9a0443e0b2ffd58714a9cae0542bdec49c7a1c21111a0cba91c45aa013014f81 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/outcomes/item"
    idf6da69a1128e56db2a075a35c90dfb9d769744b063f36a369e943680fa35795 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/resources/item"
    if2b3a1bb63b3f0f978027b12a06272583c8a6c2ebe1a551f1e1b72a751ef4004 "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/submissions/item/submittedresources/item"
)

// EducationSubmissionItemRequestBuilder provides operations to manage the submissions property of the microsoft.graph.educationAssignment entity.
type EducationSubmissionItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EducationSubmissionItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationSubmissionItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EducationSubmissionItemRequestBuilderGetQueryParameters once published, there is a submission object for each student representing their work and grade.  Read-only. Nullable.
type EducationSubmissionItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EducationSubmissionItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationSubmissionItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EducationSubmissionItemRequestBuilderGetQueryParameters
}
// EducationSubmissionItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationSubmissionItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewEducationSubmissionItemRequestBuilderInternal instantiates a new EducationSubmissionItemRequestBuilder and sets the default values.
func NewEducationSubmissionItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationSubmissionItemRequestBuilder) {
    m := &EducationSubmissionItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education/classes/{educationClass%2Did}/assignments/{educationAssignment%2Did}/submissions/{educationSubmission%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEducationSubmissionItemRequestBuilder instantiates a new EducationSubmissionItemRequestBuilder and sets the default values.
func NewEducationSubmissionItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationSubmissionItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEducationSubmissionItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property submissions for education
func (m *EducationSubmissionItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EducationSubmissionItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation once published, there is a submission object for each student representing their work and grade.  Read-only. Nullable.
func (m *EducationSubmissionItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EducationSubmissionItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property submissions in education
func (m *EducationSubmissionItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable, requestConfiguration *EducationSubmissionItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property submissions for education
func (m *EducationSubmissionItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EducationSubmissionItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get once published, there is a submission object for each student representing their work and grade.  Read-only. Nullable.
func (m *EducationSubmissionItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EducationSubmissionItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationSubmissionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable), nil
}
// Outcomes provides operations to manage the outcomes property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) Outcomes()(*ic9e6f044df077c5c49203af4cbb1cb4bdd990da9a1e6a2c50a189335ae8c62fc.OutcomesRequestBuilder) {
    return ic9e6f044df077c5c49203af4cbb1cb4bdd990da9a1e6a2c50a189335ae8c62fc.NewOutcomesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OutcomesById provides operations to manage the outcomes property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) OutcomesById(id string)(*i9a0443e0b2ffd58714a9cae0542bdec49c7a1c21111a0cba91c45aa013014f81.EducationOutcomeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationOutcome%2Did"] = id
    }
    return i9a0443e0b2ffd58714a9cae0542bdec49c7a1c21111a0cba91c45aa013014f81.NewEducationOutcomeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property submissions in education
func (m *EducationSubmissionItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable, requestConfiguration *EducationSubmissionItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationSubmissionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationSubmissionable), nil
}
// Reassign provides operations to call the reassign method.
func (m *EducationSubmissionItemRequestBuilder) Reassign()(*i6676fa5ab422c72e51792bf14b597f55c25d08008768ede327ecc59a581d7d7c.ReassignRequestBuilder) {
    return i6676fa5ab422c72e51792bf14b597f55c25d08008768ede327ecc59a581d7d7c.NewReassignRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Resources provides operations to manage the resources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) Resources()(*id25fa828362bc547da297fba40d2cdbe268a1fd50d31a4ead0dd41ffa77fdfc6.ResourcesRequestBuilder) {
    return id25fa828362bc547da297fba40d2cdbe268a1fd50d31a4ead0dd41ffa77fdfc6.NewResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResourcesById provides operations to manage the resources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) ResourcesById(id string)(*idf6da69a1128e56db2a075a35c90dfb9d769744b063f36a369e943680fa35795.EducationSubmissionResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSubmissionResource%2Did"] = id
    }
    return idf6da69a1128e56db2a075a35c90dfb9d769744b063f36a369e943680fa35795.NewEducationSubmissionResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Return_escaped provides operations to call the return method.
func (m *EducationSubmissionItemRequestBuilder) Return_escaped()(*i4897ecfe261ff150792387d2c4b3c93a4a0c1fa8d17924ec2e00dd6cf21d022d.ReturnRequestBuilder) {
    return i4897ecfe261ff150792387d2c4b3c93a4a0c1fa8d17924ec2e00dd6cf21d022d.NewReturnRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetUpResourcesFolder provides operations to call the setUpResourcesFolder method.
func (m *EducationSubmissionItemRequestBuilder) SetUpResourcesFolder()(*i3c6015cf8457d8d8ec01d019f3dbc202765651c2130e73641cf6b06caeb2d228.SetUpResourcesFolderRequestBuilder) {
    return i3c6015cf8457d8d8ec01d019f3dbc202765651c2130e73641cf6b06caeb2d228.NewSetUpResourcesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Submit provides operations to call the submit method.
func (m *EducationSubmissionItemRequestBuilder) Submit()(*i685e1486255b2d524092d2de81f1a7d1d38632cbb99a29975e93cc433375ed77.SubmitRequestBuilder) {
    return i685e1486255b2d524092d2de81f1a7d1d38632cbb99a29975e93cc433375ed77.NewSubmitRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubmittedResources provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) SubmittedResources()(*id595ffd3ed593bffbb73e849bb31881e531d536f5c57d0b37c7933bb34ad1ad0.SubmittedResourcesRequestBuilder) {
    return id595ffd3ed593bffbb73e849bb31881e531d536f5c57d0b37c7933bb34ad1ad0.NewSubmittedResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubmittedResourcesById provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) SubmittedResourcesById(id string)(*if2b3a1bb63b3f0f978027b12a06272583c8a6c2ebe1a551f1e1b72a751ef4004.EducationSubmissionResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSubmissionResource%2Did"] = id
    }
    return if2b3a1bb63b3f0f978027b12a06272583c8a6c2ebe1a551f1e1b72a751ef4004.NewEducationSubmissionResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unsubmit provides operations to call the unsubmit method.
func (m *EducationSubmissionItemRequestBuilder) Unsubmit()(*id45b66ebeec3910c7bc64c11a6669bb5fd10605cbdd08bde69bcaa49dd9ab2a7.UnsubmitRequestBuilder) {
    return id45b66ebeec3910c7bc64c11a6669bb5fd10605cbdd08bde69bcaa49dd9ab2a7.NewUnsubmitRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
