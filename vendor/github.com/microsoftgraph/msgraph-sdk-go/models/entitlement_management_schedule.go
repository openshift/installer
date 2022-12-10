package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EntitlementManagementSchedule 
type EntitlementManagementSchedule struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // When the access should expire.
    expiration ExpirationPatternable
    // The OdataType property
    odataType *string
    // For recurring access reviews.  Not used in access requests.
    recurrence PatternedRecurrenceable
    // The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
    startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
}
// NewEntitlementManagementSchedule instantiates a new entitlementManagementSchedule and sets the default values.
func NewEntitlementManagementSchedule()(*EntitlementManagementSchedule) {
    m := &EntitlementManagementSchedule{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateEntitlementManagementScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEntitlementManagementScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEntitlementManagementSchedule(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EntitlementManagementSchedule) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetExpiration gets the expiration property value. When the access should expire.
func (m *EntitlementManagementSchedule) GetExpiration()(ExpirationPatternable) {
    return m.expiration
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EntitlementManagementSchedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["expiration"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateExpirationPatternFromDiscriminatorValue , m.SetExpiration)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["recurrence"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue , m.SetRecurrence)
    res["startDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetStartDateTime)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *EntitlementManagementSchedule) GetOdataType()(*string) {
    return m.odataType
}
// GetRecurrence gets the recurrence property value. For recurring access reviews.  Not used in access requests.
func (m *EntitlementManagementSchedule) GetRecurrence()(PatternedRecurrenceable) {
    return m.recurrence
}
// GetStartDateTime gets the startDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *EntitlementManagementSchedule) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.startDateTime
}
// Serialize serializes information the current object
func (m *EntitlementManagementSchedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("expiration", m.GetExpiration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("recurrence", m.GetRecurrence())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EntitlementManagementSchedule) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetExpiration sets the expiration property value. When the access should expire.
func (m *EntitlementManagementSchedule) SetExpiration(value ExpirationPatternable)() {
    m.expiration = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *EntitlementManagementSchedule) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRecurrence sets the recurrence property value. For recurring access reviews.  Not used in access requests.
func (m *EntitlementManagementSchedule) SetRecurrence(value PatternedRecurrenceable)() {
    m.recurrence = value
}
// SetStartDateTime sets the startDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *EntitlementManagementSchedule) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.startDateTime = value
}
