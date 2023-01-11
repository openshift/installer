package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrintJob provides operations to manage the collection of agreement entities.
type PrintJob struct {
    Entity
    // The configuration property
    configuration PrintJobConfigurationable
    // The createdBy property
    createdBy UserIdentityable
    // The DateTimeOffset when the job was created. Read-only.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The documents property
    documents []PrintDocumentable
    // If true, document can be fetched by printer.
    isFetchable *bool
    // Contains the source job URL, if the job has been redirected from another printer.
    redirectedFrom *string
    // Contains the destination job URL, if the job has been redirected to another printer.
    redirectedTo *string
    // The status property
    status PrintJobStatusable
    // A list of printTasks that were triggered by this print job.
    tasks []PrintTaskable
}
// NewPrintJob instantiates a new printJob and sets the default values.
func NewPrintJob()(*PrintJob) {
    m := &PrintJob{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintJobFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrintJobFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintJob(), nil
}
// GetConfiguration gets the configuration property value. The configuration property
func (m *PrintJob) GetConfiguration()(PrintJobConfigurationable) {
    return m.configuration
}
// GetCreatedBy gets the createdBy property value. The createdBy property
func (m *PrintJob) GetCreatedBy()(UserIdentityable) {
    return m.createdBy
}
// GetCreatedDateTime gets the createdDateTime property value. The DateTimeOffset when the job was created. Read-only.
func (m *PrintJob) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDocuments gets the documents property value. The documents property
func (m *PrintJob) GetDocuments()([]PrintDocumentable) {
    return m.documents
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrintJob) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configuration"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintJobConfigurationFromDiscriminatorValue , m.SetConfiguration)
    res["createdBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserIdentityFromDiscriminatorValue , m.SetCreatedBy)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["documents"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintDocumentFromDiscriminatorValue , m.SetDocuments)
    res["isFetchable"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsFetchable)
    res["redirectedFrom"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRedirectedFrom)
    res["redirectedTo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRedirectedTo)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintJobStatusFromDiscriminatorValue , m.SetStatus)
    res["tasks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintTaskFromDiscriminatorValue , m.SetTasks)
    return res
}
// GetIsFetchable gets the isFetchable property value. If true, document can be fetched by printer.
func (m *PrintJob) GetIsFetchable()(*bool) {
    return m.isFetchable
}
// GetRedirectedFrom gets the redirectedFrom property value. Contains the source job URL, if the job has been redirected from another printer.
func (m *PrintJob) GetRedirectedFrom()(*string) {
    return m.redirectedFrom
}
// GetRedirectedTo gets the redirectedTo property value. Contains the destination job URL, if the job has been redirected to another printer.
func (m *PrintJob) GetRedirectedTo()(*string) {
    return m.redirectedTo
}
// GetStatus gets the status property value. The status property
func (m *PrintJob) GetStatus()(PrintJobStatusable) {
    return m.status
}
// GetTasks gets the tasks property value. A list of printTasks that were triggered by this print job.
func (m *PrintJob) GetTasks()([]PrintTaskable) {
    return m.tasks
}
// Serialize serializes information the current object
func (m *PrintJob) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("configuration", m.GetConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetDocuments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDocuments())
        err = writer.WriteCollectionOfObjectValues("documents", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFetchable", m.GetIsFetchable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("redirectedFrom", m.GetRedirectedFrom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("redirectedTo", m.GetRedirectedTo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    if m.GetTasks() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTasks())
        err = writer.WriteCollectionOfObjectValues("tasks", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfiguration sets the configuration property value. The configuration property
func (m *PrintJob) SetConfiguration(value PrintJobConfigurationable)() {
    m.configuration = value
}
// SetCreatedBy sets the createdBy property value. The createdBy property
func (m *PrintJob) SetCreatedBy(value UserIdentityable)() {
    m.createdBy = value
}
// SetCreatedDateTime sets the createdDateTime property value. The DateTimeOffset when the job was created. Read-only.
func (m *PrintJob) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDocuments sets the documents property value. The documents property
func (m *PrintJob) SetDocuments(value []PrintDocumentable)() {
    m.documents = value
}
// SetIsFetchable sets the isFetchable property value. If true, document can be fetched by printer.
func (m *PrintJob) SetIsFetchable(value *bool)() {
    m.isFetchable = value
}
// SetRedirectedFrom sets the redirectedFrom property value. Contains the source job URL, if the job has been redirected from another printer.
func (m *PrintJob) SetRedirectedFrom(value *string)() {
    m.redirectedFrom = value
}
// SetRedirectedTo sets the redirectedTo property value. Contains the destination job URL, if the job has been redirected to another printer.
func (m *PrintJob) SetRedirectedTo(value *string)() {
    m.redirectedTo = value
}
// SetStatus sets the status property value. The status property
func (m *PrintJob) SetStatus(value PrintJobStatusable)() {
    m.status = value
}
// SetTasks sets the tasks property value. A list of printTasks that were triggered by this print job.
func (m *PrintJob) SetTasks(value []PrintTaskable)() {
    m.tasks = value
}
