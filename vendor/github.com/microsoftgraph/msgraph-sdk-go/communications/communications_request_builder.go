package communications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i3c47b88ab2c0939d08c5ce7f7d03611ddbd6513767887f2f85b50574b229da2e "github.com/microsoftgraph/msgraph-sdk-go/communications/getpresencesbyuserid"
    i7094ab8c1a06bafeeedd3a2746844e0ec2ac0e2e7d2937ef9abdad9734836b3b "github.com/microsoftgraph/msgraph-sdk-go/communications/presences"
    i8ffdea05c8e955294fca9f00741f2c593d8a678f56d226eb2403adf492ba07b0 "github.com/microsoftgraph/msgraph-sdk-go/communications/onlinemeetings"
    i97bbd24556e25a02196df7f369bb3b5f98f59c02b80c2e99a6f1a240704122ec "github.com/microsoftgraph/msgraph-sdk-go/communications/calls"
    ia16cdc2d9b3bc568ec8cc434362cdc775f06d33edc59461753307ae81aa82be2 "github.com/microsoftgraph/msgraph-sdk-go/communications/callrecords"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i18c4da38d0089983229c741b411e29ea20d5089f44a41039b28b9300ac6b2bbc "github.com/microsoftgraph/msgraph-sdk-go/communications/presences/item"
    i477b1ae6bd5bbfab5a3df641da78a8d4e8d13bcb144b2cc255a5232c225c8e02 "github.com/microsoftgraph/msgraph-sdk-go/communications/calls/item"
    i86a5a2d6d1df2d6f16bc44eca70678c6e13a6e96de856db7943e4d10d0e565f9 "github.com/microsoftgraph/msgraph-sdk-go/communications/callrecords/item"
    ia12cb51356bed2ae549e4383f0c613f371754648330c864ba0a94eed03acad34 "github.com/microsoftgraph/msgraph-sdk-go/communications/onlinemeetings/item"
)

// CommunicationsRequestBuilder provides operations to manage the cloudCommunications singleton.
type CommunicationsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// CommunicationsRequestBuilderGetQueryParameters get communications
type CommunicationsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CommunicationsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CommunicationsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CommunicationsRequestBuilderGetQueryParameters
}
// CommunicationsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CommunicationsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CallRecords provides operations to manage the callRecords property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) CallRecords()(*ia16cdc2d9b3bc568ec8cc434362cdc775f06d33edc59461753307ae81aa82be2.CallRecordsRequestBuilder) {
    return ia16cdc2d9b3bc568ec8cc434362cdc775f06d33edc59461753307ae81aa82be2.NewCallRecordsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CallRecordsById provides operations to manage the callRecords property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) CallRecordsById(id string)(*i86a5a2d6d1df2d6f16bc44eca70678c6e13a6e96de856db7943e4d10d0e565f9.CallRecordItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["callRecord%2Did"] = id
    }
    return i86a5a2d6d1df2d6f16bc44eca70678c6e13a6e96de856db7943e4d10d0e565f9.NewCallRecordItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calls provides operations to manage the calls property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) Calls()(*i97bbd24556e25a02196df7f369bb3b5f98f59c02b80c2e99a6f1a240704122ec.CallsRequestBuilder) {
    return i97bbd24556e25a02196df7f369bb3b5f98f59c02b80c2e99a6f1a240704122ec.NewCallsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CallsById provides operations to manage the calls property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) CallsById(id string)(*i477b1ae6bd5bbfab5a3df641da78a8d4e8d13bcb144b2cc255a5232c225c8e02.CallItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["call%2Did"] = id
    }
    return i477b1ae6bd5bbfab5a3df641da78a8d4e8d13bcb144b2cc255a5232c225c8e02.NewCallItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewCommunicationsRequestBuilderInternal instantiates a new CommunicationsRequestBuilder and sets the default values.
func NewCommunicationsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CommunicationsRequestBuilder) {
    m := &CommunicationsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/communications{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewCommunicationsRequestBuilder instantiates a new CommunicationsRequestBuilder and sets the default values.
func NewCommunicationsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CommunicationsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCommunicationsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get communications
func (m *CommunicationsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *CommunicationsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update communications
func (m *CommunicationsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable, requestConfiguration *CommunicationsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get communications
func (m *CommunicationsRequestBuilder) Get(ctx context.Context, requestConfiguration *CommunicationsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCloudCommunicationsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable), nil
}
// GetPresencesByUserId provides operations to call the getPresencesByUserId method.
func (m *CommunicationsRequestBuilder) GetPresencesByUserId()(*i3c47b88ab2c0939d08c5ce7f7d03611ddbd6513767887f2f85b50574b229da2e.GetPresencesByUserIdRequestBuilder) {
    return i3c47b88ab2c0939d08c5ce7f7d03611ddbd6513767887f2f85b50574b229da2e.NewGetPresencesByUserIdRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OnlineMeetings provides operations to manage the onlineMeetings property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) OnlineMeetings()(*i8ffdea05c8e955294fca9f00741f2c593d8a678f56d226eb2403adf492ba07b0.OnlineMeetingsRequestBuilder) {
    return i8ffdea05c8e955294fca9f00741f2c593d8a678f56d226eb2403adf492ba07b0.NewOnlineMeetingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OnlineMeetingsById provides operations to manage the onlineMeetings property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) OnlineMeetingsById(id string)(*ia12cb51356bed2ae549e4383f0c613f371754648330c864ba0a94eed03acad34.OnlineMeetingItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onlineMeeting%2Did"] = id
    }
    return ia12cb51356bed2ae549e4383f0c613f371754648330c864ba0a94eed03acad34.NewOnlineMeetingItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update communications
func (m *CommunicationsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable, requestConfiguration *CommunicationsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCloudCommunicationsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CloudCommunicationsable), nil
}
// Presences provides operations to manage the presences property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) Presences()(*i7094ab8c1a06bafeeedd3a2746844e0ec2ac0e2e7d2937ef9abdad9734836b3b.PresencesRequestBuilder) {
    return i7094ab8c1a06bafeeedd3a2746844e0ec2ac0e2e7d2937ef9abdad9734836b3b.NewPresencesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PresencesById provides operations to manage the presences property of the microsoft.graph.cloudCommunications entity.
func (m *CommunicationsRequestBuilder) PresencesById(id string)(*i18c4da38d0089983229c741b411e29ea20d5089f44a41039b28b9300ac6b2bbc.PresenceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["presence%2Did"] = id
    }
    return i18c4da38d0089983229c741b411e29ea20d5089f44a41039b28b9300ac6b2bbc.NewPresenceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
