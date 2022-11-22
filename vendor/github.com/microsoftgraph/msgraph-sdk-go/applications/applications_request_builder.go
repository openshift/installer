package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i3bc7881f8d53ff46b80580bf2bbeca93ee226b9e556c80e1a29474dc201f2a62 "github.com/microsoftgraph/msgraph-sdk-go/applications/validateproperties"
    i7e8a726e09ced7d4436bf674d1a81669775cb1463962eb805fe901152e21a7a1 "github.com/microsoftgraph/msgraph-sdk-go/applications/count"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    id0455efde9c006c43aa69fbd428002d35001851a48cc95647c179d7f5fe714d8 "github.com/microsoftgraph/msgraph-sdk-go/applications/getavailableextensionproperties"
    id1a6f2211b8fade308e722e65256ccffe715e872816ec9b9e12b585d86dd4f21 "github.com/microsoftgraph/msgraph-sdk-go/applications/getbyids"
    ie526c09037227f7fe6fca4d5c202671bc46cb0934b016d8249c0780baf039d18 "github.com/microsoftgraph/msgraph-sdk-go/applications/delta"
)

// ApplicationsRequestBuilder provides operations to manage the collection of application entities.
type ApplicationsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ApplicationsRequestBuilderGetQueryParameters get the list of applications in this organization.
type ApplicationsRequestBuilderGetQueryParameters struct {
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
// ApplicationsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ApplicationsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ApplicationsRequestBuilderGetQueryParameters
}
// ApplicationsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ApplicationsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewApplicationsRequestBuilderInternal instantiates a new ApplicationsRequestBuilder and sets the default values.
func NewApplicationsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ApplicationsRequestBuilder) {
    m := &ApplicationsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewApplicationsRequestBuilder instantiates a new ApplicationsRequestBuilder and sets the default values.
func NewApplicationsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ApplicationsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewApplicationsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *ApplicationsRequestBuilder) Count()(*i7e8a726e09ced7d4436bf674d1a81669775cb1463962eb805fe901152e21a7a1.CountRequestBuilder) {
    return i7e8a726e09ced7d4436bf674d1a81669775cb1463962eb805fe901152e21a7a1.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get the list of applications in this organization.
func (m *ApplicationsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ApplicationsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create a new application object.
func (m *ApplicationsRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, requestConfiguration *ApplicationsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *ApplicationsRequestBuilder) Delta()(*ie526c09037227f7fe6fca4d5c202671bc46cb0934b016d8249c0780baf039d18.DeltaRequestBuilder) {
    return ie526c09037227f7fe6fca4d5c202671bc46cb0934b016d8249c0780baf039d18.NewDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get the list of applications in this organization.
func (m *ApplicationsRequestBuilder) Get(ctx context.Context, requestConfiguration *ApplicationsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApplicationCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApplicationCollectionResponseable), nil
}
// GetAvailableExtensionProperties provides operations to call the getAvailableExtensionProperties method.
func (m *ApplicationsRequestBuilder) GetAvailableExtensionProperties()(*id0455efde9c006c43aa69fbd428002d35001851a48cc95647c179d7f5fe714d8.GetAvailableExtensionPropertiesRequestBuilder) {
    return id0455efde9c006c43aa69fbd428002d35001851a48cc95647c179d7f5fe714d8.NewGetAvailableExtensionPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetByIds provides operations to call the getByIds method.
func (m *ApplicationsRequestBuilder) GetByIds()(*id1a6f2211b8fade308e722e65256ccffe715e872816ec9b9e12b585d86dd4f21.GetByIdsRequestBuilder) {
    return id1a6f2211b8fade308e722e65256ccffe715e872816ec9b9e12b585d86dd4f21.NewGetByIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create a new application object.
func (m *ApplicationsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, requestConfiguration *ApplicationsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable), nil
}
// ValidateProperties provides operations to call the validateProperties method.
func (m *ApplicationsRequestBuilder) ValidateProperties()(*i3bc7881f8d53ff46b80580bf2bbeca93ee226b9e556c80e1a29474dc201f2a62.ValidatePropertiesRequestBuilder) {
    return i3bc7881f8d53ff46b80580bf2bbeca93ee226b9e556c80e1a29474dc201f2a62.NewValidatePropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
