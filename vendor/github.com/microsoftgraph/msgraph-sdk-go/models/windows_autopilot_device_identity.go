package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsAutopilotDeviceIdentity the windowsAutopilotDeviceIdentity resource represents a Windows Autopilot Device.
type WindowsAutopilotDeviceIdentity struct {
    Entity
    // Addressable user name.
    addressableUserName *string
    // AAD Device ID - to be deprecated
    azureActiveDirectoryDeviceId *string
    // Display Name
    displayName *string
    // The enrollmentState property
    enrollmentState *EnrollmentState
    // Group Tag of the Windows autopilot device.
    groupTag *string
    // Intune Last Contacted Date Time of the Windows autopilot device.
    lastContactedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Managed Device ID
    managedDeviceId *string
    // Oem manufacturer of the Windows autopilot device.
    manufacturer *string
    // Model name of the Windows autopilot device.
    model *string
    // Product Key of the Windows autopilot device.
    productKey *string
    // Purchase Order Identifier of the Windows autopilot device.
    purchaseOrderIdentifier *string
    // Resource Name.
    resourceName *string
    // Serial number of the Windows autopilot device.
    serialNumber *string
    // SKU Number
    skuNumber *string
    // System Family
    systemFamily *string
    // User Principal Name.
    userPrincipalName *string
}
// NewWindowsAutopilotDeviceIdentity instantiates a new windowsAutopilotDeviceIdentity and sets the default values.
func NewWindowsAutopilotDeviceIdentity()(*WindowsAutopilotDeviceIdentity) {
    m := &WindowsAutopilotDeviceIdentity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWindowsAutopilotDeviceIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsAutopilotDeviceIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsAutopilotDeviceIdentity(), nil
}
// GetAddressableUserName gets the addressableUserName property value. Addressable user name.
func (m *WindowsAutopilotDeviceIdentity) GetAddressableUserName()(*string) {
    return m.addressableUserName
}
// GetAzureActiveDirectoryDeviceId gets the azureActiveDirectoryDeviceId property value. AAD Device ID - to be deprecated
func (m *WindowsAutopilotDeviceIdentity) GetAzureActiveDirectoryDeviceId()(*string) {
    return m.azureActiveDirectoryDeviceId
}
// GetDisplayName gets the displayName property value. Display Name
func (m *WindowsAutopilotDeviceIdentity) GetDisplayName()(*string) {
    return m.displayName
}
// GetEnrollmentState gets the enrollmentState property value. The enrollmentState property
func (m *WindowsAutopilotDeviceIdentity) GetEnrollmentState()(*EnrollmentState) {
    return m.enrollmentState
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsAutopilotDeviceIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["addressableUserName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAddressableUserName)
    res["azureActiveDirectoryDeviceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAzureActiveDirectoryDeviceId)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["enrollmentState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseEnrollmentState , m.SetEnrollmentState)
    res["groupTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetGroupTag)
    res["lastContactedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastContactedDateTime)
    res["managedDeviceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetManagedDeviceId)
    res["manufacturer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetManufacturer)
    res["model"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetModel)
    res["productKey"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetProductKey)
    res["purchaseOrderIdentifier"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPurchaseOrderIdentifier)
    res["resourceName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResourceName)
    res["serialNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSerialNumber)
    res["skuNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSkuNumber)
    res["systemFamily"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSystemFamily)
    res["userPrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserPrincipalName)
    return res
}
// GetGroupTag gets the groupTag property value. Group Tag of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetGroupTag()(*string) {
    return m.groupTag
}
// GetLastContactedDateTime gets the lastContactedDateTime property value. Intune Last Contacted Date Time of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetLastContactedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastContactedDateTime
}
// GetManagedDeviceId gets the managedDeviceId property value. Managed Device ID
func (m *WindowsAutopilotDeviceIdentity) GetManagedDeviceId()(*string) {
    return m.managedDeviceId
}
// GetManufacturer gets the manufacturer property value. Oem manufacturer of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetManufacturer()(*string) {
    return m.manufacturer
}
// GetModel gets the model property value. Model name of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetModel()(*string) {
    return m.model
}
// GetProductKey gets the productKey property value. Product Key of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetProductKey()(*string) {
    return m.productKey
}
// GetPurchaseOrderIdentifier gets the purchaseOrderIdentifier property value. Purchase Order Identifier of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetPurchaseOrderIdentifier()(*string) {
    return m.purchaseOrderIdentifier
}
// GetResourceName gets the resourceName property value. Resource Name.
func (m *WindowsAutopilotDeviceIdentity) GetResourceName()(*string) {
    return m.resourceName
}
// GetSerialNumber gets the serialNumber property value. Serial number of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) GetSerialNumber()(*string) {
    return m.serialNumber
}
// GetSkuNumber gets the skuNumber property value. SKU Number
func (m *WindowsAutopilotDeviceIdentity) GetSkuNumber()(*string) {
    return m.skuNumber
}
// GetSystemFamily gets the systemFamily property value. System Family
func (m *WindowsAutopilotDeviceIdentity) GetSystemFamily()(*string) {
    return m.systemFamily
}
// GetUserPrincipalName gets the userPrincipalName property value. User Principal Name.
func (m *WindowsAutopilotDeviceIdentity) GetUserPrincipalName()(*string) {
    return m.userPrincipalName
}
// Serialize serializes information the current object
func (m *WindowsAutopilotDeviceIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("addressableUserName", m.GetAddressableUserName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("azureActiveDirectoryDeviceId", m.GetAzureActiveDirectoryDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetEnrollmentState() != nil {
        cast := (*m.GetEnrollmentState()).String()
        err = writer.WriteStringValue("enrollmentState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("groupTag", m.GetGroupTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastContactedDateTime", m.GetLastContactedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managedDeviceId", m.GetManagedDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productKey", m.GetProductKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("purchaseOrderIdentifier", m.GetPurchaseOrderIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceName", m.GetResourceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serialNumber", m.GetSerialNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuNumber", m.GetSkuNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("systemFamily", m.GetSystemFamily())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAddressableUserName sets the addressableUserName property value. Addressable user name.
func (m *WindowsAutopilotDeviceIdentity) SetAddressableUserName(value *string)() {
    m.addressableUserName = value
}
// SetAzureActiveDirectoryDeviceId sets the azureActiveDirectoryDeviceId property value. AAD Device ID - to be deprecated
func (m *WindowsAutopilotDeviceIdentity) SetAzureActiveDirectoryDeviceId(value *string)() {
    m.azureActiveDirectoryDeviceId = value
}
// SetDisplayName sets the displayName property value. Display Name
func (m *WindowsAutopilotDeviceIdentity) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetEnrollmentState sets the enrollmentState property value. The enrollmentState property
func (m *WindowsAutopilotDeviceIdentity) SetEnrollmentState(value *EnrollmentState)() {
    m.enrollmentState = value
}
// SetGroupTag sets the groupTag property value. Group Tag of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetGroupTag(value *string)() {
    m.groupTag = value
}
// SetLastContactedDateTime sets the lastContactedDateTime property value. Intune Last Contacted Date Time of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetLastContactedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastContactedDateTime = value
}
// SetManagedDeviceId sets the managedDeviceId property value. Managed Device ID
func (m *WindowsAutopilotDeviceIdentity) SetManagedDeviceId(value *string)() {
    m.managedDeviceId = value
}
// SetManufacturer sets the manufacturer property value. Oem manufacturer of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetManufacturer(value *string)() {
    m.manufacturer = value
}
// SetModel sets the model property value. Model name of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetModel(value *string)() {
    m.model = value
}
// SetProductKey sets the productKey property value. Product Key of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetProductKey(value *string)() {
    m.productKey = value
}
// SetPurchaseOrderIdentifier sets the purchaseOrderIdentifier property value. Purchase Order Identifier of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetPurchaseOrderIdentifier(value *string)() {
    m.purchaseOrderIdentifier = value
}
// SetResourceName sets the resourceName property value. Resource Name.
func (m *WindowsAutopilotDeviceIdentity) SetResourceName(value *string)() {
    m.resourceName = value
}
// SetSerialNumber sets the serialNumber property value. Serial number of the Windows autopilot device.
func (m *WindowsAutopilotDeviceIdentity) SetSerialNumber(value *string)() {
    m.serialNumber = value
}
// SetSkuNumber sets the skuNumber property value. SKU Number
func (m *WindowsAutopilotDeviceIdentity) SetSkuNumber(value *string)() {
    m.skuNumber = value
}
// SetSystemFamily sets the systemFamily property value. System Family
func (m *WindowsAutopilotDeviceIdentity) SetSystemFamily(value *string)() {
    m.systemFamily = value
}
// SetUserPrincipalName sets the userPrincipalName property value. User Principal Name.
func (m *WindowsAutopilotDeviceIdentity) SetUserPrincipalName(value *string)() {
    m.userPrincipalName = value
}
