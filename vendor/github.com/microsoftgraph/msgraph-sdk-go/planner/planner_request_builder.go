package planner

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i2e30aee32e30369c8ecc6bc38d177e6784a40fa62474eccf2569928795aaa39d "github.com/microsoftgraph/msgraph-sdk-go/planner/plans"
    i986fc893717836028d97e76e59d0d0f8dcd31663c6dae0666baf59f70e2d3160 "github.com/microsoftgraph/msgraph-sdk-go/planner/tasks"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ib9a05c969511c90134a6d709165bc74e1713b24df408cc128f1e0dd7ead93ed5 "github.com/microsoftgraph/msgraph-sdk-go/planner/buckets"
    i61fc96c998992d6a9f34876bca0e09cf8fd042f82ea4e76e5bcd7ec6488b05fb "github.com/microsoftgraph/msgraph-sdk-go/planner/buckets/item"
    ic1589f257c7e9d50ebdff1741b3e8c9aaea2fa2dbe14575d8dfdfaa82594b4be "github.com/microsoftgraph/msgraph-sdk-go/planner/tasks/item"
    if0be1a8c4a5c143221ab5041a7fe062880699c2169cf3b86e89b735af5f493e2 "github.com/microsoftgraph/msgraph-sdk-go/planner/plans/item"
)

// PlannerRequestBuilder provides operations to manage the planner singleton.
type PlannerRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PlannerRequestBuilderGetQueryParameters get planner
type PlannerRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PlannerRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PlannerRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PlannerRequestBuilderGetQueryParameters
}
// PlannerRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PlannerRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Buckets provides operations to manage the buckets property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) Buckets()(*ib9a05c969511c90134a6d709165bc74e1713b24df408cc128f1e0dd7ead93ed5.BucketsRequestBuilder) {
    return ib9a05c969511c90134a6d709165bc74e1713b24df408cc128f1e0dd7ead93ed5.NewBucketsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BucketsById provides operations to manage the buckets property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) BucketsById(id string)(*i61fc96c998992d6a9f34876bca0e09cf8fd042f82ea4e76e5bcd7ec6488b05fb.PlannerBucketItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["plannerBucket%2Did"] = id
    }
    return i61fc96c998992d6a9f34876bca0e09cf8fd042f82ea4e76e5bcd7ec6488b05fb.NewPlannerBucketItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPlannerRequestBuilderInternal instantiates a new PlannerRequestBuilder and sets the default values.
func NewPlannerRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PlannerRequestBuilder) {
    m := &PlannerRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/planner{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPlannerRequestBuilder instantiates a new PlannerRequestBuilder and sets the default values.
func NewPlannerRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PlannerRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPlannerRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get planner
func (m *PlannerRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PlannerRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update planner
func (m *PlannerRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable, requestConfiguration *PlannerRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get planner
func (m *PlannerRequestBuilder) Get(ctx context.Context, requestConfiguration *PlannerRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePlannerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable), nil
}
// Patch update planner
func (m *PlannerRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable, requestConfiguration *PlannerRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePlannerFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Plannerable), nil
}
// Plans provides operations to manage the plans property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) Plans()(*i2e30aee32e30369c8ecc6bc38d177e6784a40fa62474eccf2569928795aaa39d.PlansRequestBuilder) {
    return i2e30aee32e30369c8ecc6bc38d177e6784a40fa62474eccf2569928795aaa39d.NewPlansRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PlansById provides operations to manage the plans property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) PlansById(id string)(*if0be1a8c4a5c143221ab5041a7fe062880699c2169cf3b86e89b735af5f493e2.PlannerPlanItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["plannerPlan%2Did"] = id
    }
    return if0be1a8c4a5c143221ab5041a7fe062880699c2169cf3b86e89b735af5f493e2.NewPlannerPlanItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Tasks provides operations to manage the tasks property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) Tasks()(*i986fc893717836028d97e76e59d0d0f8dcd31663c6dae0666baf59f70e2d3160.TasksRequestBuilder) {
    return i986fc893717836028d97e76e59d0d0f8dcd31663c6dae0666baf59f70e2d3160.NewTasksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TasksById provides operations to manage the tasks property of the microsoft.graph.planner entity.
func (m *PlannerRequestBuilder) TasksById(id string)(*ic1589f257c7e9d50ebdff1741b3e8c9aaea2fa2dbe14575d8dfdfaa82594b4be.PlannerTaskItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["plannerTask%2Did"] = id
    }
    return ic1589f257c7e9d50ebdff1741b3e8c9aaea2fa2dbe14575d8dfdfaa82594b4be.NewPlannerTaskItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
