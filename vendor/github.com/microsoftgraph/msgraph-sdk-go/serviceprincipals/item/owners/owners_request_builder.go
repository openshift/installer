package owners

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i084a21f56aa1fe12f872497ee18ea095e7b49a0658852babf734baf24a2d2305 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/serviceprincipal"
    i11240f2dcce55be8b4c482531f2c78914edd81d16e91c8629b02240cd64c6ec4 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/approleassignment"
    i3fd51d7cedbd61b7eb9da3199c41a984c523da2ba70cd0061dbe5ac3af02d349 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/ref"
    i7b3b46b21232c60fd5f3f77a6cf95c2583005c86186b50804719398e6527d9d2 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/user"
    i8f09a4ea8103ae1accba37f73966dca630c46ce795b2f32304bebbd74524215f "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/endpoint"
    i93939644e4bd092f76ff55d542b68c4603adbb8079284fd4d4882297e8d691a1 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/count"
)

// OwnersRequestBuilder provides operations to manage the owners property of the microsoft.graph.servicePrincipal entity.
type OwnersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OwnersRequestBuilderGetQueryParameters directory objects that are owners of this servicePrincipal. The owners are a set of non-admin users or servicePrincipals who are allowed to modify this object. Read-only. Nullable. Supports $expand.
type OwnersRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// OwnersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OwnersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OwnersRequestBuilderGetQueryParameters
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *OwnersRequestBuilder) AppRoleAssignment()(*i11240f2dcce55be8b4c482531f2c78914edd81d16e91c8629b02240cd64c6ec4.AppRoleAssignmentRequestBuilder) {
    return i11240f2dcce55be8b4c482531f2c78914edd81d16e91c8629b02240cd64c6ec4.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOwnersRequestBuilderInternal instantiates a new OwnersRequestBuilder and sets the default values.
func NewOwnersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnersRequestBuilder) {
    m := &OwnersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/owners{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOwnersRequestBuilder instantiates a new OwnersRequestBuilder and sets the default values.
func NewOwnersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOwnersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *OwnersRequestBuilder) Count()(*i93939644e4bd092f76ff55d542b68c4603adbb8079284fd4d4882297e8d691a1.CountRequestBuilder) {
    return i93939644e4bd092f76ff55d542b68c4603adbb8079284fd4d4882297e8d691a1.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation directory objects that are owners of this servicePrincipal. The owners are a set of non-admin users or servicePrincipals who are allowed to modify this object. Read-only. Nullable. Supports $expand.
func (m *OwnersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OwnersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Endpoint casts the previous resource to endpoint.
func (m *OwnersRequestBuilder) Endpoint()(*i8f09a4ea8103ae1accba37f73966dca630c46ce795b2f32304bebbd74524215f.EndpointRequestBuilder) {
    return i8f09a4ea8103ae1accba37f73966dca630c46ce795b2f32304bebbd74524215f.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get directory objects that are owners of this servicePrincipal. The owners are a set of non-admin users or servicePrincipals who are allowed to modify this object. Read-only. Nullable. Supports $expand.
func (m *OwnersRequestBuilder) Get(ctx context.Context, requestConfiguration *OwnersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
// Ref provides operations to manage the collection of servicePrincipal entities.
func (m *OwnersRequestBuilder) Ref()(*i3fd51d7cedbd61b7eb9da3199c41a984c523da2ba70cd0061dbe5ac3af02d349.RefRequestBuilder) {
    return i3fd51d7cedbd61b7eb9da3199c41a984c523da2ba70cd0061dbe5ac3af02d349.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *OwnersRequestBuilder) ServicePrincipal()(*i084a21f56aa1fe12f872497ee18ea095e7b49a0658852babf734baf24a2d2305.ServicePrincipalRequestBuilder) {
    return i084a21f56aa1fe12f872497ee18ea095e7b49a0658852babf734baf24a2d2305.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *OwnersRequestBuilder) User()(*i7b3b46b21232c60fd5f3f77a6cf95c2583005c86186b50804719398e6527d9d2.UserRequestBuilder) {
    return i7b3b46b21232c60fd5f3f77a6cf95c2583005c86186b50804719398e6527d9d2.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
