package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i062d4a86380aad1b2a18fcdf2ae3039afe4fb3340e8c61cd8a0118f64e2961f4 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/exportjobs"
    i0a839f3beac17f684c85ef687fb80c09048b3dcae7c1ecbde1af40e48df57f93 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getdevicemanagementintentsettingsreport"
    i2f4861dbc7fdcc004c73b6f3a193e531630645ae4e3dba4a06598c5ac995763e "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getcompliancepolicynoncompliancereport"
    i3472d4fd9599a18b9522ece83696f0b21985db9b3e0bfea8e80a93d58aafd21c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getnoncompliantdevicesandsettingsreport"
    i412da535edde14d19b0f7a70e3391905918a35901cecec9071168452048c7acc "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getpolicynoncompliancereport"
    i53f2ba02a8f5254d99ef3d211f9756166f93939a36f554a78d1688617cbe7a5f "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getreportfilters"
    i63f7c0aabdb80575f559f70509dbec8f47b990899881ff53424a1b846dd52597 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/gethistoricalreport"
    i69e11ae317e8e814a0b7be5423ee2af053bfe7bd422319ac2db1f93ee18dc6f5 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getsettingnoncompliancereport"
    ia0c903425c40dceef2013d71020f9eecda163cc98d6a70d6643a96e52d28d782 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getpolicynoncompliancesummaryreport"
    ia8bd370622efb0565c68883e987555ee843b6553925fd83251129c0ab1a330bf "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getdevicemanagementintentpersettingcontributingprofiles"
    ib215a0b2b7d6a201fa215cdae39964212be5dea6ccfa14879498271f5d9820d1 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getconfigurationpolicynoncompliancereport"
    ib550761bd4a7b3f5c6fbf730879368bdf6cf1f2e2a17f428ce1aff9b9d56b7bb "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getconfigurationpolicynoncompliancesummaryreport"
    ib8a4803ec3b35355bac5ee1f09413f3cedea9e96c350c5a3b7a72a1ce2329cd8 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getpolicynoncompliancemetadata"
    ic72dce9d0b6163f9d68d5eda111ccd3bd3a8ad43383cda3febb59c81ff2bb55f "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getcachedreport"
    id4ecac28ab783e0d39a17b453b0ef71f4a50bd1409017afcb0ca9e6b503faa03 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getconfigurationsettingnoncompliancereport"
    idf9e4ac66c87266487fcd200b5d2bd9b6e5236a426d072ab78c2b7d2b87f8fcd "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getdevicenoncompliancereport"
    ie35af60f2ad5a0746fb8d0fd9890e6a8832b0296a222da44702a59e069097f08 "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getcompliancesettingnoncompliancereport"
    ie9d9c171f1948a762a7fe58817fdfd08ca15474b245a8d61a3965d150dd4a6df "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/getcompliancepolicynoncompliancesummaryreport"
    i485f9fd192f0985070e6776ebbebebbeab22aaa9df64be7131a3e622b77e1c4e "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement/reports/exportjobs/item"
)

