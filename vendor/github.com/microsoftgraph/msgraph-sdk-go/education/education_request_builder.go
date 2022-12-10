package education

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i11d420057e9890f34c6f87265f6025f3489c0f78de3b99ece6ba73f225d5a75c "github.com/microsoftgraph/msgraph-sdk-go/education/me"
    i82d352ad9c3de7536dcc103a1a2a276e7beba1b35548191dda9e846d8c21e897 "github.com/microsoftgraph/msgraph-sdk-go/education/users"
    i8803b2d448091478d216c51910e6ed01bbc8c76e5b76968b633118267314ebe9 "github.com/microsoftgraph/msgraph-sdk-go/education/schools"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    if6b46a923414b614133630b614b404c3005a11daa1e05f10bb834755518d40dd "github.com/microsoftgraph/msgraph-sdk-go/education/classes"
    i31a8a42bc05cb93b0a8b6dc7e28fb5a481e3934df6dc18de8217d0fe268d617c "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item"
    i6fc2718a806c86250bd9f9aaba5418413286c98ed345a7dbe988c558836bde8e "github.com/microsoftgraph/msgraph-sdk-go/education/users/item"
    ibcc84b6252b9c075d39d4d8641ef3cf01a04613a5da305d7d48d21c122395065 "github.com/microsoftgraph/msgraph-sdk-go/education/schools/item"
)

// EducationRequestBuilder provides operations to manage the educationRoot singleton.
type EducationRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EducationRequestBuilderGetQueryParameters get education
type EducationRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EducationRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EducationRequestBuilderGetQueryParameters
}
// EducationRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Classes provides operations to manage the classes property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) Classes()(*if6b46a923414b614133630b614b404c3005a11daa1e05f10bb834755518d40dd.ClassesRequestBuilder) {
    return if6b46a923414b614133630b614b404c3005a11daa1e05f10bb834755518d40dd.NewClassesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClassesById provides operations to manage the classes property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) ClassesById(id string)(*i31a8a42bc05cb93b0a8b6dc7e28fb5a481e3934df6dc18de8217d0fe268d617c.EducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationClass%2Did"] = id
    }
    return i31a8a42bc05cb93b0a8b6dc7e28fb5a481e3934df6dc18de8217d0fe268d617c.NewEducationClassItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewEducationRequestBuilderInternal instantiates a new EducationRequestBuilder and sets the default values.
func NewEducationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationRequestBuilder) {
    m := &EducationRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEducationRequestBuilder instantiates a new EducationRequestBuilder and sets the default values.
func NewEducationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEducationRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get education
func (m *EducationRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EducationRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update education
func (m *EducationRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable, requestConfiguration *EducationRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get education
func (m *EducationRequestBuilder) Get(ctx context.Context, requestConfiguration *EducationRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable), nil
}
// Me provides operations to manage the me property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) Me()(*i11d420057e9890f34c6f87265f6025f3489c0f78de3b99ece6ba73f225d5a75c.MeRequestBuilder) {
    return i11d420057e9890f34c6f87265f6025f3489c0f78de3b99ece6ba73f225d5a75c.NewMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update education
func (m *EducationRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable, requestConfiguration *EducationRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationRootable), nil
}
// Schools provides operations to manage the schools property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) Schools()(*i8803b2d448091478d216c51910e6ed01bbc8c76e5b76968b633118267314ebe9.SchoolsRequestBuilder) {
    return i8803b2d448091478d216c51910e6ed01bbc8c76e5b76968b633118267314ebe9.NewSchoolsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchoolsById provides operations to manage the schools property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) SchoolsById(id string)(*ibcc84b6252b9c075d39d4d8641ef3cf01a04613a5da305d7d48d21c122395065.EducationSchoolItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSchool%2Did"] = id
    }
    return ibcc84b6252b9c075d39d4d8641ef3cf01a04613a5da305d7d48d21c122395065.NewEducationSchoolItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Users provides operations to manage the users property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) Users()(*i82d352ad9c3de7536dcc103a1a2a276e7beba1b35548191dda9e846d8c21e897.UsersRequestBuilder) {
    return i82d352ad9c3de7536dcc103a1a2a276e7beba1b35548191dda9e846d8c21e897.NewUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UsersById provides operations to manage the users property of the microsoft.graph.educationRoot entity.
func (m *EducationRequestBuilder) UsersById(id string)(*i6fc2718a806c86250bd9f9aaba5418413286c98ed345a7dbe988c558836bde8e.EducationUserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationUser%2Did"] = id
    }
    return i6fc2718a806c86250bd9f9aaba5418413286c98ed345a7dbe988c558836bde8e.NewEducationUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
