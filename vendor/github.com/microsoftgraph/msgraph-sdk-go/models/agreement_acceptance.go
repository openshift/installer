package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AgreementAcceptance provides operations to manage the collection of agreementAcceptance entities.
type AgreementAcceptance struct {
    Entity
    // The identifier of the agreement file accepted by the user.
    agreementFileId *string
    // The identifier of the agreement.
    agreementId *string
    // The display name of the device used for accepting the agreement.
    deviceDisplayName *string
    // The unique identifier of the device used for accepting the agreement. Supports $filter (eq) and eq for null values.
    deviceId *string
    // The operating system used to accept the agreement.
    deviceOSType *string
    // The operating system version of the device used to accept the agreement.
    deviceOSVersion *string
    // The expiration date time of the acceptance. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and eq for null values.
    expirationDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
    recordedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The state of the agreement acceptance. Possible values are: accepted, declined. Supports $filter (eq).
    state *AgreementAcceptanceState
    // Display name of the user when the acceptance was recorded.
    userDisplayName *string
    // Email of the user when the acceptance was recorded.
    userEmail *string
    // The identifier of the user who accepted the agreement. Supports $filter (eq).
    userId *string
    // UPN of the user when the acceptance was recorded.
    userPrincipalName *string
}
// NewAgreementAcceptance instantiates a new agreementAcceptance and sets the default values.
func NewAgreementAcceptance()(*AgreementAcceptance) {
    m := &AgreementAcceptance{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAgreementAcceptanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAgreementAcceptanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementAcceptance(), nil
}
// GetAgreementFileId gets the agreementFileId property value. The identifier of the agreement file accepted by the user.
func (m *AgreementAcceptance) GetAgreementFileId()(*string) {
    return m.agreementFileId
}
// GetAgreementId gets the agreementId property value. The identifier of the agreement.
func (m *AgreementAcceptance) GetAgreementId()(*string) {
    return m.agreementId
}
// GetDeviceDisplayName gets the deviceDisplayName property value. The display name of the device used for accepting the agreement.
func (m *AgreementAcceptance) GetDeviceDisplayName()(*string) {
    return m.deviceDisplayName
}
// GetDeviceId gets the deviceId property value. The unique identifier of the device used for accepting the agreement. Supports $filter (eq) and eq for null values.
func (m *AgreementAcceptance) GetDeviceId()(*string) {
    return m.deviceId
}
// GetDeviceOSType gets the deviceOSType property value. The operating system used to accept the agreement.
func (m *AgreementAcceptance) GetDeviceOSType()(*string) {
    return m.deviceOSType
}
// GetDeviceOSVersion gets the deviceOSVersion property value. The operating system version of the device used to accept the agreement.
func (m *AgreementAcceptance) GetDeviceOSVersion()(*string) {
    return m.deviceOSVersion
}
// GetExpirationDateTime gets the expirationDateTime property value. The expiration date time of the acceptance. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and eq for null values.
func (m *AgreementAcceptance) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.expirationDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AgreementAcceptance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["agreementFileId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAgreementFileId)
    res["agreementId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAgreementId)
    res["deviceDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceDisplayName)
    res["deviceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceId)
    res["deviceOSType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceOSType)
    res["deviceOSVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceOSVersion)
    res["expirationDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetExpirationDateTime)
    res["recordedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetRecordedDateTime)
    res["state"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAgreementAcceptanceState , m.SetState)
    res["userDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserDisplayName)
    res["userEmail"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserEmail)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    res["userPrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserPrincipalName)
    return res
}
// GetRecordedDateTime gets the recordedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AgreementAcceptance) GetRecordedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.recordedDateTime
}
// GetState gets the state property value. The state of the agreement acceptance. Possible values are: accepted, declined. Supports $filter (eq).
func (m *AgreementAcceptance) GetState()(*AgreementAcceptanceState) {
    return m.state
}
// GetUserDisplayName gets the userDisplayName property value. Display name of the user when the acceptance was recorded.
func (m *AgreementAcceptance) GetUserDisplayName()(*string) {
    return m.userDisplayName
}
// GetUserEmail gets the userEmail property value. Email of the user when the acceptance was recorded.
func (m *AgreementAcceptance) GetUserEmail()(*string) {
    return m.userEmail
}
// GetUserId gets the userId property value. The identifier of the user who accepted the agreement. Supports $filter (eq).
func (m *AgreementAcceptance) GetUserId()(*string) {
    return m.userId
}
// GetUserPrincipalName gets the userPrincipalName property value. UPN of the user when the acceptance was recorded.
func (m *AgreementAcceptance) GetUserPrincipalName()(*string) {
    return m.userPrincipalName
}
// Serialize serializes information the current object
func (m *AgreementAcceptance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("agreementFileId", m.GetAgreementFileId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("agreementId", m.GetAgreementId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceDisplayName", m.GetDeviceDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceOSType", m.GetDeviceOSType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceOSVersion", m.GetDeviceOSVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("expirationDateTime", m.GetExpirationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("recordedDateTime", m.GetRecordedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userDisplayName", m.GetUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userEmail", m.GetUserEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
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
// SetAgreementFileId sets the agreementFileId property value. The identifier of the agreement file accepted by the user.
func (m *AgreementAcceptance) SetAgreementFileId(value *string)() {
    m.agreementFileId = value
}
// SetAgreementId sets the agreementId property value. The identifier of the agreement.
func (m *AgreementAcceptance) SetAgreementId(value *string)() {
    m.agreementId = value
}
// SetDeviceDisplayName sets the deviceDisplayName property value. The display name of the device used for accepting the agreement.
func (m *AgreementAcceptance) SetDeviceDisplayName(value *string)() {
    m.deviceDisplayName = value
}
// SetDeviceId sets the deviceId property value. The unique identifier of the device used for accepting the agreement. Supports $filter (eq) and eq for null values.
func (m *AgreementAcceptance) SetDeviceId(value *string)() {
    m.deviceId = value
}
// SetDeviceOSType sets the deviceOSType property value. The operating system used to accept the agreement.
func (m *AgreementAcceptance) SetDeviceOSType(value *string)() {
    m.deviceOSType = value
}
// SetDeviceOSVersion sets the deviceOSVersion property value. The operating system version of the device used to accept the agreement.
func (m *AgreementAcceptance) SetDeviceOSVersion(value *string)() {
    m.deviceOSVersion = value
}
// SetExpirationDateTime sets the expirationDateTime property value. The expiration date time of the acceptance. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and eq for null values.
func (m *AgreementAcceptance) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.expirationDateTime = value
}
// SetRecordedDateTime sets the recordedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *AgreementAcceptance) SetRecordedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.recordedDateTime = value
}
// SetState sets the state property value. The state of the agreement acceptance. Possible values are: accepted, declined. Supports $filter (eq).
func (m *AgreementAcceptance) SetState(value *AgreementAcceptanceState)() {
    m.state = value
}
// SetUserDisplayName sets the userDisplayName property value. Display name of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserDisplayName(value *string)() {
    m.userDisplayName = value
}
// SetUserEmail sets the userEmail property value. Email of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserEmail(value *string)() {
    m.userEmail = value
}
// SetUserId sets the userId property value. The identifier of the user who accepted the agreement. Supports $filter (eq).
func (m *AgreementAcceptance) SetUserId(value *string)() {
    m.userId = value
}
// SetUserPrincipalName sets the userPrincipalName property value. UPN of the user when the acceptance was recorded.
func (m *AgreementAcceptance) SetUserPrincipalName(value *string)() {
    m.userPrincipalName = value
}
