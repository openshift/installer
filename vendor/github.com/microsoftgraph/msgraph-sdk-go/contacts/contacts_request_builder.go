package contacts

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i3c25de93fe491a45d1028201d9a97a3867c994bf94910fbb23c75589f3064367 "github.com/microsoftgraph/msgraph-sdk-go/contacts/count"
    i56415aaf9483a2e973df3b1a2fe488e7b2d8844c5d37e84bfa9789655c1de234 "github.com/microsoftgraph/msgraph-sdk-go/contacts/validateproperties"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ibd6b6137a854fdc0082d98413fd3f4475189c1b69ecd839ed41d01e72ef32af8 "github.com/microsoftgraph/msgraph-sdk-go/contacts/getbyids"
    icd7c2e08bbb98a59371e3c25cbc42e41894aefcd3c07a9dda51cde0fbcbda0a1 "github.com/microsoftgraph/msgraph-sdk-go/contacts/getavailableextensionproperties"
    id6a5f8a255008661c4f3891226b4705b4315a990d81bf1bdc10e110af7ce62f9 "github.com/microsoftgraph/msgraph-sdk-go/contacts/delta"
)

// ContactsRequestBuilder provides operations to manage the collection of orgContact entities.
type ContactsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ContactsRequestBuilderGetQueryParameters get the list of organizational contacts for this organization.
type ContactsRequestBuilderGetQueryParameters struct {
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
// ContactsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContactsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ContactsRequestBuilderGetQueryParameters
}
// ContactsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContactsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewContactsRequestBuilderInternal instantiates a new ContactsRequestBuilder and sets the default values.
func NewContactsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContactsRequestBuilder) {
    m := &ContactsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/contacts{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewContactsRequestBuilder instantiates a new ContactsRequestBuilder and sets the default values.
func NewContactsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContactsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewContactsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *ContactsRequestBuilder) Count()(*i3c25de93fe491a45d1028201d9a97a3867c994bf94910fbb23c75589f3064367.CountRequestBuilder) {
    return i3c25de93fe491a45d1028201d9a97a3867c994bf94910fbb23c75589f3064367.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get the list of organizational contacts for this organization.
func (m *ContactsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ContactsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation add new entity to contacts
func (m *ContactsRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, requestConfiguration *ContactsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Delta provides operations to call the delta method.
func (m *ContactsRequestBuilder) Delta()(*id6a5f8a255008661c4f3891226b4705b4315a990d81bf1bdc10e110af7ce62f9.DeltaRequestBuilder) {
    return id6a5f8a255008661c4f3891226b4705b4315a990d81bf1bdc10e110af7ce62f9.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get the list of organizational contacts for this organization.
func (m *ContactsRequestBuilder) Get(ctx context.Context, requestConfiguration *ContactsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOrgContactCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactCollectionResponseable), nil
}
// GetAvailableExtensionProperties provides operations to call the getAvailableExtensionProperties method.
func (m *ContactsRequestBuilder) GetAvailableExtensionProperties()(*icd7c2e08bbb98a59371e3c25cbc42e41894aefcd3c07a9dda51cde0fbcbda0a1.GetAvailableExtensionPropertiesRequestBuilder) {
    return icd7c2e08bbb98a59371e3c25cbc42e41894aefcd3c07a9dda51cde0fbcbda0a1.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *ContactsRequestBuilder) GetByIds()(*ibd6b6137a854fdc0082d98413fd3f4475189c1b69ecd839ed41d01e72ef32af8.GetByIdsRequestBuilder) {
    return ibd6b6137a854fdc0082d98413fd3f4475189c1b69ecd839ed41d01e72ef32af8.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post add new entity to contacts
func (m *ContactsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, requestConfiguration *ContactsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOrgContactFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OrgContactable), nil
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *ContactsRequestBuilder) ValidateProperties()(*i56415aaf9483a2e973df3b1a2fe488e7b2d8844c5d37e84bfa9789655c1de234.ValidatePropertiesRequestBuilder) {
    return i56415aaf9483a2e973df3b1a2fe488e7b2d8844c5d37e84bfa9789655c1de234.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
