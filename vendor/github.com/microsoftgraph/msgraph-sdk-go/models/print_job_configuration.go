package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrintJobConfiguration 
type PrintJobConfiguration struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Whether the printer should collate pages wehen printing multiple copies of a multi-page document.
    collate *bool
    // The color mode the printer should use to print the job. Valid values are described in the table below. Read-only.
    colorMode *PrintColorMode
    // The number of copies that should be printed. Read-only.
    copies *int32
    // The resolution to use when printing the job, expressed in dots per inch (DPI). Read-only.
    dpi *int32
    // The duplex mode the printer should use when printing the job. Valid values are described in the table below. Read-only.
    duplexMode *PrintDuplexMode
    // The orientation to use when feeding media into the printer. Valid values are described in the following table. Read-only.
    feedOrientation *PrinterFeedOrientation
    // Finishing processes to use when printing.
    finishings []PrintFinishing
    // The fitPdfToPage property
    fitPdfToPage *bool
    // The input bin (tray) to use when printing. See the printer's capabilities for a list of supported input bins.
    inputBin *string
    // The margin settings to use when printing.
    margin PrintMarginable
    // The media size to use when printing. Supports standard size names for ISO and ANSI media sizes. Valid values listed in the printerCapabilities topic.
    mediaSize *string
    // The default media (such as paper) type to print the document on.
    mediaType *string
    // The direction to lay out pages when multiple pages are being printed per sheet. Valid values are described in the following table.
    multipageLayout *PrintMultipageLayout
    // The OdataType property
    odataType *string
    // The orientation setting the printer should use when printing the job. Valid values are described in the following table.
    orientation *PrintOrientation
    // The output bin to place completed prints into. See the printer's capabilities for a list of supported output bins.
    outputBin *string
    // The page ranges to print. Read-only.
    pageRanges []IntegerRangeable
    // The number of document pages to print on each sheet.
    pagesPerSheet *int32
    // The print quality to use when printing the job. Valid values are described in the table below. Read-only.
    quality *PrintQuality
    // Specifies how the printer should scale the document data to fit the requested media. Valid values are described in the following table.
    scaling *PrintScaling
}
// NewPrintJobConfiguration instantiates a new printJobConfiguration and sets the default values.
func NewPrintJobConfiguration()(*PrintJobConfiguration) {
    m := &PrintJobConfiguration{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePrintJobConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrintJobConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintJobConfiguration(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PrintJobConfiguration) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCollate gets the collate property value. Whether the printer should collate pages wehen printing multiple copies of a multi-page document.
func (m *PrintJobConfiguration) GetCollate()(*bool) {
    return m.collate
}
// GetColorMode gets the colorMode property value. The color mode the printer should use to print the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) GetColorMode()(*PrintColorMode) {
    return m.colorMode
}
// GetCopies gets the copies property value. The number of copies that should be printed. Read-only.
func (m *PrintJobConfiguration) GetCopies()(*int32) {
    return m.copies
}
// GetDpi gets the dpi property value. The resolution to use when printing the job, expressed in dots per inch (DPI). Read-only.
func (m *PrintJobConfiguration) GetDpi()(*int32) {
    return m.dpi
}
// GetDuplexMode gets the duplexMode property value. The duplex mode the printer should use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) GetDuplexMode()(*PrintDuplexMode) {
    return m.duplexMode
}
// GetFeedOrientation gets the feedOrientation property value. The orientation to use when feeding media into the printer. Valid values are described in the following table. Read-only.
func (m *PrintJobConfiguration) GetFeedOrientation()(*PrinterFeedOrientation) {
    return m.feedOrientation
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrintJobConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["collate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCollate)
    res["colorMode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintColorMode , m.SetColorMode)
    res["copies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCopies)
    res["dpi"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDpi)
    res["duplexMode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintDuplexMode , m.SetDuplexMode)
    res["feedOrientation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrinterFeedOrientation , m.SetFeedOrientation)
    res["finishings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParsePrintFinishing , m.SetFinishings)
    res["fitPdfToPage"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetFitPdfToPage)
    res["inputBin"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInputBin)
    res["margin"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintMarginFromDiscriminatorValue , m.SetMargin)
    res["mediaSize"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMediaSize)
    res["mediaType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMediaType)
    res["multipageLayout"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintMultipageLayout , m.SetMultipageLayout)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["orientation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintOrientation , m.SetOrientation)
    res["outputBin"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOutputBin)
    res["pageRanges"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIntegerRangeFromDiscriminatorValue , m.SetPageRanges)
    res["pagesPerSheet"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPagesPerSheet)
    res["quality"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintQuality , m.SetQuality)
    res["scaling"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePrintScaling , m.SetScaling)
    return res
}
// GetFinishings gets the finishings property value. Finishing processes to use when printing.
func (m *PrintJobConfiguration) GetFinishings()([]PrintFinishing) {
    return m.finishings
}
// GetFitPdfToPage gets the fitPdfToPage property value. The fitPdfToPage property
func (m *PrintJobConfiguration) GetFitPdfToPage()(*bool) {
    return m.fitPdfToPage
}
// GetInputBin gets the inputBin property value. The input bin (tray) to use when printing. See the printer's capabilities for a list of supported input bins.
func (m *PrintJobConfiguration) GetInputBin()(*string) {
    return m.inputBin
}
// GetMargin gets the margin property value. The margin settings to use when printing.
func (m *PrintJobConfiguration) GetMargin()(PrintMarginable) {
    return m.margin
}
// GetMediaSize gets the mediaSize property value. The media size to use when printing. Supports standard size names for ISO and ANSI media sizes. Valid values listed in the printerCapabilities topic.
func (m *PrintJobConfiguration) GetMediaSize()(*string) {
    return m.mediaSize
}
// GetMediaType gets the mediaType property value. The default media (such as paper) type to print the document on.
func (m *PrintJobConfiguration) GetMediaType()(*string) {
    return m.mediaType
}
// GetMultipageLayout gets the multipageLayout property value. The direction to lay out pages when multiple pages are being printed per sheet. Valid values are described in the following table.
func (m *PrintJobConfiguration) GetMultipageLayout()(*PrintMultipageLayout) {
    return m.multipageLayout
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PrintJobConfiguration) GetOdataType()(*string) {
    return m.odataType
}
// GetOrientation gets the orientation property value. The orientation setting the printer should use when printing the job. Valid values are described in the following table.
func (m *PrintJobConfiguration) GetOrientation()(*PrintOrientation) {
    return m.orientation
}
// GetOutputBin gets the outputBin property value. The output bin to place completed prints into. See the printer's capabilities for a list of supported output bins.
func (m *PrintJobConfiguration) GetOutputBin()(*string) {
    return m.outputBin
}
// GetPageRanges gets the pageRanges property value. The page ranges to print. Read-only.
func (m *PrintJobConfiguration) GetPageRanges()([]IntegerRangeable) {
    return m.pageRanges
}
// GetPagesPerSheet gets the pagesPerSheet property value. The number of document pages to print on each sheet.
func (m *PrintJobConfiguration) GetPagesPerSheet()(*int32) {
    return m.pagesPerSheet
}
// GetQuality gets the quality property value. The print quality to use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) GetQuality()(*PrintQuality) {
    return m.quality
}
// GetScaling gets the scaling property value. Specifies how the printer should scale the document data to fit the requested media. Valid values are described in the following table.
func (m *PrintJobConfiguration) GetScaling()(*PrintScaling) {
    return m.scaling
}
// Serialize serializes information the current object
func (m *PrintJobConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("collate", m.GetCollate())
        if err != nil {
            return err
        }
    }
    if m.GetColorMode() != nil {
        cast := (*m.GetColorMode()).String()
        err := writer.WriteStringValue("colorMode", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("copies", m.GetCopies())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("dpi", m.GetDpi())
        if err != nil {
            return err
        }
    }
    if m.GetDuplexMode() != nil {
        cast := (*m.GetDuplexMode()).String()
        err := writer.WriteStringValue("duplexMode", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFeedOrientation() != nil {
        cast := (*m.GetFeedOrientation()).String()
        err := writer.WriteStringValue("feedOrientation", &cast)
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
    {
        err := writer.WriteBoolValue("fitPdfToPage", m.GetFitPdfToPage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("inputBin", m.GetInputBin())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("margin", m.GetMargin())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("mediaSize", m.GetMediaSize())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("mediaType", m.GetMediaType())
        if err != nil {
            return err
        }
    }
    if m.GetMultipageLayout() != nil {
        cast := (*m.GetMultipageLayout()).String()
        err := writer.WriteStringValue("multipageLayout", &cast)
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
    if m.GetOrientation() != nil {
        cast := (*m.GetOrientation()).String()
        err := writer.WriteStringValue("orientation", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("outputBin", m.GetOutputBin())
        if err != nil {
            return err
        }
    }
    if m.GetPageRanges() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPageRanges())
        err := writer.WriteCollectionOfObjectValues("pageRanges", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("pagesPerSheet", m.GetPagesPerSheet())
        if err != nil {
            return err
        }
    }
    if m.GetQuality() != nil {
        cast := (*m.GetQuality()).String()
        err := writer.WriteStringValue("quality", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetScaling() != nil {
        cast := (*m.GetScaling()).String()
        err := writer.WriteStringValue("scaling", &cast)
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
func (m *PrintJobConfiguration) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCollate sets the collate property value. Whether the printer should collate pages wehen printing multiple copies of a multi-page document.
func (m *PrintJobConfiguration) SetCollate(value *bool)() {
    m.collate = value
}
// SetColorMode sets the colorMode property value. The color mode the printer should use to print the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetColorMode(value *PrintColorMode)() {
    m.colorMode = value
}
// SetCopies sets the copies property value. The number of copies that should be printed. Read-only.
func (m *PrintJobConfiguration) SetCopies(value *int32)() {
    m.copies = value
}
// SetDpi sets the dpi property value. The resolution to use when printing the job, expressed in dots per inch (DPI). Read-only.
func (m *PrintJobConfiguration) SetDpi(value *int32)() {
    m.dpi = value
}
// SetDuplexMode sets the duplexMode property value. The duplex mode the printer should use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetDuplexMode(value *PrintDuplexMode)() {
    m.duplexMode = value
}
// SetFeedOrientation sets the feedOrientation property value. The orientation to use when feeding media into the printer. Valid values are described in the following table. Read-only.
func (m *PrintJobConfiguration) SetFeedOrientation(value *PrinterFeedOrientation)() {
    m.feedOrientation = value
}
// SetFinishings sets the finishings property value. Finishing processes to use when printing.
func (m *PrintJobConfiguration) SetFinishings(value []PrintFinishing)() {
    m.finishings = value
}
// SetFitPdfToPage sets the fitPdfToPage property value. The fitPdfToPage property
func (m *PrintJobConfiguration) SetFitPdfToPage(value *bool)() {
    m.fitPdfToPage = value
}
// SetInputBin sets the inputBin property value. The input bin (tray) to use when printing. See the printer's capabilities for a list of supported input bins.
func (m *PrintJobConfiguration) SetInputBin(value *string)() {
    m.inputBin = value
}
// SetMargin sets the margin property value. The margin settings to use when printing.
func (m *PrintJobConfiguration) SetMargin(value PrintMarginable)() {
    m.margin = value
}
// SetMediaSize sets the mediaSize property value. The media size to use when printing. Supports standard size names for ISO and ANSI media sizes. Valid values listed in the printerCapabilities topic.
func (m *PrintJobConfiguration) SetMediaSize(value *string)() {
    m.mediaSize = value
}
// SetMediaType sets the mediaType property value. The default media (such as paper) type to print the document on.
func (m *PrintJobConfiguration) SetMediaType(value *string)() {
    m.mediaType = value
}
// SetMultipageLayout sets the multipageLayout property value. The direction to lay out pages when multiple pages are being printed per sheet. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetMultipageLayout(value *PrintMultipageLayout)() {
    m.multipageLayout = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrintJobConfiguration) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOrientation sets the orientation property value. The orientation setting the printer should use when printing the job. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetOrientation(value *PrintOrientation)() {
    m.orientation = value
}
// SetOutputBin sets the outputBin property value. The output bin to place completed prints into. See the printer's capabilities for a list of supported output bins.
func (m *PrintJobConfiguration) SetOutputBin(value *string)() {
    m.outputBin = value
}
// SetPageRanges sets the pageRanges property value. The page ranges to print. Read-only.
func (m *PrintJobConfiguration) SetPageRanges(value []IntegerRangeable)() {
    m.pageRanges = value
}
// SetPagesPerSheet sets the pagesPerSheet property value. The number of document pages to print on each sheet.
func (m *PrintJobConfiguration) SetPagesPerSheet(value *int32)() {
    m.pagesPerSheet = value
}
// SetQuality sets the quality property value. The print quality to use when printing the job. Valid values are described in the table below. Read-only.
func (m *PrintJobConfiguration) SetQuality(value *PrintQuality)() {
    m.quality = value
}
// SetScaling sets the scaling property value. Specifies how the printer should scale the document data to fit the requested media. Valid values are described in the following table.
func (m *PrintJobConfiguration) SetScaling(value *PrintScaling)() {
    m.scaling = value
}
