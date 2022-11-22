package accessreviews

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i3ed202ff0cbdf2b2edbf6f858dd5dd49845b59ab278423c07797fe43a9ccd40f "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/accessreviews/historydefinitions"
    i9a0fc83be8a156b608969476f16987d6f0c56857cc62f3dec40c0a7819428e17 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/accessreviews/definitions"
    i0fca5fbe98a7c54d24af60b510c16cbc28027514fdcbe3896a60732f5c572e2b "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/accessreviews/historydefinitions/item"
    idcb18c7708e010766ed0434cebbaf8978d559fd3ae0bed76f0413cdd312a9c3f "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/accessreviews/definitions/item"
)

// AccessReviewsRequestBuilder provides operations to manage the accessReviews property of the microsoft.graph.identityGovernance entity.
type AccessReviewsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AccessReviewsRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AccessReviewsRequestBuilderGetQueryParameters get accessReviews from identityGovernance
type AccessReviewsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// AccessReviewsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AccessReviewsRequestBuilderGetQueryParameters
}
// AccessReviewsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AccessReviewsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAccessReviewsRequestBuilderInternal instantiates a new AccessReviewsRequestBuilder and sets the default values.
func NewAccessReviewsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsRequestBuilder) {
    m := &AccessReviewsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/identityGovernance/accessReviews{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewAccessReviewsRequestBuilder instantiates a new AccessReviewsRequestBuilder and sets the default values.
func NewAccessReviewsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AccessReviewsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAccessReviewsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property accessReviews for identityGovernance
func (m *AccessReviewsRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *AccessReviewsRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get accessReviews from identityGovernance
func (m *AccessReviewsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *AccessReviewsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property accessReviews in identityGovernance
func (m *AccessReviewsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable, requestConfiguration *AccessReviewsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Definitions provides operations to manage the definitions property of the microsoft.graph.accessReviewSet entity.
func (m *AccessReviewsRequestBuilder) Definitions()(*i9a0fc83be8a156b608969476f16987d6f0c56857cc62f3dec40c0a7819428e17.DefinitionsRequestBuilder) {
    return i9a0fc83be8a156b608969476f16987d6f0c56857cc62f3dec40c0a7819428e17.NewDefinitionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DefinitionsById provides operations to manage the definitions property of the microsoft.graph.accessReviewSet entity.
func (m *AccessReviewsRequestBuilder) DefinitionsById(id string)(*idcb18c7708e010766ed0434cebbaf8978d559fd3ae0bed76f0413cdd312a9c3f.AccessReviewScheduleDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessReviewScheduleDefinition%2Did"] = id
    }
    return idcb18c7708e010766ed0434cebbaf8978d559fd3ae0bed76f0413cdd312a9c3f.NewAccessReviewScheduleDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Delete delete navigation property accessReviews for identityGovernance
func (m *AccessReviewsRequestBuilder) Delete(ctx context.Context, requestConfiguration *AccessReviewsRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get accessReviews from identityGovernance
func (m *AccessReviewsRequestBuilder) Get(ctx context.Context, requestConfiguration *AccessReviewsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessReviewSetFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable), nil
}
// HistoryDefinitions provides operations to manage the historyDefinitions property of the microsoft.graph.accessReviewSet entity.
func (m *AccessReviewsRequestBuilder) HistoryDefinitions()(*i3ed202ff0cbdf2b2edbf6f858dd5dd49845b59ab278423c07797fe43a9ccd40f.HistoryDefinitionsRequestBuilder) {
    return i3ed202ff0cbdf2b2edbf6f858dd5dd49845b59ab278423c07797fe43a9ccd40f.NewHistoryDefinitionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// HistoryDefinitionsById provides operations to manage the historyDefinitions property of the microsoft.graph.accessReviewSet entity.
func (m *AccessReviewsRequestBuilder) HistoryDefinitionsById(id string)(*i0fca5fbe98a7c54d24af60b510c16cbc28027514fdcbe3896a60732f5c572e2b.AccessReviewHistoryDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessReviewHistoryDefinition%2Did"] = id
    }
    return i0fca5fbe98a7c54d24af60b510c16cbc28027514fdcbe3896a60732f5c572e2b.NewAccessReviewHistoryDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property accessReviews in identityGovernance
func (m *AccessReviewsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable, requestConfiguration *AccessReviewsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessReviewSetFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessReviewSetable), nil
}
