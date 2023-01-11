package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i13bd20348dd8380732465deeb248cdc3f18bf1b888ea708c2518379965331b7d "github.com/microsoftgraph/msgraph-sdk-go/security/securescores"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ic55f4c3693b3458d37f169972d41697131ecfdb31a625c4fae6e024f02167222 "github.com/microsoftgraph/msgraph-sdk-go/security/alerts"
    ie4328a6e5437e25d06c50bee3076698796ad30481e9e114d704c631aad9b6f5b "github.com/microsoftgraph/msgraph-sdk-go/security/securescorecontrolprofiles"
    ie7d038ee19b31eb009a453ff789ed94976aee3a535728cdf937b5d42a5963d93 "github.com/microsoftgraph/msgraph-sdk-go/security/attacksimulation"
    ifd8140bc5c892f85e2edd387a9c4f1aa051c63a62f5cbae31abf9d061b8b7b59 "github.com/microsoftgraph/msgraph-sdk-go/security/cases"
    i341438a13fe09ad4284380c5f498c991d5debd415f26d3ad4445315402edc168 "github.com/microsoftgraph/msgraph-sdk-go/security/securescorecontrolprofiles/item"
    i935eca4ee55c30c68a1ba6558e201935fcaf58fd82c8f8076c82de8f806aac5b "github.com/microsoftgraph/msgraph-sdk-go/security/securescores/item"
    i9f4cc7fe9832e758c23f51e8095ff0946d511066f5a7d110c1662df7efadb000 "github.com/microsoftgraph/msgraph-sdk-go/security/alerts/item"
)

// SecurityRequestBuilder provides operations to manage the security singleton.
type SecurityRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// SecurityRequestBuilderGetQueryParameters get security
type SecurityRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SecurityRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SecurityRequestBuilderGetQueryParameters
}
// SecurityRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Alerts provides operations to manage the alerts property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) Alerts()(*ic55f4c3693b3458d37f169972d41697131ecfdb31a625c4fae6e024f02167222.AlertsRequestBuilder) {
    return ic55f4c3693b3458d37f169972d41697131ecfdb31a625c4fae6e024f02167222.NewAlertsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AlertsById provides operations to manage the alerts property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) AlertsById(id string)(*i9f4cc7fe9832e758c23f51e8095ff0946d511066f5a7d110c1662df7efadb000.AlertItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["alert%2Did"] = id
    }
    return i9f4cc7fe9832e758c23f51e8095ff0946d511066f5a7d110c1662df7efadb000.NewAlertItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AttackSimulation provides operations to manage the attackSimulation property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) AttackSimulation()(*ie7d038ee19b31eb009a453ff789ed94976aee3a535728cdf937b5d42a5963d93.AttackSimulationRequestBuilder) {
    return ie7d038ee19b31eb009a453ff789ed94976aee3a535728cdf937b5d42a5963d93.NewAttackSimulationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Cases provides operations to manage the cases property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) Cases()(*ifd8140bc5c892f85e2edd387a9c4f1aa051c63a62f5cbae31abf9d061b8b7b59.CasesRequestBuilder) {
    return ifd8140bc5c892f85e2edd387a9c4f1aa051c63a62f5cbae31abf9d061b8b7b59.NewCasesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewSecurityRequestBuilderInternal instantiates a new SecurityRequestBuilder and sets the default values.
func NewSecurityRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityRequestBuilder) {
    m := &SecurityRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/security{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewSecurityRequestBuilder instantiates a new SecurityRequestBuilder and sets the default values.
func NewSecurityRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSecurityRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get security
func (m *SecurityRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *SecurityRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update security
func (m *SecurityRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, requestConfiguration *SecurityRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Get get security
func (m *SecurityRequestBuilder) Get(ctx context.Context, requestConfiguration *SecurityRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecurityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable), nil
}
// Patch update security
func (m *SecurityRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, requestConfiguration *SecurityRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecurityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable), nil
}
// SecureScoreControlProfiles provides operations to manage the secureScoreControlProfiles property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) SecureScoreControlProfiles()(*ie4328a6e5437e25d06c50bee3076698796ad30481e9e114d704c631aad9b6f5b.SecureScoreControlProfilesRequestBuilder) {
    return ie4328a6e5437e25d06c50bee3076698796ad30481e9e114d704c631aad9b6f5b.NewSecureScoreControlProfilesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SecureScoreControlProfilesById provides operations to manage the secureScoreControlProfiles property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) SecureScoreControlProfilesById(id string)(*i341438a13fe09ad4284380c5f498c991d5debd415f26d3ad4445315402edc168.SecureScoreControlProfileItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["secureScoreControlProfile%2Did"] = id
    }
    return i341438a13fe09ad4284380c5f498c991d5debd415f26d3ad4445315402edc168.NewSecureScoreControlProfileItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SecureScores provides operations to manage the secureScores property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) SecureScores()(*i13bd20348dd8380732465deeb248cdc3f18bf1b888ea708c2518379965331b7d.SecureScoresRequestBuilder) {
    return i13bd20348dd8380732465deeb248cdc3f18bf1b888ea708c2518379965331b7d.NewSecureScoresRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SecureScoresById provides operations to manage the secureScores property of the microsoft.graph.security entity.
func (m *SecurityRequestBuilder) SecureScoresById(id string)(*i935eca4ee55c30c68a1ba6558e201935fcaf58fd82c8f8076c82de8f806aac5b.SecureScoreItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["secureScore%2Did"] = id
    }
    return i935eca4ee55c30c68a1ba6558e201935fcaf58fd82c8f8076c82de8f806aac5b.NewSecureScoreItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
