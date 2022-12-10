package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LocationConstraint 
type LocationConstraint struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The client requests the service to include in the response a meeting location for the meeting. If this is true and all the resources are busy, findMeetingTimes will not return any meeting time suggestions. If this is false and all the resources are busy, findMeetingTimes would still look for meeting times without locations.
    isRequired *bool
    // Constraint information for one or more locations that the client requests for the meeting.
    locations []LocationConstraintItemable
    // The OdataType property
    odataType *string
    // The client requests the service to suggest one or more meeting locations.
    suggestLocation *bool
}
// NewLocationConstraint instantiates a new locationConstraint and sets the default values.
func NewLocationConstraint()(*LocationConstraint) {
    m := &LocationConstraint{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateLocationConstraintFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLocationConstraintFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocationConstraint(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *LocationConstraint) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LocationConstraint) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsRequired)
    res["locations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateLocationConstraintItemFromDiscriminatorValue , m.SetLocations)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["suggestLocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSuggestLocation)
    return res
}
// GetIsRequired gets the isRequired property value. The client requests the service to include in the response a meeting location for the meeting. If this is true and all the resources are busy, findMeetingTimes will not return any meeting time suggestions. If this is false and all the resources are busy, findMeetingTimes would still look for meeting times without locations.
func (m *LocationConstraint) GetIsRequired()(*bool) {
    return m.isRequired
}
// GetLocations gets the locations property value. Constraint information for one or more locations that the client requests for the meeting.
func (m *LocationConstraint) GetLocations()([]LocationConstraintItemable) {
    return m.locations
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *LocationConstraint) GetOdataType()(*string) {
    return m.odataType
}
// GetSuggestLocation gets the suggestLocation property value. The client requests the service to suggest one or more meeting locations.
func (m *LocationConstraint) GetSuggestLocation()(*bool) {
    return m.suggestLocation
}
// Serialize serializes information the current object
func (m *LocationConstraint) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isRequired", m.GetIsRequired())
        if err != nil {
            return err
        }
    }
    if m.GetLocations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLocations())
        err := writer.WriteCollectionOfObjectValues("locations", cast)
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
        err := writer.WriteBoolValue("suggestLocation", m.GetSuggestLocation())
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
func (m *LocationConstraint) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsRequired sets the isRequired property value. The client requests the service to include in the response a meeting location for the meeting. If this is true and all the resources are busy, findMeetingTimes will not return any meeting time suggestions. If this is false and all the resources are busy, findMeetingTimes would still look for meeting times without locations.
func (m *LocationConstraint) SetIsRequired(value *bool)() {
    m.isRequired = value
}
// SetLocations sets the locations property value. Constraint information for one or more locations that the client requests for the meeting.
func (m *LocationConstraint) SetLocations(value []LocationConstraintItemable)() {
    m.locations = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *LocationConstraint) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSuggestLocation sets the suggestLocation property value. The client requests the service to suggest one or more meeting locations.
func (m *LocationConstraint) SetSuggestLocation(value *bool)() {
    m.suggestLocation = value
}
