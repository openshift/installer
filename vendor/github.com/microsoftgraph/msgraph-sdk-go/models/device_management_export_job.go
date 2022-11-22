package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementExportJob entity representing a job to export a report
type DeviceManagementExportJob struct {
    Entity
    // Time that the exported report expires
    expirationDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Filters applied on the report
    filter *string
    // Possible values for the file format of a report
    format *DeviceManagementReportFileFormat
    // Configures how the requested export job is localized
    localizationType *DeviceManagementExportJobLocalizationType
    // Name of the report
    reportName *string
    // Time that the exported report was requested
    requestDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Columns selected from the report
    select_escaped []string
    // A snapshot is an identifiable subset of the dataset represented by the ReportName. A sessionId or CachedReportConfiguration id can be used here. If a sessionId is specified, Filter, Select, and OrderBy are applied to the data represented by the sessionId. Filter, Select, and OrderBy cannot be specified together with a CachedReportConfiguration id.
    snapshotId *string
    // Possible statuses associated with a generated report
    status *DeviceManagementReportStatus
    // Temporary location of the exported report
    url *string
}
// NewDeviceManagementExportJob instantiates a new deviceManagementExportJob and sets the default values.
func NewDeviceManagementExportJob()(*DeviceManagementExportJob) {
    m := &DeviceManagementExportJob{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementExportJobFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceManagementExportJobFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementExportJob(), nil
}
// GetExpirationDateTime gets the expirationDateTime property value. Time that the exported report expires
func (m *DeviceManagementExportJob) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.expirationDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceManagementExportJob) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["expirationDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetExpirationDateTime)
    res["filter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetFilter)
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementReportFileFormat , m.SetFormat)
    res["localizationType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementExportJobLocalizationType , m.SetLocalizationType)
    res["reportName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReportName)
    res["requestDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetRequestDateTime)
    res["select"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetSelect)
    res["snapshotId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSnapshotId)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementReportStatus , m.SetStatus)
    res["url"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUrl)
    return res
}
// GetFilter gets the filter property value. Filters applied on the report
func (m *DeviceManagementExportJob) GetFilter()(*string) {
    return m.filter
}
// GetFormat gets the format property value. Possible values for the file format of a report
func (m *DeviceManagementExportJob) GetFormat()(*DeviceManagementReportFileFormat) {
    return m.format
}
// GetLocalizationType gets the localizationType property value. Configures how the requested export job is localized
func (m *DeviceManagementExportJob) GetLocalizationType()(*DeviceManagementExportJobLocalizationType) {
    return m.localizationType
}
// GetReportName gets the reportName property value. Name of the report
func (m *DeviceManagementExportJob) GetReportName()(*string) {
    return m.reportName
}
// GetRequestDateTime gets the requestDateTime property value. Time that the exported report was requested
func (m *DeviceManagementExportJob) GetRequestDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.requestDateTime
}
// GetSelect gets the select property value. Columns selected from the report
func (m *DeviceManagementExportJob) GetSelect()([]string) {
    return m.select_escaped
}
// GetSnapshotId gets the snapshotId property value. A snapshot is an identifiable subset of the dataset represented by the ReportName. A sessionId or CachedReportConfiguration id can be used here. If a sessionId is specified, Filter, Select, and OrderBy are applied to the data represented by the sessionId. Filter, Select, and OrderBy cannot be specified together with a CachedReportConfiguration id.
func (m *DeviceManagementExportJob) GetSnapshotId()(*string) {
    return m.snapshotId
}
// GetStatus gets the status property value. Possible statuses associated with a generated report
func (m *DeviceManagementExportJob) GetStatus()(*DeviceManagementReportStatus) {
    return m.status
}
// GetUrl gets the url property value. Temporary location of the exported report
func (m *DeviceManagementExportJob) GetUrl()(*string) {
    return m.url
}
// Serialize serializes information the current object
func (m *DeviceManagementExportJob) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("expirationDateTime", m.GetExpirationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("filter", m.GetFilter())
        if err != nil {
            return err
        }
    }
    if m.GetFormat() != nil {
        cast := (*m.GetFormat()).String()
        err = writer.WriteStringValue("format", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetLocalizationType() != nil {
        cast := (*m.GetLocalizationType()).String()
        err = writer.WriteStringValue("localizationType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("reportName", m.GetReportName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("requestDateTime", m.GetRequestDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetSelect() != nil {
        err = writer.WriteCollectionOfStringValues("select", m.GetSelect())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("snapshotId", m.GetSnapshotId())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("url", m.GetUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExpirationDateTime sets the expirationDateTime property value. Time that the exported report expires
func (m *DeviceManagementExportJob) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.expirationDateTime = value
}
// SetFilter sets the filter property value. Filters applied on the report
func (m *DeviceManagementExportJob) SetFilter(value *string)() {
    m.filter = value
}
// SetFormat sets the format property value. Possible values for the file format of a report
func (m *DeviceManagementExportJob) SetFormat(value *DeviceManagementReportFileFormat)() {
    m.format = value
}
// SetLocalizationType sets the localizationType property value. Configures how the requested export job is localized
func (m *DeviceManagementExportJob) SetLocalizationType(value *DeviceManagementExportJobLocalizationType)() {
    m.localizationType = value
}
// SetReportName sets the reportName property value. Name of the report
func (m *DeviceManagementExportJob) SetReportName(value *string)() {
    m.reportName = value
}
// SetRequestDateTime sets the requestDateTime property value. Time that the exported report was requested
func (m *DeviceManagementExportJob) SetRequestDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.requestDateTime = value
}
// SetSelect sets the select property value. Columns selected from the report
func (m *DeviceManagementExportJob) SetSelect(value []string)() {
    m.select_escaped = value
}
// SetSnapshotId sets the snapshotId property value. A snapshot is an identifiable subset of the dataset represented by the ReportName. A sessionId or CachedReportConfiguration id can be used here. If a sessionId is specified, Filter, Select, and OrderBy are applied to the data represented by the sessionId. Filter, Select, and OrderBy cannot be specified together with a CachedReportConfiguration id.
func (m *DeviceManagementExportJob) SetSnapshotId(value *string)() {
    m.snapshotId = value
}
// SetStatus sets the status property value. Possible statuses associated with a generated report
func (m *DeviceManagementExportJob) SetStatus(value *DeviceManagementReportStatus)() {
    m.status = value
}
// SetUrl sets the url property value. Temporary location of the exported report
func (m *DeviceManagementExportJob) SetUrl(value *string)() {
    m.url = value
}
