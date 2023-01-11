package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUniversalAppX 
type WindowsUniversalAppX struct {
    MobileLobApp
    // Contains properties for Windows architecture.
    applicableArchitectures *WindowsArchitecture
    // Contains properties for Windows device type.
    applicableDeviceTypes *WindowsDeviceType
    // The Identity Name.
    identityName *string
    // The Identity Publisher Hash.
    identityPublisherHash *string
    // The Identity Resource Identifier.
    identityResourceIdentifier *string
    // The identity version.
    identityVersion *string
    // Whether or not the app is a bundle.
    isBundle *bool
    // The minimum operating system required for a Windows mobile app.
    minimumSupportedOperatingSystem WindowsMinimumOperatingSystemable
}
// NewWindowsUniversalAppX instantiates a new WindowsUniversalAppX and sets the default values.
func NewWindowsUniversalAppX()(*WindowsUniversalAppX) {
    m := &WindowsUniversalAppX{
        MobileLobApp: *NewMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsUniversalAppX";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindowsUniversalAppXFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsUniversalAppXFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUniversalAppX(), nil
}
// GetApplicableArchitectures gets the applicableArchitectures property value. Contains properties for Windows architecture.
func (m *WindowsUniversalAppX) GetApplicableArchitectures()(*WindowsArchitecture) {
    return m.applicableArchitectures
}
// GetApplicableDeviceTypes gets the applicableDeviceTypes property value. Contains properties for Windows device type.
func (m *WindowsUniversalAppX) GetApplicableDeviceTypes()(*WindowsDeviceType) {
    return m.applicableDeviceTypes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsUniversalAppX) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileLobApp.GetFieldDeserializers()
    res["applicableArchitectures"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsArchitecture , m.SetApplicableArchitectures)
    res["applicableDeviceTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsDeviceType , m.SetApplicableDeviceTypes)
    res["identityName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIdentityName)
    res["identityPublisherHash"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIdentityPublisherHash)
    res["identityResourceIdentifier"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIdentityResourceIdentifier)
    res["identityVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIdentityVersion)
    res["isBundle"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsBundle)
    res["minimumSupportedOperatingSystem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWindowsMinimumOperatingSystemFromDiscriminatorValue , m.SetMinimumSupportedOperatingSystem)
    return res
}
// GetIdentityName gets the identityName property value. The Identity Name.
func (m *WindowsUniversalAppX) GetIdentityName()(*string) {
    return m.identityName
}
// GetIdentityPublisherHash gets the identityPublisherHash property value. The Identity Publisher Hash.
func (m *WindowsUniversalAppX) GetIdentityPublisherHash()(*string) {
    return m.identityPublisherHash
}
// GetIdentityResourceIdentifier gets the identityResourceIdentifier property value. The Identity Resource Identifier.
func (m *WindowsUniversalAppX) GetIdentityResourceIdentifier()(*string) {
    return m.identityResourceIdentifier
}
// GetIdentityVersion gets the identityVersion property value. The identity version.
func (m *WindowsUniversalAppX) GetIdentityVersion()(*string) {
    return m.identityVersion
}
// GetIsBundle gets the isBundle property value. Whether or not the app is a bundle.
func (m *WindowsUniversalAppX) GetIsBundle()(*bool) {
    return m.isBundle
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. The minimum operating system required for a Windows mobile app.
func (m *WindowsUniversalAppX) GetMinimumSupportedOperatingSystem()(WindowsMinimumOperatingSystemable) {
    return m.minimumSupportedOperatingSystem
}
// Serialize serializes information the current object
func (m *WindowsUniversalAppX) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileLobApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApplicableArchitectures() != nil {
        cast := (*m.GetApplicableArchitectures()).String()
        err = writer.WriteStringValue("applicableArchitectures", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetApplicableDeviceTypes() != nil {
        cast := (*m.GetApplicableDeviceTypes()).String()
        err = writer.WriteStringValue("applicableDeviceTypes", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityName", m.GetIdentityName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityPublisherHash", m.GetIdentityPublisherHash())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityResourceIdentifier", m.GetIdentityResourceIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityVersion", m.GetIdentityVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isBundle", m.GetIsBundle())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("minimumSupportedOperatingSystem", m.GetMinimumSupportedOperatingSystem())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableArchitectures sets the applicableArchitectures property value. Contains properties for Windows architecture.
func (m *WindowsUniversalAppX) SetApplicableArchitectures(value *WindowsArchitecture)() {
    m.applicableArchitectures = value
}
// SetApplicableDeviceTypes sets the applicableDeviceTypes property value. Contains properties for Windows device type.
func (m *WindowsUniversalAppX) SetApplicableDeviceTypes(value *WindowsDeviceType)() {
    m.applicableDeviceTypes = value
}
// SetIdentityName sets the identityName property value. The Identity Name.
func (m *WindowsUniversalAppX) SetIdentityName(value *string)() {
    m.identityName = value
}
// SetIdentityPublisherHash sets the identityPublisherHash property value. The Identity Publisher Hash.
func (m *WindowsUniversalAppX) SetIdentityPublisherHash(value *string)() {
    m.identityPublisherHash = value
}
// SetIdentityResourceIdentifier sets the identityResourceIdentifier property value. The Identity Resource Identifier.
func (m *WindowsUniversalAppX) SetIdentityResourceIdentifier(value *string)() {
    m.identityResourceIdentifier = value
}
// SetIdentityVersion sets the identityVersion property value. The identity version.
func (m *WindowsUniversalAppX) SetIdentityVersion(value *string)() {
    m.identityVersion = value
}
// SetIsBundle sets the isBundle property value. Whether or not the app is a bundle.
func (m *WindowsUniversalAppX) SetIsBundle(value *bool)() {
    m.isBundle = value
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. The minimum operating system required for a Windows mobile app.
func (m *WindowsUniversalAppX) SetMinimumSupportedOperatingSystem(value WindowsMinimumOperatingSystemable)() {
    m.minimumSupportedOperatingSystem = value
}
