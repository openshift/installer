package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i221d3f1c8592f7a9636476245f2e4a9e82296b6e5cc50822548a6744c13ec682 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/user"
    i3ff603f1915775aa8fd74380549e8fdf33447260d93eee635198c7c3ce55ef41 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments"
    ibaab995f04768ff7b941456360e066c6fd775e37b26f2c71367a3f6264b2a4c3 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/classes"
    icc9d55c2563db563fec8aa8615c0e6e26a0a9a03820b47d4b01f7ef2f53ee2e6 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/schools"
    if19c2a0949f5f27bde752502b1831a530b8da87a55d1240bdad90f98dc585e1e "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/rubrics"
    ifc2ed629d6ca53bbfc33acc06329c54a103a8e43f74a0e018639ae7bafc6020b "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/taughtclasses"
    i142cc91605bccdf7d1947e2356d3e675fabe42483abf58b911ef3417b99a7658 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/assignments/item"
    i6c43b628928247a4ca6ee2535b4cc6a110e39b8ff8155ce809de4582b7428450 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/taughtclasses/item"
    ia6a8d66cb9ce8a7587ed39000e658b45c988146e3baee26d141c9c4829ae1818 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/rubrics/item"
    ia9b5ce3f1a9b6628a35c3d194c8db39e4f5a678c68df1c7d38765c1e6983ee9f "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/schools/item"
    id83e01290b448c553c9d9d81fdf7d67ecfbcc1916e58827c310107bad3ba0ec7 "github.com/microsoftgraph/msgraph-sdk-go/education/users/item/classes/item"
)

// EducationUserItemRequestBuilder provides operations to manage the users property of the microsoft.graph.educationRoot entity.
type EducationUserItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EducationUserItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationUserItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EducationUserItemRequestBuilderGetQueryParameters get users from education
type EducationUserItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EducationUserItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationUserItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EducationUserItemRequestBuilderGetQueryParameters
}
// EducationUserItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EducationUserItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Assignments provides operations to manage the assignments property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) Assignments()(*i3ff603f1915775aa8fd74380549e8fdf33447260d93eee635198c7c3ce55ef41.AssignmentsRequestBuilder) {
    return i3ff603f1915775aa8fd74380549e8fdf33447260d93eee635198c7c3ce55ef41.NewAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentsById provides operations to manage the assignments property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) AssignmentsById(id string)(*i142cc91605bccdf7d1947e2356d3e675fabe42483abf58b911ef3417b99a7658.EducationAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationAssignment%2Did"] = id
    }
    return i142cc91605bccdf7d1947e2356d3e675fabe42483abf58b911ef3417b99a7658.NewEducationAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Classes provides operations to manage the classes property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) Classes()(*ibaab995f04768ff7b941456360e066c6fd775e37b26f2c71367a3f6264b2a4c3.ClassesRequestBuilder) {
    return ibaab995f04768ff7b941456360e066c6fd775e37b26f2c71367a3f6264b2a4c3.NewClassesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClassesById provides operations to manage the classes property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) ClassesById(id string)(*id83e01290b448c553c9d9d81fdf7d67ecfbcc1916e58827c310107bad3ba0ec7.EducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationClass%2Did"] = id
    }
    return id83e01290b448c553c9d9d81fdf7d67ecfbcc1916e58827c310107bad3ba0ec7.NewEducationClassItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewEducationUserItemRequestBuilderInternal instantiates a new EducationUserItemRequestBuilder and sets the default values.
func NewEducationUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationUserItemRequestBuilder) {
    m := &EducationUserItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education/users/{educationUser%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEducationUserItemRequestBuilder instantiates a new EducationUserItemRequestBuilder and sets the default values.
func NewEducationUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEducationUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property users for education
func (m *EducationUserItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EducationUserItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get users from education
func (m *EducationUserItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EducationUserItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property users in education
func (m *EducationUserItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, requestConfiguration *EducationUserItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property users for education
func (m *EducationUserItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *EducationUserItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get users from education
func (m *EducationUserItemRequestBuilder) Get(ctx context.Context, requestConfiguration *EducationUserItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable), nil
}
// Patch update the navigation property users in education
func (m *EducationUserItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, requestConfiguration *EducationUserItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEducationUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable), nil
}
// Rubrics provides operations to manage the rubrics property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) Rubrics()(*if19c2a0949f5f27bde752502b1831a530b8da87a55d1240bdad90f98dc585e1e.RubricsRequestBuilder) {
    return if19c2a0949f5f27bde752502b1831a530b8da87a55d1240bdad90f98dc585e1e.NewRubricsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RubricsById provides operations to manage the rubrics property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) RubricsById(id string)(*ia6a8d66cb9ce8a7587ed39000e658b45c988146e3baee26d141c9c4829ae1818.EducationRubricItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationRubric%2Did"] = id
    }
    return ia6a8d66cb9ce8a7587ed39000e658b45c988146e3baee26d141c9c4829ae1818.NewEducationRubricItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Schools provides operations to manage the schools property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) Schools()(*icc9d55c2563db563fec8aa8615c0e6e26a0a9a03820b47d4b01f7ef2f53ee2e6.SchoolsRequestBuilder) {
    return icc9d55c2563db563fec8aa8615c0e6e26a0a9a03820b47d4b01f7ef2f53ee2e6.NewSchoolsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchoolsById provides operations to manage the schools property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) SchoolsById(id string)(*ia9b5ce3f1a9b6628a35c3d194c8db39e4f5a678c68df1c7d38765c1e6983ee9f.EducationSchoolItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSchool%2Did"] = id
    }
    return ia9b5ce3f1a9b6628a35c3d194c8db39e4f5a678c68df1c7d38765c1e6983ee9f.NewEducationSchoolItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TaughtClasses provides operations to manage the taughtClasses property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) TaughtClasses()(*ifc2ed629d6ca53bbfc33acc06329c54a103a8e43f74a0e018639ae7bafc6020b.TaughtClassesRequestBuilder) {
    return ifc2ed629d6ca53bbfc33acc06329c54a103a8e43f74a0e018639ae7bafc6020b.NewTaughtClassesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TaughtClassesById provides operations to manage the taughtClasses property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) TaughtClassesById(id string)(*i6c43b628928247a4ca6ee2535b4cc6a110e39b8ff8155ce809de4582b7428450.EducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationClass%2Did"] = id
    }
    return i6c43b628928247a4ca6ee2535b4cc6a110e39b8ff8155ce809de4582b7428450.NewEducationClassItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// User provides operations to manage the user property of the microsoft.graph.educationUser entity.
func (m *EducationUserItemRequestBuilder) User()(*i221d3f1c8592f7a9636476245f2e4a9e82296b6e5cc50822548a6744c13ec682.UserRequestBuilder) {
    return i221d3f1c8592f7a9636476245f2e4a9e82296b6e5cc50822548a6744c13ec682.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
