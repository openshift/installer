package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i16d16fffd3d9353d2aff2410b702588904cdd095a5913737677f0d820f80b04b "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/attachmentsessions"
    i4416c2ee3fc470c00146c53c47e2f9da5d169d6ba3e81e76a033e46a4c75e39f "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/checklistitems"
    i5ae9cc2252fe39eb1f8922688cd3cd931f7af1d4c5b712b1c8dbcaee74ac0a8b "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/linkedresources"
    i5ebd0ff7b1c5a72fa996da9d25188fed331be10734b9158584fca359e3382cc8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/extensions"
    i7b63357d67e15d681025f4532a2c7287b79097fc5c22a03ddaa324a302bc779b "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/attachments"
    i4be9fed866d0b4477030fa0f50dd4b81d52fdbe2080d06e4df645edba54bcd37 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/checklistitems/item"
    i570f0cd86cb59984a251e242cc6974a8247783f01e37cad24f260d6c4655aad3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/linkedresources/item"
    i72b65362f27c7879a5a62cde2844e68df9a71e7f9744f1b9814731620237f095 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/attachmentsessions/item"
    ibded91745af4973d993c7591fd0006231f1bf43270adfb384d6f7afa7427e666 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/extensions/item"
    icff305394eb0be29ccbb725d2beb1ad07cf52bf65964524abff6848d47b12dfd "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo/lists/item/tasks/item/attachments/item"
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
func (m *TodoTaskItemRequestBuilder) Attachments()(*i7b63357d67e15d681025f4532a2c7287b79097fc5c22a03ddaa324a302bc779b.AttachmentsRequestBuilder) {
    return i7b63357d67e15d681025f4532a2c7287b79097fc5c22a03ddaa324a302bc779b.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentsById(id string)(*icff305394eb0be29ccbb725d2beb1ad07cf52bf65964524abff6848d47b12dfd.AttachmentBaseItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachmentBase%2Did"] = id
    }
    return icff305394eb0be29ccbb725d2beb1ad07cf52bf65964524abff6848d47b12dfd.NewAttachmentBaseItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AttachmentSessions provides operations to manage the attachmentSessions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentSessions()(*i16d16fffd3d9353d2aff2410b702588904cdd095a5913737677f0d820f80b04b.AttachmentSessionsRequestBuilder) {
    return i16d16fffd3d9353d2aff2410b702588904cdd095a5913737677f0d820f80b04b.NewAttachmentSessionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentSessionsById provides operations to manage the attachmentSessions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) AttachmentSessionsById(id string)(*i72b65362f27c7879a5a62cde2844e68df9a71e7f9744f1b9814731620237f095.AttachmentSessionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachmentSession%2Did"] = id
    }
    return i72b65362f27c7879a5a62cde2844e68df9a71e7f9744f1b9814731620237f095.NewAttachmentSessionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ChecklistItems provides operations to manage the checklistItems property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ChecklistItems()(*i4416c2ee3fc470c00146c53c47e2f9da5d169d6ba3e81e76a033e46a4c75e39f.ChecklistItemsRequestBuilder) {
    return i4416c2ee3fc470c00146c53c47e2f9da5d169d6ba3e81e76a033e46a4c75e39f.NewChecklistItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChecklistItemsById provides operations to manage the checklistItems property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ChecklistItemsById(id string)(*i4be9fed866d0b4477030fa0f50dd4b81d52fdbe2080d06e4df645edba54bcd37.ChecklistItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["checklistItem%2Did"] = id
    }
    return i4be9fed866d0b4477030fa0f50dd4b81d52fdbe2080d06e4df645edba54bcd37.NewChecklistItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewTodoTaskItemRequestBuilderInternal instantiates a new TodoTaskItemRequestBuilder and sets the default values.
func NewTodoTaskItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TodoTaskItemRequestBuilder) {
    m := &TodoTaskItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/todo/lists/{todoTaskList%2Did}/tasks/{todoTask%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property tasks for users
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
// CreatePatchRequestInformation update the navigation property tasks in users
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
// Delete delete navigation property tasks for users
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
func (m *TodoTaskItemRequestBuilder) Extensions()(*i5ebd0ff7b1c5a72fa996da9d25188fed331be10734b9158584fca359e3382cc8.ExtensionsRequestBuilder) {
    return i5ebd0ff7b1c5a72fa996da9d25188fed331be10734b9158584fca359e3382cc8.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) ExtensionsById(id string)(*ibded91745af4973d993c7591fd0006231f1bf43270adfb384d6f7afa7427e666.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ibded91745af4973d993c7591fd0006231f1bf43270adfb384d6f7afa7427e666.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *TodoTaskItemRequestBuilder) LinkedResources()(*i5ae9cc2252fe39eb1f8922688cd3cd931f7af1d4c5b712b1c8dbcaee74ac0a8b.LinkedResourcesRequestBuilder) {
    return i5ae9cc2252fe39eb1f8922688cd3cd931f7af1d4c5b712b1c8dbcaee74ac0a8b.NewLinkedResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LinkedResourcesById provides operations to manage the linkedResources property of the microsoft.graph.todoTask entity.
func (m *TodoTaskItemRequestBuilder) LinkedResourcesById(id string)(*i570f0cd86cb59984a251e242cc6974a8247783f01e37cad24f260d6c4655aad3.LinkedResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["linkedResource%2Did"] = id
    }
    return i570f0cd86cb59984a251e242cc6974a8247783f01e37cad24f260d6c4655aad3.NewLinkedResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property tasks in users
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
