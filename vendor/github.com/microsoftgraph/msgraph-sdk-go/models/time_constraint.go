package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TimeConstraint 
type TimeConstraint struct {
    // The nature of the activity, optional. The possible values are: work, personal, unrestricted, or unknown.
    activityDomain *ActivityDomain
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // The timeSlots property
    timeSlots []TimeSlotable
}
// NewTimeConstraint instantiates a new timeConstraint and sets the default values.
func NewTimeConstraint()(*TimeConstraint) {
    m := &TimeConstraint{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateTimeConstraintFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTimeConstraintFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTimeConstraint(), nil
}
// GetActivityDomain gets the activityDomain property value. The nature of the activity, optional. The possible values are: work, personal, unrestricted, or unknown.
func (m *TimeConstraint) GetActivityDomain()(*ActivityDomain) {
    return m.activityDomain
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TimeConstraint) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TimeConstraint) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["activityDomain"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseActivityDomain , m.SetActivityDomain)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["timeSlots"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTimeSlotFromDiscriminatorValue , m.SetTimeSlots)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *TimeConstraint) GetOdataType()(*string) {
    return m.odataType
}
// GetTimeSlots gets the timeSlots property value. The timeSlots property
func (m *TimeConstraint) GetTimeSlots()([]TimeSlotable) {
    return m.timeSlots
}
// Serialize serializes information the current object
func (m *TimeConstraint) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetActivityDomain() != nil {
        cast := (*m.GetActivityDomain()).String()
        err := writer.WriteStringValue("activityDomain", &cast)
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
    if m.GetTimeSlots() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTimeSlots())
        err := writer.WriteCollectionOfObjectValues("timeSlots", cast)
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
// SetActivityDomain sets the activityDomain property value. The nature of the activity, optional. The possible values are: work, personal, unrestricted, or unknown.
func (m *TimeConstraint) SetActivityDomain(value *ActivityDomain)() {
    m.activityDomain = value
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TimeConstraint) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TimeConstraint) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTimeSlots sets the timeSlots property value. The timeSlots property
func (m *TimeConstraint) SetTimeSlots(value []TimeSlotable)() {
    m.timeSlots = value
}
