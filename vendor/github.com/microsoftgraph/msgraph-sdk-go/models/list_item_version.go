package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ListItemVersion 
type ListItemVersion struct {
    BaseItemVersion
    // A collection of the fields and values for this version of the list item.
    fields FieldValueSetable
}
// NewListItemVersion instantiates a new ListItemVersion and sets the default values.
func NewListItemVersion()(*ListItemVersion) {
    m := &ListItemVersion{
        BaseItemVersion: *NewBaseItemVersion(),
    }
    odataTypeValue := "#microsoft.graph.listItemVersion";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateListItemVersionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateListItemVersionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.documentSetVersion":
                        return NewDocumentSetVersion(), nil
                }
            }
        }
    }
    return NewListItemVersion(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ListItemVersion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItemVersion.GetFieldDeserializers()
    res["fields"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFieldValueSetFromDiscriminatorValue , m.SetFields)
    return res
}
// GetFields gets the fields property value. A collection of the fields and values for this version of the list item.
func (m *ListItemVersion) GetFields()(FieldValueSetable) {
    return m.fields
}
// Serialize serializes information the current object
func (m *ListItemVersion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItemVersion.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("fields", m.GetFields())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFields sets the fields property value. A collection of the fields and values for this version of the list item.
func (m *ListItemVersion) SetFields(value FieldValueSetable)() {
    m.fields = value
}