// ReportsRequestBuilder provides operations to manage the reports property of the microsoft.graph.deviceManagement entity.
type ReportsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ReportsRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ReportsRequestBuilderGetQueryParameters reports singleton
type ReportsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ReportsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ReportsRequestBuilderGetQueryParameters
}
// ReportsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewReportsRequestBuilderInternal instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    m := &ReportsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/deviceManagement/reports{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewReportsRequestBuilder instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewReportsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property reports for deviceManagement
func (m *ReportsRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ReportsRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation reports singleton
func (m *ReportsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property reports in deviceManagement
func (m *ReportsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property reports for deviceManagement
func (m *ReportsRequestBuilder) Delete(ctx context.Context, requestConfiguration *ReportsRequestBuilderDeleteRequestConfiguration)(error) {
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
// ExportJobs provides operations to manage the exportJobs property of the microsoft.graph.deviceManagementReports entity.
func (m *ReportsRequestBuilder) ExportJobs()(*i062d4a86380aad1b2a18fcdf2ae3039afe4fb3340e8c61cd8a0118f64e2961f4.ExportJobsRequestBuilder) {
    return i062d4a86380aad1b2a18fcdf2ae3039afe4fb3340e8c61cd8a0118f64e2961f4.NewExportJobsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExportJobsById provides operations to manage the exportJobs property of the microsoft.graph.deviceManagementReports entity.
func (m *ReportsRequestBuilder) ExportJobsById(id string)(*i485f9fd192f0985070e6776ebbebebbeab22aaa9df64be7131a3e622b77e1c4e.DeviceManagementExportJobItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceManagementExportJob%2Did"] = id
    }
    return i485f9fd192f0985070e6776ebbebebbeab22aaa9df64be7131a3e622b77e1c4e.NewDeviceManagementExportJobItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get reports singleton
func (m *ReportsRequestBuilder) Get(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementReportsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable), nil
}
// GetCachedReport provides operations to call the getCachedReport method.
func (m *ReportsRequestBuilder) GetCachedReport()(*ic72dce9d0b6163f9d68d5eda111ccd3bd3a8ad43383cda3febb59c81ff2bb55f.GetCachedReportRequestBuilder) {
    return ic72dce9d0b6163f9d68d5eda111ccd3bd3a8ad43383cda3febb59c81ff2bb55f.NewGetCachedReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetCompliancePolicyNonComplianceReport provides operations to call the getCompliancePolicyNonComplianceReport method.
func (m *ReportsRequestBuilder) GetCompliancePolicyNonComplianceReport()(*i2f4861dbc7fdcc004c73b6f3a193e531630645ae4e3dba4a06598c5ac995763e.GetCompliancePolicyNonComplianceReportRequestBuilder) {
    return i2f4861dbc7fdcc004c73b6f3a193e531630645ae4e3dba4a06598c5ac995763e.NewGetCompliancePolicyNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetCompliancePolicyNonComplianceSummaryReport provides operations to call the getCompliancePolicyNonComplianceSummaryReport method.
func (m *ReportsRequestBuilder) GetCompliancePolicyNonComplianceSummaryReport()(*ie9d9c171f1948a762a7fe58817fdfd08ca15474b245a8d61a3965d150dd4a6df.GetCompliancePolicyNonComplianceSummaryReportRequestBuilder) {
    return ie9d9c171f1948a762a7fe58817fdfd08ca15474b245a8d61a3965d150dd4a6df.NewGetCompliancePolicyNonComplianceSummaryReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetComplianceSettingNonComplianceReport provides operations to call the getComplianceSettingNonComplianceReport method.
func (m *ReportsRequestBuilder) GetComplianceSettingNonComplianceReport()(*ie35af60f2ad5a0746fb8d0fd9890e6a8832b0296a222da44702a59e069097f08.GetComplianceSettingNonComplianceReportRequestBuilder) {
    return ie35af60f2ad5a0746fb8d0fd9890e6a8832b0296a222da44702a59e069097f08.NewGetComplianceSettingNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetConfigurationPolicyNonComplianceReport provides operations to call the getConfigurationPolicyNonComplianceReport method.
func (m *ReportsRequestBuilder) GetConfigurationPolicyNonComplianceReport()(*ib215a0b2b7d6a201fa215cdae39964212be5dea6ccfa14879498271f5d9820d1.GetConfigurationPolicyNonComplianceReportRequestBuilder) {
    return ib215a0b2b7d6a201fa215cdae39964212be5dea6ccfa14879498271f5d9820d1.NewGetConfigurationPolicyNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetConfigurationPolicyNonComplianceSummaryReport provides operations to call the getConfigurationPolicyNonComplianceSummaryReport method.
func (m *ReportsRequestBuilder) GetConfigurationPolicyNonComplianceSummaryReport()(*ib550761bd4a7b3f5c6fbf730879368bdf6cf1f2e2a17f428ce1aff9b9d56b7bb.GetConfigurationPolicyNonComplianceSummaryReportRequestBuilder) {
    return ib550761bd4a7b3f5c6fbf730879368bdf6cf1f2e2a17f428ce1aff9b9d56b7bb.NewGetConfigurationPolicyNonComplianceSummaryReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetConfigurationSettingNonComplianceReport provides operations to call the getConfigurationSettingNonComplianceReport method.
func (m *ReportsRequestBuilder) GetConfigurationSettingNonComplianceReport()(*id4ecac28ab783e0d39a17b453b0ef71f4a50bd1409017afcb0ca9e6b503faa03.GetConfigurationSettingNonComplianceReportRequestBuilder) {
    return id4ecac28ab783e0d39a17b453b0ef71f4a50bd1409017afcb0ca9e6b503faa03.NewGetConfigurationSettingNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetDeviceManagementIntentPerSettingContributingProfiles provides operations to call the getDeviceManagementIntentPerSettingContributingProfiles method.
func (m *ReportsRequestBuilder) GetDeviceManagementIntentPerSettingContributingProfiles()(*ia8bd370622efb0565c68883e987555ee843b6553925fd83251129c0ab1a330bf.GetDeviceManagementIntentPerSettingContributingProfilesRequestBuilder) {
    return ia8bd370622efb0565c68883e987555ee843b6553925fd83251129c0ab1a330bf.NewGetDeviceManagementIntentPerSettingContributingProfilesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetDeviceManagementIntentSettingsReport provides operations to call the getDeviceManagementIntentSettingsReport method.
func (m *ReportsRequestBuilder) GetDeviceManagementIntentSettingsReport()(*i0a839f3beac17f684c85ef687fb80c09048b3dcae7c1ecbde1af40e48df57f93.GetDeviceManagementIntentSettingsReportRequestBuilder) {
    return i0a839f3beac17f684c85ef687fb80c09048b3dcae7c1ecbde1af40e48df57f93.NewGetDeviceManagementIntentSettingsReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetDeviceNonComplianceReport provides operations to call the getDeviceNonComplianceReport method.
func (m *ReportsRequestBuilder) GetDeviceNonComplianceReport()(*idf9e4ac66c87266487fcd200b5d2bd9b6e5236a426d072ab78c2b7d2b87f8fcd.GetDeviceNonComplianceReportRequestBuilder) {
    return idf9e4ac66c87266487fcd200b5d2bd9b6e5236a426d072ab78c2b7d2b87f8fcd.NewGetDeviceNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetHistoricalReport provides operations to call the getHistoricalReport method.
func (m *ReportsRequestBuilder) GetHistoricalReport()(*i63f7c0aabdb80575f559f70509dbec8f47b990899881ff53424a1b846dd52597.GetHistoricalReportRequestBuilder) {
    return i63f7c0aabdb80575f559f70509dbec8f47b990899881ff53424a1b846dd52597.NewGetHistoricalReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetNoncompliantDevicesAndSettingsReport provides operations to call the getNoncompliantDevicesAndSettingsReport method.
func (m *ReportsRequestBuilder) GetNoncompliantDevicesAndSettingsReport()(*i3472d4fd9599a18b9522ece83696f0b21985db9b3e0bfea8e80a93d58aafd21c.GetNoncompliantDevicesAndSettingsReportRequestBuilder) {
    return i3472d4fd9599a18b9522ece83696f0b21985db9b3e0bfea8e80a93d58aafd21c.NewGetNoncompliantDevicesAndSettingsReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetPolicyNonComplianceMetadata provides operations to call the getPolicyNonComplianceMetadata method.
func (m *ReportsRequestBuilder) GetPolicyNonComplianceMetadata()(*ib8a4803ec3b35355bac5ee1f09413f3cedea9e96c350c5a3b7a72a1ce2329cd8.GetPolicyNonComplianceMetadataRequestBuilder) {
    return ib8a4803ec3b35355bac5ee1f09413f3cedea9e96c350c5a3b7a72a1ce2329cd8.NewGetPolicyNonComplianceMetadataRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetPolicyNonComplianceReport provides operations to call the getPolicyNonComplianceReport method.
func (m *ReportsRequestBuilder) GetPolicyNonComplianceReport()(*i412da535edde14d19b0f7a70e3391905918a35901cecec9071168452048c7acc.GetPolicyNonComplianceReportRequestBuilder) {
    return i412da535edde14d19b0f7a70e3391905918a35901cecec9071168452048c7acc.NewGetPolicyNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetPolicyNonComplianceSummaryReport provides operations to call the getPolicyNonComplianceSummaryReport method.
func (m *ReportsRequestBuilder) GetPolicyNonComplianceSummaryReport()(*ia0c903425c40dceef2013d71020f9eecda163cc98d6a70d6643a96e52d28d782.GetPolicyNonComplianceSummaryReportRequestBuilder) {
    return ia0c903425c40dceef2013d71020f9eecda163cc98d6a70d6643a96e52d28d782.NewGetPolicyNonComplianceSummaryReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetReportFilters provides operations to call the getReportFilters method.
func (m *ReportsRequestBuilder) GetReportFilters()(*i53f2ba02a8f5254d99ef3d211f9756166f93939a36f554a78d1688617cbe7a5f.GetReportFiltersRequestBuilder) {
    return i53f2ba02a8f5254d99ef3d211f9756166f93939a36f554a78d1688617cbe7a5f.NewGetReportFiltersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetSettingNonComplianceReport provides operations to call the getSettingNonComplianceReport method.
func (m *ReportsRequestBuilder) GetSettingNonComplianceReport()(*i69e11ae317e8e814a0b7be5423ee2af053bfe7bd422319ac2db1f93ee18dc6f5.GetSettingNonComplianceReportRequestBuilder) {
    return i69e11ae317e8e814a0b7be5423ee2af053bfe7bd422319ac2db1f93ee18dc6f5.NewGetSettingNonComplianceReportRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property reports in deviceManagement
func (m *ReportsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementReportsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementReportsable), nil
}
