package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1548a312864f6a8093d6647950c34611e4ac02220716d0d28a7ed6ea68563e15 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/setupresourcesfolder"
    i3b9a3bc2cbb04d5774a8df86388bde86947a7b7d7c62e9d673dc25fa1f129f54 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/submittedresources"
    i4e04ce4811df69ad23b6cc639626ab84d384f66c5cbdc3157ea4eafa7746a09c "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/return_escaped"
    ia8cb4af48769798156b6d22780d66d72a7e7275c6bc9e7bc60a98007f81716fa "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/resources"
    ia9b02e54f8f0a390b69bcaf5d845b0f358c5f9201a3ec774316fdc87c9d21ca7 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/unsubmit"
    id00761bb21207f8c90cc6b9e657e4d6592705ed171fec38e1c0ce6a4e6b0fa7d "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/submit"
    id2c2532be9783e0b7091664a606705196b5b520c9fd05db1ea40f6cdf4bc5944 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/outcomes"
    id76ff41c49808e9ab3c419e39394c8b5a49d6375f73b75e7363fb291151a2bed "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/reassign"
    i363dd35cda83aa9a020591fb0f4fd184033fd6ab5dbc6a251a3af8f414725817 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/submittedresources/item"
    i71aa111a74a29d833b5fd0653f7a2d786fe91774ffbbc1f1ed0064bc02923991 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/outcomes/item"
    ia7af5a2c3efa196597d6805e5921ce27fb2948e0a38a819bfc0b0c6cf8eef9da "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item/submissions/item/resources/item"
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
    m.urlTemplate = "{+baseurl}/education/users/{educationUser%2Did}/assignments/{educationAssignment%2Did}/submissions/{educationSubmission%2Did}{?%24select,%24expand}";
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
func (m *EducationSubmissionItemRequestBuilder) Outcomes()(*id2c2532be9783e0b7091664a606705196b5b520c9fd05db1ea40f6cdf4bc5944.OutcomesRequestBuilder) {
    return id2c2532be9783e0b7091664a606705196b5b520c9fd05db1ea40f6cdf4bc5944.NewOutcomesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OutcomesById provides operations to manage the outcomes property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) OutcomesById(id string)(*i71aa111a74a29d833b5fd0653f7a2d786fe91774ffbbc1f1ed0064bc02923991.EducationOutcomeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationOutcome%2Did"] = id
    }
    return i71aa111a74a29d833b5fd0653f7a2d786fe91774ffbbc1f1ed0064bc02923991.NewEducationOutcomeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *EducationSubmissionItemRequestBuilder) Reassign()(*id76ff41c49808e9ab3c419e39394c8b5a49d6375f73b75e7363fb291151a2bed.ReassignRequestBuilder) {
    return id76ff41c49808e9ab3c419e39394c8b5a49d6375f73b75e7363fb291151a2bed.NewReassignRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Resources provides operations to manage the resources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) Resources()(*ia8cb4af48769798156b6d22780d66d72a7e7275c6bc9e7bc60a98007f81716fa.ResourcesRequestBuilder) {
    return ia8cb4af48769798156b6d22780d66d72a7e7275c6bc9e7bc60a98007f81716fa.NewResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResourcesById provides operations to manage the resources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) ResourcesById(id string)(*ia7af5a2c3efa196597d6805e5921ce27fb2948e0a38a819bfc0b0c6cf8eef9da.EducationSubmissionResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSubmissionResource%2Did"] = id
    }
    return ia7af5a2c3efa196597d6805e5921ce27fb2948e0a38a819bfc0b0c6cf8eef9da.NewEducationSubmissionResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Return_escaped provides operations to call the return method.
func (m *EducationSubmissionItemRequestBuilder) Return_escaped()(*i4e04ce4811df69ad23b6cc639626ab84d384f66c5cbdc3157ea4eafa7746a09c.ReturnRequestBuilder) {
    return i4e04ce4811df69ad23b6cc639626ab84d384f66c5cbdc3157ea4eafa7746a09c.NewReturnRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetUpResourcesFolder provides operations to call the setUpResourcesFolder method.
func (m *EducationSubmissionItemRequestBuilder) SetUpResourcesFolder()(*i1548a312864f6a8093d6647950c34611e4ac02220716d0d28a7ed6ea68563e15.SetUpResourcesFolderRequestBuilder) {
    return i1548a312864f6a8093d6647950c34611e4ac02220716d0d28a7ed6ea68563e15.NewSetUpResourcesFolderRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Submit provides operations to call the submit method.
func (m *EducationSubmissionItemRequestBuilder) Submit()(*id00761bb21207f8c90cc6b9e657e4d6592705ed171fec38e1c0ce6a4e6b0fa7d.SubmitRequestBuilder) {
    return id00761bb21207f8c90cc6b9e657e4d6592705ed171fec38e1c0ce6a4e6b0fa7d.NewSubmitRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubmittedResources provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) SubmittedResources()(*i3b9a3bc2cbb04d5774a8df86388bde86947a7b7d7c62e9d673dc25fa1f129f54.SubmittedResourcesRequestBuilder) {
    return i3b9a3bc2cbb04d5774a8df86388bde86947a7b7d7c62e9d673dc25fa1f129f54.NewSubmittedResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubmittedResourcesById provides operations to manage the submittedResources property of the microsoft.graph.educationSubmission entity.
func (m *EducationSubmissionItemRequestBuilder) SubmittedResourcesById(id string)(*i363dd35cda83aa9a020591fb0f4fd184033fd6ab5dbc6a251a3af8f414725817.EducationSubmissionResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSubmissionResource%2Did"] = id
    }
    return i363dd35cda83aa9a020591fb0f4fd184033fd6ab5dbc6a251a3af8f414725817.NewEducationSubmissionResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Unsubmit provides operations to call the unsubmit method.
func (m *EducationSubmissionItemRequestBuilder) Unsubmit()(*ia9b02e54f8f0a390b69bcaf5d845b0f358c5f9201a3ec774316fdc87c9d21ca7.UnsubmitRequestBuilder) {
    return ia9b02e54f8f0a390b69bcaf5d845b0f358c5f9201a3ec774316fdc87c9d21ca7.NewUnsubmitRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
