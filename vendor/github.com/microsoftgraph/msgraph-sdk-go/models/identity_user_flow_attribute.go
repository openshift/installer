package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IdentityUserFlowAttribute 
type IdentityUserFlowAttribute struct {
    Entity
    // The dataType property
    dataType *IdentityUserFlowAttributeDataType
    // The description of the user flow attribute that's shown to the user at the time of sign-up.
    description *string
    // The display name of the user flow attribute.
    displayName *string
    // The userFlowAttributeType property
    userFlowAttributeType *IdentityUserFlowAttributeType
}
// NewIdentityUserFlowAttribute instantiates a new identityUserFlowAttribute and sets the default values.
func NewIdentityUserFlowAttribute()(*IdentityUserFlowAttribute) {
    m := &IdentityUserFlowAttribute{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityUserFlowAttributeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIdentityUserFlowAttributeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.identityBuiltInUserFlowAttribute":
                        return NewIdentityBuiltInUserFlowAttribute(), nil
                    case "#microsoft.graph.identityCustomUserFlowAttribute":
                        return NewIdentityCustomUserFlowAttribute(), nil
                }
            }
        }
    }
    return NewIdentityUserFlowAttribute(), nil
}
// GetDataType gets the dataType property value. The dataType property
func (m *IdentityUserFlowAttribute) GetDataType()(*IdentityUserFlowAttributeDataType) {
    return m.dataType
}
// GetDescription gets the description property value. The description of the user flow attribute that's shown to the user at the time of sign-up.
func (m *IdentityUserFlowAttribute) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The display name of the user flow attribute.
func (m *IdentityUserFlowAttribute) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IdentityUserFlowAttribute) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["dataType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseIdentityUserFlowAttributeDataType , m.SetDataType)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["userFlowAttributeType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseIdentityUserFlowAttributeType , m.SetUserFlowAttributeType)
    return res
}
// GetUserFlowAttributeType gets the userFlowAttributeType property value. The userFlowAttributeType property
func (m *IdentityUserFlowAttribute) GetUserFlowAttributeType()(*IdentityUserFlowAttributeType) {
    return m.userFlowAttributeType
}
// Serialize serializes information the current object
func (m *IdentityUserFlowAttribute) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDataType() != nil {
        cast := (*m.GetDataType()).String()
        err = writer.WriteStringValue("dataType", &cast)
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
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetUserFlowAttributeType() != nil {
        cast := (*m.GetUserFlowAttributeType()).String()
        err = writer.WriteStringValue("userFlowAttributeType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDataType sets the dataType property value. The dataType property
func (m *IdentityUserFlowAttribute) SetDataType(value *IdentityUserFlowAttributeDataType)() {
    m.dataType = value
}
// SetDescription sets the description property value. The description of the user flow attribute that's shown to the user at the time of sign-up.
func (m *IdentityUserFlowAttribute) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The display name of the user flow attribute.
func (m *IdentityUserFlowAttribute) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetUserFlowAttributeType sets the userFlowAttributeType property value. The userFlowAttributeType property
func (m *IdentityUserFlowAttribute) SetUserFlowAttributeType(value *IdentityUserFlowAttributeType)() {
    m.userFlowAttributeType = value
}
