package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SubjectSet 
type SubjectSet struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
}
// NewSubjectSet instantiates a new subjectSet and sets the default values.
func NewSubjectSet()(*SubjectSet) {
    m := &SubjectSet{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSubjectSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSubjectSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.attributeRuleMembers":
                        return NewAttributeRuleMembers(), nil
                    case "#microsoft.graph.connectedOrganizationMembers":
                        return NewConnectedOrganizationMembers(), nil
                    case "#microsoft.graph.externalSponsors":
                        return NewExternalSponsors(), nil
                    case "#microsoft.graph.groupMembers":
                        return NewGroupMembers(), nil
                    case "#microsoft.graph.internalSponsors":
                        return NewInternalSponsors(), nil
                    case "#microsoft.graph.requestorManager":
                        return NewRequestorManager(), nil
                    case "#microsoft.graph.singleServicePrincipal":
                        return NewSingleServicePrincipal(), nil
                    case "#microsoft.graph.singleUser":
                        return NewSingleUser(), nil
                    case "#microsoft.graph.targetApplicationOwners":
                        return NewTargetApplicationOwners(), nil
                    case "#microsoft.graph.targetManager":
                        return NewTargetManager(), nil
                }
            }
        }
    }
    return NewSubjectSet(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SubjectSet) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SubjectSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SubjectSet) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SubjectSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
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
func (m *SubjectSet) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SubjectSet) SetOdataType(value *string)() {
    m.odataType = value
}
