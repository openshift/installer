package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0c19321a6ceef8817fa8eb99e841918f9dd2745ece282552b7513734613d3b92 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/checklistitems"
    i4bdbca5ae870f2ef282f2d3b21ef1ad8e59e0bd0730fc09f21c3157bf9515d00 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/attachments"
    i94f761dcaaceed0dc77f9f41ede8dab2e332f6dce89c821f046f7df98a8f1329 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/attachmentsessions"
    ib328057042b40aca0435e29217cc75d0b307a22c3fa6a95d56b669cd767d0505 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/extensions"
    idb041647a8374ed00dc0466469e256d948e2d774e8f6c2d867a2a02b3bc54399 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/linkedresources"
    i21c749bd06a5a477f77836251cc6fc9de0784ec9712ca06407cbed36c02fe0d8 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/attachments/item"
    i3731eec189d2b156f71a80f052e8a31ed9d925de97432bc122f165e4c13dc655 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/checklistitems/item"
    ia6e8ba8256cc0f982965ab6b7cc3a994f177bbc1e73b05c133e11c04133cbbaf "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/extensions/item"
    ie07af192062a7164b8c25ab2158765e4eb3925d07dcd66ef69fed8ceb2679be9 "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/linkedresources/item"
    if2b06c39d58538fa2faf9595143504d73e2a666d471f945daeb14af441e11bee "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks/item/attachmentsessions/item"
)

// TodoTaskItemRequestBuilder provides operations to manage the tasks property of the microsoft.graph.todoTaskList entity.
type TodoTaskItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TodoTaskItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TodoTaskItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TodoTaskItemRequestBuilderGetQueryParameters the tasks in this task list. Read-only. Nullable.
type TodoTaskItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TodoTaskItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TodoTaskItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TodoTaskItemRequestBuilderGetQueryParameters
}
// TodoTaskItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TodoTaskItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) Attachments()(*i4bdbca5ae870f2ef282f2d3b21ef1ad8e59e0bd0730fc09f21c3157bf9515d00.AttachmentsRequestBuilder) {
    return i4bdbca5ae870f2ef282f2d3b21ef1ad8e59e0bd0730fc09f21c3157bf9515d00.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentsById(id string)(*i21c749bd06a5a477f77836251cc6fc9de0784ec9712ca06407cbed36c02fe0d8.AttachmentBaseItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachmentBase%2Did"] = id
    }
    return i21c749bd06a5a477f77836251cc6fc9de0784ec9712ca06407cbed36c02fe0d8.NewAttachmentBaseItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AttachmentSessions provides operations to manage the attachmentSessions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentSessions()(*i94f761dcaaceed0dc77f9f41ede8dab2e332f6dce89c821f046f7df98a8f1329.AttachmentSessionsRequestBuilder) {
    return i94f761dcaaceed0dc77f9f41ede8dab2e332f6dce89c821f046f7df98a8f1329.NewAttachmentSessionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentSessionsById provides operations to manage the attachmentSessions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentSessionsById(id string)(*if2b06c39d58538fa2faf9595143504d73e2a666d471f945daeb14af441e11bee.AttachmentSessionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachmentSession%2Did"] = id
    }
    return if2b06c39d58538fa2faf9595143504d73e2a666d471f945daeb14af441e11bee.NewAttachmentSessionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ChecklistItems provides operations to manage the checklistItems property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ChecklistItems()(*i0c19321a6ceef8817fa8eb99e841918f9dd2745ece282552b7513734613d3b92.ChecklistItemsRequestBuilder) {
    return i0c19321a6ceef8817fa8eb99e841918f9dd2745ece282552b7513734613d3b92.NewChecklistItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChecklistItemsById provides operations to manage the checklistItems property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ChecklistItemsById(id string)(*i3731eec189d2b156f71a80f052e8a31ed9d925de97432bc122f165e4c13dc655.ChecklistItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["checklistItem%2Did"] = id
    }
    return i3731eec189d2b156f71a80f052e8a31ed9d925de97432bc122f165e4c13dc655.NewChecklistItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewTodoTaskItemRequestBuilderInternal instantiates a new TodoTaskItemRequestBuilder and sets the default values.
func NewTodoTaskItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TodoTaskItemRequestBuilder) {
    m := &TodoTaskItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/todo/lists/{todoTaskList%2Did}/tasks/{todoTask%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTodoTaskItemRequestBuilder instantiates a new TodoTaskItemRequestBuilder and sets the default values.
func NewTodoTaskItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TodoTaskItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTodoTaskItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property tasks for me
func (m *TodoTaskItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *TodoTaskItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the tasks in this task list. Read-only. Nullable.
func (m *TodoTaskItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TodoTaskItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property tasks in me
func (m *TodoTaskItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable, requestConfiguration *TodoTaskItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property tasks for me
func (m *TodoTaskItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TodoTaskItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Extensions provides operations to manage the extensions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) Extensions()(*ib328057042b40aca0435e29217cc75d0b307a22c3fa6a95d56b669cd767d0505.ExtensionsRequestBuilder) {
    return ib328057042b40aca0435e29217cc75d0b307a22c3fa6a95d56b669cd767d0505.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ExtensionsById(id string)(*ia6e8ba8256cc0f982965ab6b7cc3a994f177bbc1e73b05c133e11c04133cbbaf.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ia6e8ba8256cc0f982965ab6b7cc3a994f177bbc1e73b05c133e11c04133cbbaf.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get the tasks in this task list. Read-only. Nullable.
func (m *TodoTaskItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TodoTaskItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTodoTaskFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable), nil
}
// LinkedResources provides operations to manage the linkedResources property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) LinkedResources()(*idb041647a8374ed00dc0466469e256d948e2d774e8f6c2d867a2a02b3bc54399.LinkedResourcesRequestBuilder) {
    return idb041647a8374ed00dc0466469e256d948e2d774e8f6c2d867a2a02b3bc54399.NewLinkedResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LinkedResourcesById provides operations to manage the linkedResources property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) LinkedResourcesById(id string)(*ie07af192062a7164b8c25ab2158765e4eb3925d07dcd66ef69fed8ceb2679be9.LinkedResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["linkedResource%2Did"] = id
    }
    return ie07af192062a7164b8c25ab2158765e4eb3925d07dcd66ef69fed8ceb2679be9.NewLinkedResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property tasks in me
func (m *TodoTaskItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable, requestConfiguration *TodoTaskItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTodoTaskFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TodoTaskable), nil
}
