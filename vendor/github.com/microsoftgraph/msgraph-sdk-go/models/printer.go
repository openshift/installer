package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Printer 
type Printer struct {
    PrinterBase
    // The connectors that are associated with the printer.
    connectors []PrintConnectorable
    // True if the printer has a physical device for printing. Read-only.
    hasPhysicalDevice *bool
    // True if the printer is shared; false otherwise. Read-only.
    isShared *bool
    // The most recent dateTimeOffset when a printer interacted with Universal Print. Read-only.
    lastSeenDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The DateTimeOffset when the printer was registered. Read-only.
    registeredDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The list of printerShares that are associated with the printer. Currently, only one printerShare can be associated with the printer. Read-only. Nullable.
    shares []PrinterShareable
    // A list of task triggers that are associated with the printer.
    taskTriggers []PrintTaskTriggerable
}
// NewPrinter instantiates a new Printer and sets the default values.
func NewPrinter()(*Printer) {
    m := &Printer{
        PrinterBase: *NewPrinterBase(),
    }
    odataTypeValue := "#microsoft.graph.printer";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreatePrinterFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrinterFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrinter(), nil
}
// GetConnectors gets the connectors property value. The connectors that are associated with the printer.
func (m *Printer) GetConnectors()([]PrintConnectorable) {
    return m.connectors
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Printer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PrinterBase.GetFieldDeserializers()
    res["connectors"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintConnectorFromDiscriminatorValue , m.SetConnectors)
    res["hasPhysicalDevice"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasPhysicalDevice)
    res["isShared"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsShared)
    res["lastSeenDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastSeenDateTime)
    res["registeredDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetRegisteredDateTime)
    res["shares"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrinterShareFromDiscriminatorValue , m.SetShares)
    res["taskTriggers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintTaskTriggerFromDiscriminatorValue , m.SetTaskTriggers)
    return res
}
// GetHasPhysicalDevice gets the hasPhysicalDevice property value. True if the printer has a physical device for printing. Read-only.
func (m *Printer) GetHasPhysicalDevice()(*bool) {
    return m.hasPhysicalDevice
}
// GetIsShared gets the isShared property value. True if the printer is shared; false otherwise. Read-only.
func (m *Printer) GetIsShared()(*bool) {
    return m.isShared
}
// GetLastSeenDateTime gets the lastSeenDateTime property value. The most recent dateTimeOffset when a printer interacted with Universal Print. Read-only.
func (m *Printer) GetLastSeenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastSeenDateTime
}
// GetRegisteredDateTime gets the registeredDateTime property value. The DateTimeOffset when the printer was registered. Read-only.
func (m *Printer) GetRegisteredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.registeredDateTime
}
// GetShares gets the shares property value. The list of printerShares that are associated with the printer. Currently, only one printerShare can be associated with the printer. Read-only. Nullable.
func (m *Printer) GetShares()([]PrinterShareable) {
    return m.shares
}
// GetTaskTriggers gets the taskTriggers property value. A list of task triggers that are associated with the printer.
func (m *Printer) GetTaskTriggers()([]PrintTaskTriggerable) {
    return m.taskTriggers
}
// Serialize serializes information the current object
func (m *Printer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PrinterBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetConnectors() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetConnectors())
        err = writer.WriteCollectionOfObjectValues("connectors", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasPhysicalDevice", m.GetHasPhysicalDevice())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isShared", m.GetIsShared())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastSeenDateTime", m.GetLastSeenDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("registeredDateTime", m.GetRegisteredDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetShares() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetShares())
        err = writer.WriteCollectionOfObjectValues("shares", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTaskTriggers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTaskTriggers())
        err = writer.WriteCollectionOfObjectValues("taskTriggers", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConnectors sets the connectors property value. The connectors that are associated with the printer.
func (m *Printer) SetConnectors(value []PrintConnectorable)() {
    m.connectors = value
}
// SetHasPhysicalDevice sets the hasPhysicalDevice property value. True if the printer has a physical device for printing. Read-only.
func (m *Printer) SetHasPhysicalDevice(value *bool)() {
    m.hasPhysicalDevice = value
}
// SetIsShared sets the isShared property value. True if the printer is shared; false otherwise. Read-only.
func (m *Printer) SetIsShared(value *bool)() {
    m.isShared = value
}
// SetLastSeenDateTime sets the lastSeenDateTime property value. The most recent dateTimeOffset when a printer interacted with Universal Print. Read-only.
func (m *Printer) SetLastSeenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastSeenDateTime = value
}
// SetRegisteredDateTime sets the registeredDateTime property value. The DateTimeOffset when the printer was registered. Read-only.
func (m *Printer) SetRegisteredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.registeredDateTime = value
}
// SetShares sets the shares property value. The list of printerShares that are associated with the printer. Currently, only one printerShare can be associated with the printer. Read-only. Nullable.
func (m *Printer) SetShares(value []PrinterShareable)() {
    m.shares = value
}
// SetTaskTriggers sets the taskTriggers property value. A list of task triggers that are associated with the printer.
func (m *Printer) SetTaskTriggers(value []PrintTaskTriggerable)() {
    m.taskTriggers = value
}
