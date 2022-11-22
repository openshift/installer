package acceptedsenders

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1416a0d3df7d323804da0d19070664abc8bbc637704ec2c7963124726f667c53 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/acceptedsenders/count"
    i71095e1f0098ac4b3abd1b0c194fbaf1aa7f6d49d939ed32f129d5f348a336d2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/acceptedsenders/ref"
)

// AcceptedSendersRequestBuilder provides operations to manage the acceptedSenders property of the microsoft.graph.group entity.
type AcceptedSendersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AcceptedSendersRequestBuilderGetQueryParameters users in the accepted senders list can post to conversations of the group (identified in the GET request URL).Make sure you do not specify the same user or group in the accepted senders and rejected senders lists, otherwise you will get an error.
type AcceptedSendersRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// AcceptedSendersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AcceptedSendersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AcceptedSendersRequestBuilderGetQueryParameters
}
// NewAcceptedSendersRequestBuilderInternal instantiates a new AcceptedSendersRequestBuilder and sets the default values.
func NewAcceptedSendersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AcceptedSendersRequestBuilder) {
    m := &AcceptedSendersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/acceptedSenders{?%24top,%24skip,%24filter,%24count,%24orderby,%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewAcceptedSendersRequestBuilder instantiates a new AcceptedSendersRequestBuilder and sets the default values.
func NewAcceptedSendersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AcceptedSendersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAcceptedSendersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *AcceptedSendersRequestBuilder) Count()(*i1416a0d3df7d323804da0d19070664abc8bbc637704ec2c7963124726f667c53.CountRequestBuilder) {
    return i1416a0d3df7d323804da0d19070664abc8bbc637704ec2c7963124726f667c53.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation users in the accepted senders list can post to conversations of the group (identified in the GET request URL).Make sure you do not specify the same user or group in the accepted senders and rejected senders lists, otherwise you will get an error.
func (m *AcceptedSendersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *AcceptedSendersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get users in the accepted senders list can post to conversations of the group (identified in the GET request URL).Make sure you do not specify the same user or group in the accepted senders and rejected senders lists, otherwise you will get an error.
func (m *AcceptedSendersRequestBuilder) Get(ctx context.Context, requestConfiguration *AcceptedSendersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
// Ref provides operations to manage the collection of group entities.
func (m *AcceptedSendersRequestBuilder) Ref()(*i71095e1f0098ac4b3abd1b0c194fbaf1aa7f6d49d939ed32f129d5f348a336d2.RefRequestBuilder) {
    return i71095e1f0098ac4b3abd1b0c194fbaf1aa7f6d49d939ed32f129d5f348a336d2.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
