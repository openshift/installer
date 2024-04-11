package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrinterShare 
type PrinterShare struct {
    PrinterBase
}
// NewPrinterShare instantiates a new printerShare and sets the default values.
func NewPrinterShare()(*PrinterShare) {
    m := &PrinterShare{
        PrinterBase: *NewPrinterBase(),
    }
    odataTypeValue := "#microsoft.graph.printerShare"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePrinterShareFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrinterShareFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrinterShare(), nil
}
// GetAllowAllUsers gets the allowAllUsers property value. If true, all users and groups will be granted access to this printer share. This supersedes the allow lists defined by the allowedUsers and allowedGroups navigation properties.
func (m *PrinterShare) GetAllowAllUsers()(*bool) {
    val, err := m.GetBackingStore().Get("allowAllUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedGroups gets the allowedGroups property value. The groups whose users have access to print using the printer.
func (m *PrinterShare) GetAllowedGroups()([]Groupable) {
    val, err := m.GetBackingStore().Get("allowedGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Groupable)
    }
    return nil
}
// GetAllowedUsers gets the allowedUsers property value. The users who have access to print using the printer.
func (m *PrinterShare) GetAllowedUsers()([]Userable) {
    val, err := m.GetBackingStore().Get("allowedUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Userable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The DateTimeOffset when the printer share was created. Read-only.
func (m *PrinterShare) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrinterShare) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PrinterBase.GetFieldDeserializers()
    res["allowAllUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowAllUsers(val)
        }
        return nil
    }
    res["allowedGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Groupable, len(val))
            for i, v := range val {
                res[i] = v.(Groupable)
            }
            m.SetAllowedGroups(res)
        }
        return nil
    }
    res["allowedUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Userable, len(val))
            for i, v := range val {
                res[i] = v.(Userable)
            }
            m.SetAllowedUsers(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["printer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrinterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrinter(val.(Printerable))
        }
        return nil
    }
    return res
}
// GetPrinter gets the printer property value. The printer that this printer share is related to.
func (m *PrinterShare) GetPrinter()(Printerable) {
    val, err := m.GetBackingStore().Get("printer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Printerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrinterShare) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PrinterBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowAllUsers", m.GetAllowAllUsers())
        if err != nil {
            return err
        }
    }
    if m.GetAllowedGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllowedGroups()))
        for i, v := range m.GetAllowedGroups() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("allowedGroups", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAllowedUsers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllowedUsers()))
        for i, v := range m.GetAllowedUsers() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("allowedUsers", cast)
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
    {
        err = writer.WriteObjectValue("printer", m.GetPrinter())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowAllUsers sets the allowAllUsers property value. If true, all users and groups will be granted access to this printer share. This supersedes the allow lists defined by the allowedUsers and allowedGroups navigation properties.
func (m *PrinterShare) SetAllowAllUsers(value *bool)() {
    err := m.GetBackingStore().Set("allowAllUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedGroups sets the allowedGroups property value. The groups whose users have access to print using the printer.
func (m *PrinterShare) SetAllowedGroups(value []Groupable)() {
    err := m.GetBackingStore().Set("allowedGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedUsers sets the allowedUsers property value. The users who have access to print using the printer.
func (m *PrinterShare) SetAllowedUsers(value []Userable)() {
    err := m.GetBackingStore().Set("allowedUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The DateTimeOffset when the printer share was created. Read-only.
func (m *PrinterShare) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPrinter sets the printer property value. The printer that this printer share is related to.
func (m *PrinterShare) SetPrinter(value Printerable)() {
    err := m.GetBackingStore().Set("printer", value)
    if err != nil {
        panic(err)
    }
}
// PrinterShareable 
type PrinterShareable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PrinterBaseable
    GetAllowAllUsers()(*bool)
    GetAllowedGroups()([]Groupable)
    GetAllowedUsers()([]Userable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPrinter()(Printerable)
    SetAllowAllUsers(value *bool)()
    SetAllowedGroups(value []Groupable)()
    SetAllowedUsers(value []Userable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPrinter(value Printerable)()
}
