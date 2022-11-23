package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerTaskDetails 
type PlannerTaskDetails struct {
    Entity
    // The collection of checklist items on the task.
    checklist PlannerChecklistItemsable
    // Description of the task.
    description *string
    // This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference. When set to automatic the displayed preview is chosen by the app viewing the task.
    previewType *PlannerPreviewType
    // The collection of references on the task.
    references PlannerExternalReferencesable
}
// NewPlannerTaskDetails instantiates a new plannerTaskDetails and sets the default values.
func NewPlannerTaskDetails()(*PlannerTaskDetails) {
    m := &PlannerTaskDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerTaskDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerTaskDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerTaskDetails(), nil
}
// GetChecklist gets the checklist property value. The collection of checklist items on the task.
func (m *PlannerTaskDetails) GetChecklist()(PlannerChecklistItemsable) {
    return m.checklist
}
// GetDescription gets the description property value. Description of the task.
func (m *PlannerTaskDetails) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerTaskDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["checklist"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerChecklistItemsFromDiscriminatorValue , m.SetChecklist)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["previewType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePlannerPreviewType , m.SetPreviewType)
    res["references"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerExternalReferencesFromDiscriminatorValue , m.SetReferences)
    return res
}
// GetPreviewType gets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference. When set to automatic the displayed preview is chosen by the app viewing the task.
func (m *PlannerTaskDetails) GetPreviewType()(*PlannerPreviewType) {
    return m.previewType
}
// GetReferences gets the references property value. The collection of references on the task.
func (m *PlannerTaskDetails) GetReferences()(PlannerExternalReferencesable) {
    return m.references
}
// Serialize serializes information the current object
func (m *PlannerTaskDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("checklist", m.GetChecklist())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetPreviewType() != nil {
        cast := (*m.GetPreviewType()).String()
        err = writer.WriteStringValue("previewType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("references", m.GetReferences())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChecklist sets the checklist property value. The collection of checklist items on the task.
func (m *PlannerTaskDetails) SetChecklist(value PlannerChecklistItemsable)() {
    m.checklist = value
}
// SetDescription sets the description property value. Description of the task.
func (m *PlannerTaskDetails) SetDescription(value *string)() {
    m.description = value
}
// SetPreviewType sets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference. When set to automatic the displayed preview is chosen by the app viewing the task.
func (m *PlannerTaskDetails) SetPreviewType(value *PlannerPreviewType)() {
    m.previewType = value
}
// SetReferences sets the references property value. The collection of references on the task.
func (m *PlannerTaskDetails) SetReferences(value PlannerExternalReferencesable)() {
    m.references = value
}
