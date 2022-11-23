package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppProductCodeRule 
type Win32LobAppProductCodeRule struct {
    Win32LobAppRule
    // The product code of the app.
    productCode *string
    // The product version comparison value.
    productVersion *string
    // Contains properties for detection operator.
    productVersionOperator *Win32LobAppRuleOperator
}
// NewWin32LobAppProductCodeRule instantiates a new Win32LobAppProductCodeRule and sets the default values.
func NewWin32LobAppProductCodeRule()(*Win32LobAppProductCodeRule) {
    m := &Win32LobAppProductCodeRule{
        Win32LobAppRule: *NewWin32LobAppRule(),
    }
    odataTypeValue := "#microsoft.graph.win32LobAppProductCodeRule";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWin32LobAppProductCodeRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWin32LobAppProductCodeRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppProductCodeRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Win32LobAppProductCodeRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Win32LobAppRule.GetFieldDeserializers()
    res["productCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetProductCode)
    res["productVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetProductVersion)
    res["productVersionOperator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWin32LobAppRuleOperator , m.SetProductVersionOperator)
    return res
}
// GetProductCode gets the productCode property value. The product code of the app.
func (m *Win32LobAppProductCodeRule) GetProductCode()(*string) {
    return m.productCode
}
// GetProductVersion gets the productVersion property value. The product version comparison value.
func (m *Win32LobAppProductCodeRule) GetProductVersion()(*string) {
    return m.productVersion
}
// GetProductVersionOperator gets the productVersionOperator property value. Contains properties for detection operator.
func (m *Win32LobAppProductCodeRule) GetProductVersionOperator()(*Win32LobAppRuleOperator) {
    return m.productVersionOperator
}
// Serialize serializes information the current object
func (m *Win32LobAppProductCodeRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Win32LobAppRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("productCode", m.GetProductCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productVersion", m.GetProductVersion())
        if err != nil {
            return err
        }
    }
    if m.GetProductVersionOperator() != nil {
        cast := (*m.GetProductVersionOperator()).String()
        err = writer.WriteStringValue("productVersionOperator", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetProductCode sets the productCode property value. The product code of the app.
func (m *Win32LobAppProductCodeRule) SetProductCode(value *string)() {
    m.productCode = value
}
// SetProductVersion sets the productVersion property value. The product version comparison value.
func (m *Win32LobAppProductCodeRule) SetProductVersion(value *string)() {
    m.productVersion = value
}
// SetProductVersionOperator sets the productVersionOperator property value. Contains properties for detection operator.
func (m *Win32LobAppProductCodeRule) SetProductVersionOperator(value *Win32LobAppRuleOperator)() {
    m.productVersionOperator = value
}
