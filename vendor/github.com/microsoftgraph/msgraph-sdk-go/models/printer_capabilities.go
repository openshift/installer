package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrinterCapabilities 
type PrinterCapabilities struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // A list of supported bottom margins(in microns) for the printer.
    bottomMargins []int32
    // True if the printer supports collating when printing muliple copies of a multi-page document; false otherwise.
    collation *bool
    // The color modes supported by the printer. Valid values are described in the following table.
    colorModes []PrintColorMode
    // A list of supported content (MIME) types that the printer supports. It is not guaranteed that the Universal Print service supports printing all of these MIME types.
    contentTypes []string
    // The range of copies per job supported by the printer.
    copiesPerJob IntegerRangeable
    // The list of print resolutions in DPI that are supported by the printer.
    dpis []int32
    // The list of duplex modes that are supported by the printer. Valid values are described in the following table.
    duplexModes []PrintDuplexMode
    // The list of feed orientations that are supported by the printer.
    feedOrientations []PrinterFeedOrientation
    // Finishing processes the printer supports for a printed document.
    finishings []PrintFinishing
    // Supported input bins for the printer.
    inputBins []string
    // True if color printing is supported by the printer; false otherwise. Read-only.
    isColorPrintingSupported *bool
    // True if the printer supports printing by page ranges; false otherwise.
    isPageRangeSupported *bool
    // A list of supported left margins(in microns) for the printer.
    leftMargins []int32
    // The media (i.e., paper) colors supported by the printer.
    mediaColors []string
    // The media sizes supported by the printer. Supports standard size names for ISO and ANSI media sizes. Valid values are in the following table.
    mediaSizes []string
    // The media types supported by the printer.
    mediaTypes []string
    // The presentation directions supported by the printer. Supported values are described in the following table.
    multipageLayouts []PrintMultipageLayout
    // The OdataType property
    odataType *string
    // The print orientations supported by the printer. Valid values are described in the following table.
    orientations []PrintOrientation
    // The printer's supported output bins (trays).
    outputBins []string
    // Supported number of Input Pages to impose upon a single Impression.
    pagesPerSheet []int32
    // The print qualities supported by the printer.
    qualities []PrintQuality
    // A list of supported right margins(in microns) for the printer.
    rightMargins []int32
    // Supported print scalings.
    scalings []PrintScaling
    // True if the printer supports scaling PDF pages to match the print media size; false otherwise.
    supportsFitPdfToPage *bool
    // A list of supported top margins(in microns) for the printer.
    topMargins []int32
}
// NewPrinterCapabilities instantiates a new printerCapabilities and sets the default values.
func NewPrinterCapabilities()(*PrinterCapabilities) {
    m := &PrinterCapabilities{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePrinterCapabilitiesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrinterCapabilitiesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrinterCapabilities(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PrinterCapabilities) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBottomMargins gets the bottomMargins property value. A list of supported bottom margins(in microns) for the printer.
func (m *PrinterCapabilities) GetBottomMargins()([]int32) {
    return m.bottomMargins
}
// GetCollation gets the collation property value. True if the printer supports collating when printing muliple copies of a multi-page document; false otherwise.
func (m *PrinterCapabilities) GetCollation()(*bool) {
    return m.collation
}
// GetColorModes gets the colorModes property value. The color modes supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) GetColorModes()([]PrintColorMode) {
    return m.colorModes
}
// GetContentTypes gets the contentTypes property value. A list of supported content (MIME) types that the printer supports. It is not guaranteed that the Universal Print service supports printing all of these MIME types.
func (m *PrinterCapabilities) GetContentTypes()([]string) {
    return m.contentTypes
}
// GetCopiesPerJob gets the copiesPerJob property value. The range of copies per job supported by the printer.
func (m *PrinterCapabilities) GetCopiesPerJob()(IntegerRangeable) {
    return m.copiesPerJob
}
// GetDpis gets the dpis property value. The list of print resolutions in DPI that are supported by the printer.
func (m *PrinterCapabilities) GetDpis()([]int32) {
    return m.dpis
}
// GetDuplexModes gets the duplexModes property value. The list of duplex modes that are supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) GetDuplexModes()([]PrintDuplexMode) {
    return m.duplexModes
}
// GetFeedOrientations gets the feedOrientations property value. The list of feed orientations that are supported by the printer.
func (m *PrinterCapabilities) GetFeedOrientations()([]PrinterFeedOrientation) {
    return m.feedOrientations
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrinterCapabilities) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bottomMargins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetBottomMargins)
    res["collation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCollation)
    res["colorModes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintColorMode , m.SetColorModes)
    res["contentTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetContentTypes)
    res["copiesPerJob"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIntegerRangeFromDiscriminatorValue , m.SetCopiesPerJob)
    res["dpis"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetDpis)
    res["duplexModes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintDuplexMode , m.SetDuplexModes)
    res["feedOrientations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrinterFeedOrientation , m.SetFeedOrientations)
    res["finishings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintFinishing , m.SetFinishings)
    res["inputBins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetInputBins)
    res["isColorPrintingSupported"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsColorPrintingSupported)
    res["isPageRangeSupported"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsPageRangeSupported)
    res["leftMargins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetLeftMargins)
    res["mediaColors"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetMediaColors)
    res["mediaSizes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetMediaSizes)
    res["mediaTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetMediaTypes)
    res["multipageLayouts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintMultipageLayout , m.SetMultipageLayouts)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["orientations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintOrientation , m.SetOrientations)
    res["outputBins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetOutputBins)
    res["pagesPerSheet"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetPagesPerSheet)
    res["qualities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintQuality , m.SetQualities)
    res["rightMargins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetRightMargins)
    res["scalings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintScaling , m.SetScalings)
    res["supportsFitPdfToPage"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSupportsFitPdfToPage)
    res["topMargins"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("int32" , m.SetTopMargins)
    return res
}
// GetFinishings gets the finishings property value. Finishing processes the printer supports for a printed document.
func (m *PrinterCapabilities) GetFinishings()([]PrintFinishing) {
    return m.finishings
}
// GetInputBins gets the inputBins property value. Supported input bins for the printer.
func (m *PrinterCapabilities) GetInputBins()([]string) {
    return m.inputBins
}
// GetIsColorPrintingSupported gets the isColorPrintingSupported property value. True if color printing is supported by the printer; false otherwise. Read-only.
func (m *PrinterCapabilities) GetIsColorPrintingSupported()(*bool) {
    return m.isColorPrintingSupported
}
// GetIsPageRangeSupported gets the isPageRangeSupported property value. True if the printer supports printing by page ranges; false otherwise.
func (m *PrinterCapabilities) GetIsPageRangeSupported()(*bool) {
    return m.isPageRangeSupported
}
// GetLeftMargins gets the leftMargins property value. A list of supported left margins(in microns) for the printer.
func (m *PrinterCapabilities) GetLeftMargins()([]int32) {
    return m.leftMargins
}
// GetMediaColors gets the mediaColors property value. The media (i.e., paper) colors supported by the printer.
func (m *PrinterCapabilities) GetMediaColors()([]string) {
    return m.mediaColors
}
// GetMediaSizes gets the mediaSizes property value. The media sizes supported by the printer. Supports standard size names for ISO and ANSI media sizes. Valid values are in the following table.
func (m *PrinterCapabilities) GetMediaSizes()([]string) {
    return m.mediaSizes
}
// GetMediaTypes gets the mediaTypes property value. The media types supported by the printer.
func (m *PrinterCapabilities) GetMediaTypes()([]string) {
    return m.mediaTypes
}
// GetMultipageLayouts gets the multipageLayouts property value. The presentation directions supported by the printer. Supported values are described in the following table.
func (m *PrinterCapabilities) GetMultipageLayouts()([]PrintMultipageLayout) {
    return m.multipageLayouts
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PrinterCapabilities) GetOdataType()(*string) {
    return m.odataType
}
// GetOrientations gets the orientations property value. The print orientations supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) GetOrientations()([]PrintOrientation) {
    return m.orientations
}
// GetOutputBins gets the outputBins property value. The printer's supported output bins (trays).
func (m *PrinterCapabilities) GetOutputBins()([]string) {
    return m.outputBins
}
// GetPagesPerSheet gets the pagesPerSheet property value. Supported number of Input Pages to impose upon a single Impression.
func (m *PrinterCapabilities) GetPagesPerSheet()([]int32) {
    return m.pagesPerSheet
}
// GetQualities gets the qualities property value. The print qualities supported by the printer.
func (m *PrinterCapabilities) GetQualities()([]PrintQuality) {
    return m.qualities
}
// GetRightMargins gets the rightMargins property value. A list of supported right margins(in microns) for the printer.
func (m *PrinterCapabilities) GetRightMargins()([]int32) {
    return m.rightMargins
}
// GetScalings gets the scalings property value. Supported print scalings.
func (m *PrinterCapabilities) GetScalings()([]PrintScaling) {
    return m.scalings
}
// GetSupportsFitPdfToPage gets the supportsFitPdfToPage property value. True if the printer supports scaling PDF pages to match the print media size; false otherwise.
func (m *PrinterCapabilities) GetSupportsFitPdfToPage()(*bool) {
    return m.supportsFitPdfToPage
}
// GetTopMargins gets the topMargins property value. A list of supported top margins(in microns) for the printer.
func (m *PrinterCapabilities) GetTopMargins()([]int32) {
    return m.topMargins
}
// Serialize serializes information the current object
func (m *PrinterCapabilities) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetBottomMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("bottomMargins", m.GetBottomMargins())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("collation", m.GetCollation())
        if err != nil {
            return err
        }
    }
    if m.GetColorModes() != nil {
        err := writer.WriteCollectionOfStringValues("colorModes", SerializePrintColorMode(m.GetColorModes()))
        if err != nil {
            return err
        }
    }
    if m.GetContentTypes() != nil {
        err := writer.WriteCollectionOfStringValues("contentTypes", m.GetContentTypes())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("copiesPerJob", m.GetCopiesPerJob())
        if err != nil {
            return err
        }
    }
    if m.GetDpis() != nil {
        err := writer.WriteCollectionOfInt32Values("dpis", m.GetDpis())
        if err != nil {
            return err
        }
    }
    if m.GetDuplexModes() != nil {
        err := writer.WriteCollectionOfStringValues("duplexModes", SerializePrintDuplexMode(m.GetDuplexModes()))
        if err != nil {
            return err
        }
    }
    if m.GetFeedOrientations() != nil {
        err := writer.WriteCollectionOfStringValues("feedOrientations", SerializePrinterFeedOrientation(m.GetFeedOrientations()))
        if err != nil {
            return err
        }
    }
    if m.GetFinishings() != nil {
        err := writer.WriteCollectionOfStringValues("finishings", SerializePrintFinishing(m.GetFinishings()))
        if err != nil {
            return err
        }
    }
    if m.GetInputBins() != nil {
        err := writer.WriteCollectionOfStringValues("inputBins", m.GetInputBins())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isColorPrintingSupported", m.GetIsColorPrintingSupported())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isPageRangeSupported", m.GetIsPageRangeSupported())
        if err != nil {
            return err
        }
    }
    if m.GetLeftMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("leftMargins", m.GetLeftMargins())
        if err != nil {
            return err
        }
    }
    if m.GetMediaColors() != nil {
        err := writer.WriteCollectionOfStringValues("mediaColors", m.GetMediaColors())
        if err != nil {
            return err
        }
    }
    if m.GetMediaSizes() != nil {
        err := writer.WriteCollectionOfStringValues("mediaSizes", m.GetMediaSizes())
        if err != nil {
            return err
        }
    }
    if m.GetMediaTypes() != nil {
        err := writer.WriteCollectionOfStringValues("mediaTypes", m.GetMediaTypes())
        if err != nil {
            return err
        }
    }
    if m.GetMultipageLayouts() != nil {
        err := writer.WriteCollectionOfStringValues("multipageLayouts", SerializePrintMultipageLayout(m.GetMultipageLayouts()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetOrientations() != nil {
        err := writer.WriteCollectionOfStringValues("orientations", SerializePrintOrientation(m.GetOrientations()))
        if err != nil {
            return err
        }
    }
    if m.GetOutputBins() != nil {
        err := writer.WriteCollectionOfStringValues("outputBins", m.GetOutputBins())
        if err != nil {
            return err
        }
    }
    if m.GetPagesPerSheet() != nil {
        err := writer.WriteCollectionOfInt32Values("pagesPerSheet", m.GetPagesPerSheet())
        if err != nil {
            return err
        }
    }
    if m.GetQualities() != nil {
        err := writer.WriteCollectionOfStringValues("qualities", SerializePrintQuality(m.GetQualities()))
        if err != nil {
            return err
        }
    }
    if m.GetRightMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("rightMargins", m.GetRightMargins())
        if err != nil {
            return err
        }
    }
    if m.GetScalings() != nil {
        err := writer.WriteCollectionOfStringValues("scalings", SerializePrintScaling(m.GetScalings()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("supportsFitPdfToPage", m.GetSupportsFitPdfToPage())
        if err != nil {
            return err
        }
    }
    if m.GetTopMargins() != nil {
        err := writer.WriteCollectionOfInt32Values("topMargins", m.GetTopMargins())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PrinterCapabilities) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBottomMargins sets the bottomMargins property value. A list of supported bottom margins(in microns) for the printer.
func (m *PrinterCapabilities) SetBottomMargins(value []int32)() {
    m.bottomMargins = value
}
// SetCollation sets the collation property value. True if the printer supports collating when printing muliple copies of a multi-page document; false otherwise.
func (m *PrinterCapabilities) SetCollation(value *bool)() {
    m.collation = value
}
// SetColorModes sets the colorModes property value. The color modes supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetColorModes(value []PrintColorMode)() {
    m.colorModes = value
}
// SetContentTypes sets the contentTypes property value. A list of supported content (MIME) types that the printer supports. It is not guaranteed that the Universal Print service supports printing all of these MIME types.
func (m *PrinterCapabilities) SetContentTypes(value []string)() {
    m.contentTypes = value
}
// SetCopiesPerJob sets the copiesPerJob property value. The range of copies per job supported by the printer.
func (m *PrinterCapabilities) SetCopiesPerJob(value IntegerRangeable)() {
    m.copiesPerJob = value
}
// SetDpis sets the dpis property value. The list of print resolutions in DPI that are supported by the printer.
func (m *PrinterCapabilities) SetDpis(value []int32)() {
    m.dpis = value
}
// SetDuplexModes sets the duplexModes property value. The list of duplex modes that are supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetDuplexModes(value []PrintDuplexMode)() {
    m.duplexModes = value
}
// SetFeedOrientations sets the feedOrientations property value. The list of feed orientations that are supported by the printer.
func (m *PrinterCapabilities) SetFeedOrientations(value []PrinterFeedOrientation)() {
    m.feedOrientations = value
}
// SetFinishings sets the finishings property value. Finishing processes the printer supports for a printed document.
func (m *PrinterCapabilities) SetFinishings(value []PrintFinishing)() {
    m.finishings = value
}
// SetInputBins sets the inputBins property value. Supported input bins for the printer.
func (m *PrinterCapabilities) SetInputBins(value []string)() {
    m.inputBins = value
}
// SetIsColorPrintingSupported sets the isColorPrintingSupported property value. True if color printing is supported by the printer; false otherwise. Read-only.
func (m *PrinterCapabilities) SetIsColorPrintingSupported(value *bool)() {
    m.isColorPrintingSupported = value
}
// SetIsPageRangeSupported sets the isPageRangeSupported property value. True if the printer supports printing by page ranges; false otherwise.
func (m *PrinterCapabilities) SetIsPageRangeSupported(value *bool)() {
    m.isPageRangeSupported = value
}
// SetLeftMargins sets the leftMargins property value. A list of supported left margins(in microns) for the printer.
func (m *PrinterCapabilities) SetLeftMargins(value []int32)() {
    m.leftMargins = value
}
// SetMediaColors sets the mediaColors property value. The media (i.e., paper) colors supported by the printer.
func (m *PrinterCapabilities) SetMediaColors(value []string)() {
    m.mediaColors = value
}
// SetMediaSizes sets the mediaSizes property value. The media sizes supported by the printer. Supports standard size names for ISO and ANSI media sizes. Valid values are in the following table.
func (m *PrinterCapabilities) SetMediaSizes(value []string)() {
    m.mediaSizes = value
}
// SetMediaTypes sets the mediaTypes property value. The media types supported by the printer.
func (m *PrinterCapabilities) SetMediaTypes(value []string)() {
    m.mediaTypes = value
}
// SetMultipageLayouts sets the multipageLayouts property value. The presentation directions supported by the printer. Supported values are described in the following table.
func (m *PrinterCapabilities) SetMultipageLayouts(value []PrintMultipageLayout)() {
    m.multipageLayouts = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrinterCapabilities) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOrientations sets the orientations property value. The print orientations supported by the printer. Valid values are described in the following table.
func (m *PrinterCapabilities) SetOrientations(value []PrintOrientation)() {
    m.orientations = value
}
// SetOutputBins sets the outputBins property value. The printer's supported output bins (trays).
func (m *PrinterCapabilities) SetOutputBins(value []string)() {
    m.outputBins = value
}
// SetPagesPerSheet sets the pagesPerSheet property value. Supported number of Input Pages to impose upon a single Impression.
func (m *PrinterCapabilities) SetPagesPerSheet(value []int32)() {
    m.pagesPerSheet = value
}
// SetQualities sets the qualities property value. The print qualities supported by the printer.
func (m *PrinterCapabilities) SetQualities(value []PrintQuality)() {
    m.qualities = value
}
// SetRightMargins sets the rightMargins property value. A list of supported right margins(in microns) for the printer.
func (m *PrinterCapabilities) SetRightMargins(value []int32)() {
    m.rightMargins = value
}
// SetScalings sets the scalings property value. Supported print scalings.
func (m *PrinterCapabilities) SetScalings(value []PrintScaling)() {
    m.scalings = value
}
// SetSupportsFitPdfToPage sets the supportsFitPdfToPage property value. True if the printer supports scaling PDF pages to match the print media size; false otherwise.
func (m *PrinterCapabilities) SetSupportsFitPdfToPage(value *bool)() {
    m.supportsFitPdfToPage = value
}
// SetTopMargins sets the topMargins property value. A list of supported top margins(in microns) for the printer.
func (m *PrinterCapabilities) SetTopMargins(value []int32)() {
    m.topMargins = value
}
