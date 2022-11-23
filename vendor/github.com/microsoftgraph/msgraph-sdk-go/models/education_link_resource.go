package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EducationLinkResource 
type EducationLinkResource struct {
    EducationResource
    // URL to the resource.
    link *string
}
// NewEducationLinkResource instantiates a new EducationLinkResource and sets the default values.
func NewEducationLinkResource()(*EducationLinkResource) {
    m := &EducationLinkResource{
        EducationResource: *NewEducationResource(),
    }
    odataTypeValue := "#microsoft.graph.educationLinkResource";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEducationLinkResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEducationLinkResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationLinkResource(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EducationLinkResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationResource.GetFieldDeserializers()
    res["link"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLink)
    return res
}
// GetLink gets the link property value. URL to the resource.
func (m *EducationLinkResource) GetLink()(*string) {
    return m.link
}
// Serialize serializes information the current object
func (m *EducationLinkResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("link", m.GetLink())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLink sets the link property value. URL to the resource.
func (m *EducationLinkResource) SetLink(value *string)() {
    m.link = value
}
