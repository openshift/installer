package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ReportRoot 
type ReportRoot struct {
    Entity
    // The dailyPrintUsageByPrinter property
    dailyPrintUsageByPrinter []PrintUsageByPrinterable
    // The dailyPrintUsageByUser property
    dailyPrintUsageByUser []PrintUsageByUserable
    // The monthlyPrintUsageByPrinter property
    monthlyPrintUsageByPrinter []PrintUsageByPrinterable
    // The monthlyPrintUsageByUser property
    monthlyPrintUsageByUser []PrintUsageByUserable
    // The security property
    security SecurityReportsRootable
}
// NewReportRoot instantiates a new ReportRoot and sets the default values.
func NewReportRoot()(*ReportRoot) {
    m := &ReportRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateReportRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateReportRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewReportRoot(), nil
}
// GetDailyPrintUsageByPrinter gets the dailyPrintUsageByPrinter property value. The dailyPrintUsageByPrinter property
func (m *ReportRoot) GetDailyPrintUsageByPrinter()([]PrintUsageByPrinterable) {
    return m.dailyPrintUsageByPrinter
}
// GetDailyPrintUsageByUser gets the dailyPrintUsageByUser property value. The dailyPrintUsageByUser property
func (m *ReportRoot) GetDailyPrintUsageByUser()([]PrintUsageByUserable) {
    return m.dailyPrintUsageByUser
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ReportRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["dailyPrintUsageByPrinter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintUsageByPrinterFromDiscriminatorValue , m.SetDailyPrintUsageByPrinter)
    res["dailyPrintUsageByUser"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintUsageByUserFromDiscriminatorValue , m.SetDailyPrintUsageByUser)
    res["monthlyPrintUsageByPrinter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintUsageByPrinterFromDiscriminatorValue , m.SetMonthlyPrintUsageByPrinter)
    res["monthlyPrintUsageByUser"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePrintUsageByUserFromDiscriminatorValue , m.SetMonthlyPrintUsageByUser)
    res["security"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSecurityReportsRootFromDiscriminatorValue , m.SetSecurity)
    return res
}
// GetMonthlyPrintUsageByPrinter gets the monthlyPrintUsageByPrinter property value. The monthlyPrintUsageByPrinter property
func (m *ReportRoot) GetMonthlyPrintUsageByPrinter()([]PrintUsageByPrinterable) {
    return m.monthlyPrintUsageByPrinter
}
// GetMonthlyPrintUsageByUser gets the monthlyPrintUsageByUser property value. The monthlyPrintUsageByUser property
func (m *ReportRoot) GetMonthlyPrintUsageByUser()([]PrintUsageByUserable) {
    return m.monthlyPrintUsageByUser
}
// GetSecurity gets the security property value. The security property
func (m *ReportRoot) GetSecurity()(SecurityReportsRootable) {
    return m.security
}
// Serialize serializes information the current object
func (m *ReportRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDailyPrintUsageByPrinter() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDailyPrintUsageByPrinter())
        err = writer.WriteCollectionOfObjectValues("dailyPrintUsageByPrinter", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDailyPrintUsageByUser() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDailyPrintUsageByUser())
        err = writer.WriteCollectionOfObjectValues("dailyPrintUsageByUser", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMonthlyPrintUsageByPrinter() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMonthlyPrintUsageByPrinter())
        err = writer.WriteCollectionOfObjectValues("monthlyPrintUsageByPrinter", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMonthlyPrintUsageByUser() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMonthlyPrintUsageByUser())
        err = writer.WriteCollectionOfObjectValues("monthlyPrintUsageByUser", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("security", m.GetSecurity())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDailyPrintUsageByPrinter sets the dailyPrintUsageByPrinter property value. The dailyPrintUsageByPrinter property
func (m *ReportRoot) SetDailyPrintUsageByPrinter(value []PrintUsageByPrinterable)() {
    m.dailyPrintUsageByPrinter = value
}
// SetDailyPrintUsageByUser sets the dailyPrintUsageByUser property value. The dailyPrintUsageByUser property
func (m *ReportRoot) SetDailyPrintUsageByUser(value []PrintUsageByUserable)() {
    m.dailyPrintUsageByUser = value
}
// SetMonthlyPrintUsageByPrinter sets the monthlyPrintUsageByPrinter property value. The monthlyPrintUsageByPrinter property
func (m *ReportRoot) SetMonthlyPrintUsageByPrinter(value []PrintUsageByPrinterable)() {
    m.monthlyPrintUsageByPrinter = value
}
// SetMonthlyPrintUsageByUser sets the monthlyPrintUsageByUser property value. The monthlyPrintUsageByUser property
func (m *ReportRoot) SetMonthlyPrintUsageByUser(value []PrintUsageByUserable)() {
    m.monthlyPrintUsageByUser = value
}
// SetSecurity sets the security property value. The security property
func (m *ReportRoot) SetSecurity(value SecurityReportsRootable)() {
    m.security = value
}
