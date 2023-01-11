package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ibb195e583f4b6a598781702d4590e05f04fe6498c9bcdad35c61b13587d2a5df "github.com/microsoftgraph/msgraph-sdk-go/education/classes/item/assignments/item/categories/item/ref"
)

// EducationCategoryItemRequestBuilder builds and executes requests for operations under \education\classes\{educationClass-id}\assignments\{educationAssignment-id}\categories\{educationCategory-id}
type EducationCategoryItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// NewEducationCategoryItemRequestBuilderInternal instantiates a new EducationCategoryItemRequestBuilder and sets the default values.
func NewEducationCategoryItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationCategoryItemRequestBuilder) {
    m := &EducationCategoryItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/education/classes/{educationClass%2Did}/assignments/{educationAssignment%2Did}/categories/{educationCategory%2Did}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEducationCategoryItemRequestBuilder instantiates a new EducationCategoryItemRequestBuilder and sets the default values.
func NewEducationCategoryItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EducationCategoryItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEducationCategoryItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
func (m *EducationCategoryItemRequestBuilder) Ref()(*ibb195e583f4b6a598781702d4590e05f04fe6498c9bcdad35c61b13587d2a5df.RefRequestBuilder) {
    return ibb195e583f4b6a598781702d4590e05f04fe6498c9bcdad35c61b13587d2a5df.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
