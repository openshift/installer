package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IdentityUserFlow 
type IdentityUserFlow struct {
    Entity
}
// NewIdentityUserFlow instantiates a new identityUserFlow and sets the default values.
func NewIdentityUserFlow()(*IdentityUserFlow) {
    m := &IdentityUserFlow{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityUserFlowFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIdentityUserFlowFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.b2xIdentityUserFlow":
                        return NewB2xIdentityUserFlow(), nil
                }
            }
        }
    }
    return NewIdentityUserFlow(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IdentityUserFlow) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["userFlowType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserFlowType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserFlowType(val.(*UserFlowType))
        }
        return nil
    }
    res["userFlowTypeVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserFlowTypeVersion(val)
        }
        return nil
    }
    return res
}
// GetUserFlowType gets the userFlowType property value. The userFlowType property
func (m *IdentityUserFlow) GetUserFlowType()(*UserFlowType) {
    val, err := m.GetBackingStore().Get("userFlowType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserFlowType)
    }
    return nil
}
// GetUserFlowTypeVersion gets the userFlowTypeVersion property value. The userFlowTypeVersion property
func (m *IdentityUserFlow) GetUserFlowTypeVersion()(*float32) {
    val, err := m.GetBackingStore().Get("userFlowTypeVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityUserFlow) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetUserFlowType() != nil {
        cast := (*m.GetUserFlowType()).String()
        err = writer.WriteStringValue("userFlowType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat32Value("userFlowTypeVersion", m.GetUserFlowTypeVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUserFlowType sets the userFlowType property value. The userFlowType property
func (m *IdentityUserFlow) SetUserFlowType(value *UserFlowType)() {
    err := m.GetBackingStore().Set("userFlowType", value)
    if err != nil {
        panic(err)
    }
}
// SetUserFlowTypeVersion sets the userFlowTypeVersion property value. The userFlowTypeVersion property
func (m *IdentityUserFlow) SetUserFlowTypeVersion(value *float32)() {
    err := m.GetBackingStore().Set("userFlowTypeVersion", value)
    if err != nil {
        panic(err)
    }
}
// IdentityUserFlowable 
type IdentityUserFlowable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUserFlowType()(*UserFlowType)
    GetUserFlowTypeVersion()(*float32)
    SetUserFlowType(value *UserFlowType)()
    SetUserFlowTypeVersion(value *float32)()
}
