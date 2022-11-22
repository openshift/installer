package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EducationAssignmentIndividualRecipient 
type EducationAssignmentIndividualRecipient struct {
    EducationAssignmentRecipient
    // A collection of IDs of the recipients.
    recipients []string
}
// NewEducationAssignmentIndividualRecipient instantiates a new EducationAssignmentIndividualRecipient and sets the default values.
func NewEducationAssignmentIndividualRecipient()(*EducationAssignmentIndividualRecipient) {
    m := &EducationAssignmentIndividualRecipient{
        EducationAssignmentRecipient: *NewEducationAssignmentRecipient(),
    }
    odataTypeValue := "#microsoft.graph.educationAssignmentIndividualRecipient";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEducationAssignmentIndividualRecipientFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEducationAssignmentIndividualRecipientFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationAssignmentIndividualRecipient(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EducationAssignmentIndividualRecipient) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationAssignmentRecipient.GetFieldDeserializers()
    res["recipients"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetRecipients)
    return res
}
// GetRecipients gets the recipients property value. A collection of IDs of the recipients.
func (m *EducationAssignmentIndividualRecipient) GetRecipients()([]string) {
    return m.recipients
}
// Serialize serializes information the current object
func (m *EducationAssignmentIndividualRecipient) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationAssignmentRecipient.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetRecipients() != nil {
        err = writer.WriteCollectionOfStringValues("recipients", m.GetRecipients())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRecipients sets the recipients property value. A collection of IDs of the recipients.
func (m *EducationAssignmentIndividualRecipient) SetRecipients(value []string)() {
    m.recipients = value
}
