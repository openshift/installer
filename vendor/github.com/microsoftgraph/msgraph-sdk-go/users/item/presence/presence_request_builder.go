package presence

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i8bf936b25dc8e4c915fe2810378417d98d4c925f91faf78a54fb5f6a6fef43e6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/presence/clearuserpreferredpresence"
    i8e9793b98d900ff8e5528c9b813d76bda50532e9e5b1a489b1b63409d340eb16 "github.com/microsoftgraph/msgraph-sdk-go/users/item/presence/setpresence"
    ibe4344f966fe5874901ba1623cf995c4f53b91b7a5932404a878ed5e3ad4ed46 "github.com/microsoftgraph/msgraph-sdk-go/users/item/presence/setuserpreferredpresence"
    id0d15f0c20e289b536f406687c001f0dd217ca0b8887da9f7df46f047d006b0f "github.com/microsoftgraph/msgraph-sdk-go/users/item/presence/clearpresence"
)

// PresenceRequestBuilder provides operations to manage the presence property of the microsoft.graph.user entity.
type PresenceRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PresenceRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PresenceRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// PresenceRequestBuilderGetQueryParameters get a user's presence information.
type PresenceRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PresenceRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PresenceRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PresenceRequestBuilderGetQueryParameters
}
// PresenceRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PresenceRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ClearPresence provides operations to call the clearPresence method.
func (m *PresenceRequestBuilder) ClearPresence()(*id0d15f0c20e289b536f406687c001f0dd217ca0b8887da9f7df46f047d006b0f.ClearPresenceRequestBuilder) {
    return id0d15f0c20e289b536f406687c001f0dd217ca0b8887da9f7df46f047d006b0f.NewClearPresenceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClearUserPreferredPresence provides operations to call the clearUserPreferredPresence method.
func (m *PresenceRequestBuilder) ClearUserPreferredPresence()(*i8bf936b25dc8e4c915fe2810378417d98d4c925f91faf78a54fb5f6a6fef43e6.ClearUserPreferredPresenceRequestBuilder) {
    return i8bf936b25dc8e4c915fe2810378417d98d4c925f91faf78a54fb5f6a6fef43e6.NewClearUserPreferredPresenceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewPresenceRequestBuilderInternal instantiates a new PresenceRequestBuilder and sets the default values.
func NewPresenceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PresenceRequestBuilder) {
    m := &PresenceRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/presence{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPresenceRequestBuilder instantiates a new PresenceRequestBuilder and sets the default values.
func NewPresenceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PresenceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPresenceRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property presence for users
func (m *PresenceRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *PresenceRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get a user's presence information.
func (m *PresenceRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PresenceRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property presence in users
func (m *PresenceRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable, requestConfiguration *PresenceRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property presence for users
func (m *PresenceRequestBuilder) Delete(ctx context.Context, requestConfiguration *PresenceRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get a user's presence information.
func (m *PresenceRequestBuilder) Get(ctx context.Context, requestConfiguration *PresenceRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePresenceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable), nil
}
// Patch update the navigation property presence in users
func (m *PresenceRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable, requestConfiguration *PresenceRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePresenceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Presenceable), nil
}
// SetPresence provides operations to call the setPresence method.
func (m *PresenceRequestBuilder) SetPresence()(*i8e9793b98d900ff8e5528c9b813d76bda50532e9e5b1a489b1b63409d340eb16.SetPresenceRequestBuilder) {
    return i8e9793b98d900ff8e5528c9b813d76bda50532e9e5b1a489b1b63409d340eb16.NewSetPresenceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SetUserPreferredPresence provides operations to call the setUserPreferredPresence method.
func (m *PresenceRequestBuilder) SetUserPreferredPresence()(*ibe4344f966fe5874901ba1623cf995c4f53b91b7a5932404a878ed5e3ad4ed46.SetUserPreferredPresenceRequestBuilder) {
    return ibe4344f966fe5874901ba1623cf995c4f53b91b7a5932404a878ed5e3ad4ed46.NewSetUserPreferredPresenceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
