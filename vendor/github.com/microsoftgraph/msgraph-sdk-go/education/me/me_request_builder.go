package me

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i4827ba32335285b89a92c436cf648bf167912072becfb2cf9792b7e694bfe68c "github.com/microsoftgraph/msgraph-sdk-go/education/me/taughtclasses"
    i75d26f07b5741abbb041934b7ff07c14802e8373402b2fe0175141ce976d99f1 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments"
    i8a6dd866ad5ea518840cdc5c8b7348b010903fad7cccc3473ccbc41bd0f2f5c5 "github.com/microsoftgraph/msgraph-sdk-go/education/me/rubrics"
    idc612964c10975ba0294061975f8fd346f82666fc4c4a4245e1b4ea27187f3e0 "github.com/microsoftgraph/msgraph-sdk-go/education/me/classes"
    if174dba299b5249c431f058ba35f021f7837f7d6553d65eaa9c8f39e0c8d2f25 "github.com/microsoftgraph/msgraph-sdk-go/education/me/user"
    if8f2e52b9fa261cdafc93406e1f0deaaa9c3444674bb9bac26886d2161875b88 "github.com/microsoftgraph/msgraph-sdk-go/education/me/schools"
    i05b947f4d440080ac6b413aae5547627313333329bb760c803e6c27f60db5c06 "github.com/microsoftgraph/msgraph-sdk-go/education/me/assignments/item"
    i434c17782aff494ea747f7c3919ee493165cf0cebf127d5797e6c3c1b044f21e "github.com/microsoftgraph/msgraph-sdk-go/education/me/schools/item"
    i4a61e633c19f80889fc1151239c34c623c4ada576c2613365f8634b2142f493d "github.com/microsoftgraph/msgraph-sdk-go/education/me/taughtclasses/item"
    ib528428124f1c371970dddb56f26c3bb4c42e706f192473c9d8d84f45c312a67 "github.com/microsoftgraph/msgraph-sdk-go/education/me/classes/item"
    ic96595483f91910c603469ce138a672deaf059c7a8c295304a9c2acf2ce6134b "github.com/microsoftgraph/msgraph-sdk-go/education/me/rubrics/item"
)

// MeRequestBuilder provides operations to manage the me property of the microsoft.graph.educationRoot entity.
type MeRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MeRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// MeRequestBuilderGetQueryParameters get me from education
type MeRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// MeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MeRequestBuilderGetQueryParameters
}
// MeRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MeRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Assignments provides operations to manage the assignments property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) Assignments()(*i75d26f07b5741abbb041934b7ff07c14802e8373402b2fe0175141ce976d99f1.AssignmentsRequestBuilder) {
    return i75d26f07b5741abbb041934b7ff07c14802e8373402b2fe0175141ce976d99f1.NewAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentsById provides operations to manage the assignments property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) AssignmentsById(id string)(*i05b947f4d440080ac6b413aae5547627313333329bb760c803e6c27f60db5c06.EducationAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationAssignment%2Did"] = id
    }
    return i05b947f4d440080ac6b413aae5547627313333329bb760c803e6c27f60db5c06.NewEducationAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Classes provides operations to manage the classes property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) Classes()(*idc612964c10975ba0294061975f8fd346f82666fc4c4a4245e1b4ea27187f3e0.ClassesRequestBuilder) {
    return idc612964c10975ba0294061975f8fd346f82666fc4c4a4245e1b4ea27187f3e0.NewClassesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ClassesById provides operations to manage the classes property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) ClassesById(id string)(*ib528428124f1c371970dddb56f26c3bb4c42e706f192473c9d8d84f45c312a67.EducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationClass%2Did"] = id
    }
    return ib528428124f1c371970dddb56f26c3bb4c42e706f192473c9d8d84f45c312a67.NewEducationClassItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMeRequestBuilderInternal instantiates a new MeRequestBuilder and sets the default values.
func NewMeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeRequestBuilder) {
    m := &MeRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education/me{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMeRequestBuilder instantiates a new MeRequestBuilder and sets the default values.
func NewMeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property me for education
func (m *MeRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *MeRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get me from education
func (m *MeRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property me in education
func (m *MeRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, requestConfiguration *MeRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property me for education
func (m *MeRequestBuilder) Delete(ctx context.Context, requestConfiguration *MeRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get me from education
func (m *MeRequestBuilder) Get(ctx context.Context, requestConfiguration *MeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, error) {
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
// Patch update the navigation property me in education
func (m *MeRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, requestConfiguration *MeRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EducationUserable, error) {
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
func (m *MeRequestBuilder) Rubrics()(*i8a6dd866ad5ea518840cdc5c8b7348b010903fad7cccc3473ccbc41bd0f2f5c5.RubricsRequestBuilder) {
    return i8a6dd866ad5ea518840cdc5c8b7348b010903fad7cccc3473ccbc41bd0f2f5c5.NewRubricsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RubricsById provides operations to manage the rubrics property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) RubricsById(id string)(*ic96595483f91910c603469ce138a672deaf059c7a8c295304a9c2acf2ce6134b.EducationRubricItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationRubric%2Did"] = id
    }
    return ic96595483f91910c603469ce138a672deaf059c7a8c295304a9c2acf2ce6134b.NewEducationRubricItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Schools provides operations to manage the schools property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) Schools()(*if8f2e52b9fa261cdafc93406e1f0deaaa9c3444674bb9bac26886d2161875b88.SchoolsRequestBuilder) {
    return if8f2e52b9fa261cdafc93406e1f0deaaa9c3444674bb9bac26886d2161875b88.NewSchoolsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchoolsById provides operations to manage the schools property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) SchoolsById(id string)(*i434c17782aff494ea747f7c3919ee493165cf0cebf127d5797e6c3c1b044f21e.EducationSchoolItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationSchool%2Did"] = id
    }
    return i434c17782aff494ea747f7c3919ee493165cf0cebf127d5797e6c3c1b044f21e.NewEducationSchoolItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TaughtClasses provides operations to manage the taughtClasses property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) TaughtClasses()(*i4827ba32335285b89a92c436cf648bf167912072becfb2cf9792b7e694bfe68c.TaughtClassesRequestBuilder) {
    return i4827ba32335285b89a92c436cf648bf167912072becfb2cf9792b7e694bfe68c.NewTaughtClassesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TaughtClassesById provides operations to manage the taughtClasses property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) TaughtClassesById(id string)(*i4a61e633c19f80889fc1151239c34c623c4ada576c2613365f8634b2142f493d.EducationClassItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["educationClass%2Did"] = id
    }
    return i4a61e633c19f80889fc1151239c34c623c4ada576c2613365f8634b2142f493d.NewEducationClassItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// User provides operations to manage the user property of the microsoft.graph.educationUser entity.
func (m *MeRequestBuilder) User()(*if174dba299b5249c431f058ba35f021f7837f7d6553d65eaa9c8f39e0c8d2f25.UserRequestBuilder) {
    return if174dba299b5249c431f058ba35f021f7837f7d6553d65eaa9c8f39e0c8d2f25.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
