package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookFunctionsRequestBuilder provides operations to manage the functions property of the microsoft.graph.workbook entity.
type ItemItemsItemWorkbookFunctionsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters get functions from drives
type ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters
}
// ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Abs provides operations to call the abs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Abs()(*ItemItemsItemWorkbookFunctionsAbsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAbsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AccrInt provides operations to call the accrInt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AccrInt()(*ItemItemsItemWorkbookFunctionsAccrIntRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAccrIntRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AccrIntM provides operations to call the accrIntM method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AccrIntM()(*ItemItemsItemWorkbookFunctionsAccrIntMRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAccrIntMRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Acos provides operations to call the acos method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acos()(*ItemItemsItemWorkbookFunctionsAcosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcosRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Acosh provides operations to call the acosh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acosh()(*ItemItemsItemWorkbookFunctionsAcoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcoshRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Acot provides operations to call the acot method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acot()(*ItemItemsItemWorkbookFunctionsAcotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcotRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Acoth provides operations to call the acoth method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acoth()(*ItemItemsItemWorkbookFunctionsAcothRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcothRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AmorDegrc provides operations to call the amorDegrc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AmorDegrc()(*ItemItemsItemWorkbookFunctionsAmorDegrcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAmorDegrcRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AmorLinc provides operations to call the amorLinc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AmorLinc()(*ItemItemsItemWorkbookFunctionsAmorLincRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAmorLincRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// And provides operations to call the and method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) And()(*ItemItemsItemWorkbookFunctionsAndRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAndRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Arabic provides operations to call the arabic method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Arabic()(*ItemItemsItemWorkbookFunctionsArabicRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsArabicRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Areas provides operations to call the areas method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Areas()(*ItemItemsItemWorkbookFunctionsAreasRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAreasRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Asc provides operations to call the asc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asc()(*ItemItemsItemWorkbookFunctionsAscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Asin provides operations to call the asin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asin()(*ItemItemsItemWorkbookFunctionsAsinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAsinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Asinh provides operations to call the asinh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asinh()(*ItemItemsItemWorkbookFunctionsAsinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAsinhRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Atan provides operations to call the atan method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atan()(*ItemItemsItemWorkbookFunctionsAtanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Atan2 provides operations to call the atan2 method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atan2()(*ItemItemsItemWorkbookFunctionsAtan2RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtan2RequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Atanh provides operations to call the atanh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atanh()(*ItemItemsItemWorkbookFunctionsAtanhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtanhRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AveDev provides operations to call the aveDev method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AveDev()(*ItemItemsItemWorkbookFunctionsAveDevRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAveDevRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Average provides operations to call the average method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Average()(*ItemItemsItemWorkbookFunctionsAverageRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AverageA provides operations to call the averageA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageA()(*ItemItemsItemWorkbookFunctionsAverageARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AverageIf provides operations to call the averageIf method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageIf()(*ItemItemsItemWorkbookFunctionsAverageIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageIfRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// AverageIfs provides operations to call the averageIfs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageIfs()(*ItemItemsItemWorkbookFunctionsAverageIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageIfsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// BahtText provides operations to call the bahtText method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BahtText()(*ItemItemsItemWorkbookFunctionsBahtTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBahtTextRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Base provides operations to call the base method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Base()(*ItemItemsItemWorkbookFunctionsBaseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBaseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// BesselI provides operations to call the besselI method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselI()(*ItemItemsItemWorkbookFunctionsBesselIRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselIRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// BesselJ provides operations to call the besselJ method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselJ()(*ItemItemsItemWorkbookFunctionsBesselJRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselJRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// BesselK provides operations to call the besselK method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselK()(*ItemItemsItemWorkbookFunctionsBesselKRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselKRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// BesselY provides operations to call the besselY method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselY()(*ItemItemsItemWorkbookFunctionsBesselYRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselYRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Beta_Dist provides operations to call the beta_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Beta_Dist()(*ItemItemsItemWorkbookFunctionsBeta_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBeta_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Beta_Inv provides operations to call the beta_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Beta_Inv()(*ItemItemsItemWorkbookFunctionsBeta_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBeta_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bin2Dec provides operations to call the bin2Dec method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Dec()(*ItemItemsItemWorkbookFunctionsBin2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2DecRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bin2Hex provides operations to call the bin2Hex method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Hex()(*ItemItemsItemWorkbookFunctionsBin2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2HexRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bin2Oct provides operations to call the bin2Oct method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Oct()(*ItemItemsItemWorkbookFunctionsBin2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2OctRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Binom_Dist provides operations to call the binom_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Dist()(*ItemItemsItemWorkbookFunctionsBinom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Binom_Dist_Range provides operations to call the binom_Dist_Range method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Dist_Range()(*ItemItemsItemWorkbookFunctionsBinom_Dist_RangeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_Dist_RangeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Binom_Inv provides operations to call the binom_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Inv()(*ItemItemsItemWorkbookFunctionsBinom_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bitand provides operations to call the bitand method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitand()(*ItemItemsItemWorkbookFunctionsBitandRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitandRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bitlshift provides operations to call the bitlshift method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitlshift()(*ItemItemsItemWorkbookFunctionsBitlshiftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitlshiftRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bitor provides operations to call the bitor method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitor()(*ItemItemsItemWorkbookFunctionsBitorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitorRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bitrshift provides operations to call the bitrshift method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitrshift()(*ItemItemsItemWorkbookFunctionsBitrshiftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitrshiftRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Bitxor provides operations to call the bitxor method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitxor()(*ItemItemsItemWorkbookFunctionsBitxorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitxorRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ceiling_Math provides operations to call the ceiling_Math method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ceiling_Math()(*ItemItemsItemWorkbookFunctionsCeiling_MathRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCeiling_MathRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ceiling_Precise provides operations to call the ceiling_Precise method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ceiling_Precise()(*ItemItemsItemWorkbookFunctionsCeiling_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCeiling_PreciseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Char provides operations to call the char method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Char()(*ItemItemsItemWorkbookFunctionsCharRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCharRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ChiSq_Dist provides operations to call the chiSq_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Dist()(*ItemItemsItemWorkbookFunctionsChiSq_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ChiSq_Dist_RT provides operations to call the chiSq_Dist_RT method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Dist_RT()(*ItemItemsItemWorkbookFunctionsChiSq_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_Dist_RTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ChiSq_Inv provides operations to call the chiSq_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Inv()(*ItemItemsItemWorkbookFunctionsChiSq_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ChiSq_Inv_RT provides operations to call the chiSq_Inv_RT method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Inv_RT()(*ItemItemsItemWorkbookFunctionsChiSq_Inv_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_Inv_RTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Choose provides operations to call the choose method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Choose()(*ItemItemsItemWorkbookFunctionsChooseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChooseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Clean provides operations to call the clean method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Clean()(*ItemItemsItemWorkbookFunctionsCleanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCleanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Code provides operations to call the code method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Code()(*ItemItemsItemWorkbookFunctionsCodeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCodeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Columns provides operations to call the columns method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Columns()(*ItemItemsItemWorkbookFunctionsColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Combin provides operations to call the combin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Combin()(*ItemItemsItemWorkbookFunctionsCombinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCombinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Combina provides operations to call the combina method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Combina()(*ItemItemsItemWorkbookFunctionsCombinaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCombinaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Complex provides operations to call the complex method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Complex()(*ItemItemsItemWorkbookFunctionsComplexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsComplexRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Concatenate provides operations to call the concatenate method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Concatenate()(*ItemItemsItemWorkbookFunctionsConcatenateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConcatenateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Confidence_Norm provides operations to call the confidence_Norm method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Confidence_Norm()(*ItemItemsItemWorkbookFunctionsConfidence_NormRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConfidence_NormRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Confidence_T provides operations to call the confidence_T method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Confidence_T()(*ItemItemsItemWorkbookFunctionsConfidence_TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConfidence_TRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// NewItemItemsItemWorkbookFunctionsRequestBuilderInternal instantiates a new FunctionsRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookFunctionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    m := &ItemItemsItemWorkbookFunctionsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/functions{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams
    m.requestAdapter = requestAdapter
    return m
}
// NewItemItemsItemWorkbookFunctionsRequestBuilder instantiates a new FunctionsRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookFunctionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookFunctionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Convert provides operations to call the convert method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Convert()(*ItemItemsItemWorkbookFunctionsConvertRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConvertRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Cos provides operations to call the cos method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cos()(*ItemItemsItemWorkbookFunctionsCosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCosRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Cosh provides operations to call the cosh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cosh()(*ItemItemsItemWorkbookFunctionsCoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoshRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Cot provides operations to call the cot method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cot()(*ItemItemsItemWorkbookFunctionsCotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCotRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Coth provides operations to call the coth method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Coth()(*ItemItemsItemWorkbookFunctionsCothRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCothRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Count provides operations to call the count method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Count()(*ItemItemsItemWorkbookFunctionsCountRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CountA provides operations to call the countA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountA()(*ItemItemsItemWorkbookFunctionsCountARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CountBlank provides operations to call the countBlank method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountBlank()(*ItemItemsItemWorkbookFunctionsCountBlankRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountBlankRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CountIf provides operations to call the countIf method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountIf()(*ItemItemsItemWorkbookFunctionsCountIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountIfRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CountIfs provides operations to call the countIfs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountIfs()(*ItemItemsItemWorkbookFunctionsCountIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountIfsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupDayBs provides operations to call the coupDayBs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDayBs()(*ItemItemsItemWorkbookFunctionsCoupDayBsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDayBsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupDays provides operations to call the coupDays method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDays()(*ItemItemsItemWorkbookFunctionsCoupDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDaysRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupDaysNc provides operations to call the coupDaysNc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDaysNc()(*ItemItemsItemWorkbookFunctionsCoupDaysNcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDaysNcRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupNcd provides operations to call the coupNcd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupNcd()(*ItemItemsItemWorkbookFunctionsCoupNcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupNcdRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupNum provides operations to call the coupNum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupNum()(*ItemItemsItemWorkbookFunctionsCoupNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupNumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CoupPcd provides operations to call the coupPcd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupPcd()(*ItemItemsItemWorkbookFunctionsCoupPcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupPcdRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Csc provides operations to call the csc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Csc()(*ItemItemsItemWorkbookFunctionsCscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Csch provides operations to call the csch method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Csch()(*ItemItemsItemWorkbookFunctionsCschRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCschRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CumIPmt provides operations to call the cumIPmt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CumIPmt()(*ItemItemsItemWorkbookFunctionsCumIPmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCumIPmtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// CumPrinc provides operations to call the cumPrinc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CumPrinc()(*ItemItemsItemWorkbookFunctionsCumPrincRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCumPrincRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Date provides operations to call the date method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Date()(*ItemItemsItemWorkbookFunctionsDateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Datevalue provides operations to call the datevalue method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Datevalue()(*ItemItemsItemWorkbookFunctionsDatevalueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDatevalueRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Daverage provides operations to call the daverage method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Daverage()(*ItemItemsItemWorkbookFunctionsDaverageRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDaverageRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Day provides operations to call the day method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Day()(*ItemItemsItemWorkbookFunctionsDayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDayRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Days provides operations to call the days method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Days()(*ItemItemsItemWorkbookFunctionsDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDaysRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Days360 provides operations to call the days360 method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Days360()(*ItemItemsItemWorkbookFunctionsDays360RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDays360RequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Db provides operations to call the db method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Db()(*ItemItemsItemWorkbookFunctionsDbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dbcs provides operations to call the dbcs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dbcs()(*ItemItemsItemWorkbookFunctionsDbcsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDbcsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dcount provides operations to call the dcount method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dcount()(*ItemItemsItemWorkbookFunctionsDcountRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDcountRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DcountA provides operations to call the dcountA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DcountA()(*ItemItemsItemWorkbookFunctionsDcountARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDcountARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ddb provides operations to call the ddb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ddb()(*ItemItemsItemWorkbookFunctionsDdbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDdbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dec2Bin provides operations to call the dec2Bin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Bin()(*ItemItemsItemWorkbookFunctionsDec2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2BinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dec2Hex provides operations to call the dec2Hex method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Hex()(*ItemItemsItemWorkbookFunctionsDec2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2HexRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dec2Oct provides operations to call the dec2Oct method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Oct()(*ItemItemsItemWorkbookFunctionsDec2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2OctRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Decimal provides operations to call the decimal method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Decimal()(*ItemItemsItemWorkbookFunctionsDecimalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDecimalRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Degrees provides operations to call the degrees method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Degrees()(*ItemItemsItemWorkbookFunctionsDegreesRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDegreesRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Delete delete navigation property functions for drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Delta provides operations to call the delta method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Delta()(*ItemItemsItemWorkbookFunctionsDeltaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDeltaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DevSq provides operations to call the devSq method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DevSq()(*ItemItemsItemWorkbookFunctionsDevSqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDevSqRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dget provides operations to call the dget method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dget()(*ItemItemsItemWorkbookFunctionsDgetRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDgetRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Disc provides operations to call the disc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Disc()(*ItemItemsItemWorkbookFunctionsDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDiscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dmax provides operations to call the dmax method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dmax()(*ItemItemsItemWorkbookFunctionsDmaxRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDmaxRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dmin provides operations to call the dmin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dmin()(*ItemItemsItemWorkbookFunctionsDminRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDminRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dollar provides operations to call the dollar method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dollar()(*ItemItemsItemWorkbookFunctionsDollarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DollarDe provides operations to call the dollarDe method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DollarDe()(*ItemItemsItemWorkbookFunctionsDollarDeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarDeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DollarFr provides operations to call the dollarFr method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DollarFr()(*ItemItemsItemWorkbookFunctionsDollarFrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarFrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dproduct provides operations to call the dproduct method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dproduct()(*ItemItemsItemWorkbookFunctionsDproductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDproductRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DstDev provides operations to call the dstDev method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DstDev()(*ItemItemsItemWorkbookFunctionsDstDevRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDstDevRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DstDevP provides operations to call the dstDevP method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DstDevP()(*ItemItemsItemWorkbookFunctionsDstDevPRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDstDevPRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dsum provides operations to call the dsum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dsum()(*ItemItemsItemWorkbookFunctionsDsumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDsumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Duration provides operations to call the duration method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Duration()(*ItemItemsItemWorkbookFunctionsDurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDurationRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Dvar provides operations to call the dvar method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dvar()(*ItemItemsItemWorkbookFunctionsDvarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDvarRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// DvarP provides operations to call the dvarP method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DvarP()(*ItemItemsItemWorkbookFunctionsDvarPRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDvarPRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ecma_Ceiling provides operations to call the ecma_Ceiling method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ecma_Ceiling()(*ItemItemsItemWorkbookFunctionsEcma_CeilingRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEcma_CeilingRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Edate provides operations to call the edate method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Edate()(*ItemItemsItemWorkbookFunctionsEdateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEdateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Effect provides operations to call the effect method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Effect()(*ItemItemsItemWorkbookFunctionsEffectRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEffectRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// EoMonth provides operations to call the eoMonth method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) EoMonth()(*ItemItemsItemWorkbookFunctionsEoMonthRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEoMonthRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Erf provides operations to call the erf method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Erf()(*ItemItemsItemWorkbookFunctionsErfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Erf_Precise provides operations to call the erf_Precise method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Erf_Precise()(*ItemItemsItemWorkbookFunctionsErf_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErf_PreciseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ErfC provides operations to call the erfC method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ErfC()(*ItemItemsItemWorkbookFunctionsErfCRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfCRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ErfC_Precise provides operations to call the erfC_Precise method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ErfC_Precise()(*ItemItemsItemWorkbookFunctionsErfC_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfC_PreciseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Error_Type provides operations to call the error_Type method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Error_Type()(*ItemItemsItemWorkbookFunctionsError_TypeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsError_TypeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Even provides operations to call the even method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Even()(*ItemItemsItemWorkbookFunctionsEvenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEvenRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Exact provides operations to call the exact method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Exact()(*ItemItemsItemWorkbookFunctionsExactRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExactRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Exp provides operations to call the exp method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Exp()(*ItemItemsItemWorkbookFunctionsExpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExpRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Expon_Dist provides operations to call the expon_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Expon_Dist()(*ItemItemsItemWorkbookFunctionsExpon_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExpon_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// F_Dist provides operations to call the f_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Dist()(*ItemItemsItemWorkbookFunctionsF_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// F_Dist_RT provides operations to call the f_Dist_RT method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Dist_RT()(*ItemItemsItemWorkbookFunctionsF_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_Dist_RTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// F_Inv provides operations to call the f_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Inv()(*ItemItemsItemWorkbookFunctionsF_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// F_Inv_RT provides operations to call the f_Inv_RT method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Inv_RT()(*ItemItemsItemWorkbookFunctionsF_Inv_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_Inv_RTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Fact provides operations to call the fact method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fact()(*ItemItemsItemWorkbookFunctionsFactRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFactRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// FactDouble provides operations to call the factDouble method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FactDouble()(*ItemItemsItemWorkbookFunctionsFactDoubleRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFactDoubleRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// False provides operations to call the false method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) False()(*ItemItemsItemWorkbookFunctionsFalseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFalseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Find provides operations to call the find method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Find()(*ItemItemsItemWorkbookFunctionsFindRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFindRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// FindB provides operations to call the findB method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FindB()(*ItemItemsItemWorkbookFunctionsFindBRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFindBRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Fisher provides operations to call the fisher method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fisher()(*ItemItemsItemWorkbookFunctionsFisherRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFisherRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// FisherInv provides operations to call the fisherInv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FisherInv()(*ItemItemsItemWorkbookFunctionsFisherInvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFisherInvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Fixed provides operations to call the fixed method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fixed()(*ItemItemsItemWorkbookFunctionsFixedRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFixedRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Floor_Math provides operations to call the floor_Math method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Floor_Math()(*ItemItemsItemWorkbookFunctionsFloor_MathRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFloor_MathRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Floor_Precise provides operations to call the floor_Precise method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Floor_Precise()(*ItemItemsItemWorkbookFunctionsFloor_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFloor_PreciseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Fv provides operations to call the fv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fv()(*ItemItemsItemWorkbookFunctionsFvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Fvschedule provides operations to call the fvschedule method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fvschedule()(*ItemItemsItemWorkbookFunctionsFvscheduleRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFvscheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Gamma provides operations to call the gamma method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma()(*ItemItemsItemWorkbookFunctionsGammaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Gamma_Dist provides operations to call the gamma_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma_Dist()(*ItemItemsItemWorkbookFunctionsGamma_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGamma_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Gamma_Inv provides operations to call the gamma_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma_Inv()(*ItemItemsItemWorkbookFunctionsGamma_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGamma_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// GammaLn provides operations to call the gammaLn method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GammaLn()(*ItemItemsItemWorkbookFunctionsGammaLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaLnRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// GammaLn_Precise provides operations to call the gammaLn_Precise method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GammaLn_Precise()(*ItemItemsItemWorkbookFunctionsGammaLn_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaLn_PreciseRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Gauss provides operations to call the gauss method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gauss()(*ItemItemsItemWorkbookFunctionsGaussRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGaussRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Gcd provides operations to call the gcd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gcd()(*ItemItemsItemWorkbookFunctionsGcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGcdRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// GeoMean provides operations to call the geoMean method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GeoMean()(*ItemItemsItemWorkbookFunctionsGeoMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGeoMeanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// GeStep provides operations to call the geStep method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GeStep()(*ItemItemsItemWorkbookFunctionsGeStepRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGeStepRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Get get functions from drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFunctionsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable), nil
}
// HarMean provides operations to call the harMean method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) HarMean()(*ItemItemsItemWorkbookFunctionsHarMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHarMeanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hex2Bin provides operations to call the hex2Bin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Bin()(*ItemItemsItemWorkbookFunctionsHex2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2BinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hex2Dec provides operations to call the hex2Dec method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Dec()(*ItemItemsItemWorkbookFunctionsHex2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2DecRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hex2Oct provides operations to call the hex2Oct method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Oct()(*ItemItemsItemWorkbookFunctionsHex2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2OctRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hlookup provides operations to call the hlookup method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hlookup()(*ItemItemsItemWorkbookFunctionsHlookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHlookupRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hour provides operations to call the hour method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hour()(*ItemItemsItemWorkbookFunctionsHourRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHourRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Hyperlink provides operations to call the hyperlink method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hyperlink()(*ItemItemsItemWorkbookFunctionsHyperlinkRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHyperlinkRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// HypGeom_Dist provides operations to call the hypGeom_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) HypGeom_Dist()(*ItemItemsItemWorkbookFunctionsHypGeom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHypGeom_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IfEscaped provides operations to call the if method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IfEscaped()(*ItemItemsItemWorkbookFunctionsIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIfRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImAbs provides operations to call the imAbs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImAbs()(*ItemItemsItemWorkbookFunctionsImAbsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImAbsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Imaginary provides operations to call the imaginary method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Imaginary()(*ItemItemsItemWorkbookFunctionsImaginaryRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImaginaryRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImArgument provides operations to call the imArgument method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImArgument()(*ItemItemsItemWorkbookFunctionsImArgumentRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImArgumentRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImConjugate provides operations to call the imConjugate method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImConjugate()(*ItemItemsItemWorkbookFunctionsImConjugateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImConjugateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImCos provides operations to call the imCos method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCos()(*ItemItemsItemWorkbookFunctionsImCosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCosRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImCosh provides operations to call the imCosh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCosh()(*ItemItemsItemWorkbookFunctionsImCoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCoshRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImCot provides operations to call the imCot method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCot()(*ItemItemsItemWorkbookFunctionsImCotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCotRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImCsc provides operations to call the imCsc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCsc()(*ItemItemsItemWorkbookFunctionsImCscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImCsch provides operations to call the imCsch method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCsch()(*ItemItemsItemWorkbookFunctionsImCschRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCschRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImDiv provides operations to call the imDiv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImDiv()(*ItemItemsItemWorkbookFunctionsImDivRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImDivRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImExp provides operations to call the imExp method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImExp()(*ItemItemsItemWorkbookFunctionsImExpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImExpRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImLn provides operations to call the imLn method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLn()(*ItemItemsItemWorkbookFunctionsImLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLnRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImLog10 provides operations to call the imLog10 method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLog10()(*ItemItemsItemWorkbookFunctionsImLog10RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLog10RequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImLog2 provides operations to call the imLog2 method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLog2()(*ItemItemsItemWorkbookFunctionsImLog2RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLog2RequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImPower provides operations to call the imPower method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImPower()(*ItemItemsItemWorkbookFunctionsImPowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImPowerRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImProduct provides operations to call the imProduct method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImProduct()(*ItemItemsItemWorkbookFunctionsImProductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImProductRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImReal provides operations to call the imReal method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImReal()(*ItemItemsItemWorkbookFunctionsImRealRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImRealRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSec provides operations to call the imSec method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSec()(*ItemItemsItemWorkbookFunctionsImSecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSecRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSech provides operations to call the imSech method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSech()(*ItemItemsItemWorkbookFunctionsImSechRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSechRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSin provides operations to call the imSin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSin()(*ItemItemsItemWorkbookFunctionsImSinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSinh provides operations to call the imSinh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSinh()(*ItemItemsItemWorkbookFunctionsImSinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSinhRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSqrt provides operations to call the imSqrt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSqrt()(*ItemItemsItemWorkbookFunctionsImSqrtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSqrtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSub provides operations to call the imSub method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSub()(*ItemItemsItemWorkbookFunctionsImSubRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSubRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImSum provides operations to call the imSum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSum()(*ItemItemsItemWorkbookFunctionsImSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ImTan provides operations to call the imTan method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImTan()(*ItemItemsItemWorkbookFunctionsImTanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImTanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Int provides operations to call the int method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Int()(*ItemItemsItemWorkbookFunctionsIntRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIntRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IntRate provides operations to call the intRate method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IntRate()(*ItemItemsItemWorkbookFunctionsIntRateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIntRateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ipmt provides operations to call the ipmt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ipmt()(*ItemItemsItemWorkbookFunctionsIpmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIpmtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Irr provides operations to call the irr method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Irr()(*ItemItemsItemWorkbookFunctionsIrrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIrrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsErr provides operations to call the isErr method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsErr()(*ItemItemsItemWorkbookFunctionsIsErrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsErrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsError provides operations to call the isError method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsError()(*ItemItemsItemWorkbookFunctionsIsErrorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsErrorRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsEven provides operations to call the isEven method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsEven()(*ItemItemsItemWorkbookFunctionsIsEvenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsEvenRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsFormula provides operations to call the isFormula method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsFormula()(*ItemItemsItemWorkbookFunctionsIsFormulaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsFormulaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsLogical provides operations to call the isLogical method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsLogical()(*ItemItemsItemWorkbookFunctionsIsLogicalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsLogicalRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsNA provides operations to call the isNA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNA()(*ItemItemsItemWorkbookFunctionsIsNARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsNonText provides operations to call the isNonText method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNonText()(*ItemItemsItemWorkbookFunctionsIsNonTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNonTextRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsNumber provides operations to call the isNumber method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNumber()(*ItemItemsItemWorkbookFunctionsIsNumberRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNumberRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Iso_Ceiling provides operations to call the iso_Ceiling method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Iso_Ceiling()(*ItemItemsItemWorkbookFunctionsIso_CeilingRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIso_CeilingRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsOdd provides operations to call the isOdd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsOdd()(*ItemItemsItemWorkbookFunctionsIsOddRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsOddRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsoWeekNum provides operations to call the isoWeekNum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsoWeekNum()(*ItemItemsItemWorkbookFunctionsIsoWeekNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsoWeekNumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ispmt provides operations to call the ispmt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ispmt()(*ItemItemsItemWorkbookFunctionsIspmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIspmtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Isref provides operations to call the isref method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Isref()(*ItemItemsItemWorkbookFunctionsIsrefRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsrefRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// IsText provides operations to call the isText method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsText()(*ItemItemsItemWorkbookFunctionsIsTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsTextRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Kurt provides operations to call the kurt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Kurt()(*ItemItemsItemWorkbookFunctionsKurtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsKurtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Large provides operations to call the large method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Large()(*ItemItemsItemWorkbookFunctionsLargeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLargeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Lcm provides operations to call the lcm method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lcm()(*ItemItemsItemWorkbookFunctionsLcmRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLcmRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Left provides operations to call the left method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Left()(*ItemItemsItemWorkbookFunctionsLeftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLeftRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Leftb provides operations to call the leftb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Leftb()(*ItemItemsItemWorkbookFunctionsLeftbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLeftbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Len provides operations to call the len method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Len()(*ItemItemsItemWorkbookFunctionsLenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLenRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Lenb provides operations to call the lenb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lenb()(*ItemItemsItemWorkbookFunctionsLenbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLenbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ln provides operations to call the ln method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ln()(*ItemItemsItemWorkbookFunctionsLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLnRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Log provides operations to call the log method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Log()(*ItemItemsItemWorkbookFunctionsLogRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Log10 provides operations to call the log10 method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Log10()(*ItemItemsItemWorkbookFunctionsLog10RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLog10RequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// LogNorm_Dist provides operations to call the logNorm_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) LogNorm_Dist()(*ItemItemsItemWorkbookFunctionsLogNorm_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogNorm_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// LogNorm_Inv provides operations to call the logNorm_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) LogNorm_Inv()(*ItemItemsItemWorkbookFunctionsLogNorm_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogNorm_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Lookup provides operations to call the lookup method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lookup()(*ItemItemsItemWorkbookFunctionsLookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLookupRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Lower provides operations to call the lower method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lower()(*ItemItemsItemWorkbookFunctionsLowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLowerRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Match provides operations to call the match method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Match()(*ItemItemsItemWorkbookFunctionsMatchRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMatchRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Max provides operations to call the max method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Max()(*ItemItemsItemWorkbookFunctionsMaxRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMaxRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// MaxA provides operations to call the maxA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MaxA()(*ItemItemsItemWorkbookFunctionsMaxARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMaxARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Mduration provides operations to call the mduration method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mduration()(*ItemItemsItemWorkbookFunctionsMdurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMdurationRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Median provides operations to call the median method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Median()(*ItemItemsItemWorkbookFunctionsMedianRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMedianRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Mid provides operations to call the mid method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mid()(*ItemItemsItemWorkbookFunctionsMidRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMidRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Midb provides operations to call the midb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Midb()(*ItemItemsItemWorkbookFunctionsMidbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMidbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Min provides operations to call the min method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Min()(*ItemItemsItemWorkbookFunctionsMinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// MinA provides operations to call the minA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MinA()(*ItemItemsItemWorkbookFunctionsMinARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Minute provides operations to call the minute method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Minute()(*ItemItemsItemWorkbookFunctionsMinuteRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinuteRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Mirr provides operations to call the mirr method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mirr()(*ItemItemsItemWorkbookFunctionsMirrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMirrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Mod provides operations to call the mod method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mod()(*ItemItemsItemWorkbookFunctionsModRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsModRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Month provides operations to call the month method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Month()(*ItemItemsItemWorkbookFunctionsMonthRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMonthRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Mround provides operations to call the mround method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mround()(*ItemItemsItemWorkbookFunctionsMroundRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMroundRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// MultiNomial provides operations to call the multiNomial method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MultiNomial()(*ItemItemsItemWorkbookFunctionsMultiNomialRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMultiNomialRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// N provides operations to call the n method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) N()(*ItemItemsItemWorkbookFunctionsNRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Na provides operations to call the na method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Na()(*ItemItemsItemWorkbookFunctionsNaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// NegBinom_Dist provides operations to call the negBinom_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NegBinom_Dist()(*ItemItemsItemWorkbookFunctionsNegBinom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNegBinom_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// NetworkDays provides operations to call the networkDays method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NetworkDays()(*ItemItemsItemWorkbookFunctionsNetworkDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNetworkDaysRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// NetworkDays_Intl provides operations to call the networkDays_Intl method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NetworkDays_Intl()(*ItemItemsItemWorkbookFunctionsNetworkDays_IntlRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNetworkDays_IntlRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Nominal provides operations to call the nominal method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Nominal()(*ItemItemsItemWorkbookFunctionsNominalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNominalRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Norm_Dist provides operations to call the norm_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_Dist()(*ItemItemsItemWorkbookFunctionsNorm_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Norm_Inv provides operations to call the norm_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_Inv()(*ItemItemsItemWorkbookFunctionsNorm_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Norm_S_Dist provides operations to call the norm_S_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_S_Dist()(*ItemItemsItemWorkbookFunctionsNorm_S_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_S_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Norm_S_Inv provides operations to call the norm_S_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_S_Inv()(*ItemItemsItemWorkbookFunctionsNorm_S_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_S_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Not provides operations to call the not method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Not()(*ItemItemsItemWorkbookFunctionsNotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNotRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Now provides operations to call the now method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Now()(*ItemItemsItemWorkbookFunctionsNowRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNowRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Nper provides operations to call the nper method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Nper()(*ItemItemsItemWorkbookFunctionsNperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNperRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Npv provides operations to call the npv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Npv()(*ItemItemsItemWorkbookFunctionsNpvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNpvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// NumberValue provides operations to call the numberValue method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NumberValue()(*ItemItemsItemWorkbookFunctionsNumberValueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNumberValueRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Oct2Bin provides operations to call the oct2Bin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Bin()(*ItemItemsItemWorkbookFunctionsOct2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2BinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Oct2Dec provides operations to call the oct2Dec method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Dec()(*ItemItemsItemWorkbookFunctionsOct2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2DecRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Oct2Hex provides operations to call the oct2Hex method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Hex()(*ItemItemsItemWorkbookFunctionsOct2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2HexRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Odd provides operations to call the odd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Odd()(*ItemItemsItemWorkbookFunctionsOddRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// OddFPrice provides operations to call the oddFPrice method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddFPrice()(*ItemItemsItemWorkbookFunctionsOddFPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddFPriceRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// OddFYield provides operations to call the oddFYield method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddFYield()(*ItemItemsItemWorkbookFunctionsOddFYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddFYieldRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// OddLPrice provides operations to call the oddLPrice method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddLPrice()(*ItemItemsItemWorkbookFunctionsOddLPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddLPriceRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// OddLYield provides operations to call the oddLYield method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddLYield()(*ItemItemsItemWorkbookFunctionsOddLYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddLYieldRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Or provides operations to call the or method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Or()(*ItemItemsItemWorkbookFunctionsOrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Patch update the navigation property functions in drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFunctionsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable), nil
}
// Pduration provides operations to call the pduration method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pduration()(*ItemItemsItemWorkbookFunctionsPdurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPdurationRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Percentile_Exc provides operations to call the percentile_Exc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Percentile_Exc()(*ItemItemsItemWorkbookFunctionsPercentile_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentile_ExcRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Percentile_Inc provides operations to call the percentile_Inc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Percentile_Inc()(*ItemItemsItemWorkbookFunctionsPercentile_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentile_IncRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// PercentRank_Exc provides operations to call the percentRank_Exc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PercentRank_Exc()(*ItemItemsItemWorkbookFunctionsPercentRank_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentRank_ExcRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// PercentRank_Inc provides operations to call the percentRank_Inc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PercentRank_Inc()(*ItemItemsItemWorkbookFunctionsPercentRank_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentRank_IncRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Permut provides operations to call the permut method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Permut()(*ItemItemsItemWorkbookFunctionsPermutRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPermutRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Permutationa provides operations to call the permutationa method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Permutationa()(*ItemItemsItemWorkbookFunctionsPermutationaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPermutationaRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Phi provides operations to call the phi method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Phi()(*ItemItemsItemWorkbookFunctionsPhiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPhiRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Pi provides operations to call the pi method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pi()(*ItemItemsItemWorkbookFunctionsPiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPiRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Pmt provides operations to call the pmt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pmt()(*ItemItemsItemWorkbookFunctionsPmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPmtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Poisson_Dist provides operations to call the poisson_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Poisson_Dist()(*ItemItemsItemWorkbookFunctionsPoisson_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPoisson_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Power provides operations to call the power method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Power()(*ItemItemsItemWorkbookFunctionsPowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPowerRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Ppmt provides operations to call the ppmt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ppmt()(*ItemItemsItemWorkbookFunctionsPpmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPpmtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Price provides operations to call the price method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Price()(*ItemItemsItemWorkbookFunctionsPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// PriceDisc provides operations to call the priceDisc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PriceDisc()(*ItemItemsItemWorkbookFunctionsPriceDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceDiscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// PriceMat provides operations to call the priceMat method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PriceMat()(*ItemItemsItemWorkbookFunctionsPriceMatRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceMatRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Product provides operations to call the product method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Product()(*ItemItemsItemWorkbookFunctionsProductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsProductRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Proper provides operations to call the proper method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Proper()(*ItemItemsItemWorkbookFunctionsProperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsProperRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Pv provides operations to call the pv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pv()(*ItemItemsItemWorkbookFunctionsPvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Quartile_Exc provides operations to call the quartile_Exc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quartile_Exc()(*ItemItemsItemWorkbookFunctionsQuartile_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuartile_ExcRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Quartile_Inc provides operations to call the quartile_Inc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quartile_Inc()(*ItemItemsItemWorkbookFunctionsQuartile_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuartile_IncRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Quotient provides operations to call the quotient method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quotient()(*ItemItemsItemWorkbookFunctionsQuotientRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuotientRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Radians provides operations to call the radians method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Radians()(*ItemItemsItemWorkbookFunctionsRadiansRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRadiansRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rand provides operations to call the rand method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rand()(*ItemItemsItemWorkbookFunctionsRandRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRandRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// RandBetween provides operations to call the randBetween method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RandBetween()(*ItemItemsItemWorkbookFunctionsRandBetweenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRandBetweenRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rank_Avg provides operations to call the rank_Avg method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rank_Avg()(*ItemItemsItemWorkbookFunctionsRank_AvgRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRank_AvgRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rank_Eq provides operations to call the rank_Eq method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rank_Eq()(*ItemItemsItemWorkbookFunctionsRank_EqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRank_EqRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rate provides operations to call the rate method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rate()(*ItemItemsItemWorkbookFunctionsRateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRateRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Received provides operations to call the received method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Received()(*ItemItemsItemWorkbookFunctionsReceivedRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReceivedRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Replace provides operations to call the replace method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Replace()(*ItemItemsItemWorkbookFunctionsReplaceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReplaceRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ReplaceB provides operations to call the replaceB method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ReplaceB()(*ItemItemsItemWorkbookFunctionsReplaceBRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReplaceBRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rept provides operations to call the rept method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rept()(*ItemItemsItemWorkbookFunctionsReptRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReptRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Right provides operations to call the right method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Right()(*ItemItemsItemWorkbookFunctionsRightRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRightRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rightb provides operations to call the rightb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rightb()(*ItemItemsItemWorkbookFunctionsRightbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRightbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Roman provides operations to call the roman method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Roman()(*ItemItemsItemWorkbookFunctionsRomanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRomanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Round provides operations to call the round method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Round()(*ItemItemsItemWorkbookFunctionsRoundRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// RoundDown provides operations to call the roundDown method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RoundDown()(*ItemItemsItemWorkbookFunctionsRoundDownRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundDownRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// RoundUp provides operations to call the roundUp method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RoundUp()(*ItemItemsItemWorkbookFunctionsRoundUpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundUpRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rows provides operations to call the rows method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rows()(*ItemItemsItemWorkbookFunctionsRowsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRowsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Rri provides operations to call the rri method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rri()(*ItemItemsItemWorkbookFunctionsRriRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRriRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sec provides operations to call the sec method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sec()(*ItemItemsItemWorkbookFunctionsSecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSecRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sech provides operations to call the sech method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sech()(*ItemItemsItemWorkbookFunctionsSechRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSechRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Second provides operations to call the second method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Second()(*ItemItemsItemWorkbookFunctionsSecondRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSecondRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// SeriesSum provides operations to call the seriesSum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SeriesSum()(*ItemItemsItemWorkbookFunctionsSeriesSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSeriesSumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sheet provides operations to call the sheet method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sheet()(*ItemItemsItemWorkbookFunctionsSheetRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSheetRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sheets provides operations to call the sheets method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sheets()(*ItemItemsItemWorkbookFunctionsSheetsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSheetsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sign provides operations to call the sign method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sign()(*ItemItemsItemWorkbookFunctionsSignRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSignRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sin provides operations to call the sin method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sin()(*ItemItemsItemWorkbookFunctionsSinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSinRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sinh provides operations to call the sinh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sinh()(*ItemItemsItemWorkbookFunctionsSinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSinhRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Skew provides operations to call the skew method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Skew()(*ItemItemsItemWorkbookFunctionsSkewRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSkewRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Skew_p provides operations to call the skew_p method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Skew_p()(*ItemItemsItemWorkbookFunctionsSkew_pRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSkew_pRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sln provides operations to call the sln method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sln()(*ItemItemsItemWorkbookFunctionsSlnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSlnRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Small provides operations to call the small method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Small()(*ItemItemsItemWorkbookFunctionsSmallRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSmallRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sqrt provides operations to call the sqrt method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sqrt()(*ItemItemsItemWorkbookFunctionsSqrtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSqrtRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// SqrtPi provides operations to call the sqrtPi method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SqrtPi()(*ItemItemsItemWorkbookFunctionsSqrtPiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSqrtPiRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Standardize provides operations to call the standardize method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Standardize()(*ItemItemsItemWorkbookFunctionsStandardizeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStandardizeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// StDev_P provides operations to call the stDev_P method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDev_P()(*ItemItemsItemWorkbookFunctionsStDev_PRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDev_PRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// StDev_S provides operations to call the stDev_S method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDev_S()(*ItemItemsItemWorkbookFunctionsStDev_SRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDev_SRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// StDevA provides operations to call the stDevA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDevA()(*ItemItemsItemWorkbookFunctionsStDevARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDevARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// StDevPA provides operations to call the stDevPA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDevPA()(*ItemItemsItemWorkbookFunctionsStDevPARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDevPARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Substitute provides operations to call the substitute method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Substitute()(*ItemItemsItemWorkbookFunctionsSubstituteRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSubstituteRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Subtotal provides operations to call the subtotal method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Subtotal()(*ItemItemsItemWorkbookFunctionsSubtotalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSubtotalRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Sum provides operations to call the sum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sum()(*ItemItemsItemWorkbookFunctionsSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// SumIf provides operations to call the sumIf method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumIf()(*ItemItemsItemWorkbookFunctionsSumIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumIfRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// SumIfs provides operations to call the sumIfs method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumIfs()(*ItemItemsItemWorkbookFunctionsSumIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumIfsRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// SumSq provides operations to call the sumSq method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumSq()(*ItemItemsItemWorkbookFunctionsSumSqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumSqRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Syd provides operations to call the syd method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Syd()(*ItemItemsItemWorkbookFunctionsSydRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSydRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T provides operations to call the t method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T()(*ItemItemsItemWorkbookFunctionsTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T_Dist provides operations to call the t_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist()(*ItemItemsItemWorkbookFunctionsT_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T_Dist_2T provides operations to call the t_Dist_2T method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist_2T()(*ItemItemsItemWorkbookFunctionsT_Dist_2TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Dist_2TRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T_Dist_RT provides operations to call the t_Dist_RT method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist_RT()(*ItemItemsItemWorkbookFunctionsT_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Dist_RTRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T_Inv provides operations to call the t_Inv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Inv()(*ItemItemsItemWorkbookFunctionsT_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_InvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// T_Inv_2T provides operations to call the t_Inv_2T method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Inv_2T()(*ItemItemsItemWorkbookFunctionsT_Inv_2TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Inv_2TRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Tan provides operations to call the tan method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Tan()(*ItemItemsItemWorkbookFunctionsTanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Tanh provides operations to call the tanh method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Tanh()(*ItemItemsItemWorkbookFunctionsTanhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTanhRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// TbillEq provides operations to call the tbillEq method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillEq()(*ItemItemsItemWorkbookFunctionsTbillEqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillEqRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// TbillPrice provides operations to call the tbillPrice method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillPrice()(*ItemItemsItemWorkbookFunctionsTbillPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillPriceRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// TbillYield provides operations to call the tbillYield method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillYield()(*ItemItemsItemWorkbookFunctionsTbillYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillYieldRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Text provides operations to call the text method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Text()(*ItemItemsItemWorkbookFunctionsTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTextRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Time provides operations to call the time method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Time()(*ItemItemsItemWorkbookFunctionsTimeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTimeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Timevalue provides operations to call the timevalue method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Timevalue()(*ItemItemsItemWorkbookFunctionsTimevalueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTimevalueRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Today provides operations to call the today method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Today()(*ItemItemsItemWorkbookFunctionsTodayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTodayRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// ToDeleteRequestInformation delete navigation property functions for drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// ToGetRequestInformation get functions from drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET
    requestInfo.Headers.Add("Accept", "application/json")
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// ToPatchRequestInformation update the navigation property functions in drives
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH
    requestInfo.Headers.Add("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Trim provides operations to call the trim method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Trim()(*ItemItemsItemWorkbookFunctionsTrimRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrimRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// TrimMean provides operations to call the trimMean method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TrimMean()(*ItemItemsItemWorkbookFunctionsTrimMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrimMeanRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// True provides operations to call the true method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) True()(*ItemItemsItemWorkbookFunctionsTrueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrueRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Trunc provides operations to call the trunc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Trunc()(*ItemItemsItemWorkbookFunctionsTruncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTruncRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// TypeEscaped provides operations to call the type method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TypeEscaped()(*ItemItemsItemWorkbookFunctionsTypeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTypeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Unichar provides operations to call the unichar method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Unichar()(*ItemItemsItemWorkbookFunctionsUnicharRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUnicharRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Unicode provides operations to call the unicode method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Unicode()(*ItemItemsItemWorkbookFunctionsUnicodeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUnicodeRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Upper provides operations to call the upper method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Upper()(*ItemItemsItemWorkbookFunctionsUpperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUpperRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Usdollar provides operations to call the usdollar method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Usdollar()(*ItemItemsItemWorkbookFunctionsUsdollarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUsdollarRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Value provides operations to call the value method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Value()(*ItemItemsItemWorkbookFunctionsValueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsValueRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Var_P provides operations to call the var_P method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Var_P()(*ItemItemsItemWorkbookFunctionsVar_PRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVar_PRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Var_S provides operations to call the var_S method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Var_S()(*ItemItemsItemWorkbookFunctionsVar_SRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVar_SRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// VarA provides operations to call the varA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) VarA()(*ItemItemsItemWorkbookFunctionsVarARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVarARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// VarPA provides operations to call the varPA method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) VarPA()(*ItemItemsItemWorkbookFunctionsVarPARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVarPARequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Vdb provides operations to call the vdb method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Vdb()(*ItemItemsItemWorkbookFunctionsVdbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVdbRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Vlookup provides operations to call the vlookup method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Vlookup()(*ItemItemsItemWorkbookFunctionsVlookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVlookupRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Weekday provides operations to call the weekday method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Weekday()(*ItemItemsItemWorkbookFunctionsWeekdayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeekdayRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// WeekNum provides operations to call the weekNum method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WeekNum()(*ItemItemsItemWorkbookFunctionsWeekNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeekNumRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Weibull_Dist provides operations to call the weibull_Dist method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Weibull_Dist()(*ItemItemsItemWorkbookFunctionsWeibull_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeibull_DistRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// WorkDay provides operations to call the workDay method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WorkDay()(*ItemItemsItemWorkbookFunctionsWorkDayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWorkDayRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// WorkDay_Intl provides operations to call the workDay_Intl method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WorkDay_Intl()(*ItemItemsItemWorkbookFunctionsWorkDay_IntlRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWorkDay_IntlRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Xirr provides operations to call the xirr method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xirr()(*ItemItemsItemWorkbookFunctionsXirrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXirrRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Xnpv provides operations to call the xnpv method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xnpv()(*ItemItemsItemWorkbookFunctionsXnpvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXnpvRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Xor provides operations to call the xor method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xor()(*ItemItemsItemWorkbookFunctionsXorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXorRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Year provides operations to call the year method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Year()(*ItemItemsItemWorkbookFunctionsYearRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYearRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// YearFrac provides operations to call the yearFrac method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YearFrac()(*ItemItemsItemWorkbookFunctionsYearFracRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYearFracRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Yield provides operations to call the yield method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Yield()(*ItemItemsItemWorkbookFunctionsYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// YieldDisc provides operations to call the yieldDisc method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YieldDisc()(*ItemItemsItemWorkbookFunctionsYieldDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldDiscRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// YieldMat provides operations to call the yieldMat method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YieldMat()(*ItemItemsItemWorkbookFunctionsYieldMatRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldMatRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
// Z_Test provides operations to call the z_Test method.
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Z_Test()(*ItemItemsItemWorkbookFunctionsZ_TestRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsZ_TestRequestBuilderInternal(m.pathParameters, m.requestAdapter)
}
