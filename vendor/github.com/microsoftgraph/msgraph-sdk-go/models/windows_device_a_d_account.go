package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsDeviceADAccount 
type WindowsDeviceADAccount struct {
    WindowsDeviceAccount
    // Not yet documented
    domainName *string
    // Not yet documented
    userName *string
}
// NewWindowsDeviceADAccount instantiates a new WindowsDeviceADAccount and sets the default values.
func NewWindowsDeviceADAccount()(*WindowsDeviceADAccount) {
    m := &WindowsDeviceADAccount{
        WindowsDeviceAccount: *NewWindowsDeviceAccount(),
    }
    odataTypeValue := "#microsoft.graph.windowsDeviceADAccount";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindowsDeviceADAccountFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsDeviceADAccountFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsDeviceADAccount(), nil
}
// GetDomainName gets the domainName property value. Not yet documented
func (m *WindowsDeviceADAccount) GetDomainName()(*string) {
    return m.domainName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsDeviceADAccount) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WindowsDeviceAccount.GetFieldDeserializers()
    res["domainName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDomainName)
    res["userName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserName)
    return res
}
// GetUserName gets the userName property value. Not yet documented
func (m *WindowsDeviceADAccount) GetUserName()(*string) {
    return m.userName
}
// Serialize serializes information the current object
func (m *WindowsDeviceADAccount) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WindowsDeviceAccount.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("domainName", m.GetDomainName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userName", m.GetUserName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDomainName sets the domainName property value. Not yet documented
func (m *WindowsDeviceADAccount) SetDomainName(value *string)() {
    m.domainName = value
}
// SetUserName sets the userName property value. Not yet documented
func (m *WindowsDeviceADAccount) SetUserName(value *string)() {
    m.userName = value
}
